from flask import Flask, request, session, render_template
app = Flask(__name__)

from urllib import parse
import json
import pymysql
import time
import re
import bcrypt

json_data=open('set.json').read()
data = json.loads(json_data)

conn = pymysql.connect(host = data['host'], user = data['user'], password = data['pw'], db = data['db'], charset = 'utf8')
curs = conn.cursor(pymysql.cursors.DictCursor)

app.secret_key = data['key']

def namumark(title, data):
    data = '\n' + data + '\n'

    data = re.sub('<', '&lt;', data)
    data = re.sub('>', '&gt;', data)
    data = re.sub('"', '&quot;', data)

    h0c = 0;
    h1c = 0;
    h2c = 0;
    h3c = 0;
    h4c = 0;
    h5c = 0;
    last = 0;
    rtoc = '<div id="toc"><span id="toc-name">목차</span><br><br>'
    while True:
        m = re.search('(={1,6})\s?([^=]*)\s?(?:={1,6})(?:\s+)?\n', data)
        if(m):
            result = m.groups()
            wiki = len(result[0])
            if(last < wiki):
                last = wiki
            else:
                last = wiki;
                if(wiki == 1):
                    h1c = 0
                    h2c = 0
                    h3c = 0
                    h4c = 0
                    h5c = 0
                elif(wiki == 2):
                    h2c = 0
                    h3c = 0
                    h4c = 0
                    h5c = 0
                elif(wiki == 3):
                    h3c = 0
                    h4c = 0
                    h5c = 0
                elif(wiki == 4):
                    h4c = 0
                    h5c = 0
                elif(wiki == 5):
                    h5c = 0
            if(wiki == 1):
                h0c = h0c + 1
            elif(wiki == 2):
                h1c = h1c + 1
            elif(wiki == 3):
                h2c = h2c + 1
            elif(wiki == 4):
                h3c = h3c + 1
            elif(wiki == 5):
                h4c = h4c + 1
            else:
                h5c = h5c + 1
            toc = str(h0c) + '.' + str(h1c) + '.' + str(h2c) + '.' + str(h3c) + '.' + str(h4c) + '.' + str(h5c) + '.'
            toc = re.sub("(?P<in>[0-9]0(?:[0]*)?)\.", '\g<in>#.', toc)
            toc = re.sub("0\.", '', toc)
            toc = re.sub("#\.", '.', toc)
            toc = re.sub("\.$", '', toc)
            rtoc = rtoc + '<a href="#s-' + toc + '">' + toc + '</a>. ' + result[1] + '<br>'
            data = re.sub('(={1,6})\s?([^=]*)\s?(?:={1,6})(?:\s+)?\r\n', '<h' + str(wiki) + '><a href="#toc" id="s-' + toc + '">' + toc + '.</a> ' + result[1] + '</h' + str(wiki) + '>', data, 1);
        else:
            rtoc = rtoc + '</div>'
            break
    
    data = re.sub("\[목차\]", rtoc, data)

    data = re.sub("'''(?P<in>.+?)'''(?!')", '<b>\g<in></b>', data)
    data = re.sub("''(?P<in>.+?)''(?!')", '<i>\g<in></i>', data)
    data = re.sub('~~(?P<in>.+?)~~(?!~)', '<s>\g<in></s>', data)
    data = re.sub('--(?P<in>.+?)--(?!-)', '<s>\g<in></s>', data)
    data = re.sub('__(?P<in>.+?)__(?!_)', '<u>\g<in></u>', data)
    data = re.sub('\^\^(?P<in>.+?)\^\^(?!\^)', '<sup>\g<in></sup>', data)
    data = re.sub(',,(?P<in>.+?),,(?!,)', '<sub>\g<in></sub>', data)
    
    while True:
        p = re.compile("\[youtube\(((?:(?!,|\)\]).)*)(?:,\s)?(?:width=((?:(?!,|\)\]).)*))?(?:,\s)?(?:height=((?:(?!,|\)\]).)*))?(?:,\s)?(?:width=((?:(?!,|\)\]).)*))?\)\]", re.I)
        m = p.search(data)
        if(m):
            result = m.groups()
            if(result[1]):
                if(result[2]):
                    width = result[1]
                    height = result[2]
                else:
                    width = result[1]
                    height = '315'
            elif(result[2]):
                if(result[3]):
                    height = result[2]
                    width = result[3]
                else:
                    height = result[2]
                    width = '560'
            else:
                width = '560'
                height = '315'
            data = p.sub('<iframe width="' + width + '" height="' + height + '" src="https://www.youtube.com/embed/' + result[0] + '" frameborder="0" allowfullscreen></iframe>', data, 1)
        else:
            break
                

    while True:
        m = re.search("\[\[(((?!\]\]).)*)\]\]", data)
        if(m):
            result = m.groups()
            a = re.search("(((?!\|).)*)\|(.*)", result[0])
            if(a):
                results = a.groups()
                p = re.compile("^http(?:s)?:\/\/", re.I)
                b = p.search(results[0])
                if(b):
                    data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a class="out_link" href="' + results[0] + '">' + results[2] + '</a>', data, 1)
                else:
                    if(results[0] == title):
                        data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<b>' + results[2] + '</b>', data, 1)
                    else:
                        curs.execute("select * from data where title = '" + pymysql.escape_string(results[0]) + "'")
                        rows = curs.fetchall()
                        if(rows):
                            data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a title="' + results[0] + '" href="/w/' + parse.quote(results[0]) + '">' + results[2] + '</a>', data, 1)
                        else:
                            data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a title="' + results[0] + '" class="not_thing" href="/w/' + parse.quote(results[0]) + '">' + results[2] + '</a>', data, 1)
            else:
                b = re.search("^[Hh][Tt][Tt][Pp]([Ss])?:\/\/", result[0])
                if(b):
                    data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a class="out_link" href="' + result[0] + '">' + result[0] + '</a>', data, 1)
                else:
                    if(result[0] == title):
                        data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<b>' + result[0] + '</b>', data, 1)
                    else:
                        curs.execute("select * from data where title = '" + pymysql.escape_string(result[0]) + "'")
                        rows = curs.fetchall()
                        if(rows):
                            data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a href="/w/' + parse.quote(result[0]) + '">' + result[0] + '</a>', data, 1)
                        else:
                            data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a class="not_thing" href="/w/' + parse.quote(result[0]) + '">' + result[0] + '</a>', data, 1)
        else:
            break

    data = re.sub('\n', '<br>', data)
    return data

