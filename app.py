from flask import Flask, request, send_from_directory
from flask_reggie import Reggie

from tornado.wsgi import WSGIContainer
from tornado.httpserver import HTTPServer
from tornado.ioloop import IOLoop

import bcrypt
import difflib
import shutil
import threading
import logging
import random
logging.basicConfig(level = logging.ERROR)

app = Flask(__name__)
Reggie(app)

r_ver = 'v2.6.1'

from func import *
from set_mark.mid_pas import mid_pas
from set_mark.macro import savemark

# set.json 설정 확인
try:
    json_data = open('set.json').read()
    set_data = json.loads(json_data)
except:
    while(1):
        print('DB 이름 : ', end = '')
        new_json = str(input())

        if(new_json != ''):
            with open("set.json", "w") as f:
                f.write('{ "db" : "' + new_json + '" }')
            
            json_data = open('set.json').read()
            set_data = json.loads(json_data)

            break
        else:
            print('값을 입력하세요.')
            
            pass

# 디비 연결
conn = sqlite3.connect(set_data['db'] + '.db')
curs = conn.cursor()

# 스킨 불러오기 부분
skin_check(conn)

# 셋업 부분
curs.execute("create table if not exists data(title text, data text)")
curs.execute("create table if not exists history(id text, title text, data text, date text, ip text, send text, leng text, hide text)")
curs.execute("create table if not exists rd(title text, sub text, date text)")
curs.execute("create table if not exists user(id text, pw text, acl text, date text)")
curs.execute("create table if not exists ban(block text, end text, why text, band text)")
curs.execute("create table if not exists topic(id text, title text, sub text, data text, date text, ip text, block text, top text)")
curs.execute("create table if not exists stop(title text, sub text, close text)")
curs.execute("create table if not exists rb(block text, end text, today text, blocker text, why text)")
curs.execute("create table if not exists back(title text, link text, type text)")
curs.execute("create table if not exists agreedis(title text, sub text)")
curs.execute("create table if not exists custom(user text, css text)")
curs.execute("create table if not exists other(name text, data text)")
curs.execute("create table if not exists alist(name text, acl text)")
curs.execute("create table if not exists re_admin(who text, what text, time text)")
curs.execute("create table if not exists alarm(name text, data text, date text)")
curs.execute("create table if not exists ua_d(name text, ip text, ua text, today text, sub text)")
curs.execute("create table if not exists filter(name text, regex text, sub text)")
curs.execute("create table if not exists scan(user text, title text)")
curs.execute("create table if not exists acl(title text, dec text, dis text, why text)")

curs.execute("select name from alist where name = '소유자'")
if(not curs.fetchall()):
    curs.execute("insert into alist (name, acl) values ('소유자', 'owner')")

curs.execute("select data from other where name = 'port'")
rep_data = curs.fetchall()
if(not rep_data):
    while(1):
        print('포트 : ', end = '')
        rep_port = int(input())

        if(rep_port):
            curs.execute("insert into other (name, data) values ('port', ?)", [rep_port])
            
            break
        else:
            pass
else:
    rep_port = rep_data[0][0]

curs.execute("select data from other where name = 'key'")
rep_data = curs.fetchall()
if(not rep_data):
    while(1):
        print('비밀 키 : ', end = '')
        rep_key = str(input())

        if(rep_key):
            curs.execute("insert into other (name, data) values ('key', ?)", [rep_key])
            
            break
        else:
            pass
else:
    rep_key = rep_data[0][0]

# 호환성 설정
curs.execute("drop table if exists move")

try:
    curs.execute("alter table history add hide text default ''")
    
    curs.execute('select title, re from hidhi')
    for rep in curs.fetchall():
        curs.execute("update history set hide = 'O' where title = ? and id = ?", [rep[0], rep[1]])
    curs.execute("drop table if exists hidhi")
except:
    pass

try:
    curs.execute("alter table user add date text default ''")
except:
    pass

try:
    curs.execute("alter table rb add band text default ''")
except:
    pass

try:
    curs.execute("alter table ban add login text default ''")
except:
    pass

try:
    curs.execute("select title, acl from data where acl != ''")
    for rep in curs.fetchall():
        curs.execute("insert into acl (title, dec, dis, why) values (?, ?, '', '')", [rep[0], rep[1]])
    curs.execute("alter table data drop acl")
except:
    pass

try:
    curs.execute("select data from other where name = 'robot'")
    robot_test = curs.fetchall()
    if(robot_test):
        fw_test = open('./robots.txt', 'w')
        fw_test.write(robot_test[0][0])
        fw_test.close()
except:
    pass

conn.commit()

# 이미지 폴더 생성
if(not os.path.exists('image')):
    os.makedirs('image')
    
# 스킨 폴더 생성
if(not os.path.exists('views')):
    os.makedirs('views')

# 백업 설정
def back_up():
    try:
        shutil.copyfile(set_data['db'] + '.db', 'back_' + set_data['db'] + '.db')
        print('백업 성공')
    except:
        print('백업 오류')

    threading.Timer(60 * 60 * back_time, back_up).start()

try:
    curs.execute('select data from other where name = "back_up"')
    back_up_time = curs.fetchall()
    back_time = int(back_up_time[0][0])
except:
    back_time = 0
    
# 백업 여부 확인
if(back_time != 0):
    print(str(back_time) + '시간 간격으로 백업')

    if(__name__ == '__main__'):
        back_up()
else:
    print('백업하지 않음')

@app.route('/del_alarm')
def del_alarm():
    curs.execute("delete from alarm where name = ?", [ip_check()])
    conn.commit()

    return(redirect('/alarm'))

@app.route('/alarm')
def alarm():
    ip = ip_check()
    if(re.search('(?:\.|:)', ip)):
        return(redirect('/login'))    

    da = '<ul>'    
    curs.execute("select data, date from alarm where name = ? order by date desc", [ip])
    dt = curs.fetchall()
    if(dt):
        da = '<a href="/del_alarm">(모두 삭제)</a><hr>' + da

        for do in dt:
            da += '<li>' + do[0] + ' / ' + do[1] + '</li>'
    else:
        da += '<li>알림이 없습니다.</li>'
    da += '</ul>'

    return(html_minify(template('index', 
        imp = ['알림', wiki_set(conn, 1), custom(conn), other2([0, 0])],
        data = da,
        menu = [['user', '사용자']]
    )))

@app.route('/edit_set')
@app.route('/edit_set/<int:num>', methods=['POST', 'GET'])
def edit_set(num = 0):
    if(num != 0 and admin_check(conn, None, None) != 1):
        return(re_error(conn, '/ban'))

    if(num == 0):
        li_list = ['기본 설정', '문구 관련', '전역 HEAD', 'robots.txt', '구글 관련']
        x = 0
        li_data = ''
        for li in li_list:
            x += 1
            li_data += ' * [[wiki:edit_set/' + str(x) + '|' + li + ']]\r\n'

        return(html_minify(template('index', 
            imp = ['설정 편집', wiki_set(conn, 1), custom(conn), other2([0, 0])],
            data = namumark(conn, '',   '[목차(없음)]\r\n' + \
                                        '== 메뉴 ==\r\n' + \
                                        li_data, 0, 0, 0),
            menu = [['manager', '관리자']]
        )))
    elif(num == 1):
        i_list = ['name', 'logo', 'frontpage', 'license', 'upload', 'skin', 'edit', 'reg', 'ip_view', 'back_up', 'all_title']
        n_list = ['무명위키', '', '위키:대문', 'CC 0', '2', '', 'normal', '', '', '0', '']
        if(request.method == 'POST'):
            i = 0

            for data in i_list:
                curs.execute("update other set data = ? where name = ?", [request.form.get(data, n_list[i]), data])
                i += 1

            conn.commit()
            admin_check(conn, None, 'edit_set')

            return(redirect('/edit_set/1'))
        else:
            d_list = []
            
            x = 0
            for i in i_list:
                curs.execute('select data from other where name = ?', [i])
                sql_d = curs.fetchall()
                if(sql_d):
                    d_list += [sql_d[0][0]]
                else:
                    curs.execute('insert into other (name, data) values (?, ?)', [i, n_list[x]])
                    d_list += [n_list[x]]

                x += 1
            conn.commit()

            div = ''
            if(d_list[6] == 'login'):
                div += '<option value="login">가입자</option>'
                div += '<option value="normal">일반</option>'
                div += '<option value="admin">관리자</option>'
            elif(d_list[6] == 'admin'):
                div += '<option value="admin">관리자</option>'
                div += '<option value="login">가입자</option>'
                div += '<option value="normal">일반</option>'
            else:
                div += '<option value="normal">일반</option>'
                div += '<option value="admin">관리자</option>'
                div += '<option value="login">가입자</option>'

            ch_1 = ''
            ch_2 = ''
            ch_3 = ''
            if(d_list[7]):
                ch_1 = 'checked="checked"'
            
            if(d_list[8]):
                ch_2 = 'checked="checked"'

            if(d_list[10]):
                ch_3 = 'checked="checked"'                

            return(html_minify(template('index', 
                imp = ['기본 설정', wiki_set(conn, 1), custom(conn), other2([0, 0])],
                data = '<form method="post"> \
                            <span>이름</span><br><br> \
                            <input placeholder="이름" type="text" name="name" value="' + html.escape(d_list[0]) + '"><hr> \
                            <span>로고 (HTML)</span><br><br> \
                            <input placeholder="로고" type="text" name="logo" value="' + html.escape(d_list[1]) + '"><hr> \
                            <span>대문</span><br><br> \
                            <input placeholder="대문" type="text" name="frontpage" value="' + html.escape(d_list[2]) + '"><hr> \
                            <span>라이선스 (HTML)</span><br><br> \
                            <input placeholder="라이선스" type="text" name="license" value="' + html.escape(d_list[3]) + '"><hr> \
                            <span>파일 크기 [메가]</span><br><br> \
                            <input placeholder="파일 크기" type="text" name="upload" value="' + html.escape(d_list[4]) + '"><hr> \
                            <span>스킨 {재시작 필요}</span><br><br> \
                            <input placeholder="스킨" type="text" name="skin" value="' + html.escape(d_list[5]) + '"><hr> \
                            <span>전역 ACL</span><br><br> \
                            <select name="edit">' + div + '</select><hr> \
                            <input type="checkbox" name="reg" ' + ch_1 + '> 가입불가<hr> \
                            <input type="checkbox" name="ip_view" ' + ch_2 + '> 아이피 비공개<hr> \
                            <input type="checkbox" name="all_title" ' + ch_3 + '> 모든 문서 보기 비활성화<hr> \
                            <span>백업 간격 [시간] (끄기 : 0) {재시작 필요}</span><hr> \
                            <input placeholder="백업 간격" type="text" name="back_up" value="' + html.escape(d_list[9]) + '"><hr> \
                            <button class="btn btn-primary" type="submit">저장</button> \
                        </form>',
                menu = [['edit_set', '설정']]
            )))
    elif(num == 2):
        if(request.method == 'POST'):
            curs.execute("update other set data = ? where name = ?", [request.form['contract'], 'contract'])
            curs.execute("update other set data = ? where name = ?", [request.form['no_login_warring'], 'no_login_warring'])
            conn.commit()

            admin_check(conn, None, 'edit_set')
            return(redirect('/edit_set/2'))
        else:
            i_list = ['contract', 'no_login_warring']
            n_list = ['', '']
            d_list = []
            
            x = 0
            for i in i_list:
                curs.execute('select data from other where name = ?', [i])
                sql_d = curs.fetchall()
                if(sql_d):
                    d_list += [sql_d[0][0]]
                else:
                    curs.execute('insert into other (name, data) values (?, ?)', [i, n_list[x]])
                    d_list += [n_list[x]]

                x += 1
            conn.commit()

            return(html_minify(template('index', 
                imp = ['문구 관련', wiki_set(conn, 1), custom(conn), other2([0, 0])],
                data = '<form method="post"> \
                            <span>가입 약관</span><br><br> \
                            <input placeholder="가입 약관" type="text" name="contract" value="' + html.escape(d_list[0]) + '"><hr> \
                            <span>비 로그인 경고</span><br><br> \
                            <input placeholder="비 로그인 경고" type="text" name="no_login_warring" value="' + html.escape(d_list[1]) + '"><hr> \
                            <button class="btn btn-primary" type="submit">저장</button> \
                        </form>',
                menu = [['edit_set', '설정']]
            )))
    elif(num == 3):
        if(request.method == 'POST'):
            curs.execute("select name from other where name = 'head'")
            if(curs.fetchall()):
                curs.execute("update other set data = ? where name = 'head'", [request.form['content']])
            else:
                curs.execute("insert into other (name, data) values ('head', ?)", [request.form['content']])
            conn.commit()

            admin_check(conn, None, 'edit_set')
            return(redirect('/edit_set/3'))
        else:
            curs.execute("select data from other where name = 'head'")
            head = curs.fetchall()
            if(head):
                data = head[0][0]
            else:
                data = ''

            return(html_minify(template('index', 
                imp = ['전역 HEAD', wiki_set(conn, 1), custom(conn), other2([0, 0])],
                data =  '<span>&lt;style&gt;CSS&lt;/style&gt;<br>&lt;script&gt;JS&lt;/script&gt;</span><hr> \
                        <form method="post"> \
                            <textarea rows="25" name="content">' + html.escape(data) + '</textarea><hr> \
                            <button class="btn btn-primary" type="submit">저장</button> \
                        </form>',
                menu = [['edit_set', '설정']]
            )))
    elif(num == 4):
        if(request.method == 'POST'):
            curs.execute("select name from other where name = 'robot'")
            if(curs.fetchall()):
                curs.execute("update other set data = ? where name = 'robot'", [request.form['content']])
            else:
                curs.execute("insert into other (name, data) values ('robot', ?)", [request.form['content']])
            conn.commit()

            fw = open('./robots.txt', 'w')
            fw.write(request.form['content'])
            fw.close()
            
            admin_check(conn, None, 'edit_set')
            return(redirect('/edit_set/4'))
        else:
            curs.execute("select data from other where name = 'robot'")
            robot = curs.fetchall()
            if(robot):
                data = robot[0][0]
            else:
                data = ''

            return(html_minify(template('index', 
                imp = ['robots.txt', wiki_set(conn, 1), custom(conn), other2([0, 0])],
                data =  '<a href="/robots.txt">(보기)</a><hr> \
                        <form method="post"> \
                            <textarea rows="25" name="content">' + html.escape(data) + '</textarea><hr> \
                            <button class="btn btn-primary" type="submit">저장</button> \
                        </form>',
                menu = [['edit_set', '설정']]
            )))
    elif(num == 5):
        if(request.method == 'POST'):
            curs.execute("update other set data = ? where name = 'recaptcha'", [request.form['recaptcha']])
            curs.execute("update other set data = ? where name = 'sec_re'", [request.form['sec_re']])
            conn.commit()

            admin_check(conn, None, 'edit_set')
            return(redirect('/edit_set/5'))
        else:
            i_list = ['recaptcha', 'sec_re']
            n_list = ['', '']
            d_list = []
            
            x = 0
            for i in i_list:
                curs.execute('select data from other where name = ?', [i])
                sql_d = curs.fetchall()
                if(sql_d):
                    d_list += [sql_d[0][0]]
                else:
                    curs.execute('insert into other (name, data) values (?, ?)', [i, n_list[x]])
                    d_list += [n_list[x]]

                x += 1
            conn.commit()

            return(html_minify(template('index', 
                imp = ['구글 관련', wiki_set(conn, 1), custom(conn), other2([0, 0])],
                data = '<form method="post"> \
                            <span>리캡차 (HTML)</span><br><br> \
                            <input placeholder="리캡차 (HTML)" type="text" name="recaptcha" value="' + html.escape(d_list[0]) + '"><hr> \
                            <span>리캡차 (비밀키)</span><br><br> \
                            <input placeholder="리캡차 (비밀키)" type="text" name="sec_re" value="' + html.escape(d_list[1]) + '"><hr> \
                            <button class="btn btn-primary" type="submit">저장</button> \
                        </form>',
                menu = [['edit_set', '설정']]
            )))
    else:
        return(redirect('/'))

@app.route('/not_close_topic')
def not_close_topic():
    div = '<ul>'

    curs.execute('select title, sub from rd order by date desc')
    n_list = curs.fetchall()
    for data in n_list:
        curs.execute('select * from stop where title = ? and sub = ? and close = "O"', [data[0], data[1]])
        is_close = curs.fetchall()
        if(not is_close):
            div += '<li><a href="/topic/' + url_pas(data[0]) + '/sub/' + url_pas(data[1]) + '">' + data[0] + ' (' + data[1] + ')</a></li>'
            
    div += '</ul>'

    return(html_minify(template('index', 
        imp = ['열린 토론 목록', wiki_set(conn, 1), custom(conn), other2([0, 0])],
        data = div,
        menu = [['manager', '관리자']]
    )))

