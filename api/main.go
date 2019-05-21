package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
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

		r.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		r.Run()
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
