//Command lsd lists only directories in the current directory
//or specified directories.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	showDot = flag.Bool(".", false, "show dot files")
	nulls   = flag.Bool("0", false, "print files separated by NUL instead of \\n")
	invert  = flag.Bool("v", false, "only show non-directories")
)

func opendir(path string) (*os.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", fi.Name())
	}
	return f, nil
}

func with(dir *os.File, f func(name string, isdir bool)) error {
	for {
		fis, err := dir.Readdir(100)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		for _, fi := range fis {
			name := fi.Name()
			if len(name) == 0 {
				continue
			}
			f(name, fi.IsDir())
		}
	}

	return nil
}

func warn(arg interface{}) {
	log.Println("lsd:", arg)
}

//Usage: %name %flags dirs*
func main() {
	log.SetFlags(0)
	flag.Parse()

	dirs := flag.Args()
	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	sep := "\n"
	if *nulls {
		sep = string([]rune{0})
	}

	for _, d := range dirs {
		dir, err := opendir(d)
		if err != nil {
			warn(err)
			dir.Close()
			continue
		}

		err = with(dir, func(name string, show bool) {
			if *invert {
				show = !show
			}
			if !show {
				return
			}

			if !*showDot && name[0] == '.' {
				return
			}

			fmt.Print(name, sep)
		})
		if err != nil {
			warn(err)
		}

		dir.Close()
	}
}
