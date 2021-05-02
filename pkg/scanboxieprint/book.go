package scanboxieprint

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"scanboxie/pkg/scanboxie"

	log "github.com/sirupsen/logrus"
)

type Book struct {
	template  *template.Template
	pages     []*Page
	pageTypes map[string]Pagetype
}

func NewBook(templates fs.FS, barcodeActions []*scanboxie.BarcodeAction) *Book {
	var book Book

	// book.html
	book.template = template.Must(template.New("").Funcs(funcMap).ParseFS(templates, "book.html"))

	book.pageTypes = make(map[string]Pagetype)
	book.pageTypes["3_rows_4_cols_page"] = *NewPagetype("page_3_rows_4_cols", 12, templates)
	book.pageTypes["4_rows_4_cols_page"] = *NewPagetype("page_4_rows_4_cols", 20, templates)
	book.pageTypes["single_with_cover"] = *NewPagetype("page_single_with_cover", 1, templates)
	book.pageTypes["empty"] = *NewPagetype("page_empty", 0, templates)

	pages := []*Page{}

	var currentPage *Page
	i := 0
	for _, barcodeAction := range barcodeActions {
		log.Debugf("NewBook() - Process barcode action %v\n", barcodeAction)

		barcodeActionPagetype, ok := book.pageTypes[barcodeAction.Pagetype]
		if !ok {
			fmt.Printf("Unkown PageType %s for barcode %s\n", barcodeAction.Pagetype, barcodeAction.Barcode)
			continue
		}

		if i == 0 || currentPage.Pagetype != barcodeActionPagetype {
			// Need new page of this page type
			log.Debugf("NewBook() - Create new page for pagetype %v\n", barcodeActionPagetype)
			newpage := Page{}
			newpage.Pagetype = barcodeActionPagetype
			pages = append(pages, &newpage)
			currentPage = &newpage
			i = barcodeActionPagetype.MaxBarcodesPerPage
		}

		currentPage.BarcodeActions = append(currentPage.BarcodeActions, *barcodeAction)
		i--
	}

	fmt.Printf("Count Pages: %d\n", len(pages))

	book.pages = pages

	return &book
}

func (book Book) Write(w io.Writer) error {
	err := book.template.ExecuteTemplate(w, "book.html", book)
	if err != nil {
		return err
	}

	return nil
}

// GetPages returns slice of pages in the correct order for printing
//
// 2-side-printing:
//   Reorder Pages for 2 side print
//   -> [1,2] [3,4] -> [1,3] [2,4]
func (book Book) GetPages() []*Page {
	var pages []*Page

	for i := 0; i < len(book.pages); i = i + 4 {
		pages = append(pages, book.pages[i+0])

		if (i + 2) < len(book.pages) {
			pages = append(pages, book.pages[i+2])
		} else {
			pages = append(pages, &Page{Pagetype: book.pageTypes["empty"]})
		}

		if (i + 1) < len(book.pages) {
			pages = append(pages, book.pages[i+1])
		} else {
			pages = append(pages, &Page{Pagetype: book.pageTypes["empty"]})
		}

		if (i + 3) < len(book.pages) {
			pages = append(pages, book.pages[i+3])
		} else {
			pages = append(pages, &Page{Pagetype: book.pageTypes["empty"]})
		}
	}

	return pages
}

func (book Book) GetPageCount() int {
	return len(book.pages)
}
