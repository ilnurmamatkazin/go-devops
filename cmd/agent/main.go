package main

import (
	"bytes"
	"context"
	"encoding/binary"

	"fmt"
	"io"

	"net/http"

	"time"
)

const (
	url            = "127.0.0.1"
	port           = 8080
	pollInterval   = 2
	reportInterval = 10
)

func main() {

}

func sendMetric(ctxBase context.Context, client *http.Client, typeMetric, nameMetric string, value interface{}) (err error) {
	fmt.Println("sendMetrics")

	ctx, cancel := context.WithTimeout(ctxBase, 1*time.Second)
	defer cancel()

	endpoint := fmt.Sprintf("http://%s:%d/update/%s/%s/%v", url, port, typeMetric, nameMetric, value)

	// конструируем запрос
	// запрос методом POST должен, кроме заголовков, содержать тело
	// тело должно быть источником потокового чтения io.Reader
	// в большинстве случаев отлично подходит bytes.Buffer
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, value)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	// в заголовках запроса сообщаем, что данные кодированы стандартной URL-схемой
	request.Header.Set("Content-Type", "text/plain; charset=UTF-8")

	// отправляем запрос и получаем ответ
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	// печатаем код ответа
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	// читаем поток из тела ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	// и печатаем его
	fmt.Println(string(body))

	return
}
