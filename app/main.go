package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-redis/redis"
)

var client *redis.Client
var indexHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Redis App</title>
</head>
<body>
	<form action="/set_key" method="post">
		<label>
			Key:
			<input type="text" name="key">
		</label>
		<label>
			Value:
			<input type="text" name="value">
		</label>
		<button type="submit">Set Key</button>
	</form>
	<form action="/get_key">
		<label>
			Key:
			<input type="text" name="key">
		</label>
		<button type="submit">Get Key</button>
	</form>
	<form action="/del_key" method="post">
		<label>
			Key:
			<input type="text" name="key">
		</label>
		<button type="submit">Delete Key</button>
	</form>
</body>
</html>
`

func main() {
	client = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/set_key", setKeyHandler)
	http.HandleFunc("/get_key", getKeyHandler)
	http.HandleFunc("/del_key", delKeyHandler)

	log.Fatal(http.ListenAndServe(":8089", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index").Parse(indexHTML)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func setKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key := r.FormValue("key")
	value := r.FormValue("value")

	if key == "" || value == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := client.Set(key, value, 0).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	val, err := client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprint(w, val)
}

func delKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key := r.FormValue("key")

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := client.Del(key).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
