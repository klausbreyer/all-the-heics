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

func HomeDir() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	logs := RunDir(fmt.Sprintf("%s/%s", home, "Downloads"))

	log.Printf("%+v", logs)
	return logs
}

func RunDir(dir string) []string {
	logs := []string{}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		logs = append(logs, fmt.Sprintf("Failed to Read dir: %v", err))
	}

	for _, f := range files {
		heic := fmt.Sprintf("%s/%s", dir, f.Name())
		jpg := fmt.Sprintf("%s.jpg", heic)

		if f.IsDir() {
			logs = append(logs, RunDir(heic)...)
			continue
		}

		if isWrongHeic(heic) {
			logs = append(logs, rename(heic, jpg)...)
			continue
		}

		if isCorrectHeic(heic) {
			logs = append(logs, convert(heic, jpg)...)
			logs = append(logs, modified(heic, jpg)...)
			// logs = append(logs, remove(heic)...)
			continue
		}
	}

	return logs
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

func rename(fin string, fout string) []string {
	err := os.Rename(fin, fout)
	if err != nil {
		return []string{fmt.Sprintf("Warning: Cannot rename from %s to %s: %v", fin, fout, err)}

	}

	return []string{fmt.Sprintf("Renamed %s to %s successfully", fin, fout)}
}

func modified(fin string, fout string) []string {
	file, err := os.Stat(fin)

	if err != nil {
		return []string{fmt.Sprintf("Warning: Cannot read file to set original date: %s: %v", fin, err)}

	}

	modifiedtime := file.ModTime()
	err = os.Chtimes(fout, modifiedtime, modifiedtime)

	if err != nil {
		return []string{fmt.Sprintf("Warning: Cannot change time of %s to %s: %v", fin, modifiedtime, err)}

	}

	return []string{fmt.Sprintf("Modified %s times to: %s", fout, modifiedtime)}
}

func remove(fin string) []string {
	err := os.Remove(fin)
	if err != nil {
		return []string{fmt.Sprintf("Warning: Cannot remove %s: %v", fin, err)}
	}

	return []string{fmt.Sprintf("Removed %s successfully", fin)}
}

func convert(fin string, fout string) []string {

	logs := []string{}
	fi, err := os.Open(fin)
	if err != nil {
		logs = append(logs, fmt.Sprintf("Error: failed to open %s: %v", fin, err))
	}
	defer fi.Close()

	exif, err := goheif.ExtractExif(fi)
	if err != nil {
		logs = append(logs, fmt.Sprintf("Warning: no EXIF from %s: %v", fin, err))
	}

	img, err := goheif.Decode(fi)
	if err != nil {
		logs = append(logs, fmt.Sprintf("Failed to parse %s: %v", fin, err))
		return logs

	}
	fo, err := os.OpenFile(fout, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logs = append(logs, fmt.Sprintf("Failed to create output file %s: %v", fout, err))
		return logs
	}
	defer fo.Close()

	w, _ := newWriterExif(fo, exif)
	err = jpeg.Encode(w, img, nil)

	if err != nil {
		logs = append(logs, fmt.Sprintf("Failed to encode %s: %v", fout, err))
	}

	logs = append(logs, fmt.Sprintf("Converted %s to %s successfully", fin, fout))

	return logs
}
