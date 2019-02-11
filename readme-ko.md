openNAMU
====
[![Python 3.5 이상의 버전 필요](https://img.shields.io/badge/python-3.5%20or%20higher-blue.svg)](https://python.org)
[![라이선스](https://img.shields.io/badge/license-BSD%203--Clause-lightgrey.svg)](./LICENSE)

![](https://raw.githubusercontent.com/2du/openNAMU/master/.github/logo.png)

오픈나무는 파이썬 기반의 위키 엔진입니다. 파이썬과 그 의존성 모듈만 설치하면 사용할 수 있으며, 코드를 직접 수정하여 좀 더 주제에 특화된 위키를 만들 수 있습니다.

### 목차
 * [클론](#클론)
 * [기여](#기여)
 * [라이선스](#라이선스)
 * [기여자 목록](#기여자-목록)
 * [기타](#기타)

# 시작하기
오픈나무는 파이썬 환경에서 동작하는 파이썬 애플리케이션으로, 파이썬 환경을 필요로 합니다.

쉬운 오픈나무 설치를 위해 오픈나무 가이드를 따로 생성해두었으며, [이곳](https://github.com/Make-openNAMU/guide)에서 확인하실 수 있습니다.

가이드 목록
 * [파이썬 설치](https://github.com/Make-openNAMU/guide/blob/master/articles/ko-kr/install-python.md)
 * [오픈나무 시작](https://github.com/Make-openNAMU/guide/blob/master/articles/ko-kr/start-opennamu.md)

# 클론
아래 명령을 터미널(명령 프롬프트)에 입력하여 본 리포지토리를 클론할 수 있습니다.
## 일반
 * `git clone -b stable https://github.com/2du/openNAMU.git`

## 개발중
 * `git clone -b master https://github.com/2du/openNAMU.git`

# 기여
오픈나무에는 검증되지 않은 몇가지 버그가 존재할 수 있습니다. 당신의 오픈나무 사용과 버그 발견은 오픈나무의 발전을 돕습니다.
[이슈 생성하기](https://github.com/2du/openNAMU/issues/new)

오픈나무는 완전한 오픈소스 프로젝트입니다. 새로운 기능을 추가하고 Pull Request를 생성해보세요.
[Pull Request 생성하기](https://github.com/2du/openNAMU/compare)

# 라이선스
openNAMU 프로젝트는 [BSD 3-Clause License](./LICENSE)(이하 BSD-3 라이선스)의 보호를 받고 있으며, openNAMU 프로젝트를 사용하고자 한다면 BSD-3 라이선스를 준수해야 합니다. 본 라이선스를 위반할 경우 개발자는 DMCA Takedown 등 관련 제재를 관계자에게 요청할 권리가 있으며, 그 책임은 모두 라이선스 위반 사용자에게 있습니다. 자세한 내용은 문서를 참고하세요.

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
