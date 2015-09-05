package common

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/eris-ltd/eth-client/Godeps/_workspace/src/github.com/eris-ltd/common/go/log" // so we can flush logs on exit/ifexit
)

var (
	// Convenience Directories
	GoPath   = os.Getenv("GOPATH")
	ErisLtd  = path.Join(GoPath, "src", "github.com", "eris-ltd")
	ErisGH   = "https://github.com/eris-ltd/"
	usr, _   = user.Current() // error?!
	ErisRoot = ResolveErisRoot()

	// Major Directories
	ActionsPath        = path.Join(ErisRoot, "actions")
	BlockchainsPath    = path.Join(ErisRoot, "blockchains")
	DataContainersPath = path.Join(ErisRoot, "data")
	DappsPath          = path.Join(ErisRoot, "dapps")
	FilesPath          = path.Join(ErisRoot, "files")
	KeysPath           = path.Join(ErisRoot, "keys")
	LanguagesPath      = path.Join(ErisRoot, "languages")
	ServicesPath       = path.Join(ErisRoot, "services")
	ScratchPath        = path.Join(ErisRoot, "scratch")

	// Keys
	KeysDataPath = path.Join(KeysPath, "data")
	KeyNamesPath = path.Join(KeysPath, "names")

	// Scratch Directories (globally coordinated)
	EpmScratchPath  = path.Join(ScratchPath, "epm")
	LllcScratchPath = path.Join(ScratchPath, "lllc")
	SolcScratchPath = path.Join(ScratchPath, "sol")
	SerpScratchPath = path.Join(ScratchPath, "ser")

	// Blockchains stuff
	ChainsConfigPath = path.Join(BlockchainsPath, "config")
	HEAD             = path.Join(BlockchainsPath, "HEAD")
	Refs             = path.Join(BlockchainsPath, "refs")
)

var MajorDirs = []string{
	ErisRoot, ActionsPath, BlockchainsPath, DataContainersPath, DappsPath, FilesPath, KeysPath, LanguagesPath, ServicesPath, KeysDataPath, KeyNamesPath, ScratchPath, EpmScratchPath, LllcScratchPath, SolcScratchPath, SerpScratchPath, ChainsConfigPath,
}

//---------------------------------------------
// user and process

func Usr() string {
	u, _ := user.Current()
	return u.HomeDir
}

func Exit(err error) {
	status := 0
	if err != nil {
		log.Flush()
		fmt.Println(err)
		status = 1
	}
	os.Exit(status)
}

func IfExit(err error) {
	if err != nil {
		log.Flush()
		fmt.Println(err)
		os.Exit(1)
	}
}

// user and process
//---------------------------------------------------------------------------
// filesystem

func AbsolutePath(Datadir string, filename string) string {
	if path.IsAbs(filename) {
		return filename
	}
	return path.Join(Datadir, filename)
}

func InitDataDir(Datadir string) error {
	if _, err := os.Stat(Datadir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(Datadir, 0777); err != nil {
				return err
			}
		}
	}
	return nil
}

func ResolveErisRoot() string {
	var eris string
	if os.Getenv("ERIS") != "" {
		eris = os.Getenv("ERIS")
	} else {
		if runtime.GOOS == "windows" {
			home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
			if home == "" {
				home = os.Getenv("USERPROFILE")
			}
			eris = path.Join(home, ".eris")
		} else {
			eris = path.Join(Usr(), ".eris")
		}
	}
	return eris
}

// Create the default eris tree
func InitErisDir() (err error) {
	for _, d := range MajorDirs {
		err := InitDataDir(d)
		if err != nil {
			return err
		}
	}
	if _, err = os.Stat(HEAD); err != nil {
		_, err = os.Create(HEAD)
	}
	return
}

func ClearDir(dir string) error {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range fs {
		n := f.Name()
		if f.IsDir() {
			if err := os.RemoveAll(path.Join(dir, f.Name())); err != nil {
				return err
			}
		} else {
			if err := os.Remove(path.Join(dir, n)); err != nil {
				return err
			}
		}
	}
	return nil
}

func Copy(src, dst string) error {
	f, err := os.Stat(src)
	if err != nil {
		return err
	}
	if f.IsDir() {
		return copyDir(src, dst)
	}
	return copyFile(src, dst)
}

// assumes we've done our checking
func copyDir(src, dst string) error {
	fi, err := os.Stat(src)
	if err := os.MkdirAll(dst, fi.Mode()); err != nil {
		return err
	}
	fs, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range fs {
		s := path.Join(src, f.Name())
		d := path.Join(dst, f.Name())
		if f.IsDir() {
			if err := copyDir(s, d); err != nil {
				return err
			}
		} else {
			if err := copyFile(s, d); err != nil {
				return err
			}
		}
	}
	return nil
}

// common golang, really?
func copyFile(src, dst string) error {
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}
	return nil
}

func WriteFile(data, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0775); err != nil {
		return err
	}
	writer, err := os.Create(filepath.Join(path))
	defer writer.Close()
	if err != nil {
		return err
	}
	writer.Write([]byte(data))
	return nil
}

// filesystem
//-------------------------------------------------------
// open text editors

func Editor(file string) error {
	editr := os.Getenv("EDITOR")
	if strings.Contains(editr, "/") {
		editr = path.Base(editr)
	}
	switch editr {
	case "", "vim", "vi":
		return vi(file)
	case "emacs":
		return emacs(file)
	}
	return fmt.Errorf("Unknown editor %s", editr)
}

func emacs(file string) error {
	cmd := exec.Command("emacs", file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func vi(file string) error {
	cmd := exec.Command("vim", file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
