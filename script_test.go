package suexec

import (
	. "github.com/r7kamura/gospel"
	"io/ioutil"
	"os"
	"testing"
)

func TestSuexecScript(t *testing.T) {

	Describe(t, "NewScript", func() {
		It("should be success when a path exists", func() {
			s, err := NewScript("index.pl", ".")
			Expect(s).To(Exist)
			Expect(err == nil).To(Exist)
		})

		It("should not be success when a script not exists", func() {
			s, err := NewScript("script not found", ".")
			Expect(s == nil).To(Exist)
			Expect(os.IsExist(err)).To(Equal, false)
		})
	})

	Describe(t, "IsSetuid", func() {
		It("script is not setuid", func() {
			s, _ := NewScript("index.pl", ".")
			Expect(s.IsSetuid()).To(Equal, false)
			Expect(s.IsSetgid()).To(Equal, false)
		})

		It("script is setuid", func() {
			script, _ := ioutil.TempFile("", "")
			script.Chmod(0700 | os.ModeSetuid)
			defer os.Remove(script.Name())

			s, _ := NewScript(script.Name(), ".")
			Expect(s.IsSetuid()).To(Equal, true)
		})
	})

	Describe(t, "IsSetgid", func() {

		It("script is not setgid", func() {
			script, _ := ioutil.TempFile("", "")
			script.Chmod(700)
			defer os.Remove(script.Name())

			s, _ := NewScript(script.Name(), ".")
			Expect(s.IsSetgid()).To(Equal, false)
		})

		It("script is setgid", func() {
			script, _ := ioutil.TempFile("", "")
			script.Chmod(700 | os.ModeSetgid)
			defer os.Remove(script.Name())

			s, _ := NewScript(script.Name(), ".")
			Expect(s.IsSetgid()).To(Equal, true)
		})
	})

	Describe(t, "IfOwnerMatch", func() {
		It("owner match with specified uid, gid", func() {
			s, _ := NewScript("index.pl", ".")
			uid := os.Getuid()
			gid := os.Getgid()
			Expect(s.IfOwnerMatch(uid, gid)).To(Equal, true)
		})

		It("owner does not match with specified uid, gid", func() {
			s, _ := NewScript("index.pl", ".")
			uid := os.Getuid() + 1000
			gid := os.Getgid() + 1000
			Expect(s.IfOwnerMatch(uid, gid)).To(Equal, false)
		})
	})

	Describe(t, "NewScript", func() {
		It("index.pl is secure", func() {
			s, err := NewScript("index.pl", ".")
			Expect(s).To(NotEqual, nil)
			Expect(err).To(Equal, nil)
			Expect(s.path_info).To(NotEqual, nil)
			Expect(s.cwd_info).To(NotEqual, nil)
		})

		It("index.pl is secure", func() {
			s, err := NewScript("not exists", ".")
			Expect(s).To(Exist)
			Expect(err).To(Exist)
		})

		It("index.pl is secure", func() {
			s, err := NewScript("index.pl", "./1234567890")
			Expect(s).To(Exist)
			Expect(err).To(Exist)
		})
	})

	Describe(t, "HasSecurePath", func() {
		It("index.pl is secure", func() {
			s := Script{path: "index.pl", cwd: "."}
			Expect(s.HasSecurePath()).To(Equal, true)
		})

		It("bin/index.pl is secure", func() {
			s := Script{path: "bin/index.pl", cwd: "."}
			Expect(s.HasSecurePath()).To(Equal, true)
		})

		It("../index.pl is insecure", func() {
			s := Script{path: "..//index.pl", cwd: "."}
			Expect(s.HasSecurePath()).To(Equal, false)
		})

		It("/index.pl is insecure path", func() {
			s := Script{path: "/index.pl", cwd: "."}
			Expect(s.HasSecurePath()).To(Equal, false)
		})

		It("bin/../index.pl is insecure path", func() {
			s := Script{path: "bin/../index.pl", cwd: "."}
			Expect(s.HasSecurePath()).To(Equal, false)
		})
	})

	Describe(t, "IsWritableByOthers", func() {
		It("script is not writable by others", func() {
			script, _ := ioutil.TempFile("", "")
			defer os.Remove(script.Name())
			script.Chmod(0700)

			s, _ := NewScript(script.Name(), ".")
			Expect(s.IsWritableByOthers()).To(Equal, false)
		})

		It("script is writable by others", func() {

			insecure_modes := []os.FileMode{0720, 0702, 0730, 0703, 0760, 0706, 0766, 0777}

			for _, mode := range insecure_modes {
				script, _ := ioutil.TempFile("", "")
				defer os.Remove(script.Name())
				script.Chmod(mode)

				s, _ := NewScript(script.Name(), ".")
				Expect(s.IsWritableByOthers()).To(Equal, true)
			}
		})
	})

	Describe(t, "IsExecutable", func() {

		It("script is executable", func() {
			script, _ := ioutil.TempFile("", "")
			defer os.Remove(script.Name())
			script.Chmod(0700)

			s, _ := NewScript(script.Name(), ".")
			Expect(s.IsExecutable()).To(Equal, true)
		})

		It("script is not executable", func() {
			script, _ := ioutil.TempFile("", "")
			defer os.Remove(script.Name())
			script.Chmod(0600)

			s, _ := NewScript(script.Name(), ".")
			Expect(s.IsExecutable()).To(Equal, false)
		})

	})

	Describe(t, "IsDirWritableByOthers", func() {
		It("directory is not writable by others", func() {
			tempdir, _ := ioutil.TempDir("", "")
			defer os.RemoveAll(tempdir)
			os.Chmod(tempdir, 0700)

			script, _ := ioutil.TempFile(tempdir, "")
			s, _ := NewScript(script.Name(), tempdir)
			Expect(s.IsDirWritableByOthers()).To(Equal, false)
		})

		It("directory is writable by others", func() {

			insecure_modes := []os.FileMode{0720, 0702, 0730, 0703, 0760, 0706, 0766, 0777}
			for _, mode := range insecure_modes {
				tempdir, _ := ioutil.TempDir("", "")
				defer os.RemoveAll(tempdir)
				os.Chmod(tempdir, mode)

				script, _ := ioutil.TempFile(tempdir, "")
				s, _ := NewScript(script.Name(), tempdir)
				Expect(s.IsDirWritableByOthers()).To(Equal, true)
			}
		})
	})
}
