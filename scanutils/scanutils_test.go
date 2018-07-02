package scanutils

import (
	"os"
	"testing"
	"time"
)

var expectedHrp = "(?i)voldemort|dark lord|mundane|pinocchio"
var expectedLrp = "(?i)esolutions|gangster|ugliest|destiny|shooter|plan"

func TestRiskRegex(t *testing.T) {
	var testValues = []struct {
		r   string
		out string
	}{
		{"high", expectedHrp},
		{"low", expectedLrp},
		{"fake", ""},
	}

	for _, value := range testValues {
		out, _ := RiskRegex(value.r)
		if out != value.out {
			t.Errorf("RiskRegex(%q) = %v, want %v", value.r, out, value.out)
		}
	}
}

func TestScorePhrase(t *testing.T) {
	var testValues = []struct {
		phr string
		out int
	}{
		{"Single hp match, Dark Lord", 2},
		{"Single lp match, shOOter", 1},
		{"Mix and match, mundane gangster", 3},
		{"Repeat words matched, gangster gangster", 2},
		{"Repeat words matched, dark lord mundane dark lord", 6},
		{"No matches on this one", 0},
	}

	for _, value := range testValues {
		out := ScorePhrase(value.phr, expectedHrp, expectedLrp)
		if out != value.out {
			t.Errorf("ScorePhrase(%v) = %v, want %v", value.phr, out, value.out)
		}
	}
}

func TestScoreFile(t *testing.T) {
	var testValues = []struct {
		path string
		out  int
	}{
		{"testdata/input/testinput01.txt", 3},
		{"testdata/input/testinput02.txt", 4},
		{"testdata/input/testinput03.txt", 1},
	}

	for _, value := range testValues {
		out, _ := ScoreFile(value.path, expectedHrp, expectedLrp)
		if out != value.out {
			t.Errorf("ScoreFile(%v) = %v, want %v", value.path, out, value.out)
		}
	}
}

func TestWriteFile(t *testing.T) {
	outFilename := "testdata/output/Test Text Scan " + time.Now().Format(time.RFC3339)
	outFile, _ := os.Create(outFilename)

	err := WriteFile("testdata/input", outFile, expectedHrp, expectedLrp)
	if err != nil {
		t.Errorf("WriteFile Error: %v", err)
	}
}

func TestRunScanFile(t *testing.T) {
	err, _ := RunScanFile()
	if err != nil {
		t.Errorf("RunScanFile() error: %v", err)
	}
}
