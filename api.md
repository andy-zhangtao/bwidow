# bwidow
--
    import "github.com/andy-zhangtao/bwidow"


## Usage

```go
const (
	BW_MONGO_ENDPOINT = "BW_ENV_MONGO_ENDPOINT"
	BW_MONGO_USER     = "BW_ENV_MONGO_USER"
	BW_MONGO_PASSWD   = "BW_ENV_MONGO_PASSWD"
	BW_MONGO_DB       = "BW_ENV_MONGO_DB"
)
```

```go
const (
	//DRIVER_MONGO Mongo驱动
	DRIVER_MONGO = iota
)
```

#### type BW

```go
type BW struct {
}
```


```go
var Widow *BW
```

#### func  GetWidow

```go
func GetWidow() *BW
```
GetWidow 获取当前全局Widow. 如果没有则创建

#### func (*BW) Delete

```go
func (this *BW) Delete(uPtr interface{}, field []string) (num int, err error)
```
Delete 删除命中的所有数据 uPtr 供定位记录的数据 field 用于筛选的字段

#### func (*BW) Driver

```go
func (this *BW) Driver(driver int) (err error)
```
Driver 设置使用的数据库类型 当前支持的类型为: DRIVER_MONGO - Mongo

#### func (*BW) FindAll

```go
func (this *BW) FindAll(uPtr interface{}, aPtr interface{}) (err error)
```
FindAll 通过u的字段查询所有数据 BW会解析u的字段,然后将所有非空字段作为查询条件进行查询，同时将查询到的数据赋值给a u必须为指针
a必须为array/slice类型的指针

#### func (*BW) FindAllWithSort

```go
func (this *BW) FindAllWithSort(uPtr interface{}, aArray interface{}, sortField []string) (err error)
```
FindAllWithSort 通过u的字段查询所有数据并且按照给定的条件进行排序
BW会解析u的字段,然后将所有非空字段作为查询条件进行查询，同时将查询到的数据赋值给a u 必须为指针 a 必须为array/slice类型的指针
sortField 需要排序的字段数组

#### func (*BW) FindOne

```go
func (this *BW) FindOne(uPtr interface{}) (err error)
```
FindOne 通过u的字段查询数据 BW会解析u的字段,然后将所有非空字段作为查询条件进行查询，同时将查询到的数据赋值给u u必须为指针

#### func (*BW) First

```go
func (this *BW) First(uPtr interface{}) (err error)
```
First 查询与u绑定的表中的首条记录 u 数据结构体指针

#### func (*BW) Map

```go
func (this *BW) Map(u interface{}, name string)
```
Map 将u与数据表进行绑定 u 数据结构体 name 数据表名

#### func (*BW) Save

```go
func (this *BW) Save(u interface{}) (err error)
```
Save 插入单条数据

#### func (*BW) SaveAll

```go
func (this *BW) SaveAll(uArray interface{}) (err error)
```
SaveAll 插入一批次的数据 支持数组内存在不同数据类型 u 为数组

#### func (*BW) Update

```go
func (this *BW) Update(uPtr interface{}, field []string) (num int, err error)
```
Update 更新命中的所有数据. uPtr 供定位记录的数据 field 用于筛选的字段

#### type BWDriver

```go
type BWDriver interface {
	//Check 驱动自检
	Check() error
	DriverInit() error
	Map(u interface{}, name string)
	// contains filtered or unexported methods
}
```


#### type BWMongo

```go
type BWMongo struct {
}
```


#### func (*BWMongo) Check

```go
func (this *BWMongo) Check() (err error)
```

#### func (*BWMongo) DriverInit

```go
func (this *BWMongo) DriverInit() (err error)
```

#### func (*BWMongo) Map

```go
func (this *BWMongo) Map(u interface{}, name string)
```
