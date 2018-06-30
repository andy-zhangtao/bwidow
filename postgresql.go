/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package bwidow

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/andy-zhangtao/gogather/zReflect"
	_ "github.com/lib/pq"
	"github.com/pelletier/go-toml"
)

type BWPostgresConf struct {
	Endpoint string
	DB       string
	User     string
	Password string
}

type BWPostgresql struct {
	db       *sql.DB
	tableMap map[string]string
}

var pqconf *BWPostgresConf

func (this *BWPostgresql) setDB() {

}
func (this *BWPostgresql) Check() error {
	if os.Getenv(BW_PQ_ENDPOINT) == "" {
		if _, err := os.Stat("bwidow_pq.toml"); os.IsNotExist(err) {
			//	toml 配置文件不存在,检查bwidow_mongo.json文件
			if _, err := os.Stat("bwidow_pq.json"); os.IsNotExist(err) {
				return errors.New(fmt.Sprintf("Can not find Mongo configure. Env[%s]/Toml/Json are all lost!", BW_PQ_ENDPOINT))
			}
		}

		return nil
	}

	if os.Getenv(BW_PQ_DB) == "" {
		return errors.New(fmt.Sprintf("[%s] Empty!", BW_PQ_DB))
	}

	return nil
}

func (this *BWPostgresql) DriverInit() error {
	if err := this.Check(); err != nil {
		return err
	}

	pqconf = new(BWPostgresConf)
	pqconf.User = os.Getenv(BW_PQ_USER)
	pqconf.Password = os.Getenv(BW_PQ_PASSWD)
	pqconf.Endpoint = os.Getenv(BW_PQ_ENDPOINT)
	pqconf.DB = os.Getenv(BW_PQ_DB)

	if pqconf.Endpoint == "" {
		useToml := true
		//	读取toml
		data, err := ioutil.ReadFile("bwidow_pq.toml")
		if err != nil || len(data) == 0 {
			useToml = false
			//尝试json
			data, err = ioutil.ReadFile("bwidow_pq.json")
			if err != nil || len(data) == 0 {
				return errors.New(fmt.Sprintf("Parse Configure Error %s", err.Error()))
			}
		}

		if useToml {
			err = toml.Unmarshal(data, pqconf)
			if err != nil {
				return err
			}
		} else {
			err = json.Unmarshal(data, pqconf)
			if err != nil {
				return err
			}
		}
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", pqconf.User, pqconf.Password, pqconf.Endpoint, pqconf.DB)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	if db.Ping() != nil {
		return errors.New(fmt.Sprintf("Ping [%s] Error [%s]", connStr, err))
	}

	this.db = db
	fmt.Printf("Postgresql [%s] connect success! \n", pqconf.Endpoint)

	return nil
}

func (this *BWPostgresql) Map(u interface{}, name string) {
	if this.tableMap == nil {
		this.tableMap = make(map[string]string)
	}

	this.tableMap[reflect.TypeOf(u).Name()] = name
	return
}
func (this *BWPostgresql) checkIndex(u interface{}) error {
	m := zReflect.ReflectStructInfoWithTag(u, true, "bw")

	var key []string
	for k, _ := range m {
		key = append(key, k)
	}

	table := this.tableMap[getTypeName(u)]

	sql := fmt.Sprintf("CREATE UNIQUE INDEX IF NOT EXISTS %s_unique on %s ( %s )", table, table, strings.Join(key, ","))

	_, err := this.db.Exec(sql)

	return err
}
func (this *BWPostgresql) first(uPtr interface{}) error {
	return nil
}
func (this *BWPostgresql) findOne(uPtr interface{}) error {
	return nil
}
func (this *BWPostgresql) findAll(uPtr interface{}, aPtr interface{}) error {
	return nil
}
func (this *BWPostgresql) findAllWithSort(uPtr interface{}, aArray interface{}, sortField []string) error {
	return nil
}
func (this *BWPostgresql) save(u interface{}) error {
	return nil
}
func (this *BWPostgresql) saveAll(uArray []interface{}) error {
	return nil
}
func (this *BWPostgresql) update(uPtr interface{}, field []string) (int, error) {
	return 0, nil
}
func (this *BWPostgresql) delete(uPtr interface{}, field []string) (int, error) {
	return 0, nil
}
func (this *BWPostgresql) deleteAll(uPtr interface{}) (int, error) {
	return 0, nil
}
func (this *BWPostgresql) count(uPtr interface{}) (int, error) {
	return 0, nil
}
