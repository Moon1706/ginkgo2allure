package report_test

import (
	"fmt"
	"testing"

	"github.com/Moon1706/ginkgo2allure/pkg/convert/report"
	"github.com/onsi/ginkgo/v2/types"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAllureReport(t *testing.T) {
	id := "8791ccdd-83c6-4333-b589-f3a7822166f5"
	var tests = []struct {
		name       string
		specReport types.SpecReport
		reportOpt  []report.Opt
		result     allure.Result
		existError bool
	}{{
		name: "default settings only with id label that check",
		specReport: types.SpecReport{
			LeafNodeText: "test",
			LeafNodeLabels: []string{fmt.Sprintf("%s%s%s", report.IDLabelName,
				report.DefaultLabelSpliter, id)},
		},
		reportOpt: []report.Opt{},
		result: allure.Result{
			Name:     "test",
			FullName: id,
		},
		existError: false,
	}, {
		name: "block due to mandatory labels that check",
		specReport: types.SpecReport{
			LeafNodeText: "test",
			LeafNodeLabels: []string{fmt.Sprintf("%s%s%s", report.IDLabelName,
				report.DefaultLabelSpliter, id)},
		},
		reportOpt: []report.Opt{report.WithMandatoryLabels([]string{"test"})},
		result: allure.Result{
			Name:     "",
			FullName: "",
		},
		existError: true,
	}, {
		name: "block due to incorrect uuid that check",
		specReport: types.SpecReport{
			LeafNodeText: "test",
			LeafNodeLabels: []string{fmt.Sprintf("%s%sincorrect-uuid", report.IDLabelName,
				report.DefaultLabelSpliter)},
		},
		reportOpt: []report.Opt{},
		result: allure.Result{
			Name:     "",
			FullName: "",
		},
		existError: true,
	}, {
		name: "analyze failure field that check",
		specReport: types.SpecReport{
			LeafNodeText: "test",
			LeafNodeLabels: []string{fmt.Sprintf("%s%s%s", report.IDLabelName,
				report.DefaultLabelSpliter, id)},
			Failure: types.Failure{
				Message: "test",
				TimelineLocation: types.TimelineLocation{
					Order: 1,
				},
			},
		},
		reportOpt: []report.Opt{},
		result: allure.Result{
			Name:     "test",
			FullName: id,
			StatusDetails: allure.StatusDetail{
				Message: "test",
			},
		},
		existError: false,
	}}

	for _, tt := range tests {
		r := report.NewReport(tt.specReport, tt.reportOpt...)
		result, err := r.GenerateAllureReport([]*allure.Step{})
		if tt.existError {
			assert.Error(t, err, "expected error during allure report creation")
		} else {
			assert.Empty(t, err, "allure report was created successful")
		}
		assert.Equal(t, tt.result.Name, result.Name, tt.name+" name")
		assert.Equal(t, tt.result.FullName, result.FullName, tt.name+" full name")
		assert.Equal(t, tt.result.StatusDetails.Message,
			result.StatusDetails.Message, tt.name+" status message")
	}
}
