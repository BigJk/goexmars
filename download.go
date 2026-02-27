package goexmars

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const githubReleasesLatestDownloadBase = "https://github.com/BigJk/goexmars/releases/latest/download"

// DownloadLib downloads the latest shared-library release asset for the current
// OS/arch and extracts it next to the executable.
//
// The release zip is expected to contain a top-level lib/ directory, resulting
// in <exe-dir>/lib/<platform-library>.
func DownloadLib() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve executable path: %w", err)
	}
	return downloadLibToDir(filepath.Dir(exePath), http.DefaultClient)
}

func downloadLibToDir(destDir string, client *http.Client) error {
	asset, err := releaseAssetName(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/%s", githubReleasesLatestDownloadBase, asset)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", "goexmars")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("download release asset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("download release asset: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read release asset: %w", err)
	}
	if err := unzipBytesToDir(data, destDir); err != nil {
		return err
	}

	libPath := filepath.Join(destDir, "lib", exmarsLibraryName())
	if _, err := os.Stat(libPath); err != nil {
		return fmt.Errorf("downloaded archive missing expected library %q: %w", libPath, err)
	}
	return nil
}

func releaseAssetName(goos, goarch string) (string, error) {
	switch goos {
	case "darwin":
		if goarch == "amd64" || goarch == "arm64" {
			return fmt.Sprintf("goexmars-%s-%s.zip", goos, goarch), nil
		}
	case "linux":
		if goarch == "amd64" {
			return fmt.Sprintf("goexmars-%s-%s.zip", goos, goarch), nil
		}
	case "windows":
		if goarch == "amd64" {
			return fmt.Sprintf("goexmars-%s-%s.zip", goos, goarch), nil
		}
	}
	return "", fmt.Errorf("unsupported platform for release asset: %s/%s", goos, goarch)
}

func unzipBytesToDir(zipData []byte, destDir string) error {
	readerAt := bytes.NewReader(zipData)
	zr, err := zip.NewReader(readerAt, int64(len(zipData)))
	if err != nil {
		return fmt.Errorf("open zip: %w", err)
	}

	cleanDest := filepath.Clean(destDir)
	for _, f := range zr.File {
		if err := extractZipFile(f, cleanDest); err != nil {
			return err
		}
	}
	return nil
}

func extractZipFile(f *zip.File, destDir string) error {
	name := filepath.Clean(filepath.FromSlash(f.Name))
	targetPath := filepath.Join(destDir, name)

	destPrefix := destDir + string(filepath.Separator)
	if targetPath != destDir && !strings.HasPrefix(targetPath, destPrefix) {
		return fmt.Errorf("invalid zip entry path: %q", f.Name)
	}

	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(targetPath, 0o755); err != nil {
			return fmt.Errorf("create directory %q: %w", targetPath, err)
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return fmt.Errorf("create parent directory for %q: %w", targetPath, err)
	}

	in, err := f.Open()
	if err != nil {
		return fmt.Errorf("open zip entry %q: %w", f.Name, err)
	}
	defer in.Close()

	mode := f.Mode()
	if mode == 0 {
		mode = 0o644
	}
	out, err := os.OpenFile(targetPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return fmt.Errorf("create output file %q: %w", targetPath, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("extract %q: %w", f.Name, err)
	}
	return nil
}
