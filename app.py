# 모듈 불러오기
from flask import Flask, request, send_from_directory
from flask_compress import Compress
from flask_reggie import Reggie

from tornado.wsgi import WSGIContainer
from tornado.httpserver import HTTPServer
from tornado.ioloop import IOLoop

import urllib.request

import platform
import zipfile
import bcrypt
import difflib
import shutil
import threading
import logging
import random
import sys

# 나머지 불러오기
from func import *
from set_mark.tool import savemark

# 버전 표기
r_ver = 'v3.0.4-Beta-' + time.strftime('%y%m%d', time.localtime(os.stat('app.py').st_mtime))
print('Version : ' + r_ver)

# set.json 설정 확인
try:
    json_data = open('set.json').read()
    set_data = json.loads(json_data)
except:
    while 1:
        print('DB Name : ', end = '')
        
        new_json = str(input())
        if new_json != '':
            with open("set.json", "w") as f:
                f.write('{ "db" : "' + new_json + '" }')
            
            json_data = open('set.json').read()
            set_data = json.loads(json_data)

            break
        else:
            print('Insert Values')
            
            pass

# 디비 연결
conn = sqlite3.connect(set_data['db'] + '.db', check_same_thread = False)
curs = conn.cursor()

# 보내주기
load_conn(conn)

# 기타 설정 변경
logging.basicConfig(level = logging.ERROR)
app = Flask(__name__, template_folder = './')
Reggie(app)
compress = Compress()
compress.init_app(app)

# 템플릿 설정
def md5_replace(data):
    return hashlib.md5(data.encode()).hexdigest()       

def load_lang_skin(data):
    return load_lang(data)

app.jinja_env.filters['md5_replace'] = md5_replace
app.jinja_env.filters['load_lang_skin'] = load_lang_skin

# 셋업 부분
curs.execute("create table if not exists data(title text, data text)")
curs.execute("create table if not exists history(id text, title text, data text, date text, ip text, send text, leng text, hide text)")
curs.execute("create table if not exists rd(title text, sub text, date text)")
curs.execute("create table if not exists user(id text, pw text, acl text, date text, email text, skin text)")
curs.execute("create table if not exists ban(block text, end text, why text, band text, login text)")
curs.execute("create table if not exists topic(id text, title text, sub text, data text, date text, ip text, block text, top text)")
curs.execute("create table if not exists stop(title text, sub text, close text)")
curs.execute("create table if not exists rb(block text, end text, today text, blocker text, why text, band text)")
curs.execute("create table if not exists back(title text, link text, type text)")
curs.execute("create table if not exists agreedis(title text, sub text)")
curs.execute("create table if not exists custom(user text, css text)")
curs.execute("create table if not exists other(name text, data text)")
curs.execute("create table if not exists alist(name text, acl text)")
curs.execute("create table if not exists re_admin(who text, what text, time text)")
curs.execute("create table if not exists alarm(name text, data text, date text)")
curs.execute("create table if not exists ua_d(name text, ip text, ua text, today text, sub text)")
curs.execute("create table if not exists filter(name text, regex text, sub text)")
curs.execute("create table if not exists scan(user text, title text)")
curs.execute("create table if not exists acl(title text, dec text, dis text, why text)")
curs.execute("create table if not exists inter(title text, link text)")
curs.execute("create table if not exists html_filter(html text)")

# owner 존재 확인
curs.execute("select name from alist where acl = 'owner'")
if not curs.fetchall():
    curs.execute("delete from alist where name = 'owner'")
    curs.execute("insert into alist (name, acl) values ('owner', 'owner')")

# 포트 점검
curs.execute("select data from other where name = 'port'")
rep_data = curs.fetchall()
if not rep_data:
    while 1:
        print('Port : ', end = '')
        
        rep_port = int(input())
        if rep_port:
            curs.execute("insert into other (name, data) values ('port', ?)", [rep_port])
            
            break
        else:
            pass
else:
    rep_port = rep_data[0][0]
    
    print('Port : ' + str(rep_port))

# robots.txt 점검
try:
    if not os.path.exists('robots.txt'):
        curs.execute("select data from other where name = 'robot'")
        robot_test = curs.fetchall()
        if robot_test:
            fw_test = open('./robots.txt', 'w')
            fw_test.write(re.sub('\r\n', '\n', robot_test[0][0]))
            fw_test.close()
        else:
            fw_test = open('./robots.txt', 'w')
            fw_test.write('User-agent: *\nDisallow: /\nAllow: /$\nAllow: /w/')
            fw_test.close()

            curs.execute("insert into other (name, data) values ('robot', 'User-agent: *\nDisallow: /\nAllow: /$\nAllow: /w/')")
        
        print('robots.txt create')
except:
    pass

# 비밀 키 점검
curs.execute("select data from other where name = 'key'")
rep_data = curs.fetchall()
if not rep_data:
    while 1:
        print('Secret Key : ', end = '')
        
        rep_key = str(input())
        if rep_key:
            curs.execute("insert into other (name, data) values ('key', ?)", [rep_key])
            
            break
        else:
            pass
else:
    rep_key = rep_data[0][0]

    print('Secret Key : ' + rep_key)

# 언어 점검
curs.execute("select data from other where name = 'language'")
rep_data = curs.fetchall()
if not rep_data:
    while 1:
        print('Language [ko-KR, en-US] : ', end = '')
        support_language = ['ko-KR', 'en-US']
        
        rep_language = str(input())
        if rep_language in support_language:
            curs.execute("insert into other (name, data) values ('language', ?)", [rep_language])
            
            break
        else:
            pass
else:
    rep_language = rep_data[0][0]
    
    print('Language : ' + str(rep_language))

json_data = open(os.path.join('language', rep_language + '.json'), 'rt', encoding='utf-8').read()
lang_data = json.loads(json_data)

# 한번 개행
print('')

# 호환성 설정
try:
    curs.execute("alter table history add hide text default ''")
    
    curs.execute('select title, re from hidhi')
    for rep in curs.fetchall():
        curs.execute("update history set hide = 'O' where title = ? and id = ?", [rep[0], rep[1]])

    curs.execute("drop table if exists hidhi")

    print('move table hidhi')
except:
    pass

try:
    curs.execute("alter table user add date text default ''")

    print('user table add column date')
except:
    pass

try:
    curs.execute("alter table rb add band text default ''")

    print('rb table add column band')
except:
    pass

try:
    curs.execute("alter table ban add login text default ''")

    print('ban table add column login')
except:
    pass

try:
    curs.execute("select title, acl from data where acl != ''")
    for rep in curs.fetchall():
        curs.execute("insert into acl (title, dec, dis, why) values (?, ?, '', '')", [rep[0], rep[1]])

    curs.execute("alter table data drop acl")

    print('data table delete column acl')
except:
    pass

try:
    curs.execute("alter table user add email text default ''")

    print('user table add column email')
except:
    pass

try:
    curs.execute('select name, sub from filter where sub != "X" and sub != ""')
    filter_name = curs.fetchall()
    if filter_name:
        for filter_delete in filter_name:
            if filter_delete[1] != '' or filter_delete[1] != 'X':
                curs.execute("update filter set sub = '' where name = ?", [filter_delete[0]])

        print('filter data fix')
except:
    pass

try:
    curs.execute("alter table user add skin text default ''")

    print('user table add column skin')
except:
    pass
        
conn.commit()

# 이미지 폴더 생성
if not os.path.exists('image'):
    os.makedirs('image')
    
# 스킨 폴더 생성
if not os.path.exists('views'):
    os.makedirs('views')

# 백업 설정
def back_up():
    try:
        shutil.copyfile(set_data['db'] + '.db', 'back_' + set_data['db'] + '.db')
        
        print('Back Up Ok')
    except:
        print('Back Up Error')

    threading.Timer(60 * 60 * back_time, back_up).start()

try:
    curs.execute('select data from other where name = "back_up"')
    back_up_time = curs.fetchall()
    
    back_time = int(back_up_time[0][0])
except:
    back_time = 0
    
# 백업 여부 확인
if back_time != 0:
    print(str(back_time) + '-Hours Interval Back Up')
    
    if __name__ == '__main__':
        back_up()
else:
    print('Back Up OFF')

@app.route('/del_alarm')
def del_alarm():
    curs.execute("delete from alarm where name = ?", [ip_check()])
    conn.commit()

    return redirect('/alarm')

@app.route('/alarm')
def alarm():
    if custom()[2] == 0:
        return redirect('/login')    

    data = '<ul>'    
    
    curs.execute("select data, date from alarm where name = ? order by date desc", [ip_check()])
    data_list = curs.fetchall()
    if data_list:
        data = '<a href="/del_alarm">(' + load_lang('delete') + ')</a><hr>' + data

        for data_one in data_list:
            data += '<li>' + data_one[0] + ' (' + data_one[1] + ')</li>'
    else:
        data += '<li>' + load_lang('no_alarm') + '</li>'
    
    data += '</ul>'

    return html_minify(render_template(skin_check(), 
        imp = [load_lang('alarm'), wiki_set(1), custom(), other2([0, 0])],
        data = data,
        menu = [['user', load_lang('user')]]
    ))

@app.route('/<regex("inter_wiki|html_filter"):tools>')
def inter_wiki(tools = None):
    div = ''
    admin = admin_check(None, None)

    if tools == 'inter_wiki':
        del_link = 'del_inter_wiki'
        plus_link = 'plus_inter_wiki'
        title = load_lang('interwiki') + ' ' + load_lang('list')
        div = ''

        curs.execute('select title, link from inter')
    else:
        del_link = 'del_html_filter'
        plus_link = 'plus_html_filter'
        title = 'HTML Filter ' + load_lang('list')
        div = '<ul><li>span</li><li>div</li><li>iframe</li></ul>'

        curs.execute('select html from html_filter')

    db_data = curs.fetchall()
    if db_data:
        div += '<ul>'

        for data in db_data:
            if tools == 'inter_wiki':
                div += '<li>' + data[0] + ' : <a id="out_link" href="' + data[1] + '">' + data[1] + '</a>'
            else:
                div += '<li>' + data[0]

            if admin == 1:
                div += ' <a href="/' + del_link + '/' + url_pas(data[0]) + '">(' + load_lang('delete') + ')</a>'

            div += '</li>'

        div += '</ul>'

        if admin == 1:
            div += '<hr><a href="/' + plus_link + '">(' + load_lang('plus') + ')</a>'
    else:
        if admin == 1:
            div += '<a href="/' + plus_link + '">(' + load_lang('plus') + ')</a>'

    return html_minify(render_template(skin_check(), 
        imp = [title, wiki_set(1), custom(), other2([0, 0])],
        data = div,
        menu = [['other', load_lang('other')]]
    ))

@app.route('/<regex("del_(inter_wiki|html_filter)"):tools>/<name>')
def del_inter(tools = None, name = None):
    if admin_check(None, None) == 1:
        if tools == 'del_inter_wiki':
            curs.execute("delete from inter where title = ?", [name])
        else:
            curs.execute("delete from html_filter where html = ?", [name])
        
        conn.commit()

        return redirect('/' + re.sub('^del_', '', tools))
    else:
        return re_error('/error/3')

@app.route('/<regex("plus_(inter_wiki|html_filter)"):tools>', methods=['POST', 'GET'])
def plus_inter(tools = None):
    if request.method == 'POST':
        if tools == 'plus_inter_wiki':
            curs.execute('insert into inter (title, link) values (?, ?)', [request.form.get('title', None), request.form.get('link', None)])
        else:
            curs.execute('insert into html_filter (html) values (?)', [request.form.get('title', None)])
        
        conn.commit()
        
        admin_check(None, 'inter_wiki_plus')
    
        return redirect('/' + re.sub('^plus_', '', tools))
    else:
        if tools == 'plus_inter_wiki':
            title = load_lang('interwiki') + ' ' + load_lang('plus')
            form_data = '<input placeholder="' + load_lang('name') + '" type="text" name="title"><hr><input placeholder="Link" type="text" name="link">'
        else:
            title = 'HTML Filter ' + load_lang('plus')
            form_data = '<input placeholder="HTML" type="text" name="title">'

        return html_minify(render_template(skin_check(), 
            imp = [title, wiki_set(1), custom(), other2([0, 0])],
            data = '<form method="post">' + form_data + '<hr><button type="submit">' + load_lang('plus') + '</button></form>',
            menu = [['other', load_lang('other')]]
        ))

