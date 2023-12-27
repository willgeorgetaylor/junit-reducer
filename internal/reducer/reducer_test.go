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
	outputFileName := "output/" + junitReportFileName

	if !helpers.DirExists("output") {
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

		if len(xmlTestSuites.TestSuites[i].TestCases) != len(testFile.TestSuites[i].TestCases) {
			t.Errorf("expected test suite for file '%s' to own %d child test cases", xmlTestSuites.TestSuites[i].File, len(testFile.TestSuites[i].TestCases))
		}

		for j := 0; j < len(testFile.TestSuites[i].TestCases); j++ {
			var caseFound bool = false

			for k := 0; k < len(xmlTestSuites.TestSuites[i].TestCases); k++ {
				if xmlTestSuites.TestSuites[i].TestCases[k].Name == testFile.TestSuites[i].TestCases[j].Name &&
					xmlTestSuites.TestSuites[i].TestCases[k].Classname == testFile.TestSuites[i].TestCases[j].Classname &&
					xmlTestSuites.TestSuites[i].TestCases[k].File == testFile.TestSuites[i].TestCases[j].File &&
					xmlTestSuites.TestSuites[i].TestCases[k].Line == testFile.TestSuites[i].TestCases[j].Line &&
					xmlTestSuites.TestSuites[i].TestCases[k].Assertions == testFile.TestSuites[i].TestCases[j].Assertions &&
					xmlTestSuites.TestSuites[i].TestCases[k].Time == testFile.TestSuites[i].TestCases[j].Time {
					caseFound = true
				}
			}

			if !caseFound {
				t.Errorf("expected test suite for file '%s' to own identical test case with name '%s'", xmlTestSuites.TestSuites[i].File, testFile.TestSuites[i].TestCases[j].Name)
			}
		}
	}
}

func TestBasicReduce(t *testing.T) {
	setup()
	defer tearDown()

	err := Reduce(ReduceFunctionParams{
		IncludeFilePattern:           "fixtures/*.xml",
		ExcludeFilePattern:           "",
		OutputPath:                   "output/",
		ReduceTestSuitesBy:           enums.TestSuiteFieldNameFilepath,
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
					Tests:      5,
					Failed:     0,
					Errors:     0,
					Skipped:    0,
					Assertions: 17,
					TestCases: []serialization.TestCase{
						{
							Name:       "test_should_show_each_of_the_different_values_depending_on_which_billing_option_you_select",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       90,
							Assertions: 0,
							Time:       12.873165432000008,
						},
						{
							Name:       "test_should_not_be_able_to_view_a_background_check_without_background_check_viewer_role",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       176,
							Assertions: 2,
							Time:       4.697016684999994,
						},
						{
							Name:       "test_should_not_see_Creative_Services_Inc_integration_when_removed",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       53,
							Assertions: 1,
							Time:       4.148554419999982,
						},
					},
				},
			},
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
}

