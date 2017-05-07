from flask import Flask, request, session, render_template, send_file
app = Flask(__name__)

from urllib import parse
import json
import pymysql
import time
import re
import bcrypt
import os
import difflib
import hashlib

json_data = open('set.json').read()
set_data = json.loads(json_data)

print('port : ' + set_data['port'])

import logging
log = logging.getLogger('werkzeug')
log.setLevel(logging.ERROR)

def start():
    try:
        db_ex("select * from data limit 1")
    except:
        db_ex("create table data(title text, data longtext, acl text)")
    
    try:
        db_ex("select * from history limit 1")
    except:
        db_ex("create table history(id text, title text, data longtext, date text, ip text, send text, leng text)")
    
    try:
        db_ex("select * from rd limit 1")
    except:
        db_ex("create table rd(title text, sub text, date text)")
    
    try:
        db_ex("select * from user limit 1")
    except:
        db_ex("create table user(id text, pw text, acl text)")
    
    try:
        db_ex("select * from ban limit 1")
    except:
        db_ex("create table ban(block text, end text, why text, band text)")
    
    try:
        db_ex("select * from topic limit 1")
    except:
        db_ex("create table topic(id text, title text, sub text, data longtext, date text, ip text, block text)")
    
    try:
        db_ex("select * from stop limit 1")
    except:
        db_ex("create table stop(title text, sub text, close text)")
    
    try:
        db_ex("select * from rb limit 1")
    except:
        db_ex("create table rb(block text, end text, today text, blocker text, why text)")
    
    try:
        db_ex("select * from login limit 1")
    except:
        db_ex("create table login(user text, ip text, today text)")
    
    try:
        db_ex("select * from back limit 1")
    except:
        db_ex("create table back(title text, link text, type text)")
    
    try:
        db_ex("select * from cat limit 1")
    except:
        db_ex("create table cat(title text, cat text)")
        
    try:
        db_ex("select * from hidhi limit 1")
    except:
        db_ex("create table hidhi(title text, re text)")

    try:
        db_ex("select * from distop limit 1")
    except:
        db_ex("create table distop(id text, title text, sub text)") 

    try:
        db_ex("select * from agreedis limit 1")
    except:
        db_ex("create table agreedis(title text, sub text)") 
        
conn = pymysql.connect(host = set_data['host'], user = set_data['user'], password = set_data['pw'], charset = 'utf8mb4')
curs = conn.cursor(pymysql.cursors.DictCursor)

def db_com():
    conn.commit()

def url_pas(data):
    return parse.quote(data).replace('/','%2F')
    
def db_get():
    return curs.fetchall()
    
web_render = render_template
db_ex = curs.execute
db_pas = pymysql.escape_string

try:
    db_ex("use " + set_data['db'])
except:
    db_ex("create database " + set_data['db'])
    db_ex("use " + set_data['db'])
    db_ex("alter database " + set_data['db'] + " character set = utf8mb4 collate = utf8mb4_unicode_ci")
    
start()

from func import *
from mark import *

app.secret_key = hashlib.sha512(bytes(set_data['key'], 'ascii')).hexdigest()

@app.route('/upload', methods=['GET', 'POST'])
def upload():
    app.config['MAX_CONTENT_LENGTH'] = int(set_data['upload']) * 1024 * 1024
    
    ip = ip_check()
    ban = ban_check(ip)
    
    if(request.method == 'POST'):        
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            file = request.files['file']
            if(file):
                if(re.search('^([^./\\*<>|:?"]+)\.([Jj][Pp][Gg]|[Gg][Ii][Ff]|[Jj][Pp][Ee][Gg]|[Pp][Nn][Gg])$', file.filename)):
                    filename = file.filename
                    
                    if(os.path.exists(os.path.join('image', filename))):
                        return '<meta http-equiv="refresh" content="0;url=/error/16" />'
                    else:
                        file.save(os.path.join('image', filename))
                        
                        db_ex("insert into data (title, data, acl) value ('" + db_pas('파일:' + filename) + "', '" + db_pas('[[파일:' + filename + ']][br][br]{{{[[파일:' + filename + ']]}}}') + "', '')")
                        db_com()
                        
                        history_plus('파일:' + filename, '[[파일:' + filename + ']][br][br]{{{[[파일:' + filename + ']]}}}', get_time(), ip, '파일:' + filename + ' 업로드', '0')
                        
                        return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas('파일:' + filename) + '" />'
                else:
                    return '<meta http-equiv="refresh" content="0;url=/error/15" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/14" />'
    else:        
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            return web_render('index.html', login = login_check(), logo = set_data['name'], title = '업로드', tn = 21, number = set_data['upload'])
            
@app.route('/image/<path:name>')
def image_get(name = None):
    if(os.path.exists(os.path.join('image', name))):
        return send_file(os.path.join('image', name), mimetype='image')
    else:
        return web_render('index.html', login = login_check(), logo = set_data['name'], data = '이미지 없음.', title = '이미지 보기'), 404
        
@app.route('/adminlist')
def admin_list():
    i = 0
    div = '<div>'
    
    db_ex("select * from user where acl = 'admin' or acl = 'owner'")
    user_data = db_get()
    if(user_data):
        while(True):
            try:
                a = user_data[i]
            except:
                div = div + '</div>'
                break

            if(user_data[i]['acl'] == 'owner'):
                acl = '소유자'
            else:
                acl = '관리자'

            db_ex("select * from data where title = '사용자:" + user_data[i]['id'] + "'")
            user = db_get()
            if(user):
                name = '<a href="/w/' + url_pas('사용자:' + user_data[i]['id']) + '">' + user_data[i]['id'] + '</a> (' + acl + ')'
            else:
                name = '<a class="not_thing" href="/w/' + url_pas('사용자:' + user_data[i]['id']) + '">' + user_data[i]['id'] + '</a> (' + acl + ')'

            div = div + '<li>' + str(i + 1) + '. ' + name + '</li>'
            
            i += 1
            
        return web_render('index.html', login = login_check(), logo = set_data['name'], data = div, title = '관리자 목록')
    else:
        return web_render('index.html', login = login_check(), logo = set_data['name'], title = '관리자 목록')
        
@app.route('/recentchanges')
def recent_changes():
    i = 0
    div = '<div>'
    
    db_ex("select * from history order by date desc limit 50")
    rows = db_get()
    if(rows):
        while(True):
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                break
                
            if(rows[i]['send']):
                send = rows[i]['send']
                send = re.sub('<a href="\/w\/(?P<in>[^"]*)">(?P<out>[^&]*)<\/a>', '<a href="/w/\g<in>">\g<out></a>', send)
            else:
                send = '<br>'
                
            title = rows[i]['title']
            title = re.sub('<', '&lt;', title)
            title = re.sub('>', '&gt;', title)
            
            m = re.search("\+", rows[i]['leng'])
            n = re.search("\-", rows[i]['leng'])
            
            if(m):
                leng = '<span style="color:green;">' + rows[i]['leng'] + '</span>'
            elif(n):
                leng = '<span style="color:red;">' + rows[i]['leng'] + '</span>'
            else:
                leng = '<span style="color:gray;">' + rows[i]['leng'] + '</span>'
                
            if(admin_check() == 1):
                db_ex("select * from ban where block = '" + db_pas(rows[i]['ip']) + "'")
                row = db_get()
                if(row):
                    ban = ' <a href="/ban/' + url_pas(rows[i]['ip']) + '">(해제)</a>'
                else:
                    ban = ' <a href="/ban/' + url_pas(rows[i]['ip']) + '">(차단)</a>'
            else:
                ban = ''
                
            if(re.search('\.', rows[i]['ip'])):
                ip = rows[i]['ip'] + ' <a href="/record/' + url_pas(rows[i]['ip']) + '/n/1">(기록)</a>'
            else:
                db_ex("select * from data where title = '사용자:" + db_pas(rows[i]['ip']) + "'")
                row = db_get()
                if(row):
                    ip = '<a href="/w/' + url_pas('사용자:' + rows[i]['ip']) + '">' + rows[i]['ip'] + '</a> <a href="/record/' + url_pas(rows[i]['ip']) + '/n/1">(기록)</a>'
                else:
                    ip = '<a class="not_thing" href="/w/' + url_pas('사용자:' + rows[i]['ip']) + '">' + rows[i]['ip'] + '</a> <a href="/record/' + url_pas(rows[i]['ip']) + '/n/1">(기록)</a>'
                    
            if((int(rows[i]['id']) - 1) == 0):
                revert = ''
            else:
                revert = '<a href="/revert/' + url_pas(rows[i]['title']) + '/r/' + str(int(rows[i]['id']) - 1) + '">(되돌리기)</a>'
                
            div = div + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:33.33%;"><a href="/w/' + url_pas(rows[i]['title']) + '">' + title + '</a> <a href="/history/' + url_pas(rows[i]['title']) + '/n/1">(역사)</a> ' + revert + ' (' + leng + ')</td><td style="text-align: center;width:33.33%;">' + ip + ban + '</td><td style="text-align: center;width:33.33%;">' + rows[i]['date'] + '</td></tr><tr><td colspan="3" style="text-align: center;width:100%;">' + send + '</td></tr></tbody></table>'
            
            i += 1
            
        return web_render('index.html', login = login_check(), logo = set_data['name'], rows = div, tn = 3, title = '최근 변경내역')
    else:
        return web_render('index.html', login = login_check(), logo = set_data['name'], rows = '', tn = 3, title = '최근 변경내역')
        
@app.route('/history/<path:name>/r/<int:num>/hidden')
def history_hidden(name = None, num = None):
    if(owner_check() == 1):
        db_ex("select * from hidhi where title = '" + db_pas(name) + "' and re = '" + db_pas(str(num)) + "'")
        exist = db_get()
        if(exist):
            db_ex("delete from hidhi where title = '" + db_pas(name) + "' and re = '" + db_pas(str(num)) + "'")
        else:
            db_ex("insert into hidhi (title, re) value ('" + db_pas(name) + "', '" + db_pas(str(num)) + "')")
            
        db_com()
        
        return '<meta http-equiv="refresh" content="0;url=/history/' + url_pas(name) + '/n/1" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/history/' + url_pas(name) + '/n/1" />'
        
