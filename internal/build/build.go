package build

import (
	"fmt"
	htmltemplate "html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/sxijyoti/whiskey/internal/config"
	"github.com/sxijyoti/whiskey/internal/dependency"
	"github.com/sxijyoti/whiskey/internal/fingerprint"
	"github.com/sxijyoti/whiskey/internal/graph"
	"github.com/sxijyoti/whiskey/internal/parser"
	"github.com/sxijyoti/whiskey/internal/source"
	"github.com/sxijyoti/whiskey/internal/template"
)

// BuildPage renders a single page document to its output path.
//
// Include resolution follows the Whiskey model:
//   - local: prefix → read directly from the filesystem (no workspace, no HTTP)
//   - remote refs    → read from the materialized workspace
func BuildPage(
	siteRoot string,
	cfg *config.Config,
	sourcePath string,
	doc *parser.Document,
	output string,
) error {

	resolveLocalInclude := func(path string) ([]byte, error) {
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		// Only Markdown files can contain Whiskey frontmatter.
		if filepath.Ext(path) != ".md" {
			return raw, nil
		}

		included, err := parser.ParseFrontmatter(raw)
		if err != nil {
			return nil, err
		}

		return []byte(included.Body), nil
	}

	resolvedBody, err := dependency.ResolveIncludes(
		doc.Body,
		sourcePath,
		func(ref string) ([]byte, string, error) {

			if strings.HasPrefix(ref, "local:") {
				path, err := dependency.ResolveLocalPath(sourcePath, ref)
				if err != nil {
					return nil, "", err
				}
				body, err := resolveLocalInclude(path)
				return body, path, err
			}

			normalized, err := source.NormalizeRef(ref)
			if err != nil {
				return nil, "", err
			}

			if !source.WorkspaceExists(siteRoot, normalized) {
				return nil, "", fmt.Errorf("workspace missing for %s", normalized)
			}

			body, err := source.ReadWorkspace(siteRoot, normalized)
			return body, "", err
		},
	)
	if err != nil {
		return err
	}

	resolvedBody, err = parser.ExpandShortcodes(resolvedBody)
	if err != nil {
		return err
	}

	html, err := parser.MdToHTML(resolvedBody)
	if err != nil {
		return err
	}

	layout := doc.Meta.Layout
	if layout == "" {
		layout = "page"
	}

	var date string
	if !doc.Meta.Date.IsZero() {
		date = doc.Meta.Date.Format("2006-01-02")
	}

	page, err := template.RenderPage(
		siteRoot,
		cfg.Theme,
		layout,
		template.PageData{
			Site:        cfg,
			Title:       doc.Meta.Title,
			Description: doc.Meta.Description,
			Date:        date,
			Tags:        doc.Meta.Tags,
			Content:     htmltemplate.HTML(html),
		},
	)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return err
	}

	return os.WriteFile(output, page, 0644)
}

// BuildSite performs a full site build: all pages, static assets, and
// secondary indexes. Draft pages are skipped during rendering but their
// dependency edges are already in the graph for tracking purposes.
func BuildSite(root string) error {
	return BuildSiteWithReason(root, "manual")
}

func BuildSiteWithReason(root, reason string) error {
	start := logBuildStart("Full rebuild", reason)

	if err := EnsureWorkspace(root); err != nil {
		return err
	}

	cfg, err := config.Load(root)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(filepath.Join(root, "dist")); err != nil {
		return err
	}

	contentRoot := filepath.Join(root, "content")

	pages, err := DiscoverPages(contentRoot)
	if err != nil {
		return err
	}

	g, err := graph.BuildSiteGraph(root, cfg.Theme)
	if err != nil {
		return err
	}

	manifest, err := source.LoadManifest(root)
	if err != nil {
		return err
	}

	materialized, materializeErr := MaterializeSources(root, g, manifest)
	logSources(nil, materialized)

	if err := GarbageCollectWorkspace(root, g, manifest); err != nil {
		return err
	}

	if materializeErr != nil {
		return materializeErr
	}

	var builtPages []string
	failedPages := make(map[string]bool)

	for _, page := range pages {
		raw, err := os.ReadFile(page)
		if err != nil {
			return err
		}

		doc, err := parser.ParseFrontmatter(raw)
		if err != nil {
			return err
		}

		// Draft pages are in the graph for dependency tracking but are never
		// rendered in a published build.
		if doc.Meta.Draft {
			continue
		}

		output, err := pageOutputPath(root, contentRoot, page)
		if err != nil {
			return err
		}

		slug := pageSlug(contentRoot, page)

		skip := false
		for _, dep := range g.Dependencies(page) {
			if materialized.OfflineCached[dep] {
				continue
			}
			if _, failed := materialized.Failed[dep]; failed {
				_ = os.Remove(output)
				logPageError(contentRoot, page, fmt.Errorf("workspace missing for %s", dep))
				failedPages[slug] = true
				skip = true
				break
			}
		}

		if skip {
			continue
		}

		if err := BuildPage(root, cfg, page, doc, output); err != nil {
			_ = os.Remove(output)
			logPageError(contentRoot, page, err)
			failedPages[slug] = true
			continue
		}

		builtPages = append(builtPages, page)
	}

	// Build secondary indexes from successfully built pages only.
	cleanIndex, err := BuildIndex(root, builtPages)
	if err != nil {
		return err
	}

	cfg.Nav = BuildNav(cfg, cleanIndex)

	for _, page := range builtPages {
		raw, err := os.ReadFile(page)
		if err != nil {
			return err
		}

		doc, err := parser.ParseFrontmatter(raw)
		if err != nil {
			return err
		}

		output, err := pageOutputPath(root, contentRoot, page)
		if err != nil {
			return err
		}

		if err := BuildPage(root, cfg, page, doc, output); err != nil {
			logPageError(contentRoot, page, err)
			return err
		}
	}

	logRenderCount(len(builtPages))

	if err := BuildTags(root, cfg, cleanIndex); err != nil {
		return err
	}

	if err := BuildRSS(root, cfg, cleanIndex); err != nil {
		return err
	}

	if err := BuildSitemap(root, cfg, cleanIndex); err != nil {
		return err
	}

	if err := BuildCollections(root, cfg, cleanIndex); err != nil {
		return err
	}

	if err := CopyStatic(root, cfg.Theme); err != nil {
		return err
	}

	// Initialize and save fingerprints after a successful full build so subsequent
	// builds can be incremental.
	store := make(fingerprint.Store)
	if changed, err := fingerprint.ChangedNodes(root, g, store); err == nil {
		_ = changed
		store["__config__"] = fingerprint.Entry{Hash: fingerprint.ConfigHash(cfg)}
		_ = fingerprint.Save(filepath.Join(root, ".whiskey", "fingerprints.json"), store)
	}

	if len(failedPages) > 0 {
		return fmt.Errorf("full build completed with %d failure(s)", len(failedPages))
	}

	logBuildDone("Full rebuild", start, len(builtPages), materialized)

	return nil
}
