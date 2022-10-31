package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Token   string `json:"token"`
	Port    int    `json:"port"`
	Address string `json:"address"`
}

var config Config

func init() {
	if exists("config.json") {
		configfile, _ := os.ReadFile("config.json")
		err := json.Unmarshal(configfile, &config)
		if err != nil {
			fmt.Println("Error: config.json is not a valid json file.")
			panic(err)
		}
	} else {
		fmt.Println("config.json not found")
		os.Exit(0)
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Any("/*a", func(c *gin.Context) {
		method := c.Request.Method
		hh := c.Request.Header.Get("thost")
		url := c.Request.URL.String()
		client := http.Client{}
		turl := "https://" + hh + url
		req, err := http.NewRequest(method, turl, c.Request.Body)
		if err != nil {
			panic(err)
		}
		uheader := c.Request.Header.Clone()
		tkey := uheader.Get("tkey")
		if tkey != config.Token {
			c.String(403, "502 Bad Gateway")
			return
		}
		ua := c.Request.Header.Get("User-Agent")
		log.Printf("%s | %s", hh, ua)
		uheader.Set("User-Agent", ua+" (Forwarded by APFW 0.1)")
		uheader.Del("thost")
		uheader.Del("tkey")
		req.Header = uheader
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		a, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		Hlist := []string{"Strict-Transport-Security", "Content-Type"}
		for _, v := range Hlist {
			c.Header(v, res.Header.Get(v))
		}
		c.Status(res.StatusCode)
		c.Writer.Write(a)
	})
	fmt.Println("APFW is running on port", config.Port)
	r.Run(fmt.Sprint(config.Address) + ":" + fmt.Sprint(config.Port))
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
