from bottle import request, app, template
from bottle.ext import beaker
import json
import sqlite3
from css_html_js_minify import html_minify

json_data = open('set.json').read()
set_data = json.loads(json_data)

conn = sqlite3.connect(set_data['db'] + '.db')
curs = conn.cursor()
    
session_opts = {
    'session.type': 'dbm',
    'session.data_dir': './app_session/',
    'session.auto': 1
}

app = beaker.middleware.SessionMiddleware(app(), session_opts)

from mark import *

def other2(d):
    g = ''
    session = request.environ.get('beaker.session')
    if(session.get('View_List')):
        m = re.findall('(?:(?:([^\n]+)\n))', session.get('View_List'))
        if(m):
            g = ''
            for z in m[-6:-1]:
                g += '<a href="/w/' + url_pas(z) + '">' +  html.escape(z) + '</a> / '
            g = re.sub(' / $', '', g)
            
    r = d + [g]
    return(r)

def include(title, old, new):
    if(re.search('^틀:', title)):
        old_d = mid_pas(old, 0, 1, 1)[0]
        new_d = mid_pas(new, 0, 1, 1)[0]

        m1 = re.findall("\[\[(분류:(?:(?:(?!\]\]).)*))\]\]", old_d)
        m2 = re.findall("\[\[(분류:(?:(?:(?!\]\]).)*))\]\]", new_d)

        curs.execute("select link from back where title = ? and type = 'include'", [title])
        d1 = curs.fetchall()
        for x in d1:
            for y in m1:
                curs.execute("delete from back where link = ? and type = 'cat'", [y])

            for z in m2:
                backlink_plus(x[0], z, 'cat', 1)

        conn.commit()
    

def wiki_set(num):
    if(num == 1):
        r = []

        curs.execute('select data from other where name = ?', ['name'])
        d = curs.fetchall()
        if(d):
            r += [d[0][0]]
        else:
            r += ['무명위키']

        curs.execute('select data from other where name = "license"')
        d = curs.fetchall()
        if(d):
            r += [d[0][0]]
        else:
            r += ['CC 0']

        curs.execute("select data from other where name = 'css'")
        d = curs.fetchall()
        if(d):
            r += [d[0][0]]
        else:
            r += ['']

        curs.execute("select data from other where name = 'js'")
        d = curs.fetchall()
        if(d):
            r += [d[0][0]]
        else:
            r += ['']

        curs.execute('select data from other where name = "logo"')
        d = curs.fetchall()
        if(d):
            r += [d[0][0]]
        else:
            r += r[0]

        return(r)

    if(num == 2):
        d = '위키:대문'
        curs.execute('select data from other where name = "frontpage"')
    elif(num == 3):
        d = '2'
        curs.execute('select data from other where name = "upload"')
    
    r = curs.fetchall()
    if(r):
        return(r[0][0])
    else:
        return(d)

def diff(seqm, num):
    output= []
    for opcode, a0, a1, b0, b1 in seqm.get_opcodes():
        if(opcode == 'equal' and num == 1):
            output.append(seqm.a[a0:a1])
        elif(opcode == 'insert' and num == 0):
            output.append("<span style='background:#CFC;'>" + seqm.b[b0:b1] + "</span>")
        elif(opcode == 'delete' and num == 1):
            output.append("<span style='background:#FDD;'>" + seqm.a[a0:a1] + "</span>")
        elif(opcode == 'replace'):
            if(num == 1):
                output.append("<span style='background:#FDD;'>" + seqm.a[a0:a1] + "</span>")
            else:
                output.append("<span style='background:#CFC;'>" + seqm.b[b0:b1] + "</span>")
        elif(num == 0):
            output.append(seqm.b[b0:b1])
            
    return(''.join(output))
           
def admin_check(num, what):
    ip = ip_check() 
    curs.execute("select acl from user where id = ?", [ip])
    user = curs.fetchall()
    if(user):
        reset = 0
        while(1):
            if(num == 1 and reset == 0):
                check = 'ban'
            elif(num == 2 and reset == 0):
                check = 'mdel'
            elif(num == 3 and reset == 0):
                check = 'toron'
            elif(num == 4 and reset == 0):
                check = 'check'
            elif(num == 5 and reset == 0):
                check = 'acl'
            elif(num == 6 and reset == 0):
                check = 'hidel'
            else:
                check = 'owner'

            curs.execute('select name from alist where name = ? and acl = ?', [user[0][0], check])
            acl_data = curs.fetchall()
            if(acl_data):
                if(what):
                    curs.execute("insert into re_admin (who, what, time) values (?, ?, ?)", [ip, what, get_time()])
                    conn.commit()

                return(1)
            else:
                if(reset == 0):
                    reset = 1
                else:
                    break

