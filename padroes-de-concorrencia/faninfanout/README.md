# Fan-in/Fan-out

O padrão Fan-In / Fan-Out é um parente muito próximo do Worker Pool. Na verdade, você verá que o Worker Pool é uma aplicação prática desses dois conceitos juntos.

## O que são Fan-Out e Fan-In?

Imagine que você é o gerente de uma fábrica de brinquedos.

1. Fan-Out (Distribuir/Espalhar): Você tem um grande caixote com 1000 blocos de madeira (as tarefas). Em vez de dar a caixa inteira para uma única pessoa, você a abre e distribui os blocos para 10 artesãos diferentes. Cada artesão pega um bloco e começa a trabalhar. Isso é Fan-Out: uma única fonte de trabalho é distribuída para múltiplos processadores (goroutines) que podem trabalhar em paralelo.

2. Fan-In (Consolidar/Juntar): Cada um dos 10 artesãos, ao terminar de esculpir seu bloco de madeira, o coloca em uma única esteira rolante que leva para o setor de pintura. Isso é Fan-In: os resultados de múltiplos processadores (goroutines) são coletados em um único canal, como um funil.

Resumindo:

- Fan-Out: Um canal -> Múltiplas goroutines lendo dele.
- Fan-In: Múltiplas goroutines escrevendo -> Um único canal recebendo tudo.

### Conexão com o Worker Pool

Percebeu a semelhança? O nosso exemplo de Worker Pool já usava esses padrões!

- Fan-Out: A main goroutine enviava todas as URLs para o canal jobs. Os 3 workers liam desse único canal, distribuindo o trabalho entre si.
- Fan-In: Todos os 3 workers, após verificarem a URL, enviavam seus resultados para o único canal results.

Agora, vamos criar um exemplo mais explícito para isolar e entender melhor cada parte.

## Exemplo Prático: Pipeline de Processamento de Números

Vamos criar um programa que:

1. Gera uma sequência de números (nossa fonte de trabalho).
2. Usa múltiplos workers (Fan-Out) para calcular o quadrado de cada número (uma tarefa "pesada").
3. Usa uma função de Fan-In para juntar todos os resultados em um único fluxo para serem impressos.

Este padrão é muito comum para criar pipelines de processamento de dados em Go.

Veja o arquivo main.go.

### Análise do Código

Vamos analisar o fluxo:

1. **producer**: É o nosso ponto de partida. Ele simplesmente cria uma goroutine para popular um canal com dados e depois o fecha. Ele retorna um canal somente de leitura (<-chan int), uma boa prática para evitar que quem o recebe escreva nele acidentalmente.

2. **squareWorker** (Fan-Out):

```go
// Na main...
for i := 0; i < numWorkers; i++ {
    workerChannels[i] = squareWorker(inputChan)
}
```

Aqui está o Fan-Out em ação. Criamos 3 squareWorker. É importante notar que todos eles estão lendo do mesmo inputChan. O Go se encarrega de distribuir os valores do canal entre as goroutines que estão prontas para ler. Quando um worker pega o número 3, os outros não o pegam. Cada worker retorna seu próprio canal de resultados.

3. **fanIn** (Fan-In):

```go
// Na main...
finalResults := fanIn(workerChannels...)
```

Esta é a função mais interessante. Ela recebe um slice de canais (`...<-chan int`) e sua missão é criar um funil.

- Ela cria um WaitGroup para rastrear quando cada canal de entrada foi completamente lido.
- Para cada canal de entrada, ela inicia uma goroutine (output) que simplesmente lê os dados e os repassa para o canal de saída out.
- **O truque genial está aqui**: ela inicia uma goroutine final que apenas espera (wg.Wait()) por todas as outras terminarem para então fechar o canal out. Isso garante que o consumidor final (o for range na main) saiba quando parar de esperar por resultados.

## Por que separar Fan-In/Fan-Out?

Embora o padrão Worker Pool já use esses conceitos, entendê-los separadamente permite que você construa pipelines de processamento mais complexos e flexíveis.

Imagine um pipeline mais longo: producer -> (Fan-Out) -> validadores -> (Fan-In) -> processadores -> (Fan-Out) -> salvadoresDeBD -> (Fan-In) -> notificadores

Cada estágio do pipeline pode ter seu próprio número de workers, otimizado para a tarefa específica que ele executa. Fan-In e Fan-Out são os "canos" e "juntas" que conectam esses estágios de forma concorrente e segura.