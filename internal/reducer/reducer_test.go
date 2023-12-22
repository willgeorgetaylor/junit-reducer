package reducer

import (
	"io"
	"os"
	"testing"

	"github.com/willgeorgetaylor/junit-reducer/internal/enums"
	"github.com/willgeorgetaylor/junit-reducer/internal/helpers"
	"github.com/willgeorgetaylor/junit-reducer/internal/serialization"
)

func clearOutputDir() {
	os.RemoveAll("output")
}

func setup() {
	clearOutputDir()
}

func tearDown() {
	clearOutputDir()
}

func assertTestFile(t *testing.T, testFile serialization.TestSuites) {
	junitReportFileName := testFile.TestSuites[0].FileName
	outputFileName := "./output/" + junitReportFileName

	if !helpers.DirExists("./output") {
		t.Errorf("expected 'output' directory to exist")
	}

	if !helpers.FileExists(outputFileName) {
		t.Errorf("expected output file '%s' to exist", outputFileName)
	}

	file, err := os.Open(outputFileName)

	if err != nil {
		t.Errorf("error opening output file '%s'", outputFileName)
	}

	xmlData, err := io.ReadAll(file)

	if err != nil {
		t.Errorf("error reading data from output file '%s'", junitReportFileName)
	}

	xmlTestSuites, err := serialization.UnmarshalTestSuites(xmlData, junitReportFileName)

	if err != nil {
		t.Errorf("error parsing JUnit XML from output file '%s'", junitReportFileName)
	}

	if len(xmlTestSuites.TestSuites) == 0 {
		t.Errorf("expected output file 'output/%s' to have at least one test suite", junitReportFileName)
	}

	if len(xmlTestSuites.TestSuites) != len(testFile.TestSuites) {
		t.Errorf("expected output file 'output/%s' to have %d test suites", junitReportFileName, len(testFile.TestSuites))
	}

	for i := 0; i < len(testFile.TestSuites); i++ {
		if xmlTestSuites.TestSuites[i].Name != testFile.TestSuites[i].Name {
			t.Errorf("expected test suite with name '%s' to have name '%s'", xmlTestSuites.TestSuites[i].Name, testFile.TestSuites[i].Name)
		}

		if xmlTestSuites.TestSuites[i].File != testFile.TestSuites[i].File {
			t.Errorf("expected test suite for test file '%s' to be for test file '%s'", xmlTestSuites.TestSuites[i].File, testFile.TestSuites[i].File)
		}

		if xmlTestSuites.TestSuites[i].FileName != testFile.TestSuites[i].FileName {
			t.Errorf("expected test report with file name '%s' to have file name '%s'", xmlTestSuites.TestSuites[i].FileName, testFile.TestSuites[i].FileName)
		}

		if xmlTestSuites.TestSuites[i].Time != testFile.TestSuites[i].Time {
			t.Errorf("expected test suite for file '%s' to report time of %f seconds", xmlTestSuites.TestSuites[i].File, testFile.TestSuites[i].Time)
		}

		if xmlTestSuites.TestSuites[i].Tests != testFile.TestSuites[i].Tests {
			t.Errorf("expected test suite for file '%s' to report %d tests", xmlTestSuites.TestSuites[i].File, testFile.TestSuites[i].Tests)
		}

		if xmlTestSuites.TestSuites[i].Failed != testFile.TestSuites[i].Failed {
			t.Errorf("expected test suite for file '%s' to report %d failures", xmlTestSuites.TestSuites[i].File, testFile.TestSuites[i].Failed)
		}

		if xmlTestSuites.TestSuites[i].Errors != testFile.TestSuites[i].Errors {
			t.Errorf("expected test suite for file '%s' to report %d errors", xmlTestSuites.TestSuites[i].File, testFile.TestSuites[i].Errors)
		}

		if xmlTestSuites.TestSuites[i].Skipped != testFile.TestSuites[i].Skipped {
			t.Errorf("expected test suite for file '%s' to report %d skipped tests", xmlTestSuites.TestSuites[i].File, testFile.TestSuites[i].Skipped)
		}

		if xmlTestSuites.TestSuites[i].Assertions != testFile.TestSuites[i].Assertions {
			t.Errorf("expected test suite for file '%s' to report %d assertions", xmlTestSuites.TestSuites[i].File, testFile.TestSuites[i].Assertions)
		}
	}

}

func TestBasicReduce(t *testing.T) {
	setup()
	defer tearDown()

	Reduce(ReduceFunctionParams{
		IncludeFilePattern:           "fixtures/*.xml",
		ExcludeFilePattern:           "",
		OutputPath:                   "output/",
		ReduceTestSuitesBy:           enums.TestSuiteFieldName,
		ReduceTestCasesBy:            enums.TestCaseFieldName,
		OperatorTestSuitesTests:      enums.AggregateOperationMean,
		OperatorTestSuitesFailed:     enums.AggregateOperationMean,
		OperatorTestSuitesErrors:     enums.AggregateOperationMean,
		OperatorTestSuitesSkipped:    enums.AggregateOperationMean,
		OperatorTestSuitesAssertions: enums.AggregateOperationMean,
		OperatorTestSuitesTime:       enums.AggregateOperationMean,
		OperatorTestCasesTime:        enums.AggregateOperationMean,
		RoundingMode:                 enums.RoundingModeRound,
	})

	assertTestFile(
		t,
		serialization.TestSuites{
			TestSuites: []serialization.TestSuite{
				{
					Name:       "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
					File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
					FileName:   "Sample.xml",
					Time:       49.09959481199999,
					Tests:      7,
					Failed:     0,
					Errors:     0,
					Skipped:    0,
					Assertions: 20,
					TestCases:  []serialization.TestCase{},
				},
			},
		},
	)
}
