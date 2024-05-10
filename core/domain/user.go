package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system.
type User struct {
	ID       uuid.UUID
	Name     string
	Settings map[string]interface{} // Make sure, that the JSONB submitted to the database is actually a JSON object
	LastSeen *time.Time
}

// String returns a string representation of the User object.
func (u User) String() string {
	result := fmt.Sprintf("user %s (%s)", u.Name, u.ID)
	if u.LastSeen != nil {
		result += fmt.Sprintf(", last seen at %s", u.LastSeen)
	}

	if u.Settings != nil {
		result += fmt.Sprintf(", settings %v", u.Settings)
	}

	return result
}
