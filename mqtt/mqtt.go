package mqtt

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	aws "github.com/aws/aws-sdk-go/aws/credentials"
	paho "github.com/eclipse/paho.mqtt.golang"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-lib/uhppoted"

	"github.com/uhppoted/uhppoted-mqtt/acl"
	"github.com/uhppoted/uhppoted-mqtt/auth"
	"github.com/uhppoted/uhppoted-mqtt/device"
	"github.com/uhppoted/uhppoted-mqtt/log"
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
	ACL            ACL
	EventMap       string
	Protocol       string
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
	Events   struct {
		Feed string
		Live string
	}
	System string
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

type ACL struct {
	Verify map[acl.Verification]bool
}

type fdispatch struct {
	method string
	f      func(uhppoted.IUHPPOTED, []byte) (interface{}, error)
}

type dispatcher struct {
	mqttd    *MQTTD
	uhppoted *uhppoted.UHPPOTED
	devices  []uhppote.Device
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

func (mqttd *MQTTD) Run(u uhppote.IUHPPOTE, devices []uhppote.Device, authorized []string) error {
	device.SetProtocol(mqttd.Protocol)

	api := uhppoted.UHPPOTED{
		UHPPOTE:         u,
		ListenBatchSize: 32,
	}

	dev := device.Device{
		AuthorizedCards: authorized,
	}

	acl := acl.ACL{
		UHPPOTE:     u,
		Devices:     devices,
		RSA:         mqttd.Encryption.RSA,
		Credentials: mqttd.AWS.Credentials,
		Region:      mqttd.AWS.Region,
		Verify:      mqttd.ACL.Verify,
	}

	d := dispatcher{
		mqttd:    mqttd,
		uhppoted: &api,
		devices:  devices,

		table: map[string]fdispatch{
			mqttd.Topics.Requests + "/devices:get":                 fdispatch{"get-devices", dev.GetDevices},
			mqttd.Topics.Requests + "/device:get":                  fdispatch{"get-device", dev.GetDevice},
			mqttd.Topics.Requests + "/device:reset":                fdispatch{"restore-default-parameters", dev.RestoreDefaultParameters},
			mqttd.Topics.Requests + "/device/status:get":           fdispatch{"get-status", dev.GetStatus},
			mqttd.Topics.Requests + "/device/time:get":             fdispatch{"get-time", dev.GetTime},
			mqttd.Topics.Requests + "/device/time:set":             fdispatch{"set-time", dev.SetTime},
			mqttd.Topics.Requests + "/device/door/delay:get":       fdispatch{"get-door-delay", dev.GetDoorDelay},
			mqttd.Topics.Requests + "/device/door/delay:set":       fdispatch{"set-door-delay", dev.SetDoorDelay},
			mqttd.Topics.Requests + "/device/door/control:get":     fdispatch{"get-door-control", dev.GetDoorControl},
			mqttd.Topics.Requests + "/device/door/control:set":     fdispatch{"set-door-control", dev.SetDoorControl},
			mqttd.Topics.Requests + "/device/door/passcodes:set":   fdispatch{"set-passcodes", dev.SetDoorPasscodes},
			mqttd.Topics.Requests + "/device/door/interlock:set":   fdispatch{"set-interlock", dev.SetInterlock},
			mqttd.Topics.Requests + "/device/door/keypads:set":     fdispatch{"set-keypads", dev.SetKeypads},
			mqttd.Topics.Requests + "/device/door/lock:open":       fdispatch{"open-door", dev.OpenDoor},
			mqttd.Topics.Requests + "/device/special-events:set":   fdispatch{"record-special-events", dev.RecordSpecialEvents},
			mqttd.Topics.Requests + "/device/cards:get":            fdispatch{"get-cards", dev.GetCards},
			mqttd.Topics.Requests + "/device/cards:delete":         fdispatch{"delete-cards", dev.DeleteCards},
			mqttd.Topics.Requests + "/device/card:get":             fdispatch{"get-card", dev.GetCard},
			mqttd.Topics.Requests + "/device/card:put":             fdispatch{"put-card", dev.PutCard},
			mqttd.Topics.Requests + "/device/card:delete":          fdispatch{"delete-card", dev.DeleteCard},
			mqttd.Topics.Requests + "/device/time-profile:get":     fdispatch{"get-time-profile", dev.GetTimeProfile},
			mqttd.Topics.Requests + "/device/time-profile:set":     fdispatch{"set-time-profile", dev.PutTimeProfile},
			mqttd.Topics.Requests + "/device/time-profiles:get":    fdispatch{"get-time-profiles", dev.GetTimeProfiles},
			mqttd.Topics.Requests + "/device/time-profiles:set":    fdispatch{"get-time-profiles", dev.PutTimeProfiles},
			mqttd.Topics.Requests + "/device/time-profiles:delete": fdispatch{"clear-time-profiles", dev.ClearTimeProfiles},
			mqttd.Topics.Requests + "/device/tasklist:set":         fdispatch{"set-task-list", dev.PutTaskList},
			mqttd.Topics.Requests + "/device/events:get":           fdispatch{"get-events", dev.GetEvents},
			mqttd.Topics.Requests + "/device/event:get":            fdispatch{"get-event", dev.GetEvent},

			mqttd.Topics.Requests + "/acl/card:show":    fdispatch{"acl:show", acl.Show},
			mqttd.Topics.Requests + "/acl/card:grant":   fdispatch{"acl:grant", acl.Grant},
			mqttd.Topics.Requests + "/acl/card:revoke":  fdispatch{"acl:revoke", acl.Revoke},
			mqttd.Topics.Requests + "/acl/acl:upload":   fdispatch{"acl:upload", acl.Upload},
			mqttd.Topics.Requests + "/acl/acl:download": fdispatch{"acl:download", acl.Download},
			mqttd.Topics.Requests + "/acl/acl:compare":  fdispatch{"acl:compare", acl.Compare},
		},
	}

	if client, err := mqttd.subscribeAndServe(&d); err != nil {
		return fmt.Errorf("ERROR: Error connecting to '%s': %v", mqttd.Connection.Broker, err)
	} else {
		mqttd.client = client
	}

	if err := mqttd.listen(&api, u); err != nil {
		return fmt.Errorf("ERROR: Failed to bind to MQTT listen port '%d': %v", 12345, err)
	}

	return nil
}

func (m *MQTTD) Close() {
	if m.interrupt != nil {
		close(m.interrupt)
	}

	if m.client != nil {
		infof("closing connection to %s", m.Connection.Broker)
		m.client.Disconnect(250)
		infof("closed connection to %s", m.Connection.Broker)
	}

	m.client = nil
}

func (m *MQTTD) subscribeAndServe(d *dispatcher) (paho.Client, error) {
	var handler paho.MessageHandler = func(client paho.Client, msg paho.Message) {
		d.dispatch(client, msg)
	}

	var connected paho.OnConnectHandler = func(client paho.Client) {
		options := client.OptionsReader()
		servers := options.Servers()
		for _, url := range servers {
			infof("connected to %s", url)
		}

		token := m.client.Subscribe(m.Topics.Requests+"/#", 0, handler)
		if err := token.Error(); err != nil {
			errorf("unable to subscribe to %s (%v)", m.Topics.Requests, err)
		} else {
			infof("subscribed to %s", m.Topics.Requests)
		}
	}

	var disconnected paho.ConnectionLostHandler = func(client paho.Client, err error) {
		errorf("connection to MQTT broker lost (%v)", err)

		stats.onDisconnected()

		go func() {
			time.Sleep(10 * time.Second)
			infof("retrying connection to MQTT broker %v", m.Connection.Broker)
			token := client.Connect()
			if err := token.Error(); err != nil {
				errorf("failed to reconnect to MQTT broker (%v)", err)
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

func (m *MQTTD) listen(api *uhppoted.UHPPOTED, u uhppote.IUHPPOTE) error {
	list := u.ListenAddrList()

	if list == nil {
		infof("no event listener")
	} else if len(list) == 1 {
		infof("listening on %v", list[0])
		infof("publishing events to %s", m.Topics.Events)
	} else {
		infof("listening on %v", list)
		infof("publishing events to %s", m.Topics.Events)
	}

	handler := func(e any, queue string) bool {
		event := struct {
			Event any `json:"event"`
		}{
			Event: e,
		}

		topic := m.Topics.Events.Feed
		switch queue {
		case "live":
			topic = m.Topics.Events.Live

		case "feed":
			topic = m.Topics.Events.Feed
		}

		if err := m.send(&m.Encryption.EventsKeyID, topic, nil, event, msgEvent, true); err != nil {
			warnf("%v", err)
			return false
		}

		return true
	}

	m.interrupt = make(chan os.Signal)

	device.Listen(*api, m.EventMap, handler, m.interrupt)

	return nil
}

func (d *dispatcher) dispatch(client paho.Client, msg paho.Message) {
	if fn, ok := d.table[msg.Topic()]; ok {
		msg.Ack()

		debugf("%v", string(msg.Payload()))

		go func() {
			rq, err := d.mqttd.unwrap(msg.Payload())
			if err != nil {
				warnf("%v", err)
				return
			}

			if err := d.mqttd.authorise(rq.ClientID, msg.Topic()); err != nil {
				warnf("%-20v error authorising request (%v)", fn.method, err)
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
				warnf("%-20v %v", fn.method, err)
				if response != nil {
					reply := struct {
						Error interface{} `json:"error"`
					}{
						Error: response,
					}

					if err := d.mqttd.send(rq.ClientID, replyTo, &meta, reply, msgError, false); err != nil {
						warnf("%-20v %v", fn.method, err)
					}
				}
			} else if response != nil {
				reply := struct {
					Response interface{} `json:"response"`
				}{
					Response: response,
				}

				if err := d.mqttd.send(rq.ClientID, replyTo, &meta, reply, msgReply, false); err != nil {
					warnf("%-20v %v", fn.method, err)
				}
			}
		}()
	}
}

func (m *MQTTD) authorise(clientID *string, topic string) error {
	if m.Permissions.Enabled {
		if clientID == nil {
			return errors.New("request without client-id")
		}

		match := regexp.MustCompile(`.*?/(\w+):(\w+)$`).FindStringSubmatch(topic)
		if len(match) != 3 {
			return fmt.Errorf("invalid resource:action (%s)", topic)
		}

		return m.Permissions.Validate(*clientID, match[1], match[2])
	}

	return nil
}

// TODO: add callback for published/failed
func (mqttd *MQTTD) send(destID *string, topic string, meta *metainfo, message interface{}, msgtype msgType, critical bool) error {
	if mqttd.client == nil || !mqttd.client.IsConnected() {
		return errors.New("no connection to MQTT broker")
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
}

func isBase64(request []byte) bool {
	return regex.base64.Match(request)
}

func debugf(format string, args ...any) {
	log.Debugf("mqttd", format, args...)
}

func infof(format string, args ...any) {
	log.Infof("mqttd", format, args...)
}

func warnf(format string, args ...any) {
	log.Warnf("mqttd", format, args...)
}

func errorf(format string, args ...any) {
	log.Errorf("mqttd", format, args...)
}

func fatalf(format string, args ...any) {
	log.Fatalf("mqttd", format, args...)
}
