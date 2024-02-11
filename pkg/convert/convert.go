package convert

import (
	fmngr "github.com/Moon1706/ginkgo2allure/pkg/convert/file_manager"
	"github.com/Moon1706/ginkgo2allure/pkg/convert/parser"
	"github.com/Moon1706/ginkgo2allure/pkg/convert/report"
	"github.com/onsi/ginkgo/v2/types"
	"github.com/ozontech/allure-go/pkg/allure"
)

func GinkgoToAllureReport(ginkgoReports []types.Report, parserCreation parser.CreationFunc,
	config parser.Config) ([]allure.Result, error) {
	results := []allure.Result{}
	for _, ginkgoReport := range ginkgoReports {
		suiteName := ginkgoReport.SuiteDescription
		for _, specReport := range ginkgoReport.SpecReports {
			if specReport.LeafNodeType != types.NodeTypeIt {
				continue
			}
			config.LabelsScraperOpts = append(config.LabelsScraperOpts, report.WithSuiteName(suiteName))
			p, err := parserCreation(specReport, config)
			if err != nil {
				return results, err
			}
			result, err := p.GetAllureReport()
			if err != nil {
				return results, err
			}
			results = append(results, result)
		}
	}
	return results, nil
}

func PrintAllureReports(results []allure.Result, fm fmngr.FileManager) []error {
	errs := []error{}
	for _, result := range results {
		err := fm.SaveJSONResult(result)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
