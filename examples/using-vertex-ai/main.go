package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/datasource"
)

func main() {
	// Create a new application
	a := gofr.New()

	//HTTP service with default health check endpoint
	a.AddHTTPService("anotherService", "http://localhost:9000")

	creds := a.Config.Get("SVC_ACC_CREDS")

	vertexAIClient, err := NewVertexAIClientWithKey("endless-fire-437206-j7", "us-central1",
		"us-central1-aiplatform.googleapis.com", "gemini-1.5-pro-002", creds)
	if err != nil {
		a.Logger().Fatal(err)
	}

	a.POST("/chat", func(c *gofr.Context) (interface{}, error) {
		var prompt struct {
			Prompt string `json:"prompt"`
		}

		err = c.Bind(&prompt)
		if err != nil {
			return nil, err
		}

		payload := RequestPayload{
			Contents: []Message{
				{
					Role: "user",
					Parts: []Part{
						{Text: prompt.Prompt},
					},
				},
			},
			GenerationConfig: GenerationConfig{
				Temperature:     1.0,
				MaxOutputTokens: 8192,
				TopP:            0.95,
			},
			SafetySettings: []SafetySetting{
				{Category: "HARM_CATEGORY_HATE_SPEECH", Threshold: "OFF"},
				{Category: "HARM_CATEGORY_DANGEROUS_CONTENT", Threshold: "OFF"},
				{Category: "HARM_CATEGORY_SEXUALLY_EXPLICIT", Threshold: "OFF"},
				{Category: "HARM_CATEGORY_HARASSMENT", Threshold: "OFF"},
			},
			Tools: []Tool{
				{
					Retrieval: Retrieval{
						VertexAiSearch: VertexAiSearch{
							Datastore: "projects/endless-fire-437206-j7/locations/global/collections/default_collection/dataStores/gofr-datastore_1732298621027",
						},
					},
				},
			},
		}

		respString, err := vertexAIClient.GenerateContent(payload)
		if err != nil {
			return nil, err
		}

		//return response.Raw{Data: struct {
		//	Response string `json:"response"`
		//}{resp}}, nil

		resp := make([]DataEntry, 0)

		err = json.Unmarshal([]byte(respString), &resp)
		if err != nil {
			return nil, err
		}

		return ConcatenateAllEntries(resp), nil
	})

	// Add all the routes
	a.GET("/hello", HelloHandler)
	a.GET("/error", ErrorHandler)
	a.GET("/redis", RedisHandler)
	a.GET("/trace", TraceHandler)
	a.GET("/mysql", MysqlHandler)

	// Run the application
	a.Run()
}

func HelloHandler(c *gofr.Context) (interface{}, error) {
	name := c.Param("name")
	if name == "" {
		c.Log("Name came empty")
		name = "World"
	}

	return fmt.Sprintf("Hello %s!", name), nil
}

func ErrorHandler(c *gofr.Context) (interface{}, error) {
	return nil, errors.New("some error occurred")
}

func RedisHandler(c *gofr.Context) (interface{}, error) {
	val, err := c.Redis.Get(c, "test").Result()
	if err != nil && err != redis.Nil { // If key is not found, we are not considering this an error and returning "".
		return nil, datasource.ErrorDB{Err: err, Message: "error from redis db"}
	}

	return val, nil
}

func TraceHandler(c *gofr.Context) (interface{}, error) {
	defer c.Trace("traceHandler").End()

	span2 := c.Trace("some-sample-work")
	<-time.After(time.Millisecond * 1) //nolint:wsl    // Waiting for 1ms to simulate workload
	defer span2.End()

	// Ping redis 5 times concurrently and wait.
	count := 5
	wg := sync.WaitGroup{}
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func() {
			c.Redis.Ping(c)
			wg.Done()
		}()
	}
	wg.Wait()

	//Call to Another service
	resp, err := c.GetHTTPService("anotherService").Get(c, "redis", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var data = struct {
		Data interface{} `json:"data"`
	}{}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(b, &data)

	return data.Data, nil
}

func MysqlHandler(c *gofr.Context) (interface{}, error) {
	var value int
	err := c.SQL.QueryRowContext(c, "select 2+2").Scan(&value)
	if err != nil {
		return nil, datasource.ErrorDB{Err: err, Message: "error from sql db"}
	}

	return value, nil
}
