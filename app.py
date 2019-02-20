import os
import re

for i_data in os.listdir("route"):
    f_src = re.search("(.+)\.py$", i_data)
    if f_src:
        f_src = f_src.groups()[0]

        exec("from route." + f_src + " import *")

r_ver = 'v3.0.9-master-005'
c_ver = '309002'

print('Version : ' + r_ver)

app_var = json.loads(open('data/app_variables.json', encoding='utf-8').read())

all_src = []
for i_data in os.listdir("."):
    f_src = re.search("(.+)\.db$", i_data)
    if f_src:
        all_src += [f_src.groups()[0]]

if len(all_src) == 0:
    print('DB\'s name (data) : ', end = '')
    
    db_name = input()
    if db_name == '':
        db_name = 'data'
elif len(all_src) > 1:
    db_num = 1

    for i_data in all_src:
        print(str(db_num) + ' : ' + i_data)

        db_num += 1

    print('Number : ', end = '')    
    db_name = all_src[int(number_check(input())) - 1]
else:
    db_name = all_src[0]

if len(all_src) == 1:
    print('DB\'s name : ' + db_name)
            
if os.path.exists(db_name + '.db'):
    setup_tool = 0
else:
    setup_tool = 1

conn = sqlite3.connect(db_name + '.db', check_same_thread = False)
curs = conn.cursor()

load_conn(conn)

logging.basicConfig(level = logging.ERROR)

app = flask.Flask(__name__, template_folder = './')
app.config['JSON_AS_ASCII'] = False

flask_reggie.Reggie(app)

compress = flask_compress.Compress()
compress.init_app(app)

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
    try:
        curs.execute('select data from other where name = "ver"')
        ver_set_data = curs.fetchall()
        if not ver_set_data:
            setup_tool = 1
        else:
            if c_ver > ver_set_data[0][0]:
                setup_tool = 1
    except:
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

if not os.path.exists(app_var['PATH_DATA_IMAGES']):
    os.makedirs(app_var['PATH_DATA_IMAGES'])
    
if not os.path.exists('views'):
    os.makedirs('views')

import route.tool.init as server_init

dislay_set_key = ['Host', 'Port', 'Language', 'Markup', 'Encrypt Method']
server_set_key = ['host', 'port', 'language', 'markup', 'encode']
server_set = {}

for i in range(len(server_set_key)):
    curs.execute('select data from other where name = ?', [server_set_key[i]])
    server_set_val = curs.fetchall()
    if not server_set_val:
        server_set_val = server_init.init(server_set_key[i])
        
        curs.execute('insert into other (name, data) values (?, ?)', [server_set_key[i], server_set_val])
        conn.commit()
    else:
        server_set_val = server_set_val[0][0]
    
    print(dislay_set_key[i] + ' : ' + server_set_val)
    
    server_set[server_set_key[i]] = server_set_val

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
        
        print('Engine made robots.txt')
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


curs.execute('select data from other where name = "adsense"')
adsense_result = curs.fetchall()
if not adsense_result:
    curs.execute('insert into other (name, data) values ("adsense", "False")')
    curs.execute('insert into other (name, data) values ("adsense_code", "")')


curs.execute('delete from other where name = "ver"')
curs.execute('insert into other (name, data) values ("ver", ?)', [c_ver])

def back_up():
    try:
        shutil.copyfile(db_name + '.db', 'back_' + db_name + '.db')
        
        print('Back up : OK')
    except:
        print('Back up : Error')

    threading.Timer(60 * 60 * back_time, back_up).start()

try:
    curs.execute('select data from other where name = "back_up"')
    back_up_time = curs.fetchall()
    
    back_time = int(back_up_time[0][0])
except:
    back_time = 0
    
if back_time != 0:
    print('Back up state : ' + str(back_time) + ' hours')
    
    if __name__ == '__main__':
        back_up()
else:
    print('Back up state : Turn off')

conn.commit()

@app.route('/del_alarm')
def del_alarm():
    return del_alarm_2(conn)

@app.route('/alarm')
def alarm():
    return alarm_2(conn)

@app.route('/<regex("inter_wiki|(?:edit|email|name)_filter"):tools>')
def inter_wiki(tools = None):
    return inter_wiki_2(conn, tools)

@app.route('/<regex("del_(?:inter_wiki|(?:edit|email|name)_filter)"):tools>/<name>')
def del_inter(tools = None, name = None):
    return del_inter_2(conn, tools, name)

@app.route('/<regex("plus_(?:inter_wiki|(?:edit|email|name)_filter)"):tools>', methods=['POST', 'GET'])
@app.route('/<regex("plus_edit_filter"):tools>/<name>', methods=['POST', 'GET'])
def plus_inter(tools = None, name = None):
    return plus_inter_2(conn, tools, name)

@app.route('/setting')
@app.route('/setting/<int:num>', methods=['POST', 'GET'])
def setting(num = 0):
    return setting_2(conn, num)

@app.route('/not_close_topic')
def not_close_topic():
    return not_close_topic_2(conn)

@app.route('/image/<name>')
def image_view(name = None):
    return image_view_2(conn, name)

@app.route('/acl_list')
def acl_list():
    return acl_list_2(conn)

@app.route('/admin_plus/<name>', methods=['POST', 'GET'])
def admin_plus(name = None):
    return admin_plus_2(conn)
        
@app.route('/admin_list')
def admin_list():
    return admin_list_2(conn)
        
@app.route('/hidden/<everything:name>')
def history_hidden(name = None):
    return history_hidden_2(name)
        
@app.route('/user_log')
def user_log():
    return user_log_2(conn)

@app.route('/admin_log')
def admin_log():
    num = int(number_check(flask.request.args.get('num', '1')))
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
        imp = [load_lang('authority_use_list'), wiki_set(), custom(), other2([0, 0])],
        data = list_data,
        menu = [['other', load_lang('return')]]
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
    
    list_data += '</ul><hr class=\"main_hr\"><a href="/manager/8">(' + load_lang('add') + ')</a>'

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('admin_group_list'), wiki_set(), custom(), other2([0, 0])],
        data = list_data,
        menu = [['other', load_lang('return')]]
    ))

