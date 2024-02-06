# Load
import time
import os
import platform
import urllib
import zipfile
import urllib.request

from route.tool.func import *

while True:
    data_db_load = input('Load DB (Y) [Y, N] : ')
    data_db_load = data_db_load.upper()
    if data_db_load in ('Y', 'N'):
        break

if data_db_load == 'Y':
    data_db_set = class_check_json()

    db_data_get(data_db_set['type'])
    do_db_set(data_db_set)

    load_db = get_db_connect()

    conn = load_db.__enter__()
    curs = conn.cursor()
else:
    print('You can use [9, 11, 19]')

# Main
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
print('14. Delete Main <HEAD>')
print('15. Give owner')
print('16. Delete 2FA password')
print('17. Change markup')
print('18. Change wiki access password')
print('19. Forced update')
print('20. Change domain')
print('21. Change TLS')
print('22. Delete body top')
print('23. Delete body bottom')
print('24. SQLite to MySQL')

what_i_do = input('Select : ')
if what_i_do == '1':
    go_num = input('All delete (Y) [Y, N] : ')
    if not go_num == 'N':
        curs.execute(db_change("delete from back"))
        conn.commit()

    try:
        go_num = int(input('Count (100) : '))
    except ValueError:
        go_num = 100

    num = 0

    print('Load...')

    curs.execute(
        db_change(
            "select title from data d "
            "where not exists ("
            "select title from back where link = d.title limit 1"
            ")"
            ""
        )
    )
    title = curs.fetchall()

    print('Rest : ' + str(len(title)))
    print('Start : ' + title[0][0])
    time.sleep(1)

    for name in title:
        num += 1
        if num % go_num == 0:
            print(str(num) + ' : ' + name[0])

        if num % 100 == 0:
            conn.commit()

        curs.execute(db_change("select data from data where title = ?"), [name[0]])
        data = curs.fetchall()

        class_do_render(conn).do_render(name[0], data[0][0], 'backlink', '')
elif what_i_do == '2':
    curs.execute(db_change("delete from other where name = 'recaptcha'"))
    curs.execute(db_change("delete from other where name = 'sec_re'"))
elif what_i_do == '3':
    user_data = input('IP or Name : ')

    curs.execute(
        db_change(
            "insert into rb (block, end, today, blocker, why, band, ongoing, login) "
            "values (?, ?, ?, ?, ?, ?, '', '')"
        ),
        [
            user_data,
            'release',
            get_time(),
            'tool:emergency',
            '',
            '',
        ]
    )

    curs.execute(db_change("update rb set ongoing = '' where block = ?"), [user_data])
elif what_i_do == '4':
    host = input('Host : ')

    curs.execute(db_change("update other set data = ? where name = 'host'"), [host])
elif what_i_do == '5':
    port = int(input('Port : '))

    curs.execute(db_change("update other set data = ? where name = 'port'"), [port])
elif what_i_do == '6':
    skin = input('Skin name : ')

    curs.execute(db_change("update other set data = ? where name = 'skin'"), [skin])
elif what_i_do == '7':
    user_name = input('User name : ')
    user_pw = input('User password : ')

    hashed = pw_encode(user_pw)
    curs.execute(db_change("update user_set set data = ? where id = ? and name = 'pw'"), [
        hashed,
        user_name
    ])
elif what_i_do == '8':
    new_ver = input('Insert version (0000000) : ')

    if new_ver == '':
        new_ver = '0000000'

    curs.execute(db_change("update other set data = ? where name = 'ver'"), [new_ver])
elif what_i_do == '9':
    if os.path.exists(os.path.join('data', 'set.json')):
        os.remove(os.path.join('data', 'set.json'))
elif what_i_do == '10':
    user_name = input('User name : ')
    new_name = input('New name : ')

    curs.execute(
        db_change("update user_set set id = ? where id = ?"),
        [new_name, user_name]
    )
