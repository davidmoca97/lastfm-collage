package config

import (
	"fmt"
	"image"
	"io/ioutil"
	"log"

	"github.com/davidmoca97/lastfm-collage/util"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

const DefaultGrid = 25
const LastFmApiUrl = "https://ws.audioscrobbler.com/2.0"
const DefaultPeriod = "12month"
const LastFMApiKey = ""
const AlbumCoverSize = 300
const FontFile = "./static/fonts/Lato-Medium.ttf"

var Font *truetype.Font
var DefaultAlbumCover image.Image

func init() {
	if err := initializeFont(); err != nil {
		log.Fatal("Error loading font:", err)
		return
	}

	defaultImageURL := fmt.Sprintf("https://via.placeholder.com/%dx%d?text=Unknown+album+cover", AlbumCoverSize, AlbumCoverSize)
	img, err := util.GetImageFromURL(defaultImageURL)
	if err != nil {
		log.Fatal("Error loading default album cover image")
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