@app.route('/indexing', methods=['POST', 'GET'])
def indexing():
    if admin_check() != 1:
        return re_error('/error/3')

    if flask.request.method == 'POST':
        admin_check(None, 'indexing')

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
        
        return redirect()  
    else:
        curs.execute("select name from sqlite_master where type = 'index'")
        data = curs.fetchall()
        if data:
            b_data = load_lang('delete')
        else:
            b_data = load_lang('create')

        return easy_minify(flask.render_template(skin_check(), 
            imp = [load_lang('indexing'), wiki_set(), custom(), other2([0, 0])],
            data =  '''
                    <form method="post">
                        <button type="submit">''' + b_data + '''</button>
                    </form>
                    ''',
            menu = [['manager', load_lang('return')]]
        ))       

   

@app.route('/restart', methods=['POST', 'GET'])
def restart():
    if admin_check(None, 'restart') != 1:
        return re_error('/error/3')

    if flask.request.method == 'POST':
        os.execl(sys.executable, sys.executable, *sys.argv)
    else:
        print('Restart')

        return easy_minify(flask.render_template(skin_check(), 
            imp = [load_lang('wiki_restart'), wiki_set(), custom(), other2([0, 0])],
            data =  '''
                    <form method="post">
                        <button type="submit">''' + load_lang('restart') + '''</button>
                    </form>
                    ''',
            menu = [['manager', load_lang('return')]]
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
        print('Update')

        os.system('git remote rm origin')
        os.system('git remote add origin https://github.com/2DU/opennamu.git')
        ok = os.system('git fetch origin ' + up_data)
        ok = os.system('git reset --hard origin/' + up_data)
        if ok == 0:
            return redirect('/restart')
    else:
        if platform.system() == 'Windows':
            print('Update')

            urllib.request.urlretrieve('https://github.com/2DU/opennamu/archive/' + up_data + '.zip', 'update.zip')
            zipfile.ZipFile('update.zip').extractall('')
            ok = os.system('xcopy /y /r opennamu-' + up_data + ' .')
            if ok == 0:
                print('Remove')
                os.system('rd /s /q opennamu-' + up_data)
                os.system('del update.zip')

                return redirect('/restart')

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('update'), wiki_set(), custom(), other2([0, 0])],
        data = load_lang("update_error") + ' <a href="https://github.com/2DU/opennamu">(Github)</a>',
        menu = [['manager/1', load_lang('return')]]
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
                data = '<h2>ie_no_data_required</h2>' + load_lang('ie_no_data_required'),
                menu = [['other', load_lang('return')]]
            ))

        with open(app_var['PATH_OAUTHSETTINGS'], 'r', encoding='utf-8') as f:
            legacy = json.loads(f.read())

        with open(app_var['PATH_OAUTHSETTINGS'], 'w', encoding='utf-8') as f:
            f.write("""
                {
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
                }
            """)
        
        return flask.redirect('/oauth_settings')

    oauth_supported = load_oauth('_README')['support']

    body_content = ''
    body_content += '''
        <script>
            function check_value (target) {
                target_box = document.getElementById(target.id + "_box");
                if (target.value !== "") {
                    target_box.checked = true;
                } else {
                    target_box.checked = false;
                } 
            }
        </script>
    '''

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
        imp = [load_lang('oauth_setting'), wiki_set(), custom(), other2([0, 0])],
        data = body_content,
        menu = [['other', load_lang('return')]]
    ))

@app.route('/adsense_settings', methods=['GET', 'POST'])
def adsense_settings():
    if admin_check(None, 'adsense_settings') != 1:
        return re_error('/error/3')
    
    if flask.request.method == 'POST':
        try:
            adsense_enabled = flask.request.form.get('adsense_enabled')
            adsense_code = flask.request.form['adsense_code']
        except:
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('inter_error'), wiki_set(), custom(), other2([0, 0])],
                data = '<h2>ie_no_data_required</h2>' + load_lang('ie_no_data_required'),
                menu = [['other', load_lang('return')]]
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
                    <input class="form-check-input" name="adsense_enabled" type="checkbox" {}>
                    {}
                </label>
            </div>
            <hr>
            <div class="form-group">
                <textarea class="form-control" id="adsense_code" name="adsense_code" rows="12">{}</textarea>
            </div>
            <button type="submit" value="publish">{}</button>
        </form>
    '''
    
    body_content += template.format(
        'checked' if adsense_enabled == 'True' else template.format(''),
        load_lang('adsense_enable'),
        load_lang('save'),
        adsense_code
    )

    return easy_minify(flask.render_template(skin_check(),
        imp = [load_lang('adsense_setting'), wiki_set(), custom(), other2([0, 0])],
        data = body_content,
        menu = [['other', load_lang('return')]]
    ))
        
@app.route('/xref/<everything:name>')
def xref(name = None):
    num = int(number_check(flask.request.args.get('num', '1')))
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
        menu = [['w/' + url_pas(name), load_lang('return')]]
    ))

@app.route('/please')
def please():
    num = int(number_check(flask.request.args.get('num', '1')))
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
        imp = [load_lang('need_document'), wiki_set(), custom(), other2([0, 0])],
        data = div,
        menu = [['other', load_lang('return')]]
    ))
        
@app.route('/recent_discuss')
def recent_discuss():
    div = ''
    
    if flask.request.args.get('what', 'normal') == 'normal':
        div += '<a href="/recent_discuss?what=close">(' + load_lang('close_discussion') + ')</a>'
       
        m_sub = 0
    else:
        div += '<a href="/recent_discuss">(' + load_lang('open_discussion') + ')</a>'
        
        m_sub = ' (' + load_lang('closed') + ')'

    div +=  '''
            <hr class=\"main_hr\">
            <table id="main_table_set">
                <tbody>
                    <tr>
                        <td id="main_table_width_half">''' + load_lang('discussion_name') + '''</td>
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
        imp = [load_lang('recent_discussion'), wiki_set(), custom(), other2([m_sub, 0])],
        data = div,
        menu = 0
    ))

