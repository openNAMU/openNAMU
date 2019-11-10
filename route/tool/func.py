import os
import sys
import platform

for i in range(0, 2):
    try:
        import werkzeug.routing
        import werkzeug.debug
        import flask_compress
        import flask_reggie
        import tornado.ioloop
        import tornado.httpserver
        import tornado.wsgi
        import urllib.request
        import email.mime.text
        import sqlite3
        import hashlib
        import smtplib
        import bcrypt
        import zipfile
        import difflib
        import shutil
        import threading
        import logging
        import random
        import flask
        import json
        import html
        import re

        if sys.version_info < (3, 6):
            import sha3

        from .set_mark.tool import *
        from .mark import *
    except ImportError as e:
        if i == 0:
            print('----')
            print(e)
            print('----')
            if platform.system() == 'Linux':
                ok = os.system('python3 -m pip install --user -r requirements.txt')
                if ok == 0:
                    print('----')
                    os.execl(sys.executable, sys.executable, *sys.argv)
                else:
                    raise
            elif platform.system() == 'Windows':
                ok = os.system('python -m pip install --user -r requirements.txt')
                if ok == 0:
                    print('----')
                    os.execl(sys.executable, sys.executable, *sys.argv)
                else:
                    raise
            else:
                print('----')
                print(e)
                raise
        else:
            print('----')
            print(e)
            raise

app_var = json.loads(open('data/app_var.json', encoding='utf-8').read())

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
        print('----')
        print('Error : Email send error')

def last_change(data):
    json_address = re.sub("(((?!\.|\/).)+)\.html$", "set.json", skin_check())
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
    return last_change(data)

def render_set(title = '', data = '', num = 0, s_data = 0, include = None):
    if acl_check(title, 'render') == 1:
        return 'HTTP Request 401.3'
    elif s_data == 1:
        return data
    else:
        if data != None:
            return namumark(title, data, num, include)
        else:
            return 'HTTP Request 404'

def captcha_get():
    data = ''

    if ip_or_user() != 0:
        curs.execute('select data from other where name = "recaptcha"')
        recaptcha = curs.fetchall()
        if recaptcha and recaptcha[0][0] != '':
            curs.execute('select data from other where name = "sec_re"')
            sec_re = curs.fetchall()
            if sec_re and sec_re[0][0] != '':
                data += recaptcha[0][0] + '<hr class=\"main_hr\">'

    return data

def update():
    #v3.1.2
    try:
        curs.execute('select title, dec from acl where dec != ""')
        db_data = curs.fetchall()
        for i in db_data:
            curs.execute("update acl set decu = ? where title = ?", [i[1], i[0]])

        print('Fix table acl column dec to decu')
        print('----')
    except:
        pass

    conn.commit()

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
        if ip_or_user() != 0 and captcha_get() != '':
            curs.execute('select data from other where name = "sec_re"')
            sec_re = curs.fetchall()
            if sec_re and sec_re[0][0] != '':
                try:
                    data = urllib.request.urlopen('https://www.google.com/recaptcha/api/siteverify?secret=' + sec_re[0][0] + '&response=' + re_data)
                except:
                    pass
                    
                if data and data.getcode() == 200:
                    json_data = json.loads(data.read().decode(data.headers.get_content_charset()))
                    if json_data['success'] == True:
                        return 0
                    else:
                        return 1
                else:
                    return 0
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
    oauth = json.loads(open(app_var['path_oauth_setting'], encoding='utf-8').read())

    return oauth[provider]

def update_oauth(provider, target, content):
    oauth = json.loads(open(app_var['path_oauth_setting'], encoding='utf-8').read())
    oauth[provider][target] = content

    with open(app_var['path_oauth_setting'], 'w') as f:
        f.write(json.dumps(oauth, sort_keys=True, indent=4))

    return 'Done'

def ip_or_user(data = ''):
    if data == '':
        data = ip_check()

    if re.search('(\.|:)', data):
        return 1
    else:
        return 0

def edit_button():
    insert_list = []
    
    curs.execute("select html, plus from html_filter where kind = 'edit_top'")
    db_data = curs.fetchall()
    for get_data in db_data:
        insert_list += [[get_data[1], get_data[0]]]

    data = ''
    for insert_data in insert_list:
        data += '<a href="javascript:do_insert_data(\'content\', \'' + insert_data[0] + '\')">(' + insert_data[1] + ')</a> '

    if admin_check() == 1:
        data += (' ' if data != '' else '') + '<a href="/edit_top">(' + load_lang('add') + ')</a>'

    return data + '<hr class=\"main_hr\">'

