package df

import (
	"strconv"
)

type SqlClauseMySql struct {
	BaseSqlClause
	fetchScopeSqlSuffix string
	lockSqlSuffix string
}

func (p *SqlClauseMySql) DoFetchPage() {
	p.fetchScopeSqlSuffix = " limit " + strconv.Itoa(p.fetchSize) + " offset " + strconv.Itoa(p.getPageStartIndex())
}
func (p *SqlClauseMySql) DoClearFetchPageClause() {
	p.fetchScopeSqlSuffix = ""
}
func (p *SqlClauseMySql) DoFetchFirst() {
	p.DoFetchPage()
}
func (p *SqlClauseMySql) CreateFromHint() string {
	return ""
}
func (p *SqlClauseMySql) CreateSqlSuffix() string {
	return p.fetchScopeSqlSuffix + p.lockSqlSuffix
}
func (p *SqlClauseMySql) CreateSelectHint() string {
	return ""
}

type SqlClauseSqlServer struct {
	BaseSqlClause
	fetchFirstSelectHint string
}

func (p *SqlClauseSqlServer) DoFetchPage() {
	p.fetchFirstSelectHint = " top " + strconv.Itoa(p.fetchSize)
}
func (p *SqlClauseSqlServer) DoClearFetchPageClause() {
	p.fetchFirstSelectHint = ""
}
func (p *SqlClauseSqlServer) DoFetchFirst() {
	p.fetchFirstSelectHint = " top " + strconv.Itoa(p.fetchSize)
}
func (p *SqlClauseSqlServer) CreateFromHint() string {
	return ""
}
func (p *SqlClauseSqlServer) CreateSqlSuffix() string {
	return ""
}
func (p *SqlClauseSqlServer) CreateSelectHint() string {
	return p.fetchFirstSelectHint
}

type SqlClausePostgres struct {
	BaseSqlClause
	fetchScopeSqlSuffix string
	lockSqlSuffix string
}

func (p *SqlClausePostgres) DoFetchPage() {
	p.fetchScopeSqlSuffix = " limit " + strconv.Itoa(p.fetchSize) + " offset " + strconv.Itoa(p.getPageStartIndex())
}
func (p *SqlClausePostgres) DoClearFetchPageClause() {
	p.fetchScopeSqlSuffix = ""
}
func (p *SqlClausePostgres) DoFetchFirst() {
	p.DoFetchPage()
}
func (p *SqlClausePostgres) CreateFromHint() string {
	return ""
}
func (p *SqlClausePostgres) CreateSqlSuffix() string {
	return p.fetchScopeSqlSuffix + p.lockSqlSuffix
}
func (p *SqlClausePostgres) CreateSelectHint() string {
	return ""
}