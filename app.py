# Init
from route import *

# Init-Version
version_list = json.loads(open('version.json', encoding = 'utf8').read())

# Init-DB
if os.path.exists(os.path.join('data', 'set.json')):
    db_set_list = ['db', 'db_type']
    set_data = json.loads(open(
        os.path.join('data', 'set.json'), 
        encoding = 'utf8'
    ).read())
    for i in db_set_list:
        if not i in set_data:
            print('Please delete set.json')
            print('----')
            raise

    print('DB name : ' + set_data['db'])
    print('DB type : ' + set_data['db_type'])
elif os.getenv('NAMU_DB') or os.getenv('NAMU_DB_TYPE'):
    set_data = {}

    if os.getenv('NAMU_DB'):
        set_data['db'] = os.getenv('NAMU_DB')
    else:
        set_data['db'] = 'data'

    if os.getenv('NAMU_DB_TYPE'):
        set_data['db'] = os.getenv('NAMU_DB_TYPE')
    else:
        set_data['db'] = 'sqlite'

    print('DB name : ' + set_data['db'])
    print('DB type : ' + set_data['db_type'])
else:
    set_data = {}
    normal_db_type = ['sqlite', 'mysql']

    print('DB type (' + normal_db_type[0] + ') [' + ', '.join(normal_db_type) + '] : ', end = '')
    data_get = str(input())
    if data_get == '' or not data_get in normal_db_type:
        set_data['db_type'] = 'sqlite'
    else:
        set_data['db_type'] = data_get

    all_src = []
    if set_data['db_type'] == 'sqlite':
        for i_data in os.listdir("."):
            f_src = re.search(r"(.+)\.db$", i_data)
            if f_src:
                all_src += [f_src.group(1)]

    print('DB name (data) [' + ', '.join(all_src) + '] : ', end = '')
    
    data_get = str(input())
    if data_get == '':
        set_data['db'] = 'data'
    else:
        set_data['db'] = data_get
        
    with open(os.path.join('data', 'set.json'), 'w', encoding = 'utf8') as f:
        f.write(json.dumps(set_data))

db_data_get(set_data['db_type'])

if set_data['db_type'] == 'mysql':
    if not os.path.exists(os.path.join('data', 'mysql.json')):
        db_set_list = ['user', 'password', 'host', 'port']
        set_data = json.loads(open(os.path.join('data', 'mysql.json'), encoding = 'utf8').read())
        for i in db_set_list:
            if not i in set_data:
                print('Please delete mysql.json')
                print('----')
                raise

        print('DB user ID : ', end = '')
        set_data_mysql['user'] = str(input())

        print('DB password : ', end = '')
        set_data_mysql['password'] = str(input())
                
        print('DB host (localhost) : ', end = '')
        set_data_mysql['host'] = str(input())
        if set_data_mysql['host'] == '':
            set_data_mysql['host'] = 'localhost'

        print('DB port (3306) : ', end = '')
        set_data_mysql['port'] = str(input())
        if set_data_mysql['port'] == '':
            set_data_mysql['port'] = '3306'

        with open(os.path.join('data', 'mysql.json'), 'w', encoding = 'utf8') as f:
            f.write(json.dumps(set_data_mysql))

    conn = pymysql.connect(
        host = set_data_mysql['host'] if 'host' in set_data_mysql else 'localhost',
        user = set_data_mysql['user'],
        password = set_data_mysql['password'],
        charset = 'utf8mb4',
        port = int(set_data_mysql['port']) if 'port' in set_data_mysql else 3306
    )
    curs = conn.cursor()

    try:
        curs.execute(db_change('create database ' + set_data['db'] + ' default character set utf8mb4;'))
    except:
        pass
    
    conn.select_db(set_data['db'])
else:
    conn = sqlite3.connect(set_data['db'] + '.db')
    curs = conn.cursor()

load_conn(conn)

