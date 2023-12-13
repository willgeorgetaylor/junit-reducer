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
  -h, --help                                     help for junit-reducer
      --input-path string                        Glob pattern for input JUnit XML reports (required) (default "./**/*.xml")
      --output-path string                       Output path for synthetic JUnit XML reports (required) (default "./output/")
      --operator-test-cases-time string          Operator for test cases time (default "mean")
      --operator-test-suites-assertions string   Operator for test suites assertions (default "mean")
      --operator-test-suites-errors string       Operator for test suites errors (default "mean")
      --operator-test-suites-failed string       Operator for test suites failed (default "mean")
      --operator-test-suites-skipped string      Operator for test suites skipped (default "mean")
      --operator-test-suites-tests string        Operator for test suites tests (default "mean")
      --operator-test-suites-time string         Operator for test suites time (default "mean")
      --preserve-errors string                   Preserve errors in output report (default "none")
      --preserve-failures string                 Preserve failures in output report (default "none")
      --preserve-skips string                    Preserve skips in output report (default "none")
      --reduce-test-cases-by string              Reduce test cases by name, classname, or file (default "name")
      --reduce-test-suites-by string             Reduce test suites by name or filepath (default "name")
      --rounding-mode string                     Rounding mode for integer counts (failures, errors etc.) that produce non-integer averages (default "round")
```

## Examples

### Basic usage

```bash
junit-reducer \
  --input-path="test-reports/**/*" \  # Input path for JUnit reports
  --output-path="avg-reports/"        # Output path for averaged reports
```

### Reduce by name

```bash
junit-reducer \
  --input-path="test-reports/**/*" \
  --output-path="avg-reports/" \
  --reduce-test-suites-by="name" \     # Grouping test suites by name
  --reduce-test-cases-by="classname"   # Grouping test cases by classname
```

### Reduce with non-average operations

```bash
junit-reducer \
  --input-path="test-reports/**/*" \
  --output-path="avg-reports/" \
  --test-suites-skipped="min" \        # Keeps min of skips across suites of same type
  --test-suites-failed="min" \         # Keeps min of failures across suites of same type
  --test-suites-errors="min" \         # Keeps min of errors across suites of same type
  --test-suites-tests="max" \          # Keeps max of tests across suites of same type
  --test-suites-assertions="max" \     # Keeps max of assertions across suites of same type
  --test-suites-time="mean" \          # Calculates mean of time across suites of same type
  --test-cases-time="mean"             # Calculates mean of time across cases of same type
```

### Rounding average counts

```bash
junit-reducer \
  --input-path="test-reports/**/*" \
  --output-path="avg-reports/" \
  --int-rounding="floor"               # Specifies the rounding method
```

### Preserve testcase children

```bash
junit-reducer \
  --input-path="test-reports/**/*" \
  --output-path="avg-reports/" \
  --preserve-errors="all" \            # Preserve ALL error children nodes (not recommended)
  --preserve-skips="1" \               # Preserve 1 skip child node per case of same type
  --preserve-failures="none"           # Do not preserve failure children nodes
```