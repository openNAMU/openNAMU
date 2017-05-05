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

from func import *
from mark import *

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
                        
                        return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩('파일:' + filename) + '" />'
                else:
                    return '<meta http-equiv="refresh" content="0;url=/error/15" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/14" />'
    else:        
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            return web_render('index.html', login = login_check(), logo = set_data['name'], title = '업로드', tn = 21, number = set_data['upload'])

