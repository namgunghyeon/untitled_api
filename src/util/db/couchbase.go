package db

import (
  "strconv"
  "gopkg.in/couchbase/gocb.v1"
  "model"
  "conf"
  "fmt"
)


func Search(project string, version string, searchType string, name string, limit int) []model.Project{
  config := conf.LoadCouchbase()
  cluster, _ := gocb.Connect("couchbase://" + config.Couchbase.Host)
  bucket, _ := cluster.OpenBucket(config.Couchbase.Project, "")
  where := `where project = "` + project + `" and version = "` + version + `" and type = "`+ searchType + `" and name like "`+ name + `%" limit ` + strconv.Itoa(limit)
  query := gocb.NewN1qlQuery("select project, version, type, name, `path` from project_data " + where)
  rows, _ := bucket.ExecuteN1qlQuery(query, nil)

  var row model.Project
  var projects []model.Project
  for rows.Next(&row) {
    projects = append(projects, row)
  }
  return projects
}

func KeywordIndex(name string, offset int, limit int) []model.KeywordIndex{
  config := conf.LoadCouchbase()
  cluster, _ := gocb.Connect("couchbase://" + config.Couchbase.Host)
  bucket, _ := cluster.OpenBucket("keyword_index", "")
  where := `WHERE keyword LIKE "`+ name + `%" LIMIT ` + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset)
  query := gocb.NewN1qlQuery("SELECT keyword FROM keyword_index " + where)
  rows, _ := bucket.ExecuteN1qlQuery(query, nil)
  fmt.Println("query", query)
  var row model.KeywordIndex
  var KeywordIndexes []model.KeywordIndex
  for rows.Next(&row) {
    KeywordIndexes = append(KeywordIndexes, row)
  }
  return KeywordIndexes
}

func Keyword(name string, limit int) []model.Keyword{
  config := conf.LoadCouchbase()
  cluster, _ := gocb.Connect("couchbase://" + config.Couchbase.Host)
  bucket, _ := cluster.OpenBucket("keyword", "")
  where := `WHERE keyword_index = "`+ name + `" `
  query := gocb.NewN1qlQuery("SELECT project, version, type, keyword_index, `path` FROM keyword " + where)
  rows, _ := bucket.ExecuteN1qlQuery(query, nil)
  fmt.Println("query", query)
  var row model.Keyword
  var keywords []model.Keyword
  for rows.Next(&row) {
    keywords = append(keywords, row)
  }
  return keywords
}
