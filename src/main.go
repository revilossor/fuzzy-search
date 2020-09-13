package main

import (
	"encoding/json"
	_ "expvar"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/blevesearch/bleve"
	bleveHttp "github.com/blevesearch/bleve/http"
)

var batchSize = flag.Int("batchSize", 100, "batch size for indexing")
var bindAddr = flag.String("addr", ":8094", "http listen address")
var jsonDir = flag.String("jsonDir", "data/", "json directory")
var indexPath = flag.String("index", "qna-search.bleve", "index path")
var standardIndexPath = flag.String("standardIndex", "qna-search-standard.bleve", "standard index path")
var staticEtag = flag.String("staticEtag", "", "A static etag value.")
var staticPath = flag.String("static", "static/", "Path to the static content")

func main() {

	flag.Parse()

	// open the index
	qnaIndex, err := bleve.Open(*indexPath)
	qnaStandardIndex, err := bleve.Open(*standardIndexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		log.Printf("Creating new index...")
		// create a mapping
		indexMapping, err := buildIndexMapping()
		if err != nil {
			log.Fatal(err)
		}
		qnaIndex, err = bleve.New(*indexPath, indexMapping)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Creating new standard index...")
		// create a mapping just for keywords, ie just 'job' in 'what is your job'
		standardIndexMapping, err := buildStandardIndexMapping()
		if err != nil {
			log.Fatal(err)
		}
		qnaStandardIndex, err = bleve.New(*standardIndexPath, standardIndexMapping)
		if err != nil {
			log.Fatal(err)
		}

		// index data in the background
		go func() {
			err = indexQnA(qnaIndex)
			if err != nil {
				log.Fatal(err)
			}
			err = indexQnA(qnaStandardIndex)
			if err != nil {
				log.Fatal(err)
			}
		}()
	} else if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Opening existing index...")
	}

	everything := NewIndexAlias(qnaIndex, qnaStandardIndex)

	// create a router to serve static files
	router := staticFileRouter()

	// add the API
	bleveHttp.RegisterIndexName("qna", everything)
	searchHandler := bleveHttp.NewSearchHandler("qna")
	router.Handle("/api/search", searchHandler).Methods("POST")
	listFieldsHandler := bleveHttp.NewListFieldsHandler("qna")
	router.Handle("/api/fields", listFieldsHandler).Methods("GET")

	debugHandler := bleveHttp.NewDebugDocumentHandler("qna")
	debugHandler.DocIDLookup = docIDLookup
	router.Handle("/api/debug/{docID}", debugHandler).Methods("GET")

	// handler to retrieve documents
	docGetHandler := bleveHttp.NewDocGetHandler("qna")
	docGetHandler.IndexNameLookup = indexNameLookup
	docGetHandler.DocIDLookup = docIDLookup
	router.Handle("/api/{indexName}/{docID}", docGetHandler).Methods("GET")

	// start the HTTP server
	http.Handle("/", router)
	log.Printf("Listening on %v", *bindAddr)
	log.Fatal(http.ListenAndServe(*bindAddr, nil))
}

func indexQnA(i bleve.Index) error {

	// open the directory
	dirEntries, err := ioutil.ReadDir(*jsonDir)
	if err != nil {
		return err
	}

	// walk the directory entries for indexing
	log.Printf("Indexing...")
	count := 0
	startTime := time.Now()
	batch := i.NewBatch()
	batchCount := 0
	for _, dirEntry := range dirEntries {
		filename := dirEntry.Name()
		// read the bytes
		jsonBytes, err := ioutil.ReadFile(*jsonDir + "/" + filename)
		if err != nil {
			return err
		}
		// parse bytes as json
		var jsonDoc interface{}
		err = json.Unmarshal(jsonBytes, &jsonDoc)
		if err != nil {
			return err
		}
		ext := filepath.Ext(filename)
		docID := filename[:(len(filename) - len(ext))]
		batch.Index(docID, jsonDoc)
		batchCount++

		if batchCount >= *batchSize {
			err = i.Batch(batch)
			if err != nil {
				return err
			}
			batch = i.NewBatch()
			batchCount = 0
		}
		count++
		if count%1000 == 0 {
			indexDuration := time.Since(startTime)
			indexDurationSeconds := float64(indexDuration) / float64(time.Second)
			timePerDoc := float64(indexDuration) / float64(count)
			log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
		}
	}
	// flush the last batch
	if batchCount > 0 {
		err = i.Batch(batch)
		if err != nil {
			log.Fatal(err)
		}
	}
	indexDuration := time.Since(startTime)
	indexDurationSeconds := float64(indexDuration) / float64(time.Second)
	timePerDoc := float64(indexDuration) / float64(count)
	log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
	return nil
}
