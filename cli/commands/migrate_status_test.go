package commands

import (
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/briandowns/spinner"
	"github.com/hasura/graphql-engine/cli"
	"github.com/hasura/graphql-engine/cli/migrate"
	"github.com/hasura/graphql-engine/cli/version"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func testMigrateStatus(t *testing.T, endpoint *url.URL, migrationsDir string, expectedStatus *migrate.Status) {
	logger, _ := test.NewNullLogger()
	opts := &migrateStatusOptions{
		EC: &cli.ExecutionContext{
			Logger:       logger,
			Spinner:      spinner.New(spinner.CharSets[7], 100*time.Millisecond),
			MigrationDir: migrationsDir,
			ServerConfig: &cli.ServerConfig{
				Endpoint:       endpoint.String(),
				AdminSecret:    os.Getenv("HASURA_GRAPHQL_TEST_ADMIN_SECRET"),
				ParsedEndpoint: endpoint,
			},
		},
	}

	opts.EC.Version = version.New()
	v, err := version.FetchServerVersion(opts.EC.ServerConfig.Endpoint)
	if err != nil {
		t.Fatalf("getting server version failed: %v", err)
	}
	opts.EC.Version.SetServerVersion(v)

	status, err := opts.run()
	if err != nil {
		t.Fatalf("failed fetching migration status: %v", err)
	}
	assert.Equal(t, expectedStatus, status)
}

func TestMigrateStatusWithInvalidEndpoint(t *testing.T) {
	logger, _ := test.NewNullLogger()
	opts := &migrateStatusOptions{
		EC: &cli.ExecutionContext{
			Logger:       logger,
			Spinner:      spinner.New(spinner.CharSets[7], 100*time.Millisecond),
			MigrationDir: filepath.Join(os.TempDir(), "hasura-cli-test-"+strconv.Itoa(rand.Intn(1000))),
			ServerConfig: &cli.ServerConfig{
				Endpoint:       ":",
				AdminSecret:    "",
				ParsedEndpoint: &url.URL{},
			},
		},
	}

	opts.EC.Version = version.New()
	opts.EC.Version.SetServerVersion("")
	_, err := opts.run()
	if err == nil {
		t.Fatalf("expected err not to be nil")
	}
}
