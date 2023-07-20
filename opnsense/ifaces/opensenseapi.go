package ifaces

type IOpnSenseApi interface {
	ModifyingRequest(module string, controller string, command string, data string, params []string) (string, error)
	NonModifyingRequest(module string, controller string, command string, params []string) (string, error)
	Core() ICore
}
