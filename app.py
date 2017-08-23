from bottle import route, run, template, error, request, static_file, app, BaseRequest
from bottle.ext import beaker
import bcrypt
import os
import difflib
import hashlib
import json
import sqlite3

try:
    json_data = open('set.json').read()
    set_data = json.loads(json_data)
except:
    new_json = []

    print('DB 이름 : ', end = '')
    new_json += [input()]
    
    print('위키 이름 : ', end = '')
    new_json += [input()]

    print('대문 : ', end = '')
    new_json += [input()]

    print('라이선스 : ', end = '')
    new_json += [input()]

    print('업로드 제한 (MB) : ', end = '')
    new_json += [input()]

    print('위키 포트 : ', end = '')
    new_json += [input()]

    with open("set.json", "w") as f:
        f.write('{ \
                    "db" : "' + new_json[0] + '", \
                    "name" : "' + new_json[1] + '", \
                    "frontpage" : "' + new_json[2] + '", \
                    "license" : "' + new_json[3] + '", \
                    "upload" : "' + new_json[4] + '", \
                    "port" : "' + new_json[5] + '" \
                }')
    
    json_data = open('set.json').read()
    set_data = json.loads(json_data)

conn = sqlite3.connect(set_data['db'] + '.db')
curs = conn.cursor()

session_opts = {
    'session.type': 'file',
    'session.data_dir': './app_session/',
    'session.auto': True
}

app = beaker.middleware.SessionMiddleware(app(), session_opts)

BaseRequest.MEMFILE_MAX = 1024 * 1024

def redirect(data):
    return('<meta http-equiv="refresh" content="0;url=' + data + '" />')
    
from func import *

r_ver = '2.2.1'

try:
    curs.execute('select data from other where name = "version"')
    version = curs.fetchall()
    if(version):
        t_ver = re.sub('\.', '', version[0]['data'])
        r_t_ver = re.sub('\.', '', r_ver)
        if(int(t_ver) <= int(r_t_ver)):
            curs.execute("update other set data = ? where name = 'version'", [r_ver])    
    
    conn.commit()
except:
    pass

@route('/setup', method=['GET', 'POST'])
def setup():
    curs.execute("create table data(title text, data text, acl text)")
    curs.execute("create table history(id text, title text, data text, date text, ip text, send text, leng text)")
    curs.execute("create table rd(title text, sub text, date text)")
    curs.execute("create table user(id text, pw text, acl text)")
    curs.execute("create table ban(block text, end text, why text, band text)")
    curs.execute("create table topic(id text, title text, sub text, data text, date text, ip text, block text, top text)")
    curs.execute("create table stop(title text, sub text, close text)")
    curs.execute("create table rb(block text, end text, today text, blocker text, why text)")
    curs.execute("create table login(user text, ip text, today text)")
    curs.execute("create table back(title text, link text, type text)")
    curs.execute("create table cat(title text, cat text)")
    curs.execute("create table hidhi(title text, re text)")
    curs.execute("create table agreedis(title text, sub text)")
    curs.execute("create table custom(user text, css text)")
    curs.execute("create table other(name text, data text)")
    curs.execute("create table alist(name text, acl text)")

    curs.execute("insert into alist (name, acl) values ('owner', 'owner')")
    curs.execute("insert into other (name, data) values ('version', ?)", [r_ver])
    conn.commit()

    return(redirect('/'))

@route('/not_close_topic')
def not_close_topic():
    div = ''
    i = 1

    curs.execute('select title, sub from rd order by date desc')
    n_list = curs.fetchall()
    for data in n_list:
        curs.execute('select * from stop where title = ? and sub = ? and close = "O"', [data[0], data[1]])
        is_close = curs.fetchall()
        if(not is_close):
            div += '<li>' + str(i) + '. <a href="/topic/' + url_pas(data[0]) + '/sub/' + url_pas(data[1]) + '">' + data[0] + ' (' + data[1] + ')</a></li>'
            i += 1

    return(
        template('other',
            custom = custom_css_user(),
            license = set_data['license'],
            login = login_check(),
            logo = set_data['name'],
            data = div,
            title = '안 닫힌 토론 목록'
        )
    )

@route('/image/<name:path>')
def static(name = None):
    if(os.path.exists(os.path.join('image', name))):
        return(static_file(name, root = 'image'))
    else:
        return(redirect('/'))

@route('/acl_list')
def acl_list():
    div = ''
    i = 0

    curs.execute("select title, acl from data where acl = 'admin' or acl = 'user' order by acl desc")
    list_data = curs.fetchall()
    for data in list_data:
        if(data[1] == 'admin'):
            acl = '관리자'
        else:
            acl = '로그인'

        div += '<li>' + str(i + 1) + '. <a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a> (' + acl + ')</li>'
        
        i += 1
    
    return(
        template('other', 
            custom = custom_css_user(), 
            license = set_data['license'], 
            login = login_check(), 
            logo = set_data['name'], 
            data = div, 
            title = 'ACL 문서 목록'
        )
    )
    
@route('/list_acl')
def list_acl():
    div = ''
    i = 1

    curs.execute("select name, acl from alist order by name desc")
    list_data = curs.fetchall()
    for data in list_data:
        if(data[1] == 'ban'):
            acl = '차단'
        elif(data[1] == 'mdel'):
            acl = '많은 문서 삭제'
        elif(data[1] == 'toron'):
            acl = '토론 관리'
        elif(data[1] == 'check'):
            acl = '사용자 검사'
        elif(data[1] == 'acl'):
            acl = '문서 ACL'
        elif(data[1] == 'hidel'):
            acl = '역사 숨김'
        elif(data[1] == 'owner'):
            acl = '소유자'
            
        div += '<li>' + str(i + 1) + '. <a href="/admin_plus/' + url_pas(data[0]) + '">' + data[0] + '</a> (' + acl + ')</li>'
        
        i += 1
    else:        
        div +=  '<br> \
                <a href="/manager/8">(생성)</a>'

    return(
        template('other', 
            custom = custom_css_user(), 
            license = set_data['license'], 
            login = login_check(), 
            logo = set_data['name'], 
            data = re.sub('^<br>', '', div), 
            title = 'ACL 목록'
        )
    )

@route('/admin_plus/<name:path>', method=['POST', 'GET'])
def admin_plus(name = None):
    if(admin_check(None) == 1):
        if(request.method == 'POST'):
            curs.execute("delete from alist where name = ?", [name])
            
            if(request.forms.ban):
                curs.execute("insert into alist (name, acl) values (?, 'ban')", [name])

            if(request.forms.mdel):
                curs.execute("insert into alist (name, acl) values (?, 'mdel')", [name])   

            if(request.forms.toron):
                curs.execute("insert into alist (name, acl) values (?, 'toron')", [name])
                
            if(request.forms.check):
                curs.execute("insert into alist (name, acl) values (?, 'check')", [name])

            if(request.forms.acl):
                curs.execute("insert into alist (name, acl) values (?, 'acl')", [name])

            if(request.forms.hidel):
                curs.execute("insert into alist (name, acl) values (?, 'hidel')", [name])

            if(request.forms.owner):
                curs.execute("insert into alist (name, acl) values (?, 'owner')", [name])
                
            conn.commit()
            
            return(redirect('/admin_plus/admin'))
        else:
            curs.execute('select acl from alist where name = ?', [name])
            test = curs.fetchall()
            
            data = ''
            exist_list = ['', '', '', '', '', '', '', '', '']

            for go in test:
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
                elif(go[0] == 'owner'):
                    exist_list[7] = 'checked="checked"'

            data += '<li><input type="checkbox" name="ban" ' + exist_list[0] + '> 차단</li>'
            data += '<li><input type="checkbox" name="mdel" ' + exist_list[1] + '> 많은 문서 삭제</li>'
            data += '<li><input type="checkbox" name="toron" ' + exist_list[2] + '> 토론 관리</li>'
            data += '<li><input type="checkbox" name="check" ' + exist_list[3] + '> 사용자 검사</li>'
            data += '<li><input type="checkbox" name="acl" ' + exist_list[4] + '> 문서 ACL</li>'
            data += '<li><input type="checkbox" name="hidel" ' + exist_list[5] + '> 역사 숨김</li>'
            data += '<li><input type="checkbox" name="owner" ' + exist_list[7] + '> 소유자</li>'
                  
            return(
                template('other', 
                    custom = custom_css_user(), 
                    license = set_data['license'], 
                    login = login_check(), 
                    title = '관리 그룹 추가', 
                    logo = set_data['name'], 
                    data = '<form method="post">' \
                                + data + \
                                '<div class="form-actions"> \
                                    <button class="btn btn-primary" type="submit">저장</button> \
                                </div> \
                            </form>'
                )
            )
    else:
        return(redirect('/error/3'))
        
@route('/admin_list')
def admin_list():
    i = 1
    div = ''
    
    curs.execute("select id, acl from user where not acl = 'user'")
    user_data = curs.fetchall()

    for data in user_data:
        name = ip_pas(data[0], 2) + ' (' + data[1] + ')'

        div += '<li>' + str(i) + '. ' + name + '</li>'
        
        i += 1
                
    return(
        template('other', 
            custom = custom_css_user(), 
            license = set_data['license'], 
            login = login_check(), 
            logo = set_data['name'], 
            data = div, 
            title = '관리자 목록'
        )
    )
        
