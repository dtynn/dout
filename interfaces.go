package dout

type Cache interface {
	Set(key, val interface{}) error
	Setex(key, val interface{}, expire int64) error
	Get(key interface{}) (interface{}, error)
	Del(key interface{}) error
}

type Storage interface {
	SaveMail()
	UpdateMailProcess()
}

type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
}
