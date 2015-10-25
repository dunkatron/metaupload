package providers

import ()

type Provider interface {
	// Returns the name of this provider
	Name() string

	// Takes a URL of an image to be sideloaded and returns the new URL
	SideLoadImage(url string) (*string, error)

	// Takes the reader of an image to be uploaded and returns the new URL after it's uploaded
	UploadImage(filename string, imageData []byte) (*string, error)
}
