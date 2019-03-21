package bolt

import (
	"errors"
	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type boltClient struct {
	dbPath string

	fileMode os.FileMode
}

func NewBoltClient(dbPath string, fileMode os.FileMode) *boltClient {
	return &boltClient{dbPath: dbPath, fileMode: fileMode}
}

func (client boltClient) getDB() (*bolt.DB, error) {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(client.dbPath, client.fileMode, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		log.Errorf("Open bolt db error. %v", err)
		return nil, err
	}
	return db, nil
}

func (client boltClient) Set(root string, key string, value []byte) error {
	db, dbErr := client.getDB()
	if dbErr != nil {
		log.Errorf("Get bolt db error. %v", dbErr)
		return dbErr
	}
	defer db.Close()

	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(root))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), value)
	})
	if err != nil {
		return err
	}
	return nil
}

func (client boltClient) NextSequence(root string) (uint64, error) {
	db, dbErr := client.getDB()
	if dbErr != nil {
		log.Errorf("Get bolt db error. %v", dbErr)
		return 0, dbErr
	}
	defer db.Close()
	tx, err := db.Begin(true)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	bucketSequence, err := tx.Bucket([]byte(root)).NextSequence()
	if err != nil {
		return 0, err
	}
	return bucketSequence, nil
}

func (client boltClient) Get(root, key string) (string, error) {
	db, dbErr := client.getDB()
	if dbErr != nil {
		log.Errorf("Get bolt db error. %v", dbErr)
		return "", dbErr
	}
	defer db.Close()

	var result []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(root))
		if b != nil {
			v := b.Get([]byte(key))
			result = v
			return nil
		}
		return errors.New("Don't have this data. root : " + root + "; key : " + key)
	})
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return string(result), nil
}
