package diaweb

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"sync"
	"text/template"
	"time"
)

type (
	// Server is GOD object for the local server
	Server struct {
		Host        string
		Port        uint
		Config      Configuration
		ReadConfigs bool
		wtime       int
	}
	// Configuration is data structure to save configs
	Configuration struct {
		LastUpdate time.Time
		Directory  string
		CFiles     []string
		Files      []MirrorFile
		RWMutex    *sync.RWMutex
	}
	// MirrorFile is a wrapper to allow easy syncing of files
	MirrorFile struct {
		Remote   string
		Duration int
		Local    string
	}
)

const (
	// DefaultWaitTime is the amount of seconds to wait in between config updates
	DefaultWaitTime = 600 // Seconds
)

// NewServer creates new Server instance
func NewServer(host string, port uint, dir string, config string, readc bool) *Server {
	s := &Server{
		Host: host,
		Port: port,
		Config: Configuration{
			Directory: dir,
			CFiles:    []string{config},
			Files:     nil,
			RWMutex:   new(sync.RWMutex),
		},
		ReadConfigs: readc,
		wtime:       DefaultWaitTime,
	}
	http.HandleFunc("/", s.Execute)
	http.HandleFunc("/r/", resourceHandler())
	http.HandleFunc("/tmp/", genTmpHandler(s))
	return s
}

func resourceHandler() func(http.ResponseWriter, *http.Request) {
	maincss, _ := template.ParseFiles(path.Join(GoPath, "src/github.com/auvii/pidia/diaweb/main.css"))
	return func(w http.ResponseWriter, r *http.Request) {
		f := r.URL.Path[3:]
		switch f {
		case "main.css":
			w.Header().Set("Content-type", "text/css")
			maincss.Execute(w, nil)
		default:
			fmt.Fprintf(w, "Not found %s", f)
		}
	}
}

// TODO fix genTmpHandler (2nd Mutex)
// currently it is still possible to mess up requesting files
//		1. read config #1
//		2. serve / request with config #1
//		3. s.Update() -> read config #2
//		4. serve /tmp/ request with config #2
//		result -> wrong pictures / wrong filetype
// Issue #1
func genTmpHandler(s *Server) func(http.ResponseWriter, *http.Request) {
	l := s.Config.Directory
	lock := s.Config.RWMutex
	return func(w http.ResponseWriter, r *http.Request) {
		lock.RLock() // Config can't change while handling image request
		defer lock.RUnlock()
		file := r.URL.Path[len("/tmp/"):]
		fh, err := os.Open(path.Join(l, file))
		if err != nil { // file not found / no permission
			fmt.Fprint(w, "Error", err)             // echo error to client
			fmt.Println("error request", file, err) // echo error to console
			return
		}
		defer fh.Close()
		io.Copy(w, fh)
	}
}

// Start starts configured Server instance
func (s *Server) Start() {
	go s.updateLoop()
	s.startLocal()
}

func (s *Server) startLocal() {
	port := fmt.Sprintf(":%d", s.Port)
	fmt.Println("> starting server on port", port)
	fmt.Println("> base config\n    ", s.Config.CFiles[0])
	fmt.Println("> mirror directory\n    ", s.Config.Directory)
	http.ListenAndServe(port, nil) // this is blocking
}

// Update reads the Config files & download files
func (s *Server) Update() {
	s.Config.Query(s.ReadConfigs)
}

func (s *Server) updateLoop() {
	acnum := 1                     // counter, TODO remove in production
	wsec := time.Duration(s.wtime) // time.Duration version of wtime
	for {
		fmt.Printf(" ... update Server #%d\n", acnum)
		acnum++
		s.Update()                     // re-read configuration, configure server
		time.Sleep(time.Second * wsec) // sleep #wsec seconds
	}
}

func minToSec(m int) int {
	return m * 60
}

// httpRefresh returns amount of seconds the client should wait for before refreshing
func (s *Server) httpRefresh() (wait int) {
	s.Config.RWMutex.RLock()
	defer s.Config.RWMutex.RUnlock()
	wait = minToSec(10)
	since := time.Since(s.Config.LastUpdate)
	wait = s.wtime - int(since.Seconds()) + int(5.0+rand.Float64()*20) // Client should refresh ~30 seconds after server-refresh
	return
}

// SetRefresh sets the refresh time intelligently to sync up with reading
// the config files correctly
// This is done to counteract issue #1
func (s *Server) SetRefresh(w http.ResponseWriter) {
	w.Header().Set("refresh", fmt.Sprintf("%d", s.httpRefresh()))
}

// GetRefresh is used in Templates, if SetRefresh does not work
func (s *Server) GetRefresh() int {
	return s.httpRefresh()
}
