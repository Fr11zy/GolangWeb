package main


import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"hw3/model"
)


// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	//Замена ioutil.ReadAll на bufio.Scanner, чтобы сразу читать построчно
	scanner := bufio.NewScanner(file)

	//Используем strings.Builder, чтобы исключить конкатенацию строк
	var builder strings.Builder
	seenBrowsers := []string{}

	// Делаем сразу все в одном цикле, чтобы занимало меньше времени и памяти
	i:=-1
	for scanner.Scan() {
		i++
		line := scanner.Bytes()
		user := model.User{}
		err := user.UnmarshalJSON(line)
		if err != nil {
			panic(err)
		}
		
		isAndroid := false
		isMSIE := false

		for _,browser := range user.Browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					seenBrowsers = append(seenBrowsers, browser)
				}
			}

			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					seenBrowsers = append(seenBrowsers, browser)
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
			}

		email := strings.Replace(user.Email, "@", " [at] ", -1)
		builder.WriteString(fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Fprintln(out, "found users:\n"+builder.String())
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
