package main

import (
	"strings"
)

func format (before string) (after string) {
	before = removeBackAndComma(before)
	bracket := true //返回参数是否有括号包围的标识符
	caller := ""
	baseType  := []string{
		"int8", "int16", "int32", "int64", "byte", "rune",
		"int", "uint8", "uint16", "uint32", "uint64", "uint",
		"float32", "float64", "complex64", "complex128", "bool", "string", "uintptr"}
	outputsName := make([]string, 0)
	outputsType := make([]string, 0)
	mockReturn := make([]string, 0)
	funcName := ""
	i := 5 //过滤func
	j := 5 //过滤func

	//结构体的方法函数判断
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
		//获取到调用者的名字
		if before[i] == '*' {
			caller = "(" + before[i:j] + ")."
		} else {
			caller = before[i:j] + "."
		}
	}

	//过度到输入变量的索引位置
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

	//此时过度在在输入变量后的括号上
	i = j
	for {
		j++
		if before[j] == '(' { //返回值用括号包起来的情况
			break
		}
		if before[j] == '{' { //返回值没有用括号包起来的情况
			bracket = false
			break
		}
	}
	//返回值没有用括号包起来的情况处理
	if !bracket {
		i++
		for {
			if before[i] != ' ' {
				break
			}
			i++
		} //此时i到达唯一的变量开端
		j--
		for {
			if before[j] != ' ' {
				break
			}
			j--
		} //此时j到达唯一的变量结尾
		//防止传入没有返回参数的函数，返回传入的字符串
		if i >= j+1 {
			return before
		}
		outputsType = append(outputsType, before[i:j+1])
	} else { //返回值用括号包起来的情况
		j++
		i = j //ij到达返回最开头的位置
		for {
			j++
			if before[j] == ')' {
				break
			}
		}
		//此时ij包括返回值
		outputs := make([]string, 0)
		outputs = strings.Split(before[i:j], ",")
		StringIsType2 := false
		for index, eachStr := range outputs {
			if eachStr[0] == ' ' {
				eachStr = eachStr[1:]
				outputs[index] = eachStr
			}
			if eachStr[len(eachStr)-1] == ' ' {
				eachStr = eachStr[:len(eachStr)]
				outputs[index] = eachStr
			}
			if strings.Contains(eachStr, " ") {
				StringIsType2 = true //主要有一个有空格，那此时不是全为类型名的情况，存在批量声明
				break
			}
		}

		//不是全为类型名的情况，存在批量声明：  startTime time.Time, endTime, performanceMonth time.Time,
		if StringIsType2 {
			for index, eachStr := range outputs {
				if eachStr[0] == ' ' {
					eachStr = eachStr[1:]
					outputs[index] = eachStr
				}
				if eachStr[len(eachStr)-1] == ' ' {
					eachStr = eachStr[:len(eachStr)]
					outputs[index] = eachStr
				}
				if strings.Contains(eachStr, " ") {
					withSpace := strings.Split(eachStr, " ")
					outputsName = append(outputsName, withSpace[0])
					outputsType = append(outputsType, withSpace[1])
				} else { //没有空格，就只有变量名，需要调用后面的批量
					outputsName = append(outputsName, eachStr)
					outputsType = append(outputsType,"")
				}
			}
		} else { //全为类型名的情况不用处理，直接去输出 StringIsType2 为false时
			outputsType = outputs
		}
	}
	//处理类似于func GetAccessibleDataRange(timeRange *ContentMmpPerformance.TimeRange, currentTime time.Time) (
	//	startTime, endTime , performanceMonth time.Time,
	//) { 这种批量声明类型的方式，返回值，
	// 一个返回只有类型应该是全部都只有类型，所以只要有一个既有名字又有类型就可以断定为第二种，第二种这种情况只出现一个的话那就是名字不是类型了
	// startTime time.Time, endTime, performanceMonth time.Time,
	for index, outputType := range outputsType {
		if outputType == "" {
			for _, afterOutput := range outputsType[index+1:] {
				if afterOutput == "" {
					continue
				} else {
					outputsType[index] = afterOutput
					break
				}
			}
		}
	}

	//此时获取到所有的返回类型，后面是字符串拼接的过程
	for _, output := range outputsType {
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
			mockReturn = append(mockReturn, "\"thisIsString\"")
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
			isBaseType := false
			for _, each := range baseType {
				if output[2:] == each {
					mockReturn = append(mockReturn, output +"{}")
					isBaseType = true
					break
				}
			}
			if !isBaseType {
			mockReturn = append(mockReturn, output +"{{}}")
			}
		}
	}
	//有无调用者（是否为方法函数）的拼接方式
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

//前置处理：去除所选字符串的所有回车、tab、以及回车前的"，"
func removeBackAndComma (in string) (out string) {

	for i := 0; i < len(in); i++ {
		if in[i] == '\n' || in[i] == '\t' {
			in = in[:i] + in[i+1:]
			i = 0
		}
	}
	for i := 0; i < len(in)-1; i++ {
		if in[i] == ',' && in[i+1] == ')' {
			in = in[:i] + in[i+1:]
			i = 0
		}
	}
	return in
}