@app.route('/image/<name>')
def image_view(name = None):
    if(os.path.exists(os.path.join('image', name))):
        return(send_from_directory('./image', name))
    else:
        return(redirect('/'))

@app.route('/acl_list')
def acl_list():
    div = '<ul>'

    curs.execute("select title, dec from acl where dec = 'admin' or dec = 'user' order by title desc")
    list_data = curs.fetchall()
    for data in list_data:
        if(not re.search('^사용자:', data[0]) and not re.search('^파일:', data[0])):
            if(data[1] == 'admin'):
                acl = '관리자'
            else:
                acl = '가입자'

            div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a> (' + acl + ')</li>'
        
    div += '</ul>'
    
    return(html_minify(template('index', 
        imp = ['ACL 문서 목록', wiki_set(conn, 1), custom(conn), other2([0, 0])],
        data = div,
        menu = [['other', '기타']]
    )))

@app.route('/admin_plus/<name>', methods=['POST', 'GET'])
def admin_plus(name = None):
    if(request.method == 'POST'):
        if(admin_check(conn, None, 'admin_plus (' + name + ')') != 1):
            return(re_error(conn, '/error/3'))

        curs.execute("delete from alist where name = ?", [name])
        
        if(request.form['ban']):
            curs.execute("insert into alist (name, acl) values (?, 'ban')", [name])

        if(request.form['mdel']):
            curs.execute("insert into alist (name, acl) values (?, 'mdel')", [name])   

        if(request.form['toron']):
            curs.execute("insert into alist (name, acl) values (?, 'toron')", [name])
            
        if(request.form['check']):
            curs.execute("insert into alist (name, acl) values (?, 'check')", [name])

        if(request.form['acl']):
            curs.execute("insert into alist (name, acl) values (?, 'acl')", [name])

        if(request.form['hidel']):
            curs.execute("insert into alist (name, acl) values (?, 'hidel')", [name])

        if(request.form['give']):
            curs.execute("insert into alist (name, acl) values (?, 'give')", [name])

        if(request.form['owner']):
            curs.execute("insert into alist (name, acl) values (?, 'owner')", [name])
            
        conn.commit()
        
        return(redirect('/admin_plus/' + url_pas(name)))
    else:
        curs.execute('select acl from alist where name = ?', [name])
        acl_list = curs.fetchall()
        
        data = '<ul>'
        exist_list = ['', '', '', '', '', '', '', '']

        for go in acl_list:
            if(go[0] == 'ban'):
                exist_list[0] = 'checked="checked"'
            elif(go[0] == 'mdel'):
                exist_list[1] = 'checked="checked"'
            elif(go[0] == 'toron'):
                exist_list[2] = 'checked="checked"'
            elif(go[0] == 'check'):
                exist_list[3] = 'checked="checked"'
            elif(go[0] == 'acl'):
                exist_list[4] = 'checked="checked"'
            elif(go[0] == 'hidel'):
                exist_list[5] = 'checked="checked"'
            elif(go[0] == 'give'):
                exist_list[6] = 'checked="checked"'
            elif(go[0] == 'owner'):
                exist_list[7] = 'checked="checked"'

        if(admin_check(conn, None, None) != 1):
            state = 'disabled'
        else:
            state = ''

        data += '<li><input type="checkbox" ' + state +  ' name="ban" ' + exist_list[0] + '> 차단</li>'
        data += '<li><input type="checkbox" ' + state +  ' name="mdel" ' + exist_list[1] + '> 많은 문서 삭제</li>'
        data += '<li><input type="checkbox" ' + state +  ' name="toron" ' + exist_list[2] + '> 토론 관리</li>'
        data += '<li><input type="checkbox" ' + state +  ' name="check" ' + exist_list[3] + '> 사용자 검사</li>'
        data += '<li><input type="checkbox" ' + state +  ' name="acl" ' + exist_list[4] + '> 문서 ACL</li>'
        data += '<li><input type="checkbox" ' + state +  ' name="hidel" ' + exist_list[5] + '> 역사 숨김</li>'
        data += '<li><input type="checkbox" ' + state +  ' name="give" ' + exist_list[6] + '> 권한 부여</li>'
        data += '<li><input type="checkbox" ' + state +  ' name="owner" ' + exist_list[7] + '> 소유자</li></ul>'

        return(html_minify(template('index', 
            imp = ['관리 그룹 추가', wiki_set(conn, 1), custom(conn), other2([0, 0])],
            data = '<form method="post">' + data + '<button ' + state +  ' class="btn btn-primary" type="submit">저장</button></form>',
            menu = [['manager', '관리자']]
        )))        
        
@app.route('/admin_list')
def admin_list():
    div = '<ul>'
    
    curs.execute("select id, acl, date from user where not acl = 'user' order by date desc")
    user_data = curs.fetchall()

    for data in user_data:
        name = ip_pas(conn, data[0]) + ' <a href="/admin_plus/' + url_pas(data[1]) + '">(' + data[1] + ')</a>'

        if(data[2] != ''):
            name += '(가입 : ' + data[2] + ')'

        div += '<li>' + name + '</li>'
        
    div += '</ul>'
                
    return(html_minify(template('index', 
        imp = ['관리자 목록', wiki_set(conn, 1), custom(conn), other2([0, 0])],
        data = div,
        menu = [['other', '기타']]
    )))
        
@app.route('/hidden/<path:name>')
def history_hidden(name = None):
    num = int(request.args.get('num', 0))

    if(admin_check(conn, 6, 'history_hidden (' + name + '#' + str(num) + ')') == 1):
        curs.execute("select title from history where title = ? and id = ? and hide = 'O'", [name, str(num)])
        if(curs.fetchall()):
            curs.execute("update history set hide = '' where title = ? and id = ?", [name, str(num)])
        else:
            curs.execute("update history set hide = 'O' where title = ? and id = ?", [name, str(num)])
            
        conn.commit()
    
    return(redirect('/history/' + url_pas(name)))
        
@app.route('/user_log')
def user_log():
    num = int(request.args.get('num', 1))
    if(num * 50 > 0):
        sql_num = num * 50 - 50
    else:
        sql_num = 0
        
    list_data = '<ul>'
    admin_one = admin_check(conn, 1, None)
    
    curs.execute("select id, date from user order by date desc limit ?, '50'", [str(sql_num)])
    user_list = curs.fetchall()
    for data in user_list:
        if(admin_one == 1):
            curs.execute("select block from ban where block = ?", [data[0]])
            ban_exist = curs.fetchall()
            if(ban_exist):
                ban_button = ' <a href="/ban/' + url_pas(data[0]) + '">(해제)</a>'
            else:
                ban_button = ' <a href="/ban/' + url_pas(data[0]) + '">(차단)</a>'
        else:
            ban_button = ''
            
        ip = ip_pas(conn, data[0])
        list_data += '<li>' + ip + ban_button

        if(data[1] != ''):
            list_data += ' (가입 : ' + data[1] + ')'

        list_data += '</li>'

    if(num == 1):
        curs.execute("select count(id) from user")
        user_count = curs.fetchall()
        if(user_count):
            count = user_count[0][0]
        else:
            count = 0

        list_data += '</ul><hr><ul><li>이 위키에는 ' + str(count) + '명의 사람이 있습니다.</li></ul>'

    list_data += next_fix('/user_log?num=', num, user_list)

    return(html_minify(template('index', 
        imp = ['사용자 가입 기록', wiki_set(conn, 1), custom(conn), other2([0, 0])],
        data = list_data,
        menu = [['other', '기타']]
    )))

@app.route('/admin_log')
def admin_log():
    num = int(request.args.get('num', 1))
    if(num * 50 > 0):
        sql_num = num * 50 - 50
    else:
        sql_num = 0

    list_data = '<ul>'
    
    curs.execute("select who, what, time from re_admin order by time desc limit ?, '50'", [str(sql_num)])
    get_list = curs.fetchall()
    for data in get_list:            
        ip = ip_pas(conn, data[0])
        list_data += '<li>' + ip + ' / ' + data[1] + ' / ' + data[2] + '</li>'

    list_data += '</ul><hr><ul><li>주의 : 권한 사용 안하고 열람만 해도 기록되는 경우도 있습니다.</li></ul>'
    list_data += next_fix('/admin_log?num=', num, get_list)

    return(html_minify(template('index', 
        imp = ['권한 사용 기록', wiki_set(conn, 1), custom(conn), other2([0, 0])],
        data = list_data,
        menu = [['other', '기타']]
    )))

@app.route('/give_log')
def give_log():        
    list_data = '<ul>'
    back = ''

    curs.execute("select distinct name from alist order by name asc")
    get_list = curs.fetchall()
    for data in get_list:                      
        if(back != data[0]):
            back = data[0]

        list_data += '<li><a href="/admin_plus/' + url_pas(data[0]) + '">' + data[0] + '</a></li>'
    
    list_data += '</ul><hr><a href="/manager/8">(생성)</a>'

    return(html_minify(template('index', 
        imp = ['권한 목록', wiki_set(conn, 1), custom(conn), other2([0, 0])],
        data = list_data,
        menu = [['other', '기타']]
    )))

@app.route('/indexing')
def indexing():
    if(admin_check(conn, None, 'indexing') != 1):
        return(re_error(conn, '/error/3'))

    curs.execute("select name from sqlite_master where type in ('table', 'view') and name not like 'sqlite_%' union all select name from sqlite_temp_master where type in ('table', 'view') order by 1;")
    data = curs.fetchall()
    for table in data:
        print('----- ' + table[0] + ' -----')
        curs.execute('select sql from sqlite_master where name = ?', [table[0]])
        cul = curs.fetchall()
        r_cul = re.findall('(?:([^ (]*) text)', str(cul[0]))
        for n_cul in r_cul:
            print(n_cul)
            sql = 'create index index_' + table[0] + '_' + n_cul + ' on ' + table[0] + '(' + n_cul + ')'
            try:
                curs.execute(sql)
            except:
                pass

    conn.commit()
    return(redirect('/'))        
        
@app.route('/xref/<path:name>')
def xref(name = None):
    num = int(request.args.get('num', 1))
    if(num * 50 > 0):
        sql_num = num * 50 - 50
    else:
        sql_num = 0
        
    div = '<ul>'
    
    curs.execute("select link, type from back where title = ? and not type = 'cat' and not type = 'no' order by link asc limit ?, '50'", [name, str(sql_num)])
    data_list = curs.fetchall()
    for data in data_list:
        div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a>'
        
        if(data[1]):
            if(data[1] == 'include'):
                side = '포함'
            elif(data[1] == 'file'):
                side = '파일'
            else:
                side = '넘겨주기'
                
            div += ' (' + side + ')'
        
        div += '</li>'
        
        if(re.search('^틀:', data[0])):
            div += '<li><a id="inside" href="/xref/' + url_pas(data[0]) + '">' + data[0] + '</a> (역링크)</li>'
      
    div += '</ul>' + next_fix('/xref/' + url_pas(name) + '?num=', num, data_list)
    
    return(html_minify(template('index', 
        imp = [name, wiki_set(conn, 1), custom(conn), other2([' (역링크)', 0])],
        data = div,
        menu = [['w/' + url_pas(name), '문서']]
    )))

@app.route('/please')
def please():
    num = int(request.args.get('num', 1))
    if(num * 50 > 0):
        sql_num = num * 50 - 50
    else:
        sql_num = 0
        
    div = '<ul>'
    var = ''
    
    curs.execute("select distinct title from back where type = 'no' order by title asc limit ?, '50'", [str(sql_num)])
    data_list = curs.fetchall()
    for data in data_list:
        if(var != data[0]):
            div += '<li><a class="not_thing" href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a></li>'        
            var = data[0]
        
    div += '</ul>' + next_fix('/please?num=', num, data_list)
    
    return(html_minify(template('index', 
        imp = ['필요한 문서', wiki_set(conn, 1), custom(conn), other2([0, 0])],
        data = div,
        menu = [['other', '기타']]
    )))
        
@app.route('/recent_discuss')
@app.route('/recent_discuss/<regex("close"):tools>')
def recent_discuss(tools = 'normal'):
    if(tools == 'normal' or tools == 'close'):
        div = ''
        
        if(tools == 'normal'):
            div += '<a href="/recent_discuss/close">(닫힌 토론)</a>'
            m_sub = 0
        else:
            div += '<a href="/recent_discuss">(열린 토론)</a>'
            m_sub = ' (닫힘)'

        div +=  '<hr><table style="width: 100%; text-align: center;"><tbody><tr><td style="width: 50%;">토론명</td><td style="width: 50%;">시간</td></tr>'
    else:
        return(redirect('/'))
    
    curs.execute("select title, sub, date from rd order by date desc limit 50")
    for data in curs.fetchall():
        title = html.escape(data[0])
        sub = html.escape(data[1])

        close = 0
        if(tools == 'normal'):
            curs.execute("select title from stop where title = ? and sub = ? and close = 'O'", [data[0], data[1]])
            if(curs.fetchall()):
                close = 1
        else:
            curs.execute("select title from stop where title = ? and sub = ? and close = 'O'", [data[0], data[1]])
            if(not curs.fetchall()):
                close = 1

        if(close == 0):
            div += '<tr><td><a href="/topic/' + url_pas(data[0]) + '/sub/' + url_pas(data[1]) + '">' + title + '</a> (' + sub + ')</td><td>' + data[2] + '</td></tr>'
    else:
        div += '</tbody></table>'
            
    return(html_minify(template('index', 
        imp = ['최근 토론내역', wiki_set(conn, 1), custom(conn), other2([m_sub, 0])],
        data = div,
        menu = 0
    )))

@app.route('/block_log')
@app.route('/block_log/<regex("ip|user|never_end|can_end|end|now"):tool2>')
@app.route('/<regex("block_user|block_admin"):tool>/<name>')
def block_log(name = None, tool = None, tool2 = None):
    num = int(request.args.get('num', 1))
    if(num * 50 > 0):
        sql_num = num * 50 - 50
    else:
        sql_num = 0
    
    div = '<table style="width: 100%; text-align: center;"><tbody><tr><td style="width: 33.3%;">차단자</td><td style="width: 33.3%;">관리자</td><td style="width: 33.3%;">기간</td></tr>'
    data_list = ''
    
    if(not name):
        if(not tool2):
            div = '<a href="/manager/11">(차단자)</a> <a href="/manager/12">(관리자)</a><hr><a href="/block_log/ip">(아이피)</a> <a href="/block_log/user">(가입자)</a> <a href="/block_log/never_end">(영구)</a> <a href="/block_log/can_end">(기간)</a> <a href="/block_log/end">(해제)</a> <a href="/block_log/now">(현재)</a><hr>' + div
            sub = 0
            menu = [['other', '기타']]

            curs.execute("select why, block, blocker, end, today from rb order by today desc limit ?, '50'", [str(sql_num)])
        else:
            menu = [['block_log', '일반']]

            if(tool2 == 'ip'):
                sub = ' (아이피)'

                curs.execute("select why, block, blocker, end, today from rb where (block like ? or block like ?) order by today desc limit ?, '50'", ['%.%', '%:%', str(sql_num)])
            elif(tool2 == 'user'):
                sub = ' (가입자)'

                curs.execute("select why, block, blocker, end, today from rb where not (block like ? or block like ?) order by today desc limit ?, '50'", ['%.%', '%:%', str(sql_num)])
            elif(tool2 == 'never_end'):
                sub = '(영구)'

                curs.execute("select why, block, blocker, end, today from rb where not end like ? and not end like ? order by today desc limit ?, '50'", ['%:%', '%해제%', str(sql_num)])
            elif(tool2 == 'end'):
                sub = '(해제)'

                curs.execute("select why, block, blocker, end, today from rb where end = ? order by today desc limit ?, '50'", ['해제', str(sql_num)])
            elif(tool2 == 'now'):
                sub = '(현재)'
                data_list = []

                curs.execute("select block from ban limit ?, '50'", [str(sql_num)])
                for in_data in curs.fetchall():
                    curs.execute("select why, block, blocker, end, today from rb where block = ? order by today desc limit 1", [in_data[0]])
                    data_list = [curs.fetchall()[0]] + data_list
            else:
                sub = '(기간)'

                curs.execute("select why, block, blocker, end, today from rb where end like ? order by today desc limit ?, '50'", ['%-%', str(sql_num)])
    else:
        menu = [['block_log', '일반']]

        if(tool == 'block_user'):
            sub = ' (차단자)'

            curs.execute("select why, block, blocker, end, today from rb where block = ? order by today desc limit ?, '50'", [name, str(sql_num)])
        else:
            sub = ' (관리자)'

            curs.execute("select why, block, blocker, end, today from rb where blocker = ? order by today desc limit ?, '50'", [name, str(sql_num)])

    if(data_list == ''):
        data_list = curs.fetchall()

    for data in data_list:
        why = html.escape(data[0])

        if(why == ''):
            why = '<br>'
        
        b = re.search("^([0-9]{1,3}\.[0-9]{1,3})$", data[1])
        if(b):
            ip = data[1] + ' (대역)'
        else:
            ip = ip_pas(conn, data[1])

        if(data[3] != ''):
            end = data[3]
        else:
            end = '무기한'
            
        div += '<tr><td>' + ip + '</td><td>' + ip_pas(conn, data[2]) + '</td><td>시작 : ' + data[4] + '<br>끝 : ' + end + '</td></tr>'
        div += '<tr><td colspan="3">' + why + '</td></tr>'

    div += '</tbody></table>'

    if(not name):
        if(not tool2):
            div += next_fix('/block_log?num=', num, data_list)
        else:
            div += next_fix('/block_log/' + url_pas(tool2) + '?num=', num, data_list)
    else:
        div += next_fix('/' + url_pas(tool) + '/' + url_pas(name) + '?num=', num, data_list)
                
    return(html_minify(template('index', 
        imp = ['차단 기록', wiki_set(conn, 1), custom(conn), other2([sub, 0])],
        data = div,
        menu = menu
    )))
            
