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

// OpenTube to put data
func (d *Driver) OpenTube(tube string) {
	if d.DataTubes[tube] == nil {
		d.DataTubes[tube] = &beanstalk.Tube{Conn: d.Conn, Name: tube}
	}
}

// ListenTube prepare tube to be listenned
func (d *Driver) ListenTube(tube string) {
	if d.HandleTubes[tube] == nil {
		d.HandleTubes[tube] = TubeSet(d.Conn, tube)
	}
}

// ReadTube data
func (d *Driver) ReadTube(tube string, timeout time.Duration) (uint64, []byte, error) {
	t, err := d.getHandleTube((tube))
	if err != nil {
		return 0, nil, xerrors.Errorf("reading tube: %w", err)
	}
	return t.Reserve(timeout)
}

// PutDefault - Put message on tube with default configs
func (d *Driver) PutDefault(tube string, body interface{}) (uint64, error) {
	t, err := d.getDataTube(tube)
	if err != nil {
		return 0, xerrors.Errorf("puting message on tube: %w", err)
	}
	return Put(t, PRIORITY_NORMAL, 0, DURATION_DEFAULT, body)
}

// DeleteMessage on tube
func (d *Driver) DeleteMessage(tube string, id uint64) error {
	t, err := d.getDataTube(tube)
	if err != nil {
		return xerrors.Errorf("deleting message on tube: %w", err)
	}
	err = t.Conn.Delete(id)
	if err != nil {
		return xerrors.Errorf("deleting message on tube (%v): %w", id, err)
	}
	return nil
}

// BuryMessage on tube
func (d *Driver) BuryMessage(tube string, id uint64, priority uint32) error {
	t, err := d.getDataTube(tube)
	if err != nil {
		return xerrors.Errorf("burying message on tube [%v]: %w", tube, err)
	}
	err = t.Conn.Bury(id, priority)
	if err != nil {
		return xerrors.Errorf("burying message on tube [%v][%v]: %w", tube, id, err)
	}
	return nil
}

// RenewDuration on tube
func (d *Driver) RenewDuration(tube string, id uint64) error {
	t, err := d.getDataTube(tube)
	if err != nil {
		return xerrors.Errorf("renewing message duration on tube: %w", err)
	}
	err = t.Conn.Touch(id)
	if err != nil {
		return xerrors.Errorf("renewing message duration on tube (%v): %w", id, err)
	}
	return nil
}

func (d *Driver) getHandleTube(tube string) (*beanstalk.TubeSet, error) {
	t := d.HandleTubes[tube]
	if t == nil {
		return nil, xerrors.Errorf("unregistered handle tube")
	}
	return t, nil
}

func (d *Driver) getDataTube(tube string) (*beanstalk.Tube, error) {
	t := d.DataTubes[tube]
	if t == nil {
		return nil, xerrors.Errorf("unregistered data tube")
	}
	return t, nil
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
