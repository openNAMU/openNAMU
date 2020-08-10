import time
from route.tool.func import *


# DB
version_list = json.loads(open('version.json', encoding='utf8').read())

print('Version : ' + version_list['beta']['r_ver'])
print('DB set version : ' + version_list['beta']['c_ver'])
print('Skin set version : ' + version_list['beta']['s_ver'])
print('----')

while 1:
    try:
        set_data = json.loads(open('data/set.json', encoding='utf8').read())
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
                f_src = re.search(r"(.+)\.db$", i_data)
                if f_src:
                    all_src += [f_src.group(1)]

            if all_src != [] and new_json[0] != 'mysql':
                print('DB name (data) [' + ', '.join(all_src) + '] : ', end = '')
            else:
                print('DB name (data) : ', end = '')

            new_json[1] = str(input())
            if new_json[1] == '':
                new_json[1] = 'data'

            with open('data/set.json', 'w', encoding='utf8') as f:
                f.write('{ "db" : "' + new_json[1] + '", "db_type" : "' + new_json[0] + '" }')

            set_data = json.loads(open('data/set.json', encoding='utf8').read())

            break

db_data_get(set_data['db_type'])

if set_data['db_type'] == 'mysql':
    try:
        set_data_mysql = json.loads(open('data/mysql.json', encoding='utf8').read())
    except:
        new_json = ['', '', '']

        while 1:
            print('DB user ID : ', end = '')
            new_json[0] = str(input())
            if new_json[0] != '':
                break

        while 1:
            print('DB password : ', end = '')
            new_json[1] = str(input())
            if new_json[1] != '':
                break
                
        print('DB host (localhost) : ', end = '')
        new_json[2] = str(input())
        if new_json[2] == '':
            new_json[2] == 'localhost'

        with open('data/mysql.json', 'w', encoding='utf8') as f:
            f.write('{ "user" : "' + new_json[0] + '", "password" : "' + new_json[1] + '", "host" : "' + new_json[2] + '" }')

        set_data_mysql = json.loads(open('data/mysql.json', encoding='utf8').read())

    conn = pymysql.connect(
        host = set_data_mysql['host'] if 'host' in set_data_mysql else 'localhost',
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
    conn = sqlite3.connect(set_data['db'] + '.db')
    curs = conn.cursor()

load_conn(conn)

# Main
print('----')
print('1. Backlink reset')
print('2. reCAPTCHA delete')
print('3. Ban delete')
print('4. Change host')
print('5. Change port')
print('6. Change skin')
print('7. Change password')
print('8. Change version')
print('9. Delete set.json')
print('10. Change name')
print('11. Delete mysql.json')
print('12. All title count reset')
print('13. Cache data reset')
print('14. Delete Main <HEAD>')
print('15. Give owner')
print('16. Delete 2FA password')

print('----')
print('Select : ', end = '')
what_i_do = input()

if what_i_do == '1':
    print('----')
    print('All delete (Y) [Y, N] : ', end = '')
    go_num = input()
    if not go_num == 'N':
        curs.execute(db_change("delete from back"))
        conn.commit()

    print('----')
    print('Count (100) : ', end = '')
    try:
        go_num = int(input())
    except:
        go_num = 100

    num = 0

    print('----')
    print('Load...')

    curs.execute(db_change("select title from data d where not exists (select title from back where link = d.title)"))
    title = curs.fetchall()

    print('----')
    print('Rest : ' + str(len(title)))
    print('Start : ' + title[0][0])
    time.sleep(1)
    print('----')

    for name in title:
        num += 1
        if num % go_num == 0:
            print(str(num) + ' : ' + name[0])

        if num % 100 == 0:
            conn.commit()

        curs.execute(db_change("select data from data where title = ?"), [name[0]])
        data = curs.fetchall()
        render_do(name[0], data[0][0], 3, None)
elif what_i_do == '2':
    curs.execute(db_change("delete from other where name = 'recaptcha'"))
    curs.execute(db_change("delete from other where name = 'sec_re'"))
elif what_i_do == '3':
    print('----')
    print('IP or Name : ', end = '')
    user_data = input()

    curs.execute(db_change("insert into rb (block, end, today, blocker, why, band) values (?, ?, ?, ?, ?, ?)"),
        [user_data,
        'release',
        get_time(),
        'tool:emergency',
        '',
        ''
    ])

    curs.execute(db_change("update rb set ongoing = '' where block = ?"), [user_data])
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
    print('----')
    print('Insert version (0000000) : ', end = '')
    new_ver = input()
    
    if new_ver == '':
        new_ver == '0000000'

    curs.execute(db_change("update other set data = ? where name = 'ver'"), [new_ver])
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
elif what_i_do == '12':
    curs.execute(db_change("select count(*) from data"))
    count_data = curs.fetchall()
    if count_data:
        count_data = count_data[0][0]
    else:
        count_data = 0

    curs.execute(db_change('delete from other where name = "count_all_title"'))
    curs.execute(db_change('insert into other (name, data) values ("count_all_title", ?)'), [str(count_data)])
elif what_i_do == '13':
    curs.execute(db_change('delete from cache_data'))
elif what_i_do == '14':
    curs.execute(db_change('delete from other where name = "head"'))
elif what_i_do == '15':
    print('----')
    print('User name : ', end = '')
    user_name = input()

    curs.execute(db_change("update user set acl = 'owner' where id = ?"), [user_name])
else:
    print('----')
    print('User name : ', end = '')
    user_name = input()

    curs.execute(db_change('select data from user_set where name = "2fa" and id = ?'), [user_name])
    if curs.fetchall():
        curs.execute(db_change("update user_set set data = '' where name = '2fa' and id = ?"), [user_name])

conn.commit()

print('----')
print('OK')