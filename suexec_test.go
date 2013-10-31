package suexec

import (
	. "github.com/r7kamura/gospel"
	"os"
	"testing"
)

var saved_io_writer *os.File

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

		Before(func() {
			saved_io_writer = os.Stderr
			os.Stderr = nil
		})

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

		It("too free arguments if len(args) < 4", func() {
			p := Param{
				args: []string{"suexec", "501", "501"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant/misc",
				log:  log,
			}
			Expect(Suexec(p)).To(Equal, 101)
		})

		It("too free arguments if len(args) < 4", func() {
			p := Param{
				args: []string{"suexec", "501", "501"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant/misc",
				log:  log,
			}
			Expect(Suexec(p)).To(Equal, 101)
		})

		It("too free arguments if len(args) < 4", func() {
			p := Param{
				args: []string{"suexec"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant/misc",
				log:  log,
			}
			Expect(Suexec(p)).To(Equal, 101)
		})

		It("too free arguments if len(args) < 4", func() {
			p := Param{
				args: []string{"suexec", "501", "501"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant/misc",
				log:  log,
			}
			Expect(Suexec(p)).To(Equal, 101)
		})

		It("by non existent user(999) return 102", func() {
			p := Param{
				args: []string{"suexec", "501", "501", "index.pl"},
				uid:  999, /* who? */
				cwd:  "/vagrant",
				log:  log,
			}
			Expect(Suexec(p)).To(Equal, 102)
		})

		It("by nobody(99) return 103", func() {
			p := Param{
				args: []string{"suexec", "501", "501", "index.pl"},
				uid:  99, /* nobody */
				cwd:  "/vagrant",
				log:  log,
			}
			Expect(Suexec(p)).To(Equal, 103)
		})

		It("by vagrant(500)/???(999) return 106", func() {
			p := Param{
				args: []string{"suexec", "500", "999", "index.pl"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant",
				log:  log,
			}
			Expect(Suexec(p)).To(Equal, 106)
		})

		It("by root(0)/root(0) return cannot run as forbidden uid", func() {
			p := Param{
				args: []string{"suexec", "0", "0", "index.pl"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant",
				log:  log,
			}
			Expect(Suexec(p)).To(Equal, 107)
		})

		It("by root(0)/vagrant(501) return 107", func() {
			p := Param{
				args: []string{"suexec", "0", "501", "index.pl"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant",
				log:  log,
			}
			Expect(Suexec(p)).To(Equal, 107)
		})

		It("by vagrant(501)/root(0) return 108", func() {
			p := Param{
				args: []string{"suexec", "501", "0", "index.pl"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant",
				log:  log,
			}
			Expect(Suexec(p)).To(Equal, 108)
		})

		It("by vagrant(501)/nobody(99) return 108", func() {
			p := Param{
				args: []string{"suexec", "501", "99", "index.pl"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant",
				log:  log,
			}
			Expect(Suexec(p)).To(Equal, 108)
		})

		It("by ???(999)/root(0) return 12", func() {
			p := Param{
				args: []string{"suexec", "999", "0", "index.pl"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant",
				log:  log,
			}
			Expect(Suexec(p)).To(Equal, 121)
		})

		It("not-exists-command return command not in docroot", func() {
			p := Param{
				args: []string{"suexec", "501", "501", "not-exists-command"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant",
				log:  log,
			}
			Expect(Suexec(p)).To(Equal, 114)
		})

		After(func() {
			os.Stderr = saved_io_writer
		})
	})
}
