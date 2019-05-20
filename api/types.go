package main

import (
	"encoding/hex"
	"github.com/EnsicoinDevs/eccd/utils"
	pb "github.com/EnsicoinDevs/ensicoin-explorer/api/rpc"
)

type Block struct {
	Hash       string   `json:"hash"`
	Version    uint32   `json:"version"`
	Flags      []string `json:"flags"`
	PrevBlock  string   `json:"prev_block"`
	MerkleRoot string   `json:"merkle_root"`
	Timestamp  uint64   `json:"timestamp"`
	Height     uint32   `json:"height"`
	Target     string   `json:"target"`
}

type Outpoint struct {
	Hash  string `json:"hash"`
	Index uint32 `json:"index"`
}

type TxInput struct {
	PreviousOutput *Outpoint `json:"previous_output"`
	Script         string    `json:"script"`
}

type TxOutput struct {
	Value  uint64 `json:"value"`
	Script string `json:"script"`
}

type Tx struct {
	Hash    string      `json:"hash"`
	Version uint32      `json:"version"`
	Flags   []string    `json:"flags"`
	Inputs  []*TxInput  `json:"inputs"`
	Outputs []*TxOutput `json:"outputs"`
}

func RpcBlockToBlock(rpcBlock *pb.Block) *Block {
	return &Block{
		Hash:       utils.NewHash(rpcBlock.GetHash()).String(),
		Version:    rpcBlock.GetVersion(),
		Flags:      rpcBlock.GetFlags(),
		PrevBlock:  utils.NewHash(rpcBlock.GetPrevBlock()).String(),
		MerkleRoot: utils.NewHash(rpcBlock.GetMerkleRoot()).String(),
		Timestamp:  rpcBlock.GetTimestamp(),
		Height:     rpcBlock.GetHeight(),
		Target:     utils.NewHash(rpcBlock.GetTarget()).String(),
	}
}

func RpcTxInputsToTxInputs(inputs []*pb.TxInput) []*TxInput {
	var txInputs []*TxInput

	for _, input := range inputs {
		txInputs = append(txInputs, &TxInput{
			PreviousOutput: &Outpoint{
				Hash:  utils.NewHash(input.GetPreviousOutput().GetHash()).String(),
				Index: input.GetPreviousOutput().GetIndex(),
			},
			Script: hex.EncodeToString(input.GetScript()),
		})
	}

	return txInputs
}

func RpcTxOutputsToTxOutputs(outputs []*pb.TxOutput) []*TxOutput {
	var txOutputs []*TxOutput

	for _, output := range outputs {
		txOutputs = append(txOutputs, &TxOutput{
			Value:  output.GetValue(),
			Script: hex.EncodeToString(output.GetScript()),
		})
	}

	return txOutputs
}

func RpcBlockToTxs(rpcBlock *pb.Block) []*Tx {
	var txs []*Tx

	for _, tx := range rpcBlock.GetTxs() {
		txs = append(txs, &Tx{
			Hash:    utils.NewHash(tx.GetHash()).String(),
			Version: tx.GetVersion(),
			Flags:   tx.GetFlags(),
			Inputs:  RpcTxInputsToTxInputs(tx.GetInputs()),
			Outputs: RpcTxOutputsToTxOutputs(tx.GetOutputs()),
		})
	}

	return txs
}
