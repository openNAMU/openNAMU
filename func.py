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

conn = pymysql.connect(host = set_data['host'], user = set_data['user'], password = set_data['pw'], charset = 'utf8mb4')
curs = conn.cursor(pymysql.cursors.DictCursor)

def url_pas(data):
    return(parse.quote(data).replace('/','%2F'))

def sha224(data):
    return(hashlib.sha224(bytes(data, 'utf-8')).hexdigest())
    
session_opts = {
    'session.type': 'file',
    'session.data_dir': './app_session/',
    'session.auto': True
}

app = beaker.middleware.SessionMiddleware(app(), session_opts)

db_pas = pymysql.escape_string

def diff(seqm):
    conn = pymysql.connect(user = set_data['user'], password = set_data['pw'], charset = 'utf8mb4', db = set_data['db'])
    curs = conn.cursor(pymysql.cursors.DictCursor)

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
            
    conn.close()
    return(''.join(output))
           
def admin_check(num):
    conn = pymysql.connect(user = set_data['user'], password = set_data['pw'], charset = 'utf8mb4', db = set_data['db'])
    curs = conn.cursor(pymysql.cursors.DictCursor)

    ip = ip_check() 
    curs.execute("select acl from user where id = '" + db_pas(ip) + "'")
    user = curs.fetchall()
    if(user):
        reset = False
        while(True):
            if(num == 1 and reset == False):
                curs.execute('select name from alist where name = "' + db_pas(user[0]["acl"]) + '" and acl = "ban"')
                acl_data = curs.fetchall()
                if(acl_data):
                    conn.close()
                    return(1)
                else:
                    reset = True
            elif(num == 2 and reset == False):
                curs.execute('select name from alist where name = "' + db_pas(user[0]["acl"]) + '" and acl = "mdel"')
                acl_data = curs.fetchall()
                if(acl_data):
                    conn.close()
                    return(1)
                else:
                    reset = True
            elif(num == 3 and reset == False):
                curs.execute('select name from alist where name = "' + db_pas(user[0]["acl"]) + '" and acl = "toron"')
                acl_data = curs.fetchall()
                if(acl_data):
                    conn.close()
                    return(1)
                else:
                    reset = True
            elif(num == 4 and reset == False):
                curs.execute('select name from alist where name = "' + db_pas(user[0]["acl"]) + '" and acl = "check"')
                acl_data = curs.fetchall()
                if(acl_data):
                    conn.close()
                    return(1)
                else:
                    reset = True
            elif(num == 5 and reset == False):
                curs.execute('select name from alist where name = "' + db_pas(user[0]["acl"]) + '" and acl = "acl"')
                acl_data = curs.fetchall()
                if(acl_data):
                    conn.close()
                    return(1)
                else:
                    reset = True
            elif(num == 6 and reset == False):
                curs.execute('select name from alist where name = "' + db_pas(user[0]["acl"]) + '" and acl = "hidel"')
                acl_data = curs.fetchall()
                if(acl_data):
                    conn.close()
                    return(1)
                else:
                    reset = True
            else:
                curs.execute('select name from alist where name = "' + db_pas(user[0]["acl"]) + '" and acl = "owner"')
                acl_data = curs.fetchall()
                if(acl_data):
                    conn.close()
                    return(1)
                else:
                    break
    conn.close()
                
def include_check(name, data):
    conn = pymysql.connect(user = set_data['user'], password = set_data['pw'], charset = 'utf8mb4', db = set_data['db'])
    curs = conn.cursor(pymysql.cursors.DictCursor)

    if(re.search('^틀:', name)):
        curs.execute("select * from back where title = '" + db_pas(name) + "' and type = 'include'")
        back = curs.fetchall()
        if(back):
            i = 0

            while(True):
                try:
                    namumark(back[i]['link'], data)
                except:
                    break
                    
                i += 1
    conn.close()
    
def login_check():
    session = request.environ.get('beaker.session')
    if(session.get('Now') == True):
        return(1)
    else:
        return(0)

