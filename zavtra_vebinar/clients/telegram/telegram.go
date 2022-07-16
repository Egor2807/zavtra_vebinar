package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"ZavtraVebinar/lib/e"
)

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

type Client struct {
	host     string //host api сервиса телеграм
	basePath string //префикс, с которого начинаются все запросы(host+basePath)
	client   http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
	SendPhotoMethod   = "sendPhoto"
)

func New(host string, token string) *Client { //создает новый клиент для телеграм
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

//получение и отправка сообщений пользователям

func (c *Client) Updates(offset int, limit int) (updates []Update, err error) { //получает сообщения от пользователей
	defer func() { err = e.WrapIfErr("can't get updates", err) }()

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error { //отправляет сообщения пользователям
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("Невозможно отправить сообщение", err)
	}

	return nil
}

func (c *Client) SendPhoto(chatID int) error { //Прикол, отправляет ФОТКИ
	var photo string = "https://yandex.ru/images/search?utm_source=main_stripe_big&text=%D0%93%D1%80%D0%B0%D0%BD%D0%B4%20%D0%9A%D0%B0%D0%BD%D1%8C%D0%BE%D0%BD&nl=1&source=morda&pos=2&rpt=simage&img_url=https%3A%2F%2Fst.depositphotos.com%2F2577341%2F3142%2Fi%2F950%2Fdepositphotos_31426201-stock-photo-horse-shoe-bend.jpg&lr=213"
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("photo", photo)

	_, err := c.doRequest(SendPhotoMethod, q)
	if err != nil {
		return e.Wrap("Невозможно отправить сообщение", err)
	}

	return nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("Невозможно сделать HTTP запрос", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil) //HTTP запрос
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }() //закрываем тело ответа

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
