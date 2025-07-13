package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Job representa a tarefa a ser feita, neste caso, uma URL para verificar.
type Job struct {
	ID  int
	URL string
}

// Result armazena o resultado de um Job processado.
type Result struct {
	Job    Job
	Status string
	Error  error
}

// worker é a nossa goroutine "trabalhadora".
// Ele recebe jobs do canal `jobs` e envia os resultados para o canal `results`.
// Usamos um WaitGroup para sinalizar que o worker terminou seu trabalho.
func worker(id int, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
	// Ao final da função, sinalizamos ao WaitGroup que este worker concluiu.
	defer wg.Done()

	// O worker fica em um loop, esperando por jobs.
	// Este loop termina quando o canal `jobs` é fechado e todos os valores foram recebidos.
	for job := range jobs {
		fmt.Printf("Worker %d iniciando job %d para a URL %s\n", id, job.ID, job.URL)

		// Criamos um cliente HTTP com um timeout para não ficar preso eternamente.
		client := http.Client{
			Timeout: 5 * time.Second,
		}

		resp, err := client.Get(job.URL)

		// Preparamos o resultado.
		result := Result{Job: job}
		if err != nil {
			result.Error = err
			result.Status = "ERRO"
		} else {
			// Lembre-se de fechar o corpo da resposta para liberar a conexão.
			resp.Body.Close()
			result.Status = resp.Status
		}

		// Enviamos o resultado para o canal de resultados.
		results <- result
	}
}

func main() {
	urls := []string{
		"https://golang.org",
		"https://google.com",
		"https://github.com",
		"http://uma-url-que-nao-existe.dev",
		"https://api.github.com",
		"https://www.uol.com.br",
		"https://youtube.com",
		"http://localhost:9999", // Outra URL que provavelmente falhará
	}

	// Criamos os canais para jobs e results.
	// O canal de jobs é bufferizado para que possamos enfileirar todos os jobs
	// sem esperar que um worker esteja imediatamente disponível.
	numJobs := len(urls)
	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	// Para esperar que todos os workers terminem, usamos um WaitGroup.
	var wg sync.WaitGroup

	// 1. Iniciando os Workers
	numWorkers := 3
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1) // Adicionamos 1 ao contador do WaitGroup para cada worker.
		go worker(w, &wg, jobs, results)
	}

	// 2. Enviando os Jobs
	for i, url := range urls {
		jobs <- Job{ID: i + 1, URL: url}
	}
	close(jobs) // Fechamos o canal de jobs. Isso é CRUCIAL!
	// É o sinal para os workers de que não há mais trabalho.

	// 3. Coletando os Resultados
	// Esperamos que todos os jobs sejam processados para imprimir os resultados.
	for i := 0; i < numJobs; i++ {
		result := <-results
		if result.Error != nil {
			fmt.Printf("Resultado Job %d (%s): %s - Erro: %v\n", result.Job.ID, result.Job.URL, result.Status, result.Error)
		} else {
			fmt.Printf("Resultado Job %d (%s): %s\n", result.Job.ID, result.Job.URL, result.Status)
		}
	}

	// Esperamos que todos os workers finalizem suas goroutines antes de encerrar o main.
	wg.Wait()
	fmt.Println("Todos os workers terminaram. Programa finalizado.")
}
