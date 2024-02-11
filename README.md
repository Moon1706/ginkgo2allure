# Ginkgo2Allure

CLI and library that are used to convert Ginkgo JSON reports into Allure JSON reports. Globally used for E2E and integration testing.

## Dependencies

There is no need to install additional dependencies on the user's part. The project itself relies heavily on the [OzonTech library](github.com/ozontech/allure-go).

## Installation

Simply download the binary file and use it.

## Usage

### Basic

#### CLI

```sh
# Get Ginkgo JSON report
ginkgo -r --keep-going -p --json-report=report.json ./tests/e2e/
# Create folder for Allure results
mkdir -p ./allure-results/
# See help
ginkgo2allure -h
# Convert Ginkgo report to Allure reports
ginkgo2allure ./report.json ./allure-results/
# Archive Allure results
zip -r allure-results.zip ./allure-results/
# Send zip archive to Allure server
```

#### Lib

```go
import (
    . "github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/types"
	"github.com/Moon1706/ginkgo2allure/pkg/convert"
	fmngr "github.com/Moon1706/ginkgo2allure/pkg/convert/file_manager"
	"github.com/Moon1706/ginkgo2allure/pkg/convert/parser"
)

var _ = ReportAfterSuite("allure report", func(report types.Report) {
    allureReports, err := convert.GinkgoToAllureReport(report, parser.NewDefaultParser, parser.Config{})
    if err != nil {
		panic(err)
	}

	fileManager := fmngr.NewFileManager("./allure-results")
	errs = convert.PrintAllureReports(allureReports, fileManager)
    if err != nil {
		panic(err)
	}
})
```

## Build

```sh
make bin
```

## Test

```sh
make coverage
```
