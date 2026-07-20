package config_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/WillieBam/support_copilot/backend/app/config"
)

var _ = Describe("ServerConfig", func() {
	It("should construct NewServerConfig with given name and getters", func() {
		sc := config.NewServerConfig("test-server")
		Expect(sc).NotTo(BeNil())
		Expect(sc.Name()).To(Equal("test-server"))
		Expect(sc.Port()).To(Equal("8080"))
	})
})

