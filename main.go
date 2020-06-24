package main

import (
	"unicode/utf8"
	"fmt"
	"strings"
	"time"
	"strconv"
	"github.com/eiannone/keyboard"
)
//インストールを行っておくもの
//go get -u github.com/eiannone/keyboard


//getlenght
func getLenght(val string, text string, match string, offset int, count int) string{
	
	var result string = ""
	var startPos int = strings.Index(val,text)
	var vals string = val;
	//fmt.Println(startPos)
	
	var lens int = 0
	var cc int = 0

	lens = utf8.RuneCountInString(vals)

	for i := startPos; i < lens; i++ {
		var c string = string([]rune(vals)[i:i + 1])
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
func unixTime() int64{
	return int64(time.Now().Unix()) 
}
//dayofweek
func dayOfWeek() string{
	wdays := [...]string{"日", "月", "火", "水", "木", "金", "土"}
	t := time.Now()
	return wdays[t.Weekday()]
}
//stopwatch
func stopWatch() string{
	start := time.Now();
	for i := 0; i < 10000000; i++ {
	}

	// 処理
	end := time.Now();
	var sokudo int64 = end.Sub(start).Milliseconds()
	delay := strconv.FormatInt(sokudo, 10)
	return delay
}
//global 変数
var LocalFlag bool = false
func do_thread_local() {
	fmt.Println("do_thread_localを開始します。")
    for {
		if LocalFlag == false{
			break
		}
        fmt.Println("処理しています...")
        time.Sleep(1000 * time.Millisecond)
	}
	fmt.Println("do_thread_localが終了しました。")
}
//keyboard_event
func keyboard_event(){
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
func main(){
	var result string = dayOfWeek()
	fmt.Println(result)
	stopWatch()
	LocalFlag = true
	go do_thread_local()
	keyboard_event()
	time.Sleep(10000 * time.Millisecond)
}