package beanstalkd

import (
	"encoding/json"
	"time"

	c "github.com/joaopandolfi/blackwhale/configurations"
	"github.com/kr/beanstalk"
	"golang.org/x/xerrors"
)

// Driver -
type Driver struct {
	Conn        *beanstalk.Conn
	DataTubes   map[string]*beanstalk.Tube
	HandleTubes map[string]*beanstalk.TubeSet
}

const PRIORITY_HIGHEST uint32 = 256
const PRIORITY_HIGH uint32 = 512
const PRIORITY_NORMAL uint32 = 1024
const PRIORITY_LOW uint32 = 2048
const PRIORITY_LOWEST uint32 = 4096
const DURATION_DEFAULT time.Duration = time.Second * 120
const TIMEOUT_DEFAULT time.Duration = 0

var conn *beanstalk.Conn

// TODO: MAKE RESILIENT

// New beanstalkd driver
func New() (*Driver, error) {
	if conn == nil {
		c, err := beanstalk.Dial("tcp", c.Configuration.BeanstalkdUrl)
		if err != nil {
			return nil, xerrors.Errorf("connecting on beanstalkd: %w", err)
		}
		conn = c
	}
	return &Driver{
		Conn:        conn,
		DataTubes:   map[string]*beanstalk.Tube{},
		HandleTubes: map[string]*beanstalk.TubeSet{},
	}, nil
}

// Fresh return a new connection
func Fresh() (*Driver, error) {

	c, err := beanstalk.Dial("tcp", c.Configuration.BeanstalkdUrl)
	if err != nil {
		return nil, xerrors.Errorf("connecting on beanstalkd: %w", err)
	}
	return &Driver{
		Conn:        c,
		DataTubes:   map[string]*beanstalk.Tube{},
		HandleTubes: map[string]*beanstalk.TubeSet{},
	}, nil
}

// ReadTube data
func (d *Driver) ReadTube(tube string, timeout time.Duration) (uint64, []byte, error) {
	t := d.getHandleTube(tube)
	return t.Reserve(timeout)
}

// PutDefault - Put message on tube with default configs
func (d *Driver) PutDefault(tube string, body interface{}) (uint64, error) {
	t := d.getDataTube(tube)
	return Put(t, PRIORITY_NORMAL, 0, DURATION_DEFAULT, body)
}

// DeleteMessage on tube
func (d *Driver) DeleteMessage(tube string, id uint64) error {
	t := d.getDataTube(tube)
	err := t.Conn.Delete(id)
	if err != nil {
		return xerrors.Errorf("deleting message on tube (%v): %w", id, err)
	}
	return nil
}

// BuryMessage on tube
func (d *Driver) BuryMessage(tube string, id uint64, priority uint32) error {
	t := d.getDataTube(tube)
	err := t.Conn.Bury(id, priority)
	if err != nil {
		return xerrors.Errorf("burying message on tube [%v][%v]: %w", tube, id, err)
	}
	return nil
}

// RenewDuration on tube
func (d *Driver) RenewDuration(tube string, id uint64) error {
	t := d.getDataTube(tube)
	err := t.Conn.Touch(id)
	if err != nil {
		return xerrors.Errorf("renewing message duration on tube (%v): %w", id, err)
	}
	return nil
}

func (d *Driver) getHandleTube(tube string) *beanstalk.TubeSet {
	if d.HandleTubes[tube] == nil {
		d.HandleTubes[tube] = TubeSet(d.Conn, tube)
	}
	return d.HandleTubes[tube]
}

func (d *Driver) getDataTube(tube string) *beanstalk.Tube {

	if d.DataTubes[tube] == nil {
		d.DataTubes[tube] = &beanstalk.Tube{Conn: d.Conn, Name: tube}
	}
	return d.DataTubes[tube]
}

// TubeSet Open
func TubeSet(conn *beanstalk.Conn, tube string) *beanstalk.TubeSet {
	return beanstalk.NewTubeSet(conn, tube)
}

// Put - put some data on defined tube
func Put(tube *beanstalk.Tube, priority uint32, delay, ttr time.Duration, body interface{}) (uint64, error) {
	bbody, err := json.Marshal(body)
	if err != nil {
		return 0, xerrors.Errorf("marshaling body to put on tube (%s): %w", tube.Name, err)
	}
	return tube.Put(bbody, priority, 0, ttr)
}
