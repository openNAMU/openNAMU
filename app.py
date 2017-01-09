from flask import Flask, request, session, render_template
app = Flask(__name__)

from urllib import parse
import json
import pymysql
import time
import re

json_data=open('set.json').read()
data = json.loads(json_data)

conn = pymysql.connect(host = data['host'], user = data['user'], password = data['pw'], db = data['db'], charset = 'utf8')
curs = conn.cursor(pymysql.cursors.DictCursor)

def namumark(data):
    data = '\n' + data + '\n'

    data = re.sub('<', '&lt;', data)
    data = re.sub('>', '&gt;', data)

    data = re.sub("======\s?(?P<in>[^=]*)\s?======(?:\s+)?\n", '<h6>\g<in></h6>', data)
    data = re.sub("=====\s?(?P<in>[^=]*)\s?=====(?:\s+)?\n", '<h5>\g<in></h5>', data)
    data = re.sub("====\s?(?P<in>[^=]*)\s?====(?:\s+)?\n", '<h4>\g<in></h4>', data)
    data = re.sub("===\s?(?P<in>[^=]*)\s?===(?:\s+)?\n", '<h3>\g<in></h3>', data)
    data = re.sub("==\s?(?P<in>[^=]*)\s?==(?:\s+)?\n", '<h2>\g<in></h2>', data)
    data = re.sub("=\s?(?P<in>[^=]*)\s?=(?:\s+)?\n", '<h1>\g<in></h1>', data)

    data = re.sub("'''(?P<in>.+?)'''(?!')", '<strong>\g<in></strong>', data)
    data = re.sub("''(?P<in>.+?)''(?!')", '<i>\g<in></i>', data)
    data = re.sub('~~(?P<in>.+?)~~(?!~)', '<s>\g<in></s>', data)
    data = re.sub('--(?P<in>.+?)--(?!-)', '<s>\g<in></s>', data)
    data = re.sub('__(?P<in>.+?)__(?!_)', '<u>\g<in></u>', data)
    data = re.sub('\^\^(?P<in>.+?)\^\^(?!\^)', '<sup>\g<in></sup>', data)
    data = re.sub(',,(?P<in>.+?),,(?!,)', '<sub>\g<in></sub>', data)

    data = re.sub('\n', '<br>', data)
    return data

def getip(request):
    if request.headers.getlist("X-Forwarded-For"):
        ip = request.headers.getlist("X-Forwarded-For")[0]
    else:
        ip = request.remote_addr
    return ip

def getnow():
    now = time.localtime()
    s = "%04d-%02d-%02d %02d:%02d:%02d" % (now.tm_year, now.tm_mon, now.tm_mday, now.tm_hour, now.tm_min, now.tm_sec)
    return s

def recent(title, ip, today, send, leng):
    curs.execute("insert into rc (title, date, ip, send, leng, back) value ('" + pymysql.escape_string(title) + "', '" + today + "', '" + ip + "', '" + pymysql.escape_string(send) + "', '" + leng + "', '')")
    conn.commit()

def discuss(title, sub, ip, date):
    curs.execute("insert into rd (title, sub, ip, date) value ('" + pymysql.escape_string(title) + "', '" + pymysql.escape_string(sub) + "', '" + ip + "', '" + date + "')")
    conn.commit()

def history(title, data, date, ip, send, leng):
    curs.execute("select * from history where title = '" + pymysql.escape_string(title) + "' order by id+0 desc limit 1")
    rows = curs.fetchall()
    if(rows):
        number = int(rows[0]['id']) + 1
        curs.execute("insert into history (id, title, data, date, ip, send, leng) value ('" + str(number) + "', '" + pymysql.escape_string(title) + "', '" + pymysql.escape_string(data) + "', '" + date + "', '" + ip + "', '" + pymysql.escape_string(send) + "', '" + leng + "')")
        conn.commit()
    else:
        curs.execute("insert into history (id, title, data, date, ip, send, leng) value ('1', '" + pymysql.escape_string(title) + "', '" + pymysql.escape_string(data) + "', '" + date + "', '" + ip + "', '" + pymysql.escape_string(send) + "', '" + leng + "')")
        conn.commit()

def getleng(existing, change):
    if(existing < change):
        leng = change - existing
        leng = '+' + str(leng)
    elif(change < existing):
        leng = existing - change
        leng = '-' + str(leng)
    else:
        leng = '0'
    return leng;

