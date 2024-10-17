package tool

import (
    "strings"

    "github.com/dlclark/regexp2"
)

// 모니마크 렌더러는 전체가 아니라 나무마크로 변환
func Monimark(data string) string {
    //
    r := regexp2.MustCompile(`<<((?:(?!<<|>>)))>>`, 0)
    for {
        if m, _ := r.FindStringMatch(data); m != nil {
            gps := m.Groups()
            m_string := m.String()

            data = strings.Replace(data, m_string, "[include("+gps[1].Captures[0].String()+")]", 1)
        } else {
            break
        }
    }

    return data
}
