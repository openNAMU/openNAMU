from flask import Flask, session, request

from urllib import parse
import json
import pymysql
import time
import re
import json

json_data = open('set.json').read()
data = json.loads(json_data)

conn = pymysql.connect(host = data['host'], user = data['user'], password = data['pw'], charset = 'utf8mb4')
curs = conn.cursor(pymysql.cursors.DictCursor)

def DB_갱신():
    conn.commit()

def URL_인코딩(데이터):
    return parse.quote(데이터).replace('/','%2F')
    
def DB_가져오기():
    return curs.fetchall()

DB_실행 = curs.execute
DB_인코딩 = pymysql.escape_string

DB_실행("use " + data['db'])

def 비교(seqm):
    output= []
    for opcode, a0, a1, b0, b1 in seqm.get_opcodes():
        if(opcode == 'equal'):
            output.append(seqm.a[a0:a1])
        elif(opcode == 'insert'):
            output.append("<span style='background:#CFC;'>" + seqm.b[b0:b1] + "</span>")
        elif(opcode == 'delete'):
            output.append("<span style='background:#FDD;'>" + seqm.a[a0:a1] + "</span>")
        elif(opcode == 'replace'):
            output.append("<span style='background:#CFC;'>" + seqm.b[b0:b1] + "</span><span style='background:#FDD;'>" + seqm.a[a0:a1] + "</span>")
        else:
            output.append(seqm.a[a0:a1])
    return ''.join(output)
           
def 관리자_확인():
    if(session.get('Now') == True):
        아이피 = 아이피_확인()
        DB_실행("select * from user where id = '" + DB_인코딩(아이피) + "'")
        사용자_자료 = DB_가져오기()
        if(사용자_자료):
            if(사용자_자료[0]['acl'] == 'owner' or 사용자_자료[0]['acl'] == 'admin'):
                return 1
                
def 소유자_확인():
    if(session.get('Now') == True):
        아이피 = 아이피_확인()
        DB_실행("select * from user where id = '" + DB_인코딩(아이피) + "'")
        사용자_자료 = DB_가져오기()
        if(사용자_자료):
            if(사용자_자료[0]['acl'] == 'owner'):
                return 1
                
def 틀_확인(이름, 데이터):
    if(re.search('^틀:', 이름)):
        DB_실행("select * from back where title = '" + DB_인코딩(이름) + "' and type = 'include'")
        틀_역링크 = DB_가져오기()
        if(틀_역링크):
            숫자 = 0

            while(True):
                try:
                    나무마크(틀_역링크[숫자]['link'], 데이터)
                except:
                    break
                    
                숫자 += 1
    
def 로그인_확인():
    if(session.get('Now') == True):
        return 1
    else:
        return 0

def 아이디_파싱(원래_아이디):
    있나 = re.search("([^-]*)\s\-\s(Close|Reopen|Stop|Restart|Admin|Agreement|Settlement)$", 원래_아이디)
    if(있나):
        분리 = 있나.groups()
        
        DB_실행("select * from data where title = '사용자:" + DB_인코딩(분리[0]) + "'")
        row = DB_가져오기()
        if(row):
            ip = '<a href="/w/' + URL_인코딩('사용자:' + 분리[0]) + '">' + 분리[0] + '</a> - ' + 분리[1] + ' <a href="/record/' + URL_인코딩(분리[0]) + '/n/1">(기록)</a>'
        else:
            ip = '<a class="not_thing" href="/w/' + URL_인코딩('사용자:' + 분리[0]) + '">' + 분리[0] + '</a> - ' + 분리[1] + ' <a href="/record/' + URL_인코딩(분리[0]) + '/n/1">(기록)</a>'
    elif(re.search("\.", 원래_아이디)):
        ip = 원래_아이디 + ' <a href="/record/' + URL_인코딩(원래_아이디) + '/n/1">(기록)</a>'
    else:
        DB_실행("select * from data where title = '사용자:" + DB_인코딩(원래_아이디) + "'")
        row = DB_가져오기()
        if(row):
            ip = '<a href="/w/' + URL_인코딩('사용자:' + 원래_아이디) + '">' + 원래_아이디 + '</a> <a href="/record/' + URL_인코딩(원래_아이디) + '/n/1">(기록)</a>'
        else:
            ip = '<a class="not_thing" href="/w/' + URL_인코딩('사용자:' + 원래_아이디) + '">' + 원래_아이디 + '</a> <a href="/record/' + URL_인코딩(원래_아이디) + '/n/1">(기록)</a>'

    return ip

def 아이피_확인():
    if(session.get('Now') == True):
        아이피 = format(session['DREAMER'])
    else:
        if(request.headers.getlist("X-Forwarded-For")):
            아이피 = request.headers.getlist("X-Forwarded-For")[0]
        else:
            아이피 = request.remote_addr
            
    return 아이피

