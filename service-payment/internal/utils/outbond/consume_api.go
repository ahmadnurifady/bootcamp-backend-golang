package outbond

import (
	"encoding/json"
	"fmt"
	"net/http"
	"service-payment/internal/domain/dto"
	"time"
)

func GetPaymentValidation(userName string) (*dto.BaseResponse, error) {
	url := fmt.Sprintf("https://wk3j1.wiremockapi.cloud/payment/validate?userId=%s", userName)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	var validationResponse dto.BaseResponse
	err = json.NewDecoder(resp.Body).Decode(&validationResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode HTTP response: %v", err)
	}

	return &validationResponse, nil
}

func GetPaymentValidationCredit(numberPhone string) (*dto.BaseResponse, error) {
	url := fmt.Sprintf("https://wk3j1.wiremockapi.cloud/validate/payment/kopay?numberPhone=%s", numberPhone)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	var validationResponse dto.BaseResponse
	err = json.NewDecoder(resp.Body).Decode(&validationResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode HTTP response: %v", err)
	}

	return &validationResponse, nil
}
