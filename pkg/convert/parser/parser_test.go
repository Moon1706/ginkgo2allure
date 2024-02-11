package parser_test

import (
	"errors"
	"fmt"
	"testing"

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

func TestParserGetAllureReport(t *testing.T) {
	var tests = []struct {
		name          string
		mockReport    parser.Reporter
		mockTransform parser.Transformer
		result        allure.Result
		err           error
	}{{
		name:       "correct",
		mockReport: mockReport{},
		mockTransform: mockTransform{
			Err: nil,
		},
		result: allure.Result{},
		err:    nil,
	}, {
		name:       "wrong",
		mockReport: mockReport{},
		mockTransform: mockTransform{
			Err: errTest,
		},
		result: allure.Result{},
		err:    errTest,
	}}

	for _, tt := range tests {
		p := parser.NewParser(types.SpecReport{}, tt.mockTransform, nil, tt.mockReport)
		result, err := p.GetAllureReport()
		assert.Equal(t, tt.err, err, fmt.Sprintf("got expected error (%s)", tt.name))
		assert.Equal(t, tt.result, result, fmt.Sprintf("got expected result (%s)", tt.name))
	}
}

func TestNewDefaultParser(t *testing.T) {
	_, err := parser.NewDefaultParser(types.SpecReport{}, parser.Config{})
	assert.Empty(t, err, "no error")
}
