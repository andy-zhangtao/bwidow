/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package bwidow

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/23.

var Widow *BW

const (
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
	First(u interface{}) error
	FindOne(u interface{}) error
	FindAll(u interface{}, a interface{}) error
	FindAllWithSort(u interface{}, a interface{}, sortField []string) error
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
func (this *BW) First(u interface{}) (err error) {
	return this.client[this.driver].First(u)
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
func (this *BW) FindOne(u interface{}) (err error) {
	return this.client[this.driver].FindOne(u)
}

//FindAll 通过u的字段查询所有数据
//BW会解析u的字段,然后将所有非空字段作为查询条件进行查询，同时将查询到的数据赋值给a
//u必须为指针
//a必须为array/slice类型的指针
func (this *BW) FindAll(u interface{}, a interface{}) (err error) {
	return this.client[this.driver].FindAll(u, a)
}

//FindAllWithSort 通过u的字段查询所有数据并且按照给定的条件进行排序
//BW会解析u的字段,然后将所有非空字段作为查询条件进行查询，同时将查询到的数据赋值给a
//u 必须为指针
//a 必须为array/slice类型的指针
//sortField 需要排序的字段数组
func (this *BW) FindAllWithSort(u interface{}, a interface{}, sortField []string) (err error) {
	return this.client[this.driver].FindAllWithSort(u, a, sortField)
}