@app.route('/block_log')
@app.route('/<regex("block_user|block_admin"):tool>/<name>')
def block_log(name = None, tool = None):
    num = int(number_check(flask.request.args.get('num', '1')))
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
        div =   '''
                <a href="/manager/11">(''' + load_lang('blocked') + ''')</a> <a href="/manager/12">(''' + load_lang('admin') + ''')</a>
                <hr class=\"main_hr\">
                ''' + div
        
        sub = 0
        menu = 0
        
        curs.execute("select why, block, blocker, end, today from rb order by today desc limit ?, '50'", [str(sql_num)])
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
            ip = data[1] + ' (' + load_lang('range') + ')'
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
        '''

    div += '</tbody></table>'
    
    if not name:
        div += next_fix('/block_log?num=', num, data_list)
    else:
        div += next_fix('/' + url_pas(tool) + '/' + url_pas(name) + '?num=', num, data_list)
                
    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('recent_ban'), wiki_set(), custom(), other2([sub, 0])],
        data = div,
        menu = menu
    ))
            
@app.route('/search', methods=['POST'])
def search():
    return redirect('/search/' + url_pas(flask.request.form.get('search', 'test')))

@app.route('/goto', methods=['POST'])
def goto():
    curs.execute("select title from data where title = ?", [flask.request.form.get('search', 'test')])
    data = curs.fetchall()
    if data:
        return redirect('/w/' + url_pas(flask.request.form.get('search', 'test')))
    else:
        return redirect('/search/' + url_pas(flask.request.form.get('search', 'test')))

@app.route('/search/<everything:name>')
def deep_search(name = ''):
    if name == '':
        return redirect()

    num = int(number_check(flask.request.args.get('num', '1')))
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
            num = int(number_check(num))
    
    if not sub_title and num:
        curs.execute("select title from history where title = ? and id = ? and hide = 'O'", [name, str(num)])
        if curs.fetchall() and admin_check(6) != 1:
            return re_error('/error/3')
        
        curs.execute("select data from history where title = ? and id = ?", [name, str(num)])
        
        sub += ' (r' + str(num) + ')'

        menu = [['history/' + url_pas(name), load_lang('history')]]
    elif sub_title:
        curs.execute("select data from topic where id = ? and title = ? and sub = ? and block = ''", [str(num), name, sub_title])
        
        v_name = load_lang('discussion_raw')
        sub = ' (' + str(num) + ')'

        menu = [['topic/' + url_pas(name) + '/sub/' + url_pas(sub_title) + '#' + str(num), load_lang('discussion')], ['topic/' + url_pas(name) + '/sub/' + url_pas(sub_title) + '/admin/' + str(num), load_lang('return')]]
    else:
        curs.execute("select data from data where title = ?", [name])
        
        menu = [['w/' + url_pas(name), load_lang('return')]]

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
    num = int(number_check(flask.request.args.get('num', '1')))

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
                flask.request.form.get('send', None) + ' (r' + str(num) + ')', 
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
                        <span>r''' + flask.request.args.get('num', '0') + '''</span>
                        <hr class=\"main_hr\">
                        ''' + ip_warring() + '''
                        <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                        <hr class=\"main_hr\">
                        ''' + captcha_get() + '''
                        <button type="submit">''' + load_lang('revert') + '''</button>
                    </form>
                    ''',
            menu = [['history/' + url_pas(name), load_lang('history')], ['recent_changes', load_lang('recent_change')]]
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
                        <a href="/manager/15?plus=''' + url_pas(name) + '">(' + load_lang('load') + ')</a> <a href="/edit_filter">(' + load_lang('edit_filter_rule') + ''')</a>
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
            menu = [['w/' + url_pas(name), load_lang('return')], ['delete/' + url_pas(name), load_lang('delete')], ['move/' + url_pas(name), load_lang('move')]]
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
        data =  '<a href="/edit_filter">(' + load_lang('edit_filter_rule') + ')</a>' + js_data[0] + '''
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
        menu = [['w/' + url_pas(name), load_lang('return')]]
    ))
        
@app.route('/delete/<everything:name>', methods=['POST', 'GET'])
def delete(name = None):
    return delete_2(conn, name)        
            
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
                    flask.request.form.get('send', None) + ' (marge <a>' + name + '</a> - <a>' + flask.request.form.get('title', None) + '</a> move)', 
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
                flask.request.form.get('send', None) + ' (<a>' + name + '</a> - <a>' + flask.request.form.get('title', None) + '</a> move)', 
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
                        <input placeholder="''' + load_lang('document_name') + '" value="' + name + '''" name="title" type="text">
                        <hr class=\"main_hr\">
                        <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                        <hr class=\"main_hr\">
                        ''' + captcha_get() + '''
                        <button type="submit">''' + load_lang('move') + '''</button>
                    </form>
                    ''',
            menu = [['w/' + url_pas(name), load_lang('return')]]
        ))

@app.route('/other')
def other():
    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('other_tool'), wiki_set(), custom(), other2([0, 0])],
        data =  '''
                <h2>''' + load_lang('record') + '''</h2>
                <ul>
                    <li><a href="/manager/6">''' + load_lang('edit_record') + '''</a></li>
                    <li><a href="/manager/7">''' + load_lang('discussion_record') + '''</a></li>
                </ul>
                <br>
                <h2>''' + load_lang('list') + '''</h2>
                <ul>
                    <li><a href="/admin_list">''' + load_lang('admin_list') + '''</a></li>
                    <li><a href="/give_log">''' + load_lang('admin_group_list') + '''</a></li>
                    <li><a href="/not_close_topic">''' + load_lang('open_discussion_list') + '''</a></li>
                    <li><a href="/title_index">''' + load_lang('all_document_list') + '''</a></li>
                    <li><a href="/acl_list">''' + load_lang('acl_document_list') + '''</a></li>
                    <li><a href="/please">''' + load_lang('need_document') + '''</a></li>
                    <li><a href="/block_log">''' + load_lang('recent_ban') + '''</a></li>
                    <li><a href="/user_log">''' + load_lang('member_list') + '''</a></li>
                    <li><a href="/admin_log">''' + load_lang('authority_use_list') + '''</a></li>
                </ul>
                <br>
                <h2>''' + load_lang('other') + '''</h2>
                <ul>
                    <li><a href="/upload">''' + load_lang('upload') + '''</a></li>
                    <li><a href="/manager/10">''' + load_lang('search') + '''</a></li>
                </ul>
                <br>
                <h2>''' + load_lang('admin') + '''</h2>
                <ul>
                    <li><a href="/manager/1">''' + load_lang('admin_tool') + '''</a></li>
                </ul>
                <br>
                <h2>''' + load_lang('version') + '''</h2>
                <ul>
                    <li>''' + load_lang('version') + ' : <a id="out_link" href="https://github.com/2DU/opennamu/blob/master/version.md">' + r_ver + '''</a></li>
                </ul>
                ''',
        menu = 0
    ))
    
