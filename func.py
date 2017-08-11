from bottle import request, app
from bottle.ext import beaker
from urllib import parse
import json
import pymysql
import time
import re
import hashlib

json_data = open('set.json').read()
set_data = json.loads(json_data)
    
session_opts = {
    'session.type': 'file',
    'session.data_dir': './app_session/',
    'session.auto': True
}

app = beaker.middleware.SessionMiddleware(app(), session_opts)

from mark import *

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
            
    return(''.join(output))
           
def admin_check(num):
    conn = pymysql.connect(
        user = set_data['user'], 
        password = set_data['pw'], 
        charset = 'utf8mb4', 
        db = set_data['db']
    )
    curs = conn.cursor(pymysql.cursors.DictCursor)

    ip = ip_check() 
    curs.execute("select acl from user where id = '" + pymysql.escape_string(ip) + "'")
    user = curs.fetchall()
    if(user):
        reset = 0
        while(1):
            if(num == 1 and reset == 0):
                curs.execute('select name from alist where name = "' + pymysql.escape_string(user[0]["acl"]) + '" and acl = "ban"')
                acl_data = curs.fetchall()
                if(acl_data):
                    conn.close()
                    return(1)
                else:
                    reset = 1
            elif(num == 2 and reset == 0):
                curs.execute('select name from alist where name = "' + pymysql.escape_string(user[0]["acl"]) + '" and acl = "mdel"')
                acl_data = curs.fetchall()
                if(acl_data):
                    conn.close()
                    return(1)
                else:
                    reset = 1
            elif(num == 3 and reset == 0):
                curs.execute('select name from alist where name = "' + pymysql.escape_string(user[0]["acl"]) + '" and acl = "toron"')
                acl_data = curs.fetchall()
                if(acl_data):
                    conn.close()
                    return(1)
                else:
                    reset = 1
            elif(num == 4 and reset == 0):
                curs.execute('select name from alist where name = "' + pymysql.escape_string(user[0]["acl"]) + '" and acl = "check"')
                acl_data = curs.fetchall()
                if(acl_data):
                    conn.close()
                    return(1)
                else:
                    reset = 1
            elif(num == 5 and reset == 0):
                curs.execute('select name from alist where name = "' + pymysql.escape_string(user[0]["acl"]) + '" and acl = "acl"')
                acl_data = curs.fetchall()
                if(acl_data):
                    conn.close()
                    return(1)
                else:
                    reset = 1
            elif(num == 6 and reset == 0):
                curs.execute('select name from alist where name = "' + pymysql.escape_string(user[0]["acl"]) + '" and acl = "hidel"')
                acl_data = curs.fetchall()
                if(acl_data):
                    conn.close()
                    return(1)
                else:
                    reset = 1
            else:
                curs.execute('select name from alist where name = "' + pymysql.escape_string(user[0]["acl"]) + '" and acl = "owner"')
                acl_data = curs.fetchall()
                if(acl_data):
                    conn.close()
                    return(1)
                else:
                    break
    conn.close()
                
def include_check(name, data):
    conn = pymysql.connect(
        user = set_data['user'], 
        password = set_data['pw'], 
        charset = 'utf8mb4', 
        db = set_data['db']
    )
    curs = conn.cursor(pymysql.cursors.DictCursor)

    if(re.search('^틀:', name)):
        curs.execute("select link from back where title = '" + pymysql.escape_string(name) + "' and type = 'include'")
        back = curs.fetchall()
        for backp in back:
            namumark(backp['link'], data, 1)
    
    conn.close()
    
def login_check():
    session = request.environ.get('beaker.session')
    if(session.get('Now') == 1):
        return(1)
    else:
        return(0)

def ip_pas(raw_ip, num):
    conn = pymysql.connect(
        user = set_data['user'], 
        password = set_data['pw'], 
        charset = 'utf8mb4', 
        db = set_data['db']
    )
    curs = conn.cursor(pymysql.cursors.DictCursor)
    
    if(re.search("(\.|:)", raw_ip)):
        ip = raw_ip
    else:
        curs.execute("select title from data where title = '사용자:" + pymysql.escape_string(raw_ip) + "'")
        row = curs.fetchall()
        if(row):
            ip = '<a href="/w/' + url_pas('사용자:' + raw_ip) + '">' + raw_ip + '</a>'
        else:
            ip = '<a class="not_thing" href="/w/' + url_pas('사용자:' + raw_ip) + '">' + raw_ip + '</a>'
            
    if(num == 1):
        ip += ' <a href="/user/' + url_pas(raw_ip) + '/topic">(기록)</a>'
    elif(num == 2):
        ip += ' <a href="/record/' + url_pas(raw_ip) + '">(기록)</a> <a href="/user/' + url_pas(raw_ip) + '/topic">(토론 기록)</a>'        
    else:
        ip += ' <a href="/record/' + url_pas(raw_ip) + '">(기록)</a>'
    
    conn.close()
    
    return(ip)

