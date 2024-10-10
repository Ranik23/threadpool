package task

import (
	"fmt"
)


type Task interface {
	Run()
}


type TTask struct {
	Id int
}


func (t TTask) Run() {
	fmt.Println("running", t.Id)
}
