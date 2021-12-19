# Init
import os
import re

for i_data in os.listdir("route"):
    f_src = re.search(r"(.+)\.py$", i_data)
    if f_src:
        exec("from route." + f_src.group(1) + " import *")

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

data_db_set = {}
data_db_set['name'] = set_data['db']
data_db_set['type'] = set_data['db_type']

if data_db_set['type'] == 'mysql':
    if os.path.exists(os.path.join('data', 'mysql.json')):
        db_set_list = ['user', 'password', 'host', 'port']
        set_data = json.loads(
            open(
                os.path.join('data', 'mysql.json'),
                encoding = 'utf8'
            ).read()
        )
        for i in db_set_list:
            if not i in set_data:
                os.remove(os.path.join('data', 'mysql.json'))
                
                break
                
        set_data_mysql = set_data
    
    if not os.path.exists(os.path.join('data', 'mysql.json')):
        set_data_mysql = {}
        
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

        with open(
            os.path.join('data', 'mysql.json'), 
            'w', 
            encoding = 'utf8'
        ) as f:
            f.write(json.dumps(set_data_mysql))
            
    data_db_set['mysql_user'] = set_data_mysql['user']
    data_db_set['mysql_pw'] = set_data_mysql['password']
    if 'host' in set_data_mysql:
        data_db_set['mysql_host'] = set_data_mysql['host']
    else:
        data_db_set['mysql_host'] = 'localhost'
    
    if 'port' in set_data_mysql:
        data_db_set['mysql_port'] = set_data_mysql['port']
    else:
        data_db_set['mysql_port'] = '3306'

db_data_get(data_db_set['type'])
load_db = get_db_connect(data_db_set)

conn = load_db.db_load()
curs = conn.cursor()

setup_tool = ''
try:
    curs.execute(db_change('select data from other where name = "ver"'))
except:
    setup_tool = 'init'

if setup_tool != 'init':
    ver_set_data = curs.fetchall()
    if ver_set_data:
        if int(version_list['beta']['c_ver']) > int(ver_set_data[0][0]):
            setup_tool = 'update'
        else:
            setup_tool = 'normal'
    else:
        setup_tool = 'init'

if setup_tool != 'normal':
    # Init-Create_DB
    create_data = {}

    # 폐지 예정 (data_set으로 통합)
    create_data['data'] = ['title', 'data', 'type']
    create_data['history'] = ['id', 'title', 'data', 'date', 'ip', 'send', 'leng', 'hide', 'type']
    create_data['rc'] = ['id', 'title', 'date', 'type']
    create_data['acl'] = ['title', 'data', 'type']

    # 개편 예정 (data_link로 변경)
    create_data['back'] = ['title', 'link', 'type']

    # 폐지 예정 (topic_set으로 통합) [가장 시급]
    create_data['rd'] = ['title', 'sub', 'code', 'date', 'band', 'stop', 'agree', 'acl']
    create_data['topic'] = ['id', 'data', 'date', 'ip', 'block', 'top', 'code']

    # 폐지 예정 (user_set으로 통합)
    create_data['rb'] = ['block', 'end', 'today', 'blocker', 'why', 'band', 'login', 'ongoing']
    create_data['scan'] = ['user', 'title', 'type']

    # 개편 예정 (wiki_set과 wiki_filter과 wiki_vote으로 변경)
    create_data['other'] = ['name', 'data', 'coverage']
    create_data['html_filter'] = ['html', 'kind', 'plus', 'plus_t']
    create_data['vote'] = ['name', 'id', 'subject', 'data', 'user', 'type', 'acl']

    # 개편 예정 (auth_list와 auth_log로 변경)
    create_data['alist'] = ['name', 'acl']
    create_data['re_admin'] = ['who', 'what', 'time']

    # 개편 예정 (user_notice와 user_agent로 변경)
    create_data['alarm'] = ['name', 'data', 'date']
    create_data['ua_d'] = ['name', 'ip', 'ua', 'today', 'sub']

    create_data['user_set'] = ['name', 'id', 'data']
    for create_table in create_data:
        for create in ['test'] + create_data[create_table]:
            try:
                curs.execute(db_change('select ' + create + ' from ' + create_table + ' limit 1'))
            except:
                try:
                    curs.execute(db_change('create table ' + create_table + '(test longtext default "")'))
                except:
                    curs.execute(db_change("alter table " + create_table + " add column " + create + " longtext default ''"))

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

