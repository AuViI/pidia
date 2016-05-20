package diaweb

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"
)

const (
	fileRegexpString = `([\S\.\w:]{1,}?\.\w*)\s*([\d]*)\s*?(\n|$)`
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
	c.CFiles = c.CFiles[0:1]

	err := c.addReadConfig(0)
	if err != nil {
		fmt.Println("Panic --- now: ", c.CFiles)
		panic(err)
	}
	defer func() {
		// Using the last time the server updated, combined with the amount of "wait time" set
		// in the updateLoop() function, to tell the client when to refresh.
		c.LastUpdate = time.Now()
	}()
	if read {
		rallconfig := 1
		for rallconfig < len(c.CFiles) {
			err := c.addReadConfig(rallconfig)
			if err != nil {
				fmt.Println("error reading", c.CFiles[rallconfig], err)
			}
			rallconfig++
		}
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
			os.Remove(path.Join(c.Directory, ent.Name()))
		}
	}
	for k := range c.Files {
		c.Files[k].Method(c.Directory, fmt.Sprintf("%04x", k))()
	}
}

func (c *Configuration) addReadConfig(cindex int) (err error) {
	configb, err := ioutil.ReadFile(c.CFiles[cindex])
	config := string(configb)
	cmatch := func(reg *regexp.Regexp) [][]string {
		return reg.FindAllStringSubmatch(config, -1)
	}
	for _, v := range cmatch(fileRegexp) {
		dur, err := strconv.ParseInt(v[2], 10, 32)
		if err != nil {
			dur = 10
		}
		if v[1][:7] == "http://" {
			c.Files = append(c.Files, MirrorFile{Remote: v[1], Duration: int(dur)})
			continue
		}
		if !path.IsAbs(v[1]) {
			v[1] = path.Join(path.Dir(c.CFiles[cindex]), v[1])
		}
		c.Files = append(c.Files, MirrorFile{Remote: v[1], Duration: int(dur)})
	}
	for _, v := range cmatch(dirRegexp) {
		dur, err := strconv.ParseInt(v[2], 10, 32)
		if err != nil {
			dur = 10
		}
		if !path.IsAbs(v[1]) {
			v[1] = path.Join(path.Dir(c.CFiles[cindex]), v[1])
		}
		addDir(&c.Files, v[1], int(dur), c)
	}
	return
}

func addDir(f *[]MirrorFile, dir string, dur int, c *Configuration) {
	dir = path.Clean(dir)
	info, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading dir:", dir, err)
		return
	}
fileLoop:
	for _, i := range info {
		ni := path.Join(dir, i.Name())
		if i.IsDir() {
			fmt.Println(ni)
			addDir(f, ni, dur, c)
		} else {
			if !(i.Name()[0] == '.' || i.Name()[0] == '~') {
				*f = append(*f, MirrorFile{Remote: ni, Duration: dur})
			} else if i.Name() == ".pidiarc" {
				for _, v := range c.CFiles {
					if v == ni {
						continue fileLoop
					} 
				}
				c.CFiles = append(c.CFiles, ni)
			}
		}
	}
}

// Show returns the html code for the web view
func (m *MirrorFile) Show() string {
	ex := path.Ext(m.Remote)
	switch ex {
	case ".png", ".jpeg", ".jpg", ".gif":
		return fmt.Sprintf("<img src=\"%s\">", m.Local)
	default:
		return fmt.Sprintf("<iframe src=\"%s\"></iframe>", m.Local)
	}
}
