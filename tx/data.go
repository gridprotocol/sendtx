package tx

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/grid/contracts/eth"
	"github.com/grid/contracts/eth/contracts"
	"github.com/grid/contracts/go/market"
	"github.com/grid/contracts/go/registry"
)

var (
	// chain node
	Endpoint = "HTTP://127.0.0.1:7545"

	// sepolia chain
	Endpoint2 string = "https://rpc.sepolia.ethpandaops.io"

	// // admin
	// A_SK   = "c1e763d955e6aea410e40b95702108a30efb4d25b32d419910fe2ac611c2229d"
	// A_ADDR = "0x5F7F7e31399531F08C2b47eA1919F11346405a16"

	// // sk for user, acc1
	// U_SK   = "e8cda8fe7c04afa4a0630af457972f88a645468cb90120a11911669deac5e96e"
	// U_ADDR = "0xe2198eb2e931f9306ABcA68D4F093E0Ac4823B0d"
	// // sk for provider, acc2
	// P_SK   = "ae313c1dc715026cf629641e3ae2dab06f95c7300d97b3121310d375b979f19c"
	// P_ADDR = "0xC4EAc9E1012DFCB4833165F5d35E027EBfE1f640"

	A_SK   = eth.SK0
	A_ADDR = eth.HexAddr0
	U_SK   = eth.SK1
	U_ADDR = eth.HexAddr1
	P_SK   = eth.SK2
	P_ADDR = eth.HexAddr2

	// abis
	RegABI    string
	MarketABI string
	CreditABI string

	// all contracts addresses
	Contracts = contracts.Contracts{}

	// abi file path for all contracts used
	REG_ABI_PATH = "../grid-contracts/abi/registry/Registry.abi"
	MAR_ABI_PATH = "../grid-contracts/abi/market/Market.abi"
	CRE_ABI_PATH = "../grid-contracts/abi/credit/Credit.abi"
)

// read ABI from file
func init() {
	// read registry abi from file
	_RegABI, err := os.ReadFile(REG_ABI_PATH)
	if err != nil {
		log.Fatal(err)
	}

	// read abi from file
	_MarketABI, err := os.ReadFile(MAR_ABI_PATH)
	if err != nil {
		log.Fatal(err)
	}

	// read abi from file
	_CreditABI, err := os.ReadFile(CRE_ABI_PATH)
	if err != nil {
		log.Fatal(err)
	}

	// local to global
	RegABI = string(_RegABI)
	MarketABI = string(_MarketABI)
	CreditABI = string(_CreditABI)
}

// the tx data for calling registry.register
func RegisterData() []byte {
	// abi
	registryABI, err := abi.JSON(strings.NewReader(RegABI))
	if err != nil {
		panic(err)
	}

	// func name
	functionName := "register"
	//value := big.NewInt(0)

	// make a cp
	info, err := newCP()
	if err != nil {
		panic(err)
	}

	// pack params
	method := registryABI.Methods[functionName]
	//input, err := method.Inputs.Pack("a", "b", "c", uint64(10), uint64(20), uint64(40), uint64(807), uint64(33), uint64(33), uint64(33), uint64(33))
	input, err := method.Inputs.Pack(info)
	if err != nil {
		panic(err)
	}

	// 构造完整的data字段
	data := append(method.ID, input...)
	//fmt.Printf("Data: %x\n", data)

	return data
}

// generate a test cp
func newCP() (*registry.IRegistryInfo, error) {
	// the register cp info
	info := registry.IRegistryInfo{
		Addr:   common.HexToAddress("0xEf95c72C836605203F7f66788E450Af2a4141957"),
		Name:   "cp1",
		Ip:     "123.123.123.0",
		Domain: "testdomain",
		Port:   "123",
	}

	return &info, nil
}

// for revise
func newCP2() (*registry.IRegistryInfo, error) {
	// the register cp info
	info := registry.IRegistryInfo{
		Addr:   common.HexToAddress("0xEf95c72C836605203F7f66788E450Af2a4141957"),
		Name:   "revised name",
		Ip:     "revise ip",
		Domain: "revised domain",
		Port:   "revised port",
	}

	return &info, nil
}

// the tx data for call add_node
func AddNodeData(node *registry.IRegistryNode) []byte {
	// registry abi
	registryABI, err := abi.JSON(strings.NewReader(RegABI))
	if err != nil {
		panic(err)
	}

	// func name
	functionName := "add_node"
	//value := big.NewInt(0)

	// pack params
	method := registryABI.Methods[functionName]
	//input, err := method.Inputs.Pack("a", "b", "c", uint64(10), uint64(20), uint64(40), uint64(807), uint64(33), uint64(33), uint64(33), uint64(33))
	input, err := method.Inputs.Pack(node)
	if err != nil {
		panic(err)
	}

	// 构造完整的data字段
	data := append(method.ID, input...)
	//fmt.Printf("Data: %x\n", data)

	return data
}

