package mcp_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMCPRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MCP Repository Suite")
}