@app.route('/search', methods=['POST'])
def search():
    return(redirect('/search/' + url_pas(request.form['search'])))

@app.route('/goto', methods=['POST'])
def goto():
    curs.execute("select title from data where title = ?", [request.form['search']])
    data = curs.fetchall()
    if(data):
        return(redirect('/w/' + url_pas(request.form['search'])))
    else:
        return(redirect('/search/' + url_pas(request.form['search'])))

@app.route('/search/<path:name>')
def deep_search(name = None):
    num = int(request.args.get('num', 1))
    if(num * 50 > 0):
        sql_num = num * 50 - 50
    else:
        sql_num = 0

    div = '<ul>'
    div_plus = ''
    no = 0
    start = 2
    test = ''

    curs.execute("select title from data where title = ?", [name])
    if(curs.fetchall()):
        div = '<ul><li>문서로 <a href="/w/' + url_pas(name) + '">바로가기</a></li></ul><hr><ul>'
    else:
        div = '<ul><li>문서가 없습니다. <a class="not_thing" href="/w/' + url_pas(name) + '">바로가기</a></li></ul><hr><ul>'

    curs.execute("select distinct title, case when title like ? then '제목' else '내용' " + \
                "end from data where title like ? or data like ? order by case when title like ? " + \
                "then 1 else 2 end limit ?, '50'", ['%' + name + '%', '%' + name + '%', '%' + name + '%', '%' + name + '%', str(sql_num)])
    all_list = curs.fetchall()
    if(all_list):
        test = all_list[0][1]
        for data in all_list:
            if(data[1] != test):
                div_plus += '</ul><hr><ul>'
                test = data[1]

            div_plus += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a> (' + data[1] + ')</li>'
    else:
        div += '<li>검색 결과 없음</li>'

    div += div_plus + '</ul>'
    div += next_fix('/search/' + url_pas(name) + '?num=', num, all_list)
    
    return(html_minify(template('index', 
        imp = [name, wiki_set(conn, 1), custom(conn), other2([' (검색)', 0])],
        data = div,
        menu = 0
    )))
         
@app.route('/raw/<path:name>')
@app.route('/topic/<path:name>/sub/<sub_t>/raw/<int:num>')
def raw_view(name = None, sub_t = None, num = None):
    v_name = name
    sub = ' (원본)'
    if(not num):
        num = request.args.get('num', None)
        if(num):
            num = int(num)
    
    if(not sub_t and num):
        curs.execute("select title from history where title = ? and id = ? and hide = 'O'", [name, str(num)])
        hid = curs.fetchall()
        if(hid and admin_check(conn, 6, None) != 1):
            return(re_error(conn, '/error/3'))
        
        curs.execute("select data from history where title = ? and id = ?", [name, str(num)])

        sub += ' (' + str(num) + '판)'
        menu = [['history/' + url_pas(name), '역사']]
    elif(sub_t):
        curs.execute("select data from topic where id = ? and title = ? and sub = ? and block = ''", [str(num), name, sub_t])

        v_name = '토론 원본'
        sub = ' (' + str(num) + '번)'
        menu = [['topic/' + url_pas(name) + '/sub/' + url_pas(sub_t) + '#' + str(num), '토론']]
    else:
        curs.execute("select data from data where title = ?", [name])
        
        menu = [['w/' + url_pas(name), '문서']]

    data = curs.fetchall()
    if(data):
        p_data = html.escape(data[0][0])
        
        p_data = '<textarea readonly rows="25">' + p_data + '</textarea>'
        
        return(html_minify(template('index', 
            imp = [v_name, wiki_set(conn, 1), custom(conn), other2([sub, 0])],
            data = p_data,
            menu = menu
        )))
    else:
        return(redirect('/w/' + url_pas(name)))
        
@app.route('/revert/<path:name>', methods=['POST', 'GET'])
def revert(name = None):
    ip = ip_check()
    can = acl_check(conn, name)
    today = get_time()
    num = int(request.args.get('num', 0))
    
    if(request.method == 'POST'):
        if(captcha_post(request.form.get('g-recaptcha-response', None), conn) == 1):
            return(re_error(conn, '/error/13'))
        else:
            captcha_post('', conn, 0)

        curs.execute("select title from history where title = ? and id = ? and hide = 'O'", [name, str(num)])
        if(curs.fetchall() and admin_check(conn, 6, None) != 1):
            return(re_error(conn, '/error/3'))

        if(can == 1):
            return(re_error(conn, '/ban'))

        curs.execute("delete from back where link = ?", [name])
        conn.commit()

        curs.execute("select data from history where title = ? and id = ?", [name, str(num)])
        data = curs.fetchall()
        if(data):                                
            curs.execute("select data from data where title = ?", [name])
            d = curs.fetchall()
            if(d):
                leng = leng_check(len(d[0][0]), len(data[0][0]))
                curs.execute("update data set data = ? where title = ?", [data[0][0], name])
            else:
                leng = '+' + str(len(data[0][0]))
                curs.execute("insert into data (title, data) values (?, ?)", [name, data[0][0]])
                
            history_plus(conn, name, data[0][0], today, ip, request.form['send'] + ' (' + str(num) + '판)', leng)
            
            namumark(conn, name, data[0][0], 1, 0, 0)
            conn.commit()
            
            return(redirect('/w/' + url_pas(name)))
    else:
        curs.execute("select title from history where title = ? and id = ? and hide = 'O'", [name, str(num)])
        hid = curs.fetchall()
        if(hid and admin_check(conn, 6, None) != 1):
            return(re_error(conn, '/error/3'))    
                          
        if(can == 1):
            return(re_error(conn, '/ban'))

        curs.execute("select title from history where title = ? and id = ?", [name, str(num)])
        if(not curs.fetchall()):
            return(redirect('/w/' + url_pas(name)))

        return(html_minify(template('index', 
            imp = [name, wiki_set(conn, 1), custom(conn), other2([' (되돌리기)', 0])],
            data =  '<form method="post"> \
                        <span>' + request.args.get('num', '0') + '판으로 되돌리겠습니까?</span><hr> \
                        ' + ip_warring(conn) + ' \
                        <input placeholder="사유" name="send" type="text"><hr> \
                        ' + captcha_get(conn) + ' \
                        <button class="btn btn-primary" type="submit">되돌리기</button> \
                    </form>',
            menu = [['history/' + url_pas(name), '역사'], ['recent_changes', '최근 변경']]
        )))            
                    
@app.route('/big_delete', methods=['POST', 'GET'])
def big_delete():
    if(admin_check(conn, 2, 'big_delete') != 1):
        return(re_error(conn, '/error/3'))

    if(request.method == 'POST'):
        today = get_time()
        ip = ip_check()

        data = request.form['content'] + '\r\n'
        m = re.findall('(.*)\r\n', data)
        for g in m:
            curs.execute("select data from data where title = ?", [g])
            d = curs.fetchall()
            if(d):
                curs.execute("delete from back where title = ?", [g])

                leng = '-' + str(len(d[0][0]))
                curs.execute("delete from data where title = ?", [g])
                history_plus(conn, g, '', today, ip, request.form['send'] + ' (대량 삭제)', leng)
            data = re.sub('(.*)\r\n', '', data, 1)
        conn.commit()

        return(redirect('/'))
    else:
        return(html_minify(template('index', 
            imp = ['많은 문서 삭제', wiki_set(conn, 1), custom(conn), other2([0, 0])],
            data = '<span>문서명 A<br>문서명 B<br>문서명 C<hr>이런 식으로 적으세요.</span><hr> \
                    <form method="post"> \
                        <textarea rows="25" name="content"></textarea><hr> \
                        <input placeholder="사유" name="send" type="text"><hr> \
                        <button class="btn btn-primary" type="submit">삭제</button> \
                    </form>',
            menu = [['manager', '관리자']]
        )))

@app.route('/edit_filter')
def edit_filter():
    div = '<ul>'
    
    curs.execute("select name from filter")
    data = curs.fetchall()
    for data_list in data:
        div += '<li><a href="/edit_filter/' + url_pas(data_list[0]) + '">' + data_list[0] + '</a></li>'

    div += '</ul>'
    div += '<hr><a href="/manager/9">(추가)</a>'

    return(html_minify(template('index', 
        imp = ['편집 필터 목록', wiki_set(conn, 1), custom(conn), other2([0, 0])],
        data = div,
        menu = [['manager', '관리자']]
    )))

@app.route('/edit_filter/<name>/delete', methods=['POST', 'GET'])
def delete_edit_filter(name = None):
    if(admin_check(conn, 1, 'edit_filter delete') != 1):
        return(re_error('/error/3'))

    curs.execute("delete from filter where name = ?", [name])
    conn.commit()

    return(redirect('/edit_filter'))

@app.route('/edit_filter/<name>', methods=['POST', 'GET'])
def set_edit_filter(name = None):
    if(request.method == 'POST'):
        if(admin_check(conn, 1, 'edit_filter edit') != 1):
            return(re_error('/error/3'))

        if(request.form['day'] == '00'):
            end = ''
        elif(request.form['day'] == '09'):
            end = 'X'
        else:
            end = request.form['day'] + ' ' + request.form['hour'] + ':' + request.form['minu'] 

        curs.execute("select name from filter where name = ?", [name])
        if(curs.fetchall()):
            curs.execute("update filter set regex = ?, sub = ? where name = ?", [request.form['content'], end, name])
        else:
            curs.execute("insert into filter (name, regex, sub) values (?, ?, ?)", [name, request.form['content'], end])
        conn.commit()
    
        return(redirect('/edit_filter/' + url_pas(name)))
    else:
        curs.execute("select regex, sub from filter where name = ?", [name])
        exist = curs.fetchall()
        if(exist):
            textarea = exist[0][0]

            match = re.search("^([^ ]+) ([^:]+):([^:]+)$", exist[0][1])
            if(match):
                end_data = match.groups()
            else:
                end_data = ['', '', '']
        else:
            textarea = ''
            end_data = ['', '', '']

        if(admin_check(conn, 1, None) != 1):
            stat = 'disabled'
        else:
            stat = ''

        day = '<option value="00">차단 X</option>'
        if(exist[0][1] == 'X'):
            day += '<option value="09" selected>영구</option>'
        else:
            day += '<option value="09">영구</option>'

        for i in range(0, 32):
            if(str(i) == end_data[0]):
                day += '<option value="' + str(i) + '" selected>' + str(i) + '</option>'
            else:
                day += '<option value="' + str(i) + '">' + str(i) + '</option>'

        hour = ''
        for i in range(0, 24):
            if(str(i) == end_data[1]):
                hour += '<option value="' + str(i) + '" selected>' + str(i) + '</option>'
            else:
                hour += '<option value="' + str(i) + '">' + str(i) + '</option>'

        minu = ''
        for i in range(0, 61):
            if(str(i) == end_data[2]):
                minu += '<option value="' + str(i) + '" selected>' + str(i) + '</option>'
            else:
                minu += '<option value="' + str(i) + '">' + str(i) + '</option>'

        data = '<select ' + stat + ' name="day">' + day + '</select> 일 '
        data += '<select ' + stat + ' name="hour">' + hour + '</select> 시 '
        data += '<select ' + stat + ' name="minu">' + minu + '</select> 분 동안<hr>'

        return(html_minify(template('index', 
            imp = [name, wiki_set(conn, 1), custom(conn), other2([' (편집 필터)', 0])],
            data = '<form method="post"> \
                        ' + data + ' \
                        <input ' + stat + ' placeholder="정규식" name="content" value="' + html.escape(textarea) + '" type="text"><hr> \
                        <button ' + stat + ' id="preview" class="btn btn-primary" type="submit">저장</button> \
                    </form>',
            menu = [['edit_filter', '목록'], ['edit_filter/' + url_pas(name) + '/delete', '삭제']]
        )))

@app.route('/edit/<path:name>', methods=['POST', 'GET'])
def edit(name = None):
    ip = ip_check()
    can = acl_check(conn, name)

    if(can == 1):
        return(re_error(conn, '/ban'))
    
    if(request.method == 'POST'):
        if(admin_check(conn, 1, 'edit_filter pass') != 1):
            curs.execute("select regex, sub from filter")
            data = curs.fetchall()
            for data_list in data:
                match = re.compile(data_list[0])
                if(match.search(request.form['content'])):
                    if(data_list[1] == 'X'):
                        curs.execute("insert into rb (block, end, today, blocker, why, band) values (?, ?, ?, ?, ?, '')", [ip, '', get_time(), '도구:편집 필터', '편집 필터에 의한 차단'])
                        curs.execute("insert into ban (block, end, why, band, login) values (?, '', ?, '', '')", [ip, '편집 필터에 의한 차단'])
                    elif(not data_list[1] == ''):
                        match = re.search("^([^ ]+) ([^:]+):([^:]+)$", data_list[1])
                        end_data = match.groups()

                        match = re.search("^([^-]+)-([^-]+)-([^ ]+) ([^:]+):([^:]+):(.+)$", get_time())
                        time_data = match.groups()

                        if(int(time_data[2]) + int(end_data[0]) > 29):
                            month = int(time_data[1]) + 1
                            day = int(time_data[2]) + int(end_data[0]) - 30

                            if(month > 12):
                                year = int(time_data[0]) + 1
                                month -= 12
                            else:
                                year = int(time_data[0])
                        else:
                            month = int(time_data[1])
                            day = int(time_data[2]) + int(end_data[0])
                            year = int(time_data[0])

                        if(int(time_data[3]) + int(end_data[1]) > 23):
                            day += 1
                            hour = int(time_data[3]) + int(end_data[1]) - 24
                        else:
                            hour = int(time_data[3]) + int(end_data[1])

                        if(int(time_data[4]) + int(end_data[2]) > 59):
                            hour += 1
                            minu = int(time_data[4]) + int(end_data[2]) - 60
                        else:
                            minu = int(time_data[4]) + int(end_data[2])

                        time_list = [month, day, hour, minu]
                        num = 0
                        for time_fix in time_list:
                            if(not re.search("[0-9]{2}", str(time_fix))):
                                time_list[num] = '0' + str(time_fix)   
                            else:
                                time_list[num] = str(time_fix)
                            
                            num += 1

                        end = str(year) + '-' + time_list[0] + '-' + time_list[1] + ' ' + time_list[2] + ':' + time_list[3] + ':' + time_data[5]

                        curs.execute("insert into rb (block, end, today, blocker, why, band) values (?, ?, ?, ?, ?, '')", [ip, end, get_time(), '도구:편집 필터', '편집 필터에 의한 차단'])
                        curs.execute("insert into ban (block, end, why, band, login) values (?, ?, ?, '', '')", [ip, end, '편집 필터에 의한 차단'])
                    
                    conn.commit()
                    return(re_error(conn, '/error/21'))

        if(captcha_post(request.form.get('g-recaptcha-response', None), conn) == 1):
            return(re_error(conn, '/error/13'))
        else:
            captcha_post('', conn, 0)

        if(len(request.form['send']) > 500):
            return(re_error(conn, '/error/15'))

        if(request.form['otent'] == request.form['content']):
            return(re_error(conn, '/error/18'))

        today = get_time()
        content = savemark(request.form['content'])

        curs.execute("select data from data where title = ?", [name])
        old = curs.fetchall()
        if(old):
            if(not request.args.get('section', None) and request.form['otent'] != old[0][0]):
                return(re_error(conn, '/error/12'))

            leng = leng_check(len(request.form['otent']), len(content))
            if(request.args.get('section', None)):
                content = old[0][0].replace(request.form['otent'], content)      
                
            curs.execute("update data set data = ? where title = ?", [content, name])
        else:
            leng = '+' + str(len(content))
            curs.execute("insert into data (title, data) values (?, ?)", [name, content])

        curs.execute("select user from scan where title = ?", [name])
        for user_data in curs.fetchall():
            curs.execute("insert into alarm (name, data, date) values (?, ?, ?)", [ip, ip + '님이 <a href="/w/' + url_pas(name) + '">' + name + '</a> 문서를 편집 했습니다.', today])

        history_plus(conn, name, content, today, ip, send_p(request.form['send']), leng)
        curs.execute("delete from back where link = ?", [name])
        curs.execute("delete from back where title = ? and type = 'no'", [name])
        namumark(conn, name, content, 1, 0, 0)
        conn.commit()
        
        return(redirect('/w/' + url_pas(name)))
    else:            
        curs.execute("select data from data where title = ?", [name])
        new = curs.fetchall()
        if(new):
            if(request.args.get('section', None)):
                i = 0
                j = 0
                
                data = new[0][0] + '\r\n'
                while(1):
                    m = re.search("((?:={1,6})\s?(?:[^=]*)\s?(?:={1,6})(?:\s+)?\n(?:(?:(?:(?!(?:={1,6})\s?(?:[^=]*)\s?(?:={1,6})(?:\s+)?\n).)*)(?:\n)?)+)", data)
                    if(m):
                        if(i == int(request.args.get('section', 0)) - 1):
                            g = m.groups()
                            data = re.sub("\r\n$", "", g[0])
                            
                            break
                        else:
                            data = re.sub("((?:={1,6})\s?(?:[^=]*)\s?(?:={1,6})(?:\s+)?\n(?:(?:(?:(?!(?:={1,6})\s?(?:[^=]*)\s?(?:={1,6})(?:\s+)?\n).)*)(?:\n)?)+)", "", data, 1)
                            
                            i += 1
                    else:
                        j = 1
                        
                        break
                        
                if(j == 0):
                    data = re.sub("\r\n$", "", data)
            else:
                data = new[0][0]
        else:
            data = ''
            
        data2 = data
        if(not request.args.get('section', None)):
            get_name = '<form method="post" id="get_edit" action="/edit_get/' + url_pas(name) + '"> \
                            <input placeholder="불러 올 문서" name="name" style="width: 50%;" type="text"> \
                            <button id="preview" class="btn" type="submit">불러오기</button> \
                        </form><hr>'
            action = ''
        else:
            get_name = ''
            action = '?section=' + request.args.get('section', None)
            
        if(request.args.get('froms', None)):
            curs.execute("select data from data where title = ?", [request.args.get('froms', None)])
            get_data = curs.fetchall()
            if(get_data):
                data = get_data[0][0]
                get_name = ''

        return(html_minify(template('index', 
            imp = [name, wiki_set(conn, 1), custom(conn), other2([' (수정)', 0])],
            data = get_name + ' \
                    <form method="post" action="/edit/' + url_pas(name) + action + '"> \
                        <textarea rows="25" name="content">' + html.escape(data) + '</textarea> \
                        <textarea style="display: none;" name="otent">' + html.escape(data2) + '</textarea><hr> \
                        <input placeholder="사유" name="send" type="text"><hr> \
                        ' + captcha_get(conn) + ' \
                        ' + ip_warring(conn) + ' \
                        <button id="preview" class="btn btn-primary" type="submit">저장</button> \
                        <button id="preview" class="btn" type="submit" formaction="/preview/' + url_pas(name) + action + '">미리보기</button> \
                    </form>',
            menu = [['w/' + url_pas(name), '문서']]
        )))
        
