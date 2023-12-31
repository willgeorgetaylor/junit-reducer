package reducer

import (
	"errors"
	"math"
	"os"
	"sort"

	"github.com/bmatcuk/doublestar"
	"github.com/willgeorgetaylor/junit-reducer/internal/enums"
	"github.com/willgeorgetaylor/junit-reducer/internal/helpers"
	"github.com/willgeorgetaylor/junit-reducer/internal/serialization"
)

type ReduceFunctionParams struct {
	IncludeFilePattern            string
	ExcludeFilePattern            string
	OutputPath                    string
	ReduceTestSuitesBy            enums.TestSuiteField
	ReduceTestCasesBy             enums.TestCaseField
	OperationTestSuitesTests      enums.AggregateOperation
	OperationTestSuitesFailed     enums.AggregateOperation
	OperationTestSuitesErrors     enums.AggregateOperation
	OperationTestSuitesSkipped    enums.AggregateOperation
	OperationTestSuitesAssertions enums.AggregateOperation
	OperationTestSuitesTime       enums.AggregateOperation
	OperationTestCasesTime        enums.AggregateOperation
	RoundingMode                  enums.RoundingMode
}

func Reduce(params ReduceFunctionParams) error {
	files := make(map[string]bool)
	includedReports, err := doublestar.Glob(params.IncludeFilePattern)

	if err != nil {
		helpers.FatalMsg("failed to enumerate included JUnit XML reports: %v", err)
		return err
	}
	for _, file := range includedReports {
		files[file] = true
	}

	// Exclude files
	if params.ExcludeFilePattern != "" {
		excludedFiles, err := doublestar.Glob(params.ExcludeFilePattern)

		if err != nil {
			helpers.FatalMsg("failed to enumerate excluded JUnit XML reports: %v", err)
			return err
		}
		for _, file := range excludedFiles {
			helpers.PrintMsg("excluding file: %v", file)
			delete(files, file)
		}
	}

	if (len(files)) == 0 {
		return errors.New("no files matched the provided include pattern")
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
		return err
	}

	// For now, just reduce testsuites by filepath, and average time values
	// TODO: Add support for other flags (reduceTestCasesBy, operationTestSuitesTests, etc.)

	testSuiteMap := make(map[string][]serialization.TestSuite)

	for _, testSuite := range testSuites {
		suiteKey := testSuite.Name
		if params.ReduceTestSuitesBy == enums.TestSuiteFieldNameFilepath {
			suiteKey = testSuite.File + ":" + testSuite.Name
		} else if params.ReduceTestSuitesBy == enums.TestSuiteFieldFilepath {
			suiteKey = testSuite.File
		}
		testSuiteMap[suiteKey] = append(testSuiteMap[suiteKey], testSuite)
	}

	// Reduce times and other aggregate fields
	for key, testSuiteSlice := range testSuiteMap {
		reducedTestSlice := reduceTestSuiteSlice(testSuiteSlice, params)
		testSuiteMap[key] = reducedTestSlice
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
		return err
	}

	serialization.Serialize(params.OutputPath, testSuites)

	return nil
}

type SuiteFieldExtractor func(serialization.TestSuite) float64

func SuiteTimeExtractor(ts serialization.TestSuite) float64 {
	return ts.Time
}

func SuiteTestsExtractor(ts serialization.TestSuite) float64 {
	return float64(ts.Tests)
}

func SuiteFailedExtractor(ts serialization.TestSuite) float64 {
	return float64(ts.Failed)
}

func SuiteErrorsExtractor(ts serialization.TestSuite) float64 {
	return float64(ts.Errors)
}

func SuiteSkippedExtractor(ts serialization.TestSuite) float64 {
	return float64(ts.Skipped)
}

func SuiteAssertionsExtractor(ts serialization.TestSuite) float64 {
	return float64(ts.Assertions)
}

func reduceTestSuiteSlice(testSuiteSlice []serialization.TestSuite, params ReduceFunctionParams) []serialization.TestSuite {
	testSuite := testSuiteSlice[0]

	testSuite.Time = reduceTestSuites(testSuiteSlice, SuiteTimeExtractor, params.OperationTestSuitesTime)

	// Tests count
	reducedTests := reduceTestSuites(testSuiteSlice, SuiteTestsExtractor, params.OperationTestSuitesTests)
	testSuite.Tests = roundToInt(reducedTests, params.RoundingMode)

	// Failed count
	reducedFailed := reduceTestSuites(testSuiteSlice, SuiteFailedExtractor, params.OperationTestSuitesFailed)
	testSuite.Failed = roundToInt(reducedFailed, params.RoundingMode)

	// Errors count
	reducedErrors := reduceTestSuites(testSuiteSlice, SuiteErrorsExtractor, params.OperationTestSuitesErrors)
	testSuite.Errors = roundToInt(reducedErrors, params.RoundingMode)

	// Skipped count
	reducedSkipped := reduceTestSuites(testSuiteSlice, SuiteSkippedExtractor, params.OperationTestSuitesSkipped)
	testSuite.Skipped = roundToInt(reducedSkipped, params.RoundingMode)

	// Assertions count
	reducedAssertions := reduceTestSuites(testSuiteSlice, SuiteAssertionsExtractor, params.OperationTestSuitesAssertions)
	testSuite.Assertions = roundToInt(reducedAssertions, params.RoundingMode)

	// Cases
	testSuite.TestCases = reduceTestCases(testSuiteSlice, params.ReduceTestCasesBy, params.OperationTestCasesTime)

	return []serialization.TestSuite{testSuite}
}

