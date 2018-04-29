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
conn = sqlite3.connect(set_data['db'] + '.db', check_same_thread = False)
curs = conn.cursor()

print('1. BackLink ReSet')
print('2. ReCaptcha Delete')
print('3. Ban Delete')
print('')

print('select : ', end = '')
what_i_do = input()

if what_i_do == '1':
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
        # 스레드 기반으로 처리
        t = threading.Thread(target = parser, args = [test])
        t.start()
        t.join()
elif what_i_do == '2':
    # 데이터 삭제
    curs.execute("delete from other where name = 'recaptcha'")
    curs.execute("delete from other where name = 'sec_re'")
elif what_i_do == '3':
    print('IP or User_Name : ', end = '')
    user_data = input()

    if re.search("^([0-9]{1,3}\.[0-9]{1,3})$", user_data):
        band = 'O'
    else:
        band = ''

    # 데이터 삭제
    curs.execute("insert into rb (block, end, today, blocker, why, band) values (?, ?, ?, ?, ?, ?)", [user_data, '해제', get_time(), 'localhost', '', band])
    curs.execute("delete from ban where block = ?", [user_data]))

# 커밋
conn.commit()