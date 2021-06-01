package mqtt

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	aws "github.com/aws/aws-sdk-go/aws/credentials"
	paho "github.com/eclipse/paho.mqtt.golang"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/acl"
	"github.com/uhppoted/uhppoted-mqtt/auth"
	"github.com/uhppoted/uhppoted-mqtt/device"
)

type MQTTD struct {
	ServerID       string
	Connection     Connection
	TLS            *tls.Config
	Topics         Topics
	Alerts         Alerts
	HMAC           auth.HMAC
	Encryption     Encryption
	Authentication string
	Permissions    auth.Permissions
	AWS            AWS
	EventMap       string
	Debug          bool

	client    paho.Client
	interrupt chan os.Signal
}

type Connection struct {
	Broker   string
	ClientID string
	UserName string
	Password string
}

type Topics struct {
	Requests string
	Replies  string
	Events   string
	System   string
}

type Alerts struct {
	QOS      byte
	Retained bool
}

type Encryption struct {
	SignOutgoing    bool
	EncryptOutgoing bool
	EventsKeyID     string
	SystemKeyID     string
	HOTP            *auth.HOTP
	RSA             *auth.RSA
	Nonce           auth.Nonce
}

type AWS struct {
	Credentials *aws.Credentials
	Region      string
}

type fdispatch struct {
	method string
	f      func(*uhppoted.UHPPOTED, []byte) (interface{}, error)
}

type dispatcher struct {
	mqttd    *MQTTD
	uhppoted *uhppoted.UHPPOTED
	devices  []uhppote.Device
	log      *log.Logger
	table    map[string]fdispatch
}

type request struct {
	ClientID  *string
	RequestID *string
	ReplyTo   *string
	Request   []byte
}

type metainfo struct {
	RequestID *string `json:"request-id,omitempty"`
	ClientID  *string `json:"client-id,omitempty"`
	ServerID  string  `json:"server-id,omitempty"`
	Method    string  `json:"method,omitempty"`
	Nonce     fnonce  `json:"nonce,omitempty"`
}

type fnonce func() uint64

func (f fnonce) MarshalJSON() ([]byte, error) {
	return json.Marshal(f())
}

var regex = struct {
	clean  *regexp.Regexp
	base64 *regexp.Regexp
}{
	clean:  regexp.MustCompile(`\s+`),
	base64: regexp.MustCompile(`^"[A-Za-z0-9+/]*[=]{0,2}"$`),
}

