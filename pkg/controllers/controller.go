/*
Copyright 2019-2020 The Tekton Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
		http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"time"

	"github.com/eddycharly/prowjobs-extension/pkg/controllers/prow"
	prowversioned "k8s.io/test-infra/prow/client/clientset/versioned"
	prowinformers "k8s.io/test-infra/prow/client/informers/externalversions"
)

func StartProwControllers(clientset prowversioned.Interface, resyncDur time.Duration, stopCh <-chan struct{}) {
	prowInformerFactory := prowinformers.NewSharedInformerFactory(clientset, resyncDur)
	// Add all prow controllers
	prow.NewProwjobsController(prowInformerFactory)
	// Started once all controllers have been registered
	prowInformerFactory.Start(stopCh)
}