@route('/record/<name:path>')
@route('/record/<name:path>/n/<num:int>')
@route('/recentchanges')
def recentchanges(name = None, num = 1):
    ydmin = admin_check(1)
    zdmin = admin_check(6)
    ban = ''
    send = '<br>'
    div =  '<table style="width: 100%; text-align: center;"> \
                <tbody> \
                    <tr> \
                        <td style="width: 33.3%;">문서명</td> \
                        <td style="width: 33.3%;">기여자</td> \
                        <td style="width: 33.3%;">시간</td> \
                    </tr>'
    
    if(name):
        if(num * 50 <= 0):
            v = 50
        else:
            v = num * 50
            
        i = v - 50

        curs.execute("select id, title, date, ip, send, leng from history where ip = ? order by date desc limit ?, ?", [name, str(i), str(v)])
    else:
        curs.execute("select id, title, date, ip, send, leng from history order by date desc limit 50")

    rows = curs.fetchall()

    for data in rows:         
        send = ''
        if(data[4]):
            if(not re.search("^(?: *)$", data[4])):
                send = data[4]
    
        title = re.sub('<', '&lt;', data[1])
        title = re.sub('>', '&gt;', title)
        
        if(re.search("\+", data[5])):
            leng = '<span style="color:green;">' + data[5] + '</span>'
        elif(re.search("\-", data[5])):
            leng = '<span style="color:red;">' + data[5] + '</span>'
        else:
            leng = '<span style="color:gray;">' + data[5] + '</span>'
            
        if(ydmin == 1):
            curs.execute("select * from ban where block = ?", [data[3]])
            row = curs.fetchall()
            if(row):
                ban = ' <a href="/ban/' + url_pas(data[3]) + '">(해제)</a>'
            else:
                ban = ' <a href="/ban/' + url_pas(data[3]) + '">(차단)</a>'            
            
        ip = ip_pas(data[3], None)
                
        if((int(data[0]) - 1) == 0):
            revert = ''
        else:
            revert = '<a href="/w/' + url_pas(data[1]) + '/r/' + str(int(data[0]) - 1) + '/diff/' + data[0] + '">(비교)</a> <a href="/revert/' + url_pas(data[1]) + '/r/' + str(int(data[0]) - 1) + '">(되돌리기)</a>'
        
        style = ''
        curs.execute("select * from hidhi where title = ? and re = ?", [data[1], data[0]])
        row = curs.fetchall()
        if(zdmin == 1):    
            if(row):                            
                ip += ' (숨김)'                            
                hidden = ' <a href="/history/' + url_pas(data[1]) + '/r/' + data[0] + '/hidden">(공개)'
            else:
                hidden = ' <a href="/history/' + url_pas(data[1]) + '/r/' + data[0] + '/hidden">(숨김)'
        else:
            if(row):
                ip = '숨김'
                hidden = ''
                send = '숨김'
                ban = ''
                style = 'display:none;'
            else:
                hidden = ''      
            
        div += '<tr style="' + style + '"> \
                    <td> \
                        <a href="/w/' + url_pas(data[1]) + '">' + title + '</a> (<a href="/history/' + url_pas(data[1]) + '">' + data[0] + '판</a>) ' + revert + ' (' + leng + ') \
                    </td> \
                    <td>' + ip + ban + hidden + '</td> \
                    <td>' + data[2] + '</td> \
                </tr> \
                <tr> \
                    <td colspan="3">' + send + '</td> \
                </tr>'
    else:
        div +=      '</tbody> \
                </table>'

    if(name):
        curs.execute("select end, why from ban where block = ?", [name])
        ban_it = curs.fetchall()
        if(ban_it):
            sub = '차단'
        else:
            sub = None
    else:
        sub = None
            
    return(
        template('other', 
            custom = custom_css_user(), 
            license = set_data['license'], 
            login = login_check(), 
            logo = set_data['name'], 
            data = div, 
            sub = sub,
            title = '최근 변경내역'
        )
    )
        
@route('/history/<name:path>/r/<num:int>/hidden')
def history_hidden(name = None, num = None):
    if(admin_check(6) == 1):
        curs.execute("select * from hidhi where title = ? and re = ?", [name, str(num)])
        exist = curs.fetchall()
        if(exist):
            curs.execute("delete from hidhi where title = ? and re = ?", [name, str(num)])
        else:
            curs.execute("insert into hidhi (title, re) values (?, ?)", [name, str(num)])
            
        conn.commit()
    
    return(redirect('/history/' + url_pas(name)))
        
@route('/userlog')
@route('/userlog/n/<num:int>')
def user_log(num = 1):
    if(num * 50 <= 0):
         i = 50
    else:
        i = num * 50
        
    j = i - 50
    list_data = ''
    ydmin = admin_check(1)
    
    curs.execute("select id from user limit ?, ?", [str(j), str(i)])
    user_list = curs.fetchall()
    for data in user_list:
        if(ydmin == 1):
            curs.execute("select block from ban where block = ?", [data[0]])
            ban_exist = curs.fetchall()
            if(ban_exist):
                ban_button = ' <a href="/ban/' + url_pas(data[0]) + '">(해제)</a>'
            else:
                ban_button = ' <a href="/ban/' + url_pas(data[0]) + '">(차단)</a>'
        else:
            ban_button = ''
            
        ip = ip_pas(data[0], None)
            
        list_data += '<li>' + str(j + 1) + '. ' + ip + ban_button + '</li>'
        
        j += 1
    else:
        list_data +=    '<br> \
                        <a href="/userlog/n/' + str(num - 1) + '">(이전)</a> <a href="/userlog/n/' + str(num + 1) + '">(이후)</a>'

    return(
        template('other', 
            custom = custom_css_user(), 
            license = set_data['license'], 
            login = login_check(), 
            logo = set_data['name'], 
            data = list_data, 
            title = '사용자 가입 기록'
        )
    )
        
@route('/backreset')
def xref_reset():
    if(admin_check(None) == 1):        
        curs.execute("delete from back")
        curs.execute("delete from cat")
        conn.commit()
        
        curs.execute("select title, data from data")
        data = curs.fetchall()
        for end in data:
            namumark(end[0], end[1], 1)
        
        return(redirect('/'))
    else:
        
        return(redirect('/error/3'))
        
@route('/xref/<name:path>')
@route('/xref/<name:path>/n/<num:int>')
def xref(name = None, num = 1):
    if(num * 50 <= 0):
        v = 50
    else:
        v = num * 50
        
    i = v - 50
    div = ''
    
    curs.execute("delete from back where title = ? and link = ''", [name])
    conn.commit()
    
    curs.execute("select link, type from back where title = ? order by link asc limit ?, ?", [name, str(i), str(v)])
    rows = curs.fetchall()
    for data in rows:
        div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a>'
        
        if(data[1]):
            div += ' (' + data[1] + ')'
        
        div += '</li>'
    else:        
        div += '<br> \
                <a href="/xref/' + url_pas(name) + '/n/' + str(num - 1) + '">(이전)</a> <a href="/xref/' + url_pas(name) + '/n/' + str(num + 1) + '">(이후)</a>'
    
    return(
        template('other', 
            custom = custom_css_user(), 
            license = set_data['license'], 
            login = login_check(), 
            logo = set_data['name'], 
            data = div, 
            title = name, 
            page = url_pas(name), 
            sub = '역링크'
        )
    )
        
@route('/recentdiscuss')
def recentdiscuss():
    div =   '<table style="width: 100%; text-align: center;"> \
                <tbody> \
                    <tr> \
                        <td style="width: 50%;">토론명</td> \
                        <td style="width: 50%;">시간</td> \
                    </tr>'
    
    curs.execute("select title, sub, date from rd order by date desc limit 50")
    rows = curs.fetchall()
    for data in rows:
        title = re.sub('<', '&lt;', data[0])
        title = re.sub('>', '&gt;', title)
        
        sub = re.sub('<', '&lt;', data[1])
        sub = re.sub('>', '&gt;', sub)
        
        div += '<tr> \
                    <td> \
                        <a href="/topic/' + url_pas(data[0]) + '/sub/' + url_pas(data[1]) + '">' + title + '</a> (' + sub + ') \
                    </td> \
                    <td>' + data[2] + '</td> \
                </tr>'
    else:
        div +=      '</tbody> \
                </table>'
            
    return(
        template('other', 
            custom = custom_css_user(), 
            license = set_data['license'], 
            login = login_check(), 
            logo = set_data['name'], 
            data = div, 
            title = '최근 토론내역'
        )
    )

@route('/blocklog')
@route('/blocklog/n/<number:int>')
def blocklog(num = 1):
    if(num * 50 <= 0):
        v = 50
    else:
        v = num * 50
    
    i = v - 50
    div =   '<table style="width: 100%; text-align: center;"> \
                <tbody> \
                    <tr> \
                        <td style="width: 20%;">차단자</td> \
                        <td style="width: 20%;">관리자</td> \
                        <td style="width: 20%;">기간</td> \
                        <td style="width: 20%;">이유</td> \
                        <td style="width: 20%;">시간</td> \
                    </tr>'
    
    curs.execute("select why, block, blocker, end, today from rb order by today desc limit ?, ?", [str(i), str(v)])
    rows = curs.fetchall()
    for data in rows:
        why = re.sub('<', '&lt;', data[0])
        why = re.sub('>', '&gt;', why)
        
        b = re.search("^([0-9]{1,3}\.[0-9]{1,3})$", data[1])
        if(b):
            ip = data[1] + ' (대역)'
        else:
            ip = data[1]
            
        div += '<tr> \
                    <td>' + ip + '</td> \
                    <td>' + data[2] + '</td> \
                    <td>' + data[3] + '</td> \
                    <td>' + why + '</td> \
                    <td>' + data[4] + '</td> \
                </tr>'
    else:
        div +=      '</tbody> \
                </table> \
                <br> \
                <a href="/xref/' + url_pas(name) + '/n/' + str(num - 1) + '">(이전)</a> <a href="/xref/' + url_pas(name) + '/n/' + str(num + 1) + '">(이후)</a>'
                
    return(
        template('other', 
            custom = custom_css_user(), 
            license = set_data['license'], 
            login = login_check(), 
            logo = set_data['name'], 
            data = div, 
            title = '사용자 차단 기록'
        )
    )

