package container

type Provider interface {
	Name() string
	Load(Container) error
	Exit()
}
