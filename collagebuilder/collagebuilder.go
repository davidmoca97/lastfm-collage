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

func BuildCollageFromData(albums []lastfm.Album) (*image.RGBA, error) {

	const albumCoverSize = config.AlbumCoverSize

	albumsPerRow := int(math.Sqrt(float64(config.DefaultGrid)))
	width := albumsPerRow * config.AlbumCoverSize
	height := width

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	c := make(chan lastfm.AlbumCoverDownloaderResponse)
	go lastfm.DownloadAlbumCovers(albums, c)

	for imgDownloadRespnose := range c {
		if imgDownloadRespnose.Err != nil {
			return nil, imgDownloadRespnose.Err
		}

		idx := imgDownloadRespnose.Idx

		currentRow := idx / albumsPerRow
		currentColumn := idx % albumsPerRow

		startingPoint := image.Point{currentColumn * albumCoverSize, currentRow * albumCoverSize}
		endingPoint := image.Point{albumCoverSize + (currentColumn * albumCoverSize), albumCoverSize + (currentRow * albumCoverSize)}
		r := image.Rectangle{startingPoint, endingPoint}
		draw.Draw(img, r, imgDownloadRespnose.Img, imgDownloadRespnose.Img.Bounds().Min, draw.Src)

		// Draw a shadow behind the text
		textShadow := image.Rectangle{
			Min: startingPoint,
			Max: image.Point{endingPoint.X, currentRow*config.AlbumCoverSize + 38},
		}
		shadowColor := color.NRGBA{0, 0, 0, 60}
		draw.Draw(img, textShadow, &image.Uniform{shadowColor}, image.ZP, draw.Over)

		albumLabel := albums[idx].Name
		artistlabel := albums[idx].Artist.Name
		playCountlabel := fmt.Sprintf("Plays: %d", albums[idx].PlayCount)
		addLabel(img, currentColumn*albumCoverSize+2, currentRow*albumCoverSize, albumLabel)
		addLabel(img, currentColumn*albumCoverSize+2, currentRow*albumCoverSize+12, artistlabel)
		addLabel(img, currentColumn*albumCoverSize+2, currentRow*albumCoverSize+24, playCountlabel)
	}

	return img, nil
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
