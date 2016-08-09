package main

/*
 * target.go
 * Convert targets into IPs
 * By J. Stuart McMurray
 * Created 20160808
 * Last Modified 20160808
 */

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

/* TargetDelay copies string from in to out, with a delay of Config.Delay,
plus or minus 50%. */
func TargetDelay(out chan<- string, in <-chan string) {
	/* Pull each target from in */
	for t := range in {
		/* Send it to out */
		out <- t
		/* Work out how much to really delay */
		sf := 0.5 + rand.Float64()
		d := time.Duration(sf * float64(CONFIG.Delay))
		time.Sleep(d)
	}
	/* Close the output channel when there's nothing left */
	close(out)
}

/* SeedPRNG seeds the random number generator with a cryptographically-soundly
random seed.  It panics on error */
func SeedPRNG() {
	b := make([]byte, 8)
	_, err := crand.Read(b)
	if nil != err {
		panic(fmt.Sprintf("Rand: %v", err))
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b)))
}

/* Targeter parses the target t and sends IPs to be attacked to tch.  wg's Done
method will be called before returning. */
func Targeter(tch chan<- string, t string, wg *sync.WaitGroup) {
	defer wg.Done()
	/* Maybe it's a CIDR range? */
	if _, ipnet, err := net.ParseCIDR(t); nil == err {
		targetCIDR(tch, ipnet)
		return
	}

	/* Look up the addresses associated with the target */
	ips, err := net.LookupIP(t)
	if nil != err {
		log.Printf("Unable to resolve %v: %v", t, err)
		return
	}
	/* Send them out */
	for _, ip := range ips {
		tch <- ip.String()
	}
}

/* targetCIDR iterates through the addresses in ipnet, and sends the out on
tch */
func targetCIDR(tch chan<- string, ipnet *net.IPNet) {
	for ip := ipnet.IP; ipnet.Contains(ip); inc(ip) {
		tch <- ip.String()
	}
}

/* inc increments ip */
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
