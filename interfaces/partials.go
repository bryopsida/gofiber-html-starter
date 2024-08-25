package interfaces

type IRequestContext interface {
	Locals(key interface{}, value ...interface{}) interface{}
}