app = flask.Flask(
    __name__, 
    template_folder = './'
)

app.config['JSON_AS_ASCII'] = False
app.config['JSONIFY_PRETTYPRINT_REGULAR'] = True

log = logging.getLogger('werkzeug')
log.setLevel(logging.ERROR)

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
server_set = {}
server_set_var = {
    'host' : {
        'display' : 'Host',
        'require' : 'conv',
        'default' : '0.0.0.0'
    }, 'port' : {
        'display' : 'Port',
        'require' : 'conv',
        'default' : '3000'
    }, 'language' : {
        'display' : 'Language',
        'require' : 'select',
        'default' : 'ko-KR',
        'list' : ['ko-KR', 'en-US']
    }, 'markup' : {
        'display' : 'Markup',
        'require' : 'select',
        'default' : 'namumark',
        'list' : ['namumark', 'custom', 'raw']
    }, 'encode' : {
        'display' : 'Encryption method',
        'require' : 'select',
        'default' : 'sha3',
        'list' : ['sha3', 'sha256']
    }
}
server_set_env = {
    'host' : os.getenv('NAMU_HOST'),
    'port' : os.getenv('NAMU_PORT'),
    'language' : os.getenv('NAMU_LANG'),
    'markup' : os.getenv('NAMU_MARKUP'),
    'encode' : os.getenv('NAMU_ENCRYPT')
}
for i in server_set_var:
    curs.execute(db_change('select data from other where name = ?'), [i])
    server_set_val = curs.fetchall()
    if server_set_val:
        server_set_val = server_set_val[0][0]
    elif server_set_env[i] != None:
        server_set_val = server_set_env[i]
    else:
        if 'list' in server_set_var[i]:
            print(server_set_var[i]['display'] + ' (' + server_set_var[i]['default'] + ') [' + ', '.join(server_set_var[i]['list']) + ']' + ' : ', end = '')
        else:
            print(server_set_var[i]['display'] + ' (' + server_set_var[i]['default'] + ') : ', end = '')

        server_set_val = input()
        if server_set_val == '':
            server_set_val = server_set_var[i]['default']
        elif server_set_var[i]['require'] == 'select':
            if not server_set_val in server_set_var[i]['list']:
                server_set_val = server_set_var[i]['default']

        curs.execute(db_change('insert into other (name, data) values (?, ?)'), [i, server_set_val])

    print(server_set_var[i]['display'] + ' : ' + server_set_val)

    server_set[i] = server_set_val

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

print('Now running... http://localhost:' + server_set['port'])
conn.commit()

# Init-custom
if os.path.exists('custom.py'):
    from custom import custom_run
    custom_run(load_db.db_get(), app)
    
# Func
# Func-inter_wiki
app.add_url_rule(
    rule = '/inter_wiki',
    defaults = { 
        'conn' : load_db.db_get(), 
        'tool' : 'inter_wiki' 
    }, 
    view_func = inter_wiki
)

app.add_url_rule(
    rule = '/inter_wiki/del/<name>',
    defaults = { 
        'conn' : load_db.db_get(), 
        'tool' : 'del_inter_wiki',
        'name' : 'Test'
    }, 
    view_func = inter_wiki_del
)

@app.route('/edit_top')
def inter_wiki_edit_top():
    return inter_wiki(load_db.db_get(), 'edit_top')

@app.route('/edit_top/del/<name>')
def inter_wiki_edit_top_del(name = 'Test'):
    return inter_wiki_del(load_db.db_get(), 'del_edit_top', name)

@app.route('/image_license')
def inter_wiki_image_license():
    return inter_wiki(load_db.db_get(), 'image_license')

@app.route('/image_license/del/<name>')
def inter_wiki_image_license_del(name = 'Test'):
    return inter_wiki_del(load_db.db_get(), 'del_image_license', name)

@app.route('/edit_filter')
def inter_wiki_edit_filter():
    return inter_wiki(load_db.db_get(), 'edit_filter')

@app.route('/edit_filter/del/<name>')
def inter_wiki_edit_filter_del(name = 'Test'):
    return inter_wiki_del(load_db.db_get(), 'del_edit_filter', name)

