package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}
	cleanName := filepath.Base(name)
	filePath := filepath.Join("storage", cleanName)
	err := os.MkdirAll("./storage", 0755)
	if err != nil {
		log.Printf("ERROR: failed to create storage directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			http.Error(w, "The file already exists. Do not overwrite", http.StatusConflict)
			return
		}
		log.Printf("ERROR: failed to create file: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	io.Copy(file, r.Body)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}
	cleanName := filepath.Base(name)
	filePath := filepath.Join("storage", cleanName)
	file, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "File or directory not found", http.StatusNotFound)
			return
		}
		log.Printf("ERROR: failed to create file %s: %v", filePath, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	io.Copy(w, file)
}
func listHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("storage")
	if err != nil {
		log.Printf("ERROR: failed to read storage directory: %v", err)
		http.Error(w, "Unable to read storage directory", http.StatusInternalServerError)
		return
	}
	if len(files) == 0 {
		fmt.Fprintln(w, "dir is empty")
		return
	}
	for _, file := range files {
		fmt.Fprintln(w, file.Name())
	}
}

func main() {
	http.HandleFunc("/update", uploadHandler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/list", listHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

}
