## 개요
오픈나무 정식 버전 입니다. 파이썬 플라스크 기반으로 돌아 갑니다.

## 설치법
set 폴더에 있는 set.json을 폴더 밖으로 꺼내고 json 내용을 수정하고 start.py를 파이썬 3.x 버전으로 실행하면 됩니다. (파이썬 3.6을 권장 합니다.)

등등등등

## 기타
 * [테스트 서버](https://namu.ml/)
 
## 의존성
### 기본
 * [MariaDB](https://mariadb.org/)나 [MySQL](https://www.mysql.com/)
### 윈도우
 * pip install [flask](https://pypi.python.org/pypi/Flask/0.12)
 * pip install [bcrypt](https://pypi.python.org/pypi/bcrypt/3.1.0)
 * pip install [pymysql](https://pypi.python.org/pypi/PyMySQL)
### 우분투
 * pip3 install [flask](https://pypi.python.org/pypi/Flask/0.12)
 * pip3 install [bcrypt](https://pypi.python.org/pypi/bcrypt/3.1.0)
 * pip3 install [pymysql](https://pypi.python.org/pypi/PyMySQL)
### 자세한 설명
 * [참조](http://namu.ml/w/오픈나무%2F설치법)
 
## set.json 설명
 * db = 데이터베이스 이름
 * host = 데이터베이스 호스트 (기본 : localhost)
 * user = 데이터베이스 사용자명 (기본 : root)
 * pw = 그 사용자의 비밀번호
 * owner = 소유자의 위키 내에서 사용 할 이름
 * name = 위키 이름
 * frontpage = 위키 대문
 * license = 하단에 표기 될 라이선스
 * key = 세션 키
 * upload = 업로드 제한 (메가 바이트 단위)
 * port = 위키 열 포트 (기본 : 3000)
 * help = 편집시 옆에 보여 줄 문법 도움말 문서 (공백이면 아무 문서도 보여주지 않음)
 
## 설치 설명
의존성 유틸들을 설치 합니다.

그 다음 오픈나무 set.json을 폴더 밖에 꺼낸 후 설명에 따라 커스텀 합니다.

오픈나무를 키고 회원 가입으로 가서 owner 이름으로 설정 한 이름으로 가입 합니다.

## 기타
만약 위키가 자꾸 터질시 MySQL이나 MariaDB의 my.ini나 my.cnf에 wait_timeout 값을 아주 큰 값으로 해주면 위키가 안 터집니다.

lower_case_table_names 값을 0 으로하면 대소문자 구분하고 1로 하면 대소문자를 구분하지 않습니다. (단 버그 존재 가능)
