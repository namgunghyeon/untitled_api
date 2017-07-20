package db

import (
  "fmt"
  "gopkg.in/couchbase/gocb.v1"
)

type Result struct {
  Default map[string]Project `json:"default"`
}

type Project struct {
	Project string `json:"project"`
	Version string `json:"version"`
	Type string `json:"type"`
  Path string `json:path`
  Name string `json:name`
  Count int `json:count`
}


func Search(project string, version string, searchType string, name string) {
  cluster, _ := gocb.Connect("couchbase://192.168.56.213")
  bucket, _ := cluster.OpenBucket("default", "")
  where := `where project = "` + project + `" and version = "` + version + `" and type = "`+ searchType + `" and Name like "`+ name + `%"`
  fmt.Println(where)
  query := gocb.NewN1qlQuery("select * from default " + where)
  rows, _ := bucket.ExecuteN1qlQuery(query, nil)
  var row interface{}
  //var projects []Project
  //items := []string{}

  for rows.Next(&row) {
    //items = append(items, row["Name"])

    //projects = append(projects,row)
    fmt.Printf("Row: %s", row.(map[string]interface{})["default"])
  }
}
