from flask import Flask, session, request

from urllib import parse
import json
import pymysql
import time
import re
import hashlib

json_data = open('set.json').read()
set_data = json.loads(json_data)

conn = pymysql.connect(host = set_data['host'], user = set_data['user'], password = set_data['pw'], charset = 'utf8mb4')
curs = conn.cursor(pymysql.cursors.DictCursor)

def db_com():
    conn.commit()

def url_pas(data):
    return parse.quote(data).replace('/','%2F')
    
def db_get():
    return curs.fetchall()

def sha224(data):
    return hashlib.sha3_224(bytes(data, 'utf-8')).hexdigest()
    
db_ex = curs.execute
db_pas = pymysql.escape_string

db_ex("use " + set_data['db'])

def diff(seqm):
    output= []
    for opcode, a0, a1, b0, b1 in seqm.get_opcodes():
        if(opcode == 'equal'):
            output.append(seqm.a[a0:a1])
        elif(opcode == 'insert'):
            output.append("<span style='background:#CFC;'>" + seqm.b[b0:b1] + "</span>")
        elif(opcode == 'delete'):
            output.append("<span style='background:#FDD;'>" + seqm.a[a0:a1] + "</span>")
        elif(opcode == 'replace'):
            output.append("<span style='background:#CFC;'>" + seqm.b[b0:b1] + "</span><span style='background:#FDD;'>" + seqm.a[a0:a1] + "</span>")
        else:
            output.append(seqm.a[a0:a1])
    return ''.join(output)
           
def admin_check():
    if(session.get('Now') == True):
        ip = ip_check()
        db_ex("select * from user where id = '" + db_pas(ip) + "'")
        user = db_get()
        if(user):
            if(user[0]['acl'] == 'owner' or user[0]['acl'] == 'admin'):
                return 1
                
def owner_check():
    if(session.get('Now') == True):
        ip = ip_check()
        db_ex("select * from user where id = '" + db_pas(ip) + "'")
        user = db_get()
        if(user):
            if(user[0]['acl'] == 'owner'):
                return 1
                
def include_check(name, data):
    if(re.search('^틀:', name)):
        db_ex("select * from back where title = '" + db_pas(name) + "' and type = 'include'")
        back = db_get()
        if(back):
            i = 0

            while(True):
                try:
                    namumark(back[i]['link'], data)
                except:
                    break
                    
                i += 1
    
def login_check():
    if(session.get('Now') == True):
        return 1
    else:
        return 0

def ip_pas(raw_ip):
    yes = re.search("([^-]*)\s\-\s(Close|Reopen|Stop|Restart|Admin|Agreement|Settlement)$", raw_ip)
    if(yes):
        results = yes.groups()
        
        db_ex("select title from data where title = '사용자:" + db_pas(results[0]) + "'")
        row = db_get()
        if(row):
            ip = '<a href="/w/' + url_pas('사용자:' + results[0]) + '">' + results[0] + '</a> - ' + results[1] + ' <a href="/record/' + url_pas(results[0]) + '/n/1">(기록)</a>'
        else:
            ip = '<a class="not_thing" href="/w/' + url_pas('사용자:' + results[0]) + '">' + results[0] + '</a> - ' + results[1] + ' <a href="/record/' + url_pas(results[0]) + '/n/1">(기록)</a>'
    elif(re.search("\.", raw_ip)):
        ip = raw_ip + ' <a href="/record/' + url_pas(raw_ip) + '/n/1">(기록)</a>'
    else:
        db_ex("select title from data where title = '사용자:" + db_pas(raw_ip) + "'")
        row = db_get()
        if(row):
            ip = '<a href="/w/' + url_pas('사용자:' + raw_ip) + '">' + raw_ip + '</a> <a href="/record/' + url_pas(raw_ip) + '/n/1">(기록)</a>'
        else:
            ip = '<a class="not_thing" href="/w/' + url_pas('사용자:' + raw_ip) + '">' + raw_ip + '</a> <a href="/record/' + url_pas(raw_ip) + '/n/1">(기록)</a>'

    return ip

def ip_check():
    if(session.get('Now') == True):
        ip = format(session['DREAMER'])
    else:
        if(request.headers.getlist("X-Forwarded-For")):
            ip = request.headers.getlist("X-Forwarded-For")[0]
        else:
            ip = request.remote_addr
            
    return ip

def custom_css_user():
    if(session.get('Now') == True):
        try:
            data = format(session['Daydream'])
        except:
            data = ''
    else:
        data = ''

    return data

