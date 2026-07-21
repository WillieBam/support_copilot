package classifier_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/WillieBam/support_copilot/backend/internal/classifier"
)

var _ = Describe("IntentClassifier", func() {
	var c *classifier.IntentClassifier

	BeforeEach(func() {
		c = classifier.NewIntentClassifier()
	})

	DescribeTable("should return IntentConversational for social / acknowledgement prompts",
		func(prompt string) {
			Expect(c.Classify(prompt)).To(Equal(classifier.IntentConversational))
		},
		// Acknowledgements
		Entry("ok", "ok"),
		Entry("ok byebye", "ok byebye"),
		Entry("okay", "okay"),
		Entry("alright", "alright"),
		Entry("got it", "got it"),
		Entry("sure", "sure"),
		Entry("fine", "fine"),
		Entry("noted", "noted"),
		// Sign-offs
		Entry("bye", "bye"),
		Entry("byebye", "byebye"),
		Entry("goodbye", "goodbye"),
		Entry("see you", "see you"),
		Entry("later", "later"),
		// Gratitude
		Entry("thanks", "thanks"),
		Entry("thank you", "thank you"),
		Entry("Thank You!", "Thank You!"),
		Entry("ty", "ty"),
		Entry("cheers", "cheers"),
		// Greetings
		Entry("hi", "hi"),
		Entry("hello", "hello"),
		Entry("hey", "hey"),
		Entry("good morning", "good morning"),
		// Yes/No
		Entry("yes", "yes"),
		Entry("no", "no"),
		Entry("yep", "yep"),
		Entry("nah", "nah"),
		// Completion
		Entry("done", "done"),
		Entry("stop", "stop"),
		Entry("that's all", "that's all"),
		Entry("finished", "finished"),
	)

	DescribeTable("should return IntentTask for task-oriented prompts",
		func(prompt string) {
			Expect(c.Classify(prompt)).To(Equal(classifier.IntentTask))
		},
		Entry("validate alert uuid",
			"validate alert 550e8400-e29b-41d4-a716-446655440000"),
		Entry("check alert",
			"check alert 123e4567-e89b-12d3-a456-426614174000"),
		Entry("what is the system status",
			"what is the current system status?"),
		Entry("is the service healthy",
			"is the auth-service healthy right now?"),
		Entry("alert id provided",
			"please validate 550e8400-e29b-41d4-a716-446655440000"),
		Entry("ok now validate — mixed intent should be task",
			"ok, now validate alert 550e8400-e29b-41d4-a716-446655440000"),
	)
})