@app.route('/edit_set')
@app.route('/edit_set/<int:num>', methods=['POST', 'GET'])
def edit_set(num = 0):
    if num != 0 and admin_check(None, None) != 1:
        return re_error('/ban')

    if num == 0:
        li_list = ['Normal', 'Set Text', 'Main HEAD', 'robots.txt', 'Google']
        
        x = 0
        
        li_data = ''
        
        for li in li_list:
            x += 1
            li_data += '<li><a href="/edit_set/' + str(x) + '">' + li + '</a></li>'

        return html_minify(render_template(skin_check(), 
            imp = [load_lang('setting'), wiki_set(1), custom(), other2([0, 0])],
            data = '<h2>' + load_lang('list') + '</h2><ul>' + li_data + '</ul>',
            menu = [['manager', load_lang('admin')]]
        ))
    elif num == 1:
        i_list = ['name', 'logo', 'frontpage', 'license', 'upload', 'skin', 'edit', 'reg', 'ip_view', 'back_up', 'port', 'key']
        n_list = ['Wiki', '', 'FrontPage', 'CC 0', '2', '', 'normal', '', '', '0', '3000', 'Test']
        
        if request.method == 'POST':
            i = 0
            
            for data in i_list:
                curs.execute("update other set data = ? where name = ?", [request.form.get(data, n_list[i]), data])
                i += 1

            conn.commit()

            admin_check(None, 'edit_set')

            return redirect('/edit_set/1')
        else:
            d_list = []
            
            x = 0
            
            for i in i_list:
                curs.execute('select data from other where name = ?', [i])
                sql_d = curs.fetchall()
                if sql_d:
                    d_list += [sql_d[0][0]]
                else:
                    curs.execute('insert into other (name, data) values (?, ?)', [i, n_list[x]])
                    
                    d_list += [n_list[x]]

                x += 1

            conn.commit()
            
            div = ''
            
            if d_list[6] == 'login':
                div += '<option value="login">' + load_lang('subscriber') + '</option>'
                div += '<option value="normal">' + load_lang('normal') + '</option>'
                div += '<option value="admin">' + load_lang('admin') + '</option>'
            elif d_list[6] == 'admin':
                div += '<option value="admin">' + load_lang('admin') + '</option>'
                div += '<option value="login">' + load_lang('subscriber') + '</option>'
                div += '<option value="normal">' + load_lang('normal') + '</option>'
            else:
                div += '<option value="normal">' + load_lang('normal') + '</option>'
                div += '<option value="admin">' + load_lang('admin') + '</option>'
                div += '<option value="login">' + load_lang('subscriber') + '</option>'

            ch_1 = ''
            if d_list[7]:
                ch_1 = 'checked="checked"'

            ch_2 = ''
            if d_list[8]:
                ch_2 = 'checked="checked"'
            
            div2 = ''
            for skin_data in os.listdir(os.path.abspath('views')):
                if d_list[5] == skin_data:
                    div2 = '<option value="' + skin_data + '">' + skin_data + '</option>' + div2
                else:
                    div2 += '<option value="' + skin_data + '">' + skin_data + '</option>'

            return html_minify(render_template(skin_check(), 
                imp = ['Normal', wiki_set(1), custom(), other2([0, 0])],
                data = '''
                        <form method="post">
                            <span>''' + load_lang('name') + '''</span>
                            <br>
                            <br>
                            <input placeholder="''' + load_lang('name') + '''" type="text" name="name" value="''' + html.escape(d_list[0]) + '''">
                            <hr>
                            <span>Logo (HTML)</span>
                            <br>
                            <br>
                            <input placeholder="Logo" type="text" name="logo" value="''' + html.escape(d_list[1]) + '''">
                            <hr>
                            <span>FrontPage</span>
                            <br>
                            <br>
                            <input placeholder="FrontPage" type="text" name="frontpage" value="''' + html.escape(d_list[2]) + '''">
                            <hr>
                            <span>''' + load_lang('license') + ''' (HTML)</span>
                            <br>
                            <br>
                            <input placeholder="''' + load_lang('license') + '''" type="text" name="license" value="''' + html.escape(d_list[3]) + '''">
                            <hr>
                            <span>Max File Size [MB]</span>
                            <br>
                            <br>
                            <input placeholder="Max File Size" type="text" name="upload" value="''' + html.escape(d_list[4]) + '''">
                            <hr>
                            <span>Back Up Interval [''' + load_lang('hour') + '''] (OFF : 0) {Need To Restart}</span>
                            <br>
                            <br>
                            <input placeholder="Back Up Interval" type="text" name="back_up" value="''' + html.escape(d_list[9]) + '''">
                            <hr>
                            <span>Skin</span>
                            <br>
                            <br>
                            <select name="skin">''' + div2 + '''</select>
                            <hr>
                            <span>Main ACL</span>
                            <br>
                            <br>
                            <select name="edit">''' + div + '''</select>
                            <hr>
                            <input type="checkbox" name="reg" ''' + ch_1 + '''> No Register
                            <hr>
                            <input type="checkbox" name="ip_view" ''' + ch_2 + '''> IP Hide
                            <hr>
                            <span>Port</span>
                            <br>
                            <br>
                            <input placeholder="Port" type="text" name="port" value="''' + html.escape(d_list[10]) + '''">
                            <hr>
                            <span>Secret Key</span>
                            <br>
                            <br>
                            <input placeholder="Secret Key" type="password" name="key" value="''' + html.escape(d_list[11]) + '''">
                            <hr>
                            <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        </form>
                        ''',
                menu = [['edit_set', load_lang('setting')]]
            ))
    elif num == 2:
        if request.method == 'POST':
            curs.execute("update other set data = ? where name = ?", [request.form.get('contract', None), 'contract'])
            curs.execute("update other set data = ? where name = ?", [request.form.get('no_login_warring', None), 'no_login_warring'])
            conn.commit()
            
            admin_check(None, 'edit_set')

            return redirect('/edit_set/2')
        else:
            i_list = ['contract', 'no_login_warring']
            n_list = ['', '']
            d_list = []
            
            x = 0
            
            for i in i_list:
                curs.execute('select data from other where name = ?', [i])
                sql_d = curs.fetchall()
                if sql_d:
                    d_list += [sql_d[0][0]]
                else:
                    curs.execute('insert into other (name, data) values (?, ?)', [i, n_list[x]])
                    
                    d_list += [n_list[x]]

                x += 1

            conn.commit()

            return html_minify(render_template(skin_check(), 
                imp = ['Set Text', wiki_set(1), custom(), other2([0, 0])],
                data = '''
                        <form method="post">
                            <span>Register Text</span>
                            <br>
                            <br>
                            <input placeholder="Register Text" type="text" name="contract" value="''' + html.escape(d_list[0]) + '''">
                            <hr>
                            <span>Non-Login Alert</span>
                            <br>
                            <br>
                            <input placeholder="Non-Login Alert" type="text" name="no_login_warring" value="''' + html.escape(d_list[1]) + '''">
                            <hr>
                            <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        </form>
                        ''',
                menu = [['edit_set', load_lang('setting')]]
            ))
    elif num == 3:
        if request.method == 'POST':
            curs.execute("select name from other where name = 'head'")
            if curs.fetchall():
                curs.execute("update other set data = ? where name = 'head'", [request.form.get('content', None)])
            else:
                curs.execute("insert into other (name, data) values ('head', ?)", [request.form.get('content', None)])
            
            conn.commit()

            admin_check(None, 'edit_set')

            return redirect('/edit_set/3')
        else:
            curs.execute("select data from other where name = 'head'")
            head = curs.fetchall()
            if head:
                data = head[0][0]
            else:
                data = ''

            return html_minify(render_template(skin_check(), 
                imp = ['Main HEAD', wiki_set(1), custom(), other2([0, 0])],
                data = '''
                        <form method="post">
                            <textarea rows="25" name="content">''' + html.escape(data) + '''</textarea>
                            <hr>
                            <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        </form>
                        ''',
                menu = [['edit_set', load_lang('setting')]]
            ))
    elif num == 4:
        if request.method == 'POST':
            curs.execute("select name from other where name = 'robot'")
            if curs.fetchall():
                curs.execute("update other set data = ? where name = 'robot'", [request.form.get('content', None)])
            else:
                curs.execute("insert into other (name, data) values ('robot', ?)", [request.form.get('content', None)])
            
            conn.commit()
            
            fw = open('./robots.txt', 'w')
            fw.write(re.sub('\r\n', '\n', request.form.get('content', None)))
            fw.close()
            
            admin_check(None, 'edit_set')

            return redirect('/edit_set/4')
        else:
            curs.execute("select data from other where name = 'robot'")
            robot = curs.fetchall()
            if robot:
                data = robot[0][0]
            else:
                data = ''

            f = open('./robots.txt', 'r')
            lines = f.readlines()
            f.close()

            if not data or data == '':
                data = ''.join(lines)

            return html_minify(render_template(skin_check(), 
                imp = ['robots.txt', wiki_set(1), custom(), other2([0, 0])],
                data = '''
                        <a href="/robots.txt">(View)</a>
                        <hr>
                        <form method="post">
                            <textarea rows="25" name="content">''' + html.escape(data) + '''</textarea>
                            <hr>
                            <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        </form>
                        ''',
                menu = [['edit_set', load_lang('setting')]]
            ))
    elif num == 5:
        if request.method == 'POST':
            curs.execute("update other set data = ? where name = 'recaptcha'", [request.form.get('recaptcha', None)])
            curs.execute("update other set data = ? where name = 'sec_re'", [request.form.get('sec_re', None)])
            conn.commit()
            
            admin_check(None, 'edit_set')

            return redirect('/edit_set/5')
        else:
            i_list = ['recaptcha', 'sec_re']
            n_list = ['', '']
            d_list = []
            
            x = 0
            
            for i in i_list:
                curs.execute('select data from other where name = ?', [i])
                sql_d = curs.fetchall()
                if sql_d:
                    d_list += [sql_d[0][0]]
                else:
                    curs.execute('insert into other (name, data) values (?, ?)', [i, n_list[x]])
                    
                    d_list += [n_list[x]]

                x += 1

            conn.commit()

            return html_minify(render_template(skin_check(), 
                imp = ['Google', wiki_set(1), custom(), other2([0, 0])],
                data = '''
                        <form method="post">
                            <span>reCAPTCHA (HTML)</span>
                            <br>
                            <br>
                            <input placeholder="reCAPTCHA (HTML)" type="text" name="recaptcha" value="''' + html.escape(d_list[0]) + '''">
                            <hr>
                            <span>reCAPTCHA (Secret Key)</span>
                            <br>
                            <br>
                            <input placeholder="reCAPTCHA (Secret Key)" type="text" name="sec_re" value="''' + html.escape(d_list[1]) + '''">
                            <hr>
                            <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        </form>''',
                menu = [['edit_set', load_lang('setting')]]
            ))
    else:
        return redirect('/')

@app.route('/not_close_topic')
def not_close_topic():
    div = '<ul>'
    
    curs.execute('select title, sub from rd order by date desc')
    n_list = curs.fetchall()
    for data in n_list:
        curs.execute('select * from stop where title = ? and sub = ? and close = "O"', [data[0], data[1]])
        is_close = curs.fetchall()
        if not is_close:
            div += '<li><a href="/topic/' + url_pas(data[0]) + '/sub/' + url_pas(data[1]) + '">' + data[0] + ' (' + data[1] + ')</a></li>'
            
    div += '</ul>'

    return html_minify(render_template(skin_check(), 
        imp = [load_lang('open') + ' ' + load_lang('discussion') + ' ' + load_lang('list'), wiki_set(1), custom(), other2([0, 0])],
        data = div,
        menu = [['manager', load_lang('admin')]]
    ))

@app.route('/image/<name>')
def image_view(name = None):
    if os.path.exists(os.path.join('image', name)):
        return send_from_directory('./image', name)
    else:
        return redirect('/')

@app.route('/acl_list')
def acl_list():
    div = '<ul>'
    
    curs.execute("select title, dec from acl where dec = 'admin' or dec = 'user' order by title desc")
    list_data = curs.fetchall()
    for data in list_data:
        if not re.search('^' + load_lang('user') + ':', data[0]) and not re.search('^' + load_lang('file') + ':', data[0]):
            if data[1] == 'admin':
                acl = load_lang('admin')
            else:
                acl = load_lang('subscriber')

            div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a> (' + acl + ')</li>'
        
    div += '</ul>'
    
    return html_minify(render_template(skin_check(), 
        imp = ['ACL ' + load_lang('document') + ' ' + load_lang('list'), wiki_set(1), custom(), other2([0, 0])],
        data = div,
        menu = [['other', load_lang('other')]]
    ))

@app.route('/admin_plus/<name>', methods=['POST', 'GET'])
def admin_plus(name = None):
    if request.method == 'POST':
        if admin_check(None, 'admin_plus (' + name + ')') != 1:
            return re_error('/error/3')

        curs.execute("delete from alist where name = ?", [name])
        
        if request.form.get('ban', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'ban')", [name])

        if request.form.get('mdel', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'mdel')", [name])   

        if request.form.get('toron', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'toron')", [name])
            
        if request.form.get('check', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'check')", [name])

        if request.form.get('acl', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'acl')", [name])

        if request.form.get('hidel', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'hidel')", [name])

        if request.form.get('give', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'give')", [name])

        if request.form.get('owner', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'owner')", [name])
            
        conn.commit()
        
        return redirect('/admin_plus/' + url_pas(name))
    else:        
        data = '<ul>'
        
        exist_list = ['', '', '', '', '', '', '', '']

        curs.execute('select acl from alist where name = ?', [name])
        acl_list = curs.fetchall()    
        for go in acl_list:
            if go[0] == 'ban':
                exist_list[0] = 'checked="checked"'
            elif go[0] == 'mdel':
                exist_list[1] = 'checked="checked"'
            elif go[0] == 'toron':
                exist_list[2] = 'checked="checked"'
            elif go[0] == 'check':
                exist_list[3] = 'checked="checked"'
            elif go[0] == 'acl':
                exist_list[4] = 'checked="checked"'
            elif go[0] == 'hidel':
                exist_list[5] = 'checked="checked"'
            elif go[0] == 'give':
                exist_list[6] = 'checked="checked"'
            elif go[0] == 'owner':
                exist_list[7] = 'checked="checked"'

        if admin_check(None, None) != 1:
            state = 'disabled'
        else:
            state = ''

        data += '<li><input type="checkbox" ' + state +  ' name="ban" ' + exist_list[0] + '> ' + load_lang('ban') + '</li>'
        data += '<li><input type="checkbox" ' + state +  ' name="mdel" ' + exist_list[1] + '> ' + load_lang('bulk_delete') + '</li>'
        data += '<li><input type="checkbox" ' + state +  ' name="toron" ' + exist_list[2] + '> ' + load_lang('discussion') + '</li>'
        data += '<li><input type="checkbox" ' + state +  ' name="check" ' + exist_list[3] + '> ' + load_lang('user') + ' ' + load_lang('check') + '</li>'
        data += '<li><input type="checkbox" ' + state +  ' name="acl" ' + exist_list[4] + '> ' + load_lang('document') + ' ACL</li>'
        data += '<li><input type="checkbox" ' + state +  ' name="hidel" ' + exist_list[5] + '> ' + load_lang('history') + ' ' + load_lang('hide') + '</li>'
        data += '<li><input type="checkbox" ' + state +  ' name="give" ' + exist_list[6] + '> ' + load_lang('authority') + '</li>'
        data += '<li><input type="checkbox" ' + state +  ' name="owner" ' + exist_list[7] + '> ' + load_lang('owner') + '</li></ul>'

        return html_minify(render_template(skin_check(), 
            imp = [load_lang('admin_group') + ' ' + load_lang('plus'), wiki_set(1), custom(), other2([0, 0])],
            data = '<form method="post">' + data + '<hr><button id="save" ' + state +  ' type="submit">' + load_lang('save') + '</button></form>',
            menu = [['manager', load_lang('admin')]]
        ))        
        
@app.route('/admin_list')
def admin_list():
    div = '<ul>'
    
    curs.execute("select id, acl, date from user where not acl = 'user' order by date desc")
    for data in curs.fetchall():
        name = ip_pas(data[0]) + ' <a href="/admin_plus/' + url_pas(data[1]) + '">(' + data[1] + ')</a>'
        
        if data[2] != '':
            name += '(' + load_lang('register') + ' : ' + data[2] + ')'

        div += '<li>' + name + '</li>'
        
    div += '</ul>'
                
    return html_minify(render_template(skin_check(), 
        imp = [load_lang('admin') + ' ' + load_lang('list'), wiki_set(1), custom(), other2([0, 0])],
        data = div,
        menu = [['other', load_lang('other')]]
    ))
        
@app.route('/hidden/<path:name>')
def history_hidden(name = None):
    num = int(request.args.get('num', 0))

    if admin_check(6, 'history_hidden (' + name + '#' + str(num) + ')') == 1:
        curs.execute("select title from history where title = ? and id = ? and hide = 'O'", [name, str(num)])
        if curs.fetchall():
            curs.execute("update history set hide = '' where title = ? and id = ?", [name, str(num)])
        else:
            curs.execute("update history set hide = 'O' where title = ? and id = ?", [name, str(num)])
            
        conn.commit()
    
    return redirect('/history/' + url_pas(name))
        
@app.route('/user_log')
def user_log():
    num = int(request.args.get('num', 1))
    if num * 50 > 0:
        sql_num = num * 50 - 50
    else:
        sql_num = 0
        
    list_data = '<ul>'

    admin_one = admin_check(1, None)
    
    curs.execute("select id, date from user order by date desc limit ?, '50'", [str(sql_num)])
    user_list = curs.fetchall()
    for data in user_list:
        if admin_one == 1:
            curs.execute("select block from ban where block = ?", [data[0]])
            if curs.fetchall():
                ban_button = ' <a href="/ban/' + url_pas(data[0]) + '">(' + load_lang('release') + ')</a>'
            else:
                ban_button = ' <a href="/ban/' + url_pas(data[0]) + '">(' + load_lang('ban') + ')</a>'
        else:
            ban_button = ''
            
        list_data += '<li>' + ip_pas(data[0]) + ban_button
        
        if data[1] != '':
            list_data += ' (' + load_lang('register') + ' : ' + data[1] + ')'

        list_data += '</li>'

    if num == 1:
        curs.execute("select count(id) from user")
        user_count = curs.fetchall()
        if user_count:
            count = user_count[0][0]
        else:
            count = 0

        list_data += '</ul><hr><ul><li>All : ' + str(count) + '</li></ul>'

    list_data += next_fix('/user_log?num=', num, user_list)

    return html_minify(render_template(skin_check(), 
        imp = ['' + load_lang('recent') + ' ' + load_lang('register') + '', wiki_set(1), custom(), other2([0, 0])],
        data = list_data,
        menu = 0
    ))

@app.route('/admin_log')
def admin_log():
    num = int(request.args.get('num', 1))
    if num * 50 > 0:
        sql_num = num * 50 - 50
    else:
        sql_num = 0

    list_data = '<ul>'

    curs.execute("select who, what, time from re_admin order by time desc limit ?, '50'", [str(sql_num)])
    get_list = curs.fetchall()
    for data in get_list:            
        list_data += '<li>' + ip_pas(data[0]) + ' / ' + data[1] + ' / ' + data[2] + '</li>'

    list_data += '</ul>'
    list_data += next_fix('/admin_log?num=', num, get_list)

    return html_minify(render_template(skin_check(), 
        imp = ['' + load_lang('recent') + ' ' + load_lang('authority'), wiki_set(1), custom(), other2([0, 0])],
        data = list_data,
        menu = 0
    ))

@app.route('/give_log')
def give_log():        
    list_data = '<ul>'
    back = ''

    curs.execute("select distinct name from alist order by name asc")
    for data in curs.fetchall():                      
        if back != data[0]:
            back = data[0]

        list_data += '<li><a href="/admin_plus/' + url_pas(data[0]) + '">' + data[0] + '</a></li>'
    
    list_data += '</ul><hr><a href="/manager/8">(' + load_lang('create') + ')</a>'

    return html_minify(render_template(skin_check(), 
        imp = [load_lang('admin_group') + ' ' + load_lang('list'), wiki_set(1), custom(), other2([0, 0])],
        data = list_data,
        menu = [['other', load_lang('other')]]
    ))

@app.route('/indexing')
def indexing():
    if admin_check(None, 'indexing') != 1:
        return re_error('/error/3')

    print('')

    curs.execute("select name from sqlite_master where type = 'index'")
    data = curs.fetchall()
    if data:
        for delete_index in data:
            print('Delete : ' + delete_index[0])

            sql = 'drop index if exists ' + delete_index[0]
            
            try:
                curs.execute(sql)
            except:
                pass
    else:
        curs.execute("select name from sqlite_master where type in ('table', 'view') and name not like 'sqlite_%' union all select name from sqlite_temp_master where type in ('table', 'view') order by 1;")
        for table in curs.fetchall():            
            curs.execute('select sql from sqlite_master where name = ?', [table[0]])
            cul = curs.fetchall()
            
            r_cul = re.findall('(?:([^ (]*) text)', str(cul[0]))
            
            for n_cul in r_cul:
                print('Create : index_' + table[0] + '_' + n_cul)

                sql = 'create index index_' + table[0] + '_' + n_cul + ' on ' + table[0] + '(' + n_cul + ')'
                try:
                    curs.execute(sql)
                except:
                    pass

    conn.commit()
    
    print('')

    return redirect('/')        

@app.route('/re_start')
def re_start():
    if admin_check(None, 're_start') != 1:
        return re_error('/error/3')

    print('')
    print('Re Start')
    print('')

    os.execl(sys.executable, sys.executable, *sys.argv)

@app.route('/update')
def update():
    if admin_check(None, 'update') != 1:
       return re_error('/error/3')

    if platform.system() == 'Linux':
        print('')
        print('Update')

        ok = os.system('git pull')
        if ok == 0:
            return redirect('/re_start')
    else:
        if platform.system() == 'Windows':
            print('')
            print('Download')

            urllib.request.urlretrieve('https://github.com/2DU/openNAMU/archive/stable.zip', 'update.zip')

            print('Zip Extract')
            zipfile.ZipFile('update.zip').extractall('')

            print('Move')
            ok = os.system('xcopy /y /r openNAMU-stable .')
            if ok == 0:
                print('Remove')
                os.system('rd /s /q openNAMU-stable')
                os.system('del update.zip')

                return redirect('/re_start')

    return html_minify(render_template(skin_check(), 
        imp = [load_lang('update'), wiki_set(1), custom(), other2([0, 0])],
        data = 'Auto Update Is Not Support. <a href="https://github.com/2DU/openNAMU">(GitHub)</a>',
        menu = [['manager/1', load_lang('admin')]]
    ))
        
