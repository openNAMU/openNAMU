import werkzeug.routing
import flask_compress
import flask_reggie
import tornado.ioloop
import tornado.httpserver
import tornado.wsgi
import urllib.request
import email.mime.text
import urllib.request
import sqlite3
import hashlib
import smtplib
import bcrypt
import platform
import zipfile
import difflib
import shutil
import threading
import logging
import random
import flask
import json
import html
import sys
import re
import os

try:
    import css_html_js_minify
except:
    pass

if sys.version_info < (3, 6):
    import sha3

from .set_mark.tool import *
from .mark import *

def load_conn(data):
    global conn
    global curs

    conn = data
    curs = conn.cursor()

    load_conn2(data)

def send_email(who, title, data):
    smtp = smtplib.SMTP_SSL('smtp.gmail.com', 465)

    try:
        curs.execute('select name, data from other where name = "g_email" or name = "g_pass"')
        rep_data = curs.fetchall()
        if rep_data:
            g_email = ''
            g_pass = ''
            for i in rep_data:
                if i[0] == 'g_email':
                    g_email = i[1]
                else:
                    g_pass = i[1]

            smtp.login(g_email, g_pass)

        msg = email.mime.text.MIMEText(data)
        msg['Subject'] = title
        smtp.sendmail(g_email, who, msg.as_string())

        smtp.quit()
    except:
        print('Error : Email login error')

def last_change(data):
    json_address = re.sub("\.html$", ".json", skin_check())
    try:
        json_data = json.loads(open(json_address).read())
    except:
        json_data = 0

    if json_data != 0:
        for j_data in json_data:
            if "class" in json_data[j_data]:
                if "require" in json_data[j_data]:
                    re_data = re.compile("<((?:" + j_data + ")( (?:(?!>).)*)?)>")
                    s_data = re_data.findall(data)
                    for i_data in s_data:
                        e_data = 0

                        for j_i_data in json_data[j_data]["require"]:
                            re_data_2 = re.compile("( |^)" + j_i_data + " *= *[\'\"]" + json_data[j_data]["require"][j_i_data] + "[\'\"]")
                            if not re_data_2.search(i_data[1]):
                                re_data_2 = re.compile("( |^)" + j_i_data + "=" + json_data[j_data]["require"][j_i_data] + "(?: |$)")
                                if not re_data_2.search(i_data[1]):
                                    e_data = 1

                                    break

                        if e_data == 0:
                            re_data_3 = re.compile("<" + i_data[0] + ">")
                            data = re_data_3.sub("<" + i_data[0] + " class=\"" + json_data[j_data]["class"] + "\">", data)        
                else:
                    re_data = re.compile("<(?P<in>" + j_data + "(?: (?:(?!>).)*)?)>")
                    data = re_data.sub("<\g<in> class=\"" + json_data[j_data]["class"] + "\">", data)        

    return data

def easy_minify(data, tool = None):
    try:
        if not tool:
            data = css_html_js_minify.html_minify(data)
        else:
            if tool == 'css':
                data = css_html_js_minify.css_minify(data)
            elif tool == 'js':
                data = css_html_js_minify.js_minify(data)
    except:
        data = re.sub('\n +<', '\n<', data)
        data = re.sub('>(\n| )+<', '> <', data)
    
    return last_change(data)

def render_set(title = '', data = '', num = 0, s_data = 0):
    if acl_check(title, 'render') == 1:
        return 'HTTP Request 401.3'
    elif s_data == 1:
        return data
    else:
        return namumark(title, data, num)

def captcha_get():
    data = ''

    if custom()[2] == 0:
        curs.execute('select data from other where name = "recaptcha"')
        recaptcha = curs.fetchall()
        if recaptcha and recaptcha[0][0] != '':
            curs.execute('select data from other where name = "sec_re"')
            sec_re = curs.fetchall()
            if sec_re and sec_re[0][0] != '':
                data += recaptcha[0][0] + '<hr class=\"main_hr\">'

    return data

