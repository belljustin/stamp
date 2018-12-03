package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"

	"github.com/belljustin/stamp/internal/db"
)

const (
	stampPath = "/stamp"
)

// GetStamp gets a stamp from the database and returns the details
func GetStamp(w http.ResponseWriter, req *http.Request, params httprouter.Params, dao *db.StampDAO) {
	sId := params.ByName("id")
	if sId == "" {
		log.Fatal("Must provide an id")
	}

	id, err := uuid.Parse(sId)
	if err != nil {
		http.Error(w, "Bad Request. Could not parse uuid.", 400)
	}

	stamp, err := dao.Get(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Not Found", 404)
			return
		}
		log.Panicf("Error retrieving stamp %v: %v", id, err)
	}

	w.Header().Set("Content-Type", "application/json")
	jStamp, err := json.Marshal(stamp)
	if err != nil {
		http.Error(w, "Unprocessable Entity", 422)
		return
	}
	w.Write(jStamp)
}

// NewGetStampHandler builds a new GetStamp Handler
func NewGetStampHandler(dao *db.StampDAO) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, pr httprouter.Params) {
		GetStamp(w, req, pr, dao)
	}
}
