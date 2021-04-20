package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
const Filepath string = "/usr/local/todo/todos.txt"

var (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Cyan   = "\033[36m"
	Purple = "\033[35m"
	Reset  = "\033[0m"
)

func init() {
	if runtime.GOOS == "windows" {
		Red = ""
		Green = ""
		Cyan = ""
		Reset = ""
		Purple = ""

	}
}

func main() {
	// Sub-commands
	listCommand := flag.NewFlagSet("ls", flag.ExitOnError)
	addCommand := flag.NewFlagSet("add", flag.ExitOnError)
	removeCommand := flag.NewFlagSet("rm", flag.ExitOnError)
	checkCommand := flag.NewFlagSet("check", flag.ExitOnError)
	unCheckCommand := flag.NewFlagSet("uncheck", flag.ExitOnError)

	// 'ls' sub-command flags
	doneTodos := listCommand.Bool("done", false, "✅ Show only checked TODOs")
	progressTodos := listCommand.Bool("prog", false, "🕝 Show only the in-progress TODOs")

	flag.Parse()

	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
		fmt.Println(Cyan + "Subcommand is required, Example Usage:" + Reset + "\n")
		fmt.Println("	todo ls [-done/-prog]")
		fmt.Println("	todo add")
		fmt.Println("	todo rm ")
		fmt.Println("	todo check")
		fmt.Println("	todo uncheck")

		os.Exit(1)
	}

	// Check the sub-command
	switch os.Args[1] {
	case "ls":
		listCommand.Parse(os.Args[2:])
	case "add":
		addCommand.Parse(os.Args[2:])
	case "rm":
		removeCommand.Parse(os.Args[2:])
	case "check":
		checkCommand.Parse(os.Args[2:])
	case "uncheck":
		unCheckCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if listCommand.Parsed() {
		Show(*doneTodos, *progressTodos)
	} else if addCommand.Parsed() {
		Add()
	} else if removeCommand.Parsed() {
		Remove()
	} else if checkCommand.Parsed() {
		Edit(true)
	} else if unCheckCommand.Parsed() {
		Edit(false)
	} else {
		flag.PrintDefaults()
	}

}

func Add() {
	f, err := os.OpenFile(Filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	check(err)

	defer f.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print(Purple, "📝 Enter your TODO: ", Reset)

	todo, _ := reader.ReadString('\n')

	// convert CRLF to LF
	todo = strings.Replace(todo, "\n", "", -1)

	n, err := f.WriteString("🕝 " + todo + "\n")

	check(err)

	fmt.Printf(Green + "✨ Your TODO has saved. (%s bytes)" + Reset + "\n", strings.TrimSpace(strconv.Itoa(n)))

}

func Edit(isCheck bool) {
	content, err := ioutil.ReadFile(Filepath)
	check(err)

	Show(false, false)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(Purple, "👉 Enter the TODO number(index): ", Reset)

	i, _ := reader.ReadString('\n')
	i = strings.TrimSpace(i)

	num, err := strconv.Atoi(i)
	check(err)

	lines := strings.Split(string(content), "\n")

	if isCheck {
		lines[num] = strings.Replace(lines[num], "🕝", "✅", 1)
	} else {
		lines[num] = strings.Replace(lines[num], "✅", "🕝", 1)
	}

	output := strings.Join(lines, "\n")

	err = ioutil.WriteFile(Filepath, []byte(output), 0644)

	fmt.Printf(Green + "✨ Your TODO has been edited." + Reset + "\n")

}

func Remove() {
	content, err := ioutil.ReadFile(Filepath)
	check(err)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(Purple, "👉 Enter the TODO number(index): ", Reset)

	i, _ := reader.ReadString('\n')
	i = strings.TrimSpace(i)
	num, err := strconv.Atoi(i)

	check(err)

	lines := strings.Split(string(content), "\n")

	lines[num] = ""

	lines = append(lines[:num], lines[num+1:]...)

	output := strings.Join(lines, "\n")

	err = ioutil.WriteFile(Filepath, []byte(strings.TrimSpace(output)), 0644)

	fmt.Println(lines)
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}

	return string(ln), err
}

func Show(done bool, inProgress bool) {
	f, err := os.OpenFile(Filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	check(err)

	r := bufio.NewReader(f)

	s, e := Readln(r)

	i := 0
	for e == nil {
		number := "[" + strconv.Itoa(i) + "] "
		if done {
			matched := strings.Contains(s, "✅")

			if matched {
				fmt.Println(number + s)
			}
		} else if inProgress {
			matched := strings.Contains(s, "🕝")
			if matched {
				fmt.Println(number + s)
			}
		} else {
			fmt.Println(number + s)
		}

		s, e = Readln(r)
		i++
	}
}
