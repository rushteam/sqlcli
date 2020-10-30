package suggest

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

//DbEngine ..
type DbEngine struct {
	Type   string
	DbName string
	User   string
	Pass   string
	Host   string
	Port   int
	Dbx    *sqlx.DB
}

//DSN ..
func (s *DbEngine) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s", s.User, s.Pass, s.Host, s.Port, s.DbName, "?parseTime=true&readTimeout=3s&writeTimeout=3s&timeout=3s")
}

//String ..
func (s *DbEngine) String() string {
	return fmt.Sprintf("%s %s@%s:%d(%s)", s.Type, s.User, s.Host, s.Port, s.DbName)
}

//Connect ..
func (s *DbEngine) Connect() (*sqlx.DB, error) {
	var err error
	if s.Dbx == nil {
		s.Dbx, err = sqlx.Open(s.Type, s.DSN())
	}
	return s.Dbx, err
}
