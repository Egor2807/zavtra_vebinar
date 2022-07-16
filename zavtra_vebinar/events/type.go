package events

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	ProcessNewCHATID(e Event) error
	ProcessExistingCHATIDList()
}

type Type int

const (
	Unknown Type = iota
	Message
)

type Event struct {
	Type Type        //тип события, 0-неизвестно,1 -сообщение
	Text string      //текст сообщения
	Meta interface{} //ИНТЕРФЕЙС ДЛЯ СЛУЧАЕВ,КОГДА CHATID не int, или нет UserName.Другие мессенджеры
}
