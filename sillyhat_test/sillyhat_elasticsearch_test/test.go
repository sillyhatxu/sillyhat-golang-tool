package main

import (
	elastic "sillyhat-golang-tool/sillyhat_elasticsearch"
	"log"
	"encoding/json"
	"sillyhat-golang-tool/sillyhat_database"
	"time"
	"strconv"
	"golang-cloud/tool/basic"
)

func init() {
	elastic.SetURL("http://172.28.2.22:9200")
	//elastic.SetIndex("deja_products")
	//elastic.SetType("tags")

	elastic.SetIndex("deja_products_detail")
	elastic.SetType("tags")

	//elastic.SetIndex("deja_products_image")
	//elastic.SetType("tags")

	//elastic.SetIndex("deja_products_inventory")
	//elastic.SetType("tags")
}

const (
	page_count = 100000
)

func main() {
	t1 := time.Now()
	//get("1")
	//get("5432777")
	//delete("5432777")
	deleteIndex()
	//get("5426451")
	//get("5447171")
	//bulkAll()
	//bulkAllDetail()
	//bulkAllImage()
	//bulkAllInventory()
	//indexExists()
	elapsed := time.Since(t1)
	log.Println("App elapsed: ", elapsed)
}

func get(id string) error {
	getResult,err := elastic.Get(id)
	if err != nil{
		log.Println(err.Error())
		return err
	}
	log.Println(getResult)
	if getResult.Found {
		log.Printf("Got document id[%v] version [%v] index [%v] type [%v]\n", getResult.Id, getResult.Version, getResult.Index, getResult.Type)
		resultJSON,err := getResult.Source.MarshalJSON()
		if err != nil{
			log.Println(err.Error())
			return err
		}
		var product Product
		json.Unmarshal([]byte(resultJSON), &product)
		log.Println(product)
	}
	return nil
}

func delete(id string) error {
	isDelete,err := elastic.Delete(id)
	if err != nil{
		log.Println(err.Error())
		return err
	}
	log.Printf("delete result [%v]",isDelete)
	return nil
}

func deleteIndex() error {
	isDelete,err := elastic.DeleteIndex()
	if err != nil{
		log.Println(err.Error())
		return err
	}
	log.Printf("delete result [%v]",isDelete)
	return nil
}

func index()  {

}

func bulk()  {
	dataSourceName := "deja_cloud:deja_cloud@tcp(deja-dt.ccf2gesv8s9h.ap-southeast-1.rds.amazonaws.com:3306)/shopping_bag"
	mysqlClient,err := sillyhat_database.NewClient(dataSourceName)
	if err != nil{
		panic(err)
	}
	productArray := queryProductArray(mysqlClient)
	log.Printf("productArray length : %v\n",len(productArray))

	var bulkEntityArray []elastic.BulkEntity
	//var docArray []Product
	//var idDeleteArray []string
	for _,product := range productArray{
		if !product.IsDelete && product.ValidateStatus{
			if product.CurrentPrice < product.OriginalPrice{
				product.IsDiscount = true
			}else{
				product.IsDiscount = false
			}
			if time.Now().Unix() - product.CreatedTime.Unix() <= 1209600{
				product.IsNewArrival = true
			}else{
				product.IsNewArrival = false
			}
			if product.Pattern != 0{
				product.ColorAndPattern = product.Pattern
			}else{
				product.ColorAndPattern = product.Color
			}
			product.AllSize = true
			product.Weight = 10000
			bulkEntityArray = append(bulkEntityArray,*&elastic.BulkEntity{Id:strconv.FormatInt(product.ProductId,10),Data:product,IsDelete:false})
		}else{
			bulkEntityArray = append(bulkEntityArray,*&elastic.BulkEntity{Id:strconv.FormatInt(product.ProductId,10),Data:product,IsDelete:true})
		}
	}
	response,err := elastic.Bulk(bulkEntityArray)
	if err != nil{
		log.Fatal(err.Error())
	}
	log.Println(response)
}

