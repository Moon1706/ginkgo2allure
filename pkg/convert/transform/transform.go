package transform

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/onsi/ginkgo/v2/types"
	"github.com/ozontech/allure-go/pkg/allure"
)

const (
	DefaultAnalyzeErrors         = true
	DefaultGetErrorDuringAlalyze = true
)

type Node struct {
	BeginEvent types.SpecEvent
	EndEvent   types.SpecEvent
}

type TraceFile struct {
	FileName   string
	LineNumber int
}

type (
	DefaultTransform struct {
		analyzeErrors         bool
		getErrorDuringAlalyze bool
		filterEvents          FilterEvents
		nodes                 []Node
		errNode               Node
		steps                 []*allure.Step
	}
	Opt          func(o *DefaultTransform)
	FilterEvents func(event types.SpecEvent) bool
)

func WillAnalyzeErrors(analyzeErrors, getErrorDuringAlalyze bool) Opt {
	return func(o *DefaultTransform) {
		o.analyzeErrors = analyzeErrors
		o.getErrorDuringAlalyze = getErrorDuringAlalyze
	}
}

func WithFilterEvents(filter FilterEvents) Opt {
	return func(o *DefaultTransform) {
		o.filterEvents = filter
	}
}

func NewTransform(opts ...Opt) *DefaultTransform {
	filterSuiteAndEachEvents := func(event types.SpecEvent) bool {
		return event.NodeType != types.NodeTypeInvalid &&
			event.NodeType != types.NodeTypeIt
	}
	t := &DefaultTransform{
		analyzeErrors:         DefaultAnalyzeErrors,
		getErrorDuringAlalyze: DefaultGetErrorDuringAlalyze,
		filterEvents:          filterSuiteAndEachEvents,
	}
	for _, o := range opts {
		o(t)
	}
	return t
}

func (t *DefaultTransform) AnalyzeEvents(events types.SpecEvents, failure types.Failure) error {
	t.nodes = t.findNodes(events)
	if failure.Message != "" && t.analyzeErrors {
		errNode, err := t.findErrorNode(t.nodes, failure)
		if err != nil && t.getErrorDuringAlalyze {
			return err
		}
		t.errNode = errNode
	}
	t.steps = t.getNestedSteps(t.nodes, t.errNode)
	return nil
}

func (t *DefaultTransform) GetAllureSteps() []*allure.Step {
	return t.steps
}

func (t *DefaultTransform) findNodes(events types.SpecEvents) (nodes []Node) {
	for _, event := range events {
		if t.filterEvents(event) {
			continue
		}
		if event.SpecEventType != types.SpecEventNodeStart &&
			event.SpecEventType != types.SpecEventByStart {
			continue
		}
		haveEndEvent := false
		for _, endEvent := range events {
			if endEvent.CodeLocation.FileName == event.CodeLocation.FileName &&
				endEvent.CodeLocation.LineNumber == event.CodeLocation.LineNumber &&
				endEvent.TimelineLocation.Order != event.TimelineLocation.Order {
				haveEndEvent = true
				nodes = append(nodes, Node{
					BeginEvent: event,
					EndEvent:   endEvent,
				})
			}
		}
		if !haveEndEvent {
			nodes = append(nodes, Node{
				BeginEvent: event,
				EndEvent:   event,
			})
		}
	}
	return nodes
}

func (t *DefaultTransform) getNestedSteps(nodes []Node, errNode Node) []*allure.Step {
	finalSteps := []*allure.Step{}
	steps := make([]*allure.Step, 0, len(nodes))
	for i := range nodes {
		stepName := nodes[i].BeginEvent.Message
		if nodes[i].BeginEvent.NodeType != types.NodeTypeInvalid {
			stepName = fmt.Sprintf("[%s] %s", nodes[i].BeginEvent.NodeType.String(), stepName)
		}
		stepStatus := allure.Passed
		if nodes[i].BeginEvent.CodeLocation.FileName == errNode.BeginEvent.CodeLocation.FileName &&
			nodes[i].BeginEvent.CodeLocation.LineNumber == errNode.BeginEvent.CodeLocation.LineNumber {
			stepStatus = allure.Failed
		}
		step := &allure.Step{
			Name:   stepName,
			Status: stepStatus,
			Start:  nodes[i].BeginEvent.TimelineLocation.Time.UnixMilli(),
			Stop:   nodes[i].EndEvent.TimelineLocation.Time.UnixMilli(),
		}
		steps = append(steps, step)
		for j := i; j >= 0; j-- {
			if nodes[i].BeginEvent.TimelineLocation.Order > nodes[j].BeginEvent.TimelineLocation.Order &&
				nodes[i].EndEvent.TimelineLocation.Order < nodes[j].EndEvent.TimelineLocation.Order {
				steps[j].Steps = append(steps[j].Steps, step)
				break
			}
			if j == 0 {
				finalSteps = append(finalSteps, step)
			}
		}
	}
	return finalSteps
}

func (t *DefaultTransform) findErrorNode(nodes []Node, failure types.Failure) (Node, error) {
	traceFiles, err := t.getTraceFiles(failure)
	if err != nil {
		return Node{}, err
	}
	lastEvent := len(nodes) - 1
	for i := lastEvent; i >= 0; i-- {
		if t.stepInTrace(nodes[i], traceFiles) {
			return nodes[i], nil
		}
		if nodes[i].BeginEvent.TimelineLocation.Order == nodes[i].EndEvent.TimelineLocation.Order &&
			i == lastEvent {
			for j := i; j >= 0; j-- {
				if nodes[j].EndEvent.TimelineLocation.Order > nodes[i].EndEvent.TimelineLocation.Order {
					if t.stepInTrace(nodes[j], traceFiles) {
						return nodes[i], nil
					}
				}
			}
		}
	}
	return t.findErrorRootNode(nodes, traceFiles), nil
}

func (t *DefaultTransform) getTraceFiles(failure types.Failure) (traceFiles []TraceFile, err error) {
	rawTraceLines := strings.Split(failure.Location.FullStackTrace, "\n")
	for i := 1; i < len(rawTraceLines); i += 2 {
		clearTraceLine := strings.Split(rawTraceLines[i][1:], " +0x")[0]
		tl := strings.Split(clearTraceLine, ":")
		line, errInt := strconv.Atoi(tl[1])
		if errInt != nil {
			return traceFiles, errInt
		}
		traceFiles = append(traceFiles, TraceFile{
			FileName:   tl[0],
			LineNumber: line,
		})
	}
	return traceFiles, err
}

func (t *DefaultTransform) stepInTrace(node Node, traceFiles []TraceFile) bool {
	for _, tf := range traceFiles {
		if node.BeginEvent.CodeLocation.FileName == tf.FileName &&
			node.BeginEvent.CodeLocation.LineNumber == tf.LineNumber {
			return true
		}
	}
	return false
}

func (t *DefaultTransform) findErrorRootNode(nodes []Node, traceFiles []TraceFile) (errNode Node) {
	lastTraceFile := traceFiles[len(traceFiles)-1]
	for _, node := range nodes {
		if node.BeginEvent.CodeLocation.FileName == lastTraceFile.FileName &&
			node.BeginEvent.NodeType == types.NodeTypeIt &&
			node.BeginEvent.CodeLocation.LineNumber < lastTraceFile.LineNumber {
			errNode = node
		}
	}
	return
}
