package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bryk-io/go-vanity/config"
	"gopkg.in/yaml.v2"
)

func main() {
	// Validate arguments
	file, port := getParameters()
	if file == "" {
		fmt.Println("error: a configuration file is required")
		os.Exit(-1)
	}

	// Read configuration file
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("failed to read configuration file: ", err)
		os.Exit(-1)
	}

	// Decode configuration file
	conf := config.New()
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

	// Prepare server mux
	h := newHandler(conf)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/version", func(res http.ResponseWriter, req *http.Request) {
		js, _ := json.MarshalIndent(versionInfo(), "", "  ")
		res.Header().Add("Content-Type", "application/json")
		res.Header().Add("Cache-Control", h.cache())
		res.WriteHeader(http.StatusOK)
		_, _ = res.Write(js)
	})
	mux.HandleFunc("/api/conf", func(res http.ResponseWriter, req *http.Request) {
		js, _ := json.MarshalIndent(conf, "", "  ")
		res.Header().Add("Content-Type", "application/json")
		res.Header().Add("Cache-Control", h.cache())
		res.WriteHeader(http.StatusOK)
		_, _ = res.Write(js)
	})
	mux.HandleFunc("/index.html", func(res http.ResponseWriter, req *http.Request) {
		index, err := h.getIndex()
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			_, _ = res.Write([]byte(err.Error()))
			return
		}
		res.Header().Add("Cache-Control", h.cache())
		res.WriteHeader(http.StatusOK)
		_, _ = res.Write(index)
	})
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		repo, err := h.getRepo(strings.TrimSuffix(req.RequestURI, "/"))
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			_, _ = res.Write([]byte(err.Error()))
			return
		}
		res.Header().Add("Cache-Control", h.cache())
		res.WriteHeader(http.StatusOK)
		_, _ = res.Write(repo)
	})

	// Start server
	fmt.Println("serving on port:", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
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
