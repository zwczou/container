package container

var app = New()

func App() *Container {
	return app
}

func Use(ext Provider) {
	app.Use(ext)
}

func Pre(ext Provider) {
	app.Pre(ext)
}

func Extensions() (exts []Provider) {
	return app.Extensions()
}

func Load() error {
	return app.Load()
}

func Exit() {
	app.Exit()
}
