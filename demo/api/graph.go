package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ijsnow/dgql/dgql"
	"github.com/ijsnow/dgql/dgql/client"
	"github.com/ijsnow/dgql/dgql/schema"
)

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	s, err := dgql.NewSchema(ctx, client.ClientOptions{"play.dgraph.io"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var args schema.ExecutionArgs

	err = json.NewDecoder(r.Body).Decode(&args)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := schema.Execute(ctx, s, args)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func main() {
	http.HandlerFunc(handler)
	log.Fatal(http.ListenAndServe(":3333", nil))
}
