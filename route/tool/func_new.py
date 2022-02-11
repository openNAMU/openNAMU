# Init
import os
import sys
import platform
import json
import smtplib
import zipfile
import shutil
import logging
import random

import email.mime.text
import email.utils
import email.header

import urllib.request

# Init-Version
version_list = json.loads(open('version.json', encoding = 'utf8').read())

print('Version : ' + version_list['beta']['r_ver'])
print('DB set version : ' + version_list['beta']['c_ver'])
print('Skin set version : ' + version_list['beta']['s_ver'])
print('----')

# Init-PIP_Install
data_up_date = 1
if os.path.exists(os.path.join('data', 'version.json')):
    data_load_ver = open(os.path.join('data', 'version.json'), encoding = 'utf8').read()
    if data_load_ver == version_list['beta']['r_ver']:
        data_up_date = 0

if data_up_date == 1:
    with open(os.path.join('data', 'version.json'), 'w', encoding = 'utf8') as f:
        f.write(version_list['beta']['r_ver'])
    
    if platform.system() in ('Linux', 'Windows'):
        os.system(
            'python' + ('3' if platform.system() != 'Windows' else '') + ' ' + \
            '-m pip install --upgrade --user -r requirements.txt'
        )
        
        print('----')
        try:
            os.execl(sys.executable, sys.executable, *sys.argv)
        except:
            pass

        try:
            os.execl(sys.executable, '"' + sys.executable + '"', *sys.argv)
        except:
            print('Error : restart failed')
            raise
    else:
        print('Error : automatic installation is not supported.')
        print('Help : try "python3 -m pip install -r requirements.txt"')
else:
    print('PIP check pass')
    
print('----')

# Init-Load
from .func_mark import *

from diff_match_patch import diff_match_patch

import netius.servers

import werkzeug.routing
import werkzeug.debug

import flask

import requests

import pymysql

if sys.version_info < (3, 6):
    import sha3
   
# Init-Global
global_lang = {}
global_wiki_set = {}

global_db_set = ''

data_css_ver = '118'
data_css = ''

conn = ''

# Func
# Func-main
def load_conn(data):
    global conn

    conn = data
    
def do_db_set(db_set):
    global global_db_set
    
    global_db_set = db_set
    
# Func-init
class get_db_connect_old:
    def __init__(self, db_set):
        self.db_set = db_set
        self.conn = ''
        
    def db_load(self):
        if self.db_set['type'] == 'sqlite':
            self.conn = sqlite3.connect(self.db_set['name'] + '.db')
        else:
            self.conn = pymysql.connect(
                host = self.db_set['mysql_host'],
                user = self.db_set['mysql_user'],
                password = self.db_set['mysql_pw'],
                charset = 'utf8mb4',
                port = int(self.db_set['mysql_port']),
            )
            curs = self.conn.cursor()

            try:
                curs.execute(db_change(
                    'create database ' + self.db_set['name'] + ' ' + \
                    'default character set utf8mb4;'
                ))
            except:
                pass

            self.conn.select_db(self.db_set['name'])

        load_conn(self.conn)

        return self.conn
    
    def db_get(self):
        # if self.db_set['type'] != 'sqlite':
        #     self.conn.ping(reconnect = True)
            
        return self.conn
    
class get_db_connect:
    def __init__(self):
        global global_db_set
        
        self.db_set = global_db_set
        
    def __enter__(self):
        if self.db_set['type'] == 'sqlite':
            self.conn = sqlite3.connect(self.db_set['name'] + '.db')
        else:
            self.conn = pymysql.connect(
                host = self.db_set['mysql_host'],
                user = self.db_set['mysql_user'],
                password = self.db_set['mysql_pw'],
                charset = 'utf8mb4',
                port = int(self.db_set['mysql_port']),
            )
            curs = self.conn.cursor()

            try:
                curs.execute(db_change(
                    'create database ' + self.db_set['name'] + ' ' + \
                    'default character set utf8mb4;'
                ))
            except:
                pass

            self.conn.select_db(self.db_set['name'])

        return self.conn
    
    def __exit__(self, exc_type, exc_value, traceback):
        self.conn.close()

