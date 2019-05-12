/* -------------------------------------------------
| Author: Zoueature
| Email: zoueature@gmail.com
| Date: 19-5-12
| Description:
| -------------------------------------------------
*/

package main

import (
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
)

func main() {
	args := os.Args[1:]
	config := make(map[string]string)
	var key []string
	var value []string
	for index, argv := range args {
		if index % 2 == 0 {
			key = append(key, string(index))
		} else {
			value = append(value, argv)
		}
	}
	valueLen := len(value)
	for i := 0; i < len(key); i ++ {
		configName := key[i]
		if valueLen <= i {
			log.Fatalln("Start Server Error " + configName + "No Value")
		}
		configValue := value[i]
		if configName == "-h" {
			matched, err := regexp.Match("", []byte(configValue))
			if err != nil || matched == false {
				log.Fatalln("Error Host : " + configValue)
			}
		} else if configName == "-p" {
			port, _ := strconv.Atoi(configValue)
			if port < 1024 || port > 65535 {
				log.Fatalln("Illegal Port :" + configValue)
			}
		}
		config[configName] = configValue
	}
	if value, ok := config["-h"]; value == "" || !ok {
		config["-h"] = "127.0.0.1"
	}
	if value, ok := config["-p"]; value == "" || !ok {
		config["-p"] = "4698"
	}
	listener, err := net.Listen("tcp", config["-h"] + ":" + config["-p"])
	if err != err {
		log.Fatalln("Start server error : " + err.Error())
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("error : " + err.Error())
		}
		go doTaskModify(conn)
	}
}

func doTaskModify()  {
	
}

func parseQuery(conn net.Conn) (*TaskModify, error) {
	return &TaskModify{}, nil
}