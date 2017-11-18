## 개요
오픈나무 정식 버전 입니다. 파이썬 Bottle 기반으로 돌아 갑니다.

## 설치법
밑의 의존성을 설치하고 app.py를 파이썬 3.5 이상 버전으로 실행하세요.

첫 번째 가입자에게 소유자 권한이 부여됩니다.

## 기타
 * [테스트 서버](http://namu.ml/)
 
## 의존성
 * [파이썬](https://www.python.org/downloads/) 3.5 이상
 
### 윈도우
 * pip install -r requirements.txt
 
cmd에 치면 됩니다.
### 우분투
 * sudo apt-get install python3-pip
 * pip3 install -r requirements.txt
 
터미널에 치면 됩니다.
### 자세한 설명
 * [참조](http://namu.ml/w/오픈나무%2F설치법)
 
## set.json 설명
 * db = 데이터베이스 이름
 * port = 위키 열 포트 (기본 : 3000)

set.json를 삭제하면 다시 새로 만들 수 있습니다.

## 업데이트 방법
새 오픈나무 버전의 파일 받고 오픈나무 폴더에 덮어 씌우고 app.py를 실행한 다음 소유자 계정으로 로그인 한 다음 /update 에 가서 업데이트 전 기존의 버전을 찾아 누르면 됩니다.

### MySQL 버전에서 업데이트
오픈나무를 받고 덮어 쓰기를 한 다음 update-mysql.py를 실행하고 완료하면 app.py를 키면 됩니다.

#### 으악 app.py 먼저 했어요
.db로 시작하는 파일을 지우고 다시 설명대로 하시면 됩니다.