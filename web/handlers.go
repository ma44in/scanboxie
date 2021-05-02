package web

import (
	"log"
	"net/http"
	"scanboxie/pkg/scanboxieprint"
	"text/template"
)

func (webapp *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFS(getFileSystem("templates"), "index.html")
	t.Execute(w, webapp)
}

func (webapp *App) addBarcodeActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	r.ParseForm()

	barcode := r.FormValue("barcode")
	commandSetKey := r.FormValue("commandset")
	value := r.FormValue("value")
	pagetype := r.FormValue("pagetype")

	err := webapp.BarcodeConfig.AddBarcodeAction(barcode, commandSetKey, value, pagetype)
	if err != nil {
		log.Printf("error adding barcode action: %v\n", err)
	}

}

func (webapp *App) bookHandler(w http.ResponseWriter, r *http.Request) {
	book := scanboxieprint.NewBook(getFileSystem("templates/books"), webapp.BarcodeConfig.BarcodeActions)

	err := book.Write(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
