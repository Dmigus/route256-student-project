//go:build e2e
// +build e2e

package internal

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func Test2e2(t *testing.T) {
	suite.Run(t, new(Suit))
}