def ip_pas(raw_ip):
    if(re.search("(\.|:)", raw_ip)):
        ip = raw_ip
    else:
        curs.execute("select title from data where title = ?", ['사용자:' + raw_ip])
        data = curs.fetchall()
        if(data):
            ip = '<a href="/w/' + url_pas('사용자:' + raw_ip) + '">' + raw_ip + '</a>'
        else:
            ip = '<a class="not_thing" href="/w/' + url_pas('사용자:' + raw_ip) + '">' + raw_ip + '</a>'
            
    ip += ' <a href="/record/' + url_pas(raw_ip) + '">(기록)</a>'

    return(ip)

def custom():
    session = request.environ.get('beaker.session')
    try:
        d1 = format(session['Daydream'])
    except:
        d1 = ''

    try:
        d2 = format(session['AQUARIUM'])
    except:
        d2 = ''

    if(session.get('Now') == 1):
        curs.execute('select name from alarm limit 1')
        if(curs.fetchall()):
            d3 = 2
        else:
            d3 = 1
    else:
        d3 = 0

    return([d1, d2, d3])

def acl_check(name):
    ip = ip_check()
    band = re.search("^([0-9]{1,3}\.[0-9]{1,3})", ip)
    if(band):
        band_it = band.groups()
    else:
        band_it = ['Not']

    curs.execute("select block from ban where block = ? and band = 'O'", [band_it[0]])
    band_d = curs.fetchall()

    curs.execute("select block from ban where block = ?", [ip])
    ban_d = curs.fetchall()
    if(band_d or ban_d):
        return(1)

    acl_c = re.search("^사용자:([^/]*)", name)
    if(acl_c):
        acl_n = acl_c.groups()

        if(admin_check(5, None) == 1):
            return(0)

        curs.execute("select acl from data where title = ?", ['사용자:' + acl_n[0]])
        acl_d = curs.fetchall()
        if(acl_d):
            if(acl_d[0][0] == 'all'):
                return(0)

            if(acl_d[0][0] == 'user' and not re.search("(\.|:)", ip)):
                return(0)

            if(not ip == acl_n[0] or re.search("(\.|:)", ip)):
                return(1)
        
        if(ip == acl_n[0] and not re.search("(\.|:)", ip) and not re.search("(\.|:)", acl_n[0])):
            return(0)
        else:
            return(1)

    file_c = re.search("^파일:(.*)", name)
    if(file_c and admin_check(5, 'edit (' + name + ')') != 1):
        return(1)

    curs.execute("select acl from data where title = ?", [name])
    acl_d = curs.fetchall()
    if(not acl_d):
        return(0)

    curs.execute("select acl from user where id = ?", [ip])
    user_d = curs.fetchall()

    curs.execute('select data from other where name = "edit"')
    set_d = curs.fetchall()
    if(acl_d[0][0] == 'user'):
        if(not user_d):
            return(1)

    if(acl_d[0][0] == 'admin'):
        if(not user_d):
            return(1)

        if(not admin_check(5, 'edit (' + name + ')') == 1):
            return(1)

    if(set_d):
        if(set_d[0][0] == 'user'):
            if(not user_d):
                return(1)

        if(set_d[0][0] == 'admin'):
            if(not user_d):
                return(1)

            if(not admin_check(5, None) == 1):
                return(1)

    return(0)

def ban_check():
    ip = ip_check()
    band = re.search("^([0-9]{1,3}\.[0-9]{1,3})", ip)
    if(band):
        band_it = band.groups()
    else:
        band_it = ['Not']
        
    curs.execute("select block from ban where block = ? and band = 'O'", [band_it[0]])
    band_d = curs.fetchall()

    curs.execute("select block from ban where block = ?", [ip])
    ban_d = curs.fetchall()
    if(band_d or ban_d):
        return(1)
    
    return(0)
        
def topic_check(name, sub):
    ip = ip_check()
    band = re.search("^([0-9]{1,3}\.[0-9]{1,3})", ip)
    if(band):
        band_it = band.groups()
    else:
        band_it = ['Not']
        
    curs.execute("select block from ban where block = ? and band = 'O'", [band_it[0]])
    band_d = curs.fetchall()

    curs.execute("select block from ban where block = ?", [ip])
    ban_d = curs.fetchall()
    if(band_d or ban_d):
        return(1)

    curs.execute("select title from stop where title = ? and sub = ?", [name, sub])
    topic_s = curs.fetchall()
    if(topic_s):
        return(1)

    return(0)

def rd_plus(title, sub, date):
    curs.execute("select title from rd where title = ? and sub = ?", [title, sub])
    rd = curs.fetchall()
    if(rd):
        curs.execute("update rd set date = ? where title = ? and sub = ?", [date, title, sub])
    else:
        curs.execute("insert into rd (title, sub, date) values (?, ?, ?)", [title, sub, date])
    conn.commit()
    
def rb_plus(block, end, today, blocker, why):
    curs.execute("insert into rb (block, end, today, blocker, why) values (?, ?, ?, ?, ?)", [block, end, today, blocker, why])
    conn.commit()