@app.route('/edit_get/<path:name>', methods=['POST'])
def edit_get(name = None):
    return(redirect('/edit/' + url_pas(name) + '?froms=' + url_pas(request.form['name'])))

@app.route('/preview/<path:name>', methods=['POST'])
def preview(name = None):
    ip = ip_check()
    can = acl_check(conn, name)
    
    if(can == 1):
        return(re_error(conn, '/ban'))
         
    newdata = request.form['content']
    newdata = re.sub('^#(?:redirect|넘겨주기) (?P<in>[^\n]*)', ' * [[\g<in>]] 문서로 넘겨주기', newdata)
    enddata = namumark(conn, name, newdata, 0, 0, 1)

    if(request.args.get('section', None)):
        action = '?section=' + request.args.get('section', None)
    else:
        action = ''

    return(html_minify(template('index', 
        imp = [name, wiki_set(conn, 1), custom(conn), other2([' (미리보기)', 0])],
        data = '<form method="post" action="/edit/' + url_pas(name) + action + '"> \
                    <textarea rows="25" name="content">' + html.escape(request.form['content']) + '</textarea> \
                    <textarea style="display: none;" name="otent">' + html.escape(request.form['otent']) + '</textarea><hr> \
                    <input placeholder="사유" name="send" type="text"><hr> \
                    ' + captcha_get(conn) + ' \
                    <button id="preview" class="btn btn-primary" type="submit">저장</button> \
                    <button id="preview" class="btn" type="submit" formaction="/preview/' + url_pas(name) + action + '">미리보기</button> \
                </form><hr>' + enddata,
        menu = [['w/' + url_pas(name), '문서']]
    )))
        
@app.route('/delete/<path:name>', methods=['POST', 'GET'])
def delete(name = None):
    ip = ip_check()
    can = acl_check(conn, name)

    if(can == 1):
        return(re_error(conn, '/ban'))
    
    if(request.method == 'POST'):
        if(captcha_post(request.form.get('g-recaptcha-response', None), conn) == 1):
            return(re_error(conn, '/error/13'))
        else:
            captcha_post('', conn, 0)

        curs.execute("select data from data where title = ?", [name])
        data = curs.fetchall()
        if(data):
            today = get_time()
            
            leng = '-' + str(len(data[0][0]))
            history_plus(conn, name, '', today, ip, request.form['send'] + ' (삭제)', leng)

            curs.execute("select title, link from back where title = ? and not type = 'cat' and not type = 'no'", [name])
            for data in curs.fetchall():
                curs.execute("insert into back (title, link, type) values (?, ?, 'no')", [data[0], data[1]])
            
            curs.execute("delete from back where link = ?", [name])
            curs.execute("delete from data where title = ?", [name])
            conn.commit()
            
        return(redirect('/w/' + url_pas(name)))
    else:
        curs.execute("select title from data where title = ?", [name])
        if(not curs.fetchall()):
            return(redirect('/w/' + url_pas(name)))

        return(html_minify(template('index', 
            imp = [name, wiki_set(conn, 1), custom(conn), other2([' (삭제)', 0])],
            data = '<form method="post"> \
                        ' + ip_warring(conn) + ' \
                        <input placeholder="사유" name="send" type="text"><hr> \
                        ' + captcha_get(conn) + ' \
                        <button class="btn btn-primary" type="submit">삭제</button> \
                    </form>',
            menu = [['w/' + url_pas(name), '문서']]
        )))            
            
@app.route('/move_data/<path:name>')
def move_data(name = None):    
    data = '<ul>'

    curs.execute("select send, date, ip from history where send like ? or send like ? order by date desc", ['%<a href="/w/' + url_pas(name) + '">' + name + '</a> 이동)%', '%(<a href="/w/' + url_pas(name) + '">' + name + '</a>%'])
    for for_data in curs.fetchall():
        match = re.findall('<a href="\/w\/(?:(?:(?!">).)+)">((?:(?!<\/a>).)+)<\/a>', for_data[0])
        send = re.sub('\([^\)]+\)$', '', for_data[0])
        
            

        data += '<li><a href="/move_data/' + url_pas(match[0]) + '">' + match[0] + '</a> - <a href="/move_data/' + url_pas(match[1]) + '">' + match[1] + '</a>'
        if(re.search('^( *)+$', send)):
            data += ' / ' + for_data[2] + ' / ' + for_data[1] + '</li>'
        else:
            data += ' / ' + for_data[2] + ' / ' + for_data[1] + ' / ' + send + '</li>'
    
    data += '</ul>'
    
    return(html_minify(template('index', 
        imp = [name, wiki_set(conn, 1), custom(conn), other2([' (이동 기록)', 0])],
        data = data,
        menu = [['history/' + url_pas(name), '역사']]
    )))        
            
@app.route('/move/<path:name>', methods=['POST', 'GET'])
def move(name = None):
    ip = ip_check()
    can = acl_check(conn, name)
    today = get_time()

    if(can == 1):
        return(re_error(conn, '/ban'))
    
    if(request.method == 'POST'):
        if(captcha_post(request.form.get('g-recaptcha-response', None), conn) == 1):
            return(re_error(conn, '/error/13'))
        else:
            captcha_post('', conn, 0)

        curs.execute("select title from history where title = ?", [request.form['title']])
        if(curs.fetchall()):
            return(re_error(conn, '/error/19'))
        
        curs.execute("select data from data where title = ?", [name])
        data = curs.fetchall()

        leng = '0'
        if(data):            
            curs.execute("update data set title = ? where title = ?", [request.form['title'], name])
            curs.execute("update back set link = ? where link = ?", [request.form['title'], name])
            
            d = data[0][0]
        else:
            d = ''
            
        history_plus(conn, name, d, today, ip, request.form['send'] + ' (<a href="/w/' + url_pas(name) + '">' + name + '</a> - <a href="/w/' + url_pas(request.form['title']) + '">' + request.form['title'] + '</a> 이동)', leng)

        curs.execute("select title, link from back where title = ? and not type = 'cat' and not type = 'no'", [name])
        for data in curs.fetchall():
            curs.execute("insert into back (title, link, type) values (?, ?, 'no')", [data[0], data[1]])
            
        curs.execute("update history set title = ? where title = ?", [request.form['title'], name])
        conn.commit()
        
        return(redirect('/w/' + url_pas(request.form['title'])))
    else:            
        return(html_minify(template('index', 
            imp = [name, wiki_set(conn, 1), custom(conn), other2([' (이동)', 0])],
            data = '<form method="post"> \
                        ' + ip_warring(conn) + ' \
                        <input placeholder="문서명" value="' + name + '" name="title" type="text"><hr> \
                        <input placeholder="사유" name="send" type="text"><hr> \
                        ' + captcha_get(conn) + ' \
                        <button class="btn btn-primary" type="submit">이동</button> \
                    </form>',
            menu = [['w/' + url_pas(name), '문서']]
        )))
            
@app.route('/other')
def other():
    return(html_minify(template('index', 
        imp = ['기타 메뉴', wiki_set(conn, 1), custom(conn), other2([0, 0])],
        data = namumark(conn, '', '[목차(없음)]\r\n' + \
                            '== 기록 ==\r\n' + \
                            ' * [[wiki:block_log|차단 기록]]\r\n' + \
                            ' * [[wiki:user_log|가입 기록]]\r\n' + \
                            ' * [[wiki:admin_log|권한 사용 기록]]\r\n' + \
                            ' * [[wiki:manager/6|편집 기록]]\r\n' + \
                            ' * [[wiki:manager/7|토론 기록]]\r\n' + \
                            '== 목록 ==\r\n' + \
                            ' * [[wiki:admin_list|관리자 목록]]\r\n' + \
                            ' * [[wiki:give_log|권한 목록]]\r\n' + 
                            ' * [[wiki:not_close_topic|열린 토론 목록]]\r\n' + \
                            '== 기타 ==\r\n' + \
                            ' * [[wiki:title_index|모든 문서]]\r\n' + \
                            ' * [[wiki:acl_list|ACL 문서]]\r\n' + \
                            ' * [[wiki:please|필요한 문서]]\r\n' + \
                            ' * [[wiki:upload|파일 올리기]]\r\n' + \
                            ' * [[wiki:manager/10|문서 검색]]\r\n' + \
                            '== 관리자 ==\r\n' + \
                            ' * [[wiki:manager/1|관리자 메뉴]]\r\n' + \
                            '== 버전 ==\r\n' + \
                            ' * 이 오픈나무는 [[https://github.com/2DU/openNAMU/blob/master/version.md|' + r_ver + ']] 입니다.', 0, 0, 0),
        menu = 0
    )))
    
@app.route('/manager', methods=['POST', 'GET'])
@app.route('/manager/<int:num>', methods=['POST', 'GET'])
def manager(num = 1):
    title_list = [['문서 ACL', '문서명', 'acl'], ['사용자 검사', 0, 'check'], ['사용자 차단', 0, 'ban'], ['권한 주기', 0, 'admin'], ['편집 기록', 0, 'record'], ['토론 기록', 0, 'topic_record'], ['그룹 생성', '그룹명', 'admin_plus'], ['편집 필터 생성', '필터명', 'edit_filter'], ['문서 검색', '문서명', 'search'], ['차단자 검색', 0, 'block_user'], ['관리자 검색', 0, 'block_admin'], ['주시 문서 추가', '문서명', 'watch_list']]
    if(num == 1):
        return(html_minify(template('index', 
            imp = ['관리자 메뉴', wiki_set(conn, 1), custom(conn), other2([0, 0])],
            data = namumark(conn, '',   '[목차(없음)]\r\n' + \
                                        '== 목록 ==\r\n' + \
                                        ' * [[wiki:manager/2|문서 ACL]]\r\n' + \
                                        ' * [[wiki:manager/3|사용자 검사]]\r\n' + \
                                        ' * [[wiki:manager/100|사용자 비교]]\r\n' + \
                                        ' * [[wiki:manager/4|사용자 차단]]\r\n' + \
                                        ' * [[wiki:manager/5|권한 주기]]\r\n' + \
                                        ' * [[wiki:big_delete|여러 문서 삭제]]\r\n' + \
                                        ' * [[wiki:edit_filter|편집 필터]]\r\n' + \
                                        '== 소유자 ==\r\n' + \
                                        ' * [[wiki:indexing|인덱싱]]\r\n' + \
                                        ' * [[wiki:manager/8|관리 그룹 생성]]\r\n' + \
                                        ' * [[wiki:edit_set|설정 편집]]\r\n' + \
                                        '== 기타 ==\r\n' + \
                                        ' * 이 메뉴에 없는 기능은 해당 문서의 역사나 토론에서 바로 사용 가능함', 0, 0, 0),
            menu = [['other', '기타']]
        )))
    elif(num in range(2, 14)):
        if(request.method == 'POST'):
            return(redirect('/' + title_list[(num - 2)][2] + '/' + url_pas(request.form['name'])))
        else:
            if(title_list[(num - 2)][1] == 0):
                placeholder = '사용자명'
            else:
                placeholder = title_list[(num - 2)][1]

            return(html_minify(template('index', 
                imp = [title_list[(num - 2)][0], wiki_set(conn, 1), custom(conn), other2([0, 0])],
                data = '<form method="post"> \
                            <input placeholder="' + placeholder + '" name="name" type="text"><hr> \
                            <button class="btn btn-primary" type="submit">이동</button> \
                        </form>',
                menu = [['manager', '관리자']]
            )))
    elif(num == 100):
        if(request.method == 'POST'):
            return(redirect('/check/' + url_pas(request.form['name']) + '?plus=' + url_pas(request.form['name2'])))
        else:
            return(html_minify(template('index', 
                imp = ['검사', wiki_set(conn, 1), custom(conn), other2([0, 0])],
                data = '<form method="post"> \
                            <input placeholder="사용자명" name="name" type="text"><hr> \
                            <input placeholder="비교 대상" name="name2" type="text"><hr> \
                            <button class="btn btn-primary" type="submit">이동</button> \
                        </form>',
                menu = [['manager', '관리자']]
            )))
    else:
        return(redirect('/'))
        
