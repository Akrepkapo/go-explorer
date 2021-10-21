package models

import (
	"encoding/json"
	"fmt"
	"github.com/IBAX-io/go-explorer/conf"
	"github.com/IBAX-io/go-explorer/storage"
	crypto2 "github.com/IBAX-io/go-ibax/packages/crypto"
	"github.com/oschwald/geoip2-golang"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

type NodeInfo struct {
	ApiAddress string `json:"api_address"`
	PublicKey  string `json:"public_key"`
	TcpAddress string `json:"tcp_address"`
}
type AddressComponent struct {
	Nation   string `json:"nation"`
	Province string `json:"province"`
	City     string `json:"city"`
	Locality string `json:"locality"`
}
type AdInfo struct {
	Nation   string `json:"nation"`
	Province string `json:"province"`
}
type detailedAddressInfo struct {
	Address          string           `json:"address"`
	AddressComponent AddressComponent `json:"address_component"`
	AdInfo           AdInfo           `json:"ad_info"`
}
type addressInfo struct {
	Status    int                 `json:"status"`
	Message   string              `json:"message"`
	RequestId string              `json:"request_id"`
	Result    detailedAddressInfo `json:"result"`
}

type FullNodeInfo struct {
	ID        int64   `gorm:"primary_key;not null"`
	Number    int64   `json:"number" gorm:"not null"`
	Value     string  `json:"value,omitempty" gorm:"not null"`
	Address   string  `json:"address,omitempty" gorm:"not null"`
	Latitude  float64 `json:"latitude,omitempty" gorm:"not null"`
	Longitude float64 `json:"longitude,omitempty" gorm:"not null"`
	Display   bool    `json:"display,omitempty" gorm:"not null"`
}
type infoIpInform struct {
	CityName  string
	Latitude  float64
	Longitude float64
}
type NodeMapInfo struct {
	Number    int64   `json:"number"`
	Name      string  `json:"name,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

const defaultAppKey = "WJLBZ-NYUL5-5UXIP-QTM2C-EMDO3-3QFAH"

var (
	addrList      []string
	FullnodesInfo []*storage.FullnodeModel
)

func getFullNodeInfoFromDb() {
	var err error = nil
	var nodeInfo string
	var fullNode []storage.FullnodeModel

	if err = GetDB(nil).Table("1_system_parameters").Where("name = ?", "honor_nodes").Select("value").Take(&nodeInfo).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.WithFields(log.Fields{"error": err}).Error("getFullnodeInfo Select failed")
		}
	}
	if len(nodeInfo) == 0 {
		syncNodeDisplayStatus(fullNode)
		syncFullNodeInfoToRedis()
		SyncFullNodeInfo()

		return
	}
	err = json.Unmarshal([]byte(nodeInfo), &fullNode)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("getFullnodeInfo jsonUnmared failed")
		return
	}
	syncNodeDisplayStatus(fullNode)
	value, errMar := json.Marshal(fullNode)
	if errMar != nil {
		log.WithFields(log.Fields{"error": err}).Error("getFullnodeInfo marshal failed")
		return
	}
	if err = GetNodeListInfo(string(value)); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("GetNodeListInfo failed")
		return
	for i := 0; i < len(nodeValue); i++ {
		var fullNode storage.FullnodeModel
		err := json.Unmarshal([]byte(nodeValue[i]), &fullNode)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("fullNodeDbIsExist json failed")
			return false
		}
		if fullNode.APIAddress == apiaddress {
			return true
		}
	}
	return false
}
func syncNodeDisplayStatus(fullnode []storage.FullnodeModel) (statusDiff bool) {
	var node []FullNodeInfo
	if err := GetDB(nil).Table("fullnode_info").Order("id desc").Find(&node).Error; err != nil {
		log.WithFields(log.Fields{"error": err}).Error("sync display db failed")
		return false
	}

	nodeMd := make([]storage.FullnodeModel, len(node))
	for i := 0; i < len(node); i++ {
		if err := json.Unmarshal([]byte(node[i].Value), &nodeMd[i]); err != nil {
			log.WithFields(log.Fields{"error": err}).Error("syn display json failed")
			continue
		}

	}
	if len(fullnode) > 0 {
		for i := 0; i < len(node); i++ {
			if node[i].Display == false {
				for j := 0; j < len(fullnode); j++ {
					if nodeMd[i].APIAddress == fullnode[j].APIAddress {
						node[i].Display = true
						if err := GetDB(nil).Table("fullnode_info").Where("id = ?", node[i].ID).Update("display", node[i].Display).Error; err != nil {
							log.WithFields(log.Fields{"error": err}).Error("sync display status update1 err")
							continue
						}
						statusDiff = true
					}
				}
			} else {
				statusIstrue := false
				for j := 0; j < len(fullnode); j++ {
					if nodeMd[i].APIAddress == fullnode[j].APIAddress {
						statusIstrue = true
					}
				}
				if statusIstrue == false {
					node[i].Display = false
					if err := GetDB(nil).Table("fullnode_info").Where("id = ?", node[i].ID).Update("display", node[i].Display).Error; err != nil {
						log.WithFields(log.Fields{"error": err}).Error("sync display status update1 err")
						continue
					}
					statusDiff = true
				}
			}
		}
	} else {
		for i := 0; i < len(node); i++ {
			if node[i].Display == true {
				node[i].Display = false
				if err := GetDB(nil).Table("fullnode_info").Where("id = ?", node[i].ID).Update("display", node[i].Display).Error; err != nil {
					log.WithFields(log.Fields{"error": err}).Error("sync display status update2 err")
					continue
				}
				statusDiff = true
			}
		}
	}
	return statusDiff
}
func syncFullNodeInfoToRedis() {
	var node []FullNodeInfo
	if err := GetDB(nil).Table("fullnode_info").Order("id desc").Find(&node).Error; err != nil {
		log.WithFields(log.Fields{"error": err}).Error("select fullnodeInfo failed")
		return
	}
	value, err := json.Marshal(node)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("sync fullnodeInfo json failed")
		return
	}
	rd := RedisParams{
		Key:   "fullNode",
		Value: string(value),
	}
	if err := rd.Set(); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("set FullNodeInfoToRedis failed")
		return
	}

}
func SyncFullNodeInfo() {
	FullnodesInfo = GetFullnodesInfo()
}
func FindNodeAddressInsave(list []storage.FullnodeModel) (err error) {
	err = nil
	var ipAddress string
	var fullNoddeInfoList []FullNodeInfo

	var fullNoddeInfo FullNodeInfo
	if err = fullNoddeInfo.createDataBass(); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("create table fullnode_info err")
		return err
	}
	fullNoddeInfo.getAddressList()
	for i := 0; i < len(list); i++ {
		ipAddress = getIPAddress(list[i].APIAddress)
		if isNotIp := net.ParseIP(list[i].APIAddress); isNotIp == nil {
			addr, err := net.ResolveIPAddr("ip", ipAddress)
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Error("err resolveip")
				continue
			} else {
				ipAddress = addr.String()
			}
		}
		if ipAddress == "" {
			continue
		}
		dBRord, _ := os.Getwd()
		dBRord = path.Join(dBRord, "geoip-dataBass", "GeoLite2-City.mmdb")
		info, result := findAddressFromIp(dBRord, ipAddress)
		if result == 3 {
			fullNoddeInfo.Address = strings.ToLower(info.CityName)
		} else if result == 2 {
			time.Sleep(time.Millisecond * 200)
			if address := queryAddressByLatitudeAndLongitude(info.Longitude, info.Latitude); address != "" {
				fullNoddeInfo.Address = strings.ToLower(address)
			} else {
				if ipAddress == "127.0.0.1" {
					fullNoddeInfo.Address = "china-beijing"
					fullNoddeInfo.Longitude = 116.3952912
					fullNoddeInfo.Latitude = 39.9087202
				} else {
					continue
				}
			}
		} else {
			continue
		}
		if ipAddress != "127.0.0.1" {
			fullNoddeInfo.Latitude = info.Latitude
			fullNoddeInfo.Longitude = info.Longitude
		}
		fullNoddeInfo.Number = fullNoddeInfo.getCityNum(fullNoddeInfoList)
		fullNoddeInfo.Display = true

		li, err := json.Marshal(list[i])
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("FindNodeAddress jsonerr")
			continue
		}

		fullNoddeInfo.Value = string(li)
		if !fullNoddeInfo.hasFullNode() {
			fullNoddeInfoList = append(fullNoddeInfoList, fullNoddeInfo)
		}

	}
	if len(fullNoddeInfoList) > 0 {
		for i := 0; i < len(fullNoddeInfoList); i++ {
			fullNoddeInfo = fullNoddeInfoList[i]
			//Before writing to the database, judge whether this value exists, if it exists, then judge whether the address of the ip is changed and update it, otherwise it will not be written
			if err := fullNoddeInfo.insertData(); err != nil {
				log.WithFields(log.Fields{"warn": err}).Warn("fullnodeinfo insertData err")
			}
		}
	}
	addrList = nil
	return nil
}

func GetNodeMapInfo() (nodeMap []NodeMapInfo, err error) {
	var info []FullNodeInfo

	rd := RedisParams{
		Key:   "fullNode",
		Value: "",
	}
	if err = rd.Get(); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("GetNodeMapInfo getdb err")
		return nil, err
	}

	if err = json.Unmarshal([]byte(rd.Value), &info); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("GetNodeMapInfo json err")
		return nil, err
	}

	nodeMap = make([]NodeMapInfo, 0)

	var node NodeMapInfo
	for i := 0; i < len(info); i++ {
		if info[i].Display == false {
			continue
		}
		node.Name = info[i].Address
		node.Latitude = info[i].Latitude
		node.Longitude = info[i].Longitude
		node.Number = info[i].Number
		nodeMap = append(nodeMap, node)
	}
	return nodeMap, nil
}
func (p *FullNodeInfo) GetNodeList() (node []storage.FullnodeModel, err error) {
	var info []FullNodeInfo
	rd := RedisParams{
		Key:   "fullNode",
		Value: "",
	}
	if err = rd.Get(); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("GetNodeList getdb err")
		return nil, err
	}
	if err = json.Unmarshal([]byte(rd.Value), &info); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("GetNodeList json err")
		return nil, err
	}
	info = redisOrderNode(info, "id desc")
	node = make([]storage.FullnodeModel, len(info))
	for i := 0; i < len(info); i++ {
		nodeValue := NodeInfo{}
		err := json.Unmarshal([]byte(info[i].Value), &nodeValue)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("GetNodeList jsonerr")
			continue
		}
		node[i].Enable = true
		node[i].Nodename = "node-" + getCity(info[i].Address)
		node[i].TCPAddress = nodeValue.TcpAddress
		node[i].APIAddress = nodeValue.ApiAddress
		node[i].City = info[i].Address
		node[i].Icon = getIcon(info[i].Address)
		node[i].NodePosition = int64(i + 1)
		node[i].KeyID = strconv.FormatInt(crypto2.Address([]byte(nodeValue.PublicKey)), 10)
		node[i].PublicKey = nodeValue.PublicKey
		node[i].UnbanTime = time.Now()
		node[i].Latitude = strconv.FormatFloat(info[i].Latitude, 'f', 5, 64)
		node[i].Longitude = strconv.FormatFloat(info[i].Longitude, 'f', 5, 64)
		node[i].Name = "ibax"
		//node[i].Engine = "postgres"
		//node[i].Version = "10.3"
		node[i].Display = info[i].Display
		node[i].IconUrl = getIconNationalFlag(node[i].Icon)
	}
	return node, nil
}
func getIconNationalFlag(icon string) string {
	road, _ := os.Getwd()
	road = path.Join(road, "logodir")
	var pictureName string
	picturefiles, _ := ioutil.ReadDir(road)
	for _, f := range picturefiles {
		if strings.Contains(f.Name(), ".png") {
			fn := strings.Replace(f.Name(), ".png", "", -1)
			if strings.Contains(icon, fn) {
				pictureName = conf.GetEnvConf().Url.URL + f.Name()
			} else if strings.EqualFold(strings.Replace(icon, " ", "", -1), "unitedstates") {
				if fn == "usa" {
					pictureName = conf.GetEnvConf().Url.URL + f.Name()
				}
			}
		}
	}
	return pictureName
}
func getCity(city string) string {
	if strings.Contains(city, "-") {
		if index := strings.Index(city, "-"); index != -1 {
			return city[index+1:]
		}
	}
	return city
}
func getIcon(city string) string {
	if strings.Contains(city, "-") {
		if index := strings.Index(city, "-"); index != -1 {
			return city[:index]
		}
	}
	return city
}
func (p *FullNodeInfo) getAddressList() {
	if err := GetDB(nil).Table(p.TableName()).Select("address").Order("id desc").Find(&addrList).Error; err != nil {
		log.WithFields(log.Fields{"error": err}).Error("getAddressList Find err")
		return
	}
}

func getIPAddress(addressName string) (ip string) {
	ip = addressName
	if strings.Contains(addressName, "http") {
		total := `://`
		if index1 := strings.Index(addressName, total); index1 != -1 {
			if strings.Contains(addressName[index1+len(total):], ":") {
				if index2 := strings.Index(addressName[index1+len(total):], ":"); index2 != -1 {
					ip = addressName[index1+len(total) : index1+len(total)+index2]
				}
			} else {
				ip = addressName[index1+len(total):]
			}
		}
	}
	return ip
}
func findAddressFromIp(dataBassName string, ipname string) (info infoIpInform, findResult int) {
	findResult = 1
	//findResult 1:unVaild 2:notFind 3:findout   2 or 3:vaild
	info = infoIpInform{}
	defer func() {
		if e := recover(); e != nil {
			panic(e)
		}
	}()
	db, err := geoip2.Open(dataBassName)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("findAddressFromIp open err")
		return info, findResult
	}
	defer db.Close()
	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP(ipname)
	record, err1 := db.City(ip)
	if err1 != nil {
		log.WithFields(log.Fields{"error": err1}).Error("findAddressFromIp dbcity err")
		return info, findResult
	}
	if len(record.Subdivisions) > 0 {
		cityName := record.Subdivisions[0].Names["en"]
		contry := record.Country.Names["en"]
		if cityName == "" {
			cityName = record.Subdivisions[0].Names["zh-CN"]
			if cityName == "" {
				info.Latitude = record.Location.Latitude
				info.Longitude = record.Location.Longitude
				info.CityName = strings.ToLower(strings.Replace(contry, " ", "", -1))
				if info.CityName == "" {
					findResult = 2
					return info, findResult
				}
				findResult = 3
				return info, findResult
			}
			cityName = Translate(cityName)
		}
		info.CityName = strings.ToLower(strings.Replace(contry, " ", "", -1) + "-" + cityName)
	}
	// Output:
	//Portuguese (BR) city name:
	//English subdivision name: Shanghai
	//Russian country name: China
	//Russian Continent name: 亚洲
	//ISO country code: CN
	//Time zone: Asia/Shanghai
	//Coordinates: 31.1394, 121.1001

	info.Latitude = record.Location.Latitude
	info.Longitude = record.Location.Longitude

	if info.CityName == "" {
		findResult = 2
		return info, findResult
	}
	findResult = 3
	return info, findResult
}
func queryAddressByLatitudeAndLongitude(longitude, latitude float64) string {
	lo := strconv.FormatFloat(longitude, 'f', 5, 64)
	la := strconv.FormatFloat(latitude, 'f', 5, 64)
	//test1,_:=strconv.ParseFloat(strconv.FormatFloat(record.Location.Longitude,'f',10,64),64)
	if lo == "" || la == "" {
		log.WithFields(log.Fields{"warn": "longitude or latitude"}).Warn("is null")
		return ""
	}
	addrInfo := requestAddressFromTencent(lo, la)
	if addrInfo != nil {
		nation := addrInfo.Result.AdInfo.Nation
		province := addrInfo.Result.AdInfo.Province
		nation1 := addrInfo.Result.AddressComponent.Nation
		province1 := addrInfo.Result.AddressComponent.Province
		locality := addrInfo.Result.AddressComponent.Locality
		addr := addrInfo.Result.Address
		if nation != "" && province != "" {
			return Translate(nation) + "-" + Translate(province)
		} else if nation1 != "" && province1 != "" {
			return Translate(nation1) + "-" + Translate(province1)
		} else if nation1 != "" && locality != "" {
			return Translate(nation1 + "-" + locality)
		} else if addr != "" {
			return Translate(addr)
		}
	}
	return ""
}