def getip(request):
    if(session.get('Now') == True):
        ip = format(session['DREAMER'])
    else:
        if(request.headers.getlist("X-Forwarded-For")):
            ip = request.headers.getlist("X-Forwarded-For")[0]
        else:
            ip = request.remote_addr
    return ip

def getcan(ip, name):
    curs.execute("select * from ban where block = '" + pymysql.escape_string(ip) + "'")
    rows = curs.fetchall()
    if(rows):
        return 1
    else:
        curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
        row = curs.fetchall()
        if(row):
            curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
            rows = curs.fetchall()
            if(row[0]['acl'] == 'user'):
                if(rows):
                    return 0
                else:
                    return 1
            elif(row[0]['acl'] == 'admin'):
                if(rows):
                    if(rows[0]['acl'] == 'admin' or rows[0]['acl'] == 'owner'):
                        return 0
                    else:
                        return 1
                else:
                    return 1
            else:
                return 0
        else:
            return 0

def getban(ip):
    curs.execute("select * from ban where block = '" + pymysql.escape_string(ip) + "'")
    rows = curs.fetchall()
    if(rows):
        return 1
    else:
        return 0
        
def getdiscuss(ip, name, sub):
    curs.execute("select * from ban where block = '" + pymysql.escape_string(ip) + "'")
    rows = curs.fetchall()
    if(rows):
        return 1
    else:
        curs.execute("select * from stop where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "'")
        rows = curs.fetchall()
        if(rows):
            return 1
        else:
            return 0

def getnow():
    now = time.localtime()
    s = "%04d-%02d-%02d %02d:%02d:%02d" % (now.tm_year, now.tm_mon, now.tm_mday, now.tm_hour, now.tm_min, now.tm_sec)
    return s

def recent(title, ip, today, send, leng):
    curs.execute("insert into rc (title, date, ip, send, leng, back) value ('" + pymysql.escape_string(title) + "', '" + today + "', '" + ip + "', '" + pymysql.escape_string(send) + "', '" + leng + "', '')")
    conn.commit()

def discuss(title, sub, date):
    curs.execute("select * from rd where title = '" + pymysql.escape_string(title) + "' and sub = '" + pymysql.escape_string(sub) + "'")
    rows = curs.fetchall()
    if(rows):
        curs.execute("update rd set date = '" + pymysql.escape_string(date) + "' where title = '" + pymysql.escape_string(title) + "' and sub = '" + pymysql.escape_string(sub) + "'")
    else:
        curs.execute("insert into rd (title, sub, date) value ('" + pymysql.escape_string(title) + "', '" + pymysql.escape_string(sub) + "', '" + pymysql.escape_string(date) + "')")
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
            title = rows[i]['title']
            title = re.sub('<', '&lt;', title)
            title = re.sub('>', '&gt;', title)
            div = div + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:33.33%;"><a href="/w/' + parse.quote(rows[i]['title']) + '">' + title + '</a> <a href="/history/' + parse.quote(rows[i]['title']) + '">(역사)</a> (' + rows[i]['leng'] + ')</td><td style="text-align: center;width:33.33%;">' + rows[i]['ip'] + '</td><td style="text-align: center;width:33.33%;">' + rows[i]['date'] + '</td></tr><tr><td colspan="3" style="text-align: center;width:100%;">' + send + '</td></tr></tbody></table>'
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
            title = rows[i]['title']
            title = re.sub('<', '&lt;', title)
            title = re.sub('>', '&gt;', title)
            sub = rows[i]['sub']
            sub = re.sub('<', '&lt;', sub)
            sub = re.sub('>', '&gt;', sub)
            div = div + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:50%;"><a href="/topic/' + parse.quote(rows[i]['title']) + '/sub/' + parse.quote(rows[i]['sub']) + '">' + title + '</a> (' + sub + ')</td><td style="text-align: center;width:50%;">' + rows[i]['date'] + '</td></tr></tbody></table>'
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
            div = div + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:33.33%;">r' + rows[i]['id'] + '</a> <a href="/w/' + parse.quote(rows[i]['title']) + '/r/' + rows[i]['id'] + '">(w)</a> <a href="/w/' + parse.quote(rows[i]['title']) + '/raw/' + rows[i]['id'] + '">(Raw)</a> <a href="/revert/' + parse.quote(rows[i]['title']) + '/r/' + rows[i]['id'] + '">(되돌리기)</a> (' + rows[i]['leng'] + ')</td><td style="text-align: center;width:33.33%;">' + rows[i]['ip'] + '</td><td style="text-align: center;width:33.33%;">' + rows[i]['date'] + '</td></tr><tr><td colspan="3" style="text-align: center;width:100%;">' + send + '</td></tr></tbody></table>'
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
        if(rows[0]['acl'] == 'admin'):
            acl = '(관리자)'
        elif(rows[0]['acl'] == 'user'):
            acl = '(유저)'
        else:
            acl = ''
        enddata = namumark(name, rows[0]['data'])
        m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
        if(m):
            result = m.groups()
            left = result[0]
        else:
            left = ''
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = enddata, license = data['license'], tn = 1, acl = acl, left = left)
    else:
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = '<br>문서 없음', license = data['license'], tn = 1)