@app.route('/email_filter')
def inter_wiki_email_filter():
    return inter_wiki(load_db.db_get(), 'email_filter')

@app.route('/email_filter/del/<name>')
def inter_wiki_email_filter_del(name = 'Test'):
    return inter_wiki_del(load_db.db_get(), 'del_email_filter', name)

@app.route('/file_filter')
def inter_wiki_file_filter():
    return inter_wiki(load_db.db_get(), 'file_filter')

@app.route('/file_filter/del/<name>')
def inter_wiki_file_filter_del(name = 'Test'):
    return inter_wiki_del(load_db.db_get(), 'del_file_filter', name)

@app.route('/name_filter')
def inter_wiki_name_filter():
    return inter_wiki(load_db.db_get(), 'name_filter')

@app.route('/name_filter/del/<name>')
def inter_wiki_name_filter_del(name = 'Test'):
    return inter_wiki_del(load_db.db_get(), 'del_name_filter', name)

@app.route('/extension_filter')
def inter_wiki_extension_filter():
    return inter_wiki(load_db.db_get(), 'extension_filter')

@app.route('/extension_filter/del/<name>')
def inter_wiki_extension_filter_del(name = 'Test'):
    return inter_wiki_del(load_db.db_get(), 'del_extension_filter', name)

@app.route('/<regex("(?:inter_wiki|edit_top|image_license|(?:edit|email|file|name|extension)_filter)"):tools>/add', methods = ['POST', 'GET'])
@app.route('/<regex("(?:inter_wiki|edit_top|image_license|(?:edit|email|file|name|extension)_filter)"):tools>/add/<name>', methods = ['POST', 'GET'])
def inter_wiki_plus(tools = None, name = None):
    return inter_wiki_plus_2(load_db.db_get(), 'plus_' + tools, name)

# Func-list
# /list/topic/open
@app.route('/not_close_topic')
def list_not_close_topic():
    return list_not_close_topic_2(load_db.db_get())

# /list/document/old
@app.route('/old_page')
def list_old_page():
    return list_old_page_2(load_db.db_get())

# /list/document/acl
@app.route('/acl_list')
def list_acl():
    return list_acl_2(load_db.db_get())

# /list/document/acl/add
@app.route('/acl/<everything:name>', methods = ['POST', 'GET'])
def give_acl(name = None):
    return give_acl_2(load_db.db_get(), name)

# /list/document/need
@app.route('/please')
def list_please():
    return list_please_2(load_db.db_get())

# /list/document/all
@app.route('/title_index')
def list_title_index():
    return list_title_index_2(load_db.db_get())

# /list/document/long
@app.route('/long_page')
def list_long_page():
    return list_long_page_2(load_db.db_get(), 'long_page')

# /list/document/short
@app.route('/short_page')
def list_short_page():
    return list_long_page_2(load_db.db_get(), 'short_page')

# /list/file
@app.route('/image_file_list')
def list_image_file():
    return list_image_file_2(load_db.db_get())

# /list/admin
# /list/admin/list
@app.route('/admin_list')
def list_admin():
    return list_admin_2(load_db.db_get())

# /list/admin/auth_use
@app.route('/admin_log', methods = ['POST', 'GET'])
def list_admin_use():
    return list_admin_use_2(load_db.db_get())

# /list/user
@app.route('/user_log')
def list_user():
    return list_user_2(load_db.db_get())

# /list/user/check
@app.route('/check/<name>')
def give_user_check(name = None):
    return give_user_check_2(load_db.db_get(), name)
    
# /list/user/check/delete
@app.route('/check_delete', methods = ['POST', 'GET'])
def give_user_check_delete():
    return give_user_check_delete_2(load_db.db_get())

# Func-auth
# /auth/give
# /auth/give/<name>
@app.route('/admin/<name>', methods = ['POST', 'GET'])
def give_admin(name = None):
    return give_admin_2(load_db.db_get(), name)

# /auth/give
# /auth/give/<name>
@app.route('/ban', methods = ['POST', 'GET'])
@app.route('/ban/<name>', methods = ['POST', 'GET'])
def give_user_ban(name = None):
    return give_user_ban_2(load_db.db_get(), name)

