package scanutils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var RiskLevelError = errors.New("Risk level does not exist")

func RiskRegex(r string) (string, error) {
	var regex string
	var filename string

	// convert to env variable to path
	switch {
	case r == "high":
		filename = os.Getenv("GOPATH") + "/tmp/risk_phrases/high_risk_phrases.txt"
	case r == "low":
		filename = os.Getenv("GOPATH") + "/tmp/risk_phrases/low_risk_phrases.txt"
	default:
		return regex, RiskLevelError
	}

	phr, errOpen := os.Open(filename)
	defer phr.Close()
	if errOpen != nil {
		return regex, errOpen
	}

	regex = "(?i)"

	scanner := bufio.NewScanner(phr)
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		if atEOF == true {
			regex = strings.TrimSuffix(regex, "|")
		}
		return
	}
	scanner.Split(split)

	for scanner.Scan() {
		r_phr := strings.ToLower(scanner.Text())
		regex = regex + r_phr + "|"
	}
	if errScanner := scanner.Err(); errScanner != nil {
		return regex, errScanner
	}

	return regex, nil
}

func ScorePhrase(phr string, hrp string, lrp string) int {
	var s int

	hrRegex := regexp.MustCompile(hrp)
	hrRegexMatch := hrRegex.FindAllString(phr, -1)
	s += len(hrRegexMatch) * 2

	lrRegex := regexp.MustCompile(lrp)
	lrRegexMatch := lrRegex.FindAllString(phr, -1)
	s += len(lrRegexMatch)

	return s
}

func ScoreFile(path string, hrRgx string, lrRgx string) (int, error) {
	var score int

	input, errOpenFile := os.Open(path)
	defer input.Close()
	if errOpenFile != nil {
		return score, errOpenFile
	}

	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		phr := strings.ToLower(scanner.Text())
		score += ScorePhrase(phr, hrRgx, lrRgx)
	}
	if errScanner := scanner.Err(); errScanner != nil {
		return score, errScanner
	}

	return score, nil
}

func WriteFile(inputDir string, outFile *os.File, hrp string, lrp string) error {
	w := bufio.NewWriter(outFile)

	errFilepath := filepath.Walk(inputDir, func(path string, info os.FileInfo, errWalkfn error) error {
		if errWalkfn != nil {
			return errWalkfn
		}

		if !info.IsDir() {
			score, errScore := ScoreFile(path, hrp, lrp)
			if errScore != nil {
				return errScore
			}

			line := fmt.Sprintf("%v:%d\n", path, score)
			w.WriteString(line)
		}

		return nil
	})

	if errFilepath == nil {
		return w.Flush()
	} else {
		return errFilepath
	}
}

func RunScanFile() (error, string) {
	hrp, _ := RiskRegex("high")
	lrp, _ := RiskRegex("low")
	inputDir := os.Getenv("GOPATH") + "/tmp/input_files"
	outFilename := os.Getenv("GOPATH") + "/tmp/output_files/Text Scan " + time.Now().Format(time.RFC3339) + ".txt"

	outFile, errCreateFile := os.Create(outFilename)
	if errCreateFile != nil {
		return errCreateFile, outFilename
	}

	errWriteFile := WriteFile(inputDir, outFile, hrp, lrp)

	if errCloseFile := outFile.Close(); errCloseFile != nil {
		return errCloseFile, outFilename
	}

	if errWriteFile != nil {
		return errWriteFile, outFilename
	} else {
		return nil, outFilename
	}
}
