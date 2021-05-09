package dialects

type Dialect string

func (d Dialect) String() string {
	return string(d)
}

const (
	Mysql  Dialect = "mysql"
	Sqlite Dialect = "sqlite"
)
