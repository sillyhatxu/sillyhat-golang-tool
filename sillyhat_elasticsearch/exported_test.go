package sillyhat_elasticsearch

import (
	"testing"
	"log"
)

func TestGetOrderGroupContent(t *testing.T) {
	var elasticClient = New()
	elasticClient.URL = "http://172.28.2.22:9200"
	elasticClient.ElasticIndex = "deja_products"
	elasticClient.ElasticType = "tags"
	idArray := []string{"1","5572483","5572484","5572485","5572443","5572447","5572426","5572428","5572435","5572439","5572414","5572416","5572420","5572386","5572356","5572365","5572341","5572344","5572345","5572346","5572328","5572333","5572308","5572280","5572281","5572284","5572285","5572291","5572292","5572293","5572294","5572295","5572264","5572268","5572269","5572270","5572246","5572247","5572250","5572251","5572252","5572253","5572255","5572227","5572235","5572242","5572208","5572209","5572210","5572211","5572212","5572213","5572214","5572215","5572216","5572217","5572218","5572198","5572199","5572200","5572201","5572202","5572203","5572204","5572205","5572206","5572207","5572186","5572188","5572189","5572191","5572192","5572193","5572194","5572195","5572196","5572197","5572093","5572046","5571909","5571911","5571912","5571877","5571882","5571854","5571867","5571868","5571850","5571797","5571804","5571755","5571761","5571762","5571696","5571697","5571708","5571686","5571694","5571621","5571622","5571623","5571631","5571632","5571584","5571580","5571581","5571582","5571545","5571546","5571547","5571535","5571536","5571537","5571538","5571539","5571540","5571542","5571543","5571544","5571519","5571520","5571521","5571522","5571523","5571525","5571526","5571498","5571499","5571500","5571501","5571502","5571503","5571504","5571505","5571506","5571507","5571508","5571509","5571510","5571511","5571512","5571513","5571514","5571515","5571472","5571473","5571474","5571475","5571476","5571477","5571478","5571479","5571480","5571481","5571482","5571483","5571484","5571485","5571486","5571488","5571489","5571490","5571491","5571492","5571493","5571494","5571495","5571496","5571497","5571439","5571440","5571441","5571442","5571443","5571444","5571445","5571446","5571447","5571449","5571452","5571453","5571456","5571459","5571461","5571462","5571463","5571464","5571466","5571469","5571470","5571471","5571422","5571423","5571424","5571425","5571426","5571427","5571428","5571429","5571430","5571431","5571432","5571433","5571434","5571436","5571437","5571438","5571403","5571404","5571406","5571408","5571409","5571410","5571411","5571412","5571413","5571415","5571417","5571418","5571420","5571387","5571388","5571389","5571390","5571391","5571393","5571394","5571395","5571396","5571397","5571398","5571399","5571400","5571402","5571370","5571371","5571372","5571373","5571374","5571375","5571376","5571377","5571378","5571379","5571380","5571381","5571382","5571383","5571384","5571385","5571353","5571356","5571357","5571358","5571359","5571360","5571361","5571362","5571363","5571364","5571365","5571366","5571367","5571368","5571369","5571340","5571341","5571343","5571344","5571345","5571348","5571349","5571351","5571352","5570031","5570038","5570046","5569996","5569998","5570005","5569984","5569989","5569990","5569991","5569994","5569995","5569957","5569961","5569973","5569975","5569976","5569978","5569942","5569881","5569882","5569883","5569884","5569885","5569886","5569887","5569889","5569890","5569892","5569867","5569868","5569869","5569870","5569871","5569872","5569875","5569876","5569877","5569878","5569852","5569853","5569854","5569855","5569856","5569857","5569858","5569859","5569860","5569862","5569863","5569837","5569838","5569839","5569841","5569842","5569843","5569844","5569846","5569848","5569849","5569850","5569818","5569820","5569821","5569823","5569803","5569804","5569805","5569806","5569807","5569808","5569809","5569810","5569811","5569812","5569813","5569815","5569816","5569792","5569793","5569795","5569796","5569797","5569798","5569799","5569801","5569802","5569789","5569790","5569791","5569775","5569776","5569777","5569778","5569779","5569780","5569781","5569782","5569783","5569784","5569785","5569786","5569787","5569761","5569762","5569763","5569765","5569766","5569767","5569768","5569769","5569770","5569771","5569735","5569744","5569745","5569746","5569747","5569749","5569751","5569717","5569720","5569722","5569723","5569724","5569726","5569728","5569730","5569695","5569696","5569697","5569698","5569699","5569700","5569701","5569705","5569706","5569707","5569709","5569711","5569712","5569714","5569679","5569686","5569687","5569688","5569689","5569690","5569691","5569692","5569693","5569694","5569663","5569671","5569672","5569677","5569678","5569644","5569645","5569646","5569647","5569648","5569649","5569650","5569652","5569653","5569655","5569657","5569658","5569659","5569660","5569628","5569629","5569631","5569632","5569633","5569634","5569636","5569637","5569638","5569639","5569641","5569642","5569643","5569612","5569621","5569622","5569623","5569624","5569625","5569627","5569597","5569598","5569604","5569605","5569608","5569609","5569610","5569611","5569579","5569580","5569581","5569583","5569584","5569585","5569590","5569591","5569593","5569567","5569568","5569569","5569570","5569571","5569572","5569574","5569575","5569576","5569577","5569578","5569561","5569562","5569563","5569565","5569566","5569536","5569535"}
	mgetResponse,err :=elasticClient.MultiGet(idArray)
	if err != nil{
		log.Panic(err)
	}
	log.Println(len(mgetResponse.Docs))
	log.Println(mgetResponse.Docs)
	//orderGroupContent,err := GetOrderGroupContent("12512397611785",78785)
	//assert.Nil(t, err)
	//assert.NotNil(t,orderGroupContent)
	//assert.Equal(t,1,len(orderGroupContent.ClientOrderContentList))
}