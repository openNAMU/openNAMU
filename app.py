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
                data = load_lang('oauth_setting_not_found'), 
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
            
        if os.path.exists(os.path.join(app_var['path_data_image'], e_data)):
            os.remove(os.path.join(app_var['path_data_image'], e_data))
            
            data.save(os.path.join(app_var['path_data_image'], e_data))
        else:
            data.save(os.path.join(app_var['path_data_image'], e_data))
            
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
