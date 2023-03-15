package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"

	"maxws/pkg"
)

type EventChanEnvelope struct {
	Key     string
	Payload []byte

	ProtoType int
}

var DataQueueChan chan *EventChanEnvelope

var addr = flag.String("addr", "0.0.0.0:8380", "http service address")

var upgrader = websocket.Upgrader{} // use default options

var MaxDB *pkg.BadgerDBWrap

// var MaxDB *pkg.MgDB
// var MaxDB *pkg.InfluxStore

var LOG *logrus.Logger

func failOnErr(err error, msg string) {
	if err != nil {
		if msg != "" {
			msg = msg + ": "
		}
		LOG.Fatal(err, msg, zap.Error(err))
	}
}

func InitDB() *pkg.BadgerDBWrap {
	BagerDB, err := pkg.InitBadgerDBWrapInstance("./badgerdb")
	if err != nil {
		failOnErr(err, "Could not connect to badger DB")
	}

	return BagerDB
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer c.Close()

	for {
		msgType, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		eventChan := &EventChanEnvelope{Key: string(message), Payload: message}
		DataQueueChan <- eventChan

		// log.Printf("recv: %s", message)
		err = c.WriteMessage(msgType, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func QueueChannelHandler() {

	DataList := map[string][]byte{}

	BatchSize := 1000

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case envelop := <-DataQueueChan:
			// fmt.Println(">>>", envelop.Key)
			// MaxDB.SaveMessage(envelop.Key, envelop.Payload)
			// err := MaxDB.BatchSet(envelop.Key, envelop.Payload)
			// if err != nil {
			// 	LOG.Fatal(err)
			// }

			if len(DataList) >= BatchSize {
				err := MaxDB.BatchSetList(&DataList)
				if err != nil {
					LOG.Fatal(err)
				} else {
					DataList = map[string][]byte{}
				}
			} else {
				DataList[envelop.Key] = envelop.Payload
			}

			// err := MaxDB.SaveMessage(envelop.Key, envelop.Payload)

			// if err != nil {
			// 	LOG.Fatal(err)
			// 	// fmt.Println(">>>", err.Error())
			// }

			// err := MaxDB.UpdateByte(envelop.Key, envelop.Payload)
			// if err != nil {
			// 	LOG.Fatal(err)
			// }
		case <-ticker.C:
			//
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	DataQueueChan = make(chan *EventChanEnvelope, 10000)

	var logConf = &lumberjack.Logger{
		// output path
		Filename: "./MaxServer.log",
		// max size / MB
		MaxSize: 500, // megabytes
		// max log files count
		MaxBackups: 5,
		// max expired days
		MaxAge: 28, //days
		// compress, use gzip
		Compress: true, // disabled by default
	}

	LOG = logrus.New()

	LOG.SetOutput(logConf)

	MaxDB = InitDB()
	defer MaxDB.Close()

	// MaxDB = pkg.NewMgDB()
	// MaxDB = pkg.NewInfluxDB()
	// defer MaxDB.Close()

	go QueueChannelHandler()

	http.HandleFunc("/messages", WebsocketHandler)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
