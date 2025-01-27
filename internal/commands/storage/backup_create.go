package storage

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"
)

type createBackupCommand struct {
	*commands.BaseCommand
	resolver.CachingStorage
	completion.Storage
	params createBackupParams
}

type createBackupParams struct {
	request.CreateBackupRequest
}

// CreateBackupCommand creates the "storage backup create" command
func CreateBackupCommand() commands.Command {
	return &createBackupCommand{
		BaseCommand: commands.New(
			"create",
			"Create backup of a storage",
			`upctl storage backup create 01cbea5e-eb5b-4072-b2ac-9b635120e5d8 --title "first backup"`,
			`upctl storage backup create "My Storage" --title second_backup`,
		),
	}
}

var defaultCreateBackupParams = &createBackupParams{
	CreateBackupRequest: request.CreateBackupRequest{},
}

// InitCommand implements Command.InitCommand
func (s *createBackupCommand) InitCommand() {
	s.params = createBackupParams{CreateBackupRequest: request.CreateBackupRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Title, "title", defaultCreateBackupParams.Title, "A short, informational description.")

	s.AddFlags(flagSet)
}

// Execute implements commands.MultipleArgumentCommand
func (s *createBackupCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.Storage()

	if s.params.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	msg := fmt.Sprintf("backing up %v", uuid)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()

	req := s.params.CreateBackupRequest
	req.UUID = uuid

	res, err := svc.CreateBackup(&req)
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	logline.SetMessage(fmt.Sprintf("%s: success", msg))
	logline.MarkDone()

	return output.OnlyMarshaled{Value: res}, nil
}