@route('/history/<name:path>', method=['POST', 'GET'])    
@route('/history/<name:path>/n/<num:int>', method=['POST', 'GET'])
def history_view(name = None, num = 1):
    if(request.method == 'POST'):
        return(redirect('/w/' + url_pas(name) + '/r/' + request.forms.b + '/diff/' + request.forms.a))
    else:
        select = ''
        if(num * 50 <= 0):
            i = 50
        else:
            i = num * 50
            
        j = i - 50
 
        admin1 = admin_check(1)
        admin2 = admin_check(6)
        
        div =   '<table style="width: 100%; text-align: center;"> \
                    <tbody> \
                        <tr> \
                            <td style="width: 33.3%;">판</td> \
                            <td style="width: 33.3%;">기여자</td> \
                            <td style="width: 33.3%;">시간</td> \
                        </tr>'
        
        curs.execute("select send, leng, ip, date, title, id from history where title = ? order by id + 0 desc limit ?, ?", [name, str(j), str(i)])
        all_data = curs.fetchall()
        for data in all_data:
            select += '<option value="' + data[5] + '">' + data[5] + '</option>'
            
            if(data[0]):
                send = data[0]
            else:
                send = '<br>'
                
            if(re.search("^\+", data[1])):
                leng = '<span style="color:green;">' + data[1] + '</span>'
            elif(re.search("^\-", data[1])):
                leng = '<span style="color:red;">' + data[1] + '</span>'
            else:
                leng = '<span style="color:gray;">' + data[1] + '</span>'
                
            ip = ip_pas(data[2], None)
            
            curs.execute("select block from ban where block = ?", [data[2]])
            ban_it = curs.fetchall()
            if(ban_it):
                if(admin1 == 1):
                    ban = ' <a href="/ban/' + url_pas(data[2]) + '">(해제)</a>'
                else:
                    ban = ' (X)'
            else:
                if(admin1 == 1):
                    ban = ' <a href="/ban/' + url_pas(data[2]) + '">(차단)</a>'
                else:
                    ban = ''
            
            curs.execute("select * from hidhi where title = ? and re = ?", [name, data[5]])
            hid_it = curs.fetchall()
            if(hid_it):
                if(admin2):
                    hidden = ' <a href="/history/' + url_pas(name) + '/r/' + url_pas(data[5]) + '/hidden">(공개)'
                    hid = 0
                else:
                    hid = 1
            else:
                if(admin2):
                    hidden = ' <a href="/history/' + url_pas(name) + '/r/' + url_pas(data[5]) + '/hidden">(숨김)'
                    hid = 0
                else:
                    hidden = ''
                    hid = 0
            
            if(hid == 1):
                div += '<tr> \
                            <td colspan="3">숨김</td> \
                        </tr>'
            else:
                div += '<tr> \
                            <td> \
                                ' + data[5] + '판</a> <a href="/w/' + url_pas(name) + '/r/' + url_pas(data[5]) + '">(보기)</a> \
                                    <a href="/raw/' + url_pas(name) + '/r/' + url_pas(data[5]) + '">(원본)</a> \
                                    <a href="/revert/' + url_pas(name) + '/r/' + url_pas(data[5]) + '">(되돌리기)</a> (' + leng + ') \
                            </td> \
                            <td>' + ip + ban + hidden + '</td> \
                            <td>' + data[3] + '</td> \
                        </tr> \
                        <tr> \
                            <td colspan="3">' + send + '</td> \
                        </tr>'
        else:
            div +=      '</tbody> \
                    </table> \
                    <br> \
                    <a href="/history/' + url_pas(name) + '/n/' + str(num - 1) + '">(이전)</a> <a href="/history/' + url_pas(name) + '/n/' + str(num + 1) + '">(이후)</a>'
                    
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                logo = set_data['name'], 
                data = div, 
                title = name, 
                page = url_pas(name), 
                select = select, 
                sub = '역사'
            )
        )
            
@route('/search', method=['POST'])
def search():
    return(redirect('/search/' + url_pas(request.forms.search)))

@route('/goto', method=['POST'])
def goto():
    curs.execute("select title from data where title = ?", [request.forms.search])
    data = curs.fetchall()
    if(data):
        return(redirect('/w/' + url_pas(request.forms.search)))
    else:
        return(redirect('/search/' + url_pas(request.forms.search)))

@route('/search/<name:path>')
@route('/search/<name:path>/n/<num:int>')
def deep_search(name = None, num = 1):
    if(num * 50 <= 0):
        v = num * 50
    else:
        v = 50

    i = v - 50

    div = ''
    div_plus = ''
    end = ''

    curs.execute("select title from data where title like ?", ['%' + name + '%'])
    title_list = curs.fetchall()

    curs.execute("select title from data where data like ?", ['%' + name + '%'])
    data_list = curs.fetchall()

    curs.execute("select title from data where title = ?", [name])
    exist = curs.fetchall()
    if(exist):
        div =   '<li>문서로 <a href="/w/' + url_pas(name) + '">바로가기</a></li> \
                <br>'
    else:
        div =   '<li>문서가 없습니다. <a class="not_thing" href="/w/' + url_pas(name) + '">바로가기</a></li> \
                <br>'

    if(title_list):
        no = 0

        if(data_list):
            all_list = title_list + data_list
        else:
            all_list = title_list
    else:
        if(data_list):
            no = 1

            all_list = data_list
        else:
            all_list = ''
    
    if(all_list != ''):
        for data in all_list:
            re_title = re.compile(name, re.I)
            if(re.search(re_title, data[0])):
                if(no == 0):
                    div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a> (문서명)</li>'
                else:
                    div_plus += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a> (내용)</li>'
            else:
                no = 1

                div_plus += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a> (내용)</li>'
    else:
        div += '<li>검색 결과 없음</li>'

    div += div_plus + end

    div += '<br> \
            <a href="/search/' + url_pas(name) + '/n/' + str(num - 1) + '">(이전)</a> <a href="/search/' + url_pas(name) + '/n/' + str(num + 1) + '">(이후)</a>'
    
    return(
        template('other', 
            custom = custom_css_user(), 
            license = set_data['license'], 
            login = login_check(), 
            logo = set_data['name'], 
            data = div, 
            title = name, 
            sub = '검색'
        )
    )
         
@route('/raw/<name:path>')
@route('/raw/<name:path>/r/<num:int>')
def raw_view(name = None, num = None):
    if(num):
        curs.execute("select title from hidhi where title = ? and re = ?", [name, str(num)])
        hid = curs.fetchall()
        if(hid and admin_check(6) != 1):
            return(redirect('/error/3'))
        
        curs.execute("select data from history where title = ? and id = ?", [name, str(num)])
    else:
        curs.execute("select data from data where title = ?", [name])

    rows = curs.fetchall()
    if(rows):
        enddata = re.sub('<', '&lt;', rows[0][0])
        enddata = re.sub('>', '&gt;', enddata)
        
        enddata = '<textarea style="height: 80%;">' + enddata + '</textarea>'
        
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = name, 
                logo = set_data['name'], 
                page = url_pas(name), 
                data = enddata, 
                sub = '원본'
            )
        )
    else:
        return(redirect('/w/' + url_pas(name)))
        
@route('/revert/<name:path>/r/<num:int>', method=['POST', 'GET'])
def revert(name = None, num = None):
    ip = ip_check()
    can = acl_check(ip, name)
    today = get_time()
    
    if(request.method == 'POST'):
        curs.execute("select title from hidhi where title = ? and re = ?", [name, str(num)])
        hid = curs.fetchall()
        if(hid and admin_check(6) != 1):
            return(redirect('/error/3'))

        if(can == 1):
            return(redirect('/ban'))
        else:
            curs.execute("delete from back where link = ?", [name])
            curs.execute("delete from cat where cat = ?", [name])
            conn.commit()

            curs.execute("select data from history where title = ? and id = ?", [name, str(num)])
            rows = curs.fetchall()
            if(rows):                                
                curs.execute("select data from data where title = ?", [name])
                row = curs.fetchall()
                if(row):
                    leng = leng_check(len(row[0][0]), len(rows[0][0]))
                    
                    curs.execute("update data set data = ? where title = ?", [rows[0][0], name])
                    conn.commit()
                else:
                    leng = '+' + str(len(rows[0][0]))
                    
                    curs.execute("insert into data (title, data, acl) values (?, ?, '')", [name, rows[0][0]])
                    conn.commit()
                    
                history_plus(
                    name, 
                    rows[0][0], 
                    today, 
                    ip, 
                    '문서를 ' + str(num) + '판으로 되돌렸습니다.', 
                    leng
                )
                
                return(redirect('/w/' + url_pas(name)))
    else:
        curs.execute("select title from hidhi where title = ? and re = ?", [name, str(num)])
        hid = curs.fetchall()
        if(hid and admin_check(6) != 1):
            return(redirect('/error/3'))    
                          
        if(can == 1):
            return(redirect('/ban'))
        else:
            curs.execute("select title from history where title = ? and id = ?", [name, str(num)])
            rows = curs.fetchall()
            if(rows):
                return(
                    template('revert', 
                        custom = custom_css_user(), 
                        license = set_data['license'], 
                        login = login_check(), 
                        title = name, 
                        logo = set_data['name'], 
                        page = url_pas(name), 
                        r = url_pas(str(num)), 
                        plus = '정말 되돌리시겠습니까?', 
                        sub = '되돌리기'
                    )
                )
            else:
                return(redirect('/w/' + url_pas(name)))
                    
@route('/mdel', method=['POST', 'GET'])
def many_del():
    today = get_time()
    ip = ip_check()
    if(admin_check(2) == 1):
        if(request.method == 'POST'):
            data = request.forms.content + '\r\n'
            m = re.findall('(.*)\r\n', data)
            for g in m:
                curs.execute("select data from data where title = ?", [g])
                rows = curs.fetchall()
                if(rows):
                    curs.execute("delete from back where title = ?", [g])
                    curs.execute("delete from cat where title = ?", [g])

                    leng = '-' + str(len(rows[0][0]))
                    curs.execute("delete from data where title = ?", [g])
                    history_plus(
                        g, 
                        '', 
                        today, 
                        ip, 
                        request.forms.send + ' (대량 삭제)', 
                        leng
                    )
                data = re.sub('(.*)\r\n', '', data, 1)
            conn.commit()

            return(redirect('/'))
        else:
            return(
                template('mdel', 
                    custom = custom_css_user(), 
                    license = set_data['license'], 
                    login = login_check(), 
                    title = '많은 문서 삭제', 
                    logo = set_data['name']
                )
            )
    else:
        return(redirect('/error/3'))
                
