package container

var app = New()

func App() *Container {
	return app
}

func Front(exts ...Provider) {
	app.Front(exts...)
}

func Push(exts ...Provider) {
	app.Push(exts...)
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
