package main

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/toorop/gin-logrus"
	"net/http"
	"strconv"
	"time"
)

var rootCmd = &cobra.Command{
	Use: "ensicoin-explorer",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("ensicoin explorer version 0.0.0")

		dbPath, err := cmd.Flags().GetString("dbpath")
		if err != nil {
			log.WithError(err).Fatal("fatal error reading the database path")
		}

		rpcServerAddress, err := cmd.Flags().GetString("rpcserver")
		if err != nil {
			log.WithError(err).Fatal("fatal error reading the rpc server address")
		}

		storage := NewStorage(dbPath)
		if err := storage.Open(); err != nil {
			log.WithError(err).Fatal("fatal error opening the database")
		}
		defer storage.Close()

		synchronizer := NewSynchronizer(storage, rpcServerAddress)
		if err := synchronizer.Start(); err != nil {
			log.WithError(err).Fatal("fatal error starting the synchronizer")
		}

		r := gin.Default()
		r.Use(ginlogrus.Logger(log.StandardLogger()), gin.Recovery())

		r.GET("/blocks", func(c *gin.Context) {
			rawPage := c.DefaultQuery("page", "0")
			rawLimit := c.DefaultQuery("limit", "10")

			page, err := strconv.Atoi(rawPage)
			if err != nil {
				c.Status(http.StatusBadRequest)
				return
			}

			limit, err := strconv.Atoi(rawLimit)
			if err != nil {
				c.Status(http.StatusBadRequest)
				return
			}

			_ = page
			_ = limit

			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		srv := &http.Server{
			Addr:    ":8080",
			Handler: r,
		}

		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.WithError(err).Fatalf("fatal error listening")
			}
		}()

		interruptChannel := newInterruptListener()
		<-interruptChannel

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.WithError(err).Fatal("fatal error during server shutdown")
		}

		log.Info("Good bye.")
	},
}

func init() {
	rootCmd.Flags().String("dbpath", "database/data.db", "database path")
	rootCmd.Flags().String("rpcserver", "localhost:4225", "RPC server to connect to")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("fatal error")
	}
}
