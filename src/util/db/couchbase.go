package db

import (
  "strconv"
  "gopkg.in/couchbase/gocb.v1"
  "model"
  "conf"
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
