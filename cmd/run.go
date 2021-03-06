package cmd

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	guuid "github.com/google/uuid"
	"github.com/soonkuk/zmq/lib/client"
	"github.com/soonkuk/zmq/lib/common"
	"github.com/soonkuk/zmq/lib/server"
	"github.com/spf13/cobra"
)

var clientCmd *cobra.Command
var serverCmd *cobra.Command
var flagFile string = common.DefaultConfigFilePath

func init() {
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "Run client simulator",
		Run: func(c *cobra.Command, args []string) {
			var conf common.Conf
			conf.GetConf(flagFile)
			count := &common.Counter{
				Send: 0,
				Recv: 0,
			}
			log.Printf("Run %s group's %d client simulators with %d sec duration to %s endpoint %s port!", conf.GroupID, conf.Total, conf.ReportDuration, conf.EndPoint, conf.Port)
			if checkConfirm() {
				go client_task(conf, count)
			} else {
				os.Exit(1)
			}
			cnt := 0
			for {
				cnt++
			}
		},
	}
	clientCmd.Flags().StringVar(&flagFile, "cfg-file", flagFile, "config file")
	rootCmd.AddCommand(clientCmd)

	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Run server simulator",
		Run: func(c *cobra.Command, args []string) {
			var conf common.Conf
			conf.GetConf(flagFile)
			sEndPoint := "tcp://*:" + conf.Port
			log.Printf("Run server simulator with %s endpoint %s port!", sEndPoint, conf.Port)
			if !checkConfirm() {
				os.Exit(1)
			}
			cfgServer := server.NewConfigServer(conf.Workers, sEndPoint)
			s, err := server.NewServerImpl(cfgServer)
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
	serverCmd.Flags().StringVar(&flagFile, "cfg-file", flagFile, "config file")
	rootCmd.AddCommand(serverCmd)
}

func client_task(conf common.Conf, count *common.Counter) {
	// Make a collection of failure node index randomly
	fClientMap := map[int]int{}
	for len(fClientMap) < conf.Fail {
		i := rand.Intn(conf.Total)
		if _, ok := fClientMap[i]; !ok {
			fClientMap[i] = -1
		}
	}
	for len(fClientMap) < (conf.Test + conf.Fail) {
		i := rand.Intn(conf.Total)
		if _, ok := fClientMap[i]; !ok {
			fClientMap[i] = 0
		}
	}
	for i := 0; i < int(conf.Total); i++ {
		var status client.StatusClient = client.CorrectClnt
		val, exist := fClientMap[i]
		if exist {
			if val == -1 {
				status = client.FailClnt
			} else {
				status = client.TestClnt
			}
		}
		cEndPoint := "tcp://" + conf.EndPoint + ":" + conf.Port
		cfgClient := client.NewConfigClient(conf.GroupID+"||"+guuid.New().String(), conf.DeviceType, conf.ReportDuration, status, cEndPoint, count)
		c, err := client.NewClientImpl(cfgClient)
		// defer n.Stop()
		if err != nil {
			os.Exit(1)
		}
		c.Init()
		c.Run()
	}
	// c.InitAndRun()
}

func checkConfirm() bool {
	var input string
	for {
		fmt.Print("Confirm to running (y/n)")
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