def update(ver_num, set_data):
    curs = conn.cursor()

    print('----')
    # 업데이트 하위 호환 유지 함수

    if ver_num < 3160027:
        print('Add init set')
        set_init()

    if ver_num < 3170002:
        curs.execute(db_change("select html from html_filter where kind = 'extension'"))
        if not curs.fetchall():
            for i in ['jpg', 'jpeg', 'png', 'gif', 'webp']:
                curs.execute(db_change(
                    "insert into html_filter (html, kind) values (?, 'extension')"
                ), [i])

    if ver_num < 3170400:
        curs.execute(db_change("select title, sub, code from topic where id = '1'"))
        for i in curs.fetchall():
            curs.execute(db_change(
                "update topic set code = ? where title = ? and sub = ?"
            ), [
                i[2], 
                i[0], 
                i[1]
            ])
            curs.execute(db_change(
                "update rd set code = ? where title = ? and sub = ?"
            ), [
                i[2], 
                i[0], 
                i[1]
            ])

    if ver_num < 3171800:
        curs.execute(db_change("select data from other where name = 'recaptcha'"))
        change_rec = curs.fetchall()
        if change_rec and change_rec[0][0] != '':
            new_rec = re.search(r'data-sitekey="([^"]+)"', change_rec[0][0])
            if new_rec:
                curs.execute(db_change(
                    "update other set data = ? where name = 'recaptcha'"
                ), [new_rec.group(1)])
            else:
                curs.execute(db_change("update other set data = '' where name = 'recaptcha'"))
                curs.execute(db_change("update other set data = '' where name = 'sec_re'"))
    
    if  ver_num < 3172800 and \
        set_data['db_type'] == 'mysql':
        get_data_mysql = json.loads(open('data/mysql.json', encoding = 'utf8').read())
        
        with open('data/mysql.json', 'w') as f:
            f.write('{ "user" : "' + get_data_mysql['user'] + '", "password" : "' + get_data_mysql['password'] + '", "host" : "localhost" }')

    if ver_num < 3183603:
        curs.execute(db_change("select block from ban where band = 'O'"))
        for i in curs.fetchall():
            curs.execute(db_change(
                "update ban set block = ?, band = 'regex' where block = ? and band = 'O'"
            ), [
                '^' + i[0].replace('.', '\\.'),
                i[0]
            ])

        curs.execute(db_change("select block from rb where band = 'O'"))
        for i in curs.fetchall():
            curs.execute(db_change(
                "update rb set block = ?, band = 'regex' where block = ? and band = 'O'"
            ), [
                '^' + i[0].replace('.', '\\.'),
                i[0]
            ])

    if ver_num < 3190201:
        today_time = get_time()

        curs.execute(db_change("select block, end, why, band, login from ban"))
        for i in curs.fetchall():
            curs.execute(db_change(
                "insert into rb (block, end, today, why, band, login, ongoing) " + \
                "values (?, ?, ?, ?, ?, ?, ?)"
            ), [
                i[0],
                i[1],
                today_time,
                i[2],
                i[3],
                i[4],
                '1'
            ])

    if ver_num < 3191301:
        curs.execute(db_change('' + \
            'select id, title, date from history ' + \
            'where not title like "user:%" ' + \
            'order by date desc ' + \
            'limit 50' + \
        ''))
        data_list = curs.fetchall()
        for get_data in data_list:
            curs.execute(db_change(
                "insert into rc (id, title, date, type) values (?, ?, ?, 'normal')"
            ), [
                get_data[0], 
                get_data[1],
                get_data[2]
            ])

    if ver_num < 3202400:
        curs.execute(db_change("select data from other where name = 'update'"))
        get_data = curs.fetchall()
        if get_data and get_data[0][0] == 'master':
            curs.execute(db_change("update other set data = 'beta' where name = 'update'"), [])

    if ver_num < 3202600:
        curs.execute(db_change("select name, regex, sub from filter"))
        for i in curs.fetchall():
            curs.execute(db_change(
                "insert into html_filter (html, kind, plus, plus_t) " + \
                "values (?, 'regex_filter', ?, ?)"
            ), [
                i[0], 
                i[1],
                i[2]
            ])

        curs.execute(db_change("select title, link, icon from inter"))
        for i in curs.fetchall():
            curs.execute(db_change(
                "insert into html_filter (html, kind, plus, plus_t) " + \
                "values (?, 'inter_wiki', ?, ?)"), [
                i[0], 
                i[1],
                i[2]
            ])

    if ver_num < 3203400:
        curs.execute(db_change("select user, css from custom"))
        for i in curs.fetchall():
            curs.execute(db_change(
                "insert into user_set (name, id, data) values ('custom_css', ?, ?)"
            ), [
                re.sub(r' \(head\)$', '', i[0]), 
                i[1]
            ])

    if ver_num < 3205500:
        curs.execute(db_change("select title, decu, dis, view, why from acl"))
        for i in curs.fetchall():
            curs.execute(db_change(
                "insert into acl (title, data, type) values (?, ?, ?)"
            ), [i[0], i[1], 'decu'])
            curs.execute(db_change(
                "insert into acl (title, data, type) values (?, ?, ?)"
            ), [i[0], i[2], 'dis'])
            curs.execute(db_change(
                "insert into acl (title, data, type) values (?, ?, ?)"
            ), [i[0], i[3], 'view'])
            curs.execute(db_change(
                "insert into acl (title, data, type) values (?, ?, ?)"
            ), [i[0], i[4], 'why'])

    if ver_num < 3300101:
        # 캐시 초기화
        curs.execute(db_change('delete from cache_data'))
    
    if ver_num < 3300301:
        # regex_filter 오류 해결
        curs.execute(db_change(
            'delete from html_filter where kind = "regex_filter" and html is null'
        ))
        
    if ver_num < 3302302:
        # user이랑 user_set 테이블의 통합
        curs.execute(db_change('select id, pw, acl, date, encode from user'))
        for i in curs.fetchall():
            curs.execute(db_change(
                "insert into user_set (name, id, data) values (?, ?, ?)"
            ), ['pw', i[0], i[1]])
            curs.execute(db_change(
                "insert into user_set (name, id, data) values (?, ?, ?)"
            ), ['acl', i[0], i[2]])
            curs.execute(db_change(
                "insert into user_set (name, id, data) values (?, ?, ?)"
            ), ['date', i[0], i[3]])
            curs.execute(db_change(
                "insert into user_set (name, id, data) values (?, ?, ?)"
            ), ['encode', i[0], i[4]])
            
    if ver_num < 3400101:
        # user_set이랑 user_application 테이블의 통합
        curs.execute(db_change('' + \
            'select id, pw, date, encode, question, answer, ip, ua, email ' + \
            'from user_application' + \
        ''))
        for i in curs.fetchall():
            sql_data = {}
            sql_data['id'] = i[0]
            sql_data['pw'] = i[1]
            sql_data['date'] = i[2]
            sql_data['encode'] = i[3]
            sql_data['question'] = i[4]
            sql_data['answer'] = i[5]
            sql_data['ip'] = i[6]
            sql_data['ua'] = i[7]
            sql_data['email'] = i[8]
            
            curs.execute(db_change(
                "insert into user_set (name, id, data) values (?, ?, ?)"
            ), ['application', i[0], json.dumps(sql_data)])
    
    conn.commit()

    print('Update completed')

def set_init_always(ver_num):
    curs = conn.cursor()

    curs.execute(db_change('delete from other where name = "ver"'))
    curs.execute(db_change('insert into other (name, data) values ("ver", ?)'), [ver_num])
    
    curs.execute(db_change('delete from alist where name = "owner"'))
    curs.execute(db_change('insert into alist (name, acl) values ("owner", "owner")'))

    if not os.path.exists(load_image_url()):
        os.makedirs(load_image_url())
    
    conn.commit()
    
def set_init():
    curs = conn.cursor()

    # 초기값 설정 함수    
    curs.execute(db_change("select html from html_filter where kind = 'email'"))
    if not curs.fetchall():
        for i in ['naver.com', 'gmail.com', 'daum.net', 'kakao.com']:
            curs.execute(db_change(
                "insert into html_filter (html, kind, plus, plus_t) values (?, 'email', '', '')"
            ), [i])

    curs.execute(db_change("select html from html_filter where kind = 'extension'"))
    if not curs.fetchall():
        for i in ['jpg', 'jpeg', 'png', 'gif', 'webp']:
            curs.execute(db_change(
                "insert into html_filter (html, kind, plus, plus_t) values (?, 'extension', '', '')"
            ), [i])

    curs.execute(db_change(
        'select data from other ' + \
        'where name = "smtp_server" or name = "smtp_port" or name = "smtp_security"'
    ))
    if not curs.fetchall():
        for i in [
            ['smtp_server', 'smtp.gmail.com'],
            ['smtp_port', '587'], 
            ['smtp_security', 'starttls']
        ]:
            curs.execute(db_change(
                "insert into other (name, data) values (?, ?)"
            ), [i[0], i[1]])
            
    curs.execute(db_change('select data from other where name = "key"'))
    rep_data = curs.fetchall()
    if not rep_data:
        curs.execute(db_change('insert into other (name, data) values ("key", ?)'), [load_random_key()])

    curs.execute(db_change('select data from other where name = "count_all_title"'))
    if not curs.fetchall():
        curs.execute(db_change('insert into other (name, data) values ("count_all_title", "0")'))
        
    conn.commit()

# Func-simple
## Func-simple-without_DB
def get_default_admin_group():
    return ['owner', 'ban']

