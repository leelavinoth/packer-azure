// Copyright (c) Microsoft Open Technologies, Inc.
// All Rights Reserved.
// Licensed under the Apache License, Version 2.0.
// See License.txt in the project root for license information.
package win

import (
	"fmt"
	"github.com/MSOpenTech/packer-azure/packer/builder/azure/driver_restapi/constants"
	"github.com/MSOpenTech/packer-azure/packer/builder/azure/driver_restapi/request"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

type StepCreateVm struct {
	StorageAccount   string
	StorageContainer string
	TmpVmName        string
	TmpServiceName   string
	InstanceSize     string
	Username         string
	Password         string
}

func (s *StepCreateVm) Run(state multistep.StateBag) multistep.StepAction {
	reqManager := state.Get(constants.RequestManager).(*request.Manager)
	ui := state.Get("ui").(packer.Ui)

	errorMsg := "Error Creating Temporary Azure VM: %s"
	var err error

	ui.Say("Creating Temporary Azure VM...")

	osImageName := state.Get(constants.OsImageName).(string)
	if len(osImageName) == 0 {
		err := fmt.Errorf(errorMsg, fmt.Errorf("osImageName is empty"))
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	isOSImage := state.Get(constants.IsOSImage).(bool)

	mediaLoc := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s.vhd", s.StorageAccount, s.StorageContainer, s.TmpVmName)

	requestData := reqManager.CreateVirtualMachineDeploymentWin(isOSImage, s.TmpServiceName, s.TmpVmName, s.InstanceSize, s.Username, s.Password, osImageName, mediaLoc)

	err = reqManager.ExecuteSync(requestData)

	if err != nil {
		err := fmt.Errorf(errorMsg, err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put(constants.VmExists, 1)
	state.Put(constants.DiskExists, 1)

	return multistep.ActionContinue
}

func (s *StepCreateVm) Cleanup(state multistep.StateBag) {
	// do nothing
}
