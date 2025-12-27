package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		writer, _ := gzip.NewWriterLevel(io.Discard, gzip.DefaultCompression)
		return writer
	},
}

type gzipResponseWriter struct {
	http.ResponseWriter
	writer        *gzip.Writer
	headerWritten bool
	minSize       int
	buffer        []byte
	statusCode    int
}

func newGzipResponseWriter(responseWriter http.ResponseWriter, minSize int) *gzipResponseWriter {
	gzWriter := gzipWriterPool.Get().(*gzip.Writer)
	gzWriter.Reset(responseWriter)

	return &gzipResponseWriter{
		ResponseWriter: responseWriter,
		writer:         gzWriter,
		minSize:        minSize,
		statusCode:     http.StatusOK,
	}
}

func (gzw *gzipResponseWriter) WriteHeader(statusCode int) {
	gzw.statusCode = statusCode
	if gzw.headerWritten {
		return
	}
}

func (gzw *gzipResponseWriter) Write(data []byte) (int, error) {
	if !gzw.headerWritten {
		gzw.buffer = append(gzw.buffer, data...)

		if len(gzw.buffer) >= gzw.minSize {
			gzw.flushBuffer(true)
		}
		return len(data), nil
	}

	return gzw.writer.Write(data)
}

func (gzw *gzipResponseWriter) flushBuffer(compress bool) {
	if gzw.headerWritten {
		return
	}
	gzw.headerWritten = true

	if compress && len(gzw.buffer) >= gzw.minSize {
		gzw.Header().Set("Content-Encoding", "gzip")
		gzw.Header().Del("Content-Length")
		gzw.ResponseWriter.WriteHeader(gzw.statusCode)
		gzw.writer.Write(gzw.buffer)
	} else {
		gzw.ResponseWriter.WriteHeader(gzw.statusCode)
		gzw.ResponseWriter.Write(gzw.buffer)
	}
}

func (gzw *gzipResponseWriter) Close() error {
	if !gzw.headerWritten && len(gzw.buffer) > 0 {
		gzw.flushBuffer(len(gzw.buffer) >= gzw.minSize)
	}

	if gzw.headerWritten && gzw.Header().Get("Content-Encoding") == "gzip" {
		err := gzw.writer.Close()
		gzipWriterPool.Put(gzw.writer)
		return err
	}

	gzipWriterPool.Put(gzw.writer)
	return nil
}

func (gzw *gzipResponseWriter) Flush() {
	if flusher, ok := gzw.ResponseWriter.(http.Flusher); ok {
		if gzw.headerWritten && gzw.Header().Get("Content-Encoding") == "gzip" {
			gzw.writer.Flush()
		}
		flusher.Flush()
	}
}

type CompressConfig struct {
	MinSize           int
	CompressionLevel  int
	ExcludedPaths     []string
	ExcludedMimeTypes []string
}

func DefaultCompressConfig() CompressConfig {
	return CompressConfig{
		MinSize:          1024,
		CompressionLevel: gzip.DefaultCompression,
		ExcludedPaths:    []string{"/metrics", "/health"},
		ExcludedMimeTypes: []string{
			"image/",
			"video/",
			"audio/",
			"application/octet-stream",
		},
	}
}

func Compress(config CompressConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if !shouldCompress(request, config) {
				next.ServeHTTP(writer, request)
				return
			}

			gzipWriter := newGzipResponseWriter(writer, config.MinSize)
			defer gzipWriter.Close()

			next.ServeHTTP(gzipWriter, request)
		})
	}
}

func shouldCompress(request *http.Request, config CompressConfig) bool {
	acceptEncoding := request.Header.Get("Accept-Encoding")
	if !strings.Contains(acceptEncoding, "gzip") {
		return false
	}

	for _, path := range config.ExcludedPaths {
		if strings.HasPrefix(request.URL.Path, path) {
			return false
		}
	}

	return true
}

func GzipCompress(next http.Handler) http.Handler {
	return Compress(DefaultCompressConfig())(next)
}
