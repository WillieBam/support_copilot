package classifier

import (
	"regexp"
	"strings"
)

// Intent labels the nature of a user prompt.
type Intent string

const (
	// IntentConversational indicates a social/acknowledgement message that
	// does not require tool execution (e.g. "ok", "thanks", "bye").
	IntentConversational Intent = "conversational"

	// IntentTask indicates the user is requesting an operation that may
	// require tool execution (e.g. providing an alert ID for validation).
	IntentTask Intent = "task"
)

// conversationalPatterns is the ordered list of regex rules used to detect
// conversational prompts.
//
// Design constraints:
//   - Each pattern MUST be anchored to the full trimmed string (^ and $).
//     "ok, now validate alert <uuid>".
var conversationalPatterns = []*regexp.Regexp{
	// Acknowledgements: ok, okay, alright, got it, sure, fine, noted
	// Trailing part must be short (≤20 chars) so a UUID after the keyword fails.
	regexp.MustCompile(`(?i)^(ok(ay)?|alright|got\s+it|sure|fine|noted)(\s+[\w\s'!,.]{0,20})?$`),
	// Sign-offs: bye, goodbye, see you, cya, ttyl, later
	regexp.MustCompile(`(?i)^(bye(bye)?|goodbye|see\s+you|cya|ttyl|later)(\s+[\w\s'!,.]{0,20})?$`),
	// Gratitude: thanks, thank you, ty, cheers, much appreciated
	regexp.MustCompile(`(?i)^(thanks?|thank\s+you|ty|cheers|much\s+appreciated)(\s+[\w\s'!,.]{0,20})?$`),
	// Greetings: hi, hello, hey, good morning/afternoon/evening
	regexp.MustCompile(`(?i)^(hi|hello|hey|good\s+(morning|afternoon|evening|day))(\s+[\w\s'!,.]{0,20})?$`),
	// Simple yes/no
	regexp.MustCompile(`(?i)^(yes|no|nope|yep|yeah|nah|yup)(\s+[\w\s'!,.]{0,20})?$`),
	// Completion signals: done, stop, quit, exit, finish, that's all
	regexp.MustCompile(`(?i)^(done|stop|quit|exit|finish(ed)?|that'?s?\s+all)(\s+[\w\s'!,.]{0,20})?$`),
}

// IntentClassifier classifies a prompt as conversational or task-oriented
// using regex pattern matching. It implements the interfaces.IIntentClassifier
// interface.
type IntentClassifier struct{}

// NewIntentClassifier returns a new IntentClassifier.
func NewIntentClassifier() *IntentClassifier {
	return &IntentClassifier{}
}

// Classify inspects the prompt and returns the detected Intent.
// It returns IntentConversational when the prompt matches any known social /
// acknowledgement pattern; IntentTask otherwise.
func (c *IntentClassifier) Classify(prompt string) Intent {
	s := strings.TrimSpace(prompt)
	for _, re := range conversationalPatterns {
		if re.MatchString(s) {
			return IntentConversational
		}
	}
	return IntentTask
}
