package worker

// pzct used this lib for framework

import (
	"errors"
	"fmt"
	"libs-go/config"
	"libs-go/mail"
	"libs-go/report"
	"libs-go/utils"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/chenhg5/collection"
	"github.com/op/go-logging"
)

// define global values
var (
	logger       = logging.MustGetLogger("test")
	wg           sync.WaitGroup
	ErrTimeout   = errors.New("received timeout")
	ErrInterrupt = errors.New("received interrupt")
)

// StatusResult result status num
type StatusResult struct {
	// time
	Start   time.Time     // tc start time
	End     time.Time     // tc end time
	Elapsed time.Duration // tc elapsed time
	// tc Pass|Fail|Error|Skip|Cancel iteration number
	NumPass   int
	NumFail   int
	NumError  int
	NumSkip   int
	NumCancel int
	NumTotal  int
	PassRate  string // eg: 90%
	// status | iteration | Summary
	Status            report.Status // tc total status
	Iteration         int           // tc total exec iteration number
	Summary           string        // ALL 10, Pass 8, Fail: 1, Error 1, Passing Rate:80%
	IterationSummarys []string      // report.TestResult Summarys
}

// Task : Multi-task run in order, and os.exit if return error
type Task struct {
	Fn       func() error // task function
	Name     string       // task name
	Desc     string       // task description
	RunTimes int          // run task times limit, 0: forever

	interrupt bool // received interrupt or not
	complete  chan error
	timeout   <-chan time.Time
	results   []report.TestResult
	loop      int // running task loop
}

// Worker .
type Worker struct {
	tasks   []Task
	results []report.TestResult
	wg      sync.WaitGroup
}

// NewWorker Worker
func NewWorker() *Worker {
	return &Worker{}
}

// aggregate tasks all Iteration results to Worker
func aggregate(results []report.TestResult) func(report.TestResult) {
	return func(t report.TestResult) {
		results = append(results, t)
	}
}

// SetTimeOut .
func (t *Task) SetTimeOut(d time.Duration) {
	t.timeout = time.After(d)
}

// Add test cases
func (w *Worker) Add(tasks ...Task) {
	w.tasks = append(w.tasks, tasks...)
}

// Start Exec tasks
func (w *Worker) Start() {
	for _, task := range w.tasks {
		task.init(w.results)
		// task.aggregate = aggregate(w.results)
		task.start()
		w.results = append(w.results, task.results...)
	}
	logger.Info("All Test Cases finishd\n")
}

// ============= task =============

func getPanicStack(err error) string {
	var buf [40960]byte
	n := runtime.Stack(buf[:], false)
	// logger.Criticalf("%s\n==> %s", err, string(getPanicTrace(4)))
	return fmt.Sprintf("%s\n==> %s", err, string(buf[:n]))
}

// Capture signal: Cancel
func (t *Task) listenSignal() {
	start := time.Now()
	sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	signal.Notify(sigs, os.Interrupt)

	select {
	case sig := <-sigs:
		output := fmt.Sprintf("Scripts exit with signal: %v", sig)
		logger.Warning(output)
		err := fmt.Errorf("Canceled")
		logger.Error(err)
		// <-sigs
		signal.Stop(sigs)
		t.interrupt = true
		t.RunTimes = t.loop
		end := time.Now()
		elapsed := fmt.Sprint(time.Since(start))
		status := report.CancelStatus
		summary := fmt.Sprintf("[%-6s] - %s - Iteration: %d - Elapsed Time: %s", status, t.Name, t.loop, elapsed)
		logger.Info(summary)
		res := report.TestResult{
			Name:      t.Name,
			Start:     start,
			End:       end,
			Elapsed:   elapsed,
			Output:    output,
			Status:    status,
			Iteration: t.loop,
			Summary:   summary,
		}
		t.results = append(t.results, res)
		wg.Done()
		logger.Infof("TestCase(Name:'%s', RunTimes:%d) is done\n", t.Name, t.RunTimes)
		return
	}
}

