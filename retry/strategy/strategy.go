// Package strategy provides a way to change the way that retry is performed.
//
// Copyright © 2016 Trevor N. Suarez (Rican7)
package strategy

import (
	"fmt"
	"time"

	"libs-go/retry/backoff"
	"libs-go/retry/jitter"

	"github.com/schollz/progressbar/v3"
)

// RetryArgs .
type RetryArgs struct {
	Tries   uint
	wait    time.Duration
	delay   time.Duration
	backoff time.Duration
}

// SleepProgressBar ...
func SleepProgressBar(d time.Duration, attempt, Tries uint) {
	intSeconds := int(d.Seconds())
	bar := progressbar.New(intSeconds)
	bar.Describe(fmt.Sprintf("Sleep %s before retry %d/%d -", d, attempt, Tries))
	for i := 0; i < intSeconds; i++ {
		bar.Add(1)
		time.Sleep(1 * time.Second)
	}
	fmt.Println()
}

// Strategy defines a function that Retry calls before every successive attempt
// to determine whether it should make the next attempt or not. Returning `true`
// allows for the next attempt to be made. Returning `false` halts the retrying
// process and returns the last error returned by the called Action.
//
// The strategy will be passed an "attempt" number on each successive retry
// iteration, starting with a `0` value before the first attempt is actually
// made. This allows for a pre-action delay, etc.
type Strategy func(attempt uint, retryArgs *RetryArgs) bool

// Limit creates a Strategy that limits the number of attempts that Retry will
// make.
func Limit(attemptLimit uint) Strategy {
	return func(attempt uint, retryArgs *RetryArgs) bool {
		retryArgs.Tries = attemptLimit
		return (attempt <= attemptLimit)
	}
}

// Delay creates a Strategy that waits the given duration before the first
// attempt is made.
func Delay(duration time.Duration) Strategy {
	return func(attempt uint, retryArgs *RetryArgs) bool {
		retryArgs.delay = duration
		if attempt == 0 {
			SleepProgressBar(duration, attempt, retryArgs.Tries)
			// time.Sleep(duration)
		}

		return true
	}
}

// Wait creates a Strategy that waits the given durations for each attempt after
// the first. If the number of attempts is greater than the number of durations
// provided, then the strategy uses the last duration provided.
func Wait(durations ...time.Duration) Strategy {
	return func(attempt uint, retryArgs *RetryArgs) bool {
		if 1 < attempt && 0 < len(durations) {
			durationIndex := int(attempt - 1)

			if len(durations) <= durationIndex {
				durationIndex = len(durations) - 1
			}
			retryArgs.wait = durations[durationIndex]
			SleepProgressBar(durations[durationIndex], attempt, retryArgs.Tries)
			// time.Sleep(durations[durationIndex])
		}

		return true
	}
}

// Backoff creates a Strategy that waits before each attempt, with a duration as
// defined by the given backoff.Algorithm.
func Backoff(algorithm backoff.Algorithm) Strategy {
	return BackoffWithJitter(algorithm, noJitter())
}

// BackoffWithJitter creates a Strategy that waits before each attempt, with a
// duration as defined by the given backoff.Algorithm and jitter.Transformation.
func BackoffWithJitter(algorithm backoff.Algorithm, transformation jitter.Transformation) Strategy {
	return func(attempt uint, retryArgs *RetryArgs) bool {
		if 1 < attempt {
			retryArgs.backoff = transformation(algorithm(attempt))
			SleepProgressBar(transformation(algorithm(attempt)), attempt, retryArgs.Tries)
			// time.Sleep(transformation(algorithm(attempt)))
		}

		return true
	}
}

// noJitter creates a jitter.Transformation that simply returns the input.
func noJitter() jitter.Transformation {
	return func(duration time.Duration) time.Duration {
		return duration
	}
}
