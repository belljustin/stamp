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

type DocumentResponse struct {
	Id    uuid.UUID `json:"id"`
	Hash  string    `json:"hash"`
	State string    `json:"state"`
	Stamp string    `json:"stamp"`
}

func buildDocumentResponse(doc *db.Document, host string) *DocumentResponse {
	res := &DocumentResponse{
		doc.Id,
		doc.Hash,
		doc.State,
		"",
	}

	if doc.StampId != uuid.Nil {
		res.Stamp = BuildStampLocation(host, doc.StampId)
	}
	return res
}

// Create adds a new document to the database and returns the newly created id
func (dr DocumentResource) Create(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	docReq := db.Document{}
	err = json.Unmarshal(body, &docReq)
	if err != nil {
		log.Fatal(err)
	}

	id, err := dr.dao.Insert(docReq.Hash)
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

	document, err := dr.dao.Get(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Not Found", 404)
			return
		}
		log.Fatalf("Error retrieving document%v: %v", id, err)
	}

	docRes := buildDocumentResponse(document, req.Host)
	w.Header().Set("Content-Type", "application/json")
	jDocument, err := json.Marshal(docRes)
	if err != nil {
		log.Fatalf("Error serializing document %v: %v", document, err)
	}
	w.Write(jDocument)
}
