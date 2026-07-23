---
title: Deployment Guide
description: Learn how to deploy Whiskey sites to GitHub Pages, Netlify, Cloudflare Pages, Vercel, and CI/CD platforms using the official installer.
unlisted: true
layout: page
---

# Deployment Guide

Deploying a `whiskey` site is simple, fast, and platform-independent. Because `whiskey` compiles your site into standard static HTML, CSS, and assets inside the `dist/` directory, it can be hosted on any static web host or CDN.

---

## Official Installation Script

Every deployment strategy relies on the single official installer script. This installs the pre-compiled `whiskey` binary for Linux or macOS directly into `~/.local/bin`:

```bash
curl -fsSL https://raw.githubusercontent.com/sxijyoti/whiskey/main/scripts/install.sh | sh
```

You can also specify a specific version tag:

```bash
curl -fsSL https://raw.githubusercontent.com/sxijyoti/whiskey/main/scripts/install.sh | sh -s -- v0.1.0
```

---

## GitHub Pages

Deploy your Whiskey site automatically using **GitHub Actions**.

Create `.github/workflows/deploy.yml` in your repository:

```yaml
name: Deploy Whiskey Site to GitHub Pages

on:
  push:
    branches:
      - main
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: "pages"
  cancel-in-progress: false

jobs:
  build-and-deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Install Whiskey
        run: |
          curl -fsSL https://raw.githubusercontent.com/sxijyoti/whiskey/main/scripts/install.sh | sh
          echo "$HOME/.local/bin" >> $GITHUB_PATH

      - name: Build site
        run: whiskey build .

      - name: Upload Pages Artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: dist

      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

> **Note:** Ensure GitHub Pages is configured under repository **Settings → Pages → Source → GitHub Actions**.

---

## Netlify

Deploy to Netlify by creating a `netlify.toml` file in the root of your repository:

```toml
[build]
command = """
curl -fsSL https://raw.githubusercontent.com/sxijyoti/whiskey/main/scripts/install.sh | sh
export PATH="$HOME/.local/bin:$PATH"
whiskey build .
"""
publish = "dist"
```

Netlify will automatically execute the build command during deployment and publish the `dist/` directory.

---

## Cloudflare Pages

Configure Cloudflare Pages via the Cloudflare Dashboard:

1. Go to **Workers & Pages → Create Application → Pages → Connect to Git**.
2. Select your Whiskey site repository.
3. Configure the Build settings:
   - **Framework preset**: None
   - **Build command**:
     ```bash
     curl -fsSL https://raw.githubusercontent.com/sxijyoti/whiskey/main/scripts/install.sh | sh && export PATH="$HOME/.local/bin:$PATH" && whiskey build .
     ```
   - **Build output directory**: `dist`
4. Click **Save and Deploy**.

---

## Vercel

Deploy to Vercel by setting up a project or configuring `vercel.json` in your repository:

```json
{
  "buildCommand": "curl -fsSL https://raw.githubusercontent.com/sxijyoti/whiskey/main/scripts/install.sh | sh && export PATH=\"$HOME/.local/bin:$PATH\" && whiskey build .",
  "outputDirectory": "dist"
}
```

Alternatively, configure the **Build Command** and **Output Directory** in the Vercel Dashboard under **Project Settings → Build & Development Settings**.

---

## Generic CI/CD Pipelines

For any standard Linux or macOS build agent (GitLab CI, CircleCI, Bitbucket Pipelines, Jenkins, or Docker containers):

```bash
# 1. Install official Whiskey binary
curl -fsSL https://raw.githubusercontent.com/sxijyoti/whiskey/main/scripts/install.sh | sh

# 2. Ensure ~/.local/bin is in PATH
export PATH="$HOME/.local/bin:$PATH"

# 3. Build site
whiskey build .
```

The compiled static site will be generated in `dist/`. Publish or upload the `dist/` directory to your web server or hosting service.
