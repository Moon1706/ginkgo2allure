package convert_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Moon1706/ginkgo2allure/pkg/convert"
	"github.com/Moon1706/ginkgo2allure/pkg/convert/parser"
	"github.com/Moon1706/ginkgo2allure/pkg/convert/report"
	"github.com/onsi/ginkgo/v2/types"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/stretchr/testify/assert"
)

var (
	errTest = errors.New("test error")
)

type (
	mockReport    struct{}
	mockTransform struct {
		Err error
	}
	mockFileManager struct {
		SaveErr error
	}
)

func (m mockReport) GenerateAllureReport(_ []*allure.Step) (allure.Result, error) {
	return allure.Result{}, nil
}
func (m mockReport) SetLabelsScraper(_ report.LabelScraper) {}

func (m mockTransform) AnalyzeEvents(_ types.SpecEvents, _ types.Failure) error {
	return m.Err
}
func (m mockTransform) GetAllureSteps() []*allure.Step {
	return []*allure.Step{}
}
func (m mockFileManager) SaveJSONResult(_ allure.Result) error {
	return m.SaveErr
}

func TestConvertGinkgoToAllureReport(t *testing.T) {
	ginkgoReports := []types.Report{{
		SuiteDescription: "test",
		SpecReports: types.SpecReports{types.SpecReport{
			LeafNodeType: types.NodeTypeBeforeAll,
		}, types.SpecReport{
			LeafNodeType: types.NodeTypeIt,
		}},
	}}
	var tests = []struct {
		name       string
		createFunc parser.CreationFunc
		results    []allure.Result
		err        error
	}{{
		name: "correct",
		createFunc: func(specReport types.SpecReport, _ parser.Config) (*parser.Parser, error) {
			return parser.NewParser(specReport, mockTransform{Err: nil}, nil, mockReport{}), nil
		},
		results: []allure.Result{{}},
		err:     nil,
	}, {
		name: "error in parser creation function",
		createFunc: func(specReport types.SpecReport, _ parser.Config) (*parser.Parser, error) {
			return &parser.Parser{}, errTest
		},
		results: []allure.Result{},
		err:     errTest,
	}, {
		name: "error in parser",
		createFunc: func(specReport types.SpecReport, _ parser.Config) (*parser.Parser, error) {
			return parser.NewParser(specReport, mockTransform{Err: errTest}, nil, mockReport{}), nil
		},
		results: []allure.Result{},
		err:     errTest,
	}}

	for _, tt := range tests {
		results, err := convert.GinkgoToAllureReport(ginkgoReports, tt.createFunc, parser.Config{})
		assert.Equal(t, tt.err, err, fmt.Sprintf("got expected error (%s)", tt.name))
		assert.Equal(t, tt.results, results, fmt.Sprintf("got expected results (%s)", tt.name))
	}
}

func TestConvertPrintAllureReports(t *testing.T) {
	var tests = []struct {
		name            string
		mockFileManager mockFileManager
		errs            []error
	}{{
		name: "correct",
		mockFileManager: mockFileManager{
			SaveErr: nil,
		},
		errs: []error{},
	}, {
		name: "wrong",
		mockFileManager: mockFileManager{
			SaveErr: errTest,
		},
		errs: []error{errTest},
	}}

	for _, tt := range tests {
		errs := convert.PrintAllureReports([]allure.Result{{}}, tt.mockFileManager)
		assert.Equal(t, tt.errs, errs, fmt.Sprintf("got expected errors (%s)", tt.name))
	}
}
