package main

import (
	"github.com/kristenfelch/pkgindexer/data"
	"github.com/kristenfelch/pkgindexer/input"
	"github.com/kristenfelch/pkgindexer/operation"
	"strings"
	"flag"
	"github.com/kristenfelch/pkgindexer/logging"
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
		response, err = s.remover.Remove(input.Package)

	case "INDEX":
		var splitDeps []string
		if len(input.Dependencies) > 0 {
			splitDeps = strings.Split(input.Dependencies, ",")
		} else {
			splitDeps = make([]string, 0)
		}
		response, err = s.indexer.Index(input.Package, splitDeps)

	case "QUERY":
		response, err = s.querier.Query(input.Package)
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

// Main method reads input parameters throttle/logLevel, and starts up our service.
func main() {
	throttle := flag.Int("throttle", 0, "limit on max messages/second from each given")
	logLevel := flag.String("logLevel", "INFO", "log level")
	flag.Parse()
	logger := logging.NewIndexLogger(logLevel)

	if (*throttle > 10000) {
		logger.Debug("Throttle must be less than 10000 requests/second, defaulting to MAX value")
		maxThrottle := 10000
		throttle = &maxThrottle
	}

	var store = data.NewIndexStore(logger)
	service := &SimpleIndexService{
		operation.NewRemover(store, logger),
		operation.NewIndexer(store, logger),
		operation.NewQuerier(store),
		data.NewLock(),
		input.NewMessageGateway(throttle, logger),
	}
	logger.Info("Indexing service starting on port 8080...")

	service.StartIndexing()
}
