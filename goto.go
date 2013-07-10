package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

var db = make(map[string]string)

// Command line args def
var port = flag.Int("p", 80, "http port to run")
var devMode = flag.Bool("dev", false, "run in development mode")

var templates *template.Template

var lazySave = debounce(save, 5*time.Second)

func init() {
	fmt.Println("Initializing the mapping database...")
	load()

	// Parse all templates
	templates = template.Must(template.New("app").ParseGlob("web/tmpl/*.html"))
}

// Home handler
func homeHandler(w http.ResponseWriter, req *http.Request) {
	if *devMode {
		// Compile templates on every request
		templates = template.Must(template.New("app").ParseGlob("web/tmpl/*.html"))
	}

	templates.ExecuteTemplate(w, "index.html", req.FormValue("m"))
}

// The core intent of the tool
func redirectHandler(w http.ResponseWriter, req *http.Request) {
	// Get it from the datastore and redirect
	entry := mux.Vars(req)["entry"]
	url := db[entry]
	if url != "" {
		http.Redirect(w, req, url, http.StatusFound)
	} else {
		templates.ExecuteTemplate(w, "404.html", entry)
	}
}

func entriesHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		templates.ExecuteTemplate(w, "listing.html", db)
	default:
		handleErr(w, req, http.StatusMethodNotAllowed)
	}
}

// Handles the mapping entry details.
func entryHandler(w http.ResponseWriter, req *http.Request) {
	entry := mux.Vars(req)["entry"]
	switch req.Method {
	case "GET":
		url := db[entry]
		if url == "" {
			handleErr(w, req, http.StatusNotFound)
		} else {
			fmt.Fprintf(w, "%s\r\n", url)
		}
	case "PUT":
		url := req.FormValue("url")
		if url == "" {
			handleErr(w, req, http.StatusBadRequest)
		} else {
			// TODO: db map needs lock syncronization (specially here)
			if db[entry] == "" {
				// Didn't exist before => 201 (Created)
				// See: http://www.w3.org/Protocols/rfc2616/rfc2616-sec9.html#sec9.6
				w.WriteHeader(http.StatusCreated)
			}
			db[entry] = url
			lazySave()
		}
	case "DELETE":
		delete(db, entry)
		lazySave()
	default:
		handleErr(w, req, http.StatusMethodNotAllowed)
	}
}

// Formats the given status in a standard. Any status would be managed in
// the same way whether is an error or not (does not enforce)
func handleErr(w http.ResponseWriter, req *http.Request, errStatus int) {
	err := fmt.Sprintf("%d %s", errStatus, http.StatusText(errStatus))
	http.Error(w, err, errStatus)
}

func main() {
	flag.Parse()

	// Create the mux router
	router := mux.NewRouter()

	// Static resources - resource ending in common know web file formats
	// (css, html, jpg, etc.) get handled directly by the fileServer
	router.Handle("/{static-res:(.+\\.)(js|css|jpg|png|ico|gif)$}", http.FileServer(http.Dir("web/")))

	// Home handler
	router.HandleFunc("/", homeHandler)

	// Mapping handling
	router.HandleFunc("/mappings", entriesHandler)
	router.HandleFunc("/mappings/{entry}", entryHandler)

	// Redirect (core function)
	router.HandleFunc("/{entry}", redirectHandler)

	// Hook it with http pkg
	http.Handle("/", router)

	host := fmt.Sprintf(":%d", *port)
	fmt.Printf("Server up and listening on %s\n", host)
	log.Fatal(http.ListenAndServe(host, nil))
}

const dbfile = "db/mappings.csv"

func load() {
	// Open the db file
	f, err := os.Open(dbfile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Read the containing records.
	reader := csv.NewReader(f)
	rec, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Load them into the map
	for _, row := range rec {
		db[row[0]] = row[1]
	}
}

// Requests a save. Debounces multiple calls
// Very quick and dirty. Revise with channels
// since update and save are concurrent matters.
func save() {
	f, err := os.Create(dbfile)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()
	fmt.Println("Persisting mappings.")
	for k, v := range db {
		err = writer.Write([]string{k, v})
		if err != nil {
			log.Printf("Error saving mapping: %s => %s", k, v)
		}
	}
}

// Utility functions

func debounce(fn func(), d time.Duration) func() {
	active := false
	return func() {
		if !active {
			active = true
			time.AfterFunc(d, func() {
				active = false
				fn() // execute
			})
		}
	}
}
