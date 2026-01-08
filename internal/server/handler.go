package server

import (
	"net/http"

	goauth "github.com/abbot/go-http-auth"
)

// FileServer wraps http.FileServer with the directory to serve
type FileServer struct {
	Dir string
}

// NewFileServer creates a new FileServer for the specified directory
func NewFileServer(dir string) *FileServer {
	return &FileServer{Dir: dir}
}

// Handle serves static files from the configured directory
func (fs *FileServer) Handle(w http.ResponseWriter, r *goauth.AuthenticatedRequest) {
	http.FileServer(http.Dir(fs.Dir)).ServeHTTP(w, &r.Request)
}
