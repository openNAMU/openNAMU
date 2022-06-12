# Init
import os
import re

from route.tool.func import *
# from route import *

for i_data in os.listdir("route"):
    f_src = re.search(r"(.+)\.py$", i_data)
    f_src = f_src.group(1) if f_src else ""
    
    if not f_src in ('', '__init__'):
        try:
            exec(
                "from route." + f_src + " " + 
                "import " + f_src
            )
        except:
            try:
                exec(
                    "from route." + f_src + " " + 
                    "import " + f_src + "_2"
                )
            except:
                pass

# Init-Version
version_list = json.loads(open(
    'version.json', 
    encoding = 'utf8'
).read())

# Init-DB
data_db_set = class_check_json()

db_data_get(data_db_set['type'])
do_db_set(data_db_set)
load_db = get_db_connect_old(data_db_set)

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
    create_data = get_db_table_list()
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
    regex = r'.*?'

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

log = logging.getLogger('waitress')
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
server_set_var = get_init_set_list()
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
if data_db_set['type'] == 'sqlite':
    def back_up(back_time, back_up_where):
        print('----')

        try:
            shutil.copyfile(
                data_db_set['name'] + '.db', 
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
            back_up_where = 'back_' + data_db_set['name'] + '.db'

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
app.route('/inter_wiki', defaults = { 'tool' : 'inter_wiki' })(filter_inter_wiki)
app.route('/inter_wiki/del/<everything:name>', defaults = { 'tool' : 'del_inter_wiki' })(filter_inter_wiki_delete)
app.route('/inter_wiki/add', methods = ['POST', 'GET'], defaults = { 'tool' : 'plus_inter_wiki' })(filter_inter_wiki_add)
app.route('/inter_wiki/add/<everything:name>', methods = ['POST', 'GET'], defaults = { 'tool' : 'plus_inter_wiki' })(filter_inter_wiki_add)

app.route('/filter/document/list')(filter_document)
app.route('/filter/document/add/<everything:name>', methods = ['POST', 'GET'])(filter_document_add)
app.route('/filter/document/add', methods = ['POST', 'GET'])(filter_document_add)
app.route('/filter/document/del/<name>')(filter_document_delete)

app.route('/edit_top', defaults = { 'tool' : 'edit_top' })(filter_inter_wiki)
app.route('/edit_top/del/<everything:name>', defaults = { 'tool' : 'del_edit_top' })(filter_inter_wiki_delete)
app.route('/edit_top/add', methods = ['POST', 'GET'], defaults = { 'tool' : 'plus_edit_top' })(filter_inter_wiki_add)

app.route('/image_license', defaults = { 'tool' : 'image_license' })(filter_inter_wiki)
app.route('/image_license/del/<everything:name>', defaults = { 'tool' : 'del_image_license' })(filter_inter_wiki_delete)
app.route('/image_license/add', methods = ['POST', 'GET'], defaults = { 'tool' : 'plus_image_license' })(filter_inter_wiki_add)

app.route('/edit_filter', defaults = { 'tool' : 'edit_filter' })(filter_inter_wiki)
app.route('/edit_filter/del/<everything:name>', defaults = { 'tool' : 'del_edit_filter' })(filter_inter_wiki_delete)
app.route('/edit_filter/add', methods = ['POST', 'GET'], defaults = { 'tool' : 'plus_edit_filter' })(filter_inter_wiki_add)
app.route('/edit_filter/add/<everything:name>', methods = ['POST', 'GET'], defaults = { 'tool' : 'plus_edit_filter' })(filter_inter_wiki_add)

app.route('/email_filter', defaults = { 'tool' : 'email_filter' })(filter_inter_wiki)
app.route('/email_filter/del/<everything:name>', defaults = { 'tool' : 'del_email_filter' })(filter_inter_wiki_delete)
app.route('/email_filter/add', methods = ['POST', 'GET'], defaults = { 'tool' : 'plus_email_filter' })(filter_inter_wiki_add)

app.route('/file_filter', defaults = { 'tool' : 'file_filter' })(filter_inter_wiki)
app.route('/file_filter/del/<everything:name>', defaults = { 'tool' : 'del_file_filter' })(filter_inter_wiki_delete)
app.route('/file_filter/add', methods = ['POST', 'GET'], defaults = { 'tool' : 'plus_file_filter' })(filter_inter_wiki_add)

app.route('/name_filter', defaults = { 'tool' : 'name_filter' })(filter_inter_wiki)
app.route('/name_filter/del/<everything:name>', defaults = { 'tool' : 'del_name_filter' })(filter_inter_wiki_delete)
app.route('/name_filter/add', methods = ['POST', 'GET'], defaults = { 'tool' : 'plus_name_filter' })(filter_inter_wiki_add)

app.route('/extension_filter', defaults = { 'tool' : 'extension_filter' })(filter_inter_wiki)
app.route('/extension_filter/del/<everything:name>', defaults = { 'tool' : 'del_extension_filter' })(filter_inter_wiki_delete)
app.route('/extension_filter/add', methods = ['POST', 'GET'], defaults = { 'tool' : 'plus_extension_filter' })(filter_inter_wiki_add)

# Func-list
# /list/document/old
app.route('/old_page')(list_old_page)

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

@app.route('/app_submit', methods = ['POST', 'GET'])
def recent_app_submit():
    return recent_app_submit_2(load_db.db_get())

# /auth/history
# ongoing 반영 필요
@app.route('/block_log')
@app.route('/block_log/<regex("user"):tool>/<name>')
@app.route('/block_log/<regex("admin"):tool>/<name>')
def recent_block(name = 'Test', tool = 'all'):
    return recent_block_2(load_db.db_get(), name, tool)

# Func-history
app.route('/recent_change')(recent_change)
app.route('/recent_changes')(recent_change)

app.route('/record/<name>', defaults = { 'tool' : 'record' })(recent_change)
app.route('/record/reset/<name>', methods = ['POST', 'GET'])(recent_record_reset)
app.route('/record/topic/<name>')(recent_record_topic)

app.route('/history/<everything:name>', defaults = { 'tool' : 'history' }, methods = ['POST', 'GET'])(recent_change)
app.route('/history_tool/<int(signed = True):rev>/<everything:name>')(recent_history_tool)
app.route('/history_delete/<int(signed = True):rev>/<everything:name>', methods = ['POST', 'GET'])(recent_history_delete)
app.route('/history_hidden/<int(signed = True):rev>/<everything:name>')(recent_history_hidden)
app.route('/history_send/<int(signed = True):rev>/<everything:name>', methods = ['POST', 'GET'])(recent_history_send)
app.route('/history_reset/<everything:name>', methods = ['POST', 'GET'])(recent_history_reset)
app.route('/history_add/<everything:name>', methods = ['POST', 'GET'])(recent_history_add)

# Func-view
app.route('/xref/<everything:name>')(view_xref)
app.route('/xref_this/<everything:name>', defaults = { 'xref_type' : 2 })(view_xref)

app.route('/raw/<everything:name>')(view_raw_2)
app.route('/raw_acl/<everything:name>', defaults = { 'doc_acl' : 1 })(view_raw_2)
app.route('/raw_rev/<int:num>/<everything:name>')(view_raw_2)

app.route('/diff/<int(signed = True):num_a>/<int(signed = True):num_b>/<everything:name>')(view_diff)

app.route('/down/<everything:name>')(view_down)

# everything 다음에 추가 붙은 경우에 대해서 재검토 필요 (진행중)
app.route('/w_rev/<int(signed = True):doc_rev>/<everything:name>')(view_read)
app.route('/w_from/<everything:name>', defaults = { 'do_type' : 'from' })(view_read)
app.route('/w/<everything:name>')(view_read)

app.route('/random')(main_func_random)

# Func-edit
app.route('/edit/<everything:name>', methods = ['POST', 'GET'])(edit)
app.route('/edit/<everything:name>/doc_from/<everything:name_load>', methods = ['POST', 'GET'])(edit)
app.route('/edit/<everything:name>/doc_section/<int:section>', methods = ['POST', 'GET'])(edit)

app.route('/upload', methods = ['POST', 'GET'])(main_func_upload)

# 개편 예정
app.route('/xref_reset/<everything:name>')(edit_backlink_reset)

app.route('/delete/<everything:name>', methods = ['POST', 'GET'])(edit_delete)
app.route('/delete_file/<everything:name>', methods = ['POST', 'GET'])(edit_delete_file)
app.route('/delete_mutiple', methods = ['POST', 'GET'])(edit_delete_mutiple)

app.route('/revert/<int:num>/<everything:name>', methods = ['POST', 'GET'])(edit_revert)

app.route('/move/<everything:name>', methods = ['POST', 'GET'])(edit_move)

# Func-topic
app.route('/recent_discuss', defaults = { 'tool' : 'normal' })(recent_discuss)
app.route('/recent_discuss/close', defaults = { 'tool' : 'close' })(recent_discuss)
app.route('/recent_discuss/open', defaults = { 'tool' : 'open' })(recent_discuss)

app.route('/thread/<int:topic_num>', methods = ['POST', 'GET'])(topic)
app.route('/topic/<everything:name>', methods = ['POST', 'GET'])(topic_list)

app.route('/thread/<int:topic_num>/tool')(topic_tool)
app.route('/thread/<int:topic_num>/setting', methods = ['POST', 'GET'])(topic_tool_setting)
app.route('/thread/<int:topic_num>/acl', methods = ['POST', 'GET'])(topic_tool_acl)
app.route('/thread/<int:topic_num>/delete', methods = ['POST', 'GET'])(topic_tool_delete)
app.route('/thread/<int:topic_num>/change', methods = ['POST', 'GET'])(topic_tool_change)

app.route('/thread/<int:topic_num>/comment/<int:num>/tool')(topic_comment_tool)
app.route('/thread/<int:topic_num>/comment/<int:num>/notice')(topic_comment_notice)
app.route('/thread/<int:topic_num>/comment/<int:num>/blind')(topic_comment_blind)
app.route('/thread/<int:topic_num>/comment/<int:num>/raw')(view_raw_2)
app.route('/thread/<int:topic_num>/comment/<int:num>/delete', methods = ['POST', 'GET'])(topic_comment_delete)

# Func-user
app.route('/change', methods = ['POST', 'GET'])(user_setting)
app.route('/change/key')(user_setting_key)
app.route('/change/key/delete')(user_setting_key_delete)
app.route('/change/pw', methods = ['POST', 'GET'])(user_setting_pw)
app.route('/change/head', methods=['GET', 'POST'])(user_setting_head)
app.route('/change/skin_set')(user_setting_skin_set)
app.route('/change/skin_set/main')(user_setting_skin_set)

app.route('/user')(user_info)
app.route('/user/<name>')(user_info)

app.route('/challenge')(user_challenge)

app.route('/count')(user_count)
app.route('/count/<name>')(user_count)

app.route('/alarm')(user_alarm)
app.route('/alarm/delete')(user_alarm_delete)

app.route('/watch_list', defaults = { 'tool' : 'watch_list' })(user_watch_list)
app.route('/watch_list/<everything:name>', defaults = { 'tool' : 'watch_list' })(user_watch_list_name)

app.route('/star_doc', defaults = { 'tool' : 'star_doc' })(user_watch_list)
app.route('/star_doc/<everything:name>', defaults = { 'tool' : 'star_doc' })(user_watch_list_name)

# 하위 호환용 S
# /change/skin_set
app.route('/skin_set')(user_setting_skin_set)
# 하위 호환용 E

# 개편 보류중 S
@app.route('/change/email', methods = ['POST', 'GET'])
def user_setting_email():
    return user_setting_email_2(load_db.db_get())

app.route('/change/email/delete')(user_setting_email_delete)

@app.route('/change/email/check', methods = ['POST', 'GET'])
def user_setting_email_check():
    return user_setting_email_check_2(load_db.db_get())
# 개편 보류중 E

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

app.route('/login/find')(login_find)
app.route('/login/find/key', methods = ['POST', 'GET'])(login_find_key)
app.route('/login/find/email', methods = ['POST', 'GET'], defaults = { 'tool' : 'pass_find' })(login_find_email)
app.route('/login/find/email/check', methods = ['POST', 'GET'], defaults = { 'tool' : 'check_key' })(login_find_email_check)
app.route('/logout')(login_logout)

# Func-vote
app.route('/vote/<int:num>', methods = ['POST', 'GET'])(vote_select)
app.route('/vote/end/<int:num>')(vote_end)
app.route('/vote/close/<int:num>')(vote_close)
app.route('/vote', defaults = { 'list_type' : 'normal' })(vote_list)
app.route('/vote/list', defaults = { 'list_type' : 'normal' })(vote_list)
app.route('/vote/list/<int:num>', defaults = { 'list_type' : 'normal' })(vote_list)
app.route('/vote/list/close', defaults = { 'list_type' : 'close' })(vote_list)
app.route('/vote/list/close/<int:num>', defaults = { 'list_type' : 'close' })(vote_list)
app.route('/vote/add', methods = ['POST', 'GET'])(vote_add)

# Func-api
app.route('/api/w/<everything:name>/doc_tool/<tool>/doc_rev/<int(signed = True):rev>')(api_w)
app.route('/api/w/<everything:name>/doc_tool/<tool>', methods = ['POST', 'GET'])(api_w)
app.route('/api/w/<everything:name>', methods = ['GET', 'POST'])(api_w)
app.route('/api/raw/<everything:name>')(api_raw)

app.route('/api/version', defaults = { 'version_list' : version_list })(api_version)
app.route('/api/skin_info')(api_skin_info)
app.route('/api/skin_info/<name>')(api_skin_info)
app.route('/api/markup')(api_markup)
app.route('/api/user_info/<name>', methods = ['POST', 'GET'])(api_user_info)
app.route('/api/setting/<name>')(api_setting)

app.route('/api/thread/<int:topic_num>/<tool>/<int:num>')(api_topic_sub)
app.route('/api/thread/<int:topic_num>/<tool>')(api_topic_sub)
app.route('/api/thread/<int:topic_num>')(api_topic_sub)

app.route('/api/search/<everything:name>/doc_num/<int:num>/<int:page>')(api_search)
app.route('/api/search/<everything:name>')(api_search)

app.route('/api/recent_change/<int:num>')(api_recent_change)
app.route('/api/recent_change')(api_recent_change)
# recent_changes -> recent_change
app.route('/api/recent_changes')(api_recent_change)

app.route('/api/recent_discuss/<get_type>/<int:num>')(api_recent_discuss)
app.route('/api/recent_discuss/<int:num>')(api_recent_discuss)
app.route('/api/recent_discuss')(api_recent_discuss)

app.route('/api/sha224/<everything:data>', methods = ['POST', 'GET'])(api_sha224)
app.route('/api/title_index')(api_title_index)
app.route('/api/image/<everything:name>', methods = ['POST', 'GET'])(api_image_view)
# 이건 API 영역이 아닌 것 같아서 고심 중
app.route('/api/sitemap.xml')(api_sitemap)

# Func-main
# 여기도 전반적인 조정 시행 예정
app.route('/other')(main_tool_other)
app.route('/manager', methods = ['POST', 'GET'])(main_tool_admin)
app.route('/manager/<int:num>', methods = ['POST', 'GET'])(main_tool_admin)
app.route('/manager/<int:num>/<add_2>', methods = ['POST', 'GET'])(main_tool_admin)

app.route('/search', methods=['POST'])(main_search)
app.route('/search/<everything:name>')(main_search_deep)
app.route('/goto', methods=['POST'])(main_search_goto)
app.route('/goto/<everything:name>', methods=['POST'])(main_search_goto)

app.route('/setting')(main_func_setting)
app.route('/setting/main', defaults = { 'db_set' : data_db_set['type'] }, methods = ['POST', 'GET'])(main_func_setting_main)
app.route('/setting/main/logo', methods = ['POST', 'GET'])(main_func_setting_main_logo)
app.route('/setting/phrase', methods = ['POST', 'GET'])(main_func_setting_phrase)
app.route('/setting/head', defaults = { 'num' : 3 }, methods = ['POST', 'GET'])(main_func_setting_head)
app.route('/setting/head/<skin_name>', defaults = { 'num' : 3 }, methods = ['POST', 'GET'])(main_func_setting_head)
app.route('/setting/body/top', defaults = { 'num' : 4 }, methods = ['POST', 'GET'])(main_func_setting_head)
app.route('/setting/body/top/<skin_name>', defaults = { 'num' : 4 }, methods = ['POST', 'GET'])(main_func_setting_head)
app.route('/setting/body/bottom', defaults = { 'num' : 7 }, methods = ['POST', 'GET'])(main_func_setting_head)
app.route('/setting/body/bottom/<skin_name>', defaults = { 'num' : 7 }, methods = ['POST', 'GET'])(main_func_setting_head)
app.route('/setting/robot', methods = ['POST', 'GET'])(main_func_setting_robot)
app.route('/setting/external', methods = ['POST', 'GET'])(main_func_setting_external)
app.route('/setting/acl', methods = ['POST', 'GET'])(main_func_setting_acl)

# views -> view
app.route('/view/<everything:name>')(main_view)
app.route('/views/<everything:name>')(main_view)
app.route('/image/<everything:name>')(main_view_image)
# 조정 계획 중
app.route('/<regex("[^.]+\.(?:txt|xml)"):data>')(main_view_file)

app.route('/shutdown', methods = ['POST', 'GET'])(main_sys_shutdown)
app.route('/restart', methods = ['POST', 'GET'])(main_sys_restart)
app.route('/update', methods = ['POST', 'GET'])(main_sys_update)

app.errorhandler(404)(main_error_404)

if __name__ == "__main__":
    waitress.serve(
        app,
        host = server_set['host'],
        port = int(server_set['port']),
        threads = 1
    )