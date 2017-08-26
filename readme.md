## 개요
오픈나무 정식 버전 입니다. 파이썬 Bottle 기반으로 돌아 갑니다.

## 설치법
app.py를 파이썬 3.x 버전으로 실행하고 완료되면 ../setup 에 들어가면 됩니다. (파이썬 3.6을 권장 합니다.)

첫 번째 가입자에게 소유자 권한이 부여됩니다.

## 기타
 * [테스트 서버](http://namu.ml/)
 
## 의존성
### 윈도우
 * pip install [bottle](https://pypi.python.org/pypi/bottle)
 * pip install [bottle-beaker](https://pypi.python.org/pypi/bottle-beaker)
 * pip install [tornado](https://pypi.python.org/pypi/tornado)
 * pip install [bcrypt](https://pypi.python.org/pypi/bcrypt)
### 우분투
 * pip3 install [bottle](https://pypi.python.org/pypi/bottle)
 * pip3 install [bottle-beaker](https://pypi.python.org/pypi/bottle-beaker)
 * pip3 install [tornado](https://pypi.python.org/pypi/tornado)
 * pip3 install [bcrypt](https://pypi.python.org/pypi/bcrypt)
### 자세한 설명
 * [참조](http://namu.ml/w/오픈나무%2F설치법)
 
## set.json 설명
 * db = 데이터베이스 이름
 * port = 위키 열 포트 (기본 : 3000)

## 업데이트 방법
새 오픈나무 버전의 파일 받고 오픈나무 폴더에 덮어 씌우고 app.py를 실행한 다음 소유자 계정으로 로그인 한 다음 /update 에 가서 업데이트 전 기존의 버전을 찾아 누르면 됩니다.

## MySQL 버전에서 업데이트
오픈나무를 받고 덮어 쓰기를 한 다음 update-mysql.py를 실행하고 완료하면 app.py를 키면 됩니다.
