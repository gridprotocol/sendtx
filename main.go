package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rockiecn/sendtx/tx"

	"github.com/grid/contracts/eth"
	"github.com/grid/contracts/eth/contracts"
)

func main() {
	var txType uint
	flag.UintVar(&txType, "tx", 1, "1=register 2=approve 3=createOrder 4=revise 5=userConfirm 6=userCancel")
	chain := flag.String("chain", "local", "local:local chain, sepo:sepolia test chain")
	auto := flag.Bool("auto", false, "auto send the tx to chain")

	flag.Parse()

	fmt.Println("chain: ", *chain)

	fmt.Println("type:", txType)

	var endpoint string
	switch *chain {
	case "local":
		endpoint = tx.Endpoint

		// load contracts
		local := contracts.Local{}
		local.LoadPath("../grid-contracts/eth/contracts/local.json")
		tx.Contracts = local.Contracts

		fmt.Println("contract addresses:", tx.Contracts)
	case "sepo":
		endpoint = tx.Endpoint2

		// load contracts
		sepo := contracts.Sepo{}
		sepo.LoadPath("../grid-contracts/eth/contracts/sepo.json")
		tx.Contracts = sepo.Contracts

		fmt.Println("contract addresses:", tx.Contracts)
	}

	// connect to an eth client
	log.Println("connecting client")
	c, err := ethclient.Dial(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	switch txType {
	case 1:
		// signed register tx for send to chain directly
		regTx, js, err := makeRegisterTx(c)
		if err != nil {
			log.Fatal(err)
		}
		_ = regTx
		_ = js

		log.Printf("signedTx for [registcp]: \n%s\n", js)

		if *auto {
			sendTx(endpoint, c, regTx)
		}
	case 2:
		// tx for send to chain directly
		apprTx, js, err := makeApproveTx(c)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("signedTx for [approve] credit to market: \n%s\n", js)

		if *auto {
			sendTx(endpoint, c, apprTx)
		}
	case 3:
		// signed market.createorder tx for send to chain directly
		coTx, js, err := makeCreateOrderTx(c)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("signedTx for [createorder]: \n%s\n", js)

		if *auto {
			sendTx(endpoint, c, coTx)
		}
	case 4:
		// signed registry.revise tx for send to chain directly
		reviseTx, js, err := makeReviseTx(c)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("signedTx for [revise]: \n%s\n", js)

		if *auto {
			sendTx(endpoint, c, reviseTx)
		}

	case 5:
		// signed market.userconfirm tx for send to chain directly
		confirmTx, js, err := makeUserConfirmTx(c)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("signedTx for [userconfirm]: \n%s\n", js)

		if *auto {
			sendTx(endpoint, c, confirmTx)
		}

	case 6:
		// signed market.userconfirm tx for send to chain directly
		cancelTx, js, err := makeUserCancelTx(c)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("signedTx for [usercancel]: \n%s\n", js)

		if *auto {
			sendTx(endpoint, c, cancelTx)
		}

	}
}

// make tx for register cp
func makeRegisterTx(c *ethclient.Client) (*types.Transaction, []byte, error) {
	// tx data
	data := tx.RegisterData()

	log.Println("making signed register tx")
	// make a signed tx
	fmt.Println("cp sk: ", tx.P_SK)
	signedTx, err := tx.MakeSignedTx(c, tx.P_SK, common.HexToAddress(tx.Contracts.Registry), nil, 1000000, data)
	if err != nil {
		return nil, nil, err
	}

	// marshal tx into json
	js, err := signedTx.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	return signedTx, js, nil
}

// make tx for approving credit to market
func makeApproveTx(c *ethclient.Client) (*types.Transaction, []byte, error) {
	// data for tx
	data := tx.ApproveData()

	log.Println("making approve tx")
	// make a signed tx for approve to credit
	signedTx, err := tx.MakeSignedTx(c, tx.U_SK, common.HexToAddress(tx.Contracts.Credit), nil, 1000000, data)
	if err != nil {
		log.Fatal(err)
	}

	// marshal tx into json
	js, err := signedTx.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	return signedTx, js, nil
}

// make tx for create order
func makeCreateOrderTx(c *ethclient.Client) (*types.Transaction, []byte, error) {
	// data for tx
	data := tx.CreateOrderData()

	log.Println("making createorder tx")
	// make a signed tx for createorder, sender must be user
	signedTx, err := tx.MakeSignedTx(c, tx.U_SK, common.HexToAddress(tx.Contracts.Market), nil, 1000000, data)
	if err != nil {
		log.Fatal(err)
	}

	// marshal tx into json
	js, err := signedTx.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	return signedTx, js, nil
}

// make tx for calling registry.revise
func makeReviseTx(c *ethclient.Client) (*types.Transaction, []byte, error) {
	// data for tx
	data := tx.ReviseData()

	log.Println("making registry.revise tx")
	// make a signed tx for revise, sender must be provider
	signedTx, err := tx.MakeSignedTx(c, tx.P_SK, common.HexToAddress(tx.Contracts.Registry), nil, 1000000, data)
	if err != nil {
		log.Fatal(err)
	}

	// marshal tx into json
	js, err := signedTx.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	return signedTx, js, nil
}

// make tx for user confirm
func makeUserConfirmTx(c *ethclient.Client) (*types.Transaction, []byte, error) {
	// data for tx
	data := tx.UserConfirmData()

	log.Println("making user confirm tx")
	// make a signed tx for createorder, sender must be user
	signedTx, err := tx.MakeSignedTx(c, tx.U_SK, common.HexToAddress(tx.Contracts.Market), nil, 1000000, data)
	if err != nil {
		log.Fatal(err)
	}

	// marshal tx into json
	js, err := signedTx.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	return signedTx, js, nil
}

func makeUserCancelTx(c *ethclient.Client) (*types.Transaction, []byte, error) {
	// data for tx
	data := tx.UserCancelData()

	log.Println("making user cancel tx")
	// make a signed tx for createorder, sender must be user
	signedTx, err := tx.MakeSignedTx(c, tx.U_SK, common.HexToAddress(tx.Contracts.Market), nil, 1000000, data)
	if err != nil {
		log.Fatal(err)
	}

	// marshal tx into json
	js, err := signedTx.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	return signedTx, js, nil
}

// send tx to chain
func sendTx(ep string, c *ethclient.Client, tx *types.Transaction) {
	log.Printf("sending signed tx")

	// send the tx to client
	if err := c.SendTransaction(context.Background(), tx); err != nil {
		fmt.Println("send tx failed:", err.Error())
		return
	}

	// wait tx ok
	fmt.Println("waiting for tx to be ok")
	err := eth.CheckTx(ep, tx.Hash(), "")
	if err != nil {
		fmt.Println("tx failed:", err.Error())
		return
	}

	fmt.Println("tx ok")
}
