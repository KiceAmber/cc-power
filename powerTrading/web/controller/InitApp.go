package controller

import "powerTrading/service"

var ProducerApp Application
var ConsumerApp Application

func InitApp(consumer *service.ServiceSetup, producer *service.ServiceSetup) {
	ConsumerApp = Application{
		Setup: consumer,
	}
	ProducerApp = Application{
		Setup: producer,
	}
}
