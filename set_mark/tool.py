import flask
import urllib.parse
import datetime
import re
import hashlib

def get_time():
    return str(datetime.datetime.today().strftime("%Y-%m-%d %H:%M:%S"))
    
def ip_check():
    if flask.session and ('state' and 'id') in flask.session and flask.session['state'] == 1:
        ip = flask.session['id']
    else:
        try:
            ip = flask.request.environ.get('HTTP_X_REAL_IP', flask.request.environ.get('HTTP_X_FORWARDED_FOR', flask.request.remote_addr))
            
            if ip == ('::1' or '127.0.0.1'):
                ip = flask.request.environ.get('HTTP_X_FORWARDED_FOR', flask.request.remote_addr)
        except:
            ip = '-'
            
    return str(ip)

def savemark(data):
    data = re.sub("\[date\(now\)\]", get_time(), data)
    
    ip = ip_check()
    if not re.search("\.", ip):
        name = '[[user:' + ip + '|' + ip + ']]'
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

def xss_protect(curs, data, ok_list = []):
    curs.execute('select html from html_filter where kind = ""')
    html_db = curs.fetchall()

    src_list = ["www.youtube.com", "serviceapi.nmv.naver.com", "tv.kakao.com", "www.google.com", "serviceapi.rmcnmv.naver.com"]
    html_list = ['div', 'span', 'embed', 'iframe', 'ruby', 'rp', 'rt'] + ok_list
    
    html_data = re.findall('&lt;(\/)?((?:(?!&gt;| ).)+)( (?:(?:(?!&gt;).)+)?)?&gt;', data)
    for in_data in html_data:
        if in_data[0] == '':
            if in_data[1] in html_list or (html_db and in_data[1] in html_db[0]):
                if re.search('&lt;\/' + in_data[1] + '&gt;', data):
                    src = re.search('src=([^ ]*)', in_data[2])
                    if src:
                        v_src = re.search('http(?:s)?:\/\/([^/\'" ]*)', src.groups()[0])
                        if v_src:
                            if not v_src.groups()[0] in src_list:
                                and_data = re.sub('&#x27;', '\'', re.sub('&quot;', '"', re.sub('src=([^ ]*)', '', in_data[2])))
                            else:
                                and_data = re.sub('&#x27;', '\'', re.sub('&quot;', '"', in_data[2]))
                        else:
                            and_data = re.sub('&#x27;', '\'', re.sub('&quot;', '"', re.sub('src=([^ ]*)', '', in_data[2])))
                    else:
                        and_data = re.sub('&#x27;', '\'', re.sub('&quot;', '"', in_data[2]))
                        
                    data = data.replace('&lt;' + in_data[1] + in_data[2] + '&gt;', '<' + in_data[1] + and_data + '>', 1)
                    data = re.sub('&lt;\/' + in_data[1] + '&gt;', '</' + in_data[1] + '>', data, 1)

    position = re.compile('position', re.I)
    data = position.sub('', data)

    return data