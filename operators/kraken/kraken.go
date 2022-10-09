package kraken

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/joho/godotenv"
)

var api = krakenapi.New(GoDotEnvVariable("KRAKEN_API_KEY"), GoDotEnvVariable("KRAKEN_API_SECRET"))

func GetBalance() (string, error) {
	result, err := api.Query("Balance", map[string]string{})
	if err != nil {
		log.Println("Unexpected error fetching Kraken balance", err)
		return "", err
	}
	res := result.(map[string]interface{})
	return fmt.Sprint(res["XXBT"]), nil
}

func Withdraw(amount string) (interface{}, error) {
	result, err := api.Query("Withdraw", map[string]string{
		"asset":  "xbt",
		"key":    GoDotEnvVariable("KRAKEN_WITHDRAW_ADDRESS_KEY"),
		"amount": amount,
	})
	if err != nil {
		log.Println("Unexpected error performing Kraken withdrawal", err)
		return nil, err
	}
	return result, nil
}

// Receives an amount defined in BTC, returns an invoice
// NOTE: The first time you run this, you need
func GetAddress(amount string) (invoice string) {
	// create chrome instance
	ctx, cancel := chromedp.NewExecAllocator(
		context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("incognito", false),
			chromedp.Flag("headless", false), // Sorry, doesn't work headless
			chromedp.UserDataDir(GoDotEnvVariable("CHROME_PROFILE_PATH")))...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(
		ctx,
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 300*time.Second)
	defer cancel()

	// navigate to a page, wait for an element, click
	err := chromedp.Run(ctx,
		browser.SetPermission(&browser.PermissionDescriptor{Name: "clipboard-read"}, browser.PermissionSettingGranted).WithOrigin("https://www.kraken.com"),

		// YOU NEED TO UNCOMMENT THE LOGIN LOGIC THE FIRST TIME YOU RUN THIS AND CONFIRM DEVICE VIA EMAIL
		/*
			chromedp.Navigate(`https://www.kraken.com/u/funding/deposit?asset=BTC&method=1`), // will redirect us to login page
			// wait for footer element is visible (ie, page is loaded)
			chromedp.WaitVisible(`//input[@name="username"]`),
			chromedp.SendKeys(`//input[@name="username"]`, GoDotEnvVariable("KRAKEN_USERNAME")),
			chromedp.WaitVisible(`//input[@name="password"]`),
			chromedp.SendKeys(`//input[@name="password"]`, GoDotEnvVariable("KRAKEN_PASSWORD")),
			chromedp.Sleep(3*time.Second),
			chromedp.SendKeys(`//input[@name="password"]`, kb.Enter),
			// find and click body > reach-portal:nth-child(37) > div:nth-child(3) > div > div > div > div > div.tr.mt3 > button.Button_button__caA8R.Button_primary__c5lrD.Button_large__T4YrY.no-tab-highlight
			chromedp.Sleep(3*time.Second),
			chromedp.SendKeys(`//input[@name="tfa"]`, GoDotEnvVariable("KRAKEN_OTP_SECRET"))),
			chromedp.Sleep(1*time.Second),
			chromedp.SendKeys(`//input[@name="tfa"]`, kb.Enter),
			chromedp.Sleep(30*time.Second),
			// GO CONFIRM YOUR DEVICE VIA EMAIL, COMMENT THIS OUT AGAIN AND RESTART SCRIPT
		*/

		chromedp.Navigate(`https://www.kraken.com/u/funding/deposit?asset=BTC&method=1`),
		chromedp.Sleep(10*time.Second),
		chromedp.Click(`div:nth-child(3) > div > div > div > div > div.tr.mt3 > button.Button_button__caA8R.Button_primary__c5lrD.Button_large__T4YrY.no-tab-highlight`, chromedp.ByQueryAll),
		chromedp.Sleep(2*time.Second),
		chromedp.SendKeys(`//input[@name="amount"]`, amount),
		chromedp.Sleep(1*time.Second),
		chromedp.Click(`#__next > div > main > div > div.container > div > div.FundingTransactionPage_form__OGaKV > div > div > div:nth-child(4) > div.LightningForm_callToAction__Y4b1E > button`, chromedp.NodeVisible),
		chromedp.Sleep(1*time.Second),
		chromedp.Click(`#__next > div > main > div > div.container > div > div.FundingTransactionPage_form__OGaKV > div > div > div:nth-child(4) > div:nth-child(5) > div > div > button > div`, chromedp.ByQuery),
		chromedp.Evaluate(`window.navigator.clipboard.readText()`, &invoice, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithAwaitPromise(true)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	return invoice
}

// use godot package to load/read the .env file and
// return the value of the key
func GoDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
