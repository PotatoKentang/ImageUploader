package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	go app.Post("/upload", func(c *fiber.Ctx) error {
		// Get the file from the form field "source":
		file, err := c.FormFile("source")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing or invalid 'source' field in the request body",
			})
		}

		// Open the file to get an io.Reader
		fileContent, err := file.Open()
		if err != nil {
			return err
		}
		defer fileContent.Close()

		// Read the file content into a buffer
		var buffer bytes.Buffer
		_, err = io.Copy(&buffer, fileContent)
		if err != nil {
			return err
		}

		// Create a new request to the target API
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add the file to the request
		part, err := writer.CreateFormFile("source", file.Filename)
		if err != nil {
			return err
		}
		part.Write(buffer.Bytes())

		// Add other form fields if needed
		writer.WriteField("key", "6d207e02198a847aa98d0a2a901485a5")
		writer.WriteField("action", "upload")
		writer.WriteField("format", "json")

		// Close the multipart writer
		writer.Close()

		// Create the HTTP request
		req, err := http.NewRequest("POST", "https://freeimage.host/api/1/upload", body)
		if err != nil {
			return err
		}

		// Set the content type header
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Perform the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Read the response body
		var responseBody bytes.Buffer
		_, err = io.Copy(&responseBody, resp.Body)
		if err != nil {
			return err
		}
		var parsedBody map[string]interface{}
		if err := json.Unmarshal([]byte(responseBody.String()), &parsedBody); err != nil {
			return err
		}
		// Check the response status
		if resp.StatusCode != http.StatusOK {
			return c.Status(resp.StatusCode).JSON(fiber.Map{
				"data":parsedBody,
			})
		}

		// Return a success response
		return c.JSON(fiber.Map{
			"data":parsedBody,
		})
	})

	log.Fatal(app.Listen("localhost:5757"))
}