@route('/edit/<name:path>/section/<num:int>', method=['POST', 'GET'])
def section_edit(name = None, num = None):
    ip = ip_check()
    can = acl_check(ip, name)
    
    if(request.method == 'POST'):
        if(len(request.forms.send) > 500):
            return(redirect('/error/15'))
        else:
            today = get_time()
            content = savemark(request.forms.content)

            curs.execute("delete from back where link = ?", [name])
            curs.execute("delete from cat where cat = ?", [name])
            conn.commit()

            curs.execute("select data from data where title = ?", [name])
            rows = curs.fetchall()
            if(rows):
                if(request.forms.otent == content):
                    return(redirect('/error/18'))
                else:                    
                    if(can == 1):
                        return(redirect('/ban'))
                    else:
                        leng = leng_check(len(request.forms.otent), len(content))
                        content = rows[0][0].replace(request.forms.otent, content)
                        
                        history_plus(
                            name, 
                            content, 
                            today, 
                            ip, 
                            send_p(request.forms.send), 
                            leng
                        )
                        
                        curs.execute("update data set data = ? where title = ?", [content, name])
                        conn.commit()

                    include_check(name, content)
                    
                    return(redirect('/w/' + url_pas(name)))
            else:
                return(redirect('/w/' + url_pas(name)))
    else:        
        if(can == 1):
            return(redirect('/ban'))
        else:                
            curs.execute("select data from data where title = ?", [name])
            rows = curs.fetchall()
            if(rows):
                i = 0
                j = 0
                
                gdata = rows[0][0] + '\r\n'
                while(1):
                    m = re.search("((?:={1,6})\s?(?:[^=]*)\s?(?:={1,6})(?:\s+)?\n(?:(?:(?:(?!(?:={1,6})\s?(?:[^=]*)\s?(?:={1,6})(?:\s+)?\n).)*)(?:\n)?)+)", gdata)
                    if(m):
                        if(i == num - 1):
                            g = m.groups()
                            gdata = re.sub("\r\n$", "", g[0])
                            
                            break
                        else:
                            gdata = re.sub("((?:={1,6})\s?(?:[^=]*)\s?(?:={1,6})(?:\s+)?\n(?:(?:(?:(?!(?:={1,6})\s?(?:[^=]*)\s?(?:={1,6})(?:\s+)?\n).)*)(?:\n)?)+)", "", gdata, 1)
                            
                            i += 1
                    else:
                        j = 1
                        
                        break
                        
                if(j == 0):
                    gdata = re.sub("\r\n$", "", gdata)

                    return(
                        template('edit', 
                            custom = custom_css_user(), 
                            license = set_data['license'], 
                            login = login_check(), 
                            title = name, 
                            logo = set_data['name'], 
                            page = url_pas(name), 
                            data = gdata,
                            section = 1, 
                            number = num, 
                            sub = '편집'
                        )
                    )
                else:
                    return(redirect('/w/' + url_pas(name)))
            else:
                return(redirect('/w/' + url_pas(name)))

@route('/edit/<name:path>', method=['POST', 'GET'])
def edit(name = None):
    ip = ip_check()
    can = acl_check(ip, name)
    
    if(request.method == 'POST'):
        if(len(request.forms.send) > 500):
            return(redirect('/error/15'))
        else:
            today = get_time()
            content = savemark(request.forms.content)

            if(can == 1):
                return(redirect('/ban'))

            curs.execute("select data from data where title = ?", [name])
            rows = curs.fetchall()
            if(rows):
                if(rows[0][0] == content):
                    return(redirect('/error/18'))

                leng = leng_check(len(rows[0][0]), len(content))
                curs.execute("update data set data = ? where title = ?", [content, name])
            else:
                leng = '+' + str(len(content))
                curs.execute("insert into data (title, data, acl) values (?, ?, '')", [name, content])

            curs.execute("delete from back where link = ?", [name])
            curs.execute("delete from cat where cat = ?", [name])
            conn.commit()
            
            history_plus(
                name, 
                content, 
                today, 
                ip,
                send_p(request.forms.send), 
                leng
            )
                    
            include_check(name, content)
            
            return(redirect('/w/' + url_pas(name)))
    else:        
        if(can == 1):
            return(redirect('/ban'))
        else:                
            curs.execute("select data from data where title = ?", [name])
            rows = curs.fetchall()
            if(rows):
                data = rows[0][0]
            else:
                data = ''
            
            return(
                template('edit', 
                    custom = custom_css_user(), 
                    license = set_data['license'], 
                    login = login_check(), 
                    title = name, 
                    logo = set_data['name'], 
                    page = url_pas(name), 
                    data = data,
                    sub = '편집'
                )
            )

@route('/preview/<name:path>/section/<num:int>', method=['POST'])
def section_preview(name = None, num = None):
    ip = ip_check()
    can = acl_check(ip, name)
    
    if(can == 1):
        return(redirect('/ban'))
    else:            
        newdata = request.forms.content
        newdata = re.sub('^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)', ' * [[\g<in>]] 문서로 넘겨주기', newdata)
        enddata = namumark(name, newdata, 0)
            
        return(
            template('edit', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = name, 
                logo = set_data['name'], 
                page = url_pas(name), 
                data = request.forms.content, 
                preview = 1, 
                enddata = enddata, 
                section = 1, 
                number = num, 
                odata = request.forms.otent, 
                sub = '미리보기'
            )
        )
                
@route('/preview/<name:path>', method=['POST'])
def preview(name = None):
    ip = ip_check()
    can = acl_check(ip, name)
    
    if(can == 1):
        return(redirect('/ban'))
    else:            
        newdata = request.forms.content
        newdata = re.sub('^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)', ' * [[\g<in>]] 문서로 넘겨주기', newdata)
        enddata = namumark(name, newdata, 0)
            
        return(
            template('edit', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = name, 
                logo = set_data['name'], 
                page = url_pas(name), 
                data = request.forms.content, 
                preview = 1, 
                enddata = enddata, 
                sub = '미리보기'
            )
        )
        
@route('/delete/<name:path>', method=['POST', 'GET'])
def delete(name = None):
    ip = ip_check()
    can = acl_check(ip, name)
    
    if(request.method == 'POST'):
        curs.execute("select data from data where title = ?", [name])
        rows = curs.fetchall()
        if(rows):
            if(can == 1):
                return(redirect('/ban'))

            today = get_time()

            curs.execute("delete from back where link = ?", [name])
            curs.execute("delete from cat where cat = ?", [name])
            
            leng = '-' + str(len(rows[0][0]))
            history_plus(
                name, 
                '', 
                today, 
                ip, 
                request.forms.send + ' (삭제)', 
                leng
            )
            
            curs.execute("delete from data where title = ?", [name])
            conn.commit()
            
        return(redirect('/w/' + url_pas(name)))
    else:
        curs.execute("select title from data where title = ?", [name])
        rows = curs.fetchall()
        if(rows):
            if(can == 1):
                return(redirect('/ban'))

            return(
                template('del', 
                    custom = custom_css_user(), 
                    license = set_data['license'], 
                    login = login_check(),
                    title = name, 
                    logo = set_data['name'], 
                    page = url_pas(name), 
                    plus = '정말 삭제 하시겠습니까?', 
                    sub = '삭제'
                )
            )
        else:
            return(redirect('/w/' + url_pas(name)))
            
@route('/move/<name:path>', method=['POST', 'GET'])
def move(name = None):
    ip = ip_check()
    can = acl_check(ip, name)
    today = get_time()

    if(can == 1):
        return(redirect('/ban'))
    
    if(request.method == 'POST'):
        curs.execute("select data from data where title = ?", [name])
        rows = curs.fetchall()

        leng = '0'
        curs.execute("select title from history where title = ?", [request.forms.title])
        row = curs.fetchall()
        if(row):
            return(redirect('/error/19'))

        history_plus(
            name, 
            rows[0][0], 
            today, 
            ip, 
            request.forms.send + ' (<a href="/w/' + url_pas(name) + '">' + name + '</a> - <a href="/w/' + url_pas(request.forms.title) + '">' + request.forms.title + '</a> 이동)', 
            leng
        )
        
        if(rows):
            curs.execute("update data set title = ? where title = ?", [request.forms.title, name])

            curs.execute("delete from back where link = ?", [name])
            curs.execute("delete from cat where cat = ?", [name])

        curs.execute("update history set title = ? where title = ?", [request.forms.title, name])
        conn.commit()
        
        return(redirect('/w/' + url_pas(request.forms.title)))
    else:
        return(
            template('move', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = name, 
                logo = set_data['name'], 
                page = url_pas(name), 
                plus = '정말 이동 하시겠습니까?', 
                sub = '이동'
            )
        )
            
@route('/other')
def other():
    return(
        template(
            'other', 
            custom = custom_css_user(), 
            license = set_data['license'], 
            login = login_check(), 
            title = '기타 메뉴', 
            logo = set_data['name'], 
            data = '<h2 style="margin-top: 0px;">기록</h2> \
                    <li><a href="/blocklog">사용자 차단 기록</a></li> \
                    <li><a href="/userlog">사용자 가입 기록</a></li> \
                    <li><a href="/manager/6">사용자 기록</a></li> \
                    <li><a href="/manager/7">사용자 토론 기록</a></li> \
                    <h2>기타</h2> \
                    <li><a href="/titleindex">모든 문서</a></li> \
                    <li><a href="/acl_list">ACL 문서 목록</a></li> \
                    <li><a href="/admin_list">관리자 목록</a></li> \
                    <li><a href="/manager/1">관리자 메뉴</a></li> \
                    <br> \
                    이 오픈나무의 버전은 <a href="https://github.com/2DU/openNAMU/blob/normal/version.md">v' + r_ver + '</a> 입니다.'
        )
    )
    