@app.route('/xref/<path:name>')
def xref(name = None):
    num = int(request.args.get('num', 1))
    if num * 50 > 0:
        sql_num = num * 50 - 50
    else:
        sql_num = 0
        
    div = '<ul>'
    
    curs.execute("select link, type from back where title = ? and not type = 'cat' and not type = 'no' order by link asc limit ?, '50'", [name, str(sql_num)])
    data_list = curs.fetchall()
    for data in data_list:
        div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a>'
        
        if data[1]:
            if data[1] == 'include':
                side = '포함'
            elif data[1] == 'file':
                side = '' + load_lang('file') + ''
            else:
                side = '넘겨주기'
                
            div += ' (' + side + ')'
        
        div += '</li>'
        
        if re.search('^틀:', data[0]):
            div += '<li><a id="inside" href="/xref/' + url_pas(data[0]) + '">' + data[0] + '</a> (' + load_lang('backlink') + ')</li>'
      
    div += '</ul>' + next_fix('/xref/' + url_pas(name) + '?num=', num, data_list)
    
    return html_minify(render_template(skin_check(), 
        imp = [name, wiki_set(1), custom(), other2([' (' + load_lang('backlink') + ')', 0])],
        data = div,
        menu = [['w/' + url_pas(name), load_lang('document')]]
    ))

@app.route('/please')
def please():
    num = int(request.args.get('num', 1))
    if num * 50 > 0:
        sql_num = num * 50 - 50
    else:
        sql_num = 0
        
    div = '<ul>'
    var = ''
    
    curs.execute("select distinct title from back where type = 'no' order by title asc limit ?, '50'", [str(sql_num)])
    data_list = curs.fetchall()
    for data in data_list:
        if var != data[0]:
            div += '<li><a id="not_thing" href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a></li>'   

            var = data[0]
        
    div += '</ul>' + next_fix('/please?num=', num, data_list)
    
    return html_minify(render_template(skin_check(), 
        imp = [load_lang('need') + ' ' + load_lang('document'), wiki_set(1), custom(), other2([0, 0])],
        data = div,
        menu = [['other', load_lang('other')]]
    ))
        
@app.route('/recent_discuss')
@app.route('/recent_discuss/<regex("close"):tools>')
def recent_discuss(tools = 'normal'):
    if tools == 'normal' or tools == 'close':
        div = ''
        
        if tools == 'normal':
            div += '<a href="/recent_discuss/close">(' + load_lang('close') + ')</a>'
           
            m_sub = 0
        else:
            div += '<a href="/recent_discuss">(' + load_lang('open') + ')</a>'
            
            m_sub = ' (' + load_lang('close') + ')'

        div += '<hr><table style="width: 100%; text-align: center;"><tbody><tr><td style="width: 50%;">' + load_lang('discussion') + ' ' + load_lang('name') + '</td><td style="width: 50%;">' + load_lang('time') + '</td></tr>'
    else:
        return redirect('/')
    
    curs.execute("select title, sub, date from rd order by date desc limit 50")
    for data in curs.fetchall():
        title = html.escape(data[0])
        sub = html.escape(data[1])
        
        close = 0
        
        if tools == 'normal':
            curs.execute("select title from stop where title = ? and sub = ? and close = 'O'", [data[0], data[1]])
            if curs.fetchall():
                close = 1
        else:
            curs.execute("select title from stop where title = ? and sub = ? and close = 'O'", [data[0], data[1]])
            if not curs.fetchall():
                close = 1

        if close == 0:
            div += '<tr><td><a href="/topic/' + url_pas(data[0]) + '/sub/' + url_pas(data[1]) + '">' + title + '</a> (' + sub + ')</td><td>' + data[2] + '</td></tr>'
    else:
        div += '</tbody></table>'
            
    return html_minify(render_template(skin_check(), 
        imp = ['' + load_lang('recent') + ' ' + load_lang('discussion'), wiki_set(1), custom(), other2([m_sub, 0])],
        data = div,
        menu = 0
    ))

@app.route('/block_log')
@app.route('/block_log/<regex("ip|user|never_end|can_end|end|now|edit_filter"):tool2>')
@app.route('/<regex("block_user|block_admin"):tool>/<name>')
def block_log(name = None, tool = None, tool2 = None):
    num = int(request.args.get('num', 1))
    if num * 50 > 0:
        sql_num = num * 50 - 50
    else:
        sql_num = 0
    
    div = '<table style="width: 100%; text-align: center;"><tbody><tr><td style="width: 33.3%;">차단자</td><td style="width: 33.3%;">' + load_lang('admin') + '</td><td style="width: 33.3%;">기간</td></tr>'
    
    data_list = ''
    
    if not name:
        if not tool2:
            div = '''
                    <a href="/manager/11">(차단자)</a> <a href="/manager/12">(''' + load_lang('admin') + ''')</a>
                    <hr>
                    <a href="/block_log/ip">(IP)</a> <a href="/block_log/user">(''' + load_lang('subscriber') + ''')</a> <a href="/block_log/never_end">(무기한)</a> <a href="/block_log/can_end">(기간)</a> <a href="/block_log/end">(''' + load_lang('release') + ''')</a> <a href="/block_log/now">(현재)</a> <a href="/block_log/edit_filter">(''' + load_lang('edit_filter') + ''')</a>
                    <hr>
                    ''' + div
            
            sub = 0
            menu = 0
            
            curs.execute("select why, block, blocker, end, today from rb order by today desc limit ?, '50'", [str(sql_num)])
        else:
            menu = [['block_log', load_lang('normal')]]
            
            if tool2 == 'ip':
                sub = ' (IP)'
                
                curs.execute("select why, block, blocker, end, today from rb where (block like ? or block like ?) order by today desc limit ?, '50'", ['%.%', '%:%', str(sql_num)])
            elif tool2 == 'user':
                sub = ' (' + load_lang('subscriber') + ')'
                
                curs.execute("select why, block, blocker, end, today from rb where not (block like ? or block like ?) order by today desc limit ?, '50'", ['%.%', '%:%', str(sql_num)])
            elif tool2 == 'never_end':
                sub = '(무기한)'
                
                curs.execute("select why, block, blocker, end, today from rb where not end like ? and not end like ? order by today desc limit ?, '50'", ['%:%', '%' + load_lang('release') + '%', str(sql_num)])
            elif tool2 == 'end':
                sub = '(' + load_lang('release') + ')'
                
                curs.execute("select why, block, blocker, end, today from rb where end = ? order by today desc limit ?, '50'", [load_lang('release'), str(sql_num)])
            elif tool2 == 'now':
                sub = '(현재)'
                
                data_list = []
                
                curs.execute("select block from ban limit ?, '50'", [str(sql_num)])
                for in_data in curs.fetchall():
                    curs.execute("select why, block, blocker, end, today from rb where block = ? order by today desc limit 1", [in_data[0]])
                    
                    data_list = [curs.fetchall()[0]] + data_list
            elif tool2 == 'edit_filter':
                sub = '(' + load_lang('edit_filter') + ')'

                curs.execute("select why, block, blocker, end, today from rb where blocker = ? order by today desc limit ?, '50'", [load_lang('tool') + ':' + load_lang('edit_filter'), str(sql_num)])
            else:
                sub = '(기간)'
                
                curs.execute("select why, block, blocker, end, today from rb where end like ? order by today desc limit ?, '50'", ['%\-%', str(sql_num)])
    else:
        menu = [['block_log', load_lang('normal')]]
        
        if tool == 'block_user':
            sub = ' (차단자)'
            
            curs.execute("select why, block, blocker, end, today from rb where block = ? order by today desc limit ?, '50'", [name, str(sql_num)])
        else:
            sub = ' (' + load_lang('admin') + ')'
            
            curs.execute("select why, block, blocker, end, today from rb where blocker = ? order by today desc limit ?, '50'", [name, str(sql_num)])

    if data_list == '':
        data_list = curs.fetchall()

    for data in data_list:
        why = html.escape(data[0])
        if why == '':
            why = '<br>'
        
        band = re.search("^([0-9]{1,3}\.[0-9]{1,3})$", data[1])
        if band:
            ip = data[1] + ' (대역)'
        else:
            ip = ip_pas(data[1])

        if data[3] != '':
            end = data[3]
        else:
            end = '무기한'
            
        div += '<tr><td>' + ip + '</td><td>' + ip_pas(data[2]) + '</td><td>시작 : ' + data[4] + '<br>끝 : ' + end + '</td></tr>'
        div += '<tr><td colspan="3">' + why + '</td></tr>'

    div += '</tbody></table>'
    
    if not name:
        if not tool2:
            div += next_fix('/block_log?num=', num, data_list)
        else:
            div += next_fix('/block_log/' + url_pas(tool2) + '?num=', num, data_list)
    else:
        div += next_fix('/' + url_pas(tool) + '/' + url_pas(name) + '?num=', num, data_list)
                
    return html_minify(render_template(skin_check(), 
        imp = ['' + load_lang('recent') + ' ' + load_lang('ban'), wiki_set(1), custom(), other2([sub, 0])],
        data = div,
        menu = menu
    ))
            
@app.route('/search', methods=['POST'])
def search():
    return redirect('/search/' + url_pas(request.form.get('search', None)))

@app.route('/goto', methods=['POST'])
def goto():
    curs.execute("select title from data where title = ?", [request.form.get('search', None)])
    data = curs.fetchall()
    if data:
        return redirect('/w/' + url_pas(request.form.get('search', None)))
    else:
        return redirect('/search/' + url_pas(request.form.get('search', None)))

@app.route('/search/<path:name>')
def deep_search(name = None):
    num = int(request.args.get('num', 1))
    if num * 50 > 0:
        sql_num = num * 50 - 50
    else:
        sql_num = 0

    div = '<ul>'
    
    div_plus = ''
    no = 0
    start = 2
    test = ''
    
    curs.execute("select title from data where title = ?", [name])
    if curs.fetchall():
        div = '<ul><li><a href="/w/' + url_pas(name) + '">' + name + '</a></li></ul><hr><ul>'
    else:
        div = '<ul><li><a id="not_thing" href="/w/' + url_pas(name) + '">' + name + '</a></li></ul><hr><ul>'

    curs.execute("select distinct title, case when title like ? then '제목' else '내용' end from data where title like ? or data like ? order by case when title like ? then 1 else 2 end limit ?, '50'", ['%' + name + '%', '%' + name + '%', '%' + name + '%', '%' + name + '%', str(sql_num)])
    all_list = curs.fetchall()
    if all_list:
        test = all_list[0][1]
        
        for data in all_list:
            if data[1] != test:
                div_plus += '</ul><hr><ul>'
                
                test = data[1]

            div_plus += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a> (' + data[1] + ')</li>'
    else:
        div += '<li>Not Found.</li>'

    div += div_plus + '</ul>'
    div += next_fix('/search/' + url_pas(name) + '?num=', num, all_list)

    return html_minify(render_template(skin_check(), 
        imp = [name, wiki_set(1), custom(), other2([' (' + load_lang('search') + ')', 0])],
        data = div,
        menu = 0
    ))
         
@app.route('/raw/<path:name>')
@app.route('/topic/<path:name>/sub/<sub_title>/raw/<int:num>')
def raw_view(name = None, sub_title = None, num = None):
    v_name = name
    sub = ' (Raw)'
    
    if not num:
        num = request.args.get('num', None)
        if num:
            num = int(num)
    
    if not sub_title and num:
        curs.execute("select title from history where title = ? and id = ? and hide = 'O'", [name, str(num)])
        if curs.fetchall() and admin_check(6, None) != 1:
            return re_error('/error/3')
        
        curs.execute("select data from history where title = ? and id = ?", [name, str(num)])
        
        sub += ' (' + str(num) + load_lang('version') + ')'

        menu = [['history/' + url_pas(name), load_lang('history')]]
    elif sub_title:
        curs.execute("select data from topic where id = ? and title = ? and sub = ? and block = ''", [str(num), name, sub_title])
        
        v_name = load_lang('discussion') + ' Raw'
        sub = ' (' + str(num) + '번)'

        menu = [['topic/' + url_pas(name) + '/sub/' + url_pas(sub_title) + '#' + str(num), load_lang('discussion')], ['topic/' + url_pas(name) + '/sub/' + url_pas(sub_title) + '/admin/' + str(num), load_lang('tool')]]
    else:
        curs.execute("select data from data where title = ?", [name])
        
        menu = [['w/' + url_pas(name), load_lang('document')]]

    data = curs.fetchall()
    if data:
        p_data = html.escape(data[0][0])
        p_data = '<textarea readonly rows="25">' + p_data + '</textarea>'
        
        return html_minify(render_template(skin_check(), 
            imp = [v_name, wiki_set(1), custom(), other2([sub, 0])],
            data = p_data,
            menu = menu
        ))
    else:
        return redirect('/w/' + url_pas(name))
        
@app.route('/revert/<path:name>', methods=['POST', 'GET'])
def revert(name = None):    
    num = int(request.args.get('num', 0))

    curs.execute("select title from history where title = ? and id = ? and hide = 'O'", [name, str(num)])
    if curs.fetchall() and admin_check(6, None) != 1:
        return re_error('/error/3')

    if acl_check(name) == 1:
        return re_error('/ban')

    if request.method == 'POST':
        if captcha_post(request.form.get('g-recaptcha-response', None), conn) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        curs.execute("delete from back where link = ?", [name])
        conn.commit()
        
        curs.execute("select data from history where title = ? and id = ?", [name, str(num)])
        data = curs.fetchall()
        if data:                                
            curs.execute("select data from data where title = ?", [name])
            data_old = curs.fetchall()
            if data_old:
                leng = leng_check(len(data_old[0][0]), len(data[0][0]))
                curs.execute("update data set data = ? where title = ?", [data[0][0], name])
            else:
                leng = ' +' + str(len(data[0][0]))
                curs.execute("insert into data (title, data) values (?, ?)", [name, data[0][0]])
                
            history_plus(name, data[0][0], get_time(), ip_check(), request.form.get('send', None) + ' (' + str(num) + load_lang('version') + ')', leng)
            namumark(name, data[0][0], 1)
            
            conn.commit()
            
            return redirect('/w/' + url_pas(name))
    else:
        curs.execute("select title from history where title = ? and id = ?", [name, str(num)])
        if not curs.fetchall():
            return redirect('/w/' + url_pas(name))

        return html_minify(render_template(skin_check(), 
            imp = [name, wiki_set(1), custom(), other2([' (' + load_lang('revert') + ')', 0])],
            data =  '<form method="post"><span>' + request.args.get('num', '0') + load_lang('version') + '</span><hr>' + ip_warring() + '<input placeholder="' + load_lang('why') + '" name="send" type="text"><hr>' + captcha_get() + '<button type="submit">' + load_lang('revert') + '</button></form>',
            menu = [['history/' + url_pas(name), load_lang('history')], ['recent_changes', '' + load_lang('recent') + ' ' + load_lang('change') + '']]
        ))            
                    
@app.route('/big_delete', methods=['POST', 'GET'])
def big_delete():
    if admin_check(2, 'big_delete') != 1:
        return re_error('/error/3')

    if request.method == 'POST':
        today = get_time()
        ip = ip_check()
        data = request.form.get('content', None) + '\r\n'
        
        match = re.findall('(.*)\r\n', data)
        for list_one in match:
            curs.execute("select data from data where title = ?", [list_one])
            data_old = curs.fetchall()
            if data_old:
                curs.execute("delete from back where title = ?", [list_one])
                curs.execute("delete from data where title = ?", [list_one])
                
                leng = '-' + str(len(data_old[0][0]))
                
                history_plus(list_one, '', today, ip, request.form.get('send', None) + ' (' + load_lang('bulk_delete') + ')', leng)

            data = re.sub('(.*)\r\n', '', data, 1)
        
        conn.commit()

        return redirect('/')
    else:
        return html_minify(render_template(skin_check(), 
            imp = [load_lang('bulk_delete'), wiki_set(1), custom(), other2([0, 0])],
            data = '''
                    <span>
                        Title A
                        <br>
                        Title B
                        <br>
                        Title C
                    </span>
                    <hr>
                    <form method="post">
                        <textarea rows="25" name="content"></textarea>
                        <hr>
                        <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                        <hr>
                        <button type="submit">''' + load_lang('delete') + '''</button>
                    </form>
                    ''',
            menu = [['manager', load_lang('admin')]]
        ))

@app.route('/edit_filter')
def edit_filter():
    div = '<ul>'
    
    curs.execute("select name from filter")
    data = curs.fetchall()
    for data_list in data:
        div += '<li><a href="/edit_filter/' + url_pas(data_list[0]) + '">' + data_list[0] + '</a></li>'

    div += '</ul>'

    if data:
        div += '<hr><a href="/manager/9">(' + load_lang('plus') + ')</a>'
    else:
        div = '<a href="/manager/9">(' + load_lang('plus') + ')</a>'

    return html_minify(render_template(skin_check(), 
        imp = [load_lang('edit_filter') + ' ' + load_lang('list'), wiki_set(1), custom(), other2([0, 0])],
        data = div,
        menu = [['manager', load_lang('admin')]]
    ))

@app.route('/edit_filter/<name>/delete', methods=['POST', 'GET'])
def delete_edit_filter(name = None):
    if admin_check(1, 'edit_filter delete') != 1:
        return re_error('/error/3')

    curs.execute("delete from filter where name = ?", [name])
    conn.commit()

    return redirect('/edit_filter')

