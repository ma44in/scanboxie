package web

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

func (webapp *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
		"add": func(x int, y int) int {
			return x + y
		},
	}

	//t, _ := template.New("index").Funcs(funcMap).ParseFS(getFileSystem("templates"), "index.html")
	t, _ := template.New("").Funcs(funcMap).ParseFS(getFileSystem("templates"), "index.html")
	//err := t.Execute(w, webapp)
	err := t.ExecuteTemplate(w, "index.html", webapp)
	if err != nil {
		panic(err)
	}
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (webapp *App) removeBarcodeActionHandler(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	barcode := keys.Get("barcode")

	err := webapp.BarcodeConfig.RemoveBarcodeAction(barcode)
	if err != nil {
		log.Printf("error removing barcode action: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (webapp *App) moveBarcodeActionHandler(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	barcode := keys.Get("barcode")
	newIndexString := keys.Get("newIndex")

	newIndex, err := strconv.Atoi(newIndexString)
	if err != nil {
		log.Printf("newIndex could not be converted to int: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = webapp.BarcodeConfig.MoveBarcodeAction(barcode, newIndex)
	if err != nil {
		log.Printf("error moving barcode action: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (webapp *App) saveBarcodeConfigHandler(w http.ResponseWriter, r *http.Request) {
	err := webapp.BarcodeConfig.Save()
	if err != nil {
		log.Printf("error saving barcodeconfig: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (webapp *App) getLastScannedBarcodeHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, webapp.scanboxie.GetLastScannedBarcode())
}

func (webapp *App) bookHandler(w http.ResponseWriter, r *http.Request) {
	err := webapp.ScanboxieBook.Write(w, webapp.BarcodeConfig.BarcodeActions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
