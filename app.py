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
data = json.loads(json_data)

print('오픈나무 시작 포트 : ' + data['port'])

import logging
log = logging.getLogger('werkzeug')
log.setLevel(logging.ERROR)

def 시작():
    try:
        DB_실행("select * from data limit 1")
    except:
        DB_실행("create table data(title text, data longtext, acl text)")
    
    try:
        DB_실행("select * from history limit 1")
    except:
        DB_실행("create table 역사_추가(id text, title text, data longtext, date text, ip text, send text, leng text)")
    
    try:
        DB_실행("select * from rd limit 1")
    except:
        DB_실행("create table rd(title text, sub text, date text)")
    
    try:
        DB_실행("select * from user limit 1")
    except:
        DB_실행("create table user(id text, pw text, acl text)")
    
    try:
        DB_실행("select * from ban limit 1")
    except:
        DB_실행("create table ban(block text, end text, why text, band text)")
    
    try:
        DB_실행("select * from topic limit 1")
    except:
        DB_실행("create table topic(id text, title text, sub text, data longtext, date text, ip text, block text)")
    
    try:
        DB_실행("select * from stop limit 1")
    except:
        DB_실행("create table stop(title text, sub text, close text)")
    
    try:
        DB_실행("select * from rb limit 1")
    except:
        DB_실행("create table rb(block text, end text, today text, blocker text, why text)")
    
    try:
        DB_실행("select * from login limit 1")
    except:
        DB_실행("create table login(user text, ip text, today text)")
    
    try:
        DB_실행("select * from back limit 1")
    except:
        DB_실행("create table back(title text, link text, type text)")
    
    try:
        DB_실행("select * from cat limit 1")
    except:
        DB_실행("create table cat(title text, cat text)")
        
    try:
        DB_실행("select * from hidhi limit 1")
    except:
        DB_실행("create table hidhi(title text, re text)")

    try:
        DB_실행("select * from distop limit 1")
    except:
        DB_실행("create table distop(id text, title text, sub text)") 

    try:
        DB_실행("select * from agreedis limit 1")
    except:
        DB_실행("create table agreedis(title text, sub text)") 

conn = pymysql.connect(host = data['host'], user = data['user'], password = data['pw'], charset = 'utf8mb4')
curs = conn.cursor(pymysql.cursors.DictCursor)

def DB_갱신():
    conn.commit()

def URL_인코딩(데이터):
    return parse.quote(데이터).replace('/','%2F')
    
def DB_가져오기():
    return curs.fetchall()
    
웹_디자인 = render_template
DB_실행 = curs.execute
DB_인코딩 = pymysql.escape_string

try:
    DB_실행("use " + data['db'])
except:
    DB_실행("create database " + data['db'])
    DB_실행("use " + data['db'])
    DB_실행("alter database " + data['db'] + " character set = utf8mb4 collate = utf8mb4_unicode_ci")

시작()

app.secret_key = hashlib.sha512(bytes(data['key'], 'ascii')).hexdigest()

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
        ip = 아이피_확인(request)
        DB_실행("select * from user where id = '" + DB_인코딩(ip) + "'")
        rows = DB_가져오기()
        if(rows):
            if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
                return 1
                
def 소유자_확인():
    if(session.get('Now') == True):
        ip = 아이피_확인(request)
        DB_실행("select * from user where id = '" + DB_인코딩(ip) + "'")
        rows = DB_가져오기()
        if(rows):
            if(rows[0]['acl'] == 'owner'):
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

def 세이브마크(데이터):
    데이터 = re.sub("\[date\(now\)\]", 시간(), 데이터)
    if(not re.search("\.", 아이피_확인(request))):
        이름 = '[[사용자:' + 아이피_확인(request) + '|' + 아이피_확인(request) + ']]'
    else:
        이름 = 아이피_확인(request)
    데이터 = re.sub("\[name\]", 이름, 데이터)

    return 데이터
    
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
    
def HTML_파싱(데이터):
    while(True):
        있나 = re.search("<((div|span|embed|iframe)(?:[^>]*))>", 데이터)
        
        if(있나):
            분리 = 있나.groups()

            if(re.search("<(\/" + 분리[1] + ")>", 데이터) and not re.search("'", 분리[0])):
                XSS = re.search('src="http(?:s)?:\/\/([^\/]*)\/(?:[^"]*)"', 분리[0])
                
                if(XSS):
                    확인 = XSS.groups()
                    
                    if(확인[0] == "www.youtube.com" or 확인[0] == "serviceapi.nmv.naver.com" or 확인[0] == "tv.kakao.com" or 확인[0] == "tvple.com" or 확인[0] == "tvpot.daum.net"):
                        임시_저장 = 분리[0]
                    else:
                        임시_저장 = re.sub('src="([^"]*)"', '', 분리[0])
                else:
                    임시_저장 = 분리[0]
                
                임시_저장 = re.sub('"', '#.#', 임시_저장)
                데이터 = re.sub("<((?:\/)?" + 분리[1] + "(?:[^>]*))>", "[" + 임시_저장 + "]", 데이터, 1)
                데이터 = re.sub("<\/" + 분리[1] + ">", "[/" + 분리[1] + "]", 데이터, 1)
            else:
                데이터 = re.sub("<((?:\/)?" + 분리[1] + "(?:[^>]*))>", '&lt;' + 분리[0] + '&gt;', 데이터, 1)
                
                break
        else:
            break

    데이터 = re.sub('<', '&lt;', 데이터)
    데이터 = re.sub('>', '&gt;', 데이터)
    데이터 = re.sub('"', '&quot;', 데이터)
    
    데이터 = re.sub("\[(?P<in>(?:\/)?(?:div|span|embed|iframe)(?:[^\]]*))\]", "<\g<in>>", 데이터)
    데이터 = re.sub('#.#', '"', 데이터)
    
    return 데이터
    
def 중괄호_문법(데이터, 접기_숫자, 틀):
    while(True):
        문법_컴파일 = re.compile("{{{((?:(?!{{{)(?!}}}).)*)}}}", re.DOTALL)
        있나 = 문법_컴파일.search(데이터)
        
        if(있나):
            분리 = 있나.groups()
            
            크게_문법 = re.compile("^\+([1-5])\s(.*)$", re.DOTALL)
            크게 = 크게_문법.search(분리[0])
            
            작게_문법 = re.compile("^\-([1-5])\s(.*)$", re.DOTALL)
            작게 = 작게_문법.search(분리[0])
            
            색깔_문법_1 = re.compile("^(#[0-9a-f-A-F]{6})\s(.*)$", re.DOTALL)
            색깔_1 = 색깔_문법_1.search(분리[0])
            
            색깔_문법_2 = re.compile("^(#[0-9a-f-A-F]{3})\s(.*)$", re.DOTALL)
            색깔_2 = 색깔_문법_2.search(분리[0])
            
            색깔_문법_3 = re.compile("^#(\w+)\s(.*)$", re.DOTALL)
            색깔_3 = 색깔_문법_3.search(분리[0])
            
            배경색_문법_1 = re.compile("^@([0-9a-f-A-F]{6})\s(.*)$", re.DOTALL)
            배경색_1 = 배경색_문법_1.search(분리[0])
            
            배경색_문법_2 = re.compile("^@([0-9a-f-A-F]{3})\s(.*)$", re.DOTALL)
            배경색_2 = 배경색_문법_2.search(분리[0])
            
            배경색_문법_3 = re.compile("^@(\w+)\s(.*)$", re.DOTALL)
            배경색_3 = 배경색_문법_3.search(분리[0])
            
            틀_제외_문법 = re.compile("^#!noin\s(.*)$", re.DOTALL)
            틀_제외 = 틀_제외_문법.search(분리[0])
            
            DIV_문법 = re.compile("^#!wiki\sstyle=&quot;((?:(?!&quot;|\n).)*)&quot;\n?\s\n(.*)$", re.DOTALL)
            DIV = DIV_문법.search(분리[0])
            
            HTML_문법 = re.compile("^#!html\s(.*)$", re.DOTALL)
            HTML = HTML_문법.search(분리[0])
            
            접기_문법 = re.compile("^#!folding\s((?:(?!\n).)*)\n?\s\n(.*)$", re.DOTALL)
            접기 = 접기_문법.search(분리[0])
            
            if(크게):
                결과 = 크게.groups()
                데이터 = 문법_컴파일.sub('<span class="font-size-' + 결과[0] + '">' + 결과[1] + '</span>', 데이터, 1)
            elif(작게):
                결과 = 작게.groups()
                데이터 = 문법_컴파일.sub('<span class="font-size-small-' + 결과[0] + '">' + 결과[1] + '</span>', 데이터, 1)
            elif(색깔_1):
                결과 = 색깔_1.groups()
                data = 문법_컴파일.sub('<span style="color:' + 결과[0] + '">' + 결과[1] + '</span>', 데이터, 1)
            elif(색깔_2):
                결과 = 색깔_2.groups()
                데이터 = 문법_컴파일.sub('<span style="color:' + 결과[0] + '">' + 결과[1] + '</span>', 데이터, 1)
            elif(색깔_3):
                결과 = 색깔_3.groups()
                데이터 = 문법_컴파일.sub('<span style="color:' + 결과[0] + '">' + 결과[1] + '</span>', 데이터, 1)
            elif(배경색_1):
                결과 = 배경색_1.groups()
                데이터 = 문법_컴파일.sub('<span style="background:#' + 결과[0] + '">' + 결과[1] + '</span>', 데이터, 1)
            elif(배경색_2):
                결과 = 배경색_2.groups()
                데이터 = 문법_컴파일.sub('<span style="background:#' + 결과[0] + '">' + 결과[1] + '</span>', 데이터, 1)
            elif(배경색_3):
                결과 = 배경색_3.groups()
                데이터 = 문법_컴파일.sub('<span style="background:' + 결과[0] + '">' + 결과[1] + '</span>', 데이터, 1)
            elif(DIV):
                결과 = DIV.groups()
                데이터 = 문법_컴파일.sub('<div style="' + 결과[0] + '">' + 결과[1] + '</div>', 데이터, 1)
            elif(HTML):
                결과 = HTML.groups()
                데이터 = 문법_컴파일.sub(결과[0], 데이터, 1)
            elif(접기):
                결과 = 접기.groups()
                데이터 = 문법_컴파일.sub("<div>" + 결과[0] + "<span style='float:right;'><div id='folding_" + str(접기_숫자 + 1) + "' style='display:block;'>[<a href='javascript:void(0);' onclick='var f=document.getElementById(\"folding_" + str(접기_숫자) + "\");var s=f.style.display==\"block\";f.style.display=s?\"none\":\"block\";this.className=s?\"\":\"opened\";var f=document.getElementById(\"folding_" + str(접기_숫자 + 1) + "\");var s=f.style.display==\"none\";f.style.display=s?\"block\":\"none\";var f=document.getElementById(\"folding_" + str(접기_숫자 + 2) + "\");var s=f.style.display==\"block\";f.style.display=s?\"none\":\"block\";'>펼치기</a>]</div><div id='folding_" + str(접기_숫자 + 2) + "' style='display:none;'>[<a href='javascript:void(0);' onclick='var f=document.getElementById(\"folding_" + str(접기_숫자) + "\");var s=f.style.display==\"block\";f.style.display=s?\"none\":\"block\";this.className=s?\"\":\"opened\";var f=document.getElementById(\"folding_" + str(접기_숫자 + 1) + "\");var s=f.style.display==\"none\";f.style.display=s?\"block\":\"none\";var f=document.getElementById(\"folding_" + str(접기_숫자 + 2) + "\");var s=f.style.display==\"block\";f.style.display=s?\"none\":\"block\";'>접기</a>]</div></a></span><div id='folding_" + str(접기_숫자) + "' style='display:none;'><br>" + 결과[1] + "</div></div>", 데이터, 1)
                
                접기_숫자 += 3
            elif(HTML):
                결과 = HTML.groups()
                데이터 = 문법_컴파일.sub(결과[0], 데이터, 1)
            elif(틀_제외):
                if(틀 == True):
                    데이터 = 문법_컴파일.sub("", 데이터, 1)
                else:
                    결과 = 틀_제외.groups()
                    데이터 = 문법_컴파일.sub(결과[0], 데이터, 1)
            else:
                데이터 = 문법_컴파일.sub('<code>' + 분리[0] + '</code>', 데이터, 1)
        else:
            break
            
    while(True):
        문법_컴파일 = re.compile("<code>(((?!<\/code>).)*)<\/code>", re.DOTALL)
        있나 = 문법_컴파일.search(데이터)
        if(있나):
            결과 = 있나.groups()
            
            중간_데이터 = re.sub("<\/span>", "}}}", 결과[0])
            중간_데이터 = re.sub("<\/div>", "}}}", 중간_데이터)
            중간_데이터 = re.sub('<span class="font\-size\-(?P<in>[1-6])">', "{{{+\g<in> ", 중간_데이터)
            중간_데이터 = re.sub('<span class="font\-size\-small\-(?P<in>[1-6])">', "{{{-\g<in> ", 중간_데이터)
            중간_데이터 = re.sub('<span style="color:(?:#)?(?P<in>[^"]*)">', "{{{#\g<in> ", 중간_데이터)
            중간_데이터 = re.sub('<span style="background:(?:#)?(?P<in>[^"]*)">', "{{{@\g<in> ", 중간_데이터)
            중간_데이터 = re.sub('<div style="(?P<in>[^"]*)">', "{{{#!wiki style=&quot;\g<in>&quot;\n", 중간_데이터)
            중간_데이터 = re.sub("(?P<in>.)", "<span>\g<in></span>", 중간_데이터)
            
            데이터 = 문법_컴파일.sub(중간_데이터, 데이터, 1)
        else:
            break
            
    데이터 = re.sub("<span>&</span><span>l</span><span>t</span><span>;</span>", "<span>&lt;</span>", 데이터)
    데이터 = re.sub("<span>&</span><span>g</span><span>t</span><span>;</span>", "<span>&gt;</span>", 데이터)
            
    return (데이터, 접기_숫자)