@app.route('/edit_filter/<name>', methods=['POST', 'GET'])
def set_edit_filter(name = None):
    if request.method == 'POST':
        if admin_check(1, 'edit_filter edit') != 1:
            return re_error('/error/3')

        if request.form.get('ban', None):
            end = 'X'
        else:
            end = ''

        curs.execute("select name from filter where name = ?", [name])
        if curs.fetchall():
            curs.execute("update filter set regex = ?, sub = ? where name = ?", [request.form.get('content', '테스트'), end, name])
        else:
            curs.execute("insert into filter (name, regex, sub) values (?, ?, ?)", [name, request.form.get('content', '테스트'), end])

        conn.commit()
    
        return redirect('/edit_filter/' + url_pas(name))
    else:
        curs.execute("select regex, sub from filter where name = ?", [name])
        exist = curs.fetchall()
        if exist:
            textarea = exist[0][0]
            
            if exist[0][1] == 'X':
                time_data = 'checked="checked"'
            else:
                time_data = ''
        else:
            textarea = ''
            time_data = ''

        if admin_check(1, None) != 1:
            stat = 'disabled'
        else:
            stat = ''

        return html_minify(render_template(skin_check(), 
            imp = [name, wiki_set(1), custom(), other2([' (' + load_lang('edit_filter') + ')', 0])],
            data = '''
                    <form method="post">
                        <input ''' + stat + ''' type="checkbox" ''' + time_data + ''' name="ban">
                        ''' + load_lang('ban') + '''
                        <hr>
                        <input ''' + stat + ''' placeholder="정규식" name="content" value="''' + html.escape(textarea) + '''" type="text">
                        <hr>
                        <button ''' + stat + ''' id="save" type="submit">''' + load_lang('save') + '''</button>
                    </form>''',
            menu = [['edit_filter', load_lang('list')], ['edit_filter/' + url_pas(name) + '/delete', load_lang('delete')]]
        ))

@app.route('/edit/<path:name>', methods=['POST', 'GET'])
def edit(name = None):
    ip = ip_check()
    if acl_check(name) == 1:
        return re_error('/ban')
    
    if request.method == 'POST':
        if admin_check(1, 'edit_filter pass') != 1:
            curs.execute("select regex, sub from filter")
            for data_list in curs.fetchall():
                match = re.compile(data_list[0])
                if match.search(request.form.get('content', None)):
                    if data_list[1] == 'X':
                        ban_insert(ip, '', load_lang('edit_filter'), None, load_lang('tool') + ':' + load_lang('edit_filter'))
                    
                    return re_error('/error/21')

        if captcha_post(request.form.get('g-recaptcha-response', None), conn) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        if len(request.form.get('send', None)) > 500:
            return re_error('/error/15')

        if request.form.get('otent', None) == request.form.get('content', None):
            return re_error('/error/18')

        today = get_time()
        content = savemark(request.form.get('content', None))
        
        curs.execute("select data from data where title = ?", [name])
        old = curs.fetchall()
        if old:
            leng = leng_check(len(request.form.get('otent', None)), len(content))
            
            if request.args.get('section', None):
                i = 1
                
                data = re.sub('\r\n', '\n', '\r\n' + old[0][0] + '\r\n')
                while 1:
                    replace_data = re.search('\n(={1,6}) ?((?:(?!=).)+) ?={1,6}\n', data)
                    if replace_data:
                        replace_data = replace_data.groups()[0]

                        if i == int(request.args.get('section', None)):
                            data = re.sub('\n(?P<in>={1,6}) ?(?P<out>(?:(?!=).)+) ?={1,6}\n', '\n<real h' + str(len(replace_data)) + '>\g<out></real h' + str(len(replace_data)) + '>\n', data, 1)
                        else:
                            data = re.sub('\n(?P<in>={1,6}) ?(?P<out>(?:(?!=).)+) ?={1,6}\n', '\n<h' + str(len(replace_data)) + '>\g<out></h' + str(len(replace_data)) + '>\n', data, 1)

                        i += 1
                    else:
                        break

                new_data = re.sub('\r\n', '\n', '\r\n' + request.form.get('otent', None) + '\r\n')
                while 1:
                    replace_data = re.search('\n(={1,6}) ?((?:(?!=).)+) ?={1,6}\n', new_data)
                    if replace_data:
                        replace_data = replace_data.groups()[0]

                        new_data = re.sub('\n(?P<in>={1,6}) ?(?P<out>(?:(?!=).)+) ?={1,6}\n', '\n<real h' + str(len(replace_data)) + '>\g<out></real h' + str(len(replace_data)) + '>\n', new_data, 1)
                    else:
                        break

                content = data.replace(new_data, '\n' + content + '\n')

                while 1:
                    replace_data = re.search('\n<(?:real )?h([1-6])>((?:(?!<h).)+) ?<\/(?:real )?h[1-6]>\n', content)
                    if replace_data:
                        replace_data = replace_data.groups()[0]

                        content = re.sub('\n<(?:real )?h([1-6])>(?P<out>(?:(?!<h).)+) ?<\/(?:real )?h[1-6]>\n', '\n' + ('=' * int(replace_data)) + ' \g<out> ' + ('=' * int(replace_data)) + '\n', content, 1)
                    else:
                        break

                content = re.sub('^\n', '', content)
                content = re.sub('\n$', '', content)
                
            curs.execute("update data set data = ? where title = ?", [content, name])
        else:
            leng = ' +' + str(len(content))
            
            curs.execute("insert into data (title, data) values (?, ?)", [name, content])

        curs.execute("select user from scan where title = ?", [name])
        for user_data in curs.fetchall():
            curs.execute("insert into alarm (name, data, date) values (?, ?, ?)", [ip, ip + ' - <a href="/w/' + url_pas(name) + '">' + name + '</a> (Edit)', today])

        history_plus(name, content, today, ip, send_parser(request.form.get('send', None)), leng)
        
        curs.execute("delete from back where link = ?", [name])
        curs.execute("delete from back where title = ? and type = 'no'", [name])
        
        namumark(name, content, 1)
        
        conn.commit()
        
        return redirect('/w/' + url_pas(name))
    else:            
        curs.execute("select data from data where title = ?", [name])
        new = curs.fetchall()
        if new:
            if request.args.get('section', None):
                test_data = '\n' + re.sub('\r\n', '\n', new[0][0]) + '\n'   
                
                section_data = re.findall('((?:={1,6}) ?(?:(?:(?!=).)+) ?={1,6}\n(?:(?:(?!(?:={1,6}) ?(?:(?:(?!=).)+) ?={1,6}\n).)*\n*)*)', test_data)
                data = section_data[int(request.args.get('section', None)) - 1]
            else:
                data = new[0][0]
        else:
            data = ''
            
        data_old = data
        
        if not request.args.get('section', None):
            get_name = '''
                        <form method="post" id="get_edit" action="/edit_get/''' + url_pas(name) + '''">
                            <input placeholder="Load Document" name="name" style="width: 50%;" type="text">
                            <button id="come" type="submit">Load</button>
                        </form>
                        <hr>
                        '''
            action = ''
        else:
            get_name = ''
            action = '?section=' + request.args.get('section', None)
            
        if request.args.get('froms', None):
            curs.execute("select data from data where title = ?", [request.args.get('froms', None)])
            get_data = curs.fetchall()
            if get_data:
                data = get_data[0][0]
                get_name = ''

        js_data = edit_help_button()

        return html_minify(render_template(skin_check(), 
            imp = [name, wiki_set(1), custom(), other2([' (' + load_lang('edit') + ')', 0])],
            data = get_name + js_data[0] + '''
                    <form method="post" action="/edit/''' + url_pas(name) + action + '''">
                        ''' + js_data[1] + '''
                        <textarea id="content" rows="25" name="content">''' + html.escape(re.sub('\n$', '', data)) + '''</textarea>
                        <textarea style="display: none;" name="otent">''' + html.escape(re.sub('\n$', '', data_old)) + '''</textarea>
                        <hr>
                        <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                        <hr>
                        ''' + captcha_get() + ip_warring() + '''
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        <button id="preview" type="submit" formaction="/preview/''' + url_pas(name) + action + '">' + load_lang('preview') + '''</button>
                    </form>
                    ''',
            menu = [['w/' + url_pas(name), load_lang('document')], ['delete/' + url_pas(name), load_lang('delete')], ['move/' + url_pas(name), load_lang('move')]]
        ))
        
@app.route('/edit_get/<path:name>', methods=['POST'])
def edit_get(name = None):
    return redirect('/edit/' + url_pas(name) + '?froms=' + url_pas(request.form.get('name', None)))

@app.route('/preview/<path:name>', methods=['POST'])
def preview(name = None):
    ip = ip_check()
    if acl_check(name) == 1:
        return re_error('/ban')
         
    new_data = re.sub('\r\n#(?:redirect|넘겨주기) (?P<in>(?:(?!\r\n).)+)\r\n', ' * Redirect : [[\g<in>]]', '\r\n' + request.form.get('content', None) + '\r\n')
    new_data = re.sub('^\r\n', '', new_data)
    new_data = re.sub('\r\n$', '', new_data)
    
    end_data = namumark(name, new_data, 0)
    
    if request.args.get('section', None):
        action = '?section=' + request.args.get('section', None)
    else:
        action = ''

    js_data = edit_help_button()
    
    return html_minify(render_template(skin_check(), 
        imp = [name, wiki_set(1), custom(), other2([' (' + load_lang('preview') + ')', 0])],
        data = js_data[0] + '''
                <form method="post" action="/edit/''' + url_pas(name) + action + '''">
                    ''' + js_data[1] + '''
                    <textarea id="content" rows="25" name="content">''' + html.escape(request.form.get('content', None)) + '''</textarea>
                    <textarea style="display: none;" name="otent">''' + html.escape(request.form.get('otent', None)) + '''</textarea>
                    <hr>
                    <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                    <hr>
                    ''' + captcha_get() + '''
                    <button id="save" type="submit">''' + load_lang('save') + '''</button>
                    <button id="preview" type="submit" formaction="/preview/''' + url_pas(name) + action + '">' + load_lang('preview') + '''</button>
                </form>
                <hr>
                ''' + end_data,
        menu = [['w/' + url_pas(name), load_lang('document')]]
    ))
        
@app.route('/delete/<path:name>', methods=['POST', 'GET'])
def delete(name = None):
    ip = ip_check()
    if acl_check(name) == 1:
        return re_error('/ban')
    
    if request.method == 'POST':
        if captcha_post(request.form.get('g-recaptcha-response', None), conn) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        curs.execute("select data from data where title = ?", [name])
        data = curs.fetchall()
        if data:
            today = get_time()
            leng = '-' + str(len(data[0][0]))
            
            history_plus(name, '', today, ip, request.form.get('send', None) + ' (' + load_lang('delete') + ')', leng)
            
            curs.execute("select title, link from back where title = ? and not type = 'cat' and not type = 'no'", [name])
            for data in curs.fetchall():
                curs.execute("insert into back (title, link, type) values (?, ?, 'no')", [data[0], data[1]])
            
            curs.execute("delete from back where link = ?", [name])
            curs.execute("delete from data where title = ?", [name])
            conn.commit()
            
        return redirect('/w/' + url_pas(name))
    else:
        curs.execute("select title from data where title = ?", [name])
        if not curs.fetchall():
            return redirect('/w/' + url_pas(name))

        return html_minify(render_template(skin_check(), 
            imp = [name, wiki_set(1), custom(), other2([' (' + load_lang('delete') + ')', 0])],
            data = '''
                    <form method="post">
                        ''' + ip_warring() + '''
                        <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                        <hr>
                        ''' + captcha_get() + '''
                        <button type="submit">''' + load_lang('delete') + '''</button>
                    </form>
                    ''',
            menu = [['w/' + url_pas(name), load_lang('document')]]
        ))            
            
@app.route('/move_data/<path:name>')
def move_data(name = None):    
    data = '<ul>'
    
    curs.execute("select send, date, ip from history where send like ? or send like ? order by date desc", ['%<a href="/w/' + url_pas(name) + '">' + name + '</a> ' + load_lang('move') + ')%', '%(<a href="/w/' + url_pas(name) + '">' + name + '</a>%'])
    for for_data in curs.fetchall():
        match = re.findall('<a href="\/w\/(?:(?:(?!">).)+)">((?:(?!<\/a>).)+)<\/a>', for_data[0])
        send = re.sub('\([^\)]+\)$', '', for_data[0])
        data += '<li><a href="/move_data/' + url_pas(match[0]) + '">' + match[0] + '</a> - <a href="/move_data/' + url_pas(match[1]) + '">' + match[1] + '</a>'
        
        if re.search('^( *)+$', send):
            data += ' / ' + for_data[2] + ' / ' + for_data[1] + '</li>'
        else:
            data += ' / ' + for_data[2] + ' / ' + for_data[1] + ' / ' + send + '</li>'
    
    data += '</ul>'
    
    return html_minify(render_template(skin_check(), 
        imp = [name, wiki_set(1), custom(), other2([' (' + load_lang('move') + ' ' + load_lang('history') + ')', 0])],
        data = data,
        menu = [['history/' + url_pas(name), load_lang('history')]]
    ))        
            
@app.route('/move/<path:name>', methods=['POST', 'GET'])
def move(name = None):
    if acl_check(name) == 1:
        return re_error('/ban')

    if request.method == 'POST':
        if captcha_post(request.form.get('g-recaptcha-response', None), conn) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        curs.execute("select title from history where title = ?", [request.form.get('title', None)])
        if curs.fetchall():
            return re_error('/error/19')
        
        curs.execute("select data from data where title = ?", [name])
        data = curs.fetchall()
        if data:            
            curs.execute("update data set title = ? where title = ?", [request.form.get('title', None), name])
            curs.execute("update back set link = ? where link = ?", [request.form.get('title', None), name])
            
            data_in = data[0][0]
        else:
            data_in = ''
            
        history_plus(name, data_in, get_time(), ip_check(), request.form.get('send', None) + ' (<a href="/w/' + url_pas(name) + '">' + name + '</a> - <a href="/w/' + url_pas(request.form.get('title', None)) + '">' + request.form.get('title', None) + '</a> ' + load_lang('move') + ')', '0')
        
        curs.execute("select title, link from back where title = ? and not type = 'cat' and not type = 'no'", [name])
        for data in curs.fetchall():
            curs.execute("insert into back (title, link, type) values (?, ?, 'no')", [data[0], data[1]])
            
        curs.execute("update history set title = ? where title = ?", [request.form.get('title', None), name])
        conn.commit()

        return redirect('/w/' + url_pas(request.form.get('title', None)))
    else:            
        return html_minify(render_template(skin_check(), 
            imp = [name, wiki_set(1), custom(), other2([' (' + load_lang('move') + ')', 0])],
            data = '''
                    <form method="post">
                        ''' + ip_warring() + '''
                        <input placeholder="''' + load_lang('document') + ' ' + load_lang('name') + '" value="' + name + '''" name="title" type="text">
                        <hr>
                        <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                        <hr>
                        ''' + captcha_get() + '''
                        <button type="submit">''' + load_lang('move') + '''</button>
                    </form>
                    ''',
            menu = [['w/' + url_pas(name), load_lang('document')]]
        ))
            
@app.route('/other')
def other():
    return html_minify(render_template(skin_check(), 
        imp = [load_lang('other') + ' ' + load_lang('tool'), wiki_set(1), custom(), other2([0, 0])],
        data = '''
                <h2>''' + load_lang('record') + '''</h2>
                <ul>
                    <li><a href="/manager/6">''' + load_lang('edit') + '''</a></li>
                    <li><a href="/manager/7">''' + load_lang('discussion') + '''</a></li>
                </ul>
                <br>
                <h2>''' + load_lang('list') + '''</h2>
                <ul>
                    <li><a href="/admin_list">''' + load_lang('admin') + '''</a></li>
                    <li><a href="/give_log">''' + load_lang('admin_group') + '''</a></li>
                    <li><a href="/not_close_topic">''' + load_lang('open') + ' ' + load_lang('discussion') + '''</a></li>
                </ul>
                <br>
                <h2>''' + load_lang('other') + '''</h2>
                <ul>
                    <li><a href="/title_index">''' + load_lang('all') + ' ' + load_lang('document') + '''</a></li>
                    <li><a href="/acl_list">ACL ''' + load_lang('document') + '''</a></li>
                    <li><a href="/please">''' + load_lang('need') + ' ' + load_lang('document') + '''</a></li>
                    <li><a href="/upload">''' + load_lang('upload') + '''</a></li>
                    <li><a href="/manager/10">''' + load_lang('document') + ' ' + load_lang('search') + '''</a></li>
                </ul>
                <br>
                <h2>''' + load_lang('admin') + '''</h2>
                <ul>
                    <li><a href="/manager/1">''' + load_lang('admin') + ' ' + load_lang('tool') + '''</a></li>
                </ul>
                <br>
                <h2>''' + load_lang('version') + '''</h2>
                <ul>
                    <li>''' + load_lang('version') + ' : <a id="out_link" href="https://github.com/2DU/openNAMU/blob/master/version.md">' + r_ver + '''</a></li>
                </ul>''',
    menu = 0
    ))
    
