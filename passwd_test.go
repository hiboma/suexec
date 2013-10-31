package suexec

import (
	. "github.com/r7kamura/gospel"
	"os/user"
	"testing"
)

func TestSuexecPasswd(t *testing.T) {
	Describe(t, "Lookup", func() {
		It("lookup 'root' should exists", func() {
			pw, err := Lookup("root")
			Expect(pw).To(Exist)
			Expect(err).To(Equal, nil)
		})

		Describe(t, "IsUserdirEnabled", func() {
			It("If username has '~' for prefix, Userdir is enabled", func() {
				Expect(IsUserdirEnabled("~namahage")).To(Equal, true)
			})

			It("Userdir is diabled", func() {
				Expect(IsUserdirEnabled("namahage")).To(Equal, false)
			})
		})

		It("lookup 'hogehogehoge' should not exists", func() {
			pw, err := Lookup("hgoehogehgoe")
			Expect(pw).To(Exist)
			Expect(err.(user.UnknownUserError)).To(Exist)
		})

		It("lookup '1234567890' should not exists", func() {
			pw, err := Lookup("1234567890")
			Expect(pw).To(Exist)
			Expect(err.(user.UnknownUserIdError)).To(Exist)
		})
	})
}
