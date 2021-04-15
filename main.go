package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

type SSHtunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint

	Config *ssh.ClientConfig
}

func (tunnel *SSHtunnel) Start() error {
	fmt.Println("listen")
	// ssh port 10000 이상으로 체크하기
	// listener, err := net.Listen("tcp", ":60000")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	// fmt.Println("first")
	// fmt.Println(listener.Addr().String())

	// defer listener.Close()
	// fmt.Println("second")
	// fmt.Println(listener.Addr().String())
	conn, err := net.Dial("tcp", "localhost:60000")
	// conn, err := listener.Accept()
	if err != nil {
		fmt.Println("connection closed")
		return err
	}

	fmt.Println("123")

	go tunnel.forward(conn)

	return nil
}

func (tunnel *SSHtunnel) forward(localConn net.Conn) {
	fmt.Println("test")
	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		fmt.Printf("Server dial error: %s\n", err)
		return
	}

	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		fmt.Printf("Remote dial error: %s\n", err)
		return
	}

	fmt.Println("remote", remoteConn)

	copyConn := func(writer, reader net.Conn) {
		defer writer.Close()
		defer reader.Close()

		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)

	defer serverConn.Close()
	defer remoteConn.Close()

	sshClient, err := agent.NewClient(remoteConn), err
	if err != nil {
		fmt.Println("ssh", err)
		return
	}

	fmt.Println("sshclient", sshClient)

	fmt.Println(ssh.PublicKeysCallback(sshClient.Signers))

	// var waitGroup sync.WaitGroup

	// waitGroup.Add(2)

	// go func(c net.Conn) {
	// 	i := 0
	// 	for {
	// 		s := "mongodb"

	// 		r, err := c.Write([]byte(s)) // 서버로 데이터를 보냄
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			return
	// 		}

	// 		fmt.Println("send", r)

	// 		i++
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }(remoteConn)

	// go func(c net.Conn) {
	// 	data := make([]byte, 4096) // 4096 크기의 바이트 슬라이스 생성

	// 	for {
	// 		n, err := c.Read(data) // 서버에서 받은 데이터를 읽음
	// 		if err != nil {
	// 			fmt.Println("return", err)
	// 			return
	// 		}

	// 		fmt.Println("data", string(data[:n])) // 데이터 출력
	// 		log.Println("Server send : " + string(data[:n]))

	// 		time.Sleep(1 * time.Second)
	// 	}
	// }(remoteConn)

	// waitGroup.Wait()

	// fmt.Scanln()

	// ctx, _ := context.WithTimeout(context.Background(), 10)
	// clientOptions := options.Client().ApplyURI("mongodb://root:test1234@localhost:27017/?connect=direct&sslInsecure=true")

	// client, err := mongo.Connect(ctx, clientOptions)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("client", client)

	// err = client.Ping(context.Background(), nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func main() {

	var waitGroup sync.WaitGroup

	localEndpoint := &Endpoint{
		Host: "122.45.88.38",
		Port: 60000,
	}

	serverEndpoint := &Endpoint{
		Host: "117.17.189.6",
		Port: 2233,
	}

	remoteEndpoint := &Endpoint{
		Host: "117.17.189.6",
		Port: 27017,
	}

	sshConfig := &ssh.ClientConfig{
		User: "itm12",
		Auth: []ssh.AuthMethod{
			ssh.Password("itmserver"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	t := &SSHtunnel{
		Config: sshConfig,
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}

	fmt.Println("tunnel", t)
	waitGroup.Add(1)
	t.Start()
	// time.Sleep(time.Second * 3)
	fmt.Println("start connection tunneling")
	conn := fmt.Sprintf("mongodb://localhost:%d", 60000)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(conn)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	logsCollection := client.Database("log_db").Collection("logs")
	find, _ := logsCollection.Find(context.TODO(), bson.D{})
	fmt.Println(find)
	insertResult, _ := logsCollection.InsertOne(context.TODO(), bson.D{
		{"userID", "test1234"},
		{"array", bson.A{"flying", "squirrel", "dev"}},
	})
	fmt.Println("insert", insertResult)

	waitGroup.Wait()

}

// func main() {

// 	// fmt.Println(tunnel)

// 	// log for debugging
// 	tunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)

// 	// start sterver background
// 	go tunnel.Start()
// 	time.Sleep(100 * time.Millisecond)
// 	// conn := fmt.Sprintf("host=127.0.0.1 port=%d username=foo", tunnel.Local.Port)
// 	// fmt.Println(conn)
// 	// ctx, _ := context.WithTimeout(context.Background(), 10)
// 	// clientOptions := options.Client().ApplyURI("mongodb://root:test1234@117.17.189.6:27017").SetAuth(options.Credential{})

// 	// client, err := mongo.Connect(ctx, clientOptions)

// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// err = client.Ping(context.Background(), nil)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// logsCollection := client.Database("log_db").Collection("logs")
// 	// insertResult, _ := logsCollection.InsertOne(context.TODO(), bson.D{
// 	// 	{"userID", "test1234"},
// 	// 	{"array", bson.A{"flying", "squirrel", "dev"}},
// 	// })
// 	// fmt.Println(insertResult)

// }