# Init-Create_DB
create_data = {}
create_data['data'] = ['title', 'data', 'type']
create_data['history'] = ['id', 'title', 'data', 'date', 'ip', 'send', 'leng', 'hide', 'type']
create_data['rc'] = ['id', 'title', 'date', 'type']
create_data['rd'] = ['title', 'sub', 'code', 'date', 'band', 'stop', 'agree', 'acl']
create_data['user_set'] = ['name', 'id', 'data']
create_data['topic'] = ['id', 'data', 'date', 'ip', 'block', 'top', 'code']

# 폐지 예정
create_data['rb'] = ['block', 'end', 'today', 'blocker', 'why', 'band', 'login', 'ongoing']

create_data['back'] = ['title', 'link', 'type']
create_data['other'] = ['name', 'data', 'coverage']
create_data['alist'] = ['name', 'acl']
create_data['re_admin'] = ['who', 'what', 'time']
create_data['alarm'] = ['name', 'data', 'date']
create_data['ua_d'] = ['name', 'ip', 'ua', 'today', 'sub']
create_data['scan'] = ['user', 'title', 'type']
create_data['acl'] = ['title', 'data', 'type']
create_data['html_filter'] = ['html', 'kind', 'plus', 'plus_t']
create_data['vote'] = ['name', 'id', 'subject', 'data', 'user', 'type', 'acl']
for create_table in create_data:
    try:
        curs.execute(db_change('select test from ' + i + ' limit 1'))
    except:
        try:
            curs.execute(db_change('create table ' + i + '(test longtext)'))
        except:
            curs.execute(db_change("alter table " + i + " add test longtext"))
            
    for create in create_data[create_table]:
        try:
            curs.execute(db_change(
                'select ' + create + ' from ' + create_table + ' limit 1'
            ))
        except:
            try:
                curs.execute(db_change(
                    "alter table " + create_table + " add " + create + " longtext default ''"
                ))
            except:
                curs.execute(db_change(
                    "alter table " + create_table + " add " + create + " longtext"
                ))

curs.execute(db_change('select data from other where name = "ver"'))
ver_set_data = curs.fetchall()
if ver_set_data:
    if int(version_list['beta']['c_ver']) > int(ver_set_data[0][0]):
        setup_tool = 'update'
    else:
        setup_tool = 'normal'
else:
    setup_tool = 'init'

if setup_tool != 'normal':
    if setup_tool == 'update':
        update(int(ver_set_data[0][0]), set_data)
    else:
        set_init()

set_init_always(version_list['beta']['c_ver'])

# Init-Route
class EverythingConverter(werkzeug.routing.PathConverter):
    regex = '.*?'

class RegexConverter(werkzeug.routing.BaseConverter):
    def __init__(self, url_map, *items):
        super(RegexConverter, self).__init__(url_map)
        self.regex = items[0]
        
app = flask.Flask(__name__, template_folder = './')
app.config['JSON_AS_ASCII'] = False

app.logger.setLevel(logging.ERROR)

app.jinja_env.filters['md5_replace'] = md5_replace
app.jinja_env.filters['load_lang'] = load_lang
app.jinja_env.filters['cut_100'] = cut_100

app.url_map.converters['everything'] = EverythingConverter
app.url_map.converters['regex'] = RegexConverter

curs.execute(db_change('select data from other where name = "key"'))
sql_data = curs.fetchall()
app.secret_key = sql_data[0][0]

print('----')

# Init-DB_Data
dislay_set_key = ['Host', 'Port', 'Language', 'Markup', 'Encryption method']
server_set_key = ['host', 'port', 'language', 'markup', 'encode']
server_set = {}

server_init = server_init()
for i in range(len(server_set_key)):
    curs.execute(db_change('select data from other where name = ?'), [server_set_key[i]])
    server_set_val = curs.fetchall()
    if not server_set_val:
        server_set_val = server_init.init(server_set_key[i])

        curs.execute(db_change('insert into other (name, data) values (?, ?)'), [server_set_key[i], server_set_val])
    else:
        server_set_val = server_set_val[0][0]

    print(dislay_set_key[i] + ' : ' + server_set_val)

    server_set[server_set_key[i]] = server_set_val