def custom_css_user():
    session = request.environ.get('beaker.session')
    try:
        data = format(session['Daydream'])
    except:
        data = ''

    return(data)

def acl_check(ip, name):
    conn = pymysql.connect(
        user = set_data['user'], 
        password = set_data['pw'], 
        charset = 'utf8mb4', 
        db = set_data['db']
    )
    curs = conn.cursor(pymysql.cursors.DictCursor)

    m = re.search("^사용자:([^/]*)", name)
    n = re.search("^파일:(.*)", name)
    if(m):
        g = m.groups()
        if(ip == g[0]):
            if(re.search("(\.|:)", g[0])):
                conn.close()
                return(1)
            else:
                curs.execute("select * from ban where block = '" + pymysql.escape_string(ip) + "'")
                rows = curs.fetchall()
                if(rows):
                    conn.close()
                    return(1)
                else:
                    conn.close()
                    return(0)
        else:
            conn.close()
            return(1)
    elif(n):
        if(admin_check(None) != 1):
            conn.close()
            return(1)
    else:
        b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
        if(b):
            results = b.groups()
            curs.execute("select * from ban where block = '" + pymysql.escape_string(results[0]) + "' and band = 'O'")
            rowss = curs.fetchall()
            if(rowss):
                conn.close()
                return(1)
            else:
                curs.execute("select * from ban where block = '" + pymysql.escape_string(ip) + "'")
                rows = curs.fetchall()
                if(rows):
                    conn.close()
                    return(1)
                else:
                    curs.execute("select acl from data where title = '" + pymysql.escape_string(name) + "'")
                    row = curs.fetchall()
                    if(row):
                        curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
                        rows = curs.fetchall()
                        if(row[0]['acl'] == 'user'):
                            if(rows):
                                conn.close()
                                return(0)
                            else:
                                conn.close()
                                return(1)
                        elif(row[0]['acl'] == 'admin'):
                            if(rows):
                                if(rows[0]['acl'] == 'admin' or rows[0]['acl'] == 'owner'):
                                    conn.close()
                                    return(0)
                                else:
                                    conn.close()
                                    return(1)
                            else:
                                conn.close()
                                return(1)
                        else:
                            conn.close()
                            return(0)
                    else:
                        conn.close()
                        return(0)
        else:
            curs.execute("select * from ban where block = '" + pymysql.escape_string(ip) + "'")
            rows = curs.fetchall()
            if(rows):
                conn.close()
                return(1)
            else:
                curs.execute("select acl from data where title = '" + pymysql.escape_string(name) + "'")
                row = curs.fetchall()
                if(row):
                    curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
                    rows = curs.fetchall()
                    if(row[0]['acl'] == 'user'):
                        if(rows):
                            conn.close()
                            return(0)
                        else:
                            conn.close()
                            return(1)
                    elif(row[0]['acl'] == 'admin'):
                        if(rows):
                            if(rows[0]['acl'] == 'admin' or rows[0]['acl'] == 'owner'):
                                conn.close()
                                return(0)
                            else:
                                conn.close()
                                return(1)
                        else:
                            conn.close()
                            return(1)
                    else:
                        conn.close()
                        return(0)
                else:
                    conn.close()
                    return(0)
    conn.close()

def ban_check(ip):
    conn = pymysql.connect(
        user = set_data['user'], 
        password = set_data['pw'], 
        charset = 'utf8mb4', 
        db = set_data['db']
    )
    curs = conn.cursor(pymysql.cursors.DictCursor)

    b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
    if(b):
        results = b.groups()
        curs.execute("select * from ban where block = '" + pymysql.escape_string(results[0]) + "' and band = 'O'")
        rowss = curs.fetchall()
        if(rowss):
            conn.close()
            return(1)
        else:
            curs.execute("select * from ban where block = '" + pymysql.escape_string(ip) + "'")
            rows = curs.fetchall()
            if(rows):
                conn.close()
                return(1)
            else:
                conn.close()
                return(0)
    else:
        curs.execute("select * from ban where block = '" + pymysql.escape_string(ip) + "'")
        rows = curs.fetchall()
        if(rows):
            conn.close()
            return(1)
        else:
            conn.close()
            return(0)
    conn.close()
        
