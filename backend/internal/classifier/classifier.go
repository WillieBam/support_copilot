package classifier

import (
	"regexp"
	"strings"
)

// Intent labels the nature of a user prompt.
type Intent string

const (
	// IntentConversational indicates a social/acknowledgement message that
	// does not require tool execution (e.g. "ok", "thanks", "bye")
	IntentConversational Intent = "conversational"

	// IntentTask indicates the user is requesting an operation that may
	// require tool execution (e.g. providing an alert ID for validation).
	IntentTask Intent = "task"
)

// uuidPattern matches a standard v4-style UUID string
var uuidPattern = regexp.MustCompile(
	`(?i)[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`,
)

// conversationalPatterns is the ordered list of regex rules used to detect
// conversational prompts.
//
// Design constraints:
//   - Each pattern MUST be anchored to the full trimmed string (^ and $)
//   - Trailing content is limited to short social words (≤20 chars) to
//     prevent false-positive matches on mixed prompts like
//     e.g. "ok, now validate alert <uuid>"
var conversationalPatterns = []*regexp.Regexp{
	// pattern1: ok, okay, alright, got it, sure, fine, noted
	regexp.MustCompile(`(?i)^(ok(ay)?|alright|got\s+it|sure|fine|noted)(\s+[\w\s'!,.]{0,20})?$`),
	// pattern2: bye, goodbye, see you, cya, ttyl, later
	regexp.MustCompile(`(?i)^(bye(bye)?|goodbye|see\s+you|cya|ttyl|later)(\s+[\w\s'!,.]{0,20})?$`),
	// pattern3: thanks, thank you, ty, cheers, much appreciated
	regexp.MustCompile(`(?i)^(thanks?|thank\s+you|ty|cheers|much\s+appreciated)(\s+[\w\s'!,.]{0,20})?$`),
	// pattern4: hi, hello, hey, halo, hei, yo + optional trailing social phrase
	regexp.MustCompile(`(?i)^(hi+|h[ae]llo+|hey+|halo+|hei|yo)(\s+[\w\s'!,.?]{0,30})?$`),
	// pattern5: Good <time-of-day> greetings
	regexp.MustCompile(`(?i)^good\s+(morning|afternoon|evening|day)(\s+[\w\s'!,.]{0,20})?$`),
	// pattern6: yes/no
	regexp.MustCompile(`(?i)^(yes|no|nope|yep|yeah|nah|yup)(\s+[\w\s'!,.]{0,20})?$`),
	// pattern6: done, stop, quit, exit, finish, that's all
	regexp.MustCompile(`(?i)^(done|stop|quit|exit|finish(ed)?|that'?s?\s+all)(\s+[\w\s'!,.]{0,20})?$`),
	// pattern7: wellness / small-talk questions: "are you ok?", "how are you?", "you good?"
	regexp.MustCompile(`(?i)^(are\s+you\s+ok\??|how\s+are\s+you\??|you\s+good\??|you\s+ok\??)$`),
}

// embeddedToolCallPattern detects when the LLM emits a raw JSON tool-call
// object as plain text content instead of via the proper tool_calls mechanism.
// Example: {"name": "greet", "parameters": {"message": "I"}}
var embeddedToolCallPattern = regexp.MustCompile(
	`(?s)^\s*\{\s*"(name|function)"\s*:\s*"[^"]+"\s*,\s*"(parameters|arguments)"\s*:\s*\{`,
)

// LooksLikeEmbeddedToolCall returns true when content appears to be a raw
// JSON tool-call emitted by the LLM as text rather than through the
// tool_calls field. Such content should be suppressed and a fallback triggered.
func LooksLikeEmbeddedToolCall(content string) bool {
	return embeddedToolCallPattern.MatchString(strings.TrimSpace(content))
}

// IntentClassifier classifies a prompt as conversational or task-oriented
// using regex pattern matching. It implements the interfaces.IIntentClassifier
// interface.
type IntentClassifier struct{}

// NewIntentClassifier returns a new IntentClassifier.
func NewIntentClassifier() *IntentClassifier {
	return &IntentClassifier{}
}

// Classify check the prompt and returns the detected Intent.
// It returns IntentConversational when the prompt:
// 1. match any known social / acknowledgement pattern, or
// 2. is a short message (≤80 chars) that contains no UUID and none of the
// known task keywords (validate, alert, check, incident, inspect).
//
// IntentTask is returned otherwise.
func (c *IntentClassifier) Classify(prompt string) Intent {
	s := strings.TrimSpace(prompt)

	// check explicit conversational patterns first
	for _, re := range conversationalPatterns {
		if re.MatchString(s) {
			return IntentConversational
		}
	}

	//  short-message: messages under 80 chars with no UUID and none of the known task trigger keywords are treated as conversational
	if len(s) <= 80 && !uuidPattern.MatchString(s) && !containsTaskKeyword(s) {
		return IntentConversational
	}

	return IntentTask
}

// taskKeywords are lowercase substrings whose presence in the prompt strongly
// signals a task request that may need tool execution
var taskKeywords = []string{
	"validate", "alert", "incident", "check", "inspect", "anomaly",
	"metric", "status", "monitor", "error", "failure", "cpu", "memory",
	"healthy", "health", "service", "latency", "throughput", "outage",
}

func containsTaskKeyword(s string) bool {
	lower := strings.ToLower(s)
	for _, kw := range taskKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}
