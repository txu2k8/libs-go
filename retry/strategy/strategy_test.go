package strategy

import (
	"testing"
	"time"
)

// timeMarginOfError represents the acceptable amount of time that may pass for
// a time-based (sleep) unit before considering invalid.
const timeMarginOfError = time.Millisecond

func TestLimit(t *testing.T) {
	const attemptLimit = 3

	strategy := Limit(attemptLimit)
	retryArgs := RetryArgs{}

	if !strategy(1, &retryArgs) {
		t.Error("strategy expected to return true")
	}

	if !strategy(2, &retryArgs) {
		t.Error("strategy expected to return true")
	}

	if !strategy(3, &retryArgs) {
		t.Error("strategy expected to return true")
	}

	if strategy(4, &retryArgs) {
		t.Error("strategy expected to return false")
	}
}

func TestDelay(t *testing.T) {
	const delayDuration = time.Duration(10 * timeMarginOfError)

	strategy := Delay(delayDuration)
	retryArgs := RetryArgs{}

	if now := time.Now(); !strategy(0, &retryArgs) || delayDuration > time.Since(now) {
		t.Errorf(
			"strategy expected to return true in %s",
			time.Duration(delayDuration),
		)
	}

	if now := time.Now(); !strategy(5, &retryArgs) || (delayDuration/10) < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}
}

func TestWait(t *testing.T) {
	strategy := Wait()
	retryArgs := RetryArgs{}

	if now := time.Now(); !strategy(0, &retryArgs) || timeMarginOfError < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}

	if now := time.Now(); !strategy(999, &retryArgs) || timeMarginOfError < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}
}

func TestWaitWithDuration(t *testing.T) {
	const waitDuration = time.Duration(10 * timeMarginOfError)

	strategy := Wait(waitDuration)
	retryArgs := RetryArgs{}

	if now := time.Now(); !strategy(0, &retryArgs) || timeMarginOfError < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}

	if now := time.Now(); !strategy(1, &retryArgs) || waitDuration > time.Since(now) {
		t.Errorf(
			"strategy expected to return true in %s",
			time.Duration(waitDuration),
		)
	}
}

func TestWaitWithMultipleDurations(t *testing.T) {
	waitDurations := []time.Duration{
		time.Duration(10 * timeMarginOfError),
		time.Duration(20 * timeMarginOfError),
		time.Duration(30 * timeMarginOfError),
		time.Duration(40 * timeMarginOfError),
	}

	strategy := Wait(waitDurations...)
	retryArgs := RetryArgs{}

	if now := time.Now(); !strategy(0, &retryArgs) || timeMarginOfError < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}

	if now := time.Now(); !strategy(1, &retryArgs) || waitDurations[0] > time.Since(now) {
		t.Errorf(
			"strategy expected to return true in %s",
			time.Duration(waitDurations[0]),
		)
	}

	if now := time.Now(); !strategy(3, &retryArgs) || waitDurations[2] > time.Since(now) {
		t.Errorf(
			"strategy expected to return true in %s",
			waitDurations[2],
		)
	}

	if now := time.Now(); !strategy(999, &retryArgs) || waitDurations[len(waitDurations)-1] > time.Since(now) {
		t.Errorf(
			"strategy expected to return true in %s",
			waitDurations[len(waitDurations)-1],
		)
	}
}

func TestBackoff(t *testing.T) {
	const backoffDuration = time.Duration(10 * timeMarginOfError)
	const algorithmDurationBase = timeMarginOfError

	algorithm := func(attempt uint) time.Duration {
		return backoffDuration - (algorithmDurationBase * time.Duration(attempt))
	}

	strategy := Backoff(algorithm)
	retryArgs := RetryArgs{}

	if now := time.Now(); !strategy(0, &retryArgs) || timeMarginOfError < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}

	for i := uint(1); i < 10; i++ {
		expectedResult := algorithm(i)

		if now := time.Now(); !strategy(i, &retryArgs) || expectedResult > time.Since(now) {
			t.Errorf(
				"strategy expected to return true in %s",
				expectedResult,
			)
		}
	}
}

func TestBackoffWithJitter(t *testing.T) {
	const backoffDuration = time.Duration(10 * timeMarginOfError)
	const algorithmDurationBase = timeMarginOfError

	algorithm := func(attempt uint) time.Duration {
		return backoffDuration - (algorithmDurationBase * time.Duration(attempt))
	}

	transformation := func(duration time.Duration) time.Duration {
		return duration - time.Duration(10*timeMarginOfError)
	}

	strategy := BackoffWithJitter(algorithm, transformation)
	retryArgs := RetryArgs{}

	if now := time.Now(); !strategy(0, &retryArgs) || timeMarginOfError < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}

	for i := uint(1); i < 10; i++ {
		expectedResult := transformation(algorithm(i))

		if now := time.Now(); !strategy(i, &retryArgs) || expectedResult > time.Since(now) {
			t.Errorf(
				"strategy expected to return true in %s",
				expectedResult,
			)
		}
	}
}

func TestNoJitter(t *testing.T) {
	transformation := noJitter()

	for i := uint(0); i < 10; i++ {
		duration := time.Duration(i) * timeMarginOfError
		result := transformation(duration)
		expected := duration

		if result != expected {
			t.Errorf("transformation expected to return a %s duration, but received %s instead", expected, result)
		}
	}
}
