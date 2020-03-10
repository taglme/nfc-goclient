package main

import (
	"fmt"
)

func main() {

	tasker := NewTasker("127.0.0.1:3011")
	tasker.Add(TaskWriteURL("https://tagl.me"))
	tasker.Add(TaskRead())
	tasker.Add(TaskGetDump())
	tasker.Add(TaskTransmit([]byte{0x60}))
	tasker.Add(TaskFormatDefault())
	tasker.Add(TaskSetPassword([]byte{0x11, 0x11, 0x11, 0x11}))
	tasker.Add(TaskRemovePassword([]byte{0x11, 0x11, 0x11, 0x11}))

	err := tasker.Run(1)
	if err != nil {
		fmt.Printf("Tasks execution failed with error: %s", err)
	}
	return

}
