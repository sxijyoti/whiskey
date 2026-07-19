---
title: Introducing Whiskey
description: Why Whiskey was built and the ideas behind its dependency-aware approach to static site generation.
date: 2026-07-13
layout: post
collection: blog
tags:
  - architecture
  - introduction
  - philosophy
---

Most static site generators follow the same workflow: read every page, apply templates, and regenerate the entire site.

For small websites this works well. As projects grow, however, rebuild times increase even when only a single page changes.

Whiskey began with a simple question:

> **Why rebuild the entire site when only a handful of pages actually changed?**

That question became the foundation for everything that followed.

---

## The Problem

A modern documentation site is more than a collection of Markdown files.

Pages depend on:

- layouts
- partials
- images
- stylesheets
- local includes
- remote documentation
- generated indexes
- RSS feeds

These relationships form a graph.

Most static site generators either rebuild everything or rely on coarse invalidation strategies that rebuild much more than necessary.

---

## A Dependency-Aware Approach

Whiskey models a site as a dependency graph.

Every page, template, asset and include becomes a node. Relationships between them become edges.

Instead of asking:

> "Which files changed?"

Whiskey asks:

> **"Which pages depend on what changed?"**

Only those pages are rebuilt.

The goal isn't simply faster builds—it is to perform the minimum amount of work while producing exactly the same output as a clean rebuild.

---

## Remote Content as a First-Class Citizen

Modern documentation rarely lives in one repository.

A page might include:

- a README from GitHub
- API documentation
- generated Markdown
- shared documentation maintained elsewhere

Instead of requiring external scripts, Whiskey treats remote content as part of the dependency graph.

Remote sources are downloaded, cached locally and refreshed only when they change, allowing them to participate in incremental builds just like local files.

---

## Building Whiskey

Whiskey is written in Go and built around a small number of composable components:

- a dependency graph
- fingerprint-based change detection
- an incremental build planner
- a rendering pipeline
- a remote content workspace

Each component has a single responsibility, making the codebase easier to understand and extend.

---

## Dogfooding the Documentation

This website is built with Whiskey itself.

Several pages are reused through local includes, the documentation demonstrates the features it describes, and every edit exercises the incremental build system during development.

Rather than creating examples solely for documentation, Whiskey uses its own capabilities to build its documentation.

---

## Looking Ahead

Whiskey is still evolving.

Future work will continue to improve performance, expand supported content sources and simplify site authoring while preserving the project's central idea:

> Build only what changed.

Everything else follows from that principle.