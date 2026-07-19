---
title: Whiskey Architecture
description: An overview of the systems that make Whiskey's dependency-aware build pipeline possible.
date: 2026-07-13
layout: post
collection: blog
tags:
  - architecture
  - internals
---

Whiskey is built around a simple idea:

> Everything that influences a page should be represented as a dependency.

Rather than treating pages, templates, assets and remote content as separate systems, Whiskey models them as parts of a single dependency graph. That graph becomes the foundation for incremental builds, remote content, live reload and secondary outputs such as RSS feeds.

This article provides a high-level overview of the major architectural components.

---

## The Big Picture

At a high level, Whiskey consists of five major systems:

```text
          Content
             │
             ▼
    Dependency Graph
             │
             ▼
     Change Detection
             │
             ▼
   Incremental Planner
             │
             ▼
        Rendering
```

Each system has a single responsibility and passes its output to the next stage.

---

## Dependency Graph

The dependency graph is the heart of Whiskey.

Every significant object becomes a node:

- pages
- layouts
- partials
- assets
- local includes
- remote includes

Edges represent dependencies.

For example:

```text
index.md
    │
    ├── getting-started.md
    │
    └── index.html
             │
             └── base.html
                      │
                      └── header.html
```

This graph allows Whiskey to answer one question efficiently:

> Which pages depend on what changed?

---

## Change Detection

Once the graph has been constructed, Whiskey determines which nodes have changed since the previous build.

Rather than immediately rebuilding pages, Whiskey first identifies the smallest possible set of changed nodes.

Those changes are then passed to the incremental planner.

---

## Incremental Planning

The planner decides whether the build can remain incremental.

For example:

- editing a page usually rebuilds only that page
- changing a local include rebuilds every page using it
- changing a shared layout rebuilds the entire site

This keeps incremental builds both fast and correct.

---

## Remote Content

Remote Markdown participates in the build exactly like local content.

When Whiskey encounters a remote include:

```markdown
@ include https://example.com/page.md
```

it downloads the content, stores it locally and adds it to the dependency graph.

From that point onward, remote content behaves like any other source file.

---

## Rendering

Once the planner has identified the pages that need rebuilding, Whiskey renders only those pages.

Each page passes through the same pipeline:

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
Write Output
```

After rendering completes, Whiskey regenerates collections, RSS feeds, sitemaps and other derived outputs as needed.

---

## Why This Architecture?

Each subsystem has a clearly defined responsibility.

- The dependency graph records relationships.
- Change detection identifies modified nodes.
- The planner determines what must be rebuilt.
- The renderer produces HTML.

Keeping these responsibilities separate makes the codebase easier to understand while allowing new features to integrate naturally with the existing build pipeline.

For example, remote content, local includes and future source providers all become additional node types in the dependency graph rather than requiring entirely separate build systems.

---

## Learn More

The architecture described here is implemented throughout Whiskey's codebase.

If you'd like to explore the implementation in more detail, continue to the **Developer Guide**, which explains the repository layout, packages and build pipeline in depth.