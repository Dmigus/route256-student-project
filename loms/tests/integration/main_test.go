//go:build integration

package integration

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestCreateSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(Suite))
}
