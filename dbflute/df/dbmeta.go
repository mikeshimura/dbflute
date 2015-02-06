package df

import (

)

type DBMeta interface {
	GetProjectName() string
	GetCurrentDBDef() DBDef
}