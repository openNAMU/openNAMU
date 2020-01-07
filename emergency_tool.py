from route.tool.func import *

version_list = json.loads(open('version.json').read())

# DB
while 1:
    try:
        set_data = json.loads(open('data/set.json').read())
        if not 'db_type' in set_data:
            try:
                os.remove('data/set.json')
            except:
                print('Please delete set.json')
                print('----')
                raise
        else:
            print('DB name : ' + set_data['db'])
            print('DB type : ' + set_data['db_type'])

            break
    except:
        if os.getenv('NAMU_DB') != None or os.getenv('NAMU_DB_TYPE') != None:
            set_data = {
                "db" : os.getenv('NAMU_DB') if os.getenv('NAMU_DB') else 'data',
                "db_type" : os.getenv('NAMU_DB_TYPE') if os.getenv('NAMU_DB_TYPE') else 'sqlite'
            }

            print('DB name : ' + set_data['db'])
            print('DB type : ' + set_data['db_type'])

            break
        else:
            new_json = ['', '']
            normal_db_type = ['sqlite', 'mysql']

            print('DB type (sqlite) [sqlite, mysql] : ', end = '')
            new_json[0] = str(input())
            if new_json[0] == '' or not new_json[0] in normal_db_type:
                new_json[0] = 'sqlite'

            all_src = []
            for i_data in os.listdir("."):
                f_src = re.search("(.+)\.db$", i_data)
                if f_src:
                    all_src += [f_src.groups()[0]]

            if all_src != [] and new_json[0] != 'mysql':
                print('DB name (data) [' + ', '.join(all_src) + '] : ', end = '')
            else:
                print('DB name (data) : ', end = '')

            new_json[1] = str(input())
            if new_json[1] == '':
                new_json[1] = 'data'

            with open('data/set.json', 'w') as f:
                f.write('{ "db" : "' + new_json[1] + '", "db_type" : "' + new_json[0] + '" }')

            set_data = json.loads(open('data/set.json').read())

            break

db_data_get(set_data['db_type'])

if set_data['db_type'] == 'mysql':
    try:
        set_data_mysql = json.loads(open('data/mysql.json').read())
    except:
        new_json = ['', '']

        while 1:
            print('DB user id : ', end = '')
            new_json[0] = str(input())
            if new_json[0] != '':
                break

        while 1:
            print('DB password : ', end = '')
            new_json[1] = str(input())
            if new_json[1] != '':
                break

        with open('data/mysql.json', 'w') as f:
            f.write('{ "user" : "' + new_json[0] + '", "password" : "' + new_json[1] + '" }')

        set_data_mysql = json.loads(open('data/mysql.json').read())

    conn = pymysql.connect(
        host = 'localhost',
        user = set_data_mysql['user'],
        password = set_data_mysql['password'],
        charset = 'utf8mb4'
    )
    curs = conn.cursor()

    try:
        curs.execute(db_change('create database ? default character set utf8mb4;')%pymysql.escape_string(set_data['db']))
    except:
        pass

    curs.execute(db_change('use ?')%pymysql.escape_string(set_data['db']))
else:
    conn = sqlite3.connect(set_data['db'] + '.db', check_same_thread = False)
    curs = conn.cursor()

load_conn(conn)

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
    'oauth_conn',
    'user_application'
]
for i in create_data['all_data']:
    try:
        curs.execute(db_change('select test from ' + i + ' limit 1'))
    except:
        try:
            curs.execute(db_change('create table ' + i + '(test longtext)'))
        except:
            curs.execute(db_change("alter table " + i + " add test longtext default ''"))

setup_tool = 0
try:
    curs.execute(db_change('select data from other where name = "ver"'))
    ver_set_data = curs.fetchall()
    if not ver_set_data:
        setup_tool = 1
    else:
        if int(version_list['master']['c_ver']) > int(ver_set_data[0][0]):
            setup_tool = 1
except:
    setup_tool = 1

