/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package bwidow

import (
	"reflect"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/23.

/*
##### Example

```go
	//BWidow Init
	bw := bwidow.GetWidow()
	err := bw.Driver(bwidow.DRIVER_MONGO)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Init Error": err}).Errorln(ModuleName)
		return
	}
	//Setting Model And Table
	bw.Map(User{}, "devex_user_copy")
	u := User{
		CurrentAuthority: "BW",
	}

	//Delete All Matched Record
	num, err := bw.Delete(&u, []string{"currentauthority"})
	if err != nil {
		logrus.Errorln(err)
		return
	}

	logrus.WithFields(logrus.Fields{"change": num}).Info(ModuleName)
```
*/
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
	Check() error
	DriverInit() error
	Map(u interface{}, name string)
	first(uPtr interface{}) error
	findOne(uPtr interface{}) error
	findAll(uPtr interface{}, aPtr interface{}) error
	findAllWithSort(uPtr interface{}, aArray interface{}, sortField []string) error
	save(u interface{}) error
	saveAll(uArray []interface{}) error
	update(uPtr interface{}, field []string) (int, error)
	delete(uPtr interface{}, field []string) (int, error)
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
/*
##### Example

```go
	err := bw.Driver(bwidow.DRIVER_MONGO)
    if err != nil {
    	logrus.WithFields(logrus.Fields{"Init Error": err}).Errorln(ModuleName)
    	return
    }
```
*/
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
/*
#### Example

```go
	u := User{}
	err = bw.First(&u)
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.WithFields(logrus.Fields{"First": u}).Info(ModuleName)
```
*/
func (this *BW) First(uPtr interface{}) (err error) {
	return this.client[this.driver].first(uPtr)
}

//Map 将u与数据表进行绑定
//u 数据结构体
//name 数据表名
/*
##### Example

```go
	bw.Map(User{}, "devex_user_copy")
```
*/
func (this *BW) Map(u interface{}, name string) {
	this.client[this.driver].Map(u, name)
}

//FindOne 通过u的字段查询数据
//BW会解析u的字段,然后将所有非空字段作为查询条件进行查询，同时将查询到的数据赋值给u
//u必须为指针
/*
##### Example

```go
	u := User{
		Name: "ztao8607@gmail.com",
	}
	logrus.WithFields(logrus.Fields{"Query": u}).Info(ModuleName)

	err = bw.FindOne(&p)
	if err != nil {
		logrus.Error(err)
		return
	}
```
*/
func (this *BW) FindOne(uPtr interface{}) (err error) {
	return this.client[this.driver].findOne(uPtr)
}

//FindAll 通过u的字段查询所有数据
//BW会解析u的字段,然后将所有非空字段作为查询条件进行查询，同时将查询到的数据赋值给a
//u 查找记录使用的数据
//a必须为array/slice类型的指针
/*

##### Example

``` go
	u := User{
		CurrentAuthority: "dev",
	}

	var allUser []User
	err = bw.FindAll(u, &allUser)
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.WithFields(logrus.Fields{"FindAll": allUser}).Info(ModuleName)
```

*/
func (this *BW) FindAll(u interface{}, aPtr interface{}) (err error) {
	return this.client[this.driver].findAll(u, aPtr)
}

//FindAllWithSort 通过u的字段查询所有数据并且按照给定的条件进行排序
//BW会解析u的字段,然后将所有非空字段作为查询条件进行查询，同时将查询到的数据赋值给a
//u 查找记录使用的数据
//a 必须为array/slice类型的指针
//sortField 需要排序的字段数组
/*
##### Example

```go
	u := User{
		CurrentAuthority: "dev",
	}
	var allUser []User
	err = bw.FindAllWithSort(u, &allUser, []string{"+name"})
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.WithFields(logrus.Fields{"FindAllWithSort": allUser}).Info(ModuleName)
```
*/
func (this *BW) FindAllWithSort(uPtr interface{}, aArray interface{}, sortField []string) (err error) {
	return this.client[this.driver].findAllWithSort(uPtr, aArray, sortField)
}

//Save 插入单条数据
/*
##### Example

```go
	u := User{
		ID:   bson.NewObjectId(),
		Name: "fromBW@gmail.com",
	}

	err = bw.Save(u)
	if err != nil {
		logrus.Errorln(err)
		return
	}
```
*/
func (this *BW) Save(u interface{}) (err error) {
	return this.client[this.driver].save(u)
}

//SaveAll 插入一批次的数据
//支持数组内存在不同数据类型
//u 为数组
/*
##### Example

```go
	al := []User{
		{
			ID:   bson.NewObjectId(),
			Name: "fromBW1@gmail.com",
		},
		{
			ID:   bson.NewObjectId(),
			Name: "fromBW2@gmail.com",
		},
		{
			ID:   bson.NewObjectId(),
			Name: "fromBW3@gmail.com",
		},
	}

	err = bw.SaveAll(al)
	if err != nil {
		logrus.Errorln(err)
		return
	}
```
*/
func (this *BW) SaveAll(uArray interface{}) (err error) {
	value := reflect.ValueOf(uArray)

	var uu []interface{}
	for i := 0; i < value.Cap(); i++ {
		uu = append(uu, value.Index(i).Interface())
	}

	return this.client[this.driver].saveAll(uu)
}

//Update 更新命中的所有数据.
//uPtr 供定位记录的数据
//field 用于筛选的字段
/*
##### Example

```go
	u := User{
		Name:             "fromBW2@gmail.com",
		CurrentAuthority: "BW",
	}

	num, err := bw.Update(&u, []string{"name"})
	if err != nil {
		logrus.Errorln(err)
		return
	}

	logrus.WithFields(logrus.Fields{"change": num}).Info(ModuleName)
```
*/
func (this *BW) Update(uPtr interface{}, field []string) (num int, err error) {
	return this.client[this.driver].update(uPtr, field)
}

//Delete 删除命中的所有数据
//uPtr 供定位记录的数据
//field 用于筛选的字段
/*

##### Example

```go
	u := User{
		CurrentAuthority: "BW",
	}

	//Delete All Matched Record
	num, err := bw.Delete(&u, []string{"currentauthority"})
	if err != nil {
		logrus.Errorln(err)
		return
	}
```

*/
func (this *BW) Delete(uPtr interface{}, field []string) (num int, err error) {
	return this.client[this.driver].delete(uPtr, field)
}
