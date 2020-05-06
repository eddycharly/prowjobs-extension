package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	restful "github.com/emicklei/go-restful"
	"knative.dev/pkg/signals"
)

func main() {
	ctx := signals.NewContext()

	log.Print("Creating server")
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})

	RegisterWeb(wsContainer)

	server := &http.Server{Addr: ":8080", Handler: wsContainer}

	errCh := make(chan error, 1)
	defer close(errCh)
	go func() {
		// Don't forward ErrServerClosed as that indicates we're already shutting down.
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("server failed: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		log.Fatal(err)
	case <-ctx.Done():
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}
}

// RegisterWeb registers extension web bundle on the container
func RegisterWeb(container *restful.Container) {
	var handler http.Handler
	webResourcesDir := os.Getenv("WEB_RESOURCES_DIR")
	koDataPath := os.Getenv("KO_DATA_PATH")
	_, err := os.Stat(webResourcesDir)
	if err != nil {
		if os.IsNotExist(err) {
			if koDataPath != "" {
				log.Printf("WEB_RESOURCES_DIR %s not found, serving static content from KO_DATA_PATH instead.", webResourcesDir)
				handler = http.FileServer(http.Dir(koDataPath))
			} else {
				log.Printf("WEB_RESOURCES_DIR %s not found and KO_DATA_PATH not found, static resource (UI) problems to be expected.", webResourcesDir)
			}
		} else {
			log.Printf("error returned while checking for WEB_RESOURCES_DIR %s", webResourcesDir)
		}
	} else {
		log.Printf("Serving static files from WEB_RESOURCES_DIR: %s", webResourcesDir)
		handler = http.FileServer(http.Dir(webResourcesDir))
	}
	container.Handle("/web/", http.StripPrefix("/web/", handler))
}