print('----')
    
# Init-DB_care
if set_data['db_type'] == 'sqlite':
    def back_up(back_time, back_up_where):
        print('----')
        
        try:
            shutil.copyfile(
                set_data['db'] + '.db', 
                back_up_where
            )

            print('Back up : OK')
        except:
            print('Back up : Error')

        threading.Timer(
            60 * 60 * back_time, 
            back_up,
            [back_time, back_up_where]
        ).start()

    curs.execute(db_change('select data from other where name = "back_up"'))
    back_time = curs.fetchall()
    back_time = int(number_check(back_time[0][0])) if back_time else 0
    if back_time != 0:
        curs.execute(db_change('select data from other where name = "backup_where"'))
        back_up_where = curs.fetchall()
        if back_up_where and back_up_where[0][0] != '':
            back_up_where = back_up_where[0][0]
        else:
            back_up_where = 'back_' + set_data['db'] + '.db'

        print('Back up state : ' + str(back_time) + ' hours')

        back_up(back_time, back_up_where)
    else:
        print('Back up state : Turn off')
else:
    def mysql_dont_off(port):
        try:
            urllib.request.urlopen('http://localhost:' + port + '/')
        except:
            pass

        threading.Timer(
            60 * 60 * 3, 
            mysql_dont_off,
            [port]
        ).start()

    mysql_dont_off(server_set['port'])

print('Now running... http://localhost:' + server_set['port'])

if os.path.exists('custom.py'):
    from custom import custom_run

    custom_run(conn, app)

# Func
# Func-alarm
@app.route('/alarm')
def alarm():
    return alarm_2(conn)

@app.route('/del_alarm')
def alarm_del():
    return alarm_del_2(conn)

# Func-inter_wiki
@app.route('/<regex("inter_wiki|edit_top|image_license|(?:edit|email|file|name|extension)_filter"):tools>')
def inter_wiki(tools = None):
    return inter_wiki_2(conn, tools)

@app.route('/<regex("del_(?:inter_wiki|edit_top|image_license|(?:edit|email|file|name|extension)_filter)"):tools>/<name>')
def inter_wiki_del(tools = None, name = None):
    return inter_wiki_del_2(conn, tools, name)

@app.route('/<regex("plus_(?:inter_wiki|edit_top|image_license|(?:edit|email|file|name|extension)_filter)"):tools>', methods = ['POST', 'GET'])
@app.route('/<regex("plus_(?:inter_wiki|edit_top|image_license|(?:edit|email|file|name|extension)_filter)"):tools>/<name>', methods = ['POST', 'GET'])
def inter_wiki_plus(tools = None, name = None):
    return inter_wiki_plus_2(conn, tools, name)

# Func-list
@app.route('/not_close_topic')
def list_not_close_topic():
    return list_not_close_topic_2(conn)

@app.route('/old_page')
def list_old_page():
    return list_old_page_2(conn)

@app.route('/acl_list')
def list_acl():
    return list_acl_2(conn)

@app.route('/image_file_list')
def list_image_file():
    return list_image_file_2(conn)

@app.route('/admin_list')
def list_admin():
    return list_admin_2(conn)

@app.route('/user_log')
def list_user():
    return list_user_2(conn)

@app.route('/admin_log', methods = ['POST', 'GET'])
def list_admin_use():
    return list_admin_use_2(conn)

@app.route('/admin_group')
def list_admin_group():
    return list_admin_group_2(conn)

@app.route('/please')
def list_please():
    return list_please_2(conn)

@app.route('/title_index')
def list_title_index():
    return list_title_index_2(conn)

@app.route('/topic_record/<name>')
def list_user_topic(name = 'test'):
    return list_user_topic_2(conn, name)

@app.route('/<regex("long_page|short_page"):tool>')
def list_long_page(tool = 'long_page'):
    return list_long_page_2(conn, tool)

# Func-give
@app.route('/admin_plus/<name>', methods = ['POST', 'GET'])
def give_admin_groups(name = None):
    return give_admin_groups_2(conn, name)

