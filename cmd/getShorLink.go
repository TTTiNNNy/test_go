/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"challenge/pkg/proto"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// getShorLinkCmd represents the getShorLink command
var getShorLinkCmd = &cobra.Command{
	Use:   "getShorLink",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		url, _ := cmd.Flags().GetString("url")

		log.Default().Printf("connecting to : %d port with link %s", port, url)

		conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := proto.NewChallengeServiceClient(conn)
		res, err := c.MakeShortLink(context.Background(), &proto.Link{Data: url})

		map_resp := map[string]interface{}{}
		json.Unmarshal([]byte(res.GetData()), &map_resp)

		if err == nil {
			fmt.Printf("got short link: %s\n\r", map_resp["link"].(string))

		} else {
			fmt.Printf("error during server interruction: %s", err.Error())

		}

	},
}

func init() {
	clientCmd.AddCommand(getShorLinkCmd)
	getShorLinkCmd.Flags().IntP("port", "p", 8088, " server port")
	getShorLinkCmd.Flags().StringP("url", "u", "https://www.rust-lang.org/", " Link that will be shorter via https://dev.bitly.com/ service")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getShorLinkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getShorLinkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
