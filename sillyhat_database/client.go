package sillyhat_database

import (
	log "sillyhat-golang-tool/sillyhat_log/logrus"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

const (
	connMaxLifetime = 10
	maxIdleConns = 10
	maxOpenConns = 50
)

type MySQLClient struct {
	dataSourceName    string
	//MaxIdle int
	//MaxOpen int
	//User    string
	//Pwd     string
	//DB      string
	//Port    int
	db    *sql.DB
}

func NewClient(dataSourceName string) (*MySQLClient,error) {
	db,err := sql.Open("mysql", dataSourceName)
	if err != nil{
		return nil,err
	}
	//使用前 Ping, 确保 DB 连接正常
	err = db.Ping()
	if err != nil {
		log.Error("Mysql client ping error.",err)
		return nil,err
	}
	//mysqlClient.pool.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	return &MySQLClient{dataSourceName:dataSourceName,db:db},nil
}

func (mysqlClient *MySQLClient) Close() {
	mysqlClient.db.Close()
}


func (mysqlClient *MySQLClient) getConnection() *sql.DB {
	return mysqlClient.db
}

//return auto created primary key
func (mysqlClient *MySQLClient) Insert(sql string,args ...interface{}) (int64,error) {
	stm,err := mysqlClient.getConnection().Prepare(sql)
	if err != nil {
		log.Error("Mysql client get connection error.",err)
		return 0,err
	}
	result,err := stm.Exec(args...)
	stm.Close()
	if err != nil {
		log.Error("Inser data error.",err)
		return 0,err
	}
	return result.LastInsertId()
}

//return affected count
func (mysqlClient *MySQLClient) Update(sql string,args ...interface{}) (int64,error) {
	stm,err := mysqlClient.getConnection().Prepare(sql)
	if err != nil {
		log.Error("Mysql client get connection error.",err)
		return 0,err
	}
	result,err := stm.Exec(args...)
	stm.Close()
	if err != nil {
		log.Error("Update data error.",err)
		return 0,err
	}
	return result.RowsAffected()
}

//return affected count
func (mysqlClient *MySQLClient) Delete(sql string,args ...interface{}) (int64,error) {
	stm,err := mysqlClient.getConnection().Prepare(sql)
	if err != nil {
		log.Error("Mysql client get connection error.",err)
		return 0,err
	}
	result,err := stm.Exec(args...)
	stm.Close()
	if err != nil {
		log.Error("Delete data error.",err)
		return 0,err
	}
	return result.RowsAffected()
}

func (mysqlClient *MySQLClient) DeleteByPrimaryKey(sql string,id int64) (int64,error) {
	stm,err := mysqlClient.getConnection().Prepare(sql)
	if err != nil {
		log.Error("Mysql client get connection error.",err)
		return 0,err
	}
	result,err := stm.Exec(id)
	stm.Close()
	if err != nil {
		log.Error("DeleteByPrimaryKey error.",err)
		return 0,err
	}
	return result.RowsAffected()
}

func (mysqlClient *MySQLClient) Count(sql string) (int,error) {
	tx,err := mysqlClient.getConnection().Begin()
	if err != nil {
		log.Error("Mysql client get connection error.",err)
		return 0,err
	}
	defer tx.Commit()
	var count int
	countErr := tx.QueryRow(sql).Scan(&count)
	if countErr != nil{
		log.Error("Query count error.",err)
		return 0,err
	}
	return count,nil
}

func (mysqlClient *MySQLClient) QueryList(sql string) ([] map[string]interface{},error) {
	tx,err := mysqlClient.getConnection().Begin()
	if err != nil {
		log.Error("Mysql client get connection error.",err)
		return nil,err
	}
	defer tx.Commit()
	rows,err := tx.Query(sql)
	if err != nil {
		log.Error("Query error.",err)
		return nil,err
	}
	defer rows.Close()
	//读出查询出的列字段名
	columns,err := rows.Columns()
	if err != nil {
		log.Error("rows.Columns() error.",err)
		return nil,err
	}
	//values是每个列的值，这里获取到byte里
	values := make([][]byte, len(columns))
	//query.Scan的参数，因为每次查询出来的列是不定长的，用len(cols)定住当次查询的长度
	scans := make([]interface{}, len(columns))
	//让每一行数据都填充到[][]byte里面
	for i := range values {
		scans[i] = &values[i]
	}
	//最后得到的map
	var results [] map[string]interface{}
	for rows.Next() { //循环，让游标往下推
		if err := rows.Scan(scans...); err != nil { //query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
			return nil,err
		}

		row := make(map[string]interface{}) //每行数据

		for k, v := range values { //每行数据是放在values里面，现在把它挪到row里
			key := columns[k]
			row[key] = string(v)
		}
		results = append(results,row)
	}
	return results,nil
}

func (mysqlClient *MySQLClient) GetByPrimaryKey(sql string) (map[string]interface{},error) {
	tx,err := mysqlClient.getConnection().Begin()
	if err != nil {
		log.Error("Mysql client get connection error.",err)
		return nil,err
	}
	defer tx.Commit()
	rows,err := tx.Query(sql)
	if err != nil {
		return nil,err
	}
	defer rows.Close()
	//读出查询出的列字段名
	columns,err := rows.Columns()
	if err != nil {
		return nil,err
	}
	//values是每个列的值，这里获取到byte里
	values := make([][]byte, len(columns))
	//query.Scan的参数，因为每次查询出来的列是不定长的，用len(cols)定住当次查询的长度
	scans := make([]interface{}, len(columns))
	//让每一行数据都填充到[][]byte里面
	for i := range values {
		scans[i] = &values[i]
	}
	//最后得到的map
	for rows.Next() { //循环，让游标往下推
		if err := rows.Scan(scans...); err != nil { //query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
			return nil,err
		}
		row := make(map[string]interface{}) //每行数据
		for k, v := range values { //每行数据是放在values里面，现在把它挪到row里
			key := columns[k]
			row[key] = string(v)
		}
		return row,nil
	}
	return nil,nil
}

//func (mysqlClient *MySQLClient) getReflectType(input interface{},column string) reflect.Type{
//	inputType := reflect.TypeOf(input)
//	//inputValue := reflect.ValueOf(input)
//	for i := 0; i < inputType.NumField(); i++ {
//		field := inputType.Field(i)
//		//value := inputValue.Field(i).Interface()
//		//log.Printf("%s: %v = %v\n", field.Name, field.Type, value)
//		if column == field.Name {
//			return field.Type
//		}
//	}
//	return nil
//}
//
//func (mysqlClient *MySQLClient) getReflectKind(input interface{},column string) reflect.Kind{
//	inputType := reflect.TypeOf(input)
//	inputValue := reflect.ValueOf(input)
//	for i := 0; i < inputType.NumField(); i++ {
//		field := inputType.Field(i)
//		value := inputValue.Field(i)
//		//log.Printf("name: %v ;type: %v;value:%v\n", field.Name, field.Type, value)
//		if column == field.Name {
//			return value.Kind()
//		}
//	}
//	return reflect.String
//}

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


//func (mysqlClient *MySQLClient) QueryList(sql string,input interface{}) ([]interface{},error) {
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
//	var output []
//	for rows.Next(){
//		if err := rows.Scan(&input); err != nil {
//			return nil,err
//		}
//		output = append(output,input)
//	}
//	return output,nil
//}


//func (mysqlClient *MySQLClient) Query(sql string,input interface{}) ([] map[string]interface{},error) {
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
//	var result []map[string]interface{}
//	//		test
//	columns, _ := rows.Columns()
//	scanArgs := make([]interface{}, len(columns))
//	values := make([]interface{}, len(columns))
//	for i := range values {
//		scanArgs[i] = &values[i]
//	}
//
//	for rows.Next() {
//		//将行数据保存到record字典
//		err = rows.Scan(scanArgs...)
//		record := make(map[string]interface{})
//		for i, value := range values {
//			if value != nil {
//				fieldKind := mysqlClient.getReflectKind(input,columns[i])
//				switch fieldKind {
//				case reflect.Bool:
//					formatValue,err := strconv.ParseBool(string(value.([]byte)))
//					if err != nil {
//						return nil,err
//					}
//					record[columns[i]] = formatValue
//				case reflect.Int:
//					formatValue,err := strconv.Atoi(string(value.([]byte)))
//					if err != nil {
//						return nil,err
//					}
//					record[columns[i]] = formatValue
//				case reflect.Int32:
//					value,err := strconv.ParseInt(string(value.([]byte)),10,64)
//					if err != nil {
//						return nil,err
//					}
//					record[columns[i]] = value
//				case reflect.Int64:
//					value,err := strconv.ParseInt(string(value.([]byte)),10,64)
//					if err != nil {
//						return nil,err
//					}
//					record[columns[i]] = value
//				default:
//					record[columns[i]] = string(value.([]byte))
//				}
//				//if columns[i] == "id"{
//				//	id,_ := strconv.ParseInt(string(value.([]byte)),10,64)
//				//	record[columns[i]] = id
//				//}else{
//				//	record[columns[i]] = string(value.([]byte))
//				//}
//			}
//		}
//		result = append(result, record)
//		//log.Println(record)
//	}
//	resultJson,_ := json.Marshal(result)
//	log.Println(resultJson)
//	//columns, _ := rows.Columns()
//	//scanArgs := make([]interface{}, len(columns))
//	//values := make([]interface{}, len(columns))
//	//for i := range values {
//	//	scanArgs[i] = &values[i]
//	//}
//	//
//	//for rows.Next() {
//	//	//将行数据保存到record字典
//	//	err = rows.Scan(scanArgs...)
//	//	record := make(map[string]string)
//	//	for i, col := range values {
//	//		if col != nil {
//	//			record[columns[i]] = string(col.([]byte))
//	//		}
//	//	}
//	//	result = append(result, record)
//	//	//fmt.Println(record)
//	//}
//	//resultJson,_ := json.Marshal(result)
//	//log.Println(resultJson)
//
//
//
//	//for rows.Next(){
//	//	var name string
//	//	var id int
//	//	if err := rows.Scan(&id,&name); err != nil {
//	//		log.Fatal(err)
//	//	}
//	//	//fmt.Printf("name:%s ,id:is %d\n", name, id)
//	//}
//	return result,nil
//	//return result,nil
//}


//func (mysqlClient *MySQLClient) GetByPrimaryKey(sql string,input interface{}) (map[string]interface{},error) {
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
//	columns, _ := rows.Columns()
//	scanArgs := make([]interface{}, len(columns))
//	values := make([]interface{}, len(columns))
//	for i := range values {
//		scanArgs[i] = &values[i]
//	}
//
//	for rows.Next() {
//		//将行数据保存到record字典
//		err = rows.Scan(scanArgs...)
//		record := make(map[string]interface{})
//		for i, value := range values {
//			if value != nil {
//				fieldKind := mysqlClient.getReflectKind(input,columns[i])
//				switch fieldKind {
//				case reflect.Bool:
//					formatValue,err := strconv.ParseBool(string(value.([]byte)))
//					if err != nil {
//						return nil,err
//					}
//					record[columns[i]] = formatValue
//				case reflect.Int:
//					formatValue,err := strconv.Atoi(string(value.([]byte)))
//					if err != nil {
//						return nil,err
//					}
//					record[columns[i]] = formatValue
//				case reflect.Int32:
//					value,err := strconv.ParseInt(string(value.([]byte)),10,64)
//					if err != nil {
//						return nil,err
//					}
//					record[columns[i]] = value
//				case reflect.Int64:
//					value,err := strconv.ParseInt(string(value.([]byte)),10,64)
//					if err != nil {
//						return nil,err
//					}
//					record[columns[i]] = value
//				default:
//					record[columns[i]] = string(value.([]byte))
//				}
//			}
//		}
//		return record,nil
//	}
//	return nil,errors.New("cannot find the data")
//}