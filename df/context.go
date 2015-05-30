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
	"bytes"
	"github.com/mikeshimura/dbflute/log"
	"strconv"
	//"reflect"
	//	"fmt"
)

type CommandContext interface {
	addSql(sql string)
	addSqlSingle(Sql string, val interface{}, ctype string)
	addSqlBind(sql string, bindVariables *List, bindVariableTypes *StringList)
	getSql() string
	GetArgs() map[string]interface{}
	GetArgTypes() map[string]string
	getBindVariables() *List
	getBindVariableTypes() *StringList
	setEnabled(value bool)
	isEnabled() bool
	isBeginChild() bool
	getArg(name string) interface{}
	isAlreadySkippedConnector() bool
	setAlreadySkippedConnector(value bool)
}

type CommandContextCreator struct {
	argNames []string
	argTypes []string
}

func (c *CommandContextCreator) createCommandContext(args []interface{}) *CommandContext {
	//create CommandContextImpl as root 相当
	ctx := new(CommandContextImpl)
	ctx.args = make(map[string]interface{})
	ctx.argTypes = make(map[string]string)
	if args != nil {
		for i := 0; i < len(args); i++ {
			argType := ""
			if args[i] != nil {
				if i < len(c.argTypes) {
					argType = c.argTypes[i]
				} else if args[i] != nil {
					argType = GetType(args[i])
				}
			}
			if i < len(c.argNames) {
				ctx.addArg(c.argNames[i], args[i], argType)
			} else {
				ctx.addArg("$"+strconv.Itoa(i+1), args[i], argType)
			}
		}
	}
	var cc CommandContext = ctx
	return &cc
}

type CommandContextImpl struct {
	sqlSb                   *bytes.Buffer
	args                    map[string]interface{}
	argTypes                map[string]string
	bindVariables           *List
	bindVariableTypes       *StringList
	parent                  *CommandContext
	enabled                 bool
	beginChild              bool
	alreadySkippedConnector bool
}

func (c *CommandContextImpl) setAlreadySkippedConnector(value bool) {
	c.alreadySkippedConnector = value
}

func (c *CommandContextImpl) isAlreadySkippedConnector() bool {
	return c.alreadySkippedConnector
}
func (c *CommandContextImpl) isBeginChild() bool {
	return c.beginChild
}
func (c *CommandContextImpl) setEnabled(value bool) {
	c.enabled = value
}
func (c *CommandContextImpl) isEnabled() bool {
	return c.enabled
}
func (c *CommandContextImpl) GetArgs() map[string]interface{} {
	return c.args
}
func (c *CommandContextImpl) GetArgTypes() map[string]string {
	return c.argTypes
}
func (c *CommandContextImpl) getSql() string {
	return c.sqlSb.String()
}
func (c *CommandContextImpl) addSql(sql string) {
	if c.sqlSb == nil {
		c.sqlSb = new(bytes.Buffer)
	}
	c.sqlSb.WriteString(sql)
	//log.InternalDebug("sql :"+c.sqlSb.String())
}

func (c *CommandContextImpl) setupBindValiables() {
	if c.bindVariables == nil {
		c.bindVariables = new(List)
	}
	if c.bindVariableTypes == nil {
		c.bindVariableTypes = new(StringList)
	}
}
func (c *CommandContextImpl) addSqlSingle(sql string, val interface{}, ctype string) {
	c.addSql(sql)
	c.setupBindValiables()
	c.bindVariables.Add(val)
	c.bindVariableTypes.Add(ctype)
}
func (c *CommandContextImpl) addSqlBind(sql string, bindVariables *List, bindVariableTypes *StringList) {
	c.addSql(sql)
	c.setupBindValiables()
	if bindVariables != nil {
		for _, v := range bindVariables.data {
			c.bindVariables.Add(v)
		}
	}
	if bindVariableTypes != nil {
		for _, v := range bindVariableTypes.data {
			c.bindVariableTypes.Add(v)
		}
	}
}
func (c *CommandContextImpl) getBindVariables() *List {
	return c.bindVariables
}
func (c *CommandContextImpl) getBindVariableTypes() *StringList {
	return c.bindVariableTypes
}
func (c *CommandContextImpl) addArg(argName string, arg interface{}, argType string) {
	c.args[argName] = arg
	c.argTypes[argName] = argType
	log.InternalDebug("CommandContext argName:" + argName + " argType:" + argType)
}
func (c *CommandContextImpl) getArg(name string) interface{} {
	var res interface{}
	res = c.args[name]
	if res != nil {
		return res
	} else if c.parent != nil {
		return (*c.parent).getArg(name)
	} else {
		if len(c.args) == 1 {
			for k := range c.args {
				return c.args[k]
			}
		}
		return nil
	}
}

type Context struct {
	cmap map[string]interface{}
}

func CreateContext() *Context {
	ctx := new(Context)
	ctx.cmap = make(map[string]interface{})
	return ctx
}
func (p *Context) Put(key string, value interface{}) {
	p.cmap[key] = value
}
func (p *Context) Get(key string) interface{} {
	return p.cmap[key]
}
