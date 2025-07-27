# Jitter e Backoff Exponencial

Quando interagimos com serviços externos (APIs, bancos de dados), falhas são inevitáveis. O serviço pode estar temporariamente sobrecarregado, a rede pode falhar, etc. Uma abordagem ingênua seria tentar novamente imediatamente em um loop, mas isso só pioraria a situação, criando um ataque de negação de serviço (DoS) contra o serviço que já está com problemas.

A estratégia correta é usar **Backoff Exponencial com Jitter**.

## O que é Backoff Exponencial?

É uma estratégia onde você aumenta exponencialmente o tempo de espera entre as tentativas.

- **Tentativa 1**: Falhou. Espere 100ms.
- **Tentativa 2**: Falhou de novo. Espere 200ms (100ms * 2).
- **Tentativa 3**: Falhou mais uma vez. Espere 400ms (100ms * 4).
- **Tentativa 4**: Falhou. Espere 800ms (100ms * 8).
- ... e assim por diante.

Isso dá ao serviço sobrecarregado um "fôlego" cada vez maior para se recuperar. A parte "exponencial" vem do fato de que o tempo de espera é multiplicado por uma base (geralmente 2) a cada tentativa.

## Por que o Jitter é Crucial?

Imagine que você tem 1000 instâncias da sua aplicação e todas elas encontram um erro ao mesmo tempo. Com o backoff exponencial simples, todas elas esperariam exatamente o mesmo tempo e tentariam novamente no exato mesmo instante, criando uma "manada" (thundering herd) que derrubaria o serviço novamente.

**Jitter** é a adição de uma pequena quantidade de tempo aleatório a cada espera.

- **Instância A**: Espera 400ms + 15ms (aleatório)
- **Instância B**: Espera 400ms + 78ms (aleatório)
- **Instância C**: Espera 400ms + 42ms (aleatório)

Isso "espalha" as novas tentativas no tempo, suavizando a carga no serviço e aumentando drasticamente a chance de recuperação do sistema como um todo.

## Exemplo Prático

O código em `main.go` demonstra exatamente essa estratégia. Ele simula uma operação que falha e, a cada tentativa, calcula um tempo de espera maior (backoff exponencial) e adiciona um tempo aleatório (jitter) antes de tentar novamente.

## Análise do Código

1.  **`math.Pow(2, float64(attempt))`**: Esta é a implementação do crescimento exponencial. `math.Pow(base, expoente)` calcula `base` elevado à potência do `expoente`.
    - Na tentativa 0, `2^0 = 1`. O atraso é `baseDelay * 1`.
    - Na tentativa 1, `2^1 = 2`. O atraso é `baseDelay * 2`.
    - Na tentativa 2, `2^2 = 4`. O atraso é `baseDelay * 4`.
    - E assim por diante, dobrando o tempo de espera a cada falha.

2.  **`rand.Intn(100)`**: Esta linha gera o Jitter. Ela cria um número aleatório entre 0 e 99, que é então convertido para uma `time.Duration` em milissegundos.

3.  **`time.Sleep(backoff + jitter)`**: A goroutine é pausada pela soma do tempo de backoff calculado mais o jitter aleatório.

Esta combinação é um padrão fundamental e extremamente robusto para construir sistemas distribuídos resilientes.
