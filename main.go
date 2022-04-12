package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func intMin(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func str2int(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func pressExit(s string, errno int) {
	fmt.Printf("%s\n按任意键退出！\n", s)
	os.Stdin.Read(make([]byte, 1))
	os.Exit(errno)
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func download(filename string, url string) error {
	fmt.Printf("Downloading %s from url: %s\n", filename, url)
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, err = io.Copy(out, res.Body)
	if err != nil {
		return err
	}
	fmt.Println("Download complete!")
	return err
}

func cmdRun(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	buf, err := cmd.Output()
	fmt.Println(string(buf))
	return err
}

func copyFile(srcFile, dstFile string) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}

type (
	browser struct {
		branch    string
		structure string
		version   string
		url       string
		filename  string
	}
	config struct {
		Branch    string `json:"Branch"`
		Structure string `json:"Structure"`
		Version   string `json:"Version"`
	}
)

func newBrowser() *browser {
	return &browser{
		branch:    "Stable",
		structure: "x64",
		version:   "0.0.0.0",
	}
}

func (b *browser) loadcfg() {
	cfgFile := "settings.json"
	if !pathExists(cfgFile) {
		pressExit(cfgFile+" 文件不存在！", 1)
	}
	content, err := ioutil.ReadFile(cfgFile)
	if err != nil || string(content) == "" {
		pressExit(cfgFile+" 文件加载失败！", 1)
	}
	var cfg config
	json.Unmarshal(content, &cfg)
	b.branch = cfg.Branch
	b.structure = cfg.Structure
	b.version = cfg.Version
}

func (b *browser) older(newVersion string) bool {
	v1 := strings.Split(b.version, ".")
	v2 := strings.Split(newVersion, ".")
	vparts := intMin(len(v1), len(v2))
	for i := 0; i < vparts; i++ {
		if str2int(v2[i]) > str2int(v1[i]) {
			return true
		}
		if str2int(v2[i]) < str2int(v1[i]) {
			return false
		}
	}
	return false
}

func (b *browser) terminate() error {
	return cmdRun("taskkill", "/f", "/im", "msedge.exe")
}

func (b *browser) download() error {
	return download(b.filename, b.url)
}

func (b *browser) extract() []error {
	var err []error
	fmt.Println("Extracting ...")
	err = append(err, cmdRun(".\\7za.exe", "x", b.filename, "-o.", "-aoa", "-y"))
	err = append(err, os.Remove(b.filename))
	err = append(err, cmdRun(".\\7za.exe", "x", "MSEDGE.7z", "-o.", "-aoa", "-y"))
	err = append(err, os.Remove("MSEDGE.7z"))
	cmdRun("md ..\\App\\")
	err = append(err, os.Mkdir("..\\App\\", 0777))
	err = append(err, os.Rename(".\\Chrome-bin\\"+b.version, "..\\App\\"+b.version))
	err = append(err, os.Remove(".\\Chrome-bin\\"))
	err = append(err, copyFile("..\\App\\"+b.version+"\\msedge.exe", "..\\App\\msedge.exe"))
	err = append(err, os.Mkdir("..\\Data\\", 0777))
	err = append(err, copyFile(".\\说明.txt", "..\\说明.txt"))
	err = append(err, copyFile(".\\msedge.exe", "..\\msedge.exe"))
	err = append(err, copyFile(".\\msedge.ini", "..\\msedge.ini"))
	fmt.Println("Extract complete!")
	return err
}

func (b *browser) edgepp() error {
	fmt.Printf("Copying edge++ ... ")
	err := copyFile(".\\version.dll", "..\\App\\version.dll")
	fmt.Println("Complete!")
	return err
}

func (b *browser) patch() []error {
	var err []error
	fmt.Println("Injecting version.dll to Edge ...")
	err = append(err, cmdRun(".\\setdll.exe", "/d:..\\App\\version.dll", "..\\App\\msedge.exe"))
	err = append(err, os.Remove("..\\App\\msedge.exe~"))
	fmt.Println("Complete!")
	return err
}

func (b *browser) replaced(ver string) {
	cfgFile := "settings.json"
	var cfg config
	cfg.Branch = b.branch
	cfg.Structure = b.structure
	cfg.Version = ver
	bytes, _ := json.Marshal(cfg)
	err := ioutil.WriteFile(cfgFile, bytes, 0777)
	if err != nil {
		pressExit(cfgFile+" 文件写入失败！", 1)
	}
}

func main() {
	current := newBrowser()
	new := newBrowser()
	current.loadcfg()
	new.branch = current.branch
	new.structure = current.structure
	fmt.Println("checking new version ...")
	updater := newEdgeUpdate()
	updateData := updater.getInfo(new.branch, new.structure)
	new.version = updateData.FileData.Version
	new.url = updateData.FileData.Url
	new.filename = updateData.FileData.FileId
	if current.older(new.version) {
		fmt.Printf("Branch: %s\tStructure: %s\n", current.branch, current.structure)
		fmt.Printf("A newer version found, %s -> %s\n", current.version, new.version)
		fmt.Printf("Please close Edge and press enter to continue.\n")
		fmt.Scanln()
		new.terminate()
		new.download()
		new.extract()
		new.edgepp()
		new.patch()
		current.replaced(new.version)
	} else {
		fmt.Println("No updates.")
	}
}
