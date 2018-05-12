# 모듈들 불러옴
from css_html_js_minify import html_minify, js_minify, css_minify
from flask import session, render_template
from urllib import parse

import json
import sqlite3
import hashlib
import requests
import re
import html
import time
import os

# 일부 툴 불러옴
from set_mark.tool import get_time
from set_mark.tool import ip_check
from set_mark.tool import url_pas
from set_mark.tool import sha224

# 나무마크 불러옴
from mark import *

# 서브 언어팩 불러옴
json_data = open(os.path.join('language', 'en-US.json'), 'rt', encoding='utf-8').read()
else_lang = json.loads(json_data)

def load_conn(data):
    global conn
    global curs

    conn = data
    curs = conn.cursor()

    load_conn2(data)

def captcha_get():
    data = ''

    if custom()[2] == 0:
        curs.execute('select data from other where name = "recaptcha"')
        recaptcha = curs.fetchall()
        if recaptcha and recaptcha[0][0] != '':
            curs.execute('select data from other where name = "sec_re"')
            sec_re = curs.fetchall()
            if sec_re and sec_re[0][0] != '':
                data += recaptcha[0][0] + '<hr>'

    return data

def captcha_post(test, num = 1):
    if num == 1:
        if custom()[2] == 0 and captcha_get() != '':
            curs.execute('select data from other where name = "sec_re"')
            sec_re = curs.fetchall()
            if sec_re and sec_re[0][0] != '':
                data = requests.get('https://www.google.com/recaptcha/api/siteverify', params = { 'secret' : sec_re, 'response' : test })
                if not data:
                    return 0
                else:
                    json_data = data.json()
                    if data.status_code == 200 and json_data['success'] == True:
                        return 0
                    else:
                        return 1
            else:
                return 0
        else:
            return 0
    else:
        pass

def load_lang(data):
    global lang

    try:
        if lang:
            pass
    except:
        curs.execute("select data from other where name = 'language'")
        rep_data = curs.fetchall()

        json_data = open(os.path.join('language', rep_data[0][0] + '.json'), 'rt', encoding='utf-8').read()
        lang = json.loads(json_data)

    if data == 'please_all':
        return lang
    else:
        if data in lang:
            return lang[data]
        else:
            return else_lang[data]

def edit_help_button():
    # https://stackoverflow.com/questions/11076975/insert-text-into-textarea-at-cursor-position-javascript
    '''<script>
                function insertAtCursor(myField, myValue) {
                    if (document.selection) { 
                        document.getElementById(myField).focus();
                        sel = document.selection.createRange(); 
                        sel.text = myValue; 
                    } else if (document.getElementById(myField).selectionStart || document.getElementById(myField).selectionStart == '0') { 
                        var startPos = document.getElementById(myField).selectionStart; 
                        var endPos = document.getElementById(myField).selectionEnd; 
                        document.getElementById(myField).value = document.getElementById(myField).value.substring(0, startPos) + myValue + document.getElementById(myField).value.substring(endPos, document.getElementById(myField).value.length); 
                    } else { 
                        document.getElementById(myField).value += myValue;
                    }
                }
            </script>
        '''

    '<a href="javascript:void(0);" onclick="insertAtCursor(\'content\', \'[[]]\');">(링크)</a> <a href="javascript:void(0);" onclick="insertAtCursor(\'content\', \'[macro()]\');">(매크로)</a> <a href="javascript:void(0);" onclick="insertAtCursor(\'content\', \'{{{#! }}}\');">(중괄호)</a><hr>'

    return ['', '']

def ip_warring():
    if custom()[2] == 0:    
        curs.execute('select data from other where name = "no_login_warring"')
        data = curs.fetchall()
        if data and data[0][0] != '':
            text_data = '<span>' + data[0][0] + '</span><hr>'
        else:
            text_data = '<span>' + load_lang('no_login_warring') + '</span><hr>'
    else:
        text_data = ''

    return text_data