@app.route('/title_index')
def title_index():
    page = int(request.args.get('page', 1))
    num = int(request.args.get('num', 100))
    if(page * num > 0):
        sql_num = page * num - num
    else:
        sql_num = 0

    if(num != 0):
        all_list = sql_num + 1
    else:
        all_list = 1

    if(num > 1000):
        return(re_error(conn, '/error/3'))

    data = '<ul><a href="/title_index?num=0">(전체)</a> <a href="/title_index?num=250">(250)</a> <a href="/title_index?num=500">(500)</a> <a href="/title_index?num=1000">(1000)</a><hr>'

    if(num == 0):
        curs.execute("select data from other where name = 'all_title'")
        all_title_can = curs.fetchall()
        if(all_title_can and all_title_can[0][0] != ''):
            return(re_error(conn, '/error/3'))

        curs.execute("select title from data order by title asc")
    else:
        curs.execute("select title from data order by title asc limit ?, ?", [str(sql_num), str(num)])
    title_list = curs.fetchall()

    for list_data in title_list:
        data += '<li>' + str(all_list) + '. <a href="/w/' + url_pas(list_data[0]) + '">' + list_data[0] + '</a></li>'        
        all_list += 1

    if(page == 1):
        count_end = []
        curs.execute("select count(title) from data")
        count = curs.fetchall()
        if(count):
            count_end += [count[0][0]]
        else:
            count_end += [0]

        sql_list = ['틀:', '분류:', '사용자:', '파일:']

        for sql in sql_list:
            curs.execute("select count(title) from data where title like ?", [sql + '%'])
            count = curs.fetchall()
            if(count):
                count_end += [count[0][0]]
            else:
                count_end += [0]

        count_end += [count_end[0] - count_end[1]  - count_end[2]  - count_end[3]  - count_end[4]]

        data += '</ul><hr><ul><li>이 위키에는 총 ' + str(count_end[0]) + '개의 문서가 있습니다.</li></ul><hr><ul>'
        data += '<li>틀 문서는 총 ' + str(count_end[1]) + '개의 문서가 있습니다.</li>'
        data += '<li>분류 문서는 총 ' + str(count_end[2]) + '개의 문서가 있습니다.</li>'
        data += '<li>사용자 문서는 총 ' + str(count_end[3]) + '개의 문서가 있습니다.</li>'
        data += '<li>파일 문서는 총 ' + str(count_end[4]) + '개의 문서가 있습니다.</li>'
        data += '<li>나머지 문서는 총 ' + str(count_end[5]) + '개의 문서가 있습니다.</li>'

    if(num != 0):
        data += '</ul>' + next_fix('/title_index?num=' + str(num) + '&page=', page, title_list, num)
    
    if(' (' + str(num) + '개)' == ' (0개)'):
        sub = 0
    else:
        sub = ' (' + str(num) + '개)'
    
    return(html_minify(template('index', 
        imp = ['모든 문서', wiki_set(conn, 1), custom(conn), other2([sub, 0])],
        data = data,
        menu = [['other', '기타']]
    )))
        
@app.route('/topic/<path:name>/sub/<sub>/b/<int:num>')
def topic_block(name = None, sub = None, num = None):
    if(admin_check(conn, 3, 'blind (' + name + ' - ' + sub + '#' + str(num) + ')') != 1):
        return(re_error(conn, '/error/3'))

    curs.execute("select block from topic where title = ? and sub = ? and id = ?", [name, sub, str(num)])
    block = curs.fetchall()
    if(block):
        if(block[0][0] == 'O'):
            curs.execute("update topic set block = '' where title = ? and sub = ? and id = ?", [name, sub, str(num)])
        else:
            curs.execute("update topic set block = 'O' where title = ? and sub = ? and id = ?", [name, sub, str(num)])
        
        rd_plus(conn, name, sub, get_time())
        conn.commit()
        
    return(redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#' + str(num)))
        
@app.route('/topic/<path:name>/sub/<sub>/notice/<int:num>')
def topic_top(name = None, sub = None, num = None):
    if(admin_check(conn, 3, 'notice (' + name + ' - ' + sub + '#' + str(num) + ')') != 1):
        return(re_error(conn, '/error/3'))

    curs.execute("select * from topic where title = ? and sub = ? and id = ?", [name, sub, str(num)])
    topic_data = curs.fetchall()
    if(topic_data):
        curs.execute("select top from topic where id = ? and title = ? and sub = ?", [str(num), name, sub])
        top_data = curs.fetchall()
        if(top_data):
            if(top_data[0][0] == 'O'):
                curs.execute("update topic set top = '' where title = ? and sub = ? and id = ?", [name, sub, str(num)])
            else:
                curs.execute("update topic set top = 'O' where title = ? and sub = ? and id = ?", [name, sub, str(num)])
        
        rd_plus(conn, name, sub, get_time())
        conn.commit()

    return(redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#' + str(num)))        
        
@app.route('/topic/<path:name>/sub/<sub>/tool/<regex("close|stop|agree"):tool>')
def topic_stop(name = None, sub = None, tool = None):
    if(tool == 'close'):
        set_list = ['O', '', '토론 닫기', '토론 열림']
    elif(tool == 'stop'):
        set_list = ['', 'O', '토론 정지', '토론 재개']
    elif(tool == 'agree'):
        pass
    else:
        return(redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub)))

    if(admin_check(conn, 3, 'topic ' + tool + ' (' + name + ' - ' + sub + ')') != 1):
        return(re_error(conn, '/error/3'))

    ip = ip_check()
    time = get_time()
    
    curs.execute("select id from topic where title = ? and sub = ? order by id + 0 desc limit 1", [name, sub])
    topic_check = curs.fetchall()
    if(topic_check):
        if(tool == 'agree'):
            curs.execute("select title from agreedis where title = ? and sub = ?", [name, sub])
            if(curs.fetchall()):
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, '합의 결렬', ?, ?, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, time, ip])
                curs.execute("delete from agreedis where title = ? and sub = ?", [name, sub])
            else:
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, '합의 완료', ?, ?, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, time, ip])
                curs.execute("insert into agreedis (title, sub) values (?, ?)", [name, sub])
        else:
            curs.execute("select title from stop where title = ? and sub = ? and close = ?", [name, sub, set_list[0]])
            if(curs.fetchall()):
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, ?, ?, ?, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, set_list[3], time, ip])
                curs.execute("delete from stop where title = ? and sub = ? and close = ?", [name, sub, set_list[0]])
            else:
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, ?, ?, ?, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, set_list[2], time, ip])
                curs.execute("insert into stop (title, sub, close) values (?, ?, ?)", [name, sub, set_list[0]])
                curs.execute("delete from stop where title = ? and sub = ? and close = ?", [name, sub, set_list[1]])
        
        rd_plus(conn, name, sub, time)
        conn.commit()
        
    return(redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub)))    

@app.route('/topic/<path:name>/sub/<sub>/admin/<int:num>')
def topic_admin(name = None, sub = None, num = None):
    curs.execute("select block, ip, date from topic where title = ? and sub = ? and id = ?", [name, sub, str(num)])
    data = curs.fetchall()
    if(not data):
        return(redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub)))

    ban = ''
    if(admin_check(conn, 3, None) == 1):
        ban += '== 관리 도구 ==\r\n'

        is_ban = ' * [[wiki:topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/b/' + str(num) + '|'
        if(data[0][0] == 'O'):
            is_ban += '가림 해제'
        else:
            is_ban += '가림'
        is_ban += ']]\r\n'

        curs.execute("select id from topic where title = ? and sub = ? and id = ? and top = 'O'", [name, sub, str(num)])
        is_ban += ' * [[wiki:topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/notice/' + str(num) + '|'
        if(curs.fetchall()):
            is_ban += '공지 해제'
        else:
            is_ban += '공지'
        is_ban += ']]\r\n'

        curs.execute("select end from ban where block = ?", [data[0][1]])
        ban += ' * [[wiki:/ban/' + url_pas(data[0][1]) + '|'
        if(curs.fetchall()):
            ban += '차단 해제'
        else:
            ban += '차단'
        ban += ']]\r\n' + is_ban

    ban += '== 기타 도구 ==\r\n'
    ban += ' * [[wiki:/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/raw/' + str(num) + '|원본]]\r\n'
    ban = ' * 작성 시간 : ' + data[0][2] + '\r\n' + ban

    if(re.search('(\.|:)', data[0][1])):
        ban = ' * 작성인 : ' + data[0][1] + ' [[wiki:record/' + url_pas(data[0][1]) + '|(기록)]]\r\n' + ban
    else:
        ban = ' * 작성인 : [[사용자:' + data[0][1] + '|' + data[0][1] + ']] [[wiki:record/' + url_pas(data[0][1]) + '|(기록)]]\r\n' + ban

    ban = '== 정보 ==\r\n' + ban
    ban = '[목차(없음)]\r\n' + ban

    return(html_minify(template('index', 
        imp = ['토론 도구', wiki_set(conn, 1), custom(conn), other2([' (' + str(num) + '번)', 0])],
        data = namumark(conn, '', ban, 0, 0, 0),
        menu = [['topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#' + str(num), '토론']]
    )))

@app.route('/topic/<path:name>/sub/<sub>', methods=['POST', 'GET'])
def topic(name = None, sub = None):
    ban = topic_check(conn, name, sub)
    admin = admin_check(conn, 3, None)
    
    if(request.method == 'POST'):
        if(captcha_post(request.form.get('g-recaptcha-response', None), conn) == 1):
            return(re_error(conn, '/error/13'))
        else:
            captcha_post('', conn, 0)

        ip = ip_check()
        today = get_time()

        if(ban == 1 and admin != 1):
            return(re_error(conn, '/ban'))
        
        curs.execute("select id from topic where title = ? and sub = ? order by id + 0 desc limit 1", [name, sub])
        old_num = curs.fetchall()
        if(old_num):
            num = int(old_num[0][0]) + 1
        else:
            num = 1

        match = re.search('^사용자:([^/]+)', name)
        if(match):
            curs.execute('insert into alarm (name, data, date) values (?, ?, ?)', [match.groups()[0], ip + '님이 <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '">사용자 토론</a>을 시작했습니다.', today])
        
        data = re.sub("\[\[(분류:(?:(?:(?!\]\]).)*))\]\]", "[br]", request.form['content'])
        match = re.findall("(?:#([0-9]+))", data)
        for rd_data in match:
            curs.execute("select ip from topic where title = ? and sub = ? and id = ?", [name, sub, rd_data])
            ip_data = curs.fetchall()
            if(ip_data and not re.search('(\.|:)', ip_data[0][0])):
                curs.execute('insert into alarm (name, data, date) values (?, ?, ?)', [ip_data[0][0], ip + '님이 <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#' + str(num) + '">토론</a>에서 언급 했습니다.', today])
            data = re.sub("(?P<in>#(?:[0-9]+))", '[[\g<in>]]', data)

        data = savemark(data)
        rd_plus(conn, name, sub, today)
        curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, ?, ?, ?, '', '')", [str(num), name, sub, data, today, ip])
        conn.commit()
        
        return(redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub)))
    else:
        curs.execute("select title from stop where title = ? and sub = ? and close = 'O'", [name, sub])
        close_data = curs.fetchall()

        curs.execute("select title from stop where title = ? and sub = ? and close = ''", [name, sub])
        stop_data = curs.fetchall()

        curs.execute("select id from topic where title = ? and sub = ? limit 1", [name, sub])
        topic_exist = curs.fetchall()
        
        display = ''
        all_data = ''
        data = ''
        number = 1

        if(admin == 1 and topic_exist):
            if(close_data):
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/close">(열기)</a> '
            else:
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/close">(닫기)</a> '
            
            if(stop_data):
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/stop">(재개)</a> '
            else:
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/stop">(정지)</a> '

            curs.execute("select title from agreedis where title = ? and sub = ?", [name, sub])
            if(curs.fetchall()):
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/agree">(취소)</a>'
            else:
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/agree">(합의)</a>'
            
            all_data += '<hr>'
        
        if((close_data or stop_data) and admin != 1):
            display = 'display: none;'
        
        curs.execute("select data, id, date, ip, block, top from topic where title = ? and sub = ? order by id + 0 asc", [name, sub])
        topic_1 = curs.fetchall()

        curs.execute("select data, id, date, ip from topic where title = ? and sub = ? and top = 'O' order by id + 0 asc", [name, sub])
        topic_2 = curs.fetchall()

        for topic_data in topic_2:                   
            who_plus = ''
            curs.execute("select who from re_admin where what = ? order by time desc limit 1", ['notice (' + name + ' - ' + sub + '#' + topic_data[1] + ')'])
            topic_data_top = curs.fetchall()
            if(topic_data_top):
                who_plus += ' <span style="margin-right: 5px;">@' + topic_data_top[0][0] + ' </span>'
                                
            all_data += '<table id="toron"><tbody><tr><td id="toron_color_red">'
            all_data += '<a href="#' + topic_data[1] + '">#' + topic_data[1] + '</a> ' + ip_pas(conn, topic_data[3]) + who_plus + ' <span style="float: right;">' + topic_data[2] + '</span>'
            all_data += '</td></tr><tr><td>' + namumark(conn, '', topic_data[0], 0, 0, 0) + '</td></tr></tbody></table><br>'    

        for topic_data in topic_1:
            if(number == 1):
                start = topic_data[3]

            if(topic_data[4] == 'O'):
                blind_data = 'style="background: gainsboro;"'
                if(admin != 1):
                    curs.execute("select who from re_admin where what = ? order by time desc limit 1", ['blind (' + name + ' - ' + sub + '#' + str(number) + ')'])
                    who_blind = curs.fetchall()
                    if(who_blind):
                        user_write = '[[사용자:' + who_blind[0][0] + ']]님이 가림'
                    else:
                        user_write = '관리자가 가림'
            else:
                blind_data = ''

            user_write = namumark(conn, '', topic_data[0], 0, 0, 0)
            ip = ip_pas(conn, topic_data[3])

            curs.execute('select acl from user where id = ?', [topic_data[3]])
            user_acl = curs.fetchall()
            if(user_acl and user_acl[0][0] != 'user'):
                ip += ' <a href="javascript:void(0);" title="관리자">★</a>'

            if(admin == 1 or blind_data == ''):
                ip += ' <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/admin/' + str(number) + '">(도구)</a>'

            curs.execute("select end from ban where block = ?", [topic_data[3]])
            if(curs.fetchall()):
                ip += ' <a href="javascript:void(0);" title="차단자">†</a>'
                    
            if(topic_data[5] == '1'):
                color = '_blue'
            elif(topic_data[3] == start):
                color = '_green'
            else:
                color = ''
                
            if(user_write == ''):
                user_write = '<br>'
                         
            all_data += '<table id="toron"><tbody><tr><td id="toron_color' + color + '">'
            all_data += '<a href="javascript:void(0);" id="' + str(number) + '">#' + str(number) + '</a> ' + ip + '</span>'
            all_data += '</td></tr><tr ' + blind_data + '><td>' + user_write + '</td></tr></tbody></table><br>'
            number += 1

        if(ban != 1 or admin == 1):
            data += '<a id="reload" href="javascript:void(0);" onclick="location.href.endsWith(\'#reload\') ?  location.reload(true) : location.href = \'#reload\'"><i aria-hidden="true" class="fa fa-refresh"></i></a>'
            data += '<form style="' + display + '" method="post"><br><textarea style="height: 100px;" name="content"></textarea><hr>' + captcha_get(conn)
            
            if(display == ''):
                data += ip_warring(conn)

            data += '<button class="btn btn-primary" type="submit">전송</button></form>'

        return(html_minify(template('index', 
            imp = [name, wiki_set(conn, 1), custom(conn), other2([' (토론)', 0])],
            data = '<h2 id="topic_top_title">' + sub + '</h2><hr id="under_bar" style="margin-top: -5px;">' + all_data + data,
            menu = [['topic/' + url_pas(name), '목록']]
        )))
        
@app.route('/topic/<path:name>', methods=['POST', 'GET'])
@app.route('/topic/<path:name>/<regex("close|agree"):tool>', methods=['GET'])
def close_topic_list(name = None, tool = None):
    div = ''
    list_d = 0

    if(request.method == 'POST'):
        t_num = ''
        while(1):
            curs.execute("select title from topic where title = ? and sub = ? limit 1", [name, request.form['topic'] + t_num])
            if(curs.fetchall()):
                if(t_num == ''):
                    t_num = ' 2'
                else:
                    t_num = ' ' + str(int(t_num.replace(' ', '')) + 1)
            else:
                break

        return(redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(request.form['topic'] + t_num)))
    else:
        plus = ''
        menu = [['topic/' + url_pas(name), '목록']]
        if(tool == 'close'):
            curs.execute("select sub from stop where title = ? and close = 'O' order by sub asc", [name])
            sub = '닫힘'
        elif(tool == 'agree'):
            curs.execute("select sub from agreedis where title = ? order by sub asc", [name])
            sub = '합의'
        else:
            curs.execute("select sub from rd where title = ? order by date desc", [name])
            sub = '토론 목록'
            menu = [['w/' + url_pas(name), '문서']]
            plus =  '<a href="/topic/' + url_pas(name) + '/close">(닫힘)</a> <a href="/topic/' + url_pas(name) + '/agree">(합의)</a><hr> \
                    <input placeholder="토론명" class="form-control" name="topic" type="text"><hr> \
                    <button class="btn btn-primary" type="submit">만들기</button>'

        for data in curs.fetchall():
            curs.execute("select data, date, ip, block from topic where title = ? and sub = ? and id = '1'", [name, data[0]])
            if(curs.fetchall()):                
                it_p = 0
                if(sub == '토론 목록'):
                    curs.execute("select title from stop where title = ? and sub = ? and close = 'O' order by sub asc", [name, data[0]])
                    close = curs.fetchall()
                    if(close):
                        it_p = 1
                
                if(it_p != 1):
                    div += '<h2><a href="/topic/' + url_pas(name) + '/sub/' + url_pas(data[0]) + '">' + data[0] + '</a></h2><hr id="under_bar" style="margin-top: -5px;">'

        if(div == ''):
            plus = re.sub('^<br>', '', plus)
        
        return(html_minify(template('index', 
            imp = [name, wiki_set(conn, 1), custom(conn), other2([' (' + sub + ')', 0])],
            data =  '<form method="post">' + div + plus + '</form>',
            menu = menu
        )))
        
