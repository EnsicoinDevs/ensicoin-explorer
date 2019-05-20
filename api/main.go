package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	log.Info("ensicoin explorer version 0.0.0")

	s := NewStorage()
	if err := s.Open(); err != nil {
		log.WithError(err).Fatal("fatal error opening the database")
	}
	defer s.Close()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run()
}
