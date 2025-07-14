package main

import (
	"fmt"
	"sync"
	"time"
)

// 1. A fonte de trabalho: gera números e os envia para um canal.
func producer(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out) // Fecha o canal quando todos os números forem enviados.
	}()
	return out
}

// 2. O worker: recebe um número, calcula seu quadrado e envia para um canal de saída.
// Esta é a parte do FAN-OUT. Vários workers serão criados.
func squareWorker(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			// Simulando um trabalho pesado
			time.Sleep(100 * time.Millisecond)
			out <- n * n
		}
		close(out) // Fecha o canal de saída deste worker.
	}()
	return out
}

// 3. O consolidador: recebe múltiplos canais de entrada e os une em um único canal de saída.
// Esta é a parte do FAN-IN.
func fanIn(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Função para transferir dados de um canal de entrada para o de saída.
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go output(c)
	}

	// Inicia uma goroutine para fechar o canal 'out' APÓS
	// que todas as goroutines de 'output' tenham terminado.
	// Isso é crucial para evitar deadlocks.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	// Etapa 1: Gerar nossos dados de entrada.
	inputChan := producer(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	// Etapa 2: FAN-OUT
	// Distribuir o trabalho do 'inputChan' para 3 workers.
	numWorkers := 3
	workerChannels := make([]<-chan int, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workerChannels[i] = squareWorker(inputChan)
	}

	// Etapa 3: FAN-IN
	// Consolidar os resultados dos 3 workers em um único canal.
	finalResults := fanIn(workerChannels...)

	// Etapa 4: Processar os resultados finais.
	// O loop termina quando o canal 'finalResults' é fechado pela função fanIn.
	for result := range finalResults {
		fmt.Println("Resultado recebido:", result)
	}

	fmt.Println("Processamento concluído.")
}
