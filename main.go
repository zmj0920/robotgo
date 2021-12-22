package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			fmt.Println("请在要打开的程序图标上点击右键")
			ok := robotgo.AddEvent("mright")
			if ok {
				var input int
				ButtonX, ButtonY := robotgo.GetMousePos()
				fmt.Println("坐标获取成功", ButtonX, ButtonY)
				fmt.Println("1.继续操作请输入  1  ")
				fmt.Println("2.退出当前程序请输入 0")
				fmt.Scanln(&input)
				switch input {
				case 1:
					files, fileUrlPath := getFilePath(ButtonX, ButtonY)
					for index, v := range files {
						fileUrl := fileUrlPath[index]
						open(v, ButtonX, ButtonY, fileUrl)
					}
				case 0:
					wg.Done()
				}
			}

		}
	}()
	wg.Wait()
}

func getFilePath(inputTextX int, inputTextY int) ([]string, []string) {
	var files []string
	var fileUrlPath []string
	pwd, _ := os.Getwd()
	filepath.Walk(filepath.Join(pwd, "/file"), func(path string, info os.FileInfo, err error) error {
		if info.Name() == "百度视频上传3.exe" {
			files = append(files, path)
			fileUrlPath = append(fileUrlPath, path[0:strings.LastIndex(path, `\`)])
		}
		return nil
	})
	sort.Strings(files)
	return files, fileUrlPath
}

func open(filePath string, ButtonX int, ButtonY int, fileUrl string) {
	cmd := exec.Command(filePath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		return
	}
	move(ButtonX, ButtonY)
	time.Sleep(time.Second * 2)
	robotgo.MoveClick(ButtonX, ButtonY, "left", true)
	TypeStrPath(fileUrl)
	saveBitmap()
	timer := time.NewTicker(time.Second * 3)
	<-timer.C
	move(ButtonX, ButtonY)
	time.Sleep(time.Second)
	pid := robotgo.GetPID()
	path, _ := robotgo.FindPath(pid)
	if filePath == path {
		robotgo.KeyTap("f4", "alt") //关闭窗口
	}
}

func saveBitmap() {
	width, height := robotgo.GetScreenSize()
	bitmap := robotgo.CaptureScreen(0, 0, width, height)
	defer robotgo.FreeBitmap(bitmap)
	img := filepath.Join("./img/", time.Now().Format("2006-01-02-15-04-05")+".png")
	robotgo.SaveBitmap(bitmap, img)
}

// 防止在点击的时候鼠标操作其他的
func move(x int, y int) {
	robotgo.Move(x, y)
	robotgo.Move(x, y)
	robotgo.Move(x, y)
}

func TypeStrPath(path string) {
	robotgo.KeyTap("tab")
	robotgo.KeyTap("tab")
	robotgo.KeyTap("tab")
	robotgo.KeyTap("tab")
	robotgo.KeyTap("tab")
	robotgo.KeyTap("tab")
	robotgo.KeyTap("tab")
	robotgo.KeyTap("tab")
	robotgo.KeyTap(`a`, `control`)
	robotgo.TypeStr("")
	time.Sleep(time.Second)
	robotgo.TypeStr(path)
}
