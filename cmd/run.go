package cmd

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/soonkuk/zmq/lib/node"
	"github.com/soonkuk/zmq/lib/server"
	"github.com/spf13/cobra"
)

var nodeCmd *cobra.Command
var serverCmd *cobra.Command

func init() {
	nodeCmd = &cobra.Command{
		Use:   "client",
		Short: "Run client simulator",
		Run: func(c *cobra.Command, args []string) {
			nodeCount, _ := strconv.ParseInt(os.Args[2], 10, 64)
			fmt.Println("Run client simulator!")
			for i := 0; i < int(nodeCount); i++ {
				// go node_task("device"+strconv.Itoa(time.Now().Nanosecond()), "sband")
				go node_task("device-"+strconv.Itoa(i), "sband")
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			}
			count := 0
			for {
				count++
			}
		},
	}
	rootCmd.AddCommand(nodeCmd)

	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Run server simulator",
		Run: func(c *cobra.Command, args []string) {
			fmt.Println("Run server simulator!")
			s, err := server.NewServerZmq()
			defer s.Stop()
			if err != nil {
				log.Print(err)
			}
			if err = s.Init(); err != nil {
				log.Print(err)
				os.Exit(1)
			}
			s.Run()
		},
	}
	rootCmd.AddCommand(serverCmd)
}

func node_task(id string, dtype string) {
	configNode := node.NewConfigNode(id, dtype)
	n, err := node.NewNodeBand(configNode)
	// defer n.Stop()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	n.Init()
	n.Run()
}
