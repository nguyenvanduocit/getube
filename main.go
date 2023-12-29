package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/cors"
	"github.com/kkdai/youtube/v2"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	h := server.Default(server.WithHostPorts(":"+port), server.WithStreamBody(true))

	h.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	h.GET("/videos/:id", getVideo)
	h.GET("/stream/:id", streamVideo)
	h.GET("/healthz", healthCheck)

	h.Spin()
}

func healthCheck(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

func getVideo(ctx context.Context, c *app.RequestContext) {
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

	streamUrl, err := client.GetStreamURL(video, &bestFormat)
	if err != nil {
		c.Error(err)
		return
	}

	// get ip of the visitor
	userIP := string(c.GetHeader("Fly-Client-IP"))
	// parse the stream url, replace ip with the visitor's ip
	parsedStreamUrl, err := url.Parse(streamUrl)
	if err != nil {
		c.Error(err)
		return
	}
	// replace the query param with the visitor's ip
	query := parsedStreamUrl.Query()
	query.Set("ip", userIP)
	parsedStreamUrl.RawQuery = query.Encode()
	streamUrl = parsedStreamUrl.String()

	c.JSON(http.StatusOK, map[string]interface{}{
		"downloadUrl": bestFormat.URL,
		"streamUrl":   streamUrl,
		"title":       video.Title,
		"thumbnail":   video.Thumbnails[0].URL,
	})
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
