package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/n-is/go-wm/wm"
)

var ConsoleReader = bufio.NewReader(os.Stdin)

func x(p *wm.Project) {
	log.Println(*p)
}

func main() {

	ws := wm.OpenWorkSpace("hike")
	pr, _ := ws.AddNewProject("helloWorld", "./tmp")
	pr.Run(x)
	pr.Save()

	wm.RemoveWorkspace("hike")
}

func prompt(str string) string {
	fmt.Print(str)
	text, _ := ConsoleReader.ReadString('\n')
	return strings.TrimSpace(text[:len(text)-1])
}
