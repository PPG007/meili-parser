package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/meilisearch/meilisearch-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	fieldSelectorMap = map[string]string{
		"titleLevel0": ".theme-default-content h1",
		"titleLevel1": ".theme-default-content h2",
		"titleLevel2": ".theme-default-content h3",
		"titleLevel3": ".theme-default-content h4",
		"titleLevel4": ".theme-default-content h5",
		"content":     ".theme-default-content p, .theme-default-content li",
	}
)

type Document struct {
	Id          string `json:"id"`
	TitleLevel0 string `json:"titleLevel0"`
	TitleLevel1 string `json:"titleLevel1"`
	TitleLevel2 string `json:"titleLevel2"`
	TitleLevel3 string `json:"titleLevel3"`
	TitleLevel4 string `json:"titleLevel4"`
	Content     string `json:"content"`
	URL         string `json:"URL"`
}

func (d *Document) Create(index *meilisearch.Index) error {
	d.Id = primitive.NewObjectID().Hex()
	_, err := index.AddDocuments(d)
	return err
}

func (d *Document) NotEmpty() bool {
	if d.TitleLevel0 != "" {
		return true
	}
	if d.TitleLevel1 != "" {
		return true
	}
	if d.TitleLevel2 != "" {
		return true
	}
	if d.TitleLevel3 != "" {
		return true
	}
	if d.TitleLevel4 != "" {
		return true
	}
	if d.Content != "" {
		return true
	}
	return false
}

func ConvertMapToDocument(m map[string]string) Document {
	doc := Document{}
	bytes, _ := json.Marshal(m)
	json.Unmarshal(bytes, &doc)
	return doc
}

func (d *Document) ConvertToMap() map[string]interface{} {
	bytes, err := json.Marshal(d)
	if err != nil {
		return nil
	}
	result := map[string]interface{}{}
	json.Unmarshal(bytes, &result)
	return result
}

func (d *Document) GetHashValue() string {
	source := fmt.Sprintf("%s%s%s%s%s%s%s", d.TitleLevel0, d.TitleLevel1, d.TitleLevel2, d.TitleLevel3, d.TitleLevel4, d.Content, d.URL)
	return getMD5(source)
}

func getMD5(source string) string {
	h := md5.New()
	h.Write([]byte(source))
	return hex.EncodeToString(h.Sum(nil))
}
