package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"

	"github.com/belljustin/stamp/internal/db"
)

const (
	documentPath = "/document"
)

type DocumentResource struct {
	dao *db.DocumentDAO
}

// Create adds a new stamp to the database and returns the newly created id
func (dr DocumentResource) Create(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	stampReq := db.Document{}
	err = json.Unmarshal(body, &stampReq)
	if err != nil {
		log.Fatal(err)
	}

	id, err := dr.dao.Insert(stampReq.Hash)
	if err != nil {
		fmt.Fprint(w, err)
	}
	fmt.Fprint(w, id.String()+"\n")
}

// Get a document from the database and returns the details
func (dr DocumentResource) Get(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	sId := params.ByName("id")
	if sId == "" {
		http.Error(w, "Bad Request", 400)
		log.Fatal("Must provide an id")
	}

	id, err := uuid.Parse(sId)
	if err != nil {
		http.Error(w, "Bad Request. Could not parse uuid.", 400)
		return
	}

	stamp, err := dr.dao.Get(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Not Found", 404)
			return
		}
		log.Fatalf("Error retrieving stamp %v: %v", id, err)
	}

	w.Header().Set("Content-Type", "application/json")
	jStamp, err := json.Marshal(stamp)
	if err != nil {
		log.Fatalf("Error serializing document %v: %v", stamp, err)
	}
	w.Write(jStamp)
}