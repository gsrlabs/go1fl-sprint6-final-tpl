package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

const (
	// Common error messages
	ErrInternalServer      = "Internal Server Error"
	ErrParseForm           = "Unable to parse form"
	ErrRetrieveFile        = "Error retrieving the file"
	ErrReadFile            = "Error reading the file"
	ErrDataConversion      = "Data conversion error"
	ErrCreateLogDir        = "Error creating logs directory"
	ErrCreateOutputFile    = "Error creating output file"
	ErrWriteOutputFile     = "Error writing output file"
	
	// Log messages
	LogErrReadIndex       = "Error reading index.html: %v"
	LogErrCleanOldLogs    = "Error cleaning old logs: %v"
	LogErrCreateLogDir    = "Error creating logs directory: %v"
	LogErrCreateOutput    = "Error creating output file: %v"
	LogErrWriteOutput     = "Error writing to output file: %v"
	LogSuccessSave        = "Successfully saved result to %s"
	LogDeletedOldFile     = "Deleted old log file: %s"
	LogCleanupStats       = "Cleaned up %d old log files"
)

// HandleMain serves the main HTML page.
func HandleMain(w http.ResponseWriter, r *http.Request) {
    // Используйте правильный путь к файлу index.html
    data, err := os.ReadFile("index.html")
    if err != nil {
        log.Printf(LogErrReadIndex, err)
        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
        return
    }
    
    // Явно устанавливаем Content-Type как text/html
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.Write(data)
}

// HandleUpload processes file uploads and converts the content.
func HandleUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, ErrParseForm, http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, ErrRetrieveFile, http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, ErrReadFile, http.StatusBadRequest)
		return
	}

	uploadData, err := service.TextHandler(string(fileBytes))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, ErrDataConversion, http.StatusInternalServerError)
		return
	}

    logDir := "logs"
    if err := os.MkdirAll(logDir, 0755); err != nil {
        log.Printf(LogErrCreateLogDir, err)
        http.Error(w, ErrCreateLogDir, http.StatusInternalServerError)
        return
    }

	if err := cleanOldLogs(logDir, 7*24*time.Hour); err != nil {
        log.Printf(LogErrCleanOldLogs, err)
    }

    timestamp := time.Now().UTC().Format("2006-01-02_15-04-05")
    ext := filepath.Ext(handler.Filename)
    if ext == "" {
        ext = ".txt"
    }
    outputFilename := filepath.Join(logDir, "result_"+timestamp+ext)

    outputFile, err := os.Create(outputFilename)
    if err != nil {
        log.Printf(LogErrCreateOutput, err)
        http.Error(w, ErrCreateOutputFile, http.StatusInternalServerError)
        return
    }
    defer outputFile.Close()

	log.Printf(LogSuccessSave, outputFilename)

    _, err = outputFile.WriteString(uploadData)
    if err != nil {
        log.Printf(LogErrWriteOutput, err)
        http.Error(w, ErrWriteOutputFile, http.StatusInternalServerError)
        return
    }
	
	w.Write([]byte(uploadData))
}

// cleanOldLogs removes log files older than maxAge from the specified directory
func cleanOldLogs(logDir string, maxAge time.Duration) error {
    entries, err := os.ReadDir(logDir)
    if err != nil {
        return fmt.Errorf("could not read log directory: %w", err)
    }

    now := time.Now()
    var deleted int
    var errors []error

    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }

        info, err := entry.Info()
        if err != nil {
            errors = append(errors, err)
            continue
        }

        if now.Sub(info.ModTime()) > maxAge {
            path := filepath.Join(logDir, entry.Name())
            if err := os.Remove(path); err != nil {
                errors = append(errors, fmt.Errorf("could not delete %s: %w", path, err))
            } else {
                deleted++
                log.Printf(LogDeletedOldFile, path)
            }
        }
    }

    log.Printf(LogCleanupStats, deleted)
    if len(errors) > 0 {
        return fmt.Errorf("encountered %d errors during cleanup: %v", len(errors), errors)
    }
    return nil
}
