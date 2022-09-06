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

type StopTimeUpdate struct {
	stopId string
	time   int64
	// track  string
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

func (c *Client) SubwayStop(line string, stop string) map[string]StopTimeUpdate {
	feed := c.getURL(lines[line])

	resp := map[string]StopTimeUpdate{}

	for _, entity := range feed.Entity {

		if entity.TripUpdate != nil {

			tripUpdate := entity.TripUpdate

			stopTimeUpdate := tripUpdate.StopTimeUpdate

			for _, s := range stopTimeUpdate {

				if *s.StopId == stop {
					resp[*tripUpdate.Trip.TripId] = StopTimeUpdate{*s.StopId, *s.Arrival.Time}
				}
			}
		}

	}

	return resp
}

func (c *Client) LineStops(line string) map[string]string {

	// Get all of the stations at an API endpoint

	feed := c.getURL(lines[line])

	resp := map[string]string{}

	for _, entity := range feed.Entity {
		if entity.TripUpdate != nil {

			tripUpdate := entity.TripUpdate

			stopTimeUpdate := tripUpdate.StopTimeUpdate

			for _, s := range stopTimeUpdate {
				resp[*s.StopId] = *tripUpdate.Trip.RouteId
			}
		}
	}

	return resp
}

func main() {

	t := NewClient(
		"ACCESS_KEY",
	)

	stop := t.SubwayStop("ACE", "A02N")

	ls := t.LineStops("ACE")

	fmt.Println(ls)

	fmt.Println()

	fmt.Println(stop)

	// Next steps: Parse json add MTA struct and functions
}