@app.route('/')
@app.route('/w/')
def redirect():
    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']) + '" />'

@app.route('/recentchanges')
def recentchanges():
    i = 0
    div = '<div>'
    curs.execute("select * from rc order by date desc limit 50")
    rows = curs.fetchall()
    if(rows):
        while True:
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                break
            if(rows[i]['send']):
                send = rows[i]['send']
                send = re.sub('<', '&lt;', send)
                send = re.sub('>', '&gt;', send)
                send = re.sub('&lt;a href="\/w\/(?P<in>[^"]*)"&gt;(?P<out>[^&]*)&lt;\/a&gt;', '<a href="/w/\g<in>">\g<out></a>', send)
            else:
                send = '<br>'
            div = div + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:33.33%;"><a href="/w/' + parse.quote(rows[i]['title']) + '">' + rows[i]['title'] + '</a> <a href="/history/' + parse.quote(rows[i]['title']) + '">(역사)</a> (' + rows[i]['leng'] + ')</td><td style="text-align: center;width:33.33%;">' + rows[i]['ip'] + '</td><td style="text-align: center;width:33.33%;">' + rows[i]['date'] + '</td></tr><tr><td colspan="3" style="text-align: center;width:100%;">' + send + '</td></tr></tbody></table>'
            i = i + 1
        return render_template('index.html', logo = data['name'], rows = div, tn = 3, title = '최근 변경내역')
    else:
         return render_template('index.html', logo = data['name'], rows = '', tn = 3, title = '최근 변경내역')

@app.route('/recentdiscuss')
def recentdiscuss():
    i = 0
    div = '<div>'
    curs.execute("select * from rd order by date desc limit 50")
    rows = curs.fetchall()
    if(rows):
        while True:
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                break
            div = div + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:33.33%;"><a href="/topic/' + parse.quote(rows[i]['title']) + '/sub/' + parse.quote(rows[i]['sub']) + '">' + rows[i]['title'] + '</a> (' + rows[i]['sub'] + ')</td><td style="text-align: center;width:33.33%;">' + rows[i]['ip'] + '</td><td style="text-align: center;width:33.33%;">' + rows[i]['date'] + '</td></tr></tbody></table>'
            i = i + 1
        return render_template('index.html', logo = data['name'], rows = div, tn = 12, title = '최근 토론내역')
    else:
         return render_template('index.html', logo = data['name'], rows = '', tn = 12, title = '최근 토론내역')

@app.route('/history/<name>')
def gethistory(name = None):
    i = 0
    div = '<div>'
    curs.execute("select * from history where title = '" + pymysql.escape_string(name) + "' order by date desc")
    rows = curs.fetchall()
    if(rows):
        while True:
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                break
            if(rows[i]['send']):
                send = rows[i]['send']
                send = re.sub('<', '&lt;', send)
                send = re.sub('>', '&gt;', send)
                send = re.sub('&lt;a href="\/w\/(?P<in>[^"]*)"&gt;(?P<out>[^&]*)&lt;\/a&gt;', '<a href="/w/\g<in>">\g<out></a>', send)
            else:
                send = '<br>'
            div = div + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:33.33%;">r' + rows[i]['id'] + '</a> <a href="/w/' + parse.quote(rows[i]['title']) + '/r/' + rows[i]['id'] + '">(w)</a> <a href="/w/' + parse.quote(rows[i]['title']) + '/raw/' + rows[i]['id'] + '">(raw)</a> <a href="/revert/' + parse.quote(rows[i]['title']) + '/r/' + rows[i]['id'] + '">(되돌리기)</a> (' + rows[i]['leng'] + ')</td><td style="text-align: center;width:33.33%;">' + rows[i]['ip'] + '</td><td style="text-align: center;width:33.33%;">' + rows[i]['date'] + '</td></tr><tr><td colspan="3" style="text-align: center;width:100%;">' + send + '</td></tr></tbody></table>'
            i = i + 1
        return render_template('index.html', logo = data['name'], rows = div, tn = 5, title = name, page = parse.quote(name))
    else:
         return render_template('index.html', logo = data['name'], rows = '', tn = 5, title = name, page = parse.quote(name))

@app.route('/search', methods=['POST', 'GET'])
def search():
    if(request.method == 'POST'):
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(request.form["search"]) + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']) + '" />'

