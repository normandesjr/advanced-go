package main

import (
	"fmt"
	"sync"
)

// --- Funções Base (as mesmas do exemplo Fan-In/Fan-Out) ---

// producer (Estágio 0): Gera os dados iniciais.
func producer(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// fanIn: Consolida múltiplos canais em um só.
func fanIn(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

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

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// --- Estágios do Pipeline ---

// squareWorker (Worker do Estágio 1): Calcula o quadrado de um número.
func squareWorker(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			// fmt.Printf("Squarer recebendo: %d\n", n)
			out <- n * n
		}
		close(out)
	}()
	return out
}

// incrementWorker (Worker do Estágio 2): Adiciona 1 a um número.
func incrementWorker(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			// fmt.Printf("Incrementer recebendo: %d\n", n)
			out <- n + 1
		}
		close(out)
	}()
	return out
}

func main() {
	// Dados de entrada para o pipeline.
	inputData := []int{1, 2, 3, 4, 5}

	// 1. Iniciar o produtor (o início da linha de montagem).
	inputChan := producer(inputData...)

	// 2. Configurar o ESTÁGIO 1: Squarers
	// Fan-Out: distribuímos o trabalho do inputChan para 3 workers.
	numSquareWorkers := 3
	squareWorkerChans := make([]<-chan int, numSquareWorkers)
	for i := 0; i < numSquareWorkers; i++ {
		squareWorkerChans[i] = squareWorker(inputChan)
	}
	// Fan-In: consolidamos os resultados dos squarers em um único canal.
	// Este canal será a entrada para o próximo estágio.
	squareResultsChan := fanIn(squareWorkerChans...)

	// 3. Configurar o ESTÁGIO 2: Incrementers
	// Fan-Out: distribuímos o trabalho do `squareResultsChan` para 2 workers.
	numIncrementWorkers := 2
	incrementWorkerChans := make([]<-chan int, numIncrementWorkers)
	for i := 0; i < numIncrementWorkers; i++ {
		incrementWorkerChans[i] = incrementWorker(squareResultsChan)
	}
	// Fan-In: consolidamos os resultados finais.
	finalResultsChan := fanIn(incrementWorkerChans...)

	// 4. Consumir os resultados finais do pipeline.
	// O loop termina quando `finalResultsChan` é fechado.
	for result := range finalResultsChan {
		// Exemplo de resultado:
		// n = 3 -> square = 9 -> increment = 10
		fmt.Println("Resultado final do pipeline:", result)
	}

	fmt.Println("Pipeline concluído.")
}
