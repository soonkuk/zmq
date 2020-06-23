package server

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	zmq "github.com/pebbe/zmq4"
	"github.com/soonkuk/zmq/lib/common"
)

func server_worker(i int) {

	worker, _ := zmq.NewSocket(zmq.DEALER)
	defer worker.Close()
	worker.Connect("inproc://backend")

	for {
		//  The DEALER socket gives us the reply envelope and message
		msg, _ := worker.RecvMessage(0)
		var response []string
		identity, content := pop(msg)
		var data common.Query
		// msgpack.Unmarshal([]byte(content[0]), &data)
		//  Send 0..4 replies back
		json.Unmarshal([]byte(content[0]), &data)
		/*
			if data.temperature > 37.5 {
				color.Set(color.FgRed)
				log.Println(i, identity, string(d))
				color.Unset()
			} else {
				log.Println(i, identity, string(d))
			}
		*/
		log.Println(i, identity, data)
		if verifyMsg(content) {
			response = append(response, "Success")
		} else {
			response = append(response, "Fail")
		}
		// replies := rand.Intn(5)
		time.Sleep(time.Duration(rand.Intn(10)+1) * time.Millisecond)
		worker.SendMessage(identity, response)

	}
}

func pop(msg []string) (head, tail []string) {
	if msg[1] == "" {
		head = msg[:2]
		tail = msg[2:]
	} else {
		head = msg[:1]
		tail = msg[1:]
	}
	return
}

func verifyMsg(msg []string) bool {
	if len(msg) < 1 {
		return false
	} else {
		return true
	}
}
