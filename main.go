package main

import (
	"EmailBomber/internal/workers"
	"bufio"
	"fmt"
	"github.com/fatih/color"
	flag "github.com/spf13/pflag"
	"net/smtp"
	"os"
	"sync"
)

const version = "v0.1.0"

var args = struct {
	amount         *int
	margin         *int
	host           *string
	port           *string
	name           *string
	password       *string
	email          *string
	receivingEmail *string
	content        *string
	numGoroutines  *int
	version        *bool
}{
	amount:         flag.IntP("amount", "a", 1, ""),
	margin:         flag.IntP("margin", "m", 5, ""),
	host:           flag.StringP("host", "h", "smtp.gmail.com", ""),
	port:           flag.StringP("port", "p", "587", ""),
	name:           flag.StringP("name", "n", "", ""),
	password:       flag.StringP("password", "P", "", ""),
	email:          flag.StringP("email", "e", "", ""),
	receivingEmail: flag.StringP("receiving-email", "r", "", ""),
	content:        flag.StringP("content", "c", "", ""),
	numGoroutines:  flag.IntP("threads", "t", 1, ""),
	version:        flag.BoolP("version", "v", false, ""),
}

func printfInfo(format string, a ...interface{}) {
	c := color.New(color.FgHiBlue)
	_, _ = c.Printf("[Info] ")
	fmt.Printf(format, a...)
}

func printfError(format string, a ...interface{}) {
	c := color.New(color.FgHiRed)
	_, _ = c.Printf("[Error] "+format, a...)
}

func printfWarn(format string, a ...interface{}) {
	c := color.New(color.FgHiYellow)
	_, _ = c.Printf("[Warning] "+format, a...)
}

func printfSuccess(format string, a ...interface{}) {
	c := color.New(color.FgGreen)
	_, _ = c.Printf("[Success] "+format, a...)
}

// fatalf is a call to printfError followed by a call to os.Exit
func fatalf(format string, a ...interface{}) {
	printfError(format, a...)
	os.Exit(1)
}

func flagUsage() {
	fmt.Printf(`Usage: emailbomber [-n | -h] [-P] [options...]

Parameters:
	-a, --amount:           The amount of emails to be sent
	-m, --margin:           The amount of times an email will be resend in case of delay error (default 5)
	-n, --name:             The name that will be used as sender in the emails.
	-h, --host:             The smtp host for outgoing emails. (default smtp.gmail.com)
	-p, --port:             Port for smtp email sending. (default 587)
	-P, --password:         Password to the email adress used
	-e, --email:            The email adress used to send the emails.
	-r, --receiving-email:  The email adress receiving the emails
	-c, --content:          File path to content file
	-t, --threads:          Maximum amount of concurrent outgoing email attempts. (default 1)
	-v, --version:          Prints the installed version of emailbomber.

Example: ./EmailBomber -a 50 -t 20 -n YourName -e youremail@gmail.com -r theiremail@gmail.com -P "yourpassword" -c yourtextfile.txt -m 2
`,
	)
}

