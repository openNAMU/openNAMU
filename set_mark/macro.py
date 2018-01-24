from flask import session, request

from urllib import parse
import time
import datetime
import re
import json

def get_time():
    now = time.localtime()
    date = "%04d-%02d-%02d %02d:%02d:%02d" % (now.tm_year, now.tm_mon, now.tm_mday, now.tm_hour, now.tm_min, now.tm_sec)

    return(date)
    
def ip_check():
    if(('Now' and 'DREAMER') in session and session['Now'] == 1):
        ip = session['DREAMER']
    else:
        ip = request.environ.get('HTTP_X_REAL_IP', request.remote_addr)

    return(str(ip))

def savemark(data):
    data = re.sub("\[date\(now\)\]", get_time(), data)
    
    if(not re.search("\.", ip_check())):
        name = '[[사용자:' + ip_check() + '|' + ip_check() + ']]'
    else:
        name = ip_check()
        
    data = re.sub("\[name\]", name, data)

    return(data)

def macro(data):      
    data = savemark(data)
    data = re.sub("\[anchor\((?P<in>[^\[\]]*)\)\]", '<span id="\g<in>"></span>', data)          
    data = re.sub("\[nicovideo\((?P<in>[^,)]*)(?:(?:,(?:[^,)]*))+)?\)\]", "[[http://embed.nicovideo.jp/watch/\g<in>]]", data)
    data = re.sub('\[ruby\((?P<in>[^\,]*)\,\s?(?P<out>[^\)]*)\)\]', '<ruby>\g<in><rp>(</rp><rt>\g<out></rt><rp>)</rp></ruby>', data)
    data = re.sub("\[br\]", '<br>', data)
    
    while(1):
        com = re.compile("\[(youtube|kakaotv)\(([^, )]*)(,[^)]*)?\)\]")
        m = com.search(data)
        if(m):
            src = ''
            width = '560'
            height = '315'
            time = '0'
            
            result = m.groups()
            if(result[1]):
                yudt = re.search('(?:\?v=(.*)|\/([^/?]*)|^([a-zA-Z0-9\-_]*))$', result[1])
                if(yudt):
                    if(yudt.groups()[0]):
                        src = yudt.groups()[0]
                    elif(yudt.groups()[1]):
                        src = yudt.groups()[1]
                    elif(yudt.groups()[2]):
                        src = yudt.groups()[2]
                else:
                    src = ''
                    
            if(result[2]):
                mdata = re.search('width=([0-9%]*)', result[2])
                if(mdata):
                    width = mdata.groups()[0]
                
                mdata = re.search('height=([0-9%]*)', result[2])
                if(mdata):
                    height = mdata.groups()[0]
                    
                mdata = re.search('start=([0-9]*)', result[2])
                if(mdata):
                    time = mdata.groups()[0]

            if(result[0] == 'youtube'):
                data = com.sub('<iframe width="' + width + '" height="' + height + '" src="https://www.youtube.com/embed/' + src + '?start=' + time + '" frameborder="0" allowfullscreen></iframe><br>', data, 1)
            else:
                data = com.sub('<iframe width="' + width + '" height="' + height + '" src="https://tv.kakao.com/embed/player/cliplink/' + src + '?service=kakao_tv&start=' + time + '" allowfullscreen frameborder="0" scrolling="no"></iframe><br>', data, 1)
        else:
            break
    
    now_time = get_time()
    data = re.sub('\[date\]', now_time, data)
    
    time_data = re.search('^([0-9]{4}-[0-9]{2}-[0-9]{2})', now_time)
    time = time_data.groups()
    
    age_data = re.findall('\[age\(([0-9]{4}-[0-9]{2}-[0-9]{2})\)\]', data)
    for age in age_data:
        old = datetime.datetime.strptime(time[0], '%Y-%m-%d')
        will = datetime.datetime.strptime(age, '%Y-%m-%d')
        e_data = old - will

        data = re.sub('\[age\(([0-9]{4})-([0-9]{2})-([0-9]{2})\)\]', str(int(int(e_data.days) / 365)), data, 1)

    dday_data = re.findall('\[dday\(([0-9]{4}-[0-9]{2}-[0-9]{2})\)\]', data)
    for dday in dday_data:
        old = datetime.datetime.strptime(time[0], '%Y-%m-%d')
        will = datetime.datetime.strptime(dday, '%Y-%m-%d')
        e_data = old - will

        if(re.search('^-', str(e_data.days))):
            e_day = str(e_data.days)
        else:
            e_day = '+' + str(e_data.days)

        data = re.sub('\[dday\(([0-9]{4}-[0-9]{2}-[0-9]{2})\)\]', e_day, data, 1)
        
    return(data)