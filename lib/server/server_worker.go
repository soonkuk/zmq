package server

import (
	"fmt"
	"math/rand"

	zmq "github.com/pebbe/zmq4"
)

func server_worker(i int) {

	worker, _ := zmq.NewSocket(zmq.DEALER)
	defer worker.Close()
	worker.Connect("inproc://backend")

	for {
		//  The DEALER socket gives us the reply envelope and message
		msg, _ := worker.RecvMessage(0)
		identity, content := pop(msg)
		fmt.Println(i, identity, content)
		//  Send 0..4 replies back

		replies := rand.Intn(5)
		for reply := 0; reply < replies; reply++ {
			//  Sleep for some fraction of a second
			// time.Sleep(time.Duration(rand.Intn(10)+1) * time.Millisecond)
			worker.SendMessage(identity, content)
		}

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
