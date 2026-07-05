package firebase_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFirebase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Firebase Repository Suite")
}
