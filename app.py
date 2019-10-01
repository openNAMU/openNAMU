import os
import re

for i_data in os.listdir("route"):
    f_src = re.search("(.+)\.py$", i_data)
    if f_src:
        f_src = f_src.groups()[0]

        exec("from route." + f_src + " import *")

version_list = json.loads(open('version.json', encoding='utf-8').read())

r_ver = version_list['master']['r_ver']
c_ver = version_list['master']['c_ver']
s_ver = version_list['master']['s_ver']

print('Version : ' + r_ver)
print('DB set version : ' + c_ver)
print('Skin set version : ' + s_ver)
print('----')

app_var = json.loads(open('data/app_var.json', encoding='utf-8').read())

# DB
setup_tool = 0
db_type = {}
try:
    open('DB_Type.json', mode='r', encoding='utf-8')
except FileNotFoundError:
    db_type["DBMS"] = input('Select DBMS Type[mariadb, sqlite] (default: mariadb) : ')
    if db_type["DBMS"] == '':
        db_type["DBMS"] = 'mariadb'
    
    with open('DB_Type.json', mode='w', encoding='utf-8') as fileMake:
        json.dump(db_type, fileMake, ensure_ascii=False, indent="\t")
    pass

with open("DB_Type.json") as fileRead:
    db_type = json.load(fileRead)

if db_type["DBMS"] == "mariadb":
    db_data = {}

    # DB Connect(MariaDB)
    try:
        open('DB_Data.json', mode='r', encoding='utf-8')
    except FileNotFoundError:   # 예외 처리
        db_data["db_host"] = input('MariaDB Host (default: localhost) : ')
        if db_data["db_host"] == '':
            db_data["db_host"] = 'localhost'
            
        db_data["db_port"] = input('MariaDB Port (default: 3306) : ')
        if db_data["db_port"] == '':
            db_data["db_port"] = 3306
        else:
            db_data["db_port"] = int(db_data["db_port"])
        
        db_data["db_username"] = input('MariaDB Username (default: root) : ')
        if db_data["db_username"] == '':
            db_data["db_username"] = 'root'
        
        db_data["db_passwd"] = input('MariaDB Password : ')
        
        setup_tool = 1
        
        db_data["db_name"] = input('DB name (default: data) : ')
        if db_data["db_name"] == '':
            db_data["db_name"] = 'data'
        
        with open('DB_Data.json', mode='w', encoding='utf-8') as fileMake:
            json.dump(db_data, fileMake, ensure_ascii=False, indent="\t")
        pass

    with open("DB_Data.json") as fileRead:
        db_data = json.load(fileRead)

    conn = pymysql.connect( host=db_data["db_host"], port=db_data["db_port"], user=db_data["db_username"],
                                    passwd=db_data["db_passwd"], charset="utf8", autocommit=True )
    curs = conn.cursor()

    load_conn(conn)

    try:
        sqlQuery("CREATE DATABASE " + db_data["db_name"])
    except:   # 예외 처리
        setup_tool = 0
        pass

    sqlQuery("USE " + db_data["db_name"])
    curs.execute('SET @@global.sql_mode = "ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION";')
elif db_type["DBMS"] == "sqlite":
    try:
        set_data = json.loads(open('data/set.json').read())
    except:
        if os.getenv('NAMU_DB') != None:
            set_data = { "db" : os.getenv('NAMU_DB') }
        else:
            print('DB name (data) : ', end = '')
            
            new_json = str(input())
            if new_json == '':
                new_json = 'data'
                
            with open('data/set.json', 'w') as f:
                f.write('{ "db" : "' + new_json + '" }')
                
            set_data = json.loads(open('data/set.json').read())
            
    print('DB name : ' + set_data['db'])
    db_name = set_data['db']
                
    if os.path.exists(db_name + '.db'):
        setup_tool = 0
    else:
        setup_tool = 1

    conn = sqlite3.connect(db_name + '.db', check_same_thread = False)
    curs = conn.cursor()

    load_conn(conn)
else:
    print("DBMS Type Error")

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
app.jinja_env.filters['cut_100'] = cut_100

app.url_map.converters['everything'] = EverythingConverter

