package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type APIServer struct{}

type ErrorResponse struct {
	Error string
}

type SuccessResponse struct {
	Message string
}

func (apiServer APIServer) Start() {
	var wg sync.WaitGroup
	log.Println("Starting API server...")
	http.HandleFunc("/", apiServer.Endpoint)
	wg.Add(1)
	go http.ListenAndServe(":8080", nil)
	log.Println("API Server started...")
	wg.Wait()
}

var mu sync.RWMutex

func (apiServer APIServer) Endpoint(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	switch path := r.URL.Path[1:]; {
	case path == "deviceList" || path == "":
		apiServer.deviceList(w)
	case strings.HasPrefix(path, "status/"):
		apiServer.deviceStatus(w, strings.Split(path, "status/")[1])
	case strings.HasPrefix(path, "turnOn/"):
		apiServer.turnOn(w, strings.Split(path, "turnOn/")[1])
	case strings.HasPrefix(path, "turnOff/"):
		apiServer.turnOff(w, strings.Split(path, "turnOff/")[1])
	case path == "health/live" || path == "health/ready":
		fmt.Fprintf(w, "ok")
	}
	mu.Unlock()
}

func getDevice(name string) (Device, ErrorResponse) {
	var device Device
	var err ErrorResponse
	for _, d := range config.Devices {
		if d.Name == name {
			return d, err
		}
	}
	err.Error = "Device '" + name + "' not found."
	return device, err
}

func SetLastOn(deviceName string) {
	for idx, d := range config.Devices {
		if d.Name == deviceName {
			config.Devices[idx].LastOn = time.Now()
		}
	}
}

func SetLastOff(deviceName string) {
	for idx, d := range config.Devices {
		if d.Name == deviceName {
			config.Devices[idx].LastOff = time.Now()
		}
	}
}

func SetStatus(deviceName string, status string) {
	for idx, d := range config.Devices {
		if d.Name == deviceName {
			config.Devices[idx].LastStatus = time.Now()
			config.Devices[idx].Status = status
		}
	}
}

func getStatusString(device Device) string {
	status := getStatus(device.IP, device.Channel, config.Key)
	if status == 1 {
		return "on"
	} else if status == 0 {
		return "off"
	} else {
		return "Unknown status: " + strconv.Itoa(status)
	}
}

func (apiServer APIServer) deviceList(w http.ResponseWriter) {
	var devices []Device
	for _, d := range config.Devices {
		d.Status = getStatusString(d)
		devices = append(devices, d)
	}
	writeResponse(w, devices, false)
}

func (apiServer APIServer) deviceStatus(w http.ResponseWriter, deviceName string) {
	device, err := getDevice(deviceName)
	if err.Error != "" {
		writeResponse(w, err, true)
		return
	}

	if !time.Now().After(device.LastStatus.Add(config.AntiSpam)) {
		log.Println("Already got status in the last 5 seconds")
		for idx, d := range config.Devices {
			if d.Name == deviceName {
				writeResponse(w, config.Devices[idx], false)
				return
			}
		}
	}

	device.Status = getStatusString(device)
	writeResponse(w, device, false)

	SetStatus(deviceName, device.Status)
}

func (apiServer APIServer) turnOn(w http.ResponseWriter, deviceName string) {
	device, err := getDevice(deviceName)
	if err.Error != "" {
		writeResponse(w, err, true)
		return
	}

	if !time.Now().After(device.LastOn.Add(config.AntiSpam)) {
		log.Println("Already turned on in the last 5 seconds")
		writeResponse(w, device, false)
		return
	}

	turnOn(device.IP, device.Channel, config.Key)
	apiServer.deviceStatus(w, device.Name)

	SetLastOn(deviceName)
}

func (apiServer APIServer) turnOff(w http.ResponseWriter, deviceName string) {
	device, err := getDevice(deviceName)
	if err.Error != "" {
		writeResponse(w, err, true)
		return
	}

	if !time.Now().After(device.LastOff.Add(config.AntiSpam)) {
		log.Println("Already turned off in the last 5 seconds")
		writeResponse(w, device, false)
		return
	}

	turnOff(device.IP, device.Channel, config.Key)
	apiServer.deviceStatus(w, device.Name)

	SetLastOff(deviceName)
}
