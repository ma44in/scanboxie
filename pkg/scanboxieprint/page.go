package scanboxieprint

import (
	"bytes"
	"fmt"
	"html/template"
	"path"
	"path/filepath"
	"scanboxie/pkg/scanboxie"
)

var funcMap = template.FuncMap{
	"mod": func(i, j int) bool { return i%j == 0 },
	"mul": func(i, j int) int { return i * j },
	"add": func(i, j int) int { return i + j },
}

type Pagetype struct {
	Name               string
	MaxBarcodesPerPage int
	template           *template.Template
}

func NewPagetype(name string, maxBarcodesPerPage int, templateDir string) *Pagetype {
	templateFilepath := filepath.Join(templateDir, fmt.Sprintf("%s.html", name))

	tmpl := template.New(path.Base(templateFilepath)).Funcs(funcMap)
	var err error
	tmpl, err = tmpl.ParseFiles(templateFilepath)
	if err != nil {
		panic("could not parse template")
	}

	var pagetype Pagetype
	pagetype.Name = name
	pagetype.template = tmpl
	pagetype.MaxBarcodesPerPage = maxBarcodesPerPage

	return &pagetype
}

type Page struct {
	Pagetype       Pagetype
	BarcodeActions []scanboxie.BarcodeAction
}

// GetContent returns templated content of this page
func (page Page) GetContent() template.HTML {
	var tpl bytes.Buffer
	if err := page.Pagetype.template.Execute(&tpl, page); err != nil {
		panic("Template execute error:" + err.Error())
	}

	return template.HTML(tpl.String())
}

func (page Page) GetBarcodeSvgByIdx(idx int) template.HTML {
	if len(page.BarcodeActions) <= idx {
		return ""
	}

	return template.HTML(page.BarcodeActions[idx].GetBarcodeSvg())
}

func (page Page) GetBarcodeByIdx(idx int) string {
	if len(page.BarcodeActions) <= idx {
		return ""
	}

	return page.BarcodeActions[idx].Barcode
}
