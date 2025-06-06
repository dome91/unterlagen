package common

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func GenerateID() string {
	return gonanoid.Must(8)
}
