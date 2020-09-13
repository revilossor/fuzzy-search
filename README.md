install go: ```brew install golang```  

run: ```go run src/*.go```  

build: ```go build -o bin/*.go src/main.go ```  
run binary ```./bin/main```  

I've pulled in some data that q-bert generates, and theres a script to convert it to json documents that bleve can index.

```
node create-documents
```

When you start the go app, the generated documents will be indexed and an api stood up on 8094 with the following routes

 - ```GET /api/fields``` get all the fields in the documents
 - ```POST /api/search``` perform a [search](https://blevesearch.com/docs/Query/)
 - ```GET /api/qna/<document_id>``` get a document by its id

you can send requests from the postman collection

There are 2 indexes, which are combined into a single queryable index

One index has a simple text analyser that analyses all parts of the question
so "what is your name" is indexed as ["what", "is", "your", "name"]

Another index has the standard bleave analyser, but with english stop word removal
so "what is your name" is indexed as ["name"]

This approach helps guard against false positives for queries against a set of data with lots of stop word similarity, eg "what is your name" when we have indexed several "what is your [X]" items

##### Ref

https://gobyexample.com
https://blevesearch.com
https://github.com/blevesearch/beer-search
https://godoc.org/github.com/blevesearch/bleve
https://www.youtube.com/watch?v=CEfaIlzki5U
http://bleveanalysis.couchbase.com/analysis

## TODO

- dockerise
- frontend - web to speech?
- test suite checks accuracy.... bit like qbert...
- n-grams increase accuracy? of non-exact queries?
