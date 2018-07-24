package sillyhat_enviroment

import (
	log "sillyhat-golang-tool/sillyhat_log/logrus"
	"os"
	"bufio"
	"strings"
	"io"
	"sillyhat-golang-tool/sillyhat_bolt"
	"fmt"
)
const (
	enviroment_root = "sillyhat_enviroment_root"
	module_name  = "module.name"
)

var sillyhatEnviromentClient *enviromentClient

type enviromentClient struct {

	enviromentName string

	moduleName string

	dbPath string

	fileMode os.FileMode

}

func (client enviromentClient) String() string {
	return fmt.Sprintf("%b", client)
}



func InitialEnviromentConfig(enviromentName,moduleName,dbPath string,fileMode os.FileMode){
	sillyhatEnviromentClient = &enviromentClient{enviromentName:enviromentName,moduleName:moduleName,dbPath:dbPath,fileMode:fileMode}
	sillyhatEnviromentClient.readEnviromentFile()
}

func Get(key string) string {
	log.Info("sillyhatEnviromentClient : ",sillyhatEnviromentClient.String())
	boltClient := sillyhat_bolt.NewBoltClient(sillyhatEnviromentClient.dbPath,sillyhatEnviromentClient.fileMode)
	result,err := boltClient.Get(enviroment_root,key)
	if err != nil{
		log.Error(err.Error())
		return ""
	}
	return result
}

//func (sillyhatEnviromentClient enviromentClient) InitialEnviromentConfig()  {
//	sillyhatEnviromentClient.readEnviromentFile()
//}

func (sillyhatEnviromentClient enviromentClient) readEnviromentFile(){
	boltClient := sillyhat_bolt.NewBoltClient(sillyhatEnviromentClient.dbPath,sillyhatEnviromentClient.fileMode)
	log.Info("--------------- config start ---------------")
	log.Infof("enviroment name : %v \n",sillyhatEnviromentClient.enviromentName)

	//打开这个ini文件
	f, _ := os.Open(sillyhatEnviromentClient.getConfigDir(sillyhatEnviromentClient.moduleName))
	//读取文件到buffer里边
	buf := bufio.NewReader(f)
	for {
		//按照换行读取每一行
		l, err := buf.ReadString('\n')
		//相当于PHP的trim
		line := strings.TrimSpace(l)
		//判断退出循环
		if err != nil {
			if err != io.EOF {
				//return err
				panic(err)
			}
			if len(line) == 0 {
				break
			}
		}
		switch {
		case len(line) == 0:
			//匹配[db]然后存储
		case line[0] == '[' && line[len(line)-1] == ']':
			section := strings.TrimSpace(line[1 : len(line)-1])
			log.Println(section)
		default:
			//dnusername = xiaowei 这种的可以匹配存储
			i := strings.IndexAny(line, "=")
			k := strings.TrimSpace(line[0:i])
			v := strings.TrimSpace(line[i+1:])
			boltClient.Set(enviroment_root,k,v)
			log.Info(k,":", v)
		}
	}
	boltClient.Set(enviroment_root,module_name,sillyhatEnviromentClient.moduleName)
	log.Info("--------------- config end ---------------")
}



func (sillyhatEnviromentClient enviromentClient) getConfigDir(moduleName string) string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	configDir := dir + "/" + moduleName + "/" + sillyhatEnviromentClient.getConfigName()
	log.Info("configDir : [%v]\n",configDir)
	return configDir
}

func (sillyhatEnviromentClient enviromentClient) getConfigName() string {
	return "config-" + sillyhatEnviromentClient.enviromentName + ".ini"
}

