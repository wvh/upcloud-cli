package ip_address

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
)

type removeCommand struct {
	*commands.BaseCommand
	service service.IpAddress
}

func RemoveCommand(service service.IpAddress) commands.Command {
	return &removeCommand{
		BaseCommand: commands.New("remove", "Removes an ip address"),
		service:     service,
	}
}

func (s *removeCommand) InitCommand() {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(GetArgCompFn(s.service))
}

func (s *removeCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return Request{
			BuildRequest: func(address string) interface{} {
				return &request.ReleaseIPAddressRequest{IPAddress: address}
			},
			Service: s.service,
			HandleContext: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.ReleaseIPAddressRequest).IPAddress },
				MaxActions:    maxIpAddressActions,
				InteractiveUI: s.Config().InteractiveUI(),
				ActionMsg:     "Removing IP Address",
				Action: func(req interface{}) (interface{}, error) {
					return nil, s.service.ReleaseIPAddress(req.(*request.ReleaseIPAddressRequest))
				},
			},
		}.Send(args)
	}
}