# /auth/list
@app.route('/admin_group')
def list_admin_group():
    return list_admin_group_2(load_db.db_get())

# /auth/list/add/<name>
@app.route('/admin_plus/<name>', methods = ['POST', 'GET'])
def give_admin_groups(name = None):
    return give_admin_groups_2(load_db.db_get(), name)

# /auth/list/delete/<name>
@app.route('/delete_admin_group/<name>', methods = ['POST', 'GET'])
def give_delete_admin_group(name = None):
    return give_delete_admin_group_2(load_db.db_get(), name)

# /auth/history
# ongoing 반영 필요
@app.route('/block_log')
@app.route('/block_log/<regex("user"):tool>/<name>')
@app.route('/block_log/<regex("admin"):tool>/<name>')
def recent_block(name = 'Test', tool = 'all'):
    return recent_block_2(load_db.db_get(), name, tool)

# Func-history
@app.route('/recent_change')
@app.route('/recent_changes')
def recent_change(name = None):
    return recent_change_2(load_db.db_get(), name, '')

@app.route('/record/<name>')
def recent_record(name = None):
    return recent_change_2(load_db.db_get(), name, 'record')

@app.route('/history/<everything:name>', methods = ['POST', 'GET'])
def recent_history(name = None):
    return recent_change_2(load_db.db_get(), name, 'history')

@app.route('/history/tool/<int(signed=True):rev>/<everything:name>')
def recent_history_tool(name = 'Test', rev = 1):
    return recent_history_tool_2(load_db.db_get(), name, rev)

@app.route('/history/delete/<int(signed=True):rev>/<everything:name>', methods = ['POST', 'GET'])
def recent_history_delete(name = 'Test', rev = 1):
    return recent_history_delete_2(load_db.db_get(), name, rev)

@app.route('/history/hidden/<int(signed=True):rev>/<everything:name>')
def recent_history_hidden(name = 'Test', rev = 1):
    return recent_history_hidden_2(load_db.db_get(), name, rev)

@app.route('/history/send/<int(signed=True):rev>/<everything:name>', methods = ['POST', 'GET'])
def recent_history_send(name = 'Test', rev = 1):
    return recent_history_send_2(load_db.db_get(), name, rev)

@app.route('/history/reset/<everything:name>', methods = ['POST', 'GET'])
def recent_history_reset(name = 'Test'):
    return recent_history_reset_2(load_db.db_get(), name)

@app.route('/history/add/<everything:name>', methods = ['POST', 'GET'])
def recent_history_add(name = 'Test'):
    return recent_history_add_2(load_db.db_get(), name)

@app.route('/record/reset/<name>', methods = ['POST', 'GET'])
def recent_record_reset(name = 'Test'):
    return recent_record_reset_2(load_db.db_get(), name)

@app.route('/record/topic/<name>')
def recent_record_topic(name = 'Test'):
    return recent_record_topic_2(load_db.db_get(), name)

# 거처를 고심중
@app.route('/app_submit', methods = ['POST', 'GET'])
def recent_app_submit():
    return recent_app_submit_2(load_db.db_get())

# Func-search
@app.route('/search', methods=['POST'])
def search():
    return search_2(load_db.db_get())

@app.route('/goto', methods=['POST'])
@app.route('/goto/<everything:name>', methods=['POST'])
def search_goto(name = 'test'):
    return search_goto_2(load_db.db_get(), name)

@app.route('/search/<everything:name>')
def search_deep(name = 'test'):
    return search_deep_2(load_db.db_get(), name)

# Func-view
@app.route('/xref/<everything:name>')
def view_xref(name = 'Test'):
    return view_xref_2(load_db.db_get(), name)

@app.route('/xref/this/<everything:name>')
def view_xref_this(name = 'Test'):
    return view_xref_2(load_db.db_get(), name, xref_type = '2')

@app.route('/raw/<everything:name>')
@app.route('/thread/<int:topic_num>/raw/<int:num>')
def view_raw(name = None, topic_num = None, num = None):
    return view_raw_2(load_db.db_get(), name, topic_num, num)

@app.route('/diff/<int:num_a>/<int:num_b>/<everything:name>')
def view_diff(name = 'Test', num_a = 1, num_b = 1):
    return view_diff_2(load_db.db_get(), name, num_a, num_b)

