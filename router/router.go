package router

type Rest interface {
	Load(args ...interface{}) error
	Find(args ...interface{}) error
	Insert(args ...interface{}) error
	Update(args ...interface{}) error
	Delete(args ...interface{}) error
}

type FileRest interface {
	Upload(args ...interface{}) error
}
