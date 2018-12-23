openNAMU
====
![Python 3.5 이상 필요](https://img.shields.io/badge/python-%3E%3D%203.5-blue.svg)

![](./logo.png)

오픈나무는 파이썬 기반의 위키 엔진입니다. 파이썬과 그 의존성 모듈만 설치하면 사용할 수 있으며, 코드를 직접 수정하여 좀 더 주제에 특화된 위키를 만들 수 있습니다.

### 목차
 * [시작하기](#시작하기)
 * [클론](#클론)
 * [기여](#기여)
 * [라이선스](#라이선스)
 * [기여자 목록](#기여자-목록)
 * [기타](#기타)

# 시작하기
오픈나무는 파이썬 환경에서 동작하는 파이썬 애플리케이션으로, 파이썬 환경을 필요로 합니다.

## 환경 구성
### 파이썬 설치
[파이썬 설치 가이드](https://github.com/404-sdok/how-to-python/blob/master/0.md)를 참고하여 파이썬을 설치합니다.

### 릴리즈 다운로드
[릴리즈](https://github.com/2du/openNAMU/releases)에서 openNAMU의 릴리즈 판을 다운로드 받고, 압축을 해제합니다. 레포지토리를 클론하여 릴리즈를 적용하는 것도 가능합니다. 자세한 내용은 [클론](#클론)문단을 참고하세요.

### 모듈 설치
다음 명령어로 openNAMU 구성 파일이 위치한 디렉토리로 이동합니다.
```
cd [path]
```


다음 명령어로 openNAMU 실행에 필요한 모듈을 설치합니다.
```
pip install -r requirements.txt
```
리눅스 환경의 경우 다음 명령으로 실행해야 합니다.
```
pip3 install -r requirements.txt
```
## 애플리케이션 시작
### 애플리케이션 시작
openNAMU를 시작합니다.
```
python app.py
```
리눅스 환경의 경우 다음 명령으로 실행해야 합니다.
```
python3 app.py
```

* openNAMU의 첫 계정은 소유자 계정으로 설정됩니다.

### 애플리케이션 공개
현재 openNAMU에서 권장하는 애플리케이션 공개 방법은 아파치나 nginx와 같은 웹 서버 구동 소프트웨어를 통해 리버스 프록시를 설정하는 것입니다.

* ([apache 설정 파일 예시](./conf_example/apache/), [nginx 설정 파일 예시](./conf_example/nginx.conf))

만약 HTTPS 리버스 프록시를 설정하는 경우, 보안을 위해 openNAMU 설정에서 호스트를 localhost 환경으로 변경할 것을 권장합니다.


# 클론
아래 명령을 터미널(명령 프롬프트)에 입력하여 본 리포지토리를 클론할 수 있습니다.
## 일반
 * `git clone -b stable https://github.com/2du/openNAMU.git`

## 베타
 * `git clone -b master https://github.com/2du/openNAMU.git`

# 기여
오픈나무에는 검증되지 않은 몇가지 버그가 존재할 수 있습니다. 당신의 오픈나무 사용과 버그 발견은 오픈나무의 발전을 돕습니다.
[이슈 생성하기](https://github.com/2du/openNAMU/issues/new)

오픈나무는 완전한 오픈소스 프로젝트입니다. 새로운 기능을 추가하고 Pull Request를 생성해보세요.
[Pull Request 생성하기](https://github.com/2du/openNAMU/compare)

# 라이선스
오픈나무는 [BSD 3-Clause License](./LICENSE)에 의해 보호받고 있습니다. 자세한 내용은 문서를 참고하세요.

## 포함된 외부 프로젝트
 * Quotes icon - [Dave Gandy](http://www.flaticon.com/free-icon/quote-left_25672) - CC 3.0 BY
 * Syntax highlighting - [highlightjs](https://highlightjs.org/) - [BSD License](https://github.com/highlightjs/highlight.js/blob/master/LICENSE)
 * Numerical expression - [MathJax](https://www.mathjax.org/) - [Apache License 2.0](https://github.com/mathjax/MathJax/blob/master/LICENSE)

# 기여자 목록
 * [참고](https://github.com/2DU/openNAMU/graphs/contributors)

## 도움을 주신 분들
 * [Team Croatia](https://github.com/TeamCroatia)
 * Basix
 * Efrit
 * Other chat rooms

# 기타
`set.json`은 몇가지 로컬 설정을 저장하는 설정 파일입니다.
 * [filename].db = 데이터베이스 이름

`set.json`은 삭제해도 다시 새로 만들 수 있습니다.

[테스트 서버](http://namu.ml/)

첫 번째 가입자에게 소유자 권한이 부여됩니다.