@app.route('/w/<name>/redirect/<redirect>')
def redirectw(name = None, redirect = None):
    curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
    rows = curs.fetchall()
    if(rows):
        enddata = namumark(name, rows[0]['data'])
        m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
        if(m):
            result = m.groups()
            left = result[0]
        else:
            left = ''
        test = redirect
        redirect = re.sub('<', '&lt;', redirect)
        redirect = re.sub('>', '&gt;', redirect)
        redirect = re.sub('"', '&quot;', redirect)
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = enddata, license = data['license'], tn = 1, redirect = '<a href="/w/' + parse.quote(test) + '">' + redirect + '</a>에서 넘어 왔습니다.', left = left)
    else:
        test = redirect
        redirect = re.sub('<', '&lt;', redirect)
        redirect = re.sub('>', '&gt;', redirect)
        redirect = re.sub('"', '&quot;', redirect)
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = '<br>문서 없음', license = data['license'], tn = 1, redirect = '<a href="/w/' + parse.quote(test) + '">' + redirect + '</a>에서 넘어 왔습니다.')

@app.route('/w/<name>/r/<number>')
def rew(name = None, number = None):
    curs.execute("select * from history where title = '" + pymysql.escape_string(name) + "' and id = '" + number + "'")
    rows = curs.fetchall()
    if(rows):
        enddata = namumark(name, rows[0]['data'])
        m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
        if(m):
            result = m.groups()
            left = result[0]
        else:
            left = ''
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = enddata, license = data['license'], tn = 6, left = left)
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
            can = getcan(ip, name)
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
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
        ip = getip(request)
        can = getcan(ip, name)
        if(can == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
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
            can = getcan(ip, name)
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                today = getnow()
                leng = getleng(len(rows[0]['data']), len(request.form["content"]))
                recent(name, ip, today, request.form["send"], leng)
                history(name, request.form["content"], today, ip, request.form["send"], leng)
                curs.execute("update data set data = '" + pymysql.escape_string(request.form["content"]) + "' where title = '" + pymysql.escape_string(name) + "'")
                conn.commit()
        else:
            ip = getip(request)
            can = getcan(ip, name)
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                today = getnow()
                leng = '+' + str(len(request.form["content"]))
                recent(name, ip, today, request.form["send"], leng)
                history(name, request.form["content"], today, ip, request.form["send"], leng)
                curs.execute("insert into data (title, data, acl) value ('" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(request.form["content"]) + "', '')")
                conn.commit()
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name) + '" />'
    else:
        ip = getip(request)
        can = getcan(ip, name)
        if(can == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
            rows = curs.fetchall()
            if(rows):
                return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = rows[0]['data'], tn = 2)
            else:
                return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = '', tn = 2)
                
