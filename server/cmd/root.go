package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"whattofarm/server/api"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   ".",
	Short: "This main command to start server",
	Long:  "This command start server at port localhost:3001",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func startServer() {
	log.Println("Listening at port 3001...")
	server := &http.Server{
		Addr:    "localhost:3001",
		Handler: api.NewRouter(),
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	gracefulShutdown(server, sig)
}

func gracefulShutdown(server *http.Server, sig chan os.Signal) {
	<-sig
	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), time.Second*20)
	defer shutdownCtxCancel()
	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			log.Fatal("graceful shutdown timed out and forcing exit.")
		}
	}()
	err := server.Shutdown(context.Background())
	if err != nil {
		log.Fatal("server shutdown error")
	}
	log.Println("Server shutdown now.")
}
