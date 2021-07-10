package main

import (
	"github.com/atotto/clipboard"
	"log"
	"os"
	"time"
)

func main() {
	text1 := ""
	text2 := ""
	Info := log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Info.Println("\nautomock监听剪切板中…………\n如果偶遇程序崩溃请重新启动程序,Thx ^.^")
	for {
		time.Sleep(500000000)//剪切板触发为0.5s一次，节约cpu资源
		text1, _ = clipboard.ReadAll()
		if text1 == text2 {
			//fmt.Println("剪切板未更新" + text1)
			continue
		}
		//过滤规则1
		if len(text1) <= 5 {
			text2 = text1
			//fmt.Println("过滤规则1过滤")
			continue
		}
		//过滤规则2
		if text1[0:4] != "func" || text1[len(text1)-1] != '{'{
			text2 = text1
			//fmt.Println("过滤规则2过滤")
			continue
		}
		//过滤规则3
		if !filter3(text1) {
			text2 = text1
			//fmt.Println("过滤规则3过滤")
			continue
		}
		////过滤规则4
		//if !filter4(text1) {
		//	text2 = text1
		//	fmt.Println("过滤规则4过滤")
		//	continue
		//}
		//自动写入剪切板
		if text1 != text2 {
			text1 = format(text1)
			err := clipboard.WriteAll(text1)
			if err != nil {
				return
			}
			text2 = text1
			//fmt.Println(text2)
		}
	}
}

//过滤规则3
func filter3 (in string) (out bool) {
	leftB := 0
	diffB := 0
	for _, ch :=range in{
		if ch =='(' {
			leftB++
			diffB++
		}
		if ch ==')' {
			diffB--
		}
		if diffB != 0 && diffB != 1 {
			return false
		}
	}
	//只接受括号对数在3对以及以下的
	return leftB <= 3
}

//过滤规则4：改在内部去过滤较方便
//func filter4 (in string) (out bool) {
//	i := len(in)-2
//	for {
//		if in[i] == ')' {
//			return false
//		} else if in[i] == ' ' {
//			i--
//		} else {
//			return true
//		}
//	}
//}

