package conf

type Config struct {
  Couchbase CouchbaseConfig `toml:"couchbase_staging"`
  Cockroach CockroachConfig `toml:"cockroach_staging"`
}

type CouchbaseConfig struct {
	Host string `toml:"host"`
	Project string `toml:"project"`
	ProjectInfo string `toml:"projectInfo"`
}

type CockroachConfig struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
	Db string `toml:"db"`
  SSLMode string `toml:"sslMode"`
}
