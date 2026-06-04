package serve

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

type Server struct {
	addr    string
	docPath string
	mu      sync.RWMutex
	body    []byte
	subs    map[chan []byte]struct{}
}

func New(addr, docPath string) *Server {
	return &Server{
		addr:    addr,
		docPath: docPath,
		subs:    make(map[chan []byte]struct{}),
	}
}

func (s *Server) Set(body []byte) {
	s.mu.Lock()
	s.body = body
	subs := make([]chan []byte, 0, len(s.subs))
	for c := range s.subs {
		subs = append(subs, c)
	}
	s.mu.Unlock()
	for _, c := range subs {
		select {
		case c <- []byte("reload"):
		default:
		}
	}
}

func (s *Server) start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleDoc)
	mux.HandleFunc("/__reload", s.handleReload)
	return http.ListenAndServe(s.addr, mux)
}

func (s *Server) Start() error {
	return s.start()
}

func (s *Server) handleDoc(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	body := s.body
	s.mu.RUnlock()
	if body == nil {
		body = []byte("<!doctype html><meta charset=utf-8><title>mend</title><p>waiting for first build&hellip;")
	}
	if !strings.Contains(strings.ToLower(string(body)), "<head>") {
		body = injectReload([]byte(body))
	} else {
		body = injectReload(body)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Write(body)
}

func injectReload(body []byte) []byte {
	script := `<script>
(function(){
  const es = new EventSource("/__reload");
  es.onmessage = e => { if (e.data === "reload") location.reload(); };
  es.onerror = () => { setTimeout(() => location.reload(), 1500); };
})();
</script>`
	if strings.Contains(string(body), "</body>") {
		return []byte(strings.Replace(string(body), "</body>", script+"</body>", 1))
	}
	return []byte(string(body) + script)
}

func (s *Server) handleReload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", 500)
		return
	}
	ch := make(chan []byte, 4)
	s.mu.Lock()
	s.subs[ch] = struct{}{}
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		delete(s.subs, ch)
		s.mu.Unlock()
		close(ch)
	}()
	fmt.Fprintf(w, ": connected\n\n")
	flusher.Flush()
	notify := r.Context().Done()
	for {
		select {
		case <-notify:
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}
}

func DocName(path string) string {
	return filepath.Base(path)
}
