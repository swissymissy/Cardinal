package main 

import (
	"github.com/swissymissy/Cardinal/internal/database"
)

type WorkerConfig struct {
	DB *database.Queries
	SMTPHost string 
	SMTPPort int 
	SMTPUsername string 
	SMTPPassword string
}
