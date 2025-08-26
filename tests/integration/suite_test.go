package integration

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TestContainerTestSuite 運行測試套件
func TestContainerTestSuite(t *testing.T) {
	suite.Run(t, new(ContainerTestSuite))
}