func (mqttd *MQTTD) Run(u uhppote.IUHPPOTE, devices []uhppote.Device, log *log.Logger) error {
	paho.CRITICAL = log
	paho.ERROR = log
	paho.WARN = log

	if mqttd.Debug {
		paho.DEBUG = log
	}

	api := uhppoted.UHPPOTED{
		UHPPOTE:         u,
		ListenBatchSize: 32,
		Log:             log,
	}

	dev := device.Device{
		Log: log,
	}

	acl := acl.ACL{
		Devices:     devices,
		RSA:         mqttd.Encryption.RSA,
		Credentials: mqttd.AWS.Credentials,
		Region:      mqttd.AWS.Region,
		Log:         log,
		NoVerify:    false,
	}

	d := dispatcher{
		mqttd:    mqttd,
		uhppoted: &api,
		devices:  devices,
		log:      log,
		table: map[string]fdispatch{
			mqttd.Topics.Requests + "/devices:get":               fdispatch{"get-devices", dev.GetDevices},
			mqttd.Topics.Requests + "/device:get":                fdispatch{"get-device", dev.GetDevice},
			mqttd.Topics.Requests + "/device/status:get":         fdispatch{"get-status", dev.GetStatus},
			mqttd.Topics.Requests + "/device/time:get":           fdispatch{"get-time", dev.GetTime},
			mqttd.Topics.Requests + "/device/time:set":           fdispatch{"set-time", dev.SetTime},
			mqttd.Topics.Requests + "/device/door/delay:get":     fdispatch{"get-door-delay", dev.GetDoorDelay},
			mqttd.Topics.Requests + "/device/door/delay:set":     fdispatch{"set-door-delay", dev.SetDoorDelay},
			mqttd.Topics.Requests + "/device/door/control:get":   fdispatch{"get-door-control", dev.GetDoorControl},
			mqttd.Topics.Requests + "/device/door/control:set":   fdispatch{"set-door-control", dev.SetDoorControl},
			mqttd.Topics.Requests + "/device/cards:get":          fdispatch{"get-cards", dev.GetCards},
			mqttd.Topics.Requests + "/device/cards:delete":       fdispatch{"delete-cards", dev.DeleteCards},
			mqttd.Topics.Requests + "/device/card:get":           fdispatch{"get-card", dev.GetCard},
			mqttd.Topics.Requests + "/device/card:put":           fdispatch{"put-card", dev.PutCard},
			mqttd.Topics.Requests + "/device/card:delete":        fdispatch{"delete-card", dev.DeleteCard},
			mqttd.Topics.Requests + "/device/time-profile:get":   fdispatch{"get-time-profile", dev.GetTimeProfile},
			mqttd.Topics.Requests + "/device/events:get":         fdispatch{"get-events", dev.GetEvents},
			mqttd.Topics.Requests + "/device/event:get":          fdispatch{"get-event", dev.GetEvent},
			mqttd.Topics.Requests + "/device/special-events:set": fdispatch{"record-special-events", dev.RecordSpecialEvents},

			mqttd.Topics.Requests + "/acl/card:show":    fdispatch{"acl:show", acl.Show},
			mqttd.Topics.Requests + "/acl/card:grant":   fdispatch{"acl:grant", acl.Grant},
			mqttd.Topics.Requests + "/acl/card:revoke":  fdispatch{"acl:revoke", acl.Revoke},
			mqttd.Topics.Requests + "/acl/acl:upload":   fdispatch{"acl:upload", acl.Upload},
			mqttd.Topics.Requests + "/acl/acl:download": fdispatch{"acl:download", acl.Download},
			mqttd.Topics.Requests + "/acl/acl:compare":  fdispatch{"acl:compare", acl.Compare},
		},
	}

	if client, err := mqttd.subscribeAndServe(&d, log); err != nil {
		return fmt.Errorf("ERROR: Error connecting to '%s': %v", mqttd.Connection.Broker, err)
	} else {
		mqttd.client = client
	}

	if err := mqttd.listen(&api, u, log); err != nil {
		return fmt.Errorf("ERROR: Failed to bind to listen port '%d': %v", 12345, err)
	}

	return nil
}

func (m *MQTTD) Close(log *log.Logger) {
	if m.interrupt != nil {
		close(m.interrupt)
	}

	if m.client != nil {
		log.Printf("INFO  closing connection to %s", m.Connection.Broker)
		m.client.Disconnect(250)
		log.Printf("INFO  closed connection to %s", m.Connection.Broker)
	}

	m.client = nil
}

