package sillyhat_database

import (
	"database/sql"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
	"testing"
	"time"
)

type Userinfo struct {
	Id               int64     `mapstructure:"id"`
	Name             string    `mapstructure:"name"`
	Age              int       `mapstructure:"age"`
	IsDelete         bool      `mapstructure:"is_delete"`
	CreatedTime      time.Time `mapstructure:"created_date"`
	LastModifiedDate time.Time `mapstructure:"last_modified_date"`
}

const dataSourceName = `database_username:database_password@tcp(127.0.0.1:3306)/ocb_syncer`

func TestClientInsert(t *testing.T) {
	client, err := NewClient(dataSourceName)
	testSQL := "INSERT INTO userinfo (name, age,is_delete, created_date, last_modified_date) VALUES (?,?,?,now(),now())"
	assert.Nil(t, err)
	defer client.Close()
	result, err := client.BatchInsert(func(tx *sql.Tx) (int, error) {
		totalCount := 0
		for i := 1001; i <= 2000; i++ {
			_, err := tx.Exec(testSQL, "name-"+strconv.Itoa(i), 25, i%2 == 0)
			assert.Nil(t, err)
			totalCount++
		}
		return totalCount, nil
	})
	assert.Nil(t, err)
	log.Println(result)
}

func TestClientUpdate(t *testing.T) {
	client, err := NewClient(dataSourceName)
	assert.Nil(t, err)
	defer client.Close()
	count, err := client.Update("UPDATE ocb_syncer.userinfo SET name = ?, age = ?,last_modified_date = now() WHERE id = ?", "xushikuan", 29, 5)
	assert.Nil(t, err)
	log.Println(count)
}

func TestClientBatchInsert(t *testing.T) {
	client, err := NewClient(dataSourceName)
	assert.Nil(t, err)
	defer client.Close()
	count, err := client.Update("UPDATE ocb_syncer.userinfo SET name = ?, age = ?,last_modified_date = now() WHERE id = ?", "xushikuan", 29, 5)
	assert.Nil(t, err)
	log.Println(count)
}

func TestClientDelete(t *testing.T) {
	client, err := NewClient(dataSourceName)
	assert.Nil(t, err)
	defer client.Close()
	count, err := client.DeleteByPrimaryKey("DELETE FROM ocb_syncer.userinfo WHERE id = ?", 5)
	assert.Nil(t, err)
	log.Println(count)
}

func TestClientGet(t *testing.T) {
	client, err := NewClient(dataSourceName)
	assert.Nil(t, err)
	defer client.Close()
	result, err := client.GetByPrimaryKey("SELECT id,name,age,(is_delete = b'1') is_delete,created_date,last_modified_date FROM ocb_syncer.userinfo WHERE id = 10 LIMIT 20")
	assert.Nil(t, err)
	var user *Userinfo
	config := &mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeHookFunc("2006-01-02 15:04:05"),
		WeaklyTypedInput: true,
		Result:           &user,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}
	err = decoder.Decode(result)
	if err != nil {
		panic(err)
	}
	log.Print(strconv.FormatInt(user.Id, 10)+"     "+user.Name+"     "+strconv.Itoa(user.Age)+"    ", user.IsDelete, "    ", user.CreatedTime, "     ", user.LastModifiedDate)
}

func TestClientQuery(t *testing.T) {
	client, err := NewClient(dataSourceName)
	assert.Nil(t, err)
	defer client.Close()
	results, err := client.QueryList("SELECT id,name,age,(is_delete = b'1') is_delete,created_date,last_modified_date FROM ocb_syncer.userinfo LIMIT 20")
	assert.Nil(t, err)
	var userArray []Userinfo
	for _, result := range results {
		var user Userinfo
		config := &mapstructure.DecoderConfig{
			DecodeHook:       mapstructure.StringToTimeHookFunc("2006-01-02 15:04:05"),
			WeaklyTypedInput: true,
			Result:           &user,
		}

		decoder, err := mapstructure.NewDecoder(config)
		if err != nil {
			panic(err)
		}
		err = decoder.Decode(result)
		if err != nil {
			panic(err)
		}
		userArray = append(userArray, user)
	}
	log.Println(len(userArray))
	for _, user := range userArray {
		log.Print(strconv.FormatInt(user.Id, 10)+"     "+user.Name+"     "+strconv.Itoa(user.Age)+"    ", user.IsDelete, "    ", user.CreatedTime, "     ", user.LastModifiedDate)
	}
}
