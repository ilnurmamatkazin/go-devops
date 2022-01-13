package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/handlers"
)

const (
	url  = "127.0.0.1"
	port = 8080
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	h := handlers.New()
	r := h.NewRouter()

	go http.ListenAndServe(fmt.Sprintf(":%d", port), r)

	// fmt.Println("Server started...")

	<-quit

	// fmt.Println("Server shutdown")
}
