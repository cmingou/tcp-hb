/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/cmingou/tcp-hb/internal/client"
	"github.com/spf13/cobra"
)

var (
	serverAddr string
	port       int
	interval   int
	timeout    int
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "client mode",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("client mode: %s:%d, interval: %d, timeout: %d\n", serverAddr, port, interval, timeout)
		cli := client.NewHeartbeatClient(serverAddr, port, interval, timeout)
		cli.Connect()
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	clientCmd.Flags().StringVarP(&serverAddr, "server", "s", "localhost", "server address")
	clientCmd.Flags().IntVarP(&port, "port", "p", 8080, "server port")
	clientCmd.Flags().IntVarP(&interval, "interval", "i", 1, "interval in seconds")
	clientCmd.Flags().IntVarP(&timeout, "timeout", "t", 5, "timeout in seconds")
}
