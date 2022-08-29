package main

import (
	"fmt"

	pb "github.com/GlennTatum/protofiles/gtfs"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"io"
	"log"
	"net/http"
)

func main() {

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-ace", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("x-api-key", "ACCESS_KEY")
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

	m := protojson.Format(feed)

	fmt.Println(m)

}