@app.route('/w/<name>')
def w(name = None):
    curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
    rows = curs.fetchall()
    if(rows):
        enddata = namumark(rows[0]['data'])
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = enddata, license = data['license'], tn = 1)
    else:
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = '<br>문서 없음', license = data['license'], tn = 1)

@app.route('/w/<name>/r/<number>')
def rew(name = None, number = None):
    curs.execute("select * from history where title = '" + pymysql.escape_string(name) + "' and id = '" + number + "'")
    rows = curs.fetchall()
    if(rows):
        enddata = namumark(rows[0]['data'])
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = enddata, license = data['license'], tn = 6)
    else:
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = '<br>문서 없음', license = data['license'], tn = 6)

@app.route('/w/<name>/raw/<number>')
def reraw(name = None, number = None):
    curs.execute("select * from history where title = '" + pymysql.escape_string(name) + "' and id = '" + number + "'")
    rows = curs.fetchall()
    if(rows):
        enddata = re.sub("\n", '<br>', rows[0]['data'])
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = enddata, license = data['license'])
    else:
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = '<br>문서 없음', license = data['license'])

@app.route('/raw/<name>')
def raw(name = None):
    curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
    rows = curs.fetchall()
    if(rows):
        enddata = re.sub("\n", '<br>', rows[0]['data'])
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = enddata, license = data['license'], tn = 7)
    else:
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = '문서 없음', license = data['license'], tn = 7)

@app.route('/revert/<name>/r/<number>', methods=['POST', 'GET'])
def revert(name = None, number = None):
    if(request.method == 'POST'):
        curs.execute("select * from history where title = '" + pymysql.escape_string(name) + "' and id = '" + number + "'")
        rows = curs.fetchall()
        if(rows):
            ip = getip(request)
            today = getnow()
            curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
            row = curs.fetchall()
            if(row):
                leng = getleng(len(row[0]['data']), len(rows[0]['data']))
                curs.execute("update data set data = '" + pymysql.escape_string(rows[0]['data']) + "' where title = '" + pymysql.escape_string(name) + "'")
                conn.commit()
            else:
                leng = '+' + str(len(rows[0]['data']))
                curs.execute("insert into data (title, data, acl) value ('" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(rows[0]['data']) + "', '')")
                conn.commit()
            recent(name, ip, today, '문서를 ' + number + '판으로 되돌렸습니다.', leng)
            history(name, rows[0]['data'], today, ip, '문서를 ' + number + '판으로 되돌렸습니다.', leng)
            return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name) + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name) + '" />'
    else:
        curs.execute("select * from history where title = '" + pymysql.escape_string(name) + "' and id = '" + number + "'")
        rows = curs.fetchall()
        if(rows):
            return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), r = parse.quote(number), tn = 13, plus = '정말 되돌리시겠습니까?')
        else:
            return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name) + '" />'

@app.route('/edit/<name>', methods=['POST', 'GET'])
def edit(name = None):
    if(request.method == 'POST'):
        curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
        rows = curs.fetchall()
        if(rows):
            ip = getip(request)
            today = getnow()
            leng = getleng(len(rows[0]['data']), len(request.form["content"]))
            recent(name, ip, today, request.form["send"], leng)
            history(name, request.form["content"], today, ip, request.form["send"], leng)
            curs.execute("update data set data = '" + pymysql.escape_string(request.form["content"]) + "' where title = '" + pymysql.escape_string(name) + "'")
            conn.commit()
        else:
            ip = getip(request)
            today = getnow()
            leng = '+' + str(len(request.form["content"]))
            recent(name, ip, today, request.form["send"], leng)
            history(name, request.form["content"], today, ip, request.form["send"], leng)
            curs.execute("insert into data (title, data, acl) value ('" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(request.form["content"]) + "', '')")
            conn.commit()
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name) + '" />'
    else:
        curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
        rows = curs.fetchall()
        if(rows):
            return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = rows[0]['data'], tn = 2)
        else:
            return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = '', tn = 2)

@app.route('/delete/<name>', methods=['POST', 'GET'])
def delete(name = None):
    if(request.method == 'POST'):
        curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
        rows = curs.fetchall()
        if(rows):
            ip = getip(request)
            today = getnow()
            leng = '-' + str(len(rows[0]['data']))
            recent(name, ip, today, '문서를 삭제 했습니다.', leng)
            history(name, '', today, ip, '문서를 삭제 했습니다.', leng)
            curs.execute("delete from data where title = '" + pymysql.escape_string(name) + "'")
            conn.commit()
            return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name) + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name) + '" />'
    else:
        curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
        rows = curs.fetchall()
        if(rows):
            return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), tn = 8, plus = '정말 삭제 하시겠습니까?')
        else:
            return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name) + '" />'

