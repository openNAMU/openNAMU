from route.tool.func import *
from route.tool.mark import load_conn2, namumark

all_src = []
for i_data in os.listdir("."):
    f_src = re.search("(.+)\.db$", i_data)
    if f_src:
        all_src += [f_src.groups()[0]]

if len(all_src) == 0:
    exit()
elif len(all_src) > 1:
    db_num = 1

    for i_data in all_src:
        print(str(db_num) + ' : ' + i_data)
        
        db_num += 1

    print('----')
    print('Number : ', end = '')    
    db_name = all_src[int(number_check(input())) - 1]
else:
    db_name = all_src[0]

if len(all_src) == 1:
    print('----')
    print('DB\'s name : ' + db_name)

conn = sqlite3.connect(db_name + '.db', check_same_thread = False)


load_conn(conn)

print('----')
print('1. Backlink reset')
print('2. reCAPTCHA delete')
print('3. Ban delete')
print('4. Change host')
print('5. Change port')
print('6. Change skin')
print('7. Change password')
print('8. Reset version')
print('9. New DB create')

print('----')
print('Select : ', end = '')
what_i_do = input()

if what_i_do == '1':
    def parser(data):
        namumark(data[0], data[1], 1)

    sqlQuery("delete from back")
    sqlQuery("commit")

    sqlQuery("select title, data from data")
    data = sqlQuery("fetchall")
    num = 0

    for test in data:
        num += 1

        t = threading.Thread(target = parser, args = [test])
        t.start()
        t.join()

        if num % 10 == 0:
            print(num)
elif what_i_do == '2':
    sqlQuery("delete from other where name = 'recaptcha'")
    sqlQuery("delete from other where name = 'sec_re'")
elif what_i_do == '3':
    print('----')
    print('IP or Name : ', end = '')
    user_data = input()

    if re.search("^([0-9]{1,3}\.[0-9]{1,3})$", user_data):
        band = 'O'
    else:
        band = ''

        sqlQuery("insert into rb (block, end, today, blocker, why, band) values (?, ?, ?, ?, ?, ?)", 
            [user_data, 
            'release', 
            get_time(), 
            'tool:emergency', 
            '', 
            band
        ])
    sqlQuery("delete from ban where block = ?", [user_data])
elif what_i_do == '4':
    print('----')
    print('Host : ', end = '')
    host = input()

    sqlQuery("update other set data = ? where name = 'host'", [host])
elif what_i_do == '5':
    print('----')
    print('Port : ', end = '')
    port = int(input())

    sqlQuery("update other set data = ? where name = 'port'", [port])
elif what_i_do == '6':
    print('----')
    print('Skin\'s name : ', end = '')
    skin = input()

    sqlQuery("update other set data = ? where name = 'skin'", [skin])
elif what_i_do == '7':
    print('----')
    print('1. sha256')
    print('2. sha3')
    
    print('----')
    print('Select : ', end = '')
    what_i_do = int(input())

    print('----')
    print('User\'s name : ', end = '')
    user_name = input()

    print('----')
    print('User\'s password : ', end = '')
    user_pw = input()

    if what_i_do == '1':
        hashed = hashlib.sha256(bytes(user_pw, 'utf-8')).hexdigest()
    else:
        if sys.version_info < (3, 6):
            hashed = sha3.sha3_256(bytes(user_pw, 'utf-8')).hexdigest()
        else:
            hashed = hashlib.sha3_256(bytes(user_pw, 'utf-8')).hexdigest()
       
    sqlQuery("update user set pw = ? where id = ?", [hashed, user_name])
elif what_i_do == '8':
    sqlQuery("update other set data = '00000' where name = 'ver'")
else:
    print('----')
    print('DB\'s name (data) : ', end = '')
    
    db_name = input()
    if db_name == '':
        db_name = 'data'

    sqlite3.connect(db_name + '.db', check_same_thread = False)

sqlQuery("commit")

print('----')
print('OK')