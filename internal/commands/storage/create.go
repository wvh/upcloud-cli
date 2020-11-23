package storage

import (
	"fmt"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
)

func CreateCommand(service service.Storage) commands.Command {
	return &createCommand{
		BaseCommand: commands.New("create", "Create a storage"),
		service:     service,
	}
}

var DefaultCreateParams = &createParams{
	CreateStorageRequest: request.CreateStorageRequest{
		Size: 10,
		Tier: "maxiops",
		BackupRule: &upcloud.BackupRule{
			Interval:  "daily",
			Retention: 7,
		},
	},
}

func newCreateParams() createParams {
	return createParams{CreateStorageRequest: request.CreateStorageRequest{BackupRule: &upcloud.BackupRule{}}}
}

type createParams struct {
	request.CreateStorageRequest
	backupTime string
}

func (s *createParams) processParams() error {
	if s.backupTime != "" {
		tv, err := time.Parse("15:04", s.backupTime)
		if err != nil {
			return fmt.Errorf("invalid backup time %q", s.backupTime)
		}
		s.BackupRule.Time = tv.Format("1504")
	} else {
		s.BackupRule = nil
	}
	return nil
}

type createCommand struct {
	*commands.BaseCommand
	service            service.Storage
	params 						 createParams
	flagSet            *pflag.FlagSet
}

func createFlags(fs *pflag.FlagSet, dst, def *createParams) {
	fs.StringVar(&dst.Title, "title", def.Title, "Storage title")
	fs.IntVar(&dst.Size, "size", def.Size, "Size of the storage in GiB")
	fs.StringVar(&dst.Zone, "zone", def.Zone, "The zone to create the storage on")
	fs.StringVar(&dst.Tier, "tier", def.Tier, "Storage tier")
	fs.StringVar(&dst.backupTime, "backup-time", def.backupTime, "The time when to create a backup in HH:MM. Empty value means no backups.")
	fs.StringVar(&dst.BackupRule.Interval, "backup-interval", def.BackupRule.Interval, "The interval of the backup.\nAvailable: daily,mon,tue,wed,thu,fri,sat,sun")
	fs.IntVar(&dst.BackupRule.Retention, "backup-retention", def.BackupRule.Retention, "How long to store the backups in days. The accepted range is 1-1095")
}

func (s *createCommand) InitCommand() {
	s.flagSet = &pflag.FlagSet{}
	s.params = newCreateParams()
	createFlags(s.flagSet, &s.params, DefaultCreateParams)
	s.AddFlags(s.flagSet)
}

func (s *createCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if err := s.params.processParams(); err != nil {
			return nil, err
		}

		createStorages := []*request.CreateStorageRequest{&s.params.CreateStorageRequest}

		return ui.HandleContext{
			RequestId:     func(in interface{}) string { return in.(*request.CreateStorageRequest).Title },
			ResultUuid:    getStorageDetailsUuid,
			InteractiveUi: s.Config().InteractiveUI(),
			MaxActions:    maxStorageActions,
			ActionMsg:     "Creating storage",
			Action: func(req interface{}) (interface{}, error) {
				return s.service.CreateStorage(req.(*request.CreateStorageRequest))
			},
		}.HandleAction(createStorages)
	}
}
