package events

type Fetcher interface {
	Fetch(int) ([]Event, error)
}

type Processor interface {
	Process(Event) error
}
type Type int

const (
	Unknown = iota
	Message
)

type Event struct {
	Type int
	Text string
	Meta interface{}
}
