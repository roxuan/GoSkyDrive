package db

import (
	mydb "filestore-server/db/mysql"
	"fmt"
)

// 用户文件表结构体
type UserFile struct{
	UserName string
	FileHash string
	FileName string
	FileSize int64
	UploadAt string
	LastUpdated string
}

// 更新用户文件表
func OnUserFileUploadFinished(username string,filehash string,filename string,filesize int64)bool{
	stmt,err:=mydb.DBConn().Prepare(
		"insert ignore into tbl_user_file(`user_name`,`file_sha1`,`file_name`,`file_size`) values(?,?,?,?)")
	if err != nil{
		fmt.Println("failed to insert, err:"+err.Error())
		return false
	}
	defer stmt.Close()

	_,err1:=stmt.Exec(username,filehash,filename,filesize)
	if err1!=nil{
		fmt.Println("failed to insert11111, err:"+err.Error())
		return false
	}
	return true
}

// 批量获取用户文件信息
func QueryUserFileMetas(username string,limit int)([]UserFile,error){
	stmt,err:=mydb.DBConn().Prepare(
		"select file_sha1,file_name,file_size,upload_at,last_update from tbl_user_file where user_name = ? limit ?")
	if err != nil{
		fmt.Println("failed to insert1, err:"+err.Error())
		return nil,err
	}
	defer stmt.Close()

	rows,err:=stmt.Query(username,limit)
	if err!=nil{
		fmt.Println("failed to insert2, err:"+err.Error())
		return nil,err
	}
	var userFiles []UserFile
	for rows.Next(){
		ufile :=UserFile{}
		err = rows.Scan(&ufile.FileHash,&ufile.FileName,&ufile.FileSize,
			&ufile.UploadAt,&ufile.LastUpdated)
		if err != nil{
			fmt.Println(err.Error())
			break
		}
		userFiles = append(userFiles,ufile)
	}
	return userFiles,err
}
