package main

import (
	"bytes"
	"encoding/json"
	"math"
	"strconv"
)

type PaymentsSummary struct {
	PaymentsDefault  Summary `json:"default"`
	PaymentsFallback Summary `json:"fallback"`
}

type Summary struct {
	TotalRequests uint    `json:"totalRequests"`
	TotalAmount   float64 `json:"totalAmount"`
}

func (s *Summary) MarshalJSON() ([]byte, error) {
	resp := struct {
		TotalRequests uint    `json:"totalRequests"`
		TotalAmount   float64 `json:"totalAmount"`
	}{
		TotalRequests: s.TotalRequests,
		TotalAmount:   math.Round(s.TotalAmount*10) / 10,
	}

	return json.Marshal(resp)

}

func (s *Summary) MarshalJSONManual() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(`{"totalRequests":`)
	buf.WriteString(strconv.FormatUint(uint64(s.TotalRequests), 10))
	buf.WriteString(`,"totalAmount":`)
	buf.WriteString(strconv.FormatFloat(math.Round(s.TotalAmount*10)/10, 'f', 1, 64))
	buf.WriteString(`}`)
	return buf.Bytes(), nil
}