// Capture panic
func (t *Task) listenPanic() {
	var err error
	start := time.Now()
	if rc := recover(); rc != nil {
		// logger.Criticalf("Capture panic：%s", rc)
		switch x := rc.(type) {
		case string:
			err = errors.New(x)
		case error:
			err = x
		default:
			err = errors.New("Unknow panic")
		}
		output := "Panic"
		if err != nil {
			output = getPanicStack(err)
			logger.Critical(output)
			t.RunTimes = t.loop
		}
		end := time.Now()
		elapsed := fmt.Sprint(time.Since(start))
		status := report.ErrorStatus // ERROR if panic
		summary := fmt.Sprintf("[%-6s] - %s - Iteration: %d - Elapsed Time: %s", status, t.Name, t.loop, elapsed)
		logger.Info(summary)
		res := report.TestResult{
			Name:      t.Name,
			Start:     start,
			End:       end,
			Elapsed:   elapsed,
			Output:    output,
			Status:    status,
			Iteration: t.loop,
			Summary:   summary,
		}
		t.results = append(t.results, res)
		wg.Done()
		return
	}
}

func (t *Task) init(results []report.TestResult) {
	// t.interrupt = make(chan os.Signal, 1)
	t.interrupt = false
	t.complete = make(chan error)
}

// start run fn loops
func (t *Task) start() {
	go t.listenSignal()

	baton := make(chan int)
	wg.Add(1)
	go t.run(baton)
	baton <- 1 // start loop 1
	wg.Wait()
}

func (t *Task) run(baton chan int) {
	var err error
	var newLoop int
	status := report.SkipStatus
	start := time.Now()
	forever := t.RunTimes <= 0
	t.loop = <-baton
	defer t.listenPanic()
	// Create for next Loop
	if forever || t.loop < t.RunTimes {
		newLoop = t.loop + 1
		go t.run(baton)
	}

	// Exec tc
	logger.Infof("[START ] - %s - Iteration: %d", t.Name, t.loop)
	output := ""
	err = t.Fn()
	// if t.interrupt {
	// 	// Stop loop
	// 	logger.Infof("TestCase(Name:'%s', RunTimes:%d) is done\n", t.Name, t.RunTimes)
	// 	wg.Done()
	// 	return
	// }
	if err != nil {
		logger.Error(err)
		t.RunTimes = t.loop
		forever = false
		status = report.FailStatus
		output = fmt.Sprintf("%s", err)
	} else {
		status = report.PassStatus
	}
	end := time.Now()
	elapsed := fmt.Sprint(time.Since(start))
	summary := fmt.Sprintf("[%-6s] - %s - Iteration: %d - Elapsed Time: %s", status, t.Name, t.loop, elapsed)
	logger.Info(summary)
	res := report.TestResult{
		Name:      t.Name,
		Start:     start,
		End:       end,
		Elapsed:   elapsed,
		Output:    output,
		Status:    status,
		Iteration: t.loop,
		Summary:   summary,
	}
	t.results = append(t.results, res)

	if t.loop == t.RunTimes {
		// Stop loop
		logger.Infof("TestCase(Name:'%s', RunTimes:%d) is done\n", t.Name, t.RunTimes)
		wg.Done()
		return
	}

	// Start next loop
	baton <- newLoop
}

// ============= StatusResult =============

// collect results summary|status
func (sr *StatusResult) collectStatusResult(results []report.TestResult) {
	for _, res := range results {
		logger.Info(res.Summary)
		sr.IterationSummarys = append(sr.IterationSummarys, res.Summary)
		sr.Status = res.Status
		sr.NumTotal++
		switch res.Status {
		case report.PassStatus:
			sr.NumPass++
		case report.FailStatus:
			sr.NumFail++
		case report.ErrorStatus:
			sr.NumError++
		case report.SkipStatus:
			sr.NumSkip++
		case report.CancelStatus:
			sr.NumCancel++
		default:
			logger.Fatalf("Not defined test status: %s", res.Status)
		}
	}
	sr.PassRate = fmt.Sprintf("%.1f%%", float32(sr.NumPass+sr.NumCancel)/float32(sr.NumTotal-sr.NumSkip)*100)
	sr.Summary = fmt.Sprintf("ALL %d, Pass %d, Fail: %d, Error %d, Passing Rate:%s", sr.NumTotal, sr.NumPass, sr.NumFail, sr.NumError, sr.PassRate)
}

