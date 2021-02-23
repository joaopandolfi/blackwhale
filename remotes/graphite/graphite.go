package graphite

import (
	"fmt"
	"strconv"

	conf "github.com/joaopandolfi/blackwhale/configurations"
	graphite "github.com/marpaia/graphite-golang"
)

// Driver graphite
type Driver struct {
	Conn   *graphite.Graphite
	Prefix string
}

var conn *graphite.Graphite

// New Graphite driver
func New(prefix string) (*Driver, error) {
	port, _ := strconv.Atoi(conf.Configuration.GraphitePort)
	if conn == nil {
		c, err := graphite.NewGraphite(conf.Configuration.GraphiteUrl, port)
		if err != nil {
			c = graphite.NewGraphiteNop(conf.Configuration.GraphiteUrl, port)
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
