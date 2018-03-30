import json
import sqlite3

# JSON 연결
json_data = open('set.json').read()
set_data = json.loads(json_data)

# 디비 연결
conn = sqlite3.connect(set_data['db'] + '.db')
curs = conn.cursor()

# 데이터 삭제
curs.execute("delete from other where name = 'recaptcha'")
curs.execute("delete from other where name = 'sec_re'")
conn.commit()