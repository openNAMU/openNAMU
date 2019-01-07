import werkzeug.routing
import flask_compress
import flask_reggie
import tornado.ioloop
import tornado.httpserver
import tornado.wsgi
import urllib.request
import platform
import zipfile
import difflib
import shutil
import threading
import logging
import random

from func import *

r_ver = 'v3.0.8-master-100'
c_ver = ''.join(re.findall('[0-9]', r_ver))

print('version : ' + r_ver)

try:
    json_data = open('set.json').read()
    set_data = json.loads(json_data)
except:
    while 1:
        print('db name : ', end = '')
        
        new_json = str(input())
        if new_json != '':
            with open('set.json', 'w') as f:
                f.write('{ "db" : "' + new_json + '" }')
            
            json_data = open('set.json').read()
            set_data = json.loads(json_data)

            break
        else:
            print('insert value')
            
            pass

if os.path.exists(set_data['db'] + '.db'):
    setup_tool = 0
else:
    setup_tool = 1

conn = sqlite3.connect(set_data['db'] + '.db', check_same_thread = False)
curs = conn.cursor()

load_conn(conn)

logging.basicConfig(level = logging.ERROR)

app = flask.Flask(__name__, template_folder = './')
app.config['JSON_AS_ASCII'] = False

flask_reggie.Reggie(app)

compress = flask_compress.Compress()
# compress.init_app(app)

class EverythingConverter(werkzeug.routing.PathConverter):
    regex = '.*?'

app.jinja_env.filters['md5_replace'] = md5_replace
app.jinja_env.filters['load_lang'] = load_lang

app.url_map.converters['everything'] = EverythingConverter

curs.execute('create table if not exists data(test text)')
curs.execute('create table if not exists cache_data(test text)')
curs.execute('create table if not exists history(test text)')
curs.execute('create table if not exists rd(test text)')
curs.execute('create table if not exists user(test text)')
curs.execute('create table if not exists user_set(test text)')
curs.execute('create table if not exists ban(test text)')
curs.execute('create table if not exists topic(test text)')
curs.execute('create table if not exists rb(test text)')
curs.execute('create table if not exists back(test text)')
curs.execute('create table if not exists custom(test text)')
curs.execute('create table if not exists other(test text)')
curs.execute('create table if not exists alist(test text)')
curs.execute('create table if not exists re_admin(test text)')
curs.execute('create table if not exists alarm(test text)')
curs.execute('create table if not exists ua_d(test text)')
curs.execute('create table if not exists filter(test text)')
curs.execute('create table if not exists scan(test text)')
curs.execute('create table if not exists acl(test text)')
curs.execute('create table if not exists inter(test text)')
curs.execute('create table if not exists html_filter(test text)')
curs.execute('create table if not exists oauth_conn(test text)')

if setup_tool == 0:
    curs.execute('select data from other where name = "ver"')
    ver_set_data = curs.fetchall()
    if not ver_set_data:
        setup_tool = 1
    else:
        if c_ver > ver_set_data[0][0]:
            setup_tool = 1

if setup_tool != 0:
    create_data = {}

    create_data['all_data'] = [
        'data', 
        'cache_data', 
        'history', 
        'rd', 
        'user',
        'user_set',
        'ban', 
        'topic', 
        'rb', 
        'back', 
        'custom', 
        'other', 
        'alist', 
        're_admin', 
        'alarm', 
        'ua_d', 
        'filter', 
        'scan', 
        'acl', 
        'inter', 
        'html_filter',
        'oauth_conn'
    ]

    create_data['data'] = ['title', 'data']
    create_data['cache_data'] = ['title', 'data']
    create_data['history'] = ['id', 'title', 'data', 'date', 'ip', 'send', 'leng', 'hide', 'type']
    create_data['rd'] = ['title', 'sub', 'date', 'band', 'stop', 'agree']
    create_data['user'] = ['id', 'pw', 'acl', 'date', 'encode']
    create_data['user_set'] = ['name', 'id', 'data']
    create_data['ban'] = ['block', 'end', 'why', 'band', 'login']
    create_data['topic'] = ['id', 'title', 'sub', 'data', 'date', 'ip', 'block', 'top']
    create_data['rb'] = ['block', 'end', 'today', 'blocker', 'why', 'band']
    create_data['back'] = ['title', 'link', 'type']
    create_data['custom'] = ['user', 'css']
    create_data['other'] = ['name', 'data']
    create_data['alist'] = ['name', 'acl']
    create_data['re_admin'] = ['who', 'what', 'time']
    create_data['alarm'] = ['name', 'data', 'date']
    create_data['ua_d'] = ['name', 'ip', 'ua', 'today', 'sub']
    create_data['filter'] = ['name', 'regex', 'sub']
    create_data['scan'] = ['user', 'title']
    create_data['acl'] = ['title', 'dec', 'dis', 'view', 'why']
    create_data['inter'] = ['title', 'link']
    create_data['html_filter'] = ['html', 'kind']
    create_data['oauth_conn'] = ['provider', 'wiki_id', 'sns_id', 'name', 'picture']

    for create_table in create_data['all_data']:
        for create in create_data[create_table]:
            try:
                curs.execute('select ' + create + ' from ' + create_table + ' limit 1')
            except:
                curs.execute("alter table " + create_table + " add " + create + " text default ''")

    update()

curs.execute('select name from alist where acl = "owner"')
if not curs.fetchall():
    curs.execute('delete from alist where name = "owner"')
    curs.execute('insert into alist (name, acl) values ("owner", "owner")')

if not os.path.exists('image'):
    os.makedirs('image')
    
if not os.path.exists('views'):
    os.makedirs('views')

if os.getenv('NAMU_HOST') != None:
    rep_host = os.getenv('NAMU_HOST')
else:
    curs.execute('select data from other where name = "host"')
    rep_data = curs.fetchall()
    if not rep_data:
        while 1:
            print('host [0.0.0.0] : ', end = '')
            rep_host = input()
            if rep_host:
                curs.execute('insert into other (name, data) values ("host", ?)', [rep_host])
                break
            else:
                pass
    else:
        rep_host = rep_data[0][0]
    
        print('host : ' + str(rep_host))

if os.getenv('NAMU_PORT') != None:
    rep_port = os.getenv('NAMU_PORT')
else:
    curs.execute('select data from other where name = "port"')
    rep_data = curs.fetchall()
    if not rep_data:
        while 1:
            print('port : ', end = '')
            rep_port = int(input())
            if rep_port:
                curs.execute('insert into other (name, data) values ("port", ?)', [rep_port])
                break
            else:
                pass
    else:
        rep_port = rep_data[0][0]
    
        print('port : ' + str(rep_port))

try:
    if not os.path.exists('robots.txt'):
        curs.execute('select data from other where name = "robot"')
        robot_test = curs.fetchall()
        if robot_test:
            fw_test = open('./robots.txt', 'w')
            fw_test.write(re.sub('\r\n', '\n', robot_test[0][0]))
            fw_test.close()
        else:
            fw_test = open('./robots.txt', 'w')
            fw_test.write('User-agent: *\nDisallow: /\nAllow: /$\nAllow: /w/')
            fw_test.close()

            curs.execute('insert into other (name, data) values ("robot", "User-agent: *\nDisallow: /\nAllow: /$\nAllow: /w/")')
        
        print('robots.txt have created')
except:
    pass

curs.execute('select data from other where name = "key"')
rep_data = curs.fetchall()
if not rep_data:
    rep_key = ''.join(random.choice("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") for i in range(16))
    if rep_key:
        curs.execute('insert into other (name, data) values ("key", ?)', [rep_key])
else:
    rep_key = rep_data[0][0]

support_language = ['ko-KR', 'en-US']

curs.execute("select data from other where name = 'language'")
rep_data = curs.fetchall()
if not rep_data:
    if os.getenv('NAMU_LANG') != None:
        if os.getenv('NAMU_LANG') in support_language:
            curs.execute("insert into other (name, data) values ('language', ?)", [os.getenv('NAMU_LANG')])
            rep_language = os.getenv('NAMU_LANG')
        else:
            print('language ' + str(os.getenv('NAMU_LANG')) + ' is not supported!')
            rep_language = 'en-US'
    else:
        while 1:
            print('language [' + ', '.join(support_language) + '] : ', end = '')
            rep_language = str(input())
            if rep_language in support_language:
                curs.execute("insert into other (name, data) values ('language', ?)", [rep_language])
                break
            else:
                pass
else:
    rep_language = rep_data[0][0]
    
    print('language : ' + str(rep_language))

curs.execute('select data from other where name = "adsense"')
adsense_result = curs.fetchall()
if not adsense_result:
    curs.execute('insert into other (name, data) values ("adsense", "False")')
    curs.execute('insert into other (name, data) values ("adsense_code", "")')

ask_this = [[['markup', 'markup'], ['namumark']], [['encryption method', 'encode'], ['sha256', 'sha3', 'bcrypt']]]
for ask_data in ask_this:
    curs.execute('select data from other where name = ?', [ask_data[0][1]])
    rep_data = curs.fetchall()
    if not rep_data:
        while 1:
            print(ask_data[0][0] + ' [' + ', '.join(ask_data[1]) + '] : ', end = '')
        
            rep_mark = str(input())
            if rep_mark and rep_mark in ask_data[1]:
                curs.execute('insert into other (name, data) values (?, ?)', [ask_data[0][1], rep_mark])

                break
            else:
                pass
    else:
        rep_mark = rep_data[0][0]

        print(ask_data[0][1] + ' : ' + str(rep_mark))

curs.execute('delete from other where name = "ver"')
curs.execute('insert into other (name, data) values ("ver", ?)', [c_ver])

def back_up():
    try:
        shutil.copyfile(set_data['db'] + '.db', 'back_' + set_data['db'] + '.db')
        
        print('back up : ok')
    except:
        print('back up : error')

    threading.Timer(60 * 60 * back_time, back_up).start()

try:
    curs.execute('select data from other where name = "back_up"')
    back_up_time = curs.fetchall()
    
    back_time = int(back_up_time[0][0])
except:
    back_time = 0
    
if back_time != 0:
    print('back up state : ' + str(back_time) + ' hours interval')
    
    if __name__ == '__main__':
        back_up()
else:
    print('back up state : turn off')

print('\n------ daemon started ------')
conn.commit()

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
        data = '<a href="/del_alarm">(' + load_lang('delete') + ')</a><hr class=\"main_hr\">' + data

        for data_one in data_list:
            data += '<li>' + data_one[0] + ' (' + data_one[1] + ')</li>'
    else:
        data += '<li>-</li>'
    
    data += '</ul>'

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('alarm'), wiki_set(), custom(), other2([0, 0])],
        data = data,
        menu = [['user', load_lang('user')]]
    ))

@app.route('/<regex("inter_wiki|(?:edit|email|name)_filter"):tools>')
def inter_wiki(tools = None):
    div = ''
    admin = admin_check()

    if tools == 'inter_wiki':
        del_link = 'del_inter_wiki'
        plus_link = 'plus_inter_wiki'
        title = load_lang('interwiki') + ' ' + load_lang('list')
        div = ''

        curs.execute('select title, link from inter')
    elif tools == 'email_filter':
        del_link = 'del_email_filter'
        plus_link = 'plus_email_filter'
        title = 'email ' + load_lang('filter') + ' ' + load_lang('list')
        div =   '''
                <ul>
                    <li>gmail.com</li>
                    <li>naver.com</li>
                    <li>daum.net</li>
                    <li>hanmail.net</li>
                    <li>hanmail2.net</li>
                </ul>
                '''

        curs.execute("select html from html_filter where kind = 'email'")
    elif tools == 'name_filter':
        del_link = 'del_name_filter'
        plus_link = 'plus_name_filter'
        title = load_lang('id') + ' ' + load_lang('filter') + ' ' + load_lang('list')
        div = ''

        curs.execute("select html from html_filter where kind = 'name'")
    else:
        del_link = 'del_edit_filter'
        plus_link = 'manager/9'
        title = load_lang('edit') + ' ' + load_lang('filter') + ' ' + load_lang('list')
        div = ''

        curs.execute("select name from filter")

    db_data = curs.fetchall()
    if db_data:
        div += '<ul>'

        for data in db_data:
            if tools == 'inter_wiki':
                div += '<li>' + data[0] + ' : <a id="out_link" href="' + data[1] + '">' + data[1] + '</a>'
            elif tools == 'edit_filter':
                div += '<li><a href="/plus_edit_filter/' + url_pas(data[0]) + '">' + data[0] + '</a>'
            else:
                div += '<li>' + data[0]

            if admin == 1:
                div += ' <a href="/' + del_link + '/' + url_pas(data[0]) + '">(' + load_lang('delete') + ')</a>'

            div += '</li>'

        div += '</ul>'

        if admin == 1:
            div += '<hr class=\"main_hr\"><a href="/' + plus_link + '">(' + load_lang('plus') + ')</a>'
    else:
        if admin == 1:
            div += '<a href="/' + plus_link + '">(' + load_lang('plus') + ')</a>'

    return easy_minify(flask.render_template(skin_check(), 
        imp = [title, wiki_set(), custom(), other2([0, 0])],
        data = div,
        menu = [['other', load_lang('other')]]
    ))

@app.route('/<regex("del_(?:inter_wiki|(?:edit|email|name)_filter)"):tools>/<name>')
def del_inter(tools = None, name = None):
    if admin_check(None, tools) == 1:
        if tools == 'del_inter_wiki':
            curs.execute("delete from inter where title = ?", [name])
        elif tools == 'del_edit_filter':
            curs.execute("delete from filter where name = ?", [name])
        elif tools == 'del_name_filter':
            curs.execute("delete from html_filter where html = ? and kind = 'name'", [name])
        else:
            curs.execute("delete from html_filter where html = ? and kind = 'email'", [name])
        
        conn.commit()

        return redirect('/' + re.sub('^del_', '', tools))
    else:
        return re_error('/error/3')

@app.route('/<regex("plus_(?:inter_wiki|(?:edit|email|name)_filter)"):tools>', methods=['POST', 'GET'])
@app.route('/<regex("plus_edit_filter"):tools>/<name>', methods=['POST', 'GET'])
def plus_inter(tools = None, name = None):
    if flask.request.method == 'POST':
        if tools == 'plus_inter_wiki':
            curs.execute('insert into inter (title, link) values (?, ?)', [flask.request.form.get('title', None), flask.request.form.get('link', None)])
            admin_check(None, 'inter_wiki_plus')
        elif tools == 'plus_edit_filter':
            if admin_check(1, 'edit_filter edit') != 1:
                return re_error('/error/3')

            if flask.request.form.get('limitless', '') != '':
                end = 'X'
            else:
                end = flask.request.form.get('second', 'X')

            curs.execute("select name from filter where name = ?", [name])
            if curs.fetchall():
                curs.execute("update filter set regex = ?, sub = ? where name = ?", [flask.request.form.get('content', 'test'), end, name])
            else:
                curs.execute("insert into filter (name, regex, sub) values (?, ?, ?)", [name, flask.request.form.get('content', 'test'), end])
        else:
            if tools == 'plus_name_filter':
                admin_check(None, 'name_filter edit')
                type_d = 'name'
            else:
                admin_check(None, 'email_filter edit')
                type_d = 'email'
            
            curs.execute('insert into html_filter (html, kind) values (?, ?)', [flask.request.form.get('title', 'test'), type_d])
        
        conn.commit()
    
        return redirect('/' + re.sub('^plus_', '', tools))
    else:
        if admin_check(1) != 1:
            stat = 'disabled'
        else:
            stat = ''

        if tools == 'plus_inter_wiki':
            title = load_lang('interwiki') + ' ' + load_lang('plus')
            form_data = '''
                        <input placeholder="''' + load_lang('name') + '''" type="text" name="title">
                        <hr class=\"main_hr\">
                        <input placeholder="link" type="text" name="link">
                        '''
        elif tools == 'plus_edit_filter':
            curs.execute("select regex, sub from filter where name = ?", [name])
            exist = curs.fetchall()
            if exist:
                textarea = exist[0][0]
                
                if exist[0][1] == 'X':
                    time_check = 'checked="checked"'
                    time_data = ''
                else:
                    time_check = ''
                    time_data = exist[0][1]
            else:
                textarea = ''
                time_check = ''
                time_data = ''

            title = load_lang('edit') + ' ' + load_lang('filter') + ' ' + load_lang('plus')
            form_data = '''
                        <input placeholder="''' + load_lang('second') + '''" name="second" type="text" value="''' + html.escape(time_data) + '''">
                        <hr class=\"main_hr\">
                        <input ''' + stat + ''' type="checkbox" ''' + time_check + ''' name="limitless"> ''' + load_lang('limitless') + '''
                        <hr class=\"main_hr\">
                        <input ''' + stat + ''' placeholder="''' + load_lang('regex') + '''" name="content" value="''' + html.escape(textarea) + '''" type="text">
                        '''
        elif tools == 'plus_name_filter':
            title = load_lang('id') + ' ' + load_lang('filter') + ' ' + load_lang('plus')
            form_data = '<input placeholder="' + load_lang('id') + ' ' + load_lang('regex') + '" type="text" name="title">'
        else:
            title = 'email ' + load_lang('filter') + ' ' + load_lang('plus')
            form_data = '<input placeholder="email" type="text" name="title">'

        return easy_minify(flask.render_template(skin_check(), 
            imp = [title, wiki_set(), custom(), other2([0, 0])],
            data =  '''
                    <form method="post">
                        ''' + form_data + '''
                        <hr class=\"main_hr\">
                        <button ''' + stat + ''' type="submit">''' + load_lang('plus') + '''</button>
                    </form>
                    ''',
            menu = [['other', load_lang('other')], [re.sub('^plus_', '', tools), load_lang('list')]]
        ))

