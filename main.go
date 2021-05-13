package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/ssh"
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
	// ssh -L 60000:127.0.0.1:27017 -p 2233 itm12@117.17.189.6
	listener, err := net.Dial("tcp", tunnel.Local.String())
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Println("listner")
	go tunnel.forward(listener)
	// for {
	// 	fmt.Println("accpet before")
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	fmt.Println("accpet after")
	// 	go tunnel.forward(conn)
	// }
	return nil
}

func (tunnel *SSHtunnel) forward(localConn net.Conn) {
	// serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	// if err != nil {
	// 	fmt.Printf("Server dial error: %s\n", err)
	// 	return
	// }

	// remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	// if err != nil {
	// 	fmt.Printf("Remote dial error: %s\n", err)
	// 	return
	// }

	// copyConn := func(writer, reader net.Conn) {
	// 	defer writer.Close()
	// 	defer reader.Close()

	// 	_, err := io.Copy(writer, reader)
	// 	if err != nil {
	// 		fmt.Printf("io.Copy error: %s", err)
	// 	}
	// }

	// go copyConn(localConn, remoteConn)
	// go copyConn(remoteConn, localConn)

}

func main() {

	// var waitGroup sync.WaitGroup

	localEndpoint := &Endpoint{
		Host: "localhost",
		Port: 60000,
	}

	serverEndpoint := &Endpoint{
		Host: "117.17.189.6",
		Port: 2233,
	}

	remoteEndpoint := &Endpoint{
		Host: "localhost",
		Port: 27017,
	}

	sshConfig := &ssh.ClientConfig{
		User: "itm12",
		Auth: []ssh.AuthMethod{
			ssh.Password("itmserver"),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	t := &SSHtunnel{
		Config: sshConfig,
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}

	fmt.Println("tunnel", t)
	// waitGroup.Add(1)
	t.Start()
	// waitGroup.Wait()

	// time.Sleep(time.Second * 3)
	fmt.Println("start connection tunneling")
	conn := fmt.Sprintf("mongodb://localhost:%d", t.Local.Port)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	credential := options.Credential{
		Username: "root",
		Password: "test123",
	}

	clientOptions := options.Client()
	clientOptions.ApplyURI(conn).SetAuth(credential)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Client", client)

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	logsCollection := client.Database("log_db").Collection("logs")
	// find, _ := logsCollection.Find(context.TODO(), bson.D{})
	// fmt.Println("find", find)

	// var data []bson.D
	// if err = find.All(context.TODO(), &data); err != nil {
	// 	fmt.Println("find err", err)
	// }

	// fmt.Println(data[0][0].Value)

	// deleteResult, err := logsCollection.DeleteMany(context.TODO(), bson.D{})
	// if err != nil {
	// 	fmt.Println("delete err", err)

	// }

	// fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)

	// read file csv
	var waitGroup sync.WaitGroup
	files, err := ioutil.ReadDir("./bigdatasample")
	if err != nil {
		log.Fatal(err)
	}
	waitGroup.Add(len(files))
	var logs []interface{}
	for _, fileName := range files {
		go func(name string) {
			defer waitGroup.Done()
			ld := ReadCSV(name)
			logs = append(logs, ld...)
			log.Println(name, "finish")
		}(fileName.Name())
	}
	waitGroup.Wait()

	start := time.Now()

	insertManyResult, err := logsCollection.InsertMany(context.TODO(), logs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
	elapsed := time.Since(start)
	log.Printf("insertMany took %s", elapsed)
}
