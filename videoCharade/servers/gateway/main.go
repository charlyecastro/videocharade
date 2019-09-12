package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"homework-charlyecastro/servers/gateway/handlers"
	"homework-charlyecastro/servers/gateway/models/users"
	"homework-charlyecastro/servers/gateway/sessions"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
)

//main is the main entry point for the server
const maxConnRetries = 5

func main() {

	//Setting up ADDR
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	//Setting up Session Key
	sessionKey := os.Getenv("SESSIONKEY")

	//Setting up RedisStore
	redisAddr := os.Getenv("REDISADDR")
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	sessionStore := sessions.NewRedisStore(client, time.Hour)

	//Setting up SQL
	dsn := fmt.Sprintf("root:%s@tcp(%s)/%s?parseTime=true", os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_ADDR"), os.Getenv("MYSQL_DATABASE"))

	db, err := connectToSQL(dsn)
	defer db.Close()
	if err != nil {
		os.Exit(1)
		log.Fatalf("error dialing MQ: %v try again!", err)
	}

	userStore := users.NewMySQLStore(db)
	userTrie, err := userStore.LoadAllUsers()

	if err != nil {
		log.Fatalf("Didnt have enough timet o load all users! becuse %s try again!", err)
	}
	notif := handlers.NewNotifier()

	//Set up SharedResource
	ctx := handlers.NewContext(sessionKey, sessionStore, userStore, userTrie, notif)

	conn, err := connectToMQ(os.Getenv("RABBITADDR"))
	if err != nil {
		log.Fatalf("error dialing MQ: %v try again!", err)
	}
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"Messaging", // name
		true,        // durable
		false,       // delete when usused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	go ctx.Process(msgs)

	//Setting up Certificates
	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")
	defer os.Exit(1)
	if len(tlsKeyPath) == 0 || len(tlsCertPath) == 0 {
		log.Fatal("tlsCertPath or tlsKeyPath does not exist")
	}

	//Set up Proxy Servers
	messagesAddr := os.Getenv("MESSAGESADDR")
	summaryAddr := os.Getenv("SUMMARYADDR")
	charadesADDR := os.Getenv("CHARADESADDR")
	fmt.Println("charadesAddr: ", charadesADDR)

	summaryProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: summaryAddr})
	charadeProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: charadesADDR})

	//Setting up mux routes
	mux := http.NewServeMux()
	mux.Handle("/v1/summary", summaryProxy)
	mux.Handle("/v1/charades", charadeProxy)
	mux.Handle("/v1/charades/guess", charadeProxy)
	mux.Handle("/v1/charades/skip", charadeProxy)
	mux.Handle("/v1/leaderboards", charadeProxy)
	mux.HandleFunc("/v1/users", ctx.UsersHandler)
	mux.HandleFunc("/v1/users/", ctx.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", ctx.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", ctx.SpecificSessionHandler)
	mux.Handle("/v1/channels", NewServiceProxy(messagesAddr, ctx))
	mux.Handle("/v1/channels/", NewServiceProxy(messagesAddr, ctx))
	mux.Handle("/v1/messages/", NewServiceProxy(messagesAddr, ctx))
	mux.Handle("/v1/handleOffer", NewServiceProxy(messagesAddr, ctx))
	mux.HandleFunc("/ws", ctx.WebSocketConnectionHandler)

	//Serving Wrapped Mux
	wrappedMux := handlers.NewCors(mux)
	log.Printf("wrapped server is listening at https://%s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMux))
}

const headerUser = "X-User"

//NewServiceProxy returns a new proxy
func NewServiceProxy(addrs string, ctx *handlers.Context) *httputil.ReverseProxy {
	splitAddrs := strings.Split(addrs, ",")
	nextAddr := 0
	mx := sync.Mutex{}

	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			mx.Lock()
			r.URL.Host = splitAddrs[nextAddr]
			nextAddr = (nextAddr + 1) % len(splitAddrs)
			mx.Unlock()
			r.Header.Del(headerUser)

			var state handlers.SessionState
			_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, &state)
			if err != nil {
				return
			}

			user := state.User
			userJSON, _ := json.Marshal(user)
			r.Header.Set(headerUser, string(userJSON))
		},
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

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
