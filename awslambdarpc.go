/*
awslambdarpc is an utility to make requests to your local AWS Lambda.

This tool is a CLI, and running awslambdarpc help will show your options.
It uses the client package for the real interaction with AWS Lambda, you can
import and use it if you wish.

Usage:

	awslambdarpc [options]

Available options:

	-a, --address
		the address of your local running function, defaults to localhost:8080
	-e, --event
		path to the event JSON to be used as input
	-d, --data
		data passed to the function as input, in JSON format, defaults to "{}"
	-dl, --deadLine
		timeout, in timestamp format, defaults to now()+15 seconds
	help, -h, --help
		show this help

To make a request to a lambda function running at localhost:3000 and passing
the contents of a file, events/input.json as payload:

	awslambdarpc -a localhost:3000 -e events/input.json

You can do passing the data directly with the -d flag:

	awslambdarpc -a localhost:3000 -d '{"body": "Hello World!"}'
*/
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/blmayer/awslambdarpc/client"
)

const help = `awslambdarpc is an utility to make requests to your local AWS Lambda
Usage:
  awslambdarpc [options]
Available options:
  -a
  --address	the address of your local running function, defaults to localhost:8080
  -e
  --event	path to the event JSON to be used as input
  -d
  --data	data passed to the function as input, in JSON format, defaults to "{}"
  -dl 
  --deadLine
  			timeout in seconds
  help
  -h
  --help	show this help
Examples:
  awslambdarpc -a localhost:3000 -e events/input.json
  awslambdarpc -a localhost:3000 -d '{"body": "Hello World!"}'`

func main() {
	addr := "localhost:8080"
	payload := []byte("{}")
	var eventFile string
	var deadLineSeconds int64 = 15

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-a", "--address":
			i++
			addr = os.Args[i]
		case "-dl", "--deadLine":
			i++
			var err error
			deadLineSeconds, err = strconv.ParseInt(os.Args[i], 10, 0)
			if err != nil {
				os.Stderr.WriteString("Invalid deadLine value: " + err.Error() + "\n")
				os.Exit(-3)
			}
			fmt.Println("Given deadLineSeconds: ", deadLineSeconds)

		case "-e", "--event":
			i++
			eventFile = os.Args[i]

			// Read event file
			if os.Args[i] != "" {
				f, err := os.Open(eventFile)
				if err != nil {
					os.Stderr.WriteString("error opening file: " + err.Error() + "\n")
					os.Exit(-3)
				}

				fileInfo, _ := f.Stat()
				content := make([]byte, fileInfo.Size())
				n, err := f.Read(content)
				if int64(n) != fileInfo.Size() {
					os.Stderr.WriteString("error: could not read whole file" + "\n")
					os.Exit(-4)
				}
				if err != nil {
					os.Stderr.WriteString("error reading file: " + err.Error() + "\n")
					os.Exit(-5)
				}
				payload = content
			}
		case "-d", "--data":
			i++
			payload = []byte(os.Args[i])
		case "-h", "--help", "help":
			println(help)
			os.Exit(0)
		default:
			os.Stderr.WriteString("error: wrong argument\n")
			println(help)
			os.Exit(-1)
		}
	}

	res, err := client.Invoke(addr, payload, deadLineSeconds)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(-2)
	}
	println(string(res))
}
