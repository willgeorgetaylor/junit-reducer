# junit-reducer

[![GitHub CI](https://github.com/willgeorgetaylor/junit-reducer/actions/workflows/test.yml/badge.svg)](https://github.com/willgeorgetaylor/junit-reducer/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/willgeorgetaylor/junit-reducer/graph/badge.svg?token=08001J4XQH)](https://codecov.io/gh/willgeorgetaylor/junit-reducer)
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![GitHub release](https://img.shields.io/github/tag/willgeorgetaylor/junit-reducer.svg?label=release)](https://github.com/willgeorgetaylor/junit-reducer/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/willgeorgetaylor/junit-reducer)](https://goreportcard.com/report/github.com/willgeorgetaylor/junit-reducer)

JUnit Reducer is a command-line tool that aggregates the [JUnit test XML reports](https://www.ibm.com/docs/en/developer-for-zos/14.1?topic=formats-junit-xml-format) from your CI runs and reduces them to a single, lighter set of reports to be downloaded later during CI, to steer your test splitting algorithm (e.g., [split_tests](https://github.com/marketplace/actions/split-tests)). The most typical use case is to regularly update a 'running average' of your recent test reports, which can be downloaded to your test runners in [less time](https://github.com/willgeorgetaylor/junit-reducer?tab=readme-ov-file#faster-ci) and without running an [ongoing race condition risk](https://github.com/willgeorgetaylor/junit-reducer?tab=readme-ov-file#coverage-integrity).

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="./diagram-dark.png">
  <img alt="Diagram explaining how junit-reducer turns multiple sets of JUnit reports into a single set of JUnit reports." src="./diagram-light.png">
</picture>

## Quickstart

Typically, you'll be using `junit-reducer` within a scheduled [cron](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#schedule) task to take the last `X` days of JUnit XML reports, reduce them and upload the results. It's recommended to accumulate the JUnit XML reports from individual CI runs in a cloud storage service like AWS S3 or Google Cloud Storage, as opposed to the caching APIs available from the CI providers (GitHub Actions, CircleCI etc.) themselves, which are designed as _overwrite_ key-value stores.

### GitHub Actions

> [!TIP]
> If you're using GitHub Actions, check out the [Action](https://github.com/marketplace/actions/reduce-junit-xml-test-reports) on GitHub Marketplace, if you prefer that sort of thing.

```yaml
name: junit-test-report-averaging
run-name: Create Average JUnit Test Reports
on:
  schedule:
      # Run every morning at 8AM
      - cron:  '0 8 * * *'
jobs:
  reduce-reports:
    runs-on: ubuntu-latest
    steps:
      # Configure with the Cloud storage provider of your choice.
      - name: Setup AWS CLI
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.YOUR_AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.YOUR_AWS_SECRET_ACCESS_KEY }}
          aws-region: eu-west-2

      # Download all test reports from all CI runs.
      # It is recommended to set up a lifecycle rule, to remove objects older
      # than a certain age from this bucket/path. This will help to keep the test reports
      # current and keep this job from taking too long.
      - name: Download test timings
        run: |
          aws s3 cp s3://your-junit-report-bucket/ci-runs-reports/ reports/ \
            --recursive

      # Extract the binary for your target environment (assumed to be Linux). See full list
      # of releases here: https://github.com/willgeorgetaylor/junit-reducer/releases
      - name: Reduce reports
        run: |
          curl -L "https://github.com/willgeorgetaylor/junit-reducer/releases/latest/download/junit-reducer_Linux_x86_64.tar.gz" | tar -xzf -
          chmod +x junit-reducer
          ./junit-reducer \
            --include="./reports/**/*" \
            --output-path="./average-reports/"

      # Upload the reduced set of test reports to a dedicated bucket/path.
      # In your actual CI process, the CI runners will copy the contents of
      # this path locally, to be ingested by the test splitter.
      - name: Upload single set of averaged reports
        run: |
          aws s3 sync ./average-reports s3://your-junit-report-bucket/average-reports/ \
            --size-only \
            --cache-control max-age=86400
```

## Why?

As your test suite grows, you may want to start splitting tests between multiple test runners, to be **executed concurrently.** While it's relatively simple to divide up your test suites by files, using lines of code (LOC) as a proxy for test duration, the LOC metric is still just an approximation and will result in uneven individual (and therefore overall slower) test run times as your codebase and test suites change.

The preferable approach for splitting test suites accurately is to use **recently reported test times,** and the most popular format for exchanging test data (including timings) between tools is the [JUnit XML reports format](https://www.ibm.com/docs/en/developer-for-zos/14.1?topic=formats-junit-xml-format). While JUnit itself is a Java project, the schema that defines JUnit reports is equally applicable to any language and reports can be generated by most testing frameworks for JavaScript, Ruby, Python etc.

In busier projects, CI will be uploading reports frequently, so even if you take a small time window (for example, the last 24 hours), you could end up with 20MB+ of test reports. These reports need to be **downloaded to every runner in your concurrency set,** only to then perform the same splitting operation to **yield the exact same time estimates.** This means unnecessary and expensive work is being performed by each concurrent runner, potentially delaying the total test time by minutes and increasing CI costs.

### Faster CI ✅

You can solve this speed issue with `junit-reducer` by creating a set of reports that looks like the set produced by a single CI run. Importantly, the values for `time` taken by test suite (as well as other counts, like errors and tests) are reduced from the wider set of reports, typically by finding the `mean` of all of the aggregate `time` values. Other reducer operations, like `min` / `max` / `mode` / `median` / `sum`, are available to handle more non-standard distributions.

### Coverage integrity ✅

In very busy projects, there is also a more **problematic race condition possible**, with larger downloads and test runners starting at different times. As CI runs from other commits upload their reports to the same remote source that you're downloading from, if any of your concurrent runners download reports with different values, the input data is misaligned and the splitting operation is corrupted. However, because the download and splitting operation is being performed in a distributed manner (across all of the runners concurrently) this misalignment will result in some tests in your run being **skipped.**

This risk is mitigated by computing the averaged reports in one place, and updating that set as part of a scheduled job. This is exactly the approach outlined in the [quickstart](https://github.com/willgeorgetaylor/junit-reducer?tab=readme-ov-file#quickstart) section.

## Usage

Download and extract the latest build for your target environment, from the [releases page](https://github.com/willgeorgetaylor/junit-reducer/releases).

For a complete list of arguments:

```bash
$./junit-reducer --help
```

```
Flags:
  -h, --help                          help for junit-reducer
      --exclude string                Pattern to exclude from input JUnit XML reports
      --include string                Pattern to find input JUnit XML reports (default "./**/*.xml")
      --op-cases-time string          Operation for test cases time. Options: "max", "mean", "median", "min", "mode" or "sum" (default "mean")
      --op-suites-assertions string   Operation for test suites assertions. Options: "max", "mean", "median", "min", "mode" or "sum" (default "mean")
      --op-suites-errors string       Operation for test suites errors. Options: "max", "mean", "median", "min", "mode" or "sum" (default "mean")
      --op-suites-failed string       Operation for test suites failed. Options: "max", "mean", "median", "min", "mode" or "sum" (default "mean")
      --op-suites-skipped string      Operation for test suites skipped. Options: "max", "mean", "median", "min", "mode" or "sum" (default "mean")
      --op-suites-tests string        Operation for test suites tests. Options: "max", "mean", "median", "min", "mode" or "sum" (default "mean")
      --op-suites-time string         Operation for test suites time. Options: "max", "mean", "median", "min", "mode" or "sum" (default "mean")
      --output-path string            Output path for synthetic JUnit XML reports (default "./output/")
      --reduce-cases-by string        Reduce test cases by name, classname, or file. Options: "classname", "file" or "name" (default "name")
      --reduce-suites-by string       Reduce test suites by name or filepath or both. Options: "filepath", "name" or "name+filepath" (default "name+filepath")
      --rounding-mode string          Rounding mode for counts that should be integers. Options: "ceil", "floor" or "round" (default "round")
```

## Examples

### Basic usage

```bash
junit-reducer \
  --include="test-reports/**/*" \     # Input path for JUnit reports
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