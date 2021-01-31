package collagebuilder

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/davidmoca97/lastfm-collage/config"
	"github.com/davidmoca97/lastfm-collage/lastfm"
	"github.com/golang/freetype"
)

func BuildCollageFromData(albums []lastfm.Album, gridSize int, includeLabel bool) (*image.RGBA, error) {

	const albumCoverSize = config.AlbumCoverSize

	albumsPerRow := int(math.Sqrt(float64(gridSize)))
	width := albumsPerRow * config.AlbumCoverSize
	height := width

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	c := make(chan lastfm.AlbumCoverDownloaderResponse)
	go lastfm.DownloadAlbumCovers(albums, c)

	for imgDownloadRespnose := range c {

		idx := imgDownloadRespnose.Idx

		currentRow := idx / albumsPerRow
		currentColumn := idx % albumsPerRow

		startingPoint := image.Point{currentColumn * albumCoverSize, currentRow * albumCoverSize}
		endingPoint := image.Point{albumCoverSize + (currentColumn * albumCoverSize), albumCoverSize + (currentRow * albumCoverSize)}
		r := image.Rectangle{startingPoint, endingPoint}
		draw.Draw(img, r, imgDownloadRespnose.Img, imgDownloadRespnose.Img.Bounds().Min, draw.Src)

		if includeLabel {
			addLabelBackground(img, startingPoint, image.Point{endingPoint.X, currentRow*config.AlbumCoverSize + 38})

			albumLabel := albums[idx].Name
			artistlabel := albums[idx].Artist.Name
			playCountlabel := fmt.Sprintf("%d plays", albums[idx].PlayCount)
			addLabel(img, currentColumn*albumCoverSize+2, currentRow*albumCoverSize, albumLabel)
			addLabel(img, currentColumn*albumCoverSize+2, currentRow*albumCoverSize+12, artistlabel)
			addLabel(img, currentColumn*albumCoverSize+2, currentRow*albumCoverSize+24, playCountlabel)
		}
	}

	return img, nil
}

func addLabelBackground(img draw.Image, min, max image.Point) {
	textShadow := image.Rectangle{
		Min: min,
		Max: max,
	}
	shadowColor := color.NRGBA{0, 0, 0, 60}
	draw.Draw(img, textShadow, &image.Uniform{shadowColor}, image.ZP, draw.Over)
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{255, 255, 255, 255}

	f := getFontContext()
	pt := freetype.Pt(x, y+int(f.PointToFixed(12)>>6))
	f.SetClip(img.Bounds())
	f.SetDst(img)
	f.SetSrc(image.NewUniform(col))
	f.DrawString(label, pt)
}

func getFontContext() *freetype.Context {
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(config.Font)
	c.SetFontSize(12)

	return c
}