@app.route('/manager', methods=['POST', 'GET'])
@app.route('/manager/<int:num>', methods=['POST', 'GET'])
def manager(num = 1):
    title_list = [[load_lang('document') + ' ' + load_lang('name'), 'acl'], [0, 'check'], [0, 'ban'], [0, 'admin'], [0, 'record'], [0, 'topic_record'], [load_lang('name'), 'admin_plus'], [load_lang('name'), 'edit_filter'], [load_lang('document') + ' ' + load_lang('name'), 'search'], [0, 'block_user'], [0, 'block_admin'], [load_lang('document') + ' ' + load_lang('name'), 'watch_list']]
    
    if num == 1:
        return html_minify(render_template(skin_check(), 
            imp = [load_lang('admin') + ' ' + load_lang('tool'), wiki_set(1), custom(), other2([0, 0])],
            data = '''
                    <h2>''' + load_lang('admin') + '''</h2>
                    <ul>
                        <li><a href="/manager/2">''' + load_lang('document') + ''' ACL</a></li>
                        <li><a href="/manager/3">''' + load_lang('user') + ' ' + load_lang('check') + '''</a></li>
                        <li><a href="/manager/4">''' + load_lang('user') + ' ' + load_lang('ban') + '''</a></li>
                        <li><a href="/manager/5">''' + load_lang('authority') + '''</a></li>
                        <li><a href="/big_delete">''' + load_lang('bulk_delete') + '''</a></li>
                        <li><a href="/edit_filter">''' + load_lang('edit_filter') + '''</a></li>
                    </ul>
                    <br>
                    <h2>''' + load_lang('owner') + '''</h2>
                    <ul>
                        <li><a href="/indexing">Indexing (''' + load_lang('create') + ' or ' + load_lang('delete') + ''')</a></li>
                        <li><a href="/manager/8">''' + load_lang('admin_group') + ' ' + load_lang('create') + '''</a></li>
                        <li><a href="/edit_set">''' + load_lang('setting') + ' ' + load_lang('edit') + '''</a></li>
                        <li><a href="/re_start">Server Restart</a></li>
                        <li><a href="/update">''' + load_lang('update') + '''</a></li>
                        <li><a href="/inter_wiki">''' + load_lang('interwiki') + '''</a></li>
                    </ul>
                    ''',
            menu = [['other', load_lang('other')]]
        ))
    elif num in range(2, 14):
        if request.method == 'POST':
            return redirect('/' + title_list[(num - 2)][1] + '/' + url_pas(request.form.get('name', None)))
        else:
            if title_list[(num - 2)][0] == 0:
                placeholder = load_lang('user') + ' ' + load_lang('name')
            else:
                placeholder = title_list[(num - 2)][0]

            return html_minify(render_template(skin_check(), 
                imp = ['Redirect', wiki_set(1), custom(), other2([0, 0])],
                data = '<form method="post"><input placeholder="' + placeholder + '" name="name" type="text"><hr><button type="submit">' + load_lang('move') + '</button></form>',
                menu = [['manager', load_lang('admin')]]
            ))
    else:
        return redirect('/')
        
@app.route('/title_index')
def title_index():
    page = int(request.args.get('page', 1))
    num = int(request.args.get('num', 100))
    if page * num > 0:
        sql_num = page * num - num
    else:
        sql_num = 0

    all_list = sql_num + 1

    if num > 1000:
        return re_error('/error/3')

    data = '<a href="/title_index?num=250">(250)</a> <a href="/title_index?num=500">(500)</a> <a href="/title_index?num=1000">(1000)</a>'

    curs.execute("select title from data order by title asc limit ?, ?", [str(sql_num), str(num)])
    title_list = curs.fetchall()
    if title_list:
        data += '<hr><ul>'

    for list_data in title_list:
        data += '<li>' + str(all_list) + '. <a href="/w/' + url_pas(list_data[0]) + '">' + list_data[0] + '</a></li>'        
        all_list += 1

    if page == 1:
        count_end = []

        curs.execute("select count(title) from data")
        count = curs.fetchall()
        if count:
            count_end += [count[0][0]]
        else:
            count_end += [0]

        sql_list = ['틀:', '분류:', load_lang('user') + ':', '' + load_lang('file') + ':']
        for sql in sql_list:
            curs.execute("select count(title) from data where title like ?", [sql + '%'])
            count = curs.fetchall()
            if count:
                count_end += [count[0][0]]
            else:
                count_end += [0]

        count_end += [count_end[0] - count_end[1]  - count_end[2]  - count_end[3]  - count_end[4]]
        
        data += '</ul><hr><ul><li>All : ' + str(count_end[0]) + '</li></ul><hr><ul>'
        data += '<li>Template : ' + str(count_end[1]) + '</li>'
        data += '<li>Category : ' + str(count_end[2]) + '</li>'
        data += '<li>User : ' + str(count_end[3]) + '</li>'
        data += '<li>File : ' + str(count_end[4]) + '</li>'
        data += '<li>Other : ' + str(count_end[5]) + '</li>'

    data += '</ul>' + next_fix('/title_index?num=' + str(num) + '&page=', page, title_list, num)
    sub = ' (' + str(num) + '개)'
    
    return html_minify(render_template(skin_check(), 
        imp = [load_lang('all') + ' ' + load_lang('document'), wiki_set(1), custom(), other2([sub, 0])],
        data = data,
        menu = [['other', load_lang('other')]]
    ))
        
@app.route('/topic/<path:name>/sub/<sub>/b/<int:num>')
def topic_block(name = None, sub = None, num = None):
    if admin_check(3, 'blind (' + name + ' - ' + sub + '#' + str(num) + ')') != 1:
        return re_error('/error/3')

    curs.execute("select block from topic where title = ? and sub = ? and id = ?", [name, sub, str(num)])
    block = curs.fetchall()
    if block:
        if block[0][0] == 'O':
            curs.execute("update topic set block = '' where title = ? and sub = ? and id = ?", [name, sub, str(num)])
        else:
            curs.execute("update topic set block = 'O' where title = ? and sub = ? and id = ?", [name, sub, str(num)])
        
        rd_plus(name, sub, get_time())
        
        conn.commit()
        
    return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#' + str(num))
        
@app.route('/topic/<path:name>/sub/<sub>/notice/<int:num>')
def topic_top(name = None, sub = None, num = None):
    if admin_check(3, 'notice (' + name + ' - ' + sub + '#' + str(num) + ')') != 1:
        return re_error('/error/3')

    curs.execute("select title from topic where title = ? and sub = ? and id = ?", [name, sub, str(num)])
    if curs.fetchall():
        curs.execute("select top from topic where id = ? and title = ? and sub = ?", [str(num), name, sub])
        top_data = curs.fetchall()
        if top_data:
            if top_data[0][0] == 'O':
                curs.execute("update topic set top = '' where title = ? and sub = ? and id = ?", [name, sub, str(num)])
            else:
                curs.execute("update topic set top = 'O' where title = ? and sub = ? and id = ?", [name, sub, str(num)])
        
        rd_plus(name, sub, get_time())

        conn.commit()

    return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#' + str(num))        
        
@app.route('/topic/<path:name>/sub/<sub>/tool/<regex("close|stop|agree"):tool>')
def topic_stop(name = None, sub = None, tool = None):
    if tool == 'close':
        set_list = ['O', '', load_lang('discussion') + ' ' + load_lang('close') + '', load_lang('discussion') + ' ' + load_lang('open') + '']
    elif tool == 'stop':
        set_list = ['', 'O', load_lang('discussion') + ' ' + load_lang('stop') + '', load_lang('discussion') + ' ' + load_lang('restart') + '']
    elif tool == 'agree':
        pass
    else:
        return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub))

    if admin_check(3, 'topic ' + tool + ' (' + name + ' - ' + sub + ')') != 1:
        return re_error('/error/3')

    ip = ip_check()
    time = get_time()
    
    curs.execute("select id from topic where title = ? and sub = ? order by id + 0 desc limit 1", [name, sub])
    topic_check = curs.fetchall()
    if topic_check:
        if tool == 'agree':
            curs.execute("select title from agreedis where title = ? and sub = ?", [name, sub])
            if curs.fetchall():
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, '" + load_lang('agreement') + " Fail', ?, ?, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, time, ip])
                curs.execute("delete from agreedis where title = ? and sub = ?", [name, sub])
            else:
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, '" + load_lang('agreement') + " OK', ?, ?, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, time, ip])
                curs.execute("insert into agreedis (title, sub) values (?, ?)", [name, sub])
        else:
            curs.execute("select title from stop where title = ? and sub = ? and close = ?", [name, sub, set_list[0]])
            if curs.fetchall():
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, ?, ?, ?, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, set_list[3], time, ip])
                curs.execute("delete from stop where title = ? and sub = ? and close = ?", [name, sub, set_list[0]])
            else:
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, ?, ?, ?, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, set_list[2], time, ip])
                curs.execute("insert into stop (title, sub, close) values (?, ?, ?)", [name, sub, set_list[0]])
                curs.execute("delete from stop where title = ? and sub = ? and close = ?", [name, sub, set_list[1]])
        
        rd_plus(name, sub, time)
        
        conn.commit()
        
    return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub))    

@app.route('/topic/<path:name>/sub/<sub>/admin/<int:num>')
def topic_admin(name = None, sub = None, num = None):
    curs.execute("select block, ip, date from topic where title = ? and sub = ? and id = ?", [name, sub, str(num)])
    data = curs.fetchall()
    if not data:
        return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub))

    ban = ''

    if admin_check(3, None) == 1:
        ban += '</ul><br><h2>관리자 ' + load_lang('tool') + '</h2><ul>'
        is_ban = '<li><a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/b/' + str(num) + '">'

        if data[0][0] == 'O':
            is_ban += load_lang('hide') + ' ' + load_lang('release')
        else:
            is_ban += load_lang('hide')
        
        is_ban += '</a></li>'
        is_ban += '<li><a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/notice/' + str(num) + '">'

        curs.execute("select id from topic where title = ? and sub = ? and id = ? and top = 'O'", [name, sub, str(num)])
        if curs.fetchall():
            is_ban += '공지 ' + load_lang('release')
        else:
            is_ban += '공지'
        
        is_ban += '</a></li></ul>'
        ban += '<li><a href="/ban/' + url_pas(data[0][1]) + '">'

        curs.execute("select end from ban where block = ?", [data[0][1]])
        if curs.fetchall():
            ban += load_lang('ban') + ' ' + load_lang('release')
        else:
            ban += load_lang('ban')
        
        ban += '</a></li>' + is_ban

    ban += '</ul><br><h2>' + load_lang('other') + ' ' + load_lang('tool') + '</h2><ul>'
    ban += '<li><a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/raw/' + str(num) + '">Raw</a></li>'
    ban = '<li>' + load_lang('time') + ' : ' + data[0][2] + '</li>' + ban
    
    if re.search('(\.|:)', data[0][1]):
        ban = '<li>작성인 : ' + data[0][1] + ' <li><a href="/record/' + url_pas(data[0][1]) + '">(' + load_lang('record') + ')</a></li>' + ban
    else:
        ban = '<li>작성인 : <a href="/w/' + load_lang('user') + ':' + data[0][1] + '">' + data[0][1] + '</a> <a href="/record/' + url_pas(data[0][1]) + '">(' + load_lang('record') + ')</a></li>' + ban

    ban = '<h2>정보</h2><ul>' + ban

    return html_minify(render_template(skin_check(), 
        imp = [load_lang('discussion') + ' ' + load_lang('tool'), wiki_set(1), custom(), other2([' (' + str(num) + '번)', 0])],
        data = ban,
        menu = [['topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#' + str(num), load_lang('discussion')]]
    ))

@app.route('/topic/<path:name>/sub/<sub>', methods=['POST', 'GET'])
def topic(name = None, sub = None):
    ban = topic_check(name, sub)
    admin = admin_check(3, None)
    
    if request.method == 'POST':
        if captcha_post(request.form.get('g-recaptcha-response', None), conn) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        ip = ip_check()
        today = get_time()
        
        if ban == 1:
            return re_error('/ban')
        
        curs.execute("select id from topic where title = ? and sub = ? order by id + 0 desc limit 1", [name, sub])
        old_num = curs.fetchall()
        if old_num:
            num = int(old_num[0][0]) + 1
        else:
            num = 1

        match = re.search('^' + load_lang('user') + ':([^/]+)', name)
        if match:
            curs.execute('insert into alarm (name, data, date) values (?, ?, ?)', [match.groups()[0], ip + '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '">' + load_lang('user') + ' - ' + load_lang('discussion') + '</a> (My)', today])
        
        data = re.sub("\[\[(분류:(?:(?:(?!\]\]).)*))\]\]", "[br]", request.form.get('content', None))
        for rd_data in re.findall("(?:#([0-9]+))", data):
            curs.execute("select ip from topic where title = ? and sub = ? and id = ?", [name, sub, rd_data])
            ip_data = curs.fetchall()
            if ip_data and not re.search('(\.|:)', ip_data[0][0]):
                curs.execute('insert into alarm (name, data, date) values (?, ?, ?)', [ip_data[0][0], ip + ' - <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#' + str(num) + '">' + load_lang('discussion') + '</a>', today])
            
            data = re.sub("(?P<in>#(?:[0-9]+))", '[[\g<in>]]', data)

        data = savemark(data)

        rd_plus(name, sub, today)

        curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, ?, ?, ?, '', '')", [str(num), name, sub, data, today, ip])
        conn.commit()
        
        return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub))
    else:
        curs.execute("select title from stop where title = ? and sub = ? and close = 'O'", [name, sub])
        close_data = curs.fetchall()
        
        curs.execute("select title from stop where title = ? and sub = ? and close = ''", [name, sub])
        stop_data = curs.fetchall()
        
        curs.execute("select id from topic where title = ? and sub = ? limit 1", [name, sub])
        topic_exist = curs.fetchall()
        
        display = ''
        all_data = ''
        data = ''
        number = 1
        
        if admin == 1 and topic_exist:
            if close_data:
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/close">(' + load_lang('open') + ')</a> '
            else:
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/close">(' + load_lang('close') + ')</a> '
            
            if stop_data:
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/stop">(' + load_lang('restart') + ')</a> '
            else:
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/stop">(' + load_lang('stop') + ')</a> '

            curs.execute("select title from agreedis where title = ? and sub = ?", [name, sub])
            if curs.fetchall():
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/agree">(' + load_lang('release') + ')</a>'
            else:
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/agree">(' + load_lang('agreement') + ')</a>'
            
            all_data += '<hr>'
        
        if (close_data or stop_data) and admin != 1:
            display = 'display: none;'
        
        curs.execute("select data, id, date, ip, block, top from topic where title = ? and sub = ? order by id + 0 asc", [name, sub])
        topic = curs.fetchall()
        
        curs.execute("select data, id, date, ip from topic where title = ? and sub = ? and top = 'O' order by id + 0 asc", [name, sub])
        for topic_data in curs.fetchall():                   
            who_plus = ''
            
            curs.execute("select who from re_admin where what = ? order by time desc limit 1", ['notice (' + name + ' - ' + sub + '#' + topic_data[1] + ')'])
            topic_data_top = curs.fetchall()
            if topic_data_top:
                who_plus += ' <span style="margin-right: 5px;">@' + topic_data_top[0][0] + ' </span>'
                                
            all_data += '<table id="toron"><tbody><tr><td id="toron_color_red">'
            all_data += '<a href="#' + topic_data[1] + '">#' + topic_data[1] + '</a> ' + ip_pas(topic_data[3]) + who_plus + ' <span style="float: right;">' + topic_data[2] + '</span>'
            all_data += '</td></tr><tr><td>' + namumark('', topic_data[0], 0) + '</td></tr></tbody></table><br>'    

        for topic_data in topic:
            if number == 1:
                start = topic_data[3]

            if topic_data[4] == 'O':
                blind_data = 'style="background: gainsboro;"'
                
                if admin != 1:
                    curs.execute("select who from re_admin where what = ? order by time desc limit 1", ['blind (' + name + ' - ' + sub + '#' + str(number) + ')'])
                    who_blind = curs.fetchall()
                    if who_blind:
                        user_write = '[[' + load_lang('user') + ':' + who_blind[0][0] + ']] ' + load_lang('hide')
                    else:
                        user_write = load_lang('hide')
            else:
                blind_data = ''

            user_write = namumark('', topic_data[0], 0)
            ip = ip_pas(topic_data[3])
            
            curs.execute('select acl from user where id = ?', [topic_data[3]])
            user_acl = curs.fetchall()
            if user_acl and user_acl[0][0] != 'user':
                ip += ' <a href="javascript:void(0);" title="' + load_lang('admin') + '">★</a>'

            if admin == 1 or blind_data == '':
                ip += ' <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/admin/' + str(number) + '">(' + load_lang('tool') + ')</a>'

            curs.execute("select end from ban where block = ?", [topic_data[3]])
            if curs.fetchall():
                ip += ' <a href="javascript:void(0);" title="차단자">†</a>'
                    
            if topic_data[5] == '1':
                color = '_blue'
            elif topic_data[3] == start:
                color = '_green'
            else:
                color = ''
                
            if user_write == '':
                user_write = '<br>'
                         
            all_data += '<table id="toron"><tbody><tr><td id="toron_color' + color + '">'
            all_data += '<a href="javascript:void(0);" id="' + str(number) + '">#' + str(number) + '</a> ' + ip + '</span>'
            all_data += '</td></tr><tr ' + blind_data + '><td>' + user_write + '</td></tr></tbody></table><br>'
           
            number += 1

        if ban != 1 or admin == 1:
            data += '''
                    <a id="reload" href="javascript:void(0);" onclick="location.href.endsWith(\'#reload\')? location.reload(true):location.href=\'#reload\'">(Reload)</a>
                    <form style="''' + display + '''" method="post">
                    <br>
                    <textarea style="height: 100px;" name="content"></textarea>
                    <hr>
                    ''' + captcha_get()
            
            if display == '':
                data += ip_warring()

            data += '<button type="submit">Send</button></form>'

        return html_minify(render_template(skin_check(), 
            imp = [name, wiki_set(1), custom(), other2([' (' + load_lang('discussion') + ')', 0])],
            data = '<h2 id="topic_top_title">' + sub + '</h2>' + all_data + data,
            menu = [['topic/' + url_pas(name), load_lang('list')]]
        ))
        
