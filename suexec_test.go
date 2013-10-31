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

	Describe(t, "NewSuexecError", func() {
		It("NewSuexecError", func() {
			err := NewSuexecError(0, "error is %s", "one")
			Expect(err.status).To(Equal, 0)
			Expect(err.message).To(Equal, "error is one")
		})

		It("NewSuexecError", func() {
			err := NewSuexecError(0, "error is %s, %s", "one", "two")
			Expect(err.status).To(Equal, 0)
			Expect(err.message).To(Equal, "error is one, two")
		})

		It("NewSuexecError", func() {
			err := NewSuexecError(100, "error is %s", "one")
			Expect(err.status).To(Equal, 100)
			Expect(err.message).To(Equal, "error is one")
		})
	})

	Describe(t, "Suexec", func() {
		log := NewLog("/dev/null")
		It("by root(0)/root(0) return cannot run as forbidden uid", func() {
			args := []string{"suexec", "0", "0", "index.pl"}
			Expect(Suexec(args, log)).To(Equal, 107)
		})

		It("by root(0)/vagrant/(501) return cannot run as forbidden gid", func() {
			args := []string{"suexec", "0", "501", "index.pl"}
			Expect(Suexec(args, log)).To(Equal, 107)
		})

		It("by vagrant(501)/root(0) return cannot run as forbidden gid", func() {
			args := []string{"suexec", "501", "0", "index.pl"}
			Expect(Suexec(args, log)).To(Equal, 108)
		})

		It("by ???(999)/root(0) return invalid target user", func() {
			args := []string{"suexec", "999", "0", "index.pl"}
			Expect(Suexec(args, log)).To(Equal, 121)
		})

		It("by vagrant(501)/???(999) return invalid target group name", func() {
			args := []string{"suexec", "501", "999", "index.pl"}
			Expect(Suexec(args, log)).To(Equal, 106)
		})

		It("not-exists-command return command not in docroot", func() {
			args := []string{"suexec", "501", "501", "not-exists-command"}
			Expect(Suexec(args, log)).To(Equal, 114)
		})

	})
}
