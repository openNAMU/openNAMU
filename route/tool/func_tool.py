import urllib.parse
import datetime
import hashlib
import flask
import random
import string
import re

import os
import html
import json   
import sqlite3
import time
import threading

try:
  import orjson
except:
  import json as orjson

def get_time():
    return str(datetime.datetime.today().strftime("%Y-%m-%d %H:%M:%S"))

class class_temp_db:
    def __enter__(self):
        self.conn = sqlite3.connect(
            os.path.join('.', 'data', 'temp.db'),
            check_same_thread = False,
            isolation_level = None
        )

        return self.conn

    def __exit__(self, exc_type, exc_value, traceback):
        self.conn.commit()
        self.conn.close()

def db_change(data):
    with class_temp_db() as m_conn:
        m_curs = m_conn.cursor()
        
        m_curs.execute('select data from temp where name = "db_type"')
        db_data = m_curs.fetchall()
        set_data = db_data[0][0] if db_data else 'sqlite'

        if set_data == 'mysql':
            data = data.replace('random()', 'rand()')
            data = data.replace('%', '%%')
            data = data.replace('?', '%s')
            data = data.replace('collate nocase', 'collate utf8mb4_general_ci')

        return data

def ip_check(d_type = 0):
    ip = '::1'
    if d_type == 0 and (flask.session and 'id' in flask.session):
        ip = flask.session['id']
    else:
        with class_temp_db() as m_conn:
            m_curs = m_conn.cursor()

            m_curs.execute('select data from temp where name = "load_ip_select"')
            db_data = m_curs.fetchall()
            set_data = db_data[0][0] if db_data else 'default'
        
            if set_data == "default":
                ip = flask.request.environ.get('HTTP_X_REAL_IP',
                    flask.request.environ.get('HTTP_CF_CONNECTING_IP',
                        flask.request.environ.get('REMOTE_ADDR',
                            '::1'
                        )
                    )
                )
            else:
                ip = flask.request.environ.get(set_data, '::1')
            
            if ip_or_user(ip) == 0:
                ip = '::1'

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

def get_main_skin_set(conn, flask_session, set_name, ip):
    curs = conn.cursor()

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