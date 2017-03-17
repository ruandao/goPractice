package main

import (
	"gopkg.in/mgo.v2"
	"log"
	"github.com/joeshaw/envdecode"
	"github.com/bitly/go-nsq"
	"sync"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var db *mgo.Session

type poll struct {
	Options []string
}

func dialdb() error {
	var err error
	var dburl struct{
		url  string `env:"DB_URL,required"`
	}
	if err = envdecode.Decode(&dburl); err != nil {
		log.Println("can't found ENV $DB_URL, use localhost instead!")
		dburl.url = "localhost"
	}
	log.Printf("dialing mongodb: %s\n", dburl.url)
	db, err = mgo.Dial(dburl.url)
	return err
}

func closedb() {
	db.Close()
	log.Println("closed database connection")
}

func loadOption() ([]string, error) {
	var options []string
	iter := db.DB("ballots").C("polls").Find(nil).Iter()
	var p poll
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}
	iter.Close()
	return options, iter.Err()
}

func publishVotes(votes <-chan string) <-chan struct{} {
	stopchan := make(chan struct{}, 1)
	pub, err := nsq.NewProducer("localhost:4150", nsq.NewConfig())
	if err != nil {
		log.Println("create nsq NewProducer err:", err)
		stopchan <- struct{}{}
		return stopchan
	}
	go func() {
		for vote := range votes {
			pub.Publish("votes", []byte(vote))
		}
		log.Println("Publisher: Stopping")
		pub.Stop()
		log.Println("Publisher: Stopped")
		stopchan <- struct {}{}
	}()
	return stopchan
}

func main() {
	var stoplock sync.Mutex // protects stop
	stop := false
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<- signalChan
		stoplock.Lock()
		stop = true
		stoplock.Unlock()
		log.Println("Stopping...")
		stopChan <- struct{}{}
		closeConn()
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	if err := dialdb(); err != nil {
		log.Fatalln("failed to dial MongoDB:", err)
	}
	defer closedb()

	// start things
	votes := make(chan string) // chan for votes
	publisherStoppedChan := publishVotes(votes)
	twitterStoppedChan := startTwitterStream(stopChan, votes)
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			closeConn()
			stoplock.Lock()
			if stop {
				stoplock.Unlock()
				return
			}
			stoplock.Unlock()
		}
	}()

	<- twitterStoppedChan
	close(votes)
	<- publisherStoppedChan

}
