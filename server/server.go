package server

import (
	"SafeApeBot/db"
	"SafeApeBot/structs"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

var DB *db.DB = nil
var Config *structs.SConfig = nil

type (
	Server struct {
		Router *httprouter.Router
		// DB     *db.DB
	}
)

func (s *Server) StartServer() {
	port := strconv.Itoa(Config.ServerPort)
	log.Println("Starting server on port " + port)
	router := httprouter.New()
	s.Router = router
	s.Addpaths()
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func (s *Server) DumpCache() {
	// cache.Dump()
}
