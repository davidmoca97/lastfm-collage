package main

import (
	"errors"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"net/http"
	"strconv"

	"github.com/davidmoca97/lastfm-collage/collagebuilder"
	"github.com/davidmoca97/lastfm-collage/lastfm"
	"github.com/gorilla/mux"
)

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

	params, err := getAndValidateParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	topAlbums, err := lastfm.GetTopAlbums(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	img, err := collagebuilder.BuildCollageFromData(topAlbums, params.Grid, params.IncludeLabels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	png.Encode(w, img)
	w.Header().Add("Content-Type", "image/png")
}

func getAndValidateParams(r *http.Request) (lastfm.GetTopAlbumsConfig, error) {
	username := r.URL.Query().Get("username")
	grid := r.URL.Query().Get("grid")
	period := r.URL.Query().Get("period")
	includeLabels := r.URL.Query().Get("includeLabels")

	var validGridSizes = []int{4, 9, 16, 25, 36, 49, 64, 81, 100}
	var validPeriods = []string{"overall", "7day", "1month", "3month", "6month", "12month"}

	if username == "" {
		return lastfm.GetTopAlbumsConfig{}, errors.New("Invalid \"username\" param. You must provide it")
	}

	gridInt, err := strconv.Atoi(grid)
	if err != nil {
		return lastfm.GetTopAlbumsConfig{}, errors.New("Invalid \"grid\" param. It must be an integer")
	}

	valid := false
	for _, size := range validGridSizes {
		if gridInt == size {
			valid = true
			break
		}
	}
	if !valid {
		return lastfm.GetTopAlbumsConfig{}, fmt.Errorf("Invalid \"grid\" param. The value must be any of these values: %v", validGridSizes)
	}

	valid = false
	for _, p := range validPeriods {
		if period == p {
			valid = true
			break
		}
	}
	if !valid {
		return lastfm.GetTopAlbumsConfig{}, fmt.Errorf("Invalid \"period\" param. The value must be any of these values: %v", validPeriods)
	}

	includeLabelsBool, err := strconv.ParseBool(includeLabels)
	if err != nil {
		return lastfm.GetTopAlbumsConfig{}, fmt.Errorf("Invalid \"includeLabels\" param. The value must be a boolean")
	}

	return lastfm.GetTopAlbumsConfig{Username: username, Grid: gridInt, Period: period, IncludeLabels: includeLabelsBool}, nil

}