@app.route('/record/<path:name>/n/<int:num>')
def user_record(name = None, num = None):
    v = num * 50
    i = v - 50
    div = '<div>'
    
    db_ex("select * from history where ip = '" + db_pas(name) + "' order by date desc")
    rows = db_get()
    if(rows):
        while(True):
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                if(num != 1):
                    div = div + '<br><a href="/record/' + url_pas(name) + '/n/' + str(num - 1) + '">(이전)'
                break
                
            if(rows[i]['send']):
                send = rows[i]['send']
                send = re.sub('<a href="\/w\/(?P<in>[^"]*)">(?P<out>[^&]*)<\/a>', '<a href="/w/\g<in>">\g<out></a>', send)
            else:
                send = '<br>'
                
            title = rows[i]['title']
            title = re.sub('<', '&lt;', title)
            title = re.sub('>', '&gt;', title)
            
            m = re.search("\+", rows[i]['leng'])
            n = re.search("\-", rows[i]['leng'])
            
            if(m):
                leng = '<span style="color:green;">' + rows[i]['leng'] + '</span>'
            elif(n):
                leng = '<span style="color:red;">' + rows[i]['leng'] + '</span>'
            else:
                leng = '<span style="color:gray;">' + rows[i]['leng'] + '</span>'
                
            if(admin_check() == 1):
                db_ex("select * from ban where block = '" + db_pas(rows[i]['ip']) + "'")
                row = db_get()
                if(row):
                    ban = ' <a href="/ban/' + url_pas(rows[i]['ip']) + '">(해제)</a>'
                else:
                    ban = ' <a href="/ban/' + url_pas(rows[i]['ip']) + '">(차단)</a>'
            else:
                ban = ''
                
            if(re.search('\.', rows[i]['ip'])):
                ip = rows[i]['ip']
            else:
                db_ex("select * from data where title = '사용자:" + db_pas(rows[i]['ip']) + "'")
                row = db_get()
                if(row):
                    ip = '<a href="/w/' + url_pas('사용자:' + rows[i]['ip']) + '">' + rows[i]['ip'] + '</a>'
                else:
                    ip = '<a class="not_thing" href="/w/' + url_pas('사용자:' + rows[i]['ip']) + '">' + rows[i]['ip'] + '</a>'
                    
            if((int(rows[i]['id']) - 1) == 0):
                revert = ''
            else:
                revert = '<a href="/revert/' + url_pas(rows[i]['title']) + '/r/' + str(int(rows[i]['id']) - 1) + '">(되돌리기)</a>'
                
            div = div + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:33.33%;"><a href="/w/' + url_pas(rows[i]['title']) + '">' + title + '</a> (' + rows[i]['id'] + '판) <a href="/history/' + url_pas(rows[i]['title']) + '/n/1">(역사)</a> ' + revert + ' (' + leng + ')</td><td style="text-align: center;width:33.33%;">' + ip + ban +  '</td><td style="text-align: center;width:33.33%;">' + rows[i]['date'] + '</td></tr><tr><td colspan="3" style="text-align: center;width:100%;">' + send + '</td></tr></tbody></table>'
            
            if(i == v):
                div = div + '</div>'
                if(num == 1):
                    div = div + '<br><a href="/record/' + url_pas(name) + '/n/' + str(num + 1) + '">(다음)'
                else:
                    div = div + '<br><a href="/record/' + url_pas(name) + '/n/' + str(num - 1) + '">(이전) <a href="/record/' + url_pas(name) + '/n/' + str(num + 1) + '">(다음)'
                break

            i += 1
                
        return web_render('index.html', login = login_check(), logo = set_data['name'], rows = div, tn = 3, title = '유저 기록')
    else:
        return web_render('index.html', login = login_check(), logo = set_data['name'], rows = '', tn = 3, title = '유저 기록')
        
@app.route('/userlog/n/<int:number>')
def user_log(number = None):
    i = number * 50
    j = i - 50
    list = ''
    
    db_ex("select * from user")
    user_list = db_get()
    if(user_list):        
        while(True):
            try:
                a = user_list[j]
            except:
                if(number != 1):
                    list = list + '<br><a href="/userlog/n/' + str(number - 1) + '">(이전)'
                break
                
            if(admin_check() == 1):
                db_ex("select * from ban where block = '" + db_pas(user_list[j]['id']) + "'")
                ban_exist = db_get()
                if(ban_exist):
                    ban_button = ' <a href="/ban/' + url_pas(user_list[j]['id']) + '">(해제)</a>'
                else:
                    ban_button = ' <a href="/ban/' + url_pas(user_list[j]['id']) + '">(차단)</a>'
            else:
                ban_button = ''
                
            db_ex("select * from data where title = '사용자:" + db_pas(user_list[j]['id']) + "'")
            file = db_get()
            if(file):
                ip = '<a href="/w/' + url_pas('사용자:' + user_list[j]['id']) + '">' + user_list[j]['id'] + '</a> <a href="/record/' + url_pas(user_list[j]['id']) + '/n/1">(기록)</a>'
            else:
                ip = '<a class="not_thing" href="/w/' + url_pas('사용자:' + user_list[j]['id']) + '">' + user_list[j]['id'] + '</a> <a href="/record/' + url_pas(user_list[j]['id']) + '/n/1">(기록)</a>'
                
            list = list + '<li>' + str(j + 1) + '. ' + ip + ban_button + '</li>'
            
            if(j == i):
                if(number == 1):
                    list = list + '<br><a href="/userlog/n/' + str(number + 1) + '">(다음)'
                else:
                    list = list + '<br><a href="/userlog/n/' + str(number - 1) + '">(이전) <a href="/userlog/n/' + str(number + 1) + '">(다음)'
                break
            else:
                j += 1
                
        return web_render('index.html', login = login_check(), logo = set_data['name'], data = list, title = '유저 가입 기록')
    else:
        return web_render('index.html', login = login_check(), logo = set_data['name'], data = '', title = '유저 가입 기록')
        
@app.route('/backreset')
def backlink_reset():
    if(owner_check() == 1):
        i = 0
        
        db_ex("delete from back")
        db_com()
        
        db_ex("select * from data")
        all = db_get()
        if(all):
            while(True):
                try:
                    a = all[i]
                except:
                    break
                
                namumark(all[i]['title'], all[i]['data'])
                
                i += 1
        
        return web_render('index.html', login = login_check(), logo = set_data['name'], data = '에러 없음', title = '완료')
    else:
        return '<meta http-equiv="refresh" content="0;url=/error/3" />'
        
@app.route('/backlink/<path:name>/n/<int:num>')
def backlink(name = None, num = None):
    v = num * 50
    i = v - 50
    div = ''
    restart = 0
    
    db_ex("delete from back where title = '" + db_pas(name) + "' and link = ''")
    db_com()
    
    db_ex("select * from back where title = '" + db_pas(name) + "' order by link asc")
    rows = db_get()
    if(rows):        
        while(True):
            try:
                a = rows[i]
            except:
                if(num != 1):
                    div = div + '<br><a href="/backlink/n/' + str(num - 1) + '">(이전)'
                
                break
                
            if(rows[i]['type'] == 'include'):
                db_ex("select * from back where title = '" + db_pas(name) + "' and link = '" + db_pas(rows[i]['link']) + "' and type = ''")
                test = db_get()
                if(test):
                    restart = 1
                    
                    db_ex("delete from back where title = '" + db_pas(name) + "' and link = '" + db_pas(rows[i]['link']) + "' and type = ''")
                    db_com()
                
            if(not re.search('^사용자:', rows[i]['link'])):
                db_ex("select * from data where title = '" + db_pas(rows[i]['link']) + "'")
                row = db_get()
                if(row):
                    data = row[0]['data']
                    data = re.sub("(?P<in>\[include\((?P<out>(?:(?!\)\]|,).)*)((?:,\s?(?:[^)]*))+)?\)\])", "\g<in>\n\n[[\g<out>]]\n\n", data)
                    data = re.sub('^#(?:[Rr][Ee][Dd][Ii][Rr][Ee][Cc][Tt]|넘겨주기)\s(?P<in>[^\n]*)', '[[\g<in>]]', data)
                    data = namumark('', data)
                    
                    if(re.search("<a(?:(?:(?!href=).)*)?href=\"\/w\/" + url_pas(name) + "(?:\#[^\"]*)?\"(?:(?:(?!>).)*)?>([^<]*)<\/a>", data)):
                        div = div + '<li><a href="/w/' + url_pas(rows[i]['link']) + '">' + rows[i]['link'] + '</a>'
                        
                        if(rows[i]['type']):
                            div = div + ' (' + rows[i]['type'] + ')</li>'
                        else:
                            div = div + '</li>'
                            
                        if(i == v):
                            if(num == 1):
                                div = div + '<br><a href="/backlink/' + url_pas(name) + '/n/' + str(num + 1) + '">(다음)'
                            else:
                                div = div + '<br><a href="/backlink/' + url_pas(name) + '/n/' + str(num - 1) + '">(이전) <a href="/backlink/' + url_pas(name) + '/n/' + str(num + 1) + '">(다음)'
                                
                            break
                        else:
                            i += 1
                    else:
                        db_ex("delete from back where title = '" + db_pas(name) + "' and link = '" + db_pas(rows[i]['link']) + "'")
                        db_com()
                        
                        i += 1
                        v += 1
                else:
                    db_ex("delete from back where title = '" + db_pas(name) + "' and link = '" + db_pas(rows[i]['link']) + "'")
                    db_com()
                    
                    i += 1
                    v += 1
            else:
                db_ex("delete from back where title = '" + db_pas(name) + "' and link = '" + db_pas(rows[i]['link']) + "'")
                db_com()
                
                i += 1
                v += 1
                
        if(restart == 1):
            return '<meta http-equiv="refresh" content="0;url=/backlink/' + url_pas(name) + '/n/' + str(num) + '" />'
        else:    
            return web_render('index.html', login = login_check(), logo = set_data['name'], data = div, title = name, page = url_pas(name), sub = '역링크')
    else:
        return web_render('index.html', login = login_check(), logo = set_data['name'], data = '', title = name, page = url_pas(name), sub = '역링크')
        
@app.route('/recentdiscuss')
def recent_discuss():
    i = 0
    div = '<div>'
    
    db_ex("select * from rd order by date desc limit 50")
    rows = db_get()
    if(rows):
        while(True):
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
            
            div = div + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:50%;"><a href="/topic/' + url_pas(rows[i]['title']) + '/sub/' + url_pas(rows[i]['sub']) + '">' + title + '</a> (' + sub + ')</td><td style="text-align: center;width:50%;">' + rows[i]['date'] + '</td></tr></tbody></table>'
            
            i += 1
            
        return web_render('index.html', login = login_check(), logo = set_data['name'], rows = div, tn = 12, title = '최근 토론내역')
    else:
        return web_render('index.html', login = login_check(), logo = set_data['name'], rows = '', tn = 12, title = '최근 토론내역')
        
