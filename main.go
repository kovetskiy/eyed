package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docopt/docopt-go"
)

const (
	usage = `eyed 1.0

Usage:
	eyed [options]

Options:
	-l <listen>     Listen specified address [default: :80]
	-d <directory>  Use specified directory for placing reports
					[default: /var/eyed/reports/].
`
)

type ReportsHandler struct {
	directory string
}

func main() {
	args, err := docopt.Parse(usage, nil, true, "1.0", false)
	if err != nil {
		panic(err)
	}

	var (
		listenAddress    = args["-l"].(string)
		reportsDirectory = args["-d"].(string)
	)

	_, err = os.Stat(reportsDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(reportsDirectory, 0755)
		}

		if err != nil {
			log.Fatal(err)
		}
	}

	err = listen(listenAddress, reportsDirectory)
	if err != nil {
		log.Fatal(err)
	}
}

func listen(address, reportsDirectory string) error {
	handler := &ReportsHandler{
		directory: reportsDirectory,
	}

	return http.ListenAndServe(address, handler)
}

func (handler *ReportsHandler) ServeHTTP(
	response http.ResponseWriter, request *http.Request,
) {
	defer func() {
		message := recover()
		if message != nil {
			log.Println(message)
		}
	}()

	var err error
	defer func() {
		if err != nil {
			log.Println(err)
			http.Error(response, err.Error(), http.StatusInternalServerError)
		}
	}()

	hostname := strings.Trim(request.URL.Path, "/")

	if hostname == "" || strings.Contains(hostname, "/") {
		http.Error(response, "bad hostname", http.StatusBadRequest)
		return
	}

	log.Printf("got report for '%s' hostname", hostname)

	file, err := os.OpenFile(
		filepath.Join(handler.directory, hostname),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return
	}

	_, err = file.WriteString(time.Now().String() + "\n")
	if err != nil {
		return
	}

	err = file.Close()
	if err != nil {
		return
	}
}
