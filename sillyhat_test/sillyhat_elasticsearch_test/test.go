package main

import (
	elastic "sillyhat-golang-tool/sillyhat_elasticsearch"
	"log"
	"encoding/json"
	"sillyhat-golang-tool/sillyhat_database"
	"time"
)

func init() {
	elastic.SetURL("http://172.28.2.22:9200")
	elastic.SetIndex("deja_products")
	elastic.SetType("tags")
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

func index()  {

}

func bulk()  {
	dataSourceName := "deja_cloud:deja_cloud@tcp(deja-dt.ccf2gesv8s9h.ap-southeast-1.rds.amazonaws.com:3306)/shopping_bag"
	var mysqlClient sillyhat_database.MySQLClient
	mysqlClient = sillyhat_database.MySQLClient{DataSourceName:dataSourceName}
	mysqlClient.Init()
	productArray := queryProductArray(mysqlClient)
	log.Printf("productArray length : ",len(productArray))
	var docArray []Product
	var idDeleteArray []string
	for _,product := range productArray{
		if product.IsDelete || !product.ValidateStatus{
			idDeleteArray = append(idDeleteArray,product.Id)
		}else{
			docArray = append(docArray,product)
		}
	}
	elastic.Bulk(docArray,idDeleteArray)
}



func main() {
	//get("1")
	//get("5432777")
	//delete("5432777")
	//get("5426451")
	//get("5447171")
	bulk()
}




type Product struct{

	elastic.BulkableDoc

	Id string `json:"id"`

	ProductId int64 `json:"product_id"`

	ProductCode string `json:"product_code"`

	ProductGroupId string `json:"product_group_id"`

	ProductName string `json:"product_name"`

	BrandId int64 `json:"brand_id"`

	BrandName string `json:"brand_name"`

	ColorSrc string `json:"colorSrc"`

	Category int `json:"category"`

	Subcategory int `json:"subcategory"`

	Color int `json:"color"`

	Pattern int `json:"pattern"`

	OCB bool `json:"ocb"`

	IsDiscount bool `json:"is_discount"`

	IsNewArrival bool `json:"is_new_arrival"`

	IsPurchasable bool `json:"is_purchasable"`

	IsDelete bool `json:"is_delete"`

	IsRecommend bool `json:"is_recommend"`

	ImageUrl string `json:"image_url"`

	Height int64 `json:"height"`

	Width int64 `json:"width"`

	OriginalPrice int64 `json:"original_price"`

	CurrentPrice int64 `json:"current_price"`

	Currency string `json:"currency"`

	RecommendReason string `json:"recommend_reason"`

	ValidateStatus bool `json:"validate_status"`

	UpdatedTime time.Time `json:"update_time"`

	CreatedTime time.Time `json:"update_time"`

}

type ProductDetail struct{

	Id string `json:"id"`

	Description int64 `json:"description"`

	DetailDescription string `json:"detail_description"`

	SizeGuideTable string `json:"size_guide_table"`

	SizeGuideDescription string `json:"size_guide_description"`


}

func queryProductArray(mysqlClient sillyhat_database.MySQLClient) []Product {
	var productArray []Product

	tx,err := mysqlClient.GetConnection().Begin()
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	defer tx.Commit()
	rows,err := tx.Query(sql)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	defer rows.Close()

	for rows.Next(){
		var product = new(Product)
		var udatedDBTime string
		var createdDBTime string
		if err := rows.Scan(&product.ProductId,&product.ProductCode,&product.ProductName,&product.Category,&product.ColorSrc,&product.Currency,&product.CurrentPrice,
				&product.OriginalPrice,&product.RecommendReason,&product.ProductGroupId,&product.BrandId,&product.BrandName,&product.IsPurchasable,&product.Width,
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

		productArray = append(productArray,*product)
	}
	log.Println("query end")
	return productArray
}

const sql = `
SELECT
  sit.product_id         AS productId,
  sit.product_code       AS productCode,
  sit.name               AS productName,
  sit.category           AS category,
  sit.color              AS colorSrc,
  sit.currency           AS currency,
  sit.current_price      AS currentPrice,
  sit.original_price     AS originalPrice,
  ''                     AS recommendReason,
  sit.group_id           AS productGroupId,
  sit.brand_id           AS brandId,
  s.name                 AS brandName,
  (sit.is_purchasable = b'1')  AS isPurchasable,
  sii.width              AS width,
  sii.height             AS height,
  sii.image_url          AS imageUrl,
  sit.sub_category       AS subcategory,
  sit.color              AS color,
  sit.pattern            AS pattern,
  (sit.is_ocb = b'1')    AS ocb,
  (sit.is_delete = b'1') AS isDelete,
  (sit.validate_status = b'1')  AS validateStatus,
  sit.created_date       AS createdDate,
  sit.last_modified_date AS lastModifiedDate
FROM shop.product sit
  LEFT JOIN shop.product_image sii ON sit.product_id = sii.product_id AND sii.is_default = TRUE
  INNER JOIN shop.shop s ON sit.brand_id = s.id
WHERE sit.validate_status = TRUE AND sit.is_delete = FALSE AND sit.is_ocb = true
ORDER BY sit.product_id limit 10
`