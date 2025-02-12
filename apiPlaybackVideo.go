package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// rewirte all
func playbackHandler(c *gin.Context) {
	eventID := c.Param("id")

	events := map[string]Event{
		"1": {
			ID: "1",

			Timestamp:   202411,
			Description: "Car detected",
			Location:    "Camera_TEST2",
			Result:      "PID 7282 | 7420837",
		},
		"2": {
			ID:          "2",
			Image:       "/static/img/event2.jpg",
			Timestamp:   7122024,
			Description: "Motorbike detected",
			Location:    "Camera_TEST1",
			Result:      "PID 7281 | 7420472",
		},
	}

	event, exists := events[eventID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.HTML(http.StatusOK, "playback.tmpl", gin.H{
		"event": event,
	})
}

func HTTPAPIPlaybackAll(c *gin.Context) {
	starttime := c.Query("starttime")
	endtime := c.Query("endtime")

	fmt.Println("Start Time:", starttime)
	fmt.Println("End Time:", endtime)
	events := []Event{
		{
			ID:          "1",
			Image:       "/static/img/event1.jpg",
			Timestamp:   202411,
			Description: "Car detected",
			Camera:      "Camera_TEST2",
			Location:    "Location_TEST2",
			Result:      "PID 7282 | 7420837",
			FullImage:   "/static/img/event1_full.jpg",
		},
		{
			ID:          "2",
			Image:       "/static/img/event2.jpg",
			Timestamp:   2024122,
			Description: "Motorbike detected",
			Camera:      "Camera_TEST1",
			Location:    "Location_TEST1",
			Result:      "PID 7281 | 7420472",
			FullImage:   "/static/img/event2_full.jpg",
		},
	}

	c.HTML(http.StatusOK, "playback.tmpl", gin.H{
		"events":  events,
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "playback",
	})

	fmt.Println("hello")

}
