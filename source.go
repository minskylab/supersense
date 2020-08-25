package supersense

// Source is a new Event emitter
type Source interface {
	Run() error
	Identify(string) bool // That's only a temporal way to identify a source
	Dispose()
	Pipeline() <-chan Event
}
