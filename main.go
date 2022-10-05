package main

import (
	"fmt"
	"sort"
	"time"

	pb "github.com/GlennTatum/GoMTA/mta"

	"github.com/golang/protobuf/proto"

	"io"
	"log"
	"net/http"
)

var lines = map[string]string{
	"ACE": "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-ace",
}

type Client struct {
	accessToken string
	http        *http.Client
}

type Message struct {
	header         Header
	stopTimeUpdate StopTimeUpdate
}

type Header struct {
	message pb.FeedMessage
}

type StopTimeUpdate struct {
	message pb.FeedMessage
}

func NewClient(token string) *Client {
	base := &http.Client{
		Timeout: time.Second * 10,
	}

	return &Client{
		accessToken: token,
		http:        base,
	}
}

func (c *Client) recieve(url string) []byte {

	client := c.http

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("x-api-key", c.accessToken)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body

}

func (c *Client) getFeed(url string) Message {

	resp := c.recieve(url)

	feed := &pb.FeedMessage{}
	err := proto.Unmarshal(resp, feed)
	if err != nil {
		log.Fatal(err)
	}

	return Message{header: Header{*feed}, stopTimeUpdate: StopTimeUpdate{*feed}}

}

type TripId struct {
	tripid string
}

type Stop struct {
	stop string
}

type ArrivalTime struct {
	arrivalTime int
}

type DateTime struct {
	datetime string
}

func UnixtoDateTime(t int) string {

	datetime := time.Unix(int64(t), 0)

	return datetime.String()

}

func (stc *StopTimeUpdate) filter(stopid string) map[DateTime]Stop {

	feed := stc.message

	ul := make(map[ArrivalTime]Stop)

	for _, entity := range feed.Entity {

		if entity.TripUpdate != nil {

			tripUpdate := entity.TripUpdate

			stopTimeUpdate := tripUpdate.StopTimeUpdate

			for _, s := range stopTimeUpdate {

				if *s.StopId == stopid {

					ul[ArrivalTime{int(*s.Arrival.Time)}] = Stop{*s.StopId}

				}
			}

		}

	}

	ol := make(map[DateTime]Stop)

	keys := make([]int, 0, len(ul))

	for k := range ul {
		keys = append(keys, k.arrivalTime)
	}

	sort.Ints(keys)

	for i, _ := range keys {

		ol[DateTime{UnixtoDateTime(keys[i])}] = Stop{stopid}
	}

	return ol
}

func main() {

	mtaclient := NewClient(
		"ACCESS_KEY",
	)

	mta := mtaclient.getFeed("https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-ace")

	fmt.Println(mta.stopTimeUpdate.filter("A02N"))

}
