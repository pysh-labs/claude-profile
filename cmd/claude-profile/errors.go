package main

import "fmt"

type PartialFailureError struct {
	Failed, Total int
}

func (e *PartialFailureError) Error() string {
	return fmt.Sprintf("%d of %d plugin installs failed", e.Failed, e.Total)
}
