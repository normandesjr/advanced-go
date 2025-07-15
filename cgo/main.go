package main

/*
#cgo LDFLAGS: -lz

#include <stdio.h>
#include <stdlib.h>
#include <zlib.h>

// zlib_compress é uma função wrapper em C para facilitar a chamada a partir do Go.
// Ela esconde a complexidade da alocação de memória e da chamada à zlib.
// Retorna 0 em caso de sucesso, ou um código de erro zlib em caso de falha.
int zlib_compress(const char* src, size_t src_len, char** dest, size_t* dest_len) {
    // 1. Estimar o tamanho máximo que o buffer de destino precisará.
    *dest_len = compressBound(src_len);

    // 2. Alocar memória no heap do C. O chamador (Go) será responsável por liberá-la!
    *dest = (char*)malloc(*dest_len);
    if (*dest == NULL) {
        return Z_MEM_ERROR;
    }

    // 3. Chamar a função de compressão da zlib.
    int res = compress(*dest, dest_len, (const Bytef*)src, src_len);
    return res;
}
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

// compressData é o nosso wrapper Go que chama a função C.
func compressData(input []byte) ([]byte, error) {
	// --- Lado do Go: Preparação para a chamada C ---

	// 1. Converter o slice de bytes Go para um ponteiro C.
	// `C.CBytes` aloca memória (que o Go GC conhece) e copia os dados.
	// É a maneira segura de passar slices de bytes.
	inputPtr := C.CBytes(input)
	// `defer` não funciona com `C.free`, então devemos gerenciar manualmente.
	// No entanto, C.CBytes é especial e o GC do Go vai lidar com isso.
	// Mas para C.malloc, o defer é essencial.

	// Ponteiros para receber os resultados da função C.
	var destPtr *C.char
	var destLen C.size_t

	// --- A Chamada: Cruzando a Fronteira ---
	// Chamamos nossa função C wrapper.
	ret := C.zlib_compress(
		(*C.char)(inputPtr),
		C.size_t(len(input)),
		&destPtr,
		&destLen,
	)

	// --- Lado do Go: Processando o Retorno ---

	// 2. É CRUCIAL liberar a memória que foi alocada por `malloc` dentro do C.
	// Usamos `defer` para garantir que isso aconteça mesmo se houver erros depois.
	if destPtr != nil {
		defer C.free(unsafe.Pointer(destPtr))
	}

	// 3. Verificar se a chamada C foi bem-sucedida.
	if ret != C.Z_OK {
		return nil, fmt.Errorf("compressão zlib falhou com o código de erro: %d", ret)
	}
	if destPtr == nil {
		return nil, errors.New("a compressão zlib retornou um ponteiro nulo")
	}

	// 4. Converter os dados do buffer C de volta para um slice Go.
	// `C.GoBytes` copia os dados da memória C para um novo slice gerenciado pelo Go.
	output := C.GoBytes(unsafe.Pointer(destPtr), C.int(destLen))

	return output, nil
}

func main() {
	data := []byte("a_long_string_that_will_compress_well_a_long_string_that_will_compress_well_a_long_string_that_will_compress_well")
	fmt.Printf("Dados originais (%d bytes): %s\n", len(data), data)

	compressed, err := compressData(data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nDados comprimidos (%d bytes): %x\n", len(compressed), compressed)
}