@route('/manager', method=['POST', 'GET'])
@route('/manager/<num:int>', method=['POST', 'GET'])
def manager(num = 1):
    if(num == 1):
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '관리자 메뉴', 
                logo = set_data['name'], 
                data = '<h2 style="margin-top: 0px;">목록</h2> \
                        <li><a href="/manager/2">문서 ACL</a></li> \
                        <li><a href="/manager/3">사용자 체크</a></li> \
                        <li><a href="/manager/4">사용자 차단</a></li> \
                        <li><a href="/manager/5">관리자 권한 주기</a></li> \
                        <li><a href="/mdel">많은 문서 삭제</a></li> \
                        <li><a href="/not_close_topic">안 닫힌 토론 목록</a></li> \
                        <h2>소유자</h2> \
                        <li><a href="/backreset">역링크, 분류 다시 생성</a></li> \
                        <li><a href="/manager/8">관리 그룹 생성</a></li> \
                        <h2>기타</h2> \
                        <li>이 메뉴에 없는 기능은 해당 문서의 역사나 토론에서 바로 사용 가능함</li>'
            )
        )
    elif(num == 2):
        if(request.method == 'POST'):
            return(redirect('/acl/' + url_pas(request.forms.name)))
        else:
            return(
                template('other', 
                    custom = custom_css_user(), 
                    license = set_data['license'], 
                    login = login_check(), 
                    title = 'ACL 이동', 
                    logo = set_data['name'], 
                    data = '<form method="post"> \
                                <input name="name" type="text"> \
                                <br> \
                                <br> \
                                <button class="btn btn-primary" type="submit">이동</button> \
                            </form>'
                )
            )
    elif(num == 3):
        if(request.method == 'POST'):
            return(redirect('/check/' + url_pas(request.forms.name)))
        else:
            return(
                template('other', 
                    custom = custom_css_user(), 
                    license = set_data['license'], 
                    login = login_check(), 
                    title = '체크 이동', 
                    logo = set_data['name'], 
                    data = '<form method="post"> \
                                <input name="name" type="text"> \
                                <br> \
                                <br> \
                                <button class="btn btn-primary" type="submit">이동</button> \
                            </form>'
                )
            )
    elif(num == 4):
        if(request.method == 'POST'):
            return(redirect('/ban/' + url_pas(request.forms.name)))
        else:
            return(
                template('other', 
                    custom = custom_css_user(), 
                    license = set_data['license'], 
                    login = login_check(), 
                    title = '차단 이동', 
                    logo = set_data['name'], 
                    data = '<form method="post"> \
                                <input name="name" type="text"> \
                                <br> \
                                <br> \
                                <button class="btn btn-primary" type="submit">이동</button> \
                                <br> \
                                <br> \
                                <span>아이피 앞 두자리 (XXX.XXX) 입력하면 대역 차단</span> \
                            </form>'
                )
            )
    elif(num == 5):
        if(request.method == 'POST'):
            return(redirect('/admin/' + url_pas(request.forms.name)))
        else:
            return(
                template('other', 
                    custom = custom_css_user(), 
                    license = set_data['license'], 
                    login = login_check(), 
                    title = '권한 이동', 
                    logo = set_data['name'], 
                    data = '<form method="post"> \
                                <input name="name" type="text"> \
                                <br> \
                                <br> \
                                <button class="btn btn-primary" type="submit">이동</button> \
                            </form>'
                )
            )
    elif(num == 6):
        if(request.method == 'POST'):
            return(redirect('/record/' + url_pas(request.forms.name)))
        else:
            return(
                template('other', 
                    custom = custom_css_user(), 
                    license = set_data['license'], 
                    login = login_check(), 
                    title = '기록 이동', 
                    logo = set_data['name'], 
                    data = '<form method="post"> \
                                <input name="name" type="text"> \
                                <br> \
                                <br> \
                                <button class="btn btn-primary" type="submit">이동</button> \
                            </form>'
                )    
            )
    elif(num == 7):
        if(request.method == 'POST'):
            return(redirect('/user/' + url_pas(request.forms.name) + '/topic'))
        else:
            return(
                template('other', 
                    custom = custom_css_user(), 
                    license = set_data['license'], 
                    login = login_check(), 
                    title = '토론 기록 이동', 
                    logo = set_data['name'], 
                    data = '<form method="post"> \
                                <input name="name" type="text"> \
                                <br> \
                                <br> \
                                <button class="btn btn-primary" type="submit">이동</button> \
                            </form>'
                )    
            )
    elif(num == 8):
        if(request.method == 'POST'):
            return(redirect('/adminplus/' + url_pas(request.forms.name)))
        else:
            return(
                template('other', 
                    custom = custom_css_user(), 
                    license = set_data['license'], 
                    login = login_check(), 
                    title = '그룹 생성 이동', 
                    logo = set_data['name'], 
                    data = '<form method="post"> \
                                <input name="name" type="text"> \
                                <br> \
                                <br> \
                                <button class="btn btn-primary" type="submit">이동</button> \
                            </form>'
                )    
            )
    else:
        return(redirect('/'))
        
@route('/title_index')
def title_index():
    i = [0, 0, 0, 0, 0, 0]
    data = ''
    curs.execute("select title from data order by title asc")
    title_list = curs.fetchall()

    for list_data in title_list:
        data += '<li>' + str(i[0] + 1) + '. <a href="/w/' + url_pas(list_data[0]) + '">' + list_data[0] + '</a></li>'

        if(re.search('^분류:', list_data[0])):
            i[1] += 1
        elif(re.search('^사용자:', list_data[0])):
            i[2] += 1
        elif(re.search('^틀:', list_data[0])):
            i[3] += 1
        elif(re.search('^파일:', list_data[0])):
            i[4] += 1
        else:
            i[5] += 1
    
        i[0] += 1

    if(title_list):        
        data += '<br> \
                    <li>이 위키에는 총 ' + str(i[0]) + '개의 문서가 있습니다.</li> \
                    <br> \
                    <li>틀 문서는 총 ' + str(i[3]) + '개의 문서가 있습니다.</li> \
                    <li>분류 문서는 총 ' + str(i[1]) + '개의 문서가 있습니다.</li> \
                    <li>사용자 문서는 총 ' + str(i[2]) + '개의 문서가 있습니다.</li> \
                    <li>파일 문서는 총 ' + str(i[4]) + '개의 문서가 있습니다.</li> \
                    <li>나머지 문서는 총 ' + str(i[5]) + '개의 문서가 있습니다.</li>'
    
    return(
        template('other', 
            custom = custom_css_user(), 
            license = set_data['license'], 
            login = login_check(), 
            logo = set_data['name'], 
            data = data, 
            title = '모든 문서'
        )
    )
        
@route('/topic/<name:path>/sub/<sub:path>/b/<num:int>')
def topic_block(name = None, sub = None, num = None):
    if(admin_check(3) == 1):
        curs.execute("select block from topic where title = ? and sub = ? and id = ?", [name, sub, str(num)])
        block = curs.fetchall()
        if(block):
            if(block[0][0] == 'O'):
                curs.execute("update topic set block = '' where title = ? and sub = ? and id = ?", [name, sub, str(num)])
            else:
                curs.execute("update topic set block = 'O' where title = ? and sub = ? and id = ?", [name, sub, str(num)])
            conn.commit()
            
            rd_plus(
                name, 
                sub, 
                get_time()
            )
            
        return(redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub)))
    else:
        return(redirect('/error/3'))
        
@route('/topic/<name:path>/sub/<sub:path>/notice/<num:int>')
def topic_top(name = None, sub = None, num = None):
    if(admin_check(3) == 1):
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

            conn.commit()
            
            rd_plus(
                name, 
                sub, 
                get_time()
            )

        return(redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub)))
    else:
        return(redirect('/error/3'))

@route('/topic/<name:path>/sub/<sub:path>/tool/agree')
def topic_agree(name = None, sub = None):
    if(admin_check(3) == 1):
        ip = ip_check()
        
        curs.execute("select id from topic where title = ? and sub = ? order by id + 0 desc limit 1", [name, sub])
        topic_check = curs.fetchall()
        if(topic_check):
            time = get_time()
            
            curs.execute("select title from agreedis where title = ? and sub = ?", [name, sub])
            agree = curs.fetchall()
            if(agree):
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, '합의 결렬', ?, ?, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, time, ip])
                curs.execute("delete from agreedis where title = ? and sub = ?", [name, sub])
            else:
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, '합의 완료', ?, ?, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, time, ip])
                curs.execute("insert into agreedis (title, sub) values (?, ?)", [name, sub])
            conn.commit()
            
            rd_plus(
                name, 
                sub, 
                time
            )
              
        return(redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub)))
    else:
        
        return(redirect('/error/3'))
        
@route('/topic/<name:path>/sub/<sub:path>/tool/<tool:path>')
def topic_stop(name = None, sub = None, tool = None):
    if(tool == 'close'):
        close = 'O'
        n_close = ''
        data = '토론 닫음'
        n_data = '토론 다시 열기'
    elif(tool == 'stop'):
        close = ''
        n_close = 'O'
        data = '토론 정지'
        n_data = '토론 재 시작'
    else:
        return(redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub)))

    if(admin_check(3) == 1):
        ip = ip_check()
        
        curs.execute("select id from topic where title = ? and sub = ? order by id + 0 desc limit 1", [name, sub])
        topic_check = curs.fetchall()
        if(topic_check):
            time = get_time()
            
            curs.execute("select title from stop where title = ? and sub = ? and close = ?", [name, sub, close])
            stop = curs.fetchall()
            if(stop):
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, ?, ?, ?, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, n_data, time, ip])
                curs.execute("delete from stop where title = ? and sub = ? and close = ?", [name, sub, close])
            else:
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, ?, ?, ?, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, data, time, ip])
                curs.execute("insert into stop (title, sub, close) values (?, ?, ?)", [name, sub, close])
                curs.execute("delete from stop where title = ? and sub = ? and close = ?", [name, sub, n_close])
                
            conn.commit()
            
            rd_plus(
                name, 
                sub, 
                time
            )
            
        return(redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub)))
    else:
        
        return(redirect('/error/3'))

