package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/docopt/docopt-go"
)

var (
	version = `2.0`
	usage   = `eyed ` + version + `

Usage:
	eyed [options]

Options:
	-l <listen>     Listen specified address for new reports [default: :80]
	-s <listen>     Listen specified address for statistics [default: :8000]
	-d <directory>  Use specified directory for placing reports
					[default: /var/eyed/reports/].
`
)

type NewReportsHandler struct {
	directory string
}

type StatisticsHandler struct {
	directory string
}

func main() {
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		panic(err)
	}

	var (
		reportsListenAddress    = args["-l"].(string)
		statisticsListenAddress = args["-s"].(string)
		reportsDirectory        = args["-d"].(string)
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

	go func() {
		handler := &NewReportsHandler{
			directory: reportsDirectory,
		}

		err := http.ListenAndServe(reportsListenAddress, handler)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		handler := &StatisticsHandler{
			directory: reportsDirectory,
		}

		err := http.ListenAndServe(statisticsListenAddress, handler)
		if err != nil {
			log.Fatal(err)
		}
	}()

	select {}
}

func (handler *NewReportsHandler) ServeHTTP(
	response http.ResponseWriter, request *http.Request,
) {
	defer func() {
		scream := recover()
		if scream != nil {
			log.Println(scream)
		}
	}()

	hostname := strings.Trim(request.URL.Path, "/")

	log.Printf("got report for '%s' hostname", hostname)

	if hostname == "" || strings.Contains(hostname, "/") {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	err := writeReport(filepath.Join(handler.directory, hostname))
	if err != nil {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
}

func writeReport(filename string) error {
	file, err := os.OpenFile(
		filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644,
	)
	if err != nil {
		return fmt.Errorf("can't open file: %s", err)
	}

	_, err = file.WriteString(time.Now().String() + "\n")
	if err != nil {
		return fmt.Errorf("can't write file: %s", err)
	}

	err = file.Sync()
	if err != nil {
		return fmt.Errorf("can't sync file: %s", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("can't close file: %s", err)
	}

	return nil
}

func (handler *StatisticsHandler) ServeHTTP(
	response http.ResponseWriter, request *http.Request,
) {
	defer func() {
		scream := recover()
		if scream != nil {
			log.Println(scream)
		}
	}()

	response.Header().Set("Content-Type", "text/plain")

	url := request.URL.Path
	switch {
	case strings.HasPrefix(url, "/f/"):
		handler.handleFilterRequest(response, strings.TrimPrefix(url, "/f/"))

	default:
		response.WriteHeader(http.StatusNotImplemented)
	}
}

func (handler *StatisticsHandler) handleFilterRequest(
	response http.ResponseWriter, daysFilter string,
) {
	statistics, err := getStatistics(handler.directory)
	if err != nil {
		log.Printf(
			"can't get statistics for directory '%s': %s",
			handler.directory, err,
		)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	var days int
	if daysFilter != "" {
		days, err = strconv.Atoi(daysFilter)
		if err != nil {
			log.Printf("can't parse int '%s': %s", daysFilter, err)
			response.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	tabber := new(tabwriter.Writer)
	tabber.Init(response, 1, 4, 1, ' ', 0)

	for hostname, lastModification := range statistics {
		if lastModification >= days {
			_, err := fmt.Fprintf(tabber, "%s\t%d\n", hostname, lastModification)
			if err != nil {
				log.Println(err)
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}

	err = tabber.Flush()
	if err != nil {
		log.Println(err)
	}
}

func getStatistics(directory string) (map[string]int, error) {
	statistics := map[string]int{}

	err := filepath.Walk(
		directory,
		func(path string, fileinfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			path = strings.TrimPrefix(path, directory)
			path = strings.TrimPrefix(path, "/")
			if strings.Contains(path, "/") {
				return filepath.SkipDir
			}

			if path == "" {
				return nil
			}

			duration := time.Now().Sub(fileinfo.ModTime())
			days := int(duration.Hours()) / 24

			statistics[path] = days

			return nil
		},
	)

	return statistics, err
}
