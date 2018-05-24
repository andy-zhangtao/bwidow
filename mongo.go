/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package bwidow

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/andy-zhangtao/gogather/zReflect"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/23.

type BWMongo struct {
	session  *mgo.Session
	client   *mgo.Session
	db       *mgo.Database
	tableMap map[string]string
}

func (this *BWMongo) Map(u interface{}, name string) {
	if this.tableMap == nil {
		this.tableMap = make(map[string]string)
	}

	this.tableMap[reflect.TypeOf(u).Name()] = name
}

func (this *BWMongo) setDB() {
	this.db = this.session.Clone().DB(os.Getenv(BW_MONGO_DB))
}

func (this *BWMongo) First(u interface{}) (err error) {
	this.setDB()
	defer this.db.Session.Close()

	return this.db.C(this.tableMap[getTypeName(u)]).Find(nil).One(u)
}

func (this *BWMongo) FindOne(u interface{}) (err error) {
	this.setDB()
	defer this.db.Session.Close()

	m := zReflect.ReflectStructInfo(u)

	return this.db.C(this.tableMap[getTypeName(u)]).Find(bson.M(m)).One(u)
}

func (this *BWMongo) FindAll(u interface{}, a interface{}) (err error) {
	this.setDB()
	defer this.db.Session.Close()

	m := zReflect.ReflectStructInfo(u)

	return this.db.C(this.tableMap[getTypeName(u)]).Find(bson.M(m)).All(a)
}

func (this *BWMongo) FindAllWithSort(u interface{}, a interface{}, sortField []string) (err error) {
	this.setDB()
	defer this.db.Session.Close()

	m := zReflect.ReflectStructInfo(u)

	return this.db.C(this.tableMap[getTypeName(u)]).Find(bson.M(m)).Sort(sortField...).All(a)
}

func (this *BWMongo) DriverInit() (err error) {
	if err = this.Check(); err != nil {
		return
	}

	user := os.Getenv(BW_MONGO_USER)
	passwd := os.Getenv(BW_MONGO_PASSWD)

	var session *mgo.Session
	if user == "" {
		session, err = mgo.Dial(os.Getenv(BW_MONGO_ENDPOINT))
		if err != nil {
			return
		}
	} else {
		dialInfo := &mgo.DialInfo{
			Addrs:    []string{os.Getenv(BW_MONGO_ENDPOINT)},
			Database: os.Getenv(BW_MONGO_DB),
			Username: user,
			Password: passwd,
			Timeout:  10 * time.Second,
		}

		session, err = mgo.DialWithInfo(dialInfo)
		if err != nil {
			return
		}
	}

	_, err = session.BuildInfo()
	if err != nil {
		return
	}

	this.session = session
	return
}

func (this *BWMongo) Check() (err error) {

	if os.Getenv(BW_MONGO_ENDPOINT) == "" {
		return errors.New(fmt.Sprintf("[%s] Empty!", BW_MONGO_ENDPOINT))
	}

	if os.Getenv(BW_MONGO_DB) == "" {
		return errors.New(fmt.Sprintf("[%s] Empty!", BW_MONGO_DB))
	}
	return
}
