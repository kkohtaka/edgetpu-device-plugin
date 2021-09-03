module github.com/therevoman/edgetpu-device-plugin

go 1.16

require (
	github.com/fsnotify/fsnotify v1.4.9
	github.com/golang/protobuf v1.4.2 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/grpc v1.27.0
	k8s.io/klog v1.0.0
	k8s.io/kubelet v0.18.19
)

replace (
	k8s.io/api => k8s.io/api v0.18.19
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.19
	k8s.io/client-go => k8s.io/client-go v0.18.19
	k8s.io/component-base => k8s.io/component-base v0.18.19
	k8s.io/kubernetes => k8s.io/kubernetes v1.18.19
)
