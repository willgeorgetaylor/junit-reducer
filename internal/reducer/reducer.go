package reducer

import (
	"os"

	"github.com/bmatcuk/doublestar"
	"github.com/willgeorgetaylor/junit-reducer/internal/enums"
	"github.com/willgeorgetaylor/junit-reducer/internal/helpers"
	"github.com/willgeorgetaylor/junit-reducer/internal/serialization"
)

type ReduceFunctionParams struct {
	IncludeFilePattern           string
	ExcludeFilePattern           string
	OutputPath                   string
	ReduceTestSuitesBy           enums.TestSuiteField
	ReduceTestCasesBy            enums.TestCaseField
	OperatorTestSuitesTests      enums.AggregateOperation
	OperatorTestSuitesFailed     enums.AggregateOperation
	OperatorTestSuitesErrors     enums.AggregateOperation
	OperatorTestSuitesSkipped    enums.AggregateOperation
	OperatorTestSuitesAssertions enums.AggregateOperation
	OperatorTestSuitesTime       enums.AggregateOperation
	OperatorTestCasesTime        enums.AggregateOperation
	RoundingMode                 enums.RoundingMode
}

func Reduce(params ReduceFunctionParams) {
	files := make(map[string]bool)
	includedReports, err := doublestar.Glob(params.IncludeFilePattern)

	if err != nil {
		helpers.FatalMsg("failed to enumerate included JUnit XML reports in: %v", err)
		os.Exit(1)
	}
	for _, file := range includedReports {
		files[file] = true
	}

	// Exclude files (optional)
	if params.ExcludeFilePattern != "" {
		excludedFiles, err := doublestar.Glob(params.ExcludeFilePattern)

		if err != nil {
			helpers.FatalMsg("failed to enumerate excluded JUnit XML reports in: %v", err)
			os.Exit(1)
		}
		for _, file := range excludedFiles {
			delete(files, file)
		}
	}

	// Get paths to included files
	filesSlice := make([]string, 0, len(files))

	for file := range files {
		filesSlice = append(filesSlice, file)
	}

	// Order alphabetically
	helpers.SortStrings(filesSlice)

	// Deserialize
	testSuites, err := serialization.Deserialize(filesSlice)

	if err != nil {
		helpers.FatalMsg("failed to deserialize JUnit XML reports: %v", err)
		os.Exit(1)
	}

	// For now, just reduce testsuites by filepath, and average time values
	// TODO: Add support for other flags (reduceTestCasesBy, operatorTestSuitesTests, etc.)

	testSuiteMap := make(map[string][]serialization.TestSuite)

	for _, testSuite := range testSuites {
		// Key off both file and name
		combinedKey := testSuite.File + ":" + testSuite.Name
		testSuiteMap[combinedKey] = append(testSuiteMap[combinedKey], testSuite)
	}

	// Reduce times and other aggregate fields
	for key, testSuiteSlice := range testSuiteMap {
		testSuiteMap[key] = reduceTestSuiteSlice(testSuiteSlice)
	}

	// Flatten back to a set of test suites
	testSuites = make([]serialization.TestSuite, 0, len(testSuiteMap))
	for _, testSuiteSlice := range testSuiteMap {
		testSuites = append(testSuites, testSuiteSlice...)
	}

	// Create output directory if it doesn't exist
	err = os.MkdirAll(params.OutputPath, os.ModePerm)
	if err != nil {
		helpers.FatalMsg("failed to create output directory: %v", err)
		os.Exit(1)
	}

	serialization.Serialize(params.OutputPath, testSuites)
}

func reduceTestSuiteSlice(testSuiteSlice []serialization.TestSuite) []serialization.TestSuite {
	testSuite := testSuiteSlice[0]
	var totalTime float64 = 0
	for _, testSuite := range testSuiteSlice {
		totalTime += testSuite.Time
	}
	testSuite.Time = totalTime / float64(len(testSuiteSlice))
	return []serialization.TestSuite{testSuite}
}
