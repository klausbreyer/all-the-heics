package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/adrium/goheif"
)

func main() {
	flag.Parse()

	// fin, fout := flag.Arg(0), flag.Arg(1)
	// convert(fin, fout)
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	RunDir(fmt.Sprintf("%s/%s", home, "Downloads"))
}

func RunDir(dir string) {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to Read dir: %v\n", err)
	}

	for _, f := range files {
		heic := fmt.Sprintf("%s/%s", dir, f.Name())
		jpg := fmt.Sprintf("%s.jpg", heic)

		if f.IsDir() {
			RunDir(heic)
			continue
		}
		if isWrongHeic(heic) {
			rename(heic, jpg)
			continue
		}

		if isCorrectHeic(heic) {
			convert(heic, jpg)
			modified(heic, jpg)
			remove(heic)
			continue
		}
	}
}

//@todo: test.
func isWrongHeic(fin string) bool {
	if !hasEnding(fin) {
		return false
	}
	mime := getFileContentType(fin)
	return mime == "image/jpeg"
}

//@todo: test.
func isCorrectHeic(fin string) bool {
	if !hasEnding(fin) {
		return false
	}
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

func rename(fin string, fout string) {
	err := os.Rename(fin, fout)
	if err != nil {
		log.Printf("Warning: Cannot rename from %s to %s: %v\n", fin, fout, err)
	}

	log.Printf("Renamed %s to %s successfully\n", fin, fout)
}

func modified(fin string, fout string) {
	file, err := os.Stat(fin)

	if err != nil {
		log.Printf("Warning: Cannot read file to set original date: %s: %v\n", fin, err)
	}

	modifiedtime := file.ModTime()

	err = os.Chtimes(fout, modifiedtime, modifiedtime)

	if err != nil {
		log.Printf("Warning: Cannot change time of %s to %s: %v\n", fin, modifiedtime, err)
	}

	log.Printf("Modified %s times to: %s", fout, modifiedtime)
}

func remove(fin string) {
	err := os.Remove(fin)
	if err != nil {
		log.Printf("Warning: Cannot remove %s: %v\n", fin, err)
	}

	log.Printf("Removed %s successfully\n", fin)
}

func convert(fin string, fout string) {
	fi, err := os.Open(fin)
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()

	exif, err := goheif.ExtractExif(fi)
	if err != nil {
		log.Printf("Warning: no EXIF from %s: %v\n", fin, err)
	}

	img, err := goheif.Decode(fi)
	if err != nil {
		log.Printf("Failed to parse %s: %v\n", fin, err)
		return

	}
	fo, err := os.OpenFile(fout, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Failed to create output file %s: %v\n", fout, err)
		return
	}
	defer fo.Close()

	w, _ := newWriterExif(fo, exif)
	err = jpeg.Encode(w, img, nil)

	if err != nil {
		log.Printf("Failed to encode %s: %v\n", fout, err)
	}

	log.Printf("Converted %s to %s successfully\n", fin, fout)
}
