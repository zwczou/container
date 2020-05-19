package container

var app = New()

func App() *Container {
	return app
}

func Front(ext Provider) {
	app.Front(ext)
}

func Push(ext Provider) {
	app.Push(ext)
}

func Before(name string, ext Provider) {
	app.Before(name, ext)
}

func After(name string, ext Provider) {
	app.After(name, ext)
}

func All() []Provider {
	return app.All()
}

func Load() error {
	return app.Load()
}

func Exit() {
	app.Exit()
}
