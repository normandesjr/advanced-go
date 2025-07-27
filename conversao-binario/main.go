package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

// A mesma struct de antes, mas sem tags, pois não são necessárias.
type Transaction struct {
	CorrelationID string
	Amount        float64
	Default       bool
}

// MarshalBinary converte o struct para nossa representação binária customizada.
// Este é o método que "arruma a caixa de ferramentas".
func (t *Transaction) MarshalBinary() ([]byte, error) {
	// Usamos um buffer para escrever os bytes de forma eficiente.
	buf := new(bytes.Buffer)

	// 1. Escreve o tamanho do CorrelationID (4 bytes) e depois a string.
	// Usamos BigEndian como uma convenção.
	if err := binary.Write(buf, binary.BigEndian, uint32(len(t.CorrelationID))); err != nil {
		return nil, err
	}
	if _, err := buf.WriteString(t.CorrelationID); err != nil {
		return nil, err
	}

	// 2. Escreve o Amount (8 bytes).
	if err := binary.Write(buf, binary.BigEndian, t.Amount); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, t.Default); err != nil {
		return nil, err
	}

	// 3. Escreve o tamanho do Endpoint (4 bytes) e depois a string.
	// if err := binary.Write(buf, binary.BigEndian, uint32(len(t.Endpoint))); err != nil {
	// 	return nil, err
	// }
	// if _, err := buf.WriteString(t.Endpoint); err != nil {
	// 	return nil, err
	// }

	return buf.Bytes(), nil
}

// UnmarshalBinary converte os bytes de volta para o struct.
// Este é o método que "pega as ferramentas da caixa".
func (t *Transaction) UnmarshalBinary(data []byte) error {
	// Usamos um reader para ler os bytes de forma sequencial.
	buf := bytes.NewReader(data)

	// 1. Lê o tamanho do CorrelationID, e depois a string.
	var idLen uint32
	if err := binary.Read(buf, binary.BigEndian, &idLen); err != nil {
		return err
	}
	idBytes := make([]byte, idLen)
	if _, err := buf.Read(idBytes); err != nil {
		return err
	}
	t.CorrelationID = string(idBytes)

	// 2. Lê o Amount.
	if err := binary.Read(buf, binary.BigEndian, &t.Amount); err != nil {
		return err
	}

	if err := binary.Read(buf, binary.BigEndian, &t.Default); err != nil {
		return err
	}

	// 3. Lê o tamanho do Endpoint, e depois a string.
	// var endpointLen uint32
	// if err := binary.Read(buf, binary.BigEndian, &endpointLen); err != nil {
	// 	return err
	// }
	// endpointBytes := make([]byte, endpointLen)
	// if _, err := buf.Read(endpointBytes); err != nil {
	// 	return err
	// }
	// t.Endpoint = string(endpointBytes)

	return nil
}

func main() {
	tx := Transaction{
		CorrelationID: "a1b2c3d4-e5f6-7890-1234-567890abcdef",
		Amount:        19.99,
		Default:       true,
	}

	// --- SERIALIZAÇÃO (Guardar na caixa) ---
	customBytes, err := tx.MarshalBinary()
	if err != nil {
		log.Fatalf("Erro ao serializar para binário customizado: %v", err)
	}

	fmt.Println("--- Usando Binário Customizado (Sem libs externas) ---")
	fmt.Printf("Struct original: %+v\n", tx)
	fmt.Printf("Resultado binário: %x\n", customBytes)
	fmt.Printf("Tamanho dos dados: %d bytes\n", len(customBytes))
	fmt.Println("-----------------------------------------------------\n")

	// --- DESSERIALIZAÇÃO (Pegar da caixa) ---
	var newTx Transaction
	err = newTx.UnmarshalBinary(customBytes)
	if err != nil {
		log.Fatalf("Erro ao desserializar de binário customizado: %v", err)
	}

	fmt.Println("--- Resultado da Desserialização ---")
	fmt.Printf("Struct após ler os dados binários: %+v\n", newTx)
	fmt.Println("----------------------------------\n")
}
