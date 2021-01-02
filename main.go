package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"math"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
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

const DefaultGrid = 49
const LastFmApiUrl = "https://ws.audioscrobbler.com/2.0"
const DefaultPeriod = "12month"
const LastFMApiKey = ""
const albumCoverSize = 300

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/collage", getCollage).Methods(http.MethodGet)
	http.ListenAndServe(":9999", router)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"message": "hola"}`))
}

func getCollage(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	if username == "" {
		http.Error(w, "Username was not provided", http.StatusBadRequest)
		return
	}

	topAlbums, err := getTopAlbums(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// s, _ := json.MarshalIndent(&topAlbums, "", "   ")
	// w.Header().Add("Content-Type", "application/json")
	// w.Write(s)

	img, err := getCollageFromData(topAlbums)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	png.Encode(w, img)
	w.Header().Add("Content-Type", "image/png")
}

func getTopAlbums(username string) ([]Album, error) {
	var lastFMResponse LastFMResponse
	r, err := http.Get(getURL(username))
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

func getURL(username string) string {
	return fmt.Sprintf("%s/?method=user.gettopalbums&format=json&api_key=%s&user=%s&period=%s&limit=%d&page=1", LastFmApiUrl, LastFMApiKey, username, DefaultPeriod, DefaultGrid)
}

func getCollageFromData(albums []Album) (*image.RGBA, error) {
	albumsPerRow := int(math.Sqrt(float64(DefaultGrid)))
	width := albumsPerRow * albumCoverSize
	height := width

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	c := make(chan AlbumCoverDownloaderResponse)
	go downloadAlbumCovers(albums, c)

	for imgDownloadRespnose := range c {
		if imgDownloadRespnose.Err != nil {
			return nil, imgDownloadRespnose.Err
		}

		idx := imgDownloadRespnose.Idx

		currentRow := idx / albumsPerRow
		currentColumn := idx % albumsPerRow

		fmt.Printf("🎇 %s: currentRow: %d   currentColumn: %d\n image: %s \n", albums[idx].Name, currentRow, currentColumn, "NOT_FOUND")

		// for y := 0; y < albumCoverSize; y++ {
		// 	for x := 0; x < albumCoverSize; x++ {
		// 		img.Set(x+(currentColumn*albumCoverSize), y+(currentRow*albumCoverSize), imgDownloadRespnose.Img.At(x, y))
		// 	}
		// }
		startingPoint := image.Point{currentColumn * albumCoverSize, currentRow * albumCoverSize}
		endingPoint := image.Point{albumCoverSize + (currentColumn * albumCoverSize), albumCoverSize + (currentRow * albumCoverSize)}
		r := image.Rectangle{startingPoint, endingPoint}
		draw.Draw(img, r, imgDownloadRespnose.Img, imgDownloadRespnose.Img.Bounds().Min, draw.Src)

		albumLabel := albums[idx].Name
		artistlabel := albums[idx].Artist.Name
		playCountlabel := fmt.Sprintf("Plays: %d", albums[idx].PlayCount)
		addLabel(img, currentColumn*albumCoverSize+2, currentRow*albumCoverSize+10, albumLabel)
		addLabel(img, currentColumn*albumCoverSize+2, currentRow*albumCoverSize+22, artistlabel)
		addLabel(img, currentColumn*albumCoverSize+2, currentRow*albumCoverSize+34, playCountlabel)
	}

	return img, nil
}

func downloadAlbumCovers(albums []Album, c chan<- AlbumCoverDownloaderResponse) {
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

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{255, 255, 255, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}