func bulkAll()  {
	dataSourceName := "deja_cloud:deja_cloud@tcp(deja-dt.ccf2gesv8s9h.ap-southeast-1.rds.amazonaws.com:3306)/shopping_bag"
	mysqlClient,err := sillyhat_database.NewClient(dataSourceName)
	if err != nil{
		panic(err)
	}
	productArray := queryProductArray(mysqlClient)
	log.Printf("productArray length : %v\n",len(productArray))
	var bulkEntityArray []elastic.BulkEntity
	for _,product := range productArray{
		if !product.IsDelete && product.ValidateStatus{
			if product.CurrentPrice < product.OriginalPrice{
				product.IsDiscount = true
			}else{
				product.IsDiscount = false
			}
			if time.Now().Unix() - product.CreatedTime.Unix() <= 1209600{
				product.IsNewArrival = true
			}else{
				product.IsNewArrival = false
			}
			if product.Pattern != 0{
				product.ColorAndPattern = product.Pattern
			}else{
				product.ColorAndPattern = product.Color
			}
			product.AllSize = true
			product.Weight = 10000
			bulkEntityArray = append(bulkEntityArray,*&elastic.BulkEntity{Id:strconv.FormatInt(product.ProductId,10),Data:product,IsDelete:false})
		}
	}

	log.Println("bulkEntityArray length : ",len(bulkEntityArray))

	totalRecord := len(bulkEntityArray)//total record
	totalPage := (totalRecord+page_count-1)/page_count;
	for i := 0;i < totalPage;i++{
		start := i*page_count
		end := basic.MinInt((i+1)*page_count,totalRecord)
		response,err := elastic.BulkAll(bulkEntityArray[start:end])
		if err != nil{
			log.Fatal(err.Error())
		}
		log.Println("success ",len(response.Succeeded()))
	}

}


func bulkAllDetail()  {
	dataSourceName := "deja_cloud:deja_cloud@tcp(deja-dt.ccf2gesv8s9h.ap-southeast-1.rds.amazonaws.com:3306)/shop"
	mysqlClient,err := sillyhat_database.NewClient(dataSourceName)
	if err != nil{
		panic(err)
	}
	productDetailArray := queryProductDetailArray(mysqlClient)
	log.Printf("productDetailArray length : %v\n",len(productDetailArray))
	var bulkEntityArray []elastic.BulkEntity
	for _,productDetail := range productDetailArray{
		bulkEntityArray = append(bulkEntityArray,*&elastic.BulkEntity{Id:strconv.FormatInt(productDetail.ProductId,10),Data:productDetail,IsDelete:false})
	}
	totalRecord := len(productDetailArray)//total record
	totalPage := (totalRecord+page_count-1)/page_count;
	for i := 0;i < totalPage;i++{
		start := i*page_count
		end := basic.MinInt((i+1)*page_count,totalRecord)
		log.Printf("start : %v ; end : %v\n",start,end)
		response,err := elastic.BulkAll(bulkEntityArray[start:end])
		if err != nil{
			log.Fatal(err.Error())
		}
		log.Println("success ",len(response.Succeeded()))
	}
}

func bulkAllImage()  {
	dataSourceName := "deja_cloud:deja_cloud@tcp(deja-dt.ccf2gesv8s9h.ap-southeast-1.rds.amazonaws.com:3306)/shop"
	mysqlClient,err := sillyhat_database.NewClient(dataSourceName)
	if err != nil{
		panic(err)
	}
	productImageArray := queryProductImageArray(mysqlClient)
	log.Printf("productImageArray length : %v\n",len(productImageArray))
	var bulkEntityArray []elastic.BulkEntity
	for _,productImage := range productImageArray{
		bulkEntityArray = append(bulkEntityArray,*&elastic.BulkEntity{Id:productImage.Id,Data:productImage,IsDelete:false})
	}
	totalRecord := len(productImageArray)//total record
	totalPage := (totalRecord+page_count-1)/page_count;
	for i := 0;i < totalPage;i++{
		start := i*page_count
		end := basic.MinInt((i+1)*page_count,totalRecord)
		log.Printf("start : %v ; end : %v\n",start,end)
		response,err := elastic.BulkAll(bulkEntityArray[start:end])
		if err != nil{
			log.Fatal(err.Error())
		}
		log.Println("success ",len(response.Succeeded()))
	}
}

