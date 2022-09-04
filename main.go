package main

import (
	"fmt"

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

type Transit struct {
	accessToken string
}

func (t *Transit) getURL(url string) pb.FeedMessage {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("x-api-key", t.accessToken)
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

func (t *Transit) getSubwayStop(line string, stop string) {
	feed := t.getURL(lines[line])

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

func (t *Transit) getLineStops(line string) {
	feed := t.getURL(lines[line])

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

	t := Transit{"ACCESS_KEY"}

	// t.getSubwayStop("ACE", "A02N")

	t.getLineStops("ACE")

	// Next steps: Parse json add MTA struct and functions
}
