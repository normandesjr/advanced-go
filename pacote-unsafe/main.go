package main

import (
	"fmt"
	"unsafe"
)

// --- Exemplo 1: Conversão Zero-Copy entre String e Slice de Bytes ---

// bytesToString converte um []byte para uma string SEM alocar nova memória.
// AVISO: Isso é perigoso. Se o slice de bytes original for modificado,
// a string também mudará, quebrando a imutabilidade das strings em Go.
func bytesToString(b []byte) string {
	// A partir do Go 1.22, esta é a forma sancionada de fazer isso.
	// Ela cria uma string cujo conteúdo é o do slice, sem copiar.
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// stringToBytes converte uma string para um []byte SEM alocar nova memória.
// AVISO: Extremamente perigoso. Modificar o slice resultante pode causar um
// panic, pois tentará escrever em memória de dados de string, que é imutável.
func stringToBytes(s string) []byte {
	// A partir do Go 1.22, esta é a forma sancionada de fazer isso.
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// --- Exemplo 2: Aritmética de Ponteiros para Acessar Campos de Struct ---
type MyStruct struct {
	A bool    // 1 byte (com padding de memória)
	B float64 // 8 bytes
	C int32   // 4 bytes
}

func main() {
	fmt.Println("--- Exemplo 1: Conversão String <-> Slice ---")
	byteSlice := []byte{'h', 'e', 'l', 'l', 'o'}
	// Esta conversão não aloca memória nova para os dados da string.
	// Ela aponta para o mesmo array de bytes subjacente.
	str := bytesToString(byteSlice)

	fmt.Printf("Slice original: %v (%s)\n", byteSlice, byteSlice)
	fmt.Printf("String convertida: %s\n", str)

	// O PERIGO REAL: modificar o slice original...
	byteSlice[0] = 'j'
	// ...muda a string, quebrando a garantia de imutabilidade das strings em Go!
	fmt.Println("\n>> Modificando o slice original... <<")
	fmt.Printf("Slice modificado: %v (%s)\n", byteSlice, byteSlice)
	fmt.Printf("String agora é: %s <--- PERIGO! A string foi mutada!\n", str)

	fmt.Println("\n\n--- Exemplo 2: Aritmética de Ponteiros ---")
	s := MyStruct{A: true, B: 3.14, C: 123}

	fmt.Printf("Struct original: %+v\n", s)
	fmt.Printf("Tamanho da struct: %d bytes\n", unsafe.Sizeof(s))

	// Vamos ler e modificar o campo 'B' (float64) usando seu offset,
	// sem usar `s.B`.

	// 1. Obter o ponteiro para o início da struct.
	// Isso nos dá o endereço de memória onde `s` começa.
	structPtr := unsafe.Pointer(&s)

	// 2. Calcular o endereço do campo 'B'.
	//    O endereço é: (endereço da struct) + (offset do campo B).
	//    `unsafe.Offsetof(s.B)` retorna a distância em bytes do início da
	//    struct até o início do campo B.
	offsetB := unsafe.Offsetof(s.B)
	// Convertemos para uintptr para poder somar.
	addressOfB := uintptr(structPtr) + offsetB

	// 3. Converter o endereço numérico de volta para um ponteiro do tipo correto.
	// Nós dizemos ao compilador: "confie em mim, este número é um endereço
	// de memória que aponta para um float64".
	ptrToB := (*float64)(unsafe.Pointer(addressOfB))

	// 4. Ler o valor.
	fmt.Printf("\nOffset do campo B: %d bytes\n", offsetB)
	fmt.Printf("Valor de B lido via unsafe: %f\n", *ptrToB)

	// 5. Modificar o valor através do ponteiro `unsafe`.
	*ptrToB = 9.99
	fmt.Printf("Valor de B modificado via unsafe. Struct agora é: %+v\n", s)
}
