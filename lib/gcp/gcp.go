package gcp

import (
	"context"
	"os/user"
	"path/filepath"

	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/mlog"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
	ini "gopkg.in/ini.v1"
)

func GetProjectID(ctx context.Context) (string, error) {
	credentials, err := google.FindDefaultCredentials(ctx, compute.ComputeScope)
	if err != nil {
		mlog.Warn(ctx, merr.Wrap(ctx, err, "gcp_default_credentials_unavailable", nil))

		return "", err
	}

	if credentials.ProjectID != "" {
		return credentials.ProjectID, nil
	}

	u, err := user.Current()
	if err != nil {
		return "", err
	}

	path := filepath.Join(u.HomeDir, ".config", "gcloud", "configurations", "config_default")

	cfg, err := ini.Load(path)
	if err != nil {
		return "", err
	}

	projectID := cfg.Section("core").Key("project").String()

	if projectID == "" {
		return "", merr.New(ctx, "gcp_project_id_missing", merr.M{"path": path})
	}

	return projectID, nil
}