def update():
    # v3.0.8 rd, agreedis, stop 테이블 통합
    try:
        curs.execute("select title, sub, close from stop")
        for i in curs.fetchall():
            if i[2] == '':
                curs.execute("update rd set stop = 'S' where title = ? and sub = ?", [i[0], i[1]])
            else:
                curs.execute("update rd set stop = 'O' where title = ? and sub = ?", [i[0], i[1]])
    except:
        pass
        
    try:
        curs.execute("select title, sub from agreedis")
        for i in curs.fetchall():
            curs.execute("update rd set agree = 'O' where title = ? and sub = ?", [i[0], i[1]])
    except:
        pass
         
    try:
        curs.execute("drop table if exists stop")
        curs.execute("drop table if exists agreedis")
    except:
        pass

def pw_encode(data, data2 = '', type_d = ''):
    if type_d == '':
        curs.execute('select data from other where name = "encode"')
        set_data = curs.fetchall()

        type_d = set_data[0][0]

    if type_d == 'sha256':
        return hashlib.sha256(bytes(data, 'utf-8')).hexdigest()
    elif type_d == 'sha3':
        if sys.version_info < (3, 6):
            return sha3.sha3_256(bytes(data, 'utf-8')).hexdigest()
        else:
            return hashlib.sha3_256(bytes(data, 'utf-8')).hexdigest()
    else:
        if data2 != '':
            salt_data = bytes(data2, 'utf-8')
        else:
            salt_data = bcrypt.gensalt(11)
            
        return bcrypt.hashpw(bytes(data, 'utf-8'), salt_data).decode()

def pw_check(data, data2, type_d = 'no', id_d = ''):
    curs.execute('select data from other where name = "encode"')
    db_data = curs.fetchall()

    if type_d != 'no':
        if type_d == '':
            set_data = 'bcrypt'
        else:
            set_data = type_d
    else:
        set_data = db_data[0][0]
    
    while 1:
        if set_data in ['sha256', 'sha3']:
            data3 = pw_encode(data = data, type_d = set_data)
            if data3 == data2:
                re_data = 1
            else:
                re_data = 0

            break
        else:
            try:
                if pw_encode(data, data2, 'bcrypt') == data2:
                    re_data = 1
                else:
                    re_data = 0

                break
            except:
                set_data = db_data[0][0]

    if db_data[0][0] != set_data and re_data == 1 and id_d != '':
        curs.execute("update user set pw = ?, encode = ? where id = ?", [pw_encode(data), db_data[0][0], id_d])

    return re_data

def captcha_post(re_data, num = 1):
    if num == 1:
        if custom()[2] == 0 and captcha_get() != '':
            curs.execute('select data from other where name = "sec_re"')
            sec_re = curs.fetchall()
            if sec_re and sec_re[0][0] != '':
                data = urllib.request.urlopen('https://www.google.com/recaptcha/api/siteverify?secret=' + sec_re[0][0] + '&response=' + re_data)
                if not data:
                    return 0
                else:
                    json_data = data.read().decode(data.headers.get_content_charset())
                    json_data = json.loads(json_data)
                    if data.getcode() == 200 and json_data['success'] == True:
                        return 0
                    else:
                        return 1
            else:
                return 0
        else:
            return 0
    else:
        pass

def load_lang(data, num = 2, safe = 0):
    if num == 1:
        curs.execute("select data from other where name = 'language'")
        rep_data = curs.fetchall()

        json_data = open(os.path.join('language', rep_data[0][0] + '.json'), 'rt', encoding='utf-8').read()
        lang = json.loads(json_data)

        if data in lang:
            if safe == 1:
                return lang[data]
            else:
                return html.escape(lang[data])
        else:
            return html.escape(data + ' (M)')
    else:
        curs.execute('select data from user_set where name = "lang" and id = ?', [ip_check()])
        rep_data = curs.fetchall()
        if rep_data:
            try:
                json_data = open(os.path.join('language', rep_data[0][0] + '.json'), 'rt', encoding='utf-8').read()
                lang = json.loads(json_data)
            except:
                return load_lang(data, 1, safe)

            if data in lang:
                if safe == 1:
                    return lang[data]
                else:
                    return html.escape(lang[data])
            else:
                return load_lang(data, 1, safe)
        else:
            return load_lang(data, 1, safe)