def load_random_key(long = 64):
    return ''.join(
        random.choice(
            "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
        ) for i in range(long)
    )

def http_warning():
    return '''
        <div id="http_warning_text"></div>
        <script>
            if(window.location.protocol !== 'https:') {
                document.getElementById('http_warning_text').innerHTML = "''' + \
                    load_lang('http_warning') + \
                '''";
                document.getElementById('http_warning_text').style.margin = "10px 0px 0px 0px";
            }
        </script>
    '''

def next_fix(link, num, page, end = 50):
    list_data = ''

    if num == 1:
        if len(page) == end:
            list_data += '' + \
                '<hr class="main_hr">' + \
                '<a href="' + link + str(num + 1) + '">(' + load_lang('next') + ')</a>' + \
            ''
    elif len(page) != end:
        list_data += '' + \
            '<hr class="main_hr">' + \
            '<a href="' + link + str(num - 1) + '">(' + load_lang('previous') + ')</a>' + \
        ''
    else:
        list_data += '' + \
            '<hr class="main_hr">' + \
            '<a href="' + link + str(num - 1) + '">(' + load_lang('previous') + ')</a> <a href="' + link + str(num + 1) + '">(' + load_lang('next') + ')</a>' + \
        ''

    return list_data

def leng_check(A, B):
    # B -> new
    # A -> old
    return '0' if A == B else (('-' + str(A - B)) if A > B else ('+' + str(B - A)))

def number_check(data):
    try:
        int(data)
        return data
    except:
        return '1'
    
def check_int(data):
    try:
        int(data)
        return data
    except:
        return ''
    
def redirect(data = '/'):
    return flask.redirect(flask.request.host_url[:-1] + data)
    
def get_acl_list(type_d = 'normal'):
    if type_d == 'user':
        return ['', 'user', 'all']
    else:
        return ['', 'all', 'user', 'admin', 'owner', '50_edit', 'email', 'ban', 'before', '30_day', 'ban_admin']

## Func-simple-with_DB
def load_image_url():
    curs = conn.cursor()

    curs.execute(db_change('select data from other where name = "image_where"'))
    image_where = curs.fetchall()
    image_where = image_where[0][0] if image_where else os.path.join('data', 'images')
    
    return image_where

def load_domain():
    curs = conn.cursor()

    curs.execute(db_change("select data from other where name = 'domain'"))
    domain = curs.fetchall()
    domain = domain[0][0] if domain and domain[0][0] != '' else flask.request.host_url

    return domain

def edit_button(editor_display = '0'):
    curs = conn.cursor()

    insert_list = []

    curs.execute(db_change("select html, plus from html_filter where kind = 'edit_top'"))
    db_data = curs.fetchall()
    for get_data in db_data:
        insert_list += [[get_data[1], get_data[0]]]

    data = ''
    for insert_data in insert_list:
        data += '' + \
            '<a href="' + \
                'javascript:do_insert_data(\'textarea_edit_view\', \'' + insert_data[0] + '\', ' + editor_display + ')' + \
            '">(' + insert_data[1] + ')</a> ' + \
        ''

    data += (' ' if data != '' else '') + '<a href="/edit_top">(' + load_lang('add') + ')</a>'
    data += '<hr class="main_hr">'
    
    return data

def ip_warning():
    curs = conn.cursor()

    if ip_or_user() != 0:
        curs.execute(db_change('select data from other where name = "no_login_warning"'))
        data = curs.fetchall()
        if data and data[0][0] != '':
            text_data = '' + \
                '<span>' + data[0][0] + '</span>' + \
                '<hr class="main_hr">' + \
            ''
        else:
            text_data = '' + \
                '<span>' + load_lang('no_login_warning') + '</span>' + \
                '<hr class="main_hr">' + \
            ''
    else:
        text_data = ''

    return text_data
    
# Func-login    
def pw_encode(data, type_d = ''):
    curs = conn.cursor()

    if type_d == '':
        curs.execute(db_change('select data from other where name = "encode"'))
        set_data = curs.fetchall()

        type_d = set_data[0][0]

    if type_d == 'sha256':
        return hashlib.sha256(bytes(data, 'utf-8')).hexdigest()
    else:
        if sys.version_info < (3, 6):
            return sha3.sha3_256(bytes(data, 'utf-8')).hexdigest()
        else:
            return hashlib.sha3_256(bytes(data, 'utf-8')).hexdigest()

def pw_check(data, data2, type_d = 'no', id_d = ''):
    curs = conn.cursor()

    curs.execute(db_change('select data from other where name = "encode"'))
    db_data = curs.fetchall()

    if type_d != 'no':
        if type_d == '':
            set_data = 'sha3'
        else:
            set_data = type_d
    else:
        set_data = db_data[0][0]

    re_data = 1 if pw_encode(data, set_data) == data2 else 0
    if db_data[0][0] != set_data and re_data == 1 and id_d != '':
        curs.execute(db_change("update user_set set data = ? where id = ? and name = 'pw'"), [
            pw_encode(data), 
            id_d
        ])
        curs.execute(db_change("update user_set set data = ? where id = ? and name = 'encode'"), [
            db_data[0][0], 
            id_d
        ])

    return re_data
        
# Func-skin
def easy_minify(data, tool = None):
    # without_DB

    return data

def load_lang(data, safe = 0):
    curs = conn.cursor()

    global global_lang

    ip = ip_check()
    if ip_or_user(ip) == 0:
        curs.execute(db_change('select data from user_set where name = "lang" and id = ?'), [ip])
        rep_data = curs.fetchall()                    
    elif 'lang' in flask.session:
        rep_data = [[flask.session['lang']]]
    else:
        curs.execute(db_change("select data from other where name = 'language'"))
        rep_data = curs.fetchall()

    if not rep_data or rep_data[0][0] in ('', 'default'):
        curs.execute(db_change("select data from other where name = 'language'"))
        rep_data = curs.fetchall()

    if rep_data:
        lang_name = rep_data[0][0]
    else:
        lang_name = 'en-US'
        
    if lang_name in global_lang:
        lang = global_lang[lang_name]
    else:
        lang_list = os.listdir('lang')
        if (lang_name + '.json') in lang_list:
            lang = json.loads(open(
                os.path.join('lang', lang_name + '.json'), 
                encoding = 'utf8'
            ).read())
            global_lang[lang_name] = lang
        else:
            lang = {}

    if data in lang:
        if safe == 1:
            return lang[data] 
        else:
            return html.escape(lang[data])

    return html.escape(data + ' (' + lang_name + ')')