@app.route('/history/<path:name>/n/<int:num>', methods=['POST', 'GET'])
def history_view(name = None, num = None):
    if(request.method == 'POST'):
        return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(name) + '/r/' + request.form["b"] + '/diff/' + request.form["a"] + '" />'
    else:
        select = ''
        v = num * 50
        i = v - 50
        div = '<div>'
        
        db_ex("select * from history where title = '" + db_pas(name) + "' order by id+0 desc")
        rows = db_get()
        if(rows):
            while(True):
                style = ''
            
                try:
                    a = rows[i]
                except:
                    div = div + '</div>'
                    
                    if(num != 1):
                        div = div + '<br><a href="/history/' + url_pas(name) + '/n/' + str(num - 1) + '">(이전)'
                    break
                    
                select = '<option value="' + str(i + 1) + '">' + str(i + 1) + '</option>' + select
                
                if(rows[i]['send']):
                    send = rows[i]['send']
                    send = re.sub('<a href="\/w\/(?P<in>[^"]*)">(?P<out>[^&]*)<\/a>', '<a href="/w/\g<in>">\g<out></a>', send)
                else:
                    send = '<br>'
                    
                m = re.search("\+", rows[i]['leng'])
                n = re.search("\-", rows[i]['leng'])
                if(m):
                    leng = '<span style="color:green;">' + rows[i]['leng'] + '</span>'
                elif(n):
                    leng = '<span style="color:red;">' + rows[i]['leng'] + '</span>'
                else:
                    leng = '<span style="color:gray;">' + rows[i]['leng'] + '</span>'                    
                    
                if(re.search("\.", rows[i]["ip"])):
                    ip = rows[i]["ip"] + ' <a href="/record/' + url_pas(rows[i]["ip"]) + '/n/1">(기록)</a>'
                else:
                    db_ex("select * from data where title = '사용자:" + db_pas(rows[i]['ip']) + "'")
                    row = db_get()
                    if(row):
                        ip = '<a href="/w/' + url_pas('사용자:' + rows[i]['ip']) + '">' + rows[i]['ip'] + '</a> <a href="/record/' + url_pas(rows[i]["ip"]) + '/n/1">(기록)</a>'
                    else:
                        ip = '<a class="not_thing" href="/w/' + url_pas('사용자:' + rows[i]['ip']) + '">' + rows[i]['ip'] + '</a> <a href="/record/' + url_pas(rows[i]["ip"]) + '/n/1">(기록)</a>'
                        
                if(admin_check() == 1):
                    db_ex("select * from user where id = '" + db_pas(rows[i]['ip']) + "'")
                    row = db_get()
                    if(row):
                        if(row[0]['acl'] == 'owner' or row[0]['acl'] == 'admin'):
                            ban = ''
                        else:
                            db_ex("select * from ban where block = '" + db_pas(rows[i]['ip']) + "'")
                            row = db_get()
                            if(row):
                                ban = ' <a href="/ban/' + url_pas(rows[i]['ip']) + '">(해제)</a>'
                            else:
                                ban = ' <a href="/ban/' + url_pas(rows[i]['ip']) + '">(차단)</a>'
                    else:
                        db_ex("select * from ban where block = '" + db_pas(rows[i]['ip']) + "'")
                        row = db_get()
                        if(row):
                            ban = ' <a href="/ban/' + url_pas(rows[i]['ip']) + '">(해제)</a>'
                        else:
                            ban = ' <a href="/ban/' + url_pas(rows[i]['ip']) + '">(차단)</a>'
                            
                    if(owner_check() == 1):
                        db_ex("select * from hidhi where title = '" + db_pas(name) + "' and re = '" + db_pas(rows[i]['id']) + "'")
                        row = db_get()
                        if(row):                            
                            ip = ip + ' (숨김)'                            
                            hidden = ' <a href="/history/' + url_pas(name) + '/r/' + rows[i]['id'] + '/hidden">(공개)'
                        else:
                            hidden = ' <a href="/history/' + url_pas(name) + '/r/' + rows[i]['id'] + '/hidden">(숨김)'
                    else:
                        db_ex("select * from hidhi where title = '" + db_pas(name) + "' and re = '" + db_pas(rows[i]['id']) + "'")
                        row = db_get()
                        if(row):
                            ip = '숨김'
                            hidden = ''
                            send = '숨김'
                            ban = ''
                            style = 'display:none;'
                            v += 1
                        else:
                            hidden = ''
                else:
                    ban = ''
                    
                    db_ex("select * from hidhi where title = '" + db_pas(name) + "' and re = '" + db_pas(rows[i]['id']) + "'")
                    row = db_get()
                    if(row):
                        ip = '숨김'
                        hidden = ''
                        send = '숨김'
                        ban = ''
                        style = 'display:none;'
                        v += 1
                    else:
                        hidden = ''                
                        
                div = div + '<table style="width: 100%;' + style + '"><tbody><tr><td style="text-align: center;width:33.33%;">' + rows[i]['id'] + '판</a> <a href="/w/' + url_pas(rows[i]['title']) + '/r/' + rows[i]['id'] + '">(w)</a> <a href="/w/' + url_pas(rows[i]['title']) + '/raw/' + rows[i]['id'] + '">(Raw)</a> <a href="/revert/' + url_pas(rows[i]['title']) + '/r/' + rows[i]['id'] + '">(되돌리기)</a> (' + leng + ')</td><td style="text-align: center;width:33.33%;">' + ip + ban + hidden + '</td><td style="text-align: center;width:33.33%;">' + rows[i]['date'] + '</td></tr><tr><td colspan="3" style="text-align: center;width:100%;">' + send + '</td></tr></tbody></table>'
                
                if(i == v):
                    div = div + '</div>'
                    
                    if(num == 1):
                        div = div + '<br><a href="/history/' + url_pas(name) + '/n/' + str(num + 1) + '">(다음)'
                    else:
                        div = div + '<br><a href="/history/' + url_pas(name) + '/n/' + str(num - 1) + '">(이전) <a href="/history/' + url_pas(name) + '/n/' + str(num + 1) + '">(다음)'
                        
                    break
                else:
                    i += 1
                    
            return web_render('index.html', login = login_check(), logo = set_data['name'], rows = div, tn = 5, title = name, page = url_pas(name), select = select, sub = '역사')
        else:
            return web_render('index.html', login = login_check(), logo = set_data['name'], rows = '', tn = 5, title = name, page = url_pas(name), select = select, sub = '역사')
            
@app.route('/search', methods=['POST'])
def search():
    db_ex("select * from data where title = '" + db_pas(request.form["search"]) + "'")
    rows = db_get()
    if(rows):
        return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(request.form["search"]) + '" />'
    else:
        db_ex("select * from data where title like '%" + db_pas(request.form["search"]) + "%'")
        rows = db_get()
        if(rows):
            i = 0
            
            div = '<li>문서가 없습니다. <a href="/w/' + url_pas(request.form["search"]) + '">바로가기</a></li><br>'
            
            while(True):
                try:
                    div = div + '<li><a href="/w/' + url_pas(rows[i]['title']) + '">' + rows[i]['title'] + '</a></li>'
                except:
                    break
                    
                i += 1
        else:
            return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(request.form["search"]) + '" />'
            
        return web_render('index.html', login = login_check(), logo = set_data['name'], data = div, title = '검색')
        
@app.route('/w/<path:name>')
@app.route('/w/<path:name>/from/<path:redirect>')
def read_view(name = None, redirect = None):
    i = 0
    
    db_ex("select * from rd where title = '" + db_pas(name) + "' order by date asc")
    rows = db_get()
    while(True):
        try:
            a = rows[i]
        except:
            topic = ""
            break
            
        db_ex("select * from stop where title = '" + db_pas(rows[i]['title']) + "' and sub = '" + db_pas(rows[i]['sub']) + "' and close = 'O'")
        row = db_get()
        if(not row):
            topic = "open"
            
            break
        else:
            i += 1
            
    acl = ''
    
    m = re.search("^(.*)\/(.*)$", name)
    if(m):
        g = m.groups()
        uppage = g[0]
        style = ""
    else:
        uppage = ""
        style = "display:none;"
        
    if(admin_check() == 1):
        admin_memu = 'ACL'
    else:
        admin_memu = ''
    
    if(re.search("^분류:", name)):
        db_ex("select * from cat where title = '" + db_pas(name) + "' order by cat asc")
        rows = db_get()
        if(rows):
            div = ''
            i = 0
            
            while(True):
                try:
                    a = rows[i]
                except:
                    break
                    
                db_ex("select * from data where title = '" + db_pas(rows[i]['cat']) + "'")
                row = db_get()
                if(row):
                    aa = row[0]['data']                  
                    aa = namumark('', aa)
                    bb = re.search('<div style="width:100%;border: 1px solid #777;padding: 5px;margin-top: 1em;">분류:((?:(?!<\/div>).)*)<\/div>', aa)
                    if(bb):
                        cc = bb.groups()
                        
                        mm = re.search("^분류:(.*)", name)
                        if(mm):
                            ee = mm.groups()
                            
                            if(re.search("<a (class=\"not_thing\")? href=\"\/w\/" + url_pas(name) + "\">" + ee[0] + "<\/a>", cc[0])):
                                div = div + '<li><a href="/w/' + url_pas(rows[i]['cat']) + '">' + rows[i]['cat'] + '</a></li>'
                                
                                i += 1
                            else:
                                db_ex("delete from cat where title = '" + db_pas(name) + "' and cat = '" + db_pas(rows[i]['cat']) + "'")
                                db_com()
                                
                                i += 1
                        else:
                            db_ex("delete from cat where title = '" + db_pas(name) + "' and cat = '" + db_pas(rows[i]['cat']) + "'")
                            db_com()
                            
                            i += 1
                    else:
                        db_ex("delete from cat where title = '" + db_pas(name) + "' and cat = '" + db_pas(rows[i]['cat']) + "'")
                        db_com()
                        
                        i += 1
                else:
                    db_ex("delete from cat where title = '" + db_pas(name) + "' and cat = '" + db_pas(rows[i]['cat']) + "'")
                    db_com()
                    
                    i += 1
                    
            div = '<h2>분류</h2>' + div
        else:
            div = ''
    else:
        div = ''
    
    db_ex("select * from data where title = '" + db_pas(name) + "'")
    rows = db_get()
    if(rows):
        if(rows[0]['acl'] == 'admin'):
            acl = '(관리자)'
        elif(rows[0]['acl'] == 'user'):
            acl = '(유저)'
        else:
            if(not acl):
                acl = ''

        m = re.search("^사용자:(.*)", name)
        if(m):
            g = m.groups()
            
            db_ex("select * from user where id = '" + db_pas(g[0]) + "'")
            test = db_get()
            if(test):
                if(test[0]['acl'] == 'owner'):
                    acl = '(소유자)'
                elif(test[0]['acl'] == 'admin'):
                    acl = '(관리자)'

            db_ex("select * from ban where block = '" + db_pas(g[0]) + "'")
            user = db_get()
            if(user):
                elsedata = '{{{#!wiki style="border:2px solid red;padding:10px;"\r\n{{{+2 {{{#red 이 사용자는 차단 당했습니다.}}}}}}\r\n\r\n차단 해제 일 : ' + user[0]['end'] + '[br]사유 : ' + user[0]['why'] + '}}}[br]' + rows[0]['data']
            else:
                elsedata = rows[0]['data']
        else:
            elsedata = rows[0]['data']
                
        if(redirect):
            elsedata = re.sub("^#(?:[Rr][Ee][Dd][Ii][Rr][Ee][Cc][Tt]|넘겨주기)\s(?P<in>[^\n]*)", " * [[\g<in>]] 문서로 넘겨주기", elsedata)
                
        enddata = namumark(name, elsedata)
        
        m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
        if(m):
            result = m.groups()
            left = result[0]
        else:
            left = ''
            
        return web_render('index.html', login = login_check(), title = name, logo = set_data['name'], page = url_pas(name), data = enddata + div, license = set_data['license'], tn = 1, acl = acl, left = left, uppage = uppage, style = style, topic = topic, redirect = redirect, admin = admin_memu)
    else:
        m = re.search("^사용자:(.*)", name)
        if(m):
            g = m.groups()
            
            db_ex("select * from ban where block = '" + db_pas(g[0]) + "'")
            user = db_get()
            if(user):
                elsedata = '{{{#!wiki style="border:2px solid red;padding:10px;"\r\n{{{+2 {{{#red 이 사용자는 차단 당했습니다.}}}}}}\r\n\r\n차단 해제 일 : ' + user[0]['end'] + '[br]사유 : ' + user[0]['why'] + '}}}[br]' + '문서 없음'
            else:
                elsedata = '문서 없음'
        else:
            elsedata = '문서 없음'
            
        if(redirect):
            elsedata = re.sub("^#(?:[Rr][Ee][Dd][Ii][Rr][Ee][Cc][Tt]|넘겨주기)\s(?P<in>[^\n]*)", " * [[\g<in>]] 문서로 넘겨주기", elsedata)
        
        return web_render('index.html', login = login_check(), title = name, logo = set_data['name'], page = url_pas(name), data = namumark(name, elsedata) + div, license = set_data['license'], tn = 1, uppage = uppage, style = style, acl = acl, topic = topic, redirect = redirect, admin = admin_memu), 404
        
