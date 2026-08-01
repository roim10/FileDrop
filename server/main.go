package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var writer *bufio.Writer

func logEvent(operation, name, status string) {
	date := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("%v %v %v %v\n", date, operation, name, status)
	writer.WriteString(msg)
	writer.Flush()
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST method is allowed", http.StatusBadRequest)
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		logEvent("POST /upload", name, "ERROR: file name is required")
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}
	cleanName := filepath.Base(name)
	filePath := filepath.Join("storage", cleanName)
	err := os.MkdirAll("./storage", 0755)
	if err != nil {
		logEvent("POST /upload", name, "ERROR: failed to create storage directory")
		log.Printf("ERROR: failed to create storage directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			logEvent("POST /upload", name, "ERROR: the file already exists. Do not overwrite")
			http.Error(w, "The file already exists. Do not overwrite", http.StatusConflict)
			return
		}
		logEvent("POST /upload", name, "ERROR: failed to create file ")
		log.Printf("ERROR: failed to create file: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	io.Copy(file, r.Body)
	logEvent("upload", name, "OK")
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET method is allowed", http.StatusBadRequest)
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		logEvent("GET /download", name, "ERROR: file name is required")
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}
	cleanName := filepath.Base(name)
	filePath := filepath.Join("storage", cleanName)
	file, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logEvent("GET /download", name, "ERROR: file or directory not found")
			http.Error(w, "File or directory not found", http.StatusNotFound)
			return
		}
		logEvent("GET /download", name, "ERROR: failed to create file")
		log.Printf("ERROR: failed to create file %s: %v", filePath, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	io.Copy(w, file)
	logEvent("download", name, "OK")
}
func listHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("storage")
	if err != nil {
		logEvent("list", "-", "ERROR: failed to read storage directory")
		log.Printf("ERROR: failed to read storage directory: %v", err)
		http.Error(w, "Unable to read storage directory", http.StatusInternalServerError)
		return
	}
	if len(files) == 0 {
		logEvent("list", "-", "OK")
		fmt.Fprintln(w, "dir is empty")
		return
	}
	for _, file := range files {
		fmt.Fprintln(w, file.Name())
	}
	logEvent("list", "-", "OK")
}

func main() {
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	writer = bufio.NewWriter(file)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/list", listHandler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