def skin_check(set_n = 0):
    curs = conn.cursor()

    # 개편 필요?
    skin_list = load_skin('tenshi', 1)
    skin = skin_list[0]
    ip = ip_check()
    
    user_need_skin = ''
    if ip_or_user(ip) == 0:
        curs.execute(db_change('select data from user_set where name = "skin" and id = ?'), [ip])
        skin_exist = curs.fetchall()
        if skin_exist:
            user_need_skin = skin_exist[0][0]            
    else:
        if 'skin' in flask.session:
            user_need_skin = flask.session['skin']

    if user_need_skin == '':
        curs.execute(db_change('select data from other where name = "skin"'))
        skin_exist = curs.fetchall()
        if skin_exist:
            user_need_skin = skin_exist[0][0]
    
    if user_need_skin != '' and user_need_skin in skin_list:
        skin = user_need_skin

    if set_n == 0:
        return './views/' + skin + '/index.html'
    else:
        return skin
    
def wiki_css(data):
    # without_DB

    global data_css
    global data_css_ver

    data += ['' for _ in range(0, 3 - len(data))]

    if data_css == '':
        for i_data in os.listdir(os.path.join("views", "main_css", "css")):
            if i_data != 'sub':
                data_css += '<link rel="stylesheet" href="/views/main_css/css/' + i_data + '?ver=' + data_css_ver + '">'

        for i_data in os.listdir(os.path.join("views", "main_css", "js")):
            if i_data != 'sub':
                data_css += '<script src="/views/main_css/js/' + i_data + '?ver=' + data_css_ver + '"></script>'
                
    data = data[0:2] + ['', '''
        <link   rel="stylesheet"
                href="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.2.0/build/styles/default.min.css">
        <script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.2.0/build/highlight.min.js"></script>
        <link   defer rel="stylesheet"
                href="https://cdn.jsdelivr.net/npm/katex@0.13.2/dist/katex.min.css"
                integrity="sha384-Cqd8ihRLum0CCg8rz0hYKPoLZ3uw+gES2rXQXycqnL5pgVQIflxAUDS7ZSjITLb5"
                crossorigin="anonymous">
        <script src="https://cdn.jsdelivr.net/npm/katex@0.13.2/dist/katex.min.js"
                integrity="sha384-1Or6BdeNQb0ezrmtGeqQHFpppNd7a/gw29xeiSikBbsb44xu3uAo8c7FwbF5jhbd"
                crossorigin="anonymous"></script>
    ''' + data_css + '<script>window.addEventListener(\'DOMContentLoaded\', main_css_skin_load);</script>'] + data[2:]

    return data

def cut_100(data):
    # without_DB

    data = re.search(r'<pre style="display: none;" id="render_content_load">([^<>]+)<\/pre>', data)
    if data:
        data = data.group(1)
        if len(data) > 100:
            return data[0:100] + '...'
        else:
            return data[0:len(data)]
    else:
        return ''

def wiki_set(num = 1):
    curs = conn.cursor()

    if num == 1:
        skin_name = skin_check(1)
        data_list = []

        curs.execute(db_change('select data from other where name = ?'), ['name'])
        db_data = curs.fetchall()
        data_list += [db_data[0][0]] if db_data and db_data[0][0] != '' else ['Wiki']

        curs.execute(db_change('select data from other where name = "license"'))
        db_data = curs.fetchall()
        data_list += [db_data[0][0]] if db_data and db_data[0][0] != '' else ['ARR']

        data_list += ['', '']

        curs.execute(db_change('select data from other where name = "logo" and coverage = ?'), [skin_name])
        db_data = curs.fetchall()
        if db_data and db_data[0][0] != '':
            data_list += [db_data[0][0]]
        else:
            curs.execute(db_change('select data from other where name = "logo" and coverage = ""'))
            db_data = curs.fetchall()
            data_list += [db_data[0][0]] if db_data and db_data[0][0] != '' else [data_list[0]]

        head_data = ''

        curs.execute(db_change("select data from other where name = 'head' and coverage = ''"))
        db_data = curs.fetchall()
        head_data += db_data[0][0] if db_data and db_data[0][0] != '' else ''

        curs.execute(db_change("select data from other where name = 'head' and coverage = ?"), [skin_name])
        db_data = curs.fetchall()
        head_data += db_data[0][0] if db_data and db_data[0][0] != '' else ''
            
        data_list += [head_data]
    elif num == 2:
        curs.execute(db_change('select data from other where name = "frontpage"'))
        db_data = curs.fetchall()
        data_list = db_data[0][0] if db_data and db_data[0][0] != '' else 'FrontPage'
    elif num == 3:
        curs.execute(db_change('select data from other where name = "upload"'))
        db_data = curs.fetchall()
        data_list = db_data[0][0] if db_data and db_data[0][0] != '' else '2'
    else:
        data_list = ''

    return data_list

def wiki_custom():
    curs = conn.cursor()

    ip = ip_check()
    if ip_or_user(ip) == 0:
        user_icon = 1
        user_name = ip

        curs.execute(db_change("select data from user_set where id = ? and name = 'custom_css'"), [ip])
        user_head = curs.fetchall()
        user_head = user_head[0][0] if user_head else ''
        
        curs.execute(db_change('select data from user_set where name = "email" and id = ?'), [ip])
        email = curs.fetchall()
        email = email[0][0] if email else ''

        if admin_check('all') == 1:
            user_admin = '1'
            user_acl_list = []

            curs.execute(db_change("select data from user_set where id = ? and name = 'acl'"), [ip])
            curs.execute(db_change('select acl from alist where name = ?'), [curs.fetchall()[0][0]])
            user_acl = curs.fetchall()
            for i in user_acl:
                user_acl_list += [i[0]]

            user_acl_list = user_acl_list if user_acl != [] else '0'
        else:
            user_admin = '0'
            user_acl_list = '0'

        curs.execute(db_change("select count(*) from alarm where name = ?"), [ip])
        count = curs.fetchall()
        user_notice = str(count[0][0]) if count else '0'
    else:
        user_icon = 0
        user_name = load_lang('user')
        email = ''
        user_admin = '0'
        user_acl_list = '0'
        user_notice = '0'
        user_head = flask.session['head'] if 'head' in flask.session else ''

    curs.execute(db_change("select title from rd where title = ? and stop = ''"), ['user:' + ip])
    user_topic = '1' if curs.fetchall() else '0'
    
    split_path = flask.request.path.split('/')
    if len(split_path) > 1:
        split_path = split_path[1]
    else:
        split_path = 0

    return [
        '',
        '',
        user_icon,
        user_head,
        email,
        user_name,
        user_admin,
        str(ban_check()),
        user_notice,
        user_acl_list,
        ip,
        user_topic,
        split_path
    ]

