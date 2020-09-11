package retry

import (
	"errors"
	"libs-go/retry/strategy"
	"testing"
)

func TestRetry(t *testing.T) {
	action := func(attempt uint) error {
		return nil
	}

	err := Retry(action)

	if nil != err {
		t.Error("expected a nil error")
	}
}

func TestRetryRetriesUntilNoErrorReturned(t *testing.T) {
	const errorUntilAttemptNumber = 5

	var attemptsMade uint

	action := func(attempt uint) error {
		attemptsMade = attempt

		if errorUntilAttemptNumber == attempt {
			return nil
		}

		return errors.New("erroring")
	}

	err := Retry(action)

	if nil != err {
		t.Error("expected a nil error")
	}

	if errorUntilAttemptNumber != attemptsMade {
		t.Errorf(
			"expected %d attempts to be made, but %d were made instead",
			errorUntilAttemptNumber,
			attemptsMade,
		)
	}
}

func TestShouldAttempt(t *testing.T) {
	retryArgs := strategy.RetryArgs{}
	shouldAttempt := shouldAttempt(1, &retryArgs)

	if !shouldAttempt {
		t.Error("expected to return true")
	}
}

func TestShouldAttemptWithStrategy(t *testing.T) {
	retryArgs := strategy.RetryArgs{}
	const attemptNumberShouldReturnFalse = 7

	strategy := func(attempt uint, retryArgs *strategy.RetryArgs) bool {
		return (attemptNumberShouldReturnFalse != attempt)
	}

	should := shouldAttempt(1, &retryArgs, strategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(1+attemptNumberShouldReturnFalse, &retryArgs, strategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(attemptNumberShouldReturnFalse, &retryArgs, strategy)

	if should {
		t.Error("expected to return false")
	}
}

func TestShouldAttemptWithMultipleStrategies(t *testing.T) {
	retryArgs := strategy.RetryArgs{}
	trueStrategy := func(attempt uint, retryArgs *strategy.RetryArgs) bool {
		return true
	}

	falseStrategy := func(attempt uint, retryArgs *strategy.RetryArgs) bool {
		return false
	}

	should := shouldAttempt(1, &retryArgs, trueStrategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(1, &retryArgs, falseStrategy)

	if should {
		t.Error("expected to return false")
	}

	should = shouldAttempt(1, &retryArgs, trueStrategy, trueStrategy, trueStrategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(1, &retryArgs, falseStrategy, falseStrategy, falseStrategy)

	if should {
		t.Error("expected to return false")
	}

	should = shouldAttempt(1, &retryArgs, trueStrategy, trueStrategy, falseStrategy)

	if should {
		t.Error("expected to return false")
	}
}
