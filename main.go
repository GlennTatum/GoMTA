package main

import (
	"fmt"

	pb "github.com/GlennTatum/protofiles/gtfs"
	"google.golang.org/protobuf/proto"

	"io"
	"log"
	"net/http"
)

type Transit struct {
	accessToken string
}

func (t *Transit) getURL(url string) pb.FeedMessage {

	client := &http.Client{}

	fmt.Println(client)

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

func main() {

	t := Transit{"ACCESS_KEY"}

	message := t.getURL("https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-ace")

	for _, entity := range message.Entity {

		// If the field is empty while reading the protobuf a SIGSEGV event will occur

		// When looping over FeedMessage the parent field should be referenced

		if entity.TripUpdate != nil {
			fmt.Println(*entity.TripUpdate.Trip.RouteId)

		}

		if entity.Vehicle != nil {
			fmt.Println(*entity.Vehicle.StopId)
		}

	}

	// Next steps: Parse json add MTA struct and functions

}