@app.route('/setting')
@app.route('/setting/<int:num>', methods=['POST', 'GET'])
def setting(num = 0):
    if num != 0 and admin_check() != 1:
        return re_error('/ban')

    if num == 0:
        li_list = [load_lang('main'), load_lang('text') + ' ' + load_lang('setting'), load_lang('main') + ' head', load_lang('main') + ' body', 'robots.txt', 'google']
        
        x = 0
        
        li_data = ''
        
        for li in li_list:
            x += 1
            li_data += '<li><a href="/setting/' + str(x) + '">' + li + '</a></li>'

        return easy_minify(flask.render_template(skin_check(), 
            imp = [load_lang('setting'), wiki_set(), custom(), other2([0, 0])],
            data = '<h2>' + load_lang('list') + '</h2><ul>' + li_data + '</ul>',
            menu = [['manager', load_lang('admin')]]
        ))
    elif num == 1:
        i_list = ['name', 'logo', 'frontpage', 'license', 'upload', 'skin', 'edit', 'reg', 'ip_view', 'back_up', 'port', 'key', 'update', 'email_have', 'discussion', 'encode', 'host']
        n_list = ['wiki', '', 'FrontPage', 'CC 0', '2', '', 'normal', '', '', '0', '3000', 'test', 'stable', '', 'normal', 'sha256', '0.0.0.0']
        
        if flask.request.method == 'POST':
            i = 0
            
            for data in i_list:
                curs.execute("update other set data = ? where name = ?", [flask.request.form.get(data, n_list[i]), data])
                i += 1

            conn.commit()

            admin_check(None, 'edit_set')

            return redirect('/setting/1')
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
            acl_list = [[load_lang('subscriber'), 'login'], [load_lang('normal'), 'normal'], [load_lang('admin'), 'admin']]
            for i in acl_list:
                if i[1] == d_list[6]:
                    div = '<option value="' + i[1] + '">' + i[0] + '</option>' + div
                else:
                    div += '<option value="' + i[1] + '">' + i[0] + '</option>'

            div4 = ''
            for i in acl_list:
                if i[1] == d_list[14]:
                    div4 = '<option value="' + i[1] + '">' + i[0] + '</option>' + div4
                else:
                    div4 += '<option value="' + i[1] + '">' + i[0] + '</option>'

            ch_1 = ''
            if d_list[7]:
                ch_1 = 'checked="checked"'

            ch_2 = ''
            if d_list[8]:
                ch_2 = 'checked="checked"'
            
            ch_3 = ''
            if d_list[13]:
                ch_3 = 'checked="checked"'

            div2 = load_skin(d_list[5])

            div3 =''
            if d_list[12] == 'stable':
                div3 += '<option value="stable">stable</option>'
                div3 += '<option value="master">master</option>'
            else:
                div3 += '<option value="master">master</option>'
                div3 += '<option value="stable">stable</option>'
                
            div5 =''
            encode_data = ['sha256', 'sha3', 'bcrypt']
            for i in encode_data:
                if d_list[15] == i:
                    div5 = '<option value="' + i + '">' + i + '</option>' + div5
                else:
                    div5 += '<option value="' + i + '">' + i + '</option>'

            return easy_minify(flask.render_template(skin_check(), 
                imp = [load_lang('main'), wiki_set(), custom(), other2([0, 0])],
                data =  '''
                        <form method="post">
                            <span>''' + load_lang('name') + '''</span>
                            <br>
                            <br>
                            <input placeholder="''' + load_lang('name') + '''" type="text" name="name" value="''' + html.escape(d_list[0]) + '''">
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('logo') + ''' (html)</span>
                            <br>
                            <br>
                            <input placeholder="''' + load_lang('logo') + '''" type="text" name="logo" value="''' + html.escape(d_list[1]) + '''">
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('frontpage') + '''</span>
                            <br>
                            <br>
                            <input placeholder="''' + load_lang('frontpage') + '''" type="text" name="frontpage" value="''' + html.escape(d_list[2]) + '''">
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('bottom') + ' ' + load_lang('text') + ''' (html)</span>
                            <br>
                            <br>
                            <input placeholder="''' + load_lang('bottom') + ' ' + load_lang('text') + '''" type="text" name="license" value="''' + html.escape(d_list[3]) + '''">
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('max_file_size') + ''' [mb]</span>
                            <br>
                            <br>
                            <input placeholder="''' + load_lang('max_file_size') + '''" type="text" name="upload" value="''' + html.escape(d_list[4]) + '''">
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('backup_interval') + ' [' + load_lang('hour') + '''] (off : 0) {restart}</span>
                            <br>
                            <br>
                            <input placeholder="''' + load_lang('backup_interval') + '''" type="text" name="back_up" value="''' + html.escape(d_list[9]) + '''">
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('skin') + '''</span>
                            <br>
                            <br>
                            <select name="skin">''' + div2 + '''</select>
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('default') + ''' acl</span>
                            <br>
                            <br>
                            <select name="edit">''' + div + '''</select>
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('default') + ' ' + load_lang('discussion') + ''' acl</span>
                            <br>
                            <br>
                            <select name="discussion">''' + div4 + '''</select>
                            <hr class=\"main_hr\">
                            <input type="checkbox" name="reg" ''' + ch_1 + '''> ''' + load_lang('register') + ''' X
                            <hr class=\"main_hr\">
                            <input type="checkbox" name="ip_view" ''' + ch_2 + '''> ip ''' + load_lang('hide') + '''
                            <hr class=\"main_hr\">
                            <input type="checkbox" name="email_have" ''' + ch_3 + '''> must have email {<a href="/setting/5">must set google imap</a>}
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('host') + '''</span>
                            <br>
                            <br>
                            <input placeholder="''' + load_lang('host') + '''" type="text" name="host" value="''' + html.escape(d_list[16]) + '''">
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('port') + '''</span>
                            <br>
                            <br>
                            <input placeholder="''' + load_lang('port') + '''" type="text" name="port" value="''' + html.escape(d_list[10]) + '''">
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('secret_key') + '''</span>
                            <br>
                            <br>
                            <input placeholder="''' + load_lang('secret_key') + '''" type="password" name="key" value="''' + html.escape(d_list[11]) + '''">
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('update_branch') + '''</span>
                            <br>
                            <br>
                            <select name="update">''' + div3 + '''</select>
                            <hr class=\"main_hr\">
                            <span>encryption method</span>
                            <br>
                            <br>
                            <select name="encode">''' + div5 + '''</select>
                            <hr class=\"main_hr\">
                            <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        </form>
                        ''',
                menu = [['setting', load_lang('setting')]]
            ))
    elif num == 2:
        if flask.request.method == 'POST':
            curs.execute("update other set data = ? where name = ?", [flask.request.form.get('contract', None), 'contract'])
            curs.execute("update other set data = ? where name = ?", [flask.request.form.get('no_login_warring', None), 'no_login_warring'])
            conn.commit()
            
            admin_check(None, 'edit_set')

            return redirect('/setting/2')
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

            return easy_minify(flask.render_template(skin_check(), 
                imp = [load_lang('text') + ' ' + load_lang('setting'), wiki_set(), custom(), other2([0, 0])],
                data =  '''
                        <form method="post">
                            <span>''' + load_lang('register_text') + '''</span>
                            <br>
                            <br>
                            <input placeholder="''' + load_lang('register_text') + '''" type="text" name="contract" value="''' + html.escape(d_list[0]) + '''">
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('non_login_alert') + '''</span>
                            <br>
                            <br>
                            <input placeholder="''' + load_lang('non_login_alert') + '''" type="text" name="no_login_warring" value="''' + html.escape(d_list[1]) + '''">
                            <hr class=\"main_hr\">
                            <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        </form>
                        ''',
                menu = [['setting', load_lang('setting')]]
            ))
    elif num == 3 or num == 4:
        if flask.request.method == 'POST':
            if num == 4:
                info_d = 'body'
                end_r = '4'
            else:
                info_d = 'head'
                end_r = '3'
            
            curs.execute("select name from other where name = ?", [info_d])
            if curs.fetchall():
                curs.execute("update other set data = ? where name = ?", [flask.request.form.get('content', ''), info_d])
            else:
                curs.execute("insert into other (name, data) values (?, ?)", [info_d, flask.request.form.get('content', '')])
            
            conn.commit()

            admin_check(None, 'edit_set')

            return redirect('/setting/' + end_r)
        else:
            if num == 4:
                curs.execute("select data from other where name = 'body'")
                title = 'body'
            else:
                curs.execute("select data from other where name = 'head'")
                title = 'head'
                
            head = curs.fetchall()
            if head:
                data = head[0][0]
            else:
                data = ''

            return easy_minify(flask.render_template(skin_check(), 
                imp = [load_lang('main') + ' ' + title, wiki_set(), custom(), other2([0, 0])],
                data =  '''
                        <form method="post">
                            <textarea rows="25" name="content">''' + html.escape(data) + '''</textarea>
                            <hr class=\"main_hr\">
                            <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        </form>
                        ''',
                menu = [['setting', load_lang('setting')]]
            ))
    elif num == 5:
        if flask.request.method == 'POST':
            curs.execute("select name from other where name = 'robot'")
            if curs.fetchall():
                curs.execute("update other set data = ? where name = 'robot'", [flask.request.form.get('content', None)])
            else:
                curs.execute("insert into other (name, data) values ('robot', ?)", [flask.request.form.get('content', None)])
            
            conn.commit()
            
            fw = open('./robots.txt', 'w')
            fw.write(re.sub('\r\n', '\n', flask.request.form.get('content', None)))
            fw.close()
            
            admin_check(None, 'edit_set')

            return redirect('/setting/4')
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

            return easy_minify(flask.render_template(skin_check(), 
                imp = ['robots.txt', wiki_set(), custom(), other2([0, 0])],
                data =  '''
                        <a href="/robots.txt">(view)</a>
                        <hr class=\"main_hr\">
                        <form method="post">
                            <textarea rows="25" name="content">''' + html.escape(data) + '''</textarea>
                            <hr class=\"main_hr\">
                            <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        </form>
                        ''',
                menu = [['setting', load_lang('setting')]]
            ))
    elif num == 6:
        i_list = ['recaptcha', 'sec_re', 'g_email', 'g_pass']

        if flask.request.method == 'POST':
            for data in i_list:
                curs.execute("update other set data = ? where name = ?", [flask.request.form.get(data, ''), data])

            conn.commit()
            
            admin_check(None, 'edit_set')

            return redirect('/setting/5')
        else:
            n_list = ['', '', '', '']
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

            return easy_minify(flask.render_template(skin_check(), 
                imp = ['google', wiki_set(), custom(), other2([0, 0])],
                data =  '''
                        <form method="post">
                            <h2><a href="https://www.google.com/recaptcha/admin">recaptcha</a></h2>
                            <span>recaptcha (html)</span>
                            <br>
                            <br>
                            <input placeholder="recaptcha (html)" type="text" name="recaptcha" value="''' + html.escape(d_list[0]) + '''">
                            <hr class=\"main_hr\">
                            <span>recaptcha (secret key)</span>
                            <br>
                            <br>
                            <input placeholder="recaptcha (secret key)" type="text" name="sec_re" value="''' + html.escape(d_list[1]) + '''">
                            <hr class=\"main_hr\">
                            <h2><a href="https://support.google.com/mail/answer/7126229">google imap</a> {restart}</h1>
                            <span>google email</span>
                            <br>
                            <br>
                            <input placeholder="google email" type="text" name="g_email" value="''' + html.escape(d_list[2]) + '''">
                            <hr class=\"main_hr\">
                            <span><a href="https://security.google.com/settings/security/apppasswords">google app password</a></span>
                            <br>
                            <br>
                            <input placeholder="google app password" type="password" name="g_pass" value="''' + html.escape(d_list[3]) + '''">
                            <hr class=\"main_hr\">
                            <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        </form>
                        ''',
                menu = [['setting', load_lang('setting')]]
            ))
    else:
        return redirect('/')

@app.route('/not_close_topic')
def not_close_topic():
    div = '<ul>'
    
    curs.execute('select title, sub from rd where stop != "O" order by date desc')
    n_list = curs.fetchall()
    for data in n_list:
        div += '<li><a href="/topic/' + url_pas(data[0]) + '/sub/' + url_pas(data[1]) + '">' + html.escape(data[0]) + ' (' + data[1] + ')</a></li>'
            
    div += '</ul>'

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('open') + ' ' + load_lang('discussion') + ' ' + load_lang('list'), wiki_set(), custom(), other2([0, 0])],
        data = div,
        menu = [['manager', load_lang('admin')]]
    ))

@app.route('/image/<name>')
def image_view(name = None):
    if os.path.exists(os.path.join('image', name)):
        return flask.send_from_directory('./image', name)
    else:
        return redirect('/')

@app.route('/acl_list')
def acl_list():
    div =   '''
            <table id="main_table_set">
                <tbody>
                    <tr>
                        <td id="main_table_width_quarter">''' + load_lang('document') + ' ' + load_lang('name') + '''</td>
                        <td id="main_table_width_quarter">''' + load_lang('document') + ''' acl</td>
                        <td id="main_table_width_quarter">''' + load_lang('discussion') + ''' acl</td>
                        <td id="main_table_width_quarter">''' + load_lang('view') + ''' acl</td>
            '''
    
    curs.execute("select title, dec, dis, view, why from acl where dec = 'admin' or dec = 'user' or dis = 'admin' or dis = 'user' or view = 'admin' or view = 'user' order by title desc")
    list_data = curs.fetchall()
    for data in list_data:
        if not re.search('^user:', data[0]) and not re.search('^file:', data[0]):
            acl = []
            for i in range(1, 4):
                if data[i] == 'admin':
                    acl += [load_lang('admin')]
                else:
                    acl += [load_lang('subscriber')]

            div +=  '''
                    <tr>
                        <td>
                            <a href="/w/''' + url_pas(data[0]) + '">' + data[0] + '''</a>
                        </td>
                        <td>''' + acl[0] + '''</td>
                        <td>''' + acl[1] + '''</td>
                        <td>''' + acl[2] + '''</td>
                    </tr>
                    '''
        
    div +=  '''
                </tbody>
            </table>
            '''
    
    return easy_minify(flask.render_template(skin_check(), 
        imp = ['acl ' + load_lang('document') + ' ' + load_lang('list'), wiki_set(), custom(), other2([0, 0])],
        data = div,
        menu = [['other', load_lang('other')]]
    ))

@app.route('/admin_plus/<name>', methods=['POST', 'GET'])
def admin_plus(name = None):
    if flask.request.method == 'POST':
        if admin_check(None, 'admin_plus (' + name + ')') != 1:
            return re_error('/error/3')

        curs.execute("delete from alist where name = ?", [name])
        
        if flask.request.form.get('ban', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'ban')", [name])

        if flask.request.form.get('toron', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'toron')", [name])
            
        if flask.request.form.get('check', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'check')", [name])

        if flask.request.form.get('acl', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'acl')", [name])

        if flask.request.form.get('hidel', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'hidel')", [name])

        if flask.request.form.get('give', 0) != 0:
            curs.execute("insert into alist (name, acl) values (?, 'give')", [name])

        if flask.request.form.get('owner', 0) != 0:
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

        if admin_check() != 1:
            state = 'disabled'
        else:
            state = ''

        data += '''
                    <li><input type="checkbox" ''' + state +  ' name="ban" ' + exist_list[0] + '> ' + load_lang('ban') + '''</li>
                    <li><input type="checkbox" ''' + state +  ' name="toron" ' + exist_list[2] + '> ' + load_lang('discussion') + '''</li>
                    <li><input type="checkbox" ''' + state +  ' name="check" ' + exist_list[3] + '> ' + load_lang('user') + ' ' + load_lang('check') + '''</li>
                    <li><input type="checkbox" ''' + state +  ' name="acl" ' + exist_list[4] + '> ' + load_lang('document') + ''' acl</li>
                    <li><input type="checkbox" ''' + state +  ' name="hidel" ' + exist_list[5] + '> ' + load_lang('history') + ' ' + load_lang('hide') + '''</li>
                    <li><input type="checkbox" ''' + state +  ' name="give" ' + exist_list[6] + '> ' + load_lang('authority') + '''</li>
                    <li><input type="checkbox" ''' + state +  ' name="owner" ' + exist_list[7] + '> ' + load_lang('owner') + '''</li>
                </ul>
                '''

        return easy_minify(flask.render_template(skin_check(), 
            imp = [load_lang('admin_group') + ' ' + load_lang('plus'), wiki_set(), custom(), other2([0, 0])],
            data =  '''
                    <form method="post">
                        ''' + data + '''
                        <hr class=\"main_hr\">
                        <button id="save" ''' + state +  ''' type="submit">''' + load_lang('save') + '''</button>
                    </form>
                    ''',
            menu = [['manager', load_lang('admin')]]
        ))        
        
@app.route('/admin_list')
def admin_list():
    div = '<ul>'
    
    curs.execute("select id, acl, date from user where not acl = 'user' order by date desc")
    for data in curs.fetchall():
        name = ip_pas(data[0]) + ' <a href="/admin_plus/' + url_pas(data[1]) + '">(' + data[1] + ')</a>'
        
        if data[2] != '':
            name += '(' + data[2] + ')'

        div += '<li>' + name + '</li>'
        
    div += '</ul>'
                
    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('admin') + ' ' + load_lang('list'), wiki_set(), custom(), other2([0, 0])],
        data = div,
        menu = [['other', load_lang('other')]]
    ))
        
@app.route('/hidden/<everything:name>')
def history_hidden(name = None):
    num = int(flask.request.args.get('num', 0))

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
    num = int(flask.request.args.get('num', 1))
    if num * 50 > 0:
        sql_num = num * 50 - 50
    else:
        sql_num = 0
        
    list_data = '<ul>'

    admin_one = admin_check(1)
    
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
            list_data += ' (' + data[1] + ')'

        list_data += '</li>'

    if num == 1:
        curs.execute("select count(id) from user")
        user_count = curs.fetchall()
        if user_count:
            count = user_count[0][0]
        else:
            count = 0

        list_data +=    '''
                        </ul>
                        <hr class=\"main_hr\">
                        <ul>
                            <li>all : ''' + str(count) + '''</li>
                        </ul>
                        '''

    list_data += next_fix('/user_log?num=', num, user_list)

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('recent') + ' ' + load_lang('subscriber'), wiki_set(), custom(), other2([0, 0])],
        data = list_data,
        menu = 0
    ))

@app.route('/admin_log')
def admin_log():
    num = int(flask.request.args.get('num', 1))
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

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('recent') + ' ' + load_lang('authority'), wiki_set(), custom(), other2([0, 0])],
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
    
    list_data += '</ul><hr class=\"main_hr\"><a href="/manager/8">(' + load_lang('create') + ')</a>'

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('admin_group') + ' ' + load_lang('list'), wiki_set(), custom(), other2([0, 0])],
        data = list_data,
        menu = [['other', load_lang('other')]]
    ))

@app.route('/indexing')
def indexing():
    if admin_check(None, 'indexing') != 1:
        return re_error('/error/3')

    curs.execute("select name from sqlite_master where type = 'index'")
    data = curs.fetchall()
    if data:
        for delete_index in data:
            print('delete : ' + delete_index[0])

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
                print('create : index_' + table[0] + '_' + n_cul)

                sql = 'create index index_' + table[0] + '_' + n_cul + ' on ' + table[0] + '(' + n_cul + ')'
                try:
                    curs.execute(sql)
                except:
                    pass

    conn.commit()
    
    return redirect('/')     

@app.route('/restart', methods=['POST', 'GET'])
def restart():
    if admin_check(None, 'restart') != 1:
        return re_error('/error/3')

    if flask.request.method == 'POST':
        os.execl(sys.executable, sys.executable, *sys.argv)
    else:
        print('restart')

        return easy_minify(flask.render_template(skin_check(), 
            imp = [load_lang('server') + ' ' + load_lang('restart'), wiki_set(), custom(), other2([0, 0])],
            data =  '''
                    <form method="post">
                        <button type="submit">''' + load_lang('restart') + '''</button>
                    </form>
                    ''',
            menu = [['manager', load_lang('admin')]]
        ))       

@app.route('/update')
def now_update():
    if admin_check(None, 'update') != 1:
       return re_error('/error/3')

    curs.execute('select data from other where name = "update"')
    up_data = curs.fetchall()
    if up_data:
        up_data = up_data[0][0]
    else:
        up_data = 'stable'

    if platform.system() == 'Linux':
        print('update')

        os.system('git remote rm origin')
        os.system('git remote add origin https://github.com/2DU/opennamu.git')
        ok = os.system('git fetch origin ' + up_data)
        ok = os.system('git reset --hard origin/' + up_data)
        if ok == 0:
            return redirect('/restart')
    else:
        if platform.system() == 'Windows':
            print('download')

            urllib.request.urlretrieve('https://github.com/2DU/opennamu/archive/' + up_data + '.zip', 'update.zip')

            print('zip extract')
            zipfile.ZipFile('update.zip').extractall('')

            print('move')
            ok = os.system('xcopy /y /r opennamu-' + up_data + ' .')
            if ok == 0:
                print('remove')
                os.system('rd /s /q opennamu-' + up_data)
                os.system('del update.zip')

                return redirect('/restart')

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('update'), wiki_set(), custom(), other2([0, 0])],
        data = 'auto update is not support. <a href="https://github.com/2DU/opennamu">(github)</a>',
        menu = [['manager/1', load_lang('admin')]]
    ))

@app.route('/oauth_settings', methods=['GET', 'POST'])
def oauth_settings():
    if admin_check(None, 'oauth_settings') != 1:
        return re_error('/error/3')

    if flask.request.method == 'POST':
        try:
            facebook_client_id = flask.request.form['facebook_client_id']
            facebook_client_secret = flask.request.form['facebook_client_secret']
            naver_client_id = flask.request.form['naver_client_id']
            naver_client_secret = flask.request.form['naver_client_secret']
        except:
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('inter_error'), wiki_set(), custom(), other2([0, 0])],
                data =  '''
                        <p>''' + load_lang('inter_error_detail') + '''</p>
                        <hr>
                        <code>ie_no_data_required</code>
                        <p>''' + load_lang('ie_no_data_required') + '''</p>
                        ''',
                menu = [['other', load_lang('other')]]
            ))
        with open('oauthsettings.json', 'r', encoding='utf-8') as f:
            legacy = json.loads(f.read())
        with open('oauthsettings.json', 'w', encoding='utf-8') as f:
            f.write(
"""{
    "_README" : {
        "en" : \"""" + legacy['_README']['en'] + """\",
        "ko" : \"""" + legacy['_README']['ko'] + """\",
        "support" : """ + str(legacy['_README']['support']).replace("'", '"') + """
    },
    "publish_url" : \"""" + legacy['publish_url'] + """\",
    "facebook" : {
        "client_id" : \"""" + facebook_client_id + """\",
        "client_secret" : \"""" + facebook_client_secret + """\"
    },
    "naver" : {
        "client_id" : \"""" + naver_client_id + """\",
        "client_secret" : \"""" + naver_client_secret + """\"
    }
}"""
            )
        return flask.redirect('/oauth_settings')

    oauth_supported = load_oauth('_README')['support']

    body_content = ''
    body_content += '''
    <script>function check_value (target) {
        target_box = document.getElementById(target.id + "_box");
        if (target.value !== "") {
            target_box.checked = true;
        } else {
            target_box.checked = false;
        } }
    </script>'''

    init_js = ''
    body_content += '<form method="post">'

    for i in range(len(oauth_supported)):
        oauth_data = load_oauth(oauth_supported[i])
        for j in range(2):
            if j == 0:
                load_target = 'id'
            elif j == 1:
                load_target = 'secret'

            init_js += 'check_value(document.getElementById("{}_client_{}"));'.format(oauth_supported[i], load_target)

            body_content += '''
                            <input id="{}_client_{}_box" type="checkbox" disabled>
                            <input placeholder="{}_client_{}" id="{}_client_{}" name="{}_client_{}" value="{}" type="text" onChange="check_value(this)" style="width: 80%;">
                            <hr>
                            '''.format(
                                oauth_supported[i],
                                load_target,
                                oauth_supported[i], 
                                load_target, 
                                oauth_supported[i], 
                                load_target, 
                                oauth_supported[i], 
                                load_target, 
                                oauth_data['client_{}'.format(load_target)]
                            )
    
    body_content += '<button id="save" type="submit">' + load_lang('save') + '</button></form>'
    body_content += '<script>' + init_js + '</script>'
    
    return easy_minify(flask.render_template(skin_check(),
        imp = [load_lang('oauth_settings'), wiki_set(), custom(), other2([0, 0])],
        data = body_content,
        menu = [['other', load_lang('other')]]
    ))

@app.route('/adsense_settings', methods=['GET', 'POST'])
def adsense_settings():
    if admin_check(None, 'oauth_settings') != 1:
        return re_error('/error/3')
    
    if flask.request.method == 'POST':
        try:
            adsense_enabled = flask.request.form.get('adsense_enabled')
            adsense_code = flask.request.form['adsense_code']
        except:
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('inter_error'), wiki_set(), custom(), other2([0, 0])],
                data =  '''
                        <p>''' + load_lang('inter_error_detail') + '''</p>
                        <hr>
                        <code>ie_no_data_required</code>
                        <p>''' + load_lang('ie_no_data_required') + '''</p>
                        ''',
                menu = [['other', load_lang('other')]]
            ))
        
        if adsense_enabled == 'on':
            curs.execute('update other set data = "True" where name = "adsense"')
        else:
            curs.execute('update other set data = "False" where name = "adsense"')
        curs.execute('update other set data = ? where name = "adsense_code"', [adsense_code])
        conn.commit()
        return redirect('/adsense_settings')

    body_content = ''

    curs.execute('select data from other where name = "adsense"')
    adsense_enabled = curs.fetchall()[0][0]

    curs.execute('select data from other where name = "adsense_code"')
    adsense_code = curs.fetchall()[0][0]

    template = '''
        <form action="" accept-charset="utf-8" method="post">
            <div class="form-check">
                <label class="form-check-label">
                    <input class="form-check-input" name="adsense_enabled" type="checkbox" %_html:adsense_enabled_%>
                    %_lang:adsense_enabled_%
                </label>
            </div>
            <hr>
            <div class="form-group">
                <textarea class="form-control" id="adsense_code" name="adsense_code" rows="12">%_html:adsense_code_%</textarea>
            </div>
            <button type="submit" value="publish">%_lang:save_%</button>
        </form>
    '''
    
    if adsense_enabled == 'True':
        template = template.replace('%_html:adsense_enabled_%', 'checked')
    else:
        template = template.replace('%_html:adsense_enabled_%', '')
    template = template.replace('%_lang:adsense_enabled_%', load_lang('adsense') + ' ' + load_lang('enable'))
    template = template.replace('%_lang:save_%', load_lang('save'))
    template = template.replace('%_html:adsense_code_%', adsense_code)

    body_content += template

    return easy_minify(flask.render_template(skin_check(),
        imp = [load_lang('adsense') + ' ' + load_lang('setting'), wiki_set(), custom(), other2([0, 0])],
        data = body_content,
        menu = [['other', load_lang('other')]]
    ))
        
@app.route('/xref/<everything:name>')
def xref(name = None):
    num = int(flask.request.args.get('num', 1))
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
            div += ' (' + data[1] + ')'
        
        curs.execute("select title from back where title = ? and type = 'include'", [data[0]])
        db_data = curs.fetchall()
        if db_data:
            div += ' <a id="inside" href="/xref/' + url_pas(data[0]) + '">(' + load_lang('backlink') + ')</a>'

        div += '</li>'
      
    div += '</ul>' + next_fix('/xref/' + url_pas(name) + '?num=', num, data_list)
    
    return easy_minify(flask.render_template(skin_check(), 
        imp = [name, wiki_set(), custom(), other2([' (' + load_lang('backlink') + ')', 0])],
        data = div,
        menu = [['w/' + url_pas(name), load_lang('document')]]
    ))

@app.route('/please')
def please():
    num = int(flask.request.args.get('num', 1))
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
    
    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('need') + ' ' + load_lang('document'), wiki_set(), custom(), other2([0, 0])],
        data = div,
        menu = [['other', load_lang('other')]]
    ))
        
@app.route('/recent_discuss')
def recent_discuss():
    div = ''
    
    if flask.request.args.get('what', 'normal') == 'normal':
        div += '<a href="/recent_discuss?what=close">(' + load_lang('close') + ')</a>'
       
        m_sub = 0
    else:
        div += '<a href="/recent_discuss">(' + load_lang('open') + ')</a>'
        
        m_sub = ' (' + load_lang('close') + ')'

    div +=  '''
            <hr class=\"main_hr\">
            <table id="main_table_set">
                <tbody>
                    <tr>
                        <td id="main_table_width_half">''' + load_lang('discussion') + ' ' + load_lang('name') + '''</td>
                        <td id="main_table_width_half">''' + load_lang('time') + '''</td>
                    </tr>
            '''
    
    if m_sub == 0:
        curs.execute("select title, sub, date from rd where not stop = 'O' order by date desc limit 50")
    else:
        curs.execute("select title, sub, date from rd where stop = 'O' order by date desc limit 50")
        
    for data in curs.fetchall():
        title = html.escape(data[0])
        sub = html.escape(data[1])

        div += '<tr><td><a href="/topic/' + url_pas(data[0]) + '/sub/' + url_pas(data[1]) + '">' + title + '</a> (' + sub + ')</td><td>' + data[2] + '</td></tr>'
    
    div += '</tbody></table>'
            
    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('recent') + ' ' + load_lang('discussion'), wiki_set(), custom(), other2([m_sub, 0])],
        data = div,
        menu = 0
    ))

@app.route('/block_log')
@app.route('/block_log/<regex("ip|user|never_end|can_end|end|now|edit_filter"):tool2>')
@app.route('/<regex("block_user|block_admin"):tool>/<name>')
def block_log(name = None, tool = None, tool2 = None):
    num = int(flask.request.args.get('num', 1))
    if num * 50 > 0:
        sql_num = num * 50 - 50
    else:
        sql_num = 0
    
    div =   '''
            <table id="main_table_set">
                <tbody>
                    <tr>
                        <td id="main_table_width">''' + load_lang('blocked') + '''</td>
                        <td id="main_table_width">''' + load_lang('admin') + '''</td>
                        <td id="main_table_width">''' + load_lang('period') + '''</td>
                    </tr>
            '''
    
    data_list = ''
    
    if not name:
        if not tool2:
            div =   '''
                    <a href="/manager/11">(''' + load_lang('blocked') + ''')</a> <a href="/manager/12">(''' + load_lang('admin') + ''')</a>
                    <hr class=\"main_hr\">
                    <a href="/block_log/ip">(ip)</a> <a href="/block_log/user">(''' + load_lang('subscriber') + ')</a> <a href="/block_log/never_end">(' + load_lang('limitless') + ')</a> <a href="/block_log/can_end">(' + load_lang('period') + ')</a> <a href="/block_log/end">(' + load_lang('release') + ')</a> <a href="/block_log/now">(' + load_lang('now') + ')</a> <a href="/block_log/edit_filter">(' + load_lang('edit') + ' ' + load_lang('filter') + ''')</a>
                    <hr class=\"main_hr\">
                    ''' + div
            
            sub = 0
            menu = 0
            
            curs.execute("select why, block, blocker, end, today from rb order by today desc limit ?, '50'", [str(sql_num)])
        else:
            menu = [['block_log', load_lang('normal')]]
            
            if tool2 == 'ip':
                sub = ' (ip)'
                
                curs.execute("select why, block, blocker, end, today from rb where (block like ? or block like ?) order by today desc limit ?, '50'", ['%.%', '%:%', str(sql_num)])
            elif tool2 == 'user':
                sub = ' (' + load_lang('subscriber') + ')'
                
                curs.execute("select why, block, blocker, end, today from rb where not (block like ? or block like ?) order by today desc limit ?, '50'", ['%.%', '%:%', str(sql_num)])
            elif tool2 == 'never_end':
                sub = '(' + load_lang('limitless') + ')'
                
                curs.execute("select why, block, blocker, end, today from rb where not end like ? and not end like ? order by today desc limit ?, '50'", ['%:%', '%' + load_lang('release', 1) + '%', str(sql_num)])
            elif tool2 == 'end':
                sub = '(' + load_lang('release') + ')'
                
                curs.execute("select why, block, blocker, end, today from rb where end = ? order by today desc limit ?, '50'", [load_lang('release', 1), str(sql_num)])
            elif tool2 == 'now':
                sub = '(' + load_lang('now') + ')'
                
                data_list = []
                
                curs.execute("select block from ban where end > ? limit ?, '50'", [get_time(), str(sql_num)])
                for in_data in curs.fetchall():
                    curs.execute("select why, block, blocker, end, today from rb where block = ? order by today desc limit 1", [in_data[0]])
                    
                    data_list = [curs.fetchall()[0]] + data_list
            elif tool2 == 'edit_filter':
                sub = '(' + load_lang('edit') + ' ' + load_lang('filter') + ')'

                curs.execute("select why, block, blocker, end, today from rb where blocker = ? order by today desc limit ?, '50'", [load_lang('tool', 1) + ':' + load_lang('edit', 1) + ' ' + load_lang('filter', 1), str(sql_num)])
            else:
                sub = '(' + load_lang('period') + ')'
                
                curs.execute("select why, block, blocker, end, today from rb where end like ? order by today desc limit ?, '50'", ['%\-%', str(sql_num)])
    else:
        menu = [['block_log', load_lang('normal')]]
        
        if tool == 'block_user':
            sub = ' (' + load_lang('blocked') + ')'
            
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
            ip = data[1] + ' (' + load_lang('band') + ')'
        else:
            ip = ip_pas(data[1])

        if data[3] != '':
            end = data[3]
        else:
            end = load_lang('limitless') + ''
            
        div +=  '''
                        <tr>
                            <td>''' + ip + '''</td>
                            <td>''' + ip_pas(data[2]) + '''</td>
                            <td>
                                start : ''' + data[4] + '''
                                <br>
                                end : ''' + end + '''
                            </td>
                        </tr>
                        <tr>
                            <td colspan="3">''' + why + '''</td>
                        </tr>
                    </tbody>
                </table>
                '''
    
    if not name:
        if not tool2:
            div += next_fix('/block_log?num=', num, data_list)
        else:
            div += next_fix('/block_log/' + url_pas(tool2) + '?num=', num, data_list)
    else:
        div += next_fix('/' + url_pas(tool) + '/' + url_pas(name) + '?num=', num, data_list)
                
    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('recent') + ' ' + load_lang('ban'), wiki_set(), custom(), other2([sub, 0])],
        data = div,
        menu = menu
    ))
            
@app.route('/search', methods=['POST'])
def search():
    return redirect('/search/' + url_pas(flask.request.form.get('search', None)))

@app.route('/goto', methods=['POST'])
def goto():
    curs.execute("select title from data where title = ?", [flask.request.form.get('search', None)])
    data = curs.fetchall()
    if data:
        return redirect('/w/' + url_pas(flask.request.form.get('search', None)))
    else:
        return redirect('/search/' + url_pas(flask.request.form.get('search', None)))

@app.route('/search/<everything:name>')
def deep_search(name = None):
    num = int(flask.request.args.get('num', 1))
    if num * 50 > 0:
        sql_num = num * 50 - 50
    else:
        sql_num = 0

    div = '<ul>'
    
    div_plus = ''
    test = ''
    
    curs.execute("select title from data where title = ?", [name])
    if curs.fetchall():
        link_id = ''
    else:
        link_id = 'id="not_thing"'
    
    div =   '''
            <ul>
                <li>
                    <a ''' + link_id + ' href="/w/' + url_pas(name) + '">' + name + '''</a>
                </li>
            </ul>
            <hr class=\"main_hr\">
            <ul>
            '''

    curs.execute("select distinct title, case when title like ? then '제목' else '내용' end from data where title like ? or data like ? order by case when title like ? then 1 else 2 end limit ?, '50'", ['%' + name + '%', '%' + name + '%', '%' + name + '%', '%' + name + '%', str(sql_num)])
    all_list = curs.fetchall()
    if all_list:
        test = all_list[0][1]
        
        for data in all_list:
            if data[1] != test:
                div_plus += '</ul><hr class=\"main_hr\"><ul>'
                
                test = data[1]

            div_plus += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a> (' + data[1] + ')</li>'
    else:
        div += '<li>-</li>'

    div += div_plus + '</ul>'
    div += next_fix('/search/' + url_pas(name) + '?num=', num, all_list)

    return easy_minify(flask.render_template(skin_check(), 
        imp = [name, wiki_set(), custom(), other2([' (' + load_lang('search') + ')', 0])],
        data = div,
        menu = 0
    ))
         
@app.route('/raw/<everything:name>')
@app.route('/topic/<everything:name>/sub/<sub_title>/raw/<int:num>')
def raw_view(name = None, sub_title = None, num = None):
    v_name = name
    sub = ' (' + load_lang('raw') + ')'
    
    if not num:
        num = flask.request.args.get('num', None)
        if num:
            num = int(num)
    
    if not sub_title and num:
        curs.execute("select title from history where title = ? and id = ? and hide = 'O'", [name, str(num)])
        if curs.fetchall() and admin_check(6) != 1:
            return re_error('/error/3')
        
        curs.execute("select data from history where title = ? and id = ?", [name, str(num)])
        
        sub += ' (' + str(num) + load_lang('version') + ')'

        menu = [['history/' + url_pas(name), load_lang('history')]]
    elif sub_title:
        curs.execute("select data from topic where id = ? and title = ? and sub = ? and block = ''", [str(num), name, sub_title])
        
        v_name = load_lang('discussion') + ' Raw'
        sub = ' (' + str(num) + ')'

        menu = [['topic/' + url_pas(name) + '/sub/' + url_pas(sub_title) + '#' + str(num), load_lang('discussion')], ['topic/' + url_pas(name) + '/sub/' + url_pas(sub_title) + '/admin/' + str(num), load_lang('tool')]]
    else:
        curs.execute("select data from data where title = ?", [name])
        
        menu = [['w/' + url_pas(name), load_lang('document')]]

    data = curs.fetchall()
    if data:
        p_data = html.escape(data[0][0])
        p_data = '<textarea readonly rows="25">' + p_data + '</textarea>'
        
        return easy_minify(flask.render_template(skin_check(), 
            imp = [v_name, wiki_set(), custom(), other2([sub, 0])],
            data = p_data,
            menu = menu
        ))
    else:
        return redirect('/w/' + url_pas(name))
        
@app.route('/revert/<everything:name>', methods=['POST', 'GET'])
def revert(name = None):    
    num = int(flask.request.args.get('num', 0))

    curs.execute("select title from history where title = ? and id = ? and hide = 'O'", [name, str(num)])
    if curs.fetchall() and admin_check(6) != 1:
        return re_error('/error/3')

    if acl_check(name) == 1:
        return re_error('/ban')

    if flask.request.method == 'POST':
        if captcha_post(flask.request.form.get('g-recaptcha-response', '')) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)
    
        curs.execute("select data from history where title = ? and id = ?", [name, str(num)])
        data = curs.fetchall()
        if data:
            if edit_filter_do(data[0][0]) == 1:
                return re_error('/error/21')

        curs.execute("delete from back where link = ?", [name])
        conn.commit()
        
        if data:                                
            curs.execute("select data from data where title = ?", [name])
            data_old = curs.fetchall()
            if data_old:
                leng = leng_check(len(data_old[0][0]), len(data[0][0]))
                curs.execute("update data set data = ? where title = ?", [data[0][0], name])
            else:
                leng = '+' + str(len(data[0][0]))
                curs.execute("insert into data (title, data) values (?, ?)", [name, data[0][0]])
                
            history_plus(
                name, 
                data[0][0], 
                get_time(), 
                ip_check(), 
                flask.request.form.get('send', None) + ' (' + str(num) + load_lang('version', 1) + ')', 
                leng
            )

            render_set(
                title = name,
                data = data[0][0],
                num = 1
            )
            
            conn.commit()
            
        return redirect('/w/' + url_pas(name))
    else:
        curs.execute("select title from history where title = ? and id = ?", [name, str(num)])
        if not curs.fetchall():
            return redirect('/w/' + url_pas(name))

        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2([' (' + load_lang('revert') + ')', 0])],
            data =  '''
                    <form method="post">
                        <span>''' + flask.request.args.get('num', '0') + load_lang('version') + '''</span>
                        <hr class=\"main_hr\">
                        ''' + ip_warring() + '''
                        <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                        <hr class=\"main_hr\">
                        ''' + captcha_get() + '''
                        <button type="submit">''' + load_lang('revert') + '''</button>
                    </form>
                    ''',
            menu = [['history/' + url_pas(name), load_lang('history')], ['recent_changes', load_lang('recent') + ' ' + load_lang('change')]]
        ))

@app.route('/edit/<everything:name>', methods=['POST', 'GET'])
def edit(name = None):
    ip = ip_check()
    if acl_check(name) == 1:
        return re_error('/ban')
    
    if flask.request.method == 'POST':
        if captcha_post(flask.request.form.get('g-recaptcha-response', '')) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)
            
        if len(flask.request.form.get('send', None)) > 500:
            return re_error('/error/15')

        if flask.request.form.get('otent', None) == flask.request.form.get('content', None):
            return redirect('/w/' + url_pas(name))
            
        if edit_filter_do(flask.request.form.get('content', '')) == 1:
            return re_error('/error/21')

        today = get_time()
        content = savemark(flask.request.form.get('content', None))
        
        curs.execute("select data from data where title = ?", [name])
        old = curs.fetchall()
        if old:
            leng = leng_check(len(flask.request.form.get('otent', None)), len(content))
            
            if flask.request.args.get('section', None):
                content = old[0][0].replace(flask.request.form.get('otent', None), content)
                
            curs.execute("update data set data = ? where title = ?", [content, name])
        else:
            leng = '+' + str(len(content))
            
            curs.execute("insert into data (title, data) values (?, ?)", [name, content])

        curs.execute("select user from scan where title = ?", [name])
        for _ in curs.fetchall():
            curs.execute("insert into alarm (name, data, date) values (?, ?, ?)", [ip, ip + ' - <a href="/w/' + url_pas(name) + '">' + name + '</a> (Edit)', today])

        history_plus(
            name,
            content,
            today,
            ip,
            flask.request.form.get('send', None),
            leng
        )
        
        curs.execute("delete from back where link = ?", [name])
        curs.execute("delete from back where title = ? and type = 'no'", [name])
        
        render_set(
            title = name,
            data = content,
            num = 1
        )
        
        conn.commit()
        
        return redirect('/w/' + url_pas(name))
    else:            
        curs.execute("select data from data where title = ?", [name])
        new = curs.fetchall()
        if new:
            if flask.request.args.get('section', None):
                test_data = '\n' + re.sub('\r\n', '\n', new[0][0]) + '\n'   
                
                section_data = re.findall('((?:={1,6}) ?(?:(?:(?!={1,6}\n).)+) ?={1,6}\n(?:(?:(?!(?:={1,6}) ?(?:(?:(?!={1,6}\n).)+) ?={1,6}\n).)*\n*)*)', test_data)
                data = section_data[int(flask.request.args.get('section', None)) - 1]
            else:
                data = new[0][0]
        else:
            data = ''
            
        data_old = data
        
        if not flask.request.args.get('section', None):
            get_name =  '''
                        <a href="/manager/15?plus=''' + url_pas(name) + '">(' + load_lang('load') + ')</a> <a href="/edit_filter">(' + load_lang('edit') + ' ' + load_lang('filter') + ''')</a>
                        <hr class=\"main_hr\">
                        '''
            action = ''
        else:
            get_name = ''
            action = '?section=' + flask.request.args.get('section', None)
            
        if flask.request.args.get('plus', None):
            curs.execute("select data from data where title = ?", [flask.request.args.get('plus', None)])
            get_data = curs.fetchall()
            if get_data:
                data = get_data[0][0]
                get_name = ''

        js_data = edit_help_button()

        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2([' (' + load_lang('edit') + ')', 0])],
            data =  get_name + js_data[0] + '''
                    <form method="post" action="/edit/''' + url_pas(name) + action + '''">
                        ''' + js_data[1] + '''
                        <textarea id="content" rows="25" name="content">''' + html.escape(re.sub('\n$', '', data)) + '''</textarea>
                        <textarea style="display: none;" name="otent">''' + html.escape(re.sub('\n$', '', data_old)) + '''</textarea>
                        <hr class=\"main_hr\">
                        <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                        <hr class=\"main_hr\">
                        ''' + captcha_get() + ip_warring() + '''
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        <button id="preview" type="submit" formaction="/preview/''' + url_pas(name) + action + '">' + load_lang('preview') + '''</button>
                    </form>
                    ''',
            menu = [['w/' + url_pas(name), load_lang('document')], ['delete/' + url_pas(name), load_lang('delete')], ['move/' + url_pas(name), load_lang('move')]]
        ))

@app.route('/preview/<everything:name>', methods=['POST'])
def preview(name = None):
    if acl_check(name) == 1:
        return re_error('/ban')
         
    new_data = re.sub('^\r\n', '', flask.request.form.get('content', None))
    new_data = re.sub('\r\n$', '', new_data)
    
    end_data = render_set(
        title = name,
        data = new_data
    )
    
    if flask.request.args.get('section', None):
        action = '?section=' + flask.request.args.get('section', None)
    else:
        action = ''

    js_data = edit_help_button()
    
    return easy_minify(flask.render_template(skin_check(), 
        imp = [name, wiki_set(), custom(), other2([' (' + load_lang('preview') + ')', 0])],
        data =  '<a href="/edit_filter">(' + load_lang('edit') + ' ' + load_lang('filter') + ')</a>' + js_data[0] + '''
                <form method="post" action="/edit/''' + url_pas(name) + action + '''">
                    ''' + js_data[1] + '''
                    <textarea id="content" rows="25" name="content">''' + html.escape(flask.request.form.get('content', None)) + '''</textarea>
                    <textarea style="display: none;" name="otent">''' + html.escape(flask.request.form.get('otent', None)) + '''</textarea>
                    <hr class=\"main_hr\">
                    <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                    <hr class=\"main_hr\">
                    ''' + captcha_get() + '''
                    <button id="save" type="submit">''' + load_lang('save') + '''</button>
                    <button id="preview" type="submit" formaction="/preview/''' + url_pas(name) + action + '">' + load_lang('preview') + '''</button>
                </form>
                <hr class=\"main_hr\">
                ''' + end_data,
        menu = [['w/' + url_pas(name), load_lang('document')]]
    ))
        
@app.route('/delete/<everything:name>', methods=['POST', 'GET'])
def delete(name = None):
    ip = ip_check()
    if acl_check(name) == 1:
        return re_error('/ban')
    
    if flask.request.method == 'POST':
        if captcha_post(flask.request.form.get('g-recaptcha-response', '')) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        curs.execute("select data from data where title = ?", [name])
        data = curs.fetchall()
        if data:
            today = get_time()
            leng = '-' + str(len(data[0][0]))
            
            history_plus(
                name, 
                '', 
                today, 
                ip, 
                flask.request.form.get('send', None) + ' (' + load_lang('delete', 1) + ')', 
                leng
            )
            
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

        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2([' (' + load_lang('delete') + ')', 0])],
            data =  '''
                    <form method="post">
                        ''' + ip_warring() + '''
                        <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                        <hr class=\"main_hr\">
                        ''' + captcha_get() + '''
                        <button type="submit">''' + load_lang('delete') + '''</button>
                    </form>
                    ''',
            menu = [['w/' + url_pas(name), load_lang('document')]]
        ))            
            
@app.route('/move_data/<everything:name>')
def move_data(name = None):    
    data = '<ul>'
    
    curs.execute("select send, date, ip from history where send like ? or send like ? order by date desc", ['%<a>' + name + '</a> ' + load_lang('move', 1) + ')%', '%(<a>' + name + '</a>%'])
    for for_data in curs.fetchall():
        match = re.findall('<a>((?:(?!<\/a>).)+)<\/a>', for_data[0])
        send = re.sub('\([^\)]+\)$', '', for_data[0])
        data += '<li><a href="/move_data/' + url_pas(match[0]) + '">' + match[0] + '</a> - <a href="/move_data/' + url_pas(match[1]) + '">' + match[1] + '</a>'
        
        if re.search('^( *)+$', send):
            data += ' / ' + for_data[2] + ' / ' + for_data[1] + '</li>'
        else:
            data += ' / ' + for_data[2] + ' / ' + for_data[1] + ' / ' + send + '</li>'
    
    data += '</ul>'
    
    return easy_minify(flask.render_template(skin_check(), 
        imp = [name, wiki_set(), custom(), other2([' (' + load_lang('move') + ' ' + load_lang('history') + ')', 0])],
        data = data,
        menu = [['history/' + url_pas(name), load_lang('history')]]
    ))        
            
@app.route('/move/<everything:name>', methods=['POST', 'GET'])
def move(name = None):
    if acl_check(name) == 1:
        return re_error('/ban')

    if flask.request.method == 'POST':
        if captcha_post(flask.request.form.get('g-recaptcha-response', '')) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        curs.execute("select title from history where title = ?", [flask.request.form.get('title', None)])
        if curs.fetchall():
            if admin_check(None, 'merge documents') == 1:
                curs.execute("select data from data where title = ?", [flask.request.form.get('title', None)])
                data = curs.fetchall()
                if data:            
                    curs.execute("delete from data where title = ?", [flask.request.form.get('title', None)])
                    curs.execute("delete from back where link = ?", [flask.request.form.get('title', None)])
                
                curs.execute("select data from data where title = ?", [name])
                data = curs.fetchall()
                if data:            
                    curs.execute("update data set title = ? where title = ?", [flask.request.form.get('title', None), name])
                    curs.execute("update back set link = ? where link = ?", [flask.request.form.get('title', None), name])
                    
                    data_in = data[0][0]
                else:
                    data_in = ''

                history_plus(
                    name, 
                    data_in, 
                    get_time(), 
                    ip_check(), 
                    flask.request.form.get('send', None) + ' (marge <a>' + name + '</a> - <a>' + flask.request.form.get('title', None) + '</a> ' + load_lang('move', 1) + ')', 
                    '0'
                )

                curs.execute("update back set type = 'no' where title = ? and not type = 'cat' and not type = 'no'", [name])
                curs.execute("delete from back where title = ? and not type = 'cat' and type = 'no'", [flask.request.form.get('title', None)])

                curs.execute("select id from history where title = ? order by id + 0 desc limit 1", [flask.request.form.get('title', None)])
                data = curs.fetchall()
                
                num = data[0][0]

                curs.execute("select id from history where title = ? order by id + 0 asc", [name])
                data = curs.fetchall()
                for move in data:
                    curs.execute("update history set title = ?, id = ? where title = ? and id = ?", [flask.request.form.get('title', None), str(int(num) + int(move[0])), name, move[0]])

                conn.commit()

                return redirect('/w/' + url_pas(flask.request.form.get('title', None)))
            else:
                return re_error('/error/19')
        else:
            curs.execute("select data from data where title = ?", [name])
            data = curs.fetchall()
            if data:            
                curs.execute("update data set title = ? where title = ?", [flask.request.form.get('title', None), name])
                curs.execute("update back set link = ? where link = ?", [flask.request.form.get('title', None), name])
                
                data_in = data[0][0]
            else:
                data_in = ''
                
            history_plus(
                name, 
                data_in, 
                get_time(), 
                ip_check(), 
                flask.request.form.get('send', None) + ' (<a>' + name + '</a> - <a>' + flask.request.form.get('title', None) + '</a> ' + load_lang('move', 1) + ')', 
                '0'
            )
            
            curs.execute("update back set type = 'no' where title = ? and not type = 'cat' and not type = 'no'", [name])
            curs.execute("delete from back where title = ? and not type = 'cat' and type = 'no'", [flask.request.form.get('title', None)])

            curs.execute("update history set title = ? where title = ?", [flask.request.form.get('title', None), name])
            conn.commit()

            return redirect('/w/' + url_pas(flask.request.form.get('title', None)))
    else:            
        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2([' (' + load_lang('move') + ')', 0])],
            data =  '''
                    <form method="post">
                        ''' + ip_warring() + '''
                        <input placeholder="''' + load_lang('document') + ' ' + load_lang('name') + '" value="' + name + '''" name="title" type="text">
                        <hr class=\"main_hr\">
                        <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                        <hr class=\"main_hr\">
                        ''' + captcha_get() + '''
                        <button type="submit">''' + load_lang('move') + '''</button>
                    </form>
                    ''',
            menu = [['w/' + url_pas(name), load_lang('document')]]
        ))
            
@app.route('/other')
def other():
    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('other') + ' ' + load_lang('tool'), wiki_set(), custom(), other2([0, 0])],
        data =  '''
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
                    <li><a href="/title_index">''' + load_lang('all') + ' ' + load_lang('document') + '''</a></li>
                    <li><a href="/acl_list">acl ''' + load_lang('document') + '''</a></li>
                    <li><a href="/please">''' + load_lang('need') + ' ' + load_lang('document') + '''</a></li>
                </ul>
                <br>
                <h2>''' + load_lang('other') + '''</h2>
                <ul>
                    <li><a href="/upload">''' + load_lang('upload') + '''</a></li>
                    <li><a href="/manager/10">''' + load_lang('document') + ' ' + load_lang('search') + '''</a></li>
                </ul>
                <br>
                <h2>''' + load_lang('admin') + '''</h2>
                <ul>
                    <li><a href="/manager/1">''' + load_lang('admin') + ' ' + load_lang('tool') + '''</a></li>
                </ul>
                <br>
                <h2>''' + load_lang('normal_version') + '''</h2>
                <ul>
                    <li>''' + load_lang('normal_version') + ' : <a id="out_link" href="https://github.com/2DU/opennamu/blob/master/version.md">' + r_ver + '''</a></li>
                </ul>
                ''',
    menu = 0
    ))
    
@app.route('/manager', methods=['POST', 'GET'])
@app.route('/manager/<int:num>', methods=['POST', 'GET'])
def manager(num = 1):
    title_list = {
        0 : [load_lang('document') + ' ' + load_lang('name'), 'acl'], 
        1 : [0, 'check'], 
        2 : [0, 'ban'], 
        3 : [0, 'admin'], 
        4 : [0, 'record'], 
        5 : [0, 'topic_record'], 
        6 : [load_lang('name'), 'admin_plus'], 
        7 : [load_lang('name'), 'plus_edit_filter'], 
        8 : [load_lang('document') + ' ' + load_lang('name'), 'search'], 
        9 : [0, 'block_user'], 
        10 : [0, 'block_admin'], 
        11 : [load_lang('document') + ' ' + load_lang('name'), 'watch_list'], 
        12 : [load_lang('compare'), 'check'], 
        13 : [load_lang('document') + ' ' + load_lang('name'), 'edit']
    }
    
    if num == 1:
        return easy_minify(flask.render_template(skin_check(), 
            imp = [load_lang('admin') + ' ' + load_lang('tool'), wiki_set(), custom(), other2([0, 0])],
            data =  '''
                    <h2>''' + load_lang('admin') + '''</h2>
                    <ul>
                        <li><a href="/manager/2">''' + load_lang('document') + ''' acl</a></li>
                        <li><a href="/manager/3">''' + load_lang('user') + ' ' + load_lang('check') + '''</a></li>
                        <li><a href="/manager/4">''' + load_lang('user') + ' ' + load_lang('ban') + '''</a></li>
                        <li><a href="/manager/5">''' + load_lang('subscriber') + ' ' + load_lang('authority') + '''</a></li>
                        <li><a href="/edit_filter">''' + load_lang('edit') + ' ' + load_lang('filter') + '''</a></li>
                    </ul>
                    <br>
                    <h2>''' + load_lang('owner') + '''</h2>
                    <ul>
                        <li><a href="/manager/8">''' + load_lang('admin_group') + ' ' + load_lang('create') + '''</a></li>
                        <li><a href="/setting">''' + load_lang('setting') + ' ' + load_lang('edit') + '''</a></li>
                    </ul>
                    <h3>''' + load_lang('filter') + '''</h3>
                    <ul>
                        <li><a href="/inter_wiki">''' + load_lang('interwiki') + '''</a></li>
                        <li><a href="/html_filter">html ''' + load_lang('filter') + '''</a></li>
                        <li><a href="/email_filter">email ''' + load_lang('filter') + '''</a></li>
                        <li><a href="/name_filter">''' + load_lang('id') + ' ' + load_lang('filter') + '''</a></li>
                    </ul>
                    <br>
                    <h2>''' + load_lang('server') + '''</h2>
                    <ul>
                        <li><a href="/indexing">''' + load_lang('indexing') + ' (' + load_lang('create') + ' or ' + load_lang('delete') + ''')</a></li>
                        <li><a href="/restart">''' + load_lang('server') + ' ' + load_lang('restart') + '''</a></li>
                        <li><a href="/update">''' + load_lang('update') + '''</a></li>
                        <li><a href="/oauth_settings">''' + load_lang('oauth_settings') + '''</a></li>
                        <li><a href="/adsense_settings">''' + load_lang('adsense') + ' ' + load_lang('setting') + '''</a></li>
                    </ul>
                    ''',
            menu = [['other', load_lang('other')]]
        ))
    elif not num - 1 > len(title_list):
        if flask.request.method == 'POST':
            if flask.request.args.get('plus', None):
                return redirect('/' + title_list[(num - 2)][1] + '/' + url_pas(flask.request.args.get('plus', None)) + '?plus=' + flask.request.form.get('name', None))
            else:
                return redirect('/' + title_list[(num - 2)][1] + '/' + url_pas(flask.request.form.get('name', None)))
        else:
            if title_list[(num - 2)][0] == 0:
                placeholder = load_lang('user') + ' ' + load_lang('name')
            else:
                placeholder = title_list[(num - 2)][0]

            return easy_minify(flask.render_template(skin_check(), 
                imp = ['Redirect', wiki_set(), custom(), other2([0, 0])],
                data =  '''
                        <form method="post">
                            <input placeholder="''' + placeholder + '''" name="name" type="text">
                            <hr class=\"main_hr\">
                            <button type="submit">''' + load_lang('move') + '''</button>
                        </form>
                        ''',
                menu = [['manager', load_lang('admin')]]
            ))
    else:
        return redirect('/')
        
@app.route('/title_index')
def title_index():
    page = int(flask.request.args.get('page', 1))
    num = int(flask.request.args.get('num', 100))
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
        data += '<hr class=\"main_hr\"><ul>'

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

        sql_list = [load_lang('template', 1) + ':', 'category:', 'user:', 'file:']
        for sql in sql_list:
            curs.execute("select count(title) from data where title like ?", [sql + '%'])
            count = curs.fetchall()
            if count:
                count_end += [count[0][0]]
            else:
                count_end += [0]

        count_end += [count_end[0] - count_end[1]  - count_end[2]  - count_end[3]  - count_end[4]]
        
        data += '''
                </ul>
                <hr class=\"main_hr\">
                <ul>
                    <li>all : ''' + str(count_end[0]) + '''</li>
                </ul>
                <hr class=\"main_hr\">
                <ul>
                    <li>''' + load_lang('template') + ' : ' + str(count_end[1]) + '''</li>
                    <li>''' + load_lang('category') + ' : ' + str(count_end[2]) + '''</li>
                    <li>''' + load_lang('user') + ' : ' + str(count_end[3]) + '''</li>
                    <li>''' + load_lang('file') + ' : ' + str(count_end[4]) + '''</li>
                    <li>other : ''' + str(count_end[5]) + '''</li>
                '''

    data += '</ul>' + next_fix('/title_index?num=' + str(num) + '&page=', page, title_list, num)
    sub = ' (' + str(num) + ')'
    
    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('all') + ' ' + load_lang('document'), wiki_set(), custom(), other2([sub, 0])],
        data = data,
        menu = [['other', load_lang('other')]]
    ))
        
@app.route('/topic/<everything:name>/sub/<sub>/b/<int:num>')
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
        
@app.route('/topic/<everything:name>/sub/<sub>/notice/<int:num>')
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
        
@app.route('/topic/<everything:name>/sub/<sub>/tool/<regex("close|stop|agree"):tool>')
def topic_stop(name = None, sub = None, tool = None):
    if tool == 'close':
        set_list = [
            'O', 
            'S', 
            load_lang('discussion', 1) + ' ' + load_lang('close', 1), 
            load_lang('discussion', 1) + ' ' + load_lang('open', 1)
        ]
    elif tool == 'stop':
        set_list = [
            '', 
            'O', 
            load_lang('discussion', 1) + ' ' + load_lang('stop', 1), 
            load_lang('discussion', 1) + ' ' + load_lang('restart', 1)
        ]
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
            curs.execute("select title from rd where title = ? and sub = ? and agree = 'O'", [name, sub])
            if curs.fetchall():
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, '" + load_lang('agreement', 1) + " X', ?, ?, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, time, ip])
                curs.execute("update rd set agree = '' where title = ? and sub = ?", [name, sub])
            else:
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, '" + load_lang('agreement', 1) + " O', ?, ?, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, time, ip])
                curs.execute("update rd set agree = 'O' where title = ? and sub = ?", [name, sub])
        else:
            curs.execute("select title from rd where title = ? and sub = ? and stop = ?", [name, sub, set_list[0]])
            if curs.fetchall():
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, ?, ?, ?, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, set_list[3], time, ip])
                curs.execute("update rd set stop = '' where title = ? and sub = ?", [name, sub])
            else:
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, ?, ?, ?, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, set_list[2], time, ip])
                curs.execute("update rd set stop = ? where title = ? and sub = ?", [set_list[0], name, sub])
        
        rd_plus(name, sub, time)
        
        conn.commit()
        
    return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub))    

@app.route('/topic/<everything:name>/sub/<sub>/admin/<int:num>')
def topic_admin(name = None, sub = None, num = None):
    curs.execute("select block, ip, date from topic where title = ? and sub = ? and id = ?", [name, sub, str(num)])
    data = curs.fetchall()
    if not data:
        return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub))

    ban = ''

    if admin_check(3) == 1:
        ban +=  '''
                </ul>
                <br>
                <h2>''' + load_lang('admin') + ' ' + load_lang('tool') + '''</h2>
                <ul>
                '''
        is_ban = '<li><a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/b/' + str(num) + '">'

        if data[0][0] == 'O':
            is_ban += load_lang('hide') + ' ' + load_lang('release')
        else:
            is_ban += load_lang('hide')
        
        is_ban +=   '''
                        </a>
                    </li>
                    <li>
                        <a href="/topic/''' + url_pas(name) + '/sub/' + url_pas(sub) + '/notice/' + str(num) + '''">
                    '''

        curs.execute("select id from topic where title = ? and sub = ? and id = ? and top = 'O'", [name, sub, str(num)])
        if curs.fetchall():
            is_ban += load_lang('notice') + ' ' + load_lang('release')
        else:
            is_ban += load_lang('notice') + ''
        
        is_ban += '</a></li></ul>'
        ban += '<li><a href="/ban/' + url_pas(data[0][1]) + '">'

        curs.execute("select end from ban where block = ?", [data[0][1]])
        if curs.fetchall():
            ban += load_lang('ban') + ' ' + load_lang('release')
        else:
            ban += load_lang('ban')
        
        ban += '</a></li>' + is_ban

    ban +=  '''
            </ul>
            <br>
            <h2>''' + load_lang('other') + ' ' + load_lang('tool') + '''</h2>
            <ul>
                <li>
                    <a href="/topic/''' + url_pas(name) + '/sub/' + url_pas(sub) + '/raw/' + str(num) + '''">raw</a>
                </li>
            '''
    ban = '<li>' + load_lang('time') + ' : ' + data[0][2] + '</li>' + ban
    
    if ip_or_user(data[0][1]) == 1:
        ban = '<li>' + load_lang('writer') + ' : ' + data[0][1] + ' <a href="/record/' + url_pas(data[0][1]) + '">(' + load_lang('record') + ')</a></li>' + ban
    else:
        ban =   '''
                <li>
                    ''' + load_lang('writer') + ' : <a href="/w/user:' + data[0][1] + '">' + data[0][1] + '</a> <a href="/record/' + url_pas(data[0][1]) + '">(' + load_lang('record') + ''')</a>
                </li>
                ''' + ban

    ban = '<h2>' + load_lang('state') + '</h2><ul>' + ban

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('discussion') + ' ' + load_lang('tool'), wiki_set(), custom(), other2([' (' + str(num) + ')', 0])],
        data = ban,
        menu = [['topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#' + str(num), load_lang('discussion')]]
    ))

@app.route('/topic/<everything:name>/sub/<sub>', methods=['POST', 'GET'])
def topic(name = None, sub = None):
    ban = topic_check(name, sub)
    admin = admin_check(3)
    
    if flask.request.method == 'POST':
        if captcha_post(flask.request.form.get('g-recaptcha-response', '')) == 1:
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

        match = re.search('^user:([^/]+)', name)
        if match:
            curs.execute('insert into alarm (name, data, date) values (?, ?, ?)', [match.groups()[0], ip + ' - <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '">' + load_lang('user', 1) + ' ' + load_lang('discussion', 1) + '</a>', today])
        
        data = re.sub('\[\[((?:분류|category):(?:(?:(?!\]\]).)*))\]\]', '[br]', flask.request.form.get('content', None))
        for rd_data in re.findall("(?:#([0-9]+))", data):
            curs.execute("select ip from topic where title = ? and sub = ? and id = ?", [name, sub, rd_data])
            ip_data = curs.fetchall()
            if ip_data and ip_or_user(ip_data[0][0]) == 0:
                curs.execute('insert into alarm (name, data, date) values (?, ?, ?)', [ip_data[0][0], ip + ' - <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#' + str(num) + '">' + load_lang('discussion', 1) + '</a>', today])
            
        data = re.sub("(?P<in>#(?:[0-9]+))", '[[\g<in>]]', data)

        data = savemark(data)

        rd_plus(name, sub, today)

        curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, ?, ?, ?, '', '')", [str(num), name, sub, data, today, ip])
        conn.commit()
        
        return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#reload')
    else:
        curs.execute("select title from rd where title = ? and sub = ? and stop = 'O'", [name, sub])
        close_data = curs.fetchall()
        
        curs.execute("select title from rd where title = ? and sub = ? and stop = 'S'", [name, sub])
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

            curs.execute("select title from rd where title = ? and sub = ? and agree = 'O'", [name, sub])
            if curs.fetchall():
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/agree">(' + load_lang('release') + ')</a>'
            else:
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/agree">(' + load_lang('agreement') + ')</a>'
            
            all_data += '<hr class=\"main_hr\">'
        
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
                                
            all_data += '''
                        <table id="toron">
                            <tbody>
                                <tr>
                                    <td id="toron_color_red">
                                        <a href="#''' + topic_data[1] + '''">
                                            #''' + topic_data[1] + '''
                                        </a> ''' + ip_pas(topic_data[3]) + who_plus + ''' <span style="float: right;">''' + topic_data[2] + '''</span>
                                    </td>
                                </tr>
                                <tr>
                                    <td>''' + render_set(data = topic_data[0]) + '''</td>
                                </tr>
                            </tbody>
                        </table>
                        <br>
                        '''    

        for topic_data in topic:
            user_write = topic_data[0]

            if number == 1:
                start = topic_data[3]

            if topic_data[4] == 'O':
                blind_data = 'id="toron_color_grey"'
                
                if admin != 1:
                    curs.execute("select who from re_admin where what = ? order by time desc limit 1", ['blind (' + name + ' - ' + sub + '#' + str(number) + ')'])
                    who_blind = curs.fetchall()
                    if who_blind:
                        user_write = '[[user:' + who_blind[0][0] + ']] ' + load_lang('hide')
                    else:
                        user_write = load_lang('hide')
            else:
                blind_data = ''

            user_write = render_set(data = user_write)
            ip = ip_pas(topic_data[3])
            
            curs.execute('select acl from user where id = ?', [topic_data[3]])
            user_acl = curs.fetchall()
            if user_acl and user_acl[0][0] != 'user':
                ip += ' <a href="javascript:void(0);" title="' + load_lang('admin') + '">★</a>'

            if admin == 1 or blind_data == '':
                ip += ' <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/admin/' + str(number) + '">(' + load_lang('discussion') + ' ' + load_lang('tool') + ')</a>'

            curs.execute("select end from ban where block = ?", [topic_data[3]])
            if curs.fetchall():
                ip += ' <a href="javascript:void(0);" title="' + load_lang('blocked') + '">†</a>'
                    
            if topic_data[5] == '1':
                color = '_blue'
            elif topic_data[3] == start:
                color = '_green'
            else:
                color = ''
                
            if user_write == '':
                user_write = '<br>'
                         
            all_data += '''
                        <table id="toron">
                            <tbody>
                                <tr>
                                    <td id="toron_color''' + color + '''">
                                        <a href="javascript:void(0);" id="''' + str(number) + '">#' + str(number) + '</a> ' + ip + '''</span>
                                    </td>
                                </tr>
                                <tr ''' + blind_data + '''>
                                    <td>''' + user_write + '''</td>
                                </tr>
                            </tbody>
                        </table>
                        <br>
                        '''
            number += 1

        if ban != 1 or admin == 1:
            data += '''
                    <div id="plus"></div>
                    <script type="text/javascript" src="/views/main_css/topic_reload.js"></script>
                    <script>topic_load("''' + name + '''", "''' + sub + '''");</script>
                    <a id="reload" href="javascript:void(0);" onclick="location.href.endsWith(\'#reload\')? location.reload(true):location.href=\'#reload\'">(''' + load_lang('reload') + ''')</a>
                    <form style="''' + display + '''" method="post">
                    <br>
                    <textarea style="height: 100px;" name="content"></textarea>
                    <hr class=\"main_hr\">
                    ''' + captcha_get()
            
            if display == '':
                data += ip_warring()

            data += '''
                        <button type="submit">''' + load_lang('send') + '''</button>
                    </form>
                    '''

        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2([' (' + load_lang('discussion') + ')', 0])],
            data = '<h2 id="topic_top_title">' + sub + '</h2>' + all_data + data,
            menu = [['topic/' + url_pas(name), load_lang('list')]]
        ))

@app.route('/tool/<name>')
def user_tool(name = None):
    data =  '''
            <h2>''' + load_lang('tool') + '''</h2>
            <ul>
                <li><a href="/record/''' + url_pas(name) + '''">''' + load_lang('record') + '''</a></li>
            </ul>
            '''
            
    if admin_check(1) == 1:
        curs.execute("select block from ban where block = ?", [name])
        if curs.fetchall():
            ban_name = load_lang('release')
        else:
            ban_name = load_lang('ban')
    
        data += '''
                <h2>''' + load_lang('admin') + '''</h2>
                <ul>
                    <li><a href="/ban/''' + url_pas(name) + '''">''' + ban_name + '''</a></li>
                    <li><a href="/check/''' + url_pas(name) + '''">''' + load_lang('check') + '''</a></li>
                </ul>
                '''

    return easy_minify(flask.render_template(skin_check(), 
        imp = [name, wiki_set(), custom(), other2([' (' + load_lang('tool') + ')', 0])],
        data = data,
        menu = [['topic/' + url_pas(name), load_lang('list')]]
    ))
        
@app.route('/topic/<everything:name>', methods=['POST', 'GET'])
@app.route('/topic/<everything:name>/<regex("close|agree"):tool>', methods=['GET'])
def close_topic_list(name = None, tool = None):
    div = ''
    
    if flask.request.method == 'POST':
        t_num = ''
        
        while 1:
            curs.execute("select title from topic where title = ? and sub = ? limit 1", [name, flask.request.form.get('topic', None) + t_num])
            if curs.fetchall():
                if t_num == '':
                    t_num = ' 2'
                else:
                    t_num = ' ' + str(int(t_num.replace(' ', '')) + 1)
            else:
                break

        return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(flask.request.form.get('topic', None) + t_num))
    else:
        plus = ''
        menu = [['topic/' + url_pas(name), load_lang('list')]]
        
        if tool == 'close':
            curs.execute("select sub from rd where title = ? and stop = 'O' order by sub asc", [name])
            
            sub = load_lang('close') + ''
        elif tool == 'agree':
            curs.execute("select sub from rd where title = ? and agree = 'O' order by sub asc", [name])
            
            sub = load_lang('agreement') + ''
        else:
            curs.execute("select sub from rd where title = ? order by date desc", [name])
            
            sub = load_lang('discussion') + ' ' + load_lang('list')
            
            menu = [['w/' + url_pas(name), load_lang('document')]]
            
            plus =  '''
                    <a href="/topic/''' + url_pas(name) + '''/close">(''' + load_lang('close') + ''')</a> <a href="/topic/''' + url_pas(name) + '''/agree">(''' + load_lang('agreement') + ''')</a>
                    <hr class=\"main_hr\">
                    <input placeholder="''' + load_lang('discussion') + ' ' + load_lang('name') + '''" name="topic" type="text">
                    <hr class=\"main_hr\">
                    <button type="submit">''' + load_lang('open') + '''</button>
                    '''

        for data in curs.fetchall():
            curs.execute("select data, date, ip, block from topic where title = ? and sub = ? and id = '1'", [name, data[0]])
            if curs.fetchall():                
                it_p = 0
                
                if sub == load_lang('discussion') + ' ' + load_lang('list'):
                    curs.execute("select title from rd where title = ? and sub = ? and stop = 'O' order by sub asc", [name, data[0]])
                    if curs.fetchall():
                        it_p = 1
                
                if it_p != 1:
                    div += '<h2><a href="/topic/' + url_pas(name) + '/sub/' + url_pas(data[0]) + '">' + data[0] + '</a></h2>'

        if div == '':
            plus = re.sub('^<br>', '', plus)
        
        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2([' (' + sub + ')', 0])],
            data =  '<form method="post">' + div + plus + '</form>',
            menu = menu
        ))
        
@app.route('/login', methods=['POST', 'GET'])
def login():
    if custom()[2] != 0:
        return redirect('/user')
    
    if ban_check(tool = 'login') == 1:
        return re_error('/ban')
        
    if flask.request.method == 'POST':        
        if captcha_post(flask.request.form.get('g-recaptcha-response', '')) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        ip = ip_check()
        agent = flask.request.headers.get('User-Agent')

        curs.execute("select pw, encode from user where id = ?", [flask.request.form.get('id', None)])
        user = curs.fetchall()
        if not user:
            return re_error('/error/2')

        pw_check_d = pw_check(
            flask.request.form.get('pw', ''), 
            user[0][0],
            user[0][1],
            flask.request.form.get('id', None)
        )
        if pw_check_d != 1:
            return re_error('/error/10')

        flask.session['state'] = 1
        flask.session['id'] = flask.request.form.get('id', None)
        
        curs.execute("select css from custom where user = ?", [flask.request.form.get('id', None)])
        css_data = curs.fetchall()
        if css_data:
            flask.session['head'] = css_data[0][0]
        else:
            flask.session['head'] = ''

        curs.execute("insert into ua_d (name, ip, ua, today, sub) values (?, ?, ?, ?, '')", [flask.request.form.get('id', None), ip_check(1), agent, get_time()])

        conn.commit()
        
        return redirect('/user')  
    else:
        oauth_content = '<link rel="stylesheet" href="/views/main_css/oauth.css"><div class="oauth-wrapper"><ul class="oauth-list">'
        oauth_supported = load_oauth('_README')['support']
        for i in range(len(oauth_supported)):
            oauth_data = load_oauth(oauth_supported[i])
            if oauth_data['client_id'] != '' and oauth_data['client_secret'] != '':
                oauth_content +=    '''
                                    <li>
                                        <a href="/oauth/{}/init">
                                            <div class="oauth-btn oauth-btn-{}">
                                                <div class="oauth-btn-logo oauth-btn-{}"></div>
                                                {}
                                            </div>
                                        </a>
                                    </li>
                                    '''.format(
                                        oauth_supported[i], 
                                        oauth_supported[i], 
                                        oauth_supported[i], 
                                        load_lang('oauth_signin_' + oauth_supported[i])
                                    )
        
        oauth_content += '</ul></div>'
        
        return easy_minify(flask.render_template(skin_check(),    
            imp = [load_lang('login'), wiki_set(), custom(), other2([0, 0])],
            data =  '''
                    <form method="post">
                        <input placeholder="''' + load_lang('id') + '''" name="id" type="text">
                        <hr class=\"main_hr\">
                        <input placeholder="''' + load_lang('password') + '''" name="pw" type="password">
                        <hr class=\"main_hr\">
                        ''' + captcha_get() + '''
                        <button type="submit">''' + load_lang('login') + '''</button><a href="/register">''' + load_lang('register_suggest') + '''</a>
                        <hr class=\"main_hr\">
                        ''' + oauth_content + '''
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('http_warring') + '''</span>
                    </form>
                    ''',
            menu = [['user', load_lang('user')]]
        ))

