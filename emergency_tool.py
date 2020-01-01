from route.tool.func import *

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
            break
    except:
        if os.getenv('NAMU_DB') != None or os.getenv('NAMU_DB_TYPE') != None:
            set_data = {
                "db" : os.getenv('NAMU_DB') if os.getenv('NAMU_DB') else 'data',
                "db_type" : os.getenv('NAMU_DB_TYPE') if os.getenv('NAMU_DB_TYPE') else 'sqlite'
            }

            break
        else:
            new_json = ['', '']
            normal_db_type = ['sqlite', 'mysql']

            print('DB type (sqlite, mysql) : ', end = '')
            new_json[0] = str(input())
            if new_json[0] == '' or not new_json[0] in normal_db_type:
                new_json[0] = 'sqlite'

            all_src = []
            for i_data in os.listdir("."):
                f_src = re.search("(.+)\.db$", i_data)
                if f_src:
                    all_src += [f_src.groups()[0]]

            if all_src != []:
                print('DB name (' + ', '.join(all_src) + ') : ', end = '')
            else:
                print('DB name (data) : ', end = '')

            new_json[1] = str(input())
            if new_json[1] == '':
                new_json[1] = 'data'

            with open('data/set.json', 'w') as f:
                f.write('{ "db" : "' + new_json[1] + '", "db_type" : "' + new_json[0] + '" }')

            set_data = json.loads(open('data/set.json').read())

            break

print('DB name : ' + set_data['db'])
print('DB type : ' + set_data['db_type'])

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