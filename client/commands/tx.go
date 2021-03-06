// Copyright 2017 Annchain Information Technology Services Co.,Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"fmt"
	"math/big"

	agtypes "github.com/annchain/annchain/angine/types"
	"github.com/annchain/annchain/tools"
	"github.com/annchain/annchain/types"
	"github.com/annchain/anth/common"
	"github.com/annchain/anth/crypto"
	"gopkg.in/urfave/cli.v1"

	"github.com/annchain/annchain/client/commons"
	cl "github.com/annchain/annchain/module/lib/go-rpc/client"
)

var (
	TxCommands = cli.Command{
		Name:     "tx",
		Usage:    "operations for transaction",
		Category: "Transaction",
		Subcommands: []cli.Command{
			{
				Name:   "send",
				Usage:  "send a transaction",
				Action: sendTx,
				Flags: []cli.Flag{
					anntoolFlags.payload,
					anntoolFlags.privkey,
					anntoolFlags.nonce,
					anntoolFlags.to,
					anntoolFlags.value,
				},
			},
		},
	}
)

func sendTx(ctx *cli.Context) error {
	skbs := ctx.String("privkey")
	privkey, err := crypto.HexToECDSA(skbs)
	if err != nil {
		panic(err)
	}

	nonce := ctx.Uint64("nonce")
	to := common.HexToAddress(ctx.String("to"))
	value := ctx.Int64("value")
	payload := ctx.String("payload")
	data := common.Hex2Bytes(payload)

	bodyTx := types.TxEvmCommon{
		To:     to[:],
		Amount: big.NewInt(value),
		Load:   data,
	}
	fmt.Printf("%+v\n", bodyTx)
	bodyBs, err := tools.ToBytes(bodyTx)
	if err != nil {
		return cli.NewExitError(err.Error(), 127)
	}

	from := crypto.PubkeyToAddress(privkey.PublicKey)
	tx := types.NewBlockTx(big.NewInt(90000), big.NewInt(2), nonce, from[:], bodyBs)

	if tx.Signature, err = tools.SignSecp256k1(tx, crypto.FromECDSA(privkey)); err != nil {
		return cli.NewExitError(err.Error(), 127)
	}

	b, err := tools.ToBytes(tx)
	if err != nil {
		return cli.NewExitError(err.Error(), 127)
	}

	tmResult := new(agtypes.ResultBroadcastTxCommit)
	clientJSON := cl.NewClientJSONRPC(logger, commons.QueryServer)
	_, err = clientJSON.Call("broadcast_tx_commit", []interface{}{agtypes.WrapTx(types.TxTagAppEvmCommon, b)}, tmResult)
	if err != nil {
		panic(err)
	}
	//res := (*tmResult).(*types.ResultBroadcastTxCommit)

	fmt.Printf("tx result: %x\n", tools.Hash(tx))

	return nil
}
