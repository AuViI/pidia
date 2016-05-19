package diaweb

import (
	"fmt"
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
	}
	// Configuration is data structure to save configs
	Configuration struct {
		Directory string
		CFiles    []string
		Files     []MirrorFile
		RWMutex   *sync.RWMutex
	}
	// MirrorFile is a wrapper to allow easy syncing of files
	MirrorFile struct {
		Remote   string
		Duration int
		Local    string
	}
)

// NewServer creates new Server instance
func NewServer(host string, port uint, dir string, config string, readc bool) *Server {
	s := &Server{
		Host: host,
		Port: port,
		Config: Configuration{
			Directory: path.Clean(fmt.Sprintf("%s/%s", os.Getenv("PWD"), dir)),
			CFiles:    []string{config},
			Files:     nil,
			RWMutex:   new(sync.RWMutex),
		},
		ReadConfigs: readc,
	}
	http.HandleFunc("/", s.Execute)
	http.HandleFunc("/r/", resourceHandler())
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

// Start starts configured Server instance
func (s *Server) Start() {
	go s.updateLoop()
	s.startLocal()
}

func (s *Server) startLocal() {
	port := fmt.Sprintf(":%d", s.Port)
	fmt.Println("starting server on port", port)
	http.ListenAndServe(port, nil)
}

// Update reads the Config files & download files
func (s *Server) Update() {
	s.Config.Query(s.ReadConfigs)
}

func (s *Server) updateLoop() {
	acnum := 1
	wsec := time.Duration(50)
	for {
		fmt.Printf("Update Server #%d", acnum)
		acnum++
		s.Update()
		time.Sleep(time.Second * wsec)
	}
}
