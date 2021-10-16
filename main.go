package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var request = "1234"
var program = ""

var maxLoops = 16777216
var loops = 0

var data [255]byte
var pointer = 0

var response = ""

func newPointer(direction int) (int, error) {
	if pointer+direction < 0 {
		return pointer, errors.New("Data underflow")
	}
	if pointer+direction > len(data)+1 {
		return pointer, errors.New("Data overflow")
	}
	return pointer + direction, nil
}

func output(value byte) {
	response += string(value)
}

func input() string {
	if len(request) == 0 {
		return "\x03"
	}
	res := string(request[0])
	request = request[1:]
	return res
}

func lookForOpeningBracket(start int) int {
	bracketCount := 1
	index := start - 1
	for ok := true; ok; ok = bracketCount > 0 && index > 0 {
		char := program[index]
		if char == 91 {
			bracketCount--
		}
		if char == 93 {
			bracketCount++
		}
		index--
	}

	return index
}

func lookForClosingBracket(start int) int {
	bracketCount := 1
	index := start + 1
	for ok := true; ok; ok = bracketCount > 0 && index < len(program)-1 {
		char := program[index]
		if char == 91 {
			bracketCount++
		}
		if char == 93 {
			bracketCount--
		}
		index++
	}

	return index
}

func run(program string) (int, error) {
	for i := 0; i < len(program); i++ {
		char := program[i]

		// debug(char)

		switch char {
		case 44: // ,
			data[pointer] = input()[0]
			break
		case 46: // .
			output(data[pointer])
			break
		case 62: // >
			p, err := newPointer(1)
			if err != nil {
				return 0, err
			}
			pointer = p
			break
		case 60: // <
			p, err := newPointer(-1)
			if err != nil {
				return 0, err
			}
			pointer = p
			break
		case 91: // [
			if data[pointer] == 0 {
				i = lookForClosingBracket(i)
			} else {
				loops++
				if loops >= maxLoops {
					return 0, errors.New("Maximum number of allowed loops (2^24)")
				}
			}
			break
		case 93: // ]
			if data[pointer] != 0 {
				i = lookForOpeningBracket(i)
			} else {
				loops = 0
			}
			break
		case 43:
			data[pointer] = (data[pointer] + 1) % 255
		case 45:
			if data[pointer] == 0 {
				data[pointer] = 255
			} else {
				data[pointer]--
			}
		}
	}
	return 1, nil
}

func debug(char byte) {
	fmt.Printf("\nChar: %s, Pointer: %d; Value: %d", string(char), pointer, data[pointer])

	println("\nData: ")
	for d := 0; d < 10; d++ {
		fmt.Printf(" %d", data[d])
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing file argument, e.g: `bfserver ./file.bf`")
		return
	}
	filename := os.Args[1]
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		println(err)
		return
	}

	// Web server stuff

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		request = r.Method + " " + r.URL.Path + " " + r.URL.RawQuery + "\n"

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		reqBody := buf.String()

		request = request + reqBody

		println(request)

		program = string(content)
		_, err := run(program)
		if err != nil {
			fmt.Fprintf(w, "<h1>The server chrased :(</h1><p>"+err.Error()+"</p>")
		} else {
			w.Write([]byte(response))
			response = ""
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