@app.route('/preview/<name>', methods=['POST'])
def preview(name = None):
    ip = getip(request)
    can = getcan(ip, name)
    if(can == 1):
        return '<meta http-equiv="refresh" content="0;url=/ban" />'
    else:
        enddata = namumark(name, request.form["content"])
        m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
        if(m):
            result = m.groups()
            left = result[0]
        else:
            left = ''
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = request.form["content"], tn = 2, preview = 1, enddata = enddata, left = left)

@app.route('/delete/<name>', methods=['POST', 'GET'])
def delete(name = None):
    if(request.method == 'POST'):
        curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
        rows = curs.fetchall()
        if(rows):
            ip = getip(request)
            can = getcan(ip, name)
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
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
            ip = getip(request)
            can = getcan(ip, name)
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
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
            can = getcan(ip, name)
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
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
            can = getcan(ip, name)
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
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
        ip = getip(request)
        can = getcan(ip, name)
        if(can == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), tn = 9, plus = '정말 이동 하시겠습니까?')

@app.route('/setup')
def setup():
    curs.execute("create table if not exists data(title text not null, data longtext not null, acl text not null)")
    curs.execute("create table if not exists history(id text not null, title text not null, data longtext not null, date text not null, ip text not null, send text not null, leng text not null)")
    curs.execute("create table if not exists rc(title text not null, date text not null, ip text not null, send text not null, leng text not null, back text not null)")
    curs.execute("create table if not exists rd(title text not null, sub text not null, date text not null)")
    curs.execute("create table if not exists user(id text not null, pw text not null, acl text not null)")
    curs.execute("create table if not exists ban(block text not null, end text not null, why text not null, band text not null)")
    curs.execute("create table if not exists topic(id text not null, title text not null, sub text not null, data longtext not null, date text not null, ip text not null, block text not null)")
    curs.execute("create table if not exists stop(title text not null, sub text not null, close text not null)")
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
        curs.execute("select * from topic where title = '" + pymysql.escape_string(name) + "' order by sub asc")
        rows = curs.fetchall()
        while True:
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                break
            if(i == 0):
                sub = rows[i]['sub']
                curs.execute("select * from stop where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' and close = 'O'")
                row = curs.fetchall()
                if(not row):
                    div = div + '<li><a href="/topic/' + parse.quote(name) + '/sub/' + parse.quote(rows[i]['sub']) + '">' + rows[i]['sub'] + '</a></li>'
            else:
                if(not sub == rows[i]['sub']):
                    sub = rows[i]['sub']
                    curs.execute("select * from stop where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' and close = 'O'")
                    row = curs.fetchall()
                    if(not row):
                        div = div + '<li><a href="/topic/' + parse.quote(name) + '/sub/' + parse.quote(rows[i]['sub']) + '">' + rows[i]['sub'] + '</a></li>'
            i = i + 1
        return render_template('index.html', title = name, page = parse.quote(name), logo = data['name'], plus = div, tn = 10, list = 1)
        
