package automock

import (
	"github.com/atotto/clipboard"
	"time"
)

func main() {
	text1 := ""
	text2 := ""
	for {
		time.Sleep(300)
		text1, _ = clipboard.ReadAll()
		if len(text1) <= 5 {
			continue
		}
		if text1[0:4] != "func" || text1[len(text1)-1] != '{'{
			continue
		}
		leftB := 0
		diffB := 0
		for _, ch :=range text1{
			if ch =='(' {
				leftB++
				diffB++
			}
			if ch ==')' {
				diffB--
			}
			if diffB != 0 && diffB != 1 {
				continue
			}
		}
		if leftB > 3 {
			continue
		}
		if text1 != text2 {
			text1 = format(text1)
			//fmt.Println(text1)
			err := clipboard.WriteAll(text1)
			if err != nil {
				return
			}
			text2 = text1
			//fmt.Println(text2)
		}
	}
}

func format (before string) (after string) {


before = removeBackAndComma(before)
	bracket := true
	caller := ""
	outputs := make([]string, 0)
	mockReturn := make([]string, 0)
	funcName := ""
	i := 5
	j := 5
	//此时是结构体的方法函数
	if before[i] == '(' {
		for {
			j++
			if before[j] == ' ' {
				i = j+1
			}
			if before[j] == ')' {
				break
			}
		}
		if before[i] == '*' {
			caller = "(" + before[i:j] + ")."
		} else {
			caller = before[i:j] + "."
		}
	}

	j += 2
	i = j
	for {
		if (before[j] == ' ') || (before[j] == '('){
			break
		}
		j++
	}
	if caller == "" {
		funcName = before[i-2:j]
	} else {
		funcName = before[i:j]
	}

	for {
		j++
		if before[j] == ')' {
			break
		}
	}
	i = j//此时在输入变量后的括号上
	for {
		j++
		if before[j] == '(' {//返回值用括号包起来的情况
			break
		}
		if before[j] == '{' {
			bracket = false
			break
		}
	}
	if !bracket {
		i++
		for {
			if before[i] != ' ' {
				break
			}
			i++
		}//此时i到达唯一的变量开端
		j--
		for {
			if before[j] != ' ' {
				break
			}
			j--
		}//此时j到达唯一的变量结尾
		outputs = append(outputs, before[i:j+1])
	} else {
		j++
		i = j//ij到达返回类型的位置
		for {
			j++
			if before[j] == ' ' {
				i = j + 1//有返回名字的情况
			}
			if before[j] == ',' || before[j] == ')' {
				outputs = append(outputs, before[i:j])
				if before[j] == ')'{
					break
				}
				i = j + 1//已经获取到一个output类型
				if before[i] == ' '{
					i++
					j++
				}
			}
		}
	}
//此时获取到所有的返回类型
	for _, output := range outputs {
		switch output {
		case "int8", "int16", "int32", "int64", "byte", "rune", "int": {
			mockReturn = append(mockReturn, "0")
			continue
		}
		case "uint8", "uint16", "uint32", "uint64", "uint": {
			mockReturn = append(mockReturn, "0")
			continue
			}
		case "float32", "float64", "complex64", "complex128": {
			mockReturn = append(mockReturn, "0.0")
			continue
		}
		case "bool": {
			mockReturn = append(mockReturn, "true")
			continue
		}
		case "string": {
			mockReturn = append(mockReturn, "thisIsString")
			continue
		}
		case "uintptr": {
			mockReturn = append(mockReturn, "nil")
			continue
		}
		}
		if output[0] == '*' {
			mockReturn = append(mockReturn, "&" + output[1:] +"{}")
		} else if output == "error" {
			//mockReturn = append(mockReturn, "errors.New(\"mock err\")")
			mockReturn = append(mockReturn, "nil")
		} else if output[0:2] == "[]" {
			mockReturn = append(mockReturn, output + "{{}}")
		} else {
			mockReturn = append(mockReturn, output +"{}")
		}
	}
	if caller == ""{
		after = "Mock(" + funcName + ").Return("
	} else {
		after = "Mock(" + caller + funcName + ").Return("
	}

	for index, mReturn := range mockReturn {
		if index == len(mockReturn)-1 {
			after += mReturn
		} else {
			after += mReturn + ", "
		}
	}
	after += ").Build()"
	return after
}

func removeBackAndComma (in string) (out string) {
	for i := 0; i < len(in); i++ {
		if in[i] == '\n' || in[i] == '\t' {
			in = in[0:i] + in[i+1:len(in)]
			i = 0
		}
	}
	for i := 0; i < len(in)-1; i++ {
		if in[i] == ',' && in[i+1] == ')' {
			in = in[0:i] + in[i+1:len(in)]
			i = 0
		}
	}
	return in
}

