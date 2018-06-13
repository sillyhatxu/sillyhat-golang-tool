package sillyhat_database

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"reflect"
	"strconv"
	"errors"
	"log"
	"encoding/json"
)

const (
	connMaxLifetime = 10
	maxIdleConns = 10
	maxOpenConns = 50
)

type MySQLClient struct {
	DataSourceName    string
	//MaxIdle int
	//MaxOpen int
	//User    string
	//Pwd     string
	//DB      string
	//Port    int
	pool    *sql.DB
}

func (mysqlClient *MySQLClient) Init() (err error) {
	mysqlClient.pool,err = sql.Open("mysql", mysqlClient.DataSourceName)
	if err != nil{
		panic(err)
	}
	if err != nil {
		return err
	}
	//使用前 Ping, 确保 DB 连接正常
	err = mysqlClient.pool.Ping()
	if err != nil {
		return err
	}
	//mysqlClient.pool.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)
	mysqlClient.pool.SetMaxIdleConns(maxIdleConns)
	mysqlClient.pool.SetMaxOpenConns(maxOpenConns)
	return nil
}

func (mysqlClient *MySQLClient) GetConnection() *sql.DB {
	return mysqlClient.pool
}


//return auto created primary key
func (mysqlClient *MySQLClient) Insert(sql string,args ...interface{}) (int64,error) {
	stm,err := mysqlClient.pool.Prepare(sql)
	if err != nil {
		log.Println(err)
		return 0,err
	}
	result,err := stm.Exec(args...)
	stm.Close()
	if err != nil {
		log.Println(err)
		return 0,err
	}
	return result.LastInsertId()
}

//return affected count
func (mysqlClient *MySQLClient) Update(sql string,args ...interface{}) (int64,error) {
	stm,err := mysqlClient.pool.Prepare(sql)
	if err != nil {
		return 0,err
	}
	result,err := stm.Exec(args...)
	stm.Close()
	if err != nil {
		return 0,err
	}
	return result.RowsAffected()
}

//return affected count
func (mysqlClient *MySQLClient) Delete(sql string,args ...interface{}) (int64,error) {
	stm,err := mysqlClient.pool.Prepare(sql)
	if err != nil {
		return 0,err
	}
	result,err := stm.Exec(args...)
	stm.Close()
	if err != nil {
		return 0,err
	}
	return result.RowsAffected()
}

func (mysqlClient *MySQLClient) DeleteByPrimaryKey(sql string,id int64) (int64,error) {
	stm,err := mysqlClient.pool.Prepare(sql)
	if err != nil {
		return 0,err
	}
	result,err := stm.Exec(id)
	stm.Close()
	if err != nil {
		return 0,err
	}
	return result.RowsAffected()
}

//type Entity struct{
//
//}
//
//type Callback func(dbData []interface{}) ([]map[string]string,error)


//return list
//func Query(sql string,) ([] map[string]string,error) {

func (mysqlClient *MySQLClient) getReflectType(input interface{},column string) reflect.Type{
	inputType := reflect.TypeOf(input)
	//inputValue := reflect.ValueOf(input)
	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)
		//value := inputValue.Field(i).Interface()
		//log.Printf("%s: %v = %v\n", field.Name, field.Type, value)
		if column == field.Name {
			return field.Type
		}
	}
	return nil
}

func (mysqlClient *MySQLClient) getReflectKind(input interface{},column string) reflect.Kind{
	inputType := reflect.TypeOf(input)
	inputValue := reflect.ValueOf(input)
	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)
		value := inputValue.Field(i)
		//log.Printf("name: %v ;type: %v;value:%v\n", field.Name, field.Type, value)
		if column == field.Name {
			return value.Kind()
		}
	}
	return reflect.String
}

//func (mysqlClient *MySQLClient) QueryList(sql string,typ reflect.Type) ([]interface{},error) {
//	tx,err := mysqlClient.pool.Begin()
//	if err != nil {
//		return nil,err
//	}
//	defer tx.Commit()
//	rows,err := tx.Query(sql)
//	if err != nil {
//		return nil,err
//	}
//	defer rows.Close()
//
//	var slice []interface{}
//
//	for rows.Next(){
//		v := reflect.New(typ).Elem()
//		if err := json.Unmarshal(rows.Scan(&interface{}); err == nil {
//			slice = append(slice, v.Interface())
//		}
//
//		var id int64
//		var shop_item_id int64
//		var size string
//		if err := rows.Scan(&id,&shop_item_id,&size); err != nil {
//			log.Fatal(err)
//		}
//		shoppingbagArray = append(shoppingbagArray,*&shoppingbagDTO{shop_item_id:shop_item_id,id:id,size:size})
//	}
//
//	for _, hit := range r.Hits.Hits {
//		v := reflect.New(typ).Elem()
//		if hit.Source == nil {
//			slice = append(slice, v.Interface())
//			continue
//		}
//		if err := json.Unmarshal(*hit.Source, v.Addr().Interface()); err == nil {
//			slice = append(slice, v.Interface())
//		}
//	}
//	return slice
//}


