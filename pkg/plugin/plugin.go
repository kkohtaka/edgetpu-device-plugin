package plugin

import (
	"context"
	"net"
	"os"
	"path"
	"time"

	errors "golang.org/x/xerrors"

	"google.golang.org/grpc"

	"k8s.io/klog"
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"
)

const (
	resourceName = "kkohtaka.org/edgetpu"
	pluginSock   = pluginapi.DevicePluginPath + "edgetpu.sock"

	defaultDialTimeout = 5 * time.Second
)

// Service represents an interface of Edge TPU Device Plugin service.
type Service interface {
	// Serve starts the service of the Device Plugin.
	Serve() error

	// Stop stops the service of the Device Plugin.
	Stop()
}

// NewService creates a new Service instance.
func NewService() Service {
	return &service{}
}

type service struct {
	grpcServer *grpc.Server
	dpServer   *DevicePluginServer
}

// Serve implements Service.Serve.
func (s *service) Serve() error {
	if err := s.start(); err != nil {
		return errors.Errorf("start gRPC server: %w", err)
	}
	klog.Info("gRPC server started.")

	if err := s.register(); err != nil {
		s.Stop()
		return errors.Errorf("register device: %w", err)
	}
	return nil
}

// Stop implements Service.Stop.
func (s *service) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.Stop()
		s.grpcServer = nil
	}

	if s.dpServer != nil {
		s.dpServer.stopMonitoringDevices()
		s.dpServer = nil
	}
}

func (s *service) start() error {
	if err := os.Remove(pluginSock); err != nil && !os.IsNotExist(err) {
		return errors.Errorf("cleanup device plugin socket: %w", err)
	}

	sock, err := net.Listen("unix", pluginSock)
	if err != nil {
		return errors.Errorf("listen device plugin socket: %w", err)
	}
	s.grpcServer = grpc.NewServer([]grpc.ServerOption{}...)
	s.dpServer = NewDevicePluginServer()
	pluginapi.RegisterDevicePluginServer(s.grpcServer, s.dpServer)

	go s.grpcServer.Serve(sock)

	c, err := grpc.Dial(pluginSock,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(defaultDialTimeout),
		grpc.WithDialer(func(address string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", address, timeout)
		}),
	)
	if err != nil {
		return errors.Errorf("dial to device plugin socket: %w", err)
	}
	c.Close()
	klog.Info("Started gRPC service on plugin socket")

	go s.dpServer.startMonitoringDevices()
	klog.Info("Started monitoring devices")

	return nil
}

func (s *service) register() error {
	conn, err := grpc.Dial(pluginapi.KubeletSocket,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(defaultDialTimeout),
		grpc.WithDialer(func(address string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", address, timeout)
		}),
	)
	if err != nil {
		return errors.Errorf("dial to kubelet socket: %w", err)
	}
	klog.Info("Opened connection to kubelet socket")
	defer conn.Close()

	client := pluginapi.NewRegistrationClient(conn)
	if _, err = client.Register(
		context.Background(),
		&pluginapi.RegisterRequest{
			Version:      pluginapi.Version,
			Endpoint:     path.Base(pluginSock),
			ResourceName: resourceName,
		},
	); err != nil {
		return errors.Errorf("register device plugin: %w", err)
	}
	klog.Info("Registered device plugin")
	return nil
}
