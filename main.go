package main

import (
	"flag"
	"github.com/FcoManueel/imtdigest/imthash"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	var filepath, url string
	var bytesPerSecond int
	// TODO Remove default values for string flags
	flag.StringVar(&filepath, "file", "/tmp/imt_hash.txt", "Filepath for output (required)")
	flag.StringVar(&url, "url", "https://jsonplaceholder.typicode.com/todos/1", "Url to download (required)")
	flag.IntVar(&bytesPerSecond, "throttle", 0, "Max download rate in bytes per second (optional)")
	flag.Parse()

	if filepath == "" {
		log.Fatalf("Please provide a valid output file (-file flag)")
	}
	if url == "" {
		log.Fatalf("Please provide a valid url (-url flag)")
	}

	stream, err := download(url, bytesPerSecond)
	if err != nil {
		log.Fatalf("Error while attempting to fetch '%s': %v", url, err)
	}

	hasher := &imthash.Hash{}
	io.Copy(hasher, stream)

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