@app.route('/oauth/<regex("naver|facebook"):platform>/<regex("init|callback"):func>', methods=['GET', 'POST'])
def login_oauth(platform = None, func = None):
    publish_url = load_oauth('publish_url')
    oauth_data = load_oauth(platform)
    api_url = {}
    data = {
        'client_id' : oauth_data['client_id'],
        'client_secret' : oauth_data['client_secret'],
        'redirect_uri' : publish_url + '/oauth/' + platform + '/callback',
        'state' : 'RAMDOMVALUE'
    }

    if platform == 'naver':
        api_url['redirect'] = 'https://nid.naver.com/oauth2.0/authorize'
        api_url['token'] = 'https://nid.naver.com/oauth2.0/token'
        api_url['profile'] = 'https://openapi.naver.com/v1/nid/me'
    elif platform == 'facebook':
        api_url['redirect'] = 'https://www.facebook.com/v3.1/dialog/oauth'
        api_url['token'] = 'https://graph.facebook.com/v3.1/oauth/access_token'
        api_url['profile'] = 'https://graph.facebook.com/me'

    if func == 'init':
        if oauth_data['client_id'] == '' or oauth_data['client_secret'] == '':
            return easy_minify(flask.render_template(skin_check(), 
                imp = [load_lang('login'), wiki_set(), custom(), other2([0, 0])], 
                data = load_lang('oauth_disabled'), 
                menu = [['user', load_lang('user')]]
            ))
        elif publish_url == 'https://':
            return easy_minify(flask.render_template(skin_check(), 
                imp = [load_lang('login'), wiki_set(), custom(), other2([0, 0])], 
                data = load_lang('oauth_settings_not_found'), 
                menu = [['user', load_lang('user')]]
            ))

        referrer_re = re.compile(r'(?P<host>^(https?):\/\/([^\/]+))\/(?P<refer>[^\/?]+)')
        if flask.request.referrer != None:
            referrer = referrer_re.search(flask.request.referrer)
            if referrer.group('host') != load_oauth('publish_url'):
                return redirect('/')
            else:
                flask.session['referrer'] = referrer.group('refer')
        else:
            return redirect('/')

        flask.session['refer'] = flask.request.referrer

        if platform == 'naver':
            return redirect(api_url['redirect'] + '?response_type=code&client_id={}&redirect_uri={}&state={}'.format(data['client_id'], data['redirect_uri'], data['state']))
        elif platform == 'facebook':
            return redirect(api_url['redirect'] + '?client_id={}&redirect_uri={}&state={}'.format(data['client_id'], data['redirect_uri'], data['state']))

    elif func == 'callback':
        code = flask.request.args.get('code')
        state = flask.request.args.get('state')
        if code == None or state == None:
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('inter_error'), wiki_set(), custom(), other2([0, 0])],
                data =  '''
                        <p>''' + load_lang('inter_error_detail') + '''</p>
                        <hr>
                        <code>ie_wrong_callback</code>
                        <p>''' + load_lang('ie_wrong_callback') + '''</p>
                        ''',
                menu = [['user', load_lang('user')]]
            ))

        if platform == 'naver':
            token_access = api_url['token']+'?grant_type=authorization_code&client_id={}&client_secret={}&code={}&state={}'.format(data['client_id'], data['client_secret'], code, state)
            token_result = urllib.request.urlopen(token_access).read().decode('utf-8')
            token_result_json = json.loads(token_result)

            headers = {'Authorization': 'Bearer {}'.format(token_result_json['access_token'])}

            profile_access = urllib.request.Request(api_url['profile'], headers = headers)
            profile_result = urllib.request.urlopen(profile_access).read().decode('utf-8')
            profile_result_json = json.loads(profile_result)

            stand_json = {'id' : profile_result_json['response']['id'], 'name' : profile_result_json['response']['name'], 'picture' : profile_result_json['response']['profile_image']}
        elif platform == 'facebook':
            token_access = api_url['token']+'?client_id={}&redirect_uri={}&client_secret={}&code={}'.format(data['client_id'], data['redirect_uri'], data['client_secret'], code)
            token_result = urllib.request.urlopen(token_access).read().decode('utf-8')
            token_result_json = json.loads(token_result)

            profile_access = api_url['profile']+'?fields=id,name,picture&access_token={}'.format(token_result_json['access_token'])
            profile_result = urllib.request.urlopen(profile_access).read().decode('utf-8')
            profile_result_json = json.loads(profile_result)

            stand_json = {'id': profile_result_json['id'], 'name': profile_result_json['name'], 'picture': profile_result_json['picture']['data']['url']}
        
        if flask.session['referrer'][0:6] == 'change':
            curs.execute('select * from oauth_conn where wiki_id = ? and provider = ?', [flask.session['id'], platform])
            oauth_result = curs.fetchall()
            if len(oauth_result) == 0:
                curs.execute('insert into oauth_conn (provider, wiki_id, sns_id, name, picture) values(?, ?, ?, ?, ?)', [platform, flask.session['id'], stand_json['id'], stand_json['name'], stand_json['picture']])
            else:
                curs.execute('update oauth_conn set name = ? picture = ? where wiki_id = ?', [stand_json['name'], stand_json['pricture'], flask.session['id']])

            conn.commit()
        elif flask.session['referrer'][0:5] == 'login':
            curs.execute('select * from oauth_conn where provider = ? and sns_id = ?', [platform, stand_json['id']])
            curs_result = curs.fetchall()
            if len(curs_result) == 0:
                return re_error('/error/2')
            else:
                flask.session['state'] = 1
                flask.session['id'] = curs_result[0][2]
        
        return redirect(flask.session['refer'])
                
