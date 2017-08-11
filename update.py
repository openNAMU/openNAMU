import pymysql
import json
import re

json_data = open('set.json').read()
set_data = json.loads(json_data)

conn = pymysql.connect(
    user = set_data['user'], 
    password = set_data['pw'], 
    charset = 'utf8mb4', 
    db = set_data['db']
)
curs = conn.cursor(pymysql.cursors.DictCursor)   

r_ver = '2.1.7'

conn.commit()
conn.close()

print('종료 하려면 아무 키나 누르시오.')
a = input()