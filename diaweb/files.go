package diaweb

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

// Method returns method to call to download the file
// local being the directory to download into, name being the name of
// the temporary file.
// The remote address is read from m.Remote
// After calling the returned function (Closure) m.Local will be set up
// correctly
func (m *MirrorFile) Method(local, name string) func() error {
	switch {
	case m.Remote[:7] == "http://", m.Remote[:8] == "https://":
		// File can be retrieved via http
		return func() error {
			return m.httpDownload(local, name)
		}
	default: // assume that file can be copied
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

func (m *MirrorFile) httpDownload(ldir, name string) error {
	newfile := fmt.Sprintf("%s%s", name, path.Ext(m.Remote))
	abs := path.Join(ldir, newfile)
	resp, err := http.Get(m.Remote) // http Get-Request
	if err != nil {
		fmt.Println("could not fetch", m.Remote, err)
		return err
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)     // read Response Body
	ioutil.WriteFile(abs, data, os.ModePerm) // write to file
	m.Local = path.Join("/tmp/", newfile)
	return nil
}
