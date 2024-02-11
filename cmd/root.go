package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/Moon1706/ginkgo2allure/internal/app"
	"github.com/Moon1706/ginkgo2allure/pkg/convert/parser"
	"github.com/Moon1706/ginkgo2allure/pkg/convert/report"
	"github.com/Moon1706/ginkgo2allure/pkg/convert/transform"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	CountArgs              = 2
	CountAnalyzeErrorsArgs = 2

	FlagEpic            = "epic"
	FlagLabelSeparator  = "label_separator"
	FlagMandatoryLabels = "mandatory_labels"
	FlagAnalyzeErrors   = "analyze_errors"
	FlagLogLevel        = "log_level"
)

var (
	logLevel string
)

var rootCmd = &cobra.Command{
	Use:   "ginkgo2allure ./ginkgo-report.json ./save/allure/reports/folder/path/",
	Short: "Convert Ginkgo report to Allure report",
	Long: `Prototype of a tool that converts Ginkgo JSON reports to Allure JSON reports
in a separate folder allure-results.`,
	Args: cobra.MatchAll(cobra.ExactArgs(CountArgs), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		config := parser.Config{}
		logger, err := buildLogger(logLevel)
		if err != nil {
			panic(err)
		}

		epic, err := cmd.Flags().GetString(FlagEpic)
		if err == nil && epic != "" {
			config.LabelsScraperOpts = append(config.LabelsScraperOpts, report.WithEpic(epic))
		}
		labelSpliter, err := cmd.Flags().GetString(FlagLabelSeparator)
		if err == nil && labelSpliter != "" {
			config.LabelsScraperOpts = append(config.LabelsScraperOpts, report.WithLabelSpliter(labelSpliter))
		}
		mandatoryLabels, err := cmd.Flags().GetStringSlice(FlagMandatoryLabels)
		if err == nil && len(mandatoryLabels) != 0 {
			config.ReportOpts = append(config.ReportOpts, report.WithMandatoryLabels(mandatoryLabels))
		}
		analyzeErrors, err := cmd.Flags().GetBool(FlagAnalyzeErrors)
		if err == nil {
			config.TransformOpts = append(config.TransformOpts, transform.WillAnalyzeErrors(analyzeErrors, analyzeErrors))
		}
		app.StartConvertion(args[0], args[1], config, logger)
	},
}

func ExecuteContext(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP(FlagEpic, "e", report.DefaultEpic, "epic name")
	rootCmd.Flags().String(FlagLabelSeparator, report.DefaultLabelSpliter, "labels separator")
	rootCmd.Flags().StringSlice(FlagMandatoryLabels, []string{report.IDLabelName}, "allure mandatory labels")
	rootCmd.Flags().Bool(FlagAnalyzeErrors, true, "will analyze test fails in Ginkgo report or not")
	rootCmd.PersistentFlags().StringVarP(&logLevel, FlagLogLevel, "l", "info", "log level")
}

func buildLogger(logLevel string) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(logLevel)
	if err != nil {
		lvl = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}

	encconf := zap.NewProductionEncoderConfig()
	encconf.TimeKey = "@timestamp"
	encconf.EncodeTime = zapcore.RFC3339TimeEncoder

	conf := zap.NewProductionConfig()
	conf.EncoderConfig = encconf
	conf.Level = lvl

	return conf.Build()
}
