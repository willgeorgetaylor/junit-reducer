## junit-reducer

JUnit Reducer is a CLI tool that aggregates multiple [JUnit test XML reports](https://www.ibm.com/docs/en/developer-for-zos/14.1?topic=formats-junit-xml-format) into a single, averaged XML report. This helps you run faster continuous integration (CI) parallel tests by reducing data volume and normalizing test execution times.

## Use case

This tool is ideal for situations where you have to handle many JUnit reports, such as those generated in Continuous Integration (CI) systems, and need to distribute tests evenly across different runners based on their execution times. To counteract the fluctuations in time measurements across individual test runs, it's necessary to calculate an average of these times. However, downloading a full set of test reports at runtime can be time-consuming and resource-intensive. This utility addresses this issue by enabling the creation of a consolidated set of reports. These reports represent a 'running average' of test times while still adhering to the required JUnit XML format, facilitating efficient test splitting.

## Usage

Download and extract the latest build for your target environment, from the [releases page](https://github.com/willgeorgetaylor/junit-reducer/releases).

For a complete list of arguments:

```bash
$./junit-reducer --help
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
  --reduce-testsuites-by="name" \     # Grouping test suites by name
  --reduce-testcases-by="classname"   # Grouping test cases by classname
```

### Reduce with non-average operations

```bash
junit-reducer \
  --input-path="test-reports/**/*" \
  --output-path="avg-reports/" \
  --testsuites-skipped="min" \        # Keeps min of skips across suites of same type
  --testsuites-failed="min" \         # Keeps min of failures across suites of same type
  --testsuites-errors="min" \         # Keeps min of errors across suites of same type
  --testsuites-tests="max" \          # Keeps max of tests across suites of same type
  --testsuites-assertions="max" \     # Keeps max of assertions across suites of same type
  --testsuites-time="mean" \          # Calculates mean of time across suites of same type
  --testcases-time="mean"             # Calculates mean of time across cases of same type
```

### Rounding average counts

```bash
junit-reducer \
  --input-path="test-reports/**/*" \
  --output-path="avg-reports/" \
  --int-rounding="down"               # Specifies the rounding method
```

### Preserve testcase children

```bash
junit-reducer \
  --input-path="test-reports/**/*" \
  --output-path="avg-reports/" \
  --preserveErrors="all" \            # Preserve ALL error children nodes (not recommended)
  --preserveSkips=1 \                 # Preserve 1 skip child node per case of same type
  --preserveFailures="none"           # Do not preserve failure children nodes
```