func bulkAllInventory()  {
	dataSourceName := "deja_cloud:deja_cloud@tcp(deja-dt.ccf2gesv8s9h.ap-southeast-1.rds.amazonaws.com:3306)/inventory"
	mysqlClient,err := sillyhat_database.NewClient(dataSourceName)
	if err != nil{
		panic(err)
	}
	productInventoryArray := queryProductInventoryArray(mysqlClient)
	log.Printf("productInventoryArray length : %v\n",len(productInventoryArray))
	var bulkEntityArray []elastic.BulkEntity
	for _,productInventory := range productInventoryArray{
		bulkEntityArray = append(bulkEntityArray,*&elastic.BulkEntity{Id:strconv.FormatInt(productInventory.Id,10),Data:productInventory,IsDelete:false})
	}
	totalRecord := len(productInventoryArray)//total record
	totalPage := (totalRecord+page_count-1)/page_count;
	for i := 0;i < totalPage;i++{
		start := i*page_count
		end := basic.MinInt((i+1)*page_count,totalRecord)
		log.Printf("start : %v ; end : %v\n",start,end)
		response,err := elastic.BulkAll(bulkEntityArray[start:end])
		if err != nil{
			log.Fatal(err.Error())
		}
		log.Println("success ",len(response.Succeeded()))
	}
}

func indexExists()  {
	exists,err := elastic.IndexExists()
	if err != nil{
		log.Println(err.Error())
	}
	log.Println(exists)
}


type Product struct{

	ProductId int64 `json:"product_id"`

	ProductCode string `json:"product_code"`

	ProductGroupId string `json:"product_group_id"`

	ProductName string `json:"product_name"`

	BrandId int64 `json:"brand_id"`

	BrandName string `json:"brand_name"`

	ProductColor string `json:"product_color"`

	Category int `json:"category"`

	Subcategory int `json:"subcategory"`

	Color int `json:"color"`

	Pattern int `json:"pattern"`

	OCB bool `json:"ocb"`

	IsDiscount bool `json:"is_discount"`

	AllSize bool `json:"all_size"`

	IsNewArrival bool `json:"is_new_arrival"`

	IsPurchasable bool `json:"is_purchasable"`

	IsDelete bool `json:"is_delete"`

	IsRecommend bool `json:"is_recommend"`

	ImageUrl string `json:"image_url"`

	Height int64 `json:"height"`

	Width int64 `json:"width"`

	OriginalPrice int64 `json:"original_price"`

	CurrentPrice int64 `json:"current_price"`

	Weight int `json:"weight"`

	ColorAndPattern int `json:"color_and_pattern"`

	Currency string `json:"currency"`

	RecommendReason string `json:"recommend_reason"`

	ValidateStatus bool `json:"validate_status"`

	UpdatedTime time.Time `json:"update_time"`

	CreatedTime time.Time `json:"update_time"`

}

type ProductBulkDoc struct{

	elastic.BulkDoc

	insertArray [] Product

	updateArray [] Product

	deleteArray [] string

}

func toInterface(i interface{}) interface{} {
	return i
}

func (productBulkDoc *ProductBulkDoc) InsertArray() []interface{} {
	var result []interface{}
	for _,p := range productBulkDoc.insertArray{
		result = append(result,p)
	}
	return result
}