@app.route('/change', methods=['POST', 'GET'])
def change_password():
    global support_language

    if ban_check() == 1:
        return re_error('/ban')

    if custom()[2] == 0:
        return redirect('/login')

    ip = ip_check()
    user_state = flask.request.args.get('user', 'ip')
    
    if user_state == 'ip':
        if flask.request.method == 'POST':    
            if flask.request.form.get('pw4', None) and flask.request.form.get('pw2', None):
                if flask.request.form.get('pw2', None) != flask.request.form.get('pw3', None):
                    return re_error('/error/20')

                curs.execute("select pw, encode from user where id = ?", [flask.session['id']])
                user = curs.fetchall()
                if not user:
                    return re_error('/error/2')
                
                pw_check_d = pw_check(
                    flask.request.form.get('pw4', ''), 
                    user[0][0],
                    user[0][1],
                    flask.request.form.get('id', None)
                )
                if pw_check_d != 1:
                    return re_error('/error/10')

                hashed = pw_encode(flask.request.form.get('pw2', None))
                
                curs.execute("update user set pw = ? where id = ?", [hashed, flask.session['id']])

            auto_list = ['email', 'skin', 'lang']

            for auto_data in auto_list:
                curs.execute('select data from user_set where name = ? and id = ?', [auto_data, ip])
                if curs.fetchall():
                    curs.execute("update user_set set data = ? where name = ? and id = ?", [flask.request.form.get(auto_data, ''), auto_data, ip])
                else:
                    curs.execute("insert into user_set (name, id, data) values (?, ?, ?)", [auto_data, ip, flask.request.form.get(auto_data, '')])

            conn.commit()
            
            return redirect('/change')
        else:        
            curs.execute('select data from user_set where name = "email" and id = ?', [ip])
            data = curs.fetchall()
            if data:
                email = data[0][0]
            else:
                email = ''

            div2 = load_skin()
            
            div3 = ''
            var_div3 = ''

            curs.execute('select data from user_set where name = "lang" and id = ?', [flask.session['id']])
            data = curs.fetchall()

            for lang_data in support_language:
                if data and data[0][0] == lang_data:
                    div3 = '<option value="' + lang_data + '">' + lang_data + '</option>'
                else:
                    var_div3 += '<option value="' + lang_data + '">' + lang_data + '</option>'

            div3 += var_div3

            oauth_provider = load_oauth('_README')['support']
            oauth_content = '<ul>'
            for i in range(len(oauth_provider)):
                curs.execute('select name, picture from oauth_conn where wiki_id = ? and provider = ?', [flask.session['id'], oauth_provider[i]])
                oauth_data = curs.fetchall()
                if len(oauth_data) == 1:
                    oauth_content += '<li>{} - {}</li>'.format(oauth_provider[i], load_lang('connection') + load_lang('oauth_conn_done') + ': <img src="{}" width="17px" height="17px">{}'.format(oauth_data[0][1], oauth_data[0][0]))
                else:
                    oauth_content += '<li>{} - {}</li>'.format(oauth_provider[i], load_lang('connection') + load_lang('oauth_conn_not') + '. <a href="/oauth/{}/init">{}</a>'.format(oauth_provider[i], load_lang('oauth_conn_new')))
            
            oauth_content += '</ul>'

            return easy_minify(flask.render_template(skin_check(),    
                imp = [load_lang('user') + ' ' + load_lang('setting') + ' ' + load_lang('edit'), wiki_set(), custom(), other2([0, 0])],
                data =  '''
                        <form method="post">
                            <span>id : ''' + ip + '''</span>
                            <hr class=\"main_hr\">
                            <input placeholder="''' + load_lang('now') + ' ' + load_lang('password') + '''" name="pw4" type="password">
                            <br>
                            <br>
                            <input placeholder="''' + load_lang('new') + ' ' + load_lang('password') + '''" name="pw2" type="password">
                            <br>
                            <br>
                            <input placeholder="''' + load_lang('password') + ' ' + load_lang('confirm') + '''" name="pw3" type="password">
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('user') + ' ' + load_lang('skin') + '''</span>
                            <br>
                            <br>
                            <select name="skin">''' + div2 + '''</select>
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('user') + ' ' + load_lang('language') + '''</span>
                            <br>
                            <br>
                            <select name="lang">''' + div3 + '''</select>
                            <hr class=\"main_hr\">
                            <span>OAuth ''' + load_lang('connection') + '''</span>
                            ''' + oauth_content + '''
                            <hr class=\"main_hr\">
                            <button type="submit">''' + load_lang('edit') + '''</button>
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('http_warring') + '''</span>
                        </form>
                        ''',
                menu = [['user', load_lang('user')]]
            ))
    else:
        pass

@app.route('/check/<name>')
def user_check(name = None):
    curs.execute("select acl from user where id = ? or id = ?", [name, flask.request.args.get('plus', '-')])
    user = curs.fetchall()
    if user and user[0][0] != 'user':
        if admin_check() != 1:
            return re_error('/error/4')

    if admin_check(4, 'check (' + name + ')') != 1:
        return re_error('/error/3')
        
    num = int(flask.request.args.get('num', 1))
    if num * 50 > 0:
        sql_num = num * 50 - 50
    else:
        sql_num = 0
    
    if flask.request.args.get('plus', None):
        end_check = 1
    
        if ip_or_user(name) == 1:
            if ip_or_user(flask.request.args.get('plus', None)) == 1:
                curs.execute("select name, ip, ua, today from ua_d where ip = ? or ip = ? order by today desc limit ?, '50'", [name, flask.request.args.get('plus', None), sql_num])
            else:
                curs.execute("select name, ip, ua, today from ua_d where ip = ? or name = ? order by today desc limit ?, '50'", [name, flask.request.args.get('plus', None), sql_num])
        else:
            if ip_or_user(flask.request.args.get('plus', None)) == 1:
                curs.execute("select name, ip, ua, today from ua_d where name = ? or ip = ? order by today desc limit ?, '50'", [name, flask.request.args.get('plus', None), sql_num])
            else:
                curs.execute("select name, ip, ua, today from ua_d where name = ? or name = ? order by today desc limit ?, '50'", [name, flask.request.args.get('plus', None), sql_num])
    else:
        end_check = 0
        
        if ip_or_user(name) == 1:
            curs.execute("select name, ip, ua, today from ua_d where ip = ? order by today desc limit ?, '50'", [name, sql_num])
        else:
            curs.execute("select name, ip, ua, today from ua_d where name = ? order by today desc limit ?, '50'", [name, sql_num])
    
    record = curs.fetchall()
    if record:
        if not flask.request.args.get('plus', None):
            div = '<a href="/manager/14?plus=' + url_pas(name) + '">(' + load_lang('compare') + ')</a><hr class=\"main_hr\">'
        else:
            div = '<a href="/check/' + url_pas(name) + '">(' + name + ')</a> <a href="/check/' + url_pas(flask.request.args.get('plus', None)) + '">(' + flask.request.args.get('plus', None) + ')</a><hr class=\"main_hr\">'

        div +=  '''
                <table id="main_table_set">
                    <tbody>
                        <tr>
                            <td id="main_table_width">''' + load_lang('name') + '''</td>
                            <td id="main_table_width">ip</td>
                            <td id="main_table_width">''' + load_lang('time') + '''</td>
                        </tr>
                '''
        
        for data in record:
            if data[2]:
                ua = data[2]
            else:
                ua = '<br>'

            div +=  '''
                    <tr>
                        <td>''' + ip_pas(data[0]) + '''</td>
                        <td>''' + ip_pas(data[1]) + '''</td>
                        <td>''' + data[3] + '''</td>
                    </tr>
                    <tr>
                        <td colspan="3">''' + ua + '''</td>
                    </tr>
                    '''
        
        div +=  '''
                    </tbody>
                </table>
                '''
    else:
        return re_error('/error/2')
        
    if end_check == 1:
        div += next_fix('/check/' + url_pas(name) + '?plus=' + flask.request.args.get('plus', None) + '&num=', num, record)
    else:
        div += next_fix('/check/' + url_pas(name) + '?num=', num, record)
            
    return easy_minify(flask.render_template(skin_check(),    
        imp = [load_lang('check'), wiki_set(), custom(), other2([0, 0])],
        data = div,
        menu = [['manager', load_lang('admin')]]
    ))
                
@app.route('/register', methods=['POST', 'GET'])
def register():
    if ban_check() == 1:
        return re_error('/ban')

    if custom()[2] != 0:
        return redirect('/user')

    if not admin_check() == 1:
        curs.execute('select data from other where name = "reg"')
        set_d = curs.fetchall()
        if set_d and set_d[0][0] == 'on':
            return re_error('/ban')
    
    if flask.request.method == 'POST': 
        if captcha_post(flask.request.form.get('g-recaptcha-response', '')) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        if flask.request.form.get('pw', None) != flask.request.form.get('pw2', None):
            return re_error('/error/20')

        if re.search('(?:[^A-Za-zㄱ-힣0-9 ])', flask.request.form.get('id', None)):
            return re_error('/error/8')
            
        curs.execute('select html from html_filter where kind = "name"')
        set_d = curs.fetchall()
        for i in set_d:
            check_r = re.compile(i[0], re.I)
            if check_r.search(flask.request.form.get('id', None)):
                return re_error('/error/8')

        if len(flask.request.form.get('id', None)) > 32:
            return re_error('/error/7')

        curs.execute("select id from user where id = ?", [flask.request.form.get('id', None)])
        if curs.fetchall():
            return re_error('/error/6')

        hashed = pw_encode(flask.request.form.get('pw', None))
        
        curs.execute('select data from other where name = "email_have"')
        sql_data = curs.fetchall()
        if sql_data and sql_data[0][0] != '':
            flask.session['c_id'] = flask.request.form.get('id', None)
            flask.session['c_pw'] = hashed
            flask.session['c_key'] = ''.join(random.choice("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") for i in range(16))

            return redirect('/need_email')
        else:
            curs.execute('select data from other where name = "encode"')
            db_data = curs.fetchall()

            curs.execute("select id from user limit 1")
            if not curs.fetchall():
                curs.execute("insert into user (id, pw, acl, date, encode) values (?, ?, 'owner', ?, ?)", [flask.request.form.get('id', None), hashed, get_time(), db_data[0][0]])

                first = 1
            else:
                curs.execute("insert into user (id, pw, acl, date, encode) values (?, ?, 'user', ?, ?)", [flask.request.form.get('id', None), hashed, get_time(), db_data[0][0]])

                first = 0

            ip = ip_check()
            agent = flask.request.headers.get('User-Agent')

            curs.execute("insert into ua_d (name, ip, ua, today, sub) values (?, ?, ?, ?, '')", [flask.request.form.get('id', None), ip, agent, get_time()])  

            flask.session['state'] = 1
            flask.session['id'] = flask.request.form.get('id', None)
            flask.session['head'] = ''
                  
            conn.commit()
            
            if first == 0:
                return redirect('/change')
            else:
                return redirect('/setting')
    else:        
        contract = ''
        
        curs.execute('select data from other where name = "contract"')
        data = curs.fetchall()
        if data and data[0][0] != '':
            contract = data[0][0] + '<hr class=\"main_hr\">'

        return easy_minify(flask.render_template(skin_check(),    
            imp = [load_lang('register'), wiki_set(), custom(), other2([0, 0])],
            data =  '''
                    <form method="post">
                        ''' + contract + '''
                        <input placeholder="''' + load_lang('id') + '''" name="id" type="text">
                        <hr class=\"main_hr\">
                        <input placeholder="''' + load_lang('password') + '''" name="pw" type="password">
                        <hr class=\"main_hr\">
                        <input placeholder="''' + load_lang('confirm') + '''" name="pw2" type="password">
                        <hr class=\"main_hr\">
                        ''' + captcha_get() + '''
                        <button type="submit">''' + load_lang('register') + '''</button>
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('http_warring') + '''</span>
                    </form>
                    ''',
            menu = [['user', load_lang('user')]]
        ))

