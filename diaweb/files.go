package diaweb

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// Method returns method to call to download the file
func (m *MirrorFile) Method(local, name string) func() error {
	switch {
	default:
		return func() error {
			return m.localCopy(local, name)
		}
	}
}

func (m *MirrorFile) localCopy(ldir, name string) error {
	filename := fmt.Sprintf("%s%s", name, path.Ext(m.Remote))
	m.Local = path.Join(ldir, filename) // real path
	dat, e := ioutil.ReadFile(m.Remote)
	if e != nil {
		fmt.Println("could not read", m.Remote, e)
		return e
	}
	ioutil.WriteFile(m.Local, dat, os.ModePerm)
	m.Local = path.Join("/tmp/", filename) // routed path
	return nil
}
