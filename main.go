// Package main provides the go-vanity CLI tool, a basic server implementation
// capable of providing custom URLs to be used by the standard go tools.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func main() {
	// Validate arguments
	file, port := getParameters()
	if file == "" {
		fmt.Println("error: a configuration file is required")
		os.Exit(-1)
	}

	// Read configuration file
	contents, err := os.ReadFile(filepath.Clean(file))
	if err != nil {
		fmt.Println("failed to read configuration file: ", err)
		os.Exit(-1)
	}

	// Decode configuration file
	conf := NewServerConfig()
	if strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml") {
		if err := yaml.Unmarshal(contents, conf); err != nil {
			fmt.Println("failed to decode YAML configuration file: ", err)
			os.Exit(-1)
		}
	}
	if strings.HasSuffix(file, ".json") {
		if err := json.Unmarshal(contents, conf); err != nil {
			fmt.Println("failed to decode JSON configuration file: ", err)
			os.Exit(-1)
		}
	}
	if len(conf.Paths) == 0 {
		fmt.Println("no valid configuration to use")
		os.Exit(-1)
	}

	// Start server
	h := newHandler(conf)
	fmt.Println("serving on port:", port)
	srv := http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           logMiddleware(getServerMux(h)),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		ReadTimeout:       10 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("server error: ", err)
		os.Exit(-1)
	}
}

func getParameters() (string, int) {
	// Define flags
	file := ""
	port := 9090
	ff := flag.String("config", file, "configuration file")
	fp := flag.Int("port", port, "TCP port")
	flag.Parse()

	// Read file from ENV variable and flag
	if ef := os.Getenv("GOVANITY_CONFIG"); ef != "" {
		file = ef
	}
	if *ff != "" {
		file = *ff
	}

	// Read port from ENV variable and flag
	if ep := os.Getenv("GOVANITY_PORT"); ep != "" {
		var err error
		port, err = strconv.Atoi(ep)
		if err != nil {
			port = 9090
		}
	}
	if *fp != port {
		port = *fp
	}
	return file, port
}

func setHeaders(res http.ResponseWriter, ct string, cache string, code int) {
	res.Header().Add("Cache-Control", cache)
	res.Header().Add("Content-Type", ct)
	res.Header().Add("X-Content-Type-Options", "nosniff")
	res.Header().Add("X-Go-Vanity-Server-Build", buildCode)
	res.Header().Add("X-Go-Vanity-Server-Version", coreVersion)
	res.WriteHeader(code)
}

func logMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/ping" {
			addr := r.RemoteAddr
			if xf := r.Header.Get("X-Real-Ip"); xf != "" {
				addr = xf
			}
			log.Printf("%s %s %s [%s]\n", addr, r.Method, r.URL, r.UserAgent())
		}
		handler.ServeHTTP(w, r)
	})
}

func getServerMux(h *handler) *http.ServeMux {
	mux := http.NewServeMux()

	// Ping
	mux.HandleFunc("/api/ping", func(res http.ResponseWriter, req *http.Request) {
		setHeaders(res, "text/plain; charset=utf-8", h.cache(), http.StatusOK)
		_, _ = res.Write([]byte("pong"))
	})

	// Version
	mux.HandleFunc("/api/version", func(res http.ResponseWriter, req *http.Request) {
		js, _ := json.MarshalIndent(versionInfo(), "", "  ")
		setHeaders(res, "application/json", h.cache(), http.StatusOK)
		_, _ = res.Write(js)
	})

	// Configuration
	mux.HandleFunc("/api/conf", func(res http.ResponseWriter, req *http.Request) {
		js, _ := json.MarshalIndent(h.conf, "", "  ")
		setHeaders(res, "application/json", h.cache(), http.StatusOK)
		_, _ = res.Write(js)
	})

	// Main index
	mux.HandleFunc("/index.html", func(res http.ResponseWriter, req *http.Request) {
		index, err := h.getIndex()
		if err != nil {
			setHeaders(res, "text/plain; charset=utf-8", h.cache(), http.StatusInternalServerError)
			_, _ = res.Write([]byte(err.Error()))
			return
		}
		setHeaders(res, "text/html; charset=utf-8", h.cache(), http.StatusOK)
		_, _ = res.Write(index)
	})

	// Catch-all path
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		repo, err := h.getRepo(strings.TrimSuffix(req.URL.Path, "/"))
		if err != nil {
			setHeaders(res, "text/plain; charset=utf-8", h.cache(), http.StatusNotFound)
			_, _ = res.Write([]byte(err.Error()))
			return
		}
		setHeaders(res, "text/html; charset=utf-8", h.cache(), http.StatusOK)
		_, _ = res.Write(repo)
	})

	return mux
}