@app.route('/login', methods=['POST', 'GET'])
def login():
    if('Now' in session and session['Now'] == 1):
        return(re_error(conn, '/error/11'))

    ip = ip_check()
    agent = request.headers.get('User-Agent')
    
    curs.execute("select block from ban where block = ? and login = 'O'", [ip])
    if(not curs.fetchall()):
        ban = ban_check(conn)
    else:
        ban = 0

    if(ban == 1):
        return(re_error(conn, '/ban'))
        
    if(request.method == 'POST'):        
        if(captcha_post(request.form.get('g-recaptcha-response', None), conn) == 1):
            return(re_error(conn, '/error/13'))
        else:
            captcha_post('', conn, 0)

        curs.execute("select pw from user where id = ?", [request.form['id']])
        user = curs.fetchall()
        if(not user):
            return(re_error(conn, '/error/5'))

        if(not bcrypt.checkpw(bytes(request.form['pw'], 'utf-8'), bytes(user[0][0], 'utf-8'))):
            return(re_error(conn, '/error/10'))

        session['Now'] = 1
        session['DREAMER'] = request.form['id']

        curs.execute("select css from custom where user = ?", [request.form['id']])
        css_data = curs.fetchall()
        if(css_data):
            session['Daydream'] = css_data[0][0]
        else:
            session['Daydream'] = ''
        
        curs.execute("insert into ua_d (name, ip, ua, today, sub) values (?, ?, ?, ?, '')", [request.form['id'], ip, agent, get_time()])
        conn.commit()
        
        return(redirect('/user'))                            
    else:        
        return(html_minify(template('index',    
            imp = ['로그인', wiki_set(conn, 1), custom(conn), other2([0, 0])],
            data = '<form method="post"> \
                        <input placeholder="아이디" name="id" type="text"><hr> \
                        <input placeholder="비밀번호" name="pw" type="password"><hr> \
                        ' + captcha_get(conn) + ' \
                        <button class="btn btn-primary" type="submit">로그인</button><hr> \
                        <span>주의 : 만약 HTTPS 연결이 아닌 경우 데이터가 유출될 가능성이 있습니다. 이에 대해 책임지지 않습니다.</span> \
                    </form>',
            menu = [['user', '사용자']]
        )))
                
@app.route('/change', methods=['POST', 'GET'])
def change_password():
    ip = ip_check()
    ban = ban_check(conn)
    
    if(request.method == 'POST'):    
        if(request.form['pw2'] != request.form['pw3']):
            return(re_error(conn, '/error/20'))

        if(ban == 1):
            return(re_error(conn, '/ban'))

        curs.execute("select pw from user where id = ?", [session['DREAMER']])
        user = curs.fetchall()
        if(not user):
            return(re_error(conn, '/error/10'))

        if(re.search('(\.|:)', ip)):
            return(redirect('/login'))

        if(not bcrypt.checkpw(bytes(request.form['pw'], 'utf-8'), bytes(user[0][0], 'utf-8'))):
            return(re_error(conn, '/error/5'))

        hashed = bcrypt.hashpw(bytes(request.form['pw2'], 'utf-8'), bcrypt.gensalt())
        
        curs.execute("update user set pw = ? where id = ?", [hashed.decode(), session['DREAMER']])
        conn.commit()
        
        return(redirect('/user'))
    else:        
        if(ban == 1):
            return(re_error(conn, '/ban'))

        if(re.search('(\.|:)', ip)):
            return(redirect('/login'))

        return(html_minify(template('index',    
            imp = ['비밀번호 변경', wiki_set(conn, 1), custom(conn), other2([0, 0])],
            data = '<form method="post"> \
                        <input placeholder="현재 비밀번호" name="pw" type="password"><hr> \
                        <input placeholder="변경할 비밀번호" name="pw2" type="password"><hr> \
                        <input placeholder="재 확인" name="pw3" type="password"><hr> \
                        <button class="btn btn-primary" type="submit">변경</button><hr> \
                        <span>주의 : 만약 HTTPS 연결이 아닌 경우 데이터가 유출될 가능성이 있습니다. 이에 대해 책임지지 않습니다.</span> \
                    </form>',
            menu = [['user', '사용자']]
        )))
                
@app.route('/check/<name>')
def user_check(name = None):
    if(admin_check(conn, 4, 'check (' + name + ')') != 1):
        return(re_error(conn, '/error/3'))

    curs.execute("select acl from user where id = ? or id = ?", [name, request.args.get('plus', 'None-Data')])
    user = curs.fetchall()
    if(user and user[0][0] != 'user'):
        if(admin_check(conn, None, None) != 1):
            return(re_error(conn, '/error/4'))
    
    if(request.args.get('plus', None)):
        if(re.search('(?:\.|:)', name)):
            if(re.search('(?:\.|:)', request.args.get('plus', None))):
                curs.execute("select name, ip, ua, today from ua_d where ip = ? or ip = ? order by today desc", [name, request.args.get('plus', None)])
            else:
                curs.execute("select name, ip, ua, today from ua_d where ip = ? or name = ? order by today desc", [name, request.args.get('plus', None)])
        else:
            if(re.search('(?:\.|:)', request.args.get('plus', None))):
                curs.execute("select name, ip, ua, today from ua_d where name = ? or ip = ? order by today desc", [name, request.args.get('plus', None)])
            else:
                curs.execute("select name, ip, ua, today from ua_d where name = ? or name = ? order by today desc", [name, request.args.get('plus', None)])
    elif(re.search('(?:\.|:)', name)):
        curs.execute("select name, ip, ua, today from ua_d where ip = ? order by today desc", [name])
    else:
        curs.execute("select name, ip, ua, today from ua_d where name = ? order by today desc", [name])
    
    record = curs.fetchall()
    if(record):
        div = '<table style="width: 100%; text-align: center;"><tbody><tr>'
        div += '<td style="width: 33.3%;">이름</td><td style="width: 33.3%;">아이피</td><td style="width: 33.3%;">언제</td></tr>'

        for data in record:
            if(data[2]):
                ua = data[2]
            else:
                ua = '<br>'

            div += '<tr><td>' + ip_pas(conn, data[0]) + '</td><td>' + ip_pas(conn, data[1]) + '</td><td>' + data[3] + '</td></tr>'
            div += '<tr><td colspan="3">' + ua + '</td></tr>'
        
        div += '</tbody></table>'
    else:
        return(re_error(conn, '/error/5'))
            
    return(html_minify(template('index',    
        imp = ['다중 검사', wiki_set(conn, 1), custom(conn), other2([0, 0])],
        data = div,
        menu = [['manager', '관리자']]
    )))
                
@app.route('/register', methods=['POST', 'GET'])
def register():
    ip = ip_check()
    ban = ban_check(conn)

    if(ban == 1):
        return(re_error(conn, '/ban'))

    if(not admin_check(conn, None, None) == 1):
        curs.execute('select data from other where name = "reg"')
        set_d = curs.fetchall()
        if(set_d and set_d[0][0] == 'on'):
            return(re_error(conn, '/ban'))
    
    if(request.method == 'POST'): 
        if(captcha_post(request.form.get('g-recaptcha-response', None), conn) == 1):
            return(re_error(conn, '/error/13'))
        else:
            captcha_post('', conn, 0)

        if(request.form['pw'] != request.form['pw2']):
            return(re_error(conn, '/error/20'))

        if(re.search('(?:[^A-Za-zㄱ-힣0-9 ])', request.form['id'])):
            return(re_error(conn, '/error/8'))

        if(len(request.form['id']) > 32):
            return(re_error(conn, '/error/7'))

        curs.execute("select id from user where id = ?", [request.form['id']])
        if(curs.fetchall()):
            return(re_error(conn, '/error/6'))

        hashed = bcrypt.hashpw(bytes(request.form['pw'], 'utf-8'), bcrypt.gensalt())
        
        curs.execute("select id from user limit 1")
        user_ex = curs.fetchall()
        if(not user_ex):
            curs.execute("insert into user (id, pw, acl, date) values (?, ?, '소유자', ?)", [request.form['id'], hashed.decode(), get_time()])
        else:
            curs.execute("insert into user (id, pw, acl, date) values (?, ?, 'user', ?)", [request.form['id'], hashed.decode(), get_time()])
        conn.commit()
        
        return(redirect('/login'))
    else:        
        contract = ''
        curs.execute('select data from other where name = "contract"')
        data = curs.fetchall()
        if(data and data[0][0] != ''):
            contract = data[0][0] + '<hr>'

        return(html_minify(template('index',    
            imp = ['회원가입', wiki_set(conn, 1), custom(conn), other2([0, 0])],
            data = '<form method="post"> \
                        ' + contract + ' \
                        <input placeholder="아이디" name="id" type="text"><hr> \
                        <input placeholder="비밀번호" name="pw" type="password"><hr> \
                        <input placeholder="다시" name="pw2" type="password"><hr> \
                        ' + captcha_get(conn) + ' \
                        <button class="btn btn-primary" type="submit">가입</button><hr> \
                        <span>주의 : 만약 HTTPS 연결이 아닌 경우 데이터가 유출될 가능성이 있습니다. 이에 대해 책임지지 않습니다.</span> \
                    </form>',
            menu = [['user', '사용자']]
        )))
            
@app.route('/logout')
def logout():
    session['Now'] = 0
    session.pop('DREAMER', None)

    return(redirect('/user'))
    
@app.route('/ban/<name>', methods=['POST', 'GET'])
def user_ban(name = None):
    curs.execute("select acl from user where id = ?", [name])
    user = curs.fetchall()
    if(user and user[0][0] != 'user'):
        if(admin_check(conn, None, None) != 1):
            return(re_error(conn, '/error/4'))

    if(request.method == 'POST'):
        if(admin_check(conn, 1, 'ban (' + name + ')') != 1):
            return(re_error(conn, '/error/3'))

        ip = ip_check()
        time = get_time()

        time_list = [request.form['month'], request.form['day'], request.form['hour'], request.form['minu']]
        num = 0
        for time_fix in time_list:
            if(not re.search("[0-9]{2}", time_fix)):
                time_list[num] = '0' + time_fix
                
            num += 1
        
        if(request.form['year'] == '09'):
            end = ''
        else:
            end = request.form['year'] + '-' + time_list[0] + '-' + time_list[1] + ' ' + time_list[2] + ':' + time_list[3] + ':00'

        curs.execute("select block from ban where block = ?", [name])
        if(curs.fetchall()):
            curs.execute("insert into rb (block, end, today, blocker, why, band) values (?, ?, ?, ?, ?, '')", [name, '해제', time, ip, ''])
            curs.execute("delete from ban where block = ?", [name])
        else:
            if(re.search("^([0-9]{1,3}\.[0-9]{1,3})$", name)):
                band_d = 'O'
            else:
                band_d = ''

            if(request.form['login_ok'] != ''):
                login = 'O'
            else:
                login = ''

            curs.execute("insert into rb (block, end, today, blocker, why, band) values (?, ?, ?, ?, ?, ?)", [name, end, time, ip, request.form['why'], band_d])
            curs.execute("insert into ban (block, end, why, band, login) values (?, ?, ?, ?, ?)", [name, end, request.form['why'], band_d, login])

        conn.commit()
        return(redirect('/ban/' + url_pas(name)))            
    else:
        if(admin_check(conn, 1, None) != 1):
            return(re_error(conn, '/error/3'))

        curs.execute("select end from ban where block = ?", [name])
        end = curs.fetchall()
        if(end):
            now = '차단 해제'
            if(end[0][0] == ''):
                data = '영구 차단<hr>'
            else:
                data = end[0][0] + ' 까지 차단<hr>'
        else:
            if(re.search("^([0-9]{1,3}\.[0-9]{1,3})$", name)):
                now = '대역 차단'
            else:
                now = '차단'

            now_time = get_time()
            m = re.search('^([0-9]{4})-([0-9]{2})-([0-9]{2}) ([0-9]{2}):([0-9]{2}):[0-9]{2}', now_time)
            g = m.groups()

            year = '<option value="09">영구</option>'
            for i in range(int(g[0]), int(g[0]) + 11):
                if(i == int(g[0])):
                    year += '<option value="' + str(i) + '" selected>' + str(i) + '</option>'
                else:
                    year += '<option value="' + str(i) + '">' + str(i) + '</option>'

            month = ''
            for i in range(1, 13):
                if(i == int(g[1])):
                    month += '<option value="' + str(i) + '" selected>' + str(i) + '</option>'
                else:
                    month += '<option value="' + str(i) + '">' + str(i) + '</option>'
                
            day = ''
            for i in range(1, 32):
                if(i == int(g[2])):
                    day += '<option value="' + str(i) + '" selected>' + str(i) + '</option>'
                else:
                    day += '<option value="' + str(i) + '">' + str(i) + '</option>'

            hour = ''
            for i in range(0, 24):
                if(i == int(g[3])):
                    hour += '<option value="' + str(i) + '" selected>' + str(i) + '</option>'
                else:
                    hour += '<option value="' + str(i) + '">' + str(i) + '</option>'

            minu = ''
            for i in range(0, 61):
                if(i == int(g[4])):
                    minu += '<option value="' + str(i) + '" selected>' + str(i) + '</option>'
                else:
                    minu += '<option value="' + str(i) + '">' + str(i) + '</option>'

            is_it = ''
            if(re.search('(\.|:)', name)):
                plus = '<input type="checkbox" name="login_ok"> 로그인 가능<hr>'
            else:
                plus = ''
            
            data = '<select name="year">' + year + '</select> 년 '
            data += '<select name="month">' + month + '</select> 월 '
            data += '<select name="day">' + day + '</select> 일 <hr>'
            data += '<select name="hour">' + hour + '</select> 시 '
            data += '<select name="minu">' + minu + '</select> 분 까지<hr>'
            data += '<input placeholder="사유" class="form-control" name="why" type="text"><br>' + plus

        return(html_minify(template('index', 
            imp = [name, wiki_set(conn, 1), custom(conn), other2([' (' + now + ')', 0])],
            data = '<form method="post">' + data + '<button class="btn btn-primary" type="submit">' + now + '</button></form>',
            menu = [['manager', '관리자']]
        )))            
                
@app.route('/acl/<path:name>', methods=['POST', 'GET'])
def acl(name = None):
    test = re.search('^사용자:(.+)$', name)

    if(request.method == 'POST'):
        if(test):
            test = test.groups()
            ip = ip_check()

            if(re.search("(\.|:)", ip)):
                return(redirect('/login'))
            elif(test[0] != ip):
                if(admin_check(conn, 5, 'acl (' + name + ')') != 1):
                    return(re_error(conn, '/error/3'))
        else:
            if(admin_check(conn, 5, 'acl (' + name + ')') != 1):
                return(re_error(conn, '/error/3'))

        if(request.form['select'] == 'admin'):
            sql = 'admin'
        elif(request.form['select'] == 'all'):
            sql = 'all'
        elif(request.form['select'] == 'user'):
            sql = 'user'
        else:
            sql = ''
            
        if(sql == ''):
            curs.execute("delete from acl where title = ?", [name])
        else:
            curs.execute("select title from acl where title = ?", [name])
            if(curs.fetchall()):
                curs.execute("update acl set dec = ? where title = ?", [sql, name])
            else:
                curs.execute("insert into acl (title, dec, dis, why) values (?, ?, '', '')", [name, sql])

        conn.commit()
            
        return(redirect('/w/' + url_pas(name)))            
    else:
        if(test):
            test = test.groups()
            ip = ip_check()

            if(re.search("(\.|:)", ip)):
                return(redirect('/login'))
            elif(test[0] != ip):
                if(admin_check(conn, 5, 'acl (' + name + ')') != 1):
                    return(re_error(conn, '/error/3'))
        else:
            if(admin_check(conn, 5, 'acl (' + name + ')') != 1):
                return(re_error(conn, '/error/3'))

        acl_list = ['', '', '']
        curs.execute("select dec from acl where title = ?", [name])
        acl_d = curs.fetchall()
        if(acl_d):
            if(test and acl_d[0][0] == 'all'):
                now = '모두'
                acl_list[0] = 'selected="selected"'
            elif(not test and acl_d[0][0] == 'admin'):
                now = '관리자'
                acl_list[0] = 'selected="selected"'
            elif(acl_d[0][0] == 'user'):
                now = '가입자'
                acl_list[1] = 'selected="selected"'
            else:
                now = '일반'
                acl_list[2] = 'selected="selected"'
        else:
            now = '일반'
            acl_list[2] = 'selected="selected"'

        if(test):
            plus = '<option value="all" ' + acl_list[0] + '>모두</option>'
        else:
            plus = '<option value="admin" ' + acl_list[0] + '>관리자</option>'
            
        return(html_minify(template('index', 
            imp = [name, wiki_set(conn, 1), custom(conn), other2([' (ACL)', 0])],
            data = '<form method="post"> \
                        <span>현재 ACL : ' + now + '</span><hr> \
                        <select name="select"> \
                            ' + plus + ' \
                            <option value="user" ' + acl_list[1] + '>가입자</option> \
                            <option value="normal" ' + acl_list[2] + '>일반</option> \
                        </select><hr> \
                        <button class="btn btn-primary" type="submit">ACL 변경</button> \
                    </form>',
            menu = [['w/' + url_pas(name), '문서'], ['manager', '관리자']]
        )))
            
