from bottle import request, app, template
from bottle.ext import beaker
import json
import sqlite3
import hashlib
from urllib import parse
import re
import html
from css_html_js_minify import html_minify, js_minify, css_minify
import time
import os

from set_mark.macro import get_time
from set_mark.macro import ip_check
from set_mark.mark import *
from set_mark.link import url_pas
from set_mark.link import sha224
    
session_opts = {
    'session.type': 'dbm',
    'session.data_dir': './app_session/',
    'session.auto': 1
}

app = beaker.middleware.SessionMiddleware(app(), session_opts)

def captcha_get(conn):
    session = request.environ.get('beaker.session')
    curs = conn.cursor()

    data = ''
    if(re.search('\.|:', ip_check()) and session.get('Awaken') and session.get('Awaken') != 1):
        curs.execute('select data from other where name = "recaptcha"')
        recaptcha = curs.fetchall()
        if(recaptcha and recaptcha[0][0] != ''):
            data += recaptcha[0][0] + '<br>'

    return(data)

def captcha_post(conn, num = 1):
    session = request.environ.get('beaker.session')
    if(num == 1):
        if(re.search('\.|:', ip_check()) and session.get('Awaken') and session.get('Awaken') != 1 and captcha_get(conn) != ''):
            return(1)
        else:
            return(0)
    else:
        session['Awaken'] = 1

def skin_check(conn):
    curs = conn.cursor()
    skin = './views/acme/'
    try:
        curs.execute('select data from other where name = "skin"')
        skin_exist = curs.fetchall()
        if(skin_exist):
            if(os.path.exists(os.path.abspath('./views/' + skin_exist[0][0] + '/index.tpl')) == 1):
                skin = './views/' + skin_exist[0][0] + '/'
    except:
        pass

    return(skin)

def other2(origin):
    div = ''
    session = request.environ.get('beaker.session')
    if(session.get('View_List')):
        match = re.findall('(?:(?:([^\n]+)\n))', session.get('View_List'))
        if(match):
            div = ''
            for data in match[-6:-1]:
                div += '<a href="/w/' + url_pas(data) + '">' +  html.escape(data) + '</a> / '
            div = re.sub(' / $', '', div)
            
    re_data = origin + [div]
    return(re_data)    

def wiki_set(conn, num):
    curs = conn.cursor()
    if(num == 1):
        data_list = []

        curs.execute('select data from other where name = ?', ['name'])
        db_data = curs.fetchall()
        if(db_data and db_data[0][0] != ''):
            data_list += [db_data[0][0]]
        else:
            data_list += ['무명위키']

        curs.execute('select data from other where name = "license"')
        db_data = curs.fetchall()
        if(db_data and db_data[0][0] != ''):
            data_list += [db_data[0][0]]
        else:
            data_list += ['CC 0']

        data_list += ['', '']

        curs.execute('select data from other where name = "logo"')
        db_data = curs.fetchall()
        if(db_data and db_data[0][0] != ''):
            data_list += [db_data[0][0]]
        else:
            data_list += [data_list[0]]
            
        curs.execute("select data from other where name = 'head'")
        db_data = curs.fetchall()
        if(db_data and db_data[0][0] != ''):
            data_list += [db_data[0][0]]
        else:
            data_list += ['']

        return(data_list)

    if(num == 2):
        var_data = '위키:대문'
        curs.execute('select data from other where name = "frontpage"')
    elif(num == 3):
        var_data = '2'
        curs.execute('select data from other where name = "upload"')
    
    db_data = curs.fetchall()
    if(db_data and db_data[0][0] != ''):
        return(db_data[0][0])
    else:
        return(var_data)

def diff(seqm):
    output = []
    for opcode, a0, a1, b0, b1 in seqm.get_opcodes():
        if(opcode == 'equal'):
            output += [seqm.a[a0:a1]]
        elif(opcode == 'insert'):
            output += ["<span style='background:#CFC;'>" + seqm.b[b0:b1] + "</span>"]
        elif(opcode == 'delete'):
            output += ["<span style='background:#FDD;'>" + seqm.a[a0:a1] + "</span>"]
        elif(opcode == 'replace'):
            output += ["<span style='background:#FDD;'>" + seqm.a[a0:a1] + "</span>"]
            output += ["<span style='background:#CFC;'>" + seqm.b[b0:b1] + "</span>"]
            
    return(''.join(output))
           
def admin_check(conn, num, what):
    ip = ip_check() 
    curs = conn.cursor()
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
            elif(num == 7 and reset == 0):
                check = 'give'
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

