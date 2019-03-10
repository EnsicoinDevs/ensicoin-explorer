package main

import (
	"context"
	pb "github.com/EnsicoinDevs/ensicoin-explorer/api/rpc"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc"
	"os"
	"strconv"
	"time"
)

func main() {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(os.Getenv("API_NODE"), opts...)
	if err != nil {
		log.WithError(err).Fatal("unable to dial with the grpc server")
	}
	defer conn.Close()

	db, err := dbOpen("db/ensicoin_explorer.db")
	if err != nil {
		log.WithError(err).Fatal("unable to open the db")
	}

	client := pb.NewNodeClient(conn)

	r := gin.Default()

	r.GET("/api/blocks/:hash", func(c *gin.Context) {
		hash := c.Param("hash")

		block, err := dbGetBlock(db, hash)
		if _, ok := err.(*NotFoundError); ok {
			c.AbortWithStatus(404)
		} else if err != nil {
			c.AbortWithError(500, err)
		}

		c.JSON(200, gin.H{
			"block": block,
		})
	})

	r.GET("/api/blocks", func(c *gin.Context) {
		rawPage := c.DefaultQuery("page", "0")
		rawLimit := c.DefaultQuery("limit", "25")

		pageInt, _ := strconv.Atoi(rawPage)
		limitInt, _ := strconv.Atoi(rawLimit)
		page := uint32(pageInt)
		limit := uint32(limitInt)

		bestBlockHash, err := dbGetBestBlockHash(db)
		if _, ok := err.(*NotFoundError); ok {
			c.JSON(200, gin.H{
				"blocks": []*Block{},
			})

			return
		} else if err != nil {
			log.WithError(err).Fatal("unable to get the best block hash from the db")
		}

		bestBlock, _ := dbGetLiteBlock(db, bestBlockHash)

		var blocks []*LiteBlock

		for i := uint32(0); i < limit; i++ {
			height := bestBlock.Header.Height - page*limit - i

			log.WithField("height", height).Info(">")

			hash, err := dbGetHeight(db, height)

			liteBlock, err := dbGetLiteBlock(db, hash)
			if err != nil {
				log.WithError(err).Fatal("unable to get a lite block from the db")
			}

			blocks = append(blocks, liteBlock)

			if height == 0 {
				break
			}
		}

		c.JSON(200, gin.H{
			"stats": gin.H{
				"bestHeight": bestBlock.Header.Height,
			},
			"blocks": blocks,
		})
	})

	quit := startBlocksTicker(client, db)

	_ = quit

	r.Run()

	dbClose(db)
}

func startBlocksTicker(client pb.NodeClient, db *bolt.DB) chan struct{} {
	ticker := time.NewTicker(30 * time.Second)

	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				getNewBlocks(client, db)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	getNewBlocks(client, db)

	return quit
}

func getNewBlocks(client pb.NodeClient, db *bolt.DB) {
	currentBestBlockHash, err := dbGetBestBlockHash(db)
	if _, ok := err.(*NotFoundError); !ok && err != nil {
		log.WithError(err).Fatal("unable to get the current best block hash from the db")
	}

	bestBlockHashReply := rpcGetBestBlockHash(client, &pb.GetBestBlockHashRequest{})
	bestBlockHash := bestBlockHashReply.GetHash()

	if bestBlockHash != currentBestBlockHash {
		log.Info("the chain has changed")

		nextBlockHash := bestBlockHash
		hasBlock := false

		for !hasBlock {
			block := rpcGetAndConvertBlock(client, nextBlockHash)

			err = dbStoreBlock(db, block)
			if err != nil {
				log.WithError(err).Fatal("unable to store a block in the db")
			}

			nextBlockHash = block.Header.HashPrevBlock
			if nextBlockHash == "4242424242424242424242424242424242424242424242424242424242424242" {
				break
			}

			hasBlock, err = dbHasBlock(db, block.Header.HashPrevBlock)
			if err != nil {
				log.WithError(err).Fatal("unable to check if the db know a block")
			}
		}

		err = dbStoreBestBlockHash(db, bestBlockHash)
		if err != nil {
			log.WithError(err).Fatal("unable to store the best block hash in the db")
		}
	}
}

func rpcGetAndConvertBlock(client pb.NodeClient, hash string) *Block {
	blockReply := rpcGetBlock(client, &pb.GetBlockRequest{
		Hash: hash,
	})

	block := blockReply.GetBlock()

	mBlock := &Block{
		Header: &BlockHeader{
			Hash:           block.GetHash(),
			Version:        block.GetVersion(),
			Flags:          block.GetFlags(),
			HashPrevBlock:  block.GetHashPrevBlock(),
			HashMerkleRoot: block.GetHashMerkleRoot(),
			Timestamp:      block.GetTimestamp(),
			Height:         block.GetHeight(),
			Bits:           block.GetBits(),
			Nonce:          block.GetNonce(),
		},
	}

	for _, tx := range block.GetTxs() {
		mTx := &Tx{
			Hash:    tx.GetHash(),
			Version: tx.GetVersion(),
			Flags:   tx.GetFlags(),
		}

		for _, input := range tx.GetInputs() {
			mTx.Inputs = append(mTx.Inputs, &Input{
				PreviousOutput: &Outpoint{
					Hash:  input.GetPreviousOutput().GetHash(),
					Index: input.GetPreviousOutput().GetIndex(),
				},
				Script: input.GetScript(),
			})
		}

		for _, output := range tx.GetOutputs() {
			mTx.Outputs = append(mTx.Outputs, &Output{
				Value:  output.GetValue(),
				Script: output.GetScript(),
			})
		}

		mBlock.Txs = append(mBlock.Txs, mTx)
	}

	return mBlock
}

func rpcGetBlock(client pb.NodeClient, request *pb.GetBlockRequest) *pb.GetBlockReply {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	reply, err := client.GetBlock(ctx, request)
	if err != nil {
		log.WithError(err).Fatal("unable to get a block by hash")
	}

	return reply
}

func rpcGetBestBlockHash(client pb.NodeClient, request *pb.GetBestBlockHashRequest) *pb.GetBestBlockHashReply {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	reply, err := client.GetBestBlockHash(ctx, request)
	if err != nil {
		log.WithError(err).Fatal("unable to get the best block hash")
	}

	return reply
}