//TODO: upper limit: 5 times/s, 1000 times per day
func requestAddressFromTencent(long, lat string) *addressInfo {
	urls := "https://apis.map.qq.com/ws/geocoder/v1/?location=" + lat + "," + long + "&key=" + getDefaultAK()

	client := http.Client{}
	req, err := http.NewRequest("GET", urls, nil)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("newrequest err")
		return nil
	}
	response, err1 := client.Do(req)
	if err1 != nil {
		log.WithFields(log.Fields{"error": err1}).Error("http newrequest err")
		return nil
	}
	defer response.Body.Close()
	addrInfo := &addressInfo{}
	if response.StatusCode == 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("http readall err")
			return nil
		}
		if err := json.Unmarshal(body, addrInfo); err != nil {
			log.WithFields(log.Fields{"error": err}).Error("http json err")
			return nil
		} else {
			if addrInfo.Status == 0 {
				return addrInfo
			} else {
				return nil
			}
		}
	}
	return nil
}

func (p *FullNodeInfo) getCityNum(list []FullNodeInfo) (num int64) {
	num = 1
	for i := 0; i < len(list); i++ {
		if p.Address == list[i].Address {
			num = num + 1
		}
	}
	for i := 0; i < len(addrList); i++ {
		if p.Address == addrList[i] {
			num = num + 1
		}
	}

	return
}

