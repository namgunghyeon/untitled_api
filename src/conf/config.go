package conf

import (
  "fmt"
  "github.com/BurntSushi/toml"
)

func LoadIgnoreDirs() {}

func LoadCouchbase() Config{
  var conf Config
  if _, err := toml.DecodeFile("src/conf/couchbase.toml", &conf); err != nil {
    fmt.Println("error", err)
    return conf
  }
  return conf
}

func LoadCockroach() Config{
  var conf Config
  if _, err := toml.DecodeFile("src/conf/cockroach.toml", &conf); err != nil {
    fmt.Println("error", err)
    return conf
  }
  return conf
}
