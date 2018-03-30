# 주요 모듈 불러옴
import json
import sqlite3
import threading

# 기타 코드 불러옴
from func import *
from mark import namumark

# JSON 불러옴
json_data = open('set.json').read()
set_data = json.loads(json_data)

# 디비 연결
conn = sqlite3.connect(set_data['db'] + '.db')
curs = conn.cursor()

# 파싱 해주는 함수
def parser(data):
    namumark(conn, data[0], data[1], 1)

# 역링크 전부 삭제
curs.execute("delete from back")
conn.commit()

# 데이터에서 제목이랑 내용 불러옴
curs.execute("select title, data from data")
data = curs.fetchall()

# for 돌려서 처리
for test in data:
    # 제목 프린트
    print(test[0])

    # 파싱
    parser(test)

# 커밋
conn.commit()