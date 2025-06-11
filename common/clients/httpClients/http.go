package httpClients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func Get(url string) {

	// Send a GET request to the API endpoint
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	// Print the response body
	fmt.Println("Response:", string(body))
}

func Post(url string, body map[string]interface{}) ([]byte, error) {

	fmt.Println(url, body)
	requestBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error encoding request body:", err)
		return []byte{}, err
	}

	// Define the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return []byte{}, err
	}

	authKey, present := os.LookupEnv("OPENAI_API_KEY")
	if !present {
		fmt.Println("OPENAI_API_KEY missing")
	}
	bearerAuthKey := "Bearer " + (authKey)

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerAuthKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return []byte{}, err
	}
	defer resp.Body.Close()

	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return respbody, nil

}