func (p *FullNodeInfo) hasFullNode() bool {
	var valueList []FullNodeInfo
	if err := GetDB(nil).Table(p.TableName()).Select("value,address,id").Find(&valueList).Error; err != nil {
		log.WithFields(log.Fields{"error": err}).Error("hasFullNode dbfind")
		return false
	}
	for i := 0; i < len(valueList); i++ {
		if p.Value == valueList[i].Value {
			if p.Address != valueList[i].Address {
				p.ID = valueList[i].ID
				if err := GetDB(nil).Model(FullNodeInfo{}).Updates(&p).Error; err != nil {
					log.WithFields(log.Fields{"error": err}).Error("hasFullNode update err")
					return true
				}
				p.getAddressList()
				return true
			}
			return true
		}
	}

	return false
}
func (p *FullNodeInfo) TableName() string {
	return "fullnode_info"
}
func (p *FullNodeInfo) createDataBass() (err error) {
	err = nil
	if !GetDB(nil).Migrator().HasTable(p) {
		if err = GetDB(nil).Migrator().CreateTable(p); err != nil {
			return err
		}
	}
	return err
}

func (p *FullNodeInfo) insertData() (err error) {
	err = nil
	var id int64
	if err = GetDB(nil).Table(p.TableName()).Count(&id).Error; err != nil {
		return err
	}
	p.ID = id + 1
	if err = GetDB(nil).Create(&p).Error; err != nil {
		return err
	}
	return err
}

