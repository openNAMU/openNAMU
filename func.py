from bottle import request, app
from bottle.ext import beaker
from urllib import parse
import json
import sqlite3
import time
import re
import hashlib
import html

json_data = open('set.json').read()
set_data = json.loads(json_data)

conn = sqlite3.connect(set_data['db'] + '.db')
curs = conn.cursor()
    
session_opts = {
    'session.type': 'file',
    'session.data_dir': './app_session/',
    'session.auto': 1
}

app = beaker.middleware.SessionMiddleware(app(), session_opts)

from mark import *

def wiki_set(num):
    if(num == 1):
        d_data = '무명위키'
        curs.execute('select data from other where name = ?', ['name'])
    elif(num == 2):
        d_data = '위키:대문'
        curs.execute('select data from other where name = "frontpage"')
    elif(num == 3):
        d_data = 'CC 0'
        curs.execute('select data from other where name = "license"')
    elif(num == 4):
        d_data = '2'
        curs.execute('select data from other where name = "upload"')
    
    data = curs.fetchall()
    if(data):
        return(data[0][0])
    else:
        return(d_data)

def diff(seqm, num):
    output= []
    for opcode, a0, a1, b0, b1 in seqm.get_opcodes():
        if(opcode == 'equal' and num == 1):
            output.append(seqm.a[a0:a1])
        elif(opcode == 'insert' and num == 0):
            output.append("<span style='background:#CFC;'>" + seqm.b[b0:b1] + "</span>")
        elif(opcode == 'delete' and num == 1):
            output.append("<span style='background:#FDD;'>" + seqm.a[a0:a1] + "</span>")
        elif(opcode == 'replace'):
            if(num == 1):
                output.append("<span style='background:#FDD;'>" + seqm.a[a0:a1] + "</span>")
            else:
                output.append("<span style='background:#CFC;'>" + seqm.b[b0:b1] + "</span>")
        elif(num == 0):
            output.append(seqm.b[b0:b1])
            
    return(''.join(output))
           
def admin_check(num, what):
    ip = ip_check() 
    curs.execute("select acl from user where id = ?", [ip])
    user = curs.fetchall()
    if(user):
        reset = 0
        while(1):
            if(num == 1 and reset == 0):
                check = 'ban'
            elif(num == 2 and reset == 0):
                check = 'mdel'
            elif(num == 3 and reset == 0):
                check = 'toron'
            elif(num == 4 and reset == 0):
                check = 'check'
            elif(num == 5 and reset == 0):
                check = 'acl'
            elif(num == 6 and reset == 0):
                check = 'hidel'
            else:
                check = 'owner'

            curs.execute('select name from alist where name = ? and acl = ?', [user[0][0], check])
            acl_data = curs.fetchall()
            if(acl_data):
                if(what):
                    curs.execute("insert into re_admin (who, what, time) values (?, ?, ?)", [ip, what, get_time()])
                    conn.commit()

                return(1)
            else:
                if(reset == 0):
                    reset = 1
                else:
                    break
                
def include_check(name, data):
    if(re.search('^틀:', name)):
        curs.execute("select link from back where title = ? and type = 'include'", [name])
        back = curs.fetchall()
        for back_p in back:
            namumark(back_p[0], data, 1, 1)  
    
def login_check():
    session = request.environ.get('beaker.session')
    if(session.get('Now') == 1):
        return(1)
    else:
        return(0)

def ip_pas(raw_ip, num):
    if(re.search("(\.|:)", raw_ip)):
        ip = raw_ip
    else:
        curs.execute("select title from data where title = ?", ['사용자:' + raw_ip])
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

    return(ip)

def custom_css():
    session = request.environ.get('beaker.session')
    try:
        data = format(session['Daydream'])
    except:
        data = ''

    return(data)

def custom_js():
    session = request.environ.get('beaker.session')
    try:
        data = format(session['AQUARIUM'])
    except:
        data = ''

    return(data)

def acl_check(name):
    ip = ip_check()
    m = re.search("^사용자:([^/]*)", name)
    n = re.search("^파일:(.*)", name)
    if(m):
        g = m.groups()
        if(ip == g[0]):
            if(re.search("(\.|:)", g[0])):
                return(1)
            else:
                curs.execute("select block from ban where block = ?", [ip])
                rows = curs.fetchall()
                if(rows):
                    return(1)
                else:
                    return(0)
        else:
            return(1)
    elif(n and admin_check(5, 'edit (' + name + ')') != 1):
        return(1)
    else:
        b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
        if(b):
            results = b.groups()
            curs.execute("select block from ban where block = ? and band = 'O'", [results[0]])
            rowss = curs.fetchall()
            if(rowss):
                return(1)

        curs.execute("select block from ban where block = ?", [ip])
        rows = curs.fetchall()
        if(rows):
            return(1)
        else:
            curs.execute("select acl from data where title = ?", [name])
            row = curs.fetchall()
            if(row):
                curs.execute("select acl from user where id = ?", [ip])
                rows = curs.fetchall()
                if(row[0][0] == 'user'):
                    if(rows):
                        return(0)
                    else:
                        return(1)
                elif(row[0][0] == 'admin'):
                    if(rows and admin_check(5, 'edit (' + name + ')') == 1):
                        return(0)
                    else:
                        return(1)
                else:
                    return(0)
            else:
                return(0)    

def ban_check():
    ip = ip_check()
    b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
    if(b):
        results = b.groups()
        curs.execute("select block from ban where block = ? and band = 'O'", [results[0]])
        rowss = curs.fetchall()
        if(rowss):
            return(1)

    curs.execute("select block from ban where block = ?", [ip])
    rows = curs.fetchall()
    if(rows):
        return(1)
    else:
        return(0)
        
def topic_check(name, sub):
    ip = ip_check()
    b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
    if(b):
        results = b.groups()
        curs.execute("select block from ban where block = ? and band = 'O'", [results[0]])
        rowss = curs.fetchall()
        if(rowss):
            return(1)

    curs.execute("select block from ban where block = ?", [ip])
    rows = curs.fetchall()
    if(rows):
        return(1)
    else:
        curs.execute("select title from stop where title = ? and sub = ?", [name, sub])
        rows = curs.fetchall()
        if(rows):
            return(1)
        else:
            return(0)

def rd_plus(title, sub, date):
    curs.execute("select title from rd where title = ? and sub = ?", [title, sub])
    rd = curs.fetchall()
    if(rd):
        curs.execute("update rd set date = ? where title = ? and sub = ?", [date, title, sub])
    else:
        curs.execute("insert into rd (title, sub, date) values (?, ?, ?)", [title, sub, date])
    conn.commit()
    
def rb_plus(block, end, today, blocker, why):
    curs.execute("insert into rb (block, end, today, blocker, why) values (?, ?, ?, ?, ?)", [block, end, today, blocker, why])
    conn.commit()

def history_plus(title, data, date, ip, send, leng):
    curs.execute("select id from history where title = ? order by id+0 desc limit 1", [title])
    rows = curs.fetchall()
    if(rows):
        number = int(rows[0][0]) + 1
        curs.execute("insert into history (id, title, data, date, ip, send, leng) values (?, ?, ?, ?, ?, ?, ?)", [str(number), title, data, date, ip, send, leng])
    else:
        curs.execute("insert into history (id, title, data, date, ip, send, leng) values ('1', ?, ?, ?, ?, ?, ?)", [title, data, date, ip, send + ' (새 문서)', leng])
    conn.commit()

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