@app.route('/admin/<name>', methods=['POST', 'GET'])
def user_admin(name = None):
    owner = admin_check(conn, None, None)

    curs.execute("select acl from user where id = ?", [name])
    user = curs.fetchall()
    if(not user):
        return(re_error(conn, '/error/5'))
    else:
        if(owner != 1):
            curs.execute('select name from alist where name = ? and acl = "owner"', [user[0][0]])
            if(curs.fetchall()):
                return(re_error(conn, '/error/3'))

            if(ip_check() == name):
                return(re_error(conn, '/error/3'))

    if(request.method == 'POST'):
        if(admin_check(conn, 7, 'admin (' + name + ')') != 1):
            return(re_error(conn, '/error/3'))

            curs.execute('select name from alist where name = ? and acl = "owner"', [request.form['select']])
            if(curs.fetchall()):
                return(re_error(conn, '/error/3'))

        if(request.form['select'] == 'X'):
            curs.execute("update user set acl = 'user' where id = ?", [name])
        else:
            curs.execute("update user set acl = ? where id = ?", [request.form['select'], name])
        conn.commit()
        
        return(redirect('/admin/' + url_pas(name)))            
    else:
        if(admin_check(conn, 7, None) != 1):
            return(re_error(conn, '/error/3'))            

        div = '<option value="X">X</option>'
            
        curs.execute('select distinct name from alist order by name asc')
        get_alist = curs.fetchall()
        if(get_alist):
            i = 0
            name_rem = ''
            for data in get_alist:
                if(user[0][0] == data[0]):
                    div += '<option value="' + data[0] + '" selected="selected">' + data[0] + '</option>'
                else:
                    if(owner != 1):
                        curs.execute('select name from alist where name = ? and acl = "owner"', [data[0]])
                        if(not curs.fetchall()):
                            div += '<option value="' + data[0] + '">' + data[0] + '</option>'
                    else:
                        div += '<option value="' + data[0] + '">' + data[0] + '</option>'
        
        return(html_minify(template('index', 
            imp = [name, wiki_set(conn, 1), custom(conn), other2([' (권한 부여)', 0])],
            data =  '<form method="post"> \
                        <select name="select">' + div + '</select><hr> \
                        <button class="btn btn-primary" type="submit">변경</button> \
                    </form>',
            menu = [['manager', '관리자']]
        )))
    
@app.route('/diff/<path:name>')
def diff_data(name = None):
    first = request.args.get('first', '1')
    second = request.args.get('second', '1')

    curs.execute("select data from history where id = ? and title = ?", [first, name])
    first_raw_data = curs.fetchall()
    if(first_raw_data):
        curs.execute("select data from history where id = ? and title = ?", [second, name])
        second_raw_data = curs.fetchall()
        if(second_raw_data):
            first_data = html.escape(first_raw_data[0][0])            
            second_data = html.escape(second_raw_data[0][0])
            if(first == second):
                result = '내용이 같습니다.'
            else:            
                diff_data = difflib.SequenceMatcher(None, first_data, second_data)
                result = diff(diff_data)
            
            return(html_minify(template('index', 
                imp = [name, wiki_set(conn, 1), custom(conn), other2([' (비교)', 0])],
                data = '<pre>' + result + '</pre>',
                menu = [['history/' + url_pas(name), '역사']]
            )))

    return(redirect('/history/' + url_pas(name)))
        
@app.route('/down/<path:name>')
def down(name = None):
    curs.execute("select title from data where title like ?", ['%' + name + '/%'])
    under = curs.fetchall()
    
    div = '<ul>'

    for data in under:
        div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a></li>'
        
    div += '</ul>'
    
    return(html_minify(template('index', 
        imp = [name, wiki_set(conn, 1), custom(conn), other2([' (하위)', 0])],
        data = div,
        menu = [['w/' + url_pas(name), '문서']]
    )))

@app.route('/w/<path:name>')
def read_view(name = None):
    data_none = 0
    sub = ''
    acl = ''
    div = ''
    topic = 0
    num = request.args.get('num', None)

    if(num):
        num = int(num)
    
    if(not num):
        if('View_list' in session):
            m = re.findall('([^\n]+)\n', session['View_List'])
            if(m[-1] != name):
                d = re.sub(name + '\n', '', session['View_List'])
                d += name + '\n'
                if(len(m) > 50):
                    d = re.sub('([^\n]+)\n', '', d, 1)
                session['View_List'] = d
        else:
            session['View_List'] = name + '\n'


    curs.execute("select sub from rd where title = ? order by date desc", [name])
    rd = curs.fetchall()
    for data in rd:
        curs.execute("select title from stop where title = ? and sub = ? and close = 'O'", [name, data[0]])
        if(not curs.fetchall()):
            topic = 1
            break
                
    curs.execute("select title from data where title like ?", ['%' + name + '/%'])
    if(curs.fetchall()):
        down = 1
    else:
        down = 0
        
    m = re.search("^(.*)\/(.*)$", name)
    if(m):
        uppage = m.groups()[0]
    else:
        uppage = 0
        
    if(admin_check(conn, 5, None) == 1):
        admin_memu = 1
    else:
        admin_memu = 0
        
    if(re.search("^분류:", name)):        
        curs.execute("select link from back where title = ? and type='cat' order by link asc", [name])
        back = curs.fetchall()
        if(back):
            div = '[목차(없음)]\r\n== 분류 ==\r\n'
            u_div = ''
            i = 0
            
            for data in back:       
                if(re.search('^분류:', data[0])):
                    if(u_div == ''):
                        u_div = '=== 하위 분류 ===\r\n'

                    u_div += ' * [[:' + data[0] + ']]\r\n'
                elif(re.search('^틀:', data[0])):
                    curs.execute("select data from data where title = ?", [data[0]])
                    db_data = curs.fetchall()
                    if(db_data):
                        cat_data = re.sub("\[\[(분류:(?:(?:(?!\]\]|\|).)+))(?!\|include)\]\]", "", mid_pas(db_data[0][0], 0, 1, 0, 0)[0])
                        if(re.search('\[\[' + name + '|include]]', cat_data)):
                            div += ' * [[' + data[0] + ']]\r\n * [[wiki:xref/' + url_pas(data[0]) + '|' + data[0] + ']] (역링크)\r\n'
                        else:
                            div += ' * [[' + data[0] + ']]\r\n'
                else:
                    div += ' * [[' + data[0] + ']]\r\n'

            div += u_div

    if(num):
        curs.execute("select title from history where title = ? and id = ? and hide = 'O'", [name, str(num)])
        hid = curs.fetchall()
        if(hid and admin_check(conn, 6, None) != 1):
            return(redirect('/history/' + url_pas(name)))

        curs.execute("select title, data from history where title = ? and id = ?", [name, str(num)])
    else:
        curs.execute("select title, data from data where title = ?", [name])
    data = curs.fetchall()
    if(data):
        elsedata = data[0][1]
    else:
        data_none = 1
        response.status = 404
        elsedata = ''

    if(not num):
        curs.execute("select dec from acl where title = ?", [name])
        acl_d = curs.fetchall()
        if(acl_d):
            if(acl_d[0][0] == 'admin'):
                acl = ' (관리자)'
            elif(acl_d[0][0] == 'user'):
                acl = ' (가입자)'
        else:
            curs.execute('select data from other where name = "edit"')
            set_data = curs.fetchall()
            if(set_data):
                if(set_data[0][0] == 'admin'):
                    acl = ' (관리자)'
                elif(set_data[0][0] == 'user'):
                    acl = ' (가입자)'

    m = re.search("^사용자:([^/]*)", name)
    if(m):
        g = m.groups()
        
        curs.execute("select acl from user where id = ?", [g[0]])
        test = curs.fetchall()
        if(test and test[0][0] != 'user'):
            acl = ' (관리자)'
        else:
            curs.execute("select block from ban where block = ?", [g[0]])
            if(curs.fetchall()):
                sub = ' (차단)'
            else:
                acl = ''

        curs.execute("select dec from acl where title = ?", [name])
        data = curs.fetchall()
        if(data):
            if(data[0][0] == 'all'):
                acl += ' (모두)'
            elif(data[0][0] == 'user'):
                acl += ' (가입자)'
            
    if(request.args.get('froms', None)):
        elsedata = re.sub("^#(?:redirect|넘겨주기) (?P<in>[^\n]*)", " * [[\g<in>]] 문서로 넘겨주기", elsedata)
            
    enddata = namumark(conn, name, elsedata, 0, 0, 1)

    if(data_none == 1):
        menu = [['edit/' + url_pas(name), '생성'], ['topic/' + url_pas(name), topic], ['history/' + url_pas(name), '역사'], ['move/' + url_pas(name), '이동'], ['xref/' + url_pas(name), '역링크']]
    else:
        menu = [['edit/' + url_pas(name), '수정'], ['topic/' + url_pas(name), topic], ['history/' + url_pas(name), '역사'], ['delete/' + url_pas(name), '삭제'], ['move/' + url_pas(name), '이동'], ['raw/' + url_pas(name), '원본'], ['xref/' + url_pas(name), '역링크']]
        
    if(admin_memu == 1):
        menu += [['acl/' + url_pas(name), 'ACL']]

    if(request.args.get('froms', None)):
        menu += [['w/' + url_pas(name), '넘기기']]
        enddata = '<ul id="redirect"><li><a href="/w/' + url_pas(request.args.get('froms', None)) + '?froms=' + url_pas(name) + '">' + request.args.get('froms', None) + '</a>에서 넘어 왔습니다.</li></ul><br>' + enddata

    if(uppage != 0):
        menu += [['w/' + url_pas(uppage), '상위']]

    if(down):
        menu += [['down/' + url_pas(name), '하위']]

    if(num):
        menu = [['history/' + url_pas(name), '역사']]
        sub = ' (' + str(num) + '판)'
        acl = ''
        r_date = 0
    else:
        curs.execute("select date from history where title = ? order by date desc limit 1", [name])
        date = curs.fetchall()
        if(date):
            r_date = date[0][0]
        else:
            r_date = 0

    if(div != '' and enddata != ''):
        div = enddata + '<br>' + namumark(conn, name, div, 0, 0, 0)
    else:
        div = enddata + namumark(conn, name, div, 0, 0, 0)

    return(html_minify(template('index', 
        imp = [name, wiki_set(conn, 1), custom(conn), other2([sub + acl, r_date])],
        data = div,
        menu = menu
    )))

@app.route('/topic_record/<name>')
def user_topic_list(name = None):
    num = int(request.args.get('num', 1))
    if(num * 50 > 0):
        sql_num = num * 50 - 50
    else:
        sql_num = 0
    
    one_admin = admin_check(conn, 1, None)
    div = '<table style="width: 100%; text-align: center;"><tbody><tr>'
    div += '<td style="width: 33.3%;">토론명</td><td style="width: 33.3%;">작성자</td><td style="width: 33.3%;">시간</td></tr>'
    
    curs.execute("select title, id, sub, ip, date from topic where ip = ? order by date desc limit ?, '50'", [name, str(sql_num)])
    data_list = curs.fetchall()
    for data in data_list:
        title = html.escape(data[0])
        sub = html.escape(data[2])
            
        if(one_admin == 1):
            curs.execute("select * from ban where block = ?", [data[3]])
            if(curs.fetchall()):
                ban = ' <a href="/ban/' + url_pas(data[3]) + '">(해제)</a>'
            else:
                ban = ' <a href="/ban/' + url_pas(data[3]) + '">(차단)</a>'
        else:
            ban = ''
            
        ip = ip_pas(conn, data[3])
            
        div += '<tr><td><a href="/topic/' + url_pas(data[0]) + '/sub/' + url_pas(data[2]) + '#' + data[1] + '">' + title + '#' + data[1] + '</a> (' + sub + ')</td>'
        div += '<td>' + ip + ban +  '</td><td>' + data[4] + '</td></tr>'

    div += '</tbody></table>'
    div += next_fix('/topic_record/' + url_pas(name) + '?num=', num, data_list)
                
    curs.execute("select end, why from ban where block = ?", [name])
    ban_it = curs.fetchall()
    if(ban_it):
        sub = ' (차단)'
    else:
        sub = 0 
    
    return(html_minify(template('index', 
        imp = ['토론 기록', wiki_set(conn, 1), custom(conn), other2([sub, 0])],
        data = div,
        menu = [['other', '기타'], ['user', '사용자'], ['count/' + url_pas(name), '횟수'], ['record/' + url_pas(name), '편집 기록']]
    )))