func getDefaultAK() string {
	ak := defaultAppKey // tenxun
	return ak
}
func GetFullnodesInfo() []*storage.FullnodeModel {
	var fullNodeInfo FullNodeInfo
	var nodeInfo []*storage.FullnodeModel

	node, err := fullNodeInfo.GetNodeList()
	if err == nil {
		nodeInfo = make([]*storage.FullnodeModel, len(node))
		for i := 0; i < len(node); i++ {
			nodeInfo[i] = &node[i]
		}
	} else {
		log.WithFields(log.Fields{"error": err}).Error("GetFullnodesInfo err")
		return nil
	}

	return nodeInfo
}
func IsDisplay(fullnode []*storage.FullnodeModel) []*storage.FullnodeModel {
	var red []*storage.FullnodeModel
	for i := 0; i < len(fullnode); i++ {
		if fullnode[i].Display == true {
			red = append(red, fullnode[i])
		}
	}
	return red
}

func redisOrderNode(cd []FullNodeInfo, order string) (rd []FullNodeInfo) {
	if strings.Contains(order, "id desc") {
		sort.SliceStable(cd, func(i, j int) bool {
			return cd[i].ID > cd[j].ID
		})
	} else if strings.Contains(order, "id asc") {
		sort.SliceStable(cd, func(i, j int) bool {
			return cd[i].ID < cd[j].ID
		})
	} else {
		log.WithFields(log.Fields{"warn": order}).Warn("redisOrderNode not find warn")
	}
	rd = cd
	return
}
func GetNodeListInfo(nodeInfo string) (err error) {
	err = nil
	var fullNode []storage.FullnodeModel
	err = json.Unmarshal([]byte(nodeInfo), &fullNode)
	if err != nil {
		return err
	}

	var nodeValue []string
	if err := GetDB(nil).Table("fullnode_info").Select("value").Order("id desc").Find(&nodeValue).Error; err != nil {
		log.WithFields(log.Fields{"error": err}).Error("fullNodeDbIsExist json failed")
	}
	var node []storage.FullnodeModel
	for i := 0; i < len(fullNode); i++ {
		if fullNodeDbIsExist(fullNode[i].APIAddress, nodeValue) {
			continue
		} else {
			node = append(node, fullNode[i])
		}
	}
	if len(node) > 0 {
		err = FindNodeAddressInsave(node)
	}
	return err
}
func Translate(text string) string {
	urls := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=zh-cn&tl=en&dt=t&q=%s", url.QueryEscape(text))
	resp, err := http.Get(urls)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Translate get failed")
		return ""
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Translate ioutil read failed")
		return ""
	}
	//The returned json deserialization is more troublesome, direct string disassembly
	ss := string(bs)
	ss = strings.ReplaceAll(ss, "[", "")
	ss = strings.ReplaceAll(ss, "]", "")
	ss = strings.ReplaceAll(ss, "null,", "")
	ss = strings.Trim(ss, `"`)
	ps := strings.Split(ss, `","`)
	return ps[0]
}
