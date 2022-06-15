package python3

type Migrates interface {
	MigrateV11()
}

type ProjectMigrate interface {
	Migrate()
}
