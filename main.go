package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Limit request size to prevent abuse
	r.ParseMultipartForm(10 << 20) // 10MB

	// Get the uploaded file
	file, handler, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Error retrieving the video file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a new file on the server to store the uploaded video
	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("Error creating the video file:", err)
		http.Error(w, "Error creating the video file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Copy the uploaded file to the server's file
	_, err = io.Copy(f, file)
	if err != nil {
		log.Println("Error copying the video file:", err)
		http.Error(w, "Error copying the video file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	// Serve the static files
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	// Handle the upload route
	http.HandleFunc("/upload", uploadHandler)

	// Start the server
	log.Println("Server is running on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