def acl_check(ip, name):
    m = re.search("^사용자:(.*)", name)
    n = re.search("^파일:(.*)", name)
    if(m):
        g = m.groups()
        if(ip == g[0]):
            if(re.search("\.", g[0])):
                return 1
            else:
                db_ex("select * from ban where block = '" + db_pas(ip) + "'")
                rows = db_get()
                if(rows):
                    return 1
                else:
                    return 0
        else:
            return 1
    elif(n):
        if(not owner_check() == 1):
            return 1
    else:
        b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
        if(b):
            results = b.groups()
            db_ex("select * from ban where block = '" + db_pas(results[0]) + "' and band = 'O'")
            rowss = db_get()
            if(rowss):
                return 1
            else:
                db_ex("select * from ban where block = '" + db_pas(ip) + "'")
                rows = db_get()
                if(rows):
                    return 1
                else:
                    db_ex("select acl from data where title = '" + db_pas(name) + "'")
                    row = db_get()
                    if(row):
                        db_ex("select * from user where id = '" + db_pas(ip) + "'")
                        rows = db_get()
                        if(row[0]['acl'] == 'user'):
                            if(rows):
                                return 0
                            else:
                                return 1
                        elif(row[0]['acl'] == 'admin'):
                            if(rows):
                                if(rows[0]['acl'] == 'admin' or rows[0]['acl'] == 'owner'):
                                    return 0
                                else:
                                    return 1
                            else:
                                return 1
                        else:
                            return 0
                    else:
                        return 0
        else:
            db_ex("select * from ban where block = '" + db_pas(ip) + "'")
            rows = db_get()
            if(rows):
                return 1
            else:
                db_ex("select acl from data where title = '" + db_pas(name) + "'")
                row = db_get()
                if(row):
                    db_ex("select * from user where id = '" + db_pas(ip) + "'")
                    rows = db_get()
                    if(row[0]['acl'] == 'user'):
                        if(rows):
                            return 0
                        else:
                            return 1
                    elif(row[0]['acl'] == 'admin'):
                        if(rows):
                            if(rows[0]['acl'] == 'admin' or rows[0]['acl'] == 'owner'):
                                return 0
                            else:
                                return 1
                        else:
                            return 1
                    else:
                        return 0
                else:
                    return 0

def ban_check(ip):
    b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
    if(b):
        results = b.groups()
        db_ex("select * from ban where block = '" + db_pas(results[0]) + "' and band = 'O'")
        rowss = db_get()
        if(rowss):
            return 1
        else:
            db_ex("select * from ban where block = '" + db_pas(ip) + "'")
            rows = db_get()
            if(rows):
                return 1
            else:
                return 0
    else:
        db_ex("select * from ban where block = '" + db_pas(ip) + "'")
        rows = db_get()
        if(rows):
            return 1
        else:
            return 0
        
def topic_check(ip, name, sub):
    b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
    if(b):
        results = b.groups()
        db_ex("select * from ban where block = '" + db_pas(results[0]) + "' and band = 'O'")
        rowss = db_get()
        if(rowss):
            return 1
        else:
            db_ex("select * from ban where block = '" + db_pas(ip) + "'")
            rows = db_get()
            if(rows):
                return 1
            else:
                db_ex("select * from stop where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "'")
                rows = db_get()
                if(rows):
                    return 1
                else:
                    return 0
    else:
        db_ex("select * from ban where block = '" + db_pas(ip) + "'")
        rows = db_get()
        if(rows):
            return 1
        else:
            db_ex("select * from stop where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "'")
            rows = db_get()
            if(rows):
                return 1
            else:
                return 0

def get_time():
    now = time.localtime()
    date = "%04d-%02d-%02d %02d:%02d:%02d" % (now.tm_year, now.tm_mon, now.tm_mday, now.tm_hour, now.tm_min, now.tm_sec)
    
    return date

def rd_plus(title, sub, date):
    db_ex("select * from rd where title = '" + db_pas(title) + "' and sub = '" + db_pas(sub) + "'")
    rd = db_get()
    if(rd):
        db_ex("update rd set date = '" + db_pas(date) + "' where title = '" + db_pas(title) + "' and sub = '" + db_pas(sub) + "'")
    else:
        db_ex("insert into rd (title, sub, date) value ('" + db_pas(title) + "', '" + db_pas(sub) + "', '" + db_pas(date) + "')")
    db_com()
    
def rb_plus(block, end, today, blocker, why):
    db_ex("insert into rb (block, end, today, blocker, why) value ('" + db_pas(block) + "', '" + db_pas(end) + "', '" + today + "', '" + db_pas(blocker) + "', '" + db_pas(why) + "')")
    db_com()

def history_plus(title, data, date, ip, send, leng):
    db_ex("select * from history where title = '" + db_pas(title) + "' order by id+0 desc limit 1")
    rows = db_get()
    if(rows):
        number = int(rows[0]['id']) + 1
        db_ex("insert into history (id, title, data, date, ip, send, leng) value ('" + str(number) + "', '" + db_pas(title) + "', '" + db_pas(data) + "', '" + date + "', '" + db_pas(ip) + "', '" + db_pas(send) + "', '" + leng + "')")
        db_com()
    else:
        db_ex("insert into history (id, title, data, date, ip, send, leng) value ('1', '" + db_pas(title) + "', '" + db_pas(data) + "', '" + date + "', '" + db_pas(ip) + "', '" + db_pas(send + ' (새 문서)') + "', '" + leng + "')")
        db_com()

def leng_check(a, b):
    if(a < b):
        c = b - a
        c = '+' + str(c)
    elif(b < a):
        c = a - b
        c = '-' + str(c)
    else:
        c = '0'
        
    return c