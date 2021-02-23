package graphite

import (
	"fmt"
	"strconv"

	conf "github.com/joaopandolfi/blackwhale/configurations"
	graphite "github.com/marpaia/graphite-golang"
)

var _host string
var _port int

// Driver graphite
type Driver struct {
	Conn   *graphite.Graphite
	Prefix string
}

var conn *graphite.Graphite

func init() {
	_host = conf.Configuration.GraphiteUrl
	p, _ := strconv.Atoi(conf.Configuration.GraphitePort)
	_port = p
}

// SetCredentials to connection
func SetCredentials(host string, port int) {
	_host = host
	_port = port
}

// New Graphite driver
func New(prefix string) (*Driver, error) {
	if conn == nil {
		c, err := graphite.NewGraphite(_host, _port)
		if err != nil {
			c = graphite.NewGraphiteNop(_host, _port)
		}
		conn = c
	}
	return &Driver{
		Conn:   conn,
		Prefix: prefix,
	}, nil
}

// Send to graphite
func (d *Driver) Send(key, data string) error {
	return d.Conn.SimpleSend(fmt.Sprintf("stats.%s.%s", d.Prefix, key), data)
}
