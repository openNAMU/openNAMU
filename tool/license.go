package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type Module struct {
    Path    string
    Version string
    Dir     string
    Main    bool
    Replace *Module
}

func main() {
    // 1) 모듈 루트 탐색
    moduleRoot, err := findModuleRoot("..")
    must(err, "모듈 루트를 찾는 중 오류")
    fmt.Println("[info] module root:", moduleRoot)

    // 2) 출력 폴더: 모듈 루트/THIRD_PARTY_LICENSES
    licenseDir := filepath.Join(moduleRoot, "THIRD_PARTY_LICENSES")
    _ = os.MkdirAll(licenseDir, 0o755)

    // 3) 모듈 캐시 준비 (다운로드)
    must(runCmd(moduleRoot, "go", "mod", "download", "-json", "all"), "go mod download 실패")

    // 4) 모듈 목록 가져오기
    out, err := runCmdOut(moduleRoot, "go", "list", "-m", "-json", "all")
    must(err, "go list 실패")

    dec := json.NewDecoder(bytes.NewReader(out))
    var mods []Module
    for {
        var m Module
        if err := dec.Decode(&m); err != nil {
            if errors.Is(err, fs.ErrClosed) || errors.Is(err, os.ErrClosed) || err.Error() == "EOF" {
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
    for _, m := range mods {
        files := findLicenseFiles(m.Dir)
        if len(files) == 0 {
            // 부모에도 있는 경우가 있어 한 번 더 탐색
            files = findLicenseFiles(filepath.Dir(m.Dir))
        }
        if len(files) == 0 {
            missing = append(missing, fmt.Sprintf("%s@%s (Dir: %s)", m.Path, m.Version, m.Dir))
            continue
        }
        if err := saveModuleLicenses(licenseDir, m, files); err != nil {
            fmt.Fprintf(os.Stderr, "[warn] %s 저장 실패: %v\n", m.Path, err)
        }
    }

    // 6) 누락 목록 남기기
    if len(missing) > 0 {
        _ = os.WriteFile(filepath.Join(licenseDir, "_missing_licenses.txt"),
            []byte(strings.Join(missing, "\n")+"\n"), 0o644)
        fmt.Println("[warn] 라이선스 파일을 찾지 못한 모듈이 있습니다. _missing_licenses.txt 참조.")
    }

    fmt.Println("[done] THIRD_PARTY_LICENSES 폴더 생성 완료")
}

func findModuleRoot(start string) (string, error) {
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

func runCmd(dir string, name string, args ...string) error {
    cmd := exec.Command(name, args...)
    cmd.Dir = dir
    var stderr bytes.Buffer
    cmd.Stderr = &stderr
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("%v\nstderr:\n%s", err, stderr.String())
    }
    return nil
}

func runCmdOut(dir string, name string, args ...string) ([]byte, error) {
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

func findLicenseFiles(dir string) []string {
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
    for _, e := range entries {
        if e.IsDir() {
            continue
        }
        up := strings.ToUpper(e.Name())
        for _, c := range candidates {
            if up == strings.ToUpper(c) {
                res = append(res, filepath.Join(dir, e.Name()))
            }
        }
        // LICENSE-APACHE, NOTICE-MIT 등 접두 허용
        if strings.HasPrefix(up, "LICENSE") || strings.HasPrefix(up, "NOTICE") || strings.HasPrefix(up, "COPYING") {
            res = append(res, filepath.Join(dir, e.Name()))
        }
    }
    sort.Strings(res)
    return res
}

func saveModuleLicenses(licenseDir string, m Module, files []string) error {
    fileName := strings.ReplaceAll(m.Path, "/", "_") + "@" + m.Version + ".txt"
    outPath := filepath.Join(licenseDir, fileName)
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
    return os.WriteFile(outPath, b.Bytes(), 0o644)
}

func must(err error, ctx string) {
    if err != nil {
        panic(fmt.Errorf("%s: %w", ctx, err))
    }
}
