package report

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
)

const (
	EpicLabelName        = "epic"
	SuiteLabelName       = "suite"
	FeatureLabelName     = "feature"
	IDLabelName          = "id"
	DescriptionLabelName = "description"

	DefaultLabelSpliter   = "="
	DefaultEpic           = "base"
	DefaultSuiteName      = ""
	DefaultAutoGenerateID = false

	CorrectCountLabelsParts = 2
)

type (
	DefaultLabelsScraper struct {
		testName       string
		epic           string
		suiteName      string
		testCaseLabels map[string]string
		labelSpliter   string
		autogenID      bool
	}
	LabelsScraperOpt func(o *DefaultLabelsScraper)
)

func WithLabelSpliter(splitter string) LabelsScraperOpt {
	return func(o *DefaultLabelsScraper) {
		o.labelSpliter = splitter
	}
}

func WithEpic(epic string) LabelsScraperOpt {
	return func(o *DefaultLabelsScraper) {
		o.epic = epic
	}
}

func WithSuiteName(suiteName string) LabelsScraperOpt {
	return func(o *DefaultLabelsScraper) {
		o.suiteName = suiteName
	}
}

func WillAutoGenerateID(autogen bool) LabelsScraperOpt {
	return func(o *DefaultLabelsScraper) {
		o.autogenID = autogen
	}
}

func NewLabelScraper(leafNodeText string, leafNodeLabels []string,
	scraperOptions ...LabelsScraperOpt) *DefaultLabelsScraper {
	scraper := &DefaultLabelsScraper{
		testName:     leafNodeText,
		epic:         DefaultEpic,
		suiteName:    DefaultSuiteName,
		labelSpliter: DefaultLabelSpliter,
		autogenID:    DefaultAutoGenerateID,
	}
	for _, o := range scraperOptions {
		o(scraper)
	}
	testCaseLabels := scraper.getAllTestCaseLabels(leafNodeLabels)
	if scraper.epic != "" {
		testCaseLabels[EpicLabelName] = scraper.epic
	}
	if scraper.suiteName != "" {
		testCaseLabels[SuiteLabelName] = scraper.suiteName
	}
	scraper.testCaseLabels = testCaseLabels
	return scraper
}

func (ls *DefaultLabelsScraper) GetTestCaseLabels() map[string]string {
	return ls.testCaseLabels
}

func (ls *DefaultLabelsScraper) getAllTestCaseLabels(labels []string) map[string]string {
	labelsMap := map[string]string{}
	for _, label := range labels {
		labelKeyValue := strings.Split(label, ls.labelSpliter)
		if len(labelKeyValue) != CorrectCountLabelsParts {
			continue
		}
		labelsMap[labelKeyValue[0]] = labelKeyValue[1]
	}
	return labelsMap
}

func (ls *DefaultLabelsScraper) CheckMandatoryLabels(mandatoryLabels []string) error {
	for _, mandatoryLabel := range mandatoryLabels {
		if _, ok := ls.testCaseLabels[mandatoryLabel]; !ok {
			return fmt.Errorf("doesn't exist mandatory label: %s", mandatoryLabel)
		}
	}
	return nil
}

func (ls *DefaultLabelsScraper) CreateAllureLabels() (labels []*allure.Label) {
	for key, value := range ls.testCaseLabels {
		labels = append(labels, &allure.Label{
			Name:  key,
			Value: value,
		})
	}
	return labels
}

func (ls *DefaultLabelsScraper) GetID(defaultID uuid.UUID) (uuid.UUID, error) {
	id, ok := ls.testCaseLabels[IDLabelName]
	if !ok {
		if ls.autogenID {
			return defaultID, nil
		}
		return uuid.UUID{}, fmt.Errorf("test with name `%s` doesn't contain UUID", ls.testName)
	}
	allureUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, err
	}
	return allureUUID, nil
}

func (ls *DefaultLabelsScraper) GetDescription(defaultDescription string) string {
	description, ok := ls.testCaseLabels[DescriptionLabelName]
	if !ok {
		description = defaultDescription
	}
	return description
}
