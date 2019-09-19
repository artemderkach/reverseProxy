package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// Service stores the name and url to be redirected to
type Service struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

var services []Service

func main() {
	port := flag.String("port", "8080", "proxy port to listen on")
	file := flag.String("file", "", "file with hosts data")
	flag.Parse()

	if *file == "" {
		log.Fatal("Expected file with hosts as paramet")
	}
	f, err := os.Open(*file)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error opening file with hosts"))
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error reading data from file"))
	}

	err = json.Unmarshal(data, &services)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error unmarshaling data"))
	}

	reverseProxy := &ReverseProxy{
		Director: director,
	}
	fmt.Printf("starting proxy on localhost:%s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, reverseProxy))
}

// findServiceByName searches for specific name in all of services
func findServiceByName(services []Service, name string) (Service, error) {
	for _, service := range services {
		if service.Name == name {
			return service, nil
		}
	}
	return Service{}, errors.New("service not found in available services")
}

func director(r *http.Request) error {
	// get the lowest domain from host
	serviceName := strings.Split(r.Host, ".")[0]
	service, err := findServiceByName(services, serviceName)
	if err != nil {
		return errors.Wrapf(err, "error finding %s in service list", serviceName)
	}

	r.Host = service.URL
	return nil
}