// Run entry
func Run(tasks []Task) {
	startTime := time.Now()
	w := NewWorker()
	w.Add(tasks...)
	w.Start()

	endTime := time.Now()
	elapsedTime := time.Since(startTime)
	// ===================== Run Test complete =====================

	logger.Info(strings.Repeat("▔", 65))
	sr := StatusResult{
		Start:   startTime,
		End:     endTime,
		Elapsed: elapsedTime,
		Status:  report.PassStatus,
	}
	sr.collectStatusResult(w.results)

	logPath := config.GetValue("log_path").(string)
	image := config.GetValue("dpl_image")
	masterIPs := config.GetValue("master_ips")

	logger.Info("Total:", sr.NumTotal)
	logger.Info("Time Elapsed:", elapsedTime)
	logger.Info("Log:", logPath)
	logger.Info(strings.Repeat("▔", 65))

	// ================= Send report to mails  =================
	ignoreMailCmds := []string{"maint", "tools", "version"}
	for _, cmd := range os.Args[1:] {
		if collection.Collect(ignoreMailCmds).Contains(cmd) {
			return
		}
	}
	nodesEnv := config.GetValue("nodes")
	nodes := []report.Node{}
	if nodesEnv != nil {
		for _, node := range config.GetValue("nodes").([]map[string]string) {
			nodes = append(nodes, report.Node{
				Name:      node["Name"],
				IPAddress: node["IP"],
				Status:    node["Status"],
				Roles:     node["Role"],
				User:      fmt.Sprintf("%v", config.GetValue("sys_user")),
				Password:  fmt.Sprintf("%v", config.GetValue("sys_pwd")),
			})
		}
	}

	desc := report.TestDesc{
		Tester:       "txu",
		Version:      "6.0.0",
		StartTime:    fmt.Sprintf("%s", startTime),
		EndTime:      fmt.Sprintf("%s", endTime),
		ElapsedTime:  fmt.Sprintf("%s", elapsedTime),
		Summary:      sr.Summary,
		DplImage:     fmt.Sprintf("%v", image),
		MasterIPs:    fmt.Sprintf("%v", masterIPs),
		TestLocation: utils.GetLocalIP(),
		ReportPath:   logPath,
		Command:      fmt.Sprintf("%v", config.GetValue("command")),
	}
	caseArr := []string{}
	for _, t := range tasks {
		caseArr = append(caseArr, t.Name)
	}
	subject := fmt.Sprintf("%s:%v-%s", sr.Status, config.GetValue("report_subject"), strings.Join(caseArr, "_"))

	htmlR := &report.HTMLResult{
		Title:    subject,
		Nodes:    report.GenNodesHTMLTableString(nodes),
		Desc:     desc.GenHTMLTableString(),
		Pass:     sr.NumPass,
		Fail:     sr.NumFail,
		Error:    sr.NumError,
		Skip:     sr.NumSkip,
		Cancel:   sr.NumCancel,
		Total:    sr.NumTotal,
		Passrate: sr.PassRate,
		Results:  report.GenResultHTMLTableString(w.results),
	}
	body := report.ParseHTMLMail(htmlR)
	reportFile := strings.Split(logPath, ".log")[0] + "-report.html"
	file, _ := os.OpenFile(reportFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	file.WriteString(body.String())
	file.Close()
	// logger.Info(body.String())
	attachments := []string{} // reportFile
	if utils.GetFileSize(logPath) < 5*1024*1024 {
		attachments = append(attachments, logPath)
	}
	if config.GetValue("debug").(bool) == true {
		config.EmailConfig.CCers = ""
	}
	mail.SendEmail(&config.EmailConfig, subject, body.String(), attachments...)
	logger.Info(strings.Repeat("▔", 65))
}