@app.route('/down/<everything:name>')
def view_down(name = None):
    return view_down_2(load_db.db_get(), name)

@app.route('/w/<everything:name>/doc_rev/<int:doc_rev>')
@app.route('/w/<everything:name>/doc_from/<everything:doc_from>')
@app.route('/w/<everything:name>')
def view_read(name = 'Test', doc_rev = 0, doc_from = ''):
    return view_read_2(load_db.db_get(), name, doc_rev, doc_from)

# Func-edit
@app.route('/revert/<everything:name>', methods = ['POST', 'GET'])
def edit_revert(name = None):
    return edit_revert_2(load_db.db_get(), name)

@app.route('/edit/<everything:name>', methods = ['POST', 'GET'])
@app.route('/edit/<everything:name>/doc_section/<int:section>', methods = ['POST', 'GET'])
def edit(name = 'Test', section = 0):
    return edit_2(load_db.db_get(), name, section)

@app.route('/backlink_reset/<everything:name>')
def edit_backlink_reset(name = 'Test'):
    return edit_backlink_reset_2(load_db.db_get(), name)

@app.route('/delete/<everything:name>', methods = ['POST', 'GET'])
def edit_delete(name = None):
    return edit_delete_2(load_db.db_get(), name)

@app.route('/delete/doc_file/<everything:name>', methods = ['POST', 'GET'])
def edit_delete_file(name = 'test.jpg'):
    return edit_delete_file_2(load_db.db_get(), name)

@app.route('/delete/doc_mutiple', methods = ['POST', 'GET'])
def edit_delete_mutiple():
    return edit_delete_mutiple_2(load_db.db_get())

@app.route('/move/<everything:name>', methods = ['POST', 'GET'])
def edit_move(name = None):
    return edit_move_2(load_db.db_get(), name)

# Func-topic
@app.route('/recent_discuss')
def recent_discuss():
    return recent_discuss_2(load_db.db_get())

@app.route('/thread/<int:topic_num>/b/<int:num>')
def topic_block(topic_num = 1, num = 1):
    return topic_block_2(load_db.db_get(), topic_num, num)

@app.route('/thread/<int:topic_num>/notice/<int:num>')
def topic_top(topic_num = 1, num = 1):
    return topic_top_2(load_db.db_get(), topic_num, num)

@app.route('/thread/<int:topic_num>/setting', methods = ['POST', 'GET'])
def topic_stop(topic_num = 1):
    return topic_stop_2(load_db.db_get(), topic_num)

@app.route('/thread/<int:topic_num>/acl', methods = ['POST', 'GET'])
def topic_acl(topic_num = 1):
    return topic_acl_2(load_db.db_get(), topic_num)

@app.route('/thread/<int:topic_num>/delete', methods = ['POST', 'GET'])
def topic_delete(topic_num = 1):
    return topic_delete_2(load_db.db_get(), topic_num)

@app.route('/thread/<int:topic_num>/tool')
def topic_tool(topic_num = 1):
    return topic_tool_2(load_db.db_get(), topic_num)

@app.route('/thread/<int:topic_num>/change', methods = ['POST', 'GET'])
def topic_change(topic_num = 1):
    return topic_change_2(load_db.db_get(), topic_num)

@app.route('/thread/<int:topic_num>/admin/<int:num>')
def topic_admin(topic_num = 1, num = 1):
    return topic_admin_2(load_db.db_get(), topic_num, num)

@app.route('/thread/<int:topic_num>', methods = ['POST', 'GET'])
def topic(topic_num = 1):
    return topic_2(load_db.db_get(), topic_num)

@app.route('/topic/<everything:name>', methods = ['POST', 'GET'])
def topic_close_list(name = 'test'):
    return topic_close_list_2(load_db.db_get(), name)

# Func-user
@app.route('/change', methods = ['POST', 'GET'])
def user_setting():
    return user_setting_2(load_db.db_get(), server_set_var)

@app.route('/change/email', methods = ['POST', 'GET'])
def user_setting_email():
    return user_setting_email_2(load_db.db_get())

@app.route('/change/email/check', methods = ['POST', 'GET'])
def user_setting_email_check():
    return user_setting_email_check_2(load_db.db_get())

