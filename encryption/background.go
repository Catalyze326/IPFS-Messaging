package main

import (
	"errors"
	"fmt"
	"time"
)

func main() {
	keepDoingSomething()
}

// keepDoingSomething will keep trying to doSomething() until either
// we get a result from doSomething() or the timeout expires
func keepDoingSomething() (bool, error) {
	timeout := time.After(1 * time.Second)
	tick := time.Tick(500 * time.Millisecond)
	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-timeout:
			return false, errors.New("timed out")
		// Got a tick, we should check on doSomething()
		case <-tick:
			ok, err := doSomething()
			// Error from doSomething(), we should bail
			if err != nil {
				return false, err
				// doSomething() worked! let's finish up
			} else if ok {
				return true, nil
			}
			fmt.Println("hey")
			// doSomething() didn't work yet, but it didn't fail, so let's try again
			// this will exit up to the for loop
		}
	}
}

//decrypt the message that we just encrypted
func doSomething() (bool, error) {
	for true {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("lol")
	}
	return true, nil
}
