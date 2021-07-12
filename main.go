package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

type XMLAnimeTitles struct {
	AnimeList []XMLAnimeTitlesAnime `xml:"anime"`
}

type XMLAnimeTitlesAnime struct {
	ID     int                        `xml:"aid,attr" json:"id"`
	Titles []XMLAnimeTitlesAnimeTitle `xml:"title" json:"titles"`
}

type XMLAnimeTitlesAnimeTitle struct {
	Language string `xml:"lang,attr" json:"language"`
	Type     string `xml:"type,attr" json:"type"`
	Content  string `xml:",chardata" json:"title"`
}

func main() {
	home := os.Getenv("HOME")
	dumpsDir := path.Join(home, "backups", "anidb", "animetitles")

	var (
		err error
		dd  *os.File
	)

	dd, err = os.Open(dumpsDir)
	if err != nil {
		panic(err)
	}
	defer dd.Close()

	var names []string

	names, err = dd.Readdirnames(-1)
	if err != nil {
		panic(err)
	}

	alreadyProcessed := map[string]bool{}
	for _, name := range names {
		parts := strings.Split(name, ".")
		if len(parts) != 3 {
			continue
		}

		timestamp := parts[0]
		ext := parts[1]

		if ext == "json" {
			alreadyProcessed[timestamp] = true
		} else {
			alreadyProcessed[timestamp] = alreadyProcessed[timestamp]
		}
	}

	for timestamp, skip := range alreadyProcessed {
		result := func() string {
			if skip {
				return ""
			}

			fileName := fmt.Sprintf("%s.xml.xz", timestamp)

			inputFile := path.Join(dumpsDir, fileName)
			outputFile := path.Join(dumpsDir, fmt.Sprintf("%s.json", timestamp))

			cmd := exec.Command("xzcat", inputFile)

			var rawInput []byte
			rawInput, err = cmd.Output()
			if err != nil {
				panic(err)
			}

			var w *os.File
			w, err = os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
			if err != nil {
				panic(err)
			}
			defer w.Close()

			dump := &XMLAnimeTitles{}

			err = xml.Unmarshal(rawInput, &dump)
			if err != nil {
				panic(err)
			}

			for _, anime := range dump.AnimeList {
				var rawJSON []byte

				rawJSON, err = json.Marshal(anime)
				if err != nil {
					panic(err)
				}

				fmt.Fprintln(w, string(rawJSON))
			}

			return outputFile
		}()

		if result != "" {
			cmd := exec.Command("xz", "-9", result)

			err = cmd.Run()
			if err != nil {
				panic(err)
			}
		}
	}
}
