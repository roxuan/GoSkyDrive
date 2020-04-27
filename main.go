package main

import (
	"filestore-server/handle"
	"fmt"
	"net/http"
)

func main(){
	http.HandleFunc("/file/upload",handle.UploadHandler)
	http.HandleFunc("/file/upload/suc",handle.UploadSucHandler)
	http.HandleFunc("/file/meta",handle.GetFileMetaHandler)
	http.HandleFunc("/file/download",handle.DownloadHandler)
	http.HandleFunc("/file/update",handle.FileMetaUpdateHandle)
	http.HandleFunc("/file/delete",handle.FileDeleteHandle)
	http.HandleFunc("/file/query",handle.FileQueryHandler)
	http.HandleFunc("/file/fastupload",handle.TryFastUploadHandle)

	// 分块上传接口
	http.HandleFunc("/file/mpupload/init",handle.HTTPInterceptor(handle.InitialMultipartUploadHandler))
	http.HandleFunc("/file/mpupload/uppart",handle.HTTPInterceptor(handle.UploadPartHandle))
	http.HandleFunc("/file/mpupload/complete",handle.HTTPInterceptor(handle.CompleteUploadHandler))

	http.HandleFunc("/user/signup",handle.SignupHandler)
	http.HandleFunc("/user/signin",handle.SignInHandler)
	http.HandleFunc("/user/info",handle.HTTPInterceptor(handle.UserInfoHandler))

	http.Handle("/static/",
		http.StripPrefix("/static/",http.FileServer(http.Dir("./static"))))
	err := http.ListenAndServe(":8080", nil)
	if err != nil{
		fmt.Printf("Failed to start , err:%s",err.Error())
	}
}
