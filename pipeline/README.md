# Pipeline

Chegamos ao clímax da nossa jornada pelos padrões de concorrência. Os Pipelines são a aplicação mais elegante e poderosa de tudo o que vimos até agora. Eles unem os conceitos de Worker Pools e Fan-In/Fan-Out para criar fluxos de processamento de dados incrivelmente eficientes e modulares.

## O que é um Pipeline?

Pense em uma linha de montagem industrial moderna.

1. Estação 1: Um braço robótico pega uma peça de metal bruta (os dados de entrada).
2. Estação 2: Outra máquina corta e dobra a peça.
3. Estação 3: Um terceiro robô pinta a peça já moldada.
4. Estação 4: A peça final é embalada.

Cada estação é um **estágio** do nosso pipeline. A esteira que conecta as estações é o nosso **canal** em Go. A beleza disso é que todas as estações podem trabalhar ao mesmo tempo. Enquanto a Estação 3 está pintando a peça A, a Estação 2 já está moldando a peça B e a Estação 1 já está pegando a peça C.

Um **Pipeline em Go** é exatamente isso: uma série de estágios de processamento conectados por canais, onde a saída de um estágio se torna a entrada do próximo.

O mais interessante é que cada "estação" (estágio) pode ter múltiplos "robôs" (goroutines) trabalhando em paralelo. E como conectamos múltiplos robôs a uma única esteira? Com os padrões Fan-Out e Fan-In que acabamos de aprender!

## Exemplo Prático: Pipeline de Múltiplos Estágios

Vamos evoluir nosso exemplo anterior. Agora, nosso pipeline terá dois estágios de trabalho:

1. Estágio 1: Recebe números e calcula o quadrado deles (como antes).
2 Estágio 2: Recebe os resultados do primeiro estágio (os quadrados) e adiciona 1 a cada um.

Cada estágio terá seu próprio conjunto de workers.

Veja o arquivo main.go

## Análise do Código

A `main` agora é a nossa "planta da fábrica", onde montamos a linha de produção:

1. `producer`: Continua sendo o ponto de partida, colocando a matéria-prima (`int`s) na primeira esteira (`inputChan`).

2. Estágio 1 (Squarer):

  - Fan-Out: Criamos 3 squareWorkers. Todos eles leem da mesma esteira de entrada (inputChan).
  - Fan-In: Usamos a função fanIn para coletar os resultados de todos os squareWorkers em uma nova esteira, a squareResultsChan.

3. Estágio 2 (Incrementer):

  - Fan-Out: O canal de saída do estágio anterior (squareResultsChan) se torna o canal de entrada para este. Criamos 2 incrementWorkers que leem dele.
  - Fan-In: Novamente, usamos fanIn para coletar os resultados dos incrementWorkers na esteira final, a finalResultsChan.

4. Consumidor: O for range no final da main é o "controle de qualidade", que pega os produtos acabados da última esteira e os exibe.

## Por que Pipelines são tão poderosos?

1. Composição e Reutilização: Cada estágio (squareWorker, incrementWorker) é uma função simples e independente. Você pode facilmente testá-los isoladamente e combiná-los de diferentes maneiras para criar novos pipelines, como peças de Lego.

2. Controle de Concorrência Granular: Perceba que definimos um número de workers diferente para cada estágio (3 squarers e 2 incrementers). Isso é fantástico! Se o "quadrado" for uma operação muito mais lenta que o "incremento", você pode alocar mais goroutines para o primeiro estágio para evitar que ele se torne um gargalo, otimizando o uso de recursos.

3. Backpressure Natural: Este é um benefício sutil e incrível. Imagine que o Estágio 2 (incremento) se torne muito lento. O canal squareResultsChan começará a encher, pois os incrementWorkers não estão consumindo rápido o suficiente. Quando o buffer do canal (se houver) encher, os squareWorkers do Estágio 1 tentarão enviar um novo resultado (out <- n * n) e ficarão bloqueados. Isso, por sua vez, fará com que eles parem de consumir do inputChan. O efeito se propaga para trás, desacelerando toda a linha de produção automaticamente, sem que você precise escrever uma linha de código para isso. Isso previne que estágios mais rápidos sobrecarreguem os mais lentos.