@app.route('/topic/<name>/close')
def topicstoplist(name = None):
    if(request.method == 'POST'):
        return '<meta http-equiv="refresh" content="0;url=/topic/' + parse.quote(name) + '/sub/' + parse.quote(request.form["topic"]) + '" />'
    else:
        div = '<div>'
        i = 0
        curs.execute("select * from stop where title = '" + pymysql.escape_string(name) + "' and close = 'O' order by sub asc")
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
        ban = getdiscuss(ip, name, sub)
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
            rows = curs.fetchall()
            if(rows):
                if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
                    ip = ip + ' - Admin'
            today = getnow()
            discuss(name, sub, today)
            curs.execute("insert into topic (id, title, sub, data, date, ip, block) value ('" + str(number) + "', '" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(sub) + "', '" + pymysql.escape_string(request.form["content"]) + "', '" + today + "', '" + ip + "', '')")
            conn.commit()
            return '<meta http-equiv="refresh" content="0;url=/topic/' + parse.quote(name) + '/sub/' + parse.quote(sub) + '" />'
    else:
        ip = getip(request)
        ban = getdiscuss(ip, name, sub)
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
            indata = rows[i]['data']
            indata = re.sub('<', '&lt;', indata)
            indata = re.sub('>', '&gt;', indata)
            indata = re.sub('"', '&quot;', indata)
            indata = re.sub('\n', '<br>', indata)
            if(rows[i]['block'] == 'O'):
                indata = '블라인드 되었습니다.'
                block = 'style="background: gainsboro;"'
            else:
                block = ''
            if(rows[i]['ip'] == start):
                j = i + 1
                div = div + '<table id="toron"><tbody><tr><td id="toroncolorgreen"><a href="javascript:void(0);" id="' + str(j) + '">#' + str(j) + '</a> ' + rows[i]['ip'] + ' <span style="float:right;">' + rows[i]['date'] + '</span></td></tr><tr><td ' + block + '>' + indata + '</td></tr></tbody></table><br>'
            else:
                j = i + 1
                div = div + '<table id="toron"><tbody><tr><td id="toroncolor"><a href="javascript:void(0);" id="' + str(j) + '">#' + str(j) + '</a> ' + rows[i]['ip'] + ' <span style="float:right;">' + rows[i]['date'] + '</span></td></tr><tr><td ' + block + '>' + indata + '</td></tr></tbody></table><br>'
            i = i + 1
        return render_template('index.html', title = name, page = parse.quote(name), suburl = parse.quote(sub), sub = sub, logo = data['name'], rows = div, tn = 11, ban = ban)

@app.route('/topic/<name>/sub/<sub>/b/<number>')
def blind(name = None, sub = None, number = None):
    if(session.get('Now') == True):
        ip = getip(request)
        curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
        rows = curs.fetchall()
        if(rows):
            if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
                curs.execute("select * from topic where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' and id = '" + number + "'")
                row = curs.fetchall()
                if(row):
                    if(row[0]['block'] == 'O'):
                        curs.execute("update topic set block = '' where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' and id = '" + number + "'")
                    else:
                        curs.execute("update topic set block = 'O' where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' and id = '" + number + "'")
                    conn.commit()
                    return '<meta http-equiv="refresh" content="0;url=/topic/' + name + '/sub/' + sub + '" />'
                else:
                    return '<meta http-equiv="refresh" content="0;url=/topic/' + name + '/sub/' + sub + '" />'
            else:
                return render_template('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.')
        else:
            return render_template('index.html', title = '권한 오류', logo = data['name'], data = '계정이 없습니다.')
    else:
        return render_template('index.html', title = '권한 오류', logo = data['name'], data = '비 로그인 상태 입니다.')
        
