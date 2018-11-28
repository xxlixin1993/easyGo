package gracefulExit

import (
	"container/list"
	"errors"
	"strings"
)

var exitList *ExitList

type ExitInterface interface {
	// 获取退出程序名
	GetModuleName() string

	// 退出时需要执行的函数
	Stop() error
}

type ExitList struct {
	// 退出链表
	ll *list.List

	// 退出程序名称
	module map[string]*list.Element
}

// 初始化退出链表
func InitExitList() {
	exitList = &ExitList{
		ll:     list.New(),
		module: make(map[string]*list.Element),
	}
}

// 获取一个退出链表实例
func GetExitList() *ExitList {
	if exitList == nil {
		return nil
	}
	return exitList
}

// 在最前面插入一个退出事件
func (el *ExitList) UnShift(exitInterface ExitInterface) error {
	if el.module == nil {
		return errors.New("[gracefulExit] plz init ExitList first")
	}

	// Judge whether it exists or not
	moduleName := exitInterface.GetModuleName()
	if _, ok := el.module[moduleName]; ok {
		return errors.New("[gracefulExit] this module(" + moduleName + ") name is exist")
	}

	// Add value
	element := el.ll.PushFront(exitInterface)
	el.module[moduleName] = element

	return nil
}

// 在最后面插入一个退出事件
func (el *ExitList) Push(exitInterface ExitInterface) error {
	if el.module == nil {
		return errors.New("[gracefulExit] plz init ExitList first")
	}

	// Judge whether it exists or not
	moduleName := exitInterface.GetModuleName()
	if _, ok := el.module[moduleName]; ok {
		return errors.New("[gracefulExit] this module(" + moduleName + ") name is exist")
	}

	// Add value
	element := el.ll.PushBack(exitInterface)
	el.module[moduleName] = element

	return nil
}

// 平滑退出
func (el *ExitList) Stop() error {
	length := el.ll.Len()
	if length == 0 {
		return nil
	}

	errInfo := make([]string, 0)
	for i := 0; i < length; i++ {
		element := el.ll.Front()
		exitElement := element.Value.(ExitInterface)

		if err := exitElement.Stop(); err != nil {
			errInfo = append(errInfo, "[gracefulExit]: Stop this module("+exitElement.GetModuleName()+")"+err.Error())
		}

		el.ll.Remove(element)
	}

	if len(errInfo) > 0 {
		return errors.New(strings.Join(errInfo, "\n"))
	}

	return nil
}
