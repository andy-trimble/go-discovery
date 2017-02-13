package discovery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/kelseyhightower/envconfig"
	uuid "github.com/satori/go.uuid"
)

// Actor is the struct representing discovered entities
type Actor struct {
	Role string    `json:"role"`
	IP   string    `json:"ip"`
	ID   uuid.UUID `json:"id"`
}

type discovery interface {
	Start(role string)
	Shutdown()
}

type config struct {
	Interface     string        `default:"eth0"`
	Group         string        `default:"230.1.1.1"`
	Port          int           `default:"8989"`
	AnnounceCount int           `default:"5"`
	AnnounceWait  time.Duration `default:"500ms"`
}

// Discovery is the multicast discovery agent
type Discovery struct {
	actor      *Actor
	c          config
	in         *net.UDPConn
	out        *net.UDPConn
	cache      map[string]*Actor
	Discovered chan Actor
	Err        chan error
}

// Start initiates the discovery process
func (d *Discovery) Start(role string) error {
	err := envconfig.Process("DISCOVERY", &d.c)
	if err != nil {
		return err
	}
	ifc, err := net.InterfaceByName(d.c.Interface)
	if err != nil {
		return err
	}
	addrs, err := ifc.Addrs()
	if err != nil {
		return err
	}
	var ip string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	if ip == "" {
		return fmt.Errorf("Could not locate IP address")
	}

	d.cache = make(map[string]*Actor, 10)
	d.Discovered = make(chan Actor)
	d.Err = make(chan error)
	d.actor = &Actor{Role: role, IP: ip, ID: uuid.NewV4()}
	d.cache[d.actor.ID.String()] = d.actor

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", d.c.Group, d.c.Port))
	if err != nil {
		return err
	}
	go d.announce(addr)
	go d.serve(addr)

	return nil
}

func (d *Discovery) serve(addr *net.UDPAddr) {
	c, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		d.Err <- err
		return
	}
	d.in = c
	d.in.SetReadBuffer(512)
	for {
		b := make([]byte, 512)
		_, _, err := d.in.ReadFromUDP(b)
		if err != nil {
			d.Err <- err
			continue
		}
		var a Actor
		err = json.NewDecoder(bytes.NewReader(b)).Decode(&a)
		if err != nil {
			d.Err <- err
			continue
		}
		if _, ok := d.cache[a.ID.String()]; !ok {
			d.cache[a.ID.String()] = &a
			d.Discovered <- a
			err = d.sendMe()
			if err != nil {
				d.Err <- err
			}
		}
	}
}

func (d *Discovery) announce(addr *net.UDPAddr) {
	c, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		d.Err <- err
		return
	}
	d.out = c
	for i := 0; i < d.c.AnnounceCount; i++ {
		err = d.sendMe()
		if err != nil {
			d.Err <- err
			continue
		}
		time.Sleep(d.c.AnnounceWait)
	}
}

func (d *Discovery) sendMe() error {
	b, err := json.Marshal(d.actor)
	if err != nil {
		return err
	}
	d.out.Write(b)
	return nil
}

// Shutdown closes all connections and channels
func (d *Discovery) Shutdown() {
	err := d.in.Close()
	if err != nil {
		d.Err <- err
	}

	err = d.out.Close()
	if err != nil {
		d.Err <- err
	}

	close(d.Discovered)
	close(d.Err)
}
