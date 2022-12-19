package main

import (
	"context"
	"fmt"
	"io"
	"os"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	"github.com/google/uuid"
)

func PrepareVM(ctx context.Context, drivePath string, kernelPath string, vcpus int64, mem int64) (*firecracker.Machine, error) {
	id := uuid.New().String()

	vmDrivePath := fmt.Sprintf("/tmp/%s-drive.ext4", id)
	if err := CopyFile(drivePath, vmDrivePath); err != nil {
		return nil, fmt.Errorf("failed to copy drive: %w", err)
	}

	vmKernelPath := fmt.Sprintf("/tmp/%s-kernel.bin", id)
	if err := CopyFile(kernelPath, vmKernelPath); err != nil {
		return nil, fmt.Errorf("failed to copy kernel: %w", err)
	}

	drive := []models.Drive{{
		DriveID:      firecracker.String("1"),
		PathOnHost:   firecracker.String(vmDrivePath),
		IsReadOnly:   firecracker.Bool(false),
		IsRootDevice: firecracker.Bool(true),
	}}

	network := []firecracker.NetworkInterface{{
		CNIConfiguration: &firecracker.CNIConfiguration{
			NetworkName: "fcnet",
			IfName:      "veth0",
		},
	}}

	config := firecracker.Config{
		VMID:              id,
		SocketPath:        fmt.Sprintf("/tmp/%s-socket.sock", id),
		KernelImagePath:   vmKernelPath,
		Drives:            drive,
		NetworkInterfaces: network,
		MachineCfg: models.MachineConfiguration{
			VcpuCount:  firecracker.Int64(vcpus),
			MemSizeMib: firecracker.Int64(mem),
		},
	}

	vm, err := firecracker.NewMachine(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create VM: %w", err)
	}

	fmt.Printf("Successfully configured VM (%s)\n", vm.Cfg.VMID)
	return vm, nil
}

func StartVM(ctx context.Context, vm *firecracker.Machine) (*firecracker.Machine, error) {
	if err := vm.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start VM: %w", err)
	}
	return vm, nil
}

func StopVM(vm *firecracker.Machine) error {
	if err := vm.StopVMM(); err != nil {
		return fmt.Errorf("failed to shutdown VM: %w", err)
	}
	if err := os.Remove(vm.Cfg.SocketPath); err != nil {
		return fmt.Errorf("failed to remove socket: %w", err)
	}
	if err := os.Remove(vm.Cfg.KernelImagePath); err != nil {
		return fmt.Errorf("failed to remove kernel: %w", err)
	}
	if err := os.Remove(*vm.Cfg.Drives[0].PathOnHost); err != nil {
		return fmt.Errorf("failed to remove drive: %w", err)
	}
	return nil
}

func CopyFile(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}