func (productBulkDoc *ProductBulkDoc) UpdateArray() []interface{} {
	var result []interface{}
	for _,p := range productBulkDoc.insertArray{
		result = append(result,toInterface(p))
	}
	return result
}

func (productBulkDoc *ProductBulkDoc) DeleteArray() [] string {
	var result [] string
	for _,p := range productBulkDoc.insertArray{
		result = append(result,strconv.FormatInt(p.ProductId,10))
	}
	return result
}

type ProductDetail struct{

	ProductId int64 `json:"product_id"`

	Description string `json:"description"`

	DetailDescription string `json:"detail_description"`

	SizeGuideTable string `json:"size_guide_table"`

	SizeGuideDescription string `json:"size_guide_description"`
}

type ProductImage struct{

	Id string `json:"id"`

	ProductId int64 `json:"product_id"`

	ImageURL string `json:"image_url"`

	Height int `json:"height"`

	Width int `json:"width"`

	IsDefault bool `json:"is_default"`

}

type ProductInventory struct{

	Id int64 `json:"id"`

	ProductId int64 `json:"product_id"`

	Size string `json:"size"`

	Quantity int `json:"quantity"`

	AutoDeleted bool `json:"auto_deleted"`

}

func queryProductInventoryArray(mysqlClient *sillyhat_database.MySQLClient) []ProductInventory {
	var resultArray []ProductInventory

	tx,err := mysqlClient.GetConnection().Begin()
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	defer tx.Commit()
	rows,err := tx.Query(product_inventory_sql)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	defer rows.Close()

	for rows.Next(){
		var productInventory = new(ProductInventory)
		if err := rows.Scan(&productInventory.Id,&productInventory.AutoDeleted,&productInventory.Quantity,&productInventory.Size,&productInventory.ProductId); err != nil {
			log.Fatal(err)
		}
		resultArray = append(resultArray,*productInventory)
	}
	log.Println("query end")
	return resultArray
}

func queryProductImageArray(mysqlClient *sillyhat_database.MySQLClient) []ProductImage {
	var resultArray []ProductImage

	tx,err := mysqlClient.GetConnection().Begin()
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	defer tx.Commit()
	rows,err := tx.Query(product_image_sql)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	defer rows.Close()

	for rows.Next(){
		var productImage = new(ProductImage)
		if err := rows.Scan(&productImage.Id,&productImage.ImageURL,&productImage.Height,&productImage.Width,&productImage.IsDefault,&productImage.ProductId); err != nil {
			log.Fatal(err)
		}
		resultArray = append(resultArray,*productImage)
	}
	log.Println("query end")
	return resultArray
}

func queryProductDetailArray(mysqlClient *sillyhat_database.MySQLClient) []ProductDetail {
	var productDetailArray []ProductDetail

	tx,err := mysqlClient.GetConnection().Begin()
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	defer tx.Commit()
	rows,err := tx.Query(product_detail_sql)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	defer rows.Close()

	for rows.Next(){
		var productDetail = new(ProductDetail)
		if err := rows.Scan(&productDetail.ProductId,&productDetail.Description,&productDetail.DetailDescription,&productDetail.SizeGuideTable,&productDetail.SizeGuideDescription); err != nil {
			log.Fatal(err)
		}
		productDetailArray = append(productDetailArray,*productDetail)
	}
	log.Println("query end")
	return productDetailArray
}

