import urllib.parse
import datetime
import hashlib
import flask
import string
import re

import os
import html
import json
import threading

set_data = ''

def get_time():
    return str(datetime.datetime.today().strftime("%Y-%m-%d %H:%M:%S"))

def db_data_get(data):
    global set_data
    
    set_data = data

def db_change(data):
    if set_data == 'mysql':
        data = data.replace('random()', 'rand()')
        data = data.replace('%', '%%')
        data = data.replace('?', '%s')
        data = data.replace('collate nocase', 'collate utf8mb4_general_ci')

    return data

def ip_check(d_type = 0):
    ip = ''
    if d_type == 0 and (flask.session and 'id' in flask.session):
        ip = flask.session['id']
    else:        
        ip = flask.request.environ.get('HTTP_X_REAL_IP',
            flask.request.environ.get('HTTP_CF_CONNECTING_IP',
                flask.request.environ.get('REMOTE_ADDR',
                    '::1'
                )
            )
        )

    return ip

def ip_or_user(data = ''):
    # without_DB

    # 1 == ip
    # 0 == reg
    
    if data == '':
        data = ip_check()

    if re.search(r'(\.|:)', data):
        return 1
    else:
        return 0

def url_pas(data):
    data = re.sub(r'^\.', '\\\\.', data)
    data = urllib.parse.quote(data)
    data = data.replace('/','%2F')

    return data

def sha224_replace(data):
    return hashlib.sha224(bytes(data, 'utf-8')).hexdigest()

def md5_replace(data):
    return hashlib.md5(data.encode()).hexdigest()

def get_main_skin_set(curs, flask_session, set_name, ip):
    if ip_or_user(ip) == 0:
        curs.execute(db_change('select data from user_set where name = ? and id = ?'), [set_name, ip])
        db_data = curs.fetchall()
        set_data = db_data[0][0] if db_data and db_data[0][0] != '' else 'default'
    else:
        set_data = flask_session[set_name] if set_name in flask_session and flask_session[set_name] != '' else 'default'

    if set_data == 'default':
        curs.execute(db_change('select data from other where name = ?'), [set_name])
        db_data = curs.fetchall()
        set_data = db_data[0][0] if db_data and db_data[0][0] != '' else 'default'

    return set_data