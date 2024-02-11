package parser

import (
	"github.com/Moon1706/ginkgo2allure/pkg/convert/report"
	"github.com/Moon1706/ginkgo2allure/pkg/convert/transform"
	"github.com/onsi/ginkgo/v2/types"
	"github.com/ozontech/allure-go/pkg/allure"
)

type (
	Reporter interface {
		GenerateAllureReport([]*allure.Step) (allure.Result, error)
		SetLabelsScraper(ls report.LabelScraper)
	}
	Transformer interface {
		AnalyzeEvents(types.SpecEvents, types.Failure) error
		GetAllureSteps() []*allure.Step
	}
	Parser struct {
		Transformer  Transformer
		Reporter     Reporter
		LabelScraper report.LabelScraper
		specReport   types.SpecReport
	}
	Config struct {
		TransformOpts     []transform.Opt
		LabelsScraperOpts []report.LabelsScraperOpt
		ReportOpts        []report.Opt
	}
	CreationFunc func(types.SpecReport, Config) (*Parser, error)
)

func NewParser(specReport types.SpecReport, transformer Transformer,
	ls report.LabelScraper, reporter Reporter) *Parser {
	reporter.SetLabelsScraper(ls)
	return &Parser{
		Transformer:  transformer,
		Reporter:     reporter,
		LabelScraper: ls,
		specReport:   specReport,
	}
}

func NewDefaultParser(specReport types.SpecReport, config Config) (*Parser, error) {
	t := transform.NewTransform(config.TransformOpts...)
	ls := report.NewLabelScraper(specReport.LeafNodeLabels, config.LabelsScraperOpts...)
	r := report.NewReport(specReport, config.ReportOpts...)
	return NewParser(specReport, t, ls, r), nil
}

func (p *Parser) GetAllureReport() (allure.Result, error) {
	err := p.Transformer.AnalyzeEvents(p.specReport.SpecEvents, p.specReport.Failure)
	if err != nil {
		return allure.Result{}, err
	}
	steps := p.Transformer.GetAllureSteps()
	return p.Reporter.GenerateAllureReport(steps)
}
