package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/carlogit/phash"
)

const (
	maxDistance = 10
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("usage: %s <sourceFile> <searchInFolder>\n", os.Args[0])
	}

	sourceFile := os.Args[1]
	destFolder := os.Args[2]

	fmt.Printf("File's hash calculating... ")
	destHashFileMap := hashFilesInFolder(destFolder)
	fmt.Printf("%v file's hash calculated\n", len(destHashFileMap))

	fmt.Printf("Find best matching file for: %v\n", sourceFile)
	sourceHash := hash(sourceFile)
	distance, bestMatchingFile := findBestMatching(sourceHash, destHashFileMap)
	if bestMatchingFile != "" {
		fmt.Printf("Result: %v with distance = %v\n", destFolder+"/"+bestMatchingFile, distance)
	} else {
		fmt.Printf("Result: Opps, no result")
	}
}

func findBestMatching(sourceHash string, destHashFileMap map[string]string) (minDistance int, bestMatchingFile string) {
	minDistance = maxDistance
	for destName, destHash := range destHashFileMap {
		distance := phash.GetDistance(sourceHash, destHash)
		if distance < minDistance {
			minDistance = distance
			bestMatchingFile = destName
		}

		if minDistance == 0 {
			break
		}
	}

	return
}

func getListFile(folderPath string) []os.FileInfo {
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		log.Fatal(err)
		return []os.FileInfo{}
	}

	return files
}

func hashFilesInFolder(folderPath string) (fileHashMap map[string]string) {
	fileHashMap = map[string]string{}
	for _, file := range getListFile(folderPath) {
		fileHashMap[file.Name()] = hash(folderPath + "/" + file.Name())
	}
	return
}

//@Mark: https://stackoverflow.com/questions/38466804/comparing-base64-image-strings-in-golang
//hash returns a phash of the image
func hash(filename string) string {
	img, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer img.Close()

	ahash, err := phash.GetHash(img)
	if err != nil {
		log.Fatal(err)
	}
	return ahash
}