def 나무마크(title, data):
    data = HTML_파싱(data)

    접기_숫자 = 0
    임시_저장 = 중괄호_문법(data, 접기_숫자, False)
    
    data = 임시_저장[0]
    접기_숫자 = 임시_저장[1]
    
    data = re.sub("\[anchor\((?P<in>[^\[\]]*)\)\]", '<span id="\g<in>"></span>', data)
    data = re.sub('\[date\(now\)\]', 시간(), data)
    if(not re.search("\.", 아이피_확인(request))):
        name = '[[사용자:' + 아이피_확인(request) + '|' + 아이피_확인(request) + ']]'
    else:
        name = 아이피_확인(request)
    data = re.sub("\[name\]", name, data)
    
    while(True):
        m = re.search("\[include\(((?:(?!\)\]|,).)*)((?:,\s?(?:[^)]*))+)?\)\]", data)
        if(m):
            results = m.groups()
            if(results[0] == title):
                data = re.sub("\[include\(((?:(?!\)\]|,).)*)((?:,\s?(?:[^)]*))+)?\)\]", "<b>" + results[0] + "</b>", data, 1)
            else:
                DB_실행("select * from data where title = '" + DB_인코딩(results[0]) + "'")
                틀_내용 = DB_가져오기()
                
                if(틀_내용):
                    DB_실행("select * from back where title = '" + DB_인코딩(results[0]) + "' and link = '" + DB_인코딩(title) + "' and type = 'include'")
                    역링크 = DB_가져오기()
                    
                    if(not 역링크):
                        DB_실행("insert into back (title, link, type) value ('" + DB_인코딩(results[0]) + "', '" + DB_인코딩(title) + "',  'include')")
                        DB_갱신()
                        
                    틀_데이터 = 틀_내용[0]['data']
                    틀_데이터 = re.sub("\[include\(((?:(?!\)\]|,).)*)((?:,\s?(?:[^)]*))+)?\)\]", "", 틀_데이터)
                    
                    틀_데이터 = HTML_파싱(틀_데이터)
                    틀_데이터 = 중괄호_문법(틀_데이터, 접기_숫자, True)[0]
                    
                    if(results[1]):
                        a = results[1]
                        while(True):
                            g = re.search("([^= ,]*)\=([^,]*)", a)
                            if(g):
                                result = g.groups()
                                틀_데이터 = re.sub("@" + result[0] + "@", result[1], 틀_데이터)
                                a = re.sub("([^= ,]*)\=([^,]*)", "", a, 1)
                            else:
                                break                        
                    data = re.sub("\[include\(((?:(?!\)\]|,).)*)((?:,\s?(?:[^)]*))+)?\)\]", '\n<div>' + 틀_데이터 + '</div>\n', data, 1)
                else:
                    DB_실행("select * from back where title = '" + DB_인코딩(results[0]) + "' and link = '" + DB_인코딩(title) + "' and type = 'include'")
                    abb = DB_가져오기()
                    if(not abb):
                        DB_실행("insert into back (title, link, type) value ('" + DB_인코딩(results[0]) + "', '" + DB_인코딩(title) + "',  'include')")
                        DB_갱신()
                
                    data = re.sub("\[include\(((?:(?!\)\]|,).)*)((?:,\s?(?:[^)]*))+)?\)\]", "<a class=\"not_thing\" href=\"" + URL_인코딩(results[0]) + "\">" + results[0] + "</a>", data, 1)
        else:
            break
    
    while(True):
        m = re.search('^#(?:redirect|넘겨주기)\s([^\n]*)', data)
        if(m):
            results = m.groups()
            aa = re.search("^(.*)(#(?:.*))$", results[0])
            if(aa):
                results = aa.groups()
                data = re.sub('^#(?:redirect|넘겨주기)\s([^\n]*)', '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(results[0]) + '/from/' + URL_인코딩(title) + results[1] + '" />', data, 1)
            else:
                data = re.sub('^#(?:redirect|넘겨주기)\s([^\n]*)', '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(results[0]) + '/from/' + URL_인코딩(title) + '" />', data, 1)
            DB_실행("select * from back where title = '" + DB_인코딩(results[0]) + "' and link = '" + DB_인코딩(title) + "' and type = 'redirect'")
            abb = DB_가져오기()
            if(not abb):
                DB_실행("insert into back (title, link, type) value ('" + DB_인코딩(results[0]) + "', '" + DB_인코딩(title) + "',  'redirect')")
                DB_갱신()
        else:
            break
    
    
    
    data = '\n' + data + '\n'
    
    while(True):
        m = re.search("\n&gt;\s?((?:[^\n]*)(?:(?:(?:(?:\n&gt;\s?)(?:[^\n]*))+)?))", data)
        if(m):
            result = m.groups()
            blockquote = result[0]
            blockquote = re.sub("\n&gt;\s?", "\n", blockquote)
            data = re.sub("\n&gt;\s?((?:[^\n]*)(?:(?:(?:(?:\n&gt;\s?)(?:[^\n]*))+)?))", "\n<blockquote>" + blockquote + "</blockquote>", data, 1)
        else:
            break
    
    m = re.search('\[목차\]', data)
    if(not m):
        data = re.sub("(?P<in>(={1,6})\s?([^=]*)\s?(?:={1,6})(?:\s+)?\n)", "[목차]\n\g<in>", data, 1)
    
    i = 0
    h0c = 0
    h1c = 0
    h2c = 0
    h3c = 0
    h4c = 0
    h5c = 0
    last = 0
    rtoc = '<div id="toc"><span id="toc-name">목차</span><br><br>'
    while(True):
        i = i + 1
        m = re.search('(={1,6})\s?([^=]*)\s?(?:={1,6})(?:\s+)?\n', data)
        if(m):
            result = m.groups()
            wiki = len(result[0])
            if(last < wiki):
                last = wiki
            else:
                last = wiki
                if(wiki == 1):
                    h1c = 0
                    h2c = 0
                    h3c = 0
                    h4c = 0
                    h5c = 0
                elif(wiki == 2):
                    h2c = 0
                    h3c = 0
                    h4c = 0
                    h5c = 0
                elif(wiki == 3):
                    h3c = 0
                    h4c = 0
                    h5c = 0
                elif(wiki == 4):
                    h4c = 0
                    h5c = 0
                elif(wiki == 5):
                    h5c = 0
            if(wiki == 1):
                h0c = h0c + 1
            elif(wiki == 2):
                h1c = h1c + 1
            elif(wiki == 3):
                h2c = h2c + 1
            elif(wiki == 4):
                h3c = h3c + 1
            elif(wiki == 5):
                h4c = h4c + 1
            else:
                h5c = h5c + 1
            toc = str(h0c) + '.' + str(h1c) + '.' + str(h2c) + '.' + str(h3c) + '.' + str(h4c) + '.' + str(h5c) + '.'
            toc = re.sub("(?P<in>[0-9]0(?:[0]*)?)\.", '\g<in>#.', toc)
            toc = re.sub("0\.", '', toc)
            toc = re.sub("#\.", '.', toc)
            toc = re.sub("\.$", '', toc)
            rtoc = rtoc + '<a href="#s-' + toc + '">' + toc + '</a>. ' + result[1] + '<br>'
            c = re.sub(" $", "", result[1])
            d = c
            c = re.sub("\[\[(([^|]*)\|)?(?P<in>[^\]]*)\]\]", "\g<in>", c)
            data = re.sub('(={1,6})\s?([^=]*)\s?(?:={1,6})(?:\s+)?\n', '<h' + str(wiki) + ' id="' + c + '"><a href="#toc" id="s-' + toc + '">' + toc + '.</a> ' + d + ' <span style="font-size:11px;">[<a href="/edit/' + URL_인코딩(title) + '/section/' + str(i) + '">편집</a>]</span></h' + str(wiki) + '>', data, 1);
        else:
            rtoc = rtoc + '</div>'
            break
    
    data = re.sub("\[목차\]", rtoc, data)
    
    category = ''
    while(True):
        m = re.search("\[\[(분류:(?:(?:(?!\]\]).)*))\]\]", data)
        if(m):
            g = m.groups()
            
            if(not title == g[0]):
                DB_실행("select * from cat where title = '" + DB_인코딩(g[0]) + "' and cat = '" + DB_인코딩(title) + "'")
                abb = DB_가져오기()
                if(not abb):
                    DB_실행("insert into cat (title, cat) value ('" + DB_인코딩(g[0]) + "', '" + DB_인코딩(title) + "')")
                    DB_갱신()                
                    
                if(category == ''):
                    DB_실행("select * from data where title = '" + DB_인코딩(g[0]) + "'")
                    exists = DB_가져오기()
                    if(exists):
                        red = ""
                    else:
                        red = 'class="not_thing"'
                        
                    category = category + '<a ' + red + ' href="/w/' + URL_인코딩(g[0]) + '">' + re.sub("분류:", "", g[0]) + '</a>'
                else:
                    DB_실행("select * from data where title = '" + DB_인코딩(g[0]) + "'")
                    exists = DB_가져오기()
                    if(exists):
                        red = ""
                    else:
                        red = 'class="not_thing"'
                        
                    category = category + ' / ' + '<a ' + red + ' href="/w/' + URL_인코딩(g[0]) + '">' + re.sub("분류:", "", g[0]) + '</a>'
            
            data = re.sub("\[\[(분류:(?:(?:(?!\]\]).)*))\]\]", '', data, 1)
        else:
            break

    data = re.sub("'''(?P<in>.+?)'''(?!')", '<b>\g<in></b>', data)
    data = re.sub("''(?P<in>.+?)''(?!')", '<i>\g<in></i>', data)
    data = re.sub('~~(?P<in>.+?)~~(?!~)', '<s>\g<in></s>', data)
    data = re.sub('--(?P<in>.+?)--(?!-)', '<s>\g<in></s>', data)
    data = re.sub('__(?P<in>.+?)__(?!_)', '<u>\g<in></u>', data)
    data = re.sub('\^\^(?P<in>.+?)\^\^(?!\^)', '<sup>\g<in></sup>', data)
    data = re.sub(',,(?P<in>.+?),,(?!,)', '<sub>\g<in></sub>', data)
    
    data = re.sub('&lt;math&gt;(?P<in>((?!&lt;math&gt;).)*)&lt;\/math&gt;', '$\g<in>$', data)
    
    data = re.sub('{{\|(?P<in>(?:(?:(?:(?!\|}}).)*)(?:\n?))+)\|}}', '<table><tbody><tr><td>\g<in></td></tr></tbody></table>', data)
    
    data = re.sub('\[ruby\((?P<in>[^\|]*)\|(?P<out>[^\)]*)\)\]', '<ruby>\g<in><rp>(</rp><rt>\g<out></rt><rp>)</rp></ruby>', data)
    
    data = re.sub("##\s?(?P<in>[^\n]*)\n", "<div style='display:none;'>\g<in></div>", data);
    
    while(True):
        m = re.search("\[\[파일:((?:(?!\]\]|\|).)*)(?:\|((?:(?!\]\]).)*))?\]\]", data)
        if(m):
            c = m.groups()
            if(c[1]):
                n = re.search("width=([^ \n&]*)", c[1])
                e = re.search("height=([^ \n&]*)", c[1])
                if(n):
                    a = n.groups()
                    width = a[0]
                else:
                    width = ''
                if(e):
                    b = e.groups()
                    height = b[0]
                else:
                    height = ''
                img = re.sub("\.(?P<in>[Jj][Pp][Gg]|[Pp][Nn][Gg]|[Gg][Ii][Ff]|[Jj][Pp][Ee][Gg])", "#\g<in>#", c[0])
                data = re.sub("\[\[파일:((?:(?!\]\]|\?).)*)(?:\?((?:(?!\]\]).)*))?\]\]", '<a href="/w/파일:' + img + '"><img src="/image/' + img + '" width="' + width + '" height="' + height + '"></a>', data, 1)
            else:
                img = re.sub("\.(?P<in>[Jj][Pp][Gg]|[Pp][Nn][Gg]|[Gg][Ii][Ff]|[Jj][Pp][Ee][Gg])", "#\g<in>#", c[0])
                data = re.sub("\[\[파일:((?:(?!\]\]|\?).)*)(?:\?((?:(?!\]\]).)*))?\]\]", "<a href='/w/파일:" + img + "'><img src='/image/" + img + "'></a>", data, 1)
            if(not re.search("^파일:([^\n]*)", title)):
                DB_실행("select * from back where title = '" + DB_인코딩(c[0]) + "' and link = '" + DB_인코딩(title) + "' and type = 'redirect'")
                abb = DB_가져오기()
                if(not abb):
                    DB_실행("insert into back (title, link, type) value ('파일:" + DB_인코딩(c[0]) + "', '" + DB_인코딩(title) + "',  'file')")
                    DB_갱신()            
        else:
            break
    
    data = re.sub("\[br\]",'<br>', data);
    
    while(True):
        m = re.search("\[youtube\(((?:(?!,|\)\]).)*)(?:,\s)?(?:width=((?:(?!,|\)\]).)*))?(?:,\s)?(?:height=((?:(?!,|\)\]).)*))?(?:,\s)?(?:width=((?:(?!,|\)\]).)*))?\)\]", data)
        if(m):
            result = m.groups()
            if(result[1]):
                if(result[2]):
                    width = result[1]
                    height = result[2]
                else:
                    width = result[1]
                    height = '315'
            elif(result[2]):
                if(result[3]):
                    height = result[2]
                    width = result[3]
                else:
                    height = result[2]
                    width = '560'
            else:
                width = '560'
                height = '315'
            data = re.sub("\[youtube\(((?:(?!,|\)\]).)*)(?:,\s)?(?:width=((?:(?!,|\)\]).)*))?(?:,\s)?(?:height=((?:(?!,|\)\]).)*))?(?:,\s)?(?:width=((?:(?!,|\)\]).)*))?\)\]", '<iframe width="' + width + '" height="' + height + '" src="https://www.youtube.com/embed/' + result[0] + '" frameborder="0" allowfullscreen></iframe>', data, 1)
        else:
            break
     
    data = re.sub("\[\[(?::(?P<in>(?:분류|파일):(?:(?:(?!\]\]).)*)))\]\]", "[[\g<in>]]", data)
            
    while(True):
        m = re.search("\[\[(((?!\]\]).)*)\]\]", data)
        if(m):
            result = m.groups()
            a = re.search("((?:(?!\|).)*)\|(.*)", result[0])
            if(a):
                results = a.groups()
                aa = re.search("^(.*)(#(?:.*))$", results[0])
                if(aa):
                    g = results[1]
                    results = aa.groups()
                    b = re.search("^http(?:s)?:\/\/", results[0])
                    if(b):
                        data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a class="out_link" href="' + results[0] + results[1] + '">' + g + '</a>', data, 1)
                    else:
                        if(results[0] == title):
                            data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<b>' + g + '</b>', data, 1)
                        else:
                            DB_실행("select * from data where title = '" + DB_인코딩(results[0]) + "'")
                            rows = DB_가져오기()
                            if(rows):
                                data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a title="' + results[0] + results[1] + '" href="/w/' + URL_인코딩(results[0]) + results[1] + '">' + g + '</a>', data, 1)
                            else:
                                data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a title="' + results[0] + results[1] + '" class="not_thing" href="/w/' + URL_인코딩(results[0]) + results[1] + '">' + g + '</a>', data, 1)
                            DB_실행("select * from back where title = '" + DB_인코딩(results[0]) + "' and link = '" + DB_인코딩(title) + "'")
                            rows = DB_가져오기()
                            if(not rows):
                                DB_실행("insert into back (title, link, type) value ('" + DB_인코딩(results[0]) + "', '" + DB_인코딩(title) + "', '')")
                                DB_갱신()
                else:
                    b = re.search("^http(?:s)?:\/\/", results[0])
                    if(b):
                        c = re.search("(?:\.[Jj][Pp][Gg]|\.[Pp][Nn][Gg]|\.[Gg][Ii][Ff]|\.[Jj][Pp][Ee][Gg])", results[0])
                        if(c):
                            img = results[0]
                            img = re.sub("\.(?P<in>[Jj][Pp][Gg]|[Pp][Nn][Gg]|[Gg][Ii][Ff]|[Jj][Pp][Ee][Gg])", "#\g<in>#", img)
                            data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a class="out_link" href="' + img + '">' + results[1] + '</a>', data, 1)
                        else:
                            data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a class="out_link" href="' + results[0] + '">' + results[1] + '</a>', data, 1)
                    else:
                        if(results[0] == title):
                            data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<b>' + results[1] + '</b>', data, 1)
                        else:
                            DB_실행("select * from data where title = '" + DB_인코딩(results[0]) + "'")
                            rows = DB_가져오기()
                            if(rows):
                                data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a title="' + results[0] + '" href="/w/' + URL_인코딩(results[0]) + '">' + results[1] + '</a>', data, 1)
                            else:
                                data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a title="' + results[0] + '" class="not_thing" href="/w/' + URL_인코딩(results[0]) + '">' + results[1] + '</a>', data, 1)
                            DB_실행("select * from back where title = '" + DB_인코딩(results[0]) + "' and link = '" + DB_인코딩(title) + "'")
                            rows = DB_가져오기()
                            if(not rows):
                                DB_실행("insert into back (title, link, type) value ('" + DB_인코딩(results[0]) + "', '" + DB_인코딩(title) + "', '')")
                                DB_갱신()
            else:
                aa = re.search("^(.*)(#(?:.*))$", result[0])
                if(aa):
                    result = aa.groups()
                    b = re.search("^http(?:s)?:\/\/", result[0])
                    if(b):
                        data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a class="out_link" href="' + result[0] + result[1] + '">' + result[0] + result[1] + '</a>', data, 1)
                    else:
                        if(result[0] == title):
                            data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<b>' + result[0] + result[1] + '</b>', data, 1)
                        else:
                            DB_실행("select * from data where title = '" + DB_인코딩(result[0]) + "'")
                            rows = DB_가져오기()
                            if(rows):
                                data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a href="/w/' + URL_인코딩(result[0]) + result[1] + '">' + result[0] + result[1] + '</a>', data, 1)
                            else:
                                data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a class="not_thing" href="/w/' + URL_인코딩(result[0]) + result[1] + '">' + result[0] + result[1] + '</a>', data, 1)
                            DB_실행("select * from back where title = '" + DB_인코딩(result[0]) + "' and link = '" + DB_인코딩(title) + "'")
                            rows = DB_가져오기()
                            if(not rows):
                                DB_실행("insert into back (title, link, type) value ('" + DB_인코딩(result[0]) + "', '" + DB_인코딩(title) + "', '')")
                                DB_갱신()
                else:
                    b = re.search("^http(?:s)?:\/\/", result[0])
                    if(b):
                        c = re.search("(?:\.[Jj][Pp][Gg]|\.[Pp][Nn][Gg]|\.[Gg][Ii][Ff]|\.[Jj][Pp][Ee][Gg])", result[0])
                        if(c):
                            img = result[0]
                            img = re.sub("\.(?P<in>[Jj][Pp][Gg]|[Pp][Nn][Gg]|[Gg][Ii][Ff]|[Jj][Pp][Ee][Gg])", "#\g<in>#", img)
                            data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a class="out_link" href="' + img + '">' + img + '</a>', data, 1)
                        else:
                            data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a class="out_link" href="' + result[0] + '">' + result[0] + '</a>', data, 1)
                    else:
                        if(result[0] == title):
                            data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<b>' + result[0] + '</b>', data, 1)
                        else:
                            DB_실행("select * from data where title = '" + DB_인코딩(result[0]) + "'")
                            rows = DB_가져오기()
                            if(rows):
                                data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a href="/w/' + URL_인코딩(result[0]) + '">' + result[0] + '</a>', data, 1)
                            else:
                                data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a class="not_thing" href="/w/' + URL_인코딩(result[0]) + '">' + result[0] + '</a>', data, 1)
                            DB_실행("select * from back where title = '" + DB_인코딩(result[0]) + "' and link = '" + DB_인코딩(title) + "'")
                            rows = DB_가져오기()
                            if(not rows):
                                DB_실행("insert into back (title, link, type) value ('" + DB_인코딩(result[0]) + "', '" + DB_인코딩(title) + "', '')")
                                DB_갱신()
        else:
            break
            
    while(True):
        m = re.search("(http(?:s)?:\/\/(?:(?:(?:(?!\.[Jj][Pp][Gg]|\.[Pp][Nn][Gg]|\.[Gg][Ii][Ff]|\.[Jj][Pp][Ee][Gg]|#[Jj][Pp][Gg]#|#[Pp][Nn][Gg]#|#[Gg][Ii][Ff]#|#[Jj][Pp][Ee][Gg]#|<\/(?:[^>]*)>).)*)(?:\.[Jj][Pp][Gg]|\.[Pp][Nn][Gg]|\.[Gg][Ii][Ff]|\.[Jj][Pp][Ee][Gg])))(?:(?:(?:\?)width=((?:[0-9]*)(?:px|%)?))?(?:(?:\?|&)height=((?:[0-9]*)(?:px|%)?))?(?:(?:&)width=((?:[0-9]*)(?:px|%)?))?)?", data)
        if(m):
            result = m.groups()
            if(result[1]):
                if(result[2]):
                    width = result[1]
                    height = result[2]
                else:
                    width = result[1]
                    height = ''
            elif(result[2]):
                if(result[3]):
                    height = result[2]
                    width = result[3]
                else:
                    height = result[2]
                    width = ''
            else:
                width = ''
                height = ''
            c = result[0]
            c = re.sub("\.(?P<in>[Jj][Pp][Gg]|[Pp][Nn][Gg]|[Gg][Ii][Ff]|[Jj][Pp][Ee][Gg])", "#\g<in>#", c)
            data = re.sub("(http(?:s)?:\/\/(?:(?:(?:(?!\.[Jj][Pp][Gg]|\.[Pp][Nn][Gg]|\.[Gg][Ii][Ff]|\.[Jj][Pp][Ee][Gg]|#[Jj][Pp][Gg]#|#[Pp][Nn][Gg]#|#[Gg][Ii][Ff]#|#[Jj][Pp][Ee][Gg]#|<\/(?:[^>]*)>).)*)(?:\.[Jj][Pp][Gg]|\.[Pp][Nn][Gg]|\.[Gg][Ii][Ff]|\.[Jj][Pp][Ee][Gg])))(?:(?:(?:\?)width=((?:[0-9]*)(?:px|%)?))?(?:(?:\?|&)height=((?:[0-9]*)(?:px|%)?))?(?:(?:&)width=((?:[0-9]*)(?:px|%)?))?)?", "<img width='" + width + "' height='" + height + "' src='" + c + "'>", data, 1)
        else:
            break
            
    while(True):
        m = re.search("((?:(?:( +)\*\s(?:[^\n]*))\n?)+)", data)
        if(m):
            result = m.groups()
            end = str(result[0])
            while(True):
                isspace = re.search("( +)\*\s([^\n]*)", end)
                if(isspace):
                    spacebar = isspace.groups()
                    up = len(spacebar[0]) * 20
                    end = re.sub("( +)\*\s([^\n]*)", "<li style='margin-left:" + str(up) + "px'>" + spacebar[1] + "</li>", end, 1)
                else:
                    break
            end = re.sub("\n", '', end)
            data = re.sub("(?:(?:(?:( +)\*\s([^\n]*))\n?)+)", '<ul id="list">' + end + '</ul>', data, 1)
        else:
            break
    
    data = re.sub('\[date\]', 시간(), data)
    
    data = re.sub("#(?P<in>[Jj][Pp][Gg]|[Pp][Nn][Gg]|[Gg][Ii][Ff]|[Jj][Pp][Ee][Gg])#", ".\g<in>", data)
    
    data = re.sub("-{4,11}", "<hr>", data)
    
    while(True):
        b = re.search("\n( +)", data)
        if(b):
            result = b.groups()
            up = re.sub(' ', '<span id="in"></span>', result[0])
            data = re.sub("\n( +)", '<br>' + up, data, 1)
        else:
            break
    
    a = 1
    tou = "<hr id='footnote'><div class='wiki-macro-footnote'><br>"
    while(True):
        b = re.search("\[\*([^\s]*)\s(((?!\]).)*)\]", data)
        if(b):
            results = b.groups()
            if(results[0]):
                c = results[1]
                c = re.sub("<(?:[^>]*)>", '', c)
                tou = tou + "<span class='footnote-list'><a href=\"#rfn-" + str(a) + "\" id=\"fn-" + str(a) + "\">[" + results[0] + "]</a> " + results[1] + "</span><br>"
                data = re.sub("\[\*([^\s]*)\s(((?!\]).)*)\]", "<sup><a class=\"footnotes\" title=\"" + c + "\" id=\"rfn-" + str(a) + "\" href=\"#fn-" + str(a) + "\">[" + results[0] + "]</a></sup>", data, 1)
            else:
                c = results[1]
                c = re.sub("<(?:[^>]*)>", '', c)
                tou = tou + "<span class='footnote-list'><a href=\"#rfn-" + str(a) + "\" id=\"fn-" + str(a) + "\">[" + str(a) + "]</a> " + results[1] + "</span><br>"
                data = re.sub("\[\*([^\s]*)\s(((?!\]).)*)\]", '<sup><a class="footnotes" title="' + c + '" id="rfn-' + str(a) + '" href="#fn-' + str(a) + '">[' + str(a) + ']</a></sup>', data, 1)
            a = a + 1
        else:
            tou = tou + '</div>'
            if(tou == "<hr id='footnote'><div class='wiki-macro-footnote'><br></div>"):
                tou = ""
            break
    
    data = re.sub("\[각주\](?:(?:<br>| |\r|\n)+)?$", "", data)
    data = re.sub("(?:(?:<br>| |\r|\n)+)$", "", data)
    data = re.sub("\[각주\]", "<br>" + tou, data)
    data = data + tou
    
    if(category):
        data = data + '<div style="width:100%;border: 1px solid #777;padding: 5px;margin-top: 1em;">분류: ' + category + '</div>'
        
    data = re.sub("(?:\|\|\r\n)", "#table#<nobr>", data)
        
    while(True):
        있나 = re.search("(\|\|(?:(?:(?:(?:(?!\|\|).)*)(?:\n?))+))", data)
        if(있나):
            분리 = 있나.groups()
            
            중간_내용 = re.sub("\|\|", "#table#", 분리[0])
            중간_내용 = re.sub("\r\n", "<br>", 중간_내용)
            
            data = re.sub("(\|\|((?:(?:(?:(?!\|\|).)*)(?:\n?))+))", 중간_내용, data, 1)
        else:
            break
            
    data = re.sub("#table#", "||", data)
    data = re.sub("<nobr>", "\r\n", data)
    
    while(True):
        m = re.search("(\|\|(?:(?:(?:.*)\n?)\|\|)+)", data)
        if(m):
            results = m.groups()
            table = results[0]
            while(True):
                a = re.search("^(\|\|(?:(?:\|\|)+)?)((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)?", table)
                if(a):
                    row = ''
                    cel = ''
                    celstyle = ''
                    rowstyle = ''
                    alltable = ''
                    result = a.groups()
                    if(result[1]):
                        q = re.search("&lt;table\s?width=((?:(?!&gt;).)*)&gt;", result[1])
                        w = re.search("&lt;table\s?height=((?:(?!&gt;).)*)&gt;", result[1])
                        e = re.search("&lt;table\s?align=((?:(?!&gt;).)*)&gt;", result[1])
                        alltable = 'style="'
                        celstyle = 'style="'
                        rowstyle = 'style="'
                        if(q):
                            resultss = q.groups()
                            alltable = alltable + 'width:' + resultss[0] + ';'
                        if(w):
                            resultss = w.groups()
                            alltable = alltable + 'height:' + resultss[0] + ';'
                        if(e):
                            resultss = e.groups()
                            if(resultss[0] == 'right'):
                                alltable = alltable + 'margin-left:auto;'
                            elif(resultss[0] == 'center'):
                                alltable = alltable + 'margin:auto;'
                            else:
                                alltable = alltable + 'margin-right:auto;'
                                
                        ee = re.search("&lt;table\s?textalign=((?:(?!&gt;).)*)&gt;", result[1])
                        if(ee):
                            resultss = ee.groups()
                            if(resultss[0] == 'right'):
                                alltable = alltable + 'text-align:right;'
                            elif(resultss[0] == 'center'):
                                alltable = alltable + 'text-align:center;'
                            else:
                                alltable = alltable + 'text-align:left;'
                        
                        r = re.search("&lt;-((?:(?!&gt;).)*)&gt;", result[1])
                        if(r):
                            resultss = r.groups()
                            cel = 'colspan="' + resultss[0] + '"'
                        else:
                            cel = 'colspan="' + str(round(len(result[0]) / 2)) + '"'    
                        t = re.search("&lt;\|((?:(?!&gt;).)*)&gt;", result[1])
                        if(t):
                            resultss = t.groups()
                            row = 'rowspan="' + resultss[0] + '"'

                        ba = re.search("&lt;rowbgcolor=(#[0-9a-f-A-F]{6})&gt;", result[1])
                        bb = re.search("&lt;rowbgcolor=(#[0-9a-f-A-F]{3})&gt;", result[1])
                        bc = re.search("&lt;rowbgcolor=(\w+)&gt;", result[1])
                        if(ba):
                            resultss = ba.groups()
                            rowstyle = rowstyle + 'background:' + resultss[0] + ';'
                        elif(bb):
                            resultss = bb.groups()
                            rowstyle = rowstyle + 'background:' + resultss[0] + ';'
                        elif(bc):
                            resultss = bc.groups()
                            rowstyle = rowstyle + 'background:' + resultss[0] + ';'
                            
                        z = re.search("&lt;table\s?bordercolor=(#[0-9a-f-A-F]{6})&gt;", result[1])
                        x = re.search("&lt;table\s?bordercolor=(#[0-9a-f-A-F]{3})&gt;", result[1])
                        c = re.search("&lt;table\s?bordercolor=(\w+)&gt;", result[1])
                        if(z):
                            resultss = z.groups()
                            alltable = alltable + 'border:' + resultss[0] + ' 2px solid;'
                        elif(x):
                            resultss = x.groups()
                            alltable = alltable + 'border:' + resultss[0] + ' 2px solid;'
                        elif(c):
                            resultss = c.groups()
                            alltable = alltable + 'border:' + resultss[0] + ' 2px solid;'
                            
                        aq = re.search("&lt;table\s?bgcolor=(#[0-9a-f-A-F]{6})&gt;", result[1])
                        aw = re.search("&lt;table\s?bgcolor=(#[0-9a-f-A-F]{3})&gt;", result[1])
                        ae = re.search("&lt;table\s?bgcolor=(\w+)&gt;", result[1])
                        if(aq):
                            resultss = aq.groups()
                            alltable = alltable + 'background:' + resultss[0] + ';'
                        elif(aw):
                            resultss = aw.groups()
                            alltable = alltable + 'background:' + resultss[0] + ';'
                        elif(ae):
                            resultss = ae.groups()
                            alltable = alltable + 'background:' + resultss[0] + ';'
                            
                        j = re.search("&lt;bgcolor=(#[0-9a-f-A-F]{6})&gt;", result[1])
                        k = re.search("&lt;bgcolor=(#[0-9a-f-A-F]{3})&gt;", result[1])
                        l = re.search("&lt;bgcolor=(\w+)&gt;", result[1])
                        if(j):
                            resultss = j.groups()
                            celstyle = celstyle + 'background:' + resultss[0] + ';'
                        elif(k):
                            resultss = k.groups()
                            celstyle = celstyle + 'background:' + resultss[0] + ';'
                        elif(l):
                            resultss = l.groups()
                            celstyle = celstyle + 'background:' + resultss[0] + ';'
                            
                        aa = re.search("&lt;(#[0-9a-f-A-F]{6})&gt;", result[1])
                        ab = re.search("&lt;(#[0-9a-f-A-F]{3})&gt;", result[1])
                        ac = re.search("&lt;(\w+)&gt;", result[1])
                        if(aa):
                            resultss = aa.groups()
                            celstyle = celstyle + 'background:' + resultss[0] + ';'
                        elif(ab):
                            resultss = ab.groups()
                            celstyle = celstyle + 'background:' + resultss[0] + ';'
                        elif(ac):
                            resultss = ac.groups()
                            celstyle = celstyle + 'background:' + resultss[0] + ';'
                            
                        qa = re.search("&lt;width=((?:(?!&gt;).)*)&gt;", result[1])
                        qb = re.search("&lt;height=((?:(?!&gt;).)*)&gt;", result[1])
                        if(qa):
                            resultss = qa.groups()
                            celstyle = celstyle + 'width:' + resultss[0] + ';'
                        if(qb):
                            resultss = qb.groups()
                            celstyle = celstyle + 'height:' + resultss[0] + ';'
                            
                        i = re.search("&lt;\)&gt;", result[1])
                        o = re.search("&lt;:&gt;", result[1])
                        p = re.search("&lt;\(&gt;", result[1])
                        if(i):
                            celstyle = celstyle + 'text-align:right;'
                        elif(o):
                            celstyle = celstyle + 'text-align:center;'
                        elif(p):
                            celstyle = celstyle + 'text-align:left;'
                            
                        alltable = alltable + '"'
                        celstyle = celstyle + '"'
                        rowstyle = rowstyle + '"'
                            
                        table = re.sub("^(\|\|(?:(?:\|\|)+)?)((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)?", "<table " + alltable + "><tbody><tr " + rowstyle + "><td " + cel + " " + row + " " + celstyle + ">", table, 1)
                    else:
                        cel = 'colspan="' + str(round(len(result[0]) / 2)) + '"'
                        table = re.sub("^(\|\|(?:(?:\|\|)+)?)((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)?", "<table><tbody><tr><td " + cel + ">", table, 1)
                else:
                    break
                    
            table = re.sub("\|\|$", "</td></tr></tbody></table>", table)
            
            while(True):
                b = re.search("\|\|\r\n(\|\|(?:(?:\|\|)+)?)((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)?", table)
                if(b):
                    row = ''
                    cel = ''
                    celstyle = ''
                    rowstyle = ''
                    result = b.groups()
                    if(result[1]):
                        celstyle = 'style="'
                        rowstyle = 'style="'
                        
                        r = re.search("&lt;-((?:(?!&gt;).)*)&gt;", result[1])
                        if(r):
                            resultss = r.groups()
                            cel = 'colspan="' + resultss[0] + '"'
                        else:
                            cel = 'colspan="' + str(round(len(result[0]) / 2)) + '"'
                        t = re.search("&lt;\|((?:(?!&gt;).)*)&gt;", result[1])
                        if(t):
                            resultss = t.groups()
                            row = 'rowspan="' + resultss[0] + '"'
                            
                        ba = re.search("&lt;rowbgcolor=(#[0-9a-f-A-F]{6})&gt;", result[1])
                        bb = re.search("&lt;rowbgcolor=(#[0-9a-f-A-F]{3})&gt;", result[1])
                        bc = re.search("&lt;rowbgcolor=(\w+)&gt;", result[1])
                        if(ba):
                            resultss = ba.groups()
                            rowstyle = rowstyle + 'background:' + resultss[0] + ';'
                        elif(bb):
                            resultss = bb.groups()
                            rowstyle = rowstyle + 'background:' + resultss[0] + ';'
                        elif(bc):
                            resultss = bc.groups()
                            rowstyle = rowstyle + 'background:' + resultss[0] + ';'
                            
                        j = re.search("&lt;bgcolor=(#[0-9a-f-A-F]{6})&gt;", result[1])
                        k = re.search("&lt;bgcolor=(#[0-9a-f-A-F]{3})&gt;", result[1])
                        l = re.search("&lt;bgcolor=(\w+)&gt;", result[1])
                        if(j):
                            resultss = j.groups()
                            celstyle = celstyle + 'background:' + resultss[0] + ';'
                        elif(k):
                            resultss = k.groups()
                            celstyle = celstyle + 'background:' + resultss[0] + ';'
                        elif(l):
                            resultss = l.groups()
                            celstyle = celstyle + 'background:' + resultss[0] + ';'

                        aa = re.search("&lt;(#[0-9a-f-A-F]{6})&gt;", result[1])
                        ab = re.search("&lt;(#[0-9a-f-A-F]{3})&gt;", result[1])
                        ac = re.search("&lt;(\w+)&gt;", result[1])
                        if(aa):
                            resultss = aa.groups()
                            celstyle = celstyle + 'background:' + resultss[0] + ';'
                        elif(ab):
                            resultss = ab.groups()
                            celstyle = celstyle + 'background:' + resultss[0] + ';'
                        elif(ac):
                            resultss = ac.groups()
                            celstyle = celstyle + 'background:' + resultss[0] + ';'    

                        qa = re.search("&lt;width=((?:(?!&gt;).)*)&gt;", result[1])
                        qb = re.search("&lt;height=((?:(?!&gt;).)*)&gt;", result[1])
                        if(qa):
                            resultss = qa.groups()
                            celstyle = celstyle + 'width:' + resultss[0] + ';'
                        if(qb):
                            resultss = qb.groups()
                            celstyle = celstyle + 'height:' + resultss[0] + ';'
                            
                        i = re.search("&lt;\)&gt;", result[1])
                        o = re.search("&lt;:&gt;", result[1])
                        p = re.search("&lt;\(&gt;", result[1])
                        if(i):
                            celstyle = celstyle + 'text-align:right;'
                        elif(o):
                            celstyle = celstyle + 'text-align:center;'
                        elif(p):
                            celstyle = celstyle + 'text-align:left;'

                        celstyle = celstyle + '"'
                        rowstyle = rowstyle + '"'
                        
                        table = re.sub("\|\|\r\n(\|\|(?:(?:\|\|)+)?)((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)?", "</td></tr><tr " + rowstyle + "><td " + cel + " " + row + " " + celstyle + ">", table, 1)
                    else:
                        cel = 'colspan="' + str(round(len(result[0]) / 2)) + '"'
                        table = re.sub("\|\|\r\n(\|\|(?:(?:\|\|)+)?)((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)?", "</td></tr><tr><td " + cel + ">", table, 1)
                else:
                    break

            while(True):
                c = re.search("(\|\|(?:(?:\|\|)+)?)((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)?", table)
                if(c):
                    row = ''
                    cel = ''
                    celstyle = ''
                    result = c.groups()
                    if(result[1]):
                        celstyle = 'style="'
                        
                        r = re.search("&lt;-((?:(?!&gt;).)*)&gt;", result[1])
                        if(r):
                            resultss = r.groups()
                            cel = 'colspan="' + resultss[0] + '"';
                        else:
                            cel = 'colspan="' + str(round(len(result[0]) / 2)) + '"'
                        t = re.search("&lt;\|((?:(?!&gt;).)*)&gt;", result[1])
                        if(t):
                            resultss = t.groups()
                            row = 'rowspan="' + resultss[0] + '"';
                            
                        j = re.search("&lt;bgcolor=(#[0-9a-f-A-F]{6})&gt;", result[1])
                        k = re.search("&lt;bgcolor=(#[0-9a-f-A-F]{3})&gt;", result[1])
                        l = re.search("&lt;bgcolor=(\w+)&gt;", result[1])
                        if(j):
                            resultss = j.groups()
                            celstyle = celstyle + 'background:' + resultss[0] + ';'
                        elif(k):
                            resultss = k.groups()
                            celstyle = celstyle + 'background:' + resultss[0] + ';'
                        elif(l):
                            resultss = l.groups()
                            celstyle = celstyle + 'background:' + resultss[0] + ';'
                            
                        aa = re.search("&lt;(#[0-9a-f-A-F]{6})&gt;", result[1])
                        ab = re.search("&lt;(#[0-9a-f-A-F]{3})&gt;", result[1])
                        ac = re.search("&lt;(\w+)&gt;", result[1])
                        if(aa):
                            resultss = aa.groups()
                            celstyle = celstyle + 'background:' + resultss[0] + ';'
                        elif(ab):
                            resultss = ab.groups()
                            celstyle = celstyle + 'background:' + resultss[0] + ';'
                        elif(ac):
                            resultss = ac.groups()
                            celstyle = celstyle + 'background:' + resultss[0] + ';'
                            
                        qa = re.search("&lt;width=((?:(?!&gt;).)*)&gt;", result[1])
                        qb = re.search("&lt;height=((?:(?!&gt;).)*)&gt;", result[1])
                        if(qa):
                            resultss = qa.groups()
                            celstyle = celstyle + 'width:' + resultss[0] + ';'
                        if(qb):
                            resultss = qb.groups()
                            celstyle = celstyle + 'height:' + resultss[0] + ';'
                            
                        i = re.search("&lt;\)&gt;", result[1])
                        o = re.search("&lt;:&gt;", result[1])
                        p = re.search("&lt;\(&gt;", result[1])
                        if(i):
                            celstyle = celstyle + 'text-align:right;'
                        elif(o):
                            celstyle = celstyle + 'text-align:center;'
                        elif(p):
                            celstyle = celstyle + 'text-align:left;'

                        celstyle = celstyle + '"'
                            
                        table = re.sub("(\|\|(?:(?:\|\|)+)?)((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)?", "</td><td " + cel + " " + row + " " + celstyle + ">", table, 1)
                    else:
                        cel = 'colspan="' + str(round(len(result[0]) / 2)) + '"'
                        table = re.sub("(\|\|(?:(?:\|\|)+)?)((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)?", "</td><td " + cel + ">", table, 1)
                else:
                    break
            
            data = re.sub("(\|\|(?:(?:(?:.*)\n?)\|\|)+)", table, data, 1)
        else:
            break
    
    data = re.sub('<\/blockquote>((\r)?\n){2}<blockquote>', '</blockquote><br><blockquote>', data)
    data = re.sub('\n', '<br>', data)
    data = re.sub('^<br>', '', data)
    return str(data)

def 아이피_확인(request):
    if(session.get('Now') == True):
        ip = format(session['DREAMER'])
    else:
        if(request.headers.getlist("X-Forwarded-For")):
            ip = request.headers.getlist("X-Forwarded-For")[0]
        else:
            ip = request.remote_addr
    return ip

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
    
@app.route('/upload', methods=['GET', 'POST'])
def 업로드():
    app.config['MAX_CONTENT_LENGTH'] = int(data['upload']) * 1024 * 1024
    if(request.method == 'POST'):
        ip = 아이피_확인(request)
        ban = 차단_체크(ip)
        
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
                        
                        DB_실행("insert into data (title, data, acl) value ('" + DB_인코딩('파일:' + filename) + "', '" + DB_인코딩('[[파일:' + filename + ']][br][br]{{{[[파일:' + filename + ']]}}}') + "', '')")
                        DB_갱신()
                        
                        역사_추가('파일:' + filename, '[[파일:' + filename + ']][br][br]{{{[[파일:' + filename + ']]}}}', 시간(), ip, '파일:' + filename + ' 업로드', '0')
                        
                        return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩('파일:' + filename) + '" />'
                else:
                    return '<meta http-equiv="refresh" content="0;url=/error/15" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/14" />'
    else:
        ip = 아이피_확인(request)
        ban = 차단_체크(ip)
        
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            return 웹_디자인('index.html', logo = data['name'], title = '업로드', tn = 21, number = data['upload'])
    
@app.route('/image/<path:name>')
def 이미지(name = None):
    if(os.path.exists(os.path.join('image', name))):
        return send_file(os.path.join('image', name), mimetype='image')
    else:
        return 웹_디자인('index.html', logo = data['name'], data = '이미지 없음.', title = '이미지 보기'), 404

@app.route('/adminlist')
def 관리자_목록():
    i = 0
    div = '<div>'
    
    DB_실행("select * from user where acl = 'admin' or acl = 'owner'")
    rows = DB_가져오기()
    if(rows):
        while(True):
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                break

            if(rows[i]['acl'] == 'owner'):
                acl = '소유자'
            else:
                acl = '관리자'

            DB_실행("select * from data where title = '사용자:" + rows[i]['id'] + "'")
            user = DB_가져오기()
            if(user):
                name = '<a href="/w/' + URL_인코딩('사용자:' + rows[i]['id']) + '">' + rows[i]['id'] + '</a> (' + acl + ')'
            else:
                name = '<a class="not_thing" href="/w/' + URL_인코딩('사용자:' + rows[i]['id']) + '">' + rows[i]['id'] + '</a> (' + acl + ')'

            div = div + '<li>' + str(i + 1) + '. ' + name + '</li>'
            
            i = i + 1
            
        return 웹_디자인('index.html', logo = data['name'], data = div, title = '관리자 목록')
    else:
        return 웹_디자인('index.html', logo = data['name'], title = '관리자 목록')

@app.route('/recentchanges')
def 최근바뀜():
    i = 0
    div = '<div>'
    
    DB_실행("select * from history order by date desc limit 50")
    rows = DB_가져오기()
    if(rows):
        admin = 관리자_확인()
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
                
            if(admin == 1):
                DB_실행("select * from ban where block = '" + DB_인코딩(rows[i]['ip']) + "'")
                row = DB_가져오기()
                if(row):
                    ban = ' <a href="/ban/' + URL_인코딩(rows[i]['ip']) + '">(해제)</a>'
                else:
                    ban = ' <a href="/ban/' + URL_인코딩(rows[i]['ip']) + '">(차단)</a>'
            else:
                ban = ''
                
            if(re.search('\.', rows[i]['ip'])):
                ip = rows[i]['ip'] + ' <a href="/record/' + URL_인코딩(rows[i]['ip']) + '/n/1">(기록)</a>'
            else:
                DB_실행("select * from data where title = '사용자:" + DB_인코딩(rows[i]['ip']) + "'")
                row = DB_가져오기()
                if(row):
                    ip = '<a href="/w/' + URL_인코딩('사용자:' + rows[i]['ip']) + '">' + rows[i]['ip'] + '</a> <a href="/record/' + URL_인코딩(rows[i]['ip']) + '/n/1">(기록)</a>'
                else:
                    ip = '<a class="not_thing" href="/w/' + URL_인코딩('사용자:' + rows[i]['ip']) + '">' + rows[i]['ip'] + '</a> <a href="/record/' + URL_인코딩(rows[i]['ip']) + '/n/1">(기록)</a>'
                    
            if((int(rows[i]['id']) - 1) == 0):
                revert = ''
            else:
                revert = '<a href="/revert/' + URL_인코딩(rows[i]['title']) + '/r/' + str(int(rows[i]['id']) - 1) + '">(되돌리기)</a>'
                
            div = div + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:33.33%;"><a href="/w/' + URL_인코딩(rows[i]['title']) + '">' + title + '</a> <a href="/history/' + URL_인코딩(rows[i]['title']) + '/n/1">(역사)</a> ' + revert + ' (' + leng + ')</td><td style="text-align: center;width:33.33%;">' + ip + ban + '</td><td style="text-align: center;width:33.33%;">' + rows[i]['date'] + '</td></tr><tr><td colspan="3" style="text-align: center;width:100%;">' + send + '</td></tr></tbody></table>'
            
            i = i + 1
            
        return 웹_디자인('index.html', logo = data['name'], rows = div, tn = 3, title = '최근 변경내역')
    else:
        return 웹_디자인('index.html', logo = data['name'], rows = '', tn = 3, title = '최근 변경내역')
        
@app.route('/history/<path:name>/r/<int:num>/hidden')
def 역사_숨기기(name = None, num = None):
    if(소유자_확인() == 1):
        DB_실행("select * from hidhi where title = '" + DB_인코딩(name) + "' and re = '" + DB_인코딩(str(num)) + "'")
        rows = DB_가져오기()
        if(rows):
            DB_실행("delete from hidhi where title = '" + DB_인코딩(name) + "' and re = '" + DB_인코딩(str(num)) + "'")
        else:
            DB_실행("insert into hidhi (title, re) value ('" + DB_인코딩(name) + "', '" + DB_인코딩(str(num)) + "')")
        DB_갱신()
        return '<meta http-equiv="refresh" content="0;url=/history/' + URL_인코딩(name) + '/n/1" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/history/' + URL_인코딩(name) + '/n/1" />'
        
@app.route('/record/<path:name>/n/<int:number>')
def 사용자_기록(name = None, number = None):
    v = number * 50
    i = v - 50
    div = '<div>'
    
    DB_실행("select * from history where ip = '" + DB_인코딩(name) + "' order by date desc")
    rows = DB_가져오기()
    if(rows):
        admin = 관리자_확인()
        
        while(True):
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                if(number != 1):
                    div = div + '<br><a href="/record/' + URL_인코딩(name) + '/n/' + str(number - 1) + '">(이전)'
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
                
            if(admin == 1):
                DB_실행("select * from ban where block = '" + DB_인코딩(rows[i]['ip']) + "'")
                row = DB_가져오기()
                if(row):
                    ban = ' <a href="/ban/' + URL_인코딩(rows[i]['ip']) + '">(해제)</a>'
                else:
                    ban = ' <a href="/ban/' + URL_인코딩(rows[i]['ip']) + '">(차단)</a>'
            else:
                ban = ''
                
            if(re.search('\.', rows[i]['ip'])):
                ip = rows[i]['ip']
            else:
                DB_실행("select * from data where title = '사용자:" + DB_인코딩(rows[i]['ip']) + "'")
                row = DB_가져오기()
                if(row):
                    ip = '<a href="/w/' + URL_인코딩('사용자:' + rows[i]['ip']) + '">' + rows[i]['ip'] + '</a>'
                else:
                    ip = '<a class="not_thing" href="/w/' + URL_인코딩('사용자:' + rows[i]['ip']) + '">' + rows[i]['ip'] + '</a>'
                    
            if((int(rows[i]['id']) - 1) == 0):
                revert = ''
            else:
                revert = '<a href="/revert/' + URL_인코딩(rows[i]['title']) + '/r/' + str(int(rows[i]['id']) - 1) + '">(되돌리기)</a>'
                
            div = div + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:33.33%;"><a href="/w/' + URL_인코딩(rows[i]['title']) + '">' + title + '</a> r' + rows[i]['id'] + ' <a href="/history/' + URL_인코딩(rows[i]['title']) + '/n/1">(역사)</a> ' + revert + ' (' + leng + ')</td><td style="text-align: center;width:33.33%;">' + ip + ban +  '</td><td style="text-align: center;width:33.33%;">' + rows[i]['date'] + '</td></tr><tr><td colspan="3" style="text-align: center;width:100%;">' + send + '</td></tr></tbody></table>'
            
            if(i == v):
                div = div + '</div>'
                if(number == 1):
                    div = div + '<br><a href="/record/' + URL_인코딩(name) + '/n/' + str(number + 1) + '">(다음)'
                else:
                    div = div + '<br><a href="/record/' + URL_인코딩(name) + '/n/' + str(number - 1) + '">(이전) <a href="/record/' + URL_인코딩(name) + '/n/' + str(number + 1) + '">(다음)'
                break
            else:
                i = i + 1
                
        return 웹_디자인('index.html', logo = data['name'], rows = div, tn = 3, title = '유저 기록')
    else:
        return 웹_디자인('index.html', logo = data['name'], rows = '', tn = 3, title = '유저 기록')
      
@app.route('/userlog/n/<int:number>')
def 모든_사용자(number = None):
    숫자_1 = number * 50
    숫자_2 = 숫자_1 - 50
    목록 = ''
    
    DB_실행("select * from user")
    사용자_목록 = DB_가져오기()
    if(사용자_목록):
        관리자 = 관리자_확인()
        
        while(True):
            try:
                임시_변수 = 사용자_목록[숫자_2]
            except:
                if(number != 1):
                    목록 = 목록 + '<br><a href="/userlog/n/' + str(number - 1) + '">(이전)'
                break
                
            if(관리자 == 1):
                DB_실행("select * from ban where block = '" + DB_인코딩(사용자_목록[숫자_2]['id']) + "'")
                차단인가 = DB_가져오기()
                if(차단인가):
                    차단_버튼 = ' <a href="/ban/' + URL_인코딩(사용자_목록[숫자_2]['id']) + '">(해제)</a>'
                else:
                    차단_버튼 = ' <a href="/ban/' + URL_인코딩(사용자_목록[숫자_2]['id']) + '">(차단)</a>'
            else:
                차단_버튼 = ''
                
            DB_실행("select * from data where title = '사용자:" + DB_인코딩(사용자_목록[숫자_2]['id']) + "'")
            자료 = DB_가져오기()
            if(자료):
                아이피 = '<a href="/w/' + URL_인코딩('사용자:' + 사용자_목록[숫자_2]['id']) + '">' + 사용자_목록[숫자_2]['id'] + '</a> <a href="/record/' + URL_인코딩(사용자_목록[숫자_2]['id']) + '/n/1">(기록)</a>'
            else:
                아이피 = '<a class="not_thing" href="/w/' + URL_인코딩('사용자:' + 사용자_목록[숫자_2]['id']) + '">' + 사용자_목록[숫자_2]['id'] + '</a> <a href="/record/' + URL_인코딩(사용자_목록[숫자_2]['id']) + '/n/1">(기록)</a>'
                
            목록 = 목록 + '<li>' + str(숫자_2 + 1) + '. ' + 아이피 + 차단_버튼 + '</li>'
            
            if(숫자_2 == 숫자_1):
                if(number == 1):
                    목록 = 목록 + '<br><a href="/userlog/n/' + str(number + 1) + '">(다음)'
                else:
                    목록 = 목록 + '<br><a href="/userlog/n/' + str(number - 1) + '">(이전) <a href="/userlog/n/' + str(number + 1) + '">(다음)'
                break
            else:
                숫자_2 += 1
                
        return 웹_디자인('index.html', logo = data['name'], data = 목록, title = '유저 가입 기록')
    else:
        return 웹_디자인('index.html', logo = data['name'], data = '', title = '유저 가입 기록')
        
@app.route('/backlink/<path:name>/n/<int:number>')
def 역링크(name = None, number = None):
    v = number * 50
    i = v - 50
    div = ''
    restart = 0
    
    DB_실행("select * from back where title = '" + DB_인코딩(name) + "' order by link asc")
    rows = DB_가져오기()
    if(rows):        
        while(True):
            try:
                a = rows[i]
            except:
                if(number != 1):
                    div = div + '<br><a href="/backlink/n/' + str(number - 1) + '">(이전)'
                break
                
            if(rows[i]['type'] == 'include'):
                DB_실행("select * from back where title = '" + DB_인코딩(name) + "' and link = '" + DB_인코딩(rows[i]['link']) + "' and type = ''")
                test = DB_가져오기()
                if(test):
                    restart = 1
                    
                    DB_실행("delete from back where title = '" + DB_인코딩(name) + "' and link = '" + DB_인코딩(rows[i]['link']) + "' and type = ''")
                    DB_갱신()
                
            if(not re.search('^사용자:', rows[i]['link'])):
                DB_실행("select * from data where title = '" + DB_인코딩(rows[i]['link']) + "'")
                row = DB_가져오기()
                if(row):
                    aa = row[0]['data']
                    aa = re.sub("(?P<in>\[include\((?P<out>(?:(?!\)\]|,).)*)((?:,\s?(?:[^)]*))+)?\)\])", "\g<in>\n\n[[\g<out>]]\n\n", aa)
                    aa = re.sub('^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)', '[[\g<in>]]', aa)
                    aa = 나무마크('', aa)
                    
                    if(re.search("<a(?:(?:(?!href=).)*)?href=\"\/w\/" + URL_인코딩(name) + "(?:\#[^\"]*)?\">([^<]*)<\/a>", aa)):
                        div = div + '<li><a href="/w/' + URL_인코딩(rows[i]['link']) + '">' + rows[i]['link'] + '</a>'
                        
                        if(rows[i]['type']):
                            div = div + ' (' + rows[i]['type'] + ')</li>'
                        else:
                            div = div + '</li>'
                            
                        if(i == v):
                            if(number == 1):
                                div = div + '<br><a href="/backlink/' + URL_인코딩(name) + '/n/' + str(number + 1) + '">(다음)'
                            else:
                                div = div + '<br><a href="/backlink/' + URL_인코딩(name) + '/n/' + str(number - 1) + '">(이전) <a href="/backlink/' + URL_인코딩(name) + '/n/' + str(number + 1) + '">(다음)'
                                
                            break
                        else:
                            i = i + 1
                    else:
                        DB_실행("delete from back where title = '" + DB_인코딩(name) + "' and link = '" + DB_인코딩(rows[i]['link']) + "'")
                        DB_갱신()
                        
                        i = i + 1
                        v = v + 1
                else:
                    DB_실행("delete from back where title = '" + DB_인코딩(name) + "' and link = '" + DB_인코딩(rows[i]['link']) + "'")
                    DB_갱신()
                    
                    i = i + 1
                    v = v + 1
            else:
                DB_실행("delete from back where title = '" + DB_인코딩(name) + "' and link = '" + DB_인코딩(rows[i]['link']) + "'")
                DB_갱신()
                
                i = i + 1
                v = v + 1
                
        if(restart == 1):
            return '<meta http-equiv="refresh" content="0;url=/backlink/' + URL_인코딩(name) + '/n/' + str(number) + '" />'
        else:    
            return 웹_디자인('index.html', logo = data['name'], data = div, title = name, page = URL_인코딩(name), sub = '역링크')
    else:
        return 웹_디자인('index.html', logo = data['name'], data = '', title = name, page = URL_인코딩(name), sub = '역링크')

@app.route('/recentdiscuss')
def 최근_토론():
    i = 0
    div = '<div>'
    
    DB_실행("select * from rd order by date desc limit 50")
    rows = DB_가져오기()
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
            
            div = div + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:50%;"><a href="/topic/' + URL_인코딩(rows[i]['title']) + '/sub/' + URL_인코딩(rows[i]['sub']) + '">' + title + '</a> (' + sub + ')</td><td style="text-align: center;width:50%;">' + rows[i]['date'] + '</td></tr></tbody></table>'
            
            i = i + 1
            
        return 웹_디자인('index.html', logo = data['name'], rows = div, tn = 12, title = '최근 토론내역')
    else:
        return 웹_디자인('index.html', logo = data['name'], rows = '', tn = 12, title = '최근 토론내역')
         
@app.route('/blocklog/n/<int:number>')
def blocklog(number = None):
    v = number * 50
    i = v - 50
    div = '<div>'
    
    DB_실행("select * from rb order by today desc")
    rows = DB_가져오기()
    if(rows):
        while(True):
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                
                if(number != 1):
                    div = div + '<br><a href="/blocklog/n/' + str(number - 1) + '">(이전)'
                    
                break
                
            why = rows[i]['why']
            why = re.sub('<', '&lt;', why)
            why = re.sub('>', '&gt;', why)
            
            b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))$", rows[i]['block'])
            if(b):
                ip = rows[i]['block'] + ' (대역)'
            else:
                ip = rows[i]['block']
                
            div = div + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:20%;">' + ip + '</a></td><td style="text-align: center;width:20%;">' + rows[i]['blocker'] + '</td><td style="text-align: center;width:20%;">' + rows[i]['end'] + '</td><td style="text-align: center;width:20%;">' + rows[i]['why'] + '</td><td style="text-align: center;width:20%;">' + rows[i]['today'] + '</td></tr></tbody></table>'
            
            if(i == v):
                div = div + '</div>'
                
                if(number == 1):
                    div = div + '<br><a href="/blocklog/n/' + str(number + 1) + '">(다음)'
                else:
                    div = div + '<br><a href="/blocklog/n/' + str(number - 1) + '">(이전) <a href="/blocklog/n/' + str(number + 1) + '">(다음)'
                    
                break
            else:
                i = i + 1
                
        return 웹_디자인('index.html', logo = data['name'], rows = div, tn = 20, title = '유저 차단 기록')
    else:
        return 웹_디자인('index.html', logo = data['name'], rows = '', tn = 20, title = '유저 차단 기록')

