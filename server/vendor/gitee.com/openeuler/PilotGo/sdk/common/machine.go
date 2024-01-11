package common

type MachineNode struct {
	UUID        string `json:"uuid"`
	Department  string `json:"department"`
	IP          string `json:"ip"`
	CPUArch     string `json:"cpu_arch"`
	OS          string `json:"os"`
	RunStatus   string `json:"runstatus"`
	MaintStatus string `json:"maintatatus"`
}

type Batch struct {
	BatchId       int      `json:"batch_id"`
	DepartmentIDs []string `json:"department_ids"`
	MachineUUIDs  []string `json:"machine_uuids"`
}
