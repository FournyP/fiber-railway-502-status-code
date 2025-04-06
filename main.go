package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

type Request struct {
	Value string `json:"value"`
}

func main() {
	app := fiber.New()

	app.Post("/send", func(c *fiber.Ctx) error {
		remoteServer := os.Getenv("REMOTE_SERVER")

		body := &Request{
			Value: "test",
		}

		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}

		request, err := http.NewRequest(http.MethodPost, remoteServer, bytes.NewReader(bodyBytes))
		if err != nil {
			return err
		}

		request.Header.Set("X-Signature", "something")

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			return err
		}

		defer response.Body.Close()

		responseBytes, _ := io.ReadAll(response.Body)

		return c.JSON(string(responseBytes))
	})

	app.Post("/handle", func(c *fiber.Ctx) error {
		request := new(Request)
		if err := c.BodyParser(request); err != nil {
			return err
		}

		log.Printf("received request: %v", request)

		return c.SendStatus(200)
	})
}
