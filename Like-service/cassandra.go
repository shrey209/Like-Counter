package main

import (
	"log"
	"time"

	"github.com/gocql/gocql"
)

func NewSession(host string, keyspace string) *gocql.Session {
	cluster := gocql.NewCluster(host)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.ConnectTimeout = time.Second * 10

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}

	return session
}