func (m *MQTTD) subscribeAndServe(d *dispatcher, log *log.Logger) (paho.Client, error) {
	var handler paho.MessageHandler = func(client paho.Client, msg paho.Message) {
		d.dispatch(client, msg)
	}

	var connected paho.OnConnectHandler = func(client paho.Client) {
		options := client.OptionsReader()
		servers := options.Servers()
		for _, url := range servers {
			log.Printf("%-5s %-12s %v", "INFO", "mqttd", fmt.Sprintf("Connected to %s", url))
		}

		token := m.client.Subscribe(m.Topics.Requests+"/#", 0, handler)
		if err := token.Error(); err != nil {
			log.Printf("ERROR unable to subscribe to %s (%v)", m.Topics.Requests, err)
			return
		}

		log.Printf("%-5s %-12s %v", "INFO", "mqttd", fmt.Sprintf("Subscribed to %s", m.Topics.Requests))
	}

	var disconnected paho.ConnectionLostHandler = func(client paho.Client, err error) {
		log.Printf("ERROR connection to MQTT broker lost (%v)", err)
		go func() {
			time.Sleep(10 * time.Second)
			log.Printf("INFO  retrying connection to MQTT broker %v", m.Connection.Broker)
			token := client.Connect()
			if err := token.Error(); err != nil {
				log.Printf("ERROR failed to reconnect to MQTT broker (%v)", err)
			}
		}()
	}

	// NOTE: Paho auto-reconnect causes a retry storm if two MQTT clients are using the same client ID.
	//       'Theoretically' (à la Terminator Genesys) the lockfile should prevent this but careful
	//       misconfiguration is always a possibility.
	options := paho.
		NewClientOptions().
		AddBroker(m.Connection.Broker).
		SetClientID(m.Connection.ClientID).
		SetTLSConfig(m.TLS).
		SetCleanSession(false).
		SetAutoReconnect(false).
		SetConnectRetry(true).
		SetConnectRetryInterval(30 * time.Second).
		SetOnConnectHandler(connected).
		SetConnectionLostHandler(disconnected)

	if m.Connection.UserName != "" {
		options.SetUsername(m.Connection.UserName)
		if m.Connection.Password != "" {
			options.SetPassword(m.Connection.Password)
		}
	}

	client := paho.NewClient(options)
	token := client.Connect()
	if err := token.Error(); err != nil {
		return nil, err
	}

	return client, nil
}

func (m *MQTTD) listen(api *uhppoted.UHPPOTED, u uhppote.IUHPPOTE, log *log.Logger) error {
	log.Printf("%-5s %-12s %v", "INFO", "mqttd", fmt.Sprintf("Listening on %v", u.ListenAddr()))
	log.Printf("%-5s %-12s %v", "INFO", "mqttd", fmt.Sprintf("Publishing events to %s", m.Topics.Events))

	last := uhppoted.NewEventMap(m.EventMap)
	if err := last.Load(log); err != nil {
		log.Printf("WARN  Error loading event map [%v]", err)
	}

	handler := func(event uhppoted.EventMessage) bool {
		if err := m.send(&m.Encryption.EventsKeyID, m.Topics.Events, nil, event, msgEvent, true); err != nil {
			log.Printf("WARN  %-12s %v", "listen", err)
			return false
		}

		return true
	}

	m.interrupt = make(chan os.Signal)

	go func() {
		api.Listen(handler, last, m.interrupt)
	}()

	return nil
}

func (d *dispatcher) dispatch(client paho.Client, msg paho.Message) {
	ctx := context.WithValue(context.Background(), "client", client)
	ctx = context.WithValue(ctx, "log", d.log)

	if fn, ok := d.table[msg.Topic()]; ok {
		msg.Ack()

		d.log.Printf("DEBUG %-20s %s", "dispatch", string(msg.Payload()))

		go func() {
			rq, err := d.mqttd.unwrap(msg.Payload())
			if err != nil {
				d.log.Printf("WARN  %-20s %v", "dispatch", err)
				return
			}

			if err := d.mqttd.authorise(rq.ClientID, msg.Topic()); err != nil {
				d.log.Printf("WARN  %-20s %v", fn.method, fmt.Errorf("Error authorising request (%v)", err))
				return
			}

			replyTo := d.mqttd.Topics.Replies

			if rq.ClientID != nil {
				replyTo = d.mqttd.Topics.Replies + "/" + *rq.ClientID
			}

			if rq.ReplyTo != nil {
				replyTo = *rq.ReplyTo
			}

			meta := metainfo{
				RequestID: rq.RequestID,
				ClientID:  rq.ClientID,
				ServerID:  d.mqttd.ServerID,
				Method:    fn.method,
				Nonce:     func() uint64 { return d.mqttd.Encryption.Nonce.Next() },
			}

			response, err := fn.f(d.uhppoted, rq.Request)

			if err != nil {
				d.log.Printf("WARN  %-12s %v", fn.method, err)
				if response != nil {
					reply := struct {
						Error interface{} `json:"error"`
					}{
						Error: response,
					}

					if err := d.mqttd.send(rq.ClientID, replyTo, &meta, reply, msgError, false); err != nil {
						d.log.Printf("WARN  %-20s %v", fn.method, err)
					}
				}
			} else if response != nil {
				reply := struct {
					Response interface{} `json:"response"`
				}{
					Response: response,
				}

				if err := d.mqttd.send(rq.ClientID, replyTo, &meta, reply, msgReply, false); err != nil {
					d.log.Printf("WARN  %-20s %v", fn.method, err)
				}
			}
		}()
	}
}

