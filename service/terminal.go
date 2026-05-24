package service

import (
	"context"
	"dashboard-api/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

var Terminal terminal

type terminal struct{}

const END_OF_TRANSMISSON = "\u0004"

// TerminalMessage 定义了终端和容器shell交互内容的格式
// Operation 是操作类型
// Data	是具体数据内容
// Rows 和 Cols 可以理解为终端的行数和列数, 也就是宽和高
type TerminalMessage struct {
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Rows      uint16 `json:"rows"`
	Cols      uint16 `json:"cols"`
}

//定义TerminalSession结构体, 实现ptyHandler接口
//wsconn 是websocket 连接
//sizeChan 是用来定义终端输入和输出的宽和高
//doneChan 用于标记退出终端
type TerminalSession struct {
	wsConn *websocket.Conn
	sizeChan	chan	remotecommand.TerminalSize
	doneChan	chan	struct{}
}








//定义websocket的handler方法
func (t *terminal) WsHandler(w http.ResponseWriter, r *http.Request)  {
	//加载k8s 配置
	config, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	if err != nil {
		logger.Error("创建k8s配置失败. " + err.Error())
	}
	//解析form入参, 获取namespace, podName, containerName参数
	if err := r.ParseForm(); err != nil {
		return 
	}
	namespace := r.Form.Get("namespace")
	podName	  := r.Form.Get("pod_name")
	containerName := r.Form.Get("container_name")
	logger.Info("exec pod: %s, container: %s, namespace: %s\n", podName, containerName, namespace)

	//new一个TerminalSession类型的pty实例, 引用另外一个函数
	pty, err := NewTerminalSession(w, r, nil)
	if err != nil {
		logger.Error("get pty failed: %v\n", err)
	}

	//处理关闭
	defer func ()  {
		logger.Info("close session.")
		pty.Close()
	}()

	//创建一个带取消信号的上下文对象
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	//初始化pod所在的corev1资源组
	//podExecOptions struct 包括Container stdout stdin Command 等结构
	//scheme.ParameterCodec 是pod的GVK (GroupVersion & Kind) 之类的
	// https://192.168.1.11:6443/api/v1/namespaces/default/pods/nginx-wf2-778d88d7c-7rmsk/
	// exec?command=%2Fbin%2Fbash&container=nginx-wf2&stderr=true&stdin=true&stdout=true&tty=true

	req := 	K8s.ClientSet.CoreV1().RESTClient().Post().Resource("pods").
			Name(podName).
			Namespace(namespace).
			SubResource("exec").
			VersionedParams(&corev1.PodExecOptions{
				Container: containerName,
				Command: []string{"/bin/bash"},
				Stdin: true,
				Stdout: true,
				Stderr: true,
				TTY: true,
			}, scheme.ParameterCodec)

	fmt.Println(req.URL())
	
	//remotecommand 主要实现了http转SPDY, 添加X-steream-Protocol-version相关的header, 并发送请求
	executor, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return
	}

	//建立连接之后从请求的stream中发送, 读取数据
	err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin: pty,
		Stdout: pty,
		Stderr: pty,
		TerminalSizeQueue: pty,
		Tty: true,
	})

	if err != nil {
		msg := fmt.Sprintf("Exec to pod error! err:%v", err)
		logger.Info(msg)

		//将报错反回去
		pty.Write([]byte(msg))
		//标记退出stream流
		pty.Done()
	}

}

//初始化一个websocket.Upgrader类型的对象, 用http协议升级websocket协议
var	upgrader = func ()  websocket.Upgrader {
	upgrader := websocket.Upgrader{}
	upgrader.HandshakeTimeout = time.Second * 2
	upgrader.CheckOrigin	= func(r *http.Request) bool {
		return true
	}
	return upgrader
}()


//该方法用于升级http协议到websocket, 并new一个terminalSession类型的对象返回
func NewTerminalSession(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (terminalSession *TerminalSession, err error)  {
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}

	session := &TerminalSession{
		wsConn: conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
	}

	return session, nil
}

//用于读取web端的输入, 接收web端输入的指令内容
func (t *TerminalSession) Read(p []byte) (int, error) {
	_, message, err := t.wsConn.ReadMessage()
	if err != nil {
		log.Printf("read message err: %v", err)
	}

	var msg TerminalMessage
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		log.Printf("read parse message err: %v", err)
		return copy(p, []byte(END_OF_TRANSMISSON)), err
	}

	switch msg.Operation {
	case "stdin":
		return copy(p, []byte(msg.Data)), nil
	case "resize":
		return 0, nil
	case "ping":
		return 0, nil
	default:
		log.Printf("unknow message type '%s'", msg.Operation)
		return copy(p, []byte(END_OF_TRANSMISSON)), fmt.Errorf("unknow message type '%s'", msg.Operation)
	}
}

//用于向web端输出, 接收web端输入的指令内容
func (t *TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(TerminalMessage{
		Operation: "stdout",
		Data: string(p),
	})

	if err != nil {
		log.Printf("write parse messge err: %v", err)
		return 0, err
	}

	if err := t.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Printf("write messge err: %v", err)
		return 0, err
	}
	return len(p), nil
}

//关闭donechan, 关闭后触发退出终端
func (t *TerminalSession) Done()  {
	close(t.doneChan)
}

//获取web终端是否resize, 以及是否退出终端
func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <- t.sizeChan:
		return &size
	case <- t.doneChan:
		return nil
	}
}

//用于关闭websocket连接
func (t *TerminalSession) Close() error  {
	return t.wsConn.Close()
}