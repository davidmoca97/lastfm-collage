package lastfm

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"net/http"

	"github.com/davidmoca97/lastfm-collage/config"
)

type Album struct {
	Artist struct {
		URL  string `json:"url"`
		Name string `json:"name"`
	} `json:"artist"`
	Image []struct {
		Size string `json:"size"`
		URL  string `json:"#text"`
	} `json:"image"`
	PlayCount int    `json:"playcount,string"`
	URL       string `json:"url"`
	Name      string `json:"name"`
}

type LastFMResponse struct {
	TopAlbums struct {
		Album []Album `json:"album"`
	} `json:"topalbums"`
}

type AlbumCoverDownloaderResponse struct {
	Idx int
	Img image.Image
	Err error
}

type GetTopAlbumsConfig struct {
	Username string
	Grid     int
	Period   string
}

func GetTopAlbums(configuration GetTopAlbumsConfig) ([]Album, error) {
	var lastFMResponse LastFMResponse
	r, err := http.Get(getURL(configuration))
	if err != nil {
		return []Album{}, err
	}
	if r.StatusCode != http.StatusOK {
		return []Album{}, errors.New("Something wrong happened while fetching the top albums")
	}

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&lastFMResponse)
	defer r.Body.Close()
	return lastFMResponse.TopAlbums.Album, nil
}

func getURL(configuration GetTopAlbumsConfig) string {
	return fmt.Sprintf("%s/?method=user.gettopalbums&format=json&api_key=%s&user=%s&period=%s&limit=%d&page=1", config.LastFmApiUrl, config.LastFMApiKey, configuration.Username, configuration.Period, configuration.Grid)
}

func DownloadAlbumCovers(albums []Album, c chan<- AlbumCoverDownloaderResponse) {
	for idx, album := range albums {
		albumCoverURL := album.Image[len(album.Image)-1].URL
		img, err := getImageFromURL(albumCoverURL)
		c <- AlbumCoverDownloaderResponse{
			Img: img,
			Err: err,
			Idx: idx,
		}
	}
	fmt.Println("Downloaded everything")
	close(c)
}

func getImageFromURL(url string) (image.Image, error) {
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