@app.route('/manager', methods=['POST', 'GET'])
@app.route('/manager/<int:num>', methods=['POST', 'GET'])
def manager(num = 1):
    title_list = {
        0 : [load_lang('document_name'), 'acl'], 
        1 : [0, 'check'], 
        2 : [0, 'ban'], 
        3 : [0, 'admin'], 
        4 : [0, 'record'], 
        5 : [0, 'topic_record'], 
        6 : [load_lang('name'), 'admin_plus'], 
        7 : [load_lang('name'), 'plus_edit_filter'], 
        8 : [load_lang('document_name'), 'search'], 
        9 : [0, 'block_user'], 
        10 : [0, 'block_admin'], 
        11 : [load_lang('document_name'), 'watch_list'], 
        12 : [load_lang('compare_target'), 'check'], 
        13 : [load_lang('document_name'), 'edit']
    }
    
    if num == 1:
        return easy_minify(flask.render_template(skin_check(), 
            imp = [load_lang('admin_tool'), wiki_set(), custom(), other2([0, 0])],
            data =  '''
                    <h2>''' + load_lang('admin') + '''</h2>
                    <ul>
                        <li><a href="/manager/2">''' + load_lang('acl_document_list') + '''</a></li>
                        <li><a href="/manager/3">''' + load_lang('check_user') + '''</a></li>
                        <li><a href="/manager/4">''' + load_lang('ban') + '''</a></li>
                        <li><a href="/manager/5">''' + load_lang('authorize') + '''</a></li>
                        <li><a href="/edit_filter">''' + load_lang('edit_filter_list') + '''</a></li>
                    </ul>
                    <br>
                    <h2>''' + load_lang('owner') + '''</h2>
                    <ul>
                        <li><a href="/manager/8">''' + load_lang('admin_group_add') + '''</a></li>
                        <li><a href="/setting">''' + load_lang('setting') + '''</a></li>
                    </ul>
                    <h3>''' + load_lang('filter') + '''</h3>
                    <ul>
                        <li><a href="/inter_wiki">''' + load_lang('interwiki_list') + '''</a></li>
                        <li><a href="/email_filter">''' + load_lang('email_filter_list') + '''</a></li>
                        <li><a href="/name_filter">''' + load_lang('id_filter_list') + '''</a></li>
                    </ul>
                    <br>
                    <h2>''' + load_lang('server') + '''</h2>
                    <ul>
                        <li><a href="/indexing">''' + load_lang('indexing') + '''</a></li>
                        <li><a href="/restart">''' + load_lang('wiki_restart') + '''</a></li>
                        <li><a href="/update">''' + load_lang('update') + '''</a></li>
                        <li><a href="/oauth_settings">''' + load_lang('oauth_setting') + '''</a></li>
                        <li><a href="/adsense_settings">''' + load_lang('adsense_setting') + '''</a></li>
                    </ul>
                    ''',
            menu = [['other', load_lang('return')]]
        ))
    elif not num - 1 > len(title_list):
        if flask.request.method == 'POST':
            if flask.request.args.get('plus', None):
                return redirect('/' + title_list[(num - 2)][1] + '/' + url_pas(flask.request.args.get('plus', None)) + '?plus=' + flask.request.form.get('name', None))
            else:
                return redirect('/' + title_list[(num - 2)][1] + '/' + url_pas(flask.request.form.get('name', None)))
        else:
            if title_list[(num - 2)][0] == 0:
                placeholder = load_lang('user_name')
            else:
                placeholder = title_list[(num - 2)][0]

            return easy_minify(flask.render_template(skin_check(), 
                imp = ['Redirect', wiki_set(), custom(), other2([0, 0])],
                data =  '''
                        <form method="post">
                            <input placeholder="''' + placeholder + '''" name="name" type="text">
                            <hr class=\"main_hr\">
                            <button type="submit">''' + load_lang('go') + '''</button>
                        </form>
                        ''',
                menu = [['manager', load_lang('return')]]
            ))
    else:
        return redirect()
        
@app.route('/title_index')
def title_index():
    page = int(number_check(flask.request.args.get('page', '1')))
    num = int(number_check(flask.request.args.get('num', '100')))
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

        sql_list = [load_lang('template', 1).lower() + ':', 'category:', 'user:', 'file:']
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
        imp = [load_lang('all_document_list'), wiki_set(), custom(), other2([sub, 0])],
        data = data,
        menu = [['other', load_lang('return')]]
    ))
        
