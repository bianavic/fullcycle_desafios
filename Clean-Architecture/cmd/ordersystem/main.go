package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/bianavic/fullcycle_clean-architecture/configs"
	"github.com/bianavic/fullcycle_clean-architecture/internal/event/handler/order_created_handler"
	"github.com/bianavic/fullcycle_clean-architecture/internal/infra/graph"
	"github.com/bianavic/fullcycle_clean-architecture/internal/infra/grpc/pb"
	"github.com/bianavic/fullcycle_clean-architecture/internal/infra/grpc/service"
	"github.com/bianavic/fullcycle_clean-architecture/internal/infra/web/webserver"
	"github.com/bianavic/fullcycle_clean-architecture/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	fmt.Print("Connecting to database...\n")
	db, err := sql.Open(configs.DBDriver,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			configs.DBUser, configs.DBPassword,
			configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic("Database connection failed: " + err.Error())
	}

	// Initialize RabbitMQ
	rabbitMQChannel := getRabbitMQChannel()
	defer rabbitMQChannel.Close()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &order_created_handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)

	// Initialize Web Server
	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)

	// Register handler
	webserver.AddHandler("/order/create", webOrderHandler.Create)
	webserver.AddHandler("/order", webOrderHandler.List)
	fmt.Println("Starting web server on port", configs.WebServerPort)

	// Start servers
	go webserver.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *eventDispatcher)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		EventDispatcher:    *eventDispatcher,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Printf("🚀 GraphQL playground rodando em http://localhost:%s/\n", configs.GraphQLServerPort)
	err = http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
	if err != nil {
		log.Fatalf("Failed to start GraphQL server: %v", err)
	}

	go consumeMessages()
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}

func consumeMessages() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"orders", // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume(
		"orders", // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			fmt.Printf("Received a message: %s\n", d.Body)
		}
	}()

	fmt.Println("Waiting for messages. To exit press CTRL+C")
	<-forever
}
