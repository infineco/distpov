package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

var mutex = &sync.Mutex{}

func page(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	sizeStr := params["size"]
	serialized := r.URL.Query().Get("serialized")
	size, _ := strconv.Atoi(sizeStr)
	if size == 0 {
		size = 1
	}
	if serialized != "" {
		serialized = "&serialized=true"
	}
	w.Write([]byte("<a href='./4'>try with 4x4</a><br>"))
	w.Write([]byte("<a href='./16'>try with 16x16</a><br>"))
	w.Write([]byte("<a href='./32'>try with 32x32</a><br>"))
	w.Write([]byte("<table>\n"))
	for j := 0; j < size; j++ {
		w.Write([]byte("<tr>\n"))
		for i := 0; i < size; i++ {
			w.Write([]byte(fmt.Sprintf("<td><img src=image?maxI=%s%s&maxJ=%s&i=%d&j=%d></td>\n", sizeStr, serialized, sizeStr, i, j)))

		}
		w.Write([]byte("</tr>\n"))
	}
	w.Write([]byte("</table>\n"))
}

func main() {
	http.HandleFunc("/image", handler)

	router := mux.NewRouter()
	router.HandleFunc("/{size:[0-9]*}", page).Methods("GET")
	http.Handle("/", router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on localhost:%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	iStr := r.URL.Query().Get("i")
	jStr := r.URL.Query().Get("j")
	maxIStr := r.URL.Query().Get("maxI")
	maxJStr := r.URL.Query().Get("maxJ")
	isLocked := r.URL.Query().Get("serialized")

	i, err := strconv.Atoi(iStr)
	j, err := strconv.Atoi(jStr)
	maxI, err := strconv.Atoi(maxIStr)
	maxJ, _ := strconv.Atoi(maxJStr)

	maxX := 800
	maxY := 600
	stepX := maxX / maxI
	fromX := stepX * i
	toX := fromX + stepX

	stepY := maxY / maxJ
	fromY := stepY * j
	toY := fromY + stepY
	filename := fmt.Sprintf("test%d.png", rand.Int())
	log.Println("Height="+strconv.Itoa(maxY), "Width="+strconv.Itoa(maxX),
		"Start_Column="+strconv.Itoa(fromX), "End_Column="+strconv.Itoa(toX),
		"Start_Row="+strconv.Itoa(fromY), "End_Row="+strconv.Itoa(toY), "Output_File_Name="+filename)
	if isLocked != "" {
		mutex.Lock()
		defer mutex.Unlock()
	}
	cmd := exec.Command("/opt/povray/bin/povray", "Height="+strconv.Itoa(maxY), "Width="+strconv.Itoa(maxX),
		"Start_Column="+strconv.Itoa(fromX), "End_Column="+strconv.Itoa(toX),
		"Start_Row="+strconv.Itoa(fromY), "End_Row="+strconv.Itoa(toY), "Output_File_Name="+filename, "test.pov")

	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		w.WriteHeader(500)
		log.Fatal(err)
	}
	log.Println(out)

	existingImageFile, err := os.Open(filename)
	if err != nil {
		// Handle error
	}
	defer existingImageFile.Close()

	my_image, err := png.Decode(existingImageFile)
	if err != nil {
		// Handle error
	}
	my_sub_image := my_image.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(fromX, fromY, toX, toY))

	fmt.Printf("bounds %v\n", my_sub_image.Bounds())

	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, my_sub_image)
}
