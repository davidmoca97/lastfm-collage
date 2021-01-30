package lastfm

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"net/http"

	"github.com/davidmoca97/lastfm-collage/config"
	"github.com/davidmoca97/lastfm-collage/util"
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
	// fmt.Println(getURL(configuration))
	decoder.Decode(&lastFMResponse)
	// fmt.Println(lastFMResponse.TopAlbums)
	defer r.Body.Close()
	lastFMResponse.TopAlbums.Album[10].Image[len(lastFMResponse.TopAlbums.Album[10].Image)-1].URL = ""
	lastFMResponse.TopAlbums.Album[11].Image[len(lastFMResponse.TopAlbums.Album[11].Image)-1].URL = "https://www.callicoder.com/golang-pointers/"
	return lastFMResponse.TopAlbums.Album, nil
}

func getURL(configuration GetTopAlbumsConfig) string {
	return fmt.Sprintf("%s/?method=user.gettopalbums&format=json&api_key=%s&user=%s&period=%s&limit=%d&page=1", config.LastFmApiUrl, config.LastFMApiKey, configuration.Username, configuration.Period, configuration.Grid)
}

func DownloadAlbumCovers(albums []Album, c chan<- AlbumCoverDownloaderResponse) {
	for idx, album := range albums {

		var img image.Image
		var err error
		var albumCoverURL string
		if len(album.Image) > 0 {
			albumCoverURL = album.Image[len(album.Image)-1].URL
		}

		if albumCoverURL != "" {
			img, err = util.GetImageFromURL(albumCoverURL)
		} else {
			img = config.DefaultAlbumCover
		}

		c <- AlbumCoverDownloaderResponse{
			Img: img,
			Err: err,
			Idx: idx,
		}
	}
	fmt.Println("Downloaded everything")
	close(c)
}