def history_plus(title, data, date, ip, send, leng):
    curs.execute("select id from history where title = ? order by id+0 desc limit 1", [title])
    rows = curs.fetchall()
    if(rows):
        curs.execute("insert into history (id, title, data, date, ip, send, leng) values (?, ?, ?, ?, ?, ?, ?)", [str(int(rows[0][0]) + 1), title, data, date, ip, send, leng])
    else:
        curs.execute("insert into history (id, title, data, date, ip, send, leng) values ('1', ?, ?, ?, ?, ?, ?)", [title, data, date, ip, send + ' (새 문서)', leng])
    conn.commit()

def leng_check(a, b):
    if(a < b):
        c = b - a
        c = '+' + str(c)
    elif(b < a):
        c = a - b
        c = '-' + str(c)
    else:
        c = '0'
        
    return(c)

def redirect(data):
    return('<meta http-equiv="refresh" content="0;url=' + data + '" />')

def re_error(data):
    if(data == '/ban'):
        ip = ip_check()
        end = '권한이 맞지 않는 상태 입니다.'
        if(ban_check() == 1):
            curs.execute("select end, why from ban where block = ?", [ip])
            d = curs.fetchall()
            if(not d):
                m = re.search("^([0-9]{1,3}\.[0-9]{1,3})", ip)
                if(m):
                    curs.execute("select end, why from ban where block = ? and band = 'O'", [m.groups()[0]])
                    d = curs.fetchall()

            if(d):
                if(d[0][0]):
                    end = d[0][0] + ' 까지 차단 상태 입니다. / 사유 : ' + d[0][1]                

                    now = re.sub(':', '', get_time())
                    now = re.sub('\-', '', now)
                    now = int(re.sub(' ', '', now))
                    
                    day = re.sub('\-', '', d[0][0])    
                    
                    if(now >= int(day + '000000')):
                        curs.execute("delete from ban where block = ?", [ip])
                        conn.commit()
                        
                        end = '차단이 풀렸습니다. 다시 시도 해 보세요.'
                else:
                    end = '영구 차단 상태 입니다. / 사유 : ' + d[0][1]
            

        return(
            html_minify(
                template('index', 
                    imp = ['권한 오류', wiki_set(1), custom(), other2([0, 0])],
                    data = end,
                    menu = 0
                )
            )
        )

    d = re.search('\/error\/([0-9]+)', data)
    if(d):
        num = int(d.groups()[0])
        if(num == 1):
            title = '권한 오류'
            data = '비 로그인 상태 입니다.'
        elif(num == 2):
            title = '권한 오류'
            data = '이 계정이 없습니다.'
        elif(num == 3):
            title = '권한 오류'
            data = '권한이 모자랍니다.'
        elif(num == 4):
            title = '권한 오류'
            data = '관리자는 차단, 검사 할 수 없습니다.'
        elif(num == 5):
            title = '사용자 오류'
            data = '그런 계정이 없습니다.'
        elif(num == 6):
            title = '가입 오류'
            data = '동일한 아이디의 사용자가 있습니다.'
        elif(num == 7):
            title = '가입 오류'
            data = '아이디는 20글자보다 짧아야 합니다.'
        elif(num == 8):
            title = '가입 오류'
            data = '아이디에는 한글과 알파벳과 공백만 허용 됩니다.'
        elif(num == 9):
            title = '파일 올리기 오류'
            data = '파일이 없습니다.'
        elif(num == 10):
            title = '변경 오류'
            data = '비밀번호가 다릅니다.'
        elif(num == 11):
            title = '로그인 오류'
            data = '이미 로그인 되어 있습니다.'
        elif(num == 14):
            title = '파일 올리기 오류'
            data = 'jpg, gif, jpeg, png, webp만 가능 합니다.'
        elif(num == 15):
            title = '편집 오류'
            data = '편집 기록은 500자를 넘을 수 없습니다.'
        elif(num == 16):
            title = '파일 올리기 오류'
            data = '동일한 이름의 파일이 있습니다.'
        elif(num == 17):
            title = '파일 올리기 오류'
            data = '파일 용량은 ' + wiki_set(3) + 'MB를 넘길 수 없습니다.'
        elif(num == 18):
            title = '편집 오류'
            data = '내용이 원래 문서와 동일 합니다.'
        elif(num == 19):
            title = '이동 오류'
            data = '이동 하려는 곳에 문서가 이미 있습니다.'
        elif(num == 20):
            title = '비밀번호 오류'
            data = '재 확인이랑 비밀번호가 다릅니다.'

        if(title):
            return(
                html_minify(
                    template(
                        'index', 
                        imp = [title, wiki_set(1), custom(), other2([0, 0])],
                        data = data,
                        menu = 0
                    )
                )
            )
        else:
            return(redirect('/'))
    else:
        return(redirect('/'))