func TestExcludeFilePattern(t *testing.T) {
	setup()
	defer tearDown()

	err := Reduce(ReduceFunctionParams{
		IncludeFilePattern:           "fixtures/*.xml",
		ExcludeFilePattern:           "fixtures/Sample.xml",
		OutputPath:                   "output/",
		ReduceTestSuitesBy:           enums.TestSuiteFieldNameFilepath,
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

	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	assertTestFile(
		t,
		serialization.TestSuites{
			TestSuites: []serialization.TestSuite{
				{
					Name:       "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
					File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
					FileName:   "Sample2.xml",
					Time:       69.09959481199999,
					Tests:      7,
					Failed:     0,
					Errors:     0,
					Skipped:    0,
					Assertions: 20,
					TestCases: []serialization.TestCase{
						{
							Name:       "test_should_show_each_of_the_different_values_depending_on_which_billing_option_you_select",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       90,
							Assertions: 0,
							Time:       10.373165432000008,
						},
						{
							Name:       "test_should_not_be_able_to_view_a_background_check_without_background_check_viewer_role",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       176,
							Assertions: 2,
							Time:       2.697016684999994,
						},
						{
							Name:       "test_should_not_see_Creative_Services_Inc_integration_when_removed",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       53,
							Assertions: 1,
							Time:       1.1485544199999822,
						},
					},
				},
			},
		},
	)
}

func TestReduceTestSuitesByFilepath(t *testing.T) {
	setup()
	defer tearDown()

	err := Reduce(ReduceFunctionParams{
		IncludeFilePattern:           "fixtures/*.xml",
		ExcludeFilePattern:           "",
		OutputPath:                   "output/",
		ReduceTestSuitesBy:           enums.TestSuiteFieldFilepath,
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

	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	assertTestFile(
		t,
		serialization.TestSuites{
			TestSuites: []serialization.TestSuite{
				{
					Name:       "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
					File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
					FileName:   "Sample.xml",
					Time:       49.09959481199999,
					Tests:      5,
					Failed:     0,
					Errors:     0,
					Skipped:    0,
					Assertions: 17,
					TestCases: []serialization.TestCase{
						{
							Name:       "test_should_show_each_of_the_different_values_depending_on_which_billing_option_you_select",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       90,
							Assertions: 0,
							Time:       12.873165432000008,
						},
						{
							Name:       "test_should_not_be_able_to_view_a_background_check_without_background_check_viewer_role",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       176,
							Assertions: 2,
							Time:       4.697016684999994,
						},
						{
							Name:       "test_should_not_see_Creative_Services_Inc_integration_when_removed",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       53,
							Assertions: 1,
							Time:       4.148554419999982,
						},
					},
				},
			},
		},
	)
}

func TestReduceTestCasesByClassName(t *testing.T) {
	setup()
	defer tearDown()

	err := Reduce(ReduceFunctionParams{
		IncludeFilePattern:           "fixtures/*.xml",
		ExcludeFilePattern:           "",
		OutputPath:                   "output/",
		ReduceTestSuitesBy:           enums.TestSuiteFieldNameFilepath,
		ReduceTestCasesBy:            enums.TestCaseFieldClassname,
		OperatorTestSuitesTests:      enums.AggregateOperationMean,
		OperatorTestSuitesFailed:     enums.AggregateOperationMean,
		OperatorTestSuitesErrors:     enums.AggregateOperationMean,
		OperatorTestSuitesSkipped:    enums.AggregateOperationMean,
		OperatorTestSuitesAssertions: enums.AggregateOperationMean,
		OperatorTestSuitesTime:       enums.AggregateOperationMean,
		OperatorTestCasesTime:        enums.AggregateOperationMean,
		RoundingMode:                 enums.RoundingModeRound,
	})

	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	assertTestFile(
		t,
		serialization.TestSuites{
			TestSuites: []serialization.TestSuite{
				{
					Name:       "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
					File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
					FileName:   "Sample.xml",
					Time:       49.09959481199999,
					Tests:      5,
					Failed:     0,
					Errors:     0,
					Skipped:    0,
					Assertions: 17,
					TestCases: []serialization.TestCase{
						{
							Name:       "test_should_show_each_of_the_different_values_depending_on_which_billing_option_you_select",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       90,
							Assertions: 0,
							Time:       7.239578845666661,
						},
					},
				},
			},
		},
	)
}

func TestReduceTestCasesByFilename(t *testing.T) {
	setup()
	defer tearDown()

	err := Reduce(ReduceFunctionParams{
		IncludeFilePattern:           "fixtures/*.xml",
		ExcludeFilePattern:           "",
		OutputPath:                   "output/",
		ReduceTestSuitesBy:           enums.TestSuiteFieldNameFilepath,
		ReduceTestCasesBy:            enums.TestCaseFieldFile,
		OperatorTestSuitesTests:      enums.AggregateOperationMean,
		OperatorTestSuitesFailed:     enums.AggregateOperationMean,
		OperatorTestSuitesErrors:     enums.AggregateOperationMean,
		OperatorTestSuitesSkipped:    enums.AggregateOperationMean,
		OperatorTestSuitesAssertions: enums.AggregateOperationMean,
		OperatorTestSuitesTime:       enums.AggregateOperationMean,
		OperatorTestCasesTime:        enums.AggregateOperationMean,
		RoundingMode:                 enums.RoundingModeRound,
	})

	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	assertTestFile(
		t,
		serialization.TestSuites{
			TestSuites: []serialization.TestSuite{
				{
					Name:       "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
					File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
					FileName:   "Sample.xml",
					Time:       49.09959481199999,
					Tests:      5,
					Failed:     0,
					Errors:     0,
					Skipped:    0,
					Assertions: 17,
					TestCases: []serialization.TestCase{
						{
							Name:       "test_should_show_each_of_the_different_values_depending_on_which_billing_option_you_select",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       90,
							Assertions: 0,
							Time:       7.239578845666661,
						},
					},
				},
			},
		},
	)
}

func TestMaxAggOperation(t *testing.T) {
	setup()
	defer tearDown()

	err := Reduce(ReduceFunctionParams{
		IncludeFilePattern:           "fixtures/*.xml",
		ExcludeFilePattern:           "",
		OutputPath:                   "output/",
		ReduceTestSuitesBy:           enums.TestSuiteFieldNameFilepath,
		ReduceTestCasesBy:            enums.TestCaseFieldName,
		OperatorTestSuitesTests:      enums.AggregateOperationMax,
		OperatorTestSuitesFailed:     enums.AggregateOperationMax,
		OperatorTestSuitesErrors:     enums.AggregateOperationMax,
		OperatorTestSuitesSkipped:    enums.AggregateOperationMax,
		OperatorTestSuitesAssertions: enums.AggregateOperationMax,
		OperatorTestSuitesTime:       enums.AggregateOperationMax,
		OperatorTestCasesTime:        enums.AggregateOperationMax,
		RoundingMode:                 enums.RoundingModeRound,
	})

	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	assertTestFile(
		t,
		serialization.TestSuites{
			TestSuites: []serialization.TestSuite{
				{
					Name:       "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
					File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
					FileName:   "Sample.xml",
					Time:       69.09959481199999,
					Tests:      7,
					Failed:     0,
					Errors:     0,
					Skipped:    0,
					Assertions: 20,
					TestCases: []serialization.TestCase{
						{
							Name:       "test_should_show_each_of_the_different_values_depending_on_which_billing_option_you_select",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       90,
							Assertions: 0,
							Time:       15.373165432000008,
						},
						{
							Name:       "test_should_not_be_able_to_view_a_background_check_without_background_check_viewer_role",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       176,
							Assertions: 2,
							Time:       6.697016684999994,
						},
						{
							Name:       "test_should_not_see_Creative_Services_Inc_integration_when_removed",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       53,
							Assertions: 1,
							Time:       7.148554419999982,
						},
					},
				},
			},
		},
	)
}

func TestMinAggOperation(t *testing.T) {
	setup()
	defer tearDown()

	err := Reduce(ReduceFunctionParams{
		IncludeFilePattern:           "fixtures/*.xml",
		ExcludeFilePattern:           "",
		OutputPath:                   "output/",
		ReduceTestSuitesBy:           enums.TestSuiteFieldNameFilepath,
		ReduceTestCasesBy:            enums.TestCaseFieldName,
		OperatorTestSuitesTests:      enums.AggregateOperationMin,
		OperatorTestSuitesFailed:     enums.AggregateOperationMin,
		OperatorTestSuitesErrors:     enums.AggregateOperationMin,
		OperatorTestSuitesSkipped:    enums.AggregateOperationMin,
		OperatorTestSuitesAssertions: enums.AggregateOperationMin,
		OperatorTestSuitesTime:       enums.AggregateOperationMin,
		OperatorTestCasesTime:        enums.AggregateOperationMin,
		RoundingMode:                 enums.RoundingModeRound,
	})

	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	assertTestFile(
		t,
		serialization.TestSuites{
			TestSuites: []serialization.TestSuite{
				{
					Name:       "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
					File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
					FileName:   "Sample.xml",
					Time:       29.099594811999992,
					Tests:      2,
					Failed:     0,
					Errors:     0,
					Skipped:    0,
					Assertions: 14,
					TestCases: []serialization.TestCase{
						{
							Name:       "test_should_show_each_of_the_different_values_depending_on_which_billing_option_you_select",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       90,
							Assertions: 0,
							Time:       10.373165432000008,
						},
						{
							Name:       "test_should_not_be_able_to_view_a_background_check_without_background_check_viewer_role",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       176,
							Assertions: 2,
							Time:       2.697016684999994,
						},
						{
							Name:       "test_should_not_see_Creative_Services_Inc_integration_when_removed",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       53,
							Assertions: 1,
							Time:       1.1485544199999822,
						},
					},
				},
			},
		},
	)
}

func TestSumAggOperation(t *testing.T) {
	setup()
	defer tearDown()

	err := Reduce(ReduceFunctionParams{
		IncludeFilePattern:           "fixtures/*.xml",
		ExcludeFilePattern:           "",
		OutputPath:                   "output/",
		ReduceTestSuitesBy:           enums.TestSuiteFieldNameFilepath,
		ReduceTestCasesBy:            enums.TestCaseFieldName,
		OperatorTestSuitesTests:      enums.AggregateOperationSum,
		OperatorTestSuitesFailed:     enums.AggregateOperationSum,
		OperatorTestSuitesErrors:     enums.AggregateOperationSum,
		OperatorTestSuitesSkipped:    enums.AggregateOperationSum,
		OperatorTestSuitesAssertions: enums.AggregateOperationSum,
		OperatorTestSuitesTime:       enums.AggregateOperationSum,
		OperatorTestCasesTime:        enums.AggregateOperationSum,
		RoundingMode:                 enums.RoundingModeRound,
	})

	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	assertTestFile(
		t,
		serialization.TestSuites{
			TestSuites: []serialization.TestSuite{
				{
					Name:       "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
					File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
					FileName:   "Sample.xml",
					Time:       98.19918962399998,
					Tests:      9,
					Failed:     0,
					Errors:     0,
					Skipped:    0,
					Assertions: 34,
					TestCases: []serialization.TestCase{
						{
							Name:       "test_should_show_each_of_the_different_values_depending_on_which_billing_option_you_select",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       90,
							Assertions: 0,
							Time:       25.746330864000015,
						},
						{
							Name:       "test_should_not_be_able_to_view_a_background_check_without_background_check_viewer_role",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       176,
							Assertions: 2,
							Time:       9.394033369999988,
						},
						{
							Name:       "test_should_not_see_Creative_Services_Inc_integration_when_removed",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       53,
							Assertions: 1,
							Time:       8.297108839999964,
						},
					},
				},
			},
		},
	)
}

func TestMedianAggOperation(t *testing.T) {
	setup()
	defer tearDown()

	err := Reduce(ReduceFunctionParams{
		IncludeFilePattern:           "fixtures/*.xml",
		ExcludeFilePattern:           "",
		OutputPath:                   "output/",
		ReduceTestSuitesBy:           enums.TestSuiteFieldNameFilepath,
		ReduceTestCasesBy:            enums.TestCaseFieldClassname,
		OperatorTestSuitesTests:      enums.AggregateOperationMedian,
		OperatorTestSuitesFailed:     enums.AggregateOperationMedian,
		OperatorTestSuitesErrors:     enums.AggregateOperationMedian,
		OperatorTestSuitesSkipped:    enums.AggregateOperationMedian,
		OperatorTestSuitesAssertions: enums.AggregateOperationMedian,
		OperatorTestSuitesTime:       enums.AggregateOperationMedian,
		OperatorTestCasesTime:        enums.AggregateOperationMedian,
		RoundingMode:                 enums.RoundingModeRound,
	})

	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	assertTestFile(
		t,
		serialization.TestSuites{
			TestSuites: []serialization.TestSuite{
				{
					Name:       "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
					File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
					FileName:   "Sample.xml",
					Time:       29.099594811999992,
					Tests:      2,
					Failed:     0,
					Errors:     0,
					Skipped:    0,
					Assertions: 14,
					TestCases: []serialization.TestCase{
						{
							Name:       "test_should_show_each_of_the_different_values_depending_on_which_billing_option_you_select",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       90,
							Assertions: 0,
							Time:       6.697016684999994,
						},
					},
				},
			},
		},
	)
}

func TestRoundingModeCeil(t *testing.T) {
	setup()
	defer tearDown()

	err := Reduce(ReduceFunctionParams{
		IncludeFilePattern:           "fixtures/*.xml",
		ExcludeFilePattern:           "",
		OutputPath:                   "output/",
		ReduceTestSuitesBy:           enums.TestSuiteFieldNameFilepath,
		ReduceTestCasesBy:            enums.TestCaseFieldName,
		OperatorTestSuitesTests:      enums.AggregateOperationMean,
		OperatorTestSuitesFailed:     enums.AggregateOperationMean,
		OperatorTestSuitesErrors:     enums.AggregateOperationMean,
		OperatorTestSuitesSkipped:    enums.AggregateOperationMean,
		OperatorTestSuitesAssertions: enums.AggregateOperationMean,
		OperatorTestSuitesTime:       enums.AggregateOperationMean,
		OperatorTestCasesTime:        enums.AggregateOperationMean,
		RoundingMode:                 enums.RoundingModeCeil,
	})

	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	assertTestFile(
		t,
		serialization.TestSuites{
			TestSuites: []serialization.TestSuite{
				{
					Name:       "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
					File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
					FileName:   "Sample.xml",
					Time:       49.09959481199999,
					Tests:      5,
					Failed:     0,
					Errors:     0,
					Skipped:    0,
					Assertions: 17,
					TestCases: []serialization.TestCase{
						{
							Name:       "test_should_show_each_of_the_different_values_depending_on_which_billing_option_you_select",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       90,
							Assertions: 0,
							Time:       12.873165432000008,
						},
						{
							Name:       "test_should_not_be_able_to_view_a_background_check_without_background_check_viewer_role",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       176,
							Assertions: 2,
							Time:       4.697016684999994,
						},
						{
							Name:       "test_should_not_see_Creative_Services_Inc_integration_when_removed",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       53,
							Assertions: 1,
							Time:       4.148554419999982,
						},
					},
				},
			},
		},
	)
}

func TestRoundingModeFloor(t *testing.T) {
	setup()
	defer tearDown()

	err := Reduce(ReduceFunctionParams{
		IncludeFilePattern:           "fixtures/*.xml",
		ExcludeFilePattern:           "",
		OutputPath:                   "output/",
		ReduceTestSuitesBy:           enums.TestSuiteFieldNameFilepath,
		ReduceTestCasesBy:            enums.TestCaseFieldName,
		OperatorTestSuitesTests:      enums.AggregateOperationMean,
		OperatorTestSuitesFailed:     enums.AggregateOperationMean,
		OperatorTestSuitesErrors:     enums.AggregateOperationMean,
		OperatorTestSuitesSkipped:    enums.AggregateOperationMean,
		OperatorTestSuitesAssertions: enums.AggregateOperationMean,
		OperatorTestSuitesTime:       enums.AggregateOperationMean,
		OperatorTestCasesTime:        enums.AggregateOperationMean,
		RoundingMode:                 enums.RoundingModeFloor,
	})

	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	assertTestFile(
		t,
		serialization.TestSuites{
			TestSuites: []serialization.TestSuite{
				{
					Name:       "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
					File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
					FileName:   "Sample.xml",
					Time:       49.09959481199999,
					Tests:      4,
					Failed:     0,
					Errors:     0,
					Skipped:    0,
					Assertions: 17,
					TestCases: []serialization.TestCase{
						{
							Name:       "test_should_show_each_of_the_different_values_depending_on_which_billing_option_you_select",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       90,
							Assertions: 0,
							Time:       12.873165432000008,
						},
						{
							Name:       "test_should_not_be_able_to_view_a_background_check_without_background_check_viewer_role",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       176,
							Assertions: 2,
							Time:       4.697016684999994,
						},
						{
							Name:       "test_should_not_see_Creative_Services_Inc_integration_when_removed",
							Classname:  "Admin::Jobs::Applications::Actions::CreativeServicesIncBackgroundCheckTest",
							File:       "test/system/admin/jobs/applications/actions/creative_services_inc_background_check_test.rb",
							Line:       53,
							Assertions: 1,
							Time:       4.148554419999982,
						},
					},
				},
			},
		},
	)
}
