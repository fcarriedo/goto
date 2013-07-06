package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var port = flag.Int("p", 80, "http port to run")

var db = map[string]string{
	"gogl": "http://google.com",
}

func init() {
	fmt.Println("Initializing the mapping database...")
}

// The core intent of the tool
func redirectHandler(w http.ResponseWriter, req *http.Request) {
	// Get it from the datastore and redirect
	url := db[mux.Vars(req)["entry"]]
	if url != "" {
		http.Redirect(w, req, url, http.StatusFound)
	} else {
		handleErr(w, req, http.StatusNotFound)
	}
}

func entriesHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		for k, v := range db {
			fmt.Fprintf(w, "%s=%s\r\n", k, v)
		}
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
			db[entry] = url
		}
	case "DELETE":
		delete(db, entry)
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

func homeHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "You're home")
}

func main() {
	flag.Parse()

	// Create the mux router
	router := mux.NewRouter()

	// Static resources
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets/"))))
	router.PathPrefix("/favicon.ico").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {})

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