def ip_warring():
    if ip_or_user() != 0:
        curs.execute('select data from other where name = "no_login_warring"')
        data = curs.fetchall()
        if data and data[0][0] != '':
            text_data = '<span>' + data[0][0] + '</span><hr class=\"main_hr\">'
        else:
            text_data = '<span>' + load_lang('no_login_warring') + '</span><hr class=\"main_hr\">'
    else:
        text_data = ''

    return text_data

def skin_check(set_n = 0):
    skin = 'marisa'

    curs.execute('select data from other where name = "skin"')
    skin_exist = curs.fetchall()
    if skin_exist and skin_exist[0][0] != '':
        if os.path.exists(os.path.abspath('./views/' + skin_exist[0][0] + '/index.html')) == 1:
            skin = skin_exist[0][0]
    
    curs.execute('select data from user_set where name = "skin" and id = ?', [ip_check()])
    skin_exist = curs.fetchall()
    if skin_exist and skin_exist[0][0] != '':
        if os.path.exists(os.path.abspath('./views/' + skin_exist[0][0] + '/index.html')) == 1:
            skin = skin_exist[0][0]

    if set_n == 0:
        return './views/' + skin + '/index.html'
    else:
        return skin

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
    req_list = ''
    
    css_filter = {}
    for i_data in os.listdir(os.path.join("views", "main_css", "css")):
        if i_data in css_filter:
            req_list += '<link rel="stylesheet" href="/views/main_css/css/' + i_data + '?ver=' + css_filter[i_data] + '">'
        else:
            req_list += '<link rel="stylesheet" href="/views/main_css/css/' + i_data + '?ver=1">'
    
    js_filter = {
        'load_include.js' : '2',
        'render_html.js' : '2',
        'do_open_foot.js' : '4'
    }
    for i_data in os.listdir(os.path.join("views", "main_css", "js")):
        if i_data in js_filter:
            req_list += '<script src="/views/main_css/js/' + i_data + '?ver=' + js_filter[i_data] + '"></script>'
        else:
            req_list += '<script src="/views/main_css/js/' + i_data + '?ver=1"></script>'

    data += ['', '''
        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.12.0/styles/default.min.css">
        <link   rel="stylesheet"
                href="https://cdn.jsdelivr.net/npm/katex@0.10.1/dist/katex.min.css"
                integrity="sha384-dbVIfZGuN1Yq7/1Ocstc1lUEm+AT+/rCkibIcC/OmWo5f0EA48Vf8CytHzGrSwbQ"
                crossorigin="anonymous">
        <script src="https://cdn.jsdelivr.net/npm/katex@0.10.1/dist/katex.min.js"
                integrity="sha384-2BKqo+exmr9su6dir+qCw08N2ZKRucY4PrGQPPWU1A7FtlCGjmEGFqXCv5nyM5Ij"
                crossorigin="anonymous"></script>
        <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.12.0/highlight.min.js"></script>
    ''' + req_list]

    return data

def cut_100(data):
    data = re.sub('<(((?!>).)*)>', ' ', data)
    data = re.sub('\n', ' ', data)
    data = re.sub('^ +', '', data)
    data = re.sub(' +$', '', data)
    data = re.sub(' {2,}', ' ', data)

    return data[0:100] + '...'

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
            
        curs.execute("select data from other where name = 'head' and coverage = ?", [skin_check(1)])
        db_data = curs.fetchall()
        if db_data and db_data[0][0] != '':
            if len(re.findall('<', db_data[0][0])) % 2 != 1:
                data_list += [db_data[0][0]]
            else:
                data_list += ['']
        else:
            curs.execute("select data from other where name = 'head' and coverage = ''")
            db_data = curs.fetchall()
            if db_data and db_data[0][0] != '':
                if len(re.findall('<', db_data[0][0])) % 2 != 1:
                    data_list += [db_data[0][0]]
                else:
                    data_list += ['']
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

    if not re.search('\n$', end):
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
           
def admin_check(num = None, what = None, name = ''):
    if name == '':
        ip = ip_check()
    else:
        ip = name

    curs.execute("select acl from user where id = ?", [ip])
    user = curs.fetchall()
    if user:
        reset = 0

        back_num = num
        while 1:
            if num == 1:
                check = 'ban'
            elif num == 2:
                check = 'nothing'
            elif num == 3:
                check = 'toron'
            elif num == 4:
                check = 'check'
            elif num == 5:
                check = 'acl'
            elif num == 6:
                check = 'hidel'
            elif num == 7:
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
                if back_num == 'all':
                    if num == 'all':
                        num = 1
                    elif num != 8:
                        num += 1
                    else:
                        break
                elif num:
                    num = None
                else:
                    break
                    
    return 0

