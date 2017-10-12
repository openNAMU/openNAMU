import json
import sqlite3
from multiprocessing import Process
from func import *

json_data = open('set.json').read()
set_data = json.loads(json_data)

conn = sqlite3.connect(set_data['db'] + '.db')
curs = conn.cursor()

curs.execute("delete from back")
curs.execute("delete from cat")
conn.commit()

def go_namu(data):
    for end in data:
        print(end[0])
        namumark(end[0], end[1], 1, 0, 0)

curs.execute("select title, data from data")
data = curs.fetchall()
print(int(len(data) / 4))
l = int(len(data) / 4)

d1 = data[:l]
d2 = data[l:l * 2]
d3 = data[l * 2:l * 3]
d4 = data[l * 3:]

if(__name__=='__main__'):
    p1 = Process(target = go_namu, args = [d1])
    p2 = Process(target = go_namu, args = [d2])
    p3 = Process(target = go_namu, args = [d3])
    p4 = Process(target = go_namu, args = [d4])
    p1.start()
    p2.start()
    p3.start()
    p4.start()
    p1.join()
    p2.join()
    p3.join()
    p4.join()
    
conn.commit()