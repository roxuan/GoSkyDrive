package main

import (
	"filestore-serve/handle"
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
	err := http.ListenAndServe(":8000", nil)
	if err != nil{
		fmt.Printf("Failed to start , err:%s",err.Error())
	}
}