def ip_pas(raw_ip):
    hide = 0

    if ip_or_user(raw_ip) != 0:    
        curs.execute("select data from other where name = 'ip_view'")
        data = curs.fetchall()
        if data and data[0][0] != '':
            ip = re.sub('((?:(?!\.).)+)$', 'xxx', raw_ip)

            if not admin_check(1):
                hide = 1
        else:
            ip = raw_ip
    else:
        curs.execute("select title from data where title = ?", ['user:' + raw_ip])
        if curs.fetchall():
            ip = '<a href="/w/' + url_pas('user:' + raw_ip) + '">' + raw_ip + '</a>'
        else:
            ip = '<a id="not_thing" href="/w/' + url_pas('user:' + raw_ip) + '">' + raw_ip + '</a>'

        if admin_check('all', None, raw_ip) == 1:
            ip = '<b>' + ip + '</b>'

    if ban_check(raw_ip) == 1:
        ip = '<s>' + ip + '</s>'
         
    if hide == 0:
        ip += ' <a href="/tool/' + url_pas(raw_ip) + '">(' + load_lang('tool') + ')</a>'

    return ip

def custom():
    if 'head' in flask.session:
        if len(re.findall('<', flask.session['head'])) % 2 != 1:
            user_head = flask.session['head']
        else:
            user_head = ''
    else:
        user_head = ''

    ip = ip_check()
    if ip_or_user(ip) == 0:
        user_icon = 1
    else:
        user_icon = 0

    if user_icon != 0:
        curs.execute('select data from user_set where name = "email" and id = ?', [ip])
        data = curs.fetchall()
        if data:
            email = data[0][0]
        else:
            email = ''
    else:
        email = ''

    if user_icon != 0:
        user_name = ip
    else:
        user_name = load_lang('user')
        
    if admin_check('all') == 1:
        user_admin = '1'

        curs.execute("select acl from user where id = ?", [ip])
        user_acl = curs.fetchall()

        user_acl_list = []
        curs.execute('select acl from alist where name = ?', [user_acl[0][0]])
        user_acl = curs.fetchall()
        for i in user_acl:
            user_acl_list += [i[0]]

        if user_acl != []:
            user_acl_list = user_acl_list
        else:
            user_acl_list = '0'
    else:
        user_admin = '0'
        user_acl_list = '0'
        
    if ban_check() == 1:
        user_ban = '1'
    else:
        user_ban = '0'

    if user_icon == 1:
        curs.execute("select count(name) from alarm where name = ?", [ip])
        count = curs.fetchall()
        if count:
            user_notice = str(count[0][0])
        else:
            user_notice = '0'
    else:
        user_notice = '0'

    return ['', '', user_icon, user_head, email, user_name, user_admin, user_ban, user_notice, user_acl_list]

def load_skin(data = '', set_n = 0):
    div2 = ''
    div3 = []
    system_file = ['main_css', 'easter_egg.html']

    if data == '':
        ip = ip_check()

        curs.execute('select data from user_set where name = "skin" and id = ?', [ip])
        data = curs.fetchall()

        if not data:
            curs.execute('select data from other where name = "skin"')
            data = curs.fetchall()
            if not data:
                data = [['marisa']]

        if set_n == 0:
            for skin_data in os.listdir(os.path.abspath('views')):
                if not skin_data in system_file:
                    if data[0][0] == skin_data:
                        div2 = '<option value="' + skin_data + '">' + skin_data + '</option>' + div2
                    else:
                        div2 += '<option value="' + skin_data + '">' + skin_data + '</option>'
        else:
            div2 = []
            for skin_data in os.listdir(os.path.abspath('views')):
                if not skin_data in system_file:
                    if data[0][0] == skin_data:
                        div2 = [skin_data] + div2
                    else:
                        div2 += [skin_data]
    else:
        if set_n == 0:
            for skin_data in os.listdir(os.path.abspath('views')):
                if not skin_data in system_file:
                    if data == skin_data:
                        div2 = '<option value="' + skin_data + '">' + skin_data + '</option>' + div2
                    else:
                        div2 += '<option value="' + skin_data + '">' + skin_data + '</option>'
        else:
            div2 = []
            for skin_data in os.listdir(os.path.abspath('views')):
                if not skin_data in system_file:
                    if data == skin_data:
                        div2 = [skin_data] + div2
                    else:
                        div2 += [skin_data]

    return div2