@route('/topic/<name:path>/sub/<sub:path>', method=['POST', 'GET'])
def topic(name = None, sub = None):
    ip = ip_check()
    ban = topic_check(ip, name, sub)
    admin = admin_check(3)
    
    if(request.method == 'POST'):
        curs.execute("select id from topic where title = ? and sub = ? order by id + 0 desc limit 1", [name, sub])
        rows = curs.fetchall()
        if(rows):
            num = int(rows[0][0]) + 1
        else:
            num = 1
        
        if(ban == 1 and admin != 1):
            return(redirect('/ban'))
        else:                    
            today = get_time()
            rd_plus(
                name, 
                sub, 
                today
            )
            
            aa = re.sub("\[\[(분류:(?:(?:(?!\]\]).)*))\]\]", "[br]", request.forms.content)
            aa = savemark(aa)
            
            curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, ?, ?, ?, '', '')", [str(num), name, sub, aa, today, ip])
            conn.commit()
            
            return(redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub)))
    else:
        style = ''
        div = ''

        curs.execute("select title from stop where title = ? and sub = ? and close = 'O'", [name, sub])
        close = curs.fetchall()

        curs.execute("select title from stop where title = ? and sub = ? and close = ''", [name, sub])
        stop = curs.fetchall()
        
        if(admin == 1):
            if(close):
                div += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/close">(토론 열기)</a> '
            else:
                div += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/close">(토론 닫기)</a> '
            
            if(stop):
                div += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/stop">(토론 재개)</a> '
            else:
                div += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/stop">(토론 정지)</a> '

            curs.execute("select title from agreedis where title = ? and sub = ?", [name, sub])
            agree = curs.fetchall()
            if(agree):
                div += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/agree">(합의 취소)</a>'
            else:
                div += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/agree">(합의 완료)</a>'
            
            div += '<br><br>'
        
        if((stop or close) and admin != 1):
            style = 'display:none;'
        
        curs.execute("select data, id, date, ip, block, top from topic where title = ? and sub = ? order by id + 0 asc", [name, sub])
        toda = curs.fetchall()

        curs.execute("select data, id, date, ip from topic where title = ? and sub = ? and top = 'O' order by id + 0 asc", [name, sub])
        top = curs.fetchall()

        for dain in top:                     
            top_data = namumark('', dain[0], 0)
            top_data = re.sub("(?P<in>#(?:[0-9]*))", '<a href="\g<in>">\g<in></a>', top_data)
                    
            ip = ip_pas(dain[3], 1)
                                
            div += '<table id="toron"> \
                        <tbody> \
                            <tr> \
                                <td id="toron_color_red"> \
                                    <a href="#' + dain[1] + '">#' + dain[1] + '</a> ' + ip + ' <span style="float:right;">' + dain[2] + '</span> \
                                </td> \
                            </tr> \
                            <tr> \
                                <td>' + top_data + '</td> \
                            </tr> \
                        </tbody> \
                    </table> \
                    <br>'
                    
        i = 0          
        for dain in toda:
            if(i == 0):
                start = dain[3]
                
            indata = namumark('', dain[0], 0)
            indata = re.sub("(?P<in>#(?:[0-9]*))", '<a href="\g<in>">\g<in></a>', indata)
            
            if(dain[4] == 'O'):
                indata = '<br>'
                block = 'id="block"'
            else:
                block = ''

            if(admin == 1):
                if(dain[4] == 'O'):
                    isblock = ' <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/b/' + str(i + 1) + '">(해제)</a>'
                else:
                    isblock = ' <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/b/' + str(i + 1) + '">(가림)</a>'

                curs.execute("select id from topic where title = ? and sub = ? and id = ? and top = 'O'", [name, sub, str(i + 1)])
                row = curs.fetchall()
                if(row):
                    isblock = isblock + ' <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/notice/' + str(i + 1) + '">(해제)</a>'
                else:
                    isblock = isblock + ' <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/notice/' + str(i + 1) + '">(공지)</a>'
                    
                curs.execute("select end from ban where block = ?", [dain[3]])
                ban_it = curs.fetchall()
                if(ban_it):
                    ban = ' <a href="/ban/' + url_pas(dain[3]) + '">(해제)</a>' + isblock
                else:
                    ban = ' <a href="/ban/' + url_pas(dain[3]) + '">(차단)</a>' + isblock
            else:
                curs.execute("select end from ban where block = ?", [dain[3]])
                ban_it = curs.fetchall()
                if(ban_it):
                    ban = ' <a href="javascript:void(0);" title="차단자">†</a>'
                else:
                    ban = ''
            
            chad = ''
            curs.execute('select acl from user where id = ?', [dain[3]])
            adch = curs.fetchall()
            if(adch and adch[0][0] != 'user'):
                chad = ' <a href="javascript:void(0);" title="관리자">★</a>'

            ip = ip_pas(dain[3], 1)
                    
            if(dain[5] == '1'):
                color = '_blue'
            elif(dain[3] == start):
                color = '_green'
            else:
                color = ''
                         
            div += '<table id="toron"> \
                        <tbody> \
                            <tr> \
                                <td id="toron_color' + color + '"> \
                                    <a href="javascript:void(0);" id="' + str(i + 1) + '">#' + str(i + 1) + '</a> ' + ip + chad + ban + ' <span style="float:right;">' + dain[2] + '</span> \
                                </td> \
                            </tr> \
                            <tr> \
                                <td ' + block + '>' + indata + '</td> \
                            </tr> \
                        </tbody> \
                    </table> \
                    <br>'
                
            i += 1
            
        return(
            template('vstopic', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = name, 
                page = url_pas(name), 
                suburl = url_pas(sub), 
                toron = sub, 
                logo = set_data['name'], 
                rows = div, 
                ban = ban, 
                style = style, 
                sub = '토론'
            )
        )
        
@route('/topic/<name:path>', method=['POST', 'GET'])
@route('/topic/<name:path>/<tool:path>', method=['GET'])
def close_topic_list(name = None, tool = None):
    div = ''
    i = 0
    list_d = 0

    if(request.method == 'POST'):
        t_num = ''
        while(1):
            curs.execute("select title from topic where title = ? and sub = ? limit 1", [name, request.forms.topic + t_num])
            t_data = curs.fetchall()
            if(t_data):
                if(t_num == ''):
                    t_num = ' 2'
                else:
                    t_num = ' ' + str(int(t_num.replace(' ', '')) + 1)
            else:
                break

        return(redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(request.forms.topic + t_num)))
    else:
        if(tool == 'close'):
            curs.execute("select sub from stop where title = ? and close = 'O' order by sub asc", [name])
            sub = '닫힘'
        elif(tool == 'agree'):
            curs.execute("select sub from agreedis where title = ? order by sub asc", [name])
            sub = '합의'
        else:
            list_d = 1
            curs.execute("select sub from rd where title = ? order by date desc", [name])
            sub = '토론 목록'

        rows = curs.fetchall()
        for data in rows:
            curs.execute("select data, date, ip, block from topic where title = ? and sub = ? and id = '1'", [name, data[0]])
            row = curs.fetchall()
            if(row):
                indata = namumark(name, row[0][0], 0)
                
                if(row[0][3] == 'O'):
                    indata = '<br>'
                    block = 'id="block"'
                else:
                    block = ''

                ip = ip_pas(row[0][2], 1)
                    
                div += '<h2> \
                            <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(data[0]) + '">' + str((i + 1)) + '. ' + data[0] + '</a> \
                        </h2> \
                        <table id="toron"> \
                            <tbody> \
                                <tr> \
                                    <td id="toron_color_green"> \
                                        <a href="javascript:void(0);" id="1">#1</a> ' + ip + ' <span style="float:right;">' + row[0][1] + '</span> \
                                    </td> \
                                </tr> \
                                <tr> \
                                    <td ' + block + '>' + indata + '</td> \
                                </tr> \
                            </tbody> \
                        </table> \
                        <br>'
                
                i += 1
        
        return(
            template('topic', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = name, 
                page = url_pas(name), 
                logo = set_data['name'], 
                plus = div, 
                sub = sub,
                list = list_d
            )
        )
        
@route('/login', method=['POST', 'GET'])
def login():
    session = request.environ.get('beaker.session')
    ip = ip_check()
    ban = ban_check(ip)
        
    if(request.method == 'POST'):        
        if(ban == 1):
            return(redirect('/ban'))

        curs.execute("select pw from user where id = ?", [request.forms.id])
        user = curs.fetchall()
        if(user):
            if(session.get('Now') == 1):
                return(redirect('/error/11'))

            if(bcrypt.checkpw(bytes(request.forms.pw, 'utf-8'), bytes(user[0][0], 'utf-8'))):
                session['Now'] = 1
                session['DREAMER'] = request.forms.id

                curs.execute("select css from custom where user = ?", [request.forms.id])
                css_data = curs.fetchall()
                if(css_data):
                    session['Daydream'] = css_data[0][0]
                else:
                    session['Daydream'] = ''
                
                curs.execute("insert into login (user, ip, today) values (?, ?, ?)", [request.forms.id, ip, get_time()])
                conn.commit()
                
                return(redirect('/user'))
            else:
                return(redirect('/error/13'))
        else:
            return(redirect('/error/12'))
    else:        
        if(ban == 1):
            return(redirect('/ban'))

        if(session.get('Now') == 1):
            return(redirect('/error/11'))

        return(
            template('login', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '로그인', 
                enter = '로그인', 
                logo = set_data['name']
            )
        )
                
@route('/change', method=['POST', 'GET'])
def change_password():
    ip = ip_check()
    ban = ban_check(ip)
    
    if(request.method == 'POST'):      
        if(request.forms.pw2 == request.forms.pw3):
            if(ban == 1):
                return(redirect('/ban'))

            curs.execute("select id from user where id = ?", [request.forms.id])
            user = curs.fetchall()
            if(user):
                if(not re.search('(\.|:)', ip)):
                    return(redirect('/logout'))

                if(bcrypt.checkpw(bytes(request.forms.pw, 'utf-8'), bytes(user[0]['pw'], 'utf-8'))):
                    hashed = bcrypt.hashpw(bytes(request.forms.pw2, 'utf-8'), bcrypt.gensalt())
                    
                    curs.execute("update user set pw = ? where id = ?", [hashed.decode(), request.forms.id])
                    conn.commit()
                    
                    return(redirect('/login'))
                else:
                    return(redirect('/error/10'))
            else:
                return(redirect('/error/9'))
        else:
            return(redirect('/error/20'))
    else:        
        if(ban == 1):
            return(redirect('/ban'))

        if(not re.search('(\.|:)', ip)):
            return(redirect('/logout'))
                
        return(
            template('login', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '비밀번호 변경', 
                enter = '변경', 
                logo = set_data['name']
            )
        )
                
@route('/check/<name:path>')
def user_check(name = None):
    curs.execute("select acl from user where id = ?", [name])
    user = curs.fetchall()
    if(user and user[0][0] != 'user'):
        return(redirect('/error/4'))

    if(admin_check(4) == 1):
        m = re.search('^(?:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}?)$', name)
        if(m):
            sql = 'ip'
        else:
            sql = 'user'

        curs.execute("select user, ip, today from login where ? = ? order by today desc", [sql, name])
        row = curs.fetchall()
        if(row):
            c = '<div> \
                    <table style="width: 100%; text-align: center;"> \
                        <tbody> \
                            <tr> \
                                <td style="width: 33.3%;">이름</td> \
                                <td style="width: 33.3%;">아이피</td> \
                                <td style="width: 33.3%;">언제</td> \
                            </tr>'
            for data in row:
                c +=    '<tr> \
                            <td>' + data[0] + '</td> \
                            <td>' + data[1] + '</td> \
                            <td>' + data[2] + '</td> \
                        </tr>'
            else:
                c +=            '</tbody> \
                            </table> \
                        </div>'
        else:
            c = ''
                
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '다중 검사', 
                logo = set_data['name'], 
                data = c
            )
        )
    else:
        return(redirect('/error/3'))
                
@route('/register', method=['POST', 'GET'])
def register():
    ip = ip_check()
    ban = ban_check(ip)

    if(ban == 1):
        return(redirect('/ban'))
    
    if(request.method == 'POST'):        
        if(request.forms.pw == request.forms.pw2):
            m = re.search('(?:[^A-Za-zㄱ-힣0-9 ])', request.forms.id)
            if(m):
                return(redirect('/error/8'))

            if(len(request.forms.id) > 20):
                return(redirect('/error/7'))

            curs.execute("select id from user where id = ?", [request.forms.id])
            rows = curs.fetchall()
            if(rows):
                return(redirect('/error/6'))

            hashed = bcrypt.hashpw(bytes(request.forms.pw, 'utf-8'), bcrypt.gensalt())
            
            curs.execute("select id from user limit 1")
            user_ex = curs.fetchall()
            if(not user_ex):
                curs.execute("insert into user (id, pw, acl) values (?, ?, 'owner')", [request.forms.id, hashed.decode()])
            else:
                curs.execute("insert into user (id, pw, acl) values (?, ?, 'user')", [request.forms.id, hashed.decode()])
            conn.commit()
            
            return(redirect('/login'))
        else:
            return(redirect('/error/20'))
    else:        
        return(
            template('login', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '회원가입', 
                enter = '회원가입', 
                logo = set_data['name']
            )
        )
            
