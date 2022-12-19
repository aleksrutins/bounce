package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
)

func ExecCode(vm *firecracker.Machine, request []byte, response *ResponsePayload) {
	ip := vm.Cfg.NetworkInterfaces[0].StaticConfiguration.IPConfiguration.IPAddr.IP

	resp, err := http.Post(fmt.Sprintf("http://%s:3005/exec", ip), "application/json", bytes.NewBuffer(request))
	handleErr(err)

	json.NewDecoder(resp.Body).Decode(response)
	response.Success = true
}
