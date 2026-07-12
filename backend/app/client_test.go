package app

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("App Package Internal Elements", func() {
	Context("appClient initialization", func() {
		It("should initialize newAppClient with default settings", func() {
			client := newAppClient()
			Expect(client).NotTo(BeNil())
			Expect(client.httpClient).NotTo(BeNil())
			Expect(client.ollamaModel).To(Equal("llama3.2"))
			Expect(client.ollamaBase).To(Equal("http://localhost:11434"))
		})
	})

	Context("appRepository initialization", func() {
		It("should initialize NewAppRepository correctly", func() {
			client := newAppClient()
			repo := NewAppRepository(client, nil, nil)
			Expect(repo).NotTo(BeNil())
			Expect(repo.Client).To(Equal(client))
		})
	})
})
