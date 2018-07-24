package sillyhat_bolt

import (
	log "sillyhat-golang-tool/sillyhat_log/logrus"
	"github.com/boltdb/bolt"
	"os"
	"errors"
)

type boltClient struct{

	dbPath string

	fileMode os.FileMode
}

func NewBoltClient(dbPath string,fileMode os.FileMode) *boltClient {
	return &boltClient{dbPath:dbPath,fileMode:fileMode}
}

func (client boltClient)getDB() (*bolt.DB, error) {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(client.dbPath, client.fileMode, nil)
	if err != nil {
		log.Error(err)
		return nil,err
	}
	return db,nil
}

func (client boltClient) Set(root,key,value string) (error) {
	db,dbErr := client.getDB()
	if dbErr != nil {
		log.Error(dbErr)
		return dbErr
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(root))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), []byte(value))
	})
	return nil
}

func (client boltClient) Get(root,key string) (string,error) {
	db,dbErr := client.getDB()
	if dbErr != nil {
		log.Error(dbErr)
		return "",dbErr
	}
	defer db.Close()

	var result []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(root))
		if b != nil{
			v := b.Get([]byte(key))
			result = v
			return nil
		}
		return errors.New("Don't have this data. root : " + root + "; key : "+ key)
	})
	if err != nil{
		log.Println(err.Error())
		return "",err
	}
	return string(result),nil
}