# whiskey

**A dependency-aware static site generator written in Go.**

`whiskey` is a modern static site generator built around **dependency-aware incremental builds**.

Instead of rebuilding your entire site whenever something changes, `whiskey` constructs a dependency graph connecting content, layouts, partials, assets, and remote sources. When a change occurs, only the pages affected by that change are rebuilt.

Unlike traditional static site generators that fetch remote content during rendering, `whiskey` synchronizes external sources into a persistent workspace before rendering begins. This enables deterministic builds, offline support, conditional HTTP caching, and reproducible outputs.

---

## Why `whiskey`?

Most static site generators think in terms of files.

Whiskey thinks in terms of **dependencies**.

```text
Markdown
     │
     ▼
Dependency Graph
     │
     ▼
Dirty Set
     │
     ▼
Incremental Build
```

Instead of rebuilding the world, Whiskey asks:

> **"Which pages actually depend on what changed?"**

That philosophy drives every part of the architecture.

## Commands
```
Usage:
  whiskey [command]

Available Commands:
  build       Build the site
  check       Inspect site dependencies
  clean       Clean generated artifacts
  help        Help about any command
  serve       Build and serve a Whiskey site
  sync        Synchronize remote sources
  theme       Manage local themes
  version     Version information about whiskey

Flags:
  -h, --help   help for whiskey

Use "whiskey [command] --help" for more information about a command.
```

---

# Features

## Build System

* Dependency-aware incremental builds
* Graph-driven dirty propagation
* SHA256 fingerprinting
* Configuration-aware rebuilds
* Pretty URLs
* Incremental asset pipeline
* Live reload development server

---

## Content

* Markdown + YAML frontmatter
* Draft support
* Multi-page sites
* Collections
* Tags
* RSS feeds
* Sitemap generation
* Image shortcodes
* YouTube embeds

---

## Templates

* Layout inheritance
* Base templates
* Partials
* Theme fallback
* Four bundled themes (`minimal`, `paper`, `terminal`, `noir`)

---

## Remote Content

* Remote Markdown includes
* Local includes
* Persistent materialized workspace
* Offline builds
* Conditional HTTP caching
* Workspace garbage collection
* HTML → Markdown extraction

---

## Developer Experience

* Live reload
* Recursive file watching
* Automatic rebuilds
* Browser refresh
* Dependency graph visualization
* Remote synchronization

---

# Architecture

Whiskey separates **source synchronization** from **page rendering**.

```text
                Markdown
                   │
                   ▼
           Dependency Graph
                   │
                   ▼
          Dirty Set Calculation
                   │
                   ▼
        Source Materialization
                   │
                   ▼
          .whiskey/workspace/
                   │
                   ▼
           Template Rendering
                   │
                   ▼
                 dist/
```

Because rendering consumes only workspace artifacts, builds become deterministic and independent of network availability.

---

# Installation

## Requirements

* Go 1.25+

Clone the repository:

```bash
git clone https://github.com/sxijyoti/whiskey.git
cd whiskey
```

To run:
```bash
# Build:
make build

# Install:
make install
```

or simply,
```bash
make whiskey 
# this handles build and install
```

Verify:

```bash
whiskey version
```

---

# Quick Start

Create a new site:

```text
site/
├── whiskey.toml
├── content/
├── layouts/
├── static/
└── themes/
```

```bash
# Build:
whiskey build

# Serve:
whiskey serve

# Synchronize remote sources:
whiskey sync

# Offline build:
whiskey build --offline
```

---

# Project Structure

```text
site/
├── whiskey.toml
├── content/
├── layouts/
├── static/
├── themes/
│
├── dist/
│
└── .whiskey/
    ├── fingerprints.json
    ├── manifest.json
    └── workspace/
```

---

# Dependency Graph

Everything inside Whiskey is represented as graph nodes.

Current node types include:

* Pages
* Layouts
* Partials
* Assets
* Remote Sources

Example:

```text
Page
 │
 ▼
Layout
 │
 ▼
Base
 │
 ▼
Partials
```

and

```text
Page
 │
 ▼
Remote Source
```

Whenever a node changes, Whiskey traverses the graph to determine which pages become dirty.

---

# Incremental Builds

Every build input is fingerprinted.

Tracked inputs include:

* Markdown
* Layouts
* Partials
* Assets
* Configuration
* Remote Sources

Build flow:

```text
Changed Node
      │
      ▼
Dependency Graph
      │
      ▼
Affected Pages
      │
      ▼
Incremental Build
```

Only affected pages are rebuilt.

---

# Remote Includes

Include external Markdown directly inside your content.

```md
@include https://raw.githubusercontent.com/sxijyoti/whiskey/main/README.md
```

Pipeline:

```text
Remote URL
      │
      ▼
Provider
      │
      ▼
Materialization
      │
      ▼
Workspace
      │
      ▼
Rendering
```

---

# Workspace-backed Sources

Remote content is synchronized before rendering.

```text
Source Provider
        │
        ▼
Fetch
        │
        ▼
Workspace
        │
        ▼
Page Rendering
```

Pages never perform network requests directly.

This enables:

* deterministic builds
* reproducible outputs
* offline builds
* reusable cached content

---

# Conditional HTTP Requests

Whiskey minimizes unnecessary downloads through HTTP cache validation.

```text
Metadata
     │
     ▼
Conditional Request
     │
     ├── 304 Not Modified
     │         │
     │         ▼
     │   Reuse Workspace
     │
     └── 200 OK
               │
               ▼
        Update Workspace
```

Supported headers:

* ETag
* Last-Modified

---

# HTML Extraction

HTML sources are automatically normalized into Markdown.

```text
HTTP
   │
   ▼
Content-Type
   │
   ├── text/html
   │
   ▼
Defuddle
   │
   ▼
Markdown
   │
   ▼
Workspace
```

This keeps the rendering pipeline content-format agnostic.

---

# Offline Builds

Synchronize once:

```bash
whiskey sync
```

Then build anywhere.

```bash
whiskey build --offline
```

Remote sources are served entirely from the workspace.

---

# Themes

Whiskey ships with four bundled themes.

* **minimal** — Clean and minimal
* **terminal** — Hacker-inspired terminal aesthetic
* **paper** — Traditional document style
* **noir** — Dark editorial layout

Themes consist of:

* layouts
* partials
* static assets


---