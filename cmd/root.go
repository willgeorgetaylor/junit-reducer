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
	include                            string
	exclude                            string
	outputPath                         string
	reduceTestSuitesByString           string
	reduceTestCasesByString            string
	operatorTestSuitesSkippedString    string
	operatorTestSuitesFailedString     string
	operatorTestSuitesErrorsString     string
	operatorTestSuitesTestsString      string
	operatorTestSuitesAssertionsString string
	operatorTestSuitesTimeString       string
	operatorTestCasesTimeString        string
	roundingModeString                 string
)

func invalidSelectionMessage(field string, selection string, options []string) string {
	return fmt.Sprintf("Invalid selection '%s' for %s. Valid options are: %v", selection, field, options)
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
			fmt.Println(invalidSelectionMessage("reduce-test-suites-by-string", reduceTestSuitesByString, enums.GetTestSuiteFields()))
			return
		}

		reduceTestCasesBy, ok := enums.TestCaseFieldValues[reduceTestCasesByString]
		if !ok {
			fmt.Println(invalidSelectionMessage("reduce-test-cases-by-string", reduceTestCasesByString, enums.GetTestCaseFields()))
			return
		}

		operatorTestSuitesSkipped, ok := enums.AggregateOperationValues[operatorTestSuitesSkippedString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-suites-skipped-string", operatorTestSuitesSkippedString, enums.GetAggregateOperations()))
			return
		}

		operatorTestSuitesFailed, ok := enums.AggregateOperationValues[operatorTestSuitesFailedString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-suites-failed-string", operatorTestSuitesFailedString, enums.GetAggregateOperations()))
			return
		}

		operatorTestSuitesErrors, ok := enums.AggregateOperationValues[operatorTestSuitesErrorsString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-suites-errors-string", operatorTestSuitesErrorsString, enums.GetAggregateOperations()))
			return
		}

		operatorTestSuitesTests, ok := enums.AggregateOperationValues[operatorTestSuitesTestsString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-suites-tests-string", operatorTestSuitesTestsString, enums.GetAggregateOperations()))
			return
		}

		operatorTestSuitesAssertions, ok := enums.AggregateOperationValues[operatorTestSuitesAssertionsString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-suites-assertions-string", operatorTestSuitesAssertionsString, enums.GetAggregateOperations()))
			return
		}

		operatorTestSuitesTime, ok := enums.AggregateOperationValues[operatorTestSuitesTimeString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-suites-time-string", operatorTestSuitesTimeString, enums.GetAggregateOperations()))
			return
		}

		operatorTestCasesTime, ok := enums.AggregateOperationValues[operatorTestCasesTimeString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-cases-time-string", operatorTestCasesTimeString, enums.GetAggregateOperations()))
			return
		}

		roundingMode, ok := enums.RoundingModeValues[roundingModeString]
		if !ok {
			fmt.Println(invalidSelectionMessage("rounding-mode-string", roundingModeString, enums.GetRoundingModes()))
			return
		}

		err := reducer.Reduce(
			reducer.ReduceFunctionParams{
				IncludeFilePattern:           include,
				ExcludeFilePattern:           exclude,
				OutputPath:                   outputPath,
				ReduceTestSuitesBy:           reduceTestSuitesBy,
				ReduceTestCasesBy:            reduceTestCasesBy,
				OperatorTestSuitesTests:      operatorTestSuitesTests,
				OperatorTestSuitesFailed:     operatorTestSuitesFailed,
				OperatorTestSuitesErrors:     operatorTestSuitesErrors,
				OperatorTestSuitesSkipped:    operatorTestSuitesSkipped,
				OperatorTestSuitesAssertions: operatorTestSuitesAssertions,
				OperatorTestSuitesTime:       operatorTestSuitesTime,
				OperatorTestCasesTime:        operatorTestCasesTime,
				RoundingMode:                 roundingMode,
			},
		)

		if err != nil {
			fmt.Println(err)
			return
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
	rootCmd.Flags().StringVar(&include, "include", "./**/*.xml", "Pattern to find input JUnit XML reports (required)")
	rootCmd.MarkFlagRequired("include")
	rootCmd.Flags().StringVar(&outputPath, "output-path", "./output/", "Output path for synthetic JUnit XML reports (required)")
	rootCmd.MarkFlagRequired("output-path")
	rootCmd.Flags().StringVar(&exclude, "exclude", "", "Pattern to exclude from input JUnit XML reports")
	rootCmd.Flags().StringVar(&reduceTestSuitesByString, "reduce-suites-by", enums.TestSuiteFieldKeys[enums.TestSuiteFieldNameFilepath], fmt.Sprintf("Reduce test suites by name or filepath or both. Options: %s", joinOptionsString(enums.GetTestSuiteFields())))
	rootCmd.Flags().StringVar(&reduceTestCasesByString, "reduce-cases-by", enums.TestCaseFieldKeys[enums.TestCaseFieldName], fmt.Sprintf("Reduce test cases by name, classname, or file. Options: %s", joinOptionsString(enums.GetTestCaseFields())))
	rootCmd.Flags().StringVar(&operatorTestSuitesSkippedString, "op-suites-skipped", enums.AggregateOperationKeys[enums.AggregateOperationMean], fmt.Sprintf("Operator for test suites skipped. Options: %s", joinOptionsString(enums.GetAggregateOperations())))
	rootCmd.Flags().StringVar(&operatorTestSuitesFailedString, "op-suites-failed", enums.AggregateOperationKeys[enums.AggregateOperationMean], fmt.Sprintf("Operator for test suites failed. Options: %s", joinOptionsString(enums.GetAggregateOperations())))
	rootCmd.Flags().StringVar(&operatorTestSuitesErrorsString, "op-suites-errors", enums.AggregateOperationKeys[enums.AggregateOperationMean], fmt.Sprintf("Operator for test suites errors. Options: %s", joinOptionsString(enums.GetAggregateOperations())))
	rootCmd.Flags().StringVar(&operatorTestSuitesTestsString, "op-suites-tests", enums.AggregateOperationKeys[enums.AggregateOperationMean], fmt.Sprintf("Operator for test suites tests. Options: %s", joinOptionsString(enums.GetAggregateOperations())))
	rootCmd.Flags().StringVar(&operatorTestSuitesAssertionsString, "op-suites-assertions", enums.AggregateOperationKeys[enums.AggregateOperationMean], fmt.Sprintf("Operator for test suites assertions. Options: %s", joinOptionsString(enums.GetAggregateOperations())))
	rootCmd.Flags().StringVar(&operatorTestSuitesTimeString, "op-suites-time", enums.AggregateOperationKeys[enums.AggregateOperationMean], fmt.Sprintf("Operator for test suites time. Options: %s", joinOptionsString(enums.GetAggregateOperations())))
	rootCmd.Flags().StringVar(&operatorTestCasesTimeString, "op-cases-time", enums.AggregateOperationKeys[enums.AggregateOperationMean], fmt.Sprintf("Operator for test cases time. Options: %s", joinOptionsString(enums.GetAggregateOperations())))
	rootCmd.Flags().StringVar(&roundingModeString, "rounding-mode", enums.RoundingModeKeys[enums.RoundingModeRound], fmt.Sprintf("Rounding mode for counts that should be integers. Options: %s", joinOptionsString(enums.GetRoundingModes())))
}