@app.route('/history/<path:name>/n/<int:number>', methods=['POST', 'GET'])
def 역사_보기(name = None, number = None):
    if(request.method == 'POST'):
        return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(name) + '/r/' + request.form["b"] + '/diff/' + request.form["a"] + '" />'
    else:
        select = ''
        v = number * 50
        i = v - 50
        div = '<div>'
        
        DB_실행("select * from history where title = '" + DB_인코딩(name) + "' order by id+0 desc")
        rows = DB_가져오기()
        if(rows):
            admin = 관리자_확인()
            
            while(True):
                style = ''
            
                try:
                    a = rows[i]
                except:
                    div = div + '</div>'
                    
                    if(number != 1):
                        div = div + '<br><a href="/history/' + URL_인코딩(name) + '/n/' + str(number - 1) + '">(이전)'
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
                    ip = rows[i]["ip"] + ' <a href="/record/' + URL_인코딩(rows[i]["ip"]) + '/n/1">(기록)</a>'
                else:
                    DB_실행("select * from data where title = '사용자:" + DB_인코딩(rows[i]['ip']) + "'")
                    row = DB_가져오기()
                    if(row):
                        ip = '<a href="/w/' + URL_인코딩('사용자:' + rows[i]['ip']) + '">' + rows[i]['ip'] + '</a> <a href="/record/' + URL_인코딩(rows[i]["ip"]) + '/n/1">(기록)</a>'
                    else:
                        ip = '<a class="not_thing" href="/w/' + URL_인코딩('사용자:' + rows[i]['ip']) + '">' + rows[i]['ip'] + '</a> <a href="/record/' + URL_인코딩(rows[i]["ip"]) + '/n/1">(기록)</a>'
                        
                if(admin == 1):
                    DB_실행("select * from user where id = '" + DB_인코딩(rows[i]['ip']) + "'")
                    row = DB_가져오기()
                    if(row):
                        if(row[0]['acl'] == 'owner' or row[0]['acl'] == 'admin'):
                            ban = ''
                        else:
                            DB_실행("select * from ban where block = '" + DB_인코딩(rows[i]['ip']) + "'")
                            row = DB_가져오기()
                            if(row):
                                ban = ' <a href="/ban/' + URL_인코딩(rows[i]['ip']) + '">(해제)</a>'
                            else:
                                ban = ' <a href="/ban/' + URL_인코딩(rows[i]['ip']) + '">(차단)</a>'
                    else:
                        DB_실행("select * from ban where block = '" + DB_인코딩(rows[i]['ip']) + "'")
                        row = DB_가져오기()
                        if(row):
                            ban = ' <a href="/ban/' + URL_인코딩(rows[i]['ip']) + '">(해제)</a>'
                        else:
                            ban = ' <a href="/ban/' + URL_인코딩(rows[i]['ip']) + '">(차단)</a>'
                            
                    if(소유자_확인() == 1):
                        DB_실행("select * from hidhi where title = '" + DB_인코딩(name) + "' and re = '" + DB_인코딩(rows[i]['id']) + "'")
                        row = DB_가져오기()
                        if(row):                            
                            ip = ip + ' (숨김)'                            
                            hidden = ' <a href="/history/' + URL_인코딩(name) + '/r/' + rows[i]['id'] + '/hidden">(공개)'
                        else:
                            hidden = ' <a href="/history/' + URL_인코딩(name) + '/r/' + rows[i]['id'] + '/hidden">(숨김)'
                    else:
                        DB_실행("select * from hidhi where title = '" + DB_인코딩(name) + "' and re = '" + DB_인코딩(rows[i]['id']) + "'")
                        row = DB_가져오기()
                        if(row):
                            ip = '숨김'
                            hidden = ''
                            send = '숨김'
                            ban = ''
                            style = 'display:none;'
                            v = v + 1
                        else:
                            hidden = ''
                else:
                    ban = ''
                    
                    DB_실행("select * from hidhi where title = '" + DB_인코딩(name) + "' and re = '" + DB_인코딩(rows[i]['id']) + "'")
                    row = DB_가져오기()
                    if(row):
                        ip = '숨김'
                        hidden = ''
                        send = '숨김'
                        ban = ''
                        style = 'display:none;'
                        v = v + 1
                    else:
                        hidden = ''                
                        
                div = div + '<table style="width: 100%;' + style + '"><tbody><tr><td style="text-align: center;width:33.33%;">r' + rows[i]['id'] + '</a> <a href="/w/' + URL_인코딩(rows[i]['title']) + '/r/' + rows[i]['id'] + '">(w)</a> <a href="/w/' + URL_인코딩(rows[i]['title']) + '/raw/' + rows[i]['id'] + '">(Raw)</a> <a href="/revert/' + URL_인코딩(rows[i]['title']) + '/r/' + rows[i]['id'] + '">(되돌리기)</a> (' + leng + ')</td><td style="text-align: center;width:33.33%;">' + ip + ban + hidden + '</td><td style="text-align: center;width:33.33%;">' + rows[i]['date'] + '</td></tr><tr><td colspan="3" style="text-align: center;width:100%;">' + send + '</td></tr></tbody></table>'
                
                if(i == v):
                    div = div + '</div>'
                    
                    if(number == 1):
                        div = div + '<br><a href="/history/' + URL_인코딩(name) + '/n/' + str(number + 1) + '">(다음)'
                    else:
                        div = div + '<br><a href="/history/' + URL_인코딩(name) + '/n/' + str(number - 1) + '">(이전) <a href="/history/' + URL_인코딩(name) + '/n/' + str(number + 1) + '">(다음)'
                        
                    break
                else:
                    i = i + 1
                    
            return 웹_디자인('index.html', logo = data['name'], rows = div, tn = 5, title = name, page = URL_인코딩(name), select = select, sub = '역사')
        else:
            return 웹_디자인('index.html', logo = data['name'], rows = '', tn = 5, title = name, page = URL_인코딩(name), select = select, sub = '역사')

