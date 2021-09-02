package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"

	"k8s.io/klog"
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"

	"github.com/revoman/edgetpu-device-plugin/pkg/plugin"
)

func main() {
	klog.InitFlags(flag.CommandLine)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		klog.Exitf("Could not create file system watcher: %v", err)
	}

	err = watcher.Add(pluginapi.DevicePluginPath)
	if err != nil {
		klog.Exitf("Could not watch device plugin path: %v", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	var svc plugin.Service
	restart := true
L:
	for {
		if restart {
			if svc != nil {
				svc.Stop()
			}
			svc = plugin.NewService()
			if err := svc.Serve(); err != nil {
				klog.Info("Could not contact Kubelet, retrying.  Did you enable the device plugin feature gate?")
			} else {
				restart = false
			}
		}

		select {
		case ev := <-watcher.Events:
			if ev.Name == pluginapi.KubeletSocket &&
				ev.Op&fsnotify.Create == fsnotify.Create {
				klog.Info("Kubelet socket created, restarting.")
				restart = true
			}

		case err := <-watcher.Errors:
			klog.Errorf("Received an error from file system watcher: %v", err)

		case sig := <-sigCh:
			switch sig {
			case syscall.SIGHUP:
				klog.Info("Received SIGHUP, restarting.")
				restart = true
			default:
				klog.Infof("Received signal \"%v\", shutting down.", sig)
				break L
			}
		}
	}
}