@app.route('/topic/<everything:name>/sub/<sub>/b/<int:num>')
def topic_block(name = None, sub = None, num = 1):
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
def topic_top(name = None, sub = None, num = 1):
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
            load_lang('close', 1), 
            load_lang('open', 1)
        ]
    elif tool == 'stop':
        set_list = [
            '', 
            'O', 
            load_lang('stop', 1), 
            load_lang('restart', 1)
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
def topic_admin(name = None, sub = None, num = 1):
    curs.execute("select block, ip, date from topic where title = ? and sub = ? and id = ?", [name, sub, str(num)])
    data = curs.fetchall()
    if not data:
        return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub))

    ban = ''

    if admin_check(3) == 1:
        ban +=  '''
                </ul>
                <br>
                <h2>''' + load_lang('admin_tool') + '''</h2>
                <ul>
                '''
        is_ban = '<li><a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/b/' + str(num) + '">'

        if data[0][0] == 'O':
            is_ban += load_lang('hide_release')
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
            is_ban += load_lang('notice_release')
        else:
            is_ban += load_lang('notice') + ''
        
        is_ban += '</a></li></ul>'
        ban += '<li><a href="/ban/' + url_pas(data[0][1]) + '">'

        curs.execute("select end from ban where block = ?", [data[0][1]])
        if curs.fetchall():
            ban += load_lang('ban_release')
        else:
            ban += load_lang('ban')
        
        ban += '</a></li>' + is_ban

    ban +=  '''
            </ul>
            <br>
            <h2>''' + load_lang('other_tool') + '''</h2>
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
        imp = [load_lang('discussion_tool'), wiki_set(), custom(), other2([' (' + str(num) + ')', 0])],
        data = ban,
        menu = [['topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#' + str(num), load_lang('return')]]
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
            curs.execute('insert into alarm (name, data, date) values (?, ?, ?)', [match.groups()[0], ip + ' - <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '">' + load_lang('user_discussion', 1) + '</a>', today])
        
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
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/agree">(' + load_lang('destruction') + ')</a>'
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
                ip += ' <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/admin/' + str(number) + '">(' + load_lang('discussion_tool') + ')</a>'

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
            ban_name = load_lang('ban_release')
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
        menu = 0
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
        menu = [['topic/' + url_pas(name), load_lang('return')]]
        
        if tool == 'close':
            curs.execute("select sub from rd where title = ? and stop = 'O' order by sub asc", [name])
            
            sub = load_lang('close') + ''
        elif tool == 'agree':
            curs.execute("select sub from rd where title = ? and agree = 'O' order by sub asc", [name])
            
            sub = load_lang('agreement') + ''
        else:
            curs.execute("select sub from rd where title = ? order by date desc", [name])
            
            sub = load_lang('discussion_list')
            
            menu = [['w/' + url_pas(name), load_lang('document')]]
            
            plus =  '''
                    <a href="/topic/''' + url_pas(name) + '''/close">(''' + load_lang('close') + ''')</a> <a href="/topic/''' + url_pas(name) + '''/agree">(''' + load_lang('agreement') + ''')</a>
                    <hr class=\"main_hr\">
                    <input placeholder="''' + load_lang('discussion_name') + '''" name="topic" type="text">
                    <hr class=\"main_hr\">
                    <button type="submit">''' + load_lang('go') + '''</button>
                    '''

        for data in curs.fetchall():
            curs.execute("select data, date, ip, block from topic where title = ? and sub = ? and id = '1'", [name, data[0]])
            if curs.fetchall():                
                it_p = 0
                
                if sub == load_lang('discussion_list'):
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
                        <button type="submit">''' + load_lang('login') + '''</button>
                        <hr class=\"main_hr\">
                        ''' + oauth_content + '''
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('http_warring') + '''</span>
                    </form>
                    ''',
            menu = [['user', load_lang('return')]]
        ))

@app.route('/oauth/<regex("discord|naver|facebook"):platform>/<regex("init|callback"):func>', methods=['GET', 'POST'])
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

    if platform == 'discord':
        api_url['redirect'] = 'https://discordapp.com/api/oauth2/authorize'
        api_url['token'] = 'https://discordapp.com/api/oauth2/token'
        api_url['profile'] = 'https://discordapp.com/api/users/@me'
    elif platform == 'naver':
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
                imp = [load_lang('error'), wiki_set(), custom(), other2([0, 0])], 
                data = load_lang('oauth_disabled'), 
                menu = [['user', load_lang('return')]]
            ))
        elif publish_url == 'https://':
            return easy_minify(flask.render_template(skin_check(), 
                imp = [load_lang('error'), wiki_set(), custom(), other2([0, 0])], 
                data = load_lang('oauth_settings_not_found'), 
                menu = [['user', load_lang('return')]]
            ))

        referrer_re = re.compile(r'(?P<host>^(https?):\/\/([^\/]+))\/(?P<refer>[^\/?]+)')
        if flask.request.referrer != None:
            referrer = referrer_re.search(flask.request.referrer)
            if referrer.group('host') != load_oauth('publish_url'):
                return redirect()
            else:
                flask.session['referrer'] = referrer.group('refer')
        else:
            return redirect()

        flask.session['refer'] = flask.request.referrer

        if platform == 'discord':
            return redirect(api_url['redirect'] + '?client_id={}&redirect_uri={}&response_type=code&scope=identify'.format(
                data['client_id'], 
                data['redirect_uri']
            ))
        elif platform == 'naver':
            return redirect(api_url['redirect'] + '?response_type=code&client_id={}&redirect_uri={}&state={}'.format(
                data['client_id'], 
                data['redirect_uri'], 
                data['state']
            ))
        elif platform == 'facebook':
            return redirect(api_url['redirect'] + '?client_id={}&redirect_uri={}&state={}'.format(
                data['client_id'], 
                data['redirect_uri'], 
                data['state']
            ))

    elif func == 'callback':
        code = flask.request.args.get('code')
        state = flask.request.args.get('state')

        if code == None:
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('inter_error'), wiki_set(), custom(), other2([0, 0])],
                data = '<h2>ie_wrong_callback</h2>' + load_lang('ie_wrong_callback'),
                menu = [['user', load_lang('return')]]
            ))

        if platform == 'discord':
            data = {
                'client_id'     : data['client_id'],
                'client_secret' : data['client_secret'],
                'grant_type'    : 'authorization_code',
                'redirect_uri'  : data['redirect_uri'],
                'scope'         : 'identify',
                'code'          : code
            }
            headers = {
                'Content-Type': 'application/x-www-form-urlencoded',
                'User-Agent': 'Mozilla/5.0'
            }
            token_exchange = urllib.request.Request(
                'https://discordapp.com/api/oauth2/token',
                data = bytes(urllib.parse.urlencode(data).encode()),
                headers = headers
            )
            token_result = urllib.request.urlopen(token_exchange).read()
            token_json = json.loads(token_result)

            headers = {
                'User-Agent'    : 'Mozilla/5.0',
                'Authorization' : 'Bearer ' + token_json['access_token']
            }
            profile_exchange = urllib.request.Request(
                'https://discordapp.com/api/users/@me',
                headers = headers
            )
            profile_result =  urllib.request.urlopen(profile_exchange).read().decode('utf-8')
            profile_result_json = json.loads(profile_result)
            stand_json = {
                'id'        : profile_result_json['id'], 
                'name'      : profile_result_json['username'] + '#' + profile_result_json['discriminator'],
                'picture'   : profile_result_json['avatar']
            }
        elif platform == 'naver':
            token_access = api_url['token'] + '?grant_type=authorization_code&client_id={}&client_secret={}&code={}&state={}'.format(
                data['client_id'], 
                data['client_secret'], 
                code, 
                state
            )
            token_result = urllib.request.urlopen(token_access).read().decode('utf-8')
            token_result_json = json.loads(token_result)

            headers = {
                'Authorization': 'Bearer {}'.format(token_result_json['access_token'])
            }

            profile_access = urllib.request.Request(api_url['profile'], headers = headers)
            profile_result = urllib.request.urlopen(profile_access).read().decode('utf-8')
            profile_result_json = json.loads(profile_result)

            stand_json = {
                'id'        : profile_result_json['response']['id'],
                'name'      : profile_result_json['response']['name'],
                'picture'   : profile_result_json['response']['profile_image']
            }
        elif platform == 'facebook':
            token_access = api_url['token'] + '?client_id={}&redirect_uri={}&client_secret={}&code={}'.format(
                data['client_id'], 
                data['redirect_uri'], 
                data['client_secret'], 
                code
            )
            token_result = urllib.request.urlopen(token_access).read().decode('utf-8')
            token_result_json = json.loads(token_result)

            profile_access = api_url['profile'] + '?fields=id,name,picture&access_token={}'.format(token_result_json['access_token'])
            profile_result = urllib.request.urlopen(profile_access).read().decode('utf-8')
            profile_result_json = json.loads(profile_result)

            stand_json = {
                'id': profile_result_json['id'], 
                'name': profile_result_json['name'], 
                'picture': profile_result_json['picture']['data']['url']
            }
        
        if flask.session['referrer'][0:6] == 'change':
            curs.execute('select * from oauth_conn where wiki_id = ? and provider = ?', [flask.session['id'], platform])
            oauth_result = curs.fetchall()
            if len(oauth_result) == 0:
                curs.execute('insert into oauth_conn (provider, wiki_id, sns_id, name, picture) values(?, ?, ?, ?, ?)', [
                    platform, 
                    flask.session['id'], 
                    stand_json['id'], 
                    stand_json['name'], 
                    stand_json['picture']
                ])
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
    support_language = server_init.server_set_var['language']['list']

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
                    oauth_content += '<li>{} - {}</li>'.format(oauth_provider[i], load_lang('connection') + ' : <img src="{}" width="17px" height="17px">{}'.format(oauth_data[0][1], oauth_data[0][0]))
                else:
                    oauth_content += '<li>{} - {}</li>'.format(oauth_provider[i], load_lang('connection') + ' : <a href="/oauth/{}/init">{}</a>'.format(oauth_provider[i], load_lang('connect')))
            
            oauth_content += '</ul>'

            return easy_minify(flask.render_template(skin_check(),    
                imp = [load_lang('user_setting'), wiki_set(), custom(), other2([0, 0])],
                data =  '''
                        <form method="post">
                            <span>id : ''' + ip + '''</span>
                            <hr class=\"main_hr\">
                            <input placeholder="''' + load_lang('now_password') + '''" name="pw4" type="password">
                            <hr class=\"main_hr\">
                            <input placeholder="''' + load_lang('new_password') + '''" name="pw2" type="password">
                            <hr class=\"main_hr\">
                            <input placeholder="''' + load_lang('password_confirm') + '''" name="pw3" type="password">
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('skin') + '''</span>
                            <hr class=\"main_hr\">
                            <select name="skin">''' + div2 + '''</select>
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('language') + '''</span>
                            <hr class=\"main_hr\">
                            <select name="lang">''' + div3 + '''</select>
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('oauth_connection') + '''</span>
                            ''' + oauth_content + '''
                            <hr class=\"main_hr\">
                            <button type="submit">''' + load_lang('save') + '''</button>
                            <hr class=\"main_hr\">
                            <span>''' + load_lang('http_warring') + '''</span>
                        </form>
                        ''',
                menu = [['user', load_lang('return')]]
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
        
    num = int(number_check(flask.request.args.get('num', '1')))
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
        menu = [['manager', load_lang('return')]]
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
                        <input placeholder="''' + load_lang('password_confirm') + '''" name="pw2" type="password">
                        <hr class=\"main_hr\">
                        ''' + captcha_get() + '''
                        <button type="submit">''' + load_lang('save') + '''</button>
                        <hr class=\"main_hr\">
                        <span>''' + load_lang('http_warring') + '''</span>
                    </form>
                    ''',
            menu = [['user', load_lang('return')]]
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

                send_email(flask.request.form.get('email', ''), wiki_set()[0] + ' ' + load_lang('password_search') + ' key', 'key : ' + flask.session['c_key'])

                return redirect('/check_pass_key')
    else:
        if tool == 'need_email':
            return easy_minify(flask.render_template(skin_check(),    
                imp = [load_lang('email'), wiki_set(), custom(), other2([0, 0])],
                data =  '''
                        <a href="/email_filter">(''' + load_lang('email_filter_list') + ''')</a>
                        <hr class=\"main_hr\">
                        <form method="post">
                            <input placeholder="''' + load_lang('email') + '''" name="email" type="text">
                            <hr class=\"main_hr\">
                            <button type="submit">''' + load_lang('save') + '''</button>
                        </form>
                        ''',
                menu = [['user', load_lang('return')]]
            ))
        else:
            return easy_minify(flask.render_template(skin_check(),    
                imp = [load_lang('password_search'), wiki_set(), custom(), other2([0, 0])],
                data =  '''
                        <form method="post">
                            <input placeholder="''' + load_lang('id') + '''" name="id" type="text">
                            <hr class=\"main_hr\">
                            <input placeholder="email" name="email" type="text">
                            <hr class=\"main_hr\">
                            <button type="submit">''' + load_lang('save') + '''</button>
                        </form>
                        ''',
                menu = [['user', load_lang('return')]]
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
                    menu = [['user', load_lang('return')]]
                ))
            else:
                return redirect('/pass_find')
    else:
        return easy_minify(flask.render_template(skin_check(),    
            imp = ['check', wiki_set(), custom(), other2([0, 0])],
            data =  '''
                    <form method="post">
                        <input placeholder="''' + load_lang('key') + '''" name="key" type="text">
                        <hr class=\"main_hr\">
                        <button type="submit">''' + load_lang('save') + '''</button>
                    </form>
                    ''',
            menu = [['user', load_lang('return')]]
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
                data = '<ul><li>' + load_lang('limitless') + '</li>'
            else:
                data = '<ul><li>' + load_lang('period') + ' : ' + end[0][0] + '</li>'
                
            curs.execute("select block from ban where block = ? and login = 'O'", [name])
            if curs.fetchall():
                data += '<li>' + load_lang('login_able') + '</li>'

            if end[0][1] != '':
                data += '<li>' + load_lang('why') + ' : ' + end[0][1] + '</li></ul><hr class=\"main_hr\">'
            else:
                data += '</ul><hr class=\"main_hr\">'
        else:
            if re.search("^([0-9]{1,3}\.[0-9]{1,3})$", name):
                now = load_lang('band_ban')
            else:
                now = load_lang('ban')
                
            if ip_or_user(name) == 1:
                plus = '<input type="checkbox" name="login"> ' + load_lang('login_able') + '<hr class=\"main_hr\">'
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
            menu = [['manager', load_lang('return')]]
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
        if flask.request.form.get('dec', '') != flask.request.form.get('view', ''):
            dec = flask.request.form.get('view', '')
            view = flask.request.form.get('view', '')
        else:
            dec = flask.request.form.get('dec', '')
            view = flask.request.form.get('view', '')

        curs.execute("select title from acl where title = ?", [name])
        if curs.fetchall():
            curs.execute("update acl set dec = ? where title = ?", [dec, name])
            curs.execute("update acl set dis = ? where title = ?", [flask.request.form.get('dis', ''), name])
            curs.execute("update acl set why = ? where title = ?", [flask.request.form.get('why', ''), name])
            curs.execute("update acl set view = ? where title = ?", [view, name])
        else:
            curs.execute("insert into acl (title, dec, dis, why, view) values (?, ?, ?, ?, ?)", [name, dec, flask.request.form.get('dis', ''), flask.request.form.get('why', ''), view])
        
        curs.execute("select title from acl where title = ? and dec = '' and dis = ''", [name])
        if curs.fetchall():
            curs.execute("delete from acl where title = ?", [name])

        conn.commit()
            
        return redirect('/acl/' + url_pas(name))            
    else:
        data = '' + load_lang('document_acl') + '<br><br><select name="dec" ' + check_ok + '>'
    
        if re.search('^user:', name):
            acl_list = [['', load_lang('normal')], ['user', load_lang('member')], ['all', load_lang('all')]]
        else:
            acl_list = [['', load_lang('normal')], ['user', load_lang('member')], ['admin', load_lang('admin')]]
        
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
            data += '<hr class=\"main_hr\">' + load_lang('discussion_acl') + '<br><br><select name="dis" ' + check_ok + '>'
        
            curs.execute("select dis, why, view from acl where title = ?", [name])
            acl_data = curs.fetchall()
            for data_list in acl_list:
                if acl_data and acl_data[0][0] == data_list[0]:
                    check = 'selected="selected"'
                else:
                    check = ''
                    
                data += '<option value="' + data_list[0] + '" ' + check + '>' + data_list[1] + '</option>'
                
            data += '</select>'

            data += '<hr class=\"main_hr\">' + load_lang('view_acl') + '<br><br><select name="view" ' + check_ok + '>'
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
            imp = [name, wiki_set(), custom(), other2([' (' + load_lang('acl') + ')', 0])],
            data =  '''
                <form method="post">
                    ''' + data + '''
                    <hr class=\"main_hr\">
                    <button type="submit" ''' + check_ok + '''>''' + load_lang('save') + '''</button>
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
            imp = [name, wiki_set(), custom(), other2([' (' + load_lang('authorize') + ')', 0])],
            data =  '''
                    <form method="post">
                        <select name="select">''' + div + '''</select>
                        <hr class=\"main_hr\">
                        <button type="submit">''' + load_lang('save') + '''</button>
                    </form>
                    ''',
            menu = [['manager', load_lang('return')]]
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
                menu = [['history/' + url_pas(name), load_lang('return')]]
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
        imp = [name, wiki_set(), custom(), other2([' (' + load_lang('sub') + ')', 0])],
        data = div,
        menu = [['w/' + url_pas(name), load_lang('return')]]
    ))

