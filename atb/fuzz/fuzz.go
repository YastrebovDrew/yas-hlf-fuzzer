//go:build gofuzz
// +build gofuzz

// fuzz_atb.go
// -----------------------------------------------------------------------------
// Harness для go-fuzz, повторяет стиль примера «Marbles» и обслуживает
// смарт-контракт Asset-Transfer-Basic (ATB):
//
//   • работает через shimtest.MockStub (Fabric не нужен);
//   • принимает два формата входа:
//       1) JSON-массив строк          ["CreateAsset","id",…]
//       2) строка с разделителем '|'  CreateAsset|id|…
//   • ввод ограничен по размеру/кол-ву аргументов;
//   • охватывает все публичные функции контракта.
// -----------------------------------------------------------------------------

package chaincode

import (
	"encoding/json"
	"strings"
	"unicode/utf8"

	_ "github.com/dvyukov/go-fuzz/go-fuzz-dep"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	cc "github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go/chaincode"
)

// ------------------- лимиты безопасности ------------------------------------
const (
	maxInputSize   = 6 << 10 // 6 KiB максимально на всё
	maxArgs        = 16
	maxArgLenBytes = 2 << 10 // 2 KiB на один арг-т
)

// ---------- таблица «функция → требуемое кол-во аргументов» ------------------
func wantArgs(fn string) int {
	switch fn {
	case "InitLedger", "GetAllAssets":
		return 0
	case "CreateAsset", "UpdateAsset":
		return 5 // id color size owner value
	case "ReadAsset", "DeleteAsset":
		return 1
	case "TransferAsset":
		return 2
	default:
		return -1
	}
}

// ----------------------- split "a|b|c" ---------------------------------------
func splitPipe(s string) []string {
	if s == "" {
		return nil
	}
	var out []string
	start := 0
	for i, r := range s {
		if r == '|' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	out = append(out, s[start:])
	return out
}

// ------------------------------- Fuzz() --------------------------------------
func Fuzz(data []byte) int {
	if len(data) == 0 || len(data) > maxInputSize {
		return 0
	}

	// 1️⃣ JSON-массив
	var args []string
	if err := json.Unmarshal(data, &args); err != nil || len(args) == 0 {
		// 2️⃣ строка ‘|’-разделённая
		args = splitPipe(string(data))
		if len(args) == 0 {
			return 0
		}
	}

	if len(args) > maxArgs {
		return 0
	}
	for _, a := range args {
		if !utf8.ValidString(a) || len(a) > maxArgLenBytes {
			return 0
		}
	}

	fn := args[0]
	if need := wantArgs(fn); need < 0 || len(args)-1 != need {
		return 0
	}

	// ---------- инициализация MockStub ---------------------------------------
	chaincode, err := contractapi.NewChaincode(new(cc.SmartContract))
	if err != nil {
		return 0
	}
	stub := shimtest.NewMockStub("atb", chaincode)

	// InitLedger
	if res := stub.MockInit("tx_init", [][]byte{[]byte("InitLedger")}); res.Status != shim.OK {
		return 0
	}

	// ---------- основной fuzz-вызов ------------------------------------------
	bin := make([][]byte, len(args))
	for i, a := range args {
		bin[i] = []byte(a)
	}
	invokeRes := stub.MockInvoke("tx_fuzz", bin)

	// ← вот тут добавляем фильтр «Conversion error»
	if invokeRes.Status != shim.OK {
		// пропускаем ошибки преобразования параметров
		if strings.Contains(invokeRes.Message, "Conversion error") {
			return 0
		}
		// все остальные статусы считаем крэшем
		panic(invokeRes.Message)
	}

	return 1 // любой результат (OK/ERR) интересен
}
