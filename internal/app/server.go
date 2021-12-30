package app

import (
	"github.com/abhinavmsra/go-api/internal/repository"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Storage repository.Storage
	router  *gin.Engine
}

func NewServer(storage repository.Storage) *Server {
	return &Server{
		Storage: storage,
	}
}

func (s *Server) Run() error {
	r := s.Routes()
	return r.Run()
}
