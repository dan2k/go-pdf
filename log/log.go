package log

import(
	"log"
	"os"
	"github.com/joho/godotenv"
	"time"
)
const (
	DDMMYYYYhhmmss = "20060102-15:04:05"
)
var L *log.Logger
var Envs map[string]string 
 func  New() {
	envs, err := godotenv.Read(".env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }
	outfile, _ := os.Create(envs["LOGFILE"]+"-"+time.Now().UTC().Format(DDMMYYYYhhmmss)) // update path for your needs
	Envs=envs
	L          = log.New(outfile, "", 0)
} 
