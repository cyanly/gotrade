package logger

import (
	"net/http"
	"os"
	"time"

	apexlog "github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/multi"
	"github.com/apex/log/handlers/text"
	"github.com/tj/go-elastic"

	es "github.com/cyanly/gotrade/core/logger/handlers"
)

func init() {
	//Default to console logging
	apexlog.SetHandler(cli.Default)

	//Upgrade to ElasticSearch if defined in ENV
	if os.Getenv("ES_ADDR") != "" {
		esClient := elastic.New(os.Getenv("ES_ADDR")) //"http://192.168.99.101:9200"
		esClient.HTTPClient = &http.Client{
			Timeout: 5 * time.Second,
		}

		e := es.New(&es.Config{
			Client:     esClient,
			BufferSize: 100,
		})

		t := text.New(os.Stderr)

		apexlog.SetHandler(multi.New(e, t))
	}
}
