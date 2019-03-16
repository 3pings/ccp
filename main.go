package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type cluster struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	Workers int    `json:"workers"`
}

// Variable Declarations

var clusterID string
var workerCount = "3"
var clusterName = "mordor"

func main() {

	var clusters []cluster
	bs := getClusters("10.139.11.50/2/", "admin", "admin")
	err := json.Unmarshal(bs, &clusters)

	if err != nil {
		log.Fatalln("error unmarshaling", err)

	}

	for _, cluster := range clusters {
		if cluster.Name == clusterName {
			clusterID = cluster.UUID
		}
	}

	patchWorkers("10.139.11.50/2/", "admin", "admin", clusterID)

}

func getClusters(bURL, uName, pWord string) (clusterInfo []byte) {

	options := cookiejar.Options{}
	//Set TLS Parameters
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{Jar: jar, Transport: tr}

	//Postform for Login
	resp, err := client.PostForm("https://"+bURL+"system/login/", url.Values{
		"username": {uName},
		"password": {pWord},
	})
	if err != nil {
		log.Fatal(err)
	}

	//Response from GET clusters
	resp, err = client.Get("https://" + bURL + "clusters")

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func patchWorkers(bURL, uName, pWord, cID string) {

	options := cookiejar.Options{}
	//Set TLS Parameters
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{Jar: jar, Transport: tr}

	// Postform for login
	resp, err := client.PostForm("https://"+bURL+"system/login/", url.Values{
		"username": {uName},
		"password": {pWord},
	})
	if err != nil {
		log.Fatal(err)
	}

	//PUT Request to update Number of Workers
	payload := strings.NewReader("{\"workers\":" + workerCount + "}")

	req, _ := http.NewRequest("PATCH", "https://"+bURL+"clusters/"+cID, payload)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")

	resp, err = client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
}