@app.route('/w/<path:name>/r/<int:num>')
def old_view(name = None, num = None):
    db_ex("select * from hidhi where title = '" + db_pas(name) + "' and re = '" + db_pas(str(num)) + "'")
    row = db_get()
    if(row):
        if(owner_check() == 1):
            db_ex("select * from history where title = '" + db_pas(name) + "' and id = '" + str(num) + "'")
            rows = db_get()
            if(rows):
                enddata = namumark(name, rows[0]['data'])
                
                m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
                if(m):
                    result = m.groups()
                    left = result[0]
                else:
                    left = ''
                    
                return web_render('index.html', login = login_check(), title = name, logo = set_data['name'], page = url_pas(name), data = enddata, tn = 6, left = left, sub = '옛 문서')
            else:
                return '<meta http-equiv="refresh" content="0;url=/history/' + url_pas(name) + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'
    else:
        db_ex("select * from history where title = '" + db_pas(name) + "' and id = '" + str(num) + "'")
        rows = db_get()
        if(rows):
            enddata = namumark(name, rows[0]['data'])
            
            m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
            if(m):
                result = m.groups()
                left = result[0]
            else:
                left = ''
                
            return web_render('index.html', login = login_check(), title = name, logo = set_data['name'], page = url_pas(name), data = enddata, tn = 6, left = left, sub = '옛 문서')
        else:
            return '<meta http-equiv="refresh" content="0;url=/history/' + url_pas(name) + '" />'
            
@app.route('/w/<path:name>/raw/<int:num>')
def old_raw(name = None, num = None):
    db_ex("select * from hidhi where title = '" + db_pas(name) + "' and re = '" + db_pas(str(num)) + "'")
    row = db_get()
    if(row):
        if(owner_check() == 1):
            db_ex("select * from history where title = '" + db_pas(name) + "' and id = '" + str(num) + "'")
            rows = db_get()
            if(rows):
                enddata = re.sub('<', '&lt;', rows[0]['data'])
                enddata = re.sub('>', '&gt;', enddata)
                enddata = re.sub('"', '&quot;', enddata)
                
                enddata = '<pre>' + enddata + '</pre>'
                
                return web_render('index.html', login = login_check(), title = name, logo = set_data['name'], page = url_pas(name), data = enddata, sub = '옛 Raw')
            else:
                return '<meta http-equiv="refresh" content="0;url=/history/' + url_pas(name) + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'
    else:
        db_ex("select * from history where title = '" + db_pas(name) + "' and id = '" + str(num) + "'")
        rows = db_get()
        if(rows):
            enddata = re.sub('<', '&lt;', rows[0]['data'])
            enddata = re.sub('>', '&gt;', enddata)
            enddata = re.sub('"', '&quot;', enddata)
            
            enddata = '<pre>' + enddata + '</pre>'
            
            return web_render('index.html', login = login_check(), title = name, logo = set_data['name'], page = url_pas(name), data = enddata, sub = '옛 Raw')
        else:
            return '<meta http-equiv="refresh" content="0;url=/history/' + url_pas(name) + '" />'
            
@app.route('/raw/<path:name>')
def raw_view(name = None):
    db_ex("select * from data where title = '" + db_pas(name) + "'")
    rows = db_get()
    if(rows):
        enddata = re.sub('<', '&lt;', rows[0]['data'])
        enddata = re.sub('>', '&gt;', enddata)
        enddata = re.sub('"', '&quot;', enddata)
        
        enddata = '<pre>' + enddata + '</pre>'
        
        return web_render('index.html', login = login_check(), title = name, logo = set_data['name'], page = url_pas(name), data = enddata, tn = 7, sub = 'Raw')
    else:
        return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(name) + '" />'
        
@app.route('/revert/<path:name>/r/<int:num>', methods=['POST', 'GET'])
def revert(name = None, num = None):
    ip = ip_check()
    can = acl_check(ip, name)
    today = get_time()
    
    if(request.method == 'POST'):
        db_ex("select * from hidhi where title = '" + db_pas(name) + "' and re = '" + db_pas(str(num)) + "'")
        row = db_get()
        if(row):
            if(owner_check() == 1):        
                db_ex("select * from history where title = '" + db_pas(name) + "' and id = '" + str(num) + "'")
                rows = db_get()
                if(rows):
                    if(can == 1):
                        return '<meta http-equiv="refresh" content="0;url=/ban" />'
                    else:
                        db_ex("select * from data where title = '" + db_pas(name) + "'")
                        row = db_get()
                        if(row):
                            leng = leng_check(len(row[0]['data']), len(rows[0]['data']))
                            
                            db_ex("update data set data = '" + db_pas(rows[0]['data']) + "' where title = '" + db_pas(name) + "'")
                            db_com()
                        else:
                            leng = '+' + str(len(rows[0]['data']))
                            
                            db_ex("insert into data (title, data, acl) value ('" + db_pas(name) + "', '" + db_pas(rows[0]['data']) + "', '')")
                            db_com()
                            
                        history_plus(name, rows[0]['data'], today, ip, '문서를 ' + str(num) + '판으로 되돌렸습니다.', leng)
                        
                        return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(name) + '" />'
                else:
                    return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(name) + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/3" />'
        else:
            db_ex("select * from history where title = '" + db_pas(name) + "' and id = '" + str(num) + "'")
            rows = db_get()
            if(rows):                
                if(can == 1):
                    return '<meta http-equiv="refresh" content="0;url=/ban" />'
                else:                    
                    db_ex("select * from data where title = '" + db_pas(name) + "'")
                    row = db_get()
                    if(row):
                        leng = leng_check(len(row[0]['data']), len(rows[0]['data']))
                        
                        db_ex("update data set data = '" + db_pas(rows[0]['data']) + "' where title = '" + db_pas(name) + "'")
                        db_com()
                    else:
                        leng = '+' + str(len(rows[0]['data']))
                        
                        db_ex("insert into data (title, data, acl) value ('" + db_pas(name) + "', '" + db_pas(rows[0]['data']) + "', '')")
                        db_com()
                        
                    history_plus(name, rows[0]['data'], today, ip, '문서를 ' + str(num) + '판으로 되돌렸습니다.', leng)
                    
                    return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(name) + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(name) + '" />'            
    else:
        db_ex("select * from hidhi where title = '" + db_pas(name) + "' and re = '" + db_pas(str(num)) + "'")
        row = db_get()
        if(row):
            if(owner_check() == 1):                
                if(can == 1):
                    return '<meta http-equiv="refresh" content="0;url=/ban" />'
                else:
                    db_ex("select * from history where title = '" + db_pas(name) + "' and id = '" + str(num) + "'")
                    rows = db_get()
                    if(rows):
                        return web_render('index.html', login = login_check(), title = name, logo = set_data['name'], page = url_pas(name), r = url_pas(str(num)), tn = 13, plus = '정말 되돌리시겠습니까?', sub = '되돌리기')
                    else:
                        return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(name) + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/3" />'
        else:            
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                db_ex("select * from history where title = '" + db_pas(name) + "' and id = '" + str(num) + "'")
                rows = db_get()
                if(rows):
                    return web_render('index.html', login = login_check(), title = name, logo = set_data['name'], page = url_pas(name), r = url_pas(str(num)), tn = 13, plus = '정말 되돌리시겠습니까?', sub = '되돌리기')
                else:
                    return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(name) + '" />'
                    