def acl_check(name = 'test', tool = '', sub = 'test'):
    ip = ip_check()
    
    if ban_check() == 1:
        return 1

    if tool != 'topic' and tool != 'render':
        acl_c = re.search("^user:((?:(?!\/).)*)", name)
        if acl_c:
            acl_n = acl_c.groups()

            if admin_check(5) == 1:
                return 0

            curs.execute("select decu from acl where title = ?", ['user:' + acl_n[0]])
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

        if re.search("^file:", name) and admin_check(None, 'file edit (' + name + ')') != 1:
            return 1

    if tool == '':
        end = 3
    elif tool == 'topic':
        end = 2
    else:
        end = 1

    for i in range(0, end):
        if tool == '':
            if i == 0:
                curs.execute("select decu from acl where title = ?", [name])
                acl_data = curs.fetchall()
            elif i == 1:
                curs.execute('select data from other where name = "edit"')
                acl_data = curs.fetchall()
            else:
                curs.execute("select view from acl where title = ?", [name])
                acl_data = curs.fetchall()

            num = 5
        elif tool == 'topic':
            if i == 0:
                curs.execute("select dis from acl where title = ?", [name])
                acl_data = curs.fetchall()
            else:
                curs.execute('select data from other where name = "discussion"')
                acl_data = curs.fetchall()

            num = 3
        else:
            curs.execute("select view from acl where title = ?", [name])
            acl_data = curs.fetchall()

            num = 5

        if acl_data and acl_data[0][0] != 'normal':
            if acl_data[0][0] == 'user':
                if ip_or_user(ip) == 1:
                    return 1

            if acl_data[0][0] == 'admin':
                if ip_or_user(ip) == 1:
                    return 1

                if admin_check(num) != 1:
                    return 1

            if acl_data[0][0] == '50_edit':
                if ip_or_user(ip) == 1:
                    return 1
                
                if admin_check(num) != 1:
                    curs.execute("select count(title) from history where ip = ?", [ip])
                    count = curs.fetchall()
                    if count:
                        count = count[0][0]
                    else:
                        count = 0

                    if count < 50:
                        return 1

            if acl_data[0][0] == 'email':
                if ip_or_user(ip) == 1:
                    return 1
                
                if admin_check(num) != 1:
                    curs.execute("select data from user_set where id = ? and name = 'email'", [ip])
                    email = curs.fetchall()
                    if not email:
                        return 1

            if acl_data[0][0] == 'owner':
                if admin_check() != 1:
                    return 1

        if tool == 'topic':
            curs.execute("select title from rd where title = ? and sub = ? and not stop = ''", [name, sub])
            if curs.fetchall():
                if admin_check(3, 'topic (' + name + ')') != 1:
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

    curs.execute("delete from ban where (end < ? and end like '2%')", [get_time()])
    conn.commit()

    curs.execute("select login, block from ban where ((end > ? and end like '2%') or end = '') and band = 'regex'", [get_time()])
    regex_d = curs.fetchall()
    for test_r in regex_d:
        g_regex = re.compile(test_r[1])
        if g_regex.search(ip):
            if tool and tool == 'login':
                if test_r[0] != 'O':
                    return 1
            else:
                return 1
    
    curs.execute("select login from ban where ((end > ? and end like '2%') or end = '') and block = ? and band = 'O'", [get_time(), band_it])
    band_d = curs.fetchall()
    if band_d:
        if tool and tool == 'login':
            if band_d[0][0] != 'O':
                return 1
        else:
            return 1

    curs.execute("select login from ban where ((end > ? and end like '2%') or end = '') and block = ? and band = ''", [get_time(), ip])
    ban_d = curs.fetchall()
    if ban_d:
        if tool and tool == 'login':
            if ban_d[0][0] != 'O':
                return 1
        else:
            return 1

    return 0
        
def topic_check(name, sub):
    return acl_check(name, 'topic', sub)

def ban_insert(name, end, why, login, blocker, type_d = None):
    now_time = get_time()

    if type_d:
        band = type_d
    else:
        if re.search("^([0-9]{1,3}\.[0-9]{1,3})$", name):
            band = 'O'
        else:
            band = ''

    curs.execute("delete from ban where (end < ? and end like '2%')", [get_time()])

    curs.execute("select block from ban where ((end > ? and end like '2%') or end = '') and block = ? and band = ?", [get_time(), name, band])
    if curs.fetchall():
        curs.execute("insert into rb (block, end, today, blocker, why, band) values (?, ?, ?, ?, ?, ?)", [
            name, 
            'release',
            now_time, 
            blocker, 
            '', 
            band
        ])
        curs.execute("delete from ban where block = ? and band = ?", [name, band])
    else:
        if login != '':
            login = 'O'
        else:
            login = ''

        if end != '0':
            end = int(number_check(end))

            time = datetime.datetime.now()
            plus = datetime.timedelta(seconds = end)
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

    conn.commit()

