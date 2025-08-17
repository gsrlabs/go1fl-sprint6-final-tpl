package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

// HandleMain serves the main HTML page.
func HandleMain(w http.ResponseWriter, r *http.Request) {
   data, err := os.ReadFile("index.html")
	if err != nil {
		log.Fatal(err)
	}
   	w.Header().Set("Content-Type", "text/html")
    w.Write(data)
}

// HandleUpload processes file uploads and converts the content.
func HandleUpload(w http.ResponseWriter, r *http.Request) {

    err := r.ParseMultipartForm(10 << 20)
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    file, handler, err := r.FormFile("myFile")
    if err != nil {
        http.Error(w, "Error retrieving the file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    fileBytes, err := io.ReadAll(file)
    if err != nil {
        http.Error(w, "Error reading the file", http.StatusInternalServerError)
        return
    }

    uploadData, err := service.TextHandler(string(fileBytes))
    if err != nil {
        log.Println(err.Error())
        http.Error(w, "Data conversion error", http.StatusInternalServerError)
        return
    }

    timestamp := time.Now().UTC().Format("2006-01-02_15-04-05")
    ext := filepath.Ext(handler.Filename)
    if ext == "" {
        ext = ".txt"
    }
    outputFilename := "result_" + timestamp + ext

    outputFile, err := os.Create(outputFilename)
    if err != nil {
        log.Println("Error creating output file:", err)
        http.Error(w, "Error creating output file", http.StatusInternalServerError)
        return
    }
    defer outputFile.Close()

    _, err = outputFile.WriteString(uploadData)
    if err != nil {
        log.Println("Error writing to output file:", err)
        http.Error(w, "Error writing output file", http.StatusInternalServerError)
        return
    }

    w.Write([]byte(uploadData))
}