def load_oauth(provider):
    oauth = json.loads(open('oauthsettings.json', encoding='utf-8').read())

    return oauth[provider]

def update_oauth(provider, target, content):
    oauth = json.loads(open('oauthsettings.json', encoding='utf-8').read())    
    oauth[provider][target] = content

    with open('oauthsettings.json', 'w', encoding='utf-8') as f:
        json.dump(oauth, f)

    return 'Done'

def ip_or_user(data):
    if re.search('(\.|:)', data):
        return 1
    else:
        return 0

def edit_help_button():
    # https://stackoverflow.com/questions/11076975/insert-text-into-textarea-at-cursor-position-javascript
    js_data =   '''
        <script>
            function insert_data(name, data) {
                if(document.selection) { 
                    document.getElementById(name).focus();

                    sel = document.selection.createRange();
                    sel.text = data; 
                } else if(document.getElementById(name).selectionStart || document.getElementById(name).selectionStart == '0') {
                    var startPos = document.getElementById(name).selectionStart;
                    var endPos = document.getElementById(name).selectionEnd;

                    document.getElementById(name).value = document.getElementById(name).value.substring(0, startPos) + data + document.getElementById(name).value.substring(endPos, document.getElementById(name).value.length); 
                } else {
                    document.getElementById(name).value += data;
                }
            }
        </script>
    '''

    insert_list = [['[[|]]', '[[|]]'], ['[*()]', '[*()]'], ['{{{#!}}}', '{{{#!}}}'], ['||<>||', '||<>||'], ["\\'\\'\\'", "\'\'\'"]]

    data = ''
    for insert_data in insert_list:
        data += '<a href="javascript:void(0);" onclick="insert_data(\'content\', \'' + insert_data[0] + '\');">(' + insert_data[1] + ')</a> '

    return [js_data, data + '<hr class=\"main_hr\">']

def ip_warring():
    if custom()[2] == 0:    
        curs.execute('select data from other where name = "no_login_warring"')
        data = curs.fetchall()
        if data and data[0][0] != '':
            text_data = '<span>' + data[0][0] + '</span><hr class=\"main_hr\">'
        else:
            text_data = '<span>' + load_lang('no_login_warring') + '</span><hr class=\"main_hr\">'
    else:
        text_data = ''

    return text_data

def skin_check():
    skin = './views/neo_yousoro/'

    curs.execute('select data from other where name = "skin"')
    skin_exist = curs.fetchall()
    if skin_exist and skin_exist[0][0] != '':
        if os.path.exists(os.path.abspath('./views/' + skin_exist[0][0] + '/index.html')) == 1:
            skin = './views/' + skin_exist[0][0] + '/'
    
    curs.execute('select data from user_set where name = "skin" and id = ?', [ip_check()])
    skin_exist = curs.fetchall()
    if skin_exist and skin_exist[0][0] != '':
        if os.path.exists(os.path.abspath('./views/' + skin_exist[0][0] + '/index.html')) == 1:
            skin = './views/' + skin_exist[0][0] + '/'

    return skin + 'index.html'

def next_fix(link, num, page, end = 50):
    list_data = ''

    if num == 1:
        if len(page) == end:
            list_data += '<hr class=\"main_hr\"><a href="' + link + str(num + 1) + '">(' + load_lang('next') + ')</a>'
    elif len(page) != end:
        list_data += '<hr class=\"main_hr\"><a href="' + link + str(num - 1) + '">(' + load_lang('previous') + ')</a>'
    else:
        list_data += '<hr class=\"main_hr\"><a href="' + link + str(num - 1) + '">(' + load_lang('previous') + ')</a> <a href="' + link + str(num + 1) + '">(' + load_lang('next') + ')</a>'

    return list_data

def other2(data):
    return data + ['']

