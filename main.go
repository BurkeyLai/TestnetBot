package main

import (
	"fmt"
	"strings"
	"time"

	//"os"
	//"strconv"
	//"strings"
	//"time"
	"github.com/go-ini/ini"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	port            = 8080
	webAddress      = "https://venom.network/tasks"
	getTokenAddress = "https://venom.network/faucet"
)

var (
	Cfg *ini.File
)

func UpdateWindowHandles(wd selenium.WebDriver, web_idx int) {

	// 等待一些时间，确保新窗口打开
	time.Sleep(2 * time.Second)

	// 获取当前所有窗口的句柄
	handles, err := wd.WindowHandles()
	if err != nil {
		fmt.Println("Can not get window handler!")
		panic(err)
	}

	// 切换到新窗口
	//if len(handles) >= 2 {
	err = wd.SwitchWindow(handles[web_idx])
	if err != nil {
		fmt.Println("Can not switch to new window!")
		panic(err)
	}
	//}
}

func main() {
	Cfg, err := ini.Load("app.ini")
	if err != nil {
		fmt.Println("Fail to parse 'app.ini': %v", err)
	}

	user := Cfg.Section("user")

	// 获取以逗号分隔的字符串，并使用 strings.Split 分割为切片
	keyValues := user.Key("SEED_PHRASE").String()
	seed_phrase := strings.Split(keyValues, ",")
	password := user.Key("PASSWORD").String()

	// 打印切片中的值
	//fmt.Println("切片中的值:")
	//for _, value := range seed_phrase {
	//	fmt.Println(value)
	//}
	//return

	opts := []selenium.ServiceOption{
		// Enable fake XWindow session.
		// selenium.StartFrameBuffer(),
		// selenium.Output(os.Stderr), // Output debug information to STDERR, comment it if you don't want too much console information
	}

	// Enable debug info.
	// selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService("chromedriver", port, opts...)
	if err != nil {
		panic(err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{
		"browserName": "chrome",
		//"chromeOptions": map[string]interface{}{
		//	"args": []string{
		//		"--log-level=3", // 設置日誌級別，3 表示只輸出錯誤訊息
		//	},
		//},
	}

	// Add extensions
	// https://stackoverflow.com/questions/34222412/load-chrome-extension-using-selenium
	c := chrome.Capabilities{
		Args: []string{
			"--ignore-certificate-errors",
			"--ignore-ssl-errors",
		},
	}

	c.AddExtension("/Users/bur/golang/robot2/chrome/ojggmchlghnjlapmfbnjholfjkiidbch.crx") // Venom Wallet Extension (need to be downloaded previously)
	caps[chrome.CapabilitiesKey] = c

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://127.0.0.1:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()

	if err := wd.Get(webAddress); err != nil {
		panic(err)
	}

	// seed phrase sign in
	UpdateWindowHandles(wd, 1)
	phrase_signin_btn, err := wd.FindElement(selenium.ByXPATH, `//*[@id="root"]/div[1]/div/div[2]/div/div/div[3]/div/div[2]/button/div`)
	if err != nil {
		panic(err)
	}
	if err := phrase_signin_btn.Click(); err != nil {
		panic(err)
	}

	// input phrases
	phrase_inputs, err := wd.FindElements(selenium.ByTagName, "input")
	if err != nil {
		fmt.Println("Can not get 'input' elements!")
		panic(err)
	}
	for i, input := range phrase_inputs {
		if i < len(seed_phrase) {
			// 清空输入框内容
			input.Clear()

			// 在输入框中输入预先准备好的字符串
			err := input.SendKeys(seed_phrase[i])
			if err != nil {
				fmt.Println("Can not input phrase!")
				panic(err)
			}
		} else {
			fmt.Println("Not enough phrases!")
			break
		}
	}
	confirm_signin_btn, err := wd.FindElement(selenium.ByXPATH, `//*[@id="confirm"]/div`)
	if err != nil {
		panic(err)
	}
	if err := confirm_signin_btn.Click(); err != nil {
		panic(err)
	}

	// input password
	password_inputs, err := wd.FindElements(selenium.ByTagName, "input")
	if err != nil {
		fmt.Println("Can not get 'input' elements!")
		panic(err)
	}
	for _, input := range password_inputs {
		input.Clear()
		err := input.SendKeys(password)
		if err != nil {
			fmt.Println("Can not input password!")
			panic(err)
		}
	}
	confirm_signin_btn, err = wd.FindElement(selenium.ByXPATH, `//*[@id="root"]/div[1]/div/div[2]/div/div[2]/button[1]`)
	if err != nil {
		panic(err)
	}
	if err := confirm_signin_btn.Click(); err != nil {
		panic(err)
	}

	// find "Connect Wallet" Button and click it
	UpdateWindowHandles(wd, 0)
	connect_wallet_btn, err := wd.FindElement(selenium.ByXPATH, `//*[@id="root"]/div[2]/div[1]/div[2]/div[2]/span`)
	if err != nil {
		panic(err)
	}
	if err := connect_wallet_btn.Click(); err != nil {
		panic(err)
	}

	time.Sleep(2 * time.Second)
	chrome_ext_btn, err := wd.FindElement(selenium.ByXPATH, `//*[@id="VENOM_CONNECT_MODAL_ID"]/div[1]/div/div[2]/div/div[2]/div[1]/div/div[1]/a/div/div/div[2]/div`)
	if err != nil {
		panic(err)
	}
	if err := chrome_ext_btn.Click(); err != nil {
		panic(err)
	}

	UpdateWindowHandles(wd, 1)
	time.Sleep(3 * time.Second)
	wallet_connect_btn, err := wd.FindElement(selenium.ByXPATH, `//*[@id="root"]/div/div/footer/div[2]/button/div`)
	if err != nil {
		panic(err)
	}
	if err := wallet_connect_btn.Click(); err != nil {
		panic(err)
	}

	UpdateWindowHandles(wd, 0)
	time.Sleep(5 * time.Second)
	if err := wd.Get(getTokenAddress); err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)
	undo_tasks, err := wd.FindElements(selenium.ByXPATH, "//div[@class='tasks__item task  ']")
	//undo_tasks, err := wd.FindElements(selenium.ByClassName, "tasks__item task  ")
	//undo_tasks, err := wd.FindElement(selenium.ByXPATH, `//*[@id="more"]/div[2]`)
	if err != nil {
		panic(err)
	}

	fmt.Println(undo_tasks)

	/*
		// Wait for dynamic change of the page
		// https://www.softwaretestingmaterial.com/dynamic-xpath-in-selenium/
		time.Sleep(5 * time.Second)
		block_msg_a, err := wd.FindElement(selenium.ByXPATH, `//*[@id="messages"]/div[2]/a`)

		if err != nil {
			panic(err)
		}

		//fmt.Print(block_msg_a.Text())

		str, err := block_msg_a.Text()
		if err != nil {
			panic(err)
		}
		if str == "開啟此連結" {
			block_msg_a.Click()

			//time_left_str, err := wd.FindElements(selenium.ByXPATH, `//*[@id="timeLeft"]`)
			//fmt.Println("////////////////////////////////////")
			//fmt.Println(time_left_str)
			//wd.ExecuteScript(`window.alert(location.href);`, nil)
			mainhandle, err := wd.CurrentWindowHandle()
			if err != nil {
				panic(err)
			}
			fmt.Println(mainhandle)

			//查看所有網頁的handle值
			handles, err := wd.WindowHandles()
			if err != nil {
				panic(err)
			}
			for _, handle := range handles {
				//fmt.Println(handle)
				wd.SwitchWindow(handle)
				url, _ := wd.CurrentURL()
				if strings.Contains(url, "wootalk.today/verify/") {
					break
				}
			}

			//handle, err := wd.CurrentWindowHandle()
			//if err != nil {
			//	panic(err)
			//}
			//fmt.Println(handle)
			//這一行是發送警報信息，寫這一行的目的，主要是看當前主窗口是哪一個
			//wd.ExecuteScript(`window.alert(location.href);`, nil)

			time_left, err := wd.FindElement(selenium.ByXPATH, `//*[@id="timeLeft"]`)
			fmt.Println("////////////////////////////////////")
			//fmt.Println(time_left.Text())
			time_left_str, err := time_left.Text()
			if err != nil {
				panic(err)
			}
			time_left_int, err := strconv.Atoi(time_left_str)

			time.Sleep(time.Duration(time_left_int+5) * time.Second)

			btn, err = wd.FindElement(selenium.ByXPATH, "/html/body/form/input[2]")
			if err != nil {
				panic(err)
			}
			fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			fmt.Print(btn)
			fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			//we, err := wd.FindElement(selenium.ByXPATH, `//*[@id="checkbox"]`)
			we, err := wd.FindElement(selenium.ByXPATH, "/html/body/form/div/iframe")
			if err != nil {
				panic(err)
			}

			fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			fmt.Println("we: ")
			fmt.Print(we)
			fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

			i, err := we.GetAttribute("data-hcaptcha-response")
			if err != nil {
				panic(err)
			}
			fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			fmt.Print("data-hcaptcha-response: " + i)
			fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

			for i == "" {
				i, err = we.GetAttribute("data-hcaptcha-response")
				if err != nil {
					panic(err)
				}
			}
			fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			fmt.Print("data-hcaptcha-response: " + i)
			fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			if err := btn.Click(); err != nil {
				panic(err)
			}
			fmt.Println("////////////////////////////////////")

			//if block_msg_a != nil {
			//	if block_msg_link, err := block_msg_a.GetAttribute("href"); err == nil {
			//		fmt.Println("////////////////////////////////////")
			//		fmt.Print(block_msg_link)
			//		fmt.Println("////////////////////////////////////")
			//
			//		//verifyAddress := []string{"https://wootalk.today/verify/"}
			//		//fmt.Println(strings.Join(verifyAddress, block_msg_link))
			//		//if err := wd.Get(strings.Join(verifyAddress, block_msg_link)); err != nil {
			//		//	panic(err)
			//		//}
			//	}
			//} else {
			//
			//}

		} else {
			fmt.Println("Here we go!")
		}
	*/
	time.Sleep(100000 * time.Second)
}
