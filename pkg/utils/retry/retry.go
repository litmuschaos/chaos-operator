package retry

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Action defines the prototype of action function, function as a value
type Action func(attempt uint) error

// Model defines the schema, contains all the attributes need for retry
type Model struct {
	retry    uint
	waitTime time.Duration
	timeout  int64
}

// Times is used to define the retry count
// it will run if the instance of model is not present before
func Times(retry uint) *Model {
	model := Model{}
	return model.Times(retry)
}

// Times is used to define the retry count
// it will run if the instance of model is already present
func (model *Model) Times(retry uint) *Model {
	model.retry = retry
	return model
}

// Wait is used to define the wait duration after each iteration of retry
// it will run if the instance of model is not present before
func Wait(waitTime time.Duration) *Model {
	model := Model{}
	return model.Wait(waitTime)
}

// Wait is used to define the wait duration after each iteration of retry
// it will run if the instance of model is already present
func (model *Model) Wait(waitTime time.Duration) *Model {
	model.waitTime = waitTime
	return model
}

// Timeout is used to define the timeout duration for each iteration of retry
// it will run if the instance of model is not present before
func Timeout(timeout int64) *Model {
	model := Model{}
	return model.Timeout(timeout)
}

// Timeout is used to define the timeout duration for each iteration of retry
// it will run if the instance of model is already present
func (model *Model) Timeout(timeout int64) *Model {
	model.timeout = timeout
	return model
}

// Try is used to run a action with retries and some delay after each iteration
func (model Model) Try(action Action) error {
	if action == nil {
		return fmt.Errorf("no action specified")
	}

	var err error
	for attempt := uint(0); (attempt == 0 || err != nil) && attempt <= model.retry; attempt++ {
		err = action(attempt)
		if model.waitTime > 0 {
			time.Sleep(model.waitTime)
		}
		if err == errors.Errorf("container is in terminated state") {
			break
		}
	}

	return err
}

// TryWithTimeout is used to run a action with retries
// for each iteration of retry there will be some timeout
func (model Model) TryWithTimeout(action Action) error {
	if action == nil {
		return fmt.Errorf("no action specified")
	}
	var err error
	err = nil
	for attempt := uint(0); (attempt == 0 || err != nil) && attempt <= model.retry; attempt++ {
		startTime := time.Now().Unix()
		currentTime := time.Now().Unix()
		for trial := uint(0); (trial == 0 || err != nil) && currentTime < startTime+model.timeout; trial++ {
			err = action(attempt)
			if model.waitTime > 0 {
				time.Sleep(model.waitTime)
			}
			currentTime = time.Now().Unix()
		}

	}

	return err
}