def ip_pas(conn, raw_ip):
    hide = 0
    curs = conn.cursor()
    if(re.search("(\.|:)", raw_ip)):
        curs.execute("select data from other where name = 'ip_view'")
        d = curs.fetchall()
        if(d and d[0][0] != ''):
            ip = '<span style="font-size: 75%;">' + hashlib.md5(bytes(raw_ip, 'utf-8')).hexdigest() + '</span>'
            if(not admin_check(conn, 'ban', None)):
                hide = 1
        else:
            ip = raw_ip
    else:
        curs.execute("select title from data where title = ?", ['사용자:' + raw_ip])
        if(curs.fetchall()):
            ip = '<a href="/w/' + url_pas('사용자:' + raw_ip) + '">' + raw_ip + '</a>'
        else:
            ip = '<a class="not_thing" href="/w/' + url_pas('사용자:' + raw_ip) + '">' + raw_ip + '</a>'
         
    if(hide == 0):
        ip += ' <a href="/record/' + url_pas(raw_ip) + '">(기록)</a>'

    return(ip)

def custom(conn):
    curs = conn.cursor()
    session = request.environ.get('beaker.session')
    try:
        user_head = session['MyMaiToNight']
    except:
        user_head = ''

    if(session.get('Now') == 1):
        curs.execute('select name from alarm limit 1')
        if(curs.fetchall()):
            user_icon = 2
        else:
            user_icon = 1
    else:
        user_icon = 0

    return(['', '', user_icon, user_head])

def acl_check(conn, name):
    ip = ip_check()
    curs = conn.cursor()
    if(ban_check(conn) == 1):
        return(1)

    acl_c = re.search("^사용자:([^/]*)", name)
    if(acl_c):
        acl_n = acl_c.groups()

        if(admin_check(conn, 5, None) == 1):
            return(0)

        curs.execute("select acl from data where title = ?", ['사용자:' + acl_n[0]])
        acl_d = curs.fetchall()
        if(acl_d):
            if(acl_d[0][0] == 'all'):
                return(0)

            if(acl_d[0][0] == 'user' and not re.search("(\.|:)", ip)):
                return(0)

            if(not ip == acl_n[0] or re.search("(\.|:)", ip)):
                return(1)
        
        if(ip == acl_n[0] and not re.search("(\.|:)", ip) and not re.search("(\.|:)", acl_n[0])):
            return(0)
        else:
            return(1)

    file_c = re.search("^파일:(.*)", name)
    if(file_c and admin_check(conn, 5, 'edit (' + name + ')') != 1):
        return(1)

    curs.execute("select acl from user where id = ?", [ip])
    user_d = curs.fetchall()

    curs.execute("select acl from data where title = ?", [name])
    acl_d = curs.fetchall()
    if(acl_d):
        if(acl_d[0][0] == 'user'):
            if(not user_d):
                return(1)

        if(acl_d[0][0] == 'admin'):
            if(not user_d):
                return(1)

            if(not admin_check(conn, 5, 'edit (' + name + ')') == 1):
                return(1)

    curs.execute('select data from other where name = "edit"')
    set_d = curs.fetchall()
    if(set_d):
        if(set_d[0][0] == 'user'):
            if(not user_d):
                return(1)

        if(set_d[0][0] == 'admin'):
            if(not user_d):
                return(1)

            if(not admin_check(conn, 5, None) == 1):
                return(1)

    return(0)

def ban_check(conn):
    ip = ip_check()
    curs = conn.cursor()
    band = re.search("^([0-9]{1,3}\.[0-9]{1,3})", ip)
    if(band):
        band_it = band.groups()
    else:
        band_it = ['Not']
        
    curs.execute("select block from ban where block = ? and band = 'O'", [band_it[0]])
    band_d = curs.fetchall()

    curs.execute("select block from ban where block = ?", [ip])
    ban_d = curs.fetchall()
    if(band_d or ban_d):
        return(1)
    
    return(0)
        
def topic_check(conn, name, sub):
    ip = ip_check()
    curs = conn.cursor()
    if(ban_check(conn) == 1):
        return(1)

    curs.execute("select title from stop where title = ? and sub = ?", [name, sub])
    topic_s = curs.fetchall()
    if(topic_s):
        return(1)

    return(0)

def rd_plus(conn, title, sub, date):
    curs = conn.cursor()
    curs.execute("select title from rd where title = ? and sub = ?", [title, sub])
    if(curs.fetchall()):
        curs.execute("update rd set date = ? where title = ? and sub = ?", [date, title, sub])
    else:
        curs.execute("insert into rd (title, sub, date) values (?, ?, ?)", [title, sub, date])
    
def rb_plus(conn, block, end, today, blocker, why):
    curs = conn.cursor()
    curs.execute("insert into rb (block, end, today, blocker, why) values (?, ?, ?, ?, ?)", [block, end, today, blocker, why])

def history_plus(conn, title, data, date, ip, send, leng):
    curs = conn.cursor()
    curs.execute("select id from history where title = ? order by id + 0 desc limit 1", [title])
    d = curs.fetchall()
    if(d):
        curs.execute("insert into history (id, title, data, date, ip, send, leng) values (?, ?, ?, ?, ?, ?, ?)", [str(int(d[0][0]) + 1), title, data, date, ip, send, leng])
    else:
        curs.execute("insert into history (id, title, data, date, ip, send, leng) values ('1', ?, ?, ?, ?, ?, ?)", [title, data, date, ip, send + ' (새 문서)', leng])

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

