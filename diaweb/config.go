package diaweb

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
)

const (
	fileRegexpString = `([\S\.\w]{1,}?\.\w*)\s*([\d]*)\s*?(\n|$)`
	dirRegexpString  = `([\/\.\w]*?\/\w*)\s*([\d]*)\s*?(\n|$)`
)

var (
	fileRegexp, _ = regexp.Compile(fileRegexpString)
	dirRegexp, _  = regexp.Compile(dirRegexpString)
)

// Query goes through the config files, updating the index
func (c *Configuration) Query(read bool) {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()
	c.Files = *new([]MirrorFile)
	if read { // recursive read

		return
	}
	// Not recursive read
	configb, err := ioutil.ReadFile(c.CFiles[0])
	if err != nil {
		fmt.Println("no config file found")
		os.Exit(1)
	}
	config := string(configb)
	cmatch := func(reg *regexp.Regexp) [][]string {
		return reg.FindAllStringSubmatch(config, -1)
	}
	if err != nil {
		panic(err)
	}
	for _, v := range cmatch(fileRegexp) {
		dur, err := strconv.ParseInt(v[2], 10, 32)
		if err != nil {
			dur = 10
		}
		c.Files = append(c.Files, MirrorFile{Remote: v[1], Duration: int(dur)})
	}
	for _, v := range cmatch(dirRegexp) {
		dur, err := strconv.ParseInt(v[2], 10, 32)
		if err != nil {
			dur = 10
		}
		addDir(&c.Files, v[1], int(dur))
	}
	entries, lerr := ioutil.ReadDir(c.Directory)
	if lerr != nil {
		if os.IsNotExist(lerr) {
			os.Mkdir(c.Directory, os.ModePerm)
			entries = nil
		} else {
			panic(lerr)
		}
	}
	for _, ent := range entries {
		if !ent.IsDir() {
			os.Remove(path.Join(c.Directory,ent.Name()))
		}
	}
	for k := range c.Files {
		c.Files[k].Download(c.Directory, fmt.Sprintf("%04x", k))
	}
}

func addDir(f *[]MirrorFile, dir string, dur int) {
	dir = path.Clean(dir)
	info, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading dir:", dir, err)
		return
	}
	for _, i := range info {
		ni := path.Join(dir, i.Name())
		if i.IsDir() {
			fmt.Println(ni)
			addDir(f, ni, dur)
		} else {
			if !(i.Name()[0] == '.' || i.Name()[0] == '~') {
				*f = append(*f, MirrorFile{Remote: ni, Duration: dur})
			}
		}
	}
}

func (m *MirrorFile) Download(ldir string, name string) {
	// dl to /tmp
	filename := fmt.Sprintf("%s%s", name, path.Ext(m.Remote))
	m.Local = path.Join(ldir, filename)
	// cp /tmp$local /ldir$local

	// currently only local dirs/files allowed TODO online download
	dat, e := ioutil.ReadFile(m.Remote)
	if e != nil {
		fmt.Println("could not read", m.Remote,e)
		return
	}
	ioutil.WriteFile(m.Local, dat, os.ModePerm)
	m.Local = path.Join("/tmp/", filename)
}

func (m *MirrorFile) Show() string {
	ex := path.Ext(m.Remote)
	switch ex {
	case ".png", ".jpeg", ".jpg", ".gif":
		return fmt.Sprintf("<img src=\"%s\">", m.Local)
	default:
		return fmt.Sprintf("<iframe src=\"%s\"></iframe>", m.Local)
	}
}
