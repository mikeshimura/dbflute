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
//	"errors"
	"strconv"
	"strings"
)

type SqlTokenizer struct {
	sql             string
	position        int
	token           string
	tokenType       int
	nextTokenType   int
	bindVariableNum int
}

func (s *SqlTokenizer) Setup(sql string) {
	s.sql = sql
	s.tokenType = TOK_SQL
	s.nextTokenType = TOK_SQL
}

func (s *SqlTokenizer) Next() int {
	//	        if (_position >= _sql.length()) {
	if s.position >= len(s.sql) {

		//            _token = null;
		//            _tokenType = EOF;
		//            _nextTokenType = EOF;
		//            return _tokenType;
		//        }
		s.token = ""
		s.tokenType = TOK_EOF
		s.nextTokenType = TOK_EOF
		return s.tokenType
	}
	//        switch (_nextTokenType) {
	//        case SQL:
	//            parseSql();
	//            break;
	//        case COMMENT:
	//            parseComment();
	//            break;
	//        case ELSE:
	//            parseElse();
	//            break;
	//        case BIND_VARIABLE:
	//            parseBindVariable();
	//            break;
	//        default:
	//            parseEof();
	//            break;
	//        }
	//        return _tokenType;
	switch s.nextTokenType {
	case TOK_SQL:
		s.parseSql()
	case TOK_COMMENT:
		s.parseComment()
	case TOK_ELSE:
		s.parseElse()
	case TOK_BIND_VARIABLE:
		s.parseBindVariable()
	default:
		s.parseEof()

	}
	return s.tokenType
}
func (s *SqlTokenizer) parseSql() {
	//        int commentStartPos = _sql.indexOf("/*", _position);
	commentStartPos := IndexAfter(s.sql, "/*", s.position)
	//        int commentStartPos2 = _sql.indexOf("#*", _position);
	commentStartPos2 := IndexAfter(s.sql, "#*", s.position)
	//        if (0 < commentStartPos2 && commentStartPos2 < commentStartPos) {
	//            commentStartPos = commentStartPos2;
	//        }
	if 0 < commentStartPos2 && commentStartPos2 < commentStartPos {
		commentStartPos = commentStartPos2
	}
	//        int bindVariableStartPos = _sql.indexOf("?", _position);
	bindVariableStartPos := IndexAfter(s.sql, "?", s.position)
	//        int elseCommentStartPos = -1;
	elseCommentStartPos := -1
	//        int elseCommentLength = -1;
	elseCommentLength := -1
	//        int elseCommentSearchCurrentPosition = _position;
	elseCommentSearchCurrentPosition := s.position
	//        while (true) { // searching nearest next ELSE comment
	for {

		//            final int lineCommentStartPos = _sql.indexOf("--", elseCommentSearchCurrentPosition);
		lineCommentStartPos := IndexAfter(s.sql, "--", elseCommentSearchCurrentPosition)
		//            if (lineCommentStartPos < 0) {
		//                break;
		//            }
		if lineCommentStartPos < 0 {
			break
		}
		//            if (calculateNextStartPos(commentStartPos, bindVariableStartPos, -1) < lineCommentStartPos) {
		//                break;
		//            }
		if s.calculateNextStartPos(commentStartPos, bindVariableStartPos, -1) < lineCommentStartPos {
			break
		}
		//            int skipPos = skipWhitespace(lineCommentStartPos + 2);
		skipPos := s.skipWhitespace(lineCommentStartPos + 2)
		//            if (skipPos + 4 < _sql.length() && "ELSE".equals(_sql.substring(skipPos, skipPos + 4))) {
		if skipPos+4 < len(s.sql) && "ELSE" == (s.sql[skipPos:skipPos+4]) {
			//                elseCommentStartPos = lineCommentStartPos;
			elseCommentStartPos = lineCommentStartPos
			//                elseCommentLength = skipPos + 4 - lineCommentStartPos;
			//                break;
			elseCommentLength = skipPos + 4 - lineCommentStartPos
			break
		}
		//            }
		//            elseCommentSearchCurrentPosition = skipPos;
		elseCommentSearchCurrentPosition = skipPos
		//        }
	}
	//        int nextStartPos = calculateNextStartPos(commentStartPos, bindVariableStartPos, elseCommentStartPos);
	nextStartPos := s.calculateNextStartPos(commentStartPos, bindVariableStartPos, elseCommentStartPos)
	if nextStartPos < 0 {
		s.token = s.sql[s.position:]
		s.nextTokenType = TOK_EOF
		s.position = len(s.sql)
		s.tokenType = TOK_SQL
	} else {
		s.token = s.sql[s.position:nextStartPos]
		s.tokenType = TOK_SQL
		needNext := nextStartPos == s.position
		if nextStartPos == commentStartPos {
			s.nextTokenType = TOK_COMMENT
			s.position = commentStartPos + 2
		} else if nextStartPos == elseCommentStartPos {
			s.nextTokenType = TOK_ELSE
			s.position = elseCommentStartPos + elseCommentLength
		} else if nextStartPos == bindVariableStartPos {
			s.nextTokenType = TOK_BIND_VARIABLE
			s.position = bindVariableStartPos
		}
		if needNext {
			s.Next()
		}
	}
	return
}
func (s *SqlTokenizer) parseComment()  {
	commentEndPos := IndexAfter(s.sql, "*/", s.position)
	commentEndPos2 := IndexAfter(s.sql, "*#", s.position)
	if 0 < commentEndPos2 && commentEndPos2 < commentEndPos {
		commentEndPos = commentEndPos2
	}
	if commentEndPos < 0 {
		panic("CommentTerminatorNotFoundException" + s.sql[s.position:])
	}
	s.token = s.sql[s.position:commentEndPos]
	s.nextTokenType = TOK_SQL
	s.position = commentEndPos + 2
	s.tokenType = TOK_COMMENT
	return
}
func (s *SqlTokenizer) parseElse()  {
	s.token = ""
	s.nextTokenType = TOK_SQL
	s.tokenType = TOK_ELSE
	return 
}
func (s *SqlTokenizer) parseBindVariable() {
	s.token = s.nextBindVariableName()
	s.nextTokenType = TOK_SQL
	s.position += 1
	s.tokenType = TOK_BIND_VARIABLE
	return
}
func (s *SqlTokenizer) parseEof(){
	s.token = ""
	s.tokenType = TOK_EOF
	s.nextTokenType = TOK_EOF
	return
}
func (s *SqlTokenizer) calculateNextStartPos(commentStartPos int, bindVariableStartPos int, elseCommentStartPos int) int {
	//	        int nextStartPos = -1;
	nextStartPos := -1
	if commentStartPos >= 0 {
		nextStartPos = commentStartPos
	}
	if bindVariableStartPos >= 0 && (nextStartPos < 0 || bindVariableStartPos < nextStartPos) {
		nextStartPos = bindVariableStartPos
	}
	if elseCommentStartPos >= 0 && (nextStartPos < 0 || elseCommentStartPos < nextStartPos) {
		nextStartPos = elseCommentStartPos
	}
	return nextStartPos
}
func (s *SqlTokenizer) skipWhitespaceNoPos() string {
	index := s.skipWhitespace(s.position)
	s.token = s.sql[s.position:index]
	s.position = index
	return s.token
}
func (s *SqlTokenizer) skipWhitespace(position int) int {
	var whitespace = " \n\r\t"
	index := len(s.sql)
	for i := position; i < len(s.sql); i++ {
		c := string(s.sql[i])
		if strings.Index(whitespace, c) == -1 {
			index = i
			break
		}
	}
	return index
}
func (s *SqlTokenizer) nextBindVariableName() string {
	s.bindVariableNum++
	return "$" + strconv.Itoa(s.bindVariableNum)
}
func (s *SqlTokenizer) getAfter() string {
	return s.sql[s.position:]
}
func (s *SqlTokenizer) getBefore() string {
	return s.sql[0:s.position]
}
func (s *SqlTokenizer) extractDateLiteralPrefix(testValue bool, currentSql string, position int) string {
	if !testValue {
		return ""
	}
	if position >= len(currentSql) {
		return ""
	}
	firstChar := string(currentSql[position])
	if firstChar != "d" && firstChar != "D" && firstChar != "t" && firstChar != "T" {
		return ""
	}
	var rear string
	{
		tmpRear := currentSql[position:]
		maxlength := len("timestamp '")
		if len(tmpRear) > maxlength {
			// get only the quantity needed for performance
			rear = tmpRear[0:maxlength]
		} else {
			rear = tmpRear
		}
	}
	lowerRear := strings.ToLower(rear)
	literalPrefix := ""
	if strings.Index(lowerRear, "date '") == 0 {
		literalPrefix = rear[0:len("date ")]
	} else if strings.Index(lowerRear, "date'") == 0 {
		literalPrefix = rear[0:len("date")]
	} else if strings.Index(lowerRear, "timestamp '") == 0 { // has max length
		literalPrefix = rear[0:len("timestamp ")]
	} else if strings.Index(lowerRear, "timestamp'") == 0 {
		literalPrefix = rear[0:len("timestamp")]
	}
	return literalPrefix
}
func (s *SqlTokenizer) skipToken(testValue bool) string {
	index := len(s.sql) // last index as default

	dateLiteralPrefix := s.extractDateLiteralPrefix(testValue, s.sql, s.position)
	if dateLiteralPrefix != "" {
		s.position = s.position + len(dateLiteralPrefix)
	}

	var quote uint8
	{
		var firstChar uint8
		if s.position < len(s.sql) {
			firstChar = s.sql[s.position]
		} else {
			//Original '\0' but golang can't accept
			firstChar = 0
		}
		// firstChar := (_position < _sql.length() ? _sql.charAt(_position) : '\0');
		//quote = (firstChar == '(' ? ')' : firstChar);
		if firstChar == '(' {
			quote = ')'
		} else {
			quote = firstChar
		}
	}
	quoting := quote == '\'' || quote == ')'
	var spos int
	if quoting {
		spos = s.position + 1
	} else {
		spos = s.position
	}
	// for (int i = quoting ? _position + 1 : _position; i < _sql.length(); ++i) {
	for i := spos; i < len(s.sql); i++ {
		c := s.sql[i]
		if s.isNotQuoteEndPoint(quoting, c) {
			index = i
			break
		} else if s.isBlockCommentBeginPoint(s.sql, c, i) {
			index = i
			break
		} else if s.isLineCommentBeginPoint(s.sql, c, i) {
			index = i
			break
		} else if quoting && s.isSingleQuoteEndPoint(s.sql, quote, c, i) {
			index = i + 1
			break
		} else if quoting && s.isQuoteEndPoint(s.sql, quote, c, i) {
			index = i + 1
			break
		}
	}
	s.token = s.sql[s.position:index]
	if dateLiteralPrefix != "" {
		s.token = dateLiteralPrefix + s.token
	}
	s.tokenType = TOK_SQL
	s.nextTokenType = TOK_SQL
	s.position = index
	return s.token
}
func (s *SqlTokenizer) isNotQuoteEndPoint(quoting bool, c uint8) bool {
	return !quoting && (c == ' ' || c == ',' || c == ')' || c == '(' || c == '\t' || c == '\r' || c == '\n')
}
func (s *SqlTokenizer) isBlockCommentBeginPoint(currentSql string, c uint8, i int) bool {
	return c == '/' && s.isNextCharacter(currentSql, i, '*')
}
func (s *SqlTokenizer) isNextCharacter(currentSql string, i int, targetChar uint8) bool {
	return i+1 < len(currentSql) && currentSql[i+1] == targetChar
}
func (s *SqlTokenizer) isLineCommentBeginPoint(currentSql string, c uint8, i int) bool {
	return c == '-' && s.isNextCharacter(currentSql, i, '-')
}
func (s *SqlTokenizer) isSingleQuoteEndPoint(currentSql string, quote uint8, c uint8, i int) bool {
	sqlLen := len(currentSql)
	var endSqlOrNotEscapeQuote bool = (i+1 >= sqlLen || currentSql[i+1] != '\'')
	return quote == '\'' && c == '\'' && endSqlOrNotEscapeQuote
}
func (s *SqlTokenizer) isQuoteEndPoint(currentSql string, quote uint8, c uint8, i int) bool {
	return c == quote
}