func (m *MQTTD) authorise(clientID *string, topic string) error {
	if m.Permissions.Enabled {
		if clientID == nil {
			return errors.New("Request without client-id")
		}

		match := regexp.MustCompile(`.*?/(\w+):(\w+)$`).FindStringSubmatch(topic)
		if len(match) != 3 {
			return fmt.Errorf("Invalid resource:action (%s)", topic)
		}

		return m.Permissions.Validate(*clientID, match[1], match[2])
	}

	return nil
}

// TODO: add callback for published/failed
// func (mqttd *MQTTD) send(destID *string, topic string, message interface{}, msgtype msgType, critical bool) error {
// 	if mqttd.client == nil || !mqttd.client.IsConnected() {
// 		return errors.New("No connection to MQTT broker")
// 	}
//
// 	m, err := mqttd.wrap(msgtype, message, destID)
// 	if err != nil {
// 		return err
// 	} else if m == nil {
// 		return errors.New("'wrap' failed to return a publishable message")
// 	}
//
// 	qos := byte(0)
// 	retained := false
// 	if critical {
// 		qos = mqttd.Alerts.QOS
// 		retained = mqttd.Alerts.Retained
// 	}
//
// 	token := mqttd.client.Publish(topic, qos, retained, string(m))
// 	if token.Error() != nil {
// 		return token.Error()
// 	}
//
// 	return nil
// }
//
// TODO: add callback for published/failed
func (mqttd *MQTTD) send(destID *string, topic string, meta *metainfo, message interface{}, msgtype msgType, critical bool) error {
	if mqttd.client == nil || !mqttd.client.IsConnected() {
		return errors.New("No connection to MQTT broker")
	}

	content, err := compose(meta, message)
	if err != nil {
		return err
	}

	m, err := mqttd.wrap(msgtype, content, destID)
	if err != nil {
		return err
	} else if m == nil {
		return errors.New("'wrap' failed to return a publishable message")
	}

	qos := byte(0)
	retained := false
	if critical {
		qos = mqttd.Alerts.QOS
		retained = mqttd.Alerts.Retained
	}

	token := mqttd.client.Publish(topic, qos, retained, string(m))
	if token.Error() != nil {
		return token.Error()
	}

	return nil
}

func compose(meta *metainfo, content interface{}) (interface{}, error) {
	reply := make(map[string]interface{})

	s, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(s, &reply)
	if err != nil {
		return nil, err
	}

	s, err = json.Marshal(content)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(s, &reply)
	if err != nil {
		return nil, err
	}

	return reply, nil

	//	return struct {
	//		*metainfo `json:",omitempty"`
	//		Content   interface{} `json:"body"`
	//	}{
	//		metainfo: meta,
	//		Content:  content,
	//	}, nil
}

func isBase64(request []byte) bool {
	return regex.base64.Match(request)
}

func clean(s string) string {
	return regex.clean.ReplaceAllString(s, " ")
}
