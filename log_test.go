package suexec

import (
	"fmt"
	. "github.com/r7kamura/gospel"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

func TestSuexecLog(t *testing.T) {
	Describe(t, "NewLog", func() {
		It("logErr write valid format", func() {
			tempdir, _ := ioutil.TempDir("", "suexec")
			defer os.RemoveAll(tempdir)
			path := fmt.Sprintf("%s/suexec.log", tempdir)

			log := NewLog(path)
			log.LogNoErr("%s %s %s", "this", "is", "test")

			file, _ := os.Open(path)
			bytes, _ := ioutil.ReadAll(file)
			matched, _ := regexp.Match(`\[\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\]: this is test`, bytes)
			Expect(matched).To(Equal, true)
		})
	})
}