@app.route('/search', methods=['POST'])
def search():
    DB_실행("select * from data where title = '" + DB_인코딩(request.form["search"]) + "'")
    rows = DB_가져오기()
    if(rows):
        return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(request.form["search"]) + '" />'
    else:
        DB_실행("select * from data where title like '%" + DB_인코딩(request.form["search"]) + "%'")
        rows = DB_가져오기()
        if(rows):
            i = 0
            
            div = '<li>문서가 없습니다. <a href="/w/' + URL_인코딩(request.form["search"]) + '">바로가기</a></li><br>'
            
            while(True):
                try:
                    div = div + '<li><a href="/w/' + URL_인코딩(rows[i]['title']) + '">' + rows[i]['title'] + '</a></li>'
                except:
                    break
                    
                i = i + 1
        else:
            return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(request.form["search"]) + '" />'
            
        return 웹_디자인('index.html', logo = data['name'], data = div, title = '검색')

@app.route('/w/<path:name>')
@app.route('/w/<path:name>/from/<path:redirect>')
def w(name = None, redirect = None):
    i = 0
    
    DB_실행("select * from rd where title = '" + DB_인코딩(name) + "' order by date asc")
    rows = DB_가져오기()
    while(True):
        try:
            a = rows[i]
        except:
            topic = ""
            break
            
        DB_실행("select * from stop where title = '" + DB_인코딩(rows[i]['title']) + "' and sub = '" + DB_인코딩(rows[i]['sub']) + "' and close = 'O'")
        row = DB_가져오기()
        if(not row):
            topic = "open"
            
            break
        else:
            i = i + 1
            
    acl = ''
    
    m = re.search("^(.*)\/(.*)$", name)
    if(m):
        g = m.groups()
        uppage = g[0]
        style = ""
    else:
        uppage = ""
        style = "display:none;"
        
    if(re.search("^분류:", name)):
        DB_실행("select * from cat where title = '" + DB_인코딩(name) + "' order by cat asc")
        rows = DB_가져오기()
        if(rows):
            div = ''
            i = 0
            
            while(True):
                try:
                    a = rows[i]
                except:
                    break
                    
                DB_실행("select * from data where title = '" + DB_인코딩(rows[i]['cat']) + "'")
                row = DB_가져오기()
                if(row):
                    aa = row[0]['data']                  
                    aa = 나무마크('', aa)
                    bb = re.search('<div style="width:100%;border: 1px solid #777;padding: 5px;margin-top: 1em;">분류:((?:(?!<\/div>).)*)<\/div>', aa)
                    if(bb):
                        cc = bb.groups()
                        
                        mm = re.search("^분류:(.*)", name)
                        if(mm):
                            ee = mm.groups()
                            
                            if(re.search("<a (class=\"not_thing\")? href=\"\/w\/" + URL_인코딩(name) + "\">" + ee[0] + "<\/a>", cc[0])):
                                div = div + '<li><a href="/w/' + URL_인코딩(rows[i]['cat']) + '">' + rows[i]['cat'] + '</a></li>'
                                
                                i = i + 1
                            else:
                                DB_실행("delete from cat where title = '" + DB_인코딩(name) + "' and cat = '" + DB_인코딩(rows[i]['cat']) + "'")
                                DB_갱신()
                                
                                i = i + 1
                        else:
                            DB_실행("delete from cat where title = '" + DB_인코딩(name) + "' and cat = '" + DB_인코딩(rows[i]['cat']) + "'")
                            DB_갱신()
                            
                            i = i + 1
                    else:
                        DB_실행("delete from cat where title = '" + DB_인코딩(name) + "' and cat = '" + DB_인코딩(rows[i]['cat']) + "'")
                        DB_갱신()
                        
                        i = i + 1
                else:
                    DB_실행("delete from cat where title = '" + DB_인코딩(name) + "' and cat = '" + DB_인코딩(rows[i]['cat']) + "'")
                    DB_갱신()
                    
                    i = i + 1
                    
            div = '<h2>분류</h2>' + div
            
            DB_실행("select * from data where title = '" + DB_인코딩(name) + "'")
            bb = DB_가져오기()
            if(bb):
                if(bb[0]['acl'] == 'admin'):
                    acl = '(관리자)'
                elif(bb[0]['acl'] == 'user'):
                    acl = '(유저)'
                else:
                    if(not acl):
                        acl = ''
                        
                enddata = 나무마크(name, bb[0]['data'])
                
                m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
                if(m):
                    result = m.groups()
                    left = result[0]
                else:
                    left = ''
                    
                return 웹_디자인('index.html', title = name, logo = data['name'], page = URL_인코딩(name), data = enddata + '<br>' + div, license = data['license'], tn = 1, uppage = uppage, style = style, acl = acl, topic = topic, redirect = redirect, login = 로그인_확인())
            else:
                return 웹_디자인('index.html', title = name, logo = data['name'], page = URL_인코딩(name), data = div, license = data['license'], tn = 1, uppage = uppage, style = style, acl = acl, topic = topic, redirect = redirect, login = 로그인_확인())
        else:
            return 웹_디자인('index.html', title = name, logo = data['name'], page = URL_인코딩(name), data = '분류 문서 없음', license = data['license'], tn = 1, uppage = uppage, style = style, acl = acl, topic = topic, redirect = redirect, login = 로그인_확인()), 404
    else:                    
        DB_실행("select * from data where title = '" + DB_인코딩(name) + "'")
        rows = DB_가져오기()
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
                
                DB_실행("select * from user where id = '" + DB_인코딩(g[0]) + "'")
                test = DB_가져오기()
                if(test):
                    if(test[0]['acl'] == 'owner'):
                        acl = '(소유자)'
                    elif(test[0]['acl'] == 'admin'):
                        acl = '(관리자)'

                DB_실행("select * from ban where block = '" + DB_인코딩(g[0]) + "'")
                user = DB_가져오기()
                if(user):
                    elsedata = '{{{#!wiki style="border:2px solid red;padding:10px;"\r\n{{{+2 {{{#red 이 사용자는 차단 당했습니다.}}}}}}\r\n\r\n차단 해제 일 : ' + user[0]['end'] + '[br]사유 : ' + user[0]['why'] + '}}}[br]' + rows[0]['data']
                else:
                    elsedata = rows[0]['data']
            else:
                elsedata = rows[0]['data']
                    
            enddata = 나무마크(name, elsedata)
            
            m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
            if(m):
                result = m.groups()
                left = result[0]
            else:
                left = ''
                
            return 웹_디자인('index.html', title = name, logo = data['name'], page = URL_인코딩(name), data = enddata, license = data['license'], tn = 1, acl = acl, left = left, uppage = uppage, style = style, topic = topic, redirect = redirect, login = 로그인_확인())
        else:
            m = re.search("^사용자:(.*)", name)
            if(m):
                g = m.groups()
                
                DB_실행("select * from ban where block = '" + DB_인코딩(g[0]) + "'")
                user = DB_가져오기()
                if(user):
                    elsedata = '{{{#!wiki style="border:2px solid red;padding:10px;"\r\n{{{+2 {{{#red 이 사용자는 차단 당했습니다.}}}}}}\r\n\r\n차단 해제 일 : ' + user[0]['end'] + '[br]사유 : ' + user[0]['why'] + '}}}[br]' + '문서 없음'
                else:
                    elsedata = '문서 없음'
            else:
                elsedata = '문서 없음'
            
            return 웹_디자인('index.html', title = name, logo = data['name'], page = URL_인코딩(name), data = 나무마크(name, elsedata), license = data['license'], tn = 1, uppage = uppage, style = style, acl = acl, topic = topic, redirect = redirect, login = 로그인_확인()), 404