def wiki_set(num = 1):
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

    end = ''.join(output)
    end = end.replace('\r\n', '\n')
    sub = ''

    if not re.search('\n', end):
        end += '\n'

    num = 0
    left = 1
    while 1:
        data = re.search('((?:(?!\n).)*)\n', end)
        if data:
            data = data.groups()[0]
            
            left += 1
            if re.search('<span style=\'(?:(?:(?!\').)+)\'>', data):
                num += 1
                if re.search('<\/span>', data):
                    num -= 1

                sub += str(left) + ' : ' + re.sub('(?P<in>(?:(?!\n).)*)\n', '\g<in>', data, 1) + '<br>'
            else:
                if re.search('<\/span>', data):
                    num -= 1
                    sub += str(left) + ' : ' + re.sub('(?P<in>(?:(?!\n).)*)\n', '\g<in>', data, 1) + '<br>'
                else:
                    if num > 0:
                        sub += str(left) + ' : ' + re.sub('(?P<in>.*)\n', '\g<in>', data, 1) + '<br>'

            end = re.sub('((?:(?!\n).)*)\n', '', end, 1)
        else:
            break
            
    return sub
           
def admin_check(num = None, what = None):
    ip = ip_check() 

    curs.execute("select acl from user where id = ?", [ip])
    user = curs.fetchall()
    if user:
        reset = 0

        while 1:
            if num == 1 and reset == 0:
                check = 'ban'
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
                    
    return 0

def ip_pas(raw_ip):
    hide = 0

    if re.search("(\.|:)", raw_ip):
        if not re.search("^tool:", raw_ip):    
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
        curs.execute("select title from data where title = ?", ['user:' + raw_ip])
        if curs.fetchall():
            ip = '<a href="/w/' + url_pas('user:' + raw_ip) + '">' + raw_ip + '</a>'
        else:
            ip = '<a id="not_thing" href="/w/' + url_pas('user:' + raw_ip) + '">' + raw_ip + '</a>'
         
    if hide == 0:
        ip += ' <a href="/tool/' + url_pas(raw_ip) + '">(' + load_lang('tool') + ')</a>'

    return ip

def custom():
    if 'head' in flask.session:
        user_head = flask.session['head']
    else:
        user_head = ''

    if 'state' in flask.session and flask.session['state'] == 1:
        curs.execute('select name from alarm where name = ? limit 1', [ip_check()])
        if curs.fetchall():
            user_icon = 2
        else:
            user_icon = 1
    else:
        user_icon = 0

    if user_icon != 0:
        curs.execute('select data from user_set where name = "email" and id = ?', [ip_check()])
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
        user_name = load_lang('user')

    return ['', '', user_icon, user_head, email, user_name]

def load_skin(data = ''):
    div2 = ''
    system_file = ['main_css', 'easter_egg.html']

    if data == '':
        ip = ip_check()

        curs.execute('select data from user_set where name = "skin" and id = ?', [ip])
        data = curs.fetchall()
        for skin_data in os.listdir(os.path.abspath('views')):
            if not skin_data in system_file:
                if not data:
                    curs.execute('select data from other where name = "skin"')
                    sql_data = curs.fetchall()
                    if sql_data and sql_data[0][0] == skin_data:
                        div2 = '<option value="' + skin_data + '">' + skin_data + '</option>' + div2
                    else:
                        div2 += '<option value="' + skin_data + '">' + skin_data + '</option>'
                elif data[0][0] == skin_data:
                    div2 = '<option value="' + skin_data + '">' + skin_data + '</option>' + div2
                else:
                    div2 += '<option value="' + skin_data + '">' + skin_data + '</option>'
    else:
        for skin_data in os.listdir(os.path.abspath('views')):
            if not skin_data in system_file:
                if data == skin_data:
                    div2 = '<option value="' + skin_data + '">' + skin_data + '</option>' + div2
                else:
                    div2 += '<option value="' + skin_data + '">' + skin_data + '</option>'

    return div2