elif what_i_do == '11':
    if os.path.exists(os.path.join('data', 'mysql.json')):
        os.remove(os.path.join('data', 'mysql.json'))
elif what_i_do == '14':
    curs.execute(db_change('delete from other where name = "head"'))
elif what_i_do == '15':
    user_name = input('User name : ')

    curs.execute(db_change("update user_set set data = 'owner' where id = ? and name = 'acl'"), [user_name])
elif what_i_do == '16':
    user_name = input('User name : ')

    curs.execute(db_change('select data from user_set where name = "2fa" and id = ?'), [user_name])
    if curs.fetchall():
        curs.execute(db_change("update user_set set data = '' where name = '2fa' and id = ?"), [user_name])
elif what_i_do == '17':
    markup = input('Markup name : ')

    curs.execute(db_change("update other set data = ? where name = 'markup'"), [markup])
elif what_i_do == '18':
    wiki_access_password = input('Password : ')

    curs.execute(db_change("update other set data = ? where name = 'wiki_access_password'"), [wiki_access_password])
elif what_i_do == '19':
    up_data = input('Insert branch (beta) [stable, beta, dev] : ')

    if not up_data in ['stable', 'beta', 'dev']:
        up_data = 'beta'

    if platform.system() == 'Linux':
        ok = []

        ok += [os.system('git remote rm origin')]
        ok += [os.system('git remote add origin https://github.com/opennamu/opennamu.git')]
        ok += [os.system('git fetch origin ' + up_data)]
        ok += [os.system('git reset --hard origin/' + up_data)]
        if (ok[0] and ok[1] and ok[2] and ok[3]) != 0:
            print('Error : update failed')
    elif platform.system() == 'Windows':
        os.system('rd /s /q route')
        urllib.request.urlretrieve('https://github.com/opennamu/opennamu/archive/' + up_data + '.zip', 'update.zip')
        zipfile.ZipFile('update.zip').extractall('')
        ok = os.system('xcopy /y /s /r opennamu-' + up_data + ' .')
        if ok == 0:
            os.system('rd /s /q opennamu-' + up_data)
            os.system('del update.zip')
        else:
            print('Error : update failed')
elif what_i_do == '20':
    domain = input('Domain (EX : 2du.pythonanywhere.com) : ')

    curs.execute(db_change('delete from other where name = "domain"'))
    curs.execute(db_change('insert into other (name, data, coverage) values ("domain", ?, "")'), [domain])
elif what_i_do == '21':
    tls_v = input('TLS (http) [http, https] : ')
    if not tls_v in ['http', 'https']:
        tls_v = 'http'

    curs.execute(db_change('delete from other where name = "http_select"'))
    curs.execute(db_change('insert into other (name, data, coverage) values ("http_select", ?, "")'), [tls_v])
elif what_i_do == '22':
    curs.execute(db_change('delete from other where name = "body"'))
elif what_i_do == '23':
    curs.execute(db_change('delete from other where name = "bottom_body"'))
elif what_i_do == '24':
    load_db = get_db_connect('mysql')
    mysql_conn = load_db.__enter__()
    mysql_curs = mysql_conn.cursor()
    
    load_db = get_db_connect('sqlite')
    sqlite_conn = load_db.__enter__()
    sqlite_curs = sqlite_conn.cursor()

    create_data = get_db_table_list()
    for create_table in create_data:
        create = ['test'] + create_data[create_table]

        create_r = ', '.join(['%s' for _ in create])
        create = ', '.join(create)
        
        mysql_curs.execute(db_change('delete from ' + create_table))
        
        sqlite_curs.execute(db_change('select ' + create + ' from ' + create_table))
        db_data = sqlite_curs.fetchall()
        if db_data:
            mysql_curs.executemany("insert into " + create_table + " (" + create + ") values (" + create_r + ")", db_data)
else:
    raise ValueError(what_i_do)

if data_db_load == 'Y':
    try:
        conn.commit()
    except:
        pass

print('OK')