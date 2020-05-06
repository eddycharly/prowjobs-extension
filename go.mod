module github.com/eddycharly/prowjobs-extension

go 1.14

replace (
	k8s.io/api => k8s.io/api v0.17.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.3
	k8s.io/client-go => k8s.io/client-go v0.17.3
)

require (
	github.com/emicklei/go-restful v2.9.5+incompatible
	github.com/tektoncd/dashboard v0.6.1
	google.golang.org/appengine v1.6.5
	k8s.io/client-go v9.0.0+incompatible
	k8s.io/test-infra v0.0.0-20200502165826-8fe3dec9cbf7
	knative.dev/pkg v0.0.0-20200207155214-fef852970f43
)
