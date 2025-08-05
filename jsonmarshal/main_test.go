package main

import (
	"encoding/json"
	"testing"
)

func BenchmarkJsonMarshal_Standard(b *testing.B) {
	s := &Summary{
		TotalRequests: 12345,
		TotalAmount:   9876.54321,
	}

	b.ResetTimer()
	for b.Loop() {
		_, err := json.Marshal(s)
		if err != nil {
			b.Fatalf("json.Marshal failed: %v", err)
		}
	}
}

func BenchmarkJsonMarshal_Manual(b *testing.B) {
	s := &Summary{
		TotalRequests: 12345,
		TotalAmount:   9876.54321,
	}

	b.ResetTimer()
	for b.Loop() {
		s.MarshalJSONManual()
	}
}
