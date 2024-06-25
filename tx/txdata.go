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
	"github.com/grid/contracts/go/market"
)

var (
	// chain node
	Endpoint = "HTTP://127.0.0.1:7545"

	// admin
	A_SK   = "c1e763d955e6aea410e40b95702108a30efb4d25b32d419910fe2ac611c2229d"
	A_ADDR = "0x5F7F7e31399531F08C2b47eA1919F11346405a16"

	// sk for user, acc1
	U_SK   = "e8cda8fe7c04afa4a0630af457972f88a645468cb90120a11911669deac5e96e"
	U_ADDR = "0xe2198eb2e931f9306ABcA68D4F093E0Ac4823B0d"
	// sk for provider, acc2
	P_SK   = "ae313c1dc715026cf629641e3ae2dab06f95c7300d97b3121310d375b979f19b"
	P_ADDR = "0xC4EAc9E1012DFCB4833165F5d35E027EBfE1f640"

	// abis
	RegABI    string
	MarketABI string
	CreditABI string

	// all contracts addresses
	Contracts = eth.Address{}
)

// read ABI from file
func init() {
	Contracts = eth.LoadJSON()
	fmt.Println("contract addresses:", Contracts)

	// read registry abi from file
	_RegABI, err := os.ReadFile("../../grid-contracts/abi/registry/Registry.abi")
	if err != nil {
		log.Fatal(err)
	}

	// read abi from file
	_MarketABI, err := os.ReadFile("../../grid-contracts/abi/market/Market.abi")
	if err != nil {
		log.Fatal(err)
	}

	// read abi from file
	_CreditABI, err := os.ReadFile("../../grid-contracts/abi/credit/Credit.abi")
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
	// 假设我们有一个已编译的合约ABI
	registryABI, err := abi.JSON(strings.NewReader(RegABI))
	if err != nil {
		panic(err)
	}

	// 我们要调用的函数和参数
	functionName := "register"
	//value := big.NewInt(0)

	// 构造调用函数和参数的方法和输入参数
	method := registryABI.Methods[functionName]
	input, err := method.Inputs.Pack("a", "b", "c", uint64(10), uint64(20), uint64(40), uint64(807), uint64(33), uint64(33), uint64(33), uint64(33))
	if err != nil {
		panic(err)
	}

	// 构造完整的data字段
	data := append(method.ID, input...)
	//fmt.Printf("Data: %x\n", data)

	return data
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
	amount, ok := new(big.Int).SetString("2831300", 10)
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
	totalValue, ok := new(big.Int).SetString("2626954", 10)
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

		P: market.MarketPricePerHour{
			PCPU:  100,
			PGPU:  1000,
			PMEM:  10,
			PDISK: 1,
		},
		R: market.MarketResources{
			NCPU:  1,
			NGPU:  2,
			NMEM:  3,
			NDISK: 4,
		},
		// deposit 0.01 eth
		TotalValue:      totalValue,
		Remain:          totalValue,
		Remuneration:    remu,
		UserConfirm:     false,
		ProviderConfirm: false,
		ActivateTime:    new(big.Int).SetInt64(0),
		LastSettleTime:  new(big.Int).SetInt64(0),
		Probation:       new(big.Int).SetInt64(5),
		Duration:        new(big.Int).SetInt64(1231),
		Status:          0, // unactive
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

	// construct the input of this method
	input, err := method.Inputs.Pack("revised ip", "revised domain", "revised port", uint64(22), uint64(33), uint64(44), uint64(55), uint64(33), uint64(33), uint64(33), uint64(33))
	if err != nil {
		panic(err)
	}

	// the full data of tx
	data := append(method.ID, input...)

	return data
}