def load_skin(data = '', set_n = 0, default = 0):
    # without_DB

    # data -> 가장 앞에 있을 스킨 이름
    # set_n == 0 -> 스트링으로 반환
    # set_n == 1 -> 리스트로 반환
    # default == 0 -> 디폴트 미포함
    # default == 1 -> 디폴트 포함

    if set_n == 0:
        skin_return_data = ''
    else:
        skin_return_data = []

    skin_list_get = os.listdir('views')
    if default == 1:
        skin_list_get = ['default'] + skin_list_get

    for skin_data in skin_list_get:
        if skin_data != 'default':
            see_data = skin_data
        else:
            see_data = load_lang('default')

        if skin_data != 'main_css':
            if set_n == 0:
                if skin_data == data:
                    skin_return_data = '' + \
                        '<option value="' + skin_data + '">' + \
                            see_data + \
                        '</option>' + \
                    '' + skin_return_data
                else:
                    skin_return_data += '' + \
                        '<option value="' + skin_data + '">' + \
                            see_data + \
                        '</option>' + \
                    ''
            else:
                if skin_data == data:
                    skin_return_data = [skin_data] + skin_return_data
                else:
                    skin_return_data += [skin_data]                    

    return skin_return_data

# Func-markup
def render_set(doc_name = '', doc_data = '', data_type = 'view', data_in = '', doc_acl = ''):
    # without_DB

    # data_type in ['view', 'raw', 'api_view', 'backlink']
    doc_acl = acl_check(doc_name, 'render') if doc_acl == '' else doc_acl
    doc_data = 0 if doc_data == None else doc_data
        
    if doc_acl == 1:
        return 'HTTP Request 401.3'
    else:
        if data_type == 'raw':
            return doc_data
        else:
            if doc_data != 0:
                get_class_render = class_do_render(conn)
                return get_class_render.do_render(doc_name, doc_data, data_type, data_in)
            else:
                return 'HTTP Request 404'

# Func-request
def send_email(who, title, data):
    curs = conn.cursor()

    try:
        curs.execute(db_change('' + \
            'select name, data from other ' + \
            'where name = "smtp_email" or name = "smtp_pass" or name = "smtp_server" or name = "smtp_port" or name = "smtp_security"' + \
        ''))
        rep_data = curs.fetchall()

        smtp_email = ''
        smtp_pass = ''
        smtp_server = ''
        smtp_security = ''
        smtp_port = ''
        smtp = ''

        for i in rep_data:
            if i[0] == 'smtp_email':
                smtp_email = i[1]
            elif i[0] == 'smtp_pass':
                smtp_pass = i[1]
            elif i[0] == 'smtp_server':
                smtp_server = i[1]
            elif i[0] == 'smtp_security':
                smtp_security = i[1]
            elif i[0] == 'smtp_port':
                smtp_port = i[1]
        
        smtp_port = int(smtp_port)
        if smtp_security == 'plain':
            smtp = smtplib.SMTP(smtp_server, smtp_port)
        elif smtp_security == 'starttls':
            smtp = smtplib.SMTP(smtp_server, smtp_port)
            smtp.starttls()
        else:
            # if smtp_security == 'tls':
            smtp = smtplib.SMTP_SSL(smtp_server, smtp_port)
        
        smtp.login(smtp_email, smtp_pass)

        domain = load_domain()
        wiki_name = wiki_set()[0]
        
        msg = email.mime.text.MIMEText(data)
        msg['Subject'] = title
        msg['From'] = 'openNAMU <noreply@' + domain + '>'
        msg['To'] = who
        
        smtp.sendmail('openNAMU@' + domain, who, msg.as_string())
        smtp.quit()

        return 1
    except Exception as e:
        print('----')
        print('Error : email send error')
        print(e)

        return 0

def captcha_get():
    curs = conn.cursor()

    data = ''
    
    if ip_or_user() != 0:
        curs.execute(db_change('select data from other where name = "recaptcha"'))
        recaptcha = curs.fetchall()
        
        curs.execute(db_change('select data from other where name = "sec_re"'))
        sec_re = curs.fetchall()
        
        curs.execute(db_change('select data from other where name = "recaptcha_ver"'))
        rec_ver = curs.fetchall()
        if  recaptcha and recaptcha[0][0] != '' and \
            sec_re and sec_re[0][0] != '':
            if not rec_ver or rec_ver[0][0] == '':
                data += '' + \
                    '<script src="https://www.google.com/recaptcha/api.js" async defer></script>' + \
                    '<div class="g-recaptcha" data-sitekey="' + recaptcha[0][0] + '"></div>' + \
                    '<hr class="main_hr">' + \
                ''
            else:
                data += '' + \
                    '<script src="https://www.google.com/recaptcha/api.js?render=' + recaptcha[0][0] + '"></script>' + \
                    '<input type="hidden" id="g-recaptcha" name="g-recaptcha">' + \
                    '<script type="text/javascript">' + \
                        'grecaptcha.ready(function() {' + \
                            'grecaptcha.execute(\'' + recaptcha[0][0] + '\', {action: \'homepage\'}).then(function(token) {' + \
                                'document.getElementById(\'g-recaptcha\').value = token;' + \
                            '});' + \
                        '});' + \
                    '</script>' + \
                ''

    return data

def captcha_post(re_data, num = 1):
    curs = conn.cursor()

    if num == 1:
        curs.execute(db_change('select data from other where name = "sec_re"'))
        sec_re = curs.fetchall()
        if  sec_re and sec_re[0][0] != '' and \
            ip_or_user() != 0 and captcha_get() != '':
            data = requests.get(
                'https://www.google.com/recaptcha/api/siteverify' + \
                '?secret=' + sec_re[0][0] + '&response=' + re_data
            )
            if data.status_code == 200:
                json_data = json.loads(data.text)
                if json_data['success'] != True:
                    return 1

        return 0

# Func-user
def ip_or_user(data = ''):
    # without_DB

    # 1 == ip
    # 0 == reg
    
    if data == '':
        data = ip_check()

    if re.search(r'(\.|:)', data):
        return 1
    else:
        return 0

def admin_check(num = None, what = None, name = ''):
    curs = conn.cursor()

    ip = ip_check() if name == '' else name
    time_data = get_time()
    pass_ok = 0

    if ip_or_user(ip) == 0:
        curs.execute(db_change(
            "select data from user_set where id = ? and name = 'acl'"
        ), [ip])
        user_auth = curs.fetchall()
        if user_auth:
            user_auth = user_auth[0][0]
            
            check = {
                0 : 'owner',
                1 : 'ban',
                2 : 'nothing',
                3 : 'toron',
                4 : 'check',
                5 : 'acl',
                6 : 'hidel',
                7 : 'give'
            }
            if not num:
                check = check[0]
            elif num == 'all':
                check = [check[i] for i in check]
            else:
                check = check[num]

            curs.execute(db_change(
                'select name from alist where name = ? and acl = "owner"'
            ), [user_auth])
            if curs.fetchall():
                pass_ok = 1
            else:
                if num == 'all':                    
                    curs.execute(db_change(
                        'select name from alist where name = ?'
                    ), [user_auth])
                else:
                    curs.execute(db_change(
                        'select name from alist where name = ? and acl = ?'
                    ), [user_auth, check])
                    
                if curs.fetchall():
                    pass_ok = 1

                
            if pass_ok == 1:
                if what:
                    curs.execute(db_change(
                        "insert into re_admin (who, what, time) values (?, ?, ?)"
                    ), [ip, what, time_data])
                    conn.commit()

                return 1

    return 0

