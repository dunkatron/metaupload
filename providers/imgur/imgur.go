package imgur

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dunkatron/metaupload/providers"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

var imgurUrl = "https://api.imgur.com/3/image"

type imgurResponse struct {
	Data *struct {
		Error *string `json:"error"`
		Link  *string `json:"link"`
	} `json:"data"`
	Success *bool `json:"success"`
}

type provider struct {
	clientId string
}

var genericError = fmt.Errorf("Failed to upload image to Imgur.")

// Parse the returned response from Imgur and look for the image link to return.
func parseResponse(res *http.Response) (*string, error) {
	ret, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, genericError
	}

	responseObject := imgurResponse{}

	err = json.Unmarshal(ret, &responseObject)

	if err != nil {
		return nil, genericError
	}

	// Figure out what was successfully pulled out of the response JSON and use it to
	// fill out the return values
	if responseObject.Success != nil && *responseObject.Success && responseObject.Data.Link != nil {
		return responseObject.Data.Link, nil
	} else if responseObject.Data != nil && responseObject.Data.Error != nil {
		return nil, fmt.Errorf("Imgur error: %s", *responseObject.Data.Error)
	} else {
		return nil, fmt.Errorf("Illegal Imgur API response.")
	}
}

func (p provider) SideLoadImage(imageUrl string) (*string, error) {
	// Imgur API needs us to send the image URL in the "image" field of the post data

	values := url.Values{"image": {imageUrl}}
	body := strings.NewReader(values.Encode())
	req, err := http.NewRequest("POST", imgurUrl, body)

	req.Header.Add("Authorization", "Client-ID "+p.clientId)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return nil, genericError
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, genericError
	}

	return parseResponse(res)
}

func (p provider) UploadImage(filename string, imageData []byte) (*string, error) {

	// Here we stick our image into multipart form data and upload it to Imgur
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", filename)

	if err != nil {
		return nil, genericError
	}

	_, err = part.Write(imageData)

	if err != nil {
		return nil, genericError
	}

	err = writer.Close()
	if err != nil {
		return nil, genericError
	}

	req, err := http.NewRequest("POST", imgurUrl, body)

	req.Header.Add("Authorization", "Client-ID "+p.clientId)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, genericError
	}

	return parseResponse(resp)
}

func (p provider) Name() string {
	return "imgur"
}

func Provider(clientId string) providers.Provider {
	return provider{clientId: clientId}
}
