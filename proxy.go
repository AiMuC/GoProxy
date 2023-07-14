package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"log"
	"net/url"
	"steam/support"
)

func main() {
	Addr := support.Env("ADDR", "localhost:9797").(string)
	log.Println("代理服务器已启动，监听在 " + Addr)
	ln, err := reuseport.Listen("tcp4", Addr)
	if err != nil {
		log.Fatal("无法监听端口:", err)
	}
	server := &fasthttp.Server{
		Handler: proxyHandler,
	}
	if err := server.Serve(ln); err != nil {
		log.Fatal("无法启动服务器:", err)
	}
}

func proxyHandler(ctx *fasthttp.RequestCtx) {
	BaseUrl := support.Env("BASE_URL", "https://steamcommunity.com").(string)
	targetURL, err := url.Parse(BaseUrl)
	if err != nil {
		log.Fatal("无法解析目标URL:", err)
	}
	client := &fasthttp.Client{}
	targetURL.Path = "/" + string(ctx.Path())
	targetURL.RawQuery = string(ctx.QueryArgs().QueryString())
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(targetURL.String())
	req.Header.SetHost(targetURL.Host)
	req.Header.SetMethodBytes(ctx.Method())
	req.Header.SetReferer(BaseUrl)
	req.Header.Set("Access-Control-Allow-Origin", "*")
	req.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	req.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	req.SetBody(ctx.Request.Body())
	resp := fasthttp.AcquireResponse()
	if err := client.Do(req, resp); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	resp.Header.VisitAll(func(key, value []byte) {
		ctx.Response.Header.SetBytesKV(key, value)
	})
	ctx.Response.SetStatusCode(resp.StatusCode())
	ctx.Response.SetBody(resp.Body())
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
}
