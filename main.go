package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

//Defines Json strcuture for API calls
type BatchSendMessageRequest struct {
	Number       string `json:"number"`
	TemplateName string `json:"templateName"`
	Language     string `json:"language"`
	Parameters   []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"parameters"`
}

//Defines Json structure for API call's results / responses
type BatchSendMessageResponse struct {
	Response string `json:"response"`
}

//downloads CSV file
func downloadCSVFile(url string, localPath string) error {
	// Realiza una solicitud HTTP GET para descargar el archivo
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Creates a local file where the content will be copied
	file, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copies the downloaded content in the local file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

//Reads CSV file
func readCSVFile(filePath string) ([]BatchSendMessageRequest, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';' // Establece el separador de valores a punto y coma

	var rows []BatchSendMessageRequest

	// Ignorar la primera línea del encabezado
	_, _ = reader.Read()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		//Processes CSV data and creates structure BatchSendMessageRequest
		//considering API's requirement
		row := BatchSendMessageRequest{
			Number:       record[1],
			TemplateName: "template",
			Language:     "en",
		}
		row.Parameters = append(row.Parameters, struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{
			Name:  "message",
			Value: record[3],
		})

		rows = append(rows, row)
	}

	return rows, nil
}

func yourAPIPostFunction(requestJSON []byte) (string, error) {
	// Configures HTTP request for API call
	apiURL := "https://cde176f9-7913-4af7-b352-75e26f94fbe3.mock.pstmn.io/v1/businesses/f1980bf7-c7d6-40ec-b665-dbe13620bffa/send/whatsapp_templated_messages"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		return "", err
	}

	// Sets required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Key", "<API Key>")
	req.Header.Set("x-mock-response-code", "201") // Cambia el código de respuesta si es necesario

	// Calls API
	client := &http.Client{
		Timeout: 10 * time.Second, // Establece un tiempo de espera para la solicitud
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Reads API's response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Verifies HTTP response status
	if resp.StatusCode == http.StatusCreated {
		return "OK", nil
	}

	return string(body), nil
}

func sendRequest(row BatchSendMessageRequest, responses chan BatchSendMessageResponse, index int) {
	// Converts structure in JSON
	requestJSON, err := json.Marshal(row)
	if err != nil {
		fmt.Printf("Error al convertir la solicitud en JSON (línea %d): %v\n", index, err)
		responses <- BatchSendMessageResponse{Response: "Error: JSON encoding failed"}
		return
	}

	responseBody, err := yourAPIPostFunction(requestJSON)
	if err != nil {
		fmt.Printf("Error en la solicitud a la API (línea %d): %v\n", index, err)
		responses <- BatchSendMessageResponse{Response: "Error: API request failed"}
		return
	}

	response := BatchSendMessageResponse{Response: responseBody}
	responses <- response
}

func BatchSendCampaign(csvURL string, localPath string, apiKey string) error {
	
	if err := downloadCSVFile(csvURL, localPath); err != nil {
		return fmt.Errorf("Error al descargar el archivo CSV: %v", err)
	}

	fmt.Println("Archivo CSV descargado con éxito en:", localPath)


	rows, err := readCSVFile(localPath)
	if err != nil {
		fmt.Println("Error al leer el archivo CSV:", err)
		return err
	}

	// Channel to receive the responses to the API calls
	responses := make(chan BatchSendMessageResponse, len(rows))

	// AWaitgroup to guarantee all requests were completed
	var wg sync.WaitGroup

	for i, row := range rows {
		wg.Add(1)
		go func(row BatchSendMessageRequest, index int) {
			defer wg.Done()
			sendRequest(row, responses, index)
		}(row, i+1)
	}

	// awaits
	wg.Wait()

	//Closes the responses channel after all responses are available

	close(responses)

	// Creates a map to store response per row 
	responseMap := make(map[int]BatchSendMessageResponse)

	// Gathers responses in the map
	for response := range responses {
		responseMap[len(responseMap)+1] = response
	}

	// Converts the responses map to JSON
	responseJSON, err := json.Marshal(responseMap)
	if err != nil {
		fmt.Println("Error al convertir las respuestas en JSON:", err)
		return err
	}

	// Saves JSON in local file
	responsePath := "./responses.json"
	err = ioutil.WriteFile(responsePath, responseJSON, 0644)
	if err != nil {
		fmt.Println("Error al guardar el archivo de respuestas JSON:", err)
		return err
	}

	fmt.Println("Respuestas guardadas en:", responsePath)
	return nil
}

func main() {
	// URL from GIthub for testing
	csvURL := "https://github.com/gabrielolivos/campaigncsv/raw/main/campaign_csv.csv"

	localPath := "./campaign_csv.csv"

	// 
	apiKey := "<API Key>"

	// Calls Function BatchSendCampaign
	if err := BatchSendCampaign(csvURL, localPath, apiKey); err != nil {
		fmt.Println("Error al enviar la campaña:", err)
		return
	}
}