def acl_check(name, tool = ''):
    ip = ip_check()
    
    if tool == 'render':
        curs.execute("select view from acl where title = ?", [name])
        acl_data = curs.fetchall()
        if acl_data:
            if acl_data[0][0] == 'user':
                if ip_or_user(ip):
                    return 1

            if acl_data[0][0] == 'admin':
                if ip_or_user(ip):
                    return 1

                if admin_check(5, 'view (' + name + ')') != 1:
                    return 1

        return 0
    else:
        if ban_check() == 1:
            return 1

        acl_c = re.search("^user:([^/]*)", name)
        if acl_c:
            acl_n = acl_c.groups()

            if admin_check(5, None) == 1:
                return 0

            curs.execute("select dec from acl where title = ?", ['user:' + acl_n[0]])
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

        file_c = re.search("^file:(.*)", name)
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
            if set_data[0][0] == 'login':
                if not user_data:
                    return 1

            if set_data[0][0] == 'admin':
                if not user_data:
                    return 1

                if not admin_check(5, None) == 1:
                    return 1

        return 0

def ban_check(ip = None, tool = None):
    if not ip:
        ip = ip_check()

    band = re.search("^([0-9]{1,3}\.[0-9]{1,3})", ip)
    if band:
        band_it = band.groups()[0]
    else:
        band_it = '-'
    
    curs.execute("select end, login from ban where block = ?", [band_it])
    band_d = curs.fetchall()
    
    curs.execute("select end, login from ban where block = ?", [ip])
    ban_d = curs.fetchall()
    
    data = band_d or ban_d
    if data and (data[0][0] == '' or data[0][0] > get_time()):
        if tool and tool == 'login':                    
            if data[0][1] == 'O':
                return 0
                
        return 1

    return 0
        
def topic_check(name, sub):
    ip = ip_check()

    if ban_check() == 1:
        return 1
        
    curs.execute("select acl from user where id = ?", [ip])
    user_data = curs.fetchall()

    curs.execute('select data from other where name = "discussion"')
    acl_data = curs.fetchall()
    if acl_data:
        if acl_data[0][0] == 'login':
            if not user_data:
                return 1

        if acl_data[0][0] == 'admin':
            if not user_data:
                return 1

            if not admin_check(3, 'topic (' + name + ')') == 1:
                return 1

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
        
    curs.execute("select title from rd where title = ? and sub = ? and not stop = ''", [name, sub])
    if curs.fetchall():
        if not admin_check(3, 'topic (' + name + ')') == 1:
            return 1

    return 0

def ban_insert(name, end, why, login, blocker):
    now_time = get_time()

    if re.search("^([0-9]{1,3}\.[0-9]{1,3})$", name):
        band = 'O'
    else:
        band = ''

    curs.execute("select block from ban where block = ?", [name])
    if curs.fetchall():
        curs.execute("insert into rb (block, end, today, blocker, why, band) values (?, ?, ?, ?, ?, ?)", [name, load_lang('release', 1), now_time, blocker, '', band])
        curs.execute("delete from ban where block = ?", [name])
    else:
        if login != '':
            login = 'O'
        else:
            login = ''

        if end != '0':
            time = datetime.datetime.now()
            plus = datetime.timedelta(seconds = int(end))
            r_time = (time + plus).strftime("%Y-%m-%d %H:%M:%S")
        else:
            r_time = ''

        curs.execute("insert into rb (block, end, today, blocker, why, band) values (?, ?, ?, ?, ?, ?)", [name, r_time, now_time, blocker, why, band])
        curs.execute("insert into ban (block, end, why, band, login) values (?, ?, ?, ?, ?)", [name, r_time, why, band, login])
    
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
    
    curs.execute("insert into history (id, title, data, date, ip, send, leng, hide) values (?, ?, ?, ?, ?, ?, ?, '')", [str(int(id_data[0][0]) + 1) if id_data else '1', title, data, date, ip, send, leng])

def leng_check(first, second):
    if first < second:
        all_plus = '+' + str(second - first)
    elif second < first:
        all_plus = '-' + str(first - second)
    else:
        all_plus = '0'
        
    return all_plus