def topic_check(ip, name, sub):
    conn = pymysql.connect(
        user = set_data['user'], 
        password = set_data['pw'], 
        charset = 'utf8mb4', 
        db = set_data['db']
    )
    curs = conn.cursor(pymysql.cursors.DictCursor)

    b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
    if(b):
        results = b.groups()
        curs.execute("select * from ban where block = '" + pymysql.escape_string(results[0]) + "' and band = 'O'")
        rowss = curs.fetchall()
        if(rowss):
            conn.close()
            return(1)
        else:
            curs.execute("select * from ban where block = '" + pymysql.escape_string(ip) + "'")
            rows = curs.fetchall()
            if(rows):
                conn.close()
                return(1)
            else:
                curs.execute("select * from stop where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "'")
                rows = curs.fetchall()
                if(rows):
                    conn.close()
                    return(1)
                else:
                    conn.close()
                    return(0)
    else:
        curs.execute("select * from ban where block = '" + pymysql.escape_string(ip) + "'")
        rows = curs.fetchall()
        if(rows):
            conn.close()
            return(1)
        else:
            curs.execute("select * from stop where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "'")
            rows = curs.fetchall()
            if(rows):
                conn.close()
                return(1)
            else:
                conn.close()
                return(0)
    conn.close()

def rd_plus(title, sub, date):
    conn = pymysql.connect(
        user = set_data['user'], 
        password = set_data['pw'], 
        charset = 'utf8mb4', 
        db = set_data['db']
    )
    curs = conn.cursor(pymysql.cursors.DictCursor)

    curs.execute("select * from rd where title = '" + pymysql.escape_string(title) + "' and sub = '" + pymysql.escape_string(sub) + "'")
    rd = curs.fetchall()
    if(rd):
        curs.execute("update rd set date = '" + pymysql.escape_string(date) + "' where title = '" + pymysql.escape_string(title) + "' and sub = '" + pymysql.escape_string(sub) + "'")
    else:
        curs.execute("insert into rd (title, sub, date) value ('" + pymysql.escape_string(title) + "', '" + pymysql.escape_string(sub) + "', '" + pymysql.escape_string(date) + "')")
    conn.commit()
    
    conn.close()
    
def rb_plus(block, end, today, blocker, why):
    conn = pymysql.connect(
        user = set_data['user'], 
        password = set_data['pw'], 
        charset = 'utf8mb4', 
        db = set_data['db']
    )
    curs = conn.cursor(pymysql.cursors.DictCursor)

    curs.execute("insert into rb (block, end, today, blocker, why) value ('" + pymysql.escape_string(block) + "', '" + pymysql.escape_string(end) + "', '" + today + "', '" + pymysql.escape_string(blocker) + "', '" + pymysql.escape_string(why) + "')")
    conn.commit()
    
    conn.close()

def history_plus(title, data, date, ip, send, leng):
    conn = pymysql.connect(
        user = set_data['user'], 
        password = set_data['pw'], 
        charset = 'utf8mb4', 
        db = set_data['db']
    )
    curs = conn.cursor(pymysql.cursors.DictCursor)

    curs.execute("select * from history where title = '" + pymysql.escape_string(title) + "' order by id+0 desc limit 1")
    rows = curs.fetchall()
    if(rows):
        number = int(rows[0]['id']) + 1
        curs.execute("insert into history (id, title, data, date, ip, send, leng) value ('" + str(number) + "', '" + pymysql.escape_string(title) + "', '" + pymysql.escape_string(data) + "', '" + date + "', '" + pymysql.escape_string(ip) + "', '" + pymysql.escape_string(send) + "', '" + leng + "')")
    else:
        curs.execute("insert into history (id, title, data, date, ip, send, leng) value ('1', '" + pymysql.escape_string(title) + "', '" + pymysql.escape_string(data) + "', '" + date + "', '" + pymysql.escape_string(ip) + "', '" + pymysql.escape_string(send + ' (새 문서)') + "', '" + leng + "')")
    conn.commit()
    
    conn.close()

def leng_check(a, b):
    if(a < b):
        c = b - a
        c = '+' + str(c)
    elif(b < a):
        c = a - b
        c = '-' + str(c)
    else:
        c = '0'
        
    return(c)