sqlQuery('create table if not exists data(test LONGTEXT)')
sqlQuery('create table if not exists cache_data(test LONGTEXT)')
sqlQuery('create table if not exists history(test LONGTEXT)')
sqlQuery('create table if not exists rd(test LONGTEXT)')
sqlQuery('create table if not exists user(test LONGTEXT)')
sqlQuery('create table if not exists user_set(test LONGTEXT)')
sqlQuery('create table if not exists ban(test LONGTEXT)')
sqlQuery('create table if not exists topic(test LONGTEXT)')
sqlQuery('create table if not exists rb(test LONGTEXT)')
sqlQuery('create table if not exists back(test LONGTEXT)')
sqlQuery('create table if not exists custom(test LONGTEXT)')
sqlQuery('create table if not exists other(test LONGTEXT)')
sqlQuery('create table if not exists alist(test LONGTEXT)')
sqlQuery('create table if not exists re_admin(test LONGTEXT)')
sqlQuery('create table if not exists alarm(test LONGTEXT)')
sqlQuery('create table if not exists ua_d(test LONGTEXT)')
sqlQuery('create table if not exists filter(test LONGTEXT)')
sqlQuery('create table if not exists scan(test LONGTEXT)')
sqlQuery('create table if not exists acl(test LONGTEXT)')
sqlQuery('create table if not exists inter(test LONGTEXT)')
sqlQuery('create table if not exists html_filter(test LONGTEXT)')
sqlQuery('create table if not exists oauth_conn(test LONGTEXT)')
## MARIADB & SQLITE 공통 호환

if setup_tool == 0:
    try:
        sqlQuery('select data from other where name = "ver"')
        ver_set_data = sqlQuery("fetchall")
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
    create_data['other'] = ['name', 'data', 'coverage']
    create_data['alist'] = ['name', 'acl']
    create_data['re_admin'] = ['who', 'what', 'time']
    create_data['alarm'] = ['name', 'data', 'date']
    create_data['ua_d'] = ['name', 'ip', 'ua', 'today', 'sub']
    create_data['filter'] = ['name', 'regex', 'sub']
    create_data['scan'] = ['user', 'title']
    create_data['acl'] = ['title', 'decu', 'dis', 'view', 'why']
    create_data['inter'] = ['title', 'link']
    create_data['html_filter'] = ['html', 'kind', 'plus']
    create_data['oauth_conn'] = ['provider', 'wiki_id', 'sns_id', 'name', 'picture']

    for create_table in create_data['all_data']:
        for create in create_data[create_table]:
            try:
                sqlQuery('select ' + create + ' from ' + create_table + ' limit 1')
            except:
                sqlQuery("alter table " + create_table + " add " + create + " LONGTEXT default ''")

    update()

# Init
sqlQuery('select name from alist where acl = "owner"')
if not sqlQuery("fetchall"):
    sqlQuery('delete from alist where name = "owner"')
    sqlQuery('insert into alist (name, acl) values ("owner", "owner")')

if not os.path.exists(app_var['path_data_image']):
    os.makedirs(app_var['path_data_image'])
    
if not os.path.exists('views'):
    os.makedirs('views')

import route.tool.init as server_init

dislay_set_key = ['Host', 'Port', 'Language', 'Markup', 'Encryption method']
server_set_key = ['host', 'port', 'language', 'markup', 'encode']
server_set = {}

for i in range(len(server_set_key)):
    sqlQuery('select data from other where name = ?', [server_set_key[i]])
    server_set_val = sqlQuery("fetchall")
    if not server_set_val:
        server_set_val = server_init.init(server_set_key[i])
        
        sqlQuery('insert into other (name, data) values (?, ?)', [server_set_key[i], server_set_val])
        sqlQuery("commit")
    else:
        server_set_val = server_set_val[0][0]
    
    print(dislay_set_key[i] + ' : ' + server_set_val)
    
    server_set[server_set_key[i]] = server_set_val

try:
    if not os.path.exists('robots.txt'):
        sqlQuery('select data from other where name = "robot"')
        robot_test = sqlQuery("fetchall")
        if robot_test:
            fw_test = open('./robots.txt', 'w')
            fw_test.write(re.sub('\r\n', '\n', robot_test[0][0]))
            fw_test.close()
        else:
            fw_test = open('./robots.txt', 'w')
            fw_test.write('User-agent: *\nDisallow: /\nAllow: /$\nAllow: /w/')
            fw_test.close()

            sqlQuery('insert into other (name, data) values ("robot", "User-agent: *\nDisallow: /\nAllow: /$\nAllow: /w/")')
        
        print('----')
        print('Engine made robots.txt')
