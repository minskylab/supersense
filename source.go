package supersense

// Source is a new Event emmiter
type Source interface {
	Run() error
	Events() *chan Event
}
