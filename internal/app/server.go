package app

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router *gin.Engine
}

func NewServer() *Server {
	return &Server{
		Router: Router(),
	}
}

func (s *Server) Run() error {
	return s.Router.Run()
}
