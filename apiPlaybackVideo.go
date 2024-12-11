package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func playbackHandler(c *gin.Context) {
	eventID := c.Param("id")

	events := map[string]Event{
		"1": {
			ID: "1",

			Timestamp:   202411,
			Description: "Car detected",
			Location:    "PHUCHAN_TEST2",
			Result:      "PID 7282 | 7420837",
		},
		"2": {
			ID:          "2",
			Image:       "/static/img/event2.jpg",
			Timestamp:   7122024,
			Description: "Motorbike detected",
			Location:    "PHUCHAN_TEST3",
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