@app.route('/delete_admin_group/<name>', methods = ['POST', 'GET'])
def give_delete_admin_group(name = None):
    return give_delete_admin_group_2(conn, name)

@app.route('/hidden/<everything:name>')
def give_history_hidden(name = None):
    return give_history_hidden_2(conn, name)

@app.route('/add_history/<everything:name>', methods = ['POST', 'GET'])
def give_history_add(name = None):
    return give_history_add_2(conn, name)

@app.route('/check/<name>')
def give_user_check(name = None):
    return give_user_check_2(conn, name)
    
@app.route('/check_delete', methods = ['POST', 'GET'])
def give_user_check_delete():
    return give_user_check_delete_2(conn)

@app.route('/ban', methods = ['POST', 'GET'])
@app.route('/ban/<name>', methods = ['POST', 'GET'])
def give_user_ban(name = None):
    return give_user_ban_2(conn, name)

@app.route('/acl/<everything:name>', methods = ['POST', 'GET'])
def give_acl(name = None):
    return give_acl_2(conn, name)

@app.route('/admin/<name>', methods = ['POST', 'GET'])
def give_admin(name = None):
    return give_admin_2(conn, name)

# Func-view
@app.route('/xref/<everything:name>')
def view_xref(name = None):
    return view_xref_2(conn, name)

@app.route('/raw/<everything:name>')
@app.route('/thread/<int:topic_num>/raw/<int:num>')
def view_raw(name = None, topic_num = None, num = None):
    return view_raw_2(conn, name, topic_num, num)

@app.route('/diff/<everything:name>')
def view_diff_data(name = None):
    return view_diff_data_2(conn, name)

@app.route('/down/<everything:name>')
def view_down(name = None):
    return view_down_2(conn, name)

@app.route('/w/<everything:name>')
def view_read(name = None):
    return view_read_2(conn, name)

# Func-recent
@app.route('/recent_discuss')
def recent_discuss():
    return recent_discuss_2(conn)

@app.route('/block_log')
@app.route('/<regex("block_user|block_admin"):tool>/<name>')
def recent_block(name = None, tool = None):
    return recent_block_2(conn, name, tool)

@app.route('/recent_changes')
@app.route('/<regex("record"):tool>/<name>')
@app.route('/<regex("history"):tool>/<everything:name>', methods = ['POST', 'GET'])
def recent_changes(name = None, tool = 'record'):
    return recent_changes_2(conn, name, tool)

@app.route('/history_tool/<everything:name>')
def recent_history_tool(name = None):
    return recent_history_tool_2(conn, name)

@app.route('/history_delete/<everything:name>', methods = ['POST', 'GET'])
def recent_history_delete(name = None):
    return recent_history_delete_2(conn, name)

# Func-search
@app.route('/search', methods=['POST'])
def search():
    return search_2(conn)

@app.route('/goto', methods=['POST'])
@app.route('/goto/<everything:name>', methods=['POST'])
def search_goto(name = 'test'):
    return search_goto_2(conn, name)

@app.route('/search/<everything:name>')
def search_deep(name = 'test'):
    return search_deep_2(conn, name)

# Func-edit
@app.route('/revert/<everything:name>', methods = ['POST', 'GET'])
def edit_revert(name = None):
    return edit_revert_2(conn, name)

@app.route('/edit/<everything:name>', methods = ['POST', 'GET'])
def edit(name = 'Test'):
    return edit_2(conn, name)

@app.route('/backlink_reset/<everything:name>')
def edit_backlink_reset(name = 'Test'):
    return edit_backlink_reset_2(conn, name)

@app.route('/delete/<everything:name>', methods = ['POST', 'GET'])
def edit_delete(name = None):
    return edit_delete_2(conn, name)

@app.route('/many_delete', methods = ['POST', 'GET'])
def edit_many_delete(name = None):
    return edit_many_delete_2(conn)