except:
    pass

sqlQuery('select data from other where name = "key"')
rep_data = sqlQuery("fetchall")
if not rep_data:
    rep_key = ''.join(random.choice("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") for i in range(16))
    if rep_key:
        sqlQuery('insert into other (name, data) values ("key", ?)', [rep_key])
else:
    rep_key = rep_data[0][0]

sqlQuery('select data from other where name = "adsense"')
adsense_result = sqlQuery("fetchall")
if not adsense_result:
    sqlQuery('insert into other (name, data) values ("adsense", "False")')
    sqlQuery('insert into other (name, data) values ("adsense_code", "")')

sqlQuery('delete from other where name = "ver"')
sqlQuery('insert into other (name, data) values ("ver", ?)', [c_ver])

"""
def back_up():
    print('----')
    try:
        shutil.copyfile(db_name + '.db', 'back_' + db_name + '.db')
        
        print('Back up : OK')
    except:
        print('Back up : Error')

    threading.Timer(60 * 60 * back_time, back_up).start()

try:
    sqlQuery('select data from other where name = "back_up"')
    back_up_time = sqlQuery("fetchall")
    
    back_time = int(back_up_time[0][0])
except:
    back_time = 0
    
print('----')
if back_time != 0:
    print('Back up state : ' + str(back_time) + ' hours')
    
    if __name__ == '__main__':
        back_up()
else:
    print('Back up state : Turn off')
"""

sqlQuery('select data from other where name = "s_ver"')
ver_set_data = sqlQuery("fetchall")
if not ver_set_data:
    sqlQuery('insert into other (name, data) values ("s_ver", ?)', [s_ver])

    if setup_tool == 0:
        print('----')
        print('Skin update required')
else:
    if int(ver_set_data[0][0]) < int(s_ver):
        sqlQuery('delete from other where name = "s_ver"')
        sqlQuery('insert into other (name, data) values ("s_ver", ?)', [s_ver])

        print('----')
        print('Skin update required')
        
sqlQuery("commit")

# Func
@app.route('/del_alarm')
def alarm_del():
    return alarm_del_2()

@app.route('/alarm')
def alarm():
    return alarm_2()

@app.route('/<regex("inter_wiki|edit_top|image_license|(?:edit|email|file|name)_filter"):tools>')
def inter_wiki(tools = None):
    return inter_wiki_2(tools)

@app.route('/<regex("del_(?:inter_wiki|edit_top|image_license|(?:edit|email|file|name)_filter)"):tools>/<name>')
def inter_wiki_del(tools = None, name = None):
    return inter_wiki_del_2(tools, name)

@app.route('/<regex("plus_(?:inter_wiki|edit_top|image_license|(?:edit|email|file|name)_filter)"):tools>', methods=['POST', 'GET'])
@app.route('/<regex("plus_edit_filter"):tools>/<name>', methods=['POST', 'GET'])
def inter_wiki_plus(tools = None, name = None):
    return inter_wiki_plus_2(tools, name)

@app.route('/setting')
@app.route('/setting/<int:num>', methods=['POST', 'GET'])
def setting(num = 0):
    return setting_2(num)

@app.route('/not_close_topic')
def list_not_close_topic():
    return list_not_close_topic_2()

@app.route('/acl_list')
def list_acl():
    return list_acl_2()

@app.route('/admin_plus/<name>', methods=['POST', 'GET'])
def give_admin_groups(name = None):
    return give_admin_groups_2(name)
        
@app.route('/admin_list')
def list_admin():
    return list_admin_2()
        
@app.route('/hidden/<everything:name>')
def give_history_hidden(name = None):
    return give_history_hidden_2(name)
        
@app.route('/user_log')
def list_user():
    return list_user_2()

@app.route('/admin_log', methods=['POST', 'GET'])
def list_admin_use():
    return list_admin_use_2()

@app.route('/give_log')
def list_give():
    return list_give_2()

@app.route('/indexing', methods=['POST', 'GET'])
def server_indexing():
    return server_indexing_2()       

@app.route('/restart', methods=['POST', 'GET'])
def server_restart():
    return server_restart_2()

@app.route('/update', methods=['GET', 'POST'])
def server_now_update():
    return server_now_update_2()

@app.route('/oauth_setting', methods=['GET', 'POST'])
def setting_oauth():
    return setting_oauth_2()

