package meta

import(
	mydb "filestore-server/db"
)
// FileMeta:文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init(){
	fileMetas = make(map[string]FileMeta)
}

// 新增/更新文件元信息
func UpdateFileMeta(fmeta FileMeta){
	fileMetas[fmeta.FileSha1] = fmeta
}

// 新增/更新文件元信息到mysql
func UpdateFileMetaDb(fmeta FileMeta)bool{
	return mydb.OnFileUploadFinished(
		fmeta.FileSha1,fmeta.FileName,fmeta.FileSize,fmeta.Location)
}

// 通过sha1值获取文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta{
	return fileMetas[fileSha1]
}

// 从mysql获取文件元信息
func GetFileMetaDb(fileSha1 string)(FileMeta,error){
	tfile,err:=mydb.GetFileMeta(fileSha1)
	if err!=nil{
		return FileMeta{},err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return fmeta,nil
}

// 删除文件元信息
func RemoveFileMeta(fileSha1 string){
	delete(fileMetas, fileSha1)
}