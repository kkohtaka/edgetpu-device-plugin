package plugin

import (
	"context"
	"reflect"
	"time"

	"github.com/google/gousb"

	"k8s.io/klog"
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"
)

const (
	systemPath = "/sys/bus/usb"
	devicePath = "/dev/bus/usb"

	pkgSrcDir     = "/etc/edgetpu/"
	pkgInstallDir = "/opt/edgetpu/"

	checkFile = ".check"
)

var (
	vids = []gousb.ID{0x1a6e, 0x18d1}
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
	ctx := gousb.NewContext()
	defer ctx.Close()

	dps.updateCurrentDevices(nil)

	ticker := time.NewTicker(5 * time.Second)
	for {
		<-ticker.C

		func() {
			var devices []*pluginapi.Device
			if _, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
				if desc.Vendor == vids[0] || desc.Vendor == vids[1] {
					devices = append(devices, &pluginapi.Device{
						ID:     "42",
						Health: pluginapi.Healthy,
					})
					return false
				}
				return false
			}); err != nil {
				klog.Errorf("Could not find a device: %v", err)
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