@app.route('/change/pw', methods = ['POST', 'GET'])
def user_setting_pw_change():
    return user_setting_pw_change_2(load_db.db_get())

@app.route('/change/head', methods=['GET', 'POST'])
def user_setting_head():
    return user_setting_head_2(load_db.db_get())

@app.route('/user')
@app.route('/user/<name>')
def user_info(name = ''):
    return user_info_2(load_db.db_get(), name)

@app.route('/count')
@app.route('/count/<name>')
def user_count_edit(name = None):
    return user_count_edit_2(load_db.db_get(), name)
    
@app.route('/alarm')
def user_alarm():
    return user_alarm_2(load_db.db_get())

@app.route('/alarm/delete')
def user_alarm_del():
    return user_alarm_del_2(load_db.db_get())
    
@app.route('/watch_list')
def user_watch_list():
    return user_watch_list_2(load_db.db_get(), 'watch_list')

@app.route('/watch_list/<everything:name>')
def user_watch_list_name(name = 'Test'):
    return user_watch_list_name_2(load_db.db_get(), 'watch_list', name)

@app.route('/star_doc')
def user_star_doc():
    return user_watch_list_2(load_db.db_get(), 'star_doc')

@app.route('/star_doc/<everything:name>')
def user_star_doc_name(name = 'Test'):
    return user_watch_list_name_2(load_db.db_get(), 'star_doc', name)

# Func-login
# 개편 예정

# login -> login/2fa -> login/2fa/email with login_id
# register -> register/email -> regiter/email/check with reg_id
# pass_find -> pass_find/email with find_id

@app.route('/login', methods = ['POST', 'GET'])
def login_login():
    return login_login_2(load_db.db_get())

@app.route('/login/2fa', methods = ['POST', 'GET'])
def login_login_2fa():
    return login_login_2fa_2(load_db.db_get())

'''
@app.route('/login/2fa/email', methods = ['POST', 'GET'])
def login_2fa_email():
    return login_login_2fa_email_2(load_db.db_get())
'''

@app.route('/register', methods = ['POST', 'GET'])
def login_register():
    return login_register_2(load_db.db_get())

@app.route('/register/email', methods = ['POST', 'GET'])
def login_register_email():
    return login_register_email_2(load_db.db_get())

@app.route('/register/email/check', methods = ['POST', 'GET'])
def login_register_email_check():
    return login_register_email_check_2(load_db.db_get())

@app.route('/register/submit', methods = ['POST', 'GET'])
def login_register_submit():
    return login_register_submit_2(load_db.db_get())

# 개편 필요
@app.route('/pass_find', methods = ['POST', 'GET'])
def login_pass_find():
    return login_pass_find_2(load_db.db_get(), 'pass_find')

@app.route('/pass_find/email', methods = ['POST', 'GET'])
def login_pass_find_email():
    return login_pass_find_email_2(load_db.db_get(), 'check_key')

@app.route('/logout')
def login_logout():
    return login_logout_2(load_db.db_get())

# Func-vote
@app.route('/vote/<int:num>', methods = ['POST', 'GET'])
def vote_select(num = 1):
    return vote_select_2(load_db.db_get(), str(num))

@app.route('/vote/end/<int:num>')
def vote_end(num = 1):
    return vote_end_2(load_db.db_get(), str(num))

@app.route('/vote/close/<int:num>')
def vote_close(num = 1):
    return vote_close_2(load_db.db_get(), str(num))

@app.route('/vote')
@app.route('/vote/list')
@app.route('/vote/list/<int:num>')
def vote_list(num = 1):
    return vote_list_2(load_db.db_get(), 'normal', num)

@app.route('/vote/list/close')
@app.route('/vote/list/close/<int:num>')
def vote_list_close(num = 1):
    return vote_list_2(load_db.db_get(), 'close', num)

@app.route('/vote/add', methods = ['POST', 'GET'])
def vote_add():
    return vote_add_2(load_db.db_get())

# Func-api
@app.route('/api/w/<everything:name>', methods = ['POST', 'GET'])
def api_w(name = ''):
    return api_w_2(load_db.db_get(), name)

@app.route('/api/raw/<everything:name>')
def api_raw(name = ''):
    return api_raw_2(load_db.db_get(), name)

@app.route('/api/version')
def api_version():
    return api_version_2(load_db.db_get(), version_list)

