package main

import (
	"bufio"
	"fmt"
	"notification_center/logic"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	logic.StartReceivingMsgs()

	fmt.Println("Press Enter to exit")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}
