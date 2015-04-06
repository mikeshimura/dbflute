package df

import (
//	"fmt"
	"io/ioutil"
	//"reflect"
	"strings"
	"errors"
)

type DBMetaProvider struct {
	TableDbNameFlexibleMap       *StringKeyMap
	TablePropertyNameFlexibleMap *StringKeyMap
	TableDbNameInstanceMap       map[string]*DBMeta
}

func CreateDBMetaProvider() *DBMetaProvider {
	dm := new(DBMetaProvider)
	dm.TableDbNameFlexibleMap = CreateAsFlexible()
	dm.TablePropertyNameFlexibleMap = CreateAsFlexible()
	dm.TableDbNameInstanceMap = make(map[string]*DBMeta)
	return dm
}

type ResourceContext struct {
	ConditionBeanContext *ConditionBeanContext
	OutsideSqlContext    *OutsideSqlContext
	//SqlAnalyzerFactory *SqlAnalyzerFactory
}

func (b *ResourceContext) GetConditionBeanContext() *ConditionBeanContext {
	if b.ConditionBeanContext == nil {
		b.ConditionBeanContext = new(ConditionBeanContext)
	}
	return b.ConditionBeanContext
}
func (b *ResourceContext) GetOutsideSqlContext() *OutsideSqlContext {
	if b.OutsideSqlContext == nil {
		b.OutsideSqlContext = new(OutsideSqlContext)
	}
	return b.OutsideSqlContext
}

func (b *ResourceContext) CreateSqlAnalyzer(twoWaySql string, blockeNullParameter bool) *SqlAnalyzer {
	analyzer := new(SqlAnalyzer)
	analyzer.Setup(twoWaySql, blockeNullParameter)
	return analyzer
}

type ConditionBeanContext struct {
	ConditionBean interface{}
}
type OutsideSqlContext struct {
	Pmb            interface{}
	OutsideSqlPath string
}

func (o *OutsideSqlContext) readFilteredOutsideSql(suffix string) (string,error) {
	sql,err1 := o.readPlainOutsideSql(suffix)
	if err1!=nil{
		return "",err1
	}	
	rsql,err := o.replaceOutsideSqlBindCharacterOnLineComment(sql)
	if err!=nil{
		return "",err
	}
	//        if (_outsideSqlFilter != null) {
	//            sql = _outsideSqlFilter.filterReading(sql);
	//        }
	//        return sql;
	return rsql,nil
}
func (o *OutsideSqlContext) replaceOutsideSqlBindCharacterOnLineComment(sql string) (string,error) {
	//fmt.Println(sql)
	bindCharacter := "?"
	if strings.Index(sql, bindCharacter) < 0 {
		return sql,nil
	}
	//        if (sql.indexOf(bindCharacter) < 0) {
	//            return sql;
	//        }
	//        final String lineSeparator = "\n";
	//        if (sql.indexOf(lineSeparator) < 0) {
	//            return sql;
	//        }
	//        final String lineCommentMark = "--";
	//        if (sql.indexOf(lineCommentMark) < 0) {
	//            return sql;
	//        }
	//        final StringBuilder sb = new StringBuilder();
	//        final String[] lines = sql.split(lineSeparator);
	//        for (String line : lines) {
	//            final int lineCommentIndex = line.indexOf("--");
	//            if (lineCommentIndex < 0) {
	//                sb.append(line).append(lineSeparator);
	//                continue;
	//            }
	//            final String lineComment = line.substring(lineCommentIndex);
	//            if (lineComment.contains("ELSE") || !lineComment.contains(bindCharacter)) {
	//                sb.append(line).append(lineSeparator);
	//                continue;
	//            }
	//
	//            if (_log.isDebugEnabled()) {
	//                _log.debug("...Replacing bind character on line comment: " + lineComment);
	//            }
	//            final String filteredLineComment = replaceString(lineComment, bindCharacter, "Q");
	//            sb.append(line.substring(0, lineCommentIndex)).append(filteredLineComment).append(lineSeparator);
	//        }
	//        return sb.toString();
	return "",errors.New("replaceOutsideSqlBindCharacterOnLineComment")
}
func (o *OutsideSqlContext) readPlainOutsideSql(suffix string) (string,error) {
	standardPath := o.OutsideSqlPath
	readSql, err := ioutil.ReadFile(standardPath)
	var sql string = string(readSql)
	if err != nil {
		return "",errors.New("Can't read sql file:" + standardPath)
	}
	if sql == "" {
		return "",errors.New("Sql file has no content:" + standardPath)
	}
	//        String readSql = doReadPlainOutsideSql(sqlFileEncoding, dbmsSuffix, standardPath);
	//        if (readSql != null) {
	//            return readSql;
	//        }
	//        // means not found
	//        final String pureName = Srl.substringLastRear(standardPath, "/");
	//        if (pureName.contains("Bhv_")) { // retry for ApplicationBehavior
	//            final String dir = Srl.substringLastFront(standardPath, "/");
	//            final String filtered = Srl.replace(pureName, "Bhv_", "BhvAp_");
	//            final String bhvApPath = dir + "/" + filtered;
	//            readSql = doReadPlainOutsideSql(sqlFileEncoding, dbmsSuffix, bhvApPath);
	//        }
	//        if (readSql != null) {
	//            return readSql;
	//        }
	//        throwOutsideSqlNotFoundException(standardPath);
	//        return null; // unreachable
	return sql,nil
}

func (o *OutsideSqlContext) generateSpecifiedOutsideSqlUniqueKey(methodName string, path string, pmb interface{}, option *OutsideSqlOption, resultType string) string {

	pmbKey := GetType(pmb)
	resultKey := resultType
	tableDbName := option.TableDbName
	generatedUniqueKey := option.GenerateUniqueKey()
	return tableDbName + ":" + methodName + "():" + path + ":" + pmbKey + ":" + generatedUniqueKey + ":" + resultKey
}
