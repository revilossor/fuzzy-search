package main

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/analysis/analyzer/simple"
	"github.com/blevesearch/bleve/mapping"
)

func buildIndexMapping() (mapping.IndexMapping, error) {

	keywordMapping := bleve.NewTextFieldMapping()
	keywordMapping.Analyzer = keyword.Name

	simpleMapping := bleve.NewTextFieldMapping()
	simpleMapping.Analyzer = simple.Name

	qnaMapping := bleve.NewDocumentMapping()
	qnaMapping.AddFieldMappingsAt("type", keywordMapping)		// TODO dont save type?
	qnaMapping.AddFieldMappingsAt("question", simpleMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("qna", qnaMapping)

	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = simple.Name

	return indexMapping, nil
}

func buildStandardIndexMapping() (mapping.IndexMapping, error) {

	keywordMapping := bleve.NewTextFieldMapping()
	keywordMapping.Analyzer = keyword.Name

	englishMapping := bleve.NewTextFieldMapping()
	englishMapping.Analyzer = en.AnalyzerName

	qnaMapping := bleve.NewDocumentMapping()
	qnaMapping.AddFieldMappingsAt("type", keywordMapping)
	qnaMapping.AddFieldMappingsAt("question", englishMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("qna-standard", qnaMapping)

	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}
