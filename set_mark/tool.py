import flask
import urllib.parse
import datetime
import re
import hashlib

def get_time():
    return str(datetime.datetime.today().strftime("%Y-%m-%d %H:%M:%S"))
    
def ip_check():
    if flask.session and ('Now' and 'DREAMER') in flask.session and flask.session['Now'] == 1:
        ip = flask.session['DREAMER']
    else:
        try:
            ip = flask.request.environ.get('HTTP_X_REAL_IP', flask.request.remote_addr)
            
            print(flask.request.environ.get('HTTP_X_FORWARDED_FOR', flask.request.remote_addr))
            if ip == ('::1' or '127.0.0.1'):
                ip = flask.request.environ.get('HTTP_X_FORWARDED_FOR', flask.request.remote_addr)
        except:
            ip = 'Error'

    return str(ip)

def savemark(data):
    data = re.sub("\[date\(now\)\]", get_time(), data)
    
    ip = ip_check()
    if not re.search("\.", ip):
        name = '[[사용자:' + ip + '|' + ip + ']]'
    else:
        name = ip
        
    data = re.sub("\[name\]", name, data)

    return data

def url_pas(data):
    return urllib.parse.quote(data).replace('/','%2F')

def sha224(data):
    return hashlib.sha224(bytes(data, 'utf-8')).hexdigest()

def md5_replace(data):
    return hashlib.md5(data.encode()).hexdigest()