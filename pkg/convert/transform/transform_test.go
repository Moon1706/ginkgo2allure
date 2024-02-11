package transform_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Moon1706/ginkgo2allure/pkg/convert/report"
	"github.com/Moon1706/ginkgo2allure/pkg/convert/transform"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2/types"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/stretchr/testify/assert"
)

const (
	ginkgoReportFolderPath = "./ginkgo_reports_test"
	allureReportFolderPath = "./allure_reports_test"
)

func TestTransformAnalyzeEvents(t *testing.T) {
	var tests = []struct {
		fileName       string
		transformOpts  []transform.Opt
		haveError      bool
		checkStepsFunc func(assert.TestingT, []*allure.Step, error)
	}{{
		// Ginkgo code which analyzed in this test:
		//
		// It("test 1", Label("id=c57e2b09-901f-4991-a516-a22c8bb625d4", "story=story1"), func() {
		// 	By("plain 1")
		// 	By("nested 1", func() {
		// 		By("nested 2", func() {
		// 			Expect("1").To(Equal("1"))
		// 		})
		// 		By("plain 2")
		// 	})
		// })
		fileName: "basic_success",
		transformOpts: []transform.Opt{
			transform.WillAnalyzeErrors(false, false),
			transform.WithFilterEvents(func(event types.SpecEvent) bool {
				return false
			}),
		},
		haveError: false,
	}, {
		// Ginkgo code which analyzed in this test:
		// It("test 1", Label("id=c57e2b09-901f-4991-a516-a22c8bb625d4", "story=story1"), func() {
		// 	By("nested 1", func() {
		// 		Expect("1").To(Equal("1"))
		// 		By("plain 2")
		// 	})
		// 	Expect("1").To(Equal("2"))
		// })
		fileName:      "error_it_test",
		transformOpts: []transform.Opt{transform.WillAnalyzeErrors(true, true)},
		haveError:     false,
	}, {
		// Ginkgo code which analyzed in this test:
		// It("test 1", Label("id=c57e2b09-901f-4991-a516-a22c8bb625d4", "story=story1"), func() {
		// 	By("nested 1", func() {
		// 		Expect("1").To(Equal("2"))
		// 		By("plain 2")
		// 	})
		// })
		fileName:      "error_nested_step",
		transformOpts: []transform.Opt{transform.WillAnalyzeErrors(true, true)},
		haveError:     false,
	}, {
		// Ginkgo code which analyzed in this test:
		// It("test 1", Label("id=c57e2b09-901f-4991-a516-a22c8bb625d4", "story=story1"), func() {
		// 	By("nested 1", func() {
		// 		By("plain 2")
		// 		Expect("1").To(Equal("2"))
		// 	})
		// })
		fileName:      "error_unnested_by",
		transformOpts: []transform.Opt{transform.WillAnalyzeErrors(true, true)},
		haveError:     false,
	}, {
		// Ginkgo code which analyzed in this test. Also made one mistake in FullStackTrace field (66l).
		// BeforeEach(func() {
		// 	Expect("1").To(Equal("1"))
		// })
		// It("test 1", Label("id=c57e2b09-901f-4991-a516-a22c8bb625d4", "story=story1"), func() {
		// 	Expect("1").To(Equal("2"))
		// })
		fileName:      "error_trace_mistake",
		transformOpts: []transform.Opt{transform.WillAnalyzeErrors(true, true)},
		haveError:     true,
	}, {
		// Ginkgo code which analyzed in this test:
		// It("test 1", Label("id=c57e2b09-901f-4991-a516-a22c8bb625d4", "story=story1"), func() {
		// 	Expect("1").To(Equal("1"))
		// })
		// It("test 1", Label("id=c57e2b09-901f-4991-a516-a22c8bb625d5", "story=story1"), func() {
		// 	Expect("1").To(Equal("2"))
		// })
		fileName:      "error_two_it",
		transformOpts: []transform.Opt{transform.WillAnalyzeErrors(true, true)},
		haveError:     false,
	}}

	for _, tt := range tests {
		ginkoReportPath := filepath.Join(ginkgoReportFolderPath, fmt.Sprintf("%s.json", tt.fileName))
		ginkgoReports, err := readReports[[]types.Report](ginkoReportPath)
		assert.Empty(t, err, fmt.Sprintf("no error during ginkgo report unmarshaling (%s)", tt.fileName))
		for _, ginkgoReport := range ginkgoReports {
			tr := transform.NewTransform(tt.transformOpts...)
			for _, ginkgoSpecReport := range ginkgoReport.SpecReports {
				err := tr.AnalyzeEvents(ginkgoSpecReport.SpecEvents, ginkgoSpecReport.Failure)
				if tt.haveError {
					assert.Error(t, err, fmt.Sprintf("got expected error during analyze (%s)", tt.fileName))
				} else {
					assert.Empty(t, err, fmt.Sprintf("no error during analyze (%s)", tt.fileName))
				}
				steps := tr.GetAllureSteps()

				ls := report.NewLabelScraper(ginkgoSpecReport.LeafNodeText, ginkgoSpecReport.LeafNodeLabels)
				id, err := ls.GetID(uuid.New())
				assert.Empty(t, err, "id got correctly")

				allureReportPath := filepath.Join(allureReportFolderPath, tt.fileName, fmt.Sprintf("%s-result.json", id))
				allureResult, err := readReports[allure.Result](allureReportPath)
				assert.Empty(t, err, fmt.Sprintf("no error during allure report unmarshaling (%s)", tt.fileName))
				assert.Equal(t, allureResult.Steps, steps, "got expected steps")
			}
		}
	}
}

func readReports[T any](filePath string) (T, error) {
	out := new(T)
	file, err := os.ReadFile(filePath)
	if err != nil {
		return *out, err
	}
	err = json.Unmarshal(file, out)
	return *out, err
}
