package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func lastHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	n, _ := strconv.ParseInt(vars["n"], 10, 64)
	method := vars["method"]
	//method := strings.ToUpper(vars["method"])

	db, err := dbConnect()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	last, err := dbSelectLastLogs(method, n, db)
	if err != nil || last == nil {
		log.Printf("> last %+v", last)
		log.Printf("  error %v", err)
	}

	js, err := json.Marshal(last)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func startService() {
	var port = ":" + os.Getenv("OOWPORT")

	r := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET"})
	origins := handlers.AllowedOrigins([]string{"*"})

	// two subrouters (strictly uppercase)
	r.HandleFunc("/last/{n:[0-9]+}/{method:GET|POST}", lastHandler).Methods("GET")

	log.Println("Listening on " + port)
	log.Fatal(http.ListenAndServe(port, handlers.CORS(headers, methods, origins)(r)))
}
