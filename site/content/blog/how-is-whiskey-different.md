---
title: How is Whiskey Different?
description: A closer look at Whiskey's dependency-aware approach to static site generation.
date: 2026-07-13
layout: post
collection: blog
tags:
  - architecture
  - incremental
  - remote-content
---


Most static site generators follow the same workflow:

```
Read everything
        ↓
Render everything
        ↓
Write everything
```

This approach is simple and reliable, but it also means a small change can trigger work across the entire site.

Whiskey approaches the problem differently.

Instead of asking:

> **"Which files changed?"**

it asks:

> **"Which pages depend on what changed?"**

That single design decision influences almost every part of the project.

---

## Dependency-Aware Builds

A Whiskey site is represented as a dependency graph.

Pages depend on:

- layouts
- partials
- assets
- local includes
- remote includes

Rather than treating Markdown files as isolated documents, Whiskey records these relationships so they can be used during future builds.

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

If `getting-started.md` changes, only the pages that include it are rebuilt.

If `header.html` changes, every page using that layout is rebuilt.

The result is a build that performs only the work required to produce a correct site.

---

## Incremental by Default

Incremental builds are not a separate mode in Whiskey—they are the default.

Every build begins by determining what has changed since the previous build.

If only a single page changes, only that page is rebuilt.

If a shared layout changes, Whiskey automatically performs a full rebuild because every page may be affected.

This keeps the behaviour predictable while avoiding unnecessary work.

---

## Remote Content

Modern documentation often spans multiple repositories.

A page may include:

- a GitHub README
- shared documentation
- generated Markdown
- documentation maintained by another team

Rather than requiring separate download scripts, Whiskey treats remote content exactly like local content.

```markdown
@ include https://raw.githubusercontent.com/user/project/main/README.md
```

Remote sources become part of the dependency graph and participate in incremental builds just like local files.

---

## Workspace Caching

Downloaded content is cached locally inside Whiskey's workspace.

On future builds, Whiskey checks whether a remote resource has changed before downloading it again.

If nothing has changed, the cached copy is reused.

This reduces unnecessary network requests while keeping builds reproducible.

---

## Offline Development

Because remote content is cached locally, development doesn't stop when the network disappears.

Running:

```bash
whiskey build --offline
```

uses the cached workspace instead of making HTTP requests.

The same applies to the development server:

```bash
whiskey serve --offline
```

This makes remote content behave much more like local content during day-to-day development.

---

## Correctness Before Speed

The goal of Whiskey is not simply to build faster.

The goal is to perform the **minimum amount of work necessary** while producing exactly the same output as a clean rebuild.

Sometimes that means rebuilding one page.

Sometimes it means rebuilding the entire site.

The dependency graph allows Whiskey to make that decision automatically, ensuring that incremental builds remain both fast and correct.