def ACL_체크(ip, name):
    m = re.search("^사용자:(.*)", name)
    n = re.search("^파일:(.*)", name)
    if(m):
        g = m.groups()
        if(ip == g[0]):
            if(re.search("\.", g[0])):
                return 1
            else:
                DB_실행("select * from ban where block = '" + DB_인코딩(ip) + "'")
                rows = DB_가져오기()
                if(rows):
                    return 1
                else:
                    return 0
        else:
            return 1
    elif(n):
        if(not 소유자_확인() == 1):
            return 1
    else:
        b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
        if(b):
            results = b.groups()
            DB_실행("select * from ban where block = '" + DB_인코딩(results[0]) + "' and band = 'O'")
            rowss = DB_가져오기()
            if(rowss):
                return 1
            else:
                DB_실행("select * from ban where block = '" + DB_인코딩(ip) + "'")
                rows = DB_가져오기()
                if(rows):
                    return 1
                else:
                    DB_실행("select * from data where title = '" + DB_인코딩(name) + "'")
                    row = DB_가져오기()
                    if(row):
                        DB_실행("select * from user where id = '" + DB_인코딩(ip) + "'")
                        rows = DB_가져오기()
                        if(row[0]['acl'] == 'user'):
                            if(rows):
                                return 0
                            else:
                                return 1
                        elif(row[0]['acl'] == 'admin'):
                            if(rows):
                                if(rows[0]['acl'] == 'admin' or rows[0]['acl'] == 'owner'):
                                    return 0
                                else:
                                    return 1
                            else:
                                return 1
                        else:
                            return 0
                    else:
                        return 0
        else:
            DB_실행("select * from ban where block = '" + DB_인코딩(ip) + "'")
            rows = DB_가져오기()
            if(rows):
                return 1
            else:
                DB_실행("select * from data where title = '" + DB_인코딩(name) + "'")
                row = DB_가져오기()
                if(row):
                    DB_실행("select * from user where id = '" + DB_인코딩(ip) + "'")
                    rows = DB_가져오기()
                    if(row[0]['acl'] == 'user'):
                        if(rows):
                            return 0
                        else:
                            return 1
                    elif(row[0]['acl'] == 'admin'):
                        if(rows):
                            if(rows[0]['acl'] == 'admin' or rows[0]['acl'] == 'owner'):
                                return 0
                            else:
                                return 1
                        else:
                            return 1
                    else:
                        return 0
                else:
                    return 0

def 차단_체크(ip):
    b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
    if(b):
        results = b.groups()
        DB_실행("select * from ban where block = '" + DB_인코딩(results[0]) + "' and band = 'O'")
        rowss = DB_가져오기()
        if(rowss):
            return 1
        else:
            DB_실행("select * from ban where block = '" + DB_인코딩(ip) + "'")
            rows = DB_가져오기()
            if(rows):
                return 1
            else:
                return 0
    else:
        DB_실행("select * from ban where block = '" + DB_인코딩(ip) + "'")
        rows = DB_가져오기()
        if(rows):
            return 1
        else:
            return 0
        
def 토론자_체크(ip, name, sub):
    b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
    if(b):
        results = b.groups()
        DB_실행("select * from ban where block = '" + DB_인코딩(results[0]) + "' and band = 'O'")
        rowss = DB_가져오기()
        if(rowss):
            return 1
        else:
            DB_실행("select * from ban where block = '" + DB_인코딩(ip) + "'")
            rows = DB_가져오기()
            if(rows):
                return 1
            else:
                DB_실행("select * from stop where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "'")
                rows = DB_가져오기()
                if(rows):
                    return 1
                else:
                    return 0
    else:
        DB_실행("select * from ban where block = '" + DB_인코딩(ip) + "'")
        rows = DB_가져오기()
        if(rows):
            return 1
        else:
            DB_실행("select * from stop where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "'")
            rows = DB_가져오기()
            if(rows):
                return 1
            else:
                return 0

def 시간():
    now = time.localtime()
    s = "%04d-%02d-%02d %02d:%02d:%02d" % (now.tm_year, now.tm_mon, now.tm_mday, now.tm_hour, now.tm_min, now.tm_sec)
    return s

def 최근_토론_추가(title, sub, date):
    DB_실행("select * from rd where title = '" + DB_인코딩(title) + "' and sub = '" + DB_인코딩(sub) + "'")
    최근_토론 = DB_가져오기()
    if(최근_토론):
        DB_실행("update rd set date = '" + DB_인코딩(date) + "' where title = '" + DB_인코딩(title) + "' and sub = '" + DB_인코딩(sub) + "'")
    else:
        DB_실행("insert into rd (title, sub, date) value ('" + DB_인코딩(title) + "', '" + DB_인코딩(sub) + "', '" + DB_인코딩(date) + "')")
    DB_갱신()
    
def 최근_차단_추가(block, end, today, blocker, why):
    DB_실행("insert into rb (block, end, today, blocker, why) value ('" + DB_인코딩(block) + "', '" + DB_인코딩(end) + "', '" + today + "', '" + DB_인코딩(blocker) + "', '" + DB_인코딩(why) + "')")
    DB_갱신()

def 역사_추가(title, data, date, ip, send, leng):
    DB_실행("select * from history where title = '" + DB_인코딩(title) + "' order by id+0 desc limit 1")
    rows = DB_가져오기()
    if(rows):
        number = int(rows[0]['id']) + 1
        DB_실행("insert into history (id, title, data, date, ip, send, leng) value ('" + str(number) + "', '" + DB_인코딩(title) + "', '" + DB_인코딩(data) + "', '" + date + "', '" + DB_인코딩(ip) + "', '" + DB_인코딩(send) + "', '" + leng + "')")
        DB_갱신()
    else:
        DB_실행("insert into history (id, title, data, date, ip, send, leng) value ('1', '" + DB_인코딩(title) + "', '" + DB_인코딩(data) + "', '" + date + "', '" + DB_인코딩(ip) + "', '" + DB_인코딩(send + ' (새 문서)') + "', '" + leng + "')")
        DB_갱신()

def 길이_확인(기존, 바뀜):
    if(기존 < 바뀜):
        길이 = 바뀜 - 기존
        길이 = '+' + str(길이)
    elif(바뀜 < 기존):
        길이 = 기존 - 바뀜
        길이 = '-' + str(길이)
    else:
        길이 = '0'
        
    return 길이