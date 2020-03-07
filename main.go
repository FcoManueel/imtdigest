package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	// TODO Add CLI flags, parse and validation
	filepath := "/tmp/imt_hash.txt"
	bytesPerSecond := 0
	url := "https://jsonplaceholder.typicode.com/todos/1"

	stream, err := download(url, bytesPerSecond)
	if err != nil {
		log.Fatalf("Error while attempting to fetch '%s': %v", url, err)
	}

	hasher := &Hash{}
	// TODO Connect download stream stream to hasher
	content, _ := ioutil.ReadAll(stream)
	fmt.Println(string(content))
	_ = stream

	hexHash := hasher.Hex()
	err = writeFile(filepath, hexHash)
	if err != nil {
		log.Fatalf("Error while trying to write hash '%s' to file '%s'. Make sure path exists, permissions are correct and there is disk space available", hexHash, filepath)
	}
	// TODO Update README.md
}

func download(url string, bytesPerSecond int) (io.ReadCloser, error) {
	// TODO Implement throttling for download
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func writeFile(filepath, data string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close() // TODO handle potential errors occurring here
	_, err = f.WriteString(data)
	return err
}
