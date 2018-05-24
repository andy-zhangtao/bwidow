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

func (this *BWMongo) first(u interface{}) (err error) {
	this.setDB()
	defer this.db.Session.Close()

	return this.db.C(this.tableMap[getTypeName(u)]).Find(nil).One(u)
}

func (this *BWMongo) findOne(u interface{}) (err error) {
	this.setDB()
	defer this.db.Session.Close()

	m := zReflect.ReflectStructInfo(u)

	return this.db.C(this.tableMap[getTypeName(u)]).Find(bson.M(m)).One(u)
}

func (this *BWMongo) findAll(u interface{}, a interface{}) (err error) {
	this.setDB()
	defer this.db.Session.Close()

	m := zReflect.ReflectStructInfo(u)

	return this.db.C(this.tableMap[getTypeName(u)]).Find(bson.M(m)).All(a)
}

func (this *BWMongo) findAllWithSort(u interface{}, a interface{}, sortField []string) (err error) {
	this.setDB()
	defer this.db.Session.Close()

	m := zReflect.ReflectStructInfo(u)

	return this.db.C(this.tableMap[getTypeName(u)]).Find(bson.M(m)).Sort(sortField...).All(a)
}

func (this *BWMongo) save(u interface{}) (err error) {
	this.setDB()
	defer this.db.Session.Close()
	return this.db.C(this.tableMap[getTypeName(u)]).Insert(u)
}

func (this *BWMongo) saveAll(u []interface{}) (err error) {
	typeName := reflect.TypeOf(u[0]).Name()
	this.setDB()
	defer this.db.Session.Close()

	bulk := this.db.C(this.tableMap[typeName]).Bulk()
	bulk.Insert(u...)
	_, err = bulk.Run()
	return
}

func (this *BWMongo) update(uPtr interface{}, field []string) (num int, err error) {
	this.setDB()
	defer this.db.Session.Close()

	m := zReflect.ReflectStructInfo(uPtr)

	nm := make(map[string]interface{})

	for _, f := range field {
		nm[f] = m[f]
	}

	info, err := this.db.C(this.tableMap[getTypeName(uPtr)]).UpdateAll(bson.M(nm), bson.M{"$set": bson.M(m)})
	return info.Updated, err
}

func (this *BWMongo) delete(uPtr interface{}, field []string) (num int, err error) {
	this.setDB()
	defer this.db.Session.Close()

	m := zReflect.ReflectStructInfo(uPtr)

	nm := make(map[string]interface{})

	for _, f := range field {
		nm[f] = m[f]
	}

	info, err := this.db.C(this.tableMap[getTypeName(uPtr)]).RemoveAll(bson.M(nm))

	return info.Removed, err
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

func (this *BWMongo) checkIndex(uPtr interface{}) (err error) {
	this.setDB()
	defer this.db.Session.Close()

	m := zReflect.ReflectStructInfoWithTag(uPtr,true, "bw")

	fmt.Println(m)
	var key []string
	for k, _ := range m {
		key = append(key, k)
	}

	return this.db.C(this.tableMap[getTypeName(uPtr)]).EnsureIndexKey(key...)
}
