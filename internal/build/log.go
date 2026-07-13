package build

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
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

func logSources(result *MaterializationResult) {
	if result == nil {
		return
	}

	fmt.Println("[sources]")
	if len(result.Updated) == 0 && len(result.Failed) == 0 {
		fmt.Println("  No remote changes")
		return
	}

	if len(result.Updated) > 0 {
		fmt.Println("  Updated:")
		for _, ref := range result.Updated {
			fmt.Printf("    %s\n", sourceName(ref))
		}
	}

	if len(result.Failed) > 0 {
		fmt.Println("  Failed:")
		for ref := range result.Failed {
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
