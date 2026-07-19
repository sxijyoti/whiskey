---
title: Dev-Guide
description: Technical contributor documentation explaining Whiskey's architecture, package structure, and pipelines.
date: 2026-07-13
layout: page
draft: false
tags:
  - developer
  - internals
  - go
---

# Developer Guide

This guide is intended for contributors and anyone interested in Whiskey's internal architecture. It explains how the major components interact during a build, from dependency discovery to page rendering.

---

## Repository Layout

The project is structured as a standard Go monorepo:

- **`cmd/whiskey/`**: Entry point for the CLI tool.
- **`internal/`**: Core application logic, split into modular packages:
  - `cli`: Commands and flags parsed via Cobra.
  - `config`: Handles parsing, default parameters, and validation of `whiskey.toml`.
  - `source`: Workspace management, manifest files, and remote source retrieval.
  - `graph`: Handles dependency graph construction (pages, layout, partials, assets, sources).
  - `fingerprint`: Computes cryptographic node states and writes `.whiskey/fingerprints.json`.
  - `planner`: Determines which pages are dirty and need rebuilding.
  - `dependency`: Extracts and resolves local and remote include directives.
  - `parser`: Extracts YAML frontmatter and converts Markdown text to HTML.
  - `template`: Manages HTML template loading and page rendering.
  - `build`: Coordinates the incremental and full build processes.

---

## Build Pipeline

Every build—whether incremental or full—follows the same high-level pipeline.

The process begins by preparing the workspace and loading the site configuration. Whiskey then constructs the site's dependency graph, materializes any remote content, and detects changes by comparing node fingerprints from the previous build.

If a global dependency such as the configuration, layouts or partials has changed, Whiskey performs a full rebuild. Otherwise, it computes the minimal set of affected pages, copies only modified assets, rebuilds dirty pages, regenerates secondary indexes, and finally saves the updated fingerprints for the next build.

{{< image src="/images/build.png" alt="Whiskey build pipeline" class="site-main index-page" >}}


---

## Dependency Graph (`internal/graph`)

The dependency graph is the core data structure used throughout Whiskey. It enables dependency-aware incremental builds by recording relationships between content, templates, assets and remote sources.

Whiskey represents the site structure as a directed graph. Nodes can be:
- `PageNode`: Markdown source pages.
- `SourceNode`: Inlined resources (either local files or remote URLs).
- `LayoutNode`: Template files (e.g. `base.html`, `page.html`).
- `PartialNode`: Modular layout snippets (e.g. `header.html`).
- `AssetNode`: Static files under `static/` directories.

An edge is added from a node that *depends on* another node to the node it *requires* (e.g. `PageNode` $\rightarrow$ `SourceNode`, or `LayoutNode` $\rightarrow$ `PartialNode`).

---

## Fingerprints (`internal/fingerprint`)

Each node in the dependency graph is fingerprinted based on its current state.

During a build, Whiskey compares the current fingerprints with those stored from the previous build to determine which nodes have changed.

These fingerprints are stored in:

```text
.whiskey/fingerprints.json
```

## Remote Content & Materialization (`internal/source`)

Remote sources (e.g., HTTP/HTTPS files) are materialized into local workspace storage under `.whiskey/workspace/`.

1. **Resolution**: The URL is normalized into a safe workspace filename (e.g. hashing the URL).
2. **Materialization**: Whiskey sends a conditional HTTP request using stored `ETag` or `Last-Modified` headers. If the server returns `304 Not Modified`, the cached file is reused. Otherwise, the new file content is written to the workspace.
3. **Manifest**: A manifest file (`.whiskey/manifest.json`) records the mapping of remote URLs to workspace files and content hashes.
4. **Garbage Collection**: Stale workspace entries that are no longer referenced in the dependency graph are deleted during the build.

---

## Incremental Planner (`internal/planner`)

Once changed nodes have been identified, the planner determines which pages need to be rebuilt.

The process is:

```text
Changed nodes
      ↓
Full rebuild required?
      ↓
     no
      ↓
Traverse dependency graph
      ↓
Collect affected pages
      ↓
Render dirty pages
```

Changes to shared layouts or templates may require rebuilding the entire site, while edits to individual pages or includes typically rebuild only the affected pages.

---

## Rendering Pipeline

Rendering transforms Markdown into static HTML.

```text
Markdown
    ↓
Resolve Includes
    ↓
Expand Shortcodes
    ↓
Markdown → HTML
    ↓
Apply Layout
    ↓
Write HTML
```

Rendering converts Markdown files into static HTML:
1. **Include Resolution**: Reads the Markdown file and uses `dependency.ResolveIncludes` to inline `@include` contents.
2. **Shortcode Expansion**: Replaces shortcode tags (like `image` or `youtube`) with template snippets.
3. **Markdown Conversion**: Uses a Markdown-to-HTML parser.
4. **Template Execution**: Executes Go templates matching the page's specified layout, wrapping the output inside the theme's base layout.

---

## Collections & RSS Generation

Once dirty pages have been rendered, Whiskey regenerates secondary index pages:
1. **Collections**: Rebuilds list templates for folders. For pages with `collection: blog`, it gathers all matching pages and renders a sorted directory index.
2. **RSS**: Parses the `rss` section in `whiskey.toml`. It collects items in specified collections, wraps them in standard RSS XML channel formats, and writes the feed to `dist/feed.xml`.
3. **Sitemap**: Creates `dist/sitemap.xml` detailing the URLs and modification dates of all public pages.

---

## Testing

Ensure tests pass before committing changes:

```bash
# Run unit tests
go test ./...

# Run integration tests
go test -tags=integration ./...
```

---

## Contributing

1. Fork the repository on GitHub.
2. Create a feature branch (`git checkout -b feature/my-cool-feature`).
3. Implement your changes, ensuring code is formatted with `go fmt`.
4. Run static lint checks: `golangci-lint run`.
5. Open a Pull Request detailing the changes and verification steps.

---

## Further Reading

Interested in the design decisions behind Whiskey?

Continue to the [Blog](/blog/) section for deeper articles covering the architecture, dependency graph, incremental build strategy and remote content system.
