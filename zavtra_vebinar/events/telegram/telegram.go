package telegram

import (
	"errors"
	"time"

	"log"
	"os"
	"strconv"

	"ZavtraVebinar/clients/telegram"
	"ZavtraVebinar/events"
	"ZavtraVebinar/lib/e"
)

var ChatID []int //Срез ID чатов

type Processor struct { //удовлетворяет интерфейсам Fetcher и Processor
	tg     *telegram.Client
	offset int
}

type Meta struct { //Тип, который относится только к телеграм
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("неизвестный тип события")
	ErrUnknownMetaType  = errors.New("неизвестный тип меты")
)

func New(client *telegram.Client) *Processor {
	return &Processor{
		tg: client,
	}
}

func meta(event events.Event) (Meta, error) { //получаем мету
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("Невозможно получить мету", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event { //Преобразует апдейты в ивенты
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type { //возвращает тип события, 0 - неизвестно, 1- сообщение
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) { //получает события
	var z events.Event
	updates, err := p.tg.Updates(p.offset, limit) //получаем апдейты
	if err != nil {
		return nil, e.Wrap("Невозможно получить апдейты", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates { //перебираем апдейты и преобразуем в тип event
		z = event(u)

		res = append(res, z)
	}

	p.offset = updates[len(updates)-1].ID + 1 //получаем новую пачку апдейтов при следующем вызове Fetch
	log.Print(p.offset)
	return res, nil
}

func (p *Processor) ProcessExistingCHATIDList() { //РАБОТАЕТ
	ChatID := ReadChatIDList()

	for _, chatID := range ChatID {
		go p.Sender(chatID)

	}

	ChatID = nil //обнуление, чтобы не было задвоек
}

func (p *Processor) ProcessNewCHATID(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("Невозможно обработать сообщение", err)
	}

	log.Printf("Новое сообщение %s от %s", event.Text, meta.Username)

	p.AppendNewChatID(meta.ChatID)

	return nil

}

func (p *Processor) Sender(ID int) {

	for {
		T := time.Now()
		if T.Weekday().String() == "Monday" || T.Weekday().String() == "Thursday" {
			if T.Hour() == 21 {
				p.tg.SendMessage(ID, "Завтра вебинар.Начало в 09.15!!!")
				time.Sleep(time.Hour * 1)
			}

		}
	}

	//	for i := 0; i < 1; i++ {
	//		T := time.Now()
	//		log.Print(T.Second())
	//		p.tg.SendMessage(ID, "Завтра вебинар.Начало в 09.15!!!")
	//		p.tg.SendPhoto(ID)
	//		time.Sleep(time.Second * 1)
	//	}
}

func Check(chatID int, ChatID []int) bool { //Проверяет, есть ли в массиве элемент chatID

	for _, v := range ChatID {
		if v == chatID {
			return true
		}
	}
	return false
}

func (p *Processor) AppendNewChatID(chatID int) { //РАБОТАЕТ
	ChatID := ReadChatIDList()

	if !Check(chatID, ChatID) {
		file, err := os.OpenFile("ChatIDList.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		file.Write([]byte((strconv.Itoa(chatID)) + " "))

		go p.Sender(chatID)

		defer file.Close()
	}

	ChatID = nil //обнуление, чтобы не было задвоек

}

func ReadChatIDList() []int { //РАБОТАЕТ
	var Z []byte
	var C []int
	data, _ := os.ReadFile("ChatIDList.txt")
	//ПРОБЕЛ=32

	for _, v := range data {
		if v == 32 {
			c, _ := strconv.Atoi(string(Z))
			C = append(C, c)
			Z = nil
			continue
		}
		Z = append(Z, v)

	}
	return C
}

//рассортировать функции работы с chatID по функциям
//Сделать так, чтобы бот не требовал писать сообщения
//[ERR] consumer: Невозможно получить апдейты: can't get updates: Невозможно сделать HTTP запрос: Get "https://api.telegram.org/bot5006959166:AAFrCXi3qnRphI3HdwMWgTBxYJrN3kQHh1g/getUpdates?limit=100&offset=346658914": read tcp 172.40.49.153:58507->149.154.167.220:443: wsarecv: A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond.
