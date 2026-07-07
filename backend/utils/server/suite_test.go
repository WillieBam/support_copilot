package server_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestServerUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Utils Suite")
}
