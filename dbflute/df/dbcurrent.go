package df

import ()

type DBCurrent struct {
	ProjectName                     string
	ProjectPrefix                   string
	DBDef                           *DBDef
	DBWay                           *DBWay
	PagingCountLater                bool
	PagingCountLeastJoin            bool
	InnerJoinAutoDetect             bool
	ThatsBadTimingDetect            bool
	NullOrEmptyQueryAllowed         bool
	EmptyStringQueryAllowed         bool
	EmptyStringParameterAllowed     bool
	OverridingQueryAllowed          bool
	NonSpecifiedColumnAccessAllowed bool
	ColumnNullObjectAllowed         bool
	ColumnNullObjectGearedToSpecify bool
	DisableSelectIndex              bool
	QueryUpdateCountPreCheck        bool
}
