from flask import Flask, request, session, render_template
app = Flask(__name__)

from urllib import parse
import json
import pymysql
import time
import base64

json_data=open('set.json').read()
data = json.loads(json_data)

conn = pymysql.connect(host = data['host'], user = data['user'], password = data['pw'], db = data['db'], charset = 'utf8')
curs = conn.cursor(pymysql.cursors.DictCursor)

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

def recent(title, ip, today, send):
    curs.execute("insert into rc (title, date, ip, send, leng, back) value ('" + title + "', '" + today + "', '" + ip + "', '" + send + "', '', '')")
    conn.commit()

def history(number, title, data, date, ip, send, leng):
    curs.execute("insert into history (id, title, data, date, ip, send, leng) value ()")
    conn.commit()

@app.route('/')
@app.route('/w/')
def redirect():
    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']) + '" />'

@app.route('/RecentChanges')
@app.route('/recentchanges')
def recentchanges():
    i = 0
    curs.execute("select * from rc order by date desc limit 50")
    rows = curs.fetchall()
    if(rows):
        return render_template('index.html', logo = data['name'], rows = rows, tn = 3)
    else:
         return render_template('index.html', logo = data['name'], rows = '', tn = 3)

@app.route('/search', methods=['POST', 'GET'])
def search():
    if(request.method == 'POST'):
        print(request.form["search"])
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(request.form["search"]) + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']) + '" />'

@app.route('/w/<name>')
def w(name = None):
    curs.execute("select * from data where title = '" + parse.quote(name) + "'")
    rows = curs.fetchall()
    if(rows):
        for row in rows:
            return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = parse.unquote(row['data']), license = data['license'], tn = 1)
    else:
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = '문서 없음', license = data['license'], tn = 1)

@app.route('/edit/<name>', methods=['POST', 'GET'])
def edit(name = None):
    if(request.method == 'POST'):
        curs.execute("select * from data where title = '" + parse.quote(name) + "'")
        rows = curs.fetchall()
        if(rows):
            ip = getip(request)
            today = getnow()
            recent(name, ip, today, request.form["send"])
            curs.execute("update data set data = '" + parse.quote(request.form["content"]) + "' where title = '" + parse.quote(name) + "'")
            conn.commit()
        else:
            ip = getip(request)
            today = getnow()
            recent(name, ip, today, request.form["send"])
            curs.execute("insert into data (title, data, acl) value ('" + parse.quote(name) + "', '" + parse.quote(request.form["content"]) + "', '')")
            conn.commit()
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name) + '" />'
    else:
        curs.execute("select * from data where title = '" + parse.quote(name) + "'")
        rows = curs.fetchall()
        if(rows):
            for row in rows:
                return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = parse.unquote(row['data']), tn = 2)
        else:
            return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = '', tn = 2)

@app.route('/setup')
def setup():
    curs.execute("create table if not exists data(title text not null, data longtext not null, acl text not null)")
    curs.execute("create table if not exists history(id text not null, title text not null, data longtext not null, date text not null, ip text not null, send text not null, leng text not null)")
    curs.execute("create table if not exists rc(title text not null, date text not null, ip text not null, send text not null, leng text not null, back text not null)")
    curs.execute("create table if not exists rd(title text not null, sub text not null, date text not null, ip text not null)")
    curs.execute("create table if not exists user(id text not null, pw text not null, acl text not null)")
    curs.execute("create table if not exists ban(block text not null, end text not null, why text not null, band text not null)")
    curs.execute("create table if not exists topic(id text not null, title text not null, sub text not null, data longtext not null, date text not null, ip text not null, block text not null)")
    return '문제 없음'

if __name__ == '__main__':
    app.run(host = '0.0.0.0', port = 3000)