@app.route('/adsense_setting', methods=['GET', 'POST'])
def setting_adsense():
    return setting_adsense_2()
        
@app.route('/xref/<everything:name>')
def view_xref(name = None):
    return view_xref_2(name)

@app.route('/please')
def list_please():
    return list_please_2()
        
@app.route('/recent_discuss')
def recent_discuss():
    return recent_discuss_2()

@app.route('/block_log')
@app.route('/<regex("block_user|block_admin"):tool>/<name>')
def list_block(name = None, tool = None):
    return list_block_2(name, tool)
            
@app.route('/search', methods=['POST'])
def search():
    return search_2()

@app.route('/goto', methods=['POST'])
@app.route('/goto/<everything:name>', methods=['POST'])
def search_goto(name = 'test'):
    return search_goto_2(name)

@app.route('/search/<everything:name>')
def search_deep(name = 'test'):
    return search_deep_2(name)
         
@app.route('/raw/<everything:name>')
@app.route('/topic/<everything:name>/sub/<sub_title>/raw/<int:num>')
def view_raw(name = None, sub_title = None, num = None):
    return view_raw_2(name, sub_title, num)
        
@app.route('/revert/<everything:name>', methods=['POST', 'GET'])
def edit_revert(name = None):
    return edit_revert_2(name)

@app.route('/edit/<everything:name>', methods=['POST', 'GET'])
def edit(name = None):
    return edit_2(name)
        
@app.route('/delete/<everything:name>', methods=['POST', 'GET'])
def edit_delete(name = None):
    return edit_delete_2(name, app_var)        
            
@app.route('/move/<everything:name>', methods=['POST', 'GET'])
def edit_move(name = None):
    return edit_move_2(name)

@app.route('/other')
def main_other():
    return main_other_2()
    
@app.route('/manager', methods=['POST', 'GET'])
@app.route('/manager/<int:num>', methods=['POST', 'GET'])
def main_manager(num = 1):
    return main_manager_2(num, r_ver)
        
@app.route('/title_index')
def list_title_index():
    return list_title_index_2()
                
@app.route('/topic/<everything:name>/sub/<sub>/b/<int:num>')
def topic_block(name = None, sub = None, num = 1):
    return topic_block_2(name, sub, num)
        
@app.route('/topic/<everything:name>/sub/<sub>/notice/<int:num>')
def topic_top(name = None, sub = None, num = 1):
    return topic_top_2(name, sub, num)
                
@app.route('/topic/<everything:name>/sub/<sub>/setting', methods=['POST', 'GET'])
def topic_stop(name = None, sub = None):
    return topic_stop_2(name, sub)

@app.route('/topic/<everything:name>/sub/<sub>/tool')
def topic_tool(name = None, sub = None):
    return topic_tool_2(name, sub)

@app.route('/topic/<everything:name>/sub/<sub>/admin/<int:num>')
def topic_admin(name = None, sub = None, num = 1):
    return topic_admin_2(name, sub, num)

@app.route('/topic/<everything:name>/sub/<sub>', methods=['POST', 'GET'])
def topic(name = None, sub = None):
    return topic_2(name, sub)
        
@app.route('/topic/<everything:name>', methods=['POST', 'GET'])
@app.route('/topic/<everything:name>/<regex("close|agree"):tool>', methods=['GET'])
def topic_close_list(name = None, tool = None):
    return topic_close_list_2(name, tool)

@app.route('/tool/<name>')
def user_tool(name = None):
    return user_tool_2(name)
            
@app.route('/login', methods=['POST', 'GET'])
def login():
    return login_2()

@app.route('/oauth/<regex("discord|naver|facebook|kakao"):platform>/<regex("init|callback"):func>', methods=['GET', 'POST'])
def login_oauth(platform = None, func = None):
    return login_oauth_2(platform, func)
                
@app.route('/change', methods=['POST', 'GET'])
def user_setting():
    return user_setting_2(server_init)
    
@app.route('/pw_change', methods=['POST', 'GET'])
def login_pw_change():
    return login_pw_change_2()

@app.route('/check/<name>')
def give_user_check(name = None):
    return give_user_check_2(name)
                
@app.route('/register', methods=['POST', 'GET'])
def login_register():
    return login_register_2()

