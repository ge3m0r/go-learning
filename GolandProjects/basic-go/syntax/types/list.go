package types

type list interface {
	Add(index int, val any)
	Append(val any)
	Delete(index int)
}
type newlist list

//衍生类型调用不了原油的方法
