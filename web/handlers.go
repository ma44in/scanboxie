package web

import (
	"net/http"
	"scanboxie/pkg/scanboxieprint"
)

func (webapp *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	book := scanboxieprint.NewBook(webapp.TemplateDir, webapp.BarcodeConfig.BarcodeActions)

	err := book.Write(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
