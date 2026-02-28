package wiki

import (
	"os"
	"path/filepath"
	"strings"
)

func titleToInternalTitle(title string) string {
	return strings.ReplaceAll(title, " ", "_")
}

func internalTitleToTitle(filename string) string {
	return strings.ReplaceAll(filename, "_", " ")
}

// Saves a page struct to a file
// The file is named after the page's title, and the content is the just a dump of the page's body
func (p *Page) save() error {
	if err := os.MkdirAll(savedPagesPath, 0o755); err != nil {
		return err
	}
	fullpath := filepath.Join(savedPagesPath, p.InternalTitle)
	return os.WriteFile(fullpath, p.Body, 0600)
}

// Loads a page from a file.
// The file is named after the page's title, and the body is the file content as a byte array
func loadPage(internalTitle string) (*Page, error) {
	fullpath := filepath.Join(savedPagesPath, internalTitle)
	body, err := os.ReadFile(fullpath)
	if err != nil {
		return nil, err
	}
	return &Page{
		Title:         internalTitleToTitle(internalTitle),
		InternalTitle: internalTitle,
		Body:          body,
	}, nil
}