@app.route('/w/<path:name>/r/<int:number>')
def rew(name = None, number = None):
    DB_실행("select * from hidhi where title = '" + DB_인코딩(name) + "' and re = '" + DB_인코딩(str(number)) + "'")
    row = DB_가져오기()
    if(row):
        if(소유자_확인() == 1):
            DB_실행("select * from history where title = '" + DB_인코딩(name) + "' and id = '" + str(number) + "'")
            rows = DB_가져오기()
            if(rows):
                enddata = 나무마크(name, rows[0]['data'])
                
                m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
                if(m):
                    result = m.groups()
                    left = result[0]
                else:
                    left = ''
                    
                return 웹_디자인('index.html', title = name, logo = data['name'], page = URL_인코딩(name), data = enddata, license = data['license'], tn = 6, left = left, sub = '옛 문서')
            else:
                return '<meta http-equiv="refresh" content="0;url=/history/' + URL_인코딩(name) + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'
    else:
        DB_실행("select * from history where title = '" + DB_인코딩(name) + "' and id = '" + str(number) + "'")
        rows = DB_가져오기()
        if(rows):
            enddata = 나무마크(name, rows[0]['data'])
            
            m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
            if(m):
                result = m.groups()
                left = result[0]
            else:
                left = ''
                
            return 웹_디자인('index.html', title = name, logo = data['name'], page = URL_인코딩(name), data = enddata, license = data['license'], tn = 6, left = left, sub = '옛 문서')
        else:
            return '<meta http-equiv="refresh" content="0;url=/history/' + URL_인코딩(name) + '" />'

@app.route('/w/<path:name>/raw/<int:number>')
def reraw(name = None, number = None):
    DB_실행("select * from hidhi where title = '" + DB_인코딩(name) + "' and re = '" + DB_인코딩(str(number)) + "'")
    row = DB_가져오기()
    if(row):
        if(소유자_확인() == 1):
            DB_실행("select * from history where title = '" + DB_인코딩(name) + "' and id = '" + str(number) + "'")
            rows = DB_가져오기()
            if(rows):
                enddata = re.sub('<', '&lt;', rows[0]['data'])
                enddata = re.sub('>', '&gt;', enddata)
                enddata = re.sub('"', '&quot;', enddata)
                
                enddata = '<pre>' + enddata + '</pre>'
                
                return 웹_디자인('index.html', title = name, logo = data['name'], page = URL_인코딩(name), data = enddata, license = data['license'])
            else:
                return '<meta http-equiv="refresh" content="0;url=/history/' + URL_인코딩(name) + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'
    else:
        DB_실행("select * from history where title = '" + DB_인코딩(name) + "' and id = '" + str(number) + "'")
        rows = DB_가져오기()
        if(rows):
            enddata = re.sub('<', '&lt;', rows[0]['data'])
            enddata = re.sub('>', '&gt;', enddata)
            enddata = re.sub('"', '&quot;', enddata)
            
            enddata = '<pre>' + enddata + '</pre>'
            
            return 웹_디자인('index.html', title = name, logo = data['name'], page = URL_인코딩(name), data = enddata, license = data['license'])
        else:
            return '<meta http-equiv="refresh" content="0;url=/history/' + URL_인코딩(name) + '" />'

@app.route('/raw/<path:name>')
def raw(name = None):
    DB_실행("select * from data where title = '" + DB_인코딩(name) + "'")
    rows = DB_가져오기()
    if(rows):
        enddata = re.sub('<', '&lt;', rows[0]['data'])
        enddata = re.sub('>', '&gt;', enddata)
        enddata = re.sub('"', '&quot;', enddata)
        
        enddata = '<pre>' + enddata + '</pre>'
        
        return 웹_디자인('index.html', title = name, logo = data['name'], page = URL_인코딩(name), data = enddata, license = data['license'], tn = 7, sub = 'Raw')
    else:
        return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(name) + '" />'

