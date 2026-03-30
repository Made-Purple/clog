package versionfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUpdateVersionFile(t *testing.T) {
	t.Run("writes version with newline", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "VERSION")
		if err := UpdateVersionFile(path, "1.2.3"); err != nil {
			t.Fatal(err)
		}
		got, _ := os.ReadFile(path)
		if string(got) != "1.2.3\n" {
			t.Errorf("got %q, want %q", got, "1.2.3\n")
		}
	})

	t.Run("strips v prefix", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "VERSION")
		if err := UpdateVersionFile(path, "v2.0.0"); err != nil {
			t.Fatal(err)
		}
		got, _ := os.ReadFile(path)
		if string(got) != "2.0.0\n" {
			t.Errorf("got %q, want %q", got, "2.0.0\n")
		}
	})

	t.Run("overwrites existing content", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "VERSION")
		os.WriteFile(path, []byte("0.0.1\n"), 0644)
		if err := UpdateVersionFile(path, "1.0.0"); err != nil {
			t.Fatal(err)
		}
		got, _ := os.ReadFile(path)
		if string(got) != "1.0.0\n" {
			t.Errorf("got %q, want %q", got, "1.0.0\n")
		}
	})
}

func TestUpdatePackageJSON(t *testing.T) {
	t.Run("updates version field", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "package.json")
		os.WriteFile(path, []byte(`{
  "name": "my-pkg",
  "version": "0.0.1",
  "description": "test"
}
`), 0644)
		if err := UpdatePackageJSON(path, "1.2.3"); err != nil {
			t.Fatal(err)
		}
		got, _ := os.ReadFile(path)
		want := `{
  "name": "my-pkg",
  "version": "1.2.3",
  "description": "test"
}
`
		if string(got) != want {
			t.Errorf("got:\n%s\nwant:\n%s", got, want)
		}
	})

	t.Run("strips v prefix", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "package.json")
		os.WriteFile(path, []byte(`{"version": "0.0.1"}`), 0644)
		if err := UpdatePackageJSON(path, "v3.0.0"); err != nil {
			t.Fatal(err)
		}
		got, _ := os.ReadFile(path)
		if string(got) != `{"version": "3.0.0"}` {
			t.Errorf("got %q", got)
		}
	})

	t.Run("preserves formatting and field order", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "package.json")
		original := "{\n\t\"name\": \"test\",\n\t\"version\": \"1.0.0\",\n\t\"main\": \"index.js\"\n}\n"
		os.WriteFile(path, []byte(original), 0644)
		if err := UpdatePackageJSON(path, "2.0.0"); err != nil {
			t.Fatal(err)
		}
		got, _ := os.ReadFile(path)
		want := "{\n\t\"name\": \"test\",\n\t\"version\": \"2.0.0\",\n\t\"main\": \"index.js\"\n}\n"
		if string(got) != want {
			t.Errorf("got:\n%s\nwant:\n%s", got, want)
		}
	})

	t.Run("errors on missing version field", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "package.json")
		os.WriteFile(path, []byte(`{"name": "test"}`), 0644)
		err := UpdatePackageJSON(path, "1.0.0")
		if err == nil {
			t.Fatal("expected error for missing version field")
		}
	})

	t.Run("only replaces first version match", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "package.json")
		content := `{
  "version": "1.0.0",
  "dependencies": {
    "foo": {
      "version": "3.0.0"
    }
  }
}
`
		os.WriteFile(path, []byte(content), 0644)
		if err := UpdatePackageJSON(path, "2.0.0"); err != nil {
			t.Fatal(err)
		}
		got, _ := os.ReadFile(path)
		want := `{
  "version": "2.0.0",
  "dependencies": {
    "foo": {
      "version": "3.0.0"
    }
  }
}
`
		if string(got) != want {
			t.Errorf("got:\n%s\nwant:\n%s", got, want)
		}
	})
}
