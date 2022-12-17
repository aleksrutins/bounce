package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RequestPayload struct {
	Id       string `json:"id"`
	Code     string `json:"code"`
	Input    string `json:"input"`
	Language string `json:"language"`
}

type ResponsePayload struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exitCode"`
}

func main() {
	ctx := context.Background()
	vms := make([]*firecracker.Machine, 0)

	// Read config from environment variables
	vmLimit, err := strconv.ParseInt(os.Getenv("MAX_VMS"), 10, 64)
	handleErr(err)
	vcpus, err := strconv.ParseInt(os.Getenv("VCPUS"), 10, 64)
	handleErr(err)
	mem, err := strconv.ParseInt(os.Getenv("MEM_SIZE"), 10, 64)
	handleErr(err)
	drivePath := os.Getenv("DRIVE_PATH")
	kernelPath := os.Getenv("KERNEL_PATH")

	// Pre-boot VMs
	for i := int64(0); i < vmLimit; i++ {
		vm, err := PrepareVM(ctx, drivePath, kernelPath, vcpus, mem)
		handleErr(err)
		vm, err = StartVM(ctx, vm)
		handleErr(err)
		vms = append(vms, vm)
	}

	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	handleErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	handleErr(err)
	defer ch.Close()

	q, err := ch.QueueDeclare("exec_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	handleErr(err)

	err = ch.Qos(
		int(vmLimit),
		0,
		false,
	)
	handleErr(err)

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		true,
		nil,
	)
	handleErr(err)

	var forever chan struct{}

	go func() {
		for d := range msgs {
			var request RequestPayload
			var response ResponsePayload

			if err := json.Unmarshal(d.Body, &request); err != nil {
				response.Message = "Failed to parse request"
				response.Success = false
				body, err := json.Marshal(response)
				handleErr(err)

				err = ch.PublishWithContext(ctx,
					"",
					d.ReplyTo,
					false,
					false,
					amqp.Publishing{
						ContentType:   "text/json",
						CorrelationId: d.CorrelationId,
						Body:          body,
					})
				handleErr(err)
			} else {
				vm := vms[0]
				vms = vms[1:]

				go func() {
					newVm, err := PrepareVM(ctx, drivePath, kernelPath, vcpus, mem)
					handleErr(err)
					newVm, err = StartVM(ctx, newVm)
					handleErr(err)

					// Health check
					err = fmt.Errorf("healthcheck")
					for err != nil {
						_, err = http.Get(fmt.Sprintf("http://%s:3005/health", &newVm.Cfg.NetworkInterfaces[0].StaticConfiguration.IPConfiguration.IPAddr.IP))
					}

					vms = append(vms, newVm)
					d.Ack(false)
				}()

				go func() {
					ExecCode(vm, d.Body, &response)
					body, err := json.Marshal(response)
					handleErr(err)

					err = ch.PublishWithContext(ctx,
						"",
						d.ReplyTo,
						false,
						false,
						amqp.Publishing{
							ContentType:   "text/json",
							CorrelationId: d.CorrelationId,
							Body:          body,
						})
					handleErr(err)
					StopVM(vm)
				}()
			}
		}
	}()

	fmt.Println("Started worker")
	<-forever
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
