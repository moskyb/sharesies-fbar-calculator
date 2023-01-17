package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/moskyb/sharesies-fbar-calculator/login"
	"github.com/moskyb/sharesies-fbar-calculator/portfolio"
	"golang.org/x/term"
)

var (
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Purple = "\033[35m"
)

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Yellow = ""
		Purple = ""
	}
}

func main() {
	fmt.Println("Welcome to Ben Moskovitz's FBAR calculator for Sharesies!")
	fmt.Println("This tool will log into your Sharesies account, go through your portfolio history, and find the maximum portfolio value for a given year. This tool requires you to enter your password and MFA code (if applicable). It then pretends to be a browser, and logs into your account. This tool doesn't store your password or MFA code, doesn't send your password or MFA code anywhere (other than to Sharesies, who have it anyway) and it doesn't store your portfolio history, but...")
	fmt.Println()
	fmt.Println("⚠️ " + Yellow + "YOU PROBABLY STILL SHOULDN'T TRUST THIS TOOL" + Reset + " ⚠️")
	fmt.Println("I'm just some guy and I could probably be stealing your stuff, even if i've said that i won't. All this tool does is look through your portfolio history, but it could be going in and selling all your stuff or whatever. Check through the code, and make sure that it's not doing anything dodgy.")
	fmt.Println()
	fmt.Println("One more thing: this tool is not affiliated with Sharesies in any way, other than that i used to work there. It's just a tool that I made to help me with my FBAR filing. I'm not a tax professional, and I'm not responsible for any tax issues you might have as a result of using this tool.")
	fmt.Println()
	fmt.Println("This tool, and the information it produces, aren't financial advice, and i'm not a financial advisor. If you're not sure what you're doing, you should probably talk to a tax professional.")
	fmt.Println()
	fmt.Println(`If you're still here, then let's get get going! If you want to continue, type "` + Purple + "I UNDERSTAND THE RISKS" + Reset + `" (without the quotes, in all caps) and press enter.`)
	fmt.Println()

	inputReader := bufio.NewReader(os.Stdin)
	fmt.Print(Purple)
	doTheyUnderstandTheRisks, err := inputReader.ReadString('\n')
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}
	fmt.Print(Reset)

	if strings.TrimSpace(doTheyUnderstandTheRisks) != "I UNDERSTAND THE RISKS" {
		fmt.Println(`That wasn't "` + Purple + "I UNDERSTAND THE RISKS" + Reset + `". No sweat! See you later, maybe :)`)
		os.Exit(0)
	}

	fmt.Println("Awesome! Let's get this show on the road")
	fmt.Println()

	fmt.Print("What year are you trying to find the maximum portfolio value for? " + Purple)
	desiredYearS, err := inputReader.ReadString('\n')
	if err != nil {
		log.Fatalf("failed to read input for desired year: %v", err)
	}
	fmt.Print(Reset)
	desiredYearS = strings.TrimSpace(desiredYearS)

	desiredYear, err := strconv.Atoi(desiredYearS)
	if err != nil {
		log.Fatalf("failed to parse desired year: %v", err)
	}

	fmt.Print("What's the email address associated with your Sharesies account? " + Purple)
	email, err := inputReader.ReadString('\n')
	if err != nil {
		log.Fatalf("failed to read input for email: %v", err)
	}
	fmt.Print(Reset)
	email = strings.TrimSpace(email)

	fmt.Print("What's your Sharesies password? " + Purple)
	pwBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("failed to read input for password: %v", err)
	}
	fmt.Print(Reset + "\n") // terminal.ReadPassword captures the newline when we're done with the password
	pw := strings.TrimSpace(string(pwBytes))

	fmt.Print("What's your current MFA code? (if you don't have MFA enabled, just press enter) " + Purple)
	mfaCode, err := inputReader.ReadString('\n')
	if err != nil {
		log.Fatalf("failed to read input for MFA code: %v", err)
	}
	fmt.Print(Reset)
	mfaCode = strings.TrimSpace(mfaCode)

	fmt.Println()
	fmt.Println("Sweet! That should be all I need. I'll go fetch your portfolio history now, and find the maximum portfolio value for", desiredYear)
	fmt.Println()

	r, err := login.Login(login.LoginInput{
		Email:    email,
		Password: pw,
		MFAToken: mfaCode,
	})

	if err != nil {
		log.Fatalf("failed to login: %v", err)
	}

	pf, err := portfolio.Fetch(r.User.PortfolioID, r.RakaiaToken)
	if err != nil {
		log.Fatalf("failed to fetch portfolio: %v", err)
	}

	max := 0.0
	maxDate := portfolio.PortfolioDate{}

	for _, item := range pf.PortfolioHistory {
		t := time.Time(item.Date)
		if t.Year() != desiredYear {
			continue
		}
		if item.PortfolioValue > max {
			max = item.PortfolioValue
			maxDate = item.Date
		}
	}

	d := time.Time(maxDate).Format("2006-01-02")

	fmt.Printf("Max portfolio value: "+Green+"$%.2f %s"+Reset+" on %v\n", max, pf.Currency, d)
}
