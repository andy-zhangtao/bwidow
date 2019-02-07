/*
 * Copyright (c) 2019.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package bwidow

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
	Name    string `json:"name" pq:"name" bw:"name"`
	Age     int    `json:"age" pq:"age"`
	Address string `json:"address" pq:"address"`
}

const DB = "u1"

var bw *BW

func testInit() {
	os.Setenv(BW_PQ_ENDPOINT, "127.0.0.1:5432")
	os.Setenv(BW_PQ_DB, "postgres")
	os.Setenv(BW_PQ_USER, "postgres")
	os.Setenv(BW_PQ_PASSWD, "123456")

	if bw == nil {
		bw = GetWidow()
		bw.Driver(DRIVER_PQ)
	}
}

func TestBWPostgresql_Check(t *testing.T) {
	os.Setenv(BW_PQ_ENDPOINT, "127.0.0.1:5432")
	os.Setenv(BW_PQ_DB, "test")

	pg := new(BWPostgresql)

	err := pg.Check()
	assert.Nil(t, err)
}

func TestBWPostgresql_DriverInit(t *testing.T) {
	os.Setenv(BW_PQ_ENDPOINT, "127.0.0.1:5432")
	os.Setenv(BW_PQ_DB, "postgres")
	os.Setenv(BW_PQ_USER, "postgres")
	os.Setenv(BW_PQ_PASSWD, "123456")

	bw := GetWidow()
	bw.Driver(DRIVER_PQ)
	err := bw.Error()
	assert.Nil(t, err)
}

func TestBWPostgresql_Map(t *testing.T) {
	testInit()

	bw = bw.Map(User{}, DB)
	err := bw.Error()
	assert.Nil(t, err)
}

func TestBWPostgresql_CheckIndex(t *testing.T) {
	testInit()

	bw.CheckIndex(new(User))
}

func TestBWPostgresql_Save(t *testing.T) {

	u := User{
		Name:    "abc",
		Address: "unknown",
		Age:     12,
	}

	testInit()
	bw = bw.Map(User{}, DB)
	err := bw.Error()
	assert.Nil(t, err)

	err = bw.Save(u)
	assert.Nil(t, err)
}

func TestBWPostgresql_SaveAll(t *testing.T) {
	var us []User

	for i := 0; i <= 5; i++ {
		us = append(us, User{
			Name:    fmt.Sprintf("user-%d", i),
			Age:     i + 10,
			Address: fmt.Sprintf("add-%d", i),
		})
	}

	testInit()
	bw = bw.Map(User{}, DB)
	err := bw.SaveAll(us)
	assert.Nil(t, err)

}

func TestBWPostgresql_Update(t *testing.T) {
	u := User{
		Name:    "user-2",
		Age:     22,
		Address: "add-22",
	}

	testInit()
	bw = bw.Map(User{}, DB)

	_, err := bw.Update(&u, []string{"name"})
	assert.Nil(t, err)
}

func TestBWPostgresql_Count(t *testing.T) {
	testInit()
	bw = bw.Map(User{}, DB)

	number, err := bw.Count(new(User))
	assert.Nil(t, err)
	assert.Equal(t, 7, number)
}

func TestBWPostgresql_First(t *testing.T) {
	testInit()
	bw = bw.Map(User{}, DB)

	u := new(User)
	err := bw.First(u)
	assert.Nil(t, err)
	assert.Equal(t, "abc", u.Name)
	assert.Equal(t, 12, u.Age)
	assert.Equal(t, "unknown", u.Address)
}

func TestBWPostgresql_FindOne(t *testing.T) {
	testInit()
	bw = bw.Map(User{}, DB)

	u := User{
		Name: "user-0",
	}

	err := bw.FindOne(&u)
	assert.Nil(t, err)

	assert.Equal(t, "user-0", u.Name)
	assert.Equal(t, 10, u.Age)
	assert.Equal(t, "add-0", u.Address)
}

func TestBWPostgresql_FindAll(t *testing.T) {
	testInit()
	bw = bw.Map(User{}, DB)

	var us []*User

	u := User{}
	err := bw.FindAll(&u, &us)
	assert.Nil(t, err)
	assert.Equal(t, 7, len(us))

	var uss []*User
	u.Age = 10
	err = bw.FindAll(&u, &uss)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(uss))

}

func TestBWPostgresql_FindAllWithSort(t *testing.T) {
	testInit()
	bw = bw.Map(User{}, DB)

	var us []*User

	u := User{}
	err := bw.FindAllWithSort(&u, &us, []string{"age"})
	assert.Nil(t, err)
	assert.Equal(t, 7, len(us))

	assert.Equal(t, "user-0", us[0].Name)
	assert.Equal(t, 10, us[0].Age)
	assert.Equal(t, "add-0", us[0].Address)
}

func TestBWPostgresql_Delete(t *testing.T) {
	testInit()
	bw = bw.Map(User{}, DB)

	u := User{
		Age: 10,
	}
	_, err := bw.Delete(&u, []string{"age"})
	assert.Nil(t, err)
}

func TestBWPostgresql_DeleteAll(t *testing.T) {
	testInit()
	bw = bw.Map(User{}, DB)

	_, err := bw.DeleteAll(new(User))
	assert.Nil(t, err)

	number, err := bw.Count(new(User))
	assert.Nil(t, err)
	assert.Equal(t, 0, number)
}