def skin_check():
    skin = './views/acme/'
    
    try:
        curs.execute('select skin from user where id = ?', [ip_check()])
        skin_exist = curs.fetchall()
        if skin_exist and skin_exist[0][0] != '':
            if os.path.exists(os.path.abspath('./views/' + skin_exist[0][0] + '/index.html')) == 1:
                skin = './views/' + skin_exist[0][0] + '/'
        else:
            curs.execute('select data from other where name = "skin"')
            skin_exist = curs.fetchall()
            if skin_exist:
                if os.path.exists(os.path.abspath('./views/' + skin_exist[0][0] + '/index.html')) == 1:
                    skin = './views/' + skin_exist[0][0] + '/'
    except:
        pass

    return skin + 'index.html'

def next_fix(link, num, page, end = 50):
    list_data = ''

    if num == 1:
        if len(page) == end:
            list_data += '<hr><a href="' + link + str(num + 1) + '">(' + load_lang('next') + ')</a>'
    elif len(page) != end:
        list_data += '<hr><a href="' + link + str(num - 1) + '">(' + load_lang('previous') + ')</a>'
    else:
        list_data += '<hr><a href="' + link + str(num - 1) + '">(' + load_lang('previous') + ')</a> <a href="' + link + str(num + 1) + '">(' + load_lang('next') + ')</a>'

    return list_data

def other2(origin):
    return origin + ['Deleted', load_lang('please_all')]

def wiki_set(num):
    if num == 1:
        data_list = []

        curs.execute('select data from other where name = ?', ['name'])
        db_data = curs.fetchall()
        if db_data and db_data[0][0] != '':
            data_list += [db_data[0][0]]
        else:
            data_list += ['Wiki']

        curs.execute('select data from other where name = "license"')
        db_data = curs.fetchall()
        if db_data and db_data[0][0] != '':
            data_list += [db_data[0][0]]
        else:
            data_list += ['CC 0']

        data_list += ['', '']

        curs.execute('select data from other where name = "logo"')
        db_data = curs.fetchall()
        if db_data and db_data[0][0] != '':
            data_list += [db_data[0][0]]
        else:
            data_list += [data_list[0]]
            
        curs.execute("select data from other where name = 'head'")
        db_data = curs.fetchall()
        if db_data and db_data[0][0] != '':
            data_list += [db_data[0][0]]
        else:
            data_list += ['']

        return data_list

    if num == 2:
        var_data = 'FrontPage'

        curs.execute('select data from other where name = "frontpage"')
    elif num == 3:
        var_data = '2'

        curs.execute('select data from other where name = "upload"')
    
    db_data = curs.fetchall()
    if db_data and db_data[0][0] != '':
        return db_data[0][0]
    else:
        return var_data

def diff(seqm):
    output = []

    for opcode, a0, a1, b0, b1 in seqm.get_opcodes():
        if opcode == 'equal':
            output += [seqm.a[a0:a1]]
        elif opcode == 'insert':
            output += ["<span style='background:#CFC;'>" + seqm.b[b0:b1] + "</span>"]
        elif opcode == 'delete':
            output += ["<span style='background:#FDD;'>" + seqm.a[a0:a1] + "</span>"]
        elif opcode == 'replace':
            output += ["<span style='background:#FDD;'>" + seqm.a[a0:a1] + "</span>"]
            output += ["<span style='background:#CFC;'>" + seqm.b[b0:b1] + "</span>"]
            
    return ''.join(output)
           
def admin_check(num, what):
    ip = ip_check() 

    curs.execute("select acl from user where id = ?", [ip])
    user = curs.fetchall()
    if user:
        reset = 0

        while 1:
            if num == 1 and reset == 0:
                check = 'ban'
            elif num == 2 and reset == 0:
                check = 'mdel'
            elif num == 3 and reset == 0:
                check = 'toron'
            elif num == 4 and reset == 0:
                check = 'check'
            elif num == 5 and reset == 0:
                check = 'acl'
            elif num == 6 and reset == 0:
                check = 'hidel'
            elif num == 7 and reset == 0:
                check = 'give'
            else:
                check = 'owner'

            curs.execute('select name from alist where name = ? and acl = ?', [user[0][0], check])
            if curs.fetchall():
                if what:
                    curs.execute("insert into re_admin (who, what, time) values (?, ?, ?)", [ip, what, get_time()])
                    conn.commit()

                return 1
            else:
                if reset == 0:
                    reset = 1
                else:
                    break

