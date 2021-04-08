package storage

import (
	"fmt"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/resolver"
	"github.com/UpCloudLtd/cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"
)

type cloneCommand struct {
	*commands.BaseCommand
	resolver.CachingStorage
	params cloneParams
}

type cloneParams struct {
	request.CloneStorageRequest
}

// CloneCommand creates the "storage clone" command
func CloneCommand() commands.NewCommand {
	return &cloneCommand{
		BaseCommand: commands.New("clone", "Clone a storage"),
	}
}

var defaultCloneParams = &cloneParams{
	CloneStorageRequest: request.CloneStorageRequest{
		Tier: upcloud.StorageTierHDD,
	},
}

// InitCommand implements Command.InitCommand
func (s *cloneCommand) InitCommand() {
	s.SetPositionalArgHelp(positionalArgHelp)
	// s.ArgCompletion(getStorageArgumentCompletionFunction(s.service))
	s.params = cloneParams{CloneStorageRequest: request.CloneStorageRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Tier, "tier", defaultCloneParams.Tier, "The storage tier to use.")
	flagSet.StringVar(&s.params.Title, "title", defaultCloneParams.Title, "A short, informational description.")
	flagSet.StringVar(&s.params.Zone, "zone", defaultCloneParams.Zone, "The zone in which the storage will be created, e.g. fi-hel1.")

	s.AddFlags(flagSet)
}

// Execute implements Command.MakeExecuteCommand
func (s *cloneCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {

	if s.params.Zone == "" || s.params.Title == "" {
		return nil, fmt.Errorf("title and zone are required")
	}

	svc := exec.Storage()
	req := s.params.CloneStorageRequest
	req.UUID = uuid

	msg := fmt.Sprintf("Cloning storage %v", uuid)
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()

	res, err := svc.CloneStorage(&req)
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	logline.SetMessage(fmt.Sprintf("%s: success", msg))
	logline.MarkDone()

	return output.Marshaled{Value: res}, nil
}
