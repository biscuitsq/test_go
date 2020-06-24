package main

import "unicode/utf8"
import "fmt"
import "strings"
import "time"

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
func main(){
	fmt.Println(unixTime())
}