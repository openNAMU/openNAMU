# Load
import time
from route.tool.func import *

while 1:
    print('Load DB (Y) [Y, N] : ', end = '')
    data_db_load = input()
    if data_db_load in ['Y', 'N']:
        break

if data_db_load == 'Y':
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
        if not os.path.exists(os.path.join('data', 'mysql.json')):
            db_set_list = ['user', 'password', 'host', 'port']
            set_data = json.loads(
                open(
                    os.path.join('data', 'mysql.json'),
                    encoding = 'utf8'
                ).read()
            )
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
    load_db = get_db_connect_old(data_db_set)

    conn = load_db.db_load()
    curs = conn.cursor()
else:
    print('----')
    print('You can use [9, 11]')

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
print('17. Change markup')

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

    curs.execute(db_change("" + \
        "select title from data d " + \
        "where not exists (" + \
            "select title from back where link = d.title limit 1" + \
        ")" + \
    ""))
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
        
        get_class_render = class_do_render(conn)
        get_class_render.do_render(name[0], data[0][0], 'backlink', '')
elif what_i_do == '2':
    curs.execute(db_change("delete from other where name = 'recaptcha'"))
    curs.execute(db_change("delete from other where name = 'sec_re'"))
elif what_i_do == '3':
    print('----')
    print('IP or Name : ', end = '')
    user_data = input()

    curs.execute(db_change("" + \
        "insert into rb (block, end, today, blocker, why, band) " + \
        "values (?, ?, ?, ?, ?, ?)"
    ), [
        user_data,
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

    curs.execute(db_change("update user_set set data = ? where id = ? and name = 'pw'"), [
        hashed, 
        user_name
    ])
elif what_i_do == '8':
    print('----')
    print('Insert version (0000000) : ', end = '')
    new_ver = input()
    
    if new_ver == '':
        new_ver == '0000000'

    curs.execute(db_change("update other set data = ? where name = 'ver'"), [new_ver])
elif what_i_do == '9':
    if os.path.exists(os.path.join('data', 'set.json')):
        os.remove(os.path.join('data', 'set.json'))
elif what_i_do == '10':
    print('----')
    print('User name : ', end = '')
    user_name = input()

    print('----')
    print('New name : ', end = '')
    new_name = input()

    curs.execute(db_change("update user_set set id = ? where id = ?"), [
        new_name, 
        user_name
    ])
elif what_i_do == '11':
    if os.path.exists(os.path.join('data', 'mysql.json')):
        os.remove(os.path.join('data', 'mysql.json'))
elif what_i_do == '12':
    curs.execute(db_change("select count(*) from data"))
    count_data = curs.fetchall()
    if count_data:
        count_data = count_data[0][0]
    else:
        count_data = 0

    curs.execute(db_change('delete from other where name = "count_all_title"'))
    curs.execute(db_change('insert into other (name, data) values ("count_all_title", ?)'), [
        str(count_data)
    ])
elif what_i_do == '13':
    curs.execute(db_change('delete from cache_data'))
elif what_i_do == '14':
    curs.execute(db_change('delete from other where name = "head"'))
elif what_i_do == '15':
    print('----')
    print('User name : ', end = '')
    user_name = input()

    curs.execute(db_change("update user_set set data = 'owner' where id = ? and name = 'acl'"), [user_name])
elif what_i_do == '16':
    print('----')
    print('User name : ', end = '')
    user_name = input()

    curs.execute(db_change('select data from user_set where name = "2fa" and id = ?'), [user_name])
    if curs.fetchall():
        curs.execute(db_change("update user_set set data = '' where name = '2fa' and id = ?"), [user_name])
else:
    print('----')
    print('Markup name : ', end = '')
    markup = input()

    curs.execute(db_change("update other set data = ? where name = 'markup'"), [markup])

if data_db_load == 'Y':
    conn.commit()

print('----')
print('OK')
