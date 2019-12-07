package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func main() {
	var configPath string
	var force bool
	flag.StringVar(&configPath, "config", "config.yaml", "path to the configuration file")
	flag.BoolVar(&force, "force", false, "force overwriting all ")
	flag.Parse()

	cfg, err := readConfig(configPath)
	if err != nil {
		fail(err)
	}

	for src, dst := range cfg {
		if src, err = resolve(src); err != nil {
			fail(fmt.Errorf("unable to resolve path of %s: %w", src, err))
		}

		if dst, err = resolve(dst); err != nil {
			fail(fmt.Errorf("unable to resolve path of %s: %w", dst, err))
		}

		err := symlink(src, dst, force)
		if err != nil {
			if err == errExists {
				fmt.Println("Oops! I was trying to create a symlink, but the destination already exists!")
				fmt.Printf(" - from: %s\n", src)
				fmt.Printf(" - to: %s\n", dst)
				fmt.Printf("Do you want to remove its contents and create the symlink? [y/N]: ")
				reader := bufio.NewReader(os.Stdin)
				for {
					resp, err := reader.ReadString('\n')
					if err != nil {
						resp = ""
					}

					switch strings.ToLower(strings.TrimSpace(resp)) {
					case "y":
						if err := symlink(src, dst, true); err != nil {
							fail(err)
						}
						success("%s 🠆 %s", src, dst)
					case "n":
						info("%s 🠆 %s [Skipped]")
					default:
						fmt.Println("I could not understand that. Do you want to remove the contents? [y/N]: ")
						continue
					}

					break
				}
			} else {
				fail(err)
			}
		} else {
			success("%s 🠆 %s", src, dst)
		}
	}

	success("All done!")
}

func info(msg string, args ...interface{}) {
	fmt.Printf("[i] "+msg+"\n", args...)
}

func success(msg string, args ...interface{}) {
	fmt.Printf("[✓] "+msg+"\n", args...)
}

func fail(err error) {
	fmt.Printf("[✗] %s\n", err)
	os.Exit(1)
}

type config map[string]string

func readConfig(path string) (config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open config file at %s: %w", path, err)
	}

	var seen = make(map[string]struct{})
	var conf = make(config)

	for i, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid format at line %d of config file %s", i+1, path)
		}

		from := strings.TrimSpace(parts[0])
		to := strings.TrimSpace(parts[1])

		if _, ok := seen[to]; ok {
			return nil, fmt.Errorf("destination %s has been used more than once", to)
		}

		seen[to] = struct{}{}
		conf[from] = to
	}

	return conf, nil
}

var errExists = errors.New("destination already exists")

func symlink(src, dst string, force bool) error {
	_, err := os.Stat(dst)
	if err == nil {
		if !force {
			return errExists
		}

		if err := delete(dst); err != nil {
			return err
		}
	}

	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to open %s: %w", dst, err)
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("unable to create parent directories of file %s: %w", dst, err)
	}

	if err := os.Symlink(src, dst); err != nil {
		return fmt.Errorf("unable creating the symlink between %s and %s: %w", src, dst, err)
	}

	return nil
}

func delete(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot delete %s: %w", path, err)
	}

	if fi.IsDir() {
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("cannot delete folder %s: %w", path, err)
		}
	} else {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("cannot delete file %s: %w", path, err)
		}
	}
	return nil
}

func resolve(path string) (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("cannot get user: %w", err)
	}

	path, err = filepath.Abs(strings.Replace(path, "~", user.HomeDir, -1))
	if err != nil {
		return "", fmt.Errorf("cannot get absolute path: %w", err)
	}

	return filepath.Clean(path), nil
}
