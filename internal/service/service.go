package service

import (
"fmt"
)

// Common service errors
type BadRequest struct {
	Err error
}

func (b BadRequest) Error() string {
	return fmt.Sprintf("%v", b.Err)
}

func (BadRequest) BadRequest() {}

type NotFound struct {
	Err error
}

func (n NotFound) Error() string {
	return fmt.Sprintf("%v", n.Err)
}

func (NotFound) NotFound() {}
