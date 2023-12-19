package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"golang.design/x/clipboard"
)

func main() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		ch := clipboard.Watch(context.TODO(), clipboard.FmtText)
		temp := ""
		blocking := false
		go func() {
			flag := false
			for data := range ch {
				if blocking {
					flag = false
					continue
				}
				if !flag {
					output(string(data))
					println("Converted!")
					flag = true
				} else {
					flag = false
				}
			}
		}()
		for {
			fmt.Scanf("%s", &temp)

			switch strings.ToLower(temp) {
			case "s":
				blocking = !blocking
				fmt.Printf("blocking: %v\n", blocking)
			case "exit":
				return
			}
			temp = ""
		}
	} else {
		fpath := os.Args[1]

		bytes, err := os.ReadFile(fpath)
		if err != nil {
			panic(err)
		}

		output(string(bytes))
		println("Converted!")
	}
}

func output(inputText string) {
	lines := strings.Split(inputText, "\n")
	for i := 0; i < len(lines); i++ {
		lines[i] = fmt.Sprintf(`printf("%v\n");`, strings.TrimSpace(lines[i]))
	}
	text := strings.Join(lines, "\n")

	clipboard.Write(clipboard.FmtText, []byte(fmt.Sprintf(`#include <iostream>
using namespace std;

int main() {
%v

return 0;
}`, text)))
}