def ip_pas(raw_ip, num):
    conn = pymysql.connect(user = set_data['user'], password = set_data['pw'], charset = 'utf8mb4', db = set_data['db'])
    curs = conn.cursor(pymysql.cursors.DictCursor)

    if(re.search("\.", raw_ip)):
        ip = raw_ip
    else:
        curs.execute("select title from data where title = '사용자:" + db_pas(raw_ip) + "'")
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

def ip_check():
    conn = pymysql.connect(user = set_data['user'], password = set_data['pw'], charset = 'utf8mb4', db = set_data['db'])
    curs = conn.cursor(pymysql.cursors.DictCursor)

    session = request.environ.get('beaker.session')
    if(session.get('Now') == True):
        ip = format(session['DREAMER'])
    else:
        if(request.environ.get('HTTP_X_FORWARDED_FOR')):
            ip = request.environ.get('HTTP_X_FORWARDED_FOR')
        else:
            ip = request.environ.get('REMOTE_ADDR')

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
    conn = pymysql.connect(user = set_data['user'], password = set_data['pw'], charset = 'utf8mb4', db = set_data['db'])
    curs = conn.cursor(pymysql.cursors.DictCursor)

    m = re.search("^사용자:(.*)", name)
    n = re.search("^파일:(.*)", name)
    if(m):
        g = m.groups()
        if(ip == g[0]):
            if(re.search("\.", g[0])):
                conn.close()
                return(1)
            else:
                curs.execute("select * from ban where block = '" + db_pas(ip) + "'")
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
        if(not admin_check(None) == 1):
            conn.close()
            return(1)
    else:
        b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
        if(b):
            results = b.groups()
            curs.execute("select * from ban where block = '" + db_pas(results[0]) + "' and band = 'O'")
            rowss = curs.fetchall()
            if(rowss):
                conn.close()
                return(1)
            else:
                curs.execute("select * from ban where block = '" + db_pas(ip) + "'")
                rows = curs.fetchall()
                if(rows):
                    conn.close()
                    return(1)
                else:
                    curs.execute("select acl from data where title = '" + db_pas(name) + "'")
                    row = curs.fetchall()
                    if(row):
                        curs.execute("select * from user where id = '" + db_pas(ip) + "'")
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
            curs.execute("select * from ban where block = '" + db_pas(ip) + "'")
            rows = curs.fetchall()
            if(rows):
                conn.close()
                return(1)
            else:
                curs.execute("select acl from data where title = '" + db_pas(name) + "'")
                row = curs.fetchall()
                if(row):
                    curs.execute("select * from user where id = '" + db_pas(ip) + "'")
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
    conn = pymysql.connect(user = set_data['user'], password = set_data['pw'], charset = 'utf8mb4', db = set_data['db'])
    curs = conn.cursor(pymysql.cursors.DictCursor)

    b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
    if(b):
        results = b.groups()
        curs.execute("select * from ban where block = '" + db_pas(results[0]) + "' and band = 'O'")
        rowss = curs.fetchall()
        if(rowss):
            conn.close()
            return(1)
        else:
            curs.execute("select * from ban where block = '" + db_pas(ip) + "'")
            rows = curs.fetchall()
            if(rows):
                conn.close()
                return(1)
            else:
                conn.close()
                return(0)
    else:
        curs.execute("select * from ban where block = '" + db_pas(ip) + "'")
        rows = curs.fetchall()
        if(rows):
            conn.close()
            return(1)
        else:
            conn.close()
            return(0)
    conn.close()
        
