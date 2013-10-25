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

	Describe(t, "IsValidCommand", func() {
		It("index.pl is valid command", func() {
			Expect(IsValidCommand("index.pl")).To(Equal, true)
		})

		It("bin/index.pl is valid command", func() {
			Expect(IsValidCommand("bin/index.pl")).To(Equal, true)
		})

		It("../index.pl is invalid command", func() {
			Expect(IsValidCommand("../index.pl")).To(Equal, false)
		})

		It("/index.pl is invalid command", func() {
			Expect(IsValidCommand("/index.pl")).To(Equal, false)
		})

		It("bin/../index.pl is invalid command", func() {
			Expect(IsValidCommand("bin/../index.pl")).To(Equal, false)
		})
	})

	Describe(t, "IsUserdirEnabled", func() {
		It("If username has '~' for prefix, Userdir is enabled", func() {
			Expect(IsValidCommand("~namahage")).To(Equal, true)
		})

		It("Userdir is diabled", func() {
			Expect(IsValidCommand("namahage")).To(Equal, true)
		})
	})
}
