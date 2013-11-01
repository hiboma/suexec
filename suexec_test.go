package suexec

import (
	. "github.com/r7kamura/gospel"
	"os"
	"testing"
)

var saved_io_writer *os.File
var original_cwd string

func init() {
	original_cwd, _ = os.Getwd()
}

func TestConstant(t *testing.T) {

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

		Before(func() {
			saved_io_writer = os.Stderr
			os.Stderr = nil
			os.Chdir(original_cwd)
		})

		It("too free arguments if len(args) < 4", func() {
			err := Suexec(Param{
				args: []string{"suexec", "501", "501"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant/misc",
				log:  log,
			})
			Expect(err.status).To(Equal, 101)
			Expect(err.message).To(Equal, "too few arguments")
		})

		It("too free arguments if len(args) < 4", func() {
			err := Suexec(Param{
				args: []string{"suexec", "501"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant/misc",
				log:  log,
			})
			Expect(err.status).To(Equal, 101)
			Expect(err.message).To(Equal, "too few arguments")
		})

		It("too free arguments if len(args) < 4", func() {
			err := Suexec(Param{
				args: []string{"suexec"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant/misc",
				log:  log,
			})
			Expect(err.status).To(Equal, 101)
			Expect(err.message).To(Equal, "too few arguments")
		})

		It("by non existent user(999) return 102", func() {
			err := Suexec(Param{
				args: []string{"suexec", "501", "501", "index.pl"},
				uid:  999, /* who? */
				cwd:  "/vagrant",
				log:  log,
			})
			Expect(err.status).To(Equal, 102)
			Expect(err.message).To(Equal, "crit: invalid uid: (999) user: unknown userid 999\n")
		})

		It("by nobody(99) return 103", func() {
			err := Suexec(Param{
				args: []string{"suexec", "501", "501", "index.pl"},
				uid:  99, /* nobody */
				cwd:  "/vagrant",
				log:  log,
			})
			Expect(err.status).To(Equal, 103)
			Expect(err.message).To(Equal, "user mismatch (nobody instead os vagrant)")
		})

		It("by vagrant(500)/???(999) return 106", func() {
			err := Suexec(Param{
				args: []string{"suexec", "500", "999", "index.pl"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant",
				log:  log,
			})
			Expect(err.status).To(Equal, 106)
			Expect(err.message).To(Equal, "invalid target group name: (999)\n")
		})

		It("by root(0)/root(0) return cannot run as forbidden uid", func() {
			err := Suexec(Param{
				args: []string{"suexec", "0", "0", "index.pl"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant",
				log:  log,
			})
			Expect(err.status).To(Equal, 107)
			Expect(err.message).To(Equal, "cannot run as forbidden uid (0/index.pl)\n")
		})

		It("by root(0)/vagrant(501) return 107", func() {
			err := Suexec(Param{
				args: []string{"suexec", "0", "501", "index.pl"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant",
				log:  log,
			})
			Expect(err.status).To(Equal, 107)
			Expect(err.message).To(Equal, "cannot run as forbidden uid (0/index.pl)\n")
		})

		It("by vagrant(501)/root(0) return 108", func() {
			err := Suexec(Param{
				args: []string{"suexec", "501", "0", "index.pl"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant",
				log:  log,
			})
			Expect(err.status).To(Equal, 108)
			Expect(err.message).To(Equal, "cannot run as forbidden gid (0/index.pl)\n")
		})

		It("by vagrant(501)/nobody(99) return 108", func() {
			err := Suexec(Param{
				args: []string{"suexec", "501", "99", "index.pl"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant",
				log:  log,
			})
			Expect(err.status).To(Equal, 108)
			Expect(err.message).To(Equal, "cannot run as forbidden gid (99/index.pl)\n")
		})

		It("by ???(999)/root(0) return 12", func() {
			err := Suexec(Param{
				args: []string{"suexec", "999", "0", "index.pl"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant",
				log:  log,
			})
			Expect(err.status).To(Equal, 121)
			Expect(err.message).To(Equal, "invalid target user: (999)\n")
		})

		It("not-exists-command return command not in docroot", func() {
			err := Suexec(Param{
				args: []string{"suexec", "501", "501", "not-exists-command"},
				uid:  501, /* vagrant */
				cwd:  "/vagrant",
				log:  log,
			})
			Expect(err.status).To(Equal, 114)
			Expect(err.message).To(Equal, "command not in docroot (/vagrant/not-exists-command)\n")
		})

		After(func() {
			os.Stderr = saved_io_writer
		})
	})
}
