package report

import (
	"bytes"
	"path"
	"runtime"
	html_template "text/template"
	text_template "text/template"
	"time"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("test")

// Status The test case status
type Status string

// These are the valid Test Status of test summary.
const (
	PassStatus   Status = "PASS"
	FailStatus   Status = "FAIL"
	ErrorStatus  Status = "ERROR"
	SkipStatus   Status = "SKIP"
	CancelStatus Status = "CANCEL"
)

// TestDesc .
type TestDesc struct {
	Tester       string
	Version      string
	StartTime    string
	EndTime      string
	ElapsedTime  string
	Summary      string
	DplImage     string
	MasterIPs    string
	TestLocation string
	ReportPath   string
	Command      string
}

// TestResult test result struct for tc each iteration
type TestResult struct {
	JobNr     int    // tc exec number
	Name      string // tc name
	Desc      string // tc description
	Start     time.Time
	End       time.Time
	Elapsed   string // string of time.Duration
	Error     error  // error
	Output    string // Contents for output to reports
	Status    Status // Test case status
	Iteration int    // The current exec iteration number
	Summary   string // [ Status ]- TestCase - Iteration: 1 - Elapsed Time: 1m10s
}

// Node  test env nodes info
type Node struct {
	Name      string
	Status    string
	IPAddress string
	Roles     string
	User      string
	Password  string
}

// HTMLResult .
type HTMLResult struct {
	FromUser string
	ToUsers  string
	TimeDate string
	Content  string

	Title    string // Test Report title
	Desc     string // Test Description, string of TestDesc{}
	Pass     int    // Test Case Pass number
	Fail     int    // Test Case Fail number
	Error    int    // Test Case Error number
	Skip     int    // Test Case Skip number
	Cancel   int    // Test Case Cancel number
	Total    int    // Test Case Total number
	Passrate string // Test Cases Passrate
	Results  string // Result array->string  []Result
	Nodes    string // Node array->string []Node
}

// ParseHTML .
func ParseHTML(hr *HTMLResult) *bytes.Buffer {
	tmpFile := "report-template.html" // report-template-lite.html
	_, file, _, _ := runtime.Caller(0)
	t, err := html_template.ParseFiles(path.Join(path.Dir(file), tmpFile))
	if err != nil {
		logger.Fatal(err)
	}
	body := new(bytes.Buffer)
	t.Execute(body, hr)
	return body
}

// ParseHTMLMail ParseHTML with text/template, for mail body
func ParseHTMLMail(hr *HTMLResult) *bytes.Buffer {
	tmpFile := "report-template-mail.html"
	_, file, _, _ := runtime.Caller(0)
	t, err := text_template.ParseFiles(path.Join(path.Dir(file), tmpFile))
	if err != nil {
		logger.Fatal(err)
	}
	body := new(bytes.Buffer)
	t.Execute(body, hr)
	return body
}
