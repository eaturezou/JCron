/* -------------------------------------------------
| Author: Zoueature
| Email: zoueature@gmail.com
| Date: 19-5-12
| Description:
| -------------------------------------------------
*/

package main

import (
	"bufio"
	"errors"
	"learning/JCron/jcron"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var msgChan = make(chan string, 5)
var connection = make(map[string]net.Conn)

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
	if err != nil {
		log.Fatalln("Start server error : " + err.Error())
	} else {
		log.Println("Start server : " + config["-h"] + ":" + config["-p"])
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("error : " + err.Error())
		}
		go doTaskModify(conn)
	}
}

func init()  {
	go dispatcher()
}

func dispatcher() {
	for {
		select {
		case msg := <-msgChan:
			var result string
			params := strings.Split(msg, " ")
			clientId := params[0]
			paramsLength := len(params)
			if paramsLength < 7 {
				result = "System error "
				err := sendMsgToClient(clientId, result)
				if err != nil {
					log.Println(err.Error())
				}
				continue
			}
			task := &jcron.Task{
				Name: msg,
				TaskFrequency: jcron.TaskFrequency{
					Second:params[1],
					Minute:params[2],
					Hour:params[3],
					Day:params[4],
					Month:params[5],
					Week:params[6],
				},
				Command:params[7],
			}
			err := jcron.New(task)
			if err != nil {
				result = err.Error()
			}
			if err != nil {
				result = err.Error()
			} else {
				result = "success"
			}
			err = sendMsgToClient(clientId, result)
			if err != nil {
				log.Println("Error: " + err.Error())
			}
		}
	}
}

func doTaskModify(conn net.Conn)  {
	connectId := conn.RemoteAddr().String()
	connection[connectId] = conn
	reader := bufio.NewScanner(conn)
	for reader.Scan() {
		msgChan<- connectId + " " + reader.Text()
	}
}

func sendMsgToClient(clientId string, msg string) error {
	connect := connection[clientId]
	if connect == nil {
		return errors.New("Send msg error, not found client, client id " + clientId)
	}
	sendNum, err := connect.Write([]byte(msg + "\r\n"))
	if err != nil {
		return errors.New("Send error, send :" + strconv.Itoa(sendNum) + "error : " + err.Error())
	}
	return nil
}