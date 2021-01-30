package util

import (
	"fmt"
	"image"
	"net/http"
)

func GetImageFromURL(url string) (image.Image, error) {
	var img image.Image
	r, err := http.Get(url)
	if err != nil {
		return img, err
	}
	if r.StatusCode != http.StatusOK {
		return img, fmt.Errorf("Something wrong happened while getting the image from %s", url)
	}

	img, _, err = image.Decode(r.Body)
	defer r.Body.Close()
	if err != nil {
		return img, fmt.Errorf("Error while decoding the image from %s: %s", url, err)
	}

	return img, nil
}
