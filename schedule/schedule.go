package schedule

import (
	"libs-go/prettytable"
	"reflect"
	"runtime"
	"strings"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("test")

// Action defines a callable function that package retry can handle.
type Action func() error

// Enter just for print enter step， do nothing
func Enter() error { return nil }

// Scheduler .
type Scheduler interface {
	SetUp(options ...OptionFunc) error
	TearDown(options ...OptionFunc) error
	RunPhase(action Action, options ...OptionFunc) error
}

// Input .
type Input struct {
	Verbosity int           // Print Phase in Teardown If > 0
	Skip      bool          // Skip the phase if true
	Desc      string        // The phase description
	FnLevel   int           // The func level(func name as Phase name), default=1
	PanicErr  bool          // panic if fn()err != nil
	FnArgs    []interface{} // The args for Fn(args ...interface{})
}

// Phase .
type Phase struct {
	Idx    int    // The phase index, start with 1
	Name   string // The phase name
	Status string // The phase running status
	Desc   string // The phase description
}

// Schedule .
type Schedule struct {
	Input    Input   // Input args
	PhaseArr []Phase // Store the running phase list
}

// PrintPhase .
func (sc *Schedule) PrintPhase() error {
	table, _ := prettytable.NewTable(
		prettytable.Column{Header: "No."},
		prettytable.Column{Header: "Step"},
		prettytable.Column{Header: "Result"},
		prettytable.Column{Header: "Description"},
	)
	for _, p := range sc.PhaseArr {
		table.AddRow(p.Idx, p.Name, p.Status, p.Desc)
	}
	if len(table.Rows) > 0 {
		logger.Infof("Test Progress:\n%s", table.String())
	}

	return nil
}

// ApplyOptions Apply any given schedule options.
func (sc *Schedule) ApplyOptions(options ...OptionFunc) error {
	for _, fn := range options {
		if fn == nil {
			continue
		}
		if err := fn(sc); err != nil {
			return err
		}
	}
	return nil
}

// SetUp .
func (sc *Schedule) SetUp(options ...OptionFunc) error {
	// Initialize Input
	sc.Input = Input{
		Verbosity: 1,
		FnLevel:   1,
	}
	sc.ApplyOptions(options...)

	status := "START"
	if sc.Input.Skip == true {
		status = "SKIP"
	}
	description := sc.Input.Desc
	idx := len(sc.PhaseArr) + 1
	pc, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	fName := f.Name()
	fNameSplit := strings.Split(fName, ".")
	pName := fNameSplit[len(fNameSplit)-sc.Input.FnLevel]
	phase := Phase{
		Idx:    idx,
		Name:   pName,
		Status: status,
		Desc:   description,
	}

	if idx <= 1 {
		sc.PhaseArr = append(sc.PhaseArr, phase)
	} else {
		lastPhase := sc.PhaseArr[idx-2]
		if lastPhase.Name != fName && lastPhase.Status != status {
			sc.PhaseArr = append(sc.PhaseArr, phase)
		}
	}

	if sc.Input.Verbosity > 0 {
		sc.PrintPhase()
	} else {
		logger.Infof("%s: %s", status, fName)
	}

	return nil
}

// TearDown .
func (sc *Schedule) TearDown(options ...OptionFunc) error {
	// Initialize Input
	sc.Input = Input{
		Verbosity: 1,
		FnLevel:   1,
	}
	sc.ApplyOptions(options...)
	if sc.Input.Skip == true {
		return nil
	}

	pc, _, _, ok := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	fName := f.Name()
	fNameSplit := strings.Split(fName, ".")
	pName := fNameSplit[len(fNameSplit)-sc.Input.FnLevel]
	status := "PASS"
	if ok == false {
		status = "FAIL"
	}

	description := sc.Input.Desc
	idx := len(sc.PhaseArr) + 1

	phase := Phase{
		Idx:    idx,
		Name:   pName,
		Status: status,
		Desc:   description,
	}
	sc.PhaseArr = append(sc.PhaseArr, phase)

	if sc.Input.Verbosity > 0 {
		sc.PrintPhase()
	} else {
		logger.Infof("%s: %s", status, fName)
	}

	return nil
}

// RunPhase .
func (sc *Schedule) RunPhase(action Action, options ...OptionFunc) error {
	// Initialize Input
	sc.Input = Input{
		Verbosity: 1,
		FnLevel:   1,
	}
	sc.ApplyOptions(options...)
	phaseInput := sc.Input

	status := "START"
	if phaseInput.Skip == true {
		status = "SKIP"
	}
	description := "nil"
	if phaseInput.Desc != "" {
		description = phaseInput.Desc
	}
	idx := len(sc.PhaseArr) + 1
	fName := strings.TrimSuffix(runtime.FuncForPC(reflect.ValueOf(action).Pointer()).Name(), "-fm")
	fNameSplit := strings.Split(fName, ".")
	pName := fNameSplit[len(fNameSplit)-sc.Input.FnLevel]
	phase := Phase{
		Idx:    idx,
		Name:   pName,
		Status: status,
		Desc:   description,
	}

	if idx <= 1 {
		sc.PhaseArr = append(sc.PhaseArr, phase)
	} else {
		lastPhase := sc.PhaseArr[idx-2]
		if pName == "Enter" || lastPhase.Name != pName || lastPhase.Status != status {
			sc.PhaseArr = append(sc.PhaseArr, phase)
		} else {
			idx--
		}
	}

	if phaseInput.Verbosity > 0 {
		sc.PrintPhase()
	} else {
		logger.Infof("%s: %s", status, fName)
		logger.Infof("Description: %s", description)
	}
	// Run func
	if phaseInput.Skip == true || pName == "action" {
		return nil
	}
	logger.Infof("Enter %s ...", fName)
	err := action()
	status = "PASS"
	if err != nil {
		status = "FAIL"
	}
	sc.PhaseArr[idx-1].Status = status
	if phaseInput.Verbosity > 1 {
		sc.PrintPhase()
	} else {
		if status == "FAIL" {
			logger.Errorf("%s: %s", status, fName)
			logger.Errorf("Description: %s", description)
		} else {
			logger.Infof("%s: %s", status, fName)
			logger.Infof("Description: %s", description)
		}

	}

	if phaseInput.PanicErr == true && err != nil {
		logger.Panicf("Run %s Failed: %v", pName, err)
	}

	return err
}

// RunPhasePanicErr RunPhase with PanicErr=true
func (sc *Schedule) RunPhasePanicErr(action Action, options ...OptionFunc) error {
	options = append(options, PanicErr(true))
	return sc.RunPhase(action, options...)
}
