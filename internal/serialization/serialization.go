package serialization

import (
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/willgeorgetaylor/junit-reducer/internal/helpers"
)

type TestSuites struct {
	TestSuites []TestSuite `xml:"testsuite"`
}

type TestSuite struct {
	Name string `xml:"name,attr"`
	File string `xml:"filepath,attr"`
	// Aggregated fields
	Time       float64    `xml:"time,attr"`
	Tests      int        `xml:"tests,attr"`
	Failed     int        `xml:"failed,attr"`
	Errors     int        `xml:"errors,attr"`
	Skipped    int        `xml:"skipped,attr"`
	Assertions int        `xml:"assertions,attr"`
	TestCases  []TestCase `xml:"testcase"`
	// For reserialization
	FileName string `xml:"-"`
}

type TestCase struct {
	Name       string `xml:"name,attr"`
	Classname  string `xml:"classname,attr"`
	File       string `xml:"file,attr"`
	Line       int    `xml:"lineno,attr"`
	Assertions int    `xml:"assertions,attr"`
	// Aggregated fields
	Time float64 `xml:"time,attr"`
}

func unmarshalTestSuites(xmlData []byte, fileName string) *TestSuites {
	var testSuites TestSuites
	err := xml.Unmarshal(xmlData, &testSuites)
	if err != nil {
		helpers.FatalMsg("failed to unmarshal junit xml: %v\n", err)
	}
	for i := range testSuites.TestSuites {
		testSuites.TestSuites[i].FileName = fileName
	}
	return &testSuites
}

func deserializeFromReader(
	testSuites []TestSuite,
	reader io.Reader,
	fileName string,
) []TestSuite {
	xmlData, err := io.ReadAll(reader)
	if err != nil {
		helpers.FatalMsg("failed to read junit xml: %v\n", err)
	}
	xmlTestSuites := unmarshalTestSuites(xmlData, fileName)
	testSuites = append(testSuites, xmlTestSuites.TestSuites...)
	return testSuites
}

func Deserialize(
	junitFilePaths []string,
) []TestSuite {
	var testSuites []TestSuite
	for _, junitFilePath := range junitFilePaths {
		file, err := os.Open(junitFilePath)
		fileName := filepath.Base(junitFilePath)

		if err != nil {
			helpers.FatalMsg("failed to open junit xml: %v\n", err)
		}
		helpers.PrintMsg("deserializing junit xml: %v\n", junitFilePath)
		testSuites = deserializeFromReader(testSuites, file, fileName)
		file.Close()
	}
	return testSuites
}

func Serialize(outputPath string, testSuites []TestSuite) {
	testSuiteMap := make(map[string][]TestSuite)

	for _, testSuite := range testSuites {
		helpers.PrintMsg("serializing junit xml: %v\n", testSuite.FileName)
		testSuiteMap[testSuite.FileName] = append(testSuiteMap[testSuite.FileName], testSuite)
	}

	for fileName, suites := range testSuiteMap {
		testSuitesWrapper := TestSuites{TestSuites: suites}
		outputFileName := filepath.Join(outputPath, fileName)

		// Marshal to XML
		xmlBytes, err := xml.MarshalIndent(testSuitesWrapper, "", "  ")
		if err != nil {
			helpers.FatalMsg("failed to marshal junit xml: %v\n", err)
		}

		// Convert XML bytes to string for replacement
		xmlString := string(xmlBytes)

		// Marshalling XML requires the wrapper type be title-cased to be exportable
		// in Go, but we want to preserve the original casing for the XML tags.
		xmlString = strings.Replace(xmlString, "<TestSuites>", "<testsuites>", 1)
		xmlString = strings.Replace(xmlString, "</TestSuites>", "</testsuites>", 1)
		// Add XML header
		xmlString = xml.Header + xmlString

		// Write to file
		xmlBytes = []byte(xmlString)
		err = os.WriteFile(outputFileName, xmlBytes, 0644)
		if err != nil {
			helpers.FatalMsg("failed to write junit xml to file: %v\n", err)
		}
	}
}
