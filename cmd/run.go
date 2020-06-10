package cmd

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	guuid "github.com/google/uuid"
	"github.com/soonkuk/zmq/lib/common"
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
			var conf common.Conf
			conf.GetConf()
			log.Printf("Run %d client simulator with %d duration !", conf.Total, conf.ReportDuration)
			if checkConfirm() {
				go node_task(conf)
			} else {
				os.Exit(1)
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
			log.Println("Run server simulator!")
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

func node_task(conf common.Conf) {
	// Make a collection of failure node index randomly
	fNodeMap := map[int]bool{}
	for len(fNodeMap) < conf.Fail {
		i := rand.Intn(conf.Total)
		if _, ok := fNodeMap[i]; !ok {
			fNodeMap[i] = false
		}
	}
	for i := 0; i < int(conf.Total); i++ {
		var status node.NodeStatus = node.CorrectNode
		_, exist := fNodeMap[i]
		if exist {
			status = node.FailNode
		}
		configNode := node.NewConfigNode(guuid.New().String(), conf.DeviceType, conf.ReportDuration, status)
		n, err := node.NewNodeBand(configNode)
		// defer n.Stop()
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		n.Init()
		n.Run()
	}
	// n.InitAndRun()
}

func checkConfirm() bool {
	var input string
	for {
		fmt.Print("Confirm to running client (y/n)")
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Println(err)
		}
		switch strings.ToLower(input) {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Println("Pleas type y or n and press Enter")
		}
	}
}
