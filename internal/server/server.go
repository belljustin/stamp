package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"

	"github.com/belljustin/stamp/internal/configs"
	"github.com/belljustin/stamp/internal/db"
)

const (
	pingPath = "/ping"
)

type Server struct {
	router *httprouter.Router
}

func Ping(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Pong\n")
}

func NewServer(config *configs.Config, documentDAO *db.DocumentDAO, stampDAO *db.StampDAO) *Server {
	router := httprouter.New()

	router.GET(pingPath, Ping)

	dr := DocumentResource{documentDAO}
	router.POST(documentPath, dr.Create)
	router.GET(documentPath+"/:id", dr.Get)

	router.GET("/stamp/:id", NewGetStampHandler(stampDAO))

	return &Server{router}
}

func (s *Server) Start() {
	log.Fatal(http.ListenAndServe(":9000", s.router))
}
