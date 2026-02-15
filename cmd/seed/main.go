package main

import (
	"flag"
	"log"

	"github.com/Nurdiansyah15/ddd-arch/internal/config"
	"github.com/Nurdiansyah15/ddd-arch/internal/infrastructure/database/postgres"
	"github.com/joho/godotenv"
)

func main() {
	// Flag untuk migration path & action
	migrationPath := flag.String("path", "internal/infrastructure/database/postgres/seeders", "Migration files path")
	action := flag.String("action", "up", "Migration action: up, down, or fresh")
	flag.Parse()

	_ = godotenv.Load()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Init DB
	db, err := postgres.NewPostgres(cfg.DB.DSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	switch *action {
	case "up":
		if err := postgres.SeedUp(db, *migrationPath); err != nil {
			log.Fatal("Seeder UP failed:", err)
		}
		log.Println("Seeders UP completed successfully")

	case "down":
		if err := postgres.SeedDown(db, *migrationPath); err != nil {
			log.Fatal("Seeder DOWN failed:", err)
		}
		log.Println("Seeders DOWN completed successfully")

	case "fresh":
		log.Println("Dropping all tables...")
		if err := postgres.DropAllTables(db); err != nil {
			log.Fatal("Failed to drop tables:", err)
		}
		log.Println("Running seeders UP after drop...")
		if err := postgres.SeedUp(db, *migrationPath); err != nil {
			log.Fatal("Seeder UP failed:", err)
		}
		log.Println("Fresh seeding completed successfully")

	default:
		log.Fatal("Invalid action. Use -action up, down, or fresh")
	}
}
