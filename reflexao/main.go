package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// --- Nossas Structs para validação ---
type User struct {
	Name  string `validate:"required"`
	Email string `validate:"required"`
	Age   int    `validate:"required"`
}

type Product struct {
	SKU    string `validate:"len=8"` // Stock Keeping Unit
	Price  int    // Sem validação
	Status string `validate:"required"`
}

// validate é nossa função genérica que usa reflexão.
// Ela recebe um ponteiro para qualquer struct.
func validate(s interface{}) (bool, error) {
	// 1. Obter o `reflect.Value` do que foi passado.
	val := reflect.ValueOf(s)

	// --- Verificações de Segurança ---
	// A reflexão precisa de um ponteiro para poder inspecionar os campos.
	// Kind() retorna o tipo subjacente.
	if val.Kind() != reflect.Ptr {
		return false, errors.New("validate: input must be a pointer to a struct")
	}

	// val.Elem() "dereferencia" o ponteiro, nos dando acesso à struct em si.
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return false, errors.New("validate: input must be a pointer to a struct")
	}
	// --- Fim das Verificações ---

	// 2. Iterar sobre os campos da struct.
	for i := 0; i < val.NumField(); i++ {
		// `val.Field(i)` nos dá o `reflect.Value` do campo.
		fieldVal := val.Field(i)
		// `val.Type().Field(i)` nos dá o `reflect.Type` do campo, onde as tags vivem.
		fieldTyp := val.Type().Field(i)

		// 3. Ler a tag de validação.
		tag := fieldTyp.Tag.Get("validate")
		if tag == "" {
			continue // Sem tag, sem validação.
		}

		// 4. Aplicar as regras de validação.
		rules := strings.Split(tag, ",")
		for _, rule := range rules {
			if rule == "required" {
				// IsZero() é uma forma segura de verificar se um valor é o "valor zero"
				// para seu tipo (0, "", nil, etc).
				if fieldVal.IsZero() {
					return false, fmt.Errorf("field '%s' is required", fieldTyp.Name)
				}
			}

			if strings.HasPrefix(rule, "len=") {
				// Garantir que estamos aplicando esta regra a uma string.
				if fieldVal.Kind() != reflect.String {
					return false, fmt.Errorf("rule 'len' can only be applied to string fields, but '%s' is a %v", fieldTyp.Name, fieldVal.Kind())
				}

				lenStr := strings.TrimPrefix(rule, "len=")
				expectedLen, err := strconv.Atoi(lenStr)
				if err != nil {
					return false, fmt.Errorf("invalid 'len' value '%s' on field '%s'", lenStr, fieldTyp.Name)
				}

				if fieldVal.Len() != expectedLen {
					return false, fmt.Errorf("field '%s' must have a length of %d, but has %d", fieldTyp.Name, expectedLen, fieldVal.Len())
				}
			}
		}
	}

	return true, nil
}

func main() {
	fmt.Println("--- Validando Usuários ---")
	user1 := User{Name: "Alice", Email: "alice@example.com", Age: 30}
	valid, err := validate(&user1) // Passando um ponteiro!
	fmt.Printf("User1: %v -> Válido: %v, Erro: %v\n", user1, valid, err)

	user2 := User{Name: "Bob", Email: "", Age: 25} // Email está faltando
	valid, err = validate(&user2)
	fmt.Printf("User2: %v -> Válido: %v, Erro: %v\n", user2, valid, err)

	fmt.Println("\n--- Validando Produtos ---")
	prod1 := Product{SKU: "ABC-1234", Price: 100, Status: "available"} // SKU com 8 caracteres
	valid, err = validate(&prod1)
	fmt.Printf("Prod1: %v -> Válido: %v, Erro: %v\n", prod1, valid, err)

	prod2 := Product{SKU: "SHORT", Price: 200, Status: "available"} // SKU muito curto
	valid, err = validate(&prod2)
	fmt.Printf("Prod2: %v -> Válido: %v, Erro: %v\n", prod2, valid, err)

	// Exemplo de erro de uso (passando um não-ponteiro)
	fmt.Println("\n--- Testando Erro de Uso ---")
	invalidInput := "sou uma string"
	valid, err = validate(invalidInput)
	fmt.Printf("Input Inválido: Válido: %v, Erro: %v\n", valid, err)
}
