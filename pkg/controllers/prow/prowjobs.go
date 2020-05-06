package prow

import (
	"log"

	"k8s.io/client-go/tools/cache"
	prowinformers "k8s.io/test-infra/prow/client/informers/externalversions"
)

func NewProwjobsController(prowInformerFactory prowinformers.SharedInformerFactory) {
	prowInformerFactory.Prow().V1().ProwJobs().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    prowjobCreated,
		UpdateFunc: prowjobUpdated,
		DeleteFunc: prowjobDeleted,
	})
}

func prowjobCreated(obj interface{}) {
	log.Print("prowjobCreated")
}

func prowjobUpdated(oldObj, newObj interface{}) {
	log.Print("prowjobUpdated")
}

func prowjobDeleted(obj interface{}) {
	log.Print("prowjobDeleted")
}
