package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/eddycharly/prowjobs-extension/pkg/controllers"
	restful "github.com/emicklei/go-restful"
	"github.com/tektoncd/dashboard/pkg/logging"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	prowversioned "k8s.io/test-infra/prow/client/clientset/versioned"
	"knative.dev/pkg/signals"
)

func main() {
	// kubeConfigPath := "/Users/charlesbreteche/.kube/k8s.ci.agrico.tech"
	kubeConfigPath := ""

	var cfg *rest.Config
	var err error

	log.Print("Creating clients")
	if kubeConfigPath != "" {
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			log.Printf("Error building kubeconfig from %s: %s", kubeConfigPath, err.Error())
		}
	} else {
		if cfg, err = rest.InClusterConfig(); err != nil {
			log.Printf("Error building kubeconfig: %s", err.Error())
		}
	}

	prowjobsClient, err := prowversioned.NewForConfig(cfg)
	if err != nil {
		log.Printf("Error building prow clientset: %s", err.Error())
	}

	ctx := signals.NewContext()

	log.Print("Creating controllers")
	resyncDur := time.Second * 30
	controllers.StartProwControllers(prowjobsClient, resyncDur, ctx.Done())

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
		// if err := server.Shutdown(context.Background()); err != nil {
		// 	logging.Log.Fatal(err)
		// }
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
				logging.Log.Warnf("WEB_RESOURCES_DIR %s not found, serving static content from KO_DATA_PATH instead.", webResourcesDir)
				handler = http.FileServer(http.Dir(koDataPath))
			} else {
				logging.Log.Errorf("WEB_RESOURCES_DIR %s not found and KO_DATA_PATH not found, static resource (UI) problems to be expected.", webResourcesDir)
			}
		} else {
			logging.Log.Errorf("error returned while checking for WEB_RESOURCES_DIR %s", webResourcesDir)
		}
	} else {
		logging.Log.Infof("Serving static files from WEB_RESOURCES_DIR: %s", webResourcesDir)
		handler = http.FileServer(http.Dir(webResourcesDir))
	}
	container.Handle("/web/", http.StripPrefix("/web/", handler))
}