@app.route('/api/skin_info')
@app.route('/api/skin_info/<name>')
def api_skin_info(name = ''):
    return api_skin_info_2(load_db.db_get(), name)

@app.route('/api/markup')
def api_markup():
    return api_markup_2(load_db.db_get())

@app.route('/api/user_info/<name>')
def api_user_info(name = ''):
    return api_user_info_2(load_db.db_get(), name)

@app.route('/api/thread/<topic_num>')
def api_topic_sub(name = '', topic_num = 1):
    return api_topic_sub_2(load_db.db_get(), topic_num)

@app.route('/api/search/<name>')
def api_search(name = ''):
    return api_search_2(load_db.db_get(), name)

@app.route('/api/recent_changes')
def api_recent_change():
    return api_recent_change_2(load_db.db_get())

@app.route('/api/recent_discuss')
@app.route('/api/recent_discuss/<int:num>')
def api_recent_discuss(num = 10):
    return api_recent_discuss_2(load_db.db_get(), num, 'normal')

@app.route('/api/recent_discuss/stop')
@app.route('/api/recent_discuss/<int:num>/stop')
def api_recent_discuss_stop(num = 10):
    return api_recent_discuss_2(load_db.db_get(), num, 'stop')

@app.route('/api/recent_discuss/all')
@app.route('/api/recent_discuss/<int:num>/all')
def api_recent_discuss_all(num = 10):
    return api_recent_discuss_2(load_db.db_get(), num, 'all')

@app.route('/api/sha224/<everything:name>', methods = ['POST', 'GET'])
def api_sha224(name = 'test'):
    return api_sha224_2(load_db.db_get(), name)

@app.route('/api/title_index')
def api_title_index():
    return api_title_index_2(load_db.db_get())

@app.route('/api/image/<everything:name>', methods = ['POST', 'GET'])
def api_image_view(name = ''):
    return api_image_view_2(load_db.db_get(), name)

@app.route('/api/sitemap.xml')
def api_sitemap():
    return api_sitemap_2(load_db.db_get())

# Func-main
# 여기도 전반적인 조정 시행 예정
@app.route('/restart', methods = ['POST', 'GET'])
def main_restart():
    return main_restart_2(load_db.db_get())

@app.route('/update', methods=['GET', 'POST'])
def main_update():
    return main_update_2(load_db.db_get(), version_list['beta']['r_ver'])

@app.route('/random')
def main_title_random():
    return main_title_random_2(load_db.db_get())

@app.route('/upload', methods=['GET', 'POST'])
def main_upload():
    return main_upload_2(load_db.db_get())

@app.route('/setting')
@app.route('/setting/<int:num>', methods = ['POST', 'GET'])
def setting(num = 0):
    return main_setting_2(load_db.db_get(), num, set_data['db_type'])

@app.route('/other')
def main_other():
    return main_other_2(load_db.db_get())

@app.route('/manager', methods = ['POST', 'GET'])
@app.route('/manager/<int:num>', methods = ['POST', 'GET'])
def main_manager(num = 1):
    return main_manager_2(load_db.db_get(), num)

@app.route('/image/<everything:name>')
def main_image_view(name = None):
    return main_image_view_2(load_db.db_get(), name)

@app.route('/skin_set')
@app.route('/main_skin_set')
def main_skin_set():
    return main_skin_set_2(load_db.db_get())

@app.route('/views/<everything:name>')
def main_views(name = None):
    return main_views_2(load_db.db_get(), name)

@app.route('/test_func')
def main_test_func():
    return main_test_func_2(load_db.db_get())

@app.route('/shutdown', methods = ['POST', 'GET'])
def main_shutdown():
    return main_shutdown_2(load_db.db_get())

@app.route('/easter_egg.xml')
def main_easter_egg():
    return main_easter_egg_2(load_db.db_get())

@app.route('/<regex("[^.]+\.(?:txt|xml)"):data>')
def main_file(data = ''):
    return main_file_2(load_db.db_get(), data)

@app.errorhandler(404)
def main_error_404(e):
    return main_error_404_2(load_db.db_get())
    
if __name__ == "__main__":
    waitress.serve(
        app,
        host = server_set['host'],
        port = int(server_set['port']), 
    )