def acl_check(name = 'test', tool = '', topic_num = '1'):
    curs = conn.cursor()

    ip = ip_check()
    get_ban = ban_check()
    acl_c = re.search(r"^user:((?:(?!\/).)*)", name) if name else None
    if tool == '' and acl_c:
        acl_n = acl_c.groups()

        if get_ban == 1:
            return 1

        if admin_check(5) == 1:
            return 0

        curs.execute(db_change(
            "select data from acl where title = ? and type = 'decu'"
        ), ['user:' + acl_n[0]])
        acl_data = curs.fetchall()
        if acl_data:
            if acl_data[0][0] == 'all':
                return 0
            elif acl_data[0][0] == 'user' and not ip_or_user(ip) == 1:
                return 0
            elif ip == acl_n[0] and not ip_or_user(ip) == 1:
                return 0
        else:
            if ip == acl_n[0] and not ip_or_user(ip) == 1 and not ip_or_user(acl_n[0]) == 1:
                return 0

        return 1

    if tool == 'topic':
        if not name:
            curs.execute(db_change("select title from rd where code = ?"), [topic_num])
            name = curs.fetchall()
            name = name[0][0] if name else 'test'
        
        end = 3
    elif tool == 'render' or tool == '' or tool == 'vote':
        if tool == '' and acl_check(name, 'render') == 1:
            return 1

        end = 2
    else:
        end = 1

    for i in range(0, end):
        if tool == '':
            if i == 0:
                curs.execute(db_change(
                    "select data from acl where title = ? and type = 'decu'"
                ), [name])
            else:
                curs.execute(db_change('select data from other where name = "edit"'))

            num = 5
        elif tool == 'topic':
            if i == 0 and topic_num:
                curs.execute(db_change("select acl from rd where code = ?"), [topic_num])
            elif i == 1:
                curs.execute(db_change(
                    "select data from acl where title = ? and type = 'dis'"
                ), [name])
            else:
                curs.execute(db_change('select data from other where name = "discussion"'))

            num = 3
        elif tool == 'upload':
            curs.execute(db_change("select data from other where name = 'upload_acl'"))

            num = 5
        elif tool == 'many_upload':
            curs.execute(db_change("select data from other where name = 'many_upload_acl'"))

            num = 5
        elif tool == 'vote':
            if i == 0:
                curs.execute(db_change(
                    'select acl from vote where id = ? and user = ""'
                ), [topic_num])
            else:
                curs.execute(db_change('select data from other where name = "vote_acl"'))

            num = None
        else:
            # tool == 'render'
            if i == 0:
                curs.execute(db_change(
                    "select data from acl where title = ? and type = 'view'"
                ), [name])
            else:
                curs.execute(db_change("select data from other where name = 'all_view_acl'"))

            num = 5

        acl_data = curs.fetchall()
        if  (
                i == (end - 1) and \
                (not acl_data or acl_data[0][0] == '' or acl_data[0][0] == 'normal')
            ) and \
            get_ban == 1 and \
            tool != 'render':
            return 1
        elif acl_data and acl_data[0][0] != 'normal' and acl_data[0][0] != '':
            if acl_data[0][0] != 'ban' and get_ban == 1 and tool != 'render':
                return 1

            if acl_data[0][0] == 'all' or acl_data[0][0] == 'ban':
                return 0
            elif acl_data[0][0] == 'user':
                if ip_or_user(ip) != 1:
                    return 0
            elif acl_data[0][0] == 'admin':
                if ip_or_user(ip) != 1:
                    if admin_check(num) == 1:
                        return 0
            elif acl_data[0][0] == '50_edit':
                if ip_or_user(ip) != 1:
                    if admin_check(num) == 1:
                        return 0
                    else:
                        curs.execute(db_change(
                            "select count(*) from history where ip = ?"
                        ), [ip])
                        count = curs.fetchall()
                        count = count[0][0] if count else 0
                        if count >= 50:
                            return 0
            elif acl_data[0][0] == 'before':
                if ip_or_user(ip) != 1:
                    if admin_check(num) == 1:
                        return 0
                
                curs.execute(db_change(
                    "select ip from history where title = ? and ip = ?"
                ), [name, ip])
                if curs.fetchall():
                    return 0
            elif acl_data[0][0] == '30_day':
                if ip_or_user(ip) != 1:
                    if admin_check(num) == 1:
                        return 0
                    else:
                        curs.execute(db_change(
                            "select data from user_set where id = ? and name = 'date'"
                        ), [ip])
                        user_date = curs.fetchall()[0][0]
                        
                        time_1 = datetime.datetime.strptime(
                            user_date, 
                            '%Y-%m-%d %H:%M:%S'
                        ) + datetime.timedelta(days = 30)
                        time_2 = datetime.datetime.strptime(
                            get_time(), 
                            '%Y-%m-%d %H:%M:%S'
                        )
                        
                        if time_2 > time_1:
                            return 0
            elif acl_data[0][0] == 'email':
                if ip_or_user(ip) != 1:
                    if admin_check(num) == 1:
                        return 0
                    else:
                        curs.execute(db_change(
                            "select data from user_set where id = ? and name = 'email'"
                        ), [ip])
                        if curs.fetchall():
                            return 0
            elif acl_data[0][0] == 'owner':
                if admin_check() == 1:
                    return 0
            elif acl_data[0][0] == 'ban_admin':
                if admin_check(1) == 1 or ban_check() == 1:
                    return 0

            return 1
        else:
            if i == (end - 1):
                if tool == 'topic' and topic_num:
                    curs.execute(db_change(
                        "select title from rd where code = ? and stop != ''"
                    ), [topic_num])
                    if curs.fetchall():
                        if admin_check(3, 'topic (code ' + topic_num + ')') == 1:
                            return 0
                    else:
                        return 0
                else:
                    return 0

    return 1

def ban_check(ip = None, tool = ''):
    curs = conn.cursor()

    ip = ip_check() if not ip else ip
    tool = '' if not tool else tool

    if admin_check(None, None, ip) == 1:
        return 0

    curs.execute(db_change(
        "update rb set ongoing = '' " + \
        "where end < ? and end != '' and ongoing = '1'"
    ), [get_time()])
    conn.commit()

    curs.execute(db_change("" + \
        "select login, block from rb " + \
        "where band = 'regex' and ongoing = '1'" + \
    ""))
    regex_d = curs.fetchall()
    for test_r in regex_d:
        g_regex = re.compile(test_r[1])
        if g_regex.search(ip):
            if tool == 'login':
                if test_r[0] != 'O':
                    return 1
            else:
                return 1

    curs.execute(db_change("" + \
        "select login from rb " + \
        "where block = ? and band = '' and ongoing = '1'" + \
        "" + \
    ""), [ip])
    ban_d = curs.fetchall()
    if ban_d:
        if tool == 'login':
            if ban_d[0][0] != 'O':
                return 1
        else:
            return 1

    return 0