def history_plus(title, data, date, ip, send, leng, t_check = ''):
    curs.execute("select id from history where title = ? order by id + 0 desc limit 1", [title])
    id_data = curs.fetchall()

    send = re.sub('\(|\)|<|>', '', send)

    if len(send) > 128:
        send = send[:128]

    if t_check != '':
        send += ' (' + t_check + ')'

    curs.execute("insert into history (id, title, data, date, ip, send, leng, hide) values (?, ?, ?, ?, ?, ?, ?, '')", [
        str(int(id_data[0][0]) + 1) if id_data else '1',
        title,
        data,
        date,
        ip,
        send,
        leng
    ])

def leng_check(first, second):
    if first < second:
        all_plus = '+' + str(second - first)
    elif second < first:
        all_plus = '-' + str(first - second)
    else:
        all_plus = '0'
        
    return all_plus

def number_check(data):
    try:
        int(data)
        return data
    except:
        return '1'

def edit_filter_do(data):
    if admin_check(1) != 1:
        curs.execute("select regex, sub from filter where regex != ''")
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

def redirect(data = '/'):
    return flask.redirect(data)

def re_error(data):
    conn.commit()
    
    if data == '/ban':
        ip = ip_check()

        end = '<li>' + load_lang('why') + ' : ' + load_lang('authority_error') + '</li>'

        if ban_check() == 1:
            end = '<li>' + load_lang('state') + ' : ' + load_lang('ban') + '</li>'
            ok_sign = 1

            band = re.search("^([0-9]{1,3}\.[0-9]{1,3})", ip)
            if band:
                band_it = band.groups()[0]
            else:
                band_it = '-'

            curs.execute("delete from ban where (end < ? and end like '2%')", [get_time()])
            conn.commit()

            curs.execute("select login, block, end from ban where ((end > ? and end like '2%') or end = '') and band = 'regex'", [get_time()])
            regex_d = curs.fetchall()
            for test_r in regex_d:
                g_regex = re.compile(test_r[1])
                if g_regex.search(ip):
                    end += '<li>' + load_lang('type') + ' : regex ban</li>'
                    end += '<li>' + load_lang('end') + ' : ' + (test_r[2] if test_r[2] != '' else load_lang('limitless')) + '</li>'
                    if test_r[0] == 'O':
                        end += '<li>' + load_lang('login_able') + ' (' + load_lang('not_sure') + ')</li>'

                    end += '<hr class=\"main_hr\">'
            
            curs.execute("select login, end from ban where ((end > ? and end like '2%') or end = '') and block = ?", [get_time(), band_it])
            band_d = curs.fetchall()
            if band_d:
                end += '<li>' + load_lang('type') + ' : band ban</li>'
                end += '<li>' + load_lang('end') + ' : ' + (band_d[0][1] if band_d[0][1] != '' else load_lang('limitless')) + '</li>'
                if band_d[0][0] == 'O':
                    end += '<li>' + load_lang('login_able') + ' (' + load_lang('not_sure') + ')</li>'

                end += '<hr class=\"main_hr\">'

            curs.execute("select login, end from ban where ((end > ? and end like '2%') or end = '') and block = ?", [get_time(), ip])
            ban_d = curs.fetchall()
            if ban_d:
                end += '<li>' + load_lang('type') + ' : ban</li>'
                end += '<li>' + load_lang('end') + ' : ' + (ban_d[0][1] if ban_d[0][1] != '' else load_lang('limitless')) + '</li>'
                if ban_d[0][0] == 'O':
                    end += '<li>' + load_lang('login_able') + ' (' + load_lang('not_sure') + ')</li>'

                end += '<hr class=\"main_hr\">'

        return easy_minify(flask.render_template(skin_check(), 
            imp = [load_lang('error'), wiki_set(1), custom(), other2([0, 0])],
            data = '<h2>' + load_lang('error') + '</h2><ul>' + end + '</ul>',
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
            elif num == 11:
                data = load_lang('topic_long_error')
            elif num == 12:
                data = load_lang('email_error')
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
            elif num == 23:
                data = load_lang('regex_error')
            else:
                data = '???'

            return easy_minify(flask.render_template(skin_check(), 
                imp = [load_lang('error'), wiki_set(1), custom(), other2([0, 0])],
                data = '<h2>' + load_lang('error') + '</h2><ul><li>' + data + '</li></ul>',
                menu = 0
            )), 401
        else:
            return redirect('/')
