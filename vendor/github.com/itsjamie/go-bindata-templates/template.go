package binhtml

import (
	"html/template"
	"path/filepath"
)

type AssetFunc func(string) ([]byte, error)
type AssetDirFunc func(string) ([]string, error)

type BinTemplate struct {
	Asset    AssetFunc
	AssetDir AssetDirFunc
}

func New(a AssetFunc, b AssetDirFunc) *BinTemplate {
	return &BinTemplate{Asset: a, AssetDir: b}
}

func (t *BinTemplate) LoadDirectory(directory string) (*template.Template, error) {
	var tmpl *template.Template

	files, err := t.AssetDir(directory)
	if err != nil {
		return tmpl, err
	}

	for _, filePath := range files {
		contents, err := t.Asset(directory + "/" + filePath)
		if err != nil {
			return tmpl, err
		}

		name := filepath.Base(filePath)

		if tmpl == nil {
			tmpl = template.New(name)
		}

		if name != tmpl.Name() {
			tmpl = tmpl.New(name)
		}

		if _, err = tmpl.Parse(string(contents)); err != nil {
			return tmpl, err
		}
	}

	return tmpl, nil
}

func (t *BinTemplate) MustLoadDirectory(directory string) *template.Template {
	if tmpl, err := t.LoadDirectory(directory); err != nil {
		panic(err)
	} else {
		return tmpl
	}
}
