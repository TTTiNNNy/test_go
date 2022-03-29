package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	. "challenge/pkg/proto"

	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type TestRPCServer struct {
}

func (TestRPCServer) MakeShortLink(ctx context.Context, link *Link) (*Link, error) {
	println("env var: ")

	viper.SetEnvPrefix("BITLY_OAUTH")

	viper.BindEnv("LOGIN")
	//env_login := viper.Get("LOGIN").(string)

	viper.BindEnv("TOKEN")
	env_token := viper.Get("TOKEN").(string)

	// println(env_login)
	println(env_token)

	client := &http.Client{}

	jsonOutBody, _ := json.Marshal(map[string]string{
		"long_url": link.Data,
	})
	req, err := http.NewRequest("POST", "https://api-ssl.bitly.com/v4/shorten", strings.NewReader(string(jsonOutBody)))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", env_token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)

	// resp, err := http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(jsonOutBody))
	// client := &http.Client{}
	return &Link{Data: string(bodyText)}, nil
}

var timerPull = make(map[string]Timer)

func (TestRPCServer) StartTimer(timer *Timer, srv ChallengeService_StartTimerServer) error {

	if exTimer, ok := timerPull[timer.Name]; ok {
		for exTimer.GetSeconds() > 0 {
			exTimer = timerPull[timer.Name]

			if err := srv.Send(&exTimer); err != nil {
				log.Printf("send error %v", err)
				status.Errorf(codes.Unimplemented, "method StartTimer failed during Send method")
			}
			time.Sleep(time.Duration(exTimer.Frequency) * time.Second)
		}
	} else {
		newTimer := Timer{Name: timer.Name, Frequency: timer.Frequency, Seconds: timer.Seconds}
		timerPull[timer.Name] = newTimer

		for newTimer.GetSeconds() > 0 {

			timerPull[timer.Name] = newTimer
			println("sec:", newTimer.Seconds, "sec/freq: ", newTimer.GetSeconds()/newTimer.Frequency)
			if err := srv.Send(&newTimer); err != nil {
				log.Printf("send error %v", err)
				status.Errorf(codes.Unimplemented, "method StartTimer failed during Send method")
			}
			newTimer.Seconds -= newTimer.Frequency
			time.Sleep(time.Duration(timer.Frequency) * time.Second)
		}
	}
	delete(timerPull, timer.Name)

	return nil
}

func (TestRPCServer) ReadMetadata(ctx context.Context, pl *Placeholder) (*Placeholder, error) {
	out := "DefaultDataServerString"
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		val, _ := json.Marshal(md)
		out = string(val)
	}

	return &Placeholder{Data: out}, nil
}
