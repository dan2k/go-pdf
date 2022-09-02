package log

import(
	"log"
	"os"
	"github.com/joho/godotenv"
)

var L *log.Logger
var Envs map[string]string 
 func  New() {
	envs, err := godotenv.Read(".env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }
	outfile, _ := os.Create(envs["LOGFILE"]) // update path for your needs
	Envs=envs
	L          = log.New(outfile, "", 0)
} 
