package controllers

import (
	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-ibax/packages/consts"
	"github.com/IBAX-io/go-ibax/packages/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var LogoDir = "./logodir/"

func init() {
	if err := utils.MakeDirectory(LogoDir); err != nil {
		log.WithFields(log.Fields{"error": err, "type": consts.IOError, "dir": LogoDir}).Error("can't create temporary directory")
	}
}

// @Router /getmanagefile/{file} [get]
func logoHandler(c *gin.Context) {
	ret := &Response{}
	var scanout models.ScanOut
	rets, err := scanout.GetRedisdashboard()
	if err != nil {
		ret.ReturnFailureString(err.Error())
	//if !IsExist(LogoDir + fileName) {
	//	//errorResponse(w, errFileNotExists.Errorf(fileName))
	//	ret.Return(nil, model.CodeFileNotExists.String(fileName))
	//	JsonCodeResponse(w, &ret)
	//	return
	//}
	//
	//file, err := os.Open(LogoDir + fileName)
	//if err != nil {
	//	logger.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("getting file err")
	//	//errorResponse(w, err)
	//	ret.Return(nil, model.CodeFileOpen.Errorf(err))
	//	JsonCodeResponse(w, &ret)
	//	return
	//}
	//defer file.Close()
	//
	//stat, err := file.Stat() //
	//if err != nil {
	//	logger.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("getting file err")
	//	//errorResponse(w, err)
	//	ret.Return(nil, model.CodeFileOpen.Errorf(err))
	//	JsonCodeResponse(w, &ret)
	//	return
	//}
	//sz := stat.Size()
	////str := strconv.FormatInt(100,10)
	//str := strconv.FormatInt(sz, 10)
	////str :=strconv.Itoa(sz)
	//w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	//w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	//w.Header().Set("Content-length", str)
	//io.Copy(w, file)
	//
	//dat := model.MineActionResult{
	//	Action: "download",
	//	Result: fileName,
	//}
	//ret.Return(dat, model.CodeSuccess)
	//jsonResponse(w, ret)
}

////
//// logoHandler godoc
//// @Summary
//// @Description
//// @Accept mpfd
//// @Produce json
//// @Param file  path  string  true "filename"
//// @Success 200  {object} model.MineActionResponseResult
//// @Failure 400 {string} string ""
//// @Failure 404 {string} string ""
//// @Failure 500 {string} string ""
//// @host localhost：7079
//// @BasePath /api/v2
//// @Router /getmanagefile/{file} [get]
//func LoadlogoHandler(w http.ResponseWriter, r *http.Request) {
//	ret := model.Response{}
//	params := mux.Vars(r)
//
//	id, err := strconv.ParseInt(params["id"], 10, 64)
//	if err != nil {
//		ret.Return(nil, model.CodeFileOpen.Errorf(err))
//		JsonCodeResponse(w, &ret)
//		return
//	}
//	file, err := model.Loadlogo(id)
//	if err != nil {
//		ret.Return(nil, model.CodeFileOpen.Errorf(err))
//		JsonCodeResponse(w, &ret)
//		return
//	}
//	rs := conf.Config.Manage.URL + "/api/v2/logo/" + file
//	dat := model.MineActionResult{
//		Action: "logo download",
//		Result: rs,
//	}
//	ret.Return(dat, model.CodeSuccess)
//	jsonResponse(w, ret)
//}
