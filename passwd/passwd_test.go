package passwd

import (
	. "github.com/r7kamura/gospel"
	"testing"
)

func TestConstant(t *testing.T) {
	Describe(t, "Lookup", func() {
		It("lookup 'root' should exists", func() {
			pw, err := Lookup("root")
			Expect(pw).To(Exist)
			Expect(err).To(Equal, nil)
		})

		It("lookup '1234567890' should not exists", func() {
			pw, err := Lookup("1234567890")
			Expect(pw).To(Exist)
			Expect(err).To(NotEqual, nil)
		})
	})
}
