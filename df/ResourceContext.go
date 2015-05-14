/*
 * Copyright 2014-2015 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language
 * governing permissions and limitations under the License.
 */
package df

import (
//	"fmt"
	"io/ioutil"
	//"reflect"
	"strings"
//	"errors"
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
	SqlClause interface{}
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

func (o *OutsideSqlContext) readFilteredOutsideSql(suffix string) string {
	sql := o.readPlainOutsideSql(suffix)
	rsql := o.replaceOutsideSqlBindCharacterOnLineComment(sql)
	//        if (_outsideSqlFilter != null) {
	//            sql = _outsideSqlFilter.filterReading(sql);
	//        }
	//        return sql;
	return rsql
}
func (o *OutsideSqlContext) replaceOutsideSqlBindCharacterOnLineComment(sql string) string {
	//fmt.Println(sql)
	bindCharacter := "?"
	if strings.Index(sql, bindCharacter) < 0 {
		return sql
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
	panic("replaceOutsideSqlBindCharacterOnLineComment")
}
func (o *OutsideSqlContext) readPlainOutsideSql(suffix string) string {
	standardPath := o.OutsideSqlPath
	readSql, err := ioutil.ReadFile(standardPath)
	var sql string = string(readSql)
	if err != nil {
		panic("Can't read sql file:" + standardPath)
	}
	if sql == "" {
		panic("Sql file has no content:" + standardPath)
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
	return sql
}

func (o *OutsideSqlContext) generateSpecifiedOutsideSqlUniqueKey(methodName string, path string, pmb interface{}, option *OutsideSqlOption, resultType string) string {

	pmbKey := GetType(pmb)
	resultKey := resultType
	tableDbName := option.TableDbName
	generatedUniqueKey := option.GenerateUniqueKey()
	return tableDbName + ":" + methodName + "():" + path + ":" + pmbKey + ":" + generatedUniqueKey + ":" + resultKey
}
