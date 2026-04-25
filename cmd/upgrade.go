package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

var version = "dev"
var program = "summon"
var repo = fmt.Sprintf("tofunmiadewuyi/%s", program)

func replaceBinary(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// Write to a temp file in the same directory as dest, then rename into
	// place. rename() swaps directory entries atomically without opening the
	// running executable for writing, avoiding "text file busy" on Linux.
	tmp := dest + ".new"
	out, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		os.Remove(tmp)
		return err
	}
	out.Close()

	return os.Rename(tmp, dest)
}

type Release struct {
	TagName string `json:"tag_name"`
}

func upgrade() {
	// Get latest release
	resp, err := http.Get("https://api.github.com/repos/" + repo + "/releases/latest")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var rel Release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		fmt.Println("Upgrade failed: could not parse release info:", err)
		return
	}
	if rel.TagName == "" {
		fmt.Println("Upgrade failed: no release found (check https://github.com/" + repo + "/releases)")
		return
	}

	if rel.TagName == version {
		fmt.Println("Already up to date.")
		return
	}

	fmt.Printf("Update available: %s (current: %s)\n", rel.TagName, version)
	fmt.Print("Continue? (y/n): ")

	var input string
	fmt.Scanln(&input)

	if input != "y" && input != "Y" {
		fmt.Println("Aborted.")
		return
	}

	fmt.Println("Upgrading to", rel.TagName)

	osName := runtime.GOOS
	arch := runtime.GOARCH

	filename := fmt.Sprintf("%s_%s_%s_%s.zip", program, rel.TagName, osName, arch)
	url := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", repo, rel.TagName, filename)

	tmpFile := fmt.Sprintf("/tmp/%s.zip", program)

	out, err := os.Create(tmpFile)
	if err != nil {
		fmt.Println("Upgrade failed:", err)
		return
	}
	defer out.Close()

	resp, err = http.Get(url)
	if err != nil {
		fmt.Println("Upgrade failed: download error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Upgrade failed: could not download release (HTTP %d)\n", resp.StatusCode)
		return
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		fmt.Println("Upgrade failed: download error:", err)
		return
	}
	out.Close()

	// Extract
	zr, err := zip.OpenReader(tmpFile)
	if err != nil {
		fmt.Println("Upgrade failed: could not read archive:", err)
		return
	}
	defer zr.Close()

	var binPath = fmt.Sprintf("/tmp/%s_new", program)
	found := false

	for _, f := range zr.File {
		if f.Name == program {
			rc, err := f.Open()
			if err != nil {
				fmt.Println("Upgrade failed:", err)
				return
			}
			outFile, err := os.Create(binPath)
			if err != nil {
				rc.Close()
				fmt.Println("Upgrade failed:", err)
				return
			}
			io.Copy(outFile, rc)
			outFile.Close()
			rc.Close()
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Upgrade failed: binary not found in archive")
		return
	}

	os.Chmod(binPath, 0755)

	current, _ := os.Executable()

	err = replaceBinary(binPath, current)
	if err != nil {
		fallback := filepath.Join(os.Getenv("HOME"), fmt.Sprintf(".local/bin/%s"), program)
		os.MkdirAll(filepath.Dir(fallback), 0755)
		err2 := replaceBinary(binPath, fallback)
		if err2 != nil {
			fmt.Println("Upgrade failed:", err2)
			return
		}
		fmt.Printf("Installed to %s (original location was not writable).\n", fallback)
		fmt.Println("Ensure ~/.local/bin is in your $PATH.")
		return
	}

	fmt.Println("Upgrade complete.")
}
