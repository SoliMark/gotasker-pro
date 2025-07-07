//go:build tools
// +build tools

package tools

// This import ensures `mockgen` dependencies are tracked in go.mod for versioning.
import (
	_ "github.com/golang/mock/mockgen/model"
)
