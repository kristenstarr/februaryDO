package main

import (
	"github.com/kristenfelch/pkgindexer/data"
	"github.com/kristenfelch/pkgindexer/input"
	"github.com/kristenfelch/pkgindexer/operation"
	"strings"
	"flag"
	"fmt"
)

// IndexService is responsible for opening a Message Gateway and giving it a Channel
// to return validated messages through.  Messages are then distributed to Remover, Indexer,
// or Querier Services depending on the request, and response is sent back through a
// channel so that Message Gateway can respond to client.
type IndexService interface {
	StartIndexing() (started bool, err error)
}

type SimpleIndexService struct {
	remover operation.Remover
	indexer operation.Indexer
	querier operation.Querier
	lock    data.IndexLock
	gateway input.MessageGateway
}

func (s *SimpleIndexService) StartIndexing() (started bool, err error) {

	c := make(chan *input.ValidatedMessage)
	go s.gateway.Open(c)

	for {
		validMessage := <-c
		s.ProcessMessage(validMessage)
	}

}

func (s *SimpleIndexService) ProcessMessage(input *input.ValidatedMessage) {
	s.lock.Lock()
	respChan := input.ResponseChannel
	var response bool
	var err error

	switch input.Verb {
	case "REMOVE":
		response, err = s.remover.Remove(input.Library)

	case "INDEX":
		var splitDeps []string
		if len(input.Dependencies) > 0 {
			splitDeps = strings.Split(input.Dependencies, ",")
		} else {
			splitDeps = make([]string, 0)
		}
		response, err = s.indexer.Index(input.Library, splitDeps)

	case "QUERY":
		response, err = s.querier.Query(input.Library)
	}

	if err != nil {
		respChan <- "error"
	} else {
		if response {
			respChan <- "ok"
		} else {
			respChan <- "fail"
		}
	}

	s.lock.Unlock()
}

func main() {
	throttle := flag.Int("throttle", 0, "limit on max messages/second from each given")
	flag.Parse()
	if (*throttle > 10000) {
		fmt.Println("Throttle must be less than 10000 requests/second, defaulting to MAX value")
		maxThrottle := 10000
		throttle = &maxThrottle
	}

	var libs = data.New()
	service := &SimpleIndexService{
		operation.NewRemover(libs),
		operation.NewIndexer(libs),
		operation.NewQuerier(libs),
		data.NewLock(),
		input.NewMessageGateway(throttle),
	}
	service.StartIndexing()
}