@app.route('/edit/<path:name>', methods=['POST', 'GET'])
def edit(name = None):
    ip = ip_check()
    can = acl_check(ip, name)
    
    if(request.method == 'POST'):
        m = re.search('(?:[^A-Za-zㄱ-힣0-9 ])', request.form["send"])
        
        if(m):
            return '<meta http-equiv="refresh" content="0;url=/error/17" />'
        else:
            today = get_time()
            
            content = savemark(request.form["content"])
            
            db_ex("select * from data where title = '" + db_pas(name) + "'")
            rows = db_get()
            if(rows):
                if(rows[0]['data'] == content):
                    return '<meta http-equiv="refresh" content="0;url=/error/18" />'
                else:                    
                    if(can == 1):
                        return '<meta http-equiv="refresh" content="0;url=/ban" />'
                    else:                        
                        leng = leng_check(len(rows[0]['data']), len(content))
                        history_plus(name, content, today, ip, request.form["send"], leng)
                        
                        db_ex("update data set data = '" + db_pas(content) + "' where title = '" + db_pas(name) + "'")
                        db_com()
            else:                
                if(can == 1):
                    return '<meta http-equiv="refresh" content="0;url=/ban" />'
                else:
                    leng = '+' + str(len(content))
                    history_plus(name, content, today, ip, request.form["send"], leng)
                    
                    db_ex("insert into data (title, data, acl) value ('" + db_pas(name) + "', '" + db_pas(content) + "', '')")
                    db_com()
                    
            include_check(name, content)
            
            return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(name) + '" />'
    else:        
        if(can == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            db_ex("select * from data where title = '" + db_pas(set_data["help"]) + "'")
            rows = db_get()
            if(rows):
                newdata = re.sub('^#(?:[Rr][Ee][Dd][Ii][Rr][Ee][Cc][Tt]|넘겨주기)\s(?P<in>[^\n]*)', ' * [[\g<in>]] 문서로 넘겨주기', rows[0]["data"])
                left = namumark(name, newdata)
            else:
                left = ''
                
            db_ex("select * from data where title = '" + db_pas(name) + "'")
            rows = db_get()
            if(rows):
                return web_render('index.html', login = login_check(), title = name, logo = set_data['name'], page = url_pas(name), data = rows[0]['data'], tn = 2, left = left, sub = '편집')
            else:
                return web_render('index.html', login = login_check(), title = name, logo = set_data['name'], page = url_pas(name), data = '', tn = 2, left = left, sub = '편집')
                
@app.route('/edit/<path:name>/section/<int:num>', methods=['POST', 'GET'])
def section_edit(name = None, num = None):
    ip = ip_check()
    can = acl_check(ip, name)
    
    if(request.method == 'POST'):
        m = re.search('(?:[^A-Za-zㄱ-힣0-9 ])', request.form["send"])
        if(m):
            return '<meta http-equiv="refresh" content="0;url=/error/17" />'
        else:
            today = get_time()
            
            content = savemark(request.form["content"])
            
            db_ex("select * from data where title = '" + db_pas(name) + "'")
            rows = db_get()
            if(rows):
                if(request.form["otent"] == content):
                    return '<meta http-equiv="refresh" content="0;url=/error/18" />'
                else:                    
                    if(can == 1):
                        return '<meta http-equiv="refresh" content="0;url=/ban" />'
                    else:                        
                        leng = leng_check(len(request.form['otent']), len(content))
                        content = rows[0]['data'].replace(request.form['otent'], content)
                        
                        history_plus(name, content, today, ip, request.form["send"], leng)
                        
                        db_ex("update data set data = '" + db_pas(content) + "' where title = '" + db_pas(name) + "'")
                        db_com()
                        
                    include_check(name, content)
                    
                    return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(name) + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(name) + '" />'
    else:        
        if(can == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            db_ex("select * from data where title = '" + db_pas(set_data["help"]) + "'")
            rows = db_get()
            if(rows):
                newdata = re.sub('^#(?:[Rr][Ee][Dd][Ii][Rr][Ee][Cc][Tt]|넘겨주기)\s(?P<in>[^\n]*)', ' * [[\g<in>]] 문서로 넘겨주기', rows[0]["data"])
                
                left = namumark(name, newdata)
            else:
                left = ''
                
            db_ex("select * from data where title = '" + db_pas(name) + "'")
            rows = db_get()
            if(rows):
                i = 0
                j = 0
                
                gdata = rows[0]['data'] + '\r\n'
                
                while(True):
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
                    return web_render('index.html', login = login_check(), title = name, logo = set_data['name'], page = url_pas(name), data = gdata, tn = 2, left = left, section = 1, number = num, sub = '편집')
                else:
                    return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(name) + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(name) + '" />'
                
@app.route('/preview/<path:name>', methods=['POST'])
def preview(name = None):
    ip = ip_check()
    can = acl_check(ip, name)
    
    if(can == 1):
        return '<meta http-equiv="refresh" content="0;url=/ban" />'
    else:            
        newdata = request.form["content"]
        newdata = re.sub('^#(?:[Rr][Ee][Dd][Ii][Rr][Ee][Cc][Tt]|넘겨주기)\s(?P<in>[^\n]*)', ' * [[\g<in>]] 문서로 넘겨주기', newdata)
        enddata = namumark(name, newdata)
        
        db_ex("select * from data where title = '" + db_pas(data["help"]) + "'")
        rows = db_get()
        if(rows):
            newdata = re.sub('^#(?:[Rr][Ee][Dd][Ii][Rr][Ee][Cc][Tt]|넘겨주기)\s(?P<in>[^\n]*)', ' * [[\g<in>]] 문서로 넘겨주기', rows[0]["data"])
            
            left = namumark(name, newdata)
        else:
            left = ''
            
        return web_render('index.html', login = login_check(), title = name, logo = set_data['name'], page = url_pas(name), data = request.form["content"], tn = 2, preview = 1, enddata = enddata, left = left, sub = '미리보기')
        
@app.route('/preview/<path:name>/section/<int:num>', methods=['POST'])
def section_preview(name = None, num = None):
    ip = ip_check()
    can = acl_check(ip, name)
    
    if(can == 1):
        return '<meta http-equiv="refresh" content="0;url=/ban" />'
    else:            
        newdata = request.form["content"]
        newdata = re.sub('^#(?:[Rr][Ee][Dd][Ii][Rr][Ee][Cc][Tt]|넘겨주기)\s(?P<in>[^\n]*)', ' * [[\g<in>]] 문서로 넘겨주기', newdata)
        enddata = namumark(name, newdata)
        
        db_ex("select * from data where title = '" + db_pas(data["help"]) + "'")
        rows = db_get()
        if(rows):
            newdata = re.sub('^#(?:[Rr][Ee][Dd][Ii][Rr][Ee][Cc][Tt]|넘겨주기)\s(?P<in>[^\n]*)', ' * [[\g<in>]] 문서로 넘겨주기', rows[0]["data"])
            left = namumark(name, newdata)
        else:
            left = ''
            
        return web_render('index.html', login = login_check(), title = name, logo = set_data['name'], page = url_pas(name), data = request.form["content"], tn = 2, preview = 1, enddata = enddata, left = left, section = 1, number = num, odata = request.form["otent"], sub = '미리보기')
        
@app.route('/delete/<path:name>', methods=['POST', 'GET'])
def delete(name = None):
    ip = ip_check()
    can = acl_check(ip, name)
    
    if(request.method == 'POST'):
        db_ex("select * from data where title = '" + db_pas(name) + "'")
        rows = db_get()
        if(rows):
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                today = get_time()
                
                leng = '-' + str(len(rows[0]['data']))
                
                history_plus(name, '', today, ip, '문서를 삭제 했습니다.', leng)
                
                db_ex("delete from data where title = '" + db_pas(name) + "'")
                db_com()
                
                return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(name) + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(name) + '" />'
    else:
        db_ex("select * from data where title = '" + db_pas(name) + "'")
        rows = db_get()
        if(rows):
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                return web_render('index.html', login = login_check(), title = name, logo = set_data['name'], page = url_pas(name), tn = 8, plus = '정말 삭제 하시겠습니까?', sub = '삭제')
        else:
            return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(name) + '" />'
            
@app.route('/move/<path:name>', methods=['POST', 'GET'])
def move(name = None):
    ip = ip_check()
    can = acl_check(ip, name)
    today = get_time()
    
    if(request.method == 'POST'):
        db_ex("select * from data where title = '" + db_pas(name) + "'")
        rows = db_get()

        if(can == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            leng = '0'
            db_ex("select * from history where title = '" + db_pas(request.form["title"]) + "'")
            row = db_get()
            if(row):
                return '<meta http-equiv="refresh" content="0;url=/error/19" />'
            else:
                history_plus(name, rows[0]['data'], today, ip, '<a href="/w/' + url_pas(name) + '">' + name + '</a> 문서를 <a href="/w/' + url_pas(request.form["title"]) + '">' + request.form["title"] + '</a> 문서로 이동 했습니다.', leng)
                
                if(rows):
                    db_ex("update data set title = '" + db_pas(request.form["title"]) + "' where title = '" + db_pas(name) + "'")

                db_ex("update history set title = '" + db_pas(request.form["title"]) + "' where title = '" + db_pas(name) + "'")
                db_com()
                return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(request.form["title"]) + '" />'
    else:
        if(can == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            return web_render('index.html', login = login_check(), title = name, logo = set_data['name'], page = url_pas(name), tn = 9, plus = '정말 이동 하시겠습니까?', sub = '이동')
            
@app.route('/other')
def other():
    return web_render('index.html', login = login_check(), title = '기타 메뉴', logo = set_data['name'], data = '<h2 style="margin-top: 0px;">기록</h2><li><a href="/blocklog/n/1">유저 차단 기록</a></li><li><a href="/userlog/n/1">유저 가입 기록</a></li><li><a href="/manager/6">유저 기록</a></li><h2>기타</h2><li><a href="/titleindex">모든 문서</a></li><li><a href="/upload">업로드</a></li><li><a href="/adminlist">관리자 목록</a></li><li><a href="/manager/1">관리자 메뉴</a></li><br>이 오픈나무의 버전은 <a href="https://github.com/2DU/openNAMU/blob/master/version.md">1.9.1</a> 입니다.')
    
@app.route('/manager/<int:num>', methods=['POST', 'GET'])
def manager(num = None):
    if(num == 1):
        return web_render('index.html', login = login_check(), title = '관리자 메뉴', logo = set_data['name'], data = '<h2 style="margin-top: 0px;">관리자 및 소유자</h2><li><a href="/manager/2">문서 ACL</a></li><li><a href="/manager/3">유저 체크</a></li><li><a href="/manager/4">유저 차단</a></li><h2>소유자</h2><li><a href="/backreset">모든 역링크 재 생성</a></li><li><a href="/manager/5">관리자 권한 주기</a></li><h2>기타</h2><li>이 메뉴에 없는 기능은 해당 문서의 역사나 토론에서 바로 사용 가능함</li>')
    elif(num == 2):
        if(request.method == 'POST'):
            return '<meta http-equiv="refresh" content="0;url=/acl/' + url_pas(request.form["name"]) + '" />'
        else:
            return web_render('index.html', login = login_check(), title = 'ACL 이동', logo = set_data['name'], data = '<form id="usrform" method="POST" action="/manager/2"><input name="name" type="text"><br><br><button class="btn btn-primary" type="submit">이동</button></form>')
    elif(num == 3):
        if(request.method == 'POST'):
            return '<meta http-equiv="refresh" content="0;url=/check/' + url_pas(request.form["name"]) + '" />'
        else:
            return web_render('index.html', login = login_check(), title = '체크 이동', logo = set_data['name'], data = '<form id="usrform" method="POST" action="/manager/3"><input name="name" type="text"><br><br><button class="btn btn-primary" type="submit">이동</button></form>')
    elif(num == 4):
        if(request.method == 'POST'):
            return '<meta http-equiv="refresh" content="0;url=/ban/' + url_pas(request.form["name"]) + '" />'
        else:
            return web_render('index.html', login = login_check(), title = '차단 이동', logo = set_data['name'], data = '<form id="usrform" method="POST" action="/manager/4"><input name="name" type="text"><br><br><button class="btn btn-primary" type="submit">이동</button><br><br><span>아이피 앞 두자리 (XXX.XXX) 입력하면 대역 차단</span></form>')
    elif(num == 5):
        if(request.method == 'POST'):
            return '<meta http-equiv="refresh" content="0;url=/admin/' + url_pas(request.form["name"]) + '" />'
        else:
            return web_render('index.html', login = login_check(), title = '권한 이동', logo = set_data['name'], data = '<form id="usrform" method="POST" action="/manager/5"><input name="name" type="text"><br><br><button class="btn btn-primary" type="submit">이동</button></form>')   
    elif(num == 6):
        if(request.method == 'POST'):
            return '<meta http-equiv="refresh" content="0;url=/record/' + url_pas(request.form["name"]) + '/n/1" />'
        else:
            return web_render('index.html', login = login_check(), title = '기록 이동', logo = set_data['name'], data = '<form id="usrform" method="POST" action="/manager/6"><input name="name" type="text"><br><br><button class="btn btn-primary" type="submit">이동</button></form>')    
    else:
        return '<meta http-equiv="refresh" content="0;url=/" />'
        
@app.route('/titleindex')
def title_index():
    i = 0
    data = '<div>'
    db_ex("select title from data order by title asc")
    title_list = db_get()
    if(title_list):
        while(True):
            try:
                a = title_list[i]
            except:
                break

            data = data + '<li>' + str(i + 1) + '. <a href="/w/' + url_pas(title_list[i]['title']) + '">' + title_list[i]['title'] + '</a></li>'
            
            i += 1

        data = data + '</div>'

        return web_render('index.html', login = login_check(), logo = set_data['name'], rows = data + '<br><span>이 위키에는 총 ' + str(i) + '개의 문서가 있습니다.</span>', tn = 4, title = '모든 문서')
    else:
        return web_render('index.html', login = login_check(), logo = set_data['name'], rows = '', tn = 4, title = '모든 문서')
        
@app.route('/topic/<path:name>', methods=['POST', 'GET'])
def topic_list(name = None):
    if(request.method == 'POST'):
        return '<meta http-equiv="refresh" content="0;url=/topic/' + url_pas(name) + '/sub/' + url_pas(request.form["topic"]) + '" />'
    else:
        div = '<div>'
        i = 0
        j = 1
        db_ex("select * from rd where title = '" + db_pas(name) + "' order by date asc")
        rows = db_get()
        while(True):
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                break
                
            db_ex("select * from topic where title = '" + db_pas(rows[i]['title']) + "' and sub = '" + db_pas(rows[i]['sub']) + "' and id = '1' order by sub asc")
            aa = db_get()
            
            indata = namumark(name, aa[0]['data'])
            
            if(aa[0]['block'] == 'O'):
                indata = '블라인드 되었습니다.'
                block = 'style="background: gainsboro;"'
            else:
                block = ''

            ip = ip_pas(aa[0]['ip'])
                
            db_ex("select * from stop where title = '" + db_pas(rows[i]['title']) + "' and sub = '" + db_pas(rows[i]['sub']) + "' and close = 'O'")
            row = db_get()
            if(not row):
                div = div + '<h2><a href="/topic/' + url_pas(rows[i]['title']) + '/sub/' + url_pas(rows[i]['sub']) + '">' + str(j) + '. ' + rows[i]['sub'] + '</a></h2><table id="toron"><tbody><tr><td id="toroncolorgreen"><a href="javascript:void(0);" id="1">#1</a> ' + ip + ' <span style="float:right;">' + aa[0]['date'] + '</span></td></tr><tr><td ' + block + '>' + indata + '</td></tr></tbody></table><br>'
                j += 1
                
            i += 1
            
        return web_render('index.html', login = login_check(), title = name, page = url_pas(name), logo = set_data['name'], plus = div, tn = 10, list = 1, sub = '토론 목록')
        
@app.route('/topic/<path:name>/close')
def close_topic_list(name = None):
    div = '<div>'
    i = 0
    
    db_ex("select * from stop where title = '" + db_pas(name) + "' and close = 'O' order by sub asc")
    rows = db_get()
    while(True):
        try:
            a = rows[i]
        except:
            div = div + '</div>'
            break
            
        db_ex("select * from topic where title = '" + db_pas(name) + "' and sub = '" + db_pas(rows[i]['sub']) + "' and id = '1'")
        row = db_get()
        if(row):
            indata = namumark(name, row[0]['data'])
            
            if(row[0]['block'] == 'O'):
                indata = '블라인드 되었습니다.'
                block = 'style="background: gainsboro;"'
            else:
                block = ''

            ip = ip_pas(row[0]['ip'])
                
            div = div + '<h2><a href="/topic/' + url_pas(name) + '/sub/' + url_pas(rows[i]['sub']) + '">' + str((i + 1)) + '. ' + rows[i]['sub'] + '</a></h2><table id="toron"><tbody><tr><td id="toroncolorgreen"><a href="javascript:void(0);" id="1">#1</a> ' + ip + ' <span style="float:right;">' + row[0]['date'] + '</span></td></tr><tr><td ' + block + '>' + indata + '</td></tr></tbody></table><br>'
            
        i += 1
        
    return web_render('index.html', login = login_check(), title = name, page = url_pas(name), logo = set_data['name'], plus = div, tn = 10, sub = '닫힌 토론')
    
@app.route('/topic/<path:name>/agree')
def agree_topic_list(name = None):
    div = '<div>'
    i = 0
    
    db_ex("select * from agreedis where title = '" + db_pas(name) + "' order by sub asc")
    agree_list = db_get()
    while(True):
        try:
            a = agree_list[i]
        except:
            div = div + '</div>'
            break
            
        db_ex("select * from topic where title = '" + db_pas(name) + "' and sub = '" + db_pas(agree_list[i]['sub']) + "' and id = '1'")
        data = db_get()
        if(data):
            indata = namumark(name, data[0]['data'])
            
            if(data[0]['block'] == 'O'):
                indata = '블라인드 되었습니다.'
                block = 'style="background: gainsboro;"'
            else:
                block = ''

            ip = ip_pas(data[0]['ip'])
                
            div = div + '<h2><a href="/topic/' + url_pas(name) + '/sub/' + url_pas(data[i]['sub']) + '">' + str((i + 1)) + '. ' + data[i]['sub'] + '</a></h2><table id="toron"><tbody><tr><td id="toroncolorgreen"><a href="javascript:void(0);" id="1">#1</a> ' + 아이디 + ' <span style="float:right;">' + data[0]['date'] + '</span></td></tr><tr><td ' + block + '>' + indata + '</td></tr></tbody></table><br>'
            
        i += 1
        
    return web_render('index.html', login = login_check(), title = name, page = url_pas(name), logo = set_data['name'], plus = div, tn = 10, sub = '합의된 토론')
    
@app.route('/topic/<path:name>/sub/<path:sub>', methods=['POST', 'GET'])
def topic(name = None, sub = None):
    ip = ip_check()
    ban = topic_check(ip, name, sub)
    admin = admin_check()
    
    if(request.method == 'POST'):
        db_ex("select * from topic where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "' order by id+0 desc limit 1")
        rows = db_get()
        if(rows):
            number = int(rows[0]['id']) + 1
        else:
            number = 1
        
        if(ban == 1 and not admin == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            db_ex("select * from user where id = '" + db_pas(ip) + "'")
            rows = db_get()
            if(rows):
                if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
                    ip = ip + ' - Admin'
                    
            today = get_time()
            rd_plus(name, sub, today)
            
            aa = request.form["content"]
            aa = re.sub("\[\[(분류:(?:(?:(?!\]\]).)*))\]\]", "[br]", aa)
            aa = savemark(aa)
            
            db_ex("insert into topic (id, title, sub, data, date, ip, block) value ('" + str(number) + "', '" + db_pas(name) + "', '" + db_pas(sub) + "', '" + db_pas(aa) + "', '" + today + "', '" + ip + "', '')")
            db_com()
            
            return '<meta http-equiv="refresh" content="0;url=/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '" />'
    else:
        style = ''

        db_ex("select * from stop where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "' and close = 'O'")
        close = db_get()

        db_ex("select * from stop where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "' and close = ''")
        stop = db_get()
        
        if(admin == 1):
            div = '<div>'
            
            if(close):
                div = div + '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/close">(토론 열기)</a> '
            else:
                div = div + '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/close">(토론 닫기)</a> '
            
            if(stop):
                div = div + '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/stop">(토론 재개)</a> '
            else:
                div = div + '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/stop">(토론 정지)</a> '

            db_ex("select * from agreedis where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "'")
            agree = db_get()
            if(agree):
                div = div + '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/agree">(합의 취소)</a>'
            else:
                div = div + '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/agree">(합의 완료)</a>'
            
            div = div + '<br><br>'
        else:
            div = '<div>'
        
        if(stop or close):
            if(not admin == 1):
                style = 'display:none;'
        
        db_ex("select * from topic where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "' order by id+0 asc")
        rows = db_get()

        db_ex("select * from distop where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "' order by id+0 asc")
        top = db_get()

        i = 0


        if(top):
            while(True):
                try:
                    a = top[i]
                except:
                    break

                num = int(top[i]['id']) - 1

                if(i == 0):
                    start = rows[num]['ip']
                    
                top_data = namumark('', rows[num]['data'])
                top_data = re.sub("(?P<in>#(?:[0-9]*))", '<a href="\g<in>">\g<in></a>', top_data)
                        
                ip = ip_pas(rows[num]['ip'])
                                   
                div = div + '<table id="toron"><tbody><tr><td id="toroncolorred"><a href="#' + top[i]['id'] + '" id="' + top[i]['id'] + '-nt">#' + top[i]['id'] + '</a> ' + ip + ' <span style="float:right;">' + rows[num]['date'] + '</span></td></tr><tr><td>' + top_data + '</td></tr></tbody></table><br>'
                    
                i = i + 1

        i = 0
        
        while(True):
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                break
                
            if(i == 0):
                start = rows[i]['ip']
                
            indata = namumark('', rows[i]['data'])
            indata = re.sub("(?P<in>#(?:[0-9]*))", '<a href="\g<in>">\g<in></a>', indata)
            
            if(rows[i]['block'] == 'O'):
                indata = '블라인드 되었습니다.'
                block = 'style="background: gainsboro;"'
            else:
                block = ''

            m = re.search("^([^-]*)\s\-\s(Close|Reopen|Stop|Restart|Agreement|Settlement)$", rows[i]['ip'])
            if(m):
                ban = ""
            else:
                if(admin == 1):
                    if(rows[i]['block'] == 'O'):
                        isblock = ' <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/b/' + str(i + 1) + '">(해제)</a>'
                    else:
                        isblock = ' <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/b/' + str(i + 1) + '">(블라인드)</a>'

                    db_ex("select * from distop where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "' and id = '" + db_pas(str(i + 1)) + "'")
                    row = db_get()
                    if(row):
                        isblock = isblock + ' <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/notice/' + str(i + 1) + '">(해제)</a>'
                    else:
                        isblock = isblock + ' <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/notice/' + str(i + 1) + '">(공지)</a>'
                        
                    n = re.search("\- (?:Admin)$", rows[i]['ip'])
                    if(n):
                        ban = isblock
                    else:
                        db_ex("select * from ban where block = '" + db_pas(rows[i]['ip']) + "'")
                        row = db_get()
                        if(row):
                            ban = ' <a href="/ban/' + url_pas(rows[i]['ip']) + '">(해제)</a>' + isblock
                        else:
                            ban = ' <a href="/ban/' + url_pas(rows[i]['ip']) + '">(차단)</a>' + isblock
                else:
                    ban = ""

            ip = ip_pas(rows[i]['ip'])
                    
            if(rows[i]['ip'] == start):
                j = i + 1
                
                div = div + '<table id="toron"><tbody><tr><td id="toroncolorgreen"><a href="javascript:void(0);" id="' + str(j) + '">#' + str(j) + '</a> ' + ip + ban + ' <span style="float:right;">' + rows[i]['date'] + '</span></td></tr><tr><td ' + block + '>' + indata + '</td></tr></tbody></table><br>'
            else:
                j = i + 1
                
                div = div + '<table id="toron"><tbody><tr><td id="toroncolor"><a href="javascript:void(0);" id="' + str(j) + '">#' + str(j) + '</a> ' + ip + ban + ' <span style="float:right;">' + rows[i]['date'] + '</span></td></tr><tr><td ' + block + '>' + indata + '</td></tr></tbody></table><br>'
                
            i += 1
            
        return web_render('index.html', login = login_check(), title = name, page = url_pas(name), suburl = url_pas(sub), toron = sub, logo = set_data['name'], rows = div, tn = 11, ban = ban, style = style, sub = '토론')
        
@app.route('/topic/<path:name>/sub/<path:sub>/b/<int:num>')
def topic_block(name = None, sub = None, num = None):
    if(admin_check() == 1):
        db_ex("select * from topic where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "' and id = '" + str(num) + "'")
        block = db_get()
        if(block):
            if(block[0]['block'] == 'O'):
                db_ex("update topic set block = '' where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "' and id = '" + str(num) + "'")
            else:
                db_ex("update topic set block = 'O' where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "' and id = '" + str(num) + "'")
            db_com()
            
            rd_plus(name, sub, get_time())
            
            return '<meta http-equiv="refresh" content="0;url=/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/error/3" />'
        
@app.route('/topic/<path:name>/sub/<path:sub>/notice/<int:num>')
def topic_top(name = None, sub = None, num = None):
    if(admin_check() == 1):
        db_ex("select * from topic where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "' and id = '" + str(num) + "'")
        topic_data = db_get()
        if(topic_data):
            db_ex("select * from distop where id = '" + str(num) + "' and title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "'")
            top_data = db_get()
            if(top_data):
                db_ex("delete from distop where id = '" + str(num) + "' and title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "'")
            else:
                db_ex("insert into distop (id, title, sub) value ('" + db_pas(str(num)) + "', '" + db_pas(name) + "', '" + db_pas(sub) + "')")
            db_com()
            
            rd_plus(name, sub, get_time())

            return '<meta http-equiv="refresh" content="0;url=/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/error/3" />'
        
@app.route('/topic/<path:name>/sub/<path:sub>/stop')
def topic_stop(name = None, sub = None):
    if(admin_check() == 1):
        ip = ip_check()
        
        db_ex("select * from topic where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "' limit 1")
        topic_check = db_get()
        if(topic_check):
            time = get_time()
            
            db_ex("select * from stop where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "' and close = ''")
            stop = db_get()
            if(stop):
                db_ex("insert into topic (id, title, sub, data, date, ip, block) value ('" + db_pas(str(int(topic_check[0]['id']) + 1)) + "', '" + db_pas(name) + "', '" + db_pas(sub) + "', 'Restart', '" + db_pas(time) + "', '" + db_pas(ip) + " - Restart', '')")
                db_ex("delete from stop where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "' and close = ''")
            else:
                db_ex("insert into topic (id, title, sub, data, date, ip, block) value ('" + db_pas(str(int(topic_check[0]['id']) + 1)) + "', '" + db_pas(name) + "', '" + db_pas(sub) + "', 'Stop', '" + db_pas(time) + "', '" + db_pas(ip) + " - Stop', '')")
                db_ex("insert into stop (title, sub, close) value ('" + db_pas(name) + "', '" + db_pas(sub) + "', '')")
            db_com()
            
            rd_plus(name, sub, time)
            
            return '<meta http-equiv="refresh" content="0;url=/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/error/3" />'
        
@app.route('/topic/<path:name>/sub/<path:sub>/close')
def topic_close(name = None, sub = None):
    if(admin_check() == 1):
        ip = ip_check()
        
        db_ex("select * from topic where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "' order by id+0 desc limit 1")
        topic_check = db_get()
        if(topic_check):
            time = get_time()
            
            db_ex("select * from stop where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "' and close = 'O'")
            close = db_get()
            if(close):
                db_ex("insert into topic (id, title, sub, data, date, ip, block) value ('" + db_pas(str(int(topic_check[0]['id']) + 1)) + "', '" + db_pas(name) + "', '" + db_pas(sub) + "', 'Reopen', '" + db_pas(time) + "', '" + db_pas(ip) + " - Reopen', '')")
                db_ex("delete from stop where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "' and close = 'O'")
            else:
                db_ex("insert into topic (id, title, sub, data, date, ip, block) value ('" + db_pas(str(int(topic_check[0]['id']) + 1)) + "', '" + db_pas(name) + "', '" + db_pas(sub) + "', 'Close', '" + db_pas(time) + "', '" + db_pas(ip) + " - Close', '')")
                db_ex("insert into stop (title, sub, close) value ('" + db_pas(name) + "', '" + db_pas(sub) + "', 'O')")
            db_com()
            
            rd_plus(name, sub, time)
            
            return '<meta http-equiv="refresh" content="0;url=/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/error/3" />'
        
@app.route('/topic/<path:name>/sub/<path:sub>/agree')
def topic_agree(name = None, sub = None):
    if(admin_check() == 1):
        ip = ip_check()
        
        db_ex("select id from topic where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "' order by id+0 desc limit 1")
        topic_check = db_get()
        if(topic_check):
            time = get_time()
            
            db_ex("select * from agreedis where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "'")
            agree = db_get()
            if(agree):
                db_ex("insert into topic (id, title, sub, data, date, ip, block) value ('" + db_pas(str(int(topic_check[0]['id']) + 1)) + "', '" + db_pas(name) + "', '" + db_pas(sub) + "', 'Settlement', '" + db_pas(time) + "', '" + db_pas(ip) + " - Settlement', '')")
                db_ex("delete from agreedis where title = '" + db_pas(name) + "' and sub = '" + db_pas(sub) + "'")
            else:
                db_ex("insert into topic (id, title, sub, data, date, ip, block) value ('" + db_pas(str(int(topic_check[0]['id']) + 1)) + "', '" + db_pas(name) + "', '" + db_pas(sub) + "', 'Agreement', '" + db_pas(time) + "', '" + db_pas(ip) + " - Agreement', '')")
                db_ex("insert into agreedis (title, sub) value ('" + db_pas(name) + "', '" + db_pas(sub) + "')")
            db_com()
            
            rd_plus(name, sub, time)
            
            return '<meta http-equiv="refresh" content="0;url=/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/error/3" />'
        
@app.route('/login', methods=['POST', 'GET'])
def login():
    ip = ip_check()
    ban = ban_check(ip)
        
    if(request.method == 'POST'):        
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            db_ex("select * from user where id = '" + db_pas(request.form["id"]) + "'")
            user = db_get()
            if(user):
                if(session.get('Now') == True):
                    return '<meta http-equiv="refresh" content="0;url=/error/11" />'
                elif(bcrypt.checkpw(bytes(request.form["pw"], 'utf-8'), bytes(user[0]['pw'], 'utf-8'))):
                    session['Now'] = True
                    session['DREAMER'] = request.form["id"]
                    
                    db_ex("insert into login (user, ip, today) value ('" + db_pas(request.form["id"]) + "', '" + db_pas(ip) + "', '" + db_pas(get_time()) + "')")
                    db_com()
                    
                    return '<meta http-equiv="refresh" content="0;url=/user" />'
                else:
                    return '<meta http-equiv="refresh" content="0;url=/error/13" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/12" />'
    else:        
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            if(session.get('Now') == True):
                return '<meta http-equiv="refresh" content="0;url=/error/11" />'
            else:
                return web_render('index.html', login = login_check(), title = '로그인', enter = '로그인', logo = set_data['name'], tn = 15)
                
@app.route('/change', methods=['POST', 'GET'])
def change_password():
    ip = ip_check()
    ban = ban_check(ip)
    
    if(request.method == 'POST'):      
        if(request.form["pw2"] == request.form["pw3"]):
            if(ban == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                db_ex("select * from user where id = '" + db_pas(request.form["id"]) + "'")
                user = db_get()
                if(user):
                    if(session.get('Now') == True):
                        session['Now'] = False
                        session.pop('DREAMER', None)
                        return '<meta http-equiv="refresh" content="0;url=/change" />'
                    elif(bcrypt.checkpw(bytes(request.form["pw"], 'utf-8'), bytes(user[0]['pw'], 'utf-8'))):
                        hashed = bcrypt.hashpw(bytes(request.form["pw2"], 'utf-8'), bcrypt.gensalt())
                        
                        db_ex("update user set pw = '" + db_pas(hashed.decode()) + "' where id = '" + db_pas(request.form["id"]) + "'")
                        db_com()
                        
                        return '<meta http-equiv="refresh" content="0;url=/login" />'
                    else:
                        return '<meta http-equiv="refresh" content="0;url=/error/10" />'
                else:
                    return '<meta http-equiv="refresh" content="0;url=/error/9" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/20" />'
    else:        
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            if(session.get('Now') == True):
                session['Now'] = False
                session.pop('DREAMER', None)
                return '<meta http-equiv="refresh" content="0;url=/change" />'
            else:
                return web_render('index.html', login = login_check(), title = '비밀번호 변경', enter = '변경', logo = set_data['name'], tn = 15)
                
@app.route('/check/<name>')
def user_check(name = None, sub = None):
    db_ex("select * from user where id = '" + db_pas(name) + "'")
    user = db_get()
    if(user and user[0]['acl'] == 'owner' or user and user[0]['acl'] == 'admin'):
        return '<meta http-equiv="refresh" content="0;url=/error/4" />'
    else:
        if(admin_check() == 1):
            m = re.search('(?:[0-9](?:[0-9][0-9])?\.[0-9](?:[0-9][0-9])?\.[0-9](?:[0-9][0-9])?\.[0-9](?:[0-9][0-9])?)', name)
            if(m):
                db_ex("select * from login where ip = '" + db_pas(name) + "' order by today desc")
                row = db_get()
                if(row):
                    i = 0
                    c = ''
                    while(True):
                        try:
                            c = c + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:33.33%;">' + row[i]['user'] + '</td><td style="text-align: center;width:33.33%;">' + row[i]['ip'] + '</td><td style="text-align: center;width:33.33%;">' + row[i]['today'] + '</td></tr></tbody></table>'
                        except:
                            break
                            
                        i += 1
                        
                    return web_render('index.html', login = login_check(), title = '다중 검사', logo = set_data['name'], tn = 22, rows = c)
                else:
                    return web_render('index.html', login = login_check(), title = '다중 검사', logo = set_data['name'], tn = 22, rows = '')
            else:
                db_ex("select * from login where user = '" + db_pas(name) + "' order by today desc")
                row = db_get()
                if(row):
                    i = 0
                    c = ''
                    while(True):
                        try:
                            c = c + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:33.33%;">' + row[i]['user'] + '</td><td style="text-align: center;width:33.33%;">' + row[i]['ip'] + '</td><td style="text-align: center;width:33.33%;">' + row[i]['today'] + '</td></tr></tbody></table>'
                        except:
                            break
                            
                        i += 1
                        
                    return web_render('index.html', login = login_check(), title = '다중 검사', logo = set_data['name'], tn = 22, rows = c)
                else:
                    return web_render('index.html', login = login_check(), title = '다중 검사', logo = set_data['name'], tn = 22, rows = '')
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'
                
@app.route('/register', methods=['POST', 'GET'])
def register():
    ip = ip_check()
    ban = ban_check(ip)
    
    if(request.method == 'POST'):        
        if(request.form["pw"] == request.form["pw2"]):
            if(ban == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                m = re.search('(?:[^A-Za-zㄱ-힣0-9 ])', request.form["id"])
                if(m):
                    return '<meta http-equiv="refresh" content="0;url=/error/8" />'
                else:
                    if(len(request.form["id"]) > 20):
                        return '<meta http-equiv="refresh" content="0;url=/error/7" />'
                    else:
                        db_ex("select * from user where id = '" + db_pas(request.form["id"]) + "'")
                        rows = db_get()
                        if(rows):
                            return '<meta http-equiv="refresh" content="0;url=/error/6" />'
                        else:
                            hashed = bcrypt.hashpw(bytes(request.form["pw"], 'utf-8'), bcrypt.gensalt())
                            
                            if(request.form["id"] == set_data['owner']):
                                db_ex("insert into user (id, pw, acl) value ('" + db_pas(request.form["id"]) + "', '" + db_pas(hashed.decode()) + "', 'owner')")
                            else:
                                db_ex("insert into user (id, pw, acl) value ('" + db_pas(request.form["id"]) + "', '" + db_pas(hashed.decode()) + "', 'user')")
                            db_com()
                            
                            return '<meta http-equiv="refresh" content="0;url=/login" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/20" />'
    else:        
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            return web_render('index.html', login = login_check(), title = '회원가입', enter = '회원가입', logo = set_data['name'], tn = 15)
            
@app.route('/logout')
def logout():
    session['Now'] = False
    session.pop('DREAMER', None)
    return '<meta http-equiv="refresh" content="0;url=/user" />'
    
@app.route('/ban/<name>', methods=['POST', 'GET'])
def user_ban(name = None):
    db_ex("select * from user where id = '" + db_pas(name) + "'")
    user = db_get()
    if(user and user[0]['acl'] == 'owner' or user and user[0]['acl'] == 'admin'):
        return '<meta http-equiv="refresh" content="0;url=/error/4" />'
    else:
        if(request.method == 'POST'):
            if(admin_check() == 1):
                ip = ip_check()
                
                if(not re.search("[0-9]{4}-[0-9]{2}-[0-9]{2}", request.form["end"])):
                    end = ''
                else:
                    end = request.form["end"]

                db_ex("select * from ban where block = '" + db_pas(name) + "'")
                row = db_get()
                if(row):
                    rb_plus(name, '해제', 시간(), ip, '')
                    
                    db_ex("delete from ban where block = '" + db_pas(name) + "'")
                else:
                    b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))$", name)
                    if(b):
                        rb_plus(name, end, 시간(), ip, request.form["why"])
                        
                        db_ex("insert into ban (block, end, why, band) value ('" + db_pas(name) + "', '" + db_pas(end) + "', '" + db_pas(request.form["why"]) + "', 'O')")
                    else:
                        rb_plus(name, end, 시간(), ip, request.form["why"])
                        
                        db_ex("insert into ban (block, end, why, band) value ('" + db_pas(name) + "', '" + db_pas(end) + "', '" + db_pas(request.form["why"]) + "', '')")
                db_com()
                
                return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(data['frontpage']) + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/3" />'
        else:
            if(admin_check() == 1):
                db_ex("select * from ban where block = '" + db_pas(name) + "'")
                row = db_get()
                if(row):
                    now = '차단 해제'
                else:
                    b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))$", name)
                    if(b):
                        now = '대역 차단'
                    else:
                        now = '차단'
                        
                return web_render('index.html', login = login_check(), title = name, page = url_pas(name), logo = set_data['name'], tn = 16, now = now, today = get_time(), sub = '차단')
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/3" />'
                
@app.route('/acl/<path:name>', methods=['POST', 'GET'])
def acl(name = None):
    if(request.method == 'POST'):
        if(admin_check() == 1):
            db_ex("select * from data where title = '" + db_pas(name) + "'")
            row = db_get()
            if(row):
                if(request.form["select"] == 'admin'):
                   db_ex("update data set acl = 'admin' where title = '" + db_pas(name) + "'")
                elif(request.form["select"] == 'user'):
                    db_ex("update data set acl = 'user' where title = '" + db_pas(name) + "'")
                else:
                    db_ex("update data set acl = '' where title = '" + db_pas(name) + "'")
                    
                db_com()
                
            return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(name) + '" />' 
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'
    else:
        if(admin_check() == 1):
            db_ex("select * from data where title = '" + db_pas(name) + "'")
            row = db_get()
            if(row):
                if(row[0]['acl'] == 'admin'):
                    now = '관리자만'
                elif(row[0]['acl'] == 'user'):
                    now = '유저 이상'
                else:
                    now = '일반'
                    
                return web_render('index.html', login = login_check(), title = name, page = url_pas(name), logo = set_data['name'], tn = 19, now = '현재 ACL 상태는 ' + now, sub = 'ACL')
            else:
                return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(name) + '" />' 
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'
            
@app.route('/admin/<name>', methods=['POST', 'GET'])
def user_admin(name = None):
    if(request.method == 'POST'):
        if(owner_check() == 1):
            db_ex("select * from user where id = '" + db_pas(name) + "'")
            user = db_get()
            if(user):
                if(user[0]['acl'] == 'admin' or user[0]['acl'] == 'owner'):
                    db_ex("update user set acl = 'user' where id = '" + db_pas(name) + "'")
                else:
                    db_ex("update user set acl = '" + db_pas(request.form["select"]) + "' where id = '" + db_pas(name) + "'")
                db_com()
                
                return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(data['frontpage']) + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/5" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'
    else:
        if(owner_check() == 1):
            db_ex("select * from user where id = '" + db_pas(name) + "'")
            user = db_get()
            if(user):
                if(user[0]['acl'] == 'admin' or user[0]['acl'] == 'owner'):
                    now = '권한 해제'
                else:
                    now = '권한 부여'
                    
                return web_render('index.html', login = login_check(), title = name, page = url_pas(name), logo = data['name'], tn = 18, now = now, sub = '권한 부여')
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/5" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'
            
@app.route('/ban')
def are_you_ban():
    ip = ip_check()
    
    if(ban_check(ip) == 1):
        db_ex("select * from ban where block = '" + db_pas(ip) + "'")
        rows = db_get()
        if(rows):
            if(rows[0]['end']):
                end = rows[0]['end'] + ' 까지 차단 상태 입니다. / 사유 : ' + rows[0]['why']                
                now = get_time()
                
                now = re.sub(':', '', now)
                now = re.sub('\-', '', now)
                now = re.sub(' ', '', now)
                now = int(now)
                
                day = rows[0]['end']
                day = re.sub('\-', '', day)    
                
                if(now >= int(day + '000000')):
                    db_ex("delete from ban where block = '" + db_pas(ip) + "'")
                    db_com()
                    
                    end = '차단이 풀렸습니다. 다시 시도 해 보세요.'
            else:
                end = '영구 차단 상태 입니다. / 사유 : ' + rows[0]['why']
        else:
            b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
            if(b):
                results = b.groups()
                
                db_ex("select * from ban where block = '" + db_pas(results[0]) + "' and band = 'O'")
                row = db_get()
                if(row):
                    if(row[0]['end']):
                        end = row[0]['end'] + ' 까지 차단 상태 입니다. / 사유 : ' + rows[0]['why']             
                        
                        now = get_time()
                        now = re.sub(':', '', now)
                        now = re.sub('\-', '', now)
                        now = re.sub(' ', '', now)
                        now = int(now)    
                        
                        day = row[0]['end']
                        day = re.sub('\-', '', day)
                        
                        if(now >= int(day + '000000')):
                            db_ex("delete from ban where block = '" + db_pas(results[0]) + "' and band = 'O'")
                            db_com()
                            
                            end = '차단이 풀렸습니다. 다시 시도 해 보세요.'
                    else:
                        end = '영구 차단 상태 입니다. / 사유 : ' + row[0]['why']                
    else:
        end = '권한이 맞지 않는 상태 입니다.'
        
    return web_render('index.html', login = login_check(), title = '권한 오류', logo = set_data['name'], data = end), 401
    
@app.route('/w/<path:name>/r/<int:a>/diff/<int:b>')
def diff_data(name = None, a = None, b = None):
    db_ex("select * from history where id = '" + db_pas(str(a)) + "' and title = '" + db_pas(name) + "'")
    a_raw_data = db_get()
    if(a_raw_data):
        db_ex("select * from history where id = '" + db_pas(str(b)) + "' and title = '" + db_pas(name) + "'")
        b_raw_data = db_get()
        if(b_raw_data):
            a_data = re.sub('<', '&lt;', a_raw_data[0]['data'])
            a_data = re.sub('>', '&gt;', a_data)
            a_data = re.sub('"', '&quot;', a_data)
            
            b_data = re.sub('<', '&lt;', b_raw_data[0]['data'])
            b_data = re.sub('>', '&gt;', b_data)
            b_data = re.sub('"', '&quot;', b_data)
            
            diff_data = difflib.SequenceMatcher(None, a_data, b_data)
            result = diff(diff_data)
            
            result = '<pre>' + result + '</pre>'
            
            return web_render('index.html', login = login_check(), title = name, logo = set_data['name'], data = result, sub = '비교')
        else:
            return '<meta http-equiv="refresh" content="0;url=/history/' + url_pas(name) + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/history/' + url_pas(name) + '" />'
        
@app.route('/user')
def user_info():
    ip = ip_check()
    raw_ip = ip
    
    db_ex("select * from user where id = '" + db_pas(ip) + "'")
    rows = db_get()
    if(ban_check(ip) == 0):
        if(rows):
            if(rows[0]['acl'] == 'admin' or rows[0]['acl'] == 'owner'):
                if(rows[0]['acl'] == 'admin'):
                    acl = '관리자'
                else:
                    acl = '소유자'
            else:
                acl = '유저'
        else:
            acl = '일반'
    else:
        acl = '차단'
        
    ip = ip_pas(ip)
        
    return web_render('index.html', login = login_check(), title = '유저 메뉴', logo = set_data['name'], data = ip + '<br><br><span>권한 상태 : ' + acl + '<h2>로그인 관련</h2><li><a href="/login">로그인</a></li><li><a href="/logout">로그아웃</a></li><li><a href="/register">회원가입</a></li><h2>기타</h2><li><a href="/change">비밀번호 변경</a></li><li><a href="/count">기여 횟수</a></li><li><a href="/record/' + raw_ip + '/n/1">기여 목록</a></li>')
    
@app.route('/count')
def count_edit():
    db_ex("select count(title) from history where ip = '" + ip_check() + "'")
    i = db_get()
    
    if(i):
        return web_render('index.html', login = login_check(), title = '기여 횟수', logo = set_data['name'], data = "기여 횟수 : " + str(i[0]["count(title)"]))
    else:
        return web_render('index.html', login = login_check(), title = '기여 횟수', logo = set_data['name'], data = "기여 횟수 : 0")
        
@app.route('/random')
def random():
    db_ex("select * from data order by rand() limit 1")
    rows = db_get()
    if(rows):
        return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(rows[0]['title']) + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/" />'
        
@app.route('/error/<int:num>')
def 오류(num = None):
    if(num == 1):
        return web_render('index.html', login = login_check(), title = '권한 오류', logo = set_data['name'], data = '비 로그인 상태 입니다.'), 401
    elif(num == 2):
        return web_render('index.html', login = login_check(), title = '권한 오류', logo = set_data['name'], data = '이 계정이 없습니다.'), 401
    elif(num == 3):
        return web_render('index.html', login = login_check(), title = '권한 오류', logo = set_data['name'], data = '권한이 모자랍니다.'), 401
    elif(num == 4):
        return web_render('index.html', login = login_check(), title = '권한 오류', logo = set_data['name'], data = '관리자는 차단, 검사 할 수 없습니다.'), 401
    elif(num == 5):
        return web_render('index.html', login = login_check(), title = '유저 오류', logo = set_data['name'], data = '그런 계정이 없습니다.'), 401
    elif(num == 6):
        return web_render('index.html', login = login_check(), title = '가입 오류', logo = set_data['name'], data = '동일한 아이디의 유저가 있습니다.'), 401
    elif(num == 7):
        return web_render('index.html', login = login_check(), title = '가입 오류', logo = set_data['name'], data = '아이디는 20글자보다 짧아야 합니다.'), 401
    elif(num == 8):
        return web_render('index.html', login = login_check(), title = '가입 오류', logo = set_data['name'], data = '아이디에는 한글과 알파벳과 공백만 허용 됩니다.'), 401
    elif(num == 9):
        return web_render('index.html', login = login_check(), title = '변경 오류', logo = set_data['name'], data = '그런 계정이 없습니다.'), 401
    elif(num == 10):
        return web_render('index.html', login = login_check(), title = '변경 오류', logo = set_data['name'], data = '비밀번호가 다릅니다.'), 401
    elif(num == 11):
        return web_render('index.html', login = login_check(), title = '로그인 오류', logo = set_data['name'], data = '이미 로그인 되어 있습니다.'), 401
    elif(num == 12):
        return web_render('index.html', login = login_check(), title = '로그인 오류', logo = set_data['name'], data = '그런 계정이 없습니다.'), 401
    elif(num == 13):
        return web_render('index.html', login = login_check(), title = '로그인 오류', logo = set_data['name'], data = '비밀번호가 다릅니다.'), 401
    elif(num == 14):
        return web_render('index.html', login = login_check(), title = '업로드 오류', logo = set_data['name'], data = 'jpg, gif, jpeg, png만 가능 합니다.'), 401
    elif(num == 15):
        return web_render('index.html', login = login_check(), title = '업로드 오류', logo = set_data['name'], data = '파일 명에 . / \ * < > | : ? 가 들어 갈 수 없습니다.'), 401
    elif(num == 16):
        return web_render('index.html', login = login_check(), title = '업로드 오류', logo = set_data['name'], data = '동일한 이름의 파일이 있습니다.'), 401
    elif(num == 17):
        return web_render('index.html', login = login_check(), title = '편집 오류', logo = set_data['name'], data = '편집 내용 기록에는 한글과 영어와 숫자, 공백만 허용 됩니다.'), 401
    elif(num == 18):
        return web_render('index.html', login = login_check(), title = '편집 오류', logo = set_data['name'], data = '내용이 원래 문서와 동일 합니다.'), 401
    elif(num == 19):
        return web_render('index.html', login = login_check(), title = '이동 오류', logo = set_data['name'], data = '이동 하려는 곳에 문서가 이미 있습니다.'), 401
    elif(num == 20):
        return web_render('index.html', login = login_check(), title = '비밀번호 오류', logo = set_data['name'], data = '재 확인이랑 비밀번호가 다릅니다.'), 401
    else:
        return '<meta http-equiv="refresh" content="0;url=/" />'

@app.errorhandler(404)
def uncaughtError(error):
    return '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(set_data['frontpage']) + '" />'

@app.errorhandler(413)
def uncaughtError(error):
    app.config['MAX_CONTENT_LENGTH'] = (1024**3)
    return error, 401
    
if(__name__ == '__main__'):
    app.run(host = '0.0.0.0', port = int(set_data['port']))