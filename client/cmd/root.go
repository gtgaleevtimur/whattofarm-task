package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var pair string

var rootCmd = &cobra.Command{
	Use:   "GET",
	Short: "This main command to send GET to server",
	Long:  "This command sent GET to server at port localhost:3001",
	Run: func(cmd *cobra.Command, args []string) {
		rest()
	},
}

func Execute() {
	rootCmd.Flags().StringVar(&pair, "pair", "", "returned pair")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Data struct {
	EthUsdt string `json:"ETH-USDT,omitempty"`
	BtcUsdt string `json:"BTC-USDT,omitempty"`
}

func rest() {
	var request *http.Request
	var err error
	if pair == "" {
		request, err = http.NewRequest("GET", "http://localhost:3001/api/v1/rates?pairs=BTC-USDT,ETH-USDT", nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		path := "?pairs=" + pair
		request, err = http.NewRequest("GET", "http://localhost:3001/api/v1/rates"+path, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var data Data
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}
	switch pair {
	case "":
		fmt.Printf("ETH-USDT:%s, BTC-USDT:%s\n", data.EthUsdt, data.BtcUsdt)
	case "ETH-USDT":
		fmt.Printf("%s\n", data.EthUsdt)
	case "BTC-USDT":
		fmt.Printf("%s\n", data.BtcUsdt)
	}
}
