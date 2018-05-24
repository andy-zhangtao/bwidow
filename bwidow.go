/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package bwidow

import (
	"reflect"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/23.

var Widow *BW

const (
	//DRIVER_MONGO Mongo驱动
	DRIVER_MONGO = iota
)

func init() {
	Widow = new(BW)
	Widow.client = make(map[int]BWDriver)
}

type BWDriver interface {
	//Check 驱动自检
	Check() error
	DriverInit() error
	Map(u interface{}, name string)
	First(uPtr interface{}) error
	FindOne(uPtr interface{}) error
	FindAll(uPtr interface{}, aPtr interface{}) error
	FindAllWithSort(uPtr interface{}, aArray interface{}, sortField []string) error
	Save(u interface{}) error
	SaveAll(uArray []interface{}) error
	Update(uPtr interface{}, field []string) error
}

type BW struct {
	driver int
	client map[int]BWDriver
}

// GetWidow 获取当前全局Widow. 如果没有则创建
func GetWidow() (*BW) {
	if Widow == nil {
		Widow = new(BW)
	}

	return Widow
}

//Driver 设置使用的数据库类型
//当前支持的类型为:
//DRIVER_MONGO - Mongo
func (this *BW) Driver(driver int) (err error) {
	switch driver {
	case DRIVER_MONGO:
		bm := BWMongo{}
		if err = bm.DriverInit(); err != nil {
			return err
		}
		Widow.driver = DRIVER_MONGO
		Widow.client[DRIVER_MONGO] = &bm
	}

	return
}

//First 查询与u绑定的表中的首条记录
//u 数据结构体指针
func (this *BW) First(uPtr interface{}) (err error) {
	return this.client[this.driver].First(uPtr)
}

//Map 将u与数据表进行绑定
//u 数据结构体
//name 数据表名
func (this *BW) Map(u interface{}, name string) {
	this.client[this.driver].Map(u, name)
}

//FindOne 通过u的字段查询数据
//BW会解析u的字段,然后将所有非空字段作为查询条件进行查询，同时将查询到的数据赋值给u
//u必须为指针
func (this *BW) FindOne(uPtr interface{}) (err error) {
	return this.client[this.driver].FindOne(uPtr)
}

//FindAll 通过u的字段查询所有数据
//BW会解析u的字段,然后将所有非空字段作为查询条件进行查询，同时将查询到的数据赋值给a
//u必须为指针
//a必须为array/slice类型的指针
func (this *BW) FindAll(uPtr interface{}, aPtr interface{}) (err error) {
	return this.client[this.driver].FindAll(uPtr, aPtr)
}

//FindAllWithSort 通过u的字段查询所有数据并且按照给定的条件进行排序
//BW会解析u的字段,然后将所有非空字段作为查询条件进行查询，同时将查询到的数据赋值给a
//u 必须为指针
//a 必须为array/slice类型的指针
//sortField 需要排序的字段数组
func (this *BW) FindAllWithSort(uPtr interface{}, aArray interface{}, sortField []string) (err error) {
	return this.client[this.driver].FindAllWithSort(uPtr, aArray, sortField)
}

//Save 插入单条数据
func (this *BW) Save(u interface{}) (err error) {
	return this.client[this.driver].Save(u)
}

//SaveAll 插入一批次的数据
//支持数组内存在不同数据类型
//u 为数组
func (this *BW) SaveAll(uArray interface{}) (err error) {
	value := reflect.ValueOf(uArray)

	var uu []interface{}
	for i := 0; i < value.Cap(); i++ {
		uu = append(uu, value.Index(i).Interface())
	}

	return this.client[this.driver].SaveAll(uu)
}

//Update 更新命中的所有数据.
//uPtr 供定位记录的数据
//field 用于筛选的字段
func (this *BW) Update(uPtr interface{}, field []string) (err error) {
	return this.client[this.driver].Update(uPtr, field)
}