@app.route('/move/<name>', methods=['POST', 'GET'])
def move(name = None):
    if(request.method == 'POST'):
        curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
        rows = curs.fetchall()
        if(rows):
            ip = getip(request)
            today = getnow()
            leng = '0'
            curs.execute("select * from history where title = '" + pymysql.escape_string(request.form["title"]) + "'")
            row = curs.fetchall()
            if(row):
                 return render_template('index.html', title = '이동 오류', logo = data['name'], data = '이동 하려는 곳에 문서가 이미 있습니다.')
            else:
                recent(name, ip, today, '문서를 <a href="/w/' + pymysql.escape_string(parse.quote(request.form["title"])) + '">' + pymysql.escape_string(request.form["title"]) + '</a> 문서로 이동 했습니다.', leng)
                history(name, rows[0]['data'], today, ip, '<a href="/w/' + pymysql.escape_string(parse.quote(name)) + '">' + pymysql.escape_string(name) + '</a> 문서를 <a href="/w/' + pymysql.escape_string(parse.quote(request.form["title"])) + '">' + pymysql.escape_string(request.form["title"]) + '</a> 문서로 이동 했습니다.', leng)
                curs.execute("update data set title = '" + pymysql.escape_string(request.form["title"]) + "' where title = '" + pymysql.escape_string(name) + "'")
                curs.execute("update history set title = '" + pymysql.escape_string(request.form["title"]) + "' where title = '" + pymysql.escape_string(name) + "'")
                conn.commit()
                return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(request.form["title"]) + '" />'
        else:
            ip = getip(request)
            today = getnow()
            leng = '0'
            curs.execute("select * from history where title = '" + pymysql.escape_string(request.form["title"]) + "'")
            row = curs.fetchall()
            if(row):
                 return render_template('index.html', title = '이동 오류', logo = data['name'], data = '이동 하려는 곳에 문서가 이미 있습니다.')
            else:
                recent(name, ip, today, '문서를 <a href="/w/' + pymysql.escape_string(parse.quote(request.form["title"])) + '">' + pymysql.escape_string(request.form["title"]) + '</a> 문서로 이동 했습니다.', leng)
                history(name, rows[0]['data'], today, ip, '<a href="/w/' + pymysql.escape_string(parse.quote(name)) + '">' + pymysql.escape_string(name) + '</a> 문서를 <a href="/w/' + pymysql.escape_string(parse.quote(request.form["title"])) + '">' + pymysql.escape_string(request.form["title"]) + '</a> 문서로 이동 했습니다.', leng)
                curs.execute("update history set title = '" + pymysql.escape_string(request.form["title"]) + "' where title = '" + pymysql.escape_string(name) + "'")
                conn.commit()
                return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(request.form["title"]) + '" />'
    else:
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), tn = 9, plus = '정말 이동 하시겠습니까?')

@app.route('/setup')
def setup():
    curs.execute("create table if not exists data(title text not null, data longtext not null, acl text not null)")
    curs.execute("create table if not exists history(id text not null, title text not null, data longtext not null, date text not null, ip text not null, send text not null, leng text not null)")
    curs.execute("create table if not exists rc(title text not null, date text not null, ip text not null, send text not null, leng text not null, back text not null)")
    curs.execute("create table if not exists rd(title text not null, sub text not null, date text not null, ip text not null)")
    curs.execute("create table if not exists user(id text not null, pw text not null, acl text not null)")
    curs.execute("create table if not exists ban(block text not null, end text not null, why text not null, band text not null)")
    curs.execute("create table if not exists topic(id text not null, title text not null, sub text not null, data longtext not null, date text not null, ip text not null, block text not null)")
    return render_template('index.html', title = '설치 완료', logo = data['name'], data = '문제 없었음')

@app.route('/other')
def other():
    return render_template('index.html', title = '기타 메뉴', logo = data['name'], data = '<li><a href="/titleindex">모든 문서</a><li><a href="/grammar">문법 설명</a></li><li><a href="/version">버전</a></li>')

