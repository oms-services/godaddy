package domain

import (
	"context"
	"encoding/json"
	"fmt"
	result "github.com/heaptracetechnology/godaddy/result"
	godaddy "github.com/kryptoslogic/godaddy-domainclient"
	"net/http"
	"os"
)

//Arguments struct
type Arguments struct {
	Domain string `json:"domain,omitempty"`
}

//Response struct
type Response struct {
	Available  bool   `json:"available"`
	Currency   string `json:"currency"`
	Definitive bool   `json:"definitive"`
	Domain     string `json:"domain"`
	Period     int32  `json:"period"`
	Price      int32  `json:"price"`
}

//CheckDomainAvailability GoDaddy
func CheckDomainAvailability(responseWriter http.ResponseWriter, request *http.Request) {

	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")

	decoder := json.NewDecoder(request.Body)
	var param Arguments
	decodeErr := decoder.Decode(&param)
	if decodeErr != nil {
		result.WriteErrorResponseString(responseWriter, decodeErr.Error())
		return
	}

	var apiConfig = godaddy.NewConfiguration()

	apiConfig.BasePath = "https://api.godaddy.com/"

	var authString = fmt.Sprintf("sso-key %s:%s", apiKey, apiSecret)
	apiConfig.AddDefaultHeader("Authorization", authString)

	var apiClient = godaddy.NewAPIClient(apiConfig)

	var ctx context.Context
	DomainAvailableResponse, _, err := apiClient.V1Api.Available(ctx, param.Domain, nil)
	if err != nil {
		result.WriteErrorResponse(responseWriter, err)
		return
	}

	var response Response
	if !DomainAvailableResponse.Available {
		response = Response{
			Available:  DomainAvailableResponse.Available,
			Currency:   "NA",
			Definitive: DomainAvailableResponse.Definitive,
			Domain:     DomainAvailableResponse.Domain,
			Period:     0,
			Price:      0,
		}
	} else {
		response = Response{
			Available:  DomainAvailableResponse.Available,
			Currency:   DomainAvailableResponse.Currency,
			Definitive: DomainAvailableResponse.Definitive,
			Domain:     DomainAvailableResponse.Domain,
			Period:     DomainAvailableResponse.Period,
			Price:      DomainAvailableResponse.Price,
		}
	}

	bytes, _ := json.Marshal(response)
	result.WriteJsonResponse(responseWriter, bytes, http.StatusOK)
}
