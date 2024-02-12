# Ginkgo2Allure

CLI and library that are used to convert Ginkgo JSON reports into Allure JSON reports. Globally used for E2E and integration testing.

>NOTE: Was tested only for Ginkgo version `>=v2.3.0`.

## Presettings

### TestCaseID issue

For production usage, it is really important to avoid test case duplication on the Allure server. For these goals, Allure realises the `TestCaseID` mechanism. The main idea here is to **ALWAYS** use the same test cases with an identical ID. The format of `TestCaseID` is an MD5 hash. So, for implementation this behaviour, I decided to attach to each Ginkgo test `id` label with UUIDv4. This `id` will transform into `TestCaseID` and will allow us to match the current test to the previous one in Allure and also automatically sort them by start/stop runtime.

For you, it means that you **MUST** add for each test label `id`. Example:

```go
It("test", Label("id=b1f3572c-f1f0-4001-a4b6-97625206d9f9"), func() {
    ...
})
```

For test aims, you can use the flag `--auto_gen_id`, which automatically generates a UUID for tests that don't have `id` in labels. Keep in mind that it's not a consistent UUID, and it will change in the next regeneration.

### Allure labels

#### General

General test case `It` label format is `<name><separator><value>`. If you want, you can change separator used flag `--label_separator`.

You **HAVE TO** understand that all labels that were added in `It` labels will be checked in a loop, and if theirs can be split by a separator on **2** parts, they will be added to the final Allure labels. That feature allows you to add your own labels, like `owner`, `feature`, `story`, etc. Example `e2e_test.go`:

```go
It("test", Label("id=b1f3572c-f1f0-4001-a4b6-97625206d9f9", "test=test"), func() {
    ...
})
```

Result `b1f3572c-f1f0-4001-a4b6-97625206d9f9-result.json`
```json
{
    ...
    "labels": [
        {
            "name": "id",
            "value": "b1f3572c-f1f0-4001-a4b6-97625206d9f9"
        },
        {
            "name": "test",
            "value": "test"
        }
    ],
    ...
}
```

#### Default labels

By default, each Allure test adds the labels `id=<uuid>`,`suite=<suite-name>`, and `epic=base`. Label `epic` you can change with 2 options:
1. Define the `It` test (high priority).
2. Add the flag `--epic` in CLI (low priority).

#### Mandatory labels

For your own goals, you can define a list of Ginkgo labels (flag `--mandatory_labels`), which must be in **ALL** `It` tests, like `featur`,`story`, etc. By default, it's only an `id`.

### Analyse the error issue

In general, Ginkgo starts like that.

```go
import (
	"testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRunner(t *testing.T) {
	suiteConfig, reporterConfig := GinkgoConfiguration()

	RegisterFailHandler(Fail)
	RunSpecs(t, SuiteName, suiteConfig, reporterConfig)
}
```

As you can see, we use the basic Ginkgo `Fail` handler, which indeed doesn't have a lot of really important information for us (for instance, expect and actual values in a Gomega assert function). However, for compatibility, it was decided to stay with this handler and parse explicit trace output. It's a bad approach, but it will close most test cases. For you, it means that if you find any problems with this functionality, please inform me in Issue and disable it with the flag `--analyze_errors`.

### Test description

If you check [the official Ginko documentation](https://onsi.github.io/ginkgo/#adding-specs-to-a-suite), you will see that Ginkgo `Describe + Context (second Describe) + It` form simple English sentences. `Categorizing books with more than 300 pages should be a novel`. That's a basic naming rule in tests. Therefore, I decided to use this approach to write down the default description of the test case in Allure. However, I also offer the opportunity to create your own description; just append an additional label to `It`: `description=<your describe>`.

## Usage

### CLI

Now, after reading [TestCaseID issue](#TestCaseID_issue) you grasp how to prepare your code for conversion. Below is a basic CLI run.

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

### Lib

You can also use the converter exactly in your Ginkgo code.

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

### Docker

```sh
# Get Ginkgo JSON report
ginkgo -r --keep-going -p --json-report=report.json ./tests/e2e/
# Create folder for Allure results
mkdir -p ./allure-results/
# Use CLI in docker
docker run -it -v ./:/results moon1706/ginkgo2allure:0.2.0 /results/report.json /results/allure-results
```

## Build

### Go code

```sh
make bin
```

### Docker image

```sh
docker build -t ginkgo2allure:v0.0.1 .
```

## Test

```sh
make coverage
```
