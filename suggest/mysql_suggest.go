package suggest

import (
	"github.com/c-bata/go-prompt"
	"github.com/jmoiron/sqlx"
)

//NewMysqlSuggest ..
func NewMysqlSuggest(DB *sqlx.DB) *MysqlSuggest {
	return &MysqlSuggest{
		DB: DB,
	}
}

//MysqlSuggest ..
type MysqlSuggest struct {
	*sqlx.DB
	databases []prompt.Suggest // 获取db列表
	keywords  []prompt.Suggest //获取关键字列表
	functions []prompt.Suggest
}

//Keywords ..
func (s *MysqlSuggest) Keywords() ([]prompt.Suggest, error) {
	if len(s.keywords) == 0 {
		var keys = []string{"ACCESS", "ADD", "ALL", "ALTER TABLE", "AND", "ANY", "AS",
			"ASC", "AUTO_INCREMENT", "BEFORE", "BEGIN", "BETWEEN",
			"BIGINT", "BINARY", "BY", "CASE", "CHANGE MASTER TO", "CHAR",
			"CHARACTER SET", "CHECK", "COLLATE", "COLUMN", "COMMENT",
			"COMMIT", "CONSTRAINT", "CREATE", "CURRENT",
			"CURRENT_TIMESTAMP", "DATABASE", "DATE", "DECIMAL", "DEFAULT",
			"DELETE FROM", "DESC", "DESCRIBE", "DROP",
			"ELSE", "END", "ENGINE", "ESCAPE", "EXISTS", "FILE", "FLOAT",
			"FOR", "FOREIGN KEY", "FORMAT", "FROM", "FULL", "FUNCTION",
			"GRANT", "GROUP BY", "HAVING", "HOST", "IDENTIFIED", "IN",
			"INCREMENT", "INDEX", "INSERT INTO", "INT", "INTEGER",
			"INTERVAL", "INTO", "IS", "JOIN", "KEY", "LEFT", "LEVEL",
			"LIKE", "LIMIT", "LOCK", "LOGS", "LONG", "MASTER",
			"MEDIUMINT", "MODE", "MODIFY", "NOT", "NULL", "NUMBER",
			"OFFSET", "ON", "OPTION", "OR", "ORDER BY", "OUTER", "OWNER",
			"PASSWORD", "PORT", "PRIMARY", "PRIVILEGES", "PROCESSLIST",
			"PURGE", "REFERENCES", "REGEXP", "RENAME", "REPAIR", "RESET",
			"REVOKE", "RIGHT", "ROLLBACK", "ROW", "ROWS", "ROW_FORMAT",
			"SAVEPOINT", "SELECT", "SESSION", "SET", "SHARE", "SHOW",
			"SLAVE", "SMALLINT", "SMALLINT", "START", "STOP", "TABLE",
			"THEN", "TINYINT", "TO", "TRANSACTION", "TRIGGER", "TRUNCATE",
			"UNION", "UNIQUE", "UNSIGNED", "UPDATE", "USE", "USER",
			"USING", "VALUES", "VARCHAR", "VIEW", "WHEN", "WHERE", "WITH"}
		s.keywords = make([]prompt.Suggest, len(keys))
		for i, key := range keys {
			s.keywords[i] = prompt.Suggest{
				Text:        key,
				Description: "keyword: " + key,
			}
		}
	}
	return s.keywords, nil
}

//Functions ..
func (s *MysqlSuggest) Functions() ([]prompt.Suggest, error) {
	if len(s.functions) == 0 {
		var keys = []string{"AVG", "CONCAT", "COUNT", "DISTINCT", "FIRST", "FORMAT",
			"FROM_UNIXTIME", "LAST", "LCASE", "LEN", "MAX", "MID",
			"MIN", "NOW", "ROUND", "SUM", "TOP", "UCASE", "UNIX_TIMESTAMP"}
		s.functions = make([]prompt.Suggest, len(keys))
		for i, key := range keys {
			s.functions[i] = prompt.Suggest{
				Text:        key,
				Description: "Function: " + key,
			}
		}
	}
	return s.functions, nil
}

//Databases ..
func (s *MysqlSuggest) Databases() ([]prompt.Suggest, error) {
	var err error
	if len(s.databases) == 0 {
		//DbInfo ..
		type DbInfo struct {
			Schema  string `db:"TABLE_SCHEMA"`
			Table   string `db:"TABLE_NAME"`
			Comment string `db:"TABLE_COMMENT"`
		}
		_, err = s.Exec("use information_schema;")
		if err != nil {
			return s.databases, err
		}
		rows, err := s.Queryx("select TABLE_SCHEMA,TABLE_NAME,TABLE_COMMENT from tables")
		if err != nil {
			return s.databases, err
		}
		defer rows.Close()
		var list []*DbInfo
		for rows.Next() {
			var item = &DbInfo{}
			err = rows.StructScan(item)
			if err != nil {
				break
			}
			list = append(list, item)
		}
		var uniq = make(map[string]*DbInfo, 0)
		for _, item := range list {
			uniq[item.Schema] = item
		}
		for _, v := range uniq {
			s.databases = append(s.databases, prompt.Suggest{
				Text:        v.Schema,
				Description: "database: " + v.Schema,
			})
		}
	}
	return s.databases, err
}