@app.route('/topic/<path:name>', methods=['POST', 'GET'])
@app.route('/topic/<path:name>/<regex("close|agree"):tool>', methods=['GET'])
def close_topic_list(name = None, tool = None):
    div = ''
    list_d = 0
    
    if request.method == 'POST':
        t_num = ''
        
        while 1:
            curs.execute("select title from topic where title = ? and sub = ? limit 1", [name, request.form.get('topic', None) + t_num])
            if curs.fetchall():
                if t_num == '':
                    t_num = ' 2'
                else:
                    t_num = ' ' + str(int(t_num.replace(' ', '')) + 1)
            else:
                break

        return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(request.form.get('topic', None) + t_num))
    else:
        plus = ''
        menu = [['topic/' + url_pas(name), load_lang('list')]]
        
        if tool == 'close':
            curs.execute("select sub from stop where title = ? and close = 'O' order by sub asc", [name])
            
            sub = '' + load_lang('close') + ''
        elif tool == 'agree':
            curs.execute("select sub from agreedis where title = ? order by sub asc", [name])
            
            sub = '' + load_lang('agreement') + ''
        else:
            curs.execute("select sub from rd where title = ? order by date desc", [name])
            
            sub = load_lang('discussion') + ' ' + load_lang('list')
            
            menu = [['w/' + url_pas(name), load_lang('document')]]
            
            plus =  '<a href="/topic/' + url_pas(name) + '/close">(' + load_lang('close') + ')</a> <a href="/topic/' + url_pas(name) + '/agree">(' + load_lang('agreement') + ')</a><hr><input placeholder="' + load_lang('discussion') + ' ' + load_lang('name') + '" name="topic" type="text"><hr><button type="submit">' + load_lang('open') + '</button>'

        for data in curs.fetchall():
            curs.execute("select data, date, ip, block from topic where title = ? and sub = ? and id = '1'", [name, data[0]])
            if curs.fetchall():                
                it_p = 0
                
                if sub == load_lang('discussion') + ' ' + load_lang('list'):
                    curs.execute("select title from stop where title = ? and sub = ? and close = 'O' order by sub asc", [name, data[0]])
                    if curs.fetchall():
                        it_p = 1
                
                if it_p != 1:
                    div += '<h2><a href="/topic/' + url_pas(name) + '/sub/' + url_pas(data[0]) + '">' + data[0] + '</a></h2>'

        if div == '':
            plus = re.sub('^<br>', '', plus)
        
        return html_minify(render_template(skin_check(), 
            imp = [name, wiki_set(1), custom(), other2([' (' + sub + ')', 0])],
            data =  '<form method="post">' + div + plus + '</form>',
            menu = menu
        ))
        
@app.route('/login', methods=['POST', 'GET'])
def login():
    if 'Now' in session and session['Now'] == 1:
        return re_error('/error/11')

    ip = ip_check()
    agent = request.headers.get('User-Agent')
    
    curs.execute("select block from ban where block = ? and login = 'O'", [ip])
    if not curs.fetchall():
        match = re.search("^([0-9]{1,3}\.[0-9]{1,3})", ip)
        if match:
            match = match.groups()[0]
        else:
            match = 'Not'

        curs.execute("select block from ban where block = ? and login = 'O'", [match])
        if not curs.fetchall():
            ban = ban_check()
        else:
            ban = 0
    else:
        ban = 0

    if ban == 1:
        return re_error('/ban')
        
    if request.method == 'POST':        
        if captcha_post(request.form.get('g-recaptcha-response', None), conn) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        curs.execute("select pw from user where id = ?", [request.form.get('id', None)])
        user = curs.fetchall()
        if not user:
            return re_error('/error/5')

        if not bcrypt.checkpw(bytes(request.form.get('pw', None), 'utf-8'), bytes(user[0][0], 'utf-8')):
            return re_error('/error/10')

        session['Now'] = 1
        session['DREAMER'] = request.form.get('id', None)
        
        curs.execute("select css from custom where user = ?", [request.form.get('id', None)])
        css_data = curs.fetchall()
        if css_data:
            session['Daydream'] = css_data[0][0]
        else:
            session['Daydream'] = ''
        
        curs.execute("insert into ua_d (name, ip, ua, today, sub) values (?, ?, ?, ?, '')", [request.form.get('id', None), ip, agent, get_time()])
        conn.commit()
        
        return redirect('/user')  
    else:        
        return html_minify(render_template(skin_check(),    
            imp = [load_lang('login'), wiki_set(1), custom(), other2([0, 0])],
            data = '''
                    <form method="post">
                        <input placeholder="ID" name="id" type="text">
                        <hr>
                        <input placeholder="PassWord" name="pw" type="password">
                        <hr>
                        ''' + captcha_get() + '''
                        <button type="submit">''' + load_lang('login') + '''</button>
                        <hr>
                        <span>''' + load_lang('http_warring') + '''</span>
                    </form>''',
            menu = [['user', load_lang('user')]]
        ))
                
@app.route('/change', methods=['POST', 'GET'])
def change_password():
    if ban_check() == 1:
        return re_error('/ban')

    if custom()[2] == 0:
        return redirect('/login')
    
    if request.method == 'POST':    
        if request.form.get('pw', None):
            if request.form.get('pw2', None) != request.form.get('pw3', None):
                return re_error('/error/20')

            curs.execute("select pw from user where id = ?", [session['DREAMER']])
            user = curs.fetchall()
            if not user:
                return re_error('/error/10')

            if not bcrypt.checkpw(bytes(request.form.get('pw', None), 'utf-8'), bytes(user[0][0], 'utf-8')):
                return re_error('/error/5')

            hashed = bcrypt.hashpw(bytes(request.form.get('pw2', None), 'utf-8'), bcrypt.gensalt())
            
            curs.execute("update user set pw = ? where id = ?", [hashed.decode(), session['DREAMER']])
        
        curs.execute("update user set email = ? where id = ?", [request.form.get('email', ''), ip_check()])
        curs.execute("update user set skin = ? where id = ?", [request.form.get('skin', ''), ip_check()])
        conn.commit()
        
        return redirect('/change')
    else:        
        ip = ip_check()

        curs.execute('select email from user where id = ?', [ip])
        data = curs.fetchall()
        if data:
            email = data[0][0]
        else:
            email = ''

        div2 = ''

        curs.execute('select skin from user where id = ?', [ip])
        data = curs.fetchall()

        for skin_data in os.listdir(os.path.abspath('views')):
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

        return html_minify(render_template(skin_check(),    
            imp = [load_lang('my_info') + ' ' + load_lang('edit'), wiki_set(1), custom(), other2([0, 0])],
            data = '''
                    <form method="post">
                        <span>ID : ''' + ip + '''</span>
                        <hr>
                        <input placeholder="Now" name="pw" type="password">
                        <br>
                        <br>
                        <input placeholder="Change" name="pw2" type="password">
                        <br>
                        <br>
                        <input placeholder="Re" name="pw3" type="password">
                        <hr>
                        <input placeholder="Email" name="email" type="text" value="''' + email + '''">
                        <hr>
                        <span>Skin</span>
                        <br>
                        <br>
                        <select name="skin">''' + div2 + '''</select>
                        <hr>
                        <button type="submit">''' + load_lang('edit') + '''</button>
                        <hr>
                        <span>''' + load_lang('http_warring') + '''</span>
                    </form>
                    ''',
            menu = [['user', load_lang('user')]]
        ))

@app.route('/check/<name>')
def user_check(name = None):
    if admin_check(4, 'check (' + name + ')') != 1:
        return re_error('/error/3')

    curs.execute("select acl from user where id = ? or id = ?", [name, request.args.get('plus', 'None-Data')])
    user = curs.fetchall()
    if user and user[0][0] != 'user':
        if admin_check(None, None) != 1:
            return re_error('/error/4')
    
    if request.args.get('plus', None):
        if re.search('(?:\.|:)', name):
            if re.search('(?:\.|:)', request.args.get('plus', None)):
                curs.execute("select name, ip, ua, today from ua_d where ip = ? or ip = ? order by today desc", [name, request.args.get('plus', None)])
            else:
                curs.execute("select name, ip, ua, today from ua_d where ip = ? or name = ? order by today desc", [name, request.args.get('plus', None)])
        else:
            if re.search('(?:\.|:)', request.args.get('plus', None)):
                curs.execute("select name, ip, ua, today from ua_d where name = ? or ip = ? order by today desc", [name, request.args.get('plus', None)])
            else:
                curs.execute("select name, ip, ua, today from ua_d where name = ? or name = ? order by today desc", [name, request.args.get('plus', None)])
    elif re.search('(?:\.|:)', name):
        curs.execute("select name, ip, ua, today from ua_d where ip = ? order by today desc", [name])
    else:
        curs.execute("select name, ip, ua, today from ua_d where name = ? order by today desc", [name])
    
    record = curs.fetchall()
    if record:
        if not request.args.get('plus', None):
            div = '<a href="/plus_check/' + url_pas(name) + '">(' + load_lang('compare') + ')</a><hr>'
        else:
            div = '<a href="/check/' + url_pas(name) + '">(Main)</a> <a href="/check/' + url_pas(request.args.get('plus', None)) + '">(Sub)</a><hr>'

        div += '<table style="width: 100%; text-align: center;"><tbody><tr>'
        div += '<td style="width: 33.3%;">' + load_lang('name') + '</td><td style="width: 33.3%;">IP</td><td style="width: 33.3%;">언제</td></tr>'
        
        for data in record:
            if data[2]:
                ua = data[2]
            else:
                ua = '<br>'

            div += '<tr><td>' + ip_pas(data[0]) + '</td><td>' + ip_pas(data[1]) + '</td><td>' + data[3] + '</td></tr>'
            div += '<tr><td colspan="3">' + ua + '</td></tr>'
        
        div += '</tbody></table>'
    else:
        return re_error('/error/5')
            
    return html_minify(render_template(skin_check(),    
        imp = ['' + load_lang('check') + '', wiki_set(1), custom(), other2([0, 0])],
        data = div,
        menu = [['manager', load_lang('admin')]]
    ))

@app.route('/plus_check/<name>', methods=['POST', 'GET'])
def plus_check(name):
    if request.method == 'POST':
        return redirect('/check/' + url_pas(name) + '?plus=' + url_pas(request.form.get('name2', None)))
    else:
        return html_minify(render_template(skin_check(),
            imp = ['' + load_lang('plus'), wiki_set(1), custom(), other2([0, 0])],
            data = '''
                    <form method="post">
                        <input placeholder="' + load_lang('compare') + '" name="name2" type="text">
                        <hr>
                        <button type="submit">''' + load_lang('move') + '''</button>
                    </form>
                    ''',
            menu = [['manager', load_lang('admin')]]
        ))
                
@app.route('/register', methods=['POST', 'GET'])
def register():
    if ban_check() == 1:
        return re_error('/ban')

    if not admin_check(None, None) == 1:
        curs.execute('select data from other where name = "reg"')
        set_d = curs.fetchall()
        if set_d and set_d[0][0] == 'on':
            return re_error('/ban')
    
    if request.method == 'POST': 
        if captcha_post(request.form.get('g-recaptcha-response', None), conn) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        if request.form.get('pw', None) != request.form.get('pw2', None):
            return re_error('/error/20')

        if re.search('(?:[^A-Za-zㄱ-힣0-9 ])', request.form.get('id', None)):
            return re_error('/error/8')

        if len(request.form.get('id', None)) > 32:
            return re_error('/error/7')

        curs.execute("select id from user where id = ?", [request.form.get('id', None)])
        if curs.fetchall():
            return re_error('/error/6')

        hashed = bcrypt.hashpw(bytes(request.form.get('pw', None), 'utf-8'), bcrypt.gensalt())
        
        curs.execute("select id from user limit 1")
        if not curs.fetchall():
            curs.execute("insert into user (id, pw, acl, date, email) values (?, ?, 'owner', ?, ?)", [request.form.get('id', None), hashed.decode(), get_time(), request.form.get('email', '')])
        else:
            curs.execute("insert into user (id, pw, acl, date, email) values (?, ?, 'user', ?, ?)", [request.form.get('id', None), hashed.decode(), get_time(), request.form.get('email', '')])
        
        conn.commit()
        
        return redirect('/login')
    else:        
        contract = ''
        
        curs.execute('select data from other where name = "contract"')
        data = curs.fetchall()
        if data and data[0][0] != '':
            contract = data[0][0] + '<hr>'

        return html_minify(render_template(skin_check(),    
            imp = ['' + load_lang('register') + '', wiki_set(1), custom(), other2([0, 0])],
            data = '''
                    <form method="post">
                        ''' + contract + '''
                        <input placeholder="ID" name="id" type="text">
                        <hr>
                        <input placeholder="PassWord" name="pw" type="password">
                        <hr>
                        <input placeholder="Re" name="pw2" type="password">
                        <hr>
                        <input placeholder="Email (Option)" name="email" type="text">
                        <hr>
                        ''' + captcha_get() + '''
                        <button type="submit">' + load_lang('register') + '</button>
                        <hr>
                        <span>''' + load_lang('http_warring') + '''</span>
                    </form>
                    ''',
            menu = [['user', load_lang('user')]]
        ))
            
@app.route('/logout')
def logout():
    session['Now'] = 0
    session.pop('DREAMER', None)

    return redirect('/user')
    
@app.route('/ban/<name>', methods=['POST', 'GET'])
def user_ban(name = None):
    curs.execute("select acl from user where id = ?", [name])
    user = curs.fetchall()
    if user and user[0][0] != 'user':
        if admin_check(None, None) != 1:
            return re_error('/error/4')

    if request.method == 'POST':
        if admin_check(1, 'ban (' + name + ')') != 1:
            return re_error('/error/3')

        if request.form.get('year', 'no_end') == 'no_end':
            end = ''
        else:
            end = request.form.get('year', '') + '-' + request.form.get('month', '') + '-' + request.form.get('day', '')

        if end == '--':
            end = ''

        ban_insert(name, end, request.form.get('why', ''), request.form.get('login', ''), ip_check())

        return redirect('/ban/' + url_pas(name))     
    else:
        if admin_check(1, None) != 1:
            return re_error('/error/3')

        curs.execute("select end, why from ban where block = ?", [name])
        end = curs.fetchall()
        if end:
            now = load_lang('ban') + ' ' + load_lang('release')

            if end[0][0] == '':
                data = '<ul><li>무기한 ' + load_lang('ban') + '</li>'
            else:
                data = '<ul><li>' + load_lang('ban') + ' : ' + end[0][0] + '</li>'

            if end[0][1] != '':
                data += '<li>' + load_lang('why') + ' : ' + end[0][1] + '</li></ul><hr>'
            else:
                data += '</ul><hr>'
        else:
            if re.search("^([0-9]{1,3}\.[0-9]{1,3})$", name):
                now = '대역 ' + load_lang('ban')
            else:
                now = load_lang('ban')

            now_time = get_time()

            m = re.search('^([0-9]{4})-([0-9]{2})-([0-9]{2})', now_time)
            g = m.groups()

            year = '<option value="no_end">영구</option>'
            for i in range(int(g[0]), int(g[0]) + 11):
                if i == int(g[0]):
                    year += '<option value="' + str(i) + '" selected>' + str(i) + '</option>'
                else:
                    year += '<option value="' + str(i) + '">' + str(i) + '</option>'

            month = ''
            for i in range(1, 13):
                if int(i / 10) == 0:
                    num = '0' + str(i)
                else:
                    num = str(i)

                if i == int(g[1]):
                    month += '<option value="' + num + '" selected>' + num + '</option>'
                else:
                    month += '<option value="' + num + '">' + num + '</option>'
                
            day = ''
            for i in range(1, 32):
                if int(i / 10) == 0:
                    num = '0' + str(i)
                else:
                    num = str(i)

                if i == int(g[2]):
                    day += '<option value="' + num + '" selected>' + num + '</option>'
                else:
                    day += '<option value="' + num + '">' + num + '</option>'
            
            if re.search('(\.|:)', name):
                plus = '<input type="checkbox" name="login"> ' + load_lang('login') + ' ' + load_lang('able') + '<hr>'
            else:
                plus = ''

            data = '<select name="year">' + year + '</select> ' + load_lang('year') + ' '
            data += '<select name="month">' + month + '</select> ' + load_lang('month') + ' '
            data += '<select name="day">' + day + '</select> ' + load_lang('day') + ' <hr>'

            data += '<input placeholder="' + load_lang('why') + '" name="why" type="text"><hr>' + plus

        return html_minify(render_template(skin_check(), 
            imp = [name, wiki_set(1), custom(), other2([' (' + now + ')', 0])],
            data = '<form method="post">' + data + '<button type="submit">' + now + '</button></form>',
            menu = [['manager', load_lang('admin')]]
        ))            
                
