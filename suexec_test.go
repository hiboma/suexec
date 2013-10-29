package suexec

import (
	. "github.com/r7kamura/gospel"
	"testing"
)

func TestConstant(t *testing.T) {
	Describe(t, "cleanupEnv", func() {

		It("AP_HTTPD_USER", func() {
			Expect(AP_HTTPD_USER).To(Exist)
		})

		It("AP_UID_MIN", func() {
			Expect(AP_UID_MIN).To(Exist)
		})

		It("AP_GID_MIN", func() {
			Expect(AP_GID_MIN).To(Exist)
		})

		It("AP_DOC_ROOT", func() {
			Expect(AP_DOC_ROOT).To(Exist)
		})

		It("AP_LOG_EXEC", func() {
			Expect(AP_LOG_EXEC).To(Exist)
		})

	})

	Describe(t, "IsUserdirEnabled", func() {
		It("If username has '~' for prefix, Userdir is enabled", func() {
			Expect(IsUserdirEnabled("~namahage")).To(Equal, true)
		})

		It("Userdir is diabled", func() {
			Expect(IsUserdirEnabled("namahage")).To(Equal, false)
		})
	})

	Describe(t, "NewError", func() {
		It("NewError", func() {
			err := NewError(0, "error is %s", "one")
			Expect(err.status).To(Equal, 0)
			Expect(err.message).To(Equal, "error is one")
		})

		It("NewError", func() {
			err := NewError(0, "error is %s, %s", "one", "two")
			Expect(err.status).To(Equal, 0)
			Expect(err.message).To(Equal, "error is one, two")
		})

		It("NewError", func() {
			err := NewError(100, "error is %s", "one")
			Expect(err.status).To(Equal, 100)
			Expect(err.message).To(Equal, "error is one")
		})
	})
}
