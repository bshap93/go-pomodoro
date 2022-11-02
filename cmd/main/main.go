package main

import (
	"fmt"
)

type profile struct {
	intervalLen  int
	breakLen     int
	intervalNum  int
	longBreakLen int
}

func main() {
	fmt.Println("Pomodoro App")
	var myProfile profile = profile{}
	defaultProfile := profile{intervalLen: 25, breakLen: 5, intervalNum: 4, longBreakLen: 15}
	var ans string
	for ans != "y" && ans != "Y" && ans != "n" && ans != "N" {
		fmt.Println("Accept the default? [y/n]")

		// var ans string
		fmt.Scan(&ans)

		fmt.Println(len(ans))

		switch ans {
		case "y", "Y":
			myProfile = defaultProfile
			break
		case "n", "N":
			myProfile.setup()
			break
		}

	}

	fmt.Println(myProfile)

}

func (p profile) setup() {
	fmt.Println("Put your preferred interval: ")
	intervalLength, err := fmt.Scan()
	if err != nil {
		fmt.Println("invalid interval:", err)
	}

	fmt.Println(intervalLength)

}
