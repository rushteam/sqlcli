package suggest

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/vitessio/vitess/go/vt/sqlparser"
	// "github.com/pingcap/parser"
	// _ "github.com/pingcap/tidb/types/parser_driver"
)

const (
	//TipContected ..
	TipContected = "You are now connected to database \"%s\" as user \"%s\"\n"
)

//StateCommon ..
type StateCommon struct {
	Prompt  *prompt.Prompt
	Engine  *DbEngine
	Suggest *MysqlSuggest
}

func (s *StateCommon) commonSuggests() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "source", Description: "Execute commands from file."},
		{Text: "quit", Description: "Quit."},
	}
}

//Run ..
func (s *StateCommon) Run() {
	s.Prompt = prompt.New(
		s.Executor,
		s.Complete,
		prompt.OptionPrefix(fmt.Sprintf("%s >", s.Engine.String())),
		prompt.OptionInputTextColor(prompt.Yellow),
	)
	db, err := s.Engine.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	s.Suggest = NewMysqlSuggest(db)
	s.Prompt.Run()
	//fmt.Printf(TipContected, s.Engine.DbName, s.Engine.User)
	return
}

//Executor ..
func (s *StateCommon) Executor(t string) {
	// fmt.Println("执行命令:" + t)
	cmds := make(map[string]ExecutorHandler, 0)
	cmds["quit"] = QuitHandler
	cmds["exit"] = QuitHandler

	if cmd, ok := cmds[t]; ok {
		cmd(t)
		return
	}
	//执行sql
	//1.获取连接
	// db := s.Engine.Connect()
	//2.连接不可用重连
	//执行sql
	fmt.Println(t)
	return
}

//Complete ..
func (s *StateCommon) Complete(d prompt.Document) []prompt.Suggest {
	var suggests []prompt.Suggest
	var found = false
	var err error
	if d.TextBeforeCursor() == "" {
		return suggests
	}
	// _ := strings.ToLower(d.CurrentLine())
	args := strings.Split(d.TextBeforeCursor(), " ")
	// If PIPE is in text before the cursor, returns empty suggestions.
	for i := range args {
		if args[i] == "|" {
			return suggests
		}
	}
	w := d.GetWordBeforeCursor()
	// If word before the cursor starts with "\", returns system cli
	if strings.HasPrefix(w, "\\") {
		// return optionCompleter(args, strings.HasPrefix(w, "--"))
	}
	if suggests, found, err = s.CompleteUseDatabase(d); found {
		if err != nil {
			fmt.Println(err)
		}
		return suggests
	}
	if suggests, found = s.CompleteSQL(d); found {
		return suggests
	}
	// if len(prefixSug) == 1 {
	// 	// fmt.Println(d.CurrentLine())
	// }

	suggests, err = s.Suggest.Keywords()
	if err != nil {
		fmt.Println(err)
	}
	suggests = append(suggests, s.commonSuggests()...)

	return prompt.FilterHasPrefix(
		suggests,
		w,
		true,
	)
}

//CompleteUseDatabase change db
func (s *StateCommon) CompleteUseDatabase(d prompt.Document) ([]prompt.Suggest, bool, error) {
	// fmt.Println("")
	fmt.Println("TextBeforeCursor", d.TextBeforeCursor())
	fmt.Println("TextAfterCursor", d.TextAfterCursor())
	fmt.Println("GetWordAfterCursor", d.GetWordAfterCursor())
	fmt.Println("GetWordAfterCursorUntilSeparator", d.GetWordAfterCursorUntilSeparator(""))
	fmt.Println("GetWordAfterCursorUntilSeparatorIgnoreNextToCursor", d.GetWordAfterCursorUntilSeparatorIgnoreNextToCursor(" "))
	fmt.Println("GetWordAfterCursorWithSpace", d.GetWordAfterCursorWithSpace())
	fmt.Println("GetWordBeforeCursor", d.GetWordBeforeCursor())
	fmt.Println("CurrentLineBeforeCursor", d.CurrentLineBeforeCursor())
	fmt.Println("FindEndOfCurrentWord", d.FindEndOfCurrentWord())
	fmt.Println("FindEndOfCurrentWordUntilSeparator", d.FindEndOfCurrentWordUntilSeparator(" "))
	fmt.Println("FindEndOfCurrentWordUntilSeparatorIgnoreNextToCursor", d.FindEndOfCurrentWordUntilSeparatorIgnoreNextToCursor(" "))
	fmt.Println("FindStartOfPreviousWord", d.FindStartOfPreviousWord())
	fmt.Println("------")
	var suggests []prompt.Suggest
	//GetWordBeforeCursor
	//CurrentLineBeforeCursor
	if strings.ToLower(d.TextBeforeCursor()) == "use " {
		//提示db
		databases, err := s.Suggest.Databases()
		if err != nil {
			return suggests, false, err
		}
		suggests = prompt.FilterHasPrefix(
			databases,
			d.GetWordBeforeCursor(),
			true,
		)
		return suggests, true, nil
	}
	return suggests, false, nil
}

//CompleteSQL ..
func (s *StateCommon) CompleteSQL(d prompt.Document) ([]prompt.Suggest, bool) {
	// line := d.CurrentLine()
	sql := d.CurrentLineBeforeCursor()
	word := d.GetWordBeforeCursor()
	if len(word) > 0 {
		if strings.HasSuffix(word, "(") || strings.HasPrefix(word, "\\") {
			sql = sql[:len(sql)-len(word)]
		}
	}
	fmt.Println("sql:", sql)
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		fmt.Println("stmt:", err)
		return []prompt.Suggest{}, false
	}
	tokens := sqlparser.NewStringTokenizer(sql)
	fmt.Printf("stmt = %+v tokens = %+v\n", stmt, tokens)
	// switch stmt := stmt.(type) {
	// case *sqlparser.Select:
	// 	fmt.Println("select:", stmt)
	// case *sqlparser.Insert:
	// 	fmt.Println(stmt)
	// case *sqlparser.Update:
	// 	fmt.Println(stmt)
	// case *sqlparser.Delete:
	// 	fmt.Println(stmt)
	// }
	suggests := []prompt.Suggest{
		// {Text: "select * from", Description: ""},
	}
	suggests = prompt.FilterHasPrefix(
		suggests,
		d.TextBeforeCursor(),
		true,
	)
	return suggests, len(suggests) > 0
}

//ExecutorHandler ..
type ExecutorHandler func(t string) error

//QuitHandler ..
func QuitHandler(t string) error {
	fmt.Println("Bye!")
	os.Exit(0)
	return nil
}
