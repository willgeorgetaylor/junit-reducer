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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "junit-reducer",
	Short: "Aggregates and optimizes JUnit reports for CI",
	Long:  `JUnit Reducer streamlines CI testing by averaging JUnit reports for balanced test runner distribution.`,
	Run: func(cmd *cobra.Command, args []string) {
		reduceTestSuitesBy, ok := enums.TestSuiteFieldValues[reduceTestSuitesByString]
		if !ok {
			fmt.Println(invalidSelectionMessage("reduce-test-suites-by-string", reduceTestSuitesByString, []string{"name", "filepath"}))
			return
		}

		reduceTestCasesBy, ok := enums.TestCaseFieldValues[reduceTestCasesByString]
		if !ok {
			fmt.Println(invalidSelectionMessage("reduce-test-cases-by-string", reduceTestCasesByString, []string{"name", "classname", "file"}))
			return
		}

		operatorTestSuitesSkipped, ok := enums.AggregateOperationValues[operatorTestSuitesSkippedString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-suites-skipped-string", operatorTestSuitesSkippedString, []string{"mean", "mode", "median", "min", "max", "sum"}))
			return
		}

		operatorTestSuitesFailed, ok := enums.AggregateOperationValues[operatorTestSuitesFailedString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-suites-failed-string", operatorTestSuitesFailedString, []string{"mean", "mode", "median", "min", "max", "sum"}))
			return
		}

		operatorTestSuitesErrors, ok := enums.AggregateOperationValues[operatorTestSuitesErrorsString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-suites-errors-string", operatorTestSuitesErrorsString, []string{"mean", "mode", "median", "min", "max", "sum"}))
			return
		}

		operatorTestSuitesTests, ok := enums.AggregateOperationValues[operatorTestSuitesTestsString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-suites-tests-string", operatorTestSuitesTestsString, []string{"mean", "mode", "median", "min", "max", "sum"}))
			return
		}

		operatorTestSuitesAssertions, ok := enums.AggregateOperationValues[operatorTestSuitesAssertionsString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-suites-assertions-string", operatorTestSuitesAssertionsString, []string{"mean", "mode", "median", "min", "max", "sum"}))
			return
		}

		operatorTestSuitesTime, ok := enums.AggregateOperationValues[operatorTestSuitesTimeString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-suites-time-string", operatorTestSuitesTimeString, []string{"mean", "mode", "median", "min", "max", "sum"}))
			return
		}

		operatorTestCasesTime, ok := enums.AggregateOperationValues[operatorTestCasesTimeString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-cases-time-string", operatorTestCasesTimeString, []string{"mean", "mode", "median", "min", "max", "sum"}))
			return
		}

		roundingMode, ok := enums.RoundingModeValues[roundingModeString]
		if !ok {
			fmt.Println(invalidSelectionMessage("rounding-mode-string", roundingModeString, []string{"round", "ceil", "floor"}))
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

func init() {
	rootCmd.Flags().StringVar(&include, "include", "./**/*.xml", "Pattern to find input JUnit XML reports (required)")
	rootCmd.MarkFlagRequired("include")
	rootCmd.Flags().StringVar(&outputPath, "output-path", "./output/", "Output path for synthetic JUnit XML reports (required)")
	rootCmd.MarkFlagRequired("output-path")
	rootCmd.Flags().StringVar(&exclude, "exclude", "", "Pattern to exclude from input JUnit XML reports")
	rootCmd.Flags().StringVar(&reduceTestSuitesByString, "reduce-suites-by", enums.TestSuiteFieldKeys[enums.TestSuiteFieldNameFilepath], "Reduce test suites by name or filepath or both")
	rootCmd.Flags().StringVar(&reduceTestCasesByString, "reduce-cases-by", enums.TestCaseFieldKeys[enums.TestCaseFieldName], "Reduce test cases by name, classname, or file")
	rootCmd.Flags().StringVar(&operatorTestSuitesSkippedString, "op-suites-skipped", enums.AggregateOperationKeys[enums.AggregateOperationMean], "Operator for test suites skipped")
	rootCmd.Flags().StringVar(&operatorTestSuitesFailedString, "op-suites-failed", enums.AggregateOperationKeys[enums.AggregateOperationMean], "Operator for test suites failed")
	rootCmd.Flags().StringVar(&operatorTestSuitesErrorsString, "op-suites-errors", enums.AggregateOperationKeys[enums.AggregateOperationMean], "Operator for test suites errors")
	rootCmd.Flags().StringVar(&operatorTestSuitesTestsString, "op-suites-tests", enums.AggregateOperationKeys[enums.AggregateOperationMean], "Operator for test suites tests")
	rootCmd.Flags().StringVar(&operatorTestSuitesAssertionsString, "op-suites-assertions", enums.AggregateOperationKeys[enums.AggregateOperationMean], "Operator for test suites assertions")
	rootCmd.Flags().StringVar(&operatorTestSuitesTimeString, "op-suites-time", enums.AggregateOperationKeys[enums.AggregateOperationMean], "Operator for test suites time")
	rootCmd.Flags().StringVar(&operatorTestCasesTimeString, "op-cases-time", enums.AggregateOperationKeys[enums.AggregateOperationMean], "Operator for test cases time")
	rootCmd.Flags().StringVar(&roundingModeString, "rounding-mode", enums.RoundingModeKeys[enums.RoundingModeRound], "Rounding mode for integer counts (failures, errors etc.) that produce non-integer averages")
}
