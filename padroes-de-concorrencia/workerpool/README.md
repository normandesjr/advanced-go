# Advanced Go

A ideia desse repo é agrupar diversos conceitos avançados de Go.

## Worker Pool

É um conceito fundamental em concorrência, importante para aplicações robustas e performáticas.

### O que é um Worker Pool?

Imagine que você tem uma montanha de tarefas para fazer, como carimbar 1000 envelopes. Você poderia fazer tudo sozinho, um por um (sequencial). Ou, você poderia contratar 5 amigos (os workers) para te ajudar.

Para se organizar, você cria duas pilhas de bandejas:

1. **Bandeja de "Tarefas a Fazer"**: Onde você coloca todos os envelopes que precisam ser carimbados.
2. **Bandeja de "Tarefas Concluídas"**: Onde seus amigos colocam os envelopes já carimbados.

Seus amigos (workers) pegam um envelope da primeira bandeja, carimbam, e colocam na segunda. Eles repetem isso até não haver mais envelopes na bandeja de "Tarefas a Fazer".

O Worker Pool em Go é exatamente isso:

- **Tarefas**: Os dados ou trabalhos que precisam ser processados.
- **Workers**: Goroutines que ficam prontas para executar uma tarefa.
- **Bandeja de "Tarefas a Fazer"**: Um canal (chan) em Go que distribui as tarefas.
- **Bandeja de "Tarefas Concluídas"**: Outro canal (chan) que coleta os resultados.

### Por que usar Worker Pools?

O principal motivo é controlar a concorrência.

Embora seja muito fácil criar goroutines em Go (go minhaFuncao()), criar goroutines sem controle pode ser perigoso. Imagine ter que processar 1 milhão de itens. Se você criar 1 milhão de goroutines de uma vez, você pode esgotar a memória e os recursos do seu sistema, causando um crash.

O Worker Pool limita o número de goroutines ativas a um valor que você define (o tamanho da "piscina"), garantindo que sua aplicação use os recursos de forma eficiente e previsível, sem sobrecarregar o sistema ou serviços externos (como um banco de dados ou uma API).

## Exemplo Prático: Verificador de Status de Websites

Vamos criar um programa que recebe uma lista de URLs e verifica se elas estão online. Este é um exemplo clássico de tarefa I/O-bound (limitada pela rede), onde os workers passarão a maior parte do tempo esperando uma resposta HTTP.

Veja o arquivo main.go

## Análise do Código

1. Iniciando os Workers (A contratação dos amigos)

```go
numWorkers := 3
for w := 1; w <= numWorkers; w++ {
    wg.Add(1)
    go worker(w, &wg, jobs, results)
}
```

Aqui, criamos um pool com 3 workers. Cada worker é uma goroutine que fica pronta para receber trabalho do canal jobs. Usamos um sync.WaitGroup para garantir que o programa principal não termine antes que todos os workers tenham finalizado.

Preste bastante atenção nessa parte, os 3 workers vão iniciar e parar lá no "for job := range jobs" porque é um chan que precisa ser preenchido com algo antes de iniciar. Só vai sair desse for quando o "chan jobs" for fechado com "close(jobs)".

2. Enviando os Jobs (Colocando os envelopes na bandeja)

```go
for i, url := range urls {
    jobs <- Job{ID: i + 1, URL: url}
}
close(jobs)
```

Iteramos sobre nossa lista de URLs e enviamos cada uma como um Job para o canal jobs. A linha close(jobs) é talvez a mais importante aqui. Ela "fecha a bandeja de entrada" e sinaliza para as goroutines que estão lendo (range jobs) que não virão mais tarefas. Sem isso, os workers ficariam esperando para sempre, causando um deadlock.

3. Coletando os Resultados (Verificando a bandeja de concluídos)

```go
for i := 0; i < numJobs; i++ {
    result := <-results
    // ... imprime o resultado
}
```

Sabemos exatamente quantos resultados esperar (o mesmo número de jobs enviados). Então, fazemos um loop para ler cada resultado do canal results e o imprimimos na tela.

## Pontos-chave e Boas Práticas

- Controle de Concorrência: Repare que, mesmo com 8 URLs, apenas 3 requisições HTTP são feitas simultaneamente, pois temos apenas 3 workers. Isso evita sobrecarregar a rede ou a API de destino.
- Canais e Comunicação: A comunicação entre a goroutine principal (o "chefe") e os workers é feita de forma segura através de canais, evitando condições de corrida (race conditions).
- Graceful Shutdown (Encerramento Suave): O uso de sync.WaitGroup junto com o close(jobs) garante que o programa espere todos os workers terminarem seu trabalho atual antes de encerrar. Isso é fundamental para não deixar tarefas pela metade.
- Tamanho do Pool: Como definir numWorkers?
  - Tarefas CPU-Bound (cálculos intensos): Um bom ponto de partida é runtime.NumCPU(), o número de núcleos do seu processador.
  - Tarefas I/O-Bound (rede, disco, banco de dados): O número pode ser bem maior que o de CPUs, pois os workers passam a maior parte do tempo esperando. O valor ideal geralmente é encontrado através de testes e medições (benchmarking).