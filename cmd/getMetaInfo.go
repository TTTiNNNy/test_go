/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"challenge/pkg/proto"
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// getMetaInfoCmd represents the getMetaInfo command
var getMetaInfoCmd = &cobra.Command{
	Use:   "getMetaInfo",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")

		log.Default().Printf("connecting to : %d port", port)

		conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := proto.NewChallengeServiceClient(conn)
		r, err := c.ReadMetadata(context.Background(), &proto.Placeholder{Data: "Data"})

		if err == nil {
			fmt.Printf("got metaData: %s\n\r", r.GetData())

		} else {
			fmt.Printf("error during server interruction: %s", err.Error())
		}

	},
}

func init() {
	clientCmd.AddCommand(getMetaInfoCmd)
	getMetaInfoCmd.Flags().IntP("port", "p", 8088, " server port")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getMetaInfoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getMetaInfoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