func (mysqlClient *MySQLClient) QueryList(sql string) (*sql.Rows,error) {
	tx,err := mysqlClient.pool.Begin()
	if err != nil {
		return nil,err
	}
	defer tx.Commit()
	rows,err := tx.Query(sql)
	if err != nil {
		return nil,err
	}
	defer rows.Close()
	return rows,nil
}


func (mysqlClient *MySQLClient) Query(sql string,input interface{}) ([] map[string]interface{},error) {
	tx,err := mysqlClient.pool.Begin()
	if err != nil {
		return nil,err
	}
	defer tx.Commit()
	rows,err := tx.Query(sql)
	if err != nil {
		return nil,err
	}
	defer rows.Close()

	var result []map[string]interface{}
	//		test
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		record := make(map[string]interface{})
		for i, value := range values {
			if value != nil {
				fieldKind := mysqlClient.getReflectKind(input,columns[i])
				switch fieldKind {
				case reflect.Bool:
					formatValue,err := strconv.ParseBool(string(value.([]byte)))
					if err != nil {
						return nil,err
					}
					record[columns[i]] = formatValue
				case reflect.Int:
					formatValue,err := strconv.Atoi(string(value.([]byte)))
					if err != nil {
						return nil,err
					}
					record[columns[i]] = formatValue
				case reflect.Int32:
					value,err := strconv.ParseInt(string(value.([]byte)),10,64)
					if err != nil {
						return nil,err
					}
					record[columns[i]] = value
				case reflect.Int64:
					value,err := strconv.ParseInt(string(value.([]byte)),10,64)
					if err != nil {
						return nil,err
					}
					record[columns[i]] = value
				default:
					record[columns[i]] = string(value.([]byte))
				}
				//if columns[i] == "id"{
				//	id,_ := strconv.ParseInt(string(value.([]byte)),10,64)
				//	record[columns[i]] = id
				//}else{
				//	record[columns[i]] = string(value.([]byte))
				//}
			}
		}
		result = append(result, record)
		//log.Println(record)
	}
	resultJson,_ := json.Marshal(result)
	log.Println(resultJson)
	//columns, _ := rows.Columns()
	//scanArgs := make([]interface{}, len(columns))
	//values := make([]interface{}, len(columns))
	//for i := range values {
	//	scanArgs[i] = &values[i]
	//}
	//
	//for rows.Next() {
	//	//将行数据保存到record字典
	//	err = rows.Scan(scanArgs...)
	//	record := make(map[string]string)
	//	for i, col := range values {
	//		if col != nil {
	//			record[columns[i]] = string(col.([]byte))
	//		}
	//	}
	//	result = append(result, record)
	//	//fmt.Println(record)
	//}
	//resultJson,_ := json.Marshal(result)
	//log.Println(resultJson)



	//for rows.Next(){
	//	var name string
	//	var id int
	//	if err := rows.Scan(&id,&name); err != nil {
	//		log.Fatal(err)
	//	}
	//	//fmt.Printf("name:%s ,id:is %d\n", name, id)
	//}
	return result,nil
	//return result,nil
}


func (mysqlClient *MySQLClient) GetByPrimaryKey(sql string,input interface{}) (map[string]interface{},error) {
	tx,err := mysqlClient.pool.Begin()
	if err != nil {
		return nil,err
	}
	defer tx.Commit()
	rows,err := tx.Query(sql)
	if err != nil {
		return nil,err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		record := make(map[string]interface{})
		for i, value := range values {
			if value != nil {
				fieldKind := mysqlClient.getReflectKind(input,columns[i])
				switch fieldKind {
				case reflect.Bool:
					formatValue,err := strconv.ParseBool(string(value.([]byte)))
					if err != nil {
						return nil,err
					}
					record[columns[i]] = formatValue
				case reflect.Int:
					formatValue,err := strconv.Atoi(string(value.([]byte)))
					if err != nil {
						return nil,err
					}
					record[columns[i]] = formatValue
				case reflect.Int32:
					value,err := strconv.ParseInt(string(value.([]byte)),10,64)
					if err != nil {
						return nil,err
					}
					record[columns[i]] = value
				case reflect.Int64:
					value,err := strconv.ParseInt(string(value.([]byte)),10,64)
					if err != nil {
						return nil,err
					}
					record[columns[i]] = value
				default:
					record[columns[i]] = string(value.([]byte))
				}
			}
		}
		return record,nil
	}
	return nil,errors.New("cannot find the data")
}