func main() {
	flag.Usage = flagUsage
	flag.ErrHelp = nil

	flag.Parse()

	// Validate arguments.
	if *args.version {
		fmt.Printf(version + "\n")
		os.Exit(0)
	}
	if *args.host == "" {
		flagUsage()
		fmt.Printf("invalid arguments: you must provide either \"-h, --host\"\n")
		os.Exit(2)
	}
	if *args.port == "" {
		flagUsage()
		fmt.Printf("invalid arguments: you must provide \"-p, --port\"\n")
		os.Exit(2)
	}
	if *args.numGoroutines < 1 {
		flagUsage()
		fmt.Printf("invalid argument %d for \"-t, --threads\": number of threads must be at least 1\n", *args.numGoroutines)
		os.Exit(2)
	}
	if *args.numGoroutines > 20 {
		printfWarn("setting a high number of maximum concurrent connections might cause instability such as skipped emails\n")
	}
	if *args.amount < 1 {
		flagUsage()
		fmt.Printf("invalid argument %d for \"-a, --amount\": number of emails must be at least 1\n", *args.numGoroutines)
		os.Exit(2)
	}
	if *args.amount > 100 {
		printfWarn("setting a high number of emails can result in crashed and or the program breaking\n")
	}

	if *args.margin < 1 {
		flagUsage()
		fmt.Printf("invalid argument %d for \"-m, --margin\": margin must be at least 20\n", *args.numGoroutines)
		os.Exit(2)
	}
	if *args.margin > 5 {
		printfWarn("setting a high margin can result in longer loading times or a rate limit from google\n")
	}

	if *args.password == "" {
		flagUsage()
		fmt.Printf("invalid arguments: you must provide \"-P, --password\"\n")
		os.Exit(2)
	}
	if *args.content == "" {
		flagUsage()
		fmt.Printf("invalid arguments: you must provide \"-c, --content\"\n")
		os.Exit(2)
	}
	if *args.name == "" {
		flagUsage()
		fmt.Printf("invalid arguments: you must provide \"-n, --name\"\n")
		os.Exit(2)
	}
	if *args.receivingEmail == "" {
		flagUsage()
		fmt.Printf("invalid arguments: you must provide \"-e, --email\"\n")
		os.Exit(2)
	}

	if *args.numGoroutines > *args.amount {
		fatalf("invalid arguments: threads cant be larger than amount of emails")
	}
	f, err := os.Open(*args.content)
	if err != nil {
		fmt.Println("An error occurred while opening content file: ", err)
		return
	}

	defer f.Close()
	content := ""
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		content = content + "\n" + sc.Text()
	}
	if sc.Err() != nil {
		fmt.Println("An error occurred while reading content from file: ", err)
		return
	}

	pool, err := workers.NewPool(*args.numGoroutines)
	if err != nil {
		fmt.Println("An error occurred while creating worker pool: ", err)
		os.Exit(1)
	}

	workerWG := &sync.WaitGroup{}
	emailsPerWorker := *args.amount / *args.numGoroutines
	workerWG.Add(*args.numGoroutines)

	task := workers.Task{
		Fn: func(params []interface{}) {
			var (
				workerId     = params[0].(int)
				amount       = params[1].(int)
				margin       = params[2].(int)
				name         = params[3].(string)
				emailAddress = params[4].(string)
				password     = params[5].(string)
				to           = params[6].([]string)
				content      = params[7].(string)
				host         = params[8].(string)
				port         = params[9].(string)
				workerWG     = params[10].(*sync.WaitGroup)
			)
			defer workerWG.Done()

			margin = margin * amount
			for i := 0; i < amount; i++ {
				err := email(name, emailAddress, password, content, to, host, port)
				if err != nil {
					i--
					margin--
				}
				if margin < 1 {
					printfError("An error occurred while sending email: tries exceeded\n")
					break
				}
			}
			printfInfo("Worker %d has finished\n", workerId)
		},
	}

	for i := 0; i < *args.numGoroutines; i++ {
		task.Params = []interface{}{i, emailsPerWorker, *args.margin, *args.name, *args.email, *args.password, []string{*args.receivingEmail}, content, *args.host, *args.port, workerWG}
		if err := pool.Queue(task); err != nil {
			fmt.Println("An error occurred while queuing tasks: ", err)
			os.Exit(1)
		}
	}

	if err := pool.Start(); err != nil {
		fmt.Println("An error occurred while starting tasks: ", err)
		os.Exit(1)
	}
	defer pool.Close()

	workerWG.Wait()
	printfSuccess("All emails sent")
}

func email(name string, sendEmail string, password string, content string, to []string, host string, port string) error {
	auth := smtp.PlainAuth("", sendEmail, password, host)
	msg := []byte(content)

	err := smtp.SendMail(host+":"+port, auth, name, to, msg)
	if err != nil {
		return err
	}

	return nil
}
