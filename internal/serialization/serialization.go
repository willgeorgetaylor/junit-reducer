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

func UnmarshalTestSuites(xmlData []byte, fileName string) (*TestSuites, error) {
	var testSuites TestSuites
	err := xml.Unmarshal(xmlData, &testSuites)
	if err != nil {
		return nil, err
	}
	for i := range testSuites.TestSuites {
		testSuites.TestSuites[i].FileName = fileName
	}
	return &testSuites, nil
}

func DeserializeFromReader(
	testSuites []TestSuite,
	reader io.Reader,
	fileName string,
) ([]TestSuite, error) {
	xmlData, err := io.ReadAll(reader)
	if err != nil {
		helpers.FatalMsg("failed to read junit xml: %v\n", err)
		return nil, err
	}
	xmlTestSuites, err := UnmarshalTestSuites(xmlData, fileName)
	if err != nil {
		helpers.FatalMsg("failed to parse junit xml: %v\n", err)
		return nil, err
	}
	testSuites = append(testSuites, xmlTestSuites.TestSuites...)
	return testSuites, nil
}

func Deserialize(
	junitFilePaths []string,
) ([]TestSuite, error) {
	var testSuites []TestSuite
	for _, junitFilePath := range junitFilePaths {
		file, err := os.Open(junitFilePath)
		fileName := filepath.Base(junitFilePath)
		if err != nil {
			helpers.FatalMsg("failed to open junit xml: %v\n", err)
			return nil, err
		}

		helpers.PrintMsg("deserializing junit xml: %v\n", junitFilePath)

		testSuites, err = DeserializeFromReader(testSuites, file, fileName)
		file.Close()
		if err != nil {
			helpers.FatalMsg("failed to deserialize junit xml: %v\n", err)
			return nil, err
		}
	}
	return testSuites, nil
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
