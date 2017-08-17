import sqlite3
import pymysql
import json
import re

json_data = open('set.json').read()
set_data = json.loads(json_data)

conn2 = sqlite3.connect(set_data['db'] + '.db')
curs2 = conn2.cursor()

def escape(data):
    data = data.replace("'", "''")
    data = data.replace('"', '""')

    return(data)

conn = pymysql.connect(
    user = set_data['user'], 
    password = set_data['pw'], 
    charset = 'utf8mb4', 
    db = set_data['db']
)
curs = conn.cursor(pymysql.cursors.DictCursor)   

r_ver = '2.2.0'

curs2.execute("create table data(title text, data text, acl text)")
curs2.execute("create table history(id text, title text, data text, date text, ip text, send text, leng text)")
curs2.execute("create table rd(title text, sub text, date text)")
curs2.execute("create table user(id text, pw text, acl text)")
curs2.execute("create table ban(block text, end text, why text, band text)")
curs2.execute("create table topic(id text, title text, sub text, data text, date text, ip text, block text, top text)")
curs2.execute("create table stop(title text, sub text, close text)")
curs2.execute("create table rb(block text, end text, today text, blocker text, why text)")
curs2.execute("create table login(user text, ip text, today text)")
curs2.execute("create table back(title text, link text, type text)")
curs2.execute("create table cat(title text, cat text)")
curs2.execute("create table hidhi(title text, re text)")
curs2.execute("create table agreedis(title text, sub text)")
curs2.execute("create table custom(user text, css text)")
curs2.execute("create table other(name text, data text)")
curs2.execute("create table alist(name text, acl text)")

curs2.execute("insert into other (name, data) values ('version', '" + escape(r_ver) + "')")

conn2.commit()

curs.execute("select * from data")
data = curs.fetchall()
for data2 in data:
    curs2.execute("insert into data (title, data, acl) values ('" + escape(data2['title']) + "', '" + escape(data2['data']) + "', '" + escape(data2['acl']) + "')")

curs.execute("select * from history")
data = curs.fetchall()
for data2 in data:
    curs2.execute("insert into history (id, title, data, date, ip, send, leng) values ('" + escape(data2['id']) + "', '" + escape(data2['title']) + "', '" + escape(data2['data']) + "', '" + escape(data2['date']) + "', '" + escape(data2['ip']) + "', '" + escape(data2['send']) + "', '" + escape(data2['leng']) + "')")

curs.execute("select * from rd")
data = curs.fetchall()
for data2 in data:
    curs2.execute("insert into rd (title, sub, date) values ('" + escape(data2['title']) + "', '" + escape(data2['sub']) + "', '" + escape(data2['date']) + "')")

curs.execute("select * from user")
data = curs.fetchall()
for data2 in data:
    curs2.execute("insert into user (id, pw, acl) values ('" + escape(data2['id']) + "', '" + escape(data2['pw']) + "', '" + escape(data2['acl']) + "')")

curs.execute("select * from ban")
data = curs.fetchall()
for data2 in data:
    curs2.execute("insert into ban (block, end, why, band) values ('" + escape(data2['block']) + "', '" + escape(data2['end']) + "', '" + escape(data2['why']) + "', '" + escape(data2['band']) + "')")

curs.execute("select * from topic")
data = curs.fetchall()
for data2 in data:
    curs2.execute("insert into topic (id, title, sub, data, date, ip, block, top) values ('" + escape(data2['id']) + "', '" + escape(data2['title']) + "', '" + escape(data2['sub']) + "', '" + escape(data2['data']) + "', '" + escape(data2['date']) + "', '" + escape(data2['ip']) + "', '" + escape(data2['block']) + "', '" + escape(data2['top']) + "')")

curs.execute("select * from stop")
data = curs.fetchall()
for data2 in data:
    curs2.execute("insert into stop (title, sub, close) values ('" + escape(data2['title']) + "', '" + escape(data2['sub']) + "', '" + escape(data2['close']) + "')")

curs.execute("select * from rb")
data = curs.fetchall()
for data2 in data:
    curs2.execute("insert into rb (block, end, today, blocker, why) values ('" + escape(data2['block']) + "', '" + escape(data2['end']) + "', '" + escape(data2['today']) + "', '" + escape(data2['blocker']) + "', '" + escape(data2['why']) + "')")

curs.execute("select * from login")
data = curs.fetchall()
for data2 in data:
    curs2.execute("insert into login (user, ip, today) values ('" + escape(data2['user']) + "', '" + escape(data2['ip']) + "', '" + escape(data2['today']) + "')")

curs.execute("select * from cat")
data = curs.fetchall()
for data2 in data:
    curs2.execute("insert into cat (title, cat) values ('" + escape(data2['title']) + "', '" + escape(data2['cat']) + "')")

curs.execute("select * from back")
data = curs.fetchall()
for data2 in data:
    curs2.execute("insert into back (title, link, type) values ('" + escape(data2['title']) + "', '" + escape(data2['link']) + "', '" + escape(data2['type']) + "')")

curs.execute("select * from hidhi")
data = curs.fetchall()
for data2 in data:
    curs2.execute("insert into hidhi (title, re) values ('" + escape(data2['title']) + "', '" + escape(data2['re']) + "')")

curs.execute("select * from agreedis")
data = curs.fetchall()
for data2 in data:
    curs2.execute("insert into agreedis (title, sub) values ('" + escape(data2['title']) + "', '" + escape(data2['sub']) + "')")

curs.execute("select * from custom")
data = curs.fetchall()
for data2 in data:
    curs2.execute("insert into custom (user, css) values ('" + escape(data2['user']) + "', '" + escape(data2['css']) + "')")

curs.execute("select * from alist")
data = curs.fetchall()
for data2 in data:
    curs2.execute("insert into alist (name, acl) values ('" + escape(data2['name']) + "', '" + escape(data2['acl']) + "')")

conn2.commit()

print('종료 하려면 아무 키나 누르시오.')
a = input()