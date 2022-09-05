package main

import (
	"fmt"
	"time"

	pb "github.com/GlennTatum/GoMTA/mta"

	"google.golang.org/protobuf/proto"

	"io"
	"log"
	"net/http"
)

// Reading from the stops.txt file would be a better option

var lines = map[string]string{
	"ACE": "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-ace",
}

type Client struct {
	accessToken string
	http        *http.Client
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

func (c *Client) getURL(url string) pb.FeedMessage {

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

	feed := &pb.FeedMessage{}
	err = proto.Unmarshal(body, feed)
	if err != nil {
		log.Fatal(err)
	}

	// For a json response
	// j := protojson.Format(feed)
	// return j

	// For a protobuf response

	return *feed

}

func (c *Client) SubwayStop(line string, stop string) {
	feed := c.getURL(lines[line])

	for _, entity := range feed.Entity {

		if entity.TripUpdate != nil {

			tripUpdate := entity.TripUpdate

			stopTimeUpdate := tripUpdate.StopTimeUpdate

			for _, s := range stopTimeUpdate {

				if *s.StopId == stop {
					fmt.Println("Found", stop)
				}
			}
		}

	}
}

func (c *Client) LineStops(line string) {
	feed := c.getURL(lines[line])

	for _, entity := range feed.Entity {
		if entity.TripUpdate != nil {

			tripUpdate := entity.TripUpdate

			stopTimeUpdate := tripUpdate.StopTimeUpdate

			for _, s := range stopTimeUpdate {
				fmt.Println(*s.StopId)
			}
		}
	}
}

func main() {

	t := NewClient(
		"ACCESS_KEY",
	)

	t.SubwayStop("ACE", "A02N")

	t.LineStops("ACE")

	// Next steps: Parse json add MTA struct and functions
}
