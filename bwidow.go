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
	Check() error
	DriverInit() error
	Map(u interface{}, name string)
	First(u interface{})
	FindOne(u interface{})
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

func (this *BW) First(u interface{}) {
	this.client[this.driver].First(u)
}

func (this *BW) Map(u interface{}, name string) {
	this.client[this.driver].Map(u, name)
}

func (this *BW) FindOne(u interface{}) {
	this.client[this.driver].FindOne(u)
}
