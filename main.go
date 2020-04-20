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

	http.HandleFunc("/user/signup",handle.SignupHandler)
	http.HandleFunc("/user/signin",handle.SignInHandler)
	http.HandleFunc("/user/info",handle.HTTPInterceptor(handle.UserInfoHandler))

	http.Handle("/static/",
		http.StripPrefix("/static/",http.FileServer(http.Dir("./static"))))
	err := http.ListenAndServe(":8000", nil)
	if err != nil{
		fmt.Printf("Failed to start , err:%s",err.Error())
	}
}
