package setting

import (
	"flag"
	"fmt"
	"os"
)

// DB Configuration
const (
	DB_DSN            string = "user:password@tcp(127.0.0.1:3306)/shop?charset=utf8&autocommit=1&readTimeout=60s&writeTimeout=2s&timeout=1s&loc=Asia%2FShanghai"
	DB_DRIVER         string = "mysql"
	DB_MAX_IDLE_CONNS int    = 50
	DB_MAX_OPEN_CONNS int    = 100
	DB_LOGS_FILE      string = "./logs/DB.log"
)

// File Configuration
const (
	FILE_ORIGINAL_DATA string = "./logs/original.dat"
	FILE_NEW_SQL_DATA  string = "./logs/update_order_product.sql"
	FILE_LOG_DATA      string = "./logs/log.dat"
	FILE_ERR__LOG_DATA string = "./logs/err.dat"
)

// Contral Variables
var (
	WorkerNum    uint = 1
	ProbeBegin   int
	ProbeEnd     int
	NowPage      int
	TotalPage    int
	LimitPerTime int
)

// Common Confuguration
var (
	Version     string = "1.0"
	ShowVersion bool
	IsDebug     bool
	ShowHelp    bool
)

var usageStr = `
─────────────────────────────────────────────────────────────
██████████████───██████████████─██████████████─██████████████
██░░░░░░░░░░██───██░░░░░░░░░░██─██░░░░░░░░░░██─██░░░░░░░░░░██
██░░██████░░██───██░░██████████─██░░██████░░██─██░░██████████
██░░██──██░░██───██░░██─────────██░░██──██░░██─██░░██────────
██░░██████░░████─██░░██████████─██░░██████░░██─██░░██████████
██░░░░░░░░░░░░██─██░░░░░░░░░░██─██░░░░░░░░░░██─██░░░░░░░░░░██
██░░████████░░██─██░░██████████─██░░██████░░██─██████████░░██
██░░██────██░░██─██░░██─────────██░░██──██░░██─────────██░░██
██░░████████░░██─██░░██████████─██░░██──██░░██─██████████░░██
██░░░░░░░░░░░░██─██░░░░░░░░░░██─██░░██──██░░██─██░░░░░░░░░░██
████████████████─██████████████─██████──██████─██████████████
─────────────────────────────────────────────────────────────
Usage: Beans [options]

Options:
    -d <mode>            Select Operation mode
    -w <worker number>   Set workder number
    -s <start point>     Set probe begin point
    -n <now page>        Set now page for searching
    -t <total page>      Set total page for searching
    -l <limit per time>  Set search quantity per time
    -h, --help           Show help message
    -v, --version        Show version
`

// parse parameters
func ParseParams() {
	flag.BoolVar(&ShowVersion, "v", false, "Show Beas version information")
	flag.BoolVar(&ShowHelp, "h", false, "Show Beas help information")
	flag.BoolVar(&IsDebug, "d", false, "Select Operation Mode")
	flag.UintVar(&WorkerNum, "w", 5, "Set Workder Number")
	flag.IntVar(&ProbeBegin, "s", 1, "Set Probe Begin Point")
	flag.IntVar(&NowPage, "n", 0, "Set Now Page for Searching")
	flag.IntVar(&TotalPage, "t", 1, "Set Total Page for Searching")
	flag.IntVar(&LimitPerTime, "l", 500, "Set Search Quantity Per time")
	flag.Parse()
}

// Usage will print out the flag options for beas
func Usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

// Print version out  for beas
func PrintVersion() {
	fmt.Printf("Beans version %s\n", Version)
	os.Exit(0)
}
