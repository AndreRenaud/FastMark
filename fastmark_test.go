package main

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/AndreRenaud/fastmark/storage"
)

func TestListImages(t *testing.T) {
	tests := []struct {
		name  string
		files []string
		want  []string
	}{
		{
			name:  "flat",
			files: []string{"images/b.png", "images/a.jpg", "images/notes.txt"},
			want:  []string{"a.jpg", "b.png"},
		},
		{
			name:  "grouped",
			files: []string{"images/cam-2/b.jpg", "images/cam-1/a.JPEG", "images/cam-1/notes.txt"},
			want:  []string{filepath.Join("cam-1", "a.JPEG"), filepath.Join("cam-2", "b.jpg")},
		},
		{
			name:  "mixed",
			files: []string{"images/a.png", "images/cam-1/b.jpg"},
			want:  []string{"a.png", filepath.Join("cam-1", "b.jpg")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			for _, f := range tt.files {
				path := filepath.Join(dir, f)
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, nil, 0644); err != nil {
					t.Fatal(err)
				}
			}
			got, err := listImages(storage.NewStorage(dir))
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("listImages() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReadImageList(t *testing.T) {
	dir := t.TempDir()
	list := "# training subset\n" +
		"b.png\n" +
		"\n" +
		"  cam-1/a.jpg  \n" +
		"images/c.jpg\n" +
		"b.png\n"
	if err := os.MkdirAll(filepath.Join(dir, "lists"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "lists", "train.txt"), []byte(list), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := readImageList(storage.NewStorage(dir), filepath.Join("lists", "train.txt"))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"b.png", filepath.Join("cam-1", "a.jpg"), "c.jpg"}
	if !slices.Equal(got, want) {
		t.Errorf("readImageList() = %q, want %q", got, want)
	}

	if _, err := readImageList(storage.NewStorage(dir), "missing.txt"); err == nil {
		t.Error("readImageList() of a missing file succeeded, want error")
	}
}

func TestLabelFileFor(t *testing.T) {
	tests := map[string]string{
		"a.jpg":                           filepath.Join("labels", "a.txt"),
		filepath.Join("cam-1", "b.jpeg"): filepath.Join("labels", "cam-1", "b.txt"),
	}
	for image, want := range tests {
		if got := labelFileFor(image); got != want {
			t.Errorf("labelFileFor(%q) = %q, want %q", image, got, want)
		}
	}
}
