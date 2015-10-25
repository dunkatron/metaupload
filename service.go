package main

import (
	"fmt"
	"github.com/dunkatron/metaupload/providers"
	"github.com/dunkatron/metaupload/providers/imgur"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"io/ioutil"
	"mime/multipart"
)

func renderResult(data interface{}, err error, r render.Render) {
	if err != nil {
		r.JSON(200, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
			"data":    nil,
		})
		return
	} else {
		r.JSON(200, map[string]interface{}{
			"success": true,
			"error":   nil,
			"data":    data,
		})
	}
}

var genericError = fmt.Errorf("Failed to upload image.")

func registerProviderFactory(m *martini.ClassicMartini, provider providers.Provider) {

	// Sideload an image from a remote URL
	type imageSideload struct {
		Url string `form:"url" binding:"required"`
	}

	m.Post("/providers/"+provider.Name()+"/sideload", binding.Bind(imageSideload{}), func(sideload imageSideload, r render.Render) {
		url, err := provider.SideLoadImage(sideload.Url)
		renderResult(url, err, r)
	})

	// Upload an image from post data
	type imageUpload struct {
		Image *multipart.FileHeader `form:"image" binding:"required"`
	}

	m.Post("/providers/"+provider.Name()+"/upload", binding.Bind(imageUpload{}), func(upload imageUpload, r render.Render) {
		// Open the file uploaded to us
		file, err := upload.Image.Open()
		if err != nil {
			renderResult(nil, genericError, r)
			return
		}
		defer file.Close()

		contents, err := ioutil.ReadAll(file)
		if err != nil {
			renderResult(nil, genericError, r)
			return
		}

		url, err := provider.UploadImage(upload.Image.Filename, contents)
		renderResult(url, err, r)
	})
}

func main() {
	// Load service configs
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	// Configuration: supported providers
	var providerList = []providers.Provider{
		imgur.Provider(config.Imgur.ClientID),
	}

	m := martini.Classic()
	m.Use(render.Renderer())

	// Go through providers, build directory, and register each one
	// at its appropriate path.
	var providerNames = make([]string, len(providerList))
	for i, providerFactory := range providerList {
		providerNames[i] = providerFactory.Name()
		registerProviderFactory(m, providerFactory)
	}

	// Listing of providers
	m.Get("/providers", func(r render.Render) {
		renderResult(providerNames, nil, r)
	})

	m.Run()
}
