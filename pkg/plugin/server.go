package plugin

import (
	"context"
	"io/ioutil"
	"reflect"
	"strings"
	"time"

	"k8s.io/klog"
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"

	"github.com/kkohtaka/edgetpu-device-plugin/pkg/fileutil"
)

const (
	systemPath = "/sys/bus/usb"
	devicePath = "/dev/bus/usb"

	pkgSrcDir     = "/etc/edgetpu/"
	pkgInstallDir = "/opt/edgetpu/"
)

var (
	vids = []string{"1a6e", "18d1"}
)

type DevicePluginServer struct {
	devices   []*pluginapi.Device
	devicesCh chan []*pluginapi.Device
	stopCh    chan struct{}
}

// NewDevicePluginServer creates a new DevicePluginServer of Edge TPU.
func NewDevicePluginServer() *DevicePluginServer {
	return &DevicePluginServer{
		devicesCh: make(chan []*pluginapi.Device),
		stopCh:    make(chan struct{}),
	}
}

// GetDevicePluginOptions implements a part of
// (k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1).DevicePluginServer.
func (dps *DevicePluginServer) GetDevicePluginOptions(
	ctx context.Context,
	empty *pluginapi.Empty,
) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{}, nil
}

// ListAndWatch implements a part of (k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1).DevicePluginServer.
func (dps *DevicePluginServer) ListAndWatch(
	empty *pluginapi.Empty,
	server pluginapi.DevicePlugin_ListAndWatchServer,
) error {
	klog.Info("Start watching devices")
	for {
		select {
		case <-dps.stopCh:
			klog.Info("Exit watching devices")
			return nil
		case devices := <-dps.devicesCh:
			server.Send(&pluginapi.ListAndWatchResponse{
				Devices: devices,
			})
			klog.Info("Update a device list")
		}
	}
}

// Allocate implements a part of (k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1).DevicePluginServer.
func (dps *DevicePluginServer) Allocate(
	ctx context.Context,
	reqs *pluginapi.AllocateRequest,
) (*pluginapi.AllocateResponse, error) {
	var resp pluginapi.AllocateResponse
	for _, req := range reqs.GetContainerRequests() {
		klog.Infof("Allocating devices... Device IDs: %v", req.DevicesIDs)
		resp.ContainerResponses = append(
			resp.ContainerResponses,
			&pluginapi.ContainerAllocateResponse{
				Devices: []*pluginapi.DeviceSpec{
					&pluginapi.DeviceSpec{
						ContainerPath: devicePath,
						HostPath:      devicePath,
						Permissions:   "rw",
					},
				},
				Mounts: []*pluginapi.Mount{
					&pluginapi.Mount{
						ContainerPath: systemPath,
						HostPath:      systemPath,
						ReadOnly:      true,
					},
				},
			},
		)
	}
	return &resp, nil
}

// PreStartRequired implements a part of (k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1).DevicePluginServer.
func (dps *DevicePluginServer) PreStartContainer(
	ctx context.Context,
	req *pluginapi.PreStartContainerRequest,
) (*pluginapi.PreStartContainerResponse, error) {
	return &pluginapi.PreStartContainerResponse{}, nil
}

func (dps *DevicePluginServer) stopMonitoringDevices() {
	close(dps.stopCh)
}

func (dps *DevicePluginServer) updateCurrentDevices(devices []*pluginapi.Device) {
	dps.devices = devices
	dps.devicesCh <- devices
}

func (dps *DevicePluginServer) startMonitoringDevices() {
	dps.updateCurrentDevices(nil)

	ticker := time.NewTicker(5 * time.Second)
	for {
		<-ticker.C

		func() {
			var devices []*pluginapi.Device
			files, err := fileutil.FindFiles("/sys/devices", "idVendor")
			if err != nil {
				klog.Errorf("Could not find idVendor files in /sys/devices: %v", err)
			}
			for _, file := range files {
				data, err := ioutil.ReadFile(file)
				if err != nil {
					klog.Errorf("Could not read file: %v", err)
					continue
				}
				vid := strings.TrimSpace(string(data))
				if vid == vids[0] || vid == vids[1] {
					devices = append(devices, &pluginapi.Device{
						ID:     "42",
						Health: pluginapi.Healthy,
					})
					break
				}
			}

			if !reflect.DeepEqual(dps.devices, devices) {
				if len(devices) > 0 {
					klog.Info("Edge TPU became active.")
				} else {
					klog.Info("Edge TPU became inactive.")
				}
				dps.updateCurrentDevices(devices)
			}
		}()
	}
}
