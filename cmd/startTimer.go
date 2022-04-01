/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"challenge/pkg/proto"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// startTimerCmd represents the startTimer command
var startTimerCmd = &cobra.Command{
	Use:   "startTimer",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		name, _ := cmd.Flags().GetString("name")
		frequency, _ := cmd.Flags().GetInt("frequency")
		duration, _ := cmd.Flags().GetInt("duration")

		log.Default().Printf("connecting to : %d port with name %s, frequency %d and duration %d", port, name, frequency, duration)

		conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := proto.NewChallengeServiceClient(conn)

		res, err := c.StartTimer(context.Background(), &proto.Timer{Name: name, Seconds: int64(duration), Frequency: int64(frequency)})

		if err == nil {
			for {
				resp, err := res.Recv()
				println("resp:", resp, "err: ", err)
				if err == io.EOF {
					fmt.Println("Done ")
					break

				}
				if err != nil {
					log.Fatalf("cannot receive %v", err)
				}
				fmt.Printf("%+v\n", resp)

			}

		} else {
			println(err)
		}
	},
}

func init() {
	clientCmd.AddCommand(startTimerCmd)
	startTimerCmd.Flags().IntP("port", "p", 8088, " port thats used by server")
	startTimerCmd.Flags().IntP("duration", "d", 10, " timer duration")
	startTimerCmd.Flags().IntP("frequency", "f", 1, " message frequency")
	startTimerCmd.Flags().StringP("name", "n", "def_name", " name of the timer. Every new client trying to start a timer with the name of existing timer will be automatically subscribed to existing one.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startTimerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startTimerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
