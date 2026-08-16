package route

import (
	stdjson "encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const server_update_mode = "--opennamu-update"

func get_server_update_binary_name() (string, error) {
	switch runtime.GOOS {
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			return "main.amd64.bin", nil
		case "arm64":
			return "main.arm64.bin", nil
		}
	case "windows":
		switch runtime.GOARCH {
		case "amd64":
			return "main.amd64.exe", nil
		case "arm64":
			return "main.arm64.exe", nil
		}
	case "darwin":
		if runtime.GOARCH == "arm64" {
			return "main.mac.arm64.bin", nil
		}
	}

	return "", fmt.Errorf("unsupported platform: %s/%s", runtime.GOOS, runtime.GOARCH)
}

func get_server_update_url(branch string) (string, error) {
	if branch != "stable" && branch != "beta" {
		return "", fmt.Errorf("unsupported update branch: %s", branch)
	}

	version_url := "https://raw.githubusercontent.com/openNAMU/openNAMU/refs/heads/" + branch + "/version.json"
	request, err := http.NewRequest(http.MethodGet, version_url, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("User-Agent", "openNAMU")

	client := http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("version.json request failed: %s", response.Status)
	}

	version_data := map[string]any{}
	decoder := stdjson.NewDecoder(io.LimitReader(response.Body, 1<<20))
	if err := decoder.Decode(&version_data); err != nil {
		return "", err
	}

	release_tag, _ := version_data["r_ver"].(string)
	if branch == "beta" && release_tag == "" {
		if beta_data, ok := version_data["beta"].(map[string]any); ok {
			release_tag, _ = beta_data["r_ver"].(string)
		}
	}
	if release_tag == "" {
		return "", fmt.Errorf("r_ver is missing")
	}

	binary_name, err := get_server_update_binary_name()
	if err != nil {
		return "", err
	}

	return "https://github.com/openNAMU/openNAMU/releases/download/" + url.PathEscape(release_tag) + "/" + binary_name, nil
}

func download_server_update(binary_url string, temporary_path string) error {
	request, err := http.NewRequest(http.MethodGet, binary_url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("User-Agent", "openNAMU")

	client := http.Client{Timeout: 10 * time.Minute}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("binary download failed: %s", response.Status)
	}

	file_data, err := os.OpenFile(temporary_path, os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file_data.Close()

	size, err := io.Copy(file_data, response.Body)
	if err != nil {
		return err
	}
	if size == 0 {
		return fmt.Errorf("downloaded binary is empty")
	}

	if err := file_data.Chmod(0755); err != nil {
		return err
	}
	return file_data.Sync()
}

func start_server_update(branch string) error {
	executable, err := current_server_executable()
	if err != nil {
		return err
	}

	binary_url, err := get_server_update_url(branch)
	if err != nil {
		return err
	}

	temporary, err := os.CreateTemp(filepath.Dir(executable), ".opennamu-update-*")
	if err != nil {
		return err
	}
	temporary_path := temporary.Name()
	if err := temporary.Close(); err != nil {
		os.Remove(temporary_path)
		return err
	}

	if err := download_server_update(binary_url, temporary_path); err != nil {
		os.Remove(temporary_path)
		return err
	}

	working_dir, err := os.Getwd()
	if err != nil {
		os.Remove(temporary_path)
		return err
	}

	arguments := append([]string{server_update_mode, executable, temporary_path}, os.Args[1:]...)
	command := exec.Command(executable, arguments...)
	command.Dir = working_dir
	command.Env = server_process_environment()
	detach_server_process(command)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Start(); err != nil {
		os.Remove(temporary_path)
		return fmt.Errorf("start update process: %w", err)
	}
	return nil
}

func Run_server_update(arguments []string) int {
	if len(arguments) < 2 {
		log.Printf("server update failed: invalid arguments")
		return 1
	}

	executable := arguments[0]
	temporary_path := arguments[1]
	server_arguments := arguments[2:]

	for attempt := 0; attempt < 40; attempt++ {
		if err := os.Rename(temporary_path, executable); err == nil {
			if err := os.Chmod(executable, 0755); err != nil {
				log.Printf("server update failed: %v", err)
				return 1
			}

			working_dir, err := os.Getwd()
			if err != nil {
				log.Printf("server update failed: %v", err)
				return 1
			}

			command := exec.Command(executable, server_arguments...)
			command.Dir = working_dir
			command.Env = server_process_environment()
			detach_server_process(command)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			if err := command.Start(); err != nil {
				log.Printf("server update failed: %v", err)
				return 1
			}
			return 0
		}

		time.Sleep(250 * time.Millisecond)
	}

	os.Remove(temporary_path)
	log.Printf("server update failed: replacing binary timed out")
	return 1
}
