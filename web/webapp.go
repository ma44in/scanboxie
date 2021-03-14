package web

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"scanboxie/pkg/scanboxie"
	"text/template"

	"github.com/gorilla/mux"
)

// App Implement a singleton Pattern
type App struct {
	Router        *mux.Router
	StaticDir     string
	TemplateDir   string
	Templates     map[string]*template.Template
	BarcodeConfig *scanboxie.BarcodeConfig
}

func parseTemplates(templateDir string, filenames ...string) (templates map[string]*template.Template) {
	if len(filenames) == 0 {
		log.Fatal("You must pass at least one file to parseTemplates")
	}
	templates = make(map[string]*template.Template)

	for _, filename := range filenames {
		t := template.New("")
		t, _ = t.ParseFiles(filepath.Join(templateDir, filename))
		//t.Funcs()
		templates[filename] = t

		// templates[filename] = template.Must(template.ParseFiles(
		// 	filepath.Join(templateDir, filename),
		// )).Funcs(template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }})
	}
	return templates
}

// NewApp returns the app
func NewApp(barcodeDirMapFilepath string) *App {
	barcodeConfig, err := scanboxie.NewBarcodeConfig(barcodeDirMapFilepath, false)
	if err != nil {
		panic("could not read barcode config")
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	staticDir := filepath.Join(dir, "web", "static")
	templateDir := filepath.Join(dir, "web", "templates")

	webapp := &App{
		Router:        mux.NewRouter(),
		StaticDir:     staticDir,
		TemplateDir:   templateDir,
		BarcodeConfig: barcodeConfig,
	}

	webapp.Router.HandleFunc("/", webapp.indexHandler).Name("index")

	return webapp
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.Router.ServeHTTP(w, r)
}

// Serve serve the web app
func (app *App) Serve(address string) {
	http.Handle("/", app)

	// Serving static file
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir(app.StaticDir))))

	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("An error occured when trying to start server: \n", err)
	}
}
