package main

import (
	"fmt"
	"log"

	"github.com/piotrstrzalka/kvcache/pkg/client"
	"github.com/spf13/cobra"
)

type access interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}

var cacheClient access

var (
	protocol, address string
	key, value        string

	rootCmd = cobra.Command{
		Use:   "kvcli",
		Short: "Short",
		Long:  "Long long very long",

		Run: func(cmd *cobra.Command, args []string) {
			log.Println("KvCache client, choose one of available commands")
			return
		},
		PersistentPreRun: createClient,
	}
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	rootCmd.PersistentFlags().StringVarP(&protocol, "protocol", "p", "rest", "Define protocol [rest, grpc]")
	rootCmd.PersistentFlags().StringVarP(&address, "address", "a", "localhost:5010", "Define address (default: localhost:5010)")

	getCmd := &cobra.Command{
		Use: "get",
		Run: func(cmd *cobra.Command, args []string) {
			if key == "" {
				log.Fatal("Key cannot be empty")
			}

			log.Printf("Get %q\n", key)
			value, err := cacheClient.Get(key)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Value:", string(value))
		},
	}

	getCmd.Flags().StringVarP(&key, "key", "k", "", "Key for jey-value pair")

	setCmd := &cobra.Command{
		Use: "set",
		Run: func(cmd *cobra.Command, args []string) {
			if key == "" {
				log.Fatal("Key cannot be empty")
			}

			log.Printf("Set %q, %q\n", key, value)

			err := cacheClient.Set(key, []byte(value))
			if err != nil {
				log.Fatal(err)
			}

			log.Println("Stored succesfully!")
			return
		},
	}

	setCmd.Flags().StringVarP(&key, "key", "k", "", "Key for jey-value pair")
	setCmd.Flags().StringVarP(&value, "value", "v", "", "Value for jey-value pair")

	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(setCmd)
}

func createClient(cmd *cobra.Command, args []string) {
	log.Printf("protocol %q and address validation %q\n", protocol, address)
	//address regex
	var err error
	if protocol == "rest" {
		cacheClient, err = client.NewRest(address)

		if err != nil {
			log.Fatal("Cannot create Rest client")
		}
		return
	}

	if protocol == "grpc" {
		cacheClient, err = client.NewGrpc(address)
		if err != nil {
			log.Fatal("Cannot create grpc client")
		}
		return
	}

	log.Fatal("Cannot create client")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
