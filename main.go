package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CopyImageSchema struct {
	// Source: from where to copy the image
	Source string
	// Destination: to where should the image be copied
	Destination string
	// Origin: Original reference to the image
	Origin string
}

func main() {
	history, newMirroring, err := readFiles()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	diffDigestsOnly := diff(history, newMirroring)

	updateHistory(history, diffDigestsOnly)
}

func readFiles() (history map[string]struct{}, newMirroring []CopyImageSchema, err error) {
	mirroringHistoryFile, err := os.Open(os.Args[1])
	if err != nil {
		return nil, nil, err
	}
	defer mirroringHistoryFile.Close()

	newMirroringFile, err := os.Open(os.Args[2])
	if err != nil {
		return nil, nil, err
	}
	defer newMirroringFile.Close()

	historyScanner := bufio.NewScanner(mirroringHistoryFile)
	newMirroringScanner := bufio.NewScanner(newMirroringFile)

	history = make(map[string]struct{})
	newMirroring = []CopyImageSchema{}

	for historyScanner.Scan() {
		history[historyScanner.Text()] = struct{}{}
	}

	for newMirroringScanner.Scan() {
		imgRef := CopyImageSchema{Origin: newMirroringScanner.Text()}
		newMirroring = append(newMirroring, imgRef)
	}

	return history, newMirroring, nil
}

func diff(history map[string]struct{}, newMirroring []CopyImageSchema) []string {
	var diff []CopyImageSchema
	var diffDigestsOnly []string

	for _, imgRef := range newMirroring {
		pathComponents := strings.Split(imgRef.Origin, "@")
		var digest string
		if len(pathComponents) > 1 {
			if strings.Contains(pathComponents[1], ":") {
				digest = strings.Split(pathComponents[1], ":")[1]
			} else {
				digest = pathComponents[1]
			}
			if _, isPresent := history[digest]; !isPresent {
				fmt.Printf("Digest %s is not present in the history, adding it to the diff\n", digest)
				//This diff variable simulates the output of the allRelatedImages var in oc-mirror
				//We are not using it for the purpose of this PoC
				diff = append(diff, imgRef)
				diffDigestsOnly = append(diffDigestsOnly, digest)
			}
		}
	}

	return diffDigestsOnly
}

func updateHistory(history map[string]struct{}, diffDigestsOnly []string) {
	fmt.Println("")
	fmt.Println("Updating the history with the new digests added in the last mirroring")
	for _, diff := range diffDigestsOnly {
		history[diff] = struct{}{}
	}

	//Here we should write the new history to the file, but for the purpose of this PoC we are just printing it
	for digest := range history {
		fmt.Println(digest)
	}
}
