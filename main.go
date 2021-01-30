package main

import (
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"net/http"

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
	username := r.URL.Query().Get("username")

	if username == "" {
		http.Error(w, "Username was not provided", http.StatusBadRequest)
		return
	}

	topAlbums, err := lastfm.GetTopAlbums(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	img, err := collagebuilder.BuildCollageFromData(topAlbums)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	png.Encode(w, img)
	w.Header().Add("Content-Type", "image/png")
}
