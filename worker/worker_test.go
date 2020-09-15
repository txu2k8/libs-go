package worker

import (
	"errors"
	"libs-go/config"
	"libs-go/utils"
	"math/rand"
	"testing"
	"time"
)

func testFun1() error {
	if rand.Intn(2) == 0 {
		panic("sssssss")
		// return errors.New("testFun1 error")
	}
	logger.Info("tttttttttttttttestFun1")

	return nil
}

func testFun2() error {
	if rand.Intn(2) == 0 {
		return errors.New("testFun2 error")
	}
	return nil
}

func testFun3() error {
	if rand.Intn(2) == 0 {
		utils.SleepProgressBar(10 * time.Second)
		return errors.New("testFun3 error")
	}
	return nil
}

func TestWorker(t *testing.T) {
	config.SetValue("log_path", "/tmp/test.log")
	tasks := []Task{
		{
			Fn:       testFun1,
			Name:     "testFun1",
			RunTimes: 10,
		},
		// {
		// 	Fn:       testFun2,
		// 	Name:     "testFun2",
		// 	RunTimes: 10,
		// },
		{
			Fn:       testFun3,
			Name:     "testFun3",
			RunTimes: 10,
		},
	}
	Run(tasks)
}