@app.route('/<regex("need_email|pass_find"):tool>', methods=['POST', 'GET'])
def need_email(tool = 'pass_find'):
    if flask.request.method == 'POST':
        if tool == 'need_email':
            if 'c_id' in flask.session:
                main_email = ['naver.com', 'gmail.com', 'daum.net', 'hanmail.net', 'hanmail2.net']
                data = re.search('@([^@]+)$', flask.request.form.get('email', ''))
                if data:
                    data = data.groups()[0]

                    curs.execute("select html from html_filter where html = ? and kind = 'email'", [data])
                    if curs.fetchall() or (data in main_email):
                        curs.execute('select id from user_set where name = "email" and data = ?', [flask.request.form.get('email', '')])
                        if curs.fetchall():
                            flask.session.pop('c_id', None)
                            flask.session.pop('c_pw', None)
                            flask.session.pop('c_key', None)

                            return redirect('/register')
                        else:
                            send_email(flask.request.form.get('email', ''), wiki_set()[0] + ' key', 'key : ' + flask.session['c_key'])
                            flask.session['c_email'] = flask.request.form.get('email', '')

                            return redirect('/check_key')

            return redirect('/register')
        else:
            curs.execute("select id from user where id = ? and email = ?", [flask.request.form.get('id', ''), flask.request.form.get('email', '')])
            if curs.fetchall():
                flask.session['c_key'] = ''.join(random.choice("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") for i in range(16))
                flask.session['c_id'] = flask.request.form.get('id', '')

                send_email(flask.request.form.get('email', ''), wiki_set()[0] + ' ' + load_lang('password') + ' ' + load_lang('search') + ' key', 'key : ' + flask.session['c_key'])

                return redirect('/check_pass_key')
    else:
        if tool == 'need_email':
            return easy_minify(flask.render_template(skin_check(),    
                imp = ['email', wiki_set(), custom(), other2([0, 0])],
                data =  '''
                        <a href="/email_filter">(email ''' + load_lang('list') + ''')</a>
                        <hr class=\"main_hr\">
                        <form method="post">
                            <input placeholder="email" name="email" type="text">
                            <hr class=\"main_hr\">
                            <button type="submit">''' + load_lang('save') + '''</button>
                        </form>
                        ''',
                menu = [['user', load_lang('user')]]
            ))
        else:
            return easy_minify(flask.render_template(skin_check(),    
                imp = [load_lang('password') + ' ' + load_lang('search'), wiki_set(), custom(), other2([0, 0])],
                data =  '''
                        <form method="post">
                            <input placeholder="''' + load_lang('id') + '''" name="id" type="text">
                            <hr class=\"main_hr\">
                            <input placeholder="email" name="email" type="text">
                            <hr class=\"main_hr\">
                            <button type="submit">''' + load_lang('save') + '''</button>
                        </form>
                        ''',
                menu = [['user', load_lang('user')]]
            ))

