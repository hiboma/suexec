package group

import (
	. "github.com/r7kamura/gospel"
	"testing"
)

func TestSuexecGroup(t *testing.T) {
	Describe(t, "LookupGidAndName", func() {
		It("lookup 'root' should should return gid, gname", func() {
			gid, groupname, err := LookupGidAndName("root")
			Expect(gid).To(Equal, "0")
			Expect(groupname).To(Equal, "root")
			Expect(err).To(Equal, nil)
		})

		It("lookup '1234567890' should return gid, gname", func() {
			gid, groupname, err := LookupGidAndName("1234567890")
			Expect(gid).To(Equal, "1234567890")
			Expect(groupname).To(Equal, "1234567890")
			Expect(err).To(Equal, nil)
		})
	})
}
