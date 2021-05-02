package scanboxieprint

import (
	"bytes"
	"html/template"
	"io/fs"
	"scanboxie/pkg/scanboxie"

	log "github.com/sirupsen/logrus"
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

func NewPagetype(name string, maxBarcodesPerPage int, templateFS fs.FS) *Pagetype {
	tmpl, err := template.ParseFS(templateFS, name+".html")
	if err != nil {
		log.Fatalf("could not parse template for page type %s, error: %v\n", name, err)
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
	log.Debugf("GetContent() for page %v\n", page)

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

	barcode := (page.BarcodeActions[idx]).Barcode

	return GetBarcodeSvg(barcode)
}

func (page Page) GetBarcodeByIdx(idx int) string {
	if len(page.BarcodeActions) <= idx {
		return ""
	}

	return page.BarcodeActions[idx].Barcode
}
