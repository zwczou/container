package container

type Container interface {
	Add(Provider)
	Load() error
	Exit()
	Set(...any)
	Get(...any) error
	MustGet(...any)
	Pub(string, ...any)
	TryPub(string, ...any)
}

type Provider interface {
	Name() string
	Load(Container) error
	Exit()
}

type Queuer interface {
	Sub()
	Unsub()
	Listen(fn any)
}