func NewNode() (*registry.IRegistryNode, error) {
	// the register cp info
	info := registry.IRegistryNode{
		Id: new(big.Int).SetInt64(0),

		Cpu: registry.IRegistryCPU{
			Price: new(big.Int).SetInt64(10),
			Model: "i5",
		},
		Gpu: registry.IRegistryGPU{
			Price: new(big.Int).SetInt64(100),
			Model: "RTX4080",
		},
		Mem: registry.IRegistryMEM{
			Num:   new(big.Int).SetInt64(1),
			Price: new(big.Int).SetInt64(10),
		},
		Disk: registry.IRegistryDISK{
			Num:   new(big.Int).SetInt64(1),
			Price: new(big.Int).SetInt64(10),
		},
	}

	return &info, nil
}

// make a node
func NewNode2() (*registry.IRegistryNode, error) {
	// the register cp info
	info := registry.IRegistryNode{
		Id: new(big.Int).SetInt64(0),

		Cpu: registry.IRegistryCPU{
			Price: new(big.Int).SetInt64(1),
			Model: "i7",
		},
		Gpu: registry.IRegistryGPU{
			Price: new(big.Int).SetInt64(1),
			Model: "RTX4090",
		},
		Mem: registry.IRegistryMEM{
			Num:   new(big.Int).SetInt64(1),
			Price: new(big.Int).SetInt64(1),
		},
		Disk: registry.IRegistryDISK{
			Num:   new(big.Int).SetInt64(1),
			Price: new(big.Int).SetInt64(1),
		},
	}

	return &info, nil
}

// the tx data for calling credit.approve
//
//	function approve(address spender, uint256 amount) public virtual override returns (bool) {
func ApproveData() []byte {
	// contract abi
	creditABI, err := abi.JSON(strings.NewReader(CreditABI))
	if err != nil {
		panic(err)
	}

	// function to be called
	functionName := "approve"

	// method with name
	method := creditABI.Methods[functionName]

	// set the amount to approve to market contract
	amount, ok := new(big.Int).SetString("262695400", 10)
	if !ok {
		panic(fmt.Errorf("big set string failed"))
	}

	// input for calling method
	input, err := method.Inputs.Pack(common.HexToAddress(Contracts.Market), amount)
	if err != nil {
		panic(err)
	}

	// 构造完整的data字段
	data := append(method.ID, input...)

	return data
}

// the tx data for calling market.createorder
func CreateOrderData() []byte {
	// 假设我们有一个已编译的合约ABI
	marketABI, err := abi.JSON(strings.NewReader(MarketABI))
	if err != nil {
		panic(err)
	}

	// 我们要调用的函数和参数
	functionName := "createOrder"

	// 构造调用函数和参数的方法和输入参数
	method := marketABI.Methods[functionName]
	order, err := newOrder()
	if err != nil {
		panic(err)
	}

	fmt.Println("packing")
	// pack all params into input
	input, err := method.Inputs.Pack(eth.Addr2, *order)
	if err != nil {
		panic(err)
	}

	// 构造完整的data字段
	data := append(method.ID, input...)

	return data
}

// generate a test order
func newOrder() (*market.MarketOrder, error) {
	// generate an order with init data
	totalValue, ok := new(big.Int).SetString("40", 10)
	if !ok {
		return nil, fmt.Errorf("big set string failed")
	}
	remu, ok := new(big.Int).SetString("0", 10)
	if !ok {
		return nil, fmt.Errorf("big set string failed")
	}

	// make an order
	order := market.MarketOrder{
		User:     eth.Addr1,
		Provider: eth.Addr2,

		// the cp's node selected bye this order
		NodeId: new(big.Int).SetInt64(1),

		// deposit 0.01 eth
		TotalValue:     totalValue,
		Remain:         totalValue,
		Remuneration:   remu,
		ActivateTime:   new(big.Int).SetInt64(0),
		LastSettleTime: new(big.Int).SetInt64(0),
		Probation:      new(big.Int).SetInt64(5),
		Duration:       new(big.Int).SetInt64(10),
		Status:         1, // unactive
	}

	return &order, nil
}

// the tx data for calling registry.revise
func ReviseData() []byte {
	// contract abi
	registryABI, err := abi.JSON(strings.NewReader(RegABI))
	if err != nil {
		panic(err)
	}

	// function to be called
	functionName := "revise"
	method := registryABI.Methods[functionName]

	info, err := newCP2()
	if err != nil {
		panic(err)
	}
	// construct the input of this method
	input, err := method.Inputs.Pack(*info)
	if err != nil {
		panic(err)
	}

	// the full data of tx
	data := append(method.ID, input...)

	return data
}

// tx data for user confirm
func UserConfirmData() []byte {
	// contract abi
	marketABI, err := abi.JSON(strings.NewReader(MarketABI))
	if err != nil {
		panic(err)
	}

	// function to be called
	functionName := "userConfirm"
	method := marketABI.Methods[functionName]

	// construct the input of this method
	input, err := method.Inputs.Pack(eth.Addr2)
	if err != nil {
		panic(err)
	}

	// the full data of tx
	data := append(method.ID, input...)

	return data
}

func UserCancelData() []byte {
	// contract abi
	marketABI, err := abi.JSON(strings.NewReader(MarketABI))
	if err != nil {
		panic(err)
	}

	// function to be called
	functionName := "userCancel"
	method := marketABI.Methods[functionName]

	// construct the input of this method
	input, err := method.Inputs.Pack(eth.Addr2)
	if err != nil {
		panic(err)
	}

	// the full data of tx
	data := append(method.ID, input...)

	return data
}
