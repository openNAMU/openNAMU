package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type module struct {
    Path    string
    Version string
    Dir     string
    Main    bool
    Replace *module
}

func main() {
    // 1) 모듈 루트 탐색
    module_root, err := find_module_root(".")
    must(err, "모듈 루트를 찾는 중 오류")
    fmt.Println("[info] module root:", module_root)

    // 2) 출력 폴더: 모듈 루트/THIRD_PARTY_LICENSES
    license_dir := filepath.Join(module_root, "THIRD_PARTY_LICENSES")
    must(os.MkdirAll(license_dir, 0o755), "라이선스 폴더 생성 실패")

    // 3) 모듈 캐시 준비 (다운로드)
    must(run_cmd(module_root, "go", "mod", "download", "-json", "all"), "go mod download 실패")

    // 4) 모듈 목록 가져오기
    out, err := run_cmd_out(module_root, "go", "list", "-m", "-json", "all")
    must(err, "go list 실패")

    dec := json.NewDecoder(bytes.NewReader(out))
    var mods []module
    for {
        var m module
        if err := dec.Decode(&m); err != nil {
            if err == io.EOF {
                break
            }
            panic(fmt.Errorf("json decode 실패: %w", err))
        }
        if m.Replace != nil {
            m = *m.Replace
        }
        if m.Main || m.Dir == "" {
            continue
        }
        mods = append(mods, m)
    }

    // 정렬(안정적 출력)
    sort.Slice(mods, func(i, j int) bool {
        if mods[i].Path == mods[j].Path {
            return mods[i].Version < mods[j].Version
        }
        return mods[i].Path < mods[j].Path
    })

    // 5) 라이선스 수집
    var missing []string
    generated_files := make(map[string]struct{}, len(mods)+1)
    for _, m := range mods {
        files := find_license_files(m.Dir)
        if len(files) == 0 {
            // 부모에도 있는 경우가 있어 한 번 더 탐색
            files = find_license_files(filepath.Dir(m.Dir))
        }
        if len(files) == 0 {
            missing = append(missing, fmt.Sprintf("%s@%s (Dir: %s)", m.Path, m.Version, m.Dir))
            continue
        }
        if err := save_module_licenses(license_dir, m, files); err != nil {
            fmt.Fprintf(os.Stderr, "[warn] %s 저장 실패: %v\n", m.Path, err)
        } else {
            generated_files[module_license_file_name(m)] = struct{}{}
        }
    }

    // 6) 누락 목록 남기기
    if len(missing) > 0 {
        const missing_file_name = "_missing_licenses.txt"
        if err := os.WriteFile(filepath.Join(license_dir, missing_file_name),
            []byte(strings.Join(missing, "\n")+"\n"), 0o644); err != nil {
            fmt.Fprintf(os.Stderr, "[warn] 누락 목록 저장 실패: %v\n", err)
        } else {
            generated_files[missing_file_name] = struct{}{}
            fmt.Println("[warn] 라이선스 파일을 찾지 못한 모듈이 있습니다. _missing_licenses.txt 참조.")
        }
    }

    must(remove_stale_license_files(license_dir, generated_files), "이전 라이선스 파일 정리 실패")

    fmt.Println("[done] THIRD_PARTY_LICENSES 폴더 생성 완료")
}

func find_module_root(start string) (string, error) {
    dir, err := filepath.Abs(start)
    if err != nil {
        return "", err
    }
    for {
        if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
            return dir, nil
        }
        parent := filepath.Dir(dir)
        if parent == dir {
            break
        }
        dir = parent
    }
    return "", fmt.Errorf("go.mod를 찾을 수 없습니다")
}

func run_cmd(dir string, name string, args ...string) error {
    cmd := exec.Command(name, args...)
    cmd.Dir = dir
    var stderr bytes.Buffer
    cmd.Stderr = &stderr
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("%v\nstderr:\n%s", err, stderr.String())
    }
    return nil
}

func run_cmd_out(dir string, name string, args ...string) ([]byte, error) {
    cmd := exec.Command(name, args...)
    cmd.Dir = dir
    var stderr bytes.Buffer
    cmd.Stderr = &stderr
    out, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("%v\nstderr:\n%s", err, stderr.String())
    }
    return out, nil
}

func find_license_files(dir string) []string {
    candidates := []string{
        "LICENSE", "LICENSE.txt", "LICENSE.md",
        "COPYING", "COPYING.txt",
        "NOTICE", "NOTICE.txt",
        "UNLICENSE", "LICENCE", "LICENCE.txt",
    }
    entries, err := os.ReadDir(dir)
    if err != nil {
        return nil
    }
    var res []string
    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }

        upper_name := strings.ToUpper(entry.Name())
        license_file := false
        for _, candidate := range candidates {
            if upper_name == strings.ToUpper(candidate) {
                license_file = true
                break
            }
        }

        // LICENSE-APACHE, NOTICE-MIT 등 접두 허용
        if strings.HasPrefix(upper_name, "LICENSE") || strings.HasPrefix(upper_name, "NOTICE") || strings.HasPrefix(upper_name, "COPYING") {
            license_file = true
        }

        if license_file {
            res = append(res, filepath.Join(dir, entry.Name()))
        }
    }
    sort.Strings(res)
    return res
}

func save_module_licenses(license_dir string, m module, files []string) error {
    file_name := module_license_file_name(m)
    out_path := filepath.Join(license_dir, file_name)
    var b bytes.Buffer
    fmt.Fprintf(&b, "%s @ %s\n\n", m.Path, m.Version)
    for _, f := range files {
        data, err := os.ReadFile(f)
        if err != nil {
            fmt.Fprintf(&b, "[warn] read fail: %s: %v\n\n", f, err)
            continue
        }
        fmt.Fprintf(&b, "----- %s -----\n\n", filepath.Base(f))
        // 개행 정규화
        text := strings.ReplaceAll(string(data), "\r\n", "\n")
        b.WriteString(text)
        if !strings.HasSuffix(text, "\n") {
            b.WriteString("\n")
        }
        b.WriteString("\n")
    }
    return os.WriteFile(out_path, b.Bytes(), 0o644)
}

func module_license_file_name(m module) string {
    return strings.ReplaceAll(m.Path, "/", "_") + "@" + m.Version + ".txt"
}

func remove_stale_license_files(license_dir string, generated_files map[string]struct{}) error {
    entries, err := os.ReadDir(license_dir)
    if err != nil {
        return err
    }

    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }

        file_name := entry.Name()
        if _, exists := generated_files[file_name]; exists {
            continue
        }

        generated_license := strings.HasSuffix(file_name, ".txt") && strings.Contains(file_name, "@")
        if file_name != "_missing_licenses.txt" && !generated_license {
            continue
        }

        if err := os.Remove(filepath.Join(license_dir, file_name)); err != nil {
            return fmt.Errorf("%s 삭제 실패: %w", file_name, err)
        }
    }

    return nil
}

func must(err error, ctx string) {
    if err != nil {
        panic(fmt.Errorf("%s: %w", ctx, err))
    }
}
