package log

import(
	"log"
	"os"
	"github.com/joho/godotenv"
)
var L *log.Logger
var Envs map[string]string 
 func  New(runtime string) {
	envs, err := godotenv.Read(".env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }
	if _, err := os.Stat(envs["LOGDIR"]); err !=nil {
		if err := os.MkdirAll(envs["LOGDIR"], os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
	outfile, _ := os.Create(envs["LOGDIR"]+"/"+envs["LOGFILE"]+"-"+runtime) // update path for your needs
	Envs=envs
	L          = log.New(outfile, "", 0)
} 
