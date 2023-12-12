/*
Copyright © 2023 Will Taylor
*/
package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	inputPath                          string
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
	preserveErrors                     string
	preserveSkips                      string
	preserveFailures                   string
)

func invalidSelectionMessage(field string, selection string, options []string) string {
	return fmt.Sprintf("Invalid selection '%s' for %s. Valid options are: %v", selection, field, options)
}

func validatePreserveFlag(flag string) (string, bool) {
	numberRegex := regexp.MustCompile(`^\d+$`)
	if flag == "none" || flag == "all" || numberRegex.MatchString(flag) {
		return flag, true
	} else {
		return "", false
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "junit-reducer",
	Short: "Aggregates and optimizes JUnit reports for CI",
	Long:  `JUnit Reducer streamlines CI testing by averaging JUnit reports for balanced test runner distribution.`,
	Run: func(cmd *cobra.Command, args []string) {
		reduceTestSuitesBy, ok := testSuiteFieldValues[reduceTestSuitesByString]
		if !ok {
			fmt.Println(invalidSelectionMessage("reduce-test-suites-by-string", reduceTestSuitesByString, []string{"name", "filepath"}))
			return
		}

		reduceTestCasesBy, ok := testCaseFieldValues[reduceTestCasesByString]
		if !ok {
			fmt.Println(invalidSelectionMessage("reduce-test-cases-by-string", reduceTestCasesByString, []string{"name", "classname", "file"}))
			return
		}

		operatorTestSuitesSkipped, ok := aggregateOperationValues[operatorTestSuitesSkippedString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-suites-skipped-string", operatorTestSuitesSkippedString, []string{"mean", "mode", "median", "min", "max", "sum"}))
			return
		}

		operatorTestSuitesFailed, ok := aggregateOperationValues[operatorTestSuitesFailedString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-suites-failed-string", operatorTestSuitesFailedString, []string{"mean", "mode", "median", "min", "max", "sum"}))
			return
		}

		operatorTestSuitesErrors, ok := aggregateOperationValues[operatorTestSuitesErrorsString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-suites-errors-string", operatorTestSuitesErrorsString, []string{"mean", "mode", "median", "min", "max", "sum"}))
			return
		}

		operatorTestSuitesTests, ok := aggregateOperationValues[operatorTestSuitesTestsString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-suites-tests-string", operatorTestSuitesTestsString, []string{"mean", "mode", "median", "min", "max", "sum"}))
			return
		}

		operatorTestSuitesAssertions, ok := aggregateOperationValues[operatorTestSuitesAssertionsString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-suites-assertions-string", operatorTestSuitesAssertionsString, []string{"mean", "mode", "median", "min", "max", "sum"}))
			return
		}

		operatorTestSuitesTime, ok := aggregateOperationValues[operatorTestSuitesTimeString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-suites-time-string", operatorTestSuitesTimeString, []string{"mean", "mode", "median", "min", "max", "sum"}))
			return
		}

		operatorTestCasesTime, ok := aggregateOperationValues[operatorTestCasesTimeString]
		if !ok {
			fmt.Println(invalidSelectionMessage("operator-test-cases-time-string", operatorTestCasesTimeString, []string{"mean", "mode", "median", "min", "max", "sum"}))
			return
		}

		roundingMode, ok := roundingModeValues[roundingModeString]
		if !ok {
			fmt.Println(invalidSelectionMessage("rounding-mode-string", roundingModeString, []string{"round", "ceil", "floor"}))
			return
		}

		preserveErrors, ok := validatePreserveFlag(preserveErrors)
		if !ok {
			fmt.Println(invalidSelectionMessage("preserve-errors", preserveErrors, []string{"none", "all", "<number>"}))
			return
		}

		preserveSkips, ok := validatePreserveFlag(preserveSkips)
		if !ok {
			fmt.Println(invalidSelectionMessage("preserve-skips", preserveSkips, []string{"none", "all", "<number>"}))
			return
		}

		preserveFailures, ok := validatePreserveFlag(preserveFailures)
		if !ok {
			fmt.Println(invalidSelectionMessage("preserve-failures", preserveFailures, []string{"none", "all", "<number>"}))
			return
		}

		fmt.Println("junit-reducer")
		fmt.Println("Input path:", inputPath)
		fmt.Println("Output path:", outputPath)
		fmt.Println("Reduce test suites by:", reduceTestSuitesBy)
		fmt.Println("Reduce test cases by:", reduceTestCasesBy)
		fmt.Println("Operator test suites skipped:", operatorTestSuitesSkipped)
		fmt.Println("Operator test suites failed:", operatorTestSuitesFailed)
		fmt.Println("Operator test suites errors:", operatorTestSuitesErrors)
		fmt.Println("Operator test suites tests:", operatorTestSuitesTests)
		fmt.Println("Operator test suites assertions:", operatorTestSuitesAssertions)
		fmt.Println("Operator test suites time:", operatorTestSuitesTime)
		fmt.Println("Operator test cases time:", operatorTestCasesTime)
		fmt.Println("Rounding mode:", roundingMode)
		fmt.Println("Preserve errors:", preserveErrors)
		fmt.Println("Preserve skips:", preserveSkips)
		fmt.Println("Preserve failures:", preserveFailures)
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
	rootCmd.Flags().StringVar(&inputPath, "input-path", "./**/*.xml", "Glob pattern for input JUnit XML reports (required)")
	rootCmd.MarkFlagRequired("input-path")
	rootCmd.Flags().StringVar(&outputPath, "output-path", "./output/", "Output path for synthetic JUnit XML reports (required)")
	rootCmd.MarkFlagRequired("output-path")
	rootCmd.Flags().StringVar(&reduceTestSuitesByString, "reduce-test-suites-by", testSuiteFieldKeys[TestSuiteFieldName], "Reduce test suites by name or filepath")
	rootCmd.Flags().StringVar(&reduceTestCasesByString, "reduce-test-cases-by", testCaseFieldKeys[TestCaseFieldName], "Reduce test cases by name, classname, or file")
	rootCmd.Flags().StringVar(&operatorTestSuitesSkippedString, "operator-test-suites-skipped", aggregateOperationKeys[AggregateOperationMean], "Operator for test suites skipped")
	rootCmd.Flags().StringVar(&operatorTestSuitesFailedString, "operator-test-suites-failed", aggregateOperationKeys[AggregateOperationMean], "Operator for test suites failed")
	rootCmd.Flags().StringVar(&operatorTestSuitesErrorsString, "operator-test-suites-errors", aggregateOperationKeys[AggregateOperationMean], "Operator for test suites errors")
	rootCmd.Flags().StringVar(&operatorTestSuitesTestsString, "operator-test-suites-tests", aggregateOperationKeys[AggregateOperationMean], "Operator for test suites tests")
	rootCmd.Flags().StringVar(&operatorTestSuitesAssertionsString, "operator-test-suites-assertions", aggregateOperationKeys[AggregateOperationMean], "Operator for test suites assertions")
	rootCmd.Flags().StringVar(&operatorTestSuitesTimeString, "operator-test-suites-time", aggregateOperationKeys[AggregateOperationMean], "Operator for test suites time")
	rootCmd.Flags().StringVar(&operatorTestCasesTimeString, "operator-test-cases-time", aggregateOperationKeys[AggregateOperationMean], "Operator for test cases time")
	rootCmd.Flags().StringVar(&roundingModeString, "rounding-mode", roundingModeKeys[RoundingModeRound], "Rounding mode for integer counts (failures, errors etc.) that produce non-integer averages")
	rootCmd.Flags().StringVar(&preserveErrors, "preserve-errors", "none", "Preserve errors in output report")
	rootCmd.Flags().StringVar(&preserveSkips, "preserve-skips", "none", "Preserve skips in output report")
	rootCmd.Flags().StringVar(&preserveFailures, "preserve-failures", "none", "Preserve failures in output report")
}
