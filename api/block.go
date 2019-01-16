package main

type BlockHeader struct {
	Hash           string   `json:"hash"`
	Version        uint32   `json:"version"`
	Flags          []string `json:"flags"`
	HashPrevBlock  string   `json:"hash_prev_block"`
	HashMerkleRoot string   `json:"hash_merkle_root"`
	Timestamp      int64    `json:"timestamp"`
	Height         uint32   `json:"height"`
	Bits           uint32   `json:"bits"`
	Nonce          uint64   `json:"nonce"`
}

type Block struct {
	Header *BlockHeader `json:"header"`
	Txs    []*Tx        `json:"txs"`
}

type LiteBlock struct {
	Header *BlockHeader `json:"header"`
	Txs    []string     `json:"txs"`
}

type Tx struct {
	Hash    string    `json:"hash"`
	Version uint32    `json:"version"`
	Flags   []string  `json:"flags"`
	Inputs  []*Input  `json:"inputs"`
	Outputs []*Output `json:"outputs"`
}

type Input struct {
	PreviousOutput *Outpoint `json:"previousOutput"`
	Script         []byte    `json:"script"`
}

type Outpoint struct {
	Hash  string `json:"hash"`
	Index uint32 `json:"index"`
}

type Output struct {
	Value  uint64 `json:"value"`
	Script []byte `json:"script"`
}
