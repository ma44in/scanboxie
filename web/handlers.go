package web

import (
	"fmt"
	"net/http"
	"path/filepath"
	"scanboxie/pkg/scanboxie"
	"text/template"
)

type Pagetype struct {
	Name               string
	MaxBarcodesPerPage int
}

type Page struct {
	Pagetype       Pagetype
	BarcodeActions []scanboxie.BarcodeAction
}

func (page Page) GetBarcodeSvgByIdx(idx int) string {
	if len(page.BarcodeActions) <= idx {
		return ""
	}

	return page.BarcodeActions[idx].GetBarcodeSvg()
}

func (page Page) GetBarcodeByIdx(idx int) string {
	if len(page.BarcodeActions) <= idx {
		return ""
	}

	return page.BarcodeActions[idx].Barcode
}

type ViewData struct {
	Pages     []*Page
	PageCount int
}

func (webapp *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	//tmpl := webapp.Templates["index.html"]
	//tmpl := template.New("index.html")
	//tmpl, _ = tmpl.ParseFiles(filepath.Join(webapp.TemplateDir, "index.html"))

	var funcMap = template.FuncMap{
		"mod": func(i, j int) bool { return i%j == 0 },
		"mul": func(i, j int) int { return i * j },
		"add": func(i, j int) int { return i + j },
	}

	tmpl := template.Must(template.New("").Funcs(funcMap).ParseFiles(filepath.Join(webapp.TemplateDir, "index.html")))

	pageTypes := make(map[string]Pagetype)
	pageTypes["3_rows_4_cols_page"] = Pagetype{Name: "3_rows_4_cols_page", MaxBarcodesPerPage: 12}
	pageTypes["4_rows_4_cols_page"] = Pagetype{Name: "4_rows_4_cols_page", MaxBarcodesPerPage: 20}
	pageTypes["single_with_cover"] = Pagetype{Name: "single_with_cover", MaxBarcodesPerPage: 1}

	pages := []*Page{}

	var currentPage *Page
	i := 0
	for _, barcodeAction := range webapp.BarcodeConfig.BarcodeActions {
		barcodeActionPagetype, ok := pageTypes[barcodeAction.BookletPageType]
		if !ok {
			fmt.Printf("Unkown PageType %s for barcode %s\n", barcodeAction.BookletPageType, barcodeAction.Barcode)
			continue
		}

		if i == 0 || currentPage.Pagetype != barcodeActionPagetype {
			// Need new page of this page type
			newpage := Page{}
			newpage.Pagetype = barcodeActionPagetype
			pages = append(pages, &newpage)
			currentPage = &newpage
			i = barcodeActionPagetype.MaxBarcodesPerPage
		}

		fmt.Printf("Found %v\n", barcodeAction)
		currentPage.BarcodeActions = append(currentPage.BarcodeActions, barcodeAction)
		i--
	}

	fmt.Printf("Count Pages: %d\n", len(pages))

	vd := ViewData{
		Pages:     pages,
		PageCount: len(pages),
	}

	err := tmpl.ExecuteTemplate(w, "index.html", vd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