def ip_pas(raw_ip, type_d = 0):
    curs = conn.cursor()

    hide = 0
    end_ip = {}
    i = 0

    return_data = 0
    if type(raw_ip) != type([]):
        get_ip = [raw_ip]
        
        return_data = 1
    else:
        get_ip = raw_ip

    curs.execute(db_change("select data from other where name = 'ip_view'"))
    ip_view = curs.fetchall()
    ip_view = ip_view[0][0] if ip_view else ''
    ip_view = '' if admin_check(1) == 1 else ip_view
    
    get_ip = list(set(get_ip))
    
    for raw_ip in get_ip:
        change_ip = 0
        is_this_ip = ip_or_user(raw_ip)
        if is_this_ip != 0 and ip_view != '':
            ip = re.sub(r'\.([^.]*)\.([^.]*)$', '.*.*', raw_ip)
            ip = re.sub(r':([^:]*):([^:]*)$', ':*:*', ip)
                
            change_ip = 1
        else:     
            ip = raw_ip
            
        if type_d == 0:
            if is_this_ip == 0:
                ip = '<a href="/w/' + url_pas('user:' + raw_ip) + '">' + raw_ip + '</a>'
                
            if change_ip == 0:
                ip += ' <a href="/user/' + url_pas(raw_ip) + '">(' + load_lang('tool') + ')</a>'

        end_ip[raw_ip] = ip
    
    if return_data == 1:
        return end_ip[raw_ip]
    else:
        return end_ip
        
# Func-edit
def slow_edit_check():
    curs = conn.cursor()

    curs.execute(db_change("select data from other where name = 'slow_edit'"))
    slow_edit = curs.fetchall()
    if slow_edit and slow_edit != '0' and admin_check(5) != 1:
        slow_edit = slow_edit[0][0]

        curs.execute(db_change(
            "select date from history where ip = ? order by date desc limit 1"
        ), [ip_check()])
        last_edit_data = curs.fetchall()
        if last_edit_data:
            last_edit_data = int(re.sub(' |:|-', '', last_edit_data[0][0]))
            now_edit_data = int(
                (datetime.datetime.now() - datetime.timedelta(
                    seconds = int(slow_edit))
                ).strftime("%Y%m%d%H%M%S")
            )

            if last_edit_data > now_edit_data:
                return 1

    return 0

def edit_filter_do(data):
    curs = conn.cursor()

    if admin_check(1) != 1:
        curs.execute(db_change(
            "select plus, plus_t from html_filter where kind = 'regex_filter' and plus != ''"
        ))
        for data_list in curs.fetchall():
            match = re.compile(data_list[0], re.I)
            if match.search(data):
                ban_insert(
                    ip_check(),
                    '0' if data_list[1] == 'X' else data_list[1],
                    'edit filter',
                    None,
                    'tool:edit filter'
                )

                return 1

    return 0

# Func-insert
def add_alarm(who, context):
    curs = conn.cursor()

    curs.execute(db_change(
        'insert into alarm (name, data, date) values (?, ?, ?)'
    ), [who, context, get_time()])
    conn.commit()
    
def add_user(user_name, user_pw, user_email = '', user_encode = ''):
    curs = conn.cursor()

    if user_encode == '':
        user_pw_hash = pw_encode(user_pw)

        curs.execute(db_change('select data from other where name = "encode"'))
        data_encode = curs.fetchall()
        data_encode = data_encode[0][0]
    else:
        user_pw_hash = user_pw
        data_encode = user_encode

    curs.execute(db_change("select id from user_set limit 1"))
    if not curs.fetchall():
        user_auth = 'owner'
    else:
        user_auth = 'user'

    curs.execute(db_change("insert into user_set (id, name, data) values (?, 'pw', ?)"), [
        user_name,
        user_pw_hash
    ])
    curs.execute(db_change("insert into user_set (id, name, data) values (?, 'acl', ?)"), [
        user_name,
        user_auth
    ])
    curs.execute(db_change("insert into user_set (id, name, data) values (?, 'date', ?)"), [
        user_name,
        get_time()
    ])
    curs.execute(db_change("insert into user_set (id, name, data) values (?, 'encode', ?)"), [
        user_name,
        data_encode
    ])
    
    if user_email != '':
        curs.execute(db_change("insert into user_set (name, id, data) values ('email', ?, ?)"), [
            user_name,
            user_email
        ])
        
    conn.commit()
    
def ua_plus(u_id, u_ip, u_agent, time):
    curs = conn.cursor()

    curs.execute(db_change("select data from other where name = 'ua_get'"))
    rep_data = curs.fetchall()
    if rep_data and rep_data[0][0] != '':
        pass
    else:
        curs.execute(db_change(
            "insert into ua_d (name, ip, ua, today, sub) values (?, ?, ?, ?, '')"
        ), [
            u_id, 
            u_ip, 
            u_agent, 
            time
        ])
        conn.commit()

def ban_insert(name, end, why, login, blocker, type_d = None):
    curs = conn.cursor()

    now_time = get_time()
    band = type_d if type_d else ''

    curs.execute(db_change(
        "update rb set ongoing = '' where end < ? and end != '' and ongoing = '1'"
    ), [now_time])
    curs.execute(db_change("" + \
        "select block from rb " + \
        "where ((end > ? and end != '') or end = '') and block = ? and " + \
            "band = ? and ongoing = '1'" + \
    ""), [now_time, name, band])
    if curs.fetchall():
        curs.execute(db_change(
            "insert into rb (block, end, today, blocker, why, band) values (?, ?, ?, ?, ?, ?)"
        ), [
            name,
            'release',
            now_time,
            blocker,
            '',
            band
        ])
        curs.execute(db_change(
            "update rb set ongoing = '' where block = ? and band = ? and ongoing = '1'"
        ), [name, band])
    else:
        login = 'O' if login != '' else ''

        if end != '0':
            end = int(number_check(end))
            time = datetime.datetime.now()
            plus = datetime.timedelta(seconds = end)
            r_time = (time + plus).strftime("%Y-%m-%d %H:%M:%S")
        else:
            r_time = ''

        curs.execute(db_change(
            "insert into rb (block, end, today, blocker, why, band, ongoing, login) " + \
            "values (?, ?, ?, ?, ?, ?, '1', ?)"
        ), [
            name, 
            r_time, 
            now_time, 
            blocker, 
            why, 
            band,
            login
        ])

    conn.commit()

