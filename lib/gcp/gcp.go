package gcp

import (
	"context"
	"os/user"
	"path/filepath"

	"github.com/cuvva/cuvva-public-go/lib/clog"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
	ini "gopkg.in/ini.v1"
)

func GetProjectID(ctx context.Context) (string, error) {
	credentials, err := google.FindDefaultCredentials(ctx, compute.ComputeScope)
	if err != nil {
		clog.Get(ctx).WithError(err).Warn("gcp_default_credentials_unavailable")

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
		return "", merr.New("gcp_project_id_missing", merr.M{"path": path})
	}

	return projectID, nil
}
