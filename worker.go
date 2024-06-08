package main

import (
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/adrium/goheif"
)

type heic struct {
	Path    string
	Correct bool
}

func ListHome() []heic {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	logs := ListDir(fmt.Sprintf("%s/%s", home, "Downloads"))

	return logs
}

func ListDir(dir string) []heic {
	heics := []heic{}
	files, _ := ioutil.ReadDir(dir)

	for _, f := range files {
		path := fmt.Sprintf("%s/%s", dir, f.Name())

		if f.IsDir() {
			heics = append(heics, ListDir(path)...)
			continue
		}

		if hasEnding(path) {
			heics = append(heics, heic{Path: path, Correct: isCorrectHeic(path)})
		}

	}

	return heics
}

func WorkFile(path string) bool {
	jpg := fmt.Sprintf("%s.jpg", path)

	if isWrongHeic(path) {
		return rename(path, jpg)
	}

	if isCorrectHeic(path) {
		return convert(path, jpg) && modified(path, jpg) && remove(path)
	}

	return false
}

//@todo: test.
func isWrongHeic(fin string) bool {

	mime := getFileContentType(fin)
	return mime == "image/jpeg"
}

//@todo: test.
func isCorrectHeic(fin string) bool {
	mime := getFileContentType(fin)
	return mime == "image/heic" || mime == "image/heif" || mime == "application/octet-stream"

}

//@todo: test.
func hasEnding(fin string) bool {
	if len(fin) < 6 {
		return false
	}
	return strings.ToLower(fin[len(fin)-4:]) == "heic"
}

func rename(fin string, fout string) bool {
	err := os.Rename(fin, fout)
	if err != nil {
		log.Printf("Warning: Cannot rename from %s to %s: %v\n", fin, fout, err)
		return false
	}

	log.Printf("Renamed %s to %s successfully\n", fin, fout)
	return true
}

func modified(fin string, fout string) bool {
	file, err := os.Stat(fin)

	if err != nil {
		log.Printf("Warning: Cannot read file to set original date: %s: %v\n", fin, err)
		return false
	}

	modifiedtime := file.ModTime()

	err = os.Chtimes(fout, modifiedtime, modifiedtime)

	if err != nil {
		log.Printf("Warning: Cannot change time of %s to %s: %v\n", fin, modifiedtime, err)
		return false
	}

	log.Printf("Modified %s times to: %s", fout, modifiedtime)
	return true
}

func remove(fin string) bool {
	err := os.Remove(fin)
	if err != nil {
		log.Printf("Warning: Cannot remove %s: %v\n", fin, err)
		return false
	}

	log.Printf("Removed %s successfully\n", fin)
	return true
}

func convert(fin string, fout string) bool {
	fi, err := os.Open(fin)
	if err != nil {
		log.Printf("%+v", err)
		return false
	}
	defer fi.Close()

	exif, err := goheif.ExtractExif(fi)
	if err != nil {
		log.Printf("Warning: no EXIF from %s: %v\n", fin, err)
	}

	img, err := goheif.Decode(fi)
	if err != nil {
		log.Printf("Failed to parse %s: %v\n", fin, err)
		return false

	}
	fo, err := os.OpenFile(fout, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Failed to create output file %s: %v\n", fout, err)
		return false
	}
	defer fo.Close()

	w, _ := newWriterExif(fo, exif)
	err = jpeg.Encode(w, img, nil)

	if err != nil {
		log.Printf("Failed to encode %s: %v\n", fout, err)
		return false
	}

	log.Printf("Converted %s to %s successfully\n", fin, fout)
	return true
}
