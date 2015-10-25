# MetaUpload
This service allows you to upload or sideload images to a variety of image hosts using a common web service API.

## API

All responses are in JSON and follow this format:

```
{
	"success": true/false,
	"error": null if successful, otherwise contains error string,
	"data": return value of API function if successful, null otherwise
}
```

### GET `/providers`

#### Returns

JSON array of provider names for use in other API calls.

### POST `/providers/:provider_name/sideload`

#### Arguments

* *url:* the URL of the image to sideload

#### Returns
URL string of sideloaded image on image host.

### POST `/providers/:provider_name/upload`

#### Arguments

* *image:* multipart file upload containing the image you wish to upload

#### Returns
URL string of uploaded image on image host.

### Upload

## Code Structure
Each image provider is loaded as a plugin implementing the interface defined in the providers package.

## Configuration
Each image provider has its own top level initialization function that may take one or more configuration parameters.
The service loads all necessary configs from `config.json`. You can see a sample version of the config in `config.sample.json`.

### Imgur
Imgur requires uploaders to provide a Client ID, which can be obtained by registering an app with their API.