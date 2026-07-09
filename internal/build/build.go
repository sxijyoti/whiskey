package build

import (
	"fmt"
	htmltemplate "html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/sxijyoti/whiskey/internal/config"
	"github.com/sxijyoti/whiskey/internal/dependency"
	"github.com/sxijyoti/whiskey/internal/graph"
	"github.com/sxijyoti/whiskey/internal/parser"
	"github.com/sxijyoti/whiskey/internal/source"
	"github.com/sxijyoti/whiskey/internal/template"
)

func BuildPage(
	siteRoot string,
	cfg *config.Config,
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

		doc, err := parser.ParseFrontmatter(raw)
		if err != nil {
			return nil, err
		}

		return []byte(doc.Body), nil
	}

	resolvedBody, err := dependency.ResolveIncludes(
		doc.Body,
		func(ref string) ([]byte, error) {

			if strings.HasPrefix(ref, "local:") {
				path := strings.TrimPrefix(
					ref,
					"local:",
				)

				return resolveLocalInclude(path)
			}

			if !source.WorkspaceExists(
				siteRoot,
				ref,
			) {
				return nil, fmt.Errorf(
					"workspace missing for %s",
					ref,
				)
			}

			return source.ReadWorkspace(
				siteRoot,
				ref,
			)
		},
	)
	if err != nil {
		return err
	}

	resolvedBody, err = parser.ExpandShortcodes(
		resolvedBody,
	)
	if err != nil {
		return err
	}

	html, err := parser.MdToHTML(
		resolvedBody,
	)
	if err != nil {
		return err
	}

	layout := doc.Meta.Layout
	if layout == "" {
		layout = "page"
	}

	var date string

	if !doc.Meta.Date.IsZero() {
		date = doc.Meta.Date.Format(
			"2006-01-02",
		)
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

	if err := os.MkdirAll(
		filepath.Dir(output),
		0755,
	); err != nil {
		return err
	}

	return os.WriteFile(
		output,
		page,
		0644,
	)
}

func BuildSite(root string) error {
	if err := EnsureWorkspace(root); err != nil {
		return err
	}

	cfg, err := config.Load(root)
	if err != nil {
		return err
	}

	contentRoot := filepath.Join(root, "content")
	pages, err := DiscoverPages(contentRoot)
	if err != nil {
		return err
	}

	// 1. Build initial complete index for navigation context
	index, err := BuildIndex(root, pages)
	if err != nil {
		return err
	}

	cfg.Nav = BuildNav(cfg, index)

	g, err := graph.BuildSiteGraph(root, cfg.Theme)
	if err != nil {
		return err
	}

	manifest, err := source.LoadManifest(root)
	if err != nil {
		return err
	}

	materialized, err := MaterializeSources(root, g, manifest)
	if err := GarbageCollectWorkspace(root, g, manifest); err != nil {
		return err
	}

	var builtPages []string
	failedPages := make(map[string]bool) // Track failed pages explicitly

	for _, page := range pages {
		raw, err := os.ReadFile(page)
		if err != nil {
			return err
		}

		doc, err := parser.ParseFrontmatter(raw)
		if err != nil {
			return err
		}

		if doc.Meta.Draft {
			continue
		}

		output, err := pageOutputPath(root, contentRoot, page)
		if err != nil {
			return err
		}

		// Calculate relative slug matching the index keying approach
		rel, _ := filepath.Rel(contentRoot, page)
		slug := strings.TrimSuffix(rel, filepath.Ext(rel))

		skip := false
		for _, dep := range g.Dependencies(page) {
			if materialized.OfflineCached[dep] {
				continue
			}

			if _, failed := materialized.Failed[dep]; failed {
				_ = os.Remove(output)
				fmt.Printf("[build] skipped %s (missing %s)\n", page, dep)
				failedPages[slug] = true
				skip = true
				break
			}
		}

		if skip {
			continue
		}

		if err := BuildPage(root, cfg, doc, output); err != nil {
			_ = os.Remove(output)
			fmt.Printf("[build] failed %s: %v\n", page, err)
			failedPages[slug] = true
			continue
		}

		builtPages = append(builtPages, page)
	}

	// 2. Re-instantiate a clean SiteIndex passing only successful paths
	// to secondary artifact builders to completely eliminate ghost entries.
	cleanIndex, err := BuildIndex(root, builtPages)
	if err != nil {
		return err
	}

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

	return nil
}