@app.route('/revert/<path:name>/r/<int:number>', methods=['POST', 'GET'])
def revert(name = None, number = None):
    if(request.method == 'POST'):
        DB_실행("select * from hidhi where title = '" + DB_인코딩(name) + "' and re = '" + DB_인코딩(str(number)) + "'")
        row = DB_가져오기()
        if(row):
            if(소유자_확인() == 1):        
                DB_실행("select * from history where title = '" + DB_인코딩(name) + "' and id = '" + str(number) + "'")
                rows = DB_가져오기()
                if(rows):
                    ip = 아이피_확인(request)
                    can = ACL_체크(ip, name)
                    
                    if(can == 1):
                        return '<meta http-equiv="refresh" content="0;url=/ban" />'
                    else:
                        today = 시간()
                        
                        DB_실행("select * from data where title = '" + DB_인코딩(name) + "'")
                        row = DB_가져오기()
                        if(row):
                            leng = 길이_확인(len(row[0]['data']), len(rows[0]['data']))
                            
                            DB_실행("update data set data = '" + DB_인코딩(rows[0]['data']) + "' where title = '" + DB_인코딩(name) + "'")
                            DB_갱신()
                        else:
                            leng = '+' + str(len(rows[0]['data']))
                            
                            DB_실행("insert into data (title, data, acl) value ('" + DB_인코딩(name) + "', '" + DB_인코딩(rows[0]['data']) + "', '')")
                            DB_갱신()
                        역사_추가(name, rows[0]['data'], today, ip, '문서를 ' + str(number) + '판으로 되돌렸습니다.', leng)
                        return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(name) + '" />'
                else:
                    return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(name) + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/3" />'
        else:
            DB_실행("select * from history where title = '" + DB_인코딩(name) + "' and id = '" + str(number) + "'")
            rows = DB_가져오기()
            if(rows):
                ip = 아이피_확인(request)
                can = ACL_체크(ip, name)
                
                if(can == 1):
                    return '<meta http-equiv="refresh" content="0;url=/ban" />'
                else:
                    today = 시간()
                    
                    DB_실행("select * from data where title = '" + DB_인코딩(name) + "'")
                    row = DB_가져오기()
                    if(row):
                        leng = 길이_확인(len(row[0]['data']), len(rows[0]['data']))
                        
                        DB_실행("update data set data = '" + DB_인코딩(rows[0]['data']) + "' where title = '" + DB_인코딩(name) + "'")
                        DB_갱신()
                    else:
                        leng = '+' + str(len(rows[0]['data']))
                        
                        DB_실행("insert into data (title, data, acl) value ('" + DB_인코딩(name) + "', '" + DB_인코딩(rows[0]['data']) + "', '')")
                        DB_갱신()
                    역사_추가(name, rows[0]['data'], today, ip, '문서를 ' + str(number) + '판으로 되돌렸습니다.', leng)
                    return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(name) + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(name) + '" />'            
    else:
        DB_실행("select * from hidhi where title = '" + DB_인코딩(name) + "' and re = '" + DB_인코딩(str(number)) + "'")
        row = DB_가져오기()
        if(row):
            if(소유자_확인() == 1):
                ip = 아이피_확인(request)
                can = ACL_체크(ip, name)
                
                if(can == 1):
                    return '<meta http-equiv="refresh" content="0;url=/ban" />'
                else:
                    DB_실행("select * from history where title = '" + DB_인코딩(name) + "' and id = '" + str(number) + "'")
                    rows = DB_가져오기()
                    if(rows):
                        return 웹_디자인('index.html', title = name, logo = data['name'], page = URL_인코딩(name), r = URL_인코딩(str(number)), tn = 13, plus = '정말 되돌리시겠습니까?', sub = '되돌리기', login = 로그인_확인())
                    else:
                        return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(name) + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/3" />'
        else:            
            ip = 아이피_확인(request)
            can = ACL_체크(ip, name)
            
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                DB_실행("select * from history where title = '" + DB_인코딩(name) + "' and id = '" + str(number) + "'")
                rows = DB_가져오기()
                if(rows):
                    return 웹_디자인('index.html', title = name, logo = data['name'], page = URL_인코딩(name), r = URL_인코딩(str(number)), tn = 13, plus = '정말 되돌리시겠습니까?', sub = '되돌리기', login = 로그인_확인())
                else:
                    return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(name) + '" />'
                        
