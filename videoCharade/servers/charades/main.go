package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"final-project-crew/videoCharade/servers/charades/handlers"
	"final-project-crew/videoCharade/servers/charades/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
)

//main is the main entry point for the server
func main() {

	addr := os.Getenv("CHARADESADDR")

	//if it's blank, default to ":80", which means
	//listen port 80 for requests addressed to any host
	if len(addr) == 0 {
		addr = ":80"
	}

	games := make(map[int64]*models.GameState)

	// Get DSN for leaderboard
	// dbPassword := os.Getenv("MYSQL_ROOT_PASSWORD")
	// if len(dbPassword) == 0 {
	// 	dbPassword = "super-secret"
	// }

	dsn := fmt.Sprintf("root:%s@tcp(%s)/%s?parseTime=true", os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_ADDR"), os.Getenv("MYSQL_DATABASE"))
	// dsn := fmt.Sprintf("root:%s@tcp(mysqldemo:3306)/userDB", dbPassword)

	if len(dsn) == 0 {
		log.Fatal("Leaderboard environment variable not set. Exiting")
	}

	// Get DB connection for leaderboards
	//create a database object, which manages a pool of
	//network connections to the database server
	db, err := connectToSQL(dsn)
	if err != nil {
		os.Exit(1)
		log.Fatalf("error dialing MQ: %v try again!", err)
	}

	//ensure that the database gets closed when we are done
	defer db.Close()

	//for now, just ping the server to ensure we have
	//a live connection to it

	// if err := db.Ping(); err != nil {
	// 	fmt.Printf("error pinging database: %v\n", err)
	// 	log.Fatal("Could not connect to database")
	// } else {
	// 	fmt.Printf("successfully connected to MySQL database\n")
	// }

	mqName := os.Getenv("MQNAME")
	// Set up rabbit mq
	conn, err := connectToMQ(os.Getenv("RABBITADDR"))
	if err != nil {
		log.Fatalf("error dialing MQ: %v try again!", err)
	}
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		mqName, // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Make new context
	ctx := handlers.Context{Games: games, Db: db, Channel: ch, Queue: q}

	mux := http.NewServeMux()
	log.Printf("server is listening at %s...", addr)

	// Handles games
	mux.HandleFunc("/v1/charades", ctx.GamesHandler)

	// Handles guesses
	mux.HandleFunc("/v1/charades/guess", ctx.GuessHandler)

	// Handle skips
	mux.HandleFunc("/v1/charades/skip", ctx.SkipHandler)

	// Handles leaderboards
	mux.HandleFunc("/v1/leaderboards", ctx.LeaderboardHandler)

	log.Fatal(http.ListenAndServe(addr, mux))
}

const maxConnRetries = 5

// connectToMQ connects to the MQ
func connectToMQ(addr string) (*amqp.Connection, error) {
	mqURL := "amqp://" + addr
	var conn *amqp.Connection
	var err error
	for i := 1; i <= maxConnRetries; i++ {
		conn, err = amqp.Dial(mqURL)
		if err == nil {
			log.Printf("successfully connected to %s", mqURL)
			return conn, nil
		}
		log.Printf("error connecting to MQ at %s: %v", mqURL, err)
		log.Printf("will retry in %d seconds", i*2)
		time.Sleep(time.Second * time.Duration(i*2))
	}
	return nil, err
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func connectToSQL(dsn string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for i := 1; i <= maxConnRetries; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			log.Printf("successfully connected to %s", dsn)
			return db, nil
		}
		log.Printf("error connecting to MQ at %s: %v", dsn, err)
		log.Printf("will retry in %d seconds", i*2)
		time.Sleep(time.Second * time.Duration(i*2))
	}
	return nil, err
}
