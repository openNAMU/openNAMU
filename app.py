from flask import Flask, request, session, render_template
app = Flask(__name__)

from urllib import parse
import json
import pymysql

json_data=open('set.json').read()
data = json.loads(json_data)

conn = pymysql.connect(host = data['host'], user = data['user'], password = data['pw'], db = data['db'], charset = 'utf8')
curs = conn.cursor(pymysql.cursors.DictCursor)

@app.route('/')
@app.route('/w/')
def redirect():
    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']) + '" />'

@app.route('/w/<name>')
def w(name = None):
    curs.execute("select * from data where title = '" + name + "'")
    rows = curs.fetchall()
    if(rows):
        for row in rows:
            return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = row['data'])
    else:
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name), data = '문서 없음')

@app.route('/edit/<name>', methods=['POST', 'GET'])
def edit(name = None):
    if(request.method == 'POST'):
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name) + '" />'
    else:
        curs.execute("select * from data where title = '" + name + "'")
        rows = curs.fetchall()
        if(rows):
            for row in rows:
                return render_template('edit.html', title = name, logo = data['name'], page = parse.quote(name), data = row['data'])
        else:
            return render_template('edit.html', title = name, logo = data['name'], page = parse.quote(name), data = '')

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
    app.run()
