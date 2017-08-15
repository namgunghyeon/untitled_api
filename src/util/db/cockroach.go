package db

import (
  "strconv"
  _ "github.com/lib/pq"
  "database/sql"
  "model"
  "log"
  "fmt"
)

func CockroachKeywordIndex(name string, offset int, limit int) []model.KeywordIndex{
  db, err := sql.Open("postgres", "postgresql://root@104.156.238.187:26257/untitled?sslmode=disable")
  if err != nil {
      log.Fatal("error connecting to the database: ", err)
  }
  query := "SELECT keyword FROM keyword_index WHERE keyword LIKE '" + name + "%' LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset)
  rows, err := db.Query(query)
  fmt.Println("query", query)
  if err != nil {
      log.Fatal(err)
  }
  defer rows.Close()
  var KeywordIndexes []model.KeywordIndex
  for rows.Next() {
    var keyword string
    if err := rows.Scan(&keyword); err != nil {
        log.Fatal(err)
    }
    KeywordIndexes = append(KeywordIndexes, model.KeywordIndex{ Keyword: keyword })
  }
  return KeywordIndexes
}

func CockroachKeyword(name string, limit int) []model.Keyword{
  db, err := sql.Open("postgres", "postgresql://root@104.156.238.187:26257/untitled?sslmode=disable")
  if err != nil {
      log.Fatal("error connecting to the database: ", err)
  }

  query := "SELECT project, version, type, keyword_index, path FROM keyword where keyword_index = '" + name + "'";
  rows, err := db.Query(query)
  fmt.Println("query", query)
  if err != nil {
      log.Fatal(err)
  }

  var keywords []model.Keyword
  for rows.Next() {
    var row model.Keyword
    if err := rows.Scan(&row.Project, &row.Version, &row.Type, &row.KeywordIndex, &row.Path); err != nil {
        log.Fatal(err)
    }
    keywords = append(keywords, row)
  }
  return keywords
}