@app.route('/w/<everything:name>')
def read_view(name = None):
    return read_view_2(conn, name)

@app.route('/topic_record/<name>')
def user_topic_list(name = None):
    num = int(number_check(flask.request.args.get('num', '1')))
    if num * 50 > 0:
        sql_num = num * 50 - 50
    else:
        sql_num = 0
    
    one_admin = admin_check(1)

    div =   '''
            <table id="main_table_set">
                <tbody>
                    <tr>
                        <td id="main_table_width">''' + load_lang('discussion_name') + '''</td>
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
        sub = ' (' + load_lang('blocked') + ')'
    else:
        sub = 0 
    
    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('discussion_record'), wiki_set(), custom(), other2([sub, 0])],
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

        div =   '''
                <table id="main_table_set">
                    <tbody>
                        <tr>
                '''
        
        if name:
            num = int(number_check(flask.request.args.get('num', '1')))
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
                            <td id="main_table_width">''' + load_lang('document_name') + '''</td>
                            <td id="main_table_width">''' + load_lang('editor') + '''</td>
                            <td id="main_table_width">''' + load_lang('time') + '''</td>
                        </tr>
                        '''

                div = '<a href="/topic_record/' + url_pas(name) + '">(' + load_lang('discussion') + ')</a><hr class=\"main_hr\">' + div
                
                curs.execute("select id, title, date, ip, send, leng from history where ip = ? order by date desc limit ?, '50'", [name, str(sql_num)])
        else:
            num = int(number_check(flask.request.args.get('num', '1')))
            if num * 50 > 0:
                sql_num = num * 50 - 50
            else:
                sql_num = 0            
            
            div +=  '''
                        <td id="main_table_width">''' + load_lang('document_name') + '''</td>
                        <td id="main_table_width">''' + load_lang('editor') + '''</td>
                        <td id="main_table_width">''' + load_lang('time') + '''</td>
                    </tr>
                    '''
            
            curs.execute("select id, title, date, ip, send, leng from history where not title like 'user:%' order by date desc limit ?, 50", [str(sql_num)])

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
                    hidden = ' <a href="/hidden/' + url_pas(data[1]) + '?num=' + data[0] + '">(' + load_lang('hide_release') + ')'
                    
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
                title = '<a href="/w/' + url_pas(name) + '?num=' + data[0] + '">r' + data[0] + '</a> <a href="/raw/' + url_pas(name) + '?num=' + data[0] + '">(' + load_lang('raw') + ')</a> '
            else:
                title = '<a href="/w/' + url_pas(data[1]) + '">' + html.escape(data[1]) + '</a> <a href="/history/' + url_pas(data[1]) + '">(r' + data[0] + ')</a> '
                    
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
                        ''' + div
                title = name
                
                sub += ' (' + load_lang('history') + ')'
                
                menu = [['w/' + url_pas(name), load_lang('document')], ['raw/' + url_pas(name), 'raw']]
                
                div += next_fix('/history/' + url_pas(name) + '?num=', num, data_list)
            else:
                curs.execute("select end from ban where block = ?", [name])
                if curs.fetchall():
                    sub += ' (' + load_lang('blocked') + ')'

                title = load_lang('edit_record')
                
                menu = [['other', load_lang('other')], ['user', load_lang('user')], ['count/' + url_pas(name), load_lang('count')]]
                
                div += next_fix('/record/' + url_pas(name) + '?num=', num, data_list)
        else:
            menu = 0
            title = load_lang('recent_change')
                
            div += next_fix('/recent_changes?num=', num, data_list)
        
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
                lice = ip
            else:
                lice = '[[user:' + ip + ']]'
            
        if os.path.exists(os.path.join(app_var['PATH_DATA_IMAGES'], e_data)):
            os.remove(os.path.join(app_var['PATH_DATA_IMAGES'], e_data))
            
            data.save(os.path.join(app_var['PATH_DATA_IMAGES'], e_data))
        else:
            data.save(os.path.join(app_var['PATH_DATA_IMAGES'], e_data))
            
        curs.execute("select title from data where title = ?", ['file:' + name])
        if curs.fetchall(): 
            curs.execute("delete from data where title = ?", ['file:' + name])
        
        curs.execute("insert into data (title, data) values (?, ?)", ['file:' + name, '[[file:' + name + ']][br][br]{{{[[file:' + name + ']]}}}[br][br]' + lice])
        curs.execute("insert into acl (title, dec, dis, why, view) values (?, 'admin', '', '', '')", ['file:' + name])

        history_plus(
            'file:' + name, '[[file:' + name + ']][br][br]{{{[[file:' + name + ']]}}}[br][br]' + lice,
            get_time(), 
            ip, 
            '(upload)',
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
            menu = [['other', load_lang('return')]]
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
                acl = load_lang('member')
        else:
            acl = load_lang('normal')
    else:
        acl = load_lang('blocked')

        match = re.search("^([0-9]{1,3}\.[0-9]{1,3})", ip)
        if match:
            match = match.groups()[0]
        else:
            match = '-'

        curs.execute("select end, login, band from ban where block = ? or block = ?", [ip, match])
        block_data = curs.fetchall()
        if block_data:
            if block_data[0][0] != '':
                acl += ' (' + load_lang('period') + ' : ' + block_data[0][0] + ')'
            else:
                acl += ' (' + load_lang('limitless') + ')'        

            if block_data[0][1] != '':
                acl += ' (' + load_lang('login_able') + ')'

            if block_data[0][2] == 'O':
                acl += ' (' + load_lang('band_blocked') + ')'
            
    if custom()[2] != 0:
        ip_user = '<a href="/w/user:' + ip + '">' + ip + '</a>'
        
        plus =  '''
                <li><a href="/logout">''' + load_lang('logout') + '''</a></li>
                <li><a href="/change">''' + load_lang('user_setting') + '''</a></li>
                '''
        
        curs.execute('select name from alarm where name = ? limit 1', [ip_check()])
        if curs.fetchall():
            plus2 = '<li><a href="/alarm">' + load_lang('alarm') + ' (O)</a></li>'
        else:
            plus2 = '<li><a href="/alarm">' + load_lang('alarm') + '</a></li>'

        plus2 += '<li><a href="/watch_list">' + load_lang('watchlist') + '</a></li>'
        plus3 = '<li><a href="/acl/user:' + url_pas(ip) + '">' + load_lang('user_document_acl') + '</a></li>'
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
            plus += '<li><a href="/pass_find">' + load_lang('password_search') + '</a></li>'

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('user') + ' ' + load_lang('tool'), wiki_set(), custom(), other2([0, 0])],
        data =  '''
                <h2>''' + load_lang('state') + '''</h2>
                <ul>
                    <li>''' + ip_user + ''' <a href="/record/''' + url_pas(ip) + '''">(''' + load_lang('record') + ''')</a></li>
                    <li>''' + load_lang('state') + ''' : ''' + acl + '''</li>
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
                    <li><a href="/custom_head">''' + load_lang('user_head') + '''</a></li>
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

    div += '<a href="/manager/13">(' + load_lang('add') + ')</a>'

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('watchlist'), wiki_set(), custom(), other2([0, 0])],
        data = div,
        menu = [['manager', load_lang('return')]]
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
            imp = [load_lang(data = 'user_head', safe = 1), wiki_set(), custom(), other2([0, 0])],
            data =  start + '''
                    <form method="post">
                        <textarea rows="25" cols="100" name="content">''' + data + '''</textarea>
                        <hr class=\"main_hr\">
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                    ''',
            menu = [['user', load_lang('return')]]
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
                    <li><a href="/record/''' + url_pas(that) + '''">''' + load_lang('edit_record') + '''</a> : ''' + str(data) + '''</li>
                    <li><a href="/topic_record/''' + url_pas(that) + '''">''' + load_lang('discussion_record') + '''</a> : ''' + str(t_data) + '''</a></li>
                </ul>
                ''',
        menu = [['user', load_lang('return')]]
    ))
        
@app.route('/random')
def title_random():
    curs.execute("select title from data order by random() limit 1")
    data = curs.fetchall()
    if data:
        return redirect('/w/' + url_pas(data[0][0]))
    else:
        return redirect()

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
        json_data = { "title" : name, "data" : render_set(title = name, data = data[0][0], s_data = 1) }
    
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
    http_server.listen(server_set['port'], address=server_set['host'])
    tornado.ioloop.IOLoop.instance().start()