@app.route('/<regex("check_key|check_pass_key"):tool>', methods=['POST', 'GET'])
def check_key(tool = 'check_pass_key'):
    if flask.request.method == 'POST':
        if tool == 'check_key':
            if 'c_id' in flask.session and flask.session['c_key'] == flask.request.form.get('key', None):
                curs.execute('select data from other where name = "encode"')
                db_data = curs.fetchall()

                curs.execute("select id from user limit 1")
                if not curs.fetchall():
                    curs.execute("insert into user (id, pw, acl, date, encode) values (?, ?, 'owner', ?)", [flask.session['c_id'], flask.session['c_pw'], get_time(), db_data[0][0]])

                    first = 1
                else:
                    curs.execute("insert into user (id, pw, acl, date, encode) values (?, ?, 'user', ?)", [flask.session['c_id'], flask.session['c_pw'], get_time(), db_data[0][0]])

                    first = 0

                ip = ip_check()
                agent = flask.request.headers.get('User-Agent')

                curs.execute("insert into user_set (name, id, data) values ('email', ?, ?)", [flask.session['c_id'], flask.session['c_email']])
                curs.execute("insert into ua_d (name, ip, ua, today, sub) values (?, ?, ?, ?, '')", [flask.session['c_id'], ip, agent, get_time()])

                flask.session['state'] = 1
                flask.session['id'] = flask.session['c_id']
                flask.session['head'] = ''
                        
                conn.commit()
                
                flask.session.pop('c_id', None)
                flask.session.pop('c_pw', None)
                flask.session.pop('c_key', None)
                flask.session.pop('c_email', None)

                if first == 0:
                    return redirect('/change')
                else:
                    return redirect('/setting')
            else:
                flask.session.pop('c_id', None)
                flask.session.pop('c_pw', None)
                flask.session.pop('c_key', None)
                flask.session.pop('c_email', None)

                return redirect('/register')
        else:
            if 'c_id' in flask.session and flask.session['c_key'] == flask.request.form.get('key', None):
                hashed = pw_encode(flask.session['c_key'])
                curs.execute("update user set pw = ? where id = ?", [hashed, flask.session['c_id']])

                d_id = flask.session['c_id']
                pw = flask.session['c_key']

                flask.session.pop('c_id', None)
                flask.session.pop('c_key', None)

                return easy_minify(flask.render_template(skin_check(),    
                    imp = ['check', wiki_set(), custom(), other2([0, 0])],
                    data =  '''
                            ''' + load_lang('id') + ' : ' + d_id + '''
                            <br>
                            ''' + load_lang('password') + ' : ' + pw + '''
                            ''',
                    menu = [['user', load_lang('user')]]
                ))
            else:
                return redirect('/pass_find')
    else:
        return easy_minify(flask.render_template(skin_check(),    
            imp = ['check', wiki_set(), custom(), other2([0, 0])],
            data =  '''
                    <form method="post">
                        <input placeholder="key" name="key" type="text">
                        <hr class=\"main_hr\">
                        <button type="submit">''' + load_lang('save') + '''</button>
                    </form>
                    ''',
            menu = [['user', load_lang('user')]]
        ))
           
