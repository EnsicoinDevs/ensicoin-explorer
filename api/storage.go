package main

import (
	"encoding/binary"
	"encoding/json"
	bolt "go.etcd.io/bbolt"
)

const (
	BLOCKS_BUCKET  = "blocks"
	TXS_BUCKET     = "txs"
	STATS_BUCKET   = "stats"
	HEIGHTS_BUCKET = "heights"
)

type NotFoundError struct {
}

func (err *NotFoundError) Error() string {
	return "not found"
}

func dbOpen(path string) (*bolt.DB, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BLOCKS_BUCKET))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(TXS_BUCKET))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(STATS_BUCKET))
		if err != nil {

		}

		_, err = tx.CreateBucketIfNotExists([]byte(HEIGHTS_BUCKET))

		return err
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func dbStoreBestBlockHash(db *bolt.DB, hash string) error {
	return db.Update(func(tx *bolt.Tx) error {
		statsBucket := tx.Bucket([]byte(STATS_BUCKET))

		return statsBucket.Put([]byte("bestBlockHash"), []byte(hash))
	})
}

func dbGetBestBlockHash(db *bolt.DB) (hash string, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		statsBucket := tx.Bucket([]byte(STATS_BUCKET))

		hashBytes := statsBucket.Get([]byte("bestBlockHash"))
		if hashBytes == nil {
			return new(NotFoundError)
		}

		hash = string(hashBytes)

		return nil
	})

	return hash, err
}

func dbStoreBlock(db *bolt.DB, block *Block) error {
	err := db.Update(func(tx *bolt.Tx) error {
		blocksBucket := tx.Bucket([]byte(BLOCKS_BUCKET))
		txsBucket := tx.Bucket([]byte(TXS_BUCKET))

		liteBlock := &LiteBlock{
			Header: block.Header,
		}

		for _, tx := range block.Txs {
			liteBlock.Txs = append(liteBlock.Txs, tx.Hash)

			txBytes, err := json.Marshal(tx)
			if err != nil {
				return err
			}

			err = txsBucket.Put([]byte(tx.Hash), txBytes)
			if err != nil {
				return err
			}
		}

		liteBlockBytes, err := json.Marshal(liteBlock)
		if err != nil {
			return err
		}

		return blocksBucket.Put([]byte(block.Header.Hash), liteBlockBytes)
	})
	if err != nil {
		return err
	}

	return dbStoreHeight(db, block.Header.Height, block.Header.Hash)
}

func dbGetTx(db *bolt.DB, hash string) (*Tx, error) {
	tx := new(Tx)

	err := db.View(func(btx *bolt.Tx) error {
		txsBucket := btx.Bucket([]byte(TXS_BUCKET))

		txBytes := txsBucket.Get([]byte(hash))
		if txBytes == nil {
			return new(NotFoundError)
		}

		return json.Unmarshal(txBytes, &tx)
	})
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func dbGetBlock(db *bolt.DB, hash string) (*Block, error) {
	block := new(Block)

	liteBlock, err := dbGetLiteBlock(db, hash)
	if err != nil {
		return nil, err
	}

	block.Header = liteBlock.Header

	for _, txHash := range liteBlock.Txs {
		tx, err := dbGetTx(db, txHash)
		if err != nil {
			return nil, err
		}

		block.Txs = append(block.Txs, tx)
	}

	return block, nil
}

func dbGetLiteBlock(db *bolt.DB, hash string) (*LiteBlock, error) {
	liteBlock := new(LiteBlock)

	err := db.View(func(tx *bolt.Tx) error {
		blocksBucket := tx.Bucket([]byte(BLOCKS_BUCKET))

		liteBlockBytes := blocksBucket.Get([]byte(hash))
		if liteBlockBytes == nil {
			return new(NotFoundError)
		}

		return json.Unmarshal(liteBlockBytes, &liteBlock)
	})
	if err != nil {
		return nil, err
	}

	return liteBlock, nil
}

func dbHasBlock(db *bolt.DB, hash string) (bool, error) {
	_, err := dbGetLiteBlock(db, hash)
	if _, ok := err.(*NotFoundError); ok {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func dbStoreHeight(db *bolt.DB, height uint32, hash string) error {
	return db.Update(func(tx *bolt.Tx) error {
		heightsBucket := tx.Bucket([]byte(HEIGHTS_BUCKET))

		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, height)

		return heightsBucket.Put(b, []byte(hash))
	})
}

func dbGetHeight(db *bolt.DB, height uint32) (hash string, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		heightsBucket := tx.Bucket([]byte(HEIGHTS_BUCKET))

		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, height)

		hashBytes := heightsBucket.Get(b)
		if hashBytes == nil {
			return new(NotFoundError)
		}

		hash = string(hashBytes)

		return nil
	})

	return
}

func dbClose(db *bolt.DB) error {
	return db.Close()
}
