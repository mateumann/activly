package domain

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUser_String(t *testing.T) {
	lastSeen := time.Now()
	settings := make(map[string]any)
	settings["classify"] = true
	id1, id2, id3 := uuid.New(), uuid.New(), uuid.New()
	tests := []struct {
		name string
		u    User
		want string
	}{
		{
			name: "Complete User",
			u: User{
				ID:       id1,
				Name:     "TestUser",
				Settings: settings,
				LastSeen: &lastSeen,
			},
			want: "user TestUser (" + id1.String() + "), last seen at " + lastSeen.String() +
				", settings map[classify:true]",
		},
		{
			name: "User Without Last Seen and Settings",
			u: User{
				ID:   id2,
				Name: "TestUser",
			},
			want: "user TestUser (" + id2.String() + ")",
		},
		{
			name: "User Without Name",
			u: User{
				ID:       id3,
				Settings: settings,
				LastSeen: &lastSeen,
			},
			want: "user  (" + id3.String() + "), last seen at " + lastSeen.String() + ", settings map[classify:true]",
		},
		{
			name: "Empty User",
			u:    User{},
			want: "user  (00000000-0000-0000-0000-000000000000)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