@app.route('/logout')
def logout():
    flask.session['state'] = 0
    flask.session.pop('id', None)

    return redirect('/user')
    
@app.route('/ban/<name>', methods=['POST', 'GET'])
def user_ban(name = None):
    if ip_or_user(name) == 0:
        curs.execute("select acl from user where id = ?", [name])
        user = curs.fetchall()
        if not user:
            return re_error('/error/2')

        if user and user[0][0] != 'user':
            if admin_check() != 1:
                return re_error('/error/4')

    if ban_check(ip = ip_check(), tool = 'login') == 1:
        return re_error('/ban')
                
    if flask.request.method == 'POST':
        if admin_check(1, 'ban (' + name + ')') != 1:
            return re_error('/error/3')

        if flask.request.form.get('limitless', '') == '':
            end = flask.request.form.get('second', '0')
        else:
            end = '0'

        ban_insert(name, end, flask.request.form.get('why', ''), flask.request.form.get('login', ''), ip_check())

        return redirect('/ban/' + url_pas(name))     
    else:
        if admin_check(1) != 1:
            return re_error('/error/3')

        curs.execute("select end, why from ban where block = ?", [name])
        end = curs.fetchall()
        if end:
            now = load_lang('release')

            if end[0][0] == '':
                data = '<ul><li>' + load_lang('limitless') + ' ' + load_lang('ban') + '</li>'
            else:
                data = '<ul><li>' + load_lang('ban') + ' : ' + end[0][0] + '</li>'
                
            curs.execute("select block from ban where block = ? and login = 'O'", [name])
            if curs.fetchall():
                data += '<li>' + load_lang('login') + ' ' + load_lang('able') + '</li>'

            if end[0][1] != '':
                data += '<li>' + load_lang('why') + ' : ' + end[0][1] + '</li></ul><hr class=\"main_hr\">'
            else:
                data += '</ul><hr class=\"main_hr\">'
        else:
            if re.search("^([0-9]{1,3}\.[0-9]{1,3})$", name):
                now = load_lang('band') + ' ' + load_lang('ban')
            else:
                now = load_lang('ban')
                
            if ip_or_user(name) == 1:
                plus = '<input type="checkbox" name="login"> ' + load_lang('login') + ' ' + load_lang('able') + '<hr class=\"main_hr\">'
            else:
                plus = ''

            data =  '''
                    <input placeholder="''' + load_lang('second') + '''" name="second" type="text">
                    <hr class=\"main_hr\">
                    <input type="checkbox" name="limitless"> ''' + load_lang('limitless') + '''
                    <hr class=\"main_hr\">
                    <input placeholder="''' + load_lang('why') + '''" name="why" type="text">
                    <hr class=\"main_hr\">
                    ''' + plus

        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2([' (' + now + ')', 0])],
            data =  '''
                    <form method="post">
                        ''' + data + '''
                        <button type="submit">''' + now + '''</button>
                    </form>
                    ''',
            menu = [['manager', load_lang('admin')]]
        ))            
                
@app.route('/acl/<everything:name>', methods=['POST', 'GET'])
def acl(name = None):
    check_ok = ''
    
    if flask.request.method == 'POST':
        check_data = 'acl (' + name + ')'
    else:
        check_data = None
    
    user_data = re.search('^user:(.+)$', name)
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

    if flask.request.method == 'POST':
        curs.execute("select title from acl where title = ?", [name])
        if curs.fetchall():
            curs.execute("update acl set dec = ? where title = ?", [flask.request.form.get('dec', ''), name])
            curs.execute("update acl set dis = ? where title = ?", [flask.request.form.get('dis', ''), name])
            curs.execute("update acl set why = ? where title = ?", [flask.request.form.get('why', ''), name])
            curs.execute("update acl set view = ? where title = ?", [flask.request.form.get('view', ''), name])
        else:
            curs.execute("insert into acl (title, dec, dis, why, view) values (?, ?, ?, ?, ?)", [name, flask.request.form.get('dec', ''), flask.request.form.get('dis', ''), flask.request.form.get('why', ''), flask.request.form.get('view', '')])
        
        curs.execute("select title from acl where title = ? and dec = '' and dis = ''", [name])
        if curs.fetchall():
            curs.execute("delete from acl where title = ?", [name])

        conn.commit()
            
        return redirect('/acl/' + url_pas(name))            
    else:
        data = '' + load_lang('document') + ' acl<br><br><select name="dec" ' + check_ok + '>'
    
        if re.search('^user:', name):
            acl_list = [['', load_lang('normal')], ['user', load_lang('subscriber')], ['all', load_lang('all')]]
        else:
            acl_list = [['', load_lang('normal')], ['user', load_lang('subscriber')], ['admin', load_lang('admin')]]
        
        curs.execute("select dec from acl where title = ?", [name])
        acl_data = curs.fetchall()
        for data_list in acl_list:
            if acl_data and acl_data[0][0] == data_list[0]:
                check = 'selected="selected"'
            else:
                check = ''
            
            data += '<option value="' + data_list[0] + '" ' + check + '>' + data_list[1] + '</option>'
            
        data += '</select>'
        
        if not re.search('^user:', name):
            data += '<hr class=\"main_hr\">' + load_lang('discussion') + ' acl<br><br><select name="dis" ' + check_ok + '>'
        
            curs.execute("select dis, why, view from acl where title = ?", [name])
            acl_data = curs.fetchall()
            for data_list in acl_list:
                if acl_data and acl_data[0][0] == data_list[0]:
                    check = 'selected="selected"'
                else:
                    check = ''
                    
                data += '<option value="' + data_list[0] + '" ' + check + '>' + data_list[1] + '</option>'
                
            data += '</select>'

            data += '<hr class=\"main_hr\">' + load_lang('view') + ' acl<br><br><select name="view" ' + check_ok + '>'
            for data_list in acl_list:
                if acl_data and acl_data[0][2] == data_list[0]:
                    check = 'selected="selected"'
                else:
                    check = ''
                    
                data += '<option value="' + data_list[0] + '" ' + check + '>' + data_list[1] + '</option>'
                
            data += '</select>'
                
            if check_ok == '':
                if acl_data:
                    data += '<hr class=\"main_hr\"><input value="' + html.escape(acl_data[0][1]) + '" placeholder="' + load_lang('why') + '" name="why" type="text" ' + check_ok + '>'
                else:
                    data += '<hr class=\"main_hr\"><input placeholder="' + load_lang('why') + '" name="why" type="text" ' + check_ok + '>'
            
        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2([' (acl)', 0])],
            data =  '''
                    <form method="post">
                        ''' + data + '''
                        <hr class=\"main_hr\">
                        <button type="submit" ''' + check_ok + '''>acl ''' + load_lang('edit') + '''</button>
                    </form>
                    ''',
            menu = [['w/' + url_pas(name), load_lang('document')], ['manager', load_lang('admin')]]
        ))
            
@app.route('/admin/<name>', methods=['POST', 'GET'])
def user_admin(name = None):
    owner = admin_check()
    
    curs.execute("select acl from user where id = ?", [name])
    user = curs.fetchall()
    if not user:
        return re_error('/error/2')
    else:
        if owner != 1:
            curs.execute('select name from alist where name = ? and acl = "owner"', [user[0][0]])
            if curs.fetchall():
                return re_error('/error/3')

            if ip_check() == name:
                return re_error('/error/3')

    if flask.request.method == 'POST':
        if admin_check(7, 'admin (' + name + ')') != 1:
            return re_error('/error/3')

        if owner != 1:
            curs.execute('select name from alist where name = ? and acl = "owner"', [flask.request.form.get('select', None)])
            if curs.fetchall():
                return re_error('/error/3')

        if flask.request.form.get('select', None) == 'X':
            curs.execute("update user set acl = 'user' where id = ?", [name])
        else:
            curs.execute("update user set acl = ? where id = ?", [flask.request.form.get('select', None), name])
        
        conn.commit()
        
        return redirect('/admin/' + url_pas(name))            
    else:
        if admin_check(7) != 1:
            return re_error('/error/3')            

        div = '<option value="X">X</option>'
        
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
        
        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2([' (' + load_lang('authority') + ')', 0])],
            data =  '''
                    <form method="post">
                        <select name="select">''' + div + '''</select>
                        <hr class=\"main_hr\">
                        <button type="submit">''' + load_lang('edit') + '''</button>
                    </form>
                    ''',
            menu = [['manager', load_lang('admin')]]
        ))
    
@app.route('/diff/<everything:name>')
def diff_data(name = None):
    first = flask.request.args.get('first', '1')
    second = flask.request.args.get('second', '1')

    curs.execute("select data from history where id = ? and title = ?", [first, name])
    first_raw_data = curs.fetchall()
    if first_raw_data:
        curs.execute("select data from history where id = ? and title = ?", [second, name])
        second_raw_data = curs.fetchall()
        if second_raw_data:
            first_data = html.escape(first_raw_data[0][0])            
            second_data = html.escape(second_raw_data[0][0])

            if first == second:
                result = '-'
            else:            
                diff_data = difflib.SequenceMatcher(None, first_data, second_data)
                result = re.sub('\r', '', diff(diff_data))
            
            return easy_minify(flask.render_template(skin_check(), 
                imp = [name, wiki_set(), custom(), other2([' (' + load_lang('compare') + ')', 0])],
                data = '<pre>' + result + '</pre>',
                menu = [['history/' + url_pas(name), load_lang('history')]]
            ))

    return redirect('/history/' + url_pas(name))
        
@app.route('/down/<everything:name>')
def down(name = None):
    div = '<ul>'

    curs.execute("select title from data where title like ?", ['%' + name + '/%'])
    for data in curs.fetchall():
        div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a></li>'
        
    div += '</ul>'
    
    return easy_minify(flask.render_template(skin_check(), 
        imp = [name, wiki_set(), custom(), other2([' (' + load_lang('under') + ')', 0])],
        data = div,
        menu = [['w/' + url_pas(name), load_lang('document')]]
    ))

@app.route('/w/<everything:name>')
def read_view(name = None):
    data_none = 0
    sub = ''
    acl = ''
    div = ''

    num = flask.request.args.get('num', None)
    if num:
        num = int(num)
    else:
        if not flask.request.args.get('from', None):
            curs.execute("select title from back where link = ? and type = 'redirect'", [name])
            redirect_data = curs.fetchall()
            if redirect_data:
                return redirect('/w/' + redirect_data[0][0] + '?from=' + name)

    curs.execute("select sub from rd where title = ? and not stop = 'O' order by date desc", [name])
    if curs.fetchall():
        sub += ' (' + load_lang('discussion') + ')'

    curs.execute("select link from back where title = ? and type = 'cat' order by link asc", [name])
                
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
        
    if re.search('^category:', name):        
        curs.execute("select link from back where title = ? and type = 'cat' order by link asc", [name])
        back = curs.fetchall()
        if back:
            div = '<br><h2 id="cate_normal">' + load_lang('category') + '</h2><ul>'
            u_div = ''

            for data in back:    
                if re.search('^category:', data[0]):
                    u_div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a></li>'
                else:
                    curs.execute("select title from back where title = ? and type = 'include'", [data[0]])
                    db_data = curs.fetchall()
                    if db_data:
                        div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a> <a id="inside" href="/xref/' + url_pas(data[0]) + '">(' + load_lang('backlink') + ')</a></li>'
                    else: 
                        div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a></li>'

            div += '</ul>'
            
            if div == '<br><h2 id="cate_normal">' + load_lang('category') + '</h2><ul></ul>':
                div = ''
            
            if u_div != '':
                div += '<br><h2 id="cate_under">' + load_lang('under') + ' ' + load_lang('category') + '</h2><ul>' + u_div + '</ul>'


    if num:
        curs.execute("select title from history where title = ? and id = ? and hide = 'O'", [name, str(num)])
        if curs.fetchall() and admin_check(6) != 1:
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
        else_data = None

    m = re.search("^user:([^/]*)", name)
    if m:
        g = m.groups()
        
        curs.execute("select acl from user where id = ?", [g[0]])
        test = curs.fetchall()
        if test and test[0][0] != 'user':
            acl = ' (' + load_lang('admin') + ')'
        else:
            if ban_check(g[0]) == 1:
                sub += ' (' + load_lang('ban') + ')'
            else:
                acl = ''

    curs.execute("select dec from acl where title = ?", [name])
    data = curs.fetchall()
    if data:
        acl += ' (acl)'
            
    if flask.request.args.get('from', None) and else_data:
        else_data = re.sub('^\r\n', '', else_data)
        else_data = re.sub('\r\n$', '', else_data)
            
    end_data = render_set(
        title = name,
        data = else_data
    )

    if end_data == 'http request 401.3':
        response_data = 401
    
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

        menu += [['topic/' + url_pas(name), load_lang('discussion')], ['history/' + url_pas(name), load_lang('history')], ['xref/' + url_pas(name), load_lang('backlink')], ['acl/' + url_pas(name), 'acl']]

        if flask.request.args.get('from', None):
            menu += [['w/' + url_pas(name), load_lang('pass')]]
            end_data =  '''
                        <div id="redirect">
                            <a href="/w/''' + url_pas(flask.request.args.get('from', None)) + '?from=' + url_pas(name) + '">' + flask.request.args.get('from', None) + '</a> - ' + name + '''
                        </div>
                        <br>''' + end_data

        if uppage != 0:
            menu += [['w/' + url_pas(uppage), load_lang('upper')]]

        if down:
            menu += [['down/' + url_pas(name), load_lang('under')]]
    
        curs.execute("select date from history where title = ? order by date desc limit 1", [name])
        date = curs.fetchall()
        if date:
            r_date = date[0][0]
        else:
            r_date = 0

    div = end_data + div
            
    curs.execute("select data from other where name = 'adsense'")
    adsense_enabled = curs.fetchall()[0][0]
    adsense_code = '<div align="center" style="display: block; margin-bottom: 10px;">%_adsense_code_%</div>'
    if adsense_enabled == 'True':
        curs.execute("select data from other where name = 'adsense_code'")
        adsense_code = adsense_code.replace('%_adsense_code_%', curs.fetchall()[0][0])
    else:
        adsense_code = adsense_code.replace('%_adsense_code_%', '')
    curs.execute("select data from other where name = 'body'")
    body = curs.fetchall()
    if body:
        div = body[0][0] + '<hr class=\"main_hr\">' + div
    
    div = adsense_code + '<div>' + div + '</div>'
    return easy_minify(flask.render_template(skin_check(), 
        imp = [flask.request.args.get('show', name), wiki_set(), custom(), other2([sub + acl, r_date])],
        data = div,
        menu = menu
    )), response_data

@app.route('/topic_record/<name>')
def user_topic_list(name = None):
    num = int(flask.request.args.get('num', 1))
    if num * 50 > 0:
        sql_num = num * 50 - 50
    else:
        sql_num = 0
    
    one_admin = admin_check(1)

    div =   '''
            <table id="main_table_set">
                <tbody>
                    <tr>
                        <td id="main_table_width">''' + load_lang('discussion') + ' ' + load_lang('name') + '''</td>
                        <td id="main_table_width">''' + load_lang('writer') + '''</td>
                        <td id="main_table_width">''' + load_lang('time') + '''</td>
                    </tr>
            '''
    
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
    
    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('discussion') + ' ' + load_lang('record'), wiki_set(), custom(), other2([sub, 0])],
        data = div,
        menu = [['other', load_lang('other')], ['user', load_lang('user')], ['count/' + url_pas(name), load_lang('count')], ['record/' + url_pas(name), load_lang('record')]]
    ))

