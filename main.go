package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"unsafe"

	"github.com/eiannone/keyboard"
	"github.com/gorilla/websocket"
)

//インストールを行っておくもの
//go get -u github.com/eiannone/keyboard
//go get github.com/gorilla/websocket

//getlenght
func getLenght(val string, text string, match string, offset int, count int) string {
	var result string = ""
	var startPos int = strings.Index(val, text)
	var vals string = val
	//fmt.Println(startPos)

	var lens int = 0
	var cc int = 0

	lens = utf8.RuneCountInString(vals)

	for i := startPos; i < lens; i++ {
		var c string = string([]rune(vals)[i : i+1])
		//fmt.Println(c)
		if c == match {
			if cc >= count {
				break
			}
			result = ""
			cc += 1
			continue
		}
		result += c
	}

	return result
}

//unixtime
func unixTime() int64 {
	return int64(time.Now().Unix())
}

//dayofweek
func dayOfWeek() string {
	wdays := [...]string{"日", "月", "火", "水", "木", "金", "土"}
	t := time.Now()
	return wdays[t.Weekday()]
}

//stopwatch
func stopWatch() string {
	start := time.Now()
	for i := 0; i < 10000000; i++ {
	}

	// 処理
	end := time.Now()
	var sokudo int64 = end.Sub(start).Milliseconds()
	delay := strconv.FormatInt(sokudo, 10)
	return delay
}

//global 変数
var LocalFlag bool = false

func do_thread_local() {
	fmt.Println("do_thread_localを開始します。")
	for {
		if LocalFlag == false {
			break
		}
		SendMessagesToClients("o::{\"msg\":\"" + strconv.FormatInt(unixTime(), 10) + "\"}")
		fmt.Println("処理し送信ています...")
		time.Sleep(1000 * time.Millisecond)
	}
	fmt.Println("do_thread_localが終了しました。")
}

//keyboard_event
func keyboard_event() {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press ESC to quit")
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		fmt.Printf("You pressed: rune %q, key %X\r\n", char, key)
		if key == keyboard.KeyEsc {
			LocalFlag = false
			break
		}
	}
}

// WebSocket サーバーにつなぎにいくクライアント
var clients = make(map[*websocket.Conn]bool)

// クライアントから受け取るメッセージを格納
var broadcast = make(chan Message)

// WebSocket 更新用
var upgrader = websocket.Upgrader{}

// クライアントからは JSON 形式で受け取る
type Message struct {
	Message string `json:message`
}

// クライアントのハンドラ
func HandleClients(w http.ResponseWriter, r *http.Request) {
	// ゴルーチンで起動
	go broadcastMessagesToClients()
	// websocket の状態を更新
	websocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("error upgrading GET request to a websocket::", err)
	}
	// websocket を閉じる
	defer websocket.Close()

	clients[websocket] = true

	for {
		var message Message
		// メッセージ読み込み
		err := websocket.ReadJSON(&message)
		if err != nil {
			log.Printf("error occurred while reading message: %v", err)
			delete(clients, websocket)
			break
		}
		// メッセージを受け取る
		broadcast <- message
	}
}

func wsSocket() {
	// localhost:8080 でアクセスした時に index.html を読み込む
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/echo", HandleClients)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("error starting http server::", err)
		return
	}
}
func broadcastMessagesToClients() {
	for {
		// メッセージ受け取り
		message := <-broadcast
		// クライアントの数だけループ
		for client := range clients {
			//　書き込む
			err := client.WriteJSON(message)
			if err != nil {
				log.Printf("error occurred while writing message to client: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

//byte to string
func bstring(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//string to byte
func sbytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}
func SendMessagesToClients(msg string) {

	arrays := []byte(msg)
	// クライアントの数だけループ
	for client := range clients {
		//　書き込む
		err := client.WriteMessage(websocket.TextMessage, arrays)
		if err != nil {
			log.Printf("error occurred while writing message to client: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

//http get
func http_get() {
	url := "https://google.co.jp"
	req, _ := http.NewRequest("GET", url, nil)
	//ヘッダーつける場合
	//req.Header.Set("","")

	client := new(http.Client)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteArray))
}

//http post
func HttpPost(url, token, device string) error {
	jsonStr := `{"token":"` + token + `","device":"` + device + `"}`

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		return err
	}

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}
func main() {
	var result string = dayOfWeek()
	fmt.Println(result)
	LocalFlag = true
	go wsSocket()
	time.Sleep(5000 * time.Microsecond)
	go do_thread_local()
	keyboard_event()
}