def ip_pas(raw_ip):
    hide = 0

    if re.search("(\.|:)", raw_ip):
        if not re.search("^" + load_lang('tool') + ":", raw_ip):    
            curs.execute("select data from other where name = 'ip_view'")
            data = curs.fetchall()
            if data and data[0][0] != '':
                ip = '<span style="font-size: 75%;">' + hashlib.md5(bytes(raw_ip, 'utf-8')).hexdigest() + '</span>'

                if not admin_check('ban', None):
                    hide = 1
            else:
                ip = raw_ip
        else:
            ip = raw_ip
            hide = 1
    else:
        curs.execute("select title from data where title = ?", ['' + load_lang('user') + ':' + raw_ip])
        if curs.fetchall():
            ip = '<a href="/w/' + url_pas('' + load_lang('user') + ':' + raw_ip) + '">' + raw_ip + '</a>'
        else:
            ip = '<a id="not_thing" href="/w/' + url_pas('' + load_lang('user') + ':' + raw_ip) + '">' + raw_ip + '</a>'
         
    if hide == 0:
        ip += ' <a href="/record/' + url_pas(raw_ip) + '">(' + load_lang('record') + ')</a>'

    return ip

def custom():
    if 'MyMaiToNight' in session:
        user_head = session['MyMaiToNight']
    else:
        user_head = ''

    if 'Now' in session and session['Now'] == 1:
        curs.execute('select name from alarm where name = ? limit 1', [ip_check()])
        if curs.fetchall():
            user_icon = 2
        else:
            user_icon = 1
    else:
        user_icon = 0

    if user_icon != 0:
        curs.execute('select email from user where id = ?', [ip_check()])
        data = curs.fetchall()
        if data:
            email = data[0][0]
        else:
            email = ''
    else:
        email = ''

    if user_icon != 0:
        user_name = ip_check()
    else:
        user_name = '' + load_lang('user') + ''

    return ['', '', user_icon, user_head, email, user_name]

def acl_check(name):
    ip = ip_check()

    if ban_check() == 1:
        return 1

    acl_c = re.search("^" + load_lang('user') + ":([^/]*)", name)
    if acl_c:
        acl_n = acl_c.groups()

        if admin_check(5, None) == 1:
            return 0

        curs.execute("select dec from acl where title = ?", ['' + load_lang('user') + ':' + acl_n[0]])
        acl_data = curs.fetchall()
        if acl_data:
            if acl_data[0][0] == 'all':
                return 0

            if acl_data[0][0] == 'user' and not re.search("(\.|:)", ip):
                return 0

            if ip != acl_n[0] or re.search("(\.|:)", ip):
                return 1
        
        if ip == acl_n[0] and not re.search("(\.|:)", ip) and not re.search("(\.|:)", acl_n[0]):
            return 0
        else:
            return 1

    file_c = re.search("^" + load_lang('file') + ":(.*)", name)
    if file_c and admin_check(5, 'edit (' + name + ')') != 1:
        return 1

    curs.execute("select acl from user where id = ?", [ip])
    user_data = curs.fetchall()

    curs.execute("select dec from acl where title = ?", [name])
    acl_data = curs.fetchall()
    if acl_data:
        if acl_data[0][0] == 'user':
            if not user_data:
                return 1

        if acl_data[0][0] == 'admin':
            if not user_data:
                return 1

            if not admin_check(5, 'edit (' + name + ')') == 1:
                return 1

    curs.execute('select data from other where name = "edit"')
    set_data = curs.fetchall()
    if set_data:
        if set_data[0][0] == 'user':
            if not user_data:
                return 1

        if set_data[0][0] == 'admin':
            if not user_data:
                return 1

            if not admin_check(5, None) == 1:
                return 1

    return 0

