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

if not os.path.exists(app_var['path_data_image']):
    os.makedirs(app_var['path_data_image'])
    
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
    return admin_log_2(conn)

@app.route('/give_log')
def give_log():        
    return give_log_2(conn)

@app.route('/indexing', methods=['POST', 'GET'])
def indexing():
    return indexing_2(conn)       

@app.route('/restart', methods=['POST', 'GET'])
def restart():
    return restart_2(conn)

@app.route('/update')
def now_update():
    return now_update_2(conn)

@app.route('/oauth_setting', methods=['GET', 'POST'])
def oauth_setting():
    return oauth_setting_2(conn)

@app.route('/adsense_setting', methods=['GET', 'POST'])
def adsense_setting():
    return adsense_setting_2(conn)
        
@app.route('/xref/<everything:name>')
def xref(name = None):
    return xref_2(conn, name)

@app.route('/please')
def please():
    return please_2(conn)
        
@app.route('/recent_discuss')
def recent_discuss():
    return recent_discuss_2(conn)

@app.route('/block_log')
@app.route('/<regex("block_user|block_admin"):tool>/<name>')
def block_log(name = None, tool = None):
    return block_log_2(conn, name, tool)
            
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
    return deep_search_2(conn, name)
         
@app.route('/raw/<everything:name>')
@app.route('/topic/<everything:name>/sub/<sub_title>/raw/<int:num>')
def raw_view(name = None, sub_title = None, num = None):
    return raw_view_2(conn, name, sub_title, num)
        
@app.route('/revert/<everything:name>', methods=['POST', 'GET'])
def revert(name = None):    
    return revert_2(conn, name)

@app.route('/edit/<everything:name>', methods=['POST', 'GET'])
def edit(name = None):
    return edit_2(conn, name)

@app.route('/preview/<everything:name>', methods=['POST'])
def preview(name = None):
    return preview_2(conn, name)
        
@app.route('/delete/<everything:name>', methods=['POST', 'GET'])
def delete(name = None):
    return delete_2(conn, name)        
            
@app.route('/move/<everything:name>', methods=['POST', 'GET'])
def move(name = None):
    return move_2(conn, name)

@app.route('/other')
def other():
    return other_2(conn)
    
@app.route('/manager', methods=['POST', 'GET'])
@app.route('/manager/<int:num>', methods=['POST', 'GET'])
def manager(num = 1):
    return manager_2(conn, num)
        
@app.route('/title_index')
def title_index():
    return title_index_2(conn)
                
@app.route('/topic/<everything:name>/sub/<sub>/b/<int:num>')
def topic_block(name = None, sub = None, num = 1):
    return topic_block_2(conn, name, sub, num)
        
@app.route('/topic/<everything:name>/sub/<sub>/notice/<int:num>')
def topic_top(name = None, sub = None, num = 1):
    return topic_top_2(conn, name, sub, num)
                
@app.route('/topic/<everything:name>/sub/<sub>/tool/<regex("close|stop|agree"):tool>')
def topic_stop(name = None, sub = None, tool = None):
    return topic_stop_2(conn, name, sub, tool)

@app.route('/topic/<everything:name>/sub/<sub>/admin/<int:num>')
def topic_admin(name = None, sub = None, num = 1):
    return topic_admin_2(conn, name, sub, num)

@app.route('/topic/<everything:name>/sub/<sub>', methods=['POST', 'GET'])
def topic(name = None, sub = None):
    return topic_2(conn, name, sub)

@app.route('/tool/<name>')
def user_tool(name = None):
    return user_tool_2(conn, name)
        
@app.route('/topic/<everything:name>', methods=['POST', 'GET'])
@app.route('/topic/<everything:name>/<regex("close|agree"):tool>', methods=['GET'])
def close_topic_list(name = None, tool = None):
    return close_topic_list_2(conn, name, tool)
            
@app.route('/login', methods=['POST', 'GET'])
def login():
    return login_2(conn)

@app.route('/oauth/<regex("discord|naver|facebook"):platform>/<regex("init|callback"):func>', methods=['GET', 'POST'])
def login_oauth(platform = None, func = None):
    return login_oauth_2(conn, platform, func)
                
@app.route('/change', methods=['POST', 'GET'])
def change_password():
    return change_password_2(conn)

@app.route('/check/<name>')
def user_check(name = None):
    return user_check_2(conn, name)
                
@app.route('/register', methods=['POST', 'GET'])
def register():
    return register_2(conn)

@app.route('/<regex("need_email|pass_find"):tool>', methods=['POST', 'GET'])
def need_email(tool = 'pass_find'):
    return need_email_2(conn, tool)

@app.route('/<regex("check_key|check_pass_key"):tool>', methods=['POST', 'GET'])
def check_key(tool = 'check_pass_key'):
    return check_key_2(conn, tool)
           
@app.route('/logout')
def logout():
    flask.session['state'] = 0
    flask.session.pop('id', None)

    return redirect('/user')
    
@app.route('/ban/<name>', methods=['POST', 'GET'])
def user_ban(name = None):
    return user_ban_2(conn, name)         
                
@app.route('/acl/<everything:name>', methods=['POST', 'GET'])
def acl(name = None):
    return acl_2(conn, name)
            
@app.route('/admin/<name>', methods=['POST', 'GET'])
def user_admin(name = None):
    return user_admin_2(conn, name)
    
@app.route('/diff/<everything:name>')
def diff_data(name = None):
    return diff_data_2(conn, name)
        
@app.route('/down/<everything:name>')
def down(name = None):
    return down_2(conn, name)   

@app.route('/w/<everything:name>')
def read_view(name = None):
    return read_view_2(conn, name)

@app.route('/topic_record/<name>')
def user_topic_list(name = None):
    return user_topic_list_2(conn, name)

@app.route('/recent_changes')
@app.route('/<regex("record"):tool>/<name>')
@app.route('/<regex("history"):tool>/<everything:name>', methods=['POST', 'GET'])
def recent_changes(name = None, tool = 'record'):
    return recent_changes_2(conn, name, tool)
    
@app.route('/upload', methods=['GET', 'POST'])
def upload():
    return upload_2(conn)
        
@app.route('/user')
def user_info():
    return user_info_2(conn)

@app.route('/watch_list')
def watch_list():
    return watch_list_2(conn)

@app.route('/watch_list/<everything:name>')
def watch_list_name(name = None):
    return watch_list_name_2(conn, name)

@app.route('/custom_head', methods=['GET', 'POST'])
def custom_head_view():
    return custom_head_view_2(conn)

@app.route('/count')
@app.route('/count/<name>')
def count_edit(name = None):
    return count_edit_2(conn, name)
        
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