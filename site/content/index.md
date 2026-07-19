---
title: Whiskey
description: A dependency-aware static site generator that builds only what changed.
date: 2026-07-13
layout: index
draft: false
---

# Whiskey

A dependency-aware static site generator written in Go.

{{< image src="/images/logo.png" alt="Whiskey Logo" class="site-main index-page" >}}

---

## What is Whiskey?

Whiskey is a static site generator built around a simple idea: **only rebuild what changed**.

Instead of rebuilding an entire site after every edit, Whiskey constructs a dependency graph connecting pages, layouts, partials, assets and external content. When a file changes, Whiskey determines the minimal set of affected pages and rebuilds only those.

The result is faster builds, better scalability and a development workflow that remains responsive even as sites grow.

In simple terms, Whiskey asks:

> **"Which pages actually depend on what changed?"**

Then it rebuilds only those pages.

---

## Why Whiskey?

Modern static websites are no longer made up of just Markdown files. They often combine local content with shared documents, generated pages and remote resources.

Whiskey was designed to make these workflows first-class while keeping the development experience simple.

Its goals are:

- **Incremental by default** — rebuild only affected pages instead of the entire site.
- **Dependency-aware** — understand relationships between pages, templates, assets and external sources.
- **Remote Content** — treat remote Markdown as native content through built-in `@include` command and workspace caching.
- **Developer-friendly** — provide a fast CLI, live reload and deterministic builds.

---

## Features

- **Dependency-aware incremental builds**
- **Local includes for reusable content**
- **Remote includes with workspace caching**
- **Offline builds using cached sources**
- **Development server with live reload**
- **Layouts, partials and themes**
- **Collections and tags**
- **RSS feeds and sitemap generation**
- **Built-in image shortcode**

---

## Learn More

- **[Documentation](/docs/)** — Learn how to configure and use Whiskey.
- **[Developer Guide](/developer-guide/)** — Explore Whiskey's architecture and implementation.
- **[Blogs](/blog/)** — Design philosophy, architecture and implementation notes.
- **[Tags](/tags/)** — Browse topics across the documentation.

---

## Getting Started

The guide below is included directly from `getting-started.md` using Whiskey's local include feature. The source page remains a draft, allowing the same content to be reused without being published independently.

@include local:getting-started.md