func reduceTestCases(testSuiteSlice []serialization.TestSuite, reduceBy enums.TestCaseField, operation enums.AggregateOperation) []serialization.TestCase {
	groupedCases := make(map[string][]serialization.TestCase)

	for _, testSuite := range testSuiteSlice {
		for _, testCase := range testSuite.TestCases {
			key := extractKeyFromCase(testCase, reduceBy)
			groupedCases[key] = append(groupedCases[key], testCase)
		}
	}

	reducedCases := make([]serialization.TestCase, 0, len(groupedCases))

	for _, cases := range groupedCases {
		baseCase := cases[0]
		reducedTime := reduceTestCaseTimes(cases, operation)
		baseCase.Time = reducedTime
		reducedCases = append(reducedCases, baseCase)
	}

	return reducedCases
}

func extractKeyFromCase(testCase serialization.TestCase, reduceBy enums.TestCaseField) string {
	if reduceBy == enums.TestCaseFieldClassname {
		return testCase.Classname
	} else if reduceBy == enums.TestCaseFieldFile {
		return testCase.File
	} else {
		return testCase.Name
	}
}

func reduceTestCaseTimes(testCaseSlice []serialization.TestCase, operation enums.AggregateOperation) float64 {
	slice := make([]float64, 0, len(testCaseSlice))
	for _, testCase := range testCaseSlice {
		slice = append(slice, testCase.Time)
	}
	return reduce(slice, operation)
}

func reduceTestSuites(testSuiteSlice []serialization.TestSuite, extractor SuiteFieldExtractor, operation enums.AggregateOperation) float64 {
	slice := make([]float64, 0, len(testSuiteSlice))
	for _, testSuite := range testSuiteSlice {
		slice = append(slice, extractor(testSuite))
	}
	return reduce(slice, operation)
}

func reduce(slice []float64, operation enums.AggregateOperation) float64 {
	if operation == enums.AggregateOperationMax {
		return reduceMax(slice)
	} else if operation == enums.AggregateOperationMin {
		return reduceMin(slice)
	} else if operation == enums.AggregateOperationMode {
		return reduceMode(slice)
	} else if operation == enums.AggregateOperationSum {
		return reduceSum(slice)
	} else if operation == enums.AggregateOperationMedian {
		return reduceMedian(slice)
	} else {
		return reduceMean(slice)
	}
}

func reduceMax(slice []float64) float64 {
	var max float64 = 0
	for _, val := range slice {
		max = math.Max(val, max)
	}
	return max
}

func reduceMin(slice []float64) float64 {
	var min float64 = slice[0]
	for _, val := range slice {
		min = math.Min(val, min)
	}
	return min
}

func reduceMean(slice []float64) float64 {
	var total float64 = 0
	for _, val := range slice {
		total += val
	}
	mean := total / float64(len(slice))
	return mean
}

func reduceMode(slice []float64) float64 {
	freqs := make(map[float64]int)
	for _, val := range slice {
		freqs[val]++
	}
	var topVal float64 = 0
	var topFreq int = 0
	for val, freq := range freqs {
		if freq > topFreq {
			topVal = val
			topFreq = freq
		}
	}
	return topVal
}

func reduceSum(slice []float64) float64 {
	var total float64 = 0
	for _, val := range slice {
		total += val
	}
	return total
}

func reduceMedian(slice []float64) float64 {
	sortedSlice := make([]float64, len(slice))
	copy(sortedSlice, slice)
	sort.Float64s(sortedSlice)
	medianIndex := medianIndex(len(sortedSlice))
	return sortedSlice[medianIndex]
}

func medianIndex(sliceLength int) int {
	if sliceLength <= 2 {
		return 0
	} else {
		return int((sliceLength - 1) / 2)
	}
}

func roundToInt(value float64, roundingMode enums.RoundingMode) int {
	if roundingMode == enums.RoundingModeCeil {
		return int(math.Ceil(value))
	} else if roundingMode == enums.RoundingModeFloor {
		return int(math.Floor(value))
	} else {
		return int(math.Round(value))
	}
}