@route('/logout')
def logout():
    session = request.environ.get('beaker.session')
    session['Now'] = 0
    session.pop('DREAMER', None)

    return(redirect('/user'))
    
@route('/ban/<name:path>', method=['POST', 'GET'])
def user_ban(name = None):
    curs.execute("select acl from user where id = ?", [name])
    user = curs.fetchall()
    if(user and user[0][0] != 'user'):
        return(redirect('/error/4'))

    if(request.method == 'POST'):
        if(admin_check(1) == 1):
            ip = ip_check()
            
            if(not re.search("[0-9]{4}-[0-9]{2}-[0-9]{2}", request.forms.end)):
                end = ''
            else:
                end = request.forms.end

            curs.execute("select block from ban where block = ?", [name])
            row = curs.fetchall()
            if(row):
                rb_plus(name, '해제', get_time(), ip, '')
                
                curs.execute("delete from ban where block = ?", [name])
            else:
                b = re.search("^([0-9]{1,3}\.[0-9]{1,3})$", name)
                if(b):
                    band_d = 'O'
                else:
                    band_d = ''

                rb_plus(name, end, get_time(), ip, request.forms.why)

                curs.execute("insert into ban (block, end, why, band) values (?, ?, ?, ?)", [name, end, request.forms.why, band_d])
            conn.commit()

            return(redirect('/ban/' + url_pas(name)))
        else:
            return(redirect('/error/3'))
    else:
        if(admin_check(1) == 1):
            curs.execute("select * from ban where block = ?", [name])
            row = curs.fetchall()
            if(row):
                now = '차단 해제'
            else:
                b = re.search("^([0-9]{1,3}\.[0-9]{1,3})$", name)
                if(b):
                    now = '대역 차단'
                else:
                    now = '차단'

            return(
                template('ban', 
                    custom = custom_css_user(), 
                    license = set_data['license'], 
                    login = login_check(), 
                    title = name, 
                    page = url_pas(name), 
                    logo = set_data['name'], 
                    now = now, 
                    today = get_time(), 
                    sub = '차단'
                )
            )
        else:
            return(redirect('/error/3'))
                
@route('/acl/<name:path>', method=['POST', 'GET'])
def acl(name = None):
    if(request.method == 'POST'):
        if(admin_check(5) == 1):
            curs.execute("select acl from data where title = ?", [name])
            row = curs.fetchall()
            if(row):
                if(request.forms.select == 'admin'):
                   curs.execute("update data set acl = 'admin' where title = ?", [name])
                elif(request.forms.select == 'user'):
                    curs.execute("update data set acl = 'user' where title = ?", [name])
                else:
                    curs.execute("update data set acl = '' where title = ?", [name])
                    
                conn.commit()
                
            return(redirect('/w/' + url_pas(name)))
        else:
            return(redirect('/error/3'))
    else:
        if(admin_check(5) == 1):
            curs.execute("select acl from data where title = ?", [name])
            row = curs.fetchall()
            if(row):
                if(row[0][0] == 'admin'):
                    now = '관리자만'
                elif(row[0][0] == 'user'):
                    now = '로그인 이상'
                else:
                    now = '일반'
                
                return(
                    template('acl', 
                        custom = custom_css_user(), 
                        license = set_data['license'], 
                        login = login_check(), 
                        title = name, 
                        page = url_pas(name), 
                        logo = set_data['name'], 
                        now = '현재 ACL 상태는 ' + now, 
                        sub = 'ACL'
                    )
                )
            else:
                return(redirect('/w/' + url_pas(name)) )
        else:
            return(redirect('/error/3'))
            
@route('/admin/<name:path>', method=['POST', 'GET'])
def user_admin(name = None):
    if(request.method == 'POST'):
        if(admin_check(None) == 1):
            curs.execute("select acl from user where id = ?", [name])
            user = curs.fetchall()
            if(user):
                if(user[0][0] != 'user'):
                    curs.execute("update user set acl = 'user' where id = ?", [name])
                else:
                    curs.execute("update user set acl = ? where id = ?", [request.forms.select, name])
                conn.commit()
                
                return(redirect('/'))
            else:
                return(redirect('/error/5'))
        else:
            return(redirect('/error/3'))
    else:
        if(admin_check(None) == 1):
            curs.execute("select acl from user where id = ?", [name])
            user = curs.fetchall()
            if(user):
                if(user[0][0] != 'user'):
                    now = '권한 해제'
                else:
                    now = '권한 부여'
                    
                div = ''
                    
                curs.execute('select name from alist order by name asc')
                get_alist = curs.fetchall()
                if(get_alist):
                    i = 0
                    name_rem = ''
                    for data in get_alist:
                        if(name_rem != data[0]):
                            name_rem = data[0]
                            div += '<option value="' + data[0] + '" selected="selected">' + data[0] + '</option>'
                
                return(
                    template('admin', 
                        custom = custom_css_user(), 
                        license = set_data['license'], 
                        login = login_check(), 
                        title = name, 
                        page = url_pas(name), 
                        datalist = div, 
                        logo = set_data['name'], 
                        now = now, 
                        sub = '권한 부여'
                    )
                )
            else:
                return(redirect('/error/5'))
        else:
            return(redirect('/error/3'))
            
@route('/ban')
def are_you_ban():
    ip = ip_check()
    
    if(ban_check(ip) == 1):
        curs.execute("select end, why from ban where block = ?", [ip])
        rows = curs.fetchall()
        if(not rows):
            data = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
            if(data):
                results = data.groups()
                curs.execute("select end, why from ban where block = ? and band = 'O'", [results[0]])

                rows = curs.fetchall()

        if(rows):
            if(rows[0][0]):
                end = rows[0][0] + ' 까지 차단 상태 입니다. / 사유 : ' + rows[0][1]                

                now = re.sub(':', '', get_time())
                now = re.sub('\-', '', now)
                now = int(re.sub(' ', '', now))
                
                day = re.sub('\-', '', rows[0][0])    
                
                if(now >= int(day + '000000')):
                    curs.execute("delete from ban where block = ?", [ip])
                    conn.commit()
                    
                    end = '차단이 풀렸습니다. 다시 시도 해 보세요.'
            else:
                end = '영구 차단 상태 입니다. / 사유 : ' + rows[0][1]
        else:
            end = '권한이 맞지 않는 상태 입니다.'
    else:
        end = '권한이 맞지 않는 상태 입니다.'

    return(
        template('other', 
            custom = custom_css_user(), 
            license = set_data['license'], 
            login = login_check(), 
            title = '권한 오류', 
            logo = set_data['name'], 
            data = end
        )
    )
    
@route('/w/<name:path>/r/<a:int>/diff/<b:int>')
def diff_data(name = None, a = None, b = None):
    curs.execute("select data from history where id = ? and title = ?", [str(a), name])
    a_raw_data = curs.fetchall()
    if(a_raw_data):
        curs.execute("select data from history where id = ? and title = ?", [str(b), name])
        b_raw_data = curs.fetchall()
        if(b_raw_data):
            a_data = re.sub('<', '&lt;', a_raw_data[0][0])
            a_data = re.sub('>', '&gt;', a_data)
            
            b_data = re.sub('<', '&lt;', b_raw_data[0][0])
            b_data = re.sub('>', '&gt;', b_data)
            
            diff_data = difflib.SequenceMatcher(None, a_data, b_data)
            result = diff(diff_data)
            
            result = '<pre>' + result + '</pre>'
            
            return(
                template('other', 
                    custom = custom_css_user(), 
                    license = set_data['license'], 
                    login = login_check(), 
                    title = name, 
                    logo = set_data['name'], 
                    data = result, 
                    sub = '비교', 
                    page = url_pas(name)
                )
            )

    return(redirect('/history/' + url_pas(name)))
        
@route('/down/<name:path>')
def down(name = None):
    curs.execute("select title from data where title like ?", ['%' + name + '/%'])
    under = curs.fetchall()
    
    div = ''
    i = 0

    for data in under:
        div += '<li>' + str(i + 1) + '. <a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a></li>'
        i += 1
    
    return(
        template('other', 
            custom = custom_css_user(), 
            license = set_data['license'], 
            login = login_check(), 
            title = name, 
            logo = set_data['name'], 
            data = div, 
            sub = '하위 문서', 
            page = url_pas(name)
        )
    )

@route('/w/<name:path>')
@route('/w/<name:path>/r/<num:int>')
@route('/w/<name:path>/from/<redirect:path>')
def read_view(name = None, num = None, redirect = None):
    data_none = 0
    sub = 0
    acl = ''
    div = ''
    topic = ''
    
    curs.execute("select sub from rd where title = ? order by date desc", [name])
    rows = curs.fetchall()
    for data in rows:
        curs.execute("select title from stop where title = ? and sub = ? and close = 'O'", [name, data[0]])
        row = curs.fetchall()
        if(not row):
            topic = "open"
            
            break
                
    curs.execute("select title from data where title like ?", ['%' + name + '/%'])
    under = curs.fetchall()
    if(under):
        down = 1
    else:
        down = 0
        
    m = re.search("^(.*)\/(.*)$", name)
    if(m):
        uppage = m.groups()[0]
    else:
        uppage = 0
        
    if(admin_check(5) == 1):
        admin_memu = 'ACL'
    else:
        admin_memu = ''
        
    if(re.search("^분류:", name)):
        curs.execute("delete from cat where title = ? and cat = ''", [name])
        conn.commit()
        
        curs.execute("select cat from cat where title = ? order by cat asc", [name])
        rows = curs.fetchall()
        if(rows):
            div = '<h2>분류</h2>'
            i = 0
            
            for data in rows:       
                div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a></li>'

    if(num):
        curs.execute("select title from hidhi where title = ? and re = ?", [name, str(num)])
        hid = curs.fetchall()
        if(hid and admin_check(6) != 1):
            return(redirect('/history/' + url_pas(name)))

        curs.execute("select title, data from history where title = ? and id = ?", [name, str(num)])
    else:
        curs.execute("select acl, data from data where title = ?", [name])

    rows = curs.fetchall()
    if(rows):
        if(not num):
            if(rows[0][0] == 'admin'):
                acl = '(관리자)'
            elif(rows[0][0] == 'user'):
                acl = '(로그인)'
                
        elsedata = rows[0][1]
    else:
        data_none = 1
        elsedata = ''

    m = re.search("^사용자:([^/]*)", name)
    if(m):
        g = m.groups()
        
        curs.execute("select acl from user where id = ?", [g[0]])
        test = curs.fetchall()
        if(test and test[0][0] != 'user'):
            acl = '(관리자)'

        curs.execute("select block from ban where block = ?", [g[0]])
        user = curs.fetchall()
        if(user):
            sub = '차단'
            
    if(redirect):
        elsedata = re.sub("^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)", " * [[\g<in>]] 문서로 넘겨주기", elsedata)
            
    enddata = namumark(name, elsedata, 1)
    
    return(
        template('read', 
            custom = custom_css_user(), 
            license = set_data['license'], 
            login = login_check(), 
            title = name, 
            logo = set_data['name'], 
            page = url_pas(name), 
            data = enddata + div, 
            uppage = uppage, 
            acl = acl, 
            topic = topic, 
            redirect = redirect, 
            admin = admin_memu, 
            data_none = data_none, 
            sub = sub, 
            down = down
        )
    )

