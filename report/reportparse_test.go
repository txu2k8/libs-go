package report

import (
	"fmt"
	"libs-go/mail"
	"os"
	"strings"
	"testing"
	"time"
)

func mockHTMLResult() *HTMLResult {
	nodes := []Node{
		{
			Name:      "node1",
			Status:    "Ready",
			Roles:     "master",
			IPAddress: "10.112.208.67",
			User:      "root",
			Password:  "password",
		},
	}
	start := time.Now()
	time.Sleep(2 * time.Second)
	results := []TestResult{
		{
			Name:      "test1",
			Desc:      "reset nodes",
			Start:     time.Now(),
			End:       time.Now(),
			Elapsed:   fmt.Sprint(time.Since(start)),
			Status:    "PASS",
			Iteration: 1,
			Summary:   fmt.Sprintf("[PASS  ] - test1 - Iteration: 1 - Elapsed Time: %s", time.Since(start)),
		},
		{
			Name:      "test2",
			Desc:      "restart services",
			Start:     time.Now(),
			End:       time.Now(),
			Elapsed:   fmt.Sprint(time.Since(start)),
			Status:    "FAIL",
			Iteration: 2,
			Output:    "output",
		},
	}

	desc := TestDesc{Tester: "Tao.Xu"}
	// descJSON, _ := json.Marshal(desc)
	// descString := string(descJSON)

	// nodesJSON, _ := json.Marshal(nodes)
	// nodesString := string(nodesJSON)

	// resultsJSON, _ := json.Marshal(results)
	// resultsString := string(resultsJSON)

	return &HTMLResult{
		TimeDate: time.Now().Format("2006/01/02 15:04:05"),
		Content:  "Content=======",
		Title:    "Test Report Sample - 1",
		Desc:     desc.GenHTMLTableString(),
		Passrate: "90%",
		Pass:     10,
		Fail:     2,
		Error:    0,
		Skip:     1,
		Cancel:   1,
		Total:    13,
		Results:  GenResultHTMLTableString(results),
		Nodes:    GenNodesHTMLTableString(nodes),
	}
}

func TestReportHTML(T *testing.T) {
	logger.Info(strings.Repeat("▔", 50))
	logger.Info("Generate Test Report HTML ...")
	hr := mockHTMLResult()
	body := ParseHTMLMail(hr)
	reportFile := "C:\\tmp\\report.html"
	file, _ := os.OpenFile(reportFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	file.WriteString(body.String())
	file.Close()
	// logger.Debug(body.String())

	// Send results to mails
	subject := "Stress Test Results"
	mailConfig := &mail.Config{
		ServerHost: "smtp.gmail.com",
		ServerPort: 465,
		FromEmail:  "txu@xxx.com",
		FromPasswd: "Pass@word",
		FromUser:   "StressTest",
		Toers:      "txu1@xxx.com",
		CCers:      "",
	}
	mail.SendEmail(mailConfig, subject, body.String(), reportFile)
	logger.Info(strings.Repeat("▔", 50))
}

func TestMail(t *testing.T) {
	desc := TestDesc{
		Tester: "ssss",
	}
	tr := desc.GenHTMLTableString()
	logger.Info(tr)
}
