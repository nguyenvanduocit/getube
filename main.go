package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/kkdai/youtube/v2"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	h := server.Default(server.WithHostPorts(":"+port), server.WithStreamBody(true))
	h.GET("/:id", streamVideo)
	h.Spin()
}

func streamVideo(ctx context.Context, c *app.RequestContext) {
	videoID := c.Param("id")
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		c.Error(err)
		return
	}

	formats := video.Formats.WithAudioChannels()
	formats.Sort()
	bestFormat := formats[0]

	stream, length, err := client.GetStream(video, &bestFormat)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Range")
	c.Header("Access-Control-Expose-Headers", "Range")
	c.Header("Content-Type", "video/mp4")
	c.Header("Content-Disposition", "attachment; filename="+video.Title+".mp4")
	c.SetBodyStream(stream, int(length))
}