@app.route('/topic/<name>/sub/<sub>/stop')
def topicstop(name = None, sub = None):
    if(session.get('Now') == True):
        ip = getip(request)
        curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
        rows = curs.fetchall()
        if(rows):
            if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
                curs.execute("select * from topic where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' order by id+0 desc limit 1")
                row = curs.fetchall()
                if(row):
                    today = getnow()
                    curs.execute("select * from stop where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' and close = ''")
                    rows = curs.fetchall()
                    if(rows):
                        curs.execute("insert into topic (id, title, sub, data, date, ip, block) value ('" + pymysql.escape_string(str(int(row[0]['id']) + 1)) + "', '" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(sub) + "', 'Restart', '" + pymysql.escape_string(today) + "', '" + pymysql.escape_string(ip) + " - Restart', '')")
                        curs.execute("delete from stop where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' and close = ''")
                    else:
                        curs.execute("insert into topic (id, title, sub, data, date, ip, block) value ('" + pymysql.escape_string(str(int(row[0]['id']) + 1)) + "', '" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(sub) + "', 'Stop', '" + pymysql.escape_string(today) + "', '" + pymysql.escape_string(ip) + " - Stop', '')")
                        curs.execute("insert into stop (title, sub, close) value ('" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(sub) + "', '')")
                    conn.commit()
                    return '<meta http-equiv="refresh" content="0;url=/topic/' + name + '/sub/' + sub + '" />'
                else:
                    return '<meta http-equiv="refresh" content="0;url=/topic/' + name + '/sub/' + sub + '" />'
            else:
                return render_template('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.')
        else:
            return render_template('index.html', title = '권한 오류', logo = data['name'], data = '계정이 없습니다.')
    else:
        return render_template('index.html', title = '권한 오류', logo = data['name'], data = '비 로그인 상태 입니다.')
        
@app.route('/topic/<name>/sub/<sub>/close')
def topicclose(name = None, sub = None):
    if(session.get('Now') == True):
        ip = getip(request)
        curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
        rows = curs.fetchall()
        if(rows):
            if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
                curs.execute("select * from topic where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' order by id+0 desc limit 1")
                row = curs.fetchall()
                if(row):
                    today = getnow()
                    curs.execute("select * from stop where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' and close = 'O'")
                    rows = curs.fetchall()
                    if(rows):
                        curs.execute("insert into topic (id, title, sub, data, date, ip, block) value ('" + pymysql.escape_string(str(int(row[0]['id']) + 1)) + "', '" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(sub) + "', 'Reopen', '" + pymysql.escape_string(today) + "', '" + pymysql.escape_string(ip) + " - Reopen', '')")
                        curs.execute("delete from stop where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' and close = 'O'")
                    else:
                        curs.execute("insert into topic (id, title, sub, data, date, ip, block) value ('" + pymysql.escape_string(str(int(row[0]['id']) + 1)) + "', '" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(sub) + "', 'Close', '" + pymysql.escape_string(today) + "', '" + pymysql.escape_string(ip) + " - Close', '')")
                        curs.execute("insert into stop (title, sub, close) value ('" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(sub) + "', 'O')")
                    conn.commit()
                    return '<meta http-equiv="refresh" content="0;url=/topic/' + name + '/sub/' + sub + '" />'
                else:
                    return '<meta http-equiv="refresh" content="0;url=/topic/' + name + '/sub/' + sub + '" />'
            else:
                return render_template('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.')
        else:
            return render_template('index.html', title = '권한 오류', logo = data['name'], data = '계정이 없습니다.')
    else:
        return render_template('index.html', title = '권한 오류', logo = data['name'], data = '비 로그인 상태 입니다.')

@app.route('/login', methods=['POST', 'GET'])
def login():
    if(request.method == 'POST'):
        ip = getip(request)
        ban = getban(ip)
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            curs.execute("select * from user where id = '" + pymysql.escape_string(request.form["id"]) + "'")
            rows = curs.fetchall()
            if(rows):
                if(session.get('Now') == True):
                    return render_template('index.html', title = '로그인 오류', logo = data['name'], data = '이미 로그인 되어 있습니다.')
                elif(bcrypt.checkpw(bytes(request.form["pw"], 'utf-8'), bytes(rows[0]['pw'], 'utf-8'))):
                    session['Now'] = True
                    session['DREAMER'] = request.form["id"]
                    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']) + '" />'
                else:
                    return render_template('index.html', title = '로그인 오류', logo = data['name'], data = '비밀번호가 다릅니다.')
            else:
                return render_template('index.html', title = '로그인 오류', logo = data['name'], data = '없는 계정 입니다.')
    else:
        ip = getip(request)
        ban = getban(ip)
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            if(session.get('Now') == True):
                return render_template('index.html', title = '로그인 오류', logo = data['name'], data = '이미 로그인 되어 있습니다.')
            else:
                return render_template('index.html', title = '로그인', enter = '로그인', logo = data['name'], tn = 15)

