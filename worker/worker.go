package main

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"strconv"
)

func main() {

	client := NewClient()

	//fmt.Println(client.PSubscribe("insert").)

	//c := make(chan interface{})
	//
	//go (func(c chan interface{}) {
	//	for {
	//		c <- client.Subscribe("insert")
	//		fmt.Println(c)
	//	}
	//})(c)
	fmt.Println(0)

	//go func() {
	//	fmt.Println(1)
	//	insertSub := client.Subscribe("insert")
	//	fmt.Println(2)
	//	defer insertSub.Close()
	//	fmt.Println(3)
	//	for {
	//		fmt.Println(4)
	//		in, err := insertSub.Receive()  // *** SOMETIMES RETURNS EOF ERROR ***
	//		fmt.Println(5)
	//		if err != nil {
	//			fmt.Println("failed to get feedback", err)
	//			break
	//		}
	//		fmt.Println("===============message received==============")
	//
	//		switch in.(type) {
	//		case *redis.Message:
	//			fmt.Println("===============MESSAGE RECEIVED==============")
	//			//cm := comm.ControlMessageEvent{}
	//			//payload := []byte(in.(*redis.Message).Payload)
	//			//if err := json.Unmarshal(payload, &cm); err != nil {
	//			//	fmt.Println("failed to parse control message", err)
	//			//} else if err := handleIncomingEvent(&cm); err != nil {
	//			//	fmt.Println("failed to handle control message", err)
	//			//}
	//
	//		default:
	//			fmt.Println("Received unknown input over REDIS PubSub control channel", " | received: ", in)
	//		}
	//	}
	//}()

		pubsub := client.Subscribe("insert")
		defer pubsub.Close()

		if _, err := pubsub.Receive(); err != nil {
			fmt.Println("failed to receive from control PubSub: ", err)
			return
		}

		controlCh := pubsub.Channel()
		fmt.Println("start listening on control PubSub")

		// Endlessly listen to control channel,
		for msg := range controlCh {
			fmt.Println(msg.Payload)
			index, err := strconv.Atoi(msg.Payload)
			if err != nil {
				panic(err)
			}
			fibResult := fib(index)
			client.HSet("values", msg.Payload, fibResult)
		}

}

func fib(index int) int {
	if index < 2 {
		return 1
	}
	return fib(index - 1) + fib(index - 2)
}

func NewClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
	return client
}


// func InitAndListenAsync(log *log.Logger, sseHandler func(string, string) error) error {
//    rootLogger = log.With(zap.String("component", "redis-client"))
//
//    host := env.RedisHost
//    port := env.RedisPort
//    pass := env.RedisPass
//    addr := fmt.Sprintf("%s:%s", host, port)
//    tlsCfg := &tls.Config{}
//    client = redis.NewClient(&redis.Options{
//        Addr:      addr,
//        Password:  pass,
//        TLSConfig: tlsCfg,
//    })
//
//    if _, err := client.Ping().Result(); err != nil {
//        return err
//    }
//
//    go func() {
//        controlSub := client.Subscribe("control")
//        defer controlSub.Close()
//        for {
//            in, err := controlSub.Receive()  // *** SOMETIMES RETURNS EOF ERROR ***
//            if err != nil {
//                rootLogger.Error("failed to get feedback", zap.Error(err))
//                break
//            }
//            switch in.(type) {
//            case *redis.Message:
//                cm := comm.ControlMessageEvent{}
//                payload := []byte(in.(*redis.Message).Payload)
//                if err := json.Unmarshal(payload, &cm); err != nil {
//                    rootLogger.Error("failed to parse control message", zap.Error(err))
//                } else if err := handleIncomingEvent(&cm); err != nil {
//                    rootLogger.Error("failed to handle control message", zap.Error(err))
//                }
//
//            default:
//                rootLogger.Warn("Received unknown input over REDIS PubSub control channel", zap.Any("received", in))
//            }
//        }
//    }()
//    return nil
//}

//func listenToControlChannel(client *redis.Client) {
//    pubsub := client.Subscribe("control")
//    defer pubsub.Close()
//
//    if _, err := pubsub.Receive(); err != nil {
//        rootLogger.Error("failed to receive from control PubSub", zap.Error(err))
//        return
//    }
//
//    controlCh := pubsub.Channel()
//    fmt.Println("start listening on control PubSub")
//
//    // Endlessly listen to control channel,
//    for msg := range controlCh {
//        cm := ControlMessageEvent{}
//        payload := []byte(msg.Payload)
//        if err := json.Unmarshal(payload, &cm); err != nil {
//            fmt.Printf("failed to parse control message: %s\n", err.Error())
//        } else if err := handleIncomingEvent(&cm); err != nil {
//            fmt.Printf("failed to handle control message: %s\n", err.Error())
//        }
//    }
//}