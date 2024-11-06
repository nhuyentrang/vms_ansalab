package wsnode

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"vms/wsnode/jwt"

	"net/http"
	_ "net/http/pprof"

	"github.com/centrifugal/centrifuge"
)

var (
	// This should be read from database or .env file.
	hMACSecretKey string = "prlchASWacAWRewrenoCLP9ivarophiGlgop"
)

var (
	m_node              *centrifuge.Node
	m_broker            *centrifuge.MemoryBroker
	WebsocketHandler    *centrifuge.WebsocketHandler
	m_clientIDSeparator = "."
)

// map for store onSubscribe handle function for each prefix of topic name
var onSubscribeHandleFunctions = map[string]func(string, string){}

// map for store onDisconnect handle function for each prefix of topic name
var onDisconnectHandleFunctions = map[string]func(string, string){}

// map for store onPublish handle function for each prefix of topic name
var onPublishHandleFunctions = map[string]func(string, string, []byte){}

// RegisterProcessingFunction allows registration of a processing function for a specific topic prefix.
func RegisterProcessingFunction(prefix string,
	fnOnSubscribe func(string, string),
	fnOnDisconnect func(string, string),
	fnOnPublish func(string, string, []byte)) {
	onSubscribeHandleFunctions[prefix] = fnOnSubscribe
	onDisconnectHandleFunctions[prefix] = fnOnDisconnect
	onPublishHandleFunctions[prefix] = fnOnPublish
}

// Check whether topic is allowed for subscribing.
func topicSubscribeAllowed(topic string, clientID string) bool {
	// ClientID format: type_uuid
	// Client is allow to subscribe to topic name same as it's clientID
	// Client is allow to subscribe to topic name contain it's type, exp: notification:${type}
	// type: aieagent, ipc, smartnvr, cloudcam, web, mobile

	// Get topic prefix and topic ID
	parts := strings.Split(topic, m_clientIDSeparator)
	if len(parts) < 2 {
		log.Printf("topicSubscribeAllowed: %s, %s, invalid topic name", topic, clientID)
		return false
	}
	topicPrefix := parts[0]
	topicID := parts[1]

	// Get client type and client UUID
	parts = strings.Split(clientID, m_clientIDSeparator)
	if len(parts) < 1 {
		log.Printf("topicSubscribeAllowed: %s, %s, invalid clientID name", topic, clientID)
		return false
	}
	clientType := parts[0]

	// Check if the topic is allowed based on clientID
	if strings.Contains(topic, clientID) ||
		(topicPrefix == "notification" && strings.Contains(topicID, clientType)) {
		return true
	}
	return false
}

func handleLog(e centrifuge.LogEntry) {
	log.Printf("%s: %v", e.Message, e.Fields)
}

func WaitExitSignal() {
	sigCh := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		_ = m_node.Shutdown(context.Background())
		done <- true
	}()
	<-done

	log.Println("bye!")
}

