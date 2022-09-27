package utils

import "github.com/google/uuid"

func UUID() string {
	u, err := uuid.NewUUID()
	if err != nil {
		return ""
	}
	return u.String()
}
