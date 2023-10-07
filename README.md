# Connectly SDK for Go

This Go SDK allows you to easily send campaigns using the Connectly API.

## Installation

To install the Connectly SDK, you can use the following `go get` command:

```shell

go get github.com/gabrielolivos/connectly-sdk-go

```

This will download the SDK to your local environment

## Usage

Here's an example of how to use the SDK to send a campaign from a CSV file:

```

// Import the SDK
import "github.com/yourusername/connectly-sdk-go"

func main() {
    // Set your API key
    apiKey := "your-api-key"

    // URL of the CSV file
    csvURL := "https://example.com/path/to/your/campaign.csv"

    // Local path to save the CSV file
    localPath := "./campaign.csv"

    // Send the campaign using the SDK
    err := connectly.BatchSendCampaign(csvURL, localPath, apiKey)
    if err != nil {
        fmt.Println("Error sending the campaign:", err)
        return
    }

    fmt.Println("Campaign sent successfully!")
}

```

## Running the SDK

To run the SDK, simply use command ```go run .``` from directory where your main.go file is located. 

## CSV format 

This SDK was created to read CSV files with the following structure:

channel_type,external_id,{template_name}:body_1,{template_name}:body_2

For custom CSV format integration, please reach out to your technical contact at Connectly! 


