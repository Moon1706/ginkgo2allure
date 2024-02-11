package report_test

import (
	"fmt"
	"testing"

	"github.com/Moon1706/ginkgo2allure/pkg/convert/report"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/stretchr/testify/assert"
)

func TestLabelScraperGetAllTestCaseLabels(t *testing.T) {
	var tests = []struct {
		name           string
		leafNodeLabels []string
		testCaseLabels map[string]string
		scraperOpt     []report.LabelsScraperOpt
	}{{
		name:           "empty",
		leafNodeLabels: []string{},
		testCaseLabels: map[string]string{},
		scraperOpt:     []report.LabelsScraperOpt{report.WithEpic("")},
	}, {
		name:           "default values",
		leafNodeLabels: []string{},
		testCaseLabels: map[string]string{report.EpicLabelName: report.DefaultEpic},
		scraperOpt:     []report.LabelsScraperOpt{},
	}, {
		name:           "incorrect label without spliter",
		leafNodeLabels: []string{"incorrect-label"},
		testCaseLabels: map[string]string{},
		scraperOpt:     []report.LabelsScraperOpt{report.WithEpic("")},
	}, {
		name:           "incorrect label exist spliter, but many parts",
		leafNodeLabels: []string{fmt.Sprintf("multi%[1]sincorrect%[1]slabel", report.DefaultLabelSpliter)},
		testCaseLabels: map[string]string{},
		scraperOpt:     []report.LabelsScraperOpt{report.WithEpic("")},
	}, {
		name:           "correct label",
		leafNodeLabels: []string{fmt.Sprintf("correct%slabel", report.DefaultLabelSpliter)},
		testCaseLabels: map[string]string{"correct": "label"},
		scraperOpt:     []report.LabelsScraperOpt{report.WithEpic("")},
	}, {
		name:           "correct label with suite and epic labels",
		leafNodeLabels: []string{fmt.Sprintf("correct%slabel", report.DefaultLabelSpliter)},
		testCaseLabels: map[string]string{"correct": "label",
			report.EpicLabelName:  "test",
			report.SuiteLabelName: "test"},
		scraperOpt: []report.LabelsScraperOpt{report.WithSuiteName("test"), report.WithEpic("test")},
	}, {
		name:           "change label spliter",
		leafNodeLabels: []string{"correct:label"},
		testCaseLabels: map[string]string{"correct": "label",
			report.EpicLabelName:  "test",
			report.SuiteLabelName: "test"},
		scraperOpt: []report.LabelsScraperOpt{report.WithSuiteName("test"),
			report.WithEpic("test"),
			report.WithLabelSpliter(":")},
	}}

	for _, tt := range tests {
		lb := report.NewLabelScraper(tt.leafNodeLabels, tt.scraperOpt...)
		tc := lb.GetTestCaseLabels()
		assert.Equal(t, tc, tt.testCaseLabels, tt.name)
	}
}

func TestLabelScraperCheckMandatoryLabels(t *testing.T) {
	mandatoryLabelsLabels := []string{"correct"}
	lb := report.NewLabelScraper([]string{fmt.Sprintf("correct%slabel", report.DefaultLabelSpliter)})
	err := lb.CheckMandatoryLabels(mandatoryLabelsLabels)
	assert.Empty(t, err, "correct label was found in madatory labels")

	lb = report.NewLabelScraper([]string{fmt.Sprintf("incorrect%slabel", report.DefaultLabelSpliter)})
	err = lb.CheckMandatoryLabels(mandatoryLabelsLabels)
	assert.Error(t, err, "lables wasn't found in madatory labels")
}

func TestLabelScraperCreateAllureLabels(t *testing.T) {
	lb := report.NewLabelScraper([]string{fmt.Sprintf("correct%slabel", report.DefaultLabelSpliter)},
		report.WithEpic(""))
	allureLabels := lb.CreateAllureLabels()
	assert.Equal(t, []*allure.Label{{
		Name:  "correct",
		Value: "label",
	}}, allureLabels, "correctly convert map of strings to allure labels")
}

func TestLabelGetID(t *testing.T) {
	var tests = []struct {
		name       string
		labelName  string
		labelValue string
		id         uuid.UUID
	}{{
		name:       "incorrect label name and value",
		labelName:  "incorrect",
		labelValue: "incorrect",
		id:         uuid.UUID{},
	}, {
		name:       "incorrect label value",
		labelName:  report.IDLabelName,
		labelValue: "incorrect",
		id:         uuid.UUID{},
	}, {
		name:       "correct id",
		labelName:  report.IDLabelName,
		labelValue: "ad7583dc-0e3d-4640-a020-567452d84886",
		id:         uuid.MustParse("ad7583dc-0e3d-4640-a020-567452d84886"),
	}}

	for _, tt := range tests {
		lb := report.NewLabelScraper([]string{fmt.Sprintf("%s%s%s", tt.labelName,
			report.DefaultLabelSpliter, tt.labelValue)})
		id, _ := lb.GetID()
		assert.Equal(t, tt.id, id, tt.name)
	}
}

func TestLabelGetDescription(t *testing.T) {
	var tests = []struct {
		name               string
		descriptionLabel   string
		description        string
		defaultDescription string
	}{{
		name:               "default description",
		descriptionLabel:   "",
		description:        "test",
		defaultDescription: "test",
	}, {
		name: "change description",
		descriptionLabel: fmt.Sprintf("%s%stest", report.DescriptionLabelName,
			report.DefaultLabelSpliter),
		description:        "test",
		defaultDescription: "",
	}}

	for _, tt := range tests {
		lb := report.NewLabelScraper([]string{tt.descriptionLabel})
		description := lb.GetDescription(tt.defaultDescription)
		assert.Equal(t, tt.description, description, tt.name)
	}
}