@app.route('/acl/<path:name>', methods=['POST', 'GET'])
def acl(name = None):
    check_ok = ''
    
    if request.method == 'POST':
        check_data = 'acl (' + name + ')'
    else:
        check_data = None
    
    user_data = re.search('^' + load_lang('user') + ':(.+)$', name)
    if user_data:
        if check_data and custom()[2] == 0:
            return redirect('/login')
        
        if user_data.groups()[0] != ip_check():
            if admin_check(5, check_data) != 1:
                if check_data:
                    return re_error('/error/3')
                else:
                    check_ok = 'disabled'
    else:
        if admin_check(5, check_data) != 1:
            if check_data:
                return re_error('/error/3')
            else:
                check_ok = 'disabled'

    if request.method == 'POST':
        curs.execute("select title from acl where title = ?", [name])
        if curs.fetchall():
            curs.execute("update acl set dec = ? where title = ?", [request.form.get('dec', ''), name])
            curs.execute("update acl set dis = ? where title = ?", [request.form.get('dis', ''), name])
            curs.execute("update acl set why = ? where title = ?", [request.form.get('why', ''), name])
        else:
            curs.execute("insert into acl (title, dec, dis, why) values (?, ?, ?, ?)", [name, request.form.get('dec', ''), request.form.get('dis', ''), request.form.get('why', '')])
        
        curs.execute("select title from acl where title = ? and dec = '' and dis = ''", [name])
        if curs.fetchall():
            curs.execute("delete from acl where title = ?", [name])

        conn.commit()
            
        return redirect('/acl/' + url_pas(name))            
    else:
        data = '<h2>' + load_lang('document') + ' ACL</h2><select name="dec" ' + check_ok + '>'
    
        if re.search('^' + load_lang('user') + ':', name):
            acl_list = [['', load_lang('normal')], ['user', load_lang('subscriber')], ['all', '모두']]
        else:
            acl_list = [['', load_lang('normal')], ['user', load_lang('subscriber')], ['admin', '관리자']]
        
        curs.execute("select dec from acl where title = ?", [name])
        acl_data = curs.fetchall()
        for data_list in acl_list:
            if acl_data and acl_data[0][0] == data_list[0]:
                check = 'selected="selected"'
            else:
                check = ''
            
            data += '<option value="' + data_list[0] + '" ' + check + '>' + data_list[1] + '</option>'
            
        data += '</select>'
        
        if not re.search('^' + load_lang('user') + ':', name):
            data += '<br><br><h2>' + load_lang('discussion') + ' ACL</h2><select name="dis" ' + check_ok + '>'
        
            curs.execute("select dis, why from acl where title = ?", [name])
            acl_data = curs.fetchall()
            for data_list in acl_list:
                if acl_data and acl_data[0][0] == data_list[0]:
                    check = 'selected="selected"'
                else:
                    check = ''
                    
                data += '<option value="' + data_list[0] + '" ' + check + '>' + data_list[1] + '</option>'
                
            data += '</select>'
                
            if acl_data:
                data += '<hr><input value="' + html.escape(acl_data[0][1]) + '" placeholder="' + load_lang('why') + '" name="why" type="text" ' + check_ok + '>'
            
        return html_minify(render_template(skin_check(), 
            imp = [name, wiki_set(1), custom(), other2([' (ACL)', 0])],
            data = '<form method="post">' + data + '<hr><button type="submit">ACL ' + load_lang('edit') + '</button></form>',
            menu = [['w/' + url_pas(name), load_lang('document')], ['manager', load_lang('admin')]]
        ))
            
@app.route('/admin/<name>', methods=['POST', 'GET'])
def user_admin(name = None):
    owner = admin_check(None, None)
    
    curs.execute("select acl from user where id = ?", [name])
    user = curs.fetchall()
    if not user:
        return re_error('/error/5')
    else:
        if owner != 1:
            curs.execute('select name from alist where name = ? and acl = "owner"', [user[0][0]])
            if curs.fetchall():
                return re_error('/error/3')

            if ip_check() == name:
                return re_error('/error/3')

    if request.method == 'POST':
        if admin_check(7, 'admin (' + name + ')') != 1:
            return re_error('/error/3')

        if owner != 1:
            curs.execute('select name from alist where name = ? and acl = "owner"', [request.form.get('select', None)])
            if curs.fetchall():
                return re_error('/error/3')

        if request.form.get('select', None) == 'X':
            curs.execute("update user set acl = 'user' where id = ?", [name])
        else:
            curs.execute("update user set acl = ? where id = ?", [request.form.get('select', None), name])
        
        conn.commit()
        
        return redirect('/admin/' + url_pas(name))            
    else:
        if admin_check(7, None) != 1:
            return re_error('/error/3')            

        div = '<option value="X">X</option>'
        i = 0
        name_rem = ''
        
        curs.execute('select distinct name from alist order by name asc')
        for data in curs.fetchall():
            if user[0][0] == data[0]:
                div += '<option value="' + data[0] + '" selected="selected">' + data[0] + '</option>'
            else:
                if owner != 1:
                    curs.execute('select name from alist where name = ? and acl = "owner"', [data[0]])
                    if not curs.fetchall():
                        div += '<option value="' + data[0] + '">' + data[0] + '</option>'
                else:
                    div += '<option value="' + data[0] + '">' + data[0] + '</option>'
        
        return html_minify(render_template(skin_check(), 
            imp = [name, wiki_set(1), custom(), other2([' (' + load_lang('authority') + ')', 0])],
            data =  '<form method="post"><select name="select">' + div + '</select><hr><button type="submit">' + load_lang('edit') + '</button></form>',
            menu = [['manager', load_lang('admin')]]
        ))
    
@app.route('/diff/<path:name>')
def diff_data(name = None):
    first = request.args.get('first', '1')
    second = request.args.get('second', '1')

    curs.execute("select data from history where id = ? and title = ?", [first, name])
    first_raw_data = curs.fetchall()
    if first_raw_data:
        curs.execute("select data from history where id = ? and title = ?", [second, name])
        second_raw_data = curs.fetchall()
        if second_raw_data:
            first_data = html.escape(first_raw_data[0][0])            
            second_data = html.escape(second_raw_data[0][0])

            if first == second:
                result = 'Same.'
            else:            
                diff_data = difflib.SequenceMatcher(None, first_data, second_data)
                result = re.sub('\r', '', diff(diff_data))
            
            return html_minify(render_template(skin_check(), 
                imp = [name, wiki_set(1), custom(), other2([' (' + load_lang('compare') + ')', 0])],
                data = '<pre>' + result + '</pre>',
                menu = [['history/' + url_pas(name), load_lang('history')]]
            ))

    return redirect('/history/' + url_pas(name))
        
@app.route('/down/<path:name>')
def down(name = None):
    div = '<ul>'

    curs.execute("select title from data where title like ?", ['%' + name + '/%'])
    for data in curs.fetchall():
        div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a></li>'
        
    div += '</ul>'
    
    return html_minify(render_template(skin_check(), 
        imp = [name, wiki_set(1), custom(), other2([' (하위)', 0])],
        data = div,
        menu = [['w/' + url_pas(name), load_lang('document')]]
    ))

@app.route('/w/<path:name>')
def read_view(name = None):
    data_none = 0
    sub = ''
    acl = ''
    div = ''

    num = request.args.get('num', None)
    if num:
        num = int(num)

    curs.execute("select sub from rd where title = ? order by date desc", [name])
    for data in curs.fetchall():
        curs.execute("select title from stop where title = ? and sub = ? and close = 'O'", [name, data[0]])
        if not curs.fetchall():
            sub += ' (D)'

            break
                
    curs.execute("select title from data where title like ?", ['%' + name + '/%'])
    if curs.fetchall():
        down = 1
    else:
        down = 0
        
    m = re.search("^(.*)\/(.*)$", name)
    if m:
        uppage = m.groups()[0]
    else:
        uppage = 0
        
    if admin_check(5, None) == 1:
        admin_memu = 1
    else:
        admin_memu = 0
        
    if re.search("^분류:", name):        
        curs.execute("select link from back where title = ? and type = 'cat' order by link asc", [name])
        back = curs.fetchall()
        if back:
            div = '<br><h2 id="cate_normal">분류</h2><ul>'
            u_div = ''
            i = 0

            for data in back:    
                if re.search('^분류:', data[0]):
                    u_div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a></li>'
                elif re.search('^틀:', data[0]):
                    curs.execute("select data from data where title = ?", [data[0]])
                    db_data = curs.fetchall()
                    if db_data:
                        if re.search('\[\[' + name + '#include]]', db_data[0][0]):
                            div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a> <a href="/xref/' + url_pas(data[0]) + '">(' + load_lang('backlink') + ')</a></li>'
                        else:
                            div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a></li>'
                    else:
                        div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a></li>'
                else:
                    div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a></li>'

            div += '</ul>'
            
            if div == '<br><h2 id="cate_normal">분류</h2><ul></ul>':
                div = ''
            
            if u_div != '':
                div += '<br><h2 id="cate_under">하위 분류</h2><ul>' + u_div + '</ul>'


    if num:
        curs.execute("select title from history where title = ? and id = ? and hide = 'O'", [name, str(num)])
        if curs.fetchall() and admin_check(6, None) != 1:
            return redirect('/history/' + url_pas(name))

        curs.execute("select title, data from history where title = ? and id = ?", [name, str(num)])
    else:
        curs.execute("select title, data from data where title = ?", [name])
    
    data = curs.fetchall()
    if data:
        else_data = data[0][1]
        response_data = 200
    else:
        data_none = 1
        response_data = 404
        else_data = ''

    m = re.search("^" + load_lang('user') + ":([^/]*)", name)
    if m:
        g = m.groups()
        
        curs.execute("select acl from user where id = ?", [g[0]])
        test = curs.fetchall()
        if test and test[0][0] != 'user':
            acl = ' (' + load_lang('admin') + ')'
        else:
            curs.execute("select block from ban where block = ?", [g[0]])
            if curs.fetchall():
                sub += ' (' + load_lang('ban') + ')'
            else:
                acl = ''

    curs.execute("select dec from acl where title = ?", [name])
    data = curs.fetchall()
    if data:
        acl += ' (A)'
            
    if request.args.get('froms', None):
        else_data = re.sub('\r\n#(?:redirect|넘겨주기) (?P<in>(?:(?!\r\n).)+)\r\n', ' * Redirect : [[\g<in>]]', '\r\n' + else_data + '\r\n')
        else_data = re.sub('^\r\n', '', else_data)
        else_data = re.sub('\r\n$', '', else_data)
            
    end_data = namumark(name, else_data, 0)
    
    if num:
        menu = [['history/' + url_pas(name), load_lang('history')]]
        sub = ' (' + str(num) + load_lang('version') + ')'
        acl = ''
        r_date = 0
    else:
        if data_none == 1:
            menu = [['edit/' + url_pas(name), load_lang('create')]]
        else:
            menu = [['edit/' + url_pas(name), load_lang('edit')]]

        menu += [['topic/' + url_pas(name), load_lang('discussion')], ['history/' + url_pas(name), load_lang('history')], ['xref/' + url_pas(name), load_lang('backlink')], ['acl/' + url_pas(name), 'ACL']]

        if request.args.get('froms', None):
            menu += [['w/' + url_pas(name), '넘기기']]
            end_data = '<ul id="redirect"><li><a href="/w/' + url_pas(request.args.get('froms', None)) + '?froms=' + url_pas(name) + '">' + request.args.get('froms', None) + '</a>에서 넘어 왔습니다.</li></ul><br>' + end_data

        if uppage != 0:
            menu += [['w/' + url_pas(uppage), '상위']]

        if down:
            menu += [['down/' + url_pas(name), '하위']]
    
        curs.execute("select date from history where title = ? order by date desc limit 1", [name])
        date = curs.fetchall()
        if date:
            r_date = date[0][0]
        else:
            r_date = 0

    div = end_data + div

    return html_minify(render_template(skin_check(), 
        imp = [name, wiki_set(1), custom(), other2([sub + acl, r_date])],
        data = div,
        menu = menu
    )), response_data

@app.route('/topic_record/<name>')
def user_topic_list(name = None):
    num = int(request.args.get('num', 1))
    if num * 50 > 0:
        sql_num = num * 50 - 50
    else:
        sql_num = 0
    
    one_admin = admin_check(1, None)

    div = '<table style="width: 100%; text-align: center;"><tbody><tr>'
    div += '<td style="width: 33.3%;">' + load_lang('discussion') + ' ' + load_lang('name') + '</td><td style="width: 33.3%;">작성자</td><td style="width: 33.3%;">' + load_lang('time') + '</td></tr>'
    
    curs.execute("select title, id, sub, ip, date from topic where ip = ? order by date desc limit ?, '50'", [name, str(sql_num)])
    data_list = curs.fetchall()
    for data in data_list:
        title = html.escape(data[0])
        sub = html.escape(data[2])
        
        if one_admin == 1:
            curs.execute("select * from ban where block = ?", [data[3]])
            if curs.fetchall():
                ban = ' <a href="/ban/' + url_pas(data[3]) + '">(' + load_lang('release') + ')</a>'
            else:
                ban = ' <a href="/ban/' + url_pas(data[3]) + '">(' + load_lang('ban') + ')</a>'
        else:
            ban = ''
            
        div += '<tr><td><a href="/topic/' + url_pas(data[0]) + '/sub/' + url_pas(data[2]) + '#' + data[1] + '">' + title + '#' + data[1] + '</a> (' + sub + ')</td>'
        div += '<td>' + ip_pas(data[3]) + ban + '</td><td>' + data[4] + '</td></tr>'

    div += '</tbody></table>'
    div += next_fix('/topic_record/' + url_pas(name) + '?num=', num, data_list)      
    
    curs.execute("select end from ban where block = ?", [name])
    if curs.fetchall():
        sub = ' (' + load_lang('ban') + ')'
    else:
        sub = 0 
    
    return html_minify(render_template(skin_check(), 
        imp = [load_lang('discussion') + ' ' + load_lang('record'), wiki_set(1), custom(), other2([sub, 0])],
        data = div,
        menu = [['other', load_lang('other')], ['user', load_lang('user')], ['count/' + url_pas(name), '' + load_lang('count') + ''], ['record/' + url_pas(name), load_lang('record')]]
    ))

