## 개요
Ownet을 파이썬으로 재 구현한 버전 입니다.

## 설치법
set 폴더에 있는 set.json을 폴더 밖으로 꺼내고 json 내용을 수정하고 app.py를 파이썬 3.x 버전으로 실행하면 됩니다. (파이썬 3.6을 권장 합니다.)

## 기능
 * /admin/(아이디) - 어드민 부여
 * /ban/(아이디 or 아이피) - 차단

등등등등

## 기타
 * [테스트 서버](https://namu.ml/)
 
## 의존성
 * pip install [flask](https://pypi.python.org/pypi/Flask/0.12)
 * pip install [bcrypt](https://pypi.python.org/pypi/bcrypt/3.1.0)
 * pip install [pymysql](https://pypi.python.org/pypi/PyMySQL)
 * [MariaDB](https://mariadb.org/)나 [MySQL](https://www.mysql.com/)
 
## set.json 설명
 * db = 데이터베이스 이름
 * host = 데이터베이스 호스트 
 * user = 데이터베이스 사용자명
 * pw = 그 사용자의 비밀번호
 * owner = 소유자의 위키 내에서 사용 할 이름
 * name = 위키 이름
 * frontpage = 위키 대문
 * license = 하단에 표기 될 라이선스
 * key = 세션 키
 * upload = 업로드 제한 (메가 바이트 단위)
 * port = 위키 열 포트