package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"github.com/elliotchance/sshtunnel"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// COUNT,LOGTIME,DEVICE_NAME,ATTACK_NAME,RAW_PACKET
type logEntity struct {
	count       int    `bson:"count"`
	logtime     string `bson:"logtime"`
	attack_name string `bson:"attack_name"`
	raw_packet  string `bson:"raw_packet"`
}

type Endpoint struct {
	Host string
	Port int
}

type SSHtunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint

	Config *ssh.ClientConfig
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

func (tunnel *SSHtunnel) Start() error {
	listner, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		return err
	}

	defer listner.Close()

	for {
		conn, err := listner.Accept()
		if err != nil {
			return err
		}

		go tunnel.forward(conn)
	}
}

func (tunnel *SSHtunnel) forward(localConn net.Conn) {
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
}

func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func main() {
	tunnel := sshtunnel.NewSSHTunnel(
		"itm12@117.17.189.6:2233"
	)
	
	ctx, _ := context.WithTimeout(context.Background(), 10)
	clientOptions := options.Client().ApplyURI("mongodb://117.17.189.6:27017").SetAuth(options.Credential{
		AuthSource: "",
		Username:   "root",
		Password:   "test123",
	})

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("몽고 DB에 연결했습니다!")

	// logsCollection := client.Database("log_db").Collection("logs")
	// insertResult, _ := logsCollection.InsertOne(context.TODO(), bson.D{
	// 	{"userID", "test1234"},
	// 	{"array", bson.A{"flying", "squirrel", "dev"}},
	// })
	// fmt.Println(insertResult)

}
