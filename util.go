package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/meilisearch/meilisearch-go"
	"log"
	"net/http"
	"strings"
)

func parseHTML(basePath string, html *goquery.Document) []Document {
	var result []Document
	if !strings.HasSuffix(basePath, ".html") {
		return result
	}
	for k, v := range fieldSelectorMap {
		html.Find(v).Each(func(i int, selection *goquery.Selection) {
			id, ok := selection.Attr("id")
			m := map[string]string{
				k: selection.Text(),
			}
			if ok {
				m["URL"] = fmt.Sprintf("%s#%s", basePath, id)
			} else {
				m["URL"] = fmt.Sprintf("%s", basePath)
			}
			doc := ConvertMapToDocument(m)
			if doc.NotEmpty() {
				result = append(result, doc)
			}

		})
	}
	return result
}

func getAllLinksForOnePage(url string) ([]string, []Document) {
	html := transLinkToHTMLDoc(url)
	if html == nil {
		return nil, nil
	}
	var result []string
	html.Find("a").Each(func(i int, selection *goquery.Selection) {
		href, ok := selection.Attr("href")
		if ok && strings.HasPrefix(href, "/") {
			result = append(result, strings.Split(href, "#")[0])
		}
	})
	docs := parseHTML(url, html)
	logger.Info(fmt.Sprintf("Parsed %s, found %d records", url, len(docs)))
	return uniqueStrArray(result), docs
}

func parseFromStartPage(baseUrl string, index *meilisearch.Index) []string {
	links, docs := getAllLinksForOnePage(baseUrl)
	for _, doc := range docs {
		_, err := index.AddDocuments(doc)
		if err != nil {
			log.Printf("AddDocuments failed, %s", err.Error())
		}
	}
	var countedLinks []string
OUT:
	for _, link := range links {
		if strInArray(link, &countedLinks) {
			continue
		}
		countedLinks = append(countedLinks, link)
		tempLinks, tempDocs := getAllLinksForOnePage(fmt.Sprintf("%s%s", baseUrl, link))
		for _, doc := range tempDocs {
			err := doc.Create(index)
			if err != nil {
				log.Printf("AddDocuments failed, %s", err.Error())
			}
		}
		newLinks := uniqueStrArray(arrayDiff(tempLinks, links))
		if len(newLinks) == 0 {
			continue
		}
		links = append(links, newLinks...)
		goto OUT
	}
	return links
}

func uniqueStrArray(arr1 []string) []string {
	m := map[string]bool{}
	for _, v := range arr1 {
		m[v] = true
	}
	var result []string
	for k, _ := range m {
		result = append(result, k)
	}
	return result
}

func transLinkToHTMLDoc(link string) *goquery.Document {
	resp, err := http.Get(link)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil
	}
	return doc
}

func duplicateValue(arr1, arr2 []string) []string {
	m := map[string]bool{}
	for _, v := range arr1 {
		m[v] = true
	}
	var result []string
	for _, v := range arr2 {
		if m[v] {
			result = append(result, v)
		}
	}
	return result
}

func strInArray(str string, arr *[]string) bool {
	for _, v := range *arr {
		if str == v {
			return true
		}
	}
	return false
}

// arrayDiff return values in arr1 but not in arr2
func arrayDiff(arr1, arr2 []string) []string {
	m := map[string]bool{}
	for _, v := range arr2 {
		m[v] = true
	}
	var result []string
	for _, v := range arr1 {
		if !m[v] {
			result = append(result, v)
		}
	}
	return result
}
