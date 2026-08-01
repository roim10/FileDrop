package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func clientList() {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	respList, err := client.Get("http://localhost:8080/list")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer respList.Body.Close()
	io.Copy(os.Stdout, respList.Body)
}

func clientUpdate(name string) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	finalURL := fmt.Sprintf("http://localhost:8080/upload?name=%v", name)
	defer file.Close()
	respUp, err := client.Post(finalURL, "application/octet-stream", file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer respUp.Body.Close()
	io.Copy(os.Stdout, respUp.Body)
}
func clientDownload(name string) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	finalURL := fmt.Sprintf("http://localhost:8080/download?name=%v", name)
	respDow, err := client.Get(finalURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer respDow.Body.Close()
	finalName := fmt.Sprintf("downloaded_%v", name)
	file1, err := os.Create(finalName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file1.Close()
	io.Copy(file1, respDow.Body)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
MainLoop:
	for {
		fmt.Println("Выберите действие:\n 1 - загрузить файл\n 2 - скачать файл\n 3 - показать список файлов\n 4 - выйти")
		scanner.Scan()
		str := scanner.Text()
		switch str {
		case "1":
			fmt.Println("Введите имя файла, учтите, что нельзя использовать пробелы. также пожалуйста напишите формат файла через точку")
			scanner.Scan()
			name := scanner.Text()
			clientUpdate(name)
		case "2":
			fmt.Println("Введите имя файла, учтите, что нельзя использовать пробелы. также пожалуйста напишите формат файла через точку")
			scanner.Scan()
			name := scanner.Text()
			clientDownload(name)
		case "3":
			clientList()
		case "4":
			fmt.Println("спасибо что используете нашу программу")
			break MainLoop
		}
	}
}
