package df

import (
	"fmt"
	"reflect"
)

type DBWay interface {
	BuildSequenceNextValSql(s string) string
	GetIdentitySelectSql() string
	IsBlockCommentSupported() bool
	IsLineCommentSupported() bool
	IsScrollableCursorSupported() bool
	GetOriginalWildCardList() []string
	IsUniqueConstraintException(sqlState string, errorCode int64) bool
	Connect(v interface{}) string
}

//DBWay interfaceを実装
type WayOfPostgreSQL struct {
}

func (w *WayOfPostgreSQL) BuildSequenceNextValSql(sequenceName string) string {
	return "select nextval ('" + sequenceName + "')"
}
func (w *WayOfPostgreSQL) GetIdentitySelectSql() string {
	return ""
}
func (w *WayOfPostgreSQL) IsBlockCommentSupported() bool {
	return true
}
func (w *WayOfPostgreSQL) IsLineCommentSupported() bool {
	return true
}
func (w *WayOfPostgreSQL) IsScrollableCursorSupported() bool {
	return true
}
func (w *WayOfPostgreSQL) GetOriginalWildCardList() []string {
	return []string{}
}
func (w *WayOfPostgreSQL) IsUniqueConstraintException(sqlState string, errorCode int64) bool {
	return sqlState == "23505"
}
func (w *WayOfPostgreSQL) Connect(v interface{}) string {
	buf := ""
	s := reflect.ValueOf(v)
	for i := 0; i < s.Len(); i++ {
		if len(buf) > 0 {
			buf += " || "
		}
		buf += fmt.Sprint(s.Index(i).Interface())
	}
	return buf
}
