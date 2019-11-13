// Алгоритм:
// 1. Получение данных из источника данных
// 2. Вычисление порядкового номера события
// 3. Генерация уникального идентификатора события
// 4. Mapping входных данных в значения атрибутов
// 5. Проверка события на полноту данных
// 6. Отправка в очередь для сохранения события
// 7. Фиксация времени сохранения события в DWH
package main

import (
	"flag"
	"fmt"
	interfaces "git.fin-dev.ru/scm/dmp/dispatcher-interface.git"
	"git.fin-dev.ru/scm/dmp/dispatcher_rabbit_to_dwh.git/config"
	destination "git.fin-dev.ru/scm/dmp/dwh_client.git"
	crash "git.fin-dev.ru/scm/dmp/rabbitmq_client.git"
	source "git.fin-dev.ru/scm/dmp/rabbitmq_client.git"
	"github.com/go-errors/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

var (
	sourceClient      interfaces.Source
	destinationClient interfaces.Destination
	crashClient       interfaces.Crash
)

func main() {
	var confPath string
	flag.StringVar(&confPath, "conf", "./config.yaml", "configuration file")
	flag.Parse()
	f,err := ioutil.ReadFile(confPath)
	if err != nil {
		handleError("error on read config file", errors.Wrap(err, -1), nil, true)
	}
	err = config.Init(f)
	if err != nil {
		handleError("error on init configuration", err, nil, true)
	}
	for {
		c := config.GetConfig()
		timeStart := time.Now()
		wgDispatcher := sync.WaitGroup{}
		// канал ошибок
		errChannel := make(chan error,10)
		go func(errChannel <- chan error) {
			handleError("Ошибка в сервисах", <-errChannel, nil, false)
		}(errChannel)
		// ----------------- инициализация источника ----------------------------
		sourceClient = source.NewClient()
		// временный костыль, надеюсь (конфиги буду подтягивать с сервиса настроек)
		f,err = yaml.Marshal(c.Services.Source)
		err = sourceClient.SetConfig(f)
		if err != nil {
			handleError("error on init source client", err, nil, true)
		}
		err = sourceClient.OpenConnection()
		if err != nil {
			handleError("error on open source connection", err, nil, true)
		}
		// ----------------- завершение инициализации источника ----------------------------

		// ----------------- инициализация пункта назначения ----------------------------
		destinationClient = destination.NewClient()
		f,err = yaml.Marshal(c.Services.Destination)
		if err != nil {
			handleError("error on init destination client", errors.Wrap(err, -1), nil, true)
		}
		err = destinationClient.SetConfig(f)
		if err != nil {
			handleError("error on init destination client", err, nil, true)
		}
		err = destinationClient.OpenConnection()
		if err != nil {
			handleError("error on open destination connection", err, nil, true)
		}
		// ----------------- завершение инициализации пункта назначения ----------------------------

		// ----------------- инициализация аварийки ----------------------------
		crashClient = crash.NewClient()
		f,err = yaml.Marshal(c.Services.Crash)
		if err != nil {
			handleError("error on init destination client", errors.Wrap(err, -1), nil, true)
		}
		err = crashClient.SetConfig(f)
		if err != nil {
			handleError("error on init crash client", err, nil, true)
		}
		err = crashClient.OpenConnection()
		if err != nil {
			handleError("error on open crash connection", err, nil, true)
		}
		// ----------------- завершение инициализации аварийки ----------------------------

		wgDispatcher.Add(4)
		inputChannel := make(chan map[interface{}][]byte, 10)
		confirmChannel := make(chan interface{}, 10)
		crashChannel := make(chan []byte, 1)
		go func() {
			defer wgDispatcher.Done()
			sourceClient.ReadData(inputChannel,errChannel)
			close(inputChannel)
		}()

		go func() {
			defer wgDispatcher.Done()
			sourceClient.Confirm(confirmChannel,errChannel)
			err = sourceClient.CloseConnection()
			if err != nil {
				handleError("error on close source client", err, nil, true)
			}
		}()

		go func() {
			defer wgDispatcher.Done()
			destinationClient.WriteData(inputChannel, confirmChannel, crashChannel,errChannel)
			// закрываем каналы при завершении
			close(confirmChannel)
			close(crashChannel)
			err = destinationClient.CloseConnection()
			if err != nil {
				handleError("error on close destination client", err, nil, true)
			}
		}()

		go func() {
			defer wgDispatcher.Done()
			crashClient.SaveData(crashChannel,errChannel)
			err = crashClient.CloseConnection()
			if err != nil {
				handleError("error on close source client", err, nil, true)
			}
		}()

		wgDispatcher.Wait()
		fmt.Println("Время заняло - ", time.Since(timeStart).Seconds())
		time.Sleep(time.Duration(c.TimeOut) * time.Minute)
	}

}

func handleError(shortMessage string, err error, add *map[string]interface{}, fatal bool) {
	conf := config.GetConfig()
	if fatal {
		log.Fatalf(conf.Log.Format, "dispatcher_rabbit_to_dwh", "error", shortMessage, err.Error(), "")
	} else {
		log.Printf(conf.Log.Format, "dispatcher_rabbit_to_dwh", "warning", shortMessage, err.Error(), "")
	}
}