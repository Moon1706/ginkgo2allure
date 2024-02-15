package app

import (
	"encoding/json"
	"os"

	"github.com/Moon1706/ginkgo2allure/pkg/convert"
	fmngr "github.com/Moon1706/ginkgo2allure/pkg/convert/file_manager"
	"github.com/Moon1706/ginkgo2allure/pkg/convert/parser"
	"github.com/onsi/ginkgo/v2/types"
	"go.uber.org/zap"
)

func StartConvertion(ginkgoReportFile, allureReportsFolder string, config parser.Config, logger *zap.Logger) {
	sugar := logger.Sugar()

	file, err := os.ReadFile(ginkgoReportFile)
	if err != nil {
		sugar.Fatal("Error reading file ", ginkgoReportFile)
	}

	var ginkgoReport []types.Report
	err = json.Unmarshal(file, &ginkgoReport)
	if err != nil {
		sugar.Fatal("Error unmarshaling file ", ginkgoReportFile)
	}

	allureReports, err := convert.GinkgoToAllureReport(ginkgoReport, parser.NewDefaultParser,
		config)
	if err != nil {
		sugar.Fatal("Error converting report ", ginkgoReportFile, " ", err)
	}

	fileManager := fmngr.NewFileManager(allureReportsFolder)
	errs := convert.PrintAllureReports(allureReports, fileManager)
	for _, err := range errs {
		sugar.Error(err)
	}
	if len(errs) != 0 {
		sugar.Fatal("Saving errors")
	}
}
