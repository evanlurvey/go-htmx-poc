package app

var appversion string

func init() {
	// TODO: make app version not suck
	appversion = NewID()
}

func AppVersion() string {
	return appversion
}