def edit_filter_do(data):
    if admin_check(1, 'edit_filter pass') != 1:
        curs.execute("select regex, sub from filter")
        for data_list in curs.fetchall():
            match = re.compile(data_list[0], re.I)
            if match.search(data):
                ban_insert(
                    ip_check(), 
                    '0' if data_list[1] == 'X' else data_list[1], 
                    load_lang('edit', 1) + ' ' + load_lang('filter', 1), 
                    None, 
                    load_lang('tool', 1) + ':' + load_lang('edit', 1) + ' ' + load_lang('filter', 1)
                )
                
                return 1
    
    return 0

def redirect(data = '/'):
    return flask.redirect(data)

def re_error(data):
    conn.commit()
    
    if data == '/ban':
        ip = ip_check()

        end = '<li>' + load_lang('why') + ' : ' + load_lang('authority_error') + '</li>'

        if ban_check() == 1:
            curs.execute("select end, why from ban where block = ?", [ip])
            end_data = curs.fetchall()
            if not end_data:
                match = re.search("^([0-9]{1,3}\.[0-9]{1,3})", ip)
                if match:
                    curs.execute("select end, why from ban where block = ?", [match.groups()[0]])
                    end_data = curs.fetchall()
            
            if end_data:
                end = '<li>' + load_lang('state') + ' : ' + load_lang('ban') + '</li><li>'

                if end_data[0][0]:
                    now = int(re.sub('(\-| |:)', '', get_time()))
                    day = int(re.sub('(\-| |:)', '', end_data[0][0]))
                    
                    if now >= day:
                        curs.execute("delete from ban where block = ?", [ip])
                        conn.commit()

                        end += '<script>location.reload();</script>'
                    else:
                        end += 'end : ' + end_data[0][0]
                else:
                    end += load_lang('limitless')
                
                end += '</li>'

                if end_data[0][1] != '':
                    end += '<li>' + load_lang('why') + ' : ' + end_data[0][1] + '</li>'

        return easy_minify(flask.render_template(skin_check(), 
            imp = ['error', wiki_set(1), custom(), other2([0, 0])],
            data = '<h2>error</h2><ul>' + end + '</ul>',
            menu = 0
        ))
    else:
        error_data = re.search('\/error\/([0-9]+)', data)
        if error_data:
            num = int(error_data.groups()[0])
            if num == 1:
                data = load_lang('no_login_error')
            elif num == 2:
                data = load_lang('no_exist_user_error')
            elif num == 3:
                data = load_lang('authority_error')
            elif num == 4:
                data = load_lang('no_admin_block_error')
            elif num == 5:
                data = load_lang('skin_error')
            elif num == 6:
                data = load_lang('same_id_exist_error')
            elif num == 7:
                data = load_lang('long_id_error')
            elif num == 8:
                data = load_lang('id_char_error') + ' <a href="/name_filter">(' + load_lang('id') + ' ' + load_lang('filter') + ')</a>'
            elif num == 9:
                data = load_lang('file_exist_error')
            elif num == 10:
                data = load_lang('password_error')
            elif num == 13:
                data = load_lang('recaptcha_error')
            elif num == 14:
                data = load_lang('file_extension_error')
            elif num == 15:
                data = load_lang('edit_record_error')
            elif num == 16:
                data = load_lang('same_file_error')
            elif num == 17:
                data = load_lang('file_capacity_error') + ' ' + wiki_set(3)
            elif num == 19:
                data = load_lang('decument_exist_error')
            elif num == 20:
                data = load_lang('password_diffrent_error')
            elif num == 21:
                data = load_lang('edit_filter_error')
            elif num == 22:
                data = load_lang('file_name_error')
            else:
                data = '???'

            return easy_minify(flask.render_template(skin_check(), 
                imp = ['error', wiki_set(1), custom(), other2([0, 0])],
                data = '<h2>error</h2><ul><li>' + data + '</li></ul>',
                menu = 0
            ))
        else:
            return redirect('/')