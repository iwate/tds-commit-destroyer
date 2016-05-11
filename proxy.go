package proxy

import (
	"io"
	"net"
	"os"
)

const TRAN_MGR_REQ = 0x0e
const TM_COMMIT_XACT = 0x07
const SERVER_RES = 0x04
const ENV_CHANGE = 0xe3
const COMMIT_TRAN = 0x09

const PACKET_TYPE = 0
const TOKEN_TYPE = 8
const ENV_CODE = 11
const TM_TYPE = 30

// Proxy - Manages a Proxy connection, piping data between local and remote.
type Proxy struct {
	sentBytes     uint64
	receivedBytes uint64
	laddr, raddr  *net.TCPAddr
	lconn, rconn  io.ReadWriteCloser
	erred         bool
	errsig        chan bool
	
	// Settings
	Nagles    bool
	Log       Logger
	OutputHex bool
}

// New - Create a new Proxy instance. Takes over local connection passed in,
// and closes it when finished.
func New(lconn *net.TCPConn, laddr, raddr *net.TCPAddr) *Proxy {
	return &Proxy{
		lconn:  lconn,
		laddr:  laddr,
		raddr:  raddr,
		erred:  false,
		errsig: make(chan bool),
		Log:    NullLogger{},
	}
}

// Start - open connection to remote and start proxying data.
func (p *Proxy) Start() {
	defer p.lconn.Close()

	var err error
	//connect to remote
	p.rconn, err = net.DialTCP("tcp", nil, p.raddr)
	if err != nil {
		p.err("Remote connection failed: %s", err)
		return
	}
	defer p.rconn.Close()

	//display both ends
	p.Log.Info("Opened %s >>> %s", p.laddr.String(), p.raddr.String())

	//bidirectional copy
	go p.pipe(p.lconn, p.rconn)
	go p.pipe(p.rconn, p.lconn)

	//wait for close...
	<-p.errsig
	p.Log.Info("Closed (%d bytes sent, %d bytes recieved)", p.sentBytes, p.receivedBytes)
}

func (p *Proxy) err(s string, err error) {
	if p.erred {
		return
	}
	if err != io.EOF {
		p.Log.Warn(s, err)
	}
	p.errsig <- true
	p.erred = true
}

func (p *Proxy) pipe(src, dst io.ReadWriter) {
	islocal := src == p.lconn

	var dataDirection string
	if islocal {
		dataDirection = ">>> %d bytes sent%s"
	} else {
		dataDirection = "<<< %d bytes recieved%s"
	}

	var byteFormat string
	if p.OutputHex {
		byteFormat = "%x"
	} else {
		byteFormat = "%s"
	}

	//directional copy (64k buffer)
	buff := make([]byte, 0xffff)
	for {
		n, err := src.Read(buff)
		if err != nil {
			p.err("Read failed '%s'\n", err)
			return
		}
		b := buff[:n]
		
		//show output
		p.Log.Debug(dataDirection, n, "")
		p.Log.Trace(byteFormat, b)
    
    if b[PACKET_TYPE] == SERVER_RES && b[TOKEN_TYPE] == ENV_CHANGE && b[ENV_CODE] == COMMIT_TRAN {
			p.Log.Info("!!!!DESTROY!!!!!")
			os.Exit(0)
		}
    dst.Write(b)
	}
}