def rd_plus(topic_num, date, name = None, sub = None):
    curs = conn.cursor()

    curs.execute(db_change("select code from rd where code = ?"), [topic_num])
    if curs.fetchall():
        curs.execute(db_change("update rd set date = ? where code = ?"), [date, topic_num])
    else:
        curs.execute(db_change(
            "insert into rd (title, sub, code, date) values (?, ?, ?, ?)"
        ), [name, sub, topic_num, date])

    conn.commit()

def history_plus(title, data, date, ip, send, leng, t_check = '', mode = ''):
    curs = conn.cursor()

    if mode == 'add':
        curs.execute(db_change(
            "select id from history where title = ? order by id + 0 asc limit 1"
        ), [title])
        id_data = curs.fetchall()
        id_data = str(int(id_data[0][0]) - 1) if id_data else '0'
    else:
        curs.execute(db_change(
            "select id from history where title = ? order by id + 0 desc limit 1"
        ), [title])
        id_data = curs.fetchall()
        id_data = str(int(id_data[0][0]) + 1) if id_data else '1'
        
        mode = mode if not re.search('^user:', title) else 'user'

    send = re.sub(r'\(|\)|<|>', '', send)
    send = send[:128] if len(send) > 128 else send
    send = send + ' (' + t_check + ')' if t_check != '' else send

    if mode != 'add' and mode != 'user':
        curs.execute(db_change("select count(*) from rc where type = 'normal'"))
        if curs.fetchall()[0][0] >= 200:
            curs.execute(db_change(
                "select id, title from rc where type = 'normal' order by date asc limit 1"
            ))
            rc_data = curs.fetchall()
            if rc_data:
                curs.execute(db_change(
                    'delete from rc where id = ? and title = ? and type = "normal"'
                ), [
                    rc_data[0][0],
                    rc_data[0][1]
                ])
    
        curs.execute(db_change(
            "insert into rc (id, title, date, type) values (?, ?, ?, 'normal')"
        ), [
            id_data,
            title,
            date
        ])
    
    if mode != 'add':
        curs.execute(db_change("select count(*) from rc where type = ?"), [mode])
        if curs.fetchall()[0][0] >= 200:
            curs.execute(db_change(
                "select id, title from rc where type = ? order by date asc limit 1"
            ), [mode])
            rc_data = curs.fetchall()
            if rc_data:
                curs.execute(db_change(
                    'delete from rc where id = ? and title = ? and type = ?'
                ), [
                    rc_data[0][0],
                    rc_data[0][1],
                    mode
                ])
    
        curs.execute(db_change(
            "insert into rc (id, title, date, type) values (?, ?, ?, ?)"
        ), [
            id_data,
            title,
            date,
            mode
        ])
            
    curs.execute(db_change(
        "insert into history (id, title, data, date, ip, send, leng, hide, type) " + \
        "values (?, ?, ?, ?, ?, ?, ?, '', ?)"
    ), [
        id_data,
        title,
        data,
        date,
        ip,
        send,
        leng,
        mode
    ])

# Func-error
def re_error(data):
    curs = conn.cursor()

    conn.commit()

    if data == '/ban':
        if ban_check() == 1:
            end = '<div id="get_user_info"></div><script>load_user_info("' + ip_check() + '");</script>'
        else:
            end = '<ul class="inside_ul"><li>' + load_lang('authority_error') + '</li></ul>'

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('error'), wiki_set(1), wiki_custom(), wiki_css([0, 0])],
            data = '<h2>' + load_lang('error') + '</h2>' + end,
            menu = 0
        )), 401
    else:
        num = int(number_check(data.replace('/error/', '')))
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
            data = load_lang('id_char_error') + ' <a href="/name_filter">(' + load_lang('id_filter_list') + ')</a>'
        elif num == 9:
            data = load_lang('file_exist_error')
        elif num == 10:
            data = load_lang('password_error')
        elif num == 11:
            data = load_lang('topic_long_error')
        elif num == 12:
            data = load_lang('email_error')
        elif num == 13:
            data = load_lang('recaptcha_error')
        elif num == 14:
            data = load_lang('file_extension_error') + ' <a href="/extension_filter">(' + load_lang('extension_filter_list') + ')</a>'
        elif num == 15:
            data = load_lang('edit_record_error')
        elif num == 16:
            data = load_lang('same_file_error')
        elif num == 17:
            data = load_lang('file_capacity_error') + wiki_set(3)
        elif num == 18:
            data = load_lang('email_send_error')
        elif num == 19:
            data = load_lang('decument_exist_error')
        elif num == 20:
            data = load_lang('password_diffrent_error')
        elif num == 21:
            data = load_lang('edit_filter_error')
        elif num == 22:
            data = load_lang('file_name_error')
        elif num == 23:
            data = load_lang('regex_error')
        elif num == 24:
            curs.execute(db_change("select data from other where name = 'slow_edit'"))
            slow_edit = curs.fetchall()
            slow_edit = '' if not slow_edit else slow_edit[0][0]
            data = load_lang('fast_edit_error') + slow_edit
        elif num == 25:
            data = load_lang('too_many_dec_error')
        elif num == 26:
            data = load_lang('application_not_found')
        elif num == 27:
            data = load_lang("invalid_password_error")
        elif num == 28:
            data = load_lang('watchlist_overflow_error')
        elif num == 29:
            data = load_lang('copyright_disagreed')
        elif num == 30:
            data = load_lang('ie_wrong_callback')
        elif num == 33:
            data = load_lang('restart_fail_error')
        elif num == 34:
            data = load_lang("update_error") + ' <a href="https://github.com/opennamu/opennamu">(Github)</a>'
        elif num == 35:
            data = load_lang('same_email_error')
        elif num == 36:
            data = load_lang('input_email_error')
        else:
            data = '???'

        if num == 5:
            if flask.request.path != '/main_skin_set':
                title = load_lang('skin_set')
                tool = [['main_skin_set', load_lang('main_skin_set')]]
                load_skin_set = ''
            else:
                title = load_lang('main_skin_set')
                tool = [['skin_set', load_lang('skin_set')]]
                load_skin_set = '<script>main_css_skin_set();</script>'
        
            return easy_minify(flask.render_template(skin_check(),
                imp = [title, wiki_set(1), wiki_custom(), wiki_css([0, 0])],
                data = '' + \
                    '<div id="main_skin_set">' + \
                        '<h2>' + load_lang('error') + '</h2>' + \
                        '<ul class="inside_ul">' + \
                            '<li>' + data + ' <a href="/main_skin_set">(' + load_lang('main_skin_set') + ')</a></li>' + \
                        '</ul>' + \
                    '</div>' + \
                    load_skin_set,
                menu = tool
            ))
        else:
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('error'), wiki_set(1), wiki_custom(), wiki_css([0, 0])],
                data = '' + \
                     '<h2>' + load_lang('error') + '</h2>' + \
                     '<ul class="inside_ul">' + \
                         '<li>' + data + '</li>' + \
                     '</ul>' + \
                '',
                menu = 0
            )), 400