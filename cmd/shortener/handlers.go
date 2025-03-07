package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// Функция для генерации сокращённого URL
func generateShortURL(longURL string) string {
	// Генерируем уникальный идентификатор (uid) при помощи пакета go.uuid
	uuidObj := uuid.NamespaceURL
	uuidStr := uuidObj.String()
	uuidStr = strings.ReplaceAll(uuidStr, "-", "")
	uid := uuidStr[:8]
	//Создаем запись в хранилище (мапе) urlStore, ключ- uid, значение- longURL
	UrlStore[uid] = longURL
	//
	return "/" + uid
}

func Handlers() *http.ServeMux {
	mux := http.NewServeMux()
	//HandleFunc регистрирует функцию-обработчик для данного шаблона.
	//Здесь должны обработаться запросы с методом POST и путем "/"
	mux.HandleFunc("POST /{$}", func(w http.ResponseWriter, r *http.Request) {
		//Для нужной работы конечной точки будем смотреть поля структуры Request.
		//Читаем тело запроса- поле Body.
		//Поле Body имеет тип io.ReadCloser и данные имеют такой непосредственный вид:
		//&{0xc0001a8000 <nil> <nil> false true {0 0} true false false 0x762820}
		//func io.ReadAll(r io.Reader) ([]byte, error)
		param, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		//Вроде так нужно
		defer r.Body.Close()

		// Преобразуем тело запроса (тип []byte) в строку:
		longURL := string(param)

		// Генерируем сокращённый URL
		shortURL := r.Host + generateShortURL(longURL)

		// Устанавливаем статус ответа 201
		w.WriteHeader(http.StatusCreated)
		// Устанавливаем Content-Type как text/plain
		//w.Header().Set("Content-Type", "text/plain")

		// // Версию HTTP можно узнать так
		// httpVersion := r.Proto

		// Отправляем сокращённый URL в теле ответа
		fmt.Fprint(w, shortURL)
	})

	//Аккуратнее с / !!! Неделя мучений из- за "GET /{id}/" вместо "GET /{id}"
	//У этого даже название есть (!!):
	//Trailing-slash redirection
	//Чтобы слэши нормализовались в gorilla/mux используется
	//r := mux.NewRouter().StrictSlash(true)
	mux.HandleFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		longURL, err := UrlStore[id]

		//В случае ошибки, т.е. запрос не соответствует сохраненным данным, возвращаем статус 400
		if !err {
			http.Error(w, "URL not found", http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", longURL)
		//И так и так работает. Оставил первоначальный вариант.
		//http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
		w.WriteHeader(http.StatusTemporaryRedirect)

	})

	return mux
}