@app.route('/titleindex')
def titleindex():
    i = 0
    div = '<div>'
    curs.execute("select * from data")
    rows = curs.fetchall()
    if(rows):
        while True:
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                break
            div = div + '<li><a href="/w/' + parse.quote(rows[i]['title']) + '">' + rows[i]['title'] + '</a></li>'
            i = i + 1
        return render_template('index.html', logo = data['name'], rows = div, tn = 4, title = '모든 문서')
    else:
        return render_template('index.html', logo = data['name'], rows = '', tn = 4, title = '모든 문서')

@app.route('/topic/<name>', methods=['POST', 'GET'])
def topic(name = None):
    if(request.method == 'POST'):
        return '<meta http-equiv="refresh" content="0;url=/topic/' + parse.quote(name) + '/sub/' + parse.quote(request.form["topic"]) + '" />'
    else:
        div = '<div>'
        i = 0
        curs.execute("select * from topic where title = '" + pymysql.escape_string(name) + "'")
        rows = curs.fetchall()
        while True:
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                break
            if(i == 0):
                sub = rows[i]['sub']
                div = div + '<li><a href="/topic/' + parse.quote(name) + '/sub/' + parse.quote(rows[i]['sub']) + '">' + rows[i]['sub'] + '</a></li>'
            else:
                if(not sub == rows[i]['sub']):
                    sub = rows[i]['sub']
                    div = div + '<li><a href="/topic/' + parse.quote(name) + '/sub/' + parse.quote(rows[i]['sub']) + '">' + rows[i]['sub'] + '</a></li>'
            i = i + 1
        return render_template('index.html', title = name, page = parse.quote(name), logo = data['name'], plus = div, tn = 10)

@app.route('/topic/<name>/sub/<sub>', methods=['POST', 'GET'])
def sub(name = None, sub = None):
    if(request.method == 'POST'):
        curs.execute("select * from topic where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' order by id+0 desc limit 1")
        rows = curs.fetchall()
        if(rows):
            number = int(rows[0]['id']) + 1
        else:
            number = 1
        ip = getip(request)
        today = getnow()
        discuss(name, sub, ip, today)
        curs.execute("insert into topic (id, title, sub, data, date, ip, block) value ('" + str(number) + "', '" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(sub) + "', '" + pymysql.escape_string(request.form["content"]) + "', '" + today + "', '" + ip + "', '')")
        conn.commit()
        return '<meta http-equiv="refresh" content="0;url=/topic/' + parse.quote(name) + '/sub/' + parse.quote(sub) + '" />'
    else:
        div = '<div>'
        i = 0
        curs.execute("select * from topic where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' order by id+0 asc")
        rows = curs.fetchall()
        while True:
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                break
            if(i == 0):
                start = rows[i]['ip']
            if(rows[i]['ip'] == start):
                j = i + 1
                div = div + '<table id="toron"><tbody><tr><td id="toroncolorgreen"><a href="javascript:void(0);" id="' + str(j) + '">#' + str(j) + '</a> ' + rows[i]['ip'] + ' <span style="float:right;">' + rows[i]['date'] + '</span></td></tr><tr><td>' + rows[i]['data'] + '</td></tr></tbody></table><br>'
            else:
                j = i + 1
                div = div + '<table id="toron"><tbody><tr><td id="toroncolor"><a href="javascript:void(0);" id="' + str(j) + '">#' + str(j) + '</a> ' + rows[i]['ip'] + ' <span style="float:right;">' + rows[i]['date'] + '</span></td></tr><tr><td>' + rows[i]['data'] + '</td></tr></tbody></table><br>'
            i = i + 1
        return render_template('index.html', title = name, page = parse.quote(name), suburl = parse.quote(sub), sub = sub, logo = data['name'], rows = div, tn = 11)

@app.route('/grammar')
def grammar():
    return render_template('index.html', title = '문법 설명', logo = data['name'], data = '아직 없음')

@app.route('/version')
def version():
    return render_template('index.html', title = '버전', logo = data['name'], tn = 14)

@app.route('/user')
def user():
    ip = getip(request)
    return render_template('index.html', title = '유저 메뉴', logo = data['name'], data = ip + '<br><br><li><a href="/login">로그인</a></li><li><a href="/logout">로그아웃</a></li>')

@app.route('/random')
def random():
    curs.execute("select * from data order by rand() limit 1")
    rows = curs.fetchall()
    if(rows):
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(rows[0]['title']) + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']) + '" />'

if __name__ == '__main__':
    app.run(host = '0.0.0.0', port = 3000)