@app.route('/edit/<path:name>', methods=['POST', 'GET'])
def edit(name = None):
    if(request.method == 'POST'):
        m = re.search('(?:[^A-Za-zㄱ-힣0-9 ])', request.form["send"])
        
        if(m):
            return '<meta http-equiv="refresh" content="0;url=/error/17" />'
        else:
            today = 시간()
            
            content = 세이브마크(request.form["content"])
            
            DB_실행("select * from data where title = '" + DB_인코딩(name) + "'")
            rows = DB_가져오기()
            if(rows):
                if(rows[0]['data'] == content):
                    return '<meta http-equiv="refresh" content="0;url=/error/18" />'
                else:
                    ip = 아이피_확인(request)
                    can = ACL_체크(ip, name)
                    
                    if(can == 1):
                        return '<meta http-equiv="refresh" content="0;url=/ban" />'
                    else:                        
                        leng = 길이_확인(len(rows[0]['data']), len(content))
                        역사_추가(name, content, today, ip, request.form["send"], leng)
                        
                        DB_실행("update data set data = '" + DB_인코딩(content) + "' where title = '" + DB_인코딩(name) + "'")
                        DB_갱신()
            else:
                ip = 아이피_확인(request)
                can = ACL_체크(ip, name)
                
                if(can == 1):
                    return '<meta http-equiv="refresh" content="0;url=/ban" />'
                else:
                    leng = '+' + str(len(content))
                    역사_추가(name, content, today, ip, request.form["send"], leng)
                    
                    DB_실행("insert into data (title, data, acl) value ('" + DB_인코딩(name) + "', '" + DB_인코딩(content) + "', '')")
                    DB_갱신()
                    
            틀_확인(name, content)
            
            return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(name) + '" />'
    else:
        ip = 아이피_확인(request)
        can = ACL_체크(ip, name)
        
        if(can == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            DB_실행("select * from data where title = '" + DB_인코딩(data["help"]) + "'")
            rows = DB_가져오기()
            if(rows):
                newdata = re.sub('^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)', ' * \g<in> 문서로 넘겨주기', rows[0]["data"])
                left = 나무마크(name, newdata)
            else:
                left = ''
                
            DB_실행("select * from data where title = '" + DB_인코딩(name) + "'")
            rows = DB_가져오기()
            if(rows):
                return 웹_디자인('index.html', title = name, logo = data['name'], page = URL_인코딩(name), data = rows[0]['data'], tn = 2, left = left, sub = '편집', login = 로그인_확인())
            else:
                return 웹_디자인('index.html', title = name, logo = data['name'], page = URL_인코딩(name), data = '', tn = 2, left = left, sub = '편집', login = 로그인_확인())
                
@app.route('/edit/<path:name>/section/<int:number>', methods=['POST', 'GET'])
def secedit(name = None, number = None):
    if(request.method == 'POST'):
        m = re.search('(?:[^A-Za-zㄱ-힣0-9 ])', request.form["send"])
        if(m):
            return '<meta http-equiv="refresh" content="0;url=/error/17" />'
        else:
            today = 시간()
            
            content = 세이브마크(request.form["content"])
            
            DB_실행("select * from data where title = '" + DB_인코딩(name) + "'")
            rows = DB_가져오기()
            if(rows):
                if(request.form["otent"] == content):
                    return '<meta http-equiv="refresh" content="0;url=/error/18" />'
                else:
                    ip = 아이피_확인(request)
                    can = ACL_체크(ip, name)
                    
                    if(can == 1):
                        return '<meta http-equiv="refresh" content="0;url=/ban" />'
                    else:                        
                        leng = 길이_확인(len(request.form['otent']), len(content))
                        content = rows[0]['data'].replace(request.form['otent'], content)
                        역사_추가(name, content, today, ip, request.form["send"], leng)
                        
                        DB_실행("update data set data = '" + DB_인코딩(content) + "' where title = '" + DB_인코딩(name) + "'")
                        DB_갱신()
                        
                    틀_확인(name, content)
                    
                    return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(name) + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(name) + '" />'
    else:
        ip = 아이피_확인(request)
        can = ACL_체크(ip, name)
        
        if(can == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            DB_실행("select * from data where title = '" + DB_인코딩(data["help"]) + "'")
            rows = DB_가져오기()
            if(rows):
                newdata = re.sub('^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)', ' * \g<in> 문서로 넘겨주기', rows[0]["data"])
                
                left = 나무마크(name, newdata)
            else:
                left = ''
                
            DB_실행("select * from data where title = '" + DB_인코딩(name) + "'")
            rows = DB_가져오기()
            if(rows):
                i = 0
                j = 0
                
                gdata = rows[0]['data'] + '\r\n'
                
                while(True):
                    m = re.search("((?:={1,6})\s?(?:[^=]*)\s?(?:={1,6})(?:\s+)?\n(?:(?:(?:(?!(?:={1,6})\s?(?:[^=]*)\s?(?:={1,6})(?:\s+)?\n).)*)(?:\n)?)+)", gdata)
                    if(m):
                        if(i == number - 1):
                            g = m.groups()
                            
                            gdata = re.sub("\r\n$", "", g[0])
                            
                            break
                        else:
                            gdata = re.sub("((?:={1,6})\s?(?:[^=]*)\s?(?:={1,6})(?:\s+)?\n(?:(?:(?:(?!(?:={1,6})\s?(?:[^=]*)\s?(?:={1,6})(?:\s+)?\n).)*)(?:\n)?)+)", "", gdata, 1)
                            
                            i = i + 1
                    else:
                        j = 1
                        
                        break
                        
                if(j == 0):
                    return 웹_디자인('index.html', title = name, logo = data['name'], page = URL_인코딩(name), data = gdata, tn = 2, left = left, section = 1, number = number, sub = '편집', login = 로그인_확인())
                else:
                    return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(name) + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(name) + '" />'
                
@app.route('/preview/<path:name>', methods=['POST'])
def 미리보기(name = None):
    ip = 아이피_확인(request)
    can = ACL_체크(ip, name)
    if(can == 1):
        return '<meta http-equiv="refresh" content="0;url=/ban" />'
    else:            
        newdata = request.form["content"]
        newdata = re.sub('^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)', ' * \g<in> 문서로 넘겨주기', newdata)
        enddata = 나무마크(name, newdata)
        
        DB_실행("select * from data where title = '" + DB_인코딩(data["help"]) + "'")
        rows = DB_가져오기()
        if(rows):
            newdata = re.sub('^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)', ' * \g<in> 문서로 넘겨주기', rows[0]["data"])
            
            left = 나무마크(name, newdata)
        else:
            left = ''
            
        return 웹_디자인('index.html', title = name, logo = data['name'], page = URL_인코딩(name), data = request.form["content"], tn = 2, preview = 1, enddata = enddata, left = left, sub = '미리보기', login = 로그인_확인())
        
@app.route('/preview/<path:name>/section/<int:number>', methods=['POST'])
def 문단_미리보기(name = None, number = None):
    ip = 아이피_확인(request)
    can = ACL_체크(ip, name)
    if(can == 1):
        return '<meta http-equiv="refresh" content="0;url=/ban" />'
    else:
        if(re.search('\.', ip)):
            notice = '비 로그인 상태 입니다. 비 로그인으로 편집시 아이피가 역사에 기록 됩니다. 편집 시 동의 함으로 간주 됩니다.'
        else:
            notice = ''
        newdata = request.form["content"]
        newdata = re.sub('^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)', ' * \g<in> 문서로 넘겨주기', newdata)
        enddata = 나무마크(name, newdata)
        DB_실행("select * from data where title = '" + DB_인코딩(data["help"]) + "'")
        rows = DB_가져오기()
        if(rows):
            newdata = re.sub('^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)', ' * \g<in> 문서로 넘겨주기', rows[0]["data"])
            left = 나무마크(name, newdata)
        else:
            left = ''
        return 웹_디자인('index.html', title = name, logo = data['name'], page = URL_인코딩(name), data = request.form["content"], tn = 2, preview = 1, enddata = enddata, left = left, notice = notice, section = 1, number = number, odata = request.form["otent"], sub = '미리보기')

@app.route('/delete/<path:name>', methods=['POST', 'GET'])
def 문서_삭제(name = None):
    if(request.method == 'POST'):
        DB_실행("select * from data where title = '" + DB_인코딩(name) + "'")
        rows = DB_가져오기()
        if(rows):
            ip = 아이피_확인(request)
            can = ACL_체크(ip, name)
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                today = 시간()
                leng = '-' + str(len(rows[0]['data']))
                역사_추가(name, '', today, ip, '문서를 삭제 했습니다.', leng)
                DB_실행("delete from data where title = '" + DB_인코딩(name) + "'")
                DB_갱신()
                return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(name) + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(name) + '" />'
    else:
        DB_실행("select * from data where title = '" + DB_인코딩(name) + "'")
        rows = DB_가져오기()
        if(rows):
            ip = 아이피_확인(request)
            can = ACL_체크(ip, name)
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                return 웹_디자인('index.html', title = name, logo = data['name'], page = URL_인코딩(name), tn = 8, plus = '정말 삭제 하시겠습니까?', sub = '삭제', login = 로그인_확인())
        else:
            return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(name) + '" />'

@app.route('/move/<path:name>', methods=['POST', 'GET'])
def 문서_이동(name = None):
    if(request.method == 'POST'):
        DB_실행("select * from data where title = '" + DB_인코딩(name) + "'")
        rows = DB_가져오기()
        if(rows):
            ip = 아이피_확인(request)
            can = ACL_체크(ip, name)
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                today = 시간()
                leng = '0'
                DB_실행("select * from history where title = '" + DB_인코딩(request.form["title"]) + "'")
                row = DB_가져오기()
                if(row):
                    return '<meta http-equiv="refresh" content="0;url=/error/19" />'
                else:
                    역사_추가(name, rows[0]['data'], today, ip, '<a href="/w/' + URL_인코딩(name) + '">' + name + '</a> 문서를 <a href="/w/' + URL_인코딩(request.form["title"]) + '">' + request.form["title"] + '</a> 문서로 이동 했습니다.', leng)
                    DB_실행("update data set title = '" + DB_인코딩(request.form["title"]) + "' where title = '" + DB_인코딩(name) + "'")
                    DB_실행("update history set title = '" + DB_인코딩(request.form["title"]) + "' where title = '" + DB_인코딩(name) + "'")
                    DB_갱신()
                    return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(request.form["title"]) + '" />'
        else:
            ip = 아이피_확인(request)
            can = ACL_체크(ip, name)
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                today = 시간()
                leng = '0'
                DB_실행("select * from history where title = '" + DB_인코딩(request.form["title"]) + "'")
                row = DB_가져오기()
                if(row):
                     return '<meta http-equiv="refresh" content="0;url=/error/19" />'
                else:
                    역사_추가(name, '', today, ip, '<a href="/w/' + URL_인코딩(name) + '">' + name + '</a> 문서를 <a href="/w/' + URL_인코딩(request.form["title"]) + '">' + request.form["title"] + '</a> 문서로 이동 했습니다.', leng)
                    DB_실행("update history set title = '" + DB_인코딩(request.form["title"]) + "' where title = '" + DB_인코딩(name) + "'")
                    DB_갱신()
                    return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(request.form["title"]) + '" />'
    else:
        ip = 아이피_확인(request)
        can = ACL_체크(ip, name)
        if(can == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            return 웹_디자인('index.html', title = name, logo = data['name'], page = URL_인코딩(name), tn = 9, plus = '정말 이동 하시겠습니까?', sub = '이동', login = 로그인_확인())

@app.route('/other')
def 나머지():
    return 웹_디자인('index.html', title = '기타 메뉴', logo = data['name'], data = '<h2 style="margin-top: 0px;">기록</h2><li><a href="/blocklog/n/1">유저 차단 기록</a></li><li><a href="/userlog/n/1">유저 가입 기록</a></li><li><a href="/manager/6">유저 기록</a></li><h2>기타</h2><li><a href="/titleindex">모든 문서</a></li><li><a href="/upload">업로드</a></li><li><a href="/adminlist">관리자 목록</a></li><li><a href="/manager/1">관리자 메뉴</a></li><br>이 오픈나무의 버전은 <a href="https://github.com/2DU/openNAMU/blob/master/version.md">1.8.8c</a> 입니다.')
    
@app.route('/manager/<int:num>', methods=['POST', 'GET'])
def 관리_기능(num = None):
    if(num == 1):
        return 웹_디자인('index.html', title = '관리자 메뉴', logo = data['name'], data = '<h2 style="margin-top: 0px;">관리자 및 소유자</h2><li><a href="/manager/2">문서 ACL</a></li><li><a href="/manager/3">유저 체크</a></li><li><a href="/manager/4">유저 차단</a></li><h2>소유자</h2><li><a href="/manager/5">관리자 권한 주기</a></li><h2>기타</h2><li>이 메뉴에 없는 기능은 해당 문서의 역사나 토론에서 바로 사용 가능함</li>')
    elif(num == 2):
        if(request.method == 'POST'):
            return '<meta http-equiv="refresh" content="0;url=/acl/' + URL_인코딩(request.form["name"]) + '" />'
        else:
            return 웹_디자인('index.html', title = 'ACL 이동', logo = data['name'], data = '<form id="usrform" method="POST" action="/manager/2"><input name="name" type="text"><br><br><button class="btn btn-primary" type="submit">이동</button></form>')
    elif(num == 3):
        if(request.method == 'POST'):
            return '<meta http-equiv="refresh" content="0;url=/check/' + URL_인코딩(request.form["name"]) + '" />'
        else:
            return 웹_디자인('index.html', title = '체크 이동', logo = data['name'], data = '<form id="usrform" method="POST" action="/manager/3"><input name="name" type="text"><br><br><button class="btn btn-primary" type="submit">이동</button></form>')
    elif(num == 4):
        if(request.method == 'POST'):
            return '<meta http-equiv="refresh" content="0;url=/ban/' + URL_인코딩(request.form["name"]) + '" />'
        else:
            return 웹_디자인('index.html', title = '차단 이동', logo = data['name'], data = '<form id="usrform" method="POST" action="/manager/4"><input name="name" type="text"><br><br><button class="btn btn-primary" type="submit">이동</button><br><br><span>아이피 앞 두자리 (XXX.XXX) 입력하면 대역 차단</span></form>')
    elif(num == 5):
        if(request.method == 'POST'):
            return '<meta http-equiv="refresh" content="0;url=/admin/' + URL_인코딩(request.form["name"]) + '" />'
        else:
            return 웹_디자인('index.html', title = '권한 이동', logo = data['name'], data = '<form id="usrform" method="POST" action="/manager/5"><input name="name" type="text"><br><br><button class="btn btn-primary" type="submit">이동</button></form>')   
    elif(num == 6):
        if(request.method == 'POST'):
            return '<meta http-equiv="refresh" content="0;url=/record/' + URL_인코딩(request.form["name"]) + '/n/1" />'
        else:
            return 웹_디자인('index.html', title = '기록 이동', logo = data['name'], data = '<form id="usrform" method="POST" action="/manager/6"><input name="name" type="text"><br><br><button class="btn btn-primary" type="submit">이동</button></form>')    
    else:
        return '<meta http-equiv="refresh" content="0;url=/" />'
        
@app.route('/titleindex')
def 모든_문서():
    숫자 = 0
    데이터 = '<div>'
    DB_실행("select title from data order by title asc")
    문서명 = DB_가져오기()
    if(문서명):
        while(True):
            try:
                덤 = 문서명[숫자]
            except:
                break

            데이터 = 데이터 + '<li>' + str(숫자 + 1) + '. <a href="/w/' + URL_인코딩(문서명[숫자]['title']) + '">' + 문서명[숫자]['title'] + '</a></li>'
            
            숫자 += 1

        데이터 = 데이터 + '</div>'

        return 웹_디자인('index.html', logo = data['name'], rows = 데이터 + '<br><span>이 위키에는 총 ' + str(숫자) + '개의 문서가 있습니다.</span>', tn = 4, title = '모든 문서')
    else:
        return 웹_디자인('index.html', logo = data['name'], rows = '', tn = 4, title = '모든 문서')

@app.route('/topic/<path:name>', methods=['POST', 'GET'])
def 토론_목록(name = None):
    if(request.method == 'POST'):
        return '<meta http-equiv="refresh" content="0;url=/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(request.form["topic"]) + '" />'
    else:
        div = '<div>'
        i = 0
        j = 1
        DB_실행("select * from rd where title = '" + DB_인코딩(name) + "' order by date asc")
        rows = DB_가져오기()
        while(True):
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                break
                
            DB_실행("select * from topic where title = '" + DB_인코딩(rows[i]['title']) + "' and sub = '" + DB_인코딩(rows[i]['sub']) + "' and id = '1' order by sub asc")
            aa = DB_가져오기()
            
            indata = 나무마크(name, aa[0]['data'])
            
            if(aa[0]['block'] == 'O'):
                indata = '블라인드 되었습니다.'
                block = 'style="background: gainsboro;"'
            else:
                block = ''

            ip = 아이디_파싱(aa[0]['ip'])
                
            DB_실행("select * from stop where title = '" + DB_인코딩(rows[i]['title']) + "' and sub = '" + DB_인코딩(rows[i]['sub']) + "' and close = 'O'")
            row = DB_가져오기()
            if(not row):
                div = div + '<h2><a href="/topic/' + URL_인코딩(rows[i]['title']) + '/sub/' + URL_인코딩(rows[i]['sub']) + '">' + str(j) + '. ' + rows[i]['sub'] + '</a></h2><table id="toron"><tbody><tr><td id="toroncolorgreen"><a href="javascript:void(0);" id="1">#1</a> ' + ip + ' <span style="float:right;">' + aa[0]['date'] + '</span></td></tr><tr><td ' + block + '>' + indata + '</td></tr></tbody></table><br>'
                j = j + 1
                
            i = i + 1
            
        return 웹_디자인('index.html', title = name, page = URL_인코딩(name), logo = data['name'], plus = div, tn = 10, list = 1, sub = '토론 목록')
        
@app.route('/topic/<path:name>/close')
def 닫힌_토론_목록(name = None):
    div = '<div>'
    i = 0
    
    DB_실행("select * from stop where title = '" + DB_인코딩(name) + "' and close = 'O' order by sub asc")
    rows = DB_가져오기()
    while(True):
        try:
            a = rows[i]
        except:
            div = div + '</div>'
            break
            
        DB_실행("select * from topic where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(rows[i]['sub']) + "' and id = '1'")
        row = DB_가져오기()
        if(row):
            indata = 나무마크(name, row[0]['data'])
            
            if(row[0]['block'] == 'O'):
                indata = '블라인드 되었습니다.'
                block = 'style="background: gainsboro;"'
            else:
                block = ''

            아이디 = 아이디_파싱(row[0]['ip'])
                
            div = div + '<h2><a href="/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(rows[i]['sub']) + '">' + str((i + 1)) + '. ' + rows[i]['sub'] + '</a></h2><table id="toron"><tbody><tr><td id="toroncolorgreen"><a href="javascript:void(0);" id="1">#1</a> ' + 아이디 + ' <span style="float:right;">' + row[0]['date'] + '</span></td></tr><tr><td ' + block + '>' + indata + '</td></tr></tbody></table><br>'
            
        i += 1
        
    return 웹_디자인('index.html', title = name, page = URL_인코딩(name), logo = data['name'], plus = div, tn = 10, sub = '닫힌 토론')

@app.route('/topic/<path:name>/agree')
def 합의된_토론_목록(name = None):
    보여줄_내용 = '<div>'
    숫자 = 0
    
    DB_실행("select * from agreedis where title = '" + DB_인코딩(name) + "' order by sub asc")
    합의_토론 = DB_가져오기()
    while(True):
        try:
            덤 = 합의_토론[숫자]
        except:
            보여줄_내용 = 보여줄_내용 + '</div>'
            break
            
        DB_실행("select * from topic where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(합의_토론[숫자]['sub']) + "' and id = '1'")
        내용 = DB_가져오기()
        if(내용):
            내용_파싱 = 나무마크(name, 내용[0]['data'])
            
            if(내용[0]['block'] == 'O'):
                내용_파싱 = '블라인드 되었습니다.'
                가리기 = 'style="background: gainsboro;"'
            else:
                가리기 = ''

            아이디 = 아이디_파싱(내용[0]['ip'])
                
            보여줄_내용 = 보여줄_내용 + '<h2><a href="/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(내용[숫자]['sub']) + '">' + str((숫자 + 1)) + '. ' + 내용[숫자]['sub'] + '</a></h2><table id="toron"><tbody><tr><td id="toroncolorgreen"><a href="javascript:void(0);" id="1">#1</a> ' + 아이디 + ' <span style="float:right;">' + 내용[0]['date'] + '</span></td></tr><tr><td ' + 가리기 + '>' + 내용_파싱 + '</td></tr></tbody></table><br>'
            
        숫자 += 1
        
    return 웹_디자인('index.html', title = name, page = URL_인코딩(name), logo = data['name'], plus = 보여줄_내용, tn = 10, sub = '합의된 토론')

@app.route('/topic/<path:name>/sub/<path:sub>', methods=['POST', 'GET'])
def 토론(name = None, sub = None):
    if(request.method == 'POST'):
        DB_실행("select * from topic where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "' order by id+0 desc limit 1")
        rows = DB_가져오기()
        if(rows):
            number = int(rows[0]['id']) + 1
        else:
            number = 1
            
        ip = 아이피_확인(request)
        ban = 토론자_체크(ip, name, sub)
        admin = 관리자_확인()
        
        if(ban == 1 and not admin == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            DB_실행("select * from user where id = '" + DB_인코딩(ip) + "'")
            rows = DB_가져오기()
            if(rows):
                if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
                    ip = ip + ' - Admin'
                    
            today = 시간()
            최근_토론_추가(name, sub, today)
            
            aa = request.form["content"]
            aa = re.sub("\[\[(분류:(?:(?:(?!\]\]).)*))\]\]", "[br]", aa)
            aa = 세이브마크(aa)
            
            DB_실행("insert into topic (id, title, sub, data, date, ip, block) value ('" + str(number) + "', '" + DB_인코딩(name) + "', '" + DB_인코딩(sub) + "', '" + DB_인코딩(aa) + "', '" + today + "', '" + ip + "', '')")
            DB_갱신()
            
            return '<meta http-equiv="refresh" content="0;url=/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '" />'
    else:
        style = ''
        
        ip = 아이피_확인(request)
        ban = 토론자_체크(ip, name, sub)
        admin = 관리자_확인()

        DB_실행("select * from stop where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "' and close = 'O'")
        닫음 = DB_가져오기()

        DB_실행("select * from stop where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "' and close = ''")
        정지 = DB_가져오기()
        
        if(admin == 1):
            div = '<div>'
            
            if(닫음):
                div = div + '<a href="/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '/close">(토론 열기)</a> '
            else:
                div = div + '<a href="/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '/close">(토론 닫기)</a> '
            
            if(정지):
                div = div + '<a href="/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '/stop">(토론 재개)</a> '
            else:
                div = div + '<a href="/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '/stop">(토론 정지)</a> '

            DB_실행("select * from agreedis where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "'")
            합의 = DB_가져오기()
            if(합의):
                div = div + '<a href="/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '/agree">(합의 취소)</a>'
            else:
                div = div + '<a href="/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '/agree">(합의 완료)</a>'
            
            div = div + '<br><br>'
        else:
            div = '<div>'
        
        if(닫음 or 정지):
            if(not admin == 1):
                style = 'display:none;'
        
        DB_실행("select * from topic where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "' order by id+0 asc")
        rows = DB_가져오기()

        DB_실행("select * from distop where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "' order by id+0 asc")
        공지 = DB_가져오기()

        i = 0


        if(공지):
            while(True):
                try:
                    a = 공지[i]
                except:
                    break

                num = int(공지[i]['id']) - 1

                if(i == 0):
                    start = rows[num]['ip']
                    
                공지_데이터 = 나무마크('', rows[num]['data'])
                공지_데이터 = re.sub("(?P<in>#(?:[0-9]*))", '<a href="\g<in>">\g<in></a>', 공지_데이터)
                        
                ip = 아이디_파싱(rows[num]['ip'])
                                   
                div = div + '<table id="toron"><tbody><tr><td id="toroncolorred"><a href="#' + 공지[i]['id'] + '" id="' + 공지[i]['id'] + '-nt">#' + 공지[i]['id'] + '</a> ' + ip + ' <span style="float:right;">' + rows[num]['date'] + '</span></td></tr><tr><td>' + 공지_데이터 + '</td></tr></tbody></table><br>'
                    
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
                
            indata = 나무마크('', rows[i]['data'])
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
                        isblock = ' <a href="/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '/b/' + str(i + 1) + '">(해제)</a>'
                    else:
                        isblock = ' <a href="/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '/b/' + str(i + 1) + '">(블라인드)</a>'

                    DB_실행("select * from distop where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "' and id = '" + DB_인코딩(str(i + 1)) + "'")
                    row = DB_가져오기()
                    if(row):
                        isblock = isblock + ' <a href="/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '/notice/' + str(i + 1) + '">(해제)</a>'
                    else:
                        isblock = isblock + ' <a href="/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '/notice/' + str(i + 1) + '">(공지)</a>'
                        
                    n = re.search("\- (?:Admin)$", rows[i]['ip'])
                    if(n):
                        ban = isblock
                    else:
                        DB_실행("select * from ban where block = '" + DB_인코딩(rows[i]['ip']) + "'")
                        row = DB_가져오기()
                        if(row):
                            ban = ' <a href="/ban/' + URL_인코딩(rows[i]['ip']) + '">(해제)</a>' + isblock
                        else:
                            ban = ' <a href="/ban/' + URL_인코딩(rows[i]['ip']) + '">(차단)</a>' + isblock
                else:
                    ban = ""

            ip = 아이디_파싱(rows[i]['ip'])
                    
            if(rows[i]['ip'] == start):
                j = i + 1
                
                div = div + '<table id="toron"><tbody><tr><td id="toroncolorgreen"><a href="javascript:void(0);" id="' + str(j) + '">#' + str(j) + '</a> ' + ip + ban + ' <span style="float:right;">' + rows[i]['date'] + '</span></td></tr><tr><td ' + block + '>' + indata + '</td></tr></tbody></table><br>'
            else:
                j = i + 1
                
                div = div + '<table id="toron"><tbody><tr><td id="toroncolor"><a href="javascript:void(0);" id="' + str(j) + '">#' + str(j) + '</a> ' + ip + ban + ' <span style="float:right;">' + rows[i]['date'] + '</span></td></tr><tr><td ' + block + '>' + indata + '</td></tr></tbody></table><br>'
                
            i = i + 1
            
        return 웹_디자인('index.html', title = name, page = URL_인코딩(name), suburl = URL_인코딩(sub), toron = sub, logo = data['name'], rows = div, tn = 11, ban = ban, style = style, sub = '토론', login = 로그인_확인())

@app.route('/topic/<path:name>/sub/<path:sub>/b/<int:number>')
def 토론_블라인드(name = None, sub = None, number = None):
    if(관리자_확인() == 1):
        DB_실행("select * from topic where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "' and id = '" + str(number) + "'")
        가리기 = DB_가져오기()
        if(가리기):
            if(가리기[0]['block'] == 'O'):
                DB_실행("update topic set block = '' where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "' and id = '" + str(number) + "'")
            else:
                DB_실행("update topic set block = 'O' where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "' and id = '" + str(number) + "'")
            DB_갱신()
            
            최근_토론_추가(name, sub, 시간())
            
            return '<meta http-equiv="refresh" content="0;url=/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/error/3" />'

@app.route('/topic/<path:name>/sub/<path:sub>/notice/<int:number>')
def 토론_공지(name = None, sub = None, number = None):
    if(관리자_확인() == 1):
        DB_실행("select * from topic where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "' and id = '" + str(number) + "'")
        토론_내용 = DB_가져오기()
        if(토론_내용):
            DB_실행("select * from distop where id = '" + str(number) + "' and title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "'")
            공지_내용 = DB_가져오기()
            if(공지_내용):
                DB_실행("delete from distop where id = '" + str(number) + "' and title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "'")
            else:
                DB_실행("insert into distop (id, title, sub) value ('" + DB_인코딩(str(number)) + "', '" + DB_인코딩(name) + "', '" + DB_인코딩(sub) + "')")
            DB_갱신()
            
            최근_토론_추가(name, sub, 시간())

            return '<meta http-equiv="refresh" content="0;url=/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/error/3" />'
        
@app.route('/topic/<path:name>/sub/<path:sub>/stop')
def 토론_정지(name = None, sub = None):
    if(관리자_확인() == 1):
        아이피 = 아이피_확인(request)
        
        DB_실행("select * from topic where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "' limit 1")
        토론_확인 = DB_가져오기()
        if(토론_확인):
            현재_시간 = 시간()
            
            DB_실행("select * from stop where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "' and close = ''")
            정지 = DB_가져오기()
            if(정지):
                DB_실행("insert into topic (id, title, sub, data, date, ip, block) value ('" + DB_인코딩(str(int(토론_확인[0]['id']) + 1)) + "', '" + DB_인코딩(name) + "', '" + DB_인코딩(sub) + "', 'Restart', '" + DB_인코딩(현재_시간) + "', '" + DB_인코딩(아이피) + " - Restart', '')")
                DB_실행("delete from stop where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "' and close = ''")
            else:
                DB_실행("insert into topic (id, title, sub, data, date, ip, block) value ('" + DB_인코딩(str(int(토론_확인[0]['id']) + 1)) + "', '" + DB_인코딩(name) + "', '" + DB_인코딩(sub) + "', 'Stop', '" + DB_인코딩(현재_시간) + "', '" + DB_인코딩(아이피) + " - Stop', '')")
                DB_실행("insert into stop (title, sub, close) value ('" + DB_인코딩(name) + "', '" + DB_인코딩(sub) + "', '')")
            DB_갱신()
            
            최근_토론_추가(name, sub, 현재_시간)
            
            return '<meta http-equiv="refresh" content="0;url=/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/error/3" />'                
        
@app.route('/topic/<path:name>/sub/<path:sub>/close')
def 토론_닫기(name = None, sub = None):
    if(관리자_확인() == 1):
        아이피 = 아이피_확인(request)
        
        DB_실행("select * from topic where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "' order by id+0 desc limit 1")
        토론_확인 = DB_가져오기()
        if(토론_확인):
            현재_시간 = 시간()
            
            DB_실행("select * from stop where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "' and close = 'O'")
            닫기 = DB_가져오기()
            if(닫기):
                DB_실행("insert into topic (id, title, sub, data, date, ip, block) value ('" + DB_인코딩(str(int(토론_확인[0]['id']) + 1)) + "', '" + DB_인코딩(name) + "', '" + DB_인코딩(sub) + "', 'Reopen', '" + DB_인코딩(현재_시간) + "', '" + DB_인코딩(아이피) + " - Reopen', '')")
                DB_실행("delete from stop where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "' and close = 'O'")
            else:
                DB_실행("insert into topic (id, title, sub, data, date, ip, block) value ('" + DB_인코딩(str(int(토론_확인[0]['id']) + 1)) + "', '" + DB_인코딩(name) + "', '" + DB_인코딩(sub) + "', 'Close', '" + DB_인코딩(현재_시간) + "', '" + DB_인코딩(아이피) + " - Close', '')")
                DB_실행("insert into stop (title, sub, close) value ('" + DB_인코딩(name) + "', '" + DB_인코딩(sub) + "', 'O')")
            DB_갱신()
            
            최근_토론_추가(name, sub, 현재_시간)
            
            return '<meta http-equiv="refresh" content="0;url=/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/error/3" />'

@app.route('/topic/<path:name>/sub/<path:sub>/agree')
def 토론_관리자_기능(name = None, sub = None):
    if(관리자_확인() == 1):
        아이피 = 아이피_확인(request)
        
        DB_실행("select id from topic where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "' order by id+0 desc limit 1")
        토론 = DB_가져오기()
        if(토론):
            현재_시간 = 시간()
            
            DB_실행("select * from agreedis where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "'")
            합의안 = DB_가져오기()
            if(합의안):
                DB_실행("insert into topic (id, title, sub, data, date, ip, block) value ('" + DB_인코딩(str(int(토론[0]['id']) + 1)) + "', '" + DB_인코딩(name) + "', '" + DB_인코딩(sub) + "', 'Settlement', '" + DB_인코딩(현재_시간) + "', '" + DB_인코딩(아이피) + " - Settlement', '')")
                DB_실행("delete from agreedis where title = '" + DB_인코딩(name) + "' and sub = '" + DB_인코딩(sub) + "'")
            else:
                DB_실행("insert into topic (id, title, sub, data, date, ip, block) value ('" + DB_인코딩(str(int(토론[0]['id']) + 1)) + "', '" + DB_인코딩(name) + "', '" + DB_인코딩(sub) + "', 'Agreement', '" + DB_인코딩(현재_시간) + "', '" + DB_인코딩(아이피) + " - Agreement', '')")
                DB_실행("insert into agreedis (title, sub) value ('" + DB_인코딩(name) + "', '" + DB_인코딩(sub) + "')")
            DB_갱신()
            
            최근_토론_추가(name, sub, 시간())
            
            return '<meta http-equiv="refresh" content="0;url=/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/topic/' + URL_인코딩(name) + '/sub/' + URL_인코딩(sub) + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/error/3" />'

@app.route('/login', methods=['POST', 'GET'])
def 로그인():
    아이피 = 아이피_확인(request)
    차단인가 = 차단_체크(아이피)
        
    if(request.method == 'POST'):        
        if(차단인가 == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            DB_실행("select * from user where id = '" + DB_인코딩(request.form["id"]) + "'")
            사용자_정보 = DB_가져오기()
            if(사용자_정보):
                if(session.get('Now') == True):
                    return '<meta http-equiv="refresh" content="0;url=/error/11" />'
                elif(bcrypt.checkpw(bytes(request.form["pw"], 'utf-8'), bytes(사용자_정보[0]['pw'], 'utf-8'))):
                    session['Now'] = True
                    session['DREAMER'] = request.form["id"]
                    
                    DB_실행("insert into login (user, ip, today) value ('" + DB_인코딩(request.form["id"]) + "', '" + DB_인코딩(아이피) + "', '" + DB_인코딩(시간()) + "')")
                    DB_갱신()
                    
                    return '<meta http-equiv="refresh" content="0;url=/user" />'
                else:
                    return '<meta http-equiv="refresh" content="0;url=/error/13" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/12" />'
    else:        
        if(차단인가 == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            if(session.get('Now') == True):
                return '<meta http-equiv="refresh" content="0;url=/error/11" />'
            else:
                return 웹_디자인('index.html', title = '로그인', enter = '로그인', logo = data['name'], tn = 15)
                
@app.route('/change', methods=['POST', 'GET'])
def 비밀번호_변경():
    아이피 = 아이피_확인(request)
    차단인가 = 차단_체크(아이피)
    
    if(request.method == 'POST'):      
        if(request.form["pw2"] == request.form["pw3"]):
            if(차단인가 == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                DB_실행("select * from user where id = '" + DB_인코딩(request.form["id"]) + "'")
                사용자_정보 = DB_가져오기()
                if(사용자_정보):
                    if(session.get('Now') == True):
                        session['Now'] = False
                        session.pop('DREAMER', None)
                        return '<meta http-equiv="refresh" content="0;url=/change" />'
                    elif(bcrypt.checkpw(bytes(request.form["pw"], 'utf-8'), bytes(사용자_정보[0]['pw'], 'utf-8'))):
                        hashed = bcrypt.hashpw(bytes(request.form["pw2"], 'utf-8'), bcrypt.gensalt())
                        
                        DB_실행("update user set pw = '" + DB_인코딩(hashed.decode()) + "' where id = '" + DB_인코딩(request.form["id"]) + "'")
                        DB_갱신()
                        
                        return '<meta http-equiv="refresh" content="0;url=/login" />'
                    else:
                        return '<meta http-equiv="refresh" content="0;url=/error/10" />'
                else:
                    return '<meta http-equiv="refresh" content="0;url=/error/9" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/20" />'
    else:        
        if(차단인가 == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            if(session.get('Now') == True):
                session['Now'] = False
                session.pop('DREAMER', None)
                return '<meta http-equiv="refresh" content="0;url=/change" />'
            else:
                return 웹_디자인('index.html', title = '비밀번호 변경', enter = '변경', logo = data['name'], tn = 15)
                
@app.route('/check/<name>')
def 사용자_아이피_확인(name = None, sub = None, number = None):
    DB_실행("select * from user where id = '" + DB_인코딩(name) + "'")
    사용자_정보 = DB_가져오기()
    if(사용자_정보 and 사용자_정보[0]['acl'] == 'owner' or 사용자_정보 and 사용자_정보[0]['acl'] == 'admin'):
        return '<meta http-equiv="refresh" content="0;url=/error/4" />'
    else:
        if(관리자_확인() == 1):
            m = re.search('(?:[0-9](?:[0-9][0-9])?\.[0-9](?:[0-9][0-9])?\.[0-9](?:[0-9][0-9])?\.[0-9](?:[0-9][0-9])?)', name)
            if(m):
                DB_실행("select * from login where ip = '" + DB_인코딩(name) + "' order by today desc")
                row = DB_가져오기()
                if(row):
                    i = 0
                    c = ''
                    while(True):
                        try:
                            c = c + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:33.33%;">' + row[i]['user'] + '</td><td style="text-align: center;width:33.33%;">' + row[i]['ip'] + '</td><td style="text-align: center;width:33.33%;">' + row[i]['today'] + '</td></tr></tbody></table>'
                        except:
                            break
                        i = i + 1
                    return 웹_디자인('index.html', title = '다중 검사', logo = data['name'], tn = 22, rows = c)
                else:
                    return 웹_디자인('index.html', title = '다중 검사', logo = data['name'], tn = 22, rows = '')
            else:
                DB_실행("select * from login where user = '" + DB_인코딩(name) + "' order by today desc")
                row = DB_가져오기()
                if(row):
                    i = 0
                    c = ''
                    while(True):
                        try:
                            c = c + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:33.33%;">' + row[i]['user'] + '</td><td style="text-align: center;width:33.33%;">' + row[i]['ip'] + '</td><td style="text-align: center;width:33.33%;">' + row[i]['today'] + '</td></tr></tbody></table>'
                        except:
                            break
                        i = i + 1
                    return 웹_디자인('index.html', title = '다중 검사', logo = data['name'], tn = 22, rows = c)
                else:
                    return 웹_디자인('index.html', title = '다중 검사', logo = data['name'], tn = 22, rows = '')
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'

@app.route('/register', methods=['POST', 'GET'])
def 가입():
    아이피 = 아이피_확인(request)
    차단인가 = 차단_체크(아이피)
    
    if(request.method == 'POST'):        
        if(request.form["pw"] == request.form["pw2"]):
            if(차단인가 == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                m = re.search('(?:[^A-Za-zㄱ-힣0-9 ])', request.form["id"])
                if(m):
                    return '<meta http-equiv="refresh" content="0;url=/error/8" />'
                else:
                    if(len(request.form["id"]) > 20):
                        return '<meta http-equiv="refresh" content="0;url=/error/7" />'
                    else:
                        DB_실행("select * from user where id = '" + DB_인코딩(request.form["id"]) + "'")
                        rows = DB_가져오기()
                        if(rows):
                            return '<meta http-equiv="refresh" content="0;url=/error/6" />'
                        else:
                            hashed = bcrypt.hashpw(bytes(request.form["pw"], 'utf-8'), bcrypt.gensalt())
                            if(request.form["id"] == data['owner']):
                                DB_실행("insert into user (id, pw, acl) value ('" + DB_인코딩(request.form["id"]) + "', '" + DB_인코딩(hashed.decode()) + "', 'owner')")
                            else:
                                DB_실행("insert into user (id, pw, acl) value ('" + DB_인코딩(request.form["id"]) + "', '" + DB_인코딩(hashed.decode()) + "', 'user')")
                            DB_갱신()
                            return '<meta http-equiv="refresh" content="0;url=/login" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/20" />'
    else:        
        if(차단인가 == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            return 웹_디자인('index.html', title = '회원가입', enter = '회원가입', logo = data['name'], tn = 15)

@app.route('/logout')
def 로그아웃():
    session['Now'] = False
    session.pop('DREAMER', None)
    return '<meta http-equiv="refresh" content="0;url=/user" />'

@app.route('/ban/<name>', methods=['POST', 'GET'])
def 사용자_차단(name = None):
    DB_실행("select * from user where id = '" + DB_인코딩(name) + "'")
    rows = DB_가져오기()
    if(rows and rows[0]['acl'] == 'owner' or rows and rows[0]['acl'] == 'admin'):
        return '<meta http-equiv="refresh" content="0;url=/error/4" />'
    else:
        if(request.method == 'POST'):
            if(관리자_확인() == 1):
                ip = 아이피_확인(request)
                
                if(not re.search("[0-9]{4}-[0-9]{2}-[0-9]{2}", request.form["end"])):
                    end = ''
                else:
                    end = request.form["end"]

                DB_실행("select * from ban where block = '" + DB_인코딩(name) + "'")
                row = DB_가져오기()
                if(row):
                    최근_차단_추가(name, '해제', 시간(), ip, '')
                    
                    DB_실행("delete from ban where block = '" + DB_인코딩(name) + "'")
                else:
                    b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))$", name)
                    if(b):
                        최근_차단_추가(name, end, 시간(), ip, request.form["why"])
                        
                        DB_실행("insert into ban (block, end, why, band) value ('" + DB_인코딩(name) + "', '" + DB_인코딩(end) + "', '" + DB_인코딩(request.form["why"]) + "', 'O')")
                    else:
                        최근_차단_추가(name, end, 시간(), ip, request.form["why"])
                        
                        DB_실행("insert into ban (block, end, why, band) value ('" + DB_인코딩(name) + "', '" + DB_인코딩(end) + "', '" + DB_인코딩(request.form["why"]) + "', '')")
                DB_갱신()
                
                return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(data['frontpage']) + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/3" />'
        else:
            if(관리자_확인() == 1):
                DB_실행("select * from ban where block = '" + DB_인코딩(name) + "'")
                row = DB_가져오기()
                if(row):
                    now = '차단 해제'
                else:
                    b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))$", name)
                    if(b):
                        now = '대역 차단'
                    else:
                        now = '차단'
                        
                return 웹_디자인('index.html', title = name, page = URL_인코딩(name), logo = data['name'], tn = 16, now = now, today = 시간(), sub = '차단')
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/3" />'

@app.route('/acl/<path:name>', methods=['POST', 'GET'])
def ACL(name = None):
    if(request.method == 'POST'):
        if(관리자_확인() == 1):
            DB_실행("select * from data where title = '" + DB_인코딩(name) + "'")
            row = DB_가져오기()
            if(row):
                if(request.form["select"] == 'admin'):
                   DB_실행("update data set acl = 'admin' where title = '" + DB_인코딩(name) + "'")
                elif(request.form["select"] == 'user'):
                    DB_실행("update data set acl = 'user' where title = '" + DB_인코딩(name) + "'")
                else:
                    DB_실행("update data set acl = '' where title = '" + DB_인코딩(name) + "'")
                DB_갱신()
            return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(name) + '" />' 
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'
    else:
        if(관리자_확인() == 1):
            DB_실행("select * from data where title = '" + DB_인코딩(name) + "'")
            row = DB_가져오기()
            if(row):
                if(row[0]['acl'] == 'admin'):
                    now = '관리자만'
                elif(row[0]['acl'] == 'user'):
                    now = '유저 이상'
                else:
                    now = '일반'
                return 웹_디자인('index.html', title = name, page = URL_인코딩(name), logo = data['name'], tn = 19, now = '현재 ACL 상태는 ' + now, sub = 'ACL')
            else:
                return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(name) + '" />' 
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'

@app.route('/admin/<name>', methods=['POST', 'GET'])
def 관리자_부여(name = None):
    if(request.method == 'POST'):
        if(소유자_확인() == 1):
            DB_실행("select * from user where id = '" + DB_인코딩(name) + "'")
            사용자_정보 = DB_가져오기()
            if(사용자_정보):
                if(사용자_정보[0]['acl'] == 'admin' or 사용자_정보[0]['acl'] == 'owner'):
                    DB_실행("update user set acl = 'user' where id = '" + DB_인코딩(name) + "'")
                else:
                    DB_실행("update user set acl = '" + DB_인코딩(request.form["select"]) + "' where id = '" + DB_인코딩(name) + "'")
                DB_갱신()
                
                return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(data['frontpage']) + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/5" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'
    else:
        if(소유자_확인() == 1):
            DB_실행("select * from user where id = '" + DB_인코딩(name) + "'")
            사용자_정보 = DB_가져오기()
            if(사용자_정보):
                if(사용자_정보[0]['acl'] == 'admin' or 사용자_정보[0]['acl'] == 'owner'):
                    now = '권한 해제'
                else:
                    now = '권한 부여'
                return 웹_디자인('index.html', title = name, page = URL_인코딩(name), logo = data['name'], tn = 18, now = now, sub = '권한 부여')
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/5" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'

@app.route('/ban')
def 차단_확인_페이지():
    ip = 아이피_확인(request)
    
    if(차단_체크(ip) == 1):
        DB_실행("select * from ban where block = '" + DB_인코딩(ip) + "'")
        rows = DB_가져오기()
        if(rows):
            if(rows[0]['end']):
                end = rows[0]['end'] + ' 까지 차단 상태 입니다. / 사유 : ' + rows[0]['why']                
                now = 시간()
                
                now = re.sub(':', '', now)
                now = re.sub('\-', '', now)
                now = re.sub(' ', '', now)
                now = int(now)
                
                day = rows[0]['end']
                day = re.sub('\-', '', day)    
                
                if(now >= int(day + '000000')):
                    DB_실행("delete from ban where block = '" + DB_인코딩(ip) + "'")
                    DB_갱신()
                    end = '차단이 풀렸습니다. 다시 시도 해 보세요.'
            else:
                end = '영구 차단 상태 입니다. / 사유 : ' + rows[0]['why']
        else:
            b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
            if(b):
                results = b.groups()
                
                DB_실행("select * from ban where block = '" + DB_인코딩(results[0]) + "' and band = 'O'")
                row = DB_가져오기()
                if(row):
                    if(row[0]['end']):
                        end = row[0]['end'] + ' 까지 차단 상태 입니다. / 사유 : ' + rows[0]['why']             
                        
                        now = 시간()
                        now = re.sub(':', '', now)
                        now = re.sub('\-', '', now)
                        now = re.sub(' ', '', now)
                        now = int(now)    
                        
                        day = row[0]['end']
                        day = re.sub('\-', '', day)
                        
                        if(now >= int(day + '000000')):
                            DB_실행("delete from ban where block = '" + DB_인코딩(results[0]) + "' and band = 'O'")
                            DB_갱신()
                            end = '차단이 풀렸습니다. 다시 시도 해 보세요.'
                    else:
                        end = '영구 차단 상태 입니다. / 사유 : ' + row[0]['why']                
    else:
        end = '권한이 맞지 않는 상태 입니다.'
    return 웹_디자인('index.html', title = '권한 오류', logo = data['name'], data = end), 401
   
@app.route('/w/<path:name>/r/<int:a>/diff/<int:b>')
def 문서_비교(name = None, a = None, b = None):
    DB_실행("select * from history where id = '" + DB_인코딩(str(a)) + "' and title = '" + DB_인코딩(name) + "'")
    rows = DB_가져오기()
    if(rows):
        DB_실행("select * from history where id = '" + DB_인코딩(str(b)) + "' and title = '" + DB_인코딩(name) + "'")
        row = DB_가져오기()
        if(row):
            indata = re.sub('<', '&lt;', rows[0]['data'])
            indata = re.sub('>', '&gt;', indata)
            indata = re.sub('"', '&quot;', indata)
            
            enddata = re.sub('<', '&lt;', row[0]['data'])
            enddata = re.sub('>', '&gt;', enddata)
            enddata = re.sub('"', '&quot;', enddata)
            
            sm = difflib.SequenceMatcher(None, indata, enddata)
            c = 비교(sm)
            
            c = '<pre>' + c + '</pre>'
            
            return 웹_디자인('index.html', title = name, logo = data['name'], data = c, sub = '비교')
        else:
            return '<meta http-equiv="refresh" content="0;url=/history/' + URL_인코딩(name) + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/history/' + URL_인코딩(name) + '" />'
        
@app.route('/user')
def 사용자():
    ip = 아이피_확인(request)
    
    DB_실행("select * from user where id = '" + DB_인코딩(ip) + "'")
    rows = DB_가져오기()
    if(차단_체크(ip) == 0):
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
        
    if(not re.search('\.', ip)):
        DB_실행("select * from data where title = '사용자:" + DB_인코딩(ip) + "'")
        row = DB_가져오기()
        if(row):
            ip = '<a href="/w/' + URL_인코딩('사용자:' + ip) + '">' + ip + '</a>'
        else:
            ip = '<a class="not_thing" href="/w/' + URL_인코딩('사용자:' + ip) + '">' + ip + '</a>'
        
    return 웹_디자인('index.html', title = '유저 메뉴', logo = data['name'], data = ip + '<br><br><span>권한 상태 : ' + acl + '<br><br><li><a href="/login">로그인</a></li><li><a href="/logout">로그아웃</a></li><li><a href="/register">회원가입</a></li><li><a href="/change">비밀번호 변경</a></li>')

@app.route('/random')
def 무작위_문서():
    DB_실행("select * from data order by rand() limit 1")
    rows = DB_가져오기()
    if(rows):
        return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(rows[0]['title']) + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/" />'
        
@app.route('/error/<int:num>')
def 오류(num = None):
    if(num == 1):
        return 웹_디자인('index.html', title = '권한 오류', logo = data['name'], data = '비 로그인 상태 입니다.'), 401
    elif(num == 2):
        return 웹_디자인('index.html', title = '권한 오류', logo = data['name'], data = '이 계정이 없습니다.'), 401
    elif(num == 3):
        return 웹_디자인('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.'), 401
    elif(num == 4):
        return 웹_디자인('index.html', title = '권한 오류', logo = data['name'], data = '관리자는 차단, 검사 할 수 없습니다.'), 401
    elif(num == 5):
        return 웹_디자인('index.html', title = '유저 오류', logo = data['name'], data = '그런 계정이 없습니다.'), 401
    elif(num == 6):
        return 웹_디자인('index.html', title = '가입 오류', logo = data['name'], data = '동일한 아이디의 유저가 있습니다.'), 401
    elif(num == 7):
        return 웹_디자인('index.html', title = '가입 오류', logo = data['name'], data = '아이디는 20글자보다 짧아야 합니다.'), 401
    elif(num == 8):
        return 웹_디자인('index.html', title = '가입 오류', logo = data['name'], data = '아이디에는 한글과 알파벳과 공백만 허용 됩니다.'), 401
    elif(num == 9):
        return 웹_디자인('index.html', title = '변경 오류', logo = data['name'], data = '그런 계정이 없습니다.'), 401
    elif(num == 10):
        return 웹_디자인('index.html', title = '변경 오류', logo = data['name'], data = '비밀번호가 다릅니다.'), 401
    elif(num == 11):
        return 웹_디자인('index.html', title = '로그인 오류', logo = data['name'], data = '이미 로그인 되어 있습니다.'), 401
    elif(num == 12):
        return 웹_디자인('index.html', title = '로그인 오류', logo = data['name'], data = '그런 계정이 없습니다.'), 401
    elif(num == 13):
        return 웹_디자인('index.html', title = '로그인 오류', logo = data['name'], data = '비밀번호가 다릅니다.'), 401
    elif(num == 14):
        return 웹_디자인('index.html', title = '업로드 오류', logo = data['name'], data = 'jpg, gif, jpeg, png만 가능 합니다.'), 401
    elif(num == 15):
        return 웹_디자인('index.html', title = '업로드 오류', logo = data['name'], data = '파일 명에 . / \ * < > | : ? 가 들어 갈 수 없습니다.'), 401
    elif(num == 16):
        return 웹_디자인('index.html', title = '업로드 오류', logo = data['name'], data = '동일한 이름의 파일이 있습니다.'), 401
    elif(num == 17):
        return 웹_디자인('index.html', title = '편집 오류', logo = data['name'], data = '편집 내용 기록에는 한글과 영어와 숫자, 공백만 허용 됩니다.'), 401
    elif(num == 18):
        return 웹_디자인('index.html', title = '편집 오류', logo = data['name'], data = '내용이 원래 문서와 동일 합니다.'), 401
    elif(num == 19):
        return 웹_디자인('index.html', title = '이동 오류', logo = data['name'], data = '이동 하려는 곳에 문서가 이미 있습니다.'), 401
    elif(num == 20):
        return 웹_디자인('index.html', title = '비밀번호 오류', logo = data['name'], data = '재 확인이랑 비밀번호가 다릅니다.'), 401
    else:
        return '<meta http-equiv="refresh" content="0;url=/" />'

@app.errorhandler(404)
def uncaughtError(error):
    return '<meta http-equiv="refresh" content="0;url=/w/' + URL_인코딩(data['frontpage']) + '" />'

@app.errorhandler(413)
def uncaughtError(error):
    app.config['MAX_CONTENT_LENGTH'] = (1024**3)
    return error, 401

if(__name__ == '__main__'):
    app.run(host = '0.0.0.0', port = int(data['port']))