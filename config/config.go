package config

import (
	"io/ioutil"
	"log"

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

func init() {
	if err := initializeFont(); err != nil {
		log.Println("Error loading font:", err)
		return
	}
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
