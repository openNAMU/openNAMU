import json
import sqlite3
import threading
from func import *
from mark import namumark

json_data = open('set.json').read()
set_data = json.loads(json_data)

conn = sqlite3.connect(set_data['db'] + '.db', check_same_thread = False)
curs = conn.cursor()

def parser(data):
    namumark(conn, data[0], data[1], 1, 0, 0)

curs.execute("delete from back")
conn.commit()

curs.execute("select title, data from data")
data = curs.fetchall()

for test in data:
    print(test[0])
    t = threading.Thread(target = parser, args = [test])
    t.start()
    t.join()

conn.commit()