@app.route('/recent_changes')
@app.route('/<regex("record"):tool>/<name>')
@app.route('/<regex("history"):tool>/<path:name>', methods=['POST', 'GET'])
def recent_changes(name = None, tool = 'record'):
    if request.method == 'POST':
        return redirect('/diff/' + url_pas(name) + '?first=' + request.form.get('b', None) + '&second=' + request.form.get('a', None))
    else:
        one_admin = admin_check(1, None)
        six_admin = admin_check(6, None)
        
        ban = ''
        select = ''

        what = request.args.get('what', 'all')

        div = '<table style="width: 100%; text-align: center;"><tbody><tr>'
        
        if name:
            num = int(request.args.get('num', 1))
            if num * 50 > 0:
                sql_num = num * 50 - 50
            else:
                sql_num = 0      

            if tool == 'history':
                div += '<td style="width: 33.3%;">' + load_lang('version') + '</td><td style="width: 33.3%;">' + load_lang('editor') + '</td><td style="width: 33.3%;">' + load_lang('time') + '</td></tr>'
                
                curs.execute("select id, title, date, ip, send, leng from history where title = ? order by id + 0 desc limit ?, '50'", [name, str(sql_num)])
            else:
                div += '<td style="width: 33.3%;">' + load_lang('document') + ' ' + load_lang('name') + '</td><td style="width: 33.3%;">' + load_lang('editor') + '</td><td style="width: 33.3%;">' + load_lang('time') + '</td></tr>'

                if what == 'all':
                    div = '<a href="/record/' + url_pas(name) + '?what=revert">(' + load_lang('revert') + ')</a><hr>' + div
                    div = '<a href="/record/' + url_pas(name) + '?what=move">(' + load_lang('move') + ')</a> ' + div
                    div = '<a href="/record/' + url_pas(name) + '?what=delete">(' + load_lang('delete') + ')</a> ' + div
                    
                    curs.execute("select id, title, date, ip, send, leng from history where ip = ? order by date desc limit ?, '50'", [name, str(sql_num)])
                else:
                    if what == 'delete':
                        sql = '%(' + load_lang('delete') + ')'
                    elif what == 'move':
                        sql = '%' + load_lang('move') + ')'
                    elif what == 'revert':
                        sql = '%' + load_lang('version') + ')'
                    else:
                        return redirect('/')

                    curs.execute("select id, title, date, ip, send, leng from history where ip = ? and send like ? order by date desc limit ?, '50'", [name, sql, str(sql_num)])
        else:
            div += '<td style="width: 33.3%;">' + load_lang('document') + ' ' + load_lang('name') + '</td><td style="width: 33.3%;">' + load_lang('editor') + '</td><td style="width: 33.3%;">' + load_lang('time') + '</td></tr>'
            
            if what == 'all':
                div = '<a href="/recent_changes?what=revert">(' + load_lang('revert') + ')</a><hr>' + div
                div = '<a href="/recent_changes?what=move">(' + load_lang('move') + ')</a> ' + div
                div = '<a href="/recent_changes?what=delete">(' + load_lang('delete') + ')</a> ' + div

                div = '<a href="/recent_discuss">(' + load_lang('discussion') + ')</a> <a href="/block_log">(' + load_lang('ban') + ')</a> <a href="/user_log">(' + load_lang('register') + ')</a> <a href="/admin_log">(' + load_lang('authority') + ')</a><hr>' + div
                
                curs.execute("select id, title, date, ip, send, leng from history order by date desc limit 50")
            else:
                if what == 'delete':
                    sql = '%(' + load_lang('delete') + ')'
                elif what == 'move':
                    sql = '%' + load_lang('move') + ')'
                elif what == 'revert':
                    sql = '%' + load_lang('version') + ')'
                else:
                    return redirect('/')

                curs.execute("select id, title, date, ip, send, leng from history where send like ? order by date desc limit 50", [sql])

        data_list = curs.fetchall()
        for data in data_list:    
            select += '<option value="' + data[0] + '">' + data[0] + '</option>'     
            send = '<br>'
            
            if data[4]:
                if not re.search("^(?: *)$", data[4]):
                    send = data[4]
            
            if re.search("\+", data[5]):
                leng = '<span style="color:green;">(' + data[5] + ')</span>'
            elif re.search("\-", data[5]):
                leng = '<span style="color:red;">(' + data[5] + ')</span>'
            else:
                leng = '<span style="color:gray;">(' + data[5] + ')</span>'
                
            if one_admin == 1:
                curs.execute("select * from ban where block = ?", [data[3]])
                if curs.fetchall():
                    ban = ' <a href="/ban/' + url_pas(data[3]) + '">(' + load_lang('release') + ')</a>'
                else:
                    ban = ' <a href="/ban/' + url_pas(data[3]) + '">(' + load_lang('ban') + ')</a>'            
                
            ip = ip_pas(data[3])
            if int(data[0]) - 1 == 0:
                revert = ''
            else:
                revert = '<a href="/diff/' + url_pas(data[1]) + '?first=' + str(int(data[0]) - 1) + '&second=' + data[0] + '">(' + load_lang('compare') + ')</a> <a href="/revert/' + url_pas(data[1]) + '?num=' + str(int(data[0]) - 1) + '">(' + load_lang('revert') + ')</a>'
            
            style = ['', '']
            date = data[2]

            curs.execute("select title from history where title = ? and id = ? and hide = 'O'", [data[1], data[0]])
            hide = curs.fetchall()
            
            if six_admin == 1:
                if hide:                            
                    hidden = ' <a href="/hidden/' + url_pas(data[1]) + '?num=' + data[0] + '">(공개)'
                    
                    style[0] = 'background: gainsboro;'
                    style[1] = 'background: gainsboro;'
                    
                    if send == '<br>':
                        send = '(' + load_lang('hide') + ')'
                    else:
                        send += ' (' + load_lang('hide') + ')'
                else:
                    hidden = ' <a href="/hidden/' + url_pas(data[1]) + '?num=' + data[0] + '">(' + load_lang('hide') + ')'
            elif not hide:
                hidden = ''
            else:
                ip = ''
                hidden = ''
                ban = ''
                date = ''

                send = '(' + load_lang('hide') + ')'

                style[0] = 'display: none;'
                style[1] = 'background: gainsboro;'

            if tool == 'history':
                title = '<a href="/w/' + url_pas(name) + '?num=' + data[0] + '">' + data[0] + load_lang('version') + '</a> <a href="/raw/' + url_pas(name) + '?num=' + data[0] + '">(Raw)</a> '
            else:
                title = '<a href="/w/' + url_pas(data[1]) + '">' + html.escape(data[1]) + '</a> <a href="/history/' + url_pas(data[1]) + '">(' + data[0] + load_lang('version') + ')</a> '
                    
            div += '<tr style="' + style[0] + '"><td>' + title + revert + ' ' + leng + '</td>'
            div += '<td>' + ip + ban + hidden + '</td><td>' + date + '</td></tr><tr style="' + style[1] + '"><td colspan="3">' + send + '</td></tr>'

        div += '</tbody></table>'
        sub = ''

        if name:
            if tool == 'history':
                div = '<form method="post"><select name="a">' + select + '</select> <select name="b">' + select + '</select> <button type="submit">' + load_lang('compare') + '</button></form><hr>' + div
                title = name
                
                sub += ' (' + load_lang('history') + ')'
                
                menu = [['w/' + url_pas(name), load_lang('document')], ['raw/' + url_pas(name), 'Raw'], ['move_data/' + url_pas(name), load_lang('move') + ' ' + load_lang('history')]]
                
                div += next_fix('/history/' + url_pas(name) + '?num=', num, data_list)
            else:
                curs.execute("select end from ban where block = ?", [name])
                if curs.fetchall():
                    sub += ' (' + load_lang('ban') + ')'

                title = '' + load_lang('edit') + ' ' + load_lang('record')
                
                menu = [['other', load_lang('other')], ['user', load_lang('user')], ['count/' + url_pas(name), '' + load_lang('count') + ''], ['topic_record/' + url_pas(name), load_lang('discussion')]]
                
                div += next_fix('/record/' + url_pas(name) + '/' + url_pas(what) + '?num=', num, data_list)
                
                if what != 'all':
                    menu += [['record/' + url_pas(name), load_lang('normal')]]
        else:
            menu = 0
            title = '' + load_lang('recent') + ' ' + load_lang('change') + ''
            
            if what != 'all':
                menu = [['recent_changes', load_lang('normal')]]
                
        if what == 'delete':
            sub += ' (' + load_lang('delete') + ')'
        elif what == 'move':
            sub += ' (' + load_lang('move') + ')'
        elif what == 'revert':
            sub += ' (' + load_lang('revert') + ')'
        
        if sub == '':
            sub = 0
                
        return html_minify(render_template(skin_check(), 
            imp = [title, wiki_set(1), custom(), other2([sub, 0])],
            data = div,
            menu = menu
        ))
    
@app.route('/upload', methods=['GET', 'POST'])
def upload():
    if ban_check() == 1:
        return re_error('/ban')
    
    if request.method == 'POST':
        if captcha_post(request.form.get('g-recaptcha-response', None), conn) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        data = request.files.get('f_data', None)
        if not data:
            return re_error('/error/9')

        if int(wiki_set(3)) * 1024 * 1024 < request.content_length:
            return re_error('/error/17')
        
        value = os.path.splitext(data.filename)[1]
        if not value in ['.jpeg', '.jpg', '.gif', '.png', '.webp', '.JPEG', '.JPG', '.GIF', '.PNG', '.WEBP']:
            return re_error('/error/14')
    
        if request.form.get('f_name', None):
            name = request.form.get('f_name', None) + value
        else:
            name = data.filename
        
        piece = os.path.splitext(name)
        if re.search('[^ㄱ-힣0-9a-zA-Z_\- ]', piece[0]):
            return re_error('/error/22')

        e_data = sha224(piece[0]) + piece[1]

        curs.execute("select title from data where title = ?", ['' + load_lang('file') + ':' + name])
        if curs.fetchall():
            return re_error('/error/16')
            
        ip = ip_check()

        if request.form.get('f_lice', None):
            lice = request.form.get('f_lice', None)
        else:
            if custom()[2] == 0:
                lice = ip + ' Upload'
            else:
                lice = '[[' + load_lang('user') + ':' + ip + ']] Upload'
            
        if os.path.exists(os.path.join('image', e_data)):
            os.remove(os.path.join('image', e_data))
            
            data.save(os.path.join('image', e_data))
        else:
            data.save(os.path.join('image', e_data))
            
        curs.execute("select title from data where title = ?", ['' + load_lang('file') + ':' + name])
        if curs.fetchall(): 
            curs.execute("delete from data where title = ?", ['' + load_lang('file') + ':' + name])
        
        curs.execute("insert into data (title, data) values (?, ?)", ['' + load_lang('file') + ':' + name, '[[' + load_lang('file') + ':' + name + ']][br][br]{{{[[' + load_lang('file') + ':' + name + ']]}}}[br][br]' + lice])
        curs.execute("insert into acl (title, dec, dis, why) values (?, 'admin', '', '')", ['' + load_lang('file') + ':' + name])

        history_plus('' + load_lang('file') + ':' + name, '[[' + load_lang('file') + ':' + name + ']][br][br]{{{[[' + load_lang('file') + ':' + name + ']]}}}[br][br]' + lice, get_time(), ip, '(Upload)', '0')
        
        conn.commit()
        
        return redirect('/w/' + load_lang('file') + ':' + name)      
    else:
        return html_minify(render_template(skin_check(), 
            imp = [load_lang('upload'), wiki_set(1), custom(), other2([0, 0])],
            data =  '''
                    <form method="post" enctype="multipart/form-data" accept-charset="utf8">
                        <input type="file" name="f_data">
                        <hr>
                        <input placeholder="''' + load_lang('name') + '''" name="f_name" type="text">
                        <hr>
                        <input placeholder="''' + load_lang('license') + '''" name="f_lice" type="text">
                        <hr>
                        ''' + captcha_get() + '''
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                    ''',
            menu = [['other', load_lang('other')]]
        ))  
        
@app.route('/user')
def user_info():
    ip = ip_check()
    
    curs.execute("select acl from user where id = ?", [ip])
    data = curs.fetchall()
    if ban_check() == 0:
        if data:
            if data[0][0] != 'user':
                acl = data[0][0]
            else:
                acl = load_lang('subscriber')
        else:
            acl = load_lang('normal')
    else:
        acl = load_lang('ban')

        match = re.search("^([0-9]{1,3}\.[0-9]{1,3})", ip)
        if match:
            match = match.groups()[0]
        else:
            match = 'Not'

        curs.execute("select end, login, band from ban where block = ? or block = ?", [ip, match])
        block_data = curs.fetchall()
        if block_data:
            if block_data[0][0] != '':
                acl += ' (' + block_data[0][0] + '까지)'
            else:
                acl += ' (무기한)'        

            if block_data[0][1] != '':
                acl += ' (' + load_lang('login') + ' ' + load_lang('able') + ')'

            if block_data[0][2] == 'O':
                acl += ' (대역)'
            
    if custom()[2] != 0:
        ip_user = '<a href="/w/' + load_lang('user') + ':' + ip + '">' + ip + '</a>'
        
        plus = '<li><a href="/logout">로그아웃</a></li><li><a href="/change">' + load_lang('my_info') + ' ' + load_lang('edit') + '</a></li>'
        
        curs.execute('select name from alarm where name = ? limit 1', [ip_check()])
        if curs.fetchall():
            plus2 = '<li><a href="/alarm">' + load_lang('alarm') + ' (O)</a></li>'
        else:
            plus2 = '<li><a href="/alarm">' + load_lang('alarm') + '</a></li>'

        plus2 += '<li><a href="/watch_list">' + load_lang('watchlist') + '</a></li>'
    else:
        ip_user = ip
        
        plus = '<li><a href="/login">' + load_lang('login') + '</a></li>'
        plus2 = ''

    return html_minify(render_template(skin_check(), 
        imp = [load_lang('user') + ' ' + load_lang('tool'), wiki_set(1), custom(), other2([0, 0])],
        data =  '''
                <h2>상태</h2>
                <ul>
                    <li>''' + ip_user + ''' <a href="/record/''' + url_pas(ip) + '''">(''' + load_lang('record') + ''')</a></li><li>''' + load_lang('authority') + ''' : ''' + acl + '''</li>
                </ul>
                <br>
                <h2>''' + load_lang('login') + '''</h2>
                <ul>
                    ''' + plus + '''
                    <li><a href="/register">''' + load_lang('register') + '''</a></li>
                </ul>
                <br>
                <h2>''' + load_lang('tool') + '''</h2>
                <ul>
                    <li><a href="/acl/''' + load_lang('user') + ':' + url_pas(ip) + '">' + load_lang('user') + ' ' + load_lang('document') + ''' ACL</a></li>
                    <li><a href="/custom_head">''' + load_lang('user') + ''' HEAD</a></li></ul><br><h2>''' + load_lang('other') + '''</h2><ul>''' + plus2 + '''<li><a href="/count">''' + load_lang('count') + '''</a></li>
                </ul>
                ''',
        menu = 0
    ))

@app.route('/watch_list')
def watch_list():
    div = 'Limit : 10<hr>'
    
    if custom()[2] == 0:
        return redirect('/login')

    curs.execute("select title from scan where user = ?", [ip_check()])
    data = curs.fetchall()
    for data_list in data:
        div += '<li><a href="/w/' + url_pas(data_list[0]) + '">' + data_list[0] + '</a> <a href="/watch_list/' + url_pas(data_list[0]) + '">(' + load_lang('delete') + ')</a></li>'

    if data:
        div = '<ul>' + div + '</ul><hr>'

    div += '<a href="/manager/13">(' + load_lang('plus') + ')</a>'

    return html_minify(render_template(skin_check(), 
        imp = ['' + load_lang('watchlist') + ' ' + load_lang('list'), wiki_set(1), custom(), other2([0, 0])],
        data = div,
        menu = [['manager', load_lang('admin')]]
    ))

@app.route('/watch_list/<path:name>')
def watch_list_name(name = None):
    if custom()[2] == 0:
        return redirect('/login')

    ip = ip_check()

    curs.execute("select count(title) from scan where user = ?", [ip])
    count = curs.fetchall()
    if count and count[0][0] > 9:
        return redirect('/watch_list')

    curs.execute("select title from scan where user = ? and title = ?", [ip, name])
    if curs.fetchall():
        curs.execute("delete from scan where user = ? and title = ?", [ip, name])
    else:
        curs.execute("insert into scan (user, title) values (?, ?)", [ip, name])
    
    conn.commit()

    return redirect('/watch_list')

@app.route('/custom_head', methods=['GET', 'POST'])
def custom_head_view():
    ip = ip_check()

    if request.method == 'POST':
        if custom()[2] != 0:
            curs.execute("select user from custom where user = ?", [ip + ' (head)'])
            if curs.fetchall():
                curs.execute("update custom set css = ? where user = ?", [request.form.get('content', None), ip + ' (head)'])
            else:
                curs.execute("insert into custom (user, css) values (?, ?)", [ip + ' (head)', request.form.get('content', None)])
            
            conn.commit()

        session['MyMaiToNight'] = request.form.get('content', None)

        return redirect('/user')
    else:
        if custom()[2] != 0:
            start = ''

            curs.execute("select css from custom where user = ?", [ip + ' (head)'])
            head_data = curs.fetchall()
            if head_data:
                data = head_data[0][0]
            else:
                data = ''
        else:
            start = '<span>' + load_lang('user_css_warring') + '</span><hr>'
            
            if 'MyMaiToNight' in session:
                data = session['MyMaiToNight']
            else:
                data = ''

        start += '<span>&lt;style&gt;CSS&lt;/style&gt;<br>&lt;script&gt;JS&lt;/script&gt;</span><hr>'

        return html_minify(render_template(skin_check(), 
            imp = [load_lang('user') + ' HEAD', wiki_set(1), custom(), other2([0, 0])],
            data =  start + '''
                    <form method="post">
                        <textarea rows="25" cols="100" name="content">''' + data + '''</textarea>
                        <hr>
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                    ''',
            menu = [['user', load_lang('user')]]
        ))

@app.route('/count')
@app.route('/count/<name>')
def count_edit(name = None):
    if name == None:
        that = ip_check()
    else:
        that = name

    curs.execute("select count(title) from history where ip = ?", [that])
    count = curs.fetchall()
    if count:
        data = count[0][0]
    else:
        data = 0

    curs.execute("select count(title) from topic where ip = ?", [that])
    count = curs.fetchall()
    if count:
        t_data = count[0][0]
    else:
        t_data = 0

    return html_minify(render_template(skin_check(), 
        imp = ['' + load_lang('count') + '', wiki_set(1), custom(), other2([0, 0])],
        data = '''
                <ul>
                    <li><a href="/record/''' + url_pas(that) + '''">''' + load_lang('edit') + '''</a> : ''' + str(data) + '''</li>
                    <li><a href="/topic_record/''' + url_pas(that) + '''">''' + load_lang('discussion') + '''</a> : ''' + str(t_data) + '''</a></li>
                </ul>
                ''',
        menu = [['user', load_lang('user')]]
    ))
        
@app.route('/random')
def random():
    curs.execute("select title from data order by random() limit 1")
    data = curs.fetchall()
    if data:
        return redirect('/w/' + url_pas(data[0][0]))
    else:
        return redirect('/')
    
@app.route('/views/<path:name>')
def views(name = None):
    if re.search('\/', name):
        m = re.search('^(.*)\/(.*)$', name)
        if m:
            n = m.groups()
            plus = '/' + n[0]
            rename = n[1]
        else:
            plus = ''
            rename = name
    else:
        plus = ''
        rename = name

    m = re.search('\.(.+)$', name)
    if m:
        g = m.groups()
    else:
        g = ['']

    if g == 'css':
        return css_minify(send_from_directory('./views' + plus, rename))   
    elif g == 'js':
        return js_minify(send_from_directory('./views' + plus, rename))
    elif g == 'html':
        return html_minify(send_from_directory('./views' + plus, rename))   
    else:
        return send_from_directory('./views' + plus, rename)

@app.route('/<data>')
def main_file(data = None):
    if re.search('\.(txt|html)$', data):
        return send_from_directory('./', data)
    else:
        return ''

@app.errorhandler(404)
def error_404(e):
    return '''<!-- 
            나니카가 하지마룻테 코토와 오와리니 츠나가루다난테 캉가에테모 미나캇타. 이야, 캉카에타쿠나캇탄다...
            아마오토 마도오 타타쿠 소라카라 와타시노 요-나 카나시미 훗테루 토메도나쿠 이마오 누라시테
            오모이데 난테 이라나이노 코코로가 쿠루시쿠나루 다케다토 No more! September Rain No more! September Rain
            
            이츠닷테 아나타와 미짓카닷타 와자와자 키모치오 타시카메룻테 코토모 히츠요-쟈나쿠테
            시젠니 나카라요쿠 나레타카라 안신시테타노 카모시레나이네 도-시테? 나미니 토이카케루케도
            나츠노 하지마리가 츠레테키타 오모이 나츠가 오와루토키 키에챠우모노닷타 난테 시라나쿠테
            토키메이테타 아츠이 키세츠 우미베노 소라가 히캇테 토츠젠 쿠모가 나가레 오츠부노 아메 와타시노 나카노 나미다미타이
            콘나니 타노시이 나츠가 즛토 츠즈이테쿳테 신지테타요 But now... September Rain But now... September Rain
            -->''' + redirect('/w/' + url_pas(wiki_set(2)))

if __name__=="__main__":
    app.secret_key = rep_key
    http_server = HTTPServer(WSGIContainer(app))
    http_server.listen(rep_port)
    IOLoop.instance().start()
