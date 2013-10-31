package suexec

import (
	. "github.com/r7kamura/gospel"
	"testing"
)

func TestSuexecError(t *testing.T) {
	Describe(t, "NewSuexecError", func() {
		It("should return NewSuexecError", func() {
			error := NewSuexecError(1, "%s", "hoge")
			Expect(error.status).To(Equal, 1)
			Expect(error.message).To(Equal, "hoge")
		})
	})
}
