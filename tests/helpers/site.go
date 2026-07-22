package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type SiteBuilder struct {
	t     *testing.T
	Dir   string
	files map[string]string
}

// NewTestSite creates a new temporary site directory for testing.
func NewTestSite(t *testing.T) *SiteBuilder {
	return &SiteBuilder{
		t:     t,
		Dir:   t.TempDir(),
		files: make(map[string]string),
	}
}

// WriteFile stages a file content to be written to disk.
func (s *SiteBuilder) WriteFile(path, content string) *SiteBuilder {
	s.files[path] = content
	return s
}

// WriteConfig stages the whiskey.toml configuration file.
func (s *SiteBuilder) WriteConfig(content string) *SiteBuilder {
	return s.WriteFile("whiskey.toml", content)
}

// WritePage stages a content markdown page.
func (s *SiteBuilder) WritePage(path, content string) *SiteBuilder {
	return s.WriteFile(filepath.Join("content", path), content)
}

// WriteLayout stages a layout file under layout/ or themes/<theme>/layouts/.
func (s *SiteBuilder) WriteLayout(theme, name, content string) *SiteBuilder {
	if theme == "" {
		return s.WriteFile(filepath.Join("layouts", name), content)
	}
	return s.WriteFile(filepath.Join("themes", theme, "layouts", name), content)
}

// WritePartial stages a partial HTML template.
func (s *SiteBuilder) WritePartial(theme, name, content string) *SiteBuilder {
	if theme == "" {
		return s.WriteFile(filepath.Join("layouts", "partials", name), content)
	}
	return s.WriteFile(filepath.Join("themes", theme, "layouts", "partials", name), content)
}

// WriteAsset stages a static asset.
func (s *SiteBuilder) WriteAsset(theme, name, content string) *SiteBuilder {
	if theme == "" {
		return s.WriteFile(filepath.Join("static", name), content)
	}
	return s.WriteFile(filepath.Join("themes", theme, "static", name), content)
}

// WriteRemoteWorkspace simulates a cached remote source in the workspace.
// It writes the file to .whiskey/workspace/ and registers it in manifest.json.
func (s *SiteBuilder) WriteRemoteWorkspace(ref, content string) *SiteBuilder {
	refHasher := sha256.New()
	refHasher.Write([]byte(ref))
	workspaceName := hex.EncodeToString(refHasher.Sum(nil)) + ".md"
	workspacePath := filepath.Join(".whiskey", "workspace", workspaceName)

	s.files[workspacePath] = content

	// Register in manifest.json
	manifestPath := filepath.Join(".whiskey", "manifest.json")
	var manifest struct {
		Version int `json:"version"`
		Sources map[string]struct {
			Workspace   string            `json:"workspace"`
			ContentHash string            `json:"content_hash"`
			State       map[string]string `json:"state"`
		} `json:"sources"`
	}

	manifest.Version = 1
	manifest.Sources = make(map[string]struct {
		Workspace   string            `json:"workspace"`
		ContentHash string            `json:"content_hash"`
		State       map[string]string `json:"state"`
	})

	// If manifest.json was already staged, parse it first
	if existing, ok := s.files[manifestPath]; ok {
		_ = json.Unmarshal([]byte(existing), &manifest)
	}

	contentHasher := sha256.New()
	contentHasher.Write([]byte(content))
	contentHash := hex.EncodeToString(contentHasher.Sum(nil))

	manifest.Sources[ref] = struct {
		Workspace   string            `json:"workspace"`
		ContentHash string            `json:"content_hash"`
		State       map[string]string `json:"state"`
	}{
		Workspace:   workspaceName,
		ContentHash: contentHash,
		State: map[string]string{
			"etag":          `"v1"`,
			"last_modified": "Mon, 02 Jan 2006 15:04:05 MST",
		},
	}

	manifestBytes, _ := json.MarshalIndent(manifest, "", "    ")
	s.files[manifestPath] = string(manifestBytes)

	return s
}

// Build writes all staged files to disk.
func (s *SiteBuilder) Build() {
	for path, content := range s.files {
		s.WriteToDisk(path, content)
	}
}

// WriteToDisk writes a file directly to the workspace temp dir.
func (s *SiteBuilder) WriteToDisk(path, content string) {
	fullPath := filepath.Join(s.Dir, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		s.t.Fatalf("failed to create dir: %v", err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		s.t.Fatalf("failed to write file: %v", err)
	}
}

// EditPage writes a modified content page directly to disk.
func (s *SiteBuilder) EditPage(path, content string) {
	s.WriteToDisk(filepath.Join("content", path), content)
}

// DeletePage deletes a page file from disk.
func (s *SiteBuilder) DeletePage(path string) {
	err := os.Remove(filepath.Join(s.Dir, "content", path))
	if err != nil {
		s.t.Fatalf("failed to delete page: %v", err)
	}
}

// Theme sets up a minimal theme layout structure.
func (s *SiteBuilder) Theme(name string) *SiteBuilder {
	s.WriteFile(filepath.Join("themes", name, "theme.toml"), `name="`+name+`"`)
	s.WriteLayout(name, "base.html", `<!DOCTYPE html><html><body>{{ block "content" . }}{{ end }}</body></html>`)
	s.WriteLayout(name, "page.html", `{{ define "content" }}{{ .Content }}{{ end }}`)
	s.WriteLayout(name, "list.html", `
	{{ define "content" }}
	<h1>{{ .Title }}</h1>

	{{ range .Pages }}
	<article>
		<h2>{{ .Title }}</h2>
	</article>
	{{ end }}

	{{ end }}
	`)
	s.WriteLayout(name, "tags.html", `{{ define "content" }}{{ .Content }}{{ end }}`)
	s.WriteLayout(name, "index.html", `{{ define "content" }}{{ .Content }}{{ end }}`)
	return s
}
