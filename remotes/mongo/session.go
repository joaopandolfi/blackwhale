package mongo

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"

	"github.com/joaopandolfi/blackwhale/configurations"
	"github.com/joaopandolfi/blackwhale/utils"
	"gopkg.in/mgo.v2"
)

type Session struct {
	session *mgo.Session
}

var session Session
var pool []Session
var looper int

var maxPool int = configurations.Configuration.MongoPool

// Create session with ssl and ignore the validation cert (more common)
func NewSessionSsl() (s *mgo.Session, err error) {
	if session.session == nil {
		url := strings.Replace(configurations.Configuration.MongoUrl, "ssl=true", "", -1)
		url = strings.Replace(url, "readPreference=secondaryPreferred", "", -1)
		dialInfo, err := mgo.ParseURL(url)
		if err != nil {
			utils.CriticalError("[Mongo SSL] ERROR Url parsing", err)
		}
		//utils.Debug("[Mongo] - Before connection")
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			tlsConfig := &tls.Config{}
			tlsConfig.InsecureSkipVerify = true
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			if err != nil {
				utils.CriticalError("[Mongo SSL] ERROR SSL Connection ", err.Error(), addr.String())
				log.Println(err)
				panic(err)
			}
			return conn, err
		}
		s, err = mgo.DialWithInfo(dialInfo)

		if err != nil {
			return nil, err
		}

		//session.session.SetMode(mgo.SecondaryPreferred,true) // Unecessary
	}
	return s, err
}

// Create session with ssl and use sign cert
func NewSessionSSLMETHOD2() (s *Session, err error) {
	// --sslCAFile
	rootCerts := x509.NewCertPool()
	if ca, err := ioutil.ReadFile("ca.crt"); err == nil {
		rootCerts.AppendCertsFromPEM(ca)
	}

	// --sslPEMKeyFile
	clientCerts := []tls.Certificate{}
	if cert, err := tls.LoadX509KeyPair("client.crt", "client.key"); err == nil {
		clientCerts = append(clientCerts, cert)
	}

	// Dial with TLS
	sess, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs: []string{configurations.Configuration.MongoUrl},
		DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{
				RootCAs:      rootCerts,
				Certificates: clientCerts,
			})
		},
	})
	session.session = sess

	return &session, err
}

// Create session without ssl
func newSession() (s *mgo.Session, err error) {
	se, err := mgo.Dial(configurations.Configuration.MongoUrl)
	if err != nil {
		return nil, err
	}
	return se, err
}

func GetPoolSession() (s *Session, err error) {
	lenPool := len(pool)

	if lenPool <= maxPool {
		looper = lenPool
		var s *mgo.Session

		if strings.Contains(configurations.Configuration.MongoUrl, "ssl=") {
			s, err = NewSessionSsl()
		} else {
			s, err = newSession()
		}

		if err != nil {
			FlushPull()
			panic(err)
		}

		pool = append(pool, Session{session: s})
	} else {
		if looper >= maxPool {
			looper = 0
		} else {
			looper++
		}
	}

	return &pool[looper], err
}

func FlushPull() {
	for _, p := range pool {
		go p.Close()
	}
	pool = nil
	looper = 0
}

// https://godoc.org/gopkg.in/mgo.v2#Dial
func NewSession() (s *Session, err error) {
	if session.session == nil {
		if strings.Contains(configurations.Configuration.MongoUrl, "ssl=") {
			session.session, err = NewSessionSsl()
		} else {
			session.session, err = newSession()
		}
	}
	return &session, err
}

func (s *Session) Copy() *Session {
	return &Session{s.session.Copy()}
}

// Get mongo collection
func (s *Session) GetCollection(col string) *mgo.Collection {
	return s.session.DB(configurations.Configuration.MongoDb).C(col)
}

func (s *Session) Run(cmd interface{}) {
	var result interface{}
	s.session.DB(configurations.Configuration.MongoDb).Run(cmd, result) //.C(col)
	fmt.Println(result)
}

func (s *Session) Close() {
	if s.session != nil {
		s.session.Close()
	}
}
