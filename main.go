package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
)

func connectToMongoDB() {
	// Set up MongoDB connection options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the MongoDB server to check if the connection was successful
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	mongoClient = client
	log.Println("Connected to MongoDB")
}

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

	// Insert the video record into the MongoDB collection
	video := Video{
		Filename:   handler.Filename,
		Path:       f.Name(),
		UploadedAt: time.Now(),
	}
	collection := mongoClient.Database("videoapp").Collection("videos")
	_, err = collection.InsertOne(context.TODO(), video)
	if err != nil {
		log.Println("Error inserting video record into MongoDB:", err)
		http.Error(w, "Error inserting video record into MongoDB", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type Video struct {
	Filename   string    `bson:"filename"`
	Path       string    `bson:"path"`
	UploadedAt time.Time `bson:"uploadedAt"`
}

func main() {
	// Connect to MongoDB
	connectToMongoDB()

	// Serve the static files
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	// Handle the upload route
	http.HandleFunc("/upload", uploadHandler)

	// Start the server
	log.Println("Server is running on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
