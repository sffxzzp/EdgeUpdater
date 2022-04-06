package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func httpPostJson(reqURL string, postData string) []byte {
	res, _ := http.Post(reqURL, "application/json", bytes.NewBuffer([]byte(postData)))
	data, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return data
}

type (
	EdgeUpdate struct {
		// https://docs.microsoft.com/zh-cn/deployedge/microsoft-edge-update-policies
		url       string
		branch    []string
		structure []string
		// filesData []fileData
	}
	verRet struct {
		ContentId map[string]string `json:"ContentId"`
		Files     []string          `json:"Files"`
	}
	fileRetF struct {
		FileId   string            `json:"FileId"`
		Url      string            `json:"Url"`
		Size     int               `json:"SizeInBytes"`
		Hashes   map[string]string `json:"Hashes"`
		Delivery interface{}       `json:"DeliveryOptimization"`
		Version  string
	}
	fileData struct {
		Branch    string
		Structure string
		FileData  fileRetF
	}
)

func newEdgeUpdate() *EdgeUpdate {
	return &EdgeUpdate{
		url:       "https://msedge.api.cdp.microsoft.com",
		branch:    []string{"stable", "beta", "dev", "canary"},
		structure: []string{"x86", "x64", "arm64"},
	}
}

// func (u *EdgeUpdate) getInfoAll() {
// 	var filesData []fileData
// 	for _, branch := range u.branch {
// 		for _, struc := range u.structure {
// 			filesData = append(filesData, fileData{
// 				Branch:    branch,
// 				Structure: struc,
// 				FileData:  u.getLatestFile(fmt.Sprintf("msedge-%s-win-%s", branch, struc)),
// 			})
// 		}

// 	}
// 	u.filesData = filesData
// }

func (u *EdgeUpdate) getLatestFile(name string) fileRetF {
	res := u.getLatestVersion(name)
	return u.getFile(name, res.ContentId["Version"])
}

func (u *EdgeUpdate) getFile(name string, version string) fileRetF {
	res := httpPostJson(fmt.Sprintf("%s/api/v1/contents/Browser/namespaces/Default/names/%s/versions/%s/files?action=GenerateDownloadInfo", u.url, name, version), `{"targetingAttributes":{}}`)
	var fileRet []fileRetF
	json.Unmarshal(res, &fileRet)
	var fRet fileRetF
	for _, file := range fileRet {
		if len(strings.Split(file.FileId, "_")) < 4 {
			fRet = file
			fRet.Version = version
		}
	}
	if fRet.FileId == "" {
		log.Println("error")
	}
	return fRet
}

func (u *EdgeUpdate) getLatestVersion(name string) verRet {
	var verRet verRet
	res := httpPostJson(fmt.Sprintf("%s/api/v1/contents/Browser/namespaces/Default/names/%s/versions/latest?action=select", u.url, name), `{"targetingAttributes":{"IsInternalUser":true}}`)
	json.Unmarshal(res, &verRet)
	return verRet
}

func (u *EdgeUpdate) getInfo(branch string, structure string) *fileData {
	return &fileData{
		Branch:    branch,
		Structure: structure,
		FileData:  u.getLatestFile(fmt.Sprintf("msedge-%s-win-%s", strings.ToLower(branch), strings.ToLower(structure))),
	}
}