@app.route('/register', methods=['POST', 'GET'])
def register():
    if(request.method == 'POST'):
        ip = getip(request)
        ban = getban(ip)
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            p = re.compile('(?:[^A-Za-zㄱ-힣0-9 ])')
            m = p.search(request.form["id"])
            if(m):
                return render_template('index.html', title = '회원가입 오류', logo = data['name'], data = '아이디에는 한글과 알파벳 공백만 허용 됩니다.')
            else:
                curs.execute("select * from user where id = '" + pymysql.escape_string(request.form["id"]) + "'")
                rows = curs.fetchall()
                if(rows):
                    return render_template('index.html', title = '회원가입 오류', logo = data['name'], data = '동일한 아이디의 유저가 있습니다.')
                else:
                    hashed = bcrypt.hashpw(bytes(request.form["pw"], 'utf-8'), bcrypt.gensalt())
                    if(request.form["id"] == data['owner']):
                        curs.execute("insert into user (id, pw, acl) value ('" + pymysql.escape_string(request.form["id"]) + "', '" + pymysql.escape_string(hashed.decode()) + "', 'owner')")
                    else:
                        curs.execute("insert into user (id, pw, acl) value ('" + pymysql.escape_string(request.form["id"]) + "', '" + pymysql.escape_string(hashed.decode()) + "', 'user')")
                    conn.commit()
                    return '<meta http-equiv="refresh" content="0;url=/login" />'
    else:
        ip = getip(request)
        ban = getban(ip)
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            return render_template('index.html', title = '회원가입', enter = '회원가입', logo = data['name'], tn = 15)

@app.route('/logout')
def logout():
    session['Now'] = False
    session.pop('DREAMER', None)
    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']) + '" />'

@app.route('/ban/<name>', methods=['POST', 'GET'])
def ban(name = None):
    if(request.method == 'POST'):
        if(session.get('Now') == True):
            ip = getip(request)
            curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
            rows = curs.fetchall()
            if(rows):
                if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
                    curs.execute("select * from ban where block = '" + pymysql.escape_string(name) + "'")
                    row = curs.fetchall()
                    if(row):
                        curs.execute("delete from ban where block = '" + pymysql.escape_string(name) + "'")
                    else:
                        curs.execute("insert into ban (block, end, why, band) value ('" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(request.form["end"]) + "', '" + pymysql.escape_string(request.form["why"]) + "', '')")
                    conn.commit()
                    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']) + '" />'
                else:
                    return render_template('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.')
            else:
                return render_template('index.html', title = '권한 오류', logo = data['name'], data = '계정이 없습니다.')
        else:
            return render_template('index.html', title = '권한 오류', logo = data['name'], data = '비 로그인 상태 입니다.')
    else:
        if(session.get('Now') == True):
            ip = getip(request)
            curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
            rows = curs.fetchall()
            if(rows):
                if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
                    curs.execute("select * from ban where block = '" + pymysql.escape_string(name) + "'")
                    row = curs.fetchall()
                    if(row):
                        now = '차단 해제'
                    else:
                        now = '차단'
                    return render_template('index.html', title = name, page = parse.quote(name), logo = data['name'], tn = 16, now = now)
                else:
                    return render_template('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.')
            else:
                return render_template('index.html', title = '권한 오류', logo = data['name'], data = '계정이 없습니다.')
        else:
            return render_template('index.html', title = '권한 오류', logo = data['name'], data = '비 로그인 상태 입니다.')

