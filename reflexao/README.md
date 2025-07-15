# Reflection

## O que é Reflexão (Reflection)?

Imagine que você tem uma caixa preta. Você não sabe o que tem dentro. A reflexão é como ter um conjunto de ferramentas especiais (raio-x, pinças, medidores) que permitem que você, em tempo de execução, abra a caixa, examine o que está dentro, descubra seu tipo (é uma bola? um cubo?), suas propriedades (qual a sua cor? qual o seu peso?) e até mesmo modifique o que está lá.

Em Go, a **Reflexão** é a capacidade de um programa examinar e manipular seus próprios objetos em tempo de execução. Em vez de operar sobre um tipo estático conhecido (`int`, `string`, `User`), você opera sobre `interface{}` e usa o pacote `reflect` para descobrir o tipo e o valor concretos que essa interface contém.

## As Leis da Reflexão (e a Regra de Ouro)

O uso de reflect é governado por algumas "leis" e uma regra de ouro muito importante.

1. **Reflexão vai de uma `interface{}` para um objeto de reflexão**. Você pode inspecionar o tipo e o valor de qualquer variável passada como interface{}.
2. **Reflexão vai de um objeto de reflexão de volta para uma `interface{}`**. Você pode pegar um objeto de reflexão e transformá-lo de volta em uma interface{}, recuperando o valor original.
3. **Para modificar um objeto de reflexão, o valor deve ser settable (modificável)**. Este é o ponto mais crucial e fonte de muitos erros. Para que a reflexão possa modificar um valor, ela precisa do endereço desse valor. Você sempre deve passar um ponteiro para a função que usará reflexão para modificar algo.

E a regra de ouro, dita por Rob Pike, um dos criadores do Go:

| "Reflection is never clear." (Reflexão nunca é clara)

Isso nos leva à pergunta mais importante:

## Quando Usar Reflexão? (Com Muito Cuidado!)

O uso de reflexão vem com um custo:

* **Perda de Segurança de Tipos**: Você troca a verificação em tempo de compilação por panics em tempo de execução.
* **Performance**: É significativamente mais lento que o código normal.
* **Legibilidade**: O código se torna mais complexo e mais difícil de entender e manter.

**Portanto, a regra é: se você pode resolver seu problema sem usar reflexão, faça-o sem reflexão.**

Então, quando ela é realmente necessária?

* **Serialização/Deserialização**: A principal aplicação. O pacote encoding/json é o exemplo perfeito. Ele precisa receber qualquer struct, inspecionar seus campos e tags (json:"name") para saber como converter de/para JSON. Ele não conhece sua struct em tempo de compilação.
* **Frameworks Genéricos**: Bibliotecas que criam validadores, ORMs (Object-Relational Mappers) ou contêineres de injeção de dependência. Essas ferramentas são projetadas para funcionar com tipos de dados definidos pelo usuário.

## Exemplo Prático: Um Validador de Structs Genérico

Vamos criar uma função `validate` que recebe qualquer `struct` (através de um ponteiro) e valida seus campos com base em `tags`, assim como o pacote `json` faz.

Nosso validador suportará duas regras:

* `validate:"required"`: O campo não pode ser o valor zero do seu tipo (vazio para string, 0 para int).
* `validate:"len=N"`: Para strings, o comprimento deve ser exatamente `N`.

Veja o main.go

## Análise do Código: Como Usar reflect de Forma Segura

1. **Sempre Espere um Ponteiro (`reflect.Ptr`)**: A nossa função exige um ponteiro. Checamos isso com val.Kind() != reflect.Ptr. Isso é fundamental porque, sem um ponteiro, o Go passa uma cópia da struct para a função, e qualquer modificação (ou inspeção que precise de "settabilidade") não seria possível no valor original.
2. **Acesse o Valor Real com `.Elem()`**: Uma vez que temos um ponteiro, `val.Elem()` nos dá acesso ao valor para o qual ele aponta (a `struct` em si).
3. **Itere sobre os Campos**: `val.NumField()` e `val.Field(i)` são as ferramentas para percorrer cada campo da `struct`.
4. **Acesse o Tipo e as Tags**: Para obter as tags, precisamos do tipo do campo, que obtemos com `val.Type().Field(i).Tag`.
5. **Verifique os Tipos (`Kind`) antes de Operar**: Antes de aplicar a regra len=, verificamos se o campo é de fato uma string com fieldVal.Kind() != reflect.String. Se não fizéssemos isso, chamar .Len() em um int causaria um panic.
6. **Use Funções Seguras**: `fieldVal.IsZero()` é uma maneira robusta de verificar o valor padrão de um campo, independentemente do seu tipo.

## Como Tornar a Reflexão Eficiente? Com Cache!

A reflexão é lenta. Se a função `validate` for chamada milhões de vezes (por exemplo, em um endpoint de API), inspecionar a mesma `struct` repetidamente é um desperdício.

A solução é **fazer a reflexão uma vez e guardar os resultados em cache**.

A ideia é:

* Criar um mapa global (seguro para concorrência, como um `sync.Map`).
* A chave do mapa será o `reflect.Type` da struct.
* O valor será uma estrutura de dados pré-processada com as regras de validação para cada campo.
* A função `validate` primeiro verifica o cache. Se encontrar as regras, aplica-as diretamente, sem usar reflexão. Se não, ela faz a reflexão, armazena as regras no cache e então as aplica.

Essa otimização é um padrão comum em bibliotecas que usam reflexão intensivamente.

Reflexão é um tópico denso, mas entender esses princípios de segurança e eficiência é o que separa o uso casual do uso profissional e robusto da ferramenta.