if setup_tool != 0:
    create_data['data'] = ['title', 'data']
    create_data['cache_data'] = ['title', 'data']
    create_data['history'] = ['id', 'title', 'data', 'date', 'ip', 'send', 'leng', 'hide', 'type']
    create_data['rd'] = ['title', 'sub', 'date', 'band', 'stop', 'agree']
    create_data['user'] = ['id', 'pw', 'acl', 'date', 'encode']
    create_data['user_set'] = ['name', 'id', 'data']
    create_data['user_application'] = ['id', 'pw', 'date', 'encode', 'question', 'answer', 'ip', 'ua', 'token', 'email']
    create_data['ban'] = ['block', 'end', 'why', 'band', 'login']
    create_data['topic'] = ['id', 'title', 'sub', 'data', 'date', 'ip', 'block', 'top', 'code']
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
    create_data['inter'] = ['title', 'link', 'icon']
    create_data['html_filter'] = ['html', 'kind', 'plus']
    create_data['oauth_conn'] = ['provider', 'wiki_id', 'sns_id', 'name', 'picture']

    for create_table in create_data['all_data']:
        for create in create_data[create_table]:
            try:
                curs.execute(db_change('select ' + create + ' from ' + create_table + ' limit 1'))
            except:
                curs.execute(db_change("alter table " + create_table + " add " + create + " longtext default ''"))

            try:
                curs.execute(db_change('create index index_' + create_table + '_' + create + ' on ' + create_table + '(' + create + ')'))
            except:
                pass

    update()

# Main
print('----')
print('1. Backlink reset')
print('2. reCAPTCHA delete')
print('3. Ban delete')
print('4. Change host')
print('5. Change port')
print('6. Change skin')
print('7. Change password')
print('8. Reset version')
print('9. Delete set.json')
print('10. Change name')
print('11. Delete mysql.json')
print('12. All title count reset')

print('----')
print('Select : ', end = '')
what_i_do = input()

if what_i_do == '1':
    curs.execute(db_change("delete from back"))
    conn.commit()

    curs.execute(db_change("select title, data from data"))
    data = curs.fetchall()
    num = 0

    for test in data:
        num += 1
        if num % 100 == 0:
            print(num)

        render_do(test[0], test[1], 1, None)
elif what_i_do == '2':
    curs.execute(db_change("delete from other where name = 'recaptcha'"))
    curs.execute(db_change("delete from other where name = 'sec_re'"))
elif what_i_do == '3':
    print('----')
    print('IP or Name : ', end = '')
    user_data = input()

    if re.search("^([0-9]{1,3}\.[0-9]{1,3})$", user_data):
        band = 'O'
    else:
        band = ''

        curs.execute(db_change("insert into rb (block, end, today, blocker, why, band) values (?, ?, ?, ?, ?, ?)"),
            [user_data,
            'release',
            get_time(),
            'tool:emergency',
            '',
            band
        ])
    curs.execute(db_change("delete from ban where block = ?"), [user_data])
elif what_i_do == '4':
    print('----')
    print('Host : ', end = '')
    host = input()

    curs.execute(db_change("update other set data = ? where name = 'host'"), [host])
elif what_i_do == '5':
    print('----')
    print('Port : ', end = '')
    port = int(input())

    curs.execute(db_change("update other set data = ? where name = 'port'"), [port])
elif what_i_do == '6':
    print('----')
    print('Skin name : ', end = '')
    skin = input()

    curs.execute(db_change("update other set data = ? where name = 'skin'"), [skin])
elif what_i_do == '7':
    print('----')
    print('1. sha256')
    print('2. sha3')

    print('----')
    print('Select : ', end = '')
    what_i_do = int(input())

    print('----')
    print('User name : ', end = '')
    user_name = input()

    print('----')
    print('User password : ', end = '')
    user_pw = input()

    if what_i_do == '1':
        hashed = hashlib.sha256(bytes(user_pw, 'utf-8')).hexdigest()
    else:
        if sys.version_info < (3, 6):
            hashed = sha3.sha3_256(bytes(user_pw, 'utf-8')).hexdigest()
        else:
            hashed = hashlib.sha3_256(bytes(user_pw, 'utf-8')).hexdigest()

    curs.execute(db_change("update user set pw = ? where id = ?"), [hashed, user_name])
elif what_i_do == '8':
    curs.execute(db_change("update other set data = '00000' where name = 'ver'"))
elif what_i_do == '9':
    try:
        os.remove('data/set.json')
    except:
        pass
elif what_i_do == '10':
    print('----')
    print('User name : ', end = '')
    user_name = input()

    print('----')
    print('New name : ', end = '')
    new_name = input()

    curs.execute(db_change("update user set id = ? where id = ?"), [new_name, user_name])
elif what_i_do == '11':
    try:
        os.remove('data/mysql.json')
    except:
        pass
else:
    curs.execute(db_change("select count(title) from data"))
    count_data = curs.fetchall()
    if count_data:
        count_data = count_data[0][0]
    else:
        count_data = 0

    curs.execute(db_change('delete from other where name = "count_all_title"'))
    curs.execute(db_change('insert into other (name, data) values ("count_all_title", ?)'), [str(count_data)])

conn.commit()

print('----')
print('OK')