@app.route('/acl/<name>', methods=['POST', 'GET'])
def acl(name = None):
    if(request.method == 'POST'):
        if(session.get('Now') == True):
            ip = getip(request)
            curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
            rows = curs.fetchall()
            if(rows):
                if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
                    curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
                    row = curs.fetchall()
                    if(row):
                        if(request.form["select"] == 'admin'):
                           curs.execute("update data set acl = 'admin' where title = '" + pymysql.escape_string(name) + "'")
                        elif(request.form["select"] == 'user'):
                            curs.execute("update data set acl = 'user' where title = '" + pymysql.escape_string(name) + "'")
                        else:
                            curs.execute("update data set acl = '' where title = '" + pymysql.escape_string(name) + "'")
                        conn.commit()
                    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name) + '" />' 
                else:
                    return render_template('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.')
            else:
                return render_template('index.html', title = '권한 오류', logo = data['name'], data = '계정이 없습니다.')
        else:
            return render_template('index.html', title = '권한 오류', logo = data['name'], data = '비 로그인 상태 입니다.')
    else:
        if(session.get('Now') == True):
            ip = getip(request)
            curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
            rows = curs.fetchall()
            if(rows):
                if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
                    curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
                    row = curs.fetchall()
                    if(row):
                        if(row[0]['acl'] == 'admin'):
                            now = '관리자만'
                        elif(row[0]['acl'] == 'user'):
                            now = '유저 이상'
                        else:
                            now = '일반'
                        return render_template('index.html', title = name, page = parse.quote(name), logo = data['name'], tn = 19, now = '현재 ACL 상태는 ' + now)
                    else:
                        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name) + '" />' 
                else:
                    return render_template('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.')
            else:
                return render_template('index.html', title = '권한 오류', logo = data['name'], data = '계정이 없습니다.')
        else:
            return render_template('index.html', title = '권한 오류', logo = data['name'], data = '비 로그인 상태 입니다.')

@app.route('/admin/<name>', methods=['POST', 'GET'])
def admin(name = None):
    if(request.method == 'POST'):
        if(session.get('Now') == True):
            ip = getip(request)
            curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
            rows = curs.fetchall()
            if(rows):
                if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
                    curs.execute("select * from user where id = '" + pymysql.escape_string(name) + "'")
                    row = curs.fetchall()
                    if(row):
                        if(row[0]['acl'] == 'admin' or row[0]['acl'] == 'owner'):
                            curs.execute("update user set acl = 'user' where id = '" + pymysql.escape_string(name) + "'")
                        else:
                            curs.execute("update user set acl = '" + pymysql.escape_string(request.form["select"]) + "' where id = '" + pymysql.escape_string(name) + "'")
                        conn.commit()
                        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']) + '" />'
                    else:
                        return render_template('index.html', title = '사용자 오류', logo = data['name'], data = '계정이 없습니다.')
                else:
                    return render_template('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.')
            else:
                return render_template('index.html', title = '권한 오류', logo = data['name'], data = '계정이 없습니다.')
        else:
            return render_template('index.html', title = '권한 오류', logo = data['name'], data = '비 로그인 상태 입니다.')
    else:
        if(session.get('Now') == True):
            ip = getip(request)
            curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
            rows = curs.fetchall()
            if(rows):
                if(rows[0]['acl'] == 'owner'):
                    curs.execute("select * from user where id = '" + pymysql.escape_string(name) + "'")
                    row = curs.fetchall()
                    if(row):
                        if(row[0]['acl'] == 'admin' or row[0]['acl'] == 'owner'):
                            now = '권한 해제'
                        else:
                            now = '권한 부여'
                        return render_template('index.html', title = name, page = parse.quote(name), logo = data['name'], tn = 18, now = now)
                    else:
                        return render_template('index.html', title = '사용자 오류', logo = data['name'], data = '계정이 없습니다.')
                else:
                    return render_template('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.')
            else:
                return render_template('index.html', title = '권한 오류', logo = data['name'], data = '계정이 없습니다.')
        else:
            return render_template('index.html', title = '권한 오류', logo = data['name'], data = '비 로그인 상태 입니다.')

@app.route('/grammar')
def grammar():
    return render_template('index.html', title = '문법 설명', logo = data['name'], tn = 17)

@app.route('/ban')
def aban():
   return render_template('index.html', title = '권한 오류', logo = data['name'], data = '현재 차단 상태거나 ACL이 맞지 않습니다.')

@app.route('/version')
def version():
    return render_template('index.html', title = '버전', logo = data['name'], tn = 14)

@app.route('/user')
def user():
    ip = getip(request)
    return render_template('index.html', title = '유저 메뉴', logo = data['name'], data = ip + '<br><br><li><a href="/login">로그인</a></li><li><a href="/logout">로그아웃</a></li><li><a href="/register">회원가입</a></li>')

@app.route('/random')
def random():
    curs.execute("select * from data order by rand() limit 1")
    rows = curs.fetchall()
    if(rows):
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(rows[0]['title']) + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']) + '" />'

@app.errorhandler(404)
def uncaughtError(error):
    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']) + '" />'

if __name__ == '__main__':
    app.run(host = '0.0.0.0', port = 3000)
