package report

import (
	"errors"
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

	DefaultLabelSpliter = "="
	DefaultEpic         = "base"
	DefaultSuiteName    = ""

	CorrectCountLabelsParts = 2
)

var (
	ErrNoID = errors.New("test doesn't contain allure id in labels")
)

type (
	DefaultLabelsScraper struct {
		epic           string
		suiteName      string
		testCaseLabels map[string]string
		labelSpliter   string
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

func NewLabelScraper(leafNodeLabels []string, scraperOptions ...LabelsScraperOpt) *DefaultLabelsScraper {
	scraper := &DefaultLabelsScraper{
		epic:         DefaultEpic,
		suiteName:    DefaultSuiteName,
		labelSpliter: DefaultLabelSpliter,
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

func (ls *DefaultLabelsScraper) GetID() (uuid.UUID, error) {
	id, ok := ls.testCaseLabels[IDLabelName]
	if !ok {
		return uuid.UUID{}, ErrNoID
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
