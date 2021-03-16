package scanboxieprint

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"scanboxie/pkg/scanboxie"
)

type Book struct {
	template *template.Template
	pages    []*Page
}

func NewBook(templateDir string, barcodeActions []scanboxie.BarcodeAction) *Book {
	var book Book

	book.template = template.Must(template.New("").Funcs(funcMap).ParseFiles(filepath.Join(templateDir, "book.html")))

	pageTypes := make(map[string]Pagetype)
	pageTypes["3_rows_4_cols_page"] = *NewPagetype("page_3_rows_4_cols", 12, templateDir)
	pageTypes["4_rows_4_cols_page"] = *NewPagetype("page_4_rows_4_cols", 20, templateDir)
	pageTypes["single_with_cover"] = *NewPagetype("page_single_with_cover", 1, templateDir)

	pages := []*Page{}

	var currentPage *Page
	i := 0
	for _, barcodeAction := range barcodeActions {
		fmt.Printf("i=%d\n", i)

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

func (book Book) GetPages() []*Page {
	// Reorder Pages for 2 side print
	// [1,2] [3,4] -> [1,3] [2,4]

	var pages []*Page

	for i := 0; i < len(book.pages); i = i + 4 {
		pages = append(pages, book.pages[i+0])

		if (i + 2) < len(book.pages) {
			pages = append(pages, book.pages[i+2])
		} else {
			pages = append(pages, &Page{})
		}

		if (i + 1) < len(book.pages) {
			pages = append(pages, book.pages[i+1])
		} else {
			pages = append(pages, &Page{})
		}

		if (i + 3) < len(book.pages) {
			pages = append(pages, book.pages[i+3])
		} else {
			pages = append(pages, &Page{})
		}
	}

	return pages
}

func (book Book) GetPageCount() int {
	return len(book.pages)
}
