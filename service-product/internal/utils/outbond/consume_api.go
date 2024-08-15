package outbond

import (
	"encoding/json"
	"fmt"
	"net/http"
	"service-product/internal/domain/dto"
	"time"
)

func GetProductValidation(productID string) (*dto.BaseResponse, error) {
	url := fmt.Sprintf("https://wk3j1.wiremockapi.cloud/product/validate?id=%s", productID)

	// Buat HTTP client dengan timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Buat permintaan HTTP GET
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Kirim permintaan
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Periksa status kode respons
	//if resp.StatusCode != http.StatusOK {
	//	return nil, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	//}

	// Dekode respons JSON
	var validationResponse dto.BaseResponse
	err = json.NewDecoder(resp.Body).Decode(&validationResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode HTTP response: %v", err)
	}

	return &validationResponse, nil
}
