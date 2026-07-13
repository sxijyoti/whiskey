package build

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var LogNoopBuilds = true

func logBuildStart(kind, reason string) time.Time {
	fmt.Printf("[build] %s\n", kind)
	if reason != "" {
		fmt.Printf("  Reason: %s\n", reason)
	}
	return time.Now()
}

func logBuildDone(kind string, start time.Time, rendered int, sources *MaterializationResult) {
	fmt.Printf("[build] %s (%s)\n", kind, time.Since(start).Round(time.Millisecond))
	fmt.Printf("\n  Rendered: %s\n", pageCount(rendered))
	if sources != nil {
		fmt.Printf("  Sources updated: %d\n", len(sources.Updated))
	}
}

// logSources prints the [sources] diagnostic block.
//
// localChanged is a sorted list of site-root-relative paths for local
// dependency sources (e.g. @include'd files) that changed.
// result carries updated/failed remote sources from MaterializeSources.
//
// The entire block is omitted when neither local nor remote sources changed.
func logSources(localChanged []string, result *MaterializationResult) {
	remoteUpdated := 0
	remoteFailed := 0
	if result != nil {
		remoteUpdated = len(result.Updated)
		remoteFailed = len(result.Failed)
	}

	if len(localChanged) == 0 && remoteUpdated == 0 && remoteFailed == 0 {
		// Nothing changed — omit the section entirely.
		return
	}

	fmt.Println("[sources]")

	if len(localChanged) > 0 {
		fmt.Println("  Local:")
		for _, p := range localChanged {
			fmt.Printf("    %s\n", p)
		}
	}

	if remoteUpdated > 0 {
		if len(localChanged) > 0 {
			fmt.Println()
		}
		fmt.Println("  Remote:")
		for _, ref := range result.Updated {
			fmt.Printf("    %s\n", sourceName(ref))
		}
	}

	if remoteFailed > 0 {
		if len(localChanged) > 0 || remoteUpdated > 0 {
			fmt.Println()
		}
		fmt.Println("  Failed:")
		refs := make([]string, 0, len(result.Failed))
		for ref := range result.Failed {
			refs = append(refs, ref)
		}
		sort.Strings(refs)
		for _, ref := range refs {
			fmt.Printf("    %s\n", sourceName(ref))
		}
	}
}

func logDirtyPages(count int) {
	switch count {
	case 1:
		fmt.Println("  Dirty pages: 1 dirty page")
	default:
		fmt.Printf("  Dirty pages: %d dirty page(s)\n", count)
	}
}

func logDirtyAssets(count int) {
	fmt.Printf("  Dirty assets: %d\n", count)
}

func logRenderPage(contentRoot, page string) {
	fmt.Printf("[render] %s\n", pageLabel(contentRoot, page))
}

func logRenderCount(count int) {
	fmt.Printf("[render] %s\n", pageCount(count))
}

func logPageError(contentRoot, page string, err error) {
	fmt.Printf("[error] %s\n\n", pageLabel(contentRoot, page))
	msg := err.Error()
	switch {
	case strings.Contains(strings.ToLower(msg), "cycle detected"):
		fmt.Println("  Cycle detected")
	case strings.Contains(strings.ToLower(msg), "workspace missing"):
		fmt.Println("  Remote source unavailable")
		fmt.Printf("\n  %s\n", strings.TrimPrefix(msg, "workspace missing for "))
	default:
		fmt.Printf("  %s\n", msg)
	}
	fmt.Println()
}

func pageLabel(contentRoot, page string) string {
	rel, err := filepath.Rel(contentRoot, page)
	if err != nil {
		return filepath.Base(page)
	}
	rel = filepath.ToSlash(rel)
	return strings.TrimSuffix(rel, filepath.Ext(rel))
}

func pageCount(count int) string {
	if count == 1 {
		return "1 page"
	}
	return fmt.Sprintf("%d pages", count)
}

func sourceName(ref string) string {
	if strings.HasPrefix(ref, "local:") {
		return filepath.Base(strings.TrimPrefix(ref, "local:"))
	}

	parsed, err := url.Parse(ref)
	if err == nil && parsed.Path != "" {
		base := path.Base(parsed.Path)
		if base != "." && base != "/" {
			return base
		}
	}

	return ref
}
