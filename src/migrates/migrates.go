package migrates

type Migrates interface {
	MigrateV11()
}

type ProjectMigrate interface {
	Migrate()
}