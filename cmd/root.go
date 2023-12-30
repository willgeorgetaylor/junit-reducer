/*
Copyright © 2023 Will Taylor
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/willgeorgetaylor/junit-reducer/internal/enums"
	"github.com/willgeorgetaylor/junit-reducer/internal/reducer"

	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	include                             string
	exclude                             string
	outputPath                          string
	reduceTestSuitesByString            string
	reduceTestCasesByString             string
	operationTestSuitesSkippedString    string
	operationTestSuitesFailedString     string
	operationTestSuitesErrorsString     string
	operationTestSuitesTestsString      string
	operationTestSuitesAssertionsString string
	operationTestSuitesTimeString       string
	operationTestCasesTimeString        string
	roundingModeString                  string
)

func invalidSelectionMessage(field string, selection string, options []string) string {
	return fmt.Sprintf("Invalid option '%s' for %s. Valid options are: %s", selection, field, joinOptionsString(options))
}

func joinOptionsString(options []string) string {
	var finalString string = ""
	for index, option := range options {
		if index > 0 && index == (len(options)-1) {
			finalString += " or "
		} else if index > 0 {
			finalString += ", "
		}
		finalString += fmt.Sprintf("\"%s\"", option)
	}
	return finalString
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "junit-reducer",
	Short: "Aggregates and optimizes JUnit reports for CI",
	Long:  `JUnit Reducer streamlines CI testing by averaging JUnit reports for balanced test runner distribution.`,
	Run: func(cmd *cobra.Command, args []string) {
		reduceTestSuitesBy, ok := enums.TestSuiteFieldValues[reduceTestSuitesByString]
		if !ok {
			fmt.Println(invalidSelectionMessage("reduce-test-suites-by", reduceTestSuitesByString, enums.GetTestSuiteFields()))
			os.Exit(1)
		}

		reduceTestCasesBy, ok := enums.TestCaseFieldValues[reduceTestCasesByString]
		if !ok {
			fmt.Println(invalidSelectionMessage("reduce-test-cases-by", reduceTestCasesByString, enums.GetTestCaseFields()))
			os.Exit(1)
		}

		operationTestSuitesSkipped, ok := enums.AggregateOperationValues[operationTestSuitesSkippedString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operation-test-suites-skipped", operationTestSuitesSkippedString, enums.GetAggregateOperations()))
			os.Exit(1)
		}

		operationTestSuitesFailed, ok := enums.AggregateOperationValues[operationTestSuitesFailedString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operation-test-suites-failed", operationTestSuitesFailedString, enums.GetAggregateOperations()))
			os.Exit(1)
		}

		operationTestSuitesErrors, ok := enums.AggregateOperationValues[operationTestSuitesErrorsString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operation-test-suites-errors", operationTestSuitesErrorsString, enums.GetAggregateOperations()))
			os.Exit(1)
		}

		operationTestSuitesTests, ok := enums.AggregateOperationValues[operationTestSuitesTestsString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operation-test-suites-tests", operationTestSuitesTestsString, enums.GetAggregateOperations()))
			os.Exit(1)
		}

		operationTestSuitesAssertions, ok := enums.AggregateOperationValues[operationTestSuitesAssertionsString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operation-test-suites-assertions", operationTestSuitesAssertionsString, enums.GetAggregateOperations()))
			os.Exit(1)
		}

		operationTestSuitesTime, ok := enums.AggregateOperationValues[operationTestSuitesTimeString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operation-test-suites-time", operationTestSuitesTimeString, enums.GetAggregateOperations()))
			os.Exit(1)
		}

		operationTestCasesTime, ok := enums.AggregateOperationValues[operationTestCasesTimeString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operation-test-cases-time", operationTestCasesTimeString, enums.GetAggregateOperations()))
			os.Exit(1)
		}

		roundingMode, ok := enums.RoundingModeValues[roundingModeString]
		if !ok {
			fmt.Println(invalidSelectionMessage("rounding-mode", roundingModeString, enums.GetRoundingModes()))
			os.Exit(1)
		}

		err := reducer.Reduce(
			reducer.ReduceFunctionParams{
				IncludeFilePattern:            include,
				ExcludeFilePattern:            exclude,
				OutputPath:                    outputPath,
				ReduceTestSuitesBy:            reduceTestSuitesBy,
				ReduceTestCasesBy:             reduceTestCasesBy,
				OperationTestSuitesTests:      operationTestSuitesTests,
				OperationTestSuitesFailed:     operationTestSuitesFailed,
				OperationTestSuitesErrors:     operationTestSuitesErrors,
				OperationTestSuitesSkipped:    operationTestSuitesSkipped,
				OperationTestSuitesAssertions: operationTestSuitesAssertions,
				OperationTestSuitesTime:       operationTestSuitesTime,
				OperationTestCasesTime:        operationTestCasesTime,
				RoundingMode:                  roundingMode,
			},
		)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

//nolint:errcheck // Ignore errors from MarkFlagRequired
func init() {
	rootCmd.Flags().StringVar(&include, "include", "./**/*.xml", "Glob pattern to find JUnit XML reports to reduce")
	rootCmd.Flags().StringVar(&outputPath, "output-path", "./output/", "Output path for the reduced JUnit XML reports")
	rootCmd.Flags().StringVar(&exclude, "exclude", "", "Glob pattern to omit from included JUnit XML reports")
	rootCmd.Flags().StringVar(&reduceTestSuitesByString, "reduce-suites-by", enums.TestSuiteFieldKeys[enums.TestSuiteFieldNameFilepath], fmt.Sprintf("Key to group and reduce test suites by. Options: %s", joinOptionsString(enums.GetTestSuiteFields())))
	rootCmd.Flags().StringVar(&reduceTestCasesByString, "reduce-cases-by", enums.TestCaseFieldKeys[enums.TestCaseFieldName], fmt.Sprintf("Key to group and reduce test cases by. Options: %s", joinOptionsString(enums.GetTestCaseFields())))
	rootCmd.Flags().StringVar(&operationTestSuitesSkippedString, "op-suites-skipped", enums.AggregateOperationKeys[enums.AggregateOperationMean], fmt.Sprintf("Reducer operation for test suite skipped counts. Options: %s", joinOptionsString(enums.GetAggregateOperations())))
	rootCmd.Flags().StringVar(&operationTestSuitesFailedString, "op-suites-failed", enums.AggregateOperationKeys[enums.AggregateOperationMean], fmt.Sprintf("Reducer operation for test suite failure counts. Options: %s", joinOptionsString(enums.GetAggregateOperations())))
	rootCmd.Flags().StringVar(&operationTestSuitesErrorsString, "op-suites-errors", enums.AggregateOperationKeys[enums.AggregateOperationMean], fmt.Sprintf("Reducer operation for test suite error counts. Options: %s", joinOptionsString(enums.GetAggregateOperations())))
	rootCmd.Flags().StringVar(&operationTestSuitesTestsString, "op-suites-tests", enums.AggregateOperationKeys[enums.AggregateOperationMean], fmt.Sprintf("Reducer operation for test suite test counts. Options: %s", joinOptionsString(enums.GetAggregateOperations())))
	rootCmd.Flags().StringVar(&operationTestSuitesAssertionsString, "op-suites-assertions", enums.AggregateOperationKeys[enums.AggregateOperationMean], fmt.Sprintf("Reducer operation for test suite assertion counts. Options: %s", joinOptionsString(enums.GetAggregateOperations())))
	rootCmd.Flags().StringVar(&operationTestSuitesTimeString, "op-suites-time", enums.AggregateOperationKeys[enums.AggregateOperationMean], fmt.Sprintf("Reducer operation for test suite time values. Options: %s", joinOptionsString(enums.GetAggregateOperations())))
	rootCmd.Flags().StringVar(&operationTestCasesTimeString, "op-cases-time", enums.AggregateOperationKeys[enums.AggregateOperationMean], fmt.Sprintf("Reducer operation for test case time values. Options: %s", joinOptionsString(enums.GetAggregateOperations())))
	rootCmd.Flags().StringVar(&roundingModeString, "rounding-mode", enums.RoundingModeKeys[enums.RoundingModeRound], fmt.Sprintf("Rounding mode for counts that should be integers in the final result. Options: %s", joinOptionsString(enums.GetRoundingModes())))
}