@app.route('/recent_changes')
@app.route('/<regex("record"):tool>/<name>')
@app.route('/<regex("history"):tool>/<path:name>', methods=['POST', 'GET'])
def recent_changes(name = None, tool = 'record'):
    if(request.method == 'POST'):
        return(redirect('/diff/' + url_pas(name) + '?first=' + request.form['b'] + '&second=' + request.form['a']))
    else:
        one_admin = admin_check(conn, 1, None)
        six_admin = admin_check(conn, 6, None)
        ban = ''
        select = ''
        what = request.args.get('what', 'all')
        div = '<table style="width: 100%; text-align: center;"><tbody><tr>'
        
        if(name):
            num = int(request.args.get('num', 1))
            if(num * 50 > 0):
                sql_num = num * 50 - 50
            else:
                sql_num = 0      

            if(tool == 'history'):
                div += '<td style="width: 33.3%;">판</td><td style="width: 33.3%;">편집자</td><td style="width: 33.3%;">시간</td></tr>'

                curs.execute("select id, title, date, ip, send, leng from history where title = ? order by id + 0 desc limit ?, '50'", [name, str(sql_num)])
            else:
                div += '<td style="width: 33.3%;">문서명</td><td style="width: 33.3%;">편집자</td><td style="width: 33.3%;">시간</td></tr>'

                if(what == 'all'):
                    div = '<a href="/record/' + url_pas(name) + '?what=revert">(되돌리기)</a><hr>' + div
                    div = '<a href="/record/' + url_pas(name) + '?what=move">(이동)</a> ' + div
                    div = '<a href="/record/' + url_pas(name) + '?what=delete">(삭제)</a> ' + div
                
                    curs.execute("select id, title, date, ip, send, leng from history where ip = ? order by date desc limit ?, '50'", [name, str(sql_num)])
                else:
                    if(what == 'delete'):
                        sql = '%(삭제)'
                    elif(what == 'move'):
                        sql = '%이동)'
                    elif(what == 'revert'):
                        sql = '%판)'
                    else:
                        return(redirect('/'))

                    curs.execute("select id, title, date, ip, send, leng from history where ip = ? and send like ? order by date desc limit ?, '50'", [name, sql, str(sql_num)])
        else:
            div += '<td style="width: 33.3%;">문서명</td><td style="width: 33.3%;">편집자</td><td style="width: 33.3%;">시간</td></tr>'

            if(what == 'all'):
                div = '<a href="/recent_changes?what=revert">(되돌리기)</a><hr>' + div
                div = '<a href="/recent_changes?what=move">(이동)</a> ' + div
                div = '<a href="/recent_changes?what=delete">(삭제)</a> ' + div

                curs.execute("select id, title, date, ip, send, leng from history order by date desc limit 50")
            else:
                if(what == 'delete'):
                    sql = '%(삭제)'
                elif(what == 'move'):
                    sql = '%이동)'
                elif(what == 'revert'):
                    sql = '%판)'
                else:
                    return(redirect('/'))

                curs.execute("select id, title, date, ip, send, leng from history where send like ? order by date desc limit 50", [sql])

        data_list = curs.fetchall()
        for data in data_list:    
            select += '<option value="' + data[0] + '">' + data[0] + '</option>'     
            send = '<br>'
            if(data[4]):
                if(not re.search("^(?: *)$", data[4])):
                    send = data[4]
            
            if(re.search("\+", data[5])):
                leng = '<span style="color:green;">(' + data[5] + ')</span>'
            elif(re.search("\-", data[5])):
                leng = '<span style="color:red;">(' + data[5] + ')</span>'
            else:
                leng = '<span style="color:gray;">(' + data[5] + ')</span>'
                
            if(one_admin == 1):
                curs.execute("select * from ban where block = ?", [data[3]])
                if(curs.fetchall()):
                    ban = ' <a href="/ban/' + url_pas(data[3]) + '">(해제)</a>'
                else:
                    ban = ' <a href="/ban/' + url_pas(data[3]) + '">(차단)</a>'            
                
            ip = ip_pas(conn, data[3])
                    
            if((int(data[0]) - 1) == 0):
                revert = ''
            else:
                revert = '<a href="/diff/' + url_pas(data[1]) + '?first=' + str(int(data[0]) - 1) + '&second=' + data[0] + '">(비교)</a> <a href="/revert/' + url_pas(data[1]) + '?num=' + str(int(data[0]) - 1) + '">(되돌리기)</a>'
            
            style = ['', '']
            date = data[2]

            curs.execute("select title from history where title = ? and id = ? and hide = 'O'", [data[1], data[0]])
            hide = curs.fetchall()
            if(six_admin == 1):
                if(hide):                            
                    hidden = ' <a href="/hidden/' + url_pas(data[1]) + '?num=' + data[0] + '">(공개)'
                    
                    style[0] = 'background: gainsboro;'
                    style[1] = 'background: gainsboro;'

                    if(send == '<br>'):
                        send = '(숨김)'
                    else:
                        send += ' (숨김)'
                else:
                    hidden = ' <a href="/hidden/' + url_pas(data[1]) + '?num=' + data[0] + '">(숨김)'
            elif(not hide):
                hidden = ''
            else:
                ip = ''
                hidden = ''
                ban = ''
                date = ''
                send = '(숨김)'

                style[0] = 'display: none;'
                style[1] = 'background: gainsboro;'

            if(tool == 'history'):
                title = '<a href="/w/' + url_pas(name) + '?num=' + data[0] + '">' + data[0] + '판</a> <a href="/raw/' + url_pas(name) + '?num=' + data[0] + '">(원본)</a> '
            else:
                title = '<a href="/w/' + url_pas(data[1]) + '">' + html.escape(data[1]) + '</a> <a href="/history/' + url_pas(data[1]) + '">(' + data[0] + '판)</a> '
                    
            div += '<tr style="' + style[0] + '"><td>' + title + revert + ' ' + leng + '</td>'
            div += '<td>' + ip + ban + hidden + '</td><td>' + date + '</td></tr><tr style="' + style[1] + '"><td colspan="3">' + send + '</td></tr>'

        div += '</tbody></table>'
        sub = ''

        if(name):
            if(tool == 'history'):
                div = '<form method="post"><select name="a">' + select + '</select> <select name="b">' + select + '</select> <button class="btn btn-primary" type="submit">비교</button></form><hr>' + div
                title = name
                sub += ' (역사)'
                menu = [['w/' + url_pas(name), '문서'], ['move_data/' + url_pas(name), '이동 기록']]
                
                div += next_fix('/history/' + url_pas(name) + '?num=', num, data_list)
            else:
                curs.execute("select end, why from ban where block = ?", [name])
                ban_it = curs.fetchall()
                if(ban_it):
                    sub += ' (차단)'

                title = '편집 기록'
                menu = [['other', '기타'], ['user', '사용자'], ['count/' + url_pas(name), '횟수'], ['topic_record/' + url_pas(name), '토론 기록']]
                div += next_fix('/record/' + url_pas(name) + '/' + url_pas(what) + '?num=', num, data_list)

                if(what != 'all'):
                    menu += [['record/' + url_pas(name), '일반']]
        else:
            menu = 0
            title = '최근 변경내역'

            if(what != 'all'):
                menu = [['recent_changes', '일반']]
                
        if(what == 'delete'):
            sub += ' (삭제)'
        elif(what == 'move'):
            sub += ' (이동)'
        elif(what == 'revert'):
            sub += ' (되돌리기)'
        
        if(sub == ''):
            sub = 0
                
        return(html_minify(template('index', 
            imp = [title, wiki_set(conn, 1), custom(conn), other2([sub, 0])],
            data = div,
            menu = menu
        )))
    
@app.route('/upload', methods=['GET', 'POST'])
def upload():
    if(ban_check(conn) == 1):
        return(re_error(conn, '/ban'))
    
    if(request.method == 'POST'):
        if(captcha_post(request.form.get('g-recaptcha-response', None), conn) == 1):
            return(re_error(conn, '/error/13'))
        else:
            captcha_post('', conn, 0)

        data = request.files['f_data']
        if(not data):
            return(re_error(conn, '/error/9'))

        if(int(wiki_set(conn, 3)) * 1024 * 1024 < request.content_length):
            return(re_error(conn, '/error/17'))
        
        value = os.path.splitext(data.filename)[1]

        if(not value in ['.jpeg', '.jpg', '.gif', '.png', '.webp', '.JPEG', '.JPG', '.GIF', '.PNG', '.WEBP']):
            return(re_error(conn, '/error/14'))
    
        if(request.form['f_name']):
            name = request.form['f_name'] + value
        else:
            name = data.filename
        
        piece = os.path.splitext(name)

        if(re.search('[^ㄱ-힣0-9a-zA-Z_\- ]', piece[0])):
            return(re_error(conn, '/error/22'))

        e_data = sha224(piece[0]) + piece[1]

        curs.execute("select title from data where title = ?", ['파일:' + name])
        if(curs.fetchall()):
            return(re_error(conn, '/error/16'))
            
        ip = ip_check()
        if(request.form['f_lice']):
            lice = request.form['f_lice']
        else:
            if(re.search('(?:\.|:)', ip)):
                lice = ip + ' 올림'
            else:
                lice = '[[사용자:' + ip + ']] 올림'
            
        if(os.path.exists(os.path.join('image', e_data))):
            os.remove(os.path.join('image', e_data))
            data.save(os.path.join('image', e_data))
        else:
            data.save(os.path.join('image', e_data))
            
        curs.execute("select title from data where title = ?", ['파일:' + name])
        exist = curs.fetchall()
        if(exist): 
            curs.execute("delete from data where title = ?", ['파일:' + name])
        
        curs.execute("insert into data (title, data) values (?, ?)", ['파일:' + name, '[[파일:' + name + ']][br][br]{{{[[파일:' + name + ']]}}}[br][br]' + lice])
        curs.execute("insert into acl (title, dec, dis, why) values (?, 'admin', '', '')", ['파일:' + name])
        history_plus(conn, '파일:' + name, '[[파일:' + name + ']][br][br]{{{[[파일:' + name + ']]}}}[br][br]' + lice, get_time(), ip, '(파일 올림)', '0')
        conn.commit()
        
        return(redirect('/w/파일:' + name))            
    else:
        return(html_minify(template('index', 
            imp = ['파일 올리기', wiki_set(conn, 1), custom(conn), other2([0, 0])],
            data =  '<form method="post" enctype="multipart/form-data" accept-charset="utf8"> \
                        <input type="file" name="f_data"><hr> \
                        <input placeholder="파일 이름" name="f_name" type="text"><hr> \
                        <input placeholder="라이선스" name="f_lice" type="text"><hr> \
                        ' + captcha_get(conn) + ' \
                        <button class="btn btn-primary" type="submit">저장</button> \
                    </form>',
            menu = [['other', '기타']]
        )))  
        
@app.route('/user')
def user_info():
    ip = ip_check()
    
    curs.execute("select acl from user where id = ?", [ip])
    data = curs.fetchall()
    if(ban_check(conn) == 0):
        if(data):
            if(data[0][0] != 'user'):
                acl = data[0][0]
            else:
                acl = '가입자'
        else:
            acl = '일반'
    else:
        acl = '차단'
        
        curs.execute("select block from ban where block = ? and login = 'O'", [ip])
        if(curs.fetchall()):
            acl += ' (로그인 가능)'
            
    if(not re.search('(\.|:)', ip)):
        ip_user = '[[사용자:' + ip + '|' + ip + ']]'
    else:
        ip_user = ip

    if(custom(conn)[2] != 0):
        plus = ' * [[wiki:logout|로그아웃]]\r\n * [[wiki:change|비밀번호 변경]]'
    else:
        plus = ' * [[wiki:login|로그인]]'

    return(html_minify(template('index', 
        imp = ['사용자 메뉴', wiki_set(conn, 1), custom(conn), other2([0, 0])],
        data =  namumark(conn, '',  '[목차(없음)]\r\n' + \
                                    '== 상태 ==\r\n' + \
                                    ' * ' + ip_user + '\r\n'
                                    ' * 권한 상태 : ' + acl + '\r\n' + \
                                    '== 로그인 ==\r\n' + \
                                    plus + '\r\n' + \
                                    ' * [[wiki:register|회원가입]]\r\n' + \
                                    '== 사용자 기능 ==\r\n' + \
                                    ' * [[wiki:acl/사용자:' + url_pas(ip) + '|사용자 문서 ACL]]\r\n' + \
                                    ' * [[wiki:custom_head|사용자 HEAD]]\r\n' + \
                                    '== 기타 ==\r\n' + \
                                    ' * [[wiki:alarm|알림]]\r\n' + \
                                    ' * [[wiki:view_log|지나온 문서]]\r\n' + \
                                    ' * [[wiki:record/' + url_pas(ip) + '|편집 기록]]\r\n' + \
                                    ' * [[wiki:topic_record/' + url_pas(ip) + '|토론 기록]]\r\n' + \
                                    ' * [[wiki:watch_list|주시 문서]]\r\n' + \
                                    ' * [[wiki:count|활동 횟수]]\r\n', 0, 0, 0),
        menu = 0
    )))

@app.route('/watch_list')
def watch_list():
    div = '한도 : 10개<hr>'
    ip = ip_check()
    
    if(re.search('\.|:', ip)):
        return(redirect('/login'))

    curs.execute("select title from scan where user = ?", [ip])
    data = curs.fetchall()
    for data_list in data:
        div += '<li><a href="/w/' + url_pas(data_list[0]) + '">' + data_list[0] + '</a> <a href="/watch_list/' + url_pas(data_list[0]) + '">(삭제)</a></li>'

    if(data):
        div = '<ul>' + div + '</ul><hr>'

    div += '<a href="/manager/13">(추가)</a>'

    return(html_minify(template('index', 
        imp = ['편집 필터 목록', wiki_set(conn, 1), custom(conn), other2([0, 0])],
        data = div,
        menu = [['manager', '관리자']]
    )))

@app.route('/watch_list/<name>')
def watch_list_name(name = None):
    ip = ip_check()
    if(re.search('\.|:', ip)):
        return(redirect('/login'))

    curs.execute("select count(title) from scan where user = ?", [ip])
    count = curs.fetchall()
    if(count):
        if(count[0][0] > 9):
            return(redirect('/watch_list'))

    curs.execute("select title from scan where user = ? and title = ?", [ip, name])
    if(curs.fetchall()):
        curs.execute("delete from scan where user = ? and title = ?", [ip, name])
    else:
        curs.execute("insert into scan (user, title) values (?, ?)", [ip, name])
    conn.commit()

    return(redirect('/watch_list'))

@app.route('/view_log')
def view_log():
    data = '<ul>'
    if(session['View_List']):
        data += '<li>최근 50개</li><hr>'
        m = re.findall('([^\n]+)\n', session['View_List'])
        for d in m:
            data += '<li><a href="/w/' + url_pas(d) + '">' + d + '</a></li>'
    else:
        data += '<li>기록 없음</li>'
    data += '</ul>'

    return(html_minify(template('index', 
        imp = ['지나온 문서', wiki_set(conn, 1), custom(conn), other2([0, 0])],
        data = data,
        menu = [['user', '사용자']]
    )))

@app.route('/custom_head', methods=['GET', 'POST'])
def custom_head_view():
    ip = ip_check()
    if(request.method == 'POST'):
        if(not re.search('(\.|:)', ip)):
            curs.execute("select user from custom where user = ?", [ip + ' (head)'])
            if(curs.fetchall()):
                curs.execute("update custom set css = ? where user = ?", [request.form['content'], ip + ' (head)'])
            else:
                curs.execute("insert into custom (user, css) values (?, ?)", [ip + ' (head)', request.form['content']])
            conn.commit()

        session['MyMaiToNight'] = request.form['content']

        return(redirect('/user'))
    else:
        if(not re.search('(\.|:)', ip)):
            start = ''
            curs.execute("select css from custom where user = ?", [ip + ' (head)'])
            head_data = curs.fetchall()
            if(head_data):
                data = head_data[0][0]
            else:
                data = ''
        else:
            start = '<span>비 로그인의 경우에는 로그인하거나 브라우저 닫으면 날아갑니다.</span><hr>'
            try:
                data = session['MyMaiToNight']
            except:
                data = ''

        start += '<span>&lt;style&gt;CSS&lt;/style&gt;<br>&lt;script&gt;JS&lt;/script&gt;</span><hr>'

        return(html_minify(template('index', 
            imp = ['사용자 HEAD', wiki_set(conn, 1), custom(conn), other2([0, 0])],
            data =  start + ' \
                    <form method="post"> \
                        <textarea rows="25" cols="100" name="content">' + data + '</textarea><hr> \
                        <button class="btn btn-primary" type="submit">저장</button> \
                    </form>',
            menu = [['user', '사용자']]
        )))

@app.route('/count')
@app.route('/count/<name>')
def count_edit(name = None):
    if(name == None):
        that = ip_check()
    else:
        that = name

    curs.execute("select count(title) from history where ip = ?", [that])
    count = curs.fetchall()
    if(count):
        data = count[0][0]
    else:
        data = 0

    curs.execute("select count(title) from topic where ip = ?", [that])
    count = curs.fetchall()
    if(count):
        t_data = count[0][0]
    else:
        t_data = 0

    return(html_minify(template('index', 
        imp = ['활동 횟수', wiki_set(conn, 1), custom(conn), other2([0, 0])],
        data = namumark(conn, "", "[목차(없음)]\r\n== " + that + " ==\r\n||<:> 편집 횟수 ||<:> " + str(data) + "||\r\n||<:> 토론 횟수 ||<:> " + str(t_data) + "||", 0, 0, 0),
        menu = [['user', '사용자'], ['record/' + url_pas(that), '편집 기록'], ['topic_record/' + url_pas(that), '토론 기록']]
    )))
        
@app.route('/random')
def random():
    curs.execute("select title from data order by random() limit 1")
    d = curs.fetchall()
    if(d):
        return(redirect('/w/' + url_pas(d[0][0])))
    else:
        return(redirect('/'))
    
@app.route('/views/<path:name>')
def views(name = None):
    if(re.search('\/', name)):
        m = re.search('^(.*)\/(.*)$', name)
        if(m):
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
    if(m):
        g = m.groups()
    else:
        g = ['']

    if(g == 'css'):
        return(css_minify(send_from_directory('./views' + plus, rename)))   
    elif(g == 'js'):
        return(js_minify(send_from_directory('./views' + plus, rename)))
    elif(g == 'html'):
        return(html_minify(send_from_directory('./views' + plus, rename)))   
    else:
        return(send_from_directory('./views' + plus, rename))

@app.route('/robots.txt')
def robots():
    return(send_from_directory('./', request.path[1:]))

@app.errorhandler(404)
def error_404(e):
    return('<!-- 나니카가 하지마룻테 코토와 오와리니 츠나가루다난테 캉가에테모 미나캇타. 이야, 캉카에타쿠나캇탄다... 아마오토 마도오 타타쿠 소라카라 와타시노 요-나 카나시미 훗테루 토메도나쿠 이마오 누라시테 오모이데 난테 이라나이노 코코로가 쿠루시쿠나루 다케다토 No more! September Rain No more! September Rain 이츠닷테 아나타와 미짓카닷타 와자와자 키모치오 타시카메룻테 코토모 히츠요-쟈나쿠테 시젠니 나카라요쿠 나레타카라 안신시테타노 카모시레나이네 도-시테? 나미니 토이카케루케도 나츠노 하지마리가 츠레테키타 오모이 나츠가 오와루토키 키에챠우모노닷타 난테 시라나쿠테 토키메이테타 아츠이 키세츠 우미베노 소라가 히캇테 토츠젠 쿠모가 나가레 오츠부노 아메 와타시노 나카노 나미다미타이 콘나니 타노시이 나츠가 즛토 츠즈이테쿳테 신지테타요 But now... September Rain But now... September Rain -->' + redirect('/w/' + url_pas(wiki_set(conn, 2))))

if(__name__=="__main__"):
    app.secret_key = rep_key
    http_server = HTTPServer(WSGIContainer(app))
    http_server.listen(rep_port)
    IOLoop.instance().start()