@route('/user/<name:path>/topic')
@route('/user/<name:path>/topic/<num:int>')
def user_topic_list(name = None, num = 1):
    if(num * 50 <= 0):
        v = 50
    else:
        v = num * 50
    
    i = v - 50
    ydmin = admin_check(1)
    div =   '<table style="width: 100%; text-align: center;"> \
                <tbody> \
                    <tr> \
                        <td style="width: 33.3%;">토론명</td> \
                        <td style="width: 33.3%;">작성자</td> \
                        <td style="width: 33.3%;">시간</td> \
                    </tr>'
    
    curs.execute("select title, id, sub, ip, date from topic where ip = ? order by date desc limit ?, ?", [name, str(i), str(v)])
    rows = curs.fetchall()
    if(rows):
        for data in rows:
            title = re.sub('<', '&lt;', data[0])
            title = re.sub('>', '&gt;', title)

            sub = re.sub('<', '&lt;', data[2])
            sub = re.sub('>', '&gt;', sub)
                
            if(ydmin == 1):
                curs.execute("select * from ban where block = ?", [data[3]])
                row = curs.fetchall()
                if(row):
                    ban = ' <a href="/ban/' + url_pas(data[3]) + '">(해제)</a>'
                else:
                    ban = ' <a href="/ban/' + url_pas(data[3]) + '">(차단)</a>'
            else:
                ban = ''
                
            ip = ip_pas(data[3], 1)
                
            div += '<tr> \
                        <td> \
                            <a href="/topic/' + url_pas(data[0]) + '/sub/' + url_pas(data[2]) + '#' + data[1] + '">' + title + '#' + data[1] + '</a> (' + sub + ') \
                        </td> \
                        <td>' + ip + ban +  '</td> \
                        <td>' + data[4] + '</td> \
                    </tr>'
        else:
            div +=      '</tbody> \
                    </table>'
    else:
        div = ''
        
    div += '<br> \
            <a href="/user/' + url_pas(name) + '/topic/' + str(num - 1) + '">(이전)</a> <a href="/user/' + url_pas(name) + '/topic/' + str(num + 1) + '">(이후)</a>'
                
    curs.execute("select end, why from ban where block = ?", [name])
    ban_it = curs.fetchall()
    if(ban_it):
        sub = '차단'
    else:
        sub = None 
    
    return(
        template('other', 
            custom = custom_css_user(), 
            license = set_data['license'], 
            login = login_check(), 
            logo = set_data['name'], 
            data = div, 
            title = '사용자 토론 기록', 
            sub = sub
        )
    )
        
@route('/user')
def user_info():
    ip = ip_check()
    raw_ip = ip
    
    curs.execute("select acl from user where id = ?", [ip])
    rows = curs.fetchall()
    if(ban_check(ip) == 0):
        if(rows):
            if(rows[0][0] != 'user'):
                acl = rows[0][0]
            else:
                acl = '로그인'
        else:
            acl = '일반'
    else:
        acl = '차단'
        
    ip = ip_pas(ip, 2)

    return(
        template('other', 
            custom = custom_css_user(), 
            license = set_data['license'], 
            login = login_check(), 
            title = '사용자 메뉴', 
            logo = set_data['name'], 
            data =  ip + \
                    '<br> \
                    <br> \
                    <span>권한 상태 : ' + acl + '</span> \
                    <h2>로그인 관련</h2> \
                    <li><a href="/login">로그인</a></li> \
                    <li><a href="/logout">로그아웃</a></li> \
                    <li><a href="/register">회원가입</a></li> \
                    <h2>기타</h2> \
                    <li><a href="/change">비밀번호 변경</a></li> \
                    <li><a href="/count">기여 횟수</a></li> \
                    <li><a href="/custom">커스텀 CSS</a></li>'
        )
    )

@route('/custom', method=['GET', 'POST'])
def custom_css():
    session = request.environ.get('beaker.session')
    ip = ip_check()
    if(request.method == 'POST'):
        if(not re.search('(\.|:)', ip)):
            curs.execute("select * from custom where user = ?", [ip])
            css_data = curs.fetchall()
            if(css_data):
                curs.execute("update custom set css = ? where user = ?", [request.forms.content, ip])
            else:
                curs.execute("insert into custom (user, css) values (?, ?)", [ip, request.forms.content])
            conn.commit()

        session['Daydream'] = request.forms.content

        return(redirect('/user'))
    else:
        if(not re.search('(\.|:)', ip)):
            start = ''
            curs.execute("select css from custom where user = ?", [ip])
            css_data = curs.fetchall()
            if(css_data):
                data = css_data[0][0]
            else:
                data = ''
        else:
            start = '<span>비 로그인의 경우에는 로그인하면 날아갑니다.</span><br><br>'
            try:
                data = session['Daydream']
            except:
                data = ''

        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '커스텀 CSS', 
                logo = set_data['name'], 
                data =  start + \
                        '<form method="post"> \
                            <textarea rows="30" cols="100" name="content" form="usrform">'\
                                 + data + \
                            '</textarea> \
                            <br> \
                            <br> \
                            <div class="form-actions"> \
                                <button class="btn btn-primary" type="submit">저장</button> \
                            </div> \
                        </form>'
            )
        )
    
@route('/count')
def count_edit():
    curs.execute("select count(title) from history where ip = '" + ip_check() + "'")
    count = curs.fetchall()
    if(count):
        data = count[0][0]
    else:
        data = 0
    
    return(
        template('other', 
            custom = custom_css_user(), 
            license = set_data['license'], 
            login = login_check(), 
            title = '기여 횟수', 
            logo = set_data['name'], 
            data = "기여 횟수 : " + str(data)
        )
    )
        
@route('/random')
def random():
    curs.execute("select title from data order by random() limit 1")
    rows = curs.fetchall()
    if(rows):
        return(redirect('/w/' + url_pas(rows[0][0])))
    else:
        return(redirect('/'))
    
@route('/views/<name:path>')
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
        
    return(
        static_file(rename, 
            root = './views' + plus
        )
    )
        
@route('/error/<num:int>')
def error_test(num = None):
    if(num == 1):
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '권한 오류', 
                logo = set_data['name'], 
                data = '비 로그인 상태 입니다.'
            )
        )
    elif(num == 2):
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '권한 오류', 
                logo = set_data['name'], 
                data = '이 계정이 없습니다.'
            )
        )
    elif(num == 3):
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '권한 오류',
                logo = set_data['name'], 
                data = '권한이 모자랍니다.'
            )
        )
    elif(num == 4):
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '권한 오류', 
                logo = set_data['name'], 
                data = '관리자는 차단, 검사 할 수 없습니다.'
            )
        )
    elif(num == 5):
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '사용자 오류', 
                logo = set_data['name'], 
                data = '그런 계정이 없습니다.'
            )
        )
    elif(num == 6):
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '가입 오류', 
                logo = set_data['name'], 
                data = '동일한 아이디의 사용자가 있습니다.'
            )
        )
    elif(num == 7):
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '가입 오류', 
                logo = set_data['name'], 
                data = '아이디는 20글자보다 짧아야 합니다.'
            )
        )
    elif(num == 8):
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '가입 오류', 
                logo = set_data['name'], 
                data = '아이디에는 한글과 알파벳과 공백만 허용 됩니다.'
            )
        )
    elif(num == 9):
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '변경 오류', 
                logo = set_data['name'], 
                data = '그런 계정이 없습니다.'
            )
        )
    elif(num == 10):
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '변경 오류', 
                logo = set_data['name'], 
                data = '비밀번호가 다릅니다.'
            )
        )
    elif(num == 11):
        return(
            template('other',
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '로그인 오류', 
                logo = set_data['name'], 
                data = '이미 로그인 되어 있습니다.'
            )
        )
    elif(num == 12):
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '로그인 오류', 
                logo = set_data['name'], 
                data = '그런 계정이 없습니다.'
            )
        )
    elif(num == 13):
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '로그인 오류', 
                logo = set_data['name'], 
                data = '비밀번호가 다릅니다.'
            )
        )
    elif(num == 14):
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '업로드 오류', 
                logo = set_data['name'], 
                data = 'jpg, gif, jpeg, png만 가능 합니다.'
            )
        )
    elif(num == 15):
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '편집 오류', 
                logo = set_data['name'], 
                data = '편집 기록은 500자를 넘을 수 없습니다.'
            )
        )
    elif(num == 16):
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '업로드 오류', 
                logo = set_data['name'], 
                data = '동일한 이름의 파일이 있습니다.'
            )
        )
    elif(num == 17):
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '업로드 오류', 
                logo = set_data['name'], 
                data = '파일 용량은 ' + set_data['upload'] + 'MB를 넘길 수 없습니다.'
            )
        )
    elif(num == 18):
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '편집 오류', 
                logo = set_data['name'], 
                data = '내용이 원래 문서와 동일 합니다.'
            )
        )
    elif(num == 19):
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '이동 오류', 
                logo = set_data['name'], 
                data = '이동 하려는 곳에 문서가 이미 있습니다.'
            )
        )
    elif(num == 20):
        return(
            template('other', 
                custom = custom_css_user(), 
                license = set_data['license'], 
                login = login_check(), 
                title = '비밀번호 오류', 
                logo = set_data['name'],
                data = '재 확인이랑 비밀번호가 다릅니다.'
            )
        )
    else:
        return(redirect('/'))

@error(404)
def error_404(error):
    return(redirect('/w/' + url_pas(set_data['frontpage'])))
    
run(
    app = app, 
    server = 'tornado', 
    host = '0.0.0.0', 
    port = int(set_data['port'])
)
