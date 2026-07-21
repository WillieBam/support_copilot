package command_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/WillieBam/support_copilot/backend/internal/command"
	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
	"github.com/WillieBam/support_copilot/backend/types"
)

var _ = Describe("CommandInterceptor", func() {
	var (
		interceptor interfaces.ICommandInterceptor
		ctx         context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		interceptor = command.NewCommandInterceptor()
	})

	Context("Intercept /quit", func() {
		It("should intercept prompt starting with /quit", func() {
			res, err := interceptor.Intercept(ctx, "/quit")
			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
			Expect(res.Handled).To(BeTrue())
			Expect(res.Message).To(ContainSubstring("LLM processing stopped by /quit command"))
		})

		It("should intercept prompt with /quit and additional whitespace or text", func() {
			res, err := interceptor.Intercept(ctx, "  /QUIT now  ")
			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
			Expect(res.Handled).To(BeTrue())
		})

		It("should not intercept normal prompt", func() {
			res, err := interceptor.Intercept(ctx, "What is the system status?")
			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
			Expect(res.Handled).To(BeFalse())
			Expect(res.Message).To(BeEmpty())
		})
	})

	Context("Register custom command", func() {
		It("should allow registering custom slash command handlers", func() {
			// Use the concrete type to access RegisterCommand (not part of the interface)
			ci := command.NewCommandInterceptor().(*command.CommandInterceptor)
			ci.RegisterCommand("/ping", func(ctx context.Context, prompt string) (*types.CommandResult, error) {
				return &types.CommandResult{
					Handled: true,
					Message: "pong",
				}, nil
			})

			res, err := ci.Intercept(ctx, "/ping")
			Expect(err).NotTo(HaveOccurred())
			Expect(res.Handled).To(BeTrue())
			Expect(res.Message).To(Equal("pong"))
		})
	})
})