def ban_check():
    ip = ip_check()

    band = re.search("^([0-9]{1,3}\.[0-9]{1,3})", ip)
    if band:
        band_it = band.groups()[0]
    else:
        band_it = 'Not'
        
    curs.execute("select block from ban where block = ?", [band_it])
    band_d = curs.fetchall()
    
    curs.execute("select block from ban where block = ?", [ip])
    ban_d = curs.fetchall()
    if band_d or ban_d:
        return 1
    
    return 0
        
def topic_check(name, sub):
    ip = ip_check()

    if ban_check() == 1:
        return 1
        
    curs.execute("select acl from user where id = ?", [ip])
    user_data = curs.fetchall()

    curs.execute("select dis from acl where title = ?", [name])
    acl_data = curs.fetchall()
    if acl_data:
        if acl_data[0][0] == 'user':
            if not user_data:
                return 1

        if acl_data[0][0] == 'admin':
            if not user_data:
                return 1

            if not admin_check(3, 'topic (' + name + ')') == 1:
                return 1
        
    curs.execute("select title from stop where title = ? and sub = ?", [name, sub])
    if curs.fetchall():
        if not admin_check(3, 'topic (' + name + ')') == 1:
            return 1

    return 0

def ban_insert(name, end, why, login, blocker):
    time = get_time()

    if re.search("^([0-9]{1,3}\.[0-9]{1,3})$", name):
        band = 'O'
    else:
        band = ''

    curs.execute("select block from ban where block = ?", [name])
    if curs.fetchall():
        curs.execute("insert into rb (block, end, today, blocker, why, band) values (?, ?, ?, ?, ?, ?)", [name, '' + load_lang('release') + '', time, blocker, '', band])
        curs.execute("delete from ban where block = ?", [name])
    else:
        if login != '':
            login = 'O'
        else:
            login = ''

        if end != '':
            end += ' 00:00:00'

        curs.execute("insert into rb (block, end, today, blocker, why, band) values (?, ?, ?, ?, ?, ?)", [name, end, time, blocker, why, band])
        curs.execute("insert into ban (block, end, why, band, login) values (?, ?, ?, ?, ?)", [name, end, why, band, login])
    
    conn.commit()

def rd_plus(title, sub, date):
    curs.execute("select title from rd where title = ? and sub = ?", [title, sub])
    if curs.fetchall():
        curs.execute("update rd set date = ? where title = ? and sub = ?", [date, title, sub])
    else:
        curs.execute("insert into rd (title, sub, date) values (?, ?, ?)", [title, sub, date])

def history_plus(title, data, date, ip, send, leng):
    curs.execute("select id from history where title = ? order by id + 0 desc limit 1", [title])
    id_data = curs.fetchall()
    if id_data:
        curs.execute("insert into history (id, title, data, date, ip, send, leng) values (?, ?, ?, ?, ?, ?, ?)", [str(int(id_data[0][0]) + 1), title, data, date, ip, send, leng])
    else:
        curs.execute("insert into history (id, title, data, date, ip, send, leng) values ('1', ?, ?, ?, ?, ?, ?)", [title, data, date, ip, send + ' (' + load_lang('new') + ' ' + load_lang('document') + ')', leng])

def leng_check(first, second):
    if first < second:
        all_plus = '+' + str(second - first)
    elif second < first:
        all_plus = '-' + str(first - second)
    else:
        all_plus = '0'
        
    return all_plus

def redirect(data):
    return '<meta http-equiv="refresh" content="0; url=' + data + '">'