func queryProductArray(mysqlClient *sillyhat_database.MySQLClient) []Product {
	var productArray []Product

	tx,err := mysqlClient.GetConnection().Begin()
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	defer tx.Commit()
	rows,err := tx.Query(product_sql)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	defer rows.Close()

	for rows.Next(){
		var product = new(Product)
		var udatedDBTime string
		var createdDBTime string
		if err := rows.Scan(&product.ProductId,&product.ProductCode,&product.ProductName,&product.Category,&product.ProductColor,&product.Currency,&product.CurrentPrice,
				&product.OriginalPrice,&product.ProductGroupId,&product.BrandId,&product.BrandName,&product.IsPurchasable,&product.Width,
					&product.Height,&product.ImageUrl,&product.Subcategory,&product.Color,&product.Pattern,&product.OCB,&product.IsDelete,&product.ValidateStatus,&udatedDBTime,&createdDBTime); err != nil {
			log.Fatal(err)
		}
		DefaultTimeLoc := time.Local
		udatedTime, err := time.ParseInLocation("2006-01-02 15:04:05", udatedDBTime, DefaultTimeLoc)
		if err != nil {
			log.Println(err.Error())
		}
		createdTime, err := time.ParseInLocation("2006-01-02 15:04:05", createdDBTime, DefaultTimeLoc)
		if err != nil {
			log.Println(err.Error())
		}
		product.UpdatedTime = udatedTime
		product.CreatedTime = createdTime
		//log.Printf("today[%v]%v - createdTime[%v]%v",time.Now().Format("2006-01-02 15:04:05"),time.Now().Unix(),product.CreatedTime.Format("2006-01-02 15:04:05"),product.CreatedTime.Unix())
		//resultJson,_ := json.Marshal(product)
		//log.Println(string(resultJson))
		productArray = append(productArray,*product)
	}
	log.Println("query end")
	return productArray
}

const product_inventory_sql = `
SELECT
  id                    AS Id,
  (auto_deleted = b'1') AS AutoDeleted,
  quantity              AS Quantity,
  size                  AS Size,
  product_id            AS ProductId
FROM inventory.product_inventory
`

const product_image_sql = `
SELECT
  hash_id               AS Id,
  IFNULL(image_url, '') AS ImageURL,
  IFNULL(height, 0)     AS Height,
  IFNULL(width, 0)      AS Width,
  (is_default = b'1')   AS IsDefault,
  product_id            AS ProductId
FROM shop.product_image
`
const product_detail_sql = `
SELECT
  product_id                         AS ProductId,
  IFNULL(description, '')            AS Description,
  IFNULL(detail_description, '')     AS DetailDescription,
  IFNULL(size_guide_table, '')       AS SizeGuideTable,
  IFNULL(size_guide_description, '') AS SizeGuideDescription
FROM shop.product_detail
`
const product_sql = `
SELECT
  sit.product_id                AS productId,
  IFNULL(sit.product_code, '')  AS productCode,
  IFNULL(sit.name, '')          AS productName,
  sit.category                  AS category,
  IFNULL(sit.product_color, '') AS colorSrc,
  IFNULL(sit.currency, '')      AS currency,
  IFNULL(sit.current_price, 0)  AS currentPrice,
  IFNULL(sit.original_price, 0) AS originalPrice,
  IFNULL(sit.group_id, '')      AS productGroupId,
  IFNULL(sit.brand_id, 0)       AS brandId,
  IFNULL(s.name, '')            AS brandName,
  (sit.is_purchasable = b'1')   AS isPurchasable,
  IFNULL(sii.width, 0)          AS width,
  IFNULL(sii.height, 0)         AS height,
  IFNULL(sii.image_url, '')     AS height,
  sit.sub_category              AS subcategory,
  sit.color                     AS color,
  sit.pattern                   AS pattern,
  (sit.is_ocb = b'1')           AS ocb,
  (sit.is_delete = b'1')        AS isDelete,
  (sit.validate_status = b'1')  AS validateStatus,
  sit.last_modified_date        AS lastModifiedDate,  
  sit.created_date              AS createdDate
FROM shop.product sit
  LEFT JOIN shop.product_image sii ON sit.product_id = sii.product_id AND sii.is_default = TRUE
  LEFT JOIN shop.shop s ON sit.brand_id = s.id
ORDER BY sit.product_id
`
//WHERE sit.validate_status = TRUE AND sit.is_delete = FALSE AND sit.is_ocb = TRUE