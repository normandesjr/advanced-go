# Primitivas de Sincronização

Prompt: *quero aprender sobre Primitivas de Sincronização, quero dominar o uso de sync.Mutex, sync.RWMutex, sync.WaitGroup, sync.Once e as operações atômicas do pacote sync/atomic. O que recomenda, aprender separado ou daria pra juntar me um único exemplo?*

Estamos mergulhando agora no coração do controle de concorrência em Go. As primitivas de sincronização são as ferramentas de mais baixo nível que você usará para orquestrar suas goroutines e garantir que elas não "tropecem" umas nas outras.

Respondendo à sua pergunta: "aprender separado ou juntar em um único exemplo?"

Minha recomendação é: **vamos juntar tudo em um único exemplo coeso**.

Criar exemplos isolados para cada primitiva é útil, mas pode parecer abstrato. Ao construir um único programa que precisa de cada uma dessas ferramentas para resolver um problema específico, o propósito e a diferença entre elas se tornam muito mais claros. Você verá exatamente por que escolheríamos um RWMutex em vez de um Mutex ou quando uma operação atômica é a melhor opção.

Vamos construir um **Coletor de Métricas Concorrente**. Imagine um serviço web que recebe milhares de requisições por segundo e precisa registrar, em tempo real e de forma segura, quantas vezes cada página foi visitada.

## O Cenário: Coletor de Métricas

Nosso coletor terá as seguintes responsabilidades:

1. Contar as visitas para cada URL em um mapa (map[string]int).
2. Manter um contador global do total de requisições processadas.
3. Garantir que uma "conexão com o banco de dados" (que vamos simular) seja inicializada apenas uma vez.
4. O programa principal deve esperar que todas as requisições simuladas terminem antes de exibir os resultados.

Cada uma dessas responsabilidades é um pretexto perfeito para usar uma primitiva de sincronização diferente.

Vamos criar o arquivo main.go e construir nossa solução passo a passo, introduzindo cada primitiva.

Veja o main.go

## Análise e Explicação de Cada Primitiva

Vamos dissecar o código e o papel de cada ferramenta:

1. `sync.RWMutex` (O Porteiro Flexível)

* O que é? Um Read-Write Mutex (bloqueio de leitura-escrita). Ele é mais inteligente que um Mutex normal.

* Analogia: Pense em uma sala de biblioteca.
  - Muitas pessoas podem entrar para ler livros ao mesmo tempo (RLock).
  - Mas se alguém precisa escrever (modificar o catálogo), essa pessoa tranca a porta (Lock), e ninguém mais pode entrar (nem para ler, nem para escrever) até que ela termine.

* No nosso código:
  - Increment() precisa modificar o mapa visits, então ele usa m.mu.Lock() para obter acesso exclusivo.
  - GetVisits() só precisa ler o mapa. Ele usa m.mu.RLock(), permitindo que centenas de goroutines leiam os dados de visita simultaneamente, desde que nenhuma esteja escrevendo. Isso melhora muito a performance em cenários com mais leituras do que escritas.

2. sync/atomic (O Contador Especialista)

* O que é? Um pacote que fornece operações de baixo nível, à prova de concorrência, para tipos numéricos. Essas operações são executadas como uma única instrução indivisível pelo processador.
* Analogia: Um contador de "click" manual. Você pode apertar o botão para incrementar o número sem precisar de um ritual complexo de trancar e destrancar uma caixa. É uma ação única e segura.
* No nosso código:
  - atomic.AddUint64(&m.totalRequests, 1): Adiciona 1 ao nosso contador global. É muito mais rápido e eficiente do que usar um Mutex apenas para incrementar um número.
  - atomic.LoadUint64(&m.totalRequests): Lê o valor atual do contador de forma segura.

3. sync.Once (O Ritual de Iniciação)

* O que é? Um objeto que garante que uma determinada função será executada exatamente uma vez, não importa quantas goroutines tentem executá-la.
* Analogia: A cerimônia de inauguração de uma ponte. A fita é cortada apenas uma vez, na primeira vez que alguém tenta atravessar. Todas as tentativas subsequentes simplesmente usam a ponte já inaugurada.
* No nosso código:
  - m.once.Do(m.initializeDB): Na primeira vez que Increment é chamado por qualquer goroutine, initializeDB será executado. Em todas as 1000 chamadas seguintes, o Do simplesmente retornará imediatamente, sem executar a função novamente. Perfeito para inicializações caras que só precisam acontecer uma vez.

4. sync.WaitGroup (O Coordenador da Chegada)

* O que é? Um contador para esperar que uma coleção de goroutines termine sua execução.
* Analogia: O gerente de uma equipe de entregadores. Ele sabe que tem 10 entregas para fazer (wg.Add(10)). Cada vez que um entregador volta, ele avisa (wg.Done()). O gerente só vai para casa (wg.Wait() retorna) quando todos os 10 entregadores tiverem retornado.
* No nosso código:
  - wg.Add(numRequests): Informamos ao WaitGroup que estamos esperando por 1000 goroutines.
  - defer wg.Done(): No início de cada goroutine, garantimos que, ao sair, ela decrementará o contador.
  - wg.Wait(): A main goroutine para aqui e só continua quando o contador do WaitGroup chegar a zero.

## Resumo: Quando Usar o Quê?

| Primitiva | Caso de Uso Principal | Analogia |
|:---|:---|:---|
| `sync.Mutex` | Proteger uma seção crítica de código (dados complexos) contra acesso concorrente (leitura e escrita). | Uma chave de banheiro (só um por vez). |
| `sync.RWMutex` | Proteger dados que são lidos com muito mais frequência do que são escritos. | Uma biblioteca (muitos leitores, um escritor). |
| `sync/atomic` | Modificar/ler tipos numéricos simples (contadores, flags) de forma super eficiente. | Um contador de "click". |
| `sync.Once` | Garantir que uma inicialização (ex: config, conexão de BD) aconteça apenas uma vez. | Cerimônia de inauguração. |
| `sync.WaitGroup` | Esperar que um número conhecido de goroutines termine seu trabalho. | Gerente esperando a equipe voltar. |

Dominar essas cinco primitivas lhe dará um controle imenso sobre o comportamento concorrente de suas aplicações. 