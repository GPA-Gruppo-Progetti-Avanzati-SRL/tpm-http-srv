package embedstatic

import (
	"embed"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

const INDEX = "index.html"

type ServeEmbeddedFilesystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}

type embeddedFileSystem struct {
	http.FileSystem
	indexes bool
}

func EmbedFile(fs embed.FS, indexes bool) *embeddedFileSystem {
	return &embeddedFileSystem{http.FS(fs), indexes}
}

func (l *embeddedFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		stats, err := os.Stat(p)
		if err != nil {
			return false
		}
		if stats.IsDir() {
			if !l.indexes {
				return false
			}
		}
		return true
	}
	return false
}

// Static returns a middleware handler that serves static files in the given directory.
func ServeEmbedded(urlPrefix string, fs ServeEmbeddedFilesystem) gin.HandlerFunc {
	fileserver := http.FileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {
		// if fs.Exists(urlPrefix, c.Request.URL.Path) {
		fileserver.ServeHTTP(c.Writer, c.Request)
		// c.Abort()
		// }
	}
}
