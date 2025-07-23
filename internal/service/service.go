package service

import (
"fmt"

"github.com/google/uuid"
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

// Helper functions
func StringToUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}

func UUIDToString(id uuid.UUID) string {
	return id.String()
}

// Common status constants
const (
Active   = "active"
Inactive = "inactive"
Deleted  = "deleted"
)
