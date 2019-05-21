package main

import (
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"strconv"
)

var (
	statsBucket         = []byte("stats")
	blocksBucket        = []byte("blocks")
	txsBucket           = []byte("txs")
	blockToTxsBucket    = []byte("blockToTxs")
	heightToBlockBucket = []byte("heightToBlock")
)

type Storage struct {
	path string
	db   *bolt.DB
}

func NewStorage(path string) *Storage {
	return &Storage{
		path: path,
	}
}

func (s *Storage) Open() (err error) {
	s.db, err = bolt.Open(s.path, 0600, nil)
	if err != nil {
		return err
	}

	return s.bootstrap()
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) bootstrap() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(statsBucket); err != nil {
			return err
		}

		if _, err := tx.CreateBucketIfNotExists(blocksBucket); err != nil {
			return err
		}

		if _, err := tx.CreateBucketIfNotExists(txsBucket); err != nil {
			return err
		}

		if _, err := tx.CreateBucketIfNotExists(blockToTxsBucket); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(heightToBlockBucket); err != nil {
			return err
		}

		return nil
	})
}

func (s *Storage) StoreBestBlockHash(hash string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(statsBucket).Put([]byte("bestBlockHash"), []byte(hash))
	})
}

func (s *Storage) FindBestBlockHash() (string, error) {
	var hash string

	if err := s.db.View(func(tx *bolt.Tx) error {
		bestBlockHashBytes := tx.Bucket(statsBucket).Get([]byte("bestBlockHash"))
		if bestBlockHashBytes == nil {
			return fmt.Errorf("best block hash not found")
		}

		hash = string(bestBlockHashBytes)

		return nil
	}); err != nil {
		return "", err
	}

	return hash, nil
}

func (s *Storage) StoreBlock(block *Block) error {
	blockBytes, err := json.Marshal(block)
	if err != nil {
		return err
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		if err := tx.Bucket(blocksBucket).Put([]byte(block.Hash), blockBytes); err != nil {
			return err
		}

		return tx.Bucket(heightToBlockBucket).Put([]byte(strconv.Itoa(int(block.Height))), []byte(block.Hash))
	})
}

func (s *Storage) HasBlock(hash string) (bool, error) {
	var exist bool

	if err := s.db.View(func(tx *bolt.Tx) error {
		exist = tx.Bucket(blocksBucket).Get([]byte(hash)) != nil

		return nil
	}); err != nil {
		return false, err
	}

	return exist, nil
}

func (s *Storage) FindBlockByHash(hash string) (block *Block, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		blockBytes := tx.Bucket(blocksBucket).Get([]byte(hash))
		if blockBytes == nil {
			return fmt.Errorf("block not found")
		}

		return json.Unmarshal(blockBytes, &block)
	})

	return
}

func (s *Storage) FindBlockByHeight(height uint32) (block *Block, err error) {
	var blockHash string

	if err := s.db.View(func(tx *bolt.Tx) error {
		blockHashBytes := tx.Bucket(heightToBlockBucket).Get([]byte(strconv.Itoa(int(height))))
		if blockHashBytes == nil {
			return fmt.Errorf("block not found")
		}

		blockHash = string(blockHashBytes)

		return nil
	}); err != nil {
		return nil, err
	}

	return s.FindBlockByHash(blockHash)
}

func (s *Storage) StoreTx(tx *Tx) error {
	txBytes, err := json.Marshal(tx)
	if err != nil {
		return err
	}

	return s.db.Update(func(btx *bolt.Tx) error {
		return btx.Bucket(txsBucket).Put([]byte(tx.Hash), txBytes)
	})
}

func (s *Storage) StoreTxs(blockHash string, txs []*Tx) error {
	var txHashes []string

	for _, tx := range txs {
		if err := s.StoreTx(tx); err != nil {
			return err
		}

		txHashes = append(txHashes, tx.Hash)
	}

	txHashesBytes, err := json.Marshal(txHashes)
	if err != nil {
		return err
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(blockToTxsBucket).Put([]byte(blockHash), txHashesBytes)
	})
}

func (s *Storage) FindTxs(blockHash string) ([]*Tx, error) {
	var txHashes []string

	if err := s.db.View(func(tx *bolt.Tx) error {
		txHashesBytes := tx.Bucket(blockToTxsBucket).Get([]byte(blockHash))
		if txHashesBytes == nil {
			return fmt.Errorf("block not found")
		}

		return json.Unmarshal(txHashesBytes, &txHashes)
	}); err != nil {
		return nil, err
	}

	var txs []*Tx

	for _, hash := range txHashes {
		tx, err := s.FindTxByHash(hash)
		if err != nil {
			return nil, err
		}

		txs = append(txs, tx)
	}

	return txs, nil
}

func (s *Storage) FindTxByHash(hash string) (tx *Tx, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		txBytes := tx.Bucket(txsBucket).Get([]byte(hash))
		if txBytes == nil {
			return fmt.Errorf("tx not found")
		}

		return json.Unmarshal(txBytes, &tx)
	})

	return
}
