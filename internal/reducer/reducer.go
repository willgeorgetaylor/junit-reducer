package reducer

import (
	"math"
	"os"
	"sort"

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

func Reduce(params ReduceFunctionParams) error {
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
		reducedTestSlice, err := reduceTestSuiteSlice(testSuiteSlice, params)
		if err != nil {
			return err
		}
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
		os.Exit(1)
	}

	serialization.Serialize(params.OutputPath, testSuites)

	return nil
}

type SuiteFieldExtractor func(*serialization.TestSuite) float64

func SuiteTimeExtractor(ts *serialization.TestSuite) float64 {
	return ts.Time
}

func SuiteTestsExtractor(ts *serialization.TestSuite) float64 {
	return float64(ts.Tests)
}

func SuiteFailedExtractor(ts *serialization.TestSuite) float64 {
	return float64(ts.Failed)
}

func SuiteErrorsExtractor(ts *serialization.TestSuite) float64 {
	return float64(ts.Errors)
}

func SuiteSkippedExtractor(ts *serialization.TestSuite) float64 {
	return float64(ts.Skipped)
}

func SuiteAssertionsExtractor(ts *serialization.TestSuite) float64 {
	return float64(ts.Assertions)
}

func reduceTestSuiteSlice(testSuiteSlice []serialization.TestSuite, params ReduceFunctionParams) ([]serialization.TestSuite, error) {
	testSuite := testSuiteSlice[0]

	reducedTime, err := reduceTestSuites(testSuiteSlice, SuiteTimeExtractor, params.OperatorTestSuitesTime)
	if err != nil {
		return nil, err
	} else {
		testSuite.Time = reducedTime
	}

	reducedTests, err := reduceTestSuites(testSuiteSlice, SuiteTestsExtractor, params.OperatorTestSuitesTests)
	if err != nil {
		return nil, err
	} else {
		testSuite.Tests = roundToInt(reducedTests, params.RoundingMode)
	}

	reducedFailed, err := reduceTestSuites(testSuiteSlice, SuiteFailedExtractor, params.OperatorTestSuitesFailed)
	if err != nil {
		return nil, err
	} else {
		testSuite.Failed = roundToInt(reducedFailed, params.RoundingMode)
	}

	reducedErrors, err := reduceTestSuites(testSuiteSlice, SuiteErrorsExtractor, params.OperatorTestSuitesErrors)
	if err != nil {
		return nil, err
	} else {
		testSuite.Errors = roundToInt(reducedErrors, params.RoundingMode)
	}

	reducedSkipped, err := reduceTestSuites(testSuiteSlice, SuiteSkippedExtractor, params.OperatorTestSuitesSkipped)
	if err != nil {
		return nil, err
	} else {
		testSuite.Skipped = roundToInt(reducedSkipped, params.RoundingMode)
	}

	reducedAssertions, err := reduceTestSuites(testSuiteSlice, SuiteAssertionsExtractor, params.OperatorTestSuitesAssertions)
	if err != nil {
		return nil, err
	} else {
		testSuite.Assertions = roundToInt(reducedAssertions, params.RoundingMode)
	}

	return []serialization.TestSuite{testSuite}, nil
}

func reduceTestSuites(testSuiteSlice []serialization.TestSuite, extractor SuiteFieldExtractor, operation enums.AggregateOperation) (float64, error) {
	if operation == enums.AggregateOperationMean {
		return reduceByMean(testSuiteSlice, extractor), nil
	} else if operation == enums.AggregateOperationMax {
		return reduceByMax(testSuiteSlice, extractor), nil
	} else if operation == enums.AgregateOperationMin {
		return reduceByMin(testSuiteSlice, extractor), nil
	} else if operation == enums.AggregateOperationMode {
		return reduceByMode(testSuiteSlice, extractor), nil
	} else if operation == enums.AggregateOperationSum {
		return reduceBySum(testSuiteSlice, extractor), nil
	} else if operation == enums.AggregateOperationMedian {
		return reduceByMedian(testSuiteSlice, extractor), nil
	} else {
		return 0, nil
	}
}

func reduceByMax(testSuiteSlice []serialization.TestSuite, extractor SuiteFieldExtractor) float64 {
	var max float64 = 0
	for _, testSuite := range testSuiteSlice {
		max = math.Max(extractor(&testSuite), max)
	}
	return max
}

func reduceByMean(testSuiteSlice []serialization.TestSuite, extractor SuiteFieldExtractor) float64 {
	var total float64 = 0
	for _, testSuite := range testSuiteSlice {
		total += extractor(&testSuite)
	}
	mean := total / float64(len(testSuiteSlice))
	return mean
}

func reduceByMin(testSuiteSlice []serialization.TestSuite, extractor SuiteFieldExtractor) float64 {
	var min float64 = 0
	for _, testSuite := range testSuiteSlice {
		min = math.Min(extractor(&testSuite), min)
	}
	return min
}

func reduceByMode(testSuiteSlice []serialization.TestSuite, extractor SuiteFieldExtractor) float64 {
	freqs := make(map[float64]int)
	for _, testSuite := range testSuiteSlice {
		val := extractor(&testSuite)
		freqs[val] += 1
	}
	var topVal float64 = 0.0
	var topCount int = 0.0
	for val, count := range freqs {
		if count >= topCount {
			topCount = count
			topVal = val
		}
	}
	return topVal
}

func reduceBySum(testSuiteSlice []serialization.TestSuite, extractor SuiteFieldExtractor) float64 {
	var total float64 = 0
	for _, testSuite := range testSuiteSlice {
		total += extractor(&testSuite)
	}
	return total
}

func reduceByMedian(testSuiteSlice []serialization.TestSuite, extractor SuiteFieldExtractor) float64 {
	vals := make([]float64, 0, len(testSuiteSlice))
	for _, testSuite := range testSuiteSlice {
		vals = append(vals, extractor(&testSuite))
	}
	sort.Slice(vals, func(i, j int) bool {
		return vals[i] < vals[j]
	})
	medianIndex := medianIndex(len(vals))
	return vals[medianIndex]
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
