package handle

import(
	"encoding/json"
	dblayer "filestore-server/db"
	"filestore-server/meta"
	"filestore-server/util"
	"strconv"

	//dblayer "filestore-server/db"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func UploadHandler(w http.ResponseWriter, r * http.Request){
	if r.Method == "GET" {
		// 返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil{
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	}else if r.Method == "POST"{
		// 接收文件流及存储到本地目录
		file, head, err := r.FormFile("file")
		if err != nil{
			fmt.Printf("failed to get data, err:%s\n", err.Error())
			return
		}
		defer file.Close()

		fileMeta:= meta.FileMeta{
			FileName: head.Filename,
			Location: "D:/"+ head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		newFile, err := os.Create(fileMeta.Location)
		if err!= nil {
			fmt.Printf("failed to create file, err:%s\n", err.Error())
			return
		}
		defer newFile.Close()
		
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err!=nil {
			fmt.Printf("failed to save data into file, err:%s\n", err.Error())
			return
		}

		newFile.Seek(0,0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		//meta.UpdateFileMeta(fileMeta)
		_=meta.UpdateFileMetaDb(fileMeta)

		// TODO:更新用户文件表记录
		r.ParseForm()
		username := r.Form.Get("username")
		suc := dblayer.OnUserFileUploadFinished(username,fileMeta.FileSha1,fileMeta.FileName,fileMeta.FileSize)
		if suc{
			http.Redirect(w,r,"/static/view/home.html",http.StatusFound)
		}else{
			w.Write([]byte("Upload Failded."))
		}

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}

}

// 上传文件已完成
func UploadSucHandler(w http.ResponseWriter, r * http.Request){
	io.WriteString(w, "upload finished!")
}

// 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	filehash := r.Form["filehash"][0]
	//fMeta := meta.GetFileMeta(filehash)
	fMeta,err := meta.GetFileMetaDb(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 用户查询文件Hash信息
func FileQueryHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()

	limitCnt,_ :=strconv.Atoi(r.Form.Get("limit"))
	username := r.Form.Get("username")

	userFiles,err:=dblayer.QueryUserFileMetas(username,limitCnt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data,err:=json.Marshal(userFiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// DownloadHandler : 文件下载接口
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)

	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	w.Header().Set("content-disposition", "attachment; filename=\""+fm.FileName+"\"")
	w.Write(data)
}

// 更新文件元信息
func FileMetaUpdateHandle(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0"{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)

	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// 删除文件元信息
func FileDeleteHandle(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	fileSha1 := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(fileSha1)
	os.Remove(fMeta.Location)

	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK)
}

// 尝试秒传接口
func TryFastUploadHandle(w http.ResponseWriter,r *http.Request){
	r.ParseForm()

	// 1.解析请求参数
	username:=r.Form.Get("username")
	filehash:=r.Form.Get("filehash")
	filename:=r.Form.Get("filename")
	filesize, _ :=strconv.Atoi(r.Form.Get("filesize"))

	// 2.从文件表中查询相同hash的文件记录
	fileMeta,err:=meta.GetFileMetaDb(filehash)

	// 3.查不到记录则返回秒传失败
	if fileMeta == nil{
		resp:=util.RespMsg{
			Code: -1,
			Msg:  "秒传失败,请访问普通上传接口",
			Data: nil,
		}
		w.Write(resp.JSONBytes())
		return
	}
	if err!=nil{
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 4.上传过则将文件信息写入用户文件表,返回成功
	suc := dblayer.OnUserFileUploadFinished(username,filehash,filename,int64(filesize))
	if suc{
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
			Data: nil,
		}
		w.Write(resp.JSONBytes())
		return
	}else{
		resp:=util.RespMsg{
			Code: -2,
			Msg:  "秒传失败,请稍后重试",
			Data: nil,
		}
		w.Write(resp.JSONBytes())
		return
	}
}