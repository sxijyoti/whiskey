---
title: Getting Started
draft: true
---

## Installation

Clone the repository and build Whiskey:

```bash
git clone https://github.com/sxijyoti/whiskey.git
cd whiskey
make build && make install
```

Verify the installation:

```bash
whiskey --help
```

The `whiskey` binary is installed into your Go binary directory (typically `~/go/bin`). Make sure this directory is included in your shell's `PATH`.

---

## Creating a Site

Every Whiskey project consists of a site directory containing your content, assets and configuration.

A minimal project looks like this:

```text
my-site/
├── content/
│   └── index.md
├── static/
└── whiskey.toml
```

### Configuration (`whiskey.toml`)

Create a `whiskey.toml` file in the site root:

```toml
title = "My Site"
description = "Built with Whiskey"

base_url = "http://localhost:8080"

theme = "minimal"
```

---

## Project Structure

A typical Whiskey site contains:

- **`whiskey.toml`** — Site configuration.
- **`content/`** — Markdown pages.
- **`layouts/`** — Custom HTML templates.
- **`static/`** — Images, CSS, JavaScript and other static assets.
- **`themes/`** — Reusable themes.

---

## Your First Page

Create `content/index.md`:

````markdown
---
title: Welcome
layout: page
draft: false
---

# Hello, Whiskey!

Welcome to my first Whiskey site.
`````

---

## Build the Site

Build your site with:

```bash
whiskey build
```

Whiskey automatically detects what has changed and rebuilds only the pages affected by those changes.

To force a complete rebuild:

```bash
whiskey build --full
```

---

## Serve Locally

Start the built-in development server:

```bash
whiskey serve
```

By default, Whiskey serves your site at **[http://localhost:8080](http://localhost:8080)** and watches for changes to content, layouts and static assets, rebuilding affected pages automatically.

To use a different port:

```bash
whiskey serve --port 9000
```

---

## Remote Content

Whiskey can include Markdown from remote HTTP or HTTPS sources.

To fetch or update all configured remote content without building the site:

```bash
whiskey sync
```

Downloaded content is cached locally so it can be reused across builds.

---

## Offline Mode

When you're offline, Whiskey can build using cached remote content without making network requests.

Build using cached sources:

```bash
whiskey build --offline
```

Or start the development server in offline mode:

```bash
whiskey serve --offline
```

---

## Next Steps

You now have a basic Whiskey site up and running.

Continue with the [documentation](/docs/) to learn about:

* Layouts and themes
* Local and remote includes
* Collections and tags
* RSS and sitemap generation
* Dependency-aware incremental builds


