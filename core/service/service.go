// Core service infrastructure for servicing starting/stopping/SIGTERM, and heartbeating etc
package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/cyanly/gotrade/core/logger"
	proto "github.com/cyanly/gotrade/proto/service"
	messagebus "github.com/nats-io/nats"
)

type Service struct {
	Config Config
	Status proto.Heartbeat_Status

	shutdownChannel chan bool
	messageBus      *messagebus.Conn
	lastHBMsg       *proto.Heartbeat
	publishAddress  string
}

type Service struct {
	Config Config
	Status proto.Heartbeat_Status

	shutdownChannel chan bool
	messageBus      *messagebus.Conn
	lastHBMsg       *proto.Heartbeat
	publishAddress  string
}

func NewService(c Config) *Service {

	// Hardware Info
	uuid = fmt.Sprint(hostname, ":", pid)
	log.Infof("Service [%v] starting @ %v", c.ServiceName, uuid)

	// Service handle
	svc := &Service{
		Config:          c,
		Status:          proto.STARTING,
		shutdownChannel: make(chan bool),
	}

	// Messaging bus
	messageBus, err := messagebus.Connect(svc.Config.MessageBusURL)
	svc.messageBus = messageBus
	if err != nil {
		log.WithField("MessageBusURL", svc.Config.MessageBusURL).Fatal("error: Cannot connect to message bus")
	}

	//Heartbeating
	currDateTime := time.Now().UTC().Format(time.RFC3339)
	hbMsg := &proto.Heartbeat{
		Name:             svc.Config.ServiceName,
		Id:               uuid,
		Status:           proto.STARTING,
		Machine:          hostname,
		CreationDatetime: currDateTime,
		CurrentDatetime:  currDateTime,
	}
	svc.lastHBMsg = hbMsg
	hbTicker := time.NewTicker(time.Second * time.Duration(svc.Config.HeartbeatFreq))
	go func(shutdownChannel chan bool) {
		publish_address := "service.Heartbeat." + svc.Config.ServiceName

		for range hbTicker.C {
			hbMsg.CurrentDatetime = time.Now().UTC().Format(time.RFC3339)
			hbMsg.Status = svc.Status

			if data, _ := hbMsg.Marshal(); data != nil {
				messageBus.Publish(publish_address, data)
			}

			select {
			case <-shutdownChannel:
				hbTicker.Stop()

			//Publish Stop heartbeat
				if svc.Status != proto.ERROR {
					svc.Status = proto.STOPPED
				}
				hbMsg.CurrentDatetime = time.Now().UTC().Format(time.RFC3339)
				hbMsg.Status = svc.Status
				if data, _ := hbMsg.Marshal(); data != nil {
					messageBus.Publish(publish_address, data)
				}

				messageBus.Close()

				log.Info("Server Terminated")
				return
			}
		}
	}(svc.shutdownChannel)

	return svc
}

func (self *Service) Start() chan bool {
	//SIGINT or SIGTERM is caught
	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	shutdownCallerChannel := make(chan bool)
	go func() {
		<-quitChannel
		self.shutdownChannel <- true
		shutdownCallerChannel <- true
	}()

	self.Status = proto.RUNNING

	// Immediately publish heartbeat
	self.lastHBMsg.CurrentDatetime = time.Now().UTC().Format(time.RFC3339)
	self.lastHBMsg.Status = self.Status
	if data, _ := self.lastHBMsg.Marshal(); data != nil {
		self.messageBus.Publish("service.Heartbeat."+self.Config.ServiceName, data)
	}

	log.Infof("Service [%v] Started", self.Config.ServiceName)
	return shutdownCallerChannel
}

func (self *Service) Stop() {
	self.shutdownChannel <- true
}