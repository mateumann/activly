package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system.
type User struct {
	ID       uuid.UUID
	Name     string
	Settings []byte
	LastSeen *time.Time
}

// String returns a string representation of the User object.
//
// The returned string consists of the user's name and ID. If the user's
// LastSeen field is not nil, it also includes the last seen time. If the
// user's Settings field is not nil, then the settings are unmarshalled
// from JSON and included in the string representation.
//
// Example:
//
//	u := User{
//	    ID:       uuid.New(),
//	    Name:     "John Doe",
//	    Settings: []byte(`{"language": "en", "theme": "dark"}`),
//	    LastSeen: &time.Time{},
//	}
//	fmt.Println(u.String())  // Output: user John Doe (a1234567-89ab-cdef-0123-456789abcdef),
//	                         //         last seen at 0001-01-01 00:00:00 +0000 UTC,
//	                         //         settings map[language:en theme:dark]
//
// Returns:
//
//	The string representation of the User object.
func (u User) String() string {
	result := fmt.Sprintf("user %s (%s)", u.Name, u.ID)
	if u.LastSeen != nil {
		result += fmt.Sprintf(", last seen at %s", u.LastSeen)
	}

	if u.Settings != nil {
		var settings map[string]any
		_ = json.Unmarshal(u.Settings, &settings)
		result += fmt.Sprintf(", settings %v", settings)
	}

	return result
}