@app.route('/move/<everything:name>', methods = ['POST', 'GET'])
def edit_move(name = None):
    return edit_move_2(conn, name)

# Func-topic
@app.route('/thread/<int:topic_num>/b/<int:num>')
def topic_block(topic_num = 1, num = 1):
    return topic_block_2(conn, topic_num, num)

@app.route('/thread/<int:topic_num>/notice/<int:num>')
def topic_top(topic_num = 1, num = 1):
    return topic_top_2(conn, topic_num, num)

@app.route('/thread/<int:topic_num>/setting', methods = ['POST', 'GET'])
def topic_stop(topic_num = 1):
    return topic_stop_2(conn, topic_num)

@app.route('/thread/<int:topic_num>/acl', methods = ['POST', 'GET'])
def topic_acl(topic_num = 1):
    return topic_acl_2(conn, topic_num)

@app.route('/thread/<int:topic_num>/delete', methods = ['POST', 'GET'])
def topic_delete(topic_num = 1):
    return topic_delete_2(conn, topic_num)

@app.route('/thread/<int:topic_num>/tool')
def topic_tool(topic_num = 1):
    return topic_tool_2(conn, topic_num)

@app.route('/thread/<int:topic_num>/change', methods = ['POST', 'GET'])
def topic_change(topic_num = 1):
    return topic_change_2(conn, topic_num)

@app.route('/thread/<int:topic_num>/admin/<int:num>')
def topic_admin(topic_num = 1, num = 1):
    return topic_admin_2(conn, topic_num, num)

@app.route('/thread/<int:topic_num>', methods = ['POST', 'GET'])
def topic(topic_num = 1):
    return topic_2(conn, topic_num)

@app.route('/topic/<everything:name>', methods = ['POST', 'GET'])
def topic_close_list(name = 'test'):
    return topic_close_list_2(conn, name)

# Func-user
@app.route('/tool/<name>')
def user_tool(name = None):
    return user_tool_2(conn, name)

@app.route('/change', methods = ['POST', 'GET'])
def user_setting():
    return user_setting_2(conn, server_init)

@app.route('/user')
def user_info():
    return user_info_2(conn)

@app.route('/custom_head', methods=['GET', 'POST'])
def user_custom_head_view():
    return user_custom_head_view_2(conn)

@app.route('/count')
@app.route('/count/<name>')
def user_count_edit(name = None):
    return user_count_edit_2(conn, name)

# Func-login
@app.route('/2fa_login', methods = ['POST', 'GET'])
def login_2fa():
    return login_2fa_2(conn)

@app.route('/login', methods = ['POST', 'GET'])
def login():
    return login_2(conn)

@app.route('/pw_change', methods = ['POST', 'GET'])
def login_pw_change():
    return login_pw_change_2(conn)

@app.route('/register', methods = ['POST', 'GET'])
def login_register():
    return login_register_2(conn)

@app.route('/<regex("need_email|pass_find|email_change"):tool>', methods = ['POST', 'GET'])
def login_need_email(tool = 'pass_find'):
    return login_need_email_2(conn, tool)

@app.route('/<regex("check_key|check_pass_key|email_replace"):tool>', methods = ['POST', 'GET'])
def login_check_key(tool = 'check_pass_key'):
    return login_check_key_2(conn, tool)

@app.route('/logout')
def login_logout():
    return login_logout_2(conn)

# Func-watch_list
@app.route('/<regex("watch_list|star_doc"):tool>')
def watch_list(tool = 'star_doc'):
    return watch_list_2(conn, tool)

@app.route('/<regex("watch_list|star_doc"):tool>/<everything:name>')
def watch_list_name(tool = 'star_doc', name = 'Test'):
    return watch_list_name_2(conn, tool, name)

# Func-application
@app.route('/application_submitted')
def application_submitted():
    return application_submitted_2(conn)

@app.route('/applications', methods = ['POST', 'GET'])
def applications():
    return applications_2(conn)

# Func-vote
@app.route('/vote/<num>', methods = ['POST', 'GET'])
def vote_select(num = '1'):
    return vote_select_2(conn, num)

