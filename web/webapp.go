package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"scanboxie/pkg/scanboxie"
	"scanboxie/pkg/scanboxieprint"

	"github.com/gorilla/mux"
)

//go:embed static templates
var filesystemContent embed.FS

// App Implement a singleton Pattern
type App struct {
	router        *mux.Router
	imageDir      string
	scanboxie     *scanboxie.Scanboxie
	ScanboxieBook *scanboxieprint.Book
	BarcodeConfig *scanboxie.BarcodeConfig
	CommandSets   *scanboxie.CommandSets
}

// NewApp returns the app
func NewApp(scanboxie *scanboxie.Scanboxie, imageDir string) *App {
	scanboxieBook := scanboxieprint.NewBook(getFileSystem("templates/books"))

	webapp := &App{
		router:        mux.NewRouter(),
		imageDir:      imageDir,
		scanboxie:     scanboxie,
		ScanboxieBook: scanboxieBook,
		BarcodeConfig: scanboxie.BarcodeConfig,
		CommandSets:   scanboxie.CommandSets,
	}

	webapp.router.HandleFunc("/", webapp.indexHandler).Name("index")
	webapp.router.HandleFunc("/book", webapp.bookHandler).Name("book")
	webapp.router.HandleFunc("/addBarcodeAction", webapp.addBarcodeActionHandler).Name("addBarcodeAction")
	webapp.router.HandleFunc("/removeBarcodeAction", webapp.removeBarcodeActionHandler).Name("removeBarcodeAction")
	webapp.router.HandleFunc("/saveBarcodeConfig", webapp.saveBarcodeConfigHandler).Name("saveBarcodeConfig")
	webapp.router.HandleFunc("/getLastScannedBarcode", webapp.getLastScannedBarcodeHandler).Name("getLastScannedBarcode")
	webapp.router.HandleFunc("/moveBarcodeAction", webapp.moveBarcodeActionHandler).Name("moveBarcodeAction")

	return webapp
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.router.ServeHTTP(w, r)
}

// Serve serve the web app
func (app *App) Serve(address string) {
	http.Handle("/", app)

	// Serving static file
	http.Handle("/static/", http.FileServer(http.FS(getFileSystem("."))))

	if app.imageDir != "" {
		log.Printf("Handle /images/ for %s\n", app.imageDir)
		http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(app.imageDir))))
	}

	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("An error occured when trying to start server: \n", err)
	}
}

func getFileSystem(dir string) fs.FS {
	fsys, err := fs.Sub(filesystemContent, dir)
	if err != nil {
		log.Fatal(err)
	}
	return fsys
}
