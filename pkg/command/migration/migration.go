package migration

import (
	"bytes"
	"flag"
	"fmt"

	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database"
	"github.com/mattes/migrate/source/file"
	log "github.com/sirupsen/logrus"
	// Enables Postgres for migrations
	_ "github.com/mattes/migrate/database/postgres"
	// Enables stub DB for testing
	_ "github.com/mattes/migrate/database/stub"
)

// CreateMigrateCommand knows how to properly create the MigrateCommand
func CreateMigrateCommand(db database.Driver, dir string, version uint) *MigrateCommand {
	cmd := &MigrateCommand{db: db}
	cmd.flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	cmd.flags.UintVar(&cmd.version, "version", version, "DB Version to migrate to, i.e. 1")
	cmd.flags.StringVar(&cmd.dir, "dir", dir, "relative or absolute path for migration files source dir")

	return cmd
}

// MigrateCommand represents a CLI migrate command
type MigrateCommand struct {
	flags *flag.FlagSet
	db    database.Driver

	// cli args
	version uint
	dir     string
}

// Run execute the migration command
func (migrateCmd *MigrateCommand) Run(args []string) int {
	if err := migrateCmd.flags.Parse(args); err != nil {
		log.Error(err)
		return 1
	}

	if migrateCmd.dir == "" {
		log.Errorf("empty option: -dir. args: %v", args)
		return 1
	}

	if migrateCmd.version == 0 {
		log.Errorf("empty option: -version. args: %v", args)
		return 1
	}

	f := &file.File{}
	sourceDriver, err := f.Open("file://" + migrateCmd.dir)
	if err != nil {
		log.Error(err)
		return 1
	}

	// Configure migration
	migration, err := migrate.NewWithInstance(
		"file",
		sourceDriver,
		"database",
		migrateCmd.db,
	)
	if err != nil {
		log.Error(err)
		return 1
	}

	// Migrate the system to the correct version
	err = migration.Migrate(migrateCmd.version)
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Infof("Database already at version %d", migrateCmd.version)
			return 0
		}

		log.Errorf("Database migration error: %s", err)
		return 1
	}

	log.Infof("Migrated database to version %d", migrateCmd.version)
	return 0
}

// Help for the migrate command
func (migrateCmd *MigrateCommand) Help() string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "Migrate database to the given version\n\n")

	migrateCmd.flags.VisitAll(func(f *flag.Flag) {
		ftype, usage := flag.UnquoteUsage(f)
		fmt.Fprintf(&buf, "\t-%s %s\n\t\t%s\n", f.Name, ftype, usage)
	})

	return buf.String()
}

// Synopsis of the migrate command
func (migrateCmd *MigrateCommand) Synopsis() string {
	return "Migrate database version"
}
