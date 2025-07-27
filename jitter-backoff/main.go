package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	const maxRetries = 5
	// Define o tempo de espera inicial. Começaremos esperando 100ms após a primeira falha.
	baseDelay := 100 * time.Millisecond

	// Inicializa o gerador de números aleatórios. Em aplicações reais, isso deve ser feito apenas uma vez.
	// rand.Seed(time.Now().UnixNano()) // Deprecated em Go 1.20+

	fmt.Println("Simulando uma operação que falha e precisa de retries...")

	for attempt := 0; attempt < maxRetries; attempt++ {
		fmt.Printf("\n--- Tentativa %d ---\n", attempt+1)

		// Simula a falha da operação
		fmt.Println("Operação falhou. Calculando tempo de espera...")

		// Calcula o backoff exponencial
		backoff := baseDelay * time.Duration(math.Pow(2, float64(attempt)))

		// Adiciona um jitter (aleatoriedade) para evitar o problema da "manada"
		// Aqui, um jitter de até 100ms
		jitter := time.Duration(rand.Intn(100)) * time.Millisecond

		totalWait := backoff + jitter

		fmt.Printf("Backoff base: %v\nJitter: %v\nTempo total de espera: %v\n", backoff, jitter, totalWait)

		time.Sleep(totalWait)
	}

	fmt.Printf("\n--- Fim das %d tentativas ---\n", maxRetries)
}