def redirect(data):
    return('<meta http-equiv="refresh" content="0; url=' + data + '">')

def re_error(conn, data):
    curs = conn.cursor()
    if(data == '/ban'):
        ip = ip_check()
        end = '|| 사유 || 권한이 맞지 않는 상태 입니다. ||'
        if(ban_check(conn) == 1):
            curs.execute("select end, why from ban where block = ?", [ip])
            d = curs.fetchall()
            if(not d):
                m = re.search("^([0-9]{1,3}\.[0-9]{1,3})", ip)
                if(m):
                    curs.execute("select end, why from ban where block = ? and band = 'O'", [m.groups()[0]])
                    d = curs.fetchall()
            if(d):
                end = '|| 상태 ||'
                if(d[0][0]):
                    now = int(re.sub('(:|-| )', '', get_time()))
                    day = re.sub('\-', '', d[0][0])              
                              
                    if(re.search(':', day)):
                        day = re.sub('( |:)', '', day)
                    else:
                        day += '000000'
                    
                    if(now >= int(day)):
                        curs.execute("delete from ban where block = ?", [ip])
                        conn.commit()
                        
                        end += '차단이 풀렸습니다. 다시 시도 해 보세요.'
                    else:
                        end += d[0][0] + ' 까지 차단 상태 입니다.'
                else:
                    end += '영구 차단 상태 입니다.'
                
                end += '||'

                if(d[0][1] != ''):
                    end += '\r\n|| 사유 || ' + d[0][1] + ' ||'

        return(html_minify(template('index', 
            imp = ['권한 오류', wiki_set(conn, 1), custom(conn), other2([0, 0])],
            data = namumark(conn, "", "[목차(없음)]\r\n== 권한 상태 ==\r\n" + end, 0, 0, 0),
            menu = 0
        )))

    d = re.search('\/error\/([0-9]+)', data)
    if(d):
        num = int(d.groups()[0])
        if(num == 1):
            title = '권한 오류'
            data = '비 로그인 상태 입니다.'
        elif(num == 2):
            title = '권한 오류'
            data = '이 계정이 없습니다.'
        elif(num == 3):
            title = '권한 오류'
            data = '권한이 모자랍니다.'
        elif(num == 4):
            title = '권한 오류'
            data = '관리자는 차단, 검사 할 수 없습니다.'
        elif(num == 5):
            title = '사용자 오류'
            data = '그런 계정이 없습니다.'
        elif(num == 6):
            title = '가입 오류'
            data = '동일한 아이디의 사용자가 있습니다.'
        elif(num == 7):
            title = '가입 오류'
            data = '아이디는 20글자보다 짧아야 합니다.'
        elif(num == 8):
            title = '가입 오류'
            data = '아이디에는 한글과 알파벳과 공백만 허용 됩니다.'
        elif(num == 9):
            title = '파일 올리기 오류'
            data = '파일이 없습니다.'
        elif(num == 10):
            title = '변경 오류'
            data = '비밀번호가 다릅니다.'
        elif(num == 11):
            title = '로그인 오류'
            data = '이미 로그인 되어 있습니다.'
        elif(num == 12):
            title = '편집 오류'
            data = '누군가 먼저 편집 했습니다.'
        elif(num == 13):
            title = '리캡차 오류'
            data = '리캡차를 통과하세요.'
        elif(num == 14):
            title = '파일 올리기 오류'
            data = 'jpg, gif, jpeg, png, webp만 가능 합니다.'
        elif(num == 15):
            title = '편집 오류'
            data = '편집 기록은 500자를 넘을 수 없습니다.'
        elif(num == 16):
            title = '파일 올리기 오류'
            data = '동일한 이름의 파일이 있습니다.'
        elif(num == 17):
            title = '파일 올리기 오류'
            data = '파일 용량은 ' + wiki_set(conn, 3) + 'MB를 넘길 수 없습니다.'
        elif(num == 18):
            title = '편집 오류'
            data = '내용이 원래 문서와 동일 합니다.'
        elif(num == 19):
            title = '이동 오류'
            data = '이동 하려는 곳에 문서가 이미 있습니다.'
        elif(num == 20):
            title = '비밀번호 오류'
            data = '재 확인이랑 비밀번호가 다릅니다.'
        elif(num == 21):
            title = '편집 오류'
            data = '편집 필터에 의해 검열 되었습니다.'
        else:
            title = '정체 불명의 오류'
            data = '???'

        if(title):
            return(html_minify(template('index', 
                imp = [title, wiki_set(conn, 1), custom(conn), other2([0, 0])],
                data = namumark(conn, "", "[목차(없음)]\r\n== 오류 발생 ==\r\n" + data, 0, 0, 0),
                menu = 0
            )))
        else:
            return(redirect('/'))
    else:
        return(redirect('/'))