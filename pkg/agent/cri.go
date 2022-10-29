package agent

import (
  "fmt"
  "time"
  "context"
  "errors"

  "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
  runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
  "github.com/containerd/containerd/integration/remote/util"
)

// maxMsgSize use 16MB as the default message size limit.
// grpc library default is 4MB
const maxMsgSize = 1024 * 1024 * 16

type ImageService struct {
  client      runtimeapi.ImageServiceClient
  timeout     time.Duration
}

func NewImageService(endpoint string, connectTimeout time.Duration) (*ImageService, error) {
  addr, dialer, err := util.GetAddressAndDialer(endpoint)
  if err != nil {
    return nil, err
  }

  ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
  defer cancel()

  conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)),
	)

	if err != nil {
		return nil, err
	}

  return &ImageService{
    client:  runtimeapi.NewImageServiceClient(conn),
    timeout: connectTimeout,
  }, nil
}

func (r *ImageService) ListImages(filter *runtimeapi.ImageFilter, opts ...grpc.CallOption) ([]*runtimeapi.Image, error) {
  ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	resp, err := r.client.ListImages(ctx, &runtimeapi.ListImagesRequest{Filter: filter}, opts...)
	if err != nil {
		return nil, err
	}

	return resp.Images, nil
}


func (r *ImageService) ImageStatus(image *runtimeapi.ImageSpec, opts ...grpc.CallOption) (*runtimeapi.Image, error) {
  ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	resp, err := r.client.ImageStatus(ctx, &runtimeapi.ImageStatusRequest{
		Image: image,
	}, opts...)
	if err != nil {
		return nil, err
	}

	if resp.Image != nil {
		if resp.Image.Id == "" || resp.Image.Size() == 0 {
			errorMessage := fmt.Sprintf("Id or size of image %q is not set", image.Image)
			return nil, errors.New(errorMessage)
		}
	}

	return resp.Image, nil
}

// PullImage pulls an image with authentication config.
func (r *ImageService) PullImage(image *runtimeapi.ImageSpec, auth *runtimeapi.AuthConfig, podSandboxConfig *runtimeapi.PodSandboxConfig, opts ...grpc.CallOption) (string, error) {
  ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resp, err := r.client.PullImage(ctx, &runtimeapi.PullImageRequest{
		Image:         image,
		Auth:          auth,
		SandboxConfig: podSandboxConfig,
	}, opts...)
	if err != nil {
		return "", err
	}

	if resp.ImageRef == "" {
		errorMessage := fmt.Sprintf("imageRef of image %q is not set", image.Image)
		return "", errors.New(errorMessage)
	}

	return resp.ImageRef, nil
}

// RemoveImage removes the image.
func (r *ImageService) RemoveImage(image *runtimeapi.ImageSpec, opts ...grpc.CallOption) error {
  ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.client.RemoveImage(ctx, &runtimeapi.RemoveImageRequest{
		Image: image,
	}, opts...)
	if err != nil {
		return err
	}

	return nil
}

// ImageFsInfo returns information of the filesystem that is used to store images.
func (r *ImageService) ImageFsInfo(opts ...grpc.CallOption) ([]*runtimeapi.FilesystemUsage, error) {
	// Do not set timeout, because `ImageFsInfo` takes time.
	// TODO(random-liu): Should we assume runtime should cache the result, and set timeout here?
  ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resp, err := r.client.ImageFsInfo(ctx, &runtimeapi.ImageFsInfoRequest{}, opts...)
	if err != nil {
		return nil, err
	}
	return resp.GetImageFilesystems(), nil
}
