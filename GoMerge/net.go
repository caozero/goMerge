package GoMerge

import (
	"fmt"
	"code.google.com/p/go.net/websocket"
	"net/http"
	"log"
	"net"
	"strings"
	"encoding/json"
)

type Msg struct {
	Cmd string `json:"cmd"`		//回调函数
	Status int32 `json:"status"`	//状态,默认0成功,大于0为失败
	Msg string `json:"msg"`		//返回信息
	Data interface{} `json:"data"`	//返回数据
}

type Net struct {
	manager *Manager
	Socket string
	clientList []*websocket.Conn
}

func (p *Net)Start(){
	http.Handle("/", http.FileServer(http.Dir("www"))) // <-- note this line
	http.Handle("/dataSocket", websocket.Handler(p.dataSocket))
	addrs:=p.getIpAddr();
	fmt.Printf("监听网址:%+v\n",addrs)
	fmt.Printf("端口号:%v\n",p.Socket)
	go func(){
		if err := http.ListenAndServe(":"+p.Socket, nil); err != nil {
			log.Fatal("开启遥网络端口失败:", err)
		}
	}()
}
func (p *Net)Send(msg interface{}){
	m,e:=json.Marshal(msg)
	if e!=nil{
		return
	}
	s:=string(m)
	for _,v:=range p.clientList{
		p.sendToClient(v,s)
	}
}

func (p *Net)sendToClient(ws *websocket.Conn,msg string){
	websocket.Message.Send(ws,msg)
}


func (p *Net)dataSocket(ws *websocket.Conn) {
	p.clientList=append(p.clientList,ws)
	fmt.Printf("客户端连接.\n")
	var err error
	for {
		var v string
		if err = websocket.Message.Receive(ws, &v); err != nil {
			fmt.Println("遥控器端断开连接.")
			break
		}
		fmt.Printf("接受指令:%s\n", v)
		if p.manager.rount(v)==true{

		}

	}
	p.delWs(ws)
	fmt.Printf("客户端退出.\n")
}

func (p *Net)delWs(ws *websocket.Conn){
	for k,v:=range p.clientList{
		if v==ws{
			p.clientList=append(p.clientList[:k],p.clientList[k+1:]...)
			return
		}
	}
}

func (p *Net)getIpAddr()[]string{
	addrs:=[]string{}
	info, _ := net.InterfaceAddrs()
	for _, addr := range info {
		a:=strings.Split(addr.String(), "/")[0]
		if strings.Index(a,":")==-1 && a!="127.0.0.1"{
			addrs=append(addrs,a)
		}
	}

	return addrs
}