package suexec

import (
	. "github.com/r7kamura/gospel"
	"io/ioutil"
	"os"
	"testing"
)

func TestRubyVersion(t *testing.T) {

	Describe(t, "rubyVersionPath", func() {
		It("rubyVersionPath", func() {
			Expect(pathToRubyVersion("/home/vagrant")).To(Equal, "/home/vagrant/.ruby-version")
		})
	})

	Describe(t, "ReadRubyVersion", func() {

		tempdir, _ := ioutil.TempDir("", "")
		defer os.RemoveAll(tempdir)
		ruby_version := tempdir + "/.ruby-version"

		It("ReadRubyVersion", func() {
			ioutil.WriteFile(ruby_version, []byte("2.0.0-p247"), 0644)
			version, err := ReadRubyVersion(tempdir)
			Expect(err).To(Equal, nil)
			Expect(version).To(Equal, "2.0.0-p247")
		})

		It("ReadRubyVersion with EOL", func() {
			ioutil.WriteFile(ruby_version, []byte("2.0.0-p247\n"), 0644)
			version, err := ReadRubyVersion(tempdir)
			Expect(err).To(Equal, nil)
			Expect(version).To(Equal, "2.0.0-p247")
		})

		It("ReadRubyVersion with leading/trailing white space", func() {
			ioutil.WriteFile(ruby_version, []byte("  2.0.0-p247  "), 0644)
			version, err := ReadRubyVersion(tempdir)
			Expect(err).To(Equal, nil)
			Expect(version).To(Equal, "2.0.0-p247")
		})

		It("ReadRubyVersion with duplicated", func() {
			ioutil.WriteFile(ruby_version, []byte("2.0.0-p247\n1.9.3-p0"), 0644)
			version, err := ReadRubyVersion(tempdir)
			Expect(err).To(Equal, nil)
			Expect(version).To(Equal, "2.0.0-p247")
		})

		It("modifyPathEnv", func() {
			environ := InjectRubyPathEnv("/usr/local/ruby/2.0.0-p247/bin", []string{"PATH=/bin"})
			Expect(environ[0]).To(Equal, "PATH=/usr/local/ruby/2.0.0-p247/bin:/bin")
		})
	})
}
