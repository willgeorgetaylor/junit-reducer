[![GitHub CI](https://github.com/willgeorgetaylor/junit-reducer/actions/workflows/test.yml/badge.svg)](https://github.com/willgeorgetaylor/junit-reducer/actions/workflows/test.yml)
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![GitHub release](https://img.shields.io/github/tag/willgeorgetaylor/junit-reducer.svg?label=release)](https://github.com/willgeorgetaylor/junit-reducer/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/willgeorgetaylor/junit-reducer)](https://goreportcard.com/report/github.com/willgeorgetaylor/junit-reducer)

## junit-reducer

JUnit Reducer is a CLI tool that aggregates multiple sets of [JUnit test XML reports](https://www.ibm.com/docs/en/developer-for-zos/14.1?topic=formats-junit-xml-format) into a single, averaged XML report set. This helps you run faster continuous integration (CI) parallel tests by reducing data volume and normalizing test execution times.

## Use case

This tool is ideal for situations where you have to handle many JUnit reports, such as those generated in Continuous Integration (CI) systems, and need to distribute tests evenly across different runners based on their execution times. To counteract the fluctuations in time measurements across individual test runs, it's necessary to calculate an average of these times. However, downloading a full set of test reports at runtime can be time-consuming and resource-intensive. This utility addresses this issue by enabling the creation of a consolidated set of reports. These reports represent a 'running average' of test times while still adhering to the required JUnit XML format, facilitating efficient test splitting.

## Usage

Download and extract the latest build for your target environment, from the [releases page](https://github.com/willgeorgetaylor/junit-reducer/releases).

For a complete list of arguments:

```bash
$./junit-reducer --help
```

```
Flags:
  -h, --help                          help for junit-reducer
      --output-path string            Output path for synthetic JUnit XML reports (required) (default "./output/")
      --include string                Pattern to find input JUnit XML reports (required) (default "./**/*.xml")
      --exclude string                Pattern to exclude from input JUnit XML reports
      --op-cases-time string          Operator for test cases time. Options: "min", "max", "sum", "mean", "mode" or "median" (default "mean")
      --op-suites-assertions string   Operator for test suites assertions. Options: "max", "sum", "mean", "mode", "median" or "min" (default "mean")
      --op-suites-errors string       Operator for test suites errors. Options: "min", "max", "sum", "mean", "mode" or "median" (default "mean")
      --op-suites-failed string       Operator for test suites failed. Options: "median", "min", "max", "sum", "mean" or "mode" (default "mean")
      --op-suites-skipped string      Operator for test suites skipped. Options: "mode", "median", "min", "max", "sum" or "mean" (default "mean")
      --op-suites-tests string        Operator for test suites tests. Options: "mode", "median", "min", "max", "sum" or "mean" (default "mean")
      --op-suites-time string         Operator for test suites time. Options: "mode", "median", "min", "max", "sum" or "mean" (default "mean")
      --reduce-cases-by string        Reduce test cases by name, classname, or file. Options: "name", "classname" or "file" (default "name")
      --reduce-suites-by string       Reduce test suites by name or filepath or both. Options: "name", "filepath" or "name+filepath" (default "name+filepath")
      --rounding-mode string          Rounding mode for counts that should be integers. Options: "round", "ceil" or "floor" (default "round")
```

## Examples

### Basic usage

```bash
junit-reducer \
  --include="test-reports/**/*" \  # Input path for JUnit reports
  --output-path="avg-reports/"        # Output path for averaged reports
```

### Reduce by name

```bash
junit-reducer \
  --include="test-reports/**/*" \
  --output-path="avg-reports/" \
  --reduce-suites-by="name" \         # Grouping test suites by name
  --reduce-cases-by="classname"       # Grouping test cases by classname
```

### Reduce with non-average operations

```bash
junit-reducer \
  --include="test-reports/**/*" \
  --output-path="avg-reports/" \
  --op-suites-skipped="min" \         # Keeps min of skips across suites of same type
  --op-suites-failed="min" \          # Keeps min of failures across suites of same type
  --op-suites-errors="min" \          # Keeps min of errors across suites of same type
  --op-suites-tests="max" \           # Keeps max of tests across suites of same type
  --op-suites-assertions="max" \      # Keeps max of assertions across suites of same type
  --op-suites-time="mean" \           # Calculates mean of time across suites of same type
  --op-cases-time="mean"              # Calculates mean of time across cases of same type
```

### Rounding average counts

```bash
junit-reducer \
  --include="test-reports/**/*" \
  --output-path="avg-reports/" \
  --rounding-mode="floor"             # Specifies the rounding method
```
