package main

import (
	"flag"
	"github.com/FcoManueel/imtdigest/imthash"
	"github.com/FcoManueel/imtdigest/slowread"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	var filepath, url string
	var bytesPerSecond int
	flag.StringVar(&filepath, "file", "", "The path of the file where the output will be saved. (required)")
	flag.StringVar(&url, "url", "", "The URL from which the data will be fetched. (required)")
	flag.IntVar(&bytesPerSecond, "rate", 0, "Limits the max download rate. Units are in bytes/second. (optional)")
	flag.Parse()

	if filepath == "" || url == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Validate path before downloading in order to fail fast
	log.Println("Validating filepath...")
	if !isValidPath(filepath) {
		log.Fatalf("Cannot write to '%s'. Check that the directory exists, permissions are sufficient and there is not already a file with that name.", filepath)
	}

	log.Println("Starting download...")
	stream, err := download(url, bytesPerSecond)
	if err != nil {
		log.Fatalf("Error while attempting to fetch '%s': %v", url, err)
	}
	defer stream.Close()

	hasher := &imthash.Hash{}
	_, err = io.Copy(hasher, stream)
	if err != nil {
		log.Fatal("An error occurred while hashing the download: ", err.Error())
	}

	log.Println("Saving to file...")
	hexHash := hasher.Hex()
	err = writeFile(filepath, hexHash)
	if err != nil {
		log.Fatalf("Error while trying to write hash '%s' to file '%s'. Make sure path exists, permissions are correct and there is disk space available", hexHash, filepath)
	}

	log.Println("All done, have a nice day!")
}

func download(url string, bytesPerSecond int) (io.ReadCloser, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	var reader io.ReadCloser = res.Body
	if bytesPerSecond > 0 {
		log.Printf("Downloading at %d bytes per second", bytesPerSecond)
		reader = slowread.NewReader(reader, bytesPerSecond)
	}
	return reader, nil
}

// isValidPath returns true if the path can be written to and a file with that name doesn't exist
func isValidPath(filepath string) bool {
	if _, err := os.Stat(filepath); err == nil {
		return false
	}
	if err := ioutil.WriteFile(filepath, nil, 0644); err == nil {
		os.Remove(filepath) // File can be saved. Clean up
		return true
	}
	return false
}

func writeFile(filepath, data string) (err error) {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() {
		if e := f.Close(); e != nil && err == nil {
			err = e
		}
	}()
	_, err = f.WriteString(data)
	return err
}