@app.route('/end_vote/<num>')
def vote_end(num = '1'):
    return vote_end_2(conn, num)

@app.route('/close_vote/<num>')
def vote_close(num = '1'):
    return vote_close_2(conn, num)

@app.route('/vote')
def vote():
    return vote_2(conn)

@app.route('/add_vote', methods = ['POST', 'GET'])
def vote_add():
    return vote_add_2(conn)

# Func-api
@app.route('/api/w/<everything:name>', methods = ['POST', 'GET'])
def api_w(name = ''):
    return api_w_2(conn, name)

@app.route('/api/raw/<everything:name>')
def api_raw(name = ''):
    return api_raw_2(conn, name)

@app.route('/api/version')
def api_version():
    return api_version_2(conn, version_list['beta']['r_ver'], version_list['beta']['c_ver'])

@app.route('/api/skin_info')
@app.route('/api/skin_info/<name>')
def api_skin_info(name = ''):
    return api_skin_info_2(conn, name)

@app.route('/api/markup')
def api_markup():
    return api_markup_2(conn)

@app.route('/api/user_info/<name>')
def api_user_info(name = ''):
    return api_user_info_2(conn, name)

@app.route('/api/thread/<topic_num>')
def api_topic_sub(name = '', topic_num = 1):
    return api_topic_sub_2(conn, topic_num)

@app.route('/api/search/<name>')
def api_search(name = ''):
    return api_search_2(conn, name)

@app.route('/api/recent_changes')
def api_recent_change():
    return api_recent_change_2(conn)

@app.route('/api/sha224/<everything:name>', methods = ['POST', 'GET'])
def api_sha224(name = 'test'):
    return api_sha224_2(conn, name)

@app.route('/api/title_index')
def api_title_index():
    return api_title_index_2(conn)

@app.route('/api/image/<everything:name>', methods = ['POST', 'GET'])
def api_image_view(name = ''):
    return api_image_view_2(conn, name)

@app.route('/api/sitemap.xml')
def api_sitemap():
    return api_sitemap_2(conn)

# Func-main
@app.route('/restart', methods = ['POST', 'GET'])
def main_restart():
    return main_restart_2(conn)

@app.route('/update', methods=['GET', 'POST'])
def main_update():
    return main_update_2(conn, version_list['beta']['r_ver'])

@app.route('/random')
def main_title_random():
    return main_title_random_2(conn)

@app.route('/upload', methods=['GET', 'POST'])
def main_upload():
    return main_upload_2(conn)

@app.route('/setting')
@app.route('/setting/<int:num>', methods = ['POST', 'GET'])
def setting(num = 0):
    return main_setting_2(conn, num, set_data['db_type'])

@app.route('/other')
def main_other():
    return main_other_2(conn)

@app.route('/manager', methods = ['POST', 'GET'])
@app.route('/manager/<int:num>', methods = ['POST', 'GET'])
def main_manager(num = 1):
    return main_manager_2(conn, num, version_list['beta']['r_ver'])

@app.route('/image/<everything:name>')
def main_image_view(name = None):
    return main_image_view_2(conn, name)

@app.route('/skin_set')
@app.route('/main_skin_set')
def main_skin_set():
    return main_skin_set_2(conn)

@app.route('/views/<everything:name>')
def main_views(name = None):
    return main_views_2(conn, name)

@app.route('/test_func')
def main_test_func():
	return main_test_func_2(conn)

@app.route('/shutdown', methods = ['POST', 'GET'])
def main_shutdown():
    return main_shutdown_2(conn)

@app.route('/<regex("easter_egg\.html|\.(?:txt|xml)$"):data>')
def main_file(data = ''):
    return main_file_2(conn, data)

@app.errorhandler(404)
def main_error_404(e):
    return main_error_404_2(conn)

# End
conn.commit()
	
if __name__ == "__main__":
    WSGIServer((
        server_set['host'], 
        int(server_set['port'])
    ), app, log = app.logger).serve_forever()