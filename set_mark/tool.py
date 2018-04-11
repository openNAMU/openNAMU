from flask import session, request

from urllib import parse
import time
import datetime
import re
import hashlib

def get_time():
    now = time.localtime()
    date = "%04d-%02d-%02d %02d:%02d:%02d" % (now.tm_year, now.tm_mon, now.tm_mday, now.tm_hour, now.tm_min, now.tm_sec)

    return date
    
def ip_check():
    if session and ('Now' and 'DREAMER') in session and session['Now'] == 1:
        ip = session['DREAMER']
    else:
        try:
            ip = request.environ.get('HTTP_X_REAL_IP', request.remote_addr)
        except:
            ip = 'None'

    return str(ip)

def savemark(data):
    data = re.sub("\[date\(now\)\]", get_time(), data)
    
    if not re.search("\.", ip_check()):
        name = '[[사용자:' + ip_check() + '|' + ip_check() + ']]'
    else:
        name = ip_check()
        
    data = re.sub("\[name\]", name, data)

    return data

def url_pas(data):
    return parse.quote(data).replace('/','%2F')

def sha224(data):
    return hashlib.sha224(bytes(data, 'utf-8')).hexdigest()