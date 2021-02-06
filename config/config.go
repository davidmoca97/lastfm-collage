package config

import (
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"

	"github.com/davidmoca97/lastfm-collage/util"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

const LastFmApiUrl = "https://ws.audioscrobbler.com/2.0"
const AlbumCoverSize = 300
const FontFile = "./static/fonts/Lato-Medium.ttf"

var (
	LastFMApiKey      string
	Port              string
	Font              *truetype.Font
	DefaultAlbumCover image.Image
)

func init() {

	key, exists := os.LookupEnv("LAST_FM_APY_KEY")
	if !exists {
		log.Fatal("Error: No LastFM Api key was provided")
		return
	}
	LastFMApiKey = key

	key, exists = os.LookupEnv("PORT")
	if !exists {
		log.Fatal("Error: No LastFM Api key was provided")
		return
	}
	LastFMApiKey = key

	if err := initializeFont(); err != nil {
		log.Fatal("Error loading font:", err)
		return
	}

	defaultImageURL := fmt.Sprintf("https://via.placeholder.com/%dx%d?text=Unknown+album+cover", AlbumCoverSize, AlbumCoverSize)
	img, err := util.GetImageFromURL(defaultImageURL)
	if err != nil {
		log.Fatal("Error loading default album cover image. Error:", err)
		return
	}
	DefaultAlbumCover = img

}

func initializeFont() error {
	fontBytes, err := ioutil.ReadFile(FontFile)
	if err != nil {
		return err
	}
	Font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}
	return nil
}