def topic_check(ip, name, sub):
    conn = pymysql.connect(user = set_data['user'], password = set_data['pw'], charset = 'utf8mb4', db = set_data['db'])
    curs = conn.cursor(pymysql.cursors.DictCursor)

    b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
    if(b):
        results = b.groups()
        curs.execute("select * from ban where block = '" + db_pas(results[0]) + "' and band = 'O'")
        rowss = curs.fetchall()
        if(rowss):
            conn.close()
            return(1)
        else:
            curs.execute("select * from ban where block = '" + db_pas(ip) + "'")
            rows = curs.fetchall()
            if(rows):
                conn.close()
                return(1)
            else:
                curs.execute("select * from stop where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "'")
                rows = curs.fetchall()
                if(rows):
                    conn.close()
                    return(1)
                else:
                    conn.close()
                    return(0)
    else:
        curs.execute("select * from ban where block = '" + db_pas(ip) + "'")
        rows = curs.fetchall()
        if(rows):
            conn.close()
            return(1)
        else:
            curs.execute("select * from stop where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "'")
            rows = curs.fetchall()
            if(rows):
                conn.close()
                return(1)
            else:
                conn.close()
                return(0)
    conn.close()

def get_time():
    conn = pymysql.connect(user = set_data['user'], password = set_data['pw'], charset = 'utf8mb4', db = set_data['db'])
    curs = conn.cursor(pymysql.cursors.DictCursor)

    now = time.localtime()
    date = "%04d-%02d-%02d %02d:%02d:%02d" % (now.tm_year, now.tm_mon, now.tm_mday, now.tm_hour, now.tm_min, now.tm_sec)
    
    conn.close()
    return(date)

def rd_plus(title, sub, date):
    conn = pymysql.connect(user = set_data['user'], password = set_data['pw'], charset = 'utf8mb4', db = set_data['db'])
    curs = conn.cursor(pymysql.cursors.DictCursor)

    curs.execute("select * from rd where title = '" + db_pas(title) + "' and sub = '" + db_pas(sub) + "'")
    rd = curs.fetchall()
    if(rd):
        curs.execute("update rd set date = '" + db_pas(date) + "' where title = '" + db_pas(title) + "' and sub = '" + db_pas(sub) + "'")
    else:
        curs.execute("insert into rd (title, sub, date) value ('" + db_pas(title) + "', '" + db_pas(sub) + "', '" + db_pas(date) + "')")
    conn.commit()
    conn.close()
    
def rb_plus(block, end, today, blocker, why):
    conn = pymysql.connect(user = set_data['user'], password = set_data['pw'], charset = 'utf8mb4', db = set_data['db'])
    curs = conn.cursor(pymysql.cursors.DictCursor)

    curs.execute("insert into rb (block, end, today, blocker, why) value ('" + db_pas(block) + "', '" + db_pas(end) + "', '" + today + "', '" + db_pas(blocker) + "', '" + db_pas(why) + "')")
    conn.commit()
    conn.close()

def history_plus(title, data, date, ip, send, leng):
    conn = pymysql.connect(user = set_data['user'], password = set_data['pw'], charset = 'utf8mb4', db = set_data['db'])
    curs = conn.cursor(pymysql.cursors.DictCursor)

    curs.execute("select * from history where title = '" + db_pas(title) + "' order by id+0 desc limit 1")
    rows = curs.fetchall()
    if(rows):
        number = int(rows[0]['id']) + 1
        curs.execute("insert into history (id, title, data, date, ip, send, leng) value ('" + str(number) + "', '" + db_pas(title) + "', '" + db_pas(data) + "', '" + date + "', '" + db_pas(ip) + "', '" + db_pas(send) + "', '" + leng + "')")
    else:
        curs.execute("insert into history (id, title, data, date, ip, send, leng) value ('1', '" + db_pas(title) + "', '" + db_pas(data) + "', '" + date + "', '" + db_pas(ip) + "', '" + db_pas(send + ' (새 문서)') + "', '" + leng + "')")
    conn.commit()
    conn.close()

def leng_check(a, b):
    conn = pymysql.connect(user = set_data['user'], password = set_data['pw'], charset = 'utf8mb4', db = set_data['db'])
    curs = conn.cursor(pymysql.cursors.DictCursor)

    if(a < b):
        c = b - a
        c = '+' + str(c)
    elif(b < a):
        c = a - b
        c = '-' + str(c)
    else:
        c = '0'
        
    conn.close()
    return(c)