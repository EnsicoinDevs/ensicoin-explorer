package main

import (
	"context"
	"github.com/EnsicoinDevs/eccd/utils"
	pb "github.com/EnsicoinDevs/ensicoin-explorer/api/rpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Synchronizer struct {
	rpcServerAddress string
	storage          *Storage
	conn             *grpc.ClientConn
	client           pb.NodeClient

	quit chan struct{}
}

func NewSynchronizer(storage *Storage, rpcServerAddress string) *Synchronizer {
	return &Synchronizer{
		rpcServerAddress: rpcServerAddress,
		storage:          storage,

		quit: make(chan struct{}),
	}
}

func (s *Synchronizer) Start() (err error) {
	s.conn, err = grpc.Dial(s.rpcServerAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}

	s.client = pb.NewNodeClient(s.conn)

	if err = s.startInitialSynchronization(); err != nil {
		return err
	}

	go s.startHandler()

	return nil
}

func (s *Synchronizer) Stop() error {
	close(s.quit)

	return s.conn.Close()
}

func (s *Synchronizer) startInitialSynchronization() error {
	bestBlockHash, _, err := s.GetStats()
	if err != nil {
		return err
	}

	if err = s.synchronizeTo(bestBlockHash); err != nil {
		return err
	}

	return nil
}

func (s *Synchronizer) startHandler() {
	stream, err := s.client.GetBestBlocks(context.Background(), &pb.GetBestBlocksRequest{})
	if err != nil {
		log.WithError(err).Fatal("fatal error getting the best blocks")
	}

	for {
		bestBlockHash, err := stream.Recv()
		if err != nil {
			log.Fatal("fatal error receiving a block")
		}

		return bestBlockHash
	}
}

func (s *Synchronizer) synchronizeTo(bestBlockHash string) error {
	log.WithField("bestBlockHash", bestBlockHash).Info("synchronizing up to")

	_, genesisBlockHash, err := s.GetStats()
	if err != nil {
		return err
	}

	currentHash := bestBlockHash

	for currentHash != genesisBlockHash {
		exist, err := s.storage.HasBlock(currentHash)
		if err != nil {
			return err
		}

		if exist {
			break
		}

		block, txs, err := s.FindBlockByHash(currentHash)
		if err != nil {
			return err
		}

		if err = s.storage.StoreBlock(block); err != nil {
			return err
		}

		if err = s.storage.StoreTxs(block.Hash, txs); err != nil {
			return err
		}

		currentHash = block.PrevBlock
	}

	bestBlock, _ := s.storage.FindBlockByHash(bestBlockHash)

	log.WithFields(log.Fields{
		"bestBlockHash": bestBlockHash,
		"height":        bestBlock.Height,
	}).Info("synchronized")

	return s.storage.StoreBestBlockHash(bestBlockHash)
}

func (s *Synchronizer) GetStats() (string, string, error) {
	info, err := s.client.GetInfo(context.Background(), &pb.GetInfoRequest{})
	if err != nil {
		return "", "", err
	}

	return utils.NewHash(info.GetBestBlockHash()).String(), utils.NewHash(info.GetGenesisBlockHash()).String(), nil
}

func (s *Synchronizer) FindBlockByHash(hash string) (*Block, []*Tx, error) {
	log.WithField("hash", hash).Info("downloading block")

	hashBytes, _ := utils.StringToHash(hash)

	rpcBlock, err := s.client.GetBlockByHash(context.Background(), &pb.GetBlockByHashRequest{
		Hash: hashBytes.Bytes(),
	})
	if err != nil {
		return nil, nil, err
	}

	return RpcBlockToBlock(rpcBlock.GetBlock()), RpcBlockToTxs(rpcBlock.GetBlock()), nil
}
