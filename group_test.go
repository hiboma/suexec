package suexec

import (
	. "github.com/r7kamura/gospel"
	"os/user"
	"testing"
)

func TestSuexecGroup(t *testing.T) {
	Describe(t, "LookupGidAndName", func() {
		It("lookup 'root' should should return gid, gname", func() {
			gr, err := LookupGroup("root")
			Expect(gr.Gid).To(Equal, "0")
			Expect(gr.Name).To(Equal, "root")
			Expect(err).To(Equal, nil)
		})

		It("lookup '0' should should return gid, gname", func() {
			gr, err := LookupGroup("0")
			Expect(gr.Gid).To(Equal, "0")
			Expect(gr.Name).To(Equal, "root")
			Expect(err).To(Equal, nil)
		})

		It("lookup '1234567890' should return nil", func() {
			_, err := LookupGroup("1234567890")
			Expect(err.(user.UnknownGroupIdError)).To(Exist)
		})

		It("lookup 'hogehogehoge' should return nil", func() {
			_, err := LookupGroup("hogehogehoge")
			Expect(err.(user.UnknownGroupError)).To(Exist)
		})
	})
}