@app.route('/recent_changes')
@app.route('/<regex("record"):tool>/<name>')
@app.route('/<regex("history"):tool>/<everything:name>', methods=['POST', 'GET'])
def recent_changes(name = None, tool = 'record'):
    if flask.request.method == 'POST':
        return redirect('/diff/' + url_pas(name) + '?first=' + flask.request.form.get('b', None) + '&second=' + flask.request.form.get('a', None))
    else:
        one_admin = admin_check(1)
        six_admin = admin_check(6)
        
        ban = ''
        select = ''

        what = flask.request.args.get('what', 'all')

        div =   '''
                <table id="main_table_set">
                    <tbody>
                        <tr>
                '''
        
        if name:
            num = int(flask.request.args.get('num', 1))
            if num * 50 > 0:
                sql_num = num * 50 - 50
            else:
                sql_num = 0      

            if tool == 'history':
                div +=  '''
                        <td id="main_table_width">''' + load_lang('version') + '''</td>
                        <td id="main_table_width">''' + load_lang('editor') + '''</td>
                        <td id="main_table_width">''' + load_lang('time') + '''</td></tr>
                        '''
                
                curs.execute("select id, title, date, ip, send, leng from history where title = ? order by id + 0 desc limit ?, '50'", [name, str(sql_num)])
            else:
                div +=  '''
                            <td id="main_table_width">''' + load_lang('document') + ' ' + load_lang('name') + '''</td>
                            <td id="main_table_width">''' + load_lang('editor') + '''</td>
                            <td id="main_table_width">''' + load_lang('time') + '''</td>
                        </tr>
                        '''

                if what == 'all':
                    div = '<a href="/topic_record/' + url_pas(name) + '">(' + load_lang('discussion') + ')</a> <a href="/record/' + url_pas(name) + '?what=delete">(' + load_lang('delete') + ')</a> <a href="/record/' + url_pas(name) + '?what=move">(' + load_lang('move') + ')</a> <a href="/record/' + url_pas(name) + '?what=revert">(' + load_lang('revert') + ')</a><hr class=\"main_hr\">' + div
                    
                    curs.execute("select id, title, date, ip, send, leng from history where ip = ? order by date desc limit ?, '50'", [name, str(sql_num)])
                else:
                    if what == 'delete':
                        sql = '%(' + load_lang('delete', 1) + ')'
                    elif what == 'move':
                        sql = '%' + load_lang('move', 1) + ')'
                    elif what == 'revert':
                        sql = '%' + load_lang('version', 1) + ')'
                    else:
                        return redirect('/')

                    curs.execute("select id, title, date, ip, send, leng from history where ip = ? and send like ? order by date desc limit ?, 50", [name, sql, str(sql_num)])
        else:
            num = int(flask.request.args.get('num', 1))
            if num * 50 > 0:
                sql_num = num * 50 - 50
            else:
                sql_num = 0            
            
            div +=  '''
                        <td id="main_table_width">''' + load_lang('document') + ' ' + load_lang('name') + '''</td>
                        <td id="main_table_width">''' + load_lang('editor') + '''</td>
                        <td id="main_table_width">''' + load_lang('time') + '''</td>
                    </tr>
                    '''
            
            if what == 'all':
                div = '<a href="/recent_changes?what=delete">(' + load_lang('delete') + ')</a> <a href="/recent_changes?what=move">(' + load_lang('move') + ')</a> <a href="/recent_changes?what=revert">(' + load_lang('revert') + ')</a><hr class=\"main_hr\">' + div

                div = '<a href="/recent_discuss">(' + load_lang('discussion') + ')</a> <a href="/block_log">(' + load_lang('ban') + ')</a> <a href="/user_log">(' + load_lang('subscriber') + ')</a> <a href="/admin_log">(' + load_lang('authority') + ')</a><hr class=\"main_hr\">' + div
                
                curs.execute("select id, title, date, ip, send, leng from history where not title like 'user:%' order by date desc limit ?, 50", [str(sql_num)])
            else:
                if what == 'delete':
                    sql = '%(' + load_lang('delete', 1) + ')'
                elif what == 'move':
                    sql = '%' + load_lang('move', 1) + ')'
                elif what == 'revert':
                    sql = '%' + load_lang('version', 1) + ')'
                else:
                    return redirect('/')

                curs.execute("select id, title, date, ip, send, leng from history where send like ? and not title like 'user:%' order by date desc limit ?, 50", [sql, str(sql_num)])

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
                    hidden = ' <a href="/hidden/' + url_pas(data[1]) + '?num=' + data[0] + '">(' + load_lang('release') + ')'
                    
                    style[0] = 'id="toron_color_grey"'
                    style[1] = 'id="toron_color_grey"'
                    
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

                style[0] = 'style="display: none;"'
                style[1] = 'id="toron_color_grey"'

            if tool == 'history':
                title = '<a href="/w/' + url_pas(name) + '?num=' + data[0] + '">' + data[0] + load_lang('version') + '</a> <a href="/raw/' + url_pas(name) + '?num=' + data[0] + '">(' + load_lang('raw') + ')</a> '
            else:
                title = '<a href="/w/' + url_pas(data[1]) + '">' + html.escape(data[1]) + '</a> <a href="/history/' + url_pas(data[1]) + '">(' + data[0] + load_lang('version') + ')</a> '
                    
            div +=  '''
                    <tr ''' + style[0] + '''>
                        <td>''' + title + revert + ' ' + leng + '''</td>
                        <td>''' + ip + ban + hidden + '''</td>
                        <td>''' + date + '''</td>
                    </tr>
                    <tr ''' + style[1] + '''>
                        <td colspan="3">''' + send_parser(send) + '''</td>
                    </tr>
                    '''

        div +=  '''
                    </tbody>
                </table>
                '''
        sub = ''

        if name:
            if tool == 'history':
                div =   '''
                        <form method="post">
                            <select name="a">''' + select + '''</select> <select name="b">''' + select + '''</select>
                            <button type="submit">''' + load_lang('compare') + '''</button>
                        </form>
                        <hr class=\"main_hr\">
                        <a href="/move_data/''' + url_pas(name) + '''">(''' + load_lang('move') + ''')</a>
                        <hr class=\"main_hr\">
                        ''' + div
                title = name
                
                sub += ' (' + load_lang('history') + ')'
                
                menu = [['w/' + url_pas(name), load_lang('document')], ['raw/' + url_pas(name), 'raw']]
                
                div += next_fix('/history/' + url_pas(name) + '?num=', num, data_list)
            else:
                curs.execute("select end from ban where block = ?", [name])
                if curs.fetchall():
                    sub += ' (' + load_lang('ban') + ')'

                title = load_lang('edit') + ' ' + load_lang('record')
                
                menu = [['other', load_lang('other')], ['user', load_lang('user')], ['count/' + url_pas(name), load_lang('count')]]
                
                div += next_fix('/record/' + url_pas(name) + '/' + url_pas(what) + '?num=', num, data_list)
                
                if what != 'all':
                    menu += [['record/' + url_pas(name), load_lang('normal')]]
        else:
            menu = 0
            title = load_lang('recent') + ' ' + load_lang('change') + ''
            
            if what != 'all':
                menu = [['recent_changes', load_lang('normal')]]
                
            div += next_fix('/recent_changes?num=', num, data_list)
                
        if what == 'delete':
            sub += ' (' + load_lang('delete') + ')'
        elif what == 'move':
            sub += ' (' + load_lang('move') + ')'
        elif what == 'revert':
            sub += ' (' + load_lang('revert') + ')'
        
        if sub == '':
            sub = 0
                
        return easy_minify(flask.render_template(skin_check(), 
            imp = [title, wiki_set(), custom(), other2([sub, 0])],
            data = div,
            menu = menu
        ))
    
@app.route('/upload', methods=['GET', 'POST'])
def upload():
    if ban_check() == 1:
        return re_error('/ban')
    
    if flask.request.method == 'POST':
        if captcha_post(flask.request.form.get('g-recaptcha-response', '')) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        data = flask.request.files.get('f_data', None)
        if not data:
            return re_error('/error/9')

        if int(wiki_set(3)) * 1024 * 1024 < flask.request.content_length:
            return re_error('/error/17')
        
        value = os.path.splitext(data.filename)[1]
        if not value in ['.jpeg', '.jpg', '.gif', '.png', '.webp', '.JPEG', '.JPG', '.GIF', '.PNG', '.WEBP']:
            return re_error('/error/14')
    
        if flask.request.form.get('f_name', None):
            name = flask.request.form.get('f_name', None) + value
        else:
            name = data.filename
        
        piece = os.path.splitext(name)
        if re.search('[^ㄱ-힣0-9a-zA-Z_\- ]', piece[0]):
            return re_error('/error/22')

        e_data = sha224(piece[0]) + piece[1]

        curs.execute("select title from data where title = ?", ['file:' + name])
        if curs.fetchall():
            return re_error('/error/16')
            
        ip = ip_check()

        if flask.request.form.get('f_lice', None):
            lice = flask.request.form.get('f_lice', None)
        else:
            if custom()[2] == 0:
                lice = ip + ' ' + load_lang('upload', 1)
            else:
                lice = '[[user:' + ip + ']] ' + load_lang('upload', 1)
            
        if os.path.exists(os.path.join('image', e_data)):
            os.remove(os.path.join('image', e_data))
            
            data.save(os.path.join('image', e_data))
        else:
            data.save(os.path.join('image', e_data))
            
        curs.execute("select title from data where title = ?", ['file:' + name])
        if curs.fetchall(): 
            curs.execute("delete from data where title = ?", ['file:' + name])
        
        curs.execute("insert into data (title, data) values (?, ?)", ['file:' + name, '[[file:' + name + ']][br][br]{{{[[file:' + name + ']]}}}[br][br]' + lice])
        curs.execute("insert into acl (title, dec, dis, why, view) values (?, 'admin', '', '', '')", ['file:' + name])

        history_plus(
            'file:' + name, '[[file:' + name + ']][br][br]{{{[[file:' + name + ']]}}}[br][br]' + lice,
            get_time(), 
            ip, 
            load_lang('upload', 1), 
            '0'
        )
        
        conn.commit()
        
        return redirect('/w/file:' + name)      
    else:
        return easy_minify(flask.render_template(skin_check(), 
            imp = [load_lang('upload'), wiki_set(), custom(), other2([0, 0])],
            data =  '''
                    <form method="post" enctype="multipart/form-data" accept-charset="utf8">
                        <input type="file" name="f_data">
                        <hr class=\"main_hr\">
                        <input placeholder="''' + load_lang('name') + '''" name="f_name" type="text">
                        <hr class=\"main_hr\">
                        <input placeholder="''' + load_lang('license') + '''" name="f_lice" type="text">
                        <hr class=\"main_hr\">
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
            match = '-'

        curs.execute("select end, login, band from ban where block = ? or block = ?", [ip, match])
        block_data = curs.fetchall()
        if block_data:
            if block_data[0][0] != '':
                acl += ' (end : ' + block_data[0][0] + ')'
            else:
                acl += ' (' + load_lang('limitless') + ')'        

            if block_data[0][1] != '':
                acl += ' (' + load_lang('login') + ' ' + load_lang('able') + ')'

            if block_data[0][2] == 'O':
                acl += ' (' + load_lang('band') + ')'
            
    if custom()[2] != 0:
        ip_user = '<a href="/w/user:' + ip + '">' + ip + '</a>'
        
        plus =  '''
                <li><a href="/logout">''' + load_lang('logout') + '''</a></li>
                <li><a href="/change">''' + load_lang('user') + ' ' + load_lang('setting') + ' ' + load_lang('edit') + '''</a></li>
                '''
        
        curs.execute('select name from alarm where name = ? limit 1', [ip_check()])
        if curs.fetchall():
            plus2 = '<li><a href="/alarm">' + load_lang('alarm') + ' (O)</a></li>'
        else:
            plus2 = '<li><a href="/alarm">' + load_lang('alarm') + '</a></li>'

        plus2 += '<li><a href="/watch_list">' + load_lang('watchlist') + '</a></li>'
        plus3 = '<li><a href="/acl/user:' + url_pas(ip) + '">' + load_lang('user') + ' ' + load_lang('document') + ' acl</a></li>'
    else:
        ip_user = ip
        
        plus =  '''
                <li><a href="/login">''' + load_lang('login') + '''</a></li>
                <li><a href="/register">''' + load_lang('register') + '''</a></li>
                '''
        plus2 = ''
        plus3 = ''

        curs.execute("select data from other where name = 'email_have'")
        test = curs.fetchall()
        if test and test[0][0] != '':
            plus += '<li><a href="/pass_find">' + load_lang('password') + ' ' + load_lang('search') + '</a></li>'

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('user') + ' ' + load_lang('tool'), wiki_set(), custom(), other2([0, 0])],
        data =  '''
                <h2>''' + load_lang('state') + '''</h2>
                <ul>
                    <li>''' + ip_user + ''' <a href="/record/''' + url_pas(ip) + '''">(''' + load_lang('record') + ''')</a></li>
                    <li>''' + load_lang('authority') + ''' : ''' + acl + '''</li>
                </ul>
                <br>
                <h2>''' + load_lang('login') + '''</h2>
                <ul>
                    ''' + plus + '''
                </ul>
                <br>
                <h2>''' + load_lang('tool') + '''</h2>
                <ul>
                    ''' + plus3 + '''
                    <li><a href="/custom_head">''' + load_lang('user') + ''' head</a></li>
                </ul>
                <br>
                <h2>''' + load_lang('other') + '''</h2>
                <ul>
                ''' + plus2 + '''
                <li>
                    <a href="/count">''' + load_lang('count') + '''</a>
                </li>
                </ul>
                ''',
        menu = 0
    ))

@app.route('/watch_list')
def watch_list():
    div = 'limit : 10<hr class=\"main_hr\">'
    
    if custom()[2] == 0:
        return redirect('/login')

    curs.execute("select title from scan where user = ?", [ip_check()])
    data = curs.fetchall()
    for data_list in data:
        div += '<li><a href="/w/' + url_pas(data_list[0]) + '">' + data_list[0] + '</a> <a href="/watch_list/' + url_pas(data_list[0]) + '">(' + load_lang('delete') + ')</a></li>'

    if data:
        div = '<ul>' + div + '</ul><hr class=\"main_hr\">'

    div += '<a href="/manager/13">(' + load_lang('plus') + ')</a>'

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('watchlist') + ' ' + load_lang('list'), wiki_set(), custom(), other2([0, 0])],
        data = div,
        menu = [['manager', load_lang('admin')]]
    ))

@app.route('/watch_list/<everything:name>')
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

    if flask.request.method == 'POST':
        if custom()[2] != 0:
            curs.execute("select user from custom where user = ?", [ip + ' (head)'])
            if curs.fetchall():
                curs.execute("update custom set css = ? where user = ?", [flask.request.form.get('content', None), ip + ' (head)'])
            else:
                curs.execute("insert into custom (user, css) values (?, ?)", [ip + ' (head)', flask.request.form.get('content', None)])
            
            conn.commit()

        flask.session['head'] = flask.request.form.get('content', None)

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
            start = '<span>' + load_lang('user_head_warring') + '</span><hr class=\"main_hr\">'
            
            if 'head' in flask.session:
                data = flask.session['head']
            else:
                data = ''

        start += '<span>&lt;style&gt;css&lt;/style&gt;<br>&lt;script&gt;js&lt;/script&gt;</span><hr class=\"main_hr\">'

        return easy_minify(flask.render_template(skin_check(), 
            imp = [load_lang('user') + ' head', wiki_set(), custom(), other2([0, 0])],
            data =  start + '''
                    <form method="post">
                        <textarea rows="25" cols="100" name="content">''' + data + '''</textarea>
                        <hr class=\"main_hr\">
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

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('count'), wiki_set(), custom(), other2([0, 0])],
        data =  '''
                <ul>
                    <li><a href="/record/''' + url_pas(that) + '''">''' + load_lang('edit') + '''</a> : ''' + str(data) + '''</li>
                    <li><a href="/topic_record/''' + url_pas(that) + '''">''' + load_lang('discussion') + '''</a> : ''' + str(t_data) + '''</a></li>
                </ul>
                ''',
        menu = [['user', load_lang('user')]]
    ))
        
@app.route('/random')
def title_random():
    curs.execute("select title from data order by random() limit 1")
    data = curs.fetchall()
    if data:
        return redirect('/w/' + url_pas(data[0][0]))
    else:
        return redirect('/')

@app.route('/skin_set')
def skin_set():
    return re_error('/error/5')
    
@app.route('/api/w/<everything:name>')
def api_w(name = ''):
    curs.execute("select data from data where title = ?", [name])
    data = curs.fetchall()
    if data:
        json_data = { "title" : name, "data" : render_set(title = name, data = data[0][0]) }
    
        return flask.jsonify(json_data)
    else:
        return flask.jsonify({})
    
@app.route('/api/raw/<everything:name>')
def api_raw(name = ''):
    curs.execute("select data from data where title = ?", [name])
    data = curs.fetchall()
    if data:
        json_data = { "title" : name, "data" : data[0][0] }
    
        return flask.jsonify(json_data)
    else:
        return flask.jsonify({})

@app.route('/api/topic/<everything:name>/sub/<sub>')
def api_topic_sub(name = '', sub = '', time = ''):
    if flask.request.args.get('time', None):
        curs.execute("select id, data, ip from topic where title = ? and sub = ? and date >= ? order by id + 0 asc", [name, sub, flask.request.args.get('time', None)])
    else:
        curs.execute("select id, data, ip from topic where title = ? and sub = ? order by id + 0 asc", [name, sub])
    data = curs.fetchall()
    if data:
        json_data = {}
        for i in data:
            json_data[i[0]] =   {
                                    "data" : i[1],
                                    "id" : i[2]
                                }

        return flask.jsonify(json_data)
    else:
        return flask.jsonify({})
    
@app.route('/views/<everything:name>')
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
        return easy_minify(flask.send_from_directory('./views' + plus, rename), 'css')   
    elif g == 'js':
        return easy_minify(flask.send_from_directory('./views' + plus, rename), 'js')
    elif g == 'html':
        return easy_minify(flask.send_from_directory('./views' + plus, rename))   
    else:
        return flask.send_from_directory('./views' + plus, rename)

@app.route('/<data>')
def main_file(data = None):
    if re.search('\.txt$', data):
        return flask.send_from_directory('./', data)
    else:
        return redirect('/w/' + url_pas(wiki_set(2)))

@app.errorhandler(404)
def error_404(e):
    return redirect('/w/' + url_pas(wiki_set(2)))

if __name__=="__main__":
    app.secret_key = rep_key
    http_server = tornado.httpserver.HTTPServer(tornado.wsgi.WSGIContainer(app))
    http_server.listen(rep_port, address=rep_host)
    tornado.ioloop.IOLoop.instance().start()