def re_error(data):
    if data == '/ban':
        ip = ip_check()

        end = '<li>Why : 권한이 맞지 않는 상태 입니다.</li>'

        if ban_check() == 1:
            curs.execute("select end, why from ban where block = ?", [ip])
            end_data = curs.fetchall()
            if not end_data:
                match = re.search("^([0-9]{1,3}\.[0-9]{1,3})", ip)
                if match:
                    curs.execute("select end, why from ban where block = ?", [match.groups()[0]])
                    end_data = curs.fetchall()
            
            if end_data:
                end = '<li>Info : '

                if end_data[0][0]:
                    now = int(re.sub('(\-| |:)', '', get_time()))
                    day = int(re.sub('(\-| |:)', '', end_data[0][0]))
                    
                    if now >= day:
                        curs.execute("delete from ban where block = ?", [ip])
                        conn.commit()

                        end += 'Re Try.'
                    else:
                        end += 'Ban : ' + end_data[0][0]
                else:
                    end += 'Ban : No End'
                
                end += '</li>'

                if end_data[0][1] != '':
                    end += '<li>Why : ' + end_data[0][1] + '</li>'

        return html_minify(render_template(skin_check(), 
            imp = ['Authority Error', wiki_set(1), custom(), other2([0, 0])],
            data = '<h2>Info</h2><ul>' + end + '</ul>',
            menu = 0
        ))

    error_data = re.search('\/error\/([0-9]+)', data)
    if error_data:
        num = int(error_data.groups()[0])
        if num == 1:
            title = 'Authority Error'
            data = '비 로그인 상태 입니다.'
        elif num == 2:
            title = 'Authority Error'
            data = '이 계정이 없습니다.'
        elif num == 3:
            title = 'Authority Error'
            data = '권한이 모자랍니다.'
        elif num == 4:
            title = 'Authority Error'
            data = '관리자는 차단, 검사 할 수 없습니다.'
        elif num == 5:
            title = 'User Error'
            data = '그런 계정이 없습니다.'
        elif num == 6:
            title = 'Register Error'
            data = '동일한 아이디의 사용자가 있습니다.'
        elif num == 7:
            title = 'Register Error'
            data = '아이디는 20글자보다 짧아야 합니다.'
        elif num == 8:
            title = 'Register Error'
            data = '아이디에는 한글과 알파벳과 공백만 허용 됩니다.'
        elif num == 9:
            title = 'Upload Error'
            data = '파일이 없습니다.'
        elif num == 10:
            title = 'PassWord Error'
            data = '비밀번호가 다릅니다.'
        elif num == 11:
            title = 'Login Error'
            data = '이미 로그인 되어 있습니다.'
        elif num == 13:
            title = 'reCAPTCHA Error'
            data = '리캡차를 통과하세요.'
        elif num == 14:
            title = 'Upload Error'
            data = 'jpg, gif, jpeg, png, webp만 가능 합니다.'
        elif num == 15:
            title = 'Edit Error'
            data = '편집 기록은 500자를 넘을 수 없습니다.'
        elif num == 16:
            title = 'Upload Error'
            data = '동일한 이름의 파일이 있습니다.'
        elif num == 17:
            title = 'Upload Error'
            data = '파일 용량은 ' + wiki_set(3) + 'MB를 넘길 수 없습니다.'
        elif num == 18:
            title = 'Edit Error'
            data = '내용이 원래 문서와 동일 합니다.'
        elif num == 19:
            title = 'Move Error'
            data = '이동 하려는 곳에 문서가 이미 있습니다.'
        elif num == 20:
            title = 'PassWord Error'
            data = '재 확인이랑 비밀번호가 다릅니다.'
        elif num == 21:
            title = 'Edit Error'
            data = '편집 필터에 의해 검열 되었습니다.'
        elif num == 22:
            title = 'Upload Error'
            data = '파일 이름은 알파벳, 한글, 띄어쓰기, 언더바, 빼기표만 허용 됩니다.'
        else:
            title = 'Error'
            data = '???'

        if title:
            return html_minify(render_template(skin_check(), 
                imp = [title, wiki_set(1), custom(), other2([0, 0])],
                data = '<h2>Error</h2><ul><li>' + data + '</li></ul>',
                menu = 0
            ))
        else:
            return redirect('/')
    else:
        return redirect('/')