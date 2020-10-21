package database

import (
	"context"

	migrations "github.com/mjm/pi-tools/detect-presence/database/migrate"
	"github.com/mjm/pi-tools/pkg/migrate"
)

func (c *Client) MigrateIfNeeded(_ context.Context) error {
	return migrate.UpIfNeeded(c.db, migrations.Data)
}
