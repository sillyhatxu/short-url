package bolt

import (
	"errors"
	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

const ROOT_BUCKET = "ROOT_BUCKET"

type boltClient struct {
	dbPath string

	fileMode os.FileMode

	options *bolt.Options
}

func NewBoltClient(dbPath string, fileMode os.FileMode) *boltClient {
	return &boltClient{dbPath: dbPath, fileMode: fileMode, options: &bolt.Options{Timeout: 2 * time.Second}}
}

func (client boltClient) InitialBucket() error {
	db, err := bolt.Open(client.dbPath, client.fileMode, client.options)
	if err != nil {
		log.Errorf("Open bolt db error. %v", err)
		return err
	}
	defer db.Close()

	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.CreateBucketIfNotExists([]byte(ROOT_BUCKET))
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (client boltClient) ForEach() (map[string]string, error) {
	db, err := bolt.Open(client.dbPath, client.fileMode, client.options)
	if err != nil {
		log.Errorf("Open bolt db error. %v", err)
		return nil, err
	}
	defer db.Close()
	result := make(map[string]string)
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ROOT_BUCKET))
		err := b.ForEach(func(key, value []byte) error {
			result[string(key)] = string(value)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	return result, nil
}

func (client boltClient) Get(key string) (string, error) {
	db, err := bolt.Open(client.dbPath, client.fileMode, client.options)
	if err != nil {
		log.Errorf("Open bolt db error. %v", err)
		return "", err
	}
	defer db.Close()

	var result []byte
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(ROOT_BUCKET))
		if bucket == nil {
			return errors.New("Unknown bucket")
		}
		value := bucket.Get([]byte(key))
		result = value
		return nil
	})
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (client boltClient) Set(key, value string) error {
	db, err := bolt.Open(client.dbPath, client.fileMode, client.options)
	if err != nil {
		log.Errorf("Open bolt db error. %v", err)
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(ROOT_BUCKET))
		if bucket == nil {
			return errors.New("Unknown bucket")
		}
		return bucket.Put([]byte(key), []byte(value))
	})
	if err != nil {
		return err
	}
	return nil
}

func (client boltClient) NextSequence() (uint64, error) {
	db, err := bolt.Open(client.dbPath, client.fileMode, client.options)
	if err != nil {
		log.Errorf("Open bolt db error. %v", err)
		return 0, err
	}
	defer db.Close()

	tx, err := db.Begin(true)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	rootBucket := tx.Bucket([]byte(ROOT_BUCKET))
	if rootBucket == nil {
		return 0, errors.New("Unknown bucket")
	}
	bucketSequence, err := rootBucket.NextSequence()
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return bucketSequence, nil
}
