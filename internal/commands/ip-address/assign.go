package ip_address


import (
  "github.com/UpCloudLtd/cli/internal/commands"
  "github.com/UpCloudLtd/cli/internal/ui"
  "github.com/UpCloudLtd/upcloud-go-api/upcloud"
  "github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
  "github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
  "github.com/spf13/pflag"
)

type assignCommand struct {
  *commands.BaseCommand
  service service.IpAddress
  req request.AssignIPAddressRequest
  floating bool
}

func AssignCommand(service service.IpAddress) commands.Command {
  return &assignCommand{
    BaseCommand: commands.New("assign", "Assign an ip address"),
    service:     service,
  }
}

var defCreateParams = request.AssignIPAddressRequest{
  Access:     "public",
  Family:     "IPv4",
}

func (s *assignCommand) InitCommand() {
  fs := &pflag.FlagSet{}
  fs.StringVar(&s.req.Access, "access", defCreateParams.Access, "Is address for utility or public network.")
  fs.StringVar(&s.req.Family, "family", defCreateParams.Family, "The address family of new IP address.")
  fs.StringVar(&s.req.ServerUUID, "server", defCreateParams.ServerUUID, "The server the ip address is assigned to.")
  fs.StringVar(&s.req.MAC, "mac", defCreateParams.MAC, "MAC address of server interface to assign address to. Required for non-floating addresses.")
  fs.StringVar(&s.req.Zone, "zone", defCreateParams.Zone, "Zone of address, required when assigning a detached floating IP address.")
  fs.BoolVar(&s.floating, "floating", false, "Whether the address to be assigned is a floating one.")
  s.AddFlags(fs)
}

func (s *assignCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
  return func(args []string) (interface{}, error) {
    return Request{
      BuildRequest: func(ip *upcloud.IPAddress) interface{} {
        return &s.req
      },
      Service: s.service,
      HandleContext: ui.HandleContext{
        RequestID:     func(in interface{}) string { return in.(*request.AssignIPAddressRequest).MAC },
        MaxActions:    maxIpAddressActions,
        InteractiveUI: s.Config().InteractiveUI(),
        ActionMsg:     "Assigning IP Address",
        Action: func(req interface{}) (interface{}, error) {
          return s.service.AssignIPAddress(req.(*request.AssignIPAddressRequest))
        },
      },
    }.Send(args)
  }
}
