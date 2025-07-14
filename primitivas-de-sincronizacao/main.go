package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// WebsiteMetrics é a nossa estrutura que armazena dados compartilhados.
// Ela é segura para uso concorrente.
type WebsiteMetrics struct {
	// 1. Para proteger o mapa 'visits' de escritas e leituras concorrentes.
	// Usamos RWMutex porque teremos muitas leituras (para exibir dashboards)
	// e escritas (quando uma nova visita chega).
	mu     sync.RWMutex
	visits map[string]int

	// 2. Para contar o total de requisições de forma super eficiente.
	// Operações atômicas são mais rápidas que um Mutex para contadores simples.
	totalRequests uint64

	// 3. Para garantir que a inicialização do "banco de dados" ocorra apenas uma vez.
	once         sync.Once
	dbConnection string
}

// NewWebsiteMetrics cria uma nova instância de WebsiteMetrics.
func NewWebsiteMetrics() *WebsiteMetrics {
	return &WebsiteMetrics{
		visits: make(map[string]int),
	}
}

// initializeDB simula a inicialização de um recurso caro, como uma conexão de BD.
func (m *WebsiteMetrics) initializeDB() {
	fmt.Println("--- (Inicializando conexão com o banco de dados...) ---")
	time.Sleep(100 * time.Millisecond) // Simula trabalho
	m.dbConnection = "CONECTADO"
	fmt.Println("--- (Conexão com o banco de dados ESTABELECIDA) ---")
}

// Increment registra a visita a uma URL. Esta é a principal função de escrita.
func (m *WebsiteMetrics) Increment(url string) {
	// 3. sync.Once: Garante que initializeDB() seja chamada apenas uma vez,
	// não importa quantas goroutines chamem Increment() simultaneamente.
	m.once.Do(m.initializeDB)

	// 1. sync.Mutex: Bloqueamos para escrita exclusiva no mapa.
	// Nenhuma outra goroutine pode ler ou escrever no mapa 'visits' enquanto
	// esta seção estiver em execução.
	m.mu.Lock()
	m.visits[url]++
	m.mu.Unlock()

	// 2. sync/atomic: Incrementa o contador global.
	// Não precisa de lock, é uma operação segura e muito rápida.
	atomic.AddUint64(&m.totalRequests, 1)
}

// GetVisits retorna o número de visitas para uma URL específica.
// Esta é uma função de leitura.
func (m *WebsiteMetrics) GetVisits(url string) int {
	// 1. sync.RWMutex: Bloqueamos para leitura.
	// Múltiplas goroutines podem chamar GetVisits ao mesmo tempo, desde que
	// nenhuma goroutine esteja escrevendo (com m.mu.Lock()).
	m.mu.RLock()
	defer m.mu.RUnlock() // defer é perfeito para garantir que o unlock aconteça.
	return m.visits[url]
}

// TotalRequests retorna o contador global de requisições.
func (m *WebsiteMetrics) TotalRequests() uint64 {
	// 2. sync/atomic: Lê o valor do contador de forma segura.
	return atomic.LoadUint64(&m.totalRequests)
}

func main() {
	metrics := NewWebsiteMetrics()
	urls := []string{
		"/home",
		"/products",
		"/contact",
		"/about",
	}

	// 4. sync.WaitGroup: Para esperar que todas as nossas goroutines terminem.
	var wg sync.WaitGroup

	// Simulamos 1000 requisições concorrentes.
	numRequests := 1000
	wg.Add(numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done() // Sinaliza ao WaitGroup que esta goroutine terminou.

			// Escolhe uma URL aleatória para "visitar".
			randomURL := urls[rand.Intn(len(urls))]
			metrics.Increment(randomURL)
		}()
	}

	// Espera aqui até que o contador do WaitGroup chegue a zero.
	wg.Wait()

	fmt.Println("\n--- Métricas Finais ---")
	fmt.Printf("Conexão com o BD: %s\n", metrics.dbConnection)
	fmt.Printf("Total de requisições processadas: %d\n", metrics.TotalRequests())
	fmt.Println("Visitas por página:")
	for _, url := range urls {
		fmt.Printf("  - %s: %d visitas\n", url, metrics.GetVisits(url))
	}
}