func Run() {
	m_node, _ = centrifuge.New(centrifuge.Config{
		LogLevel:   centrifuge.LogLevelError,
		LogHandler: handleLog,
		// Better to keep default in production. Here we are just speeding up things a bit.
		//ClientExpiredCloseDelay: 5 * time.Second,
	})

	m_broker, _ = centrifuge.NewMemoryBroker(m_node, centrifuge.MemoryBrokerConfig{})
	m_node.SetBroker(m_broker)

	// Todo: should fix to use golang-jwt v5 from "github.com/golang-jwt/jwt/v5"
	tokenVerifier := jwt.NewTokenVerifier(jwt.TokenVerifierConfig{
		HMACSecretKey: hMACSecretKey,
	})

	m_node.OnConnecting(func(ctx context.Context, e centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
		log.Printf("client connecting with token: %s", e.Token)
		token, err := tokenVerifier.VerifyConnectToken(e.Token)
		if err != nil {
			if err == jwt.ErrTokenExpired {
				return centrifuge.ConnectReply{}, centrifuge.ErrorTokenExpired
			}
			log.Printf("error verifying token: %s", err)
			return centrifuge.ConnectReply{}, centrifuge.DisconnectInvalidToken
		}
		subs := make(map[string]centrifuge.SubscribeOptions, len(token.Channels))
		for _, ch := range token.Channels {
			subs[ch] = centrifuge.SubscribeOptions{}
		}
		return centrifuge.ConnectReply{
			Credentials: &centrifuge.Credentials{
				UserID:   token.UserID,
				ExpireAt: token.ExpireAt,
			},
			Subscriptions: subs,
		}, nil
	})

	m_node.OnConnect(func(client *centrifuge.Client) {
		//transport := client.Transport()
		//log.Printf("user %s connected via %s with protocol: %s", client.UserID(), transport.Name(), transport.Protocol())

		client.OnRefresh(func(e centrifuge.RefreshEvent, cb centrifuge.RefreshCallback) {
			log.Printf("user %s connection is going to expire, refreshing", client.UserID())
			cb(centrifuge.RefreshReply{ExpireAt: time.Now().Unix() + 60}, nil)
		})

		client.OnSubscribe(func(e centrifuge.SubscribeEvent, cb centrifuge.SubscribeCallback) {
			//log.Printf("user %s subscribes on %s", client.UserID(), e.Channel)
			if !topicSubscribeAllowed(e.Channel, client.UserID()) {
				cb(centrifuge.SubscribeReply{}, centrifuge.ErrorPermissionDenied)
				return
			}
			cb(centrifuge.SubscribeReply{}, nil)

			// Depend on prefix of clientID, we can apply different processing function
			// ClientID name format contain prefix and id, separated by dot, exp: aieagent.uuid
			parts := strings.Split(client.UserID(), m_clientIDSeparator)
			if len(parts) < 2 {
				log.Printf("invalid clientID format: %s", client.UserID())
				return
			}
			prefix := parts[0]
			if handleFunc, exists := onSubscribeHandleFunctions[prefix]; exists {
				//log.Printf("call process func for topic: %s", prefix)
				handleFunc(client.UserID(), e.Channel)
			} else {
				log.Printf("no subscribe handle function registered for topic prefix: %s", prefix)
			}
		})

		client.OnPublish(func(e centrifuge.PublishEvent, cb centrifuge.PublishCallback) {
			//log.Printf("user %s publishes into topic %s: %s", client.UserID(), e.Channel, string(e.Data))
			if !client.IsSubscribed(e.Channel) {
				cb(centrifuge.PublishReply{}, centrifuge.ErrorPermissionDenied)
				return
			}
			cb(centrifuge.PublishReply{}, nil)

			// Depend on prefix of clientID, we can apply different processing function
			// ClientID name format contain prefix and id, separated by dot, exp: aieagent.uuid
			parts := strings.Split(client.UserID(), m_clientIDSeparator)
			if len(parts) < 2 {
				log.Printf("invalid clientID format: %s", client.UserID())
				return
			}
			prefix := parts[0]
			if handleFunc, exists := onPublishHandleFunctions[prefix]; exists {
				//log.Printf("call process func for topic: %s", prefix)
				handleFunc(client.UserID(), e.Channel, e.Data)
			} else {
				log.Printf("no publish handle function registered for topic prefix: %s", prefix)
			}
		})

		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			//log.Printf("user %s disconnected, disconnect: %s", client.UserID(), e.Disconnect)

			// Depend on prefix of clientID, we can apply different processing function
			// ClientID name format contain prefix and id, separated by dot, exp: aieagent.uuid
			parts := strings.Split(client.UserID(), m_clientIDSeparator)
			if len(parts) < 2 {
				log.Printf("invalid clientID format: %s", client.UserID())
				return
			}
			prefix := parts[0]
			if handleFunc, exists := onDisconnectHandleFunctions[prefix]; exists {
				//log.Printf("call process func for topic: %s", prefix)
				handleFunc(client.UserID(), e.Disconnect.String())
			} else {
				log.Printf("no disconnect handle function registered for topic prefix: %s", prefix)
			}
		})
	})

	if err := m_node.Run(); err != nil {
		log.Fatal(err)
	}

	WebsocketHandler = centrifuge.NewWebsocketHandler(m_node, centrifuge.WebsocketConfig{
		ReadBufferSize:     1024,
		UseWriteBufferPool: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	})
}

func Publish(topic string, data []byte) error {
	_, err := m_node.Publish(topic, data)
	return err
}