@app.route('/<regex("need_email|pass_find|email_change"):tool>', methods=['POST', 'GET'])
def login_need_email(tool = 'pass_find'):
    return login_need_email_2(tool)

@app.route('/<regex("check_key|check_pass_key|email_replace"):tool>', methods=['POST', 'GET'])
def login_check_key(tool = 'check_pass_key'):
    return login_check_key_2(tool)
           
@app.route('/logout')
def login_logout():
    return login_logout_2()
    
@app.route('/ban', methods=['POST', 'GET'])
@app.route('/ban/<name>', methods=['POST', 'GET'])
def give_user_ban(name = None):
    return give_user_ban_2(name)         
                
@app.route('/acl/<everything:name>', methods=['POST', 'GET'])
def give_acl(name = None):
    return give_acl_2(name)
            
@app.route('/admin/<name>', methods=['POST', 'GET'])
def give_admin(name = None):
    return give_admin_2(name)
    
@app.route('/diff/<everything:name>')
def view_diff_data(name = None):
    return view_diff_data_2(name)
        
@app.route('/down/<everything:name>')
def view_down(name = None):
    return view_down_2(name)   

@app.route('/w/<everything:name>')
def view_read(name = None):
    return view_read_2(name)

@app.route('/topic_record/<name>')
def list_user_topic(name = None):
    return list_user_topic_2(name)

@app.route('/recent_changes')
@app.route('/<regex("record"):tool>/<name>')
@app.route('/<regex("history"):tool>/<everything:name>', methods=['POST', 'GET'])
def recent_changes(name = None, tool = 'record'):
    return recent_changes_2(name, tool)

@app.route('/history_tool/<everything:name>')
def recent_history_tool(name = None):
    return recent_history_tool_2(name)
    
@app.route('/history_delete/<everything:name>', methods=['POST', 'GET'])
def recent_history_delete(name = None):
    return recent_history_delete_2(name)

@app.route('/upload', methods=['GET', 'POST'])
def func_upload():
    return func_upload_2()
        
@app.route('/user')
def user_info():
    return user_info_2()

@app.route('/watch_list')
def watch_list():
    return watch_list_2()

@app.route('/watch_list/<everything:name>')
def watch_list_name(name = None):
    return watch_list_name_2(name)

@app.route('/custom_head', methods=['GET', 'POST'])
def user_custom_head_view():
    return user_custom_head_view_2()

@app.route('/count')
@app.route('/count/<name>')
def user_count_edit(name = None):
    return user_count_edit_2(name)
        
@app.route('/random')
def func_title_random():
    return func_title_random_2()

@app.route('/image/<name>')
def main_image_view(name = None):
    return main_image_view_2(name, app_var)

@app.route('/skin_set')
def main_skin_set():
    return main_skin_set_2()
    
# API
@app.route('/api/w/<everything:name>', methods=['POST', 'GET'])
def api_w(name = ''):
    return api_w_2(name)
    
@app.route('/api/raw/<everything:name>')
def api_raw(name = ''):
    return api_raw_2(name)

@app.route('/api/version')
def api_version():
    return api_version_2(r_ver, c_ver)

@app.route('/api/skin_info')
@app.route('/api/skin_info/<name>')
def api_skin_info(name = ''):
    return api_skin_info_2(name)

@app.route('/api/markup')
def api_markup():
    return api_markup_2()

@app.route('/api/user_info/<name>')
def api_user_info(name = ''):
    return api_user_info_2(name)

@app.route('/api/topic/<everything:name>/sub/<sub>')
def api_topic_sub(name = '', sub = '', time = ''):
    return api_topic_sub_2(name, sub, time)
    
# File
@app.route('/views/easter_egg.html')
def main_easter_egg():
    return main_easter_egg_2()

@app.route('/views/<everything:name>')
def main_views(name = None):
    return main_views_2(name)

@app.route('/<data>')
def main_file(data = None):
    return main_file_2(data)

# End
@app.errorhandler(404)
def main_error_404(e):
    return main_error_404_2()
    
app.secret_key = rep_key
app.wsgi_app = werkzeug.debug.DebuggedApplication(app.wsgi_app, True)
app.debug = True

if __name__ == "__main__":
    http_server = tornado.httpserver.HTTPServer(tornado.wsgi.WSGIContainer(app))
    http_server.listen(server_set['port'], address = server_set['host'])
    
    tornado.ioloop.IOLoop.instance().start()
