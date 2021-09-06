package server_test

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands/server"
)

func TestStopCommand(t *testing.T) {
	t.Parallel()
	targetMethod := "StopServer"

	Server1 := upcloud.Server{
		State: upcloud.ServerStateMaintenance,
		Title: "server-1-title",
		UUID:  "1fdfda29-ead1-4855-b71f-1e33eb2ca9de",
	}

	servers := &upcloud.Servers{
		Servers: []upcloud.Server{
			Server1,
		},
	}

	details := upcloud.ServerDetails{
		Server: Server1,
	}

	details2 := upcloud.ServerDetails{
		Server: upcloud.Server{
			State: upcloud.ServerStateStopped,
			Title: "server-1-title",
			UUID:  "1fdfda29-ead1-4855-b71f-1e33eb2ca9de",
		},
	}

	for _, test := range []struct {
		name    string
		args    []string
		stopReq request.StopServerRequest
	}{
		{
			name: "flags mapped to the correct field",
			args: []string{
				"--type", "hard",
			},
			stopReq: request.StopServerRequest{
				UUID:     Server1.UUID,
				StopType: upcloud.StopTypeHard,
			},
		},
	} {
		// grab a local reference for parallel tests
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			server.SimpleServerCommandTests(t, server.StopCommand(), servers, Server1, details2, targetMethod, &test.stopReq, &details, test.args)
		})
	}
}
