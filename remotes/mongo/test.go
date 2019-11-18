package mongo

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net"
)

func teste() {
	dialInfo, err := mgo.ParseURL("")
	if err != nil {
		log.Println(err)
	}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		tlsConfig := &tls.Config{}
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		fmt.Println("CONEXAO: ",conn)
		if err != nil {
			log.Println("ERRO: ",err)
		}
		return conn, err
	}
	s, err := mgo.DialWithInfo(dialInfo)
	fmt.Println("SESSION: ",s)
	if(err != nil) {
		fmt.Println("ERRO: ", err)
	}

	collection := s.DB("testing").C("numbers")
	collection.Insert(bson.M{"pis":222222})
}

func main(){
	teste()
}