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

def start():
    try:
        curs.execute("select * from data limit 1")
    except:
        curs.execute("create table data(title text, data longtext, acl text)")
    
    try:
        curs.execute("select * from history limit 1")
    except:
        curs.execute("create table history(id text, title text, data longtext, date text, ip text, send text, leng text)")
    
    try:
        curs.execute("select * from rd limit 1")
    except:
        curs.execute("create table rd(title text, sub text, date text)")
    
    try:
        curs.execute("select * from user limit 1")
    except:
        curs.execute("create table user(id text, pw text, acl text)")
    
    try:
        curs.execute("select * from ban limit 1")
    except:
        curs.execute("create table ban(block text, end text, why text, band text)")
    
    try:
        curs.execute("select * from topic limit 1")
    except:
        curs.execute("create table topic(id text, title text, sub text, data longtext, date text, ip text, block text)")
    
    try:
        curs.execute("select * from stop limit 1")
    except:
        curs.execute("create table stop(title text, sub text, close text)")
    
    try:
        curs.execute("select * from rb limit 1")
    except:
        curs.execute("create table rb(block text, end text, today text, blocker text, why text)")
    
    try:
        curs.execute("select * from login limit 1")
    except:
        curs.execute("create table login(user text, ip text, today text)")
    
    try:
        curs.execute("select * from back limit 1")
    except:
        curs.execute("create table back(title text, link text, type text)")
    
    try:
        curs.execute("select * from cat limit 1")
    except:
        curs.execute("create table cat(title text, cat text)")
        
    try:
        curs.execute("select * from hidhi limit 1")
    except:
        curs.execute("create table hidhi(title text, re text)")    

try:
    conn = pymysql.connect(host = data['host'], user = data['user'], password = data['pw'], db = data['db'], charset = 'utf8mb4')
    curs = conn.cursor(pymysql.cursors.DictCursor)
    
    start()
except:
    conn = pymysql.connect(host = data['host'], user = data['user'], password = data['pw'], charset = 'utf8mb4')
    curs = conn.cursor(pymysql.cursors.DictCursor)
    
    curs.execute("create database " + data['db'])
    curs.execute("alter database " + data['db'] + " character set = utf8mb4 collate = utf8mb4_unicode_ci")
    curs.execute("use " + data['db'])    
    
    start()

app.secret_key = hashlib.sha512(bytes(data['key'], 'ascii')).hexdigest()

def show_diff(seqm):
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
           
def admincheck():
    if(session.get('Now') == True):
        ip = getip(request)
        curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
        rows = curs.fetchall()
        if(rows):
            if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
                return 1
                
def ownercheck():
    if(session.get('Now') == True):
        ip = getip(request)
        curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
        rows = curs.fetchall()
        if(rows):
            if(rows[0]['acl'] == 'owner'):
                return 1
                
def isin(name):
    if(re.search('^틀:', name)):
        curs.execute("select * from back where title = '" + pymysql.escape_string(name) + "' and type = 'include'")
        include = curs.fetchall()
        if(include):
            i = 0
            
            while(True):
                try:
                    curs.execute("select * from data where title = '" + pymysql.escape_string(include[i]['link']) + "'")
                    file = curs.fetchall()
                    if(file):
                        namumark(include[i]['link'], file[0]['data'])
                except:
                    break
                    
                i = i + 1

def namumark(title, data):
    while(True):
        m = re.search("<((div|span|embed|iframe)(?:[^>]*))>", data)
        if(m):
            g = m.groups()

            if(re.search("<(\/" + g[1] + ")>", data)):
                n = re.search('src="http(?:s)?:\/\/([^\/]*)\/(?:[^"]*)"', g[0])
                if(n):
                    gm = n.groups()
                    if(gm[0] == "www.youtube.com" or gm[0] == "serviceapi.nmv.naver.com" or gm[0] == "tv.kakao.com" or gm[0] == "tvple.com" or gm[0] == "tvpot.daum.net"):
                        r = g[0]
                    else:
                        r = re.sub('src="([^"]*)"', '', g[0])
                else:
                    r = g[0]
                
                r = re.sub('"', '#.#', r)
                data = re.sub("<((?:\/)?" + g[1] + "(?:[^>]*))>", "[" + r + "]", data, 2)
            else:
                break
        else:
            break

    data = re.sub('<', '&lt;', data)
    data = re.sub('>', '&gt;', data)
    data = re.sub('"', '&quot;', data)
    
    data = re.sub("\[(?P<in>(?:\/)?(?:div|span|embed|iframe)(?:[^\]]*))\]", "<\g<in>>", data)
    data = re.sub('#.#', '"', data)
    
    jjjj = 0
    while(True):
        p = re.compile("{{{((?:(?!{{{)(?!}}}).)*)}}}", re.DOTALL)
        m = p.search(data)
        if(m):
            results = m.groups()
            q = re.compile("^\+([1-5])\s(.*)$", re.DOTALL)
            n = q.search(results[0])
            
            w = re.compile("^\-([1-5])\s(.*)$", re.DOTALL)
            a = w.search(results[0])
            
            e = re.compile("^(#[0-9a-f-A-F]{6})\s(.*)$", re.DOTALL)
            b = e.search(results[0])
            
            r = re.compile("^(#[0-9a-f-A-F]{3})\s(.*)$", re.DOTALL)
            c = r.search(results[0])
            
            t = re.compile("^#(\w+)\s(.*)$", re.DOTALL)
            d = t.search(results[0])
            
            qqq = re.compile("^@([0-9a-f-A-F]{6})\s(.*)$", re.DOTALL)
            qqe = qqq.search(results[0])
            
            qqw = re.compile("^@([0-9a-f-A-F]{3})\s(.*)$", re.DOTALL)
            qqa = qqw.search(results[0])
            
            qwe = re.compile("^@(\w+)\s(.*)$", re.DOTALL)
            qsd = qwe.search(results[0])
            
            qawe = re.compile("^#!noin\s(.*)$", re.DOTALL)
            qasd = qawe.search(results[0])
            
            y = re.compile("^#!wiki\sstyle=&quot;((?:(?!&quot;|\n).)*)&quot;\n?\s\n(.*)$", re.DOTALL)
            l = y.search(results[0])
            
            html = re.compile("^#!html\s(.*)$", re.DOTALL)
            htmlr = html.search(results[0])
            
            ppp = re.compile("^#!folding\s((?:(?!\n).)*)\n?\s\n(.*)$", re.DOTALL)
            ooo = ppp.search(results[0])
            
            if(n):
                result = n.groups()
                data = p.sub('<span class="font-size-' + result[0] + '">' + result[1] + '</span>', data, 1)
            elif(a):
                result = a.groups()
                data = p.sub('<span class="font-size-small-' + result[0] + '">' + result[1] + '</span>', data, 1)
            elif(b):
                result = b.groups()
                data = p.sub('<span style="color:' + result[0] + '">' + result[1] + '</span>', data, 1)
            elif(c):
                result = c.groups()
                data = p.sub('<span style="color:' + result[0] + '">' + result[1] + '</span>', data, 1)
            elif(d):
                result = d.groups()
                data = p.sub('<span style="color:' + result[0] + '">' + result[1] + '</span>', data, 1)
            elif(qqe):
                result = qqe.groups()
                data = p.sub('<span style="background:#' + result[0] + '">' + result[1] + '</span>', data, 1)
            elif(qqa):
                result = qqa.groups()
                data = p.sub('<span style="background:#' + result[0] + '">' + result[1] + '</span>', data, 1)
            elif(qsd):
                result = qsd.groups()
                data = p.sub('<span style="background:' + result[0] + '">' + result[1] + '</span>', data, 1)
            elif(l):
                result = l.groups()
                data = p.sub('<div style="' + result[0] + '">' + result[1] + '</div>', data, 1)
            elif(htmlr):
                result = htmlr.groups()
                data = p.sub(result[0], data, 1)
            elif(ooo):
                result = ooo.groups()
                data = p.sub("<div>" + result[0] + "<span style='float:right;'><div id='folding_" + str(jjjj + 1) + "' style='display:block;'>[<a href='javascript:void(0);' onclick='var f=document.getElementById(\"folding_" + str(jjjj) + "\");var s=f.style.display==\"block\";f.style.display=s?\"none\":\"block\";this.className=s?\"\":\"opened\";var f=document.getElementById(\"folding_" + str(jjjj + 1) + "\");var s=f.style.display==\"none\";f.style.display=s?\"block\":\"none\";var f=document.getElementById(\"folding_" + str(jjjj + 2) + "\");var s=f.style.display==\"block\";f.style.display=s?\"none\":\"block\";'>펼치기</a>]</div><div id='folding_" + str(jjjj + 2) + "' style='display:none;'>[<a href='javascript:void(0);' onclick='var f=document.getElementById(\"folding_" + str(jjjj) + "\");var s=f.style.display==\"block\";f.style.display=s?\"none\":\"block\";this.className=s?\"\":\"opened\";var f=document.getElementById(\"folding_" + str(jjjj + 1) + "\");var s=f.style.display==\"none\";f.style.display=s?\"block\":\"none\";var f=document.getElementById(\"folding_" + str(jjjj + 2) + "\");var s=f.style.display==\"block\";f.style.display=s?\"none\":\"block\";'>접기</a>]</div></a></span><div id='folding_" + str(jjjj) + "' style='display:none;'><br>" + result[1] + "</div></div>", data, 1)
                jjjj = jjjj + 3
            elif(qasd):
                result = qasd.groups()
                data = p.sub(result[0], data, 1)
            else:
                data = p.sub('<code>' + results[0] + '</code>', data, 1)
        else:
            break
    
    while(True):
        a = re.compile("<code>(((?!<\/code>).)*)<\/code>", re.DOTALL)
        m = a.search(data)
        if(m):
            g = m.groups()
            j = re.sub("<\/span>", "}}}", g[0])
            j = re.sub("<\/div>", "}}}", j)
            j = re.sub('<span class="font\-size\-(?P<in>[1-6])">', "{{{+\g<in> ", j)
            j = re.sub('<span class="font\-size\-small\-(?P<in>[1-6])">', "{{{-\g<in> ", j)
            j = re.sub('<span style="color:(?:#)?(?P<in>[^"]*)">', "{{{#\g<in> ", j)
            j = re.sub('<span style="background:(?:#)?(?P<in>[^"]*)">', "{{{@\g<in> ", j)
            j = re.sub('<div style="(?P<in>[^"]*)">', "{{{#!wiki style=&quot;\g<in>&quot;\n", j)
            j = re.sub("(?P<in>.)", "<span>\g<in></span>", j)
            data = a.sub(j, data, 1)
        else:
            break
            
    data = re.sub("<span>&</span><span>l</span><span>t</span><span>;</span>", "<span>&lt;</span>", data)
    data = re.sub("<span>&</span><span>g</span><span>t</span><span>;</span>", "<span>&gt;</span>", data)
    
    data = re.sub("\[anchor\((?P<in>[^\[\]]*)\)\]", '<span id="\g<in>"></span>', data)
    data = re.sub('\[date\(now\)\]', getnow(), data)
    
    while(True):
        m = re.search("\[include\(((?:(?!\)\]|,).)*)((?:,\s?(?:[^)]*))+)?\)\]", data)
        if(m):
            results = m.groups()
            if(results[0] == title):
                data = re.sub("\[include\(((?:(?!\)\]|,).)*)((?:,\s?(?:[^)]*))+)?\)\]", "<b>" + results[0] + "</b>", data, 1)
            else:
                curs.execute("select * from data where title = '" + pymysql.escape_string(results[0]) + "'")
                rows = curs.fetchall()
                if(rows):
                    curs.execute("select * from back where title = '" + pymysql.escape_string(results[0]) + "' and link = '" + pymysql.escape_string(title) + "' and type = 'include'")
                    abb = curs.fetchall()
                    if(not abb):
                        curs.execute("insert into back (title, link, type) value ('" + pymysql.escape_string(results[0]) + "', '" + pymysql.escape_string(title) + "',  'include')")
                        conn.commit()
                        
                    enddata = rows[0]['data']
                    enddata = re.sub("\[include\(((?:(?!\)\]|,).)*)((?:,\s?(?:[^)]*))+)?\)\]", "", enddata)
                    
                    while(True):
                        m = re.search("<((?:\/)?(?:div|span|embed|iframe)(?:[^>]*))>", enddata)
                        if(m):
                            g = m.groups()
                            
                            n = re.search('src="http(?:s)?:\/\/([^\/]*)\/(?:[^"]*)"', g[0])
                            if(n):
                                gm = n.groups()
                                if(gm[0] == "www.youtube.com" or gm[0] == "serviceapi.nmv.naver.com" or gm[0] == "tv.kakao.com" or gm[0] == "tvple.com" or gm[0] == "tvpot.daum.net"):
                                    r = g[0]
                                else:
                                    r = re.sub('src="([^"]*)"', '', g[0])
                            else:
                                r = g[0]
                            
                            r = re.sub('"', '#.#', r)
                            enddata = re.sub("<((?:\/)?(?:div|span|embed|iframe)(?:[^>]*))>", "[" + r + "]", enddata, 1)
                        else:
                            break

                    enddata = re.sub('<', '&lt;', enddata)
                    enddata = re.sub('>', '&gt;', enddata)
                    enddata = re.sub('"', '&quot;', enddata)
                    
                    enddata = re.sub("\[(?P<in>(?:\/)?(?:div|span|embed|iframe)(?:[^\]]*))\]", "<\g<in>>", enddata)
                    enddata = re.sub('#.#', '"', enddata)
                    
                    while(True):
                        p = re.compile("{{{((?:(?!{)(?!}).)*)}}}", re.DOTALL)
                        m = p.search(enddata)
                        if(m):
                            nnn = m.groups()
                            q = re.compile("^\+([1-5])\s(.*)$", re.DOTALL)
                            n = q.search(nnn[0])
                            
                            w = re.compile("^\-([1-5])\s(.*)$", re.DOTALL)
                            a = w.search(nnn[0])
                            
                            e = re.compile("^(#[0-9a-f-A-F]{6})\s(.*)$", re.DOTALL)
                            b = e.search(nnn[0])
                            
                            r = re.compile("^(#[0-9a-f-A-F]{3})\s(.*)$", re.DOTALL)
                            c = r.search(nnn[0])
                            
                            t = re.compile("^#(\w+)\s(.*)$", re.DOTALL)
                            d = t.search(nnn[0])
                            
                            qqq = re.compile("^@([0-9a-f-A-F]{6})\s(.*)$", re.DOTALL)
                            qqe = qqq.search(nnn[0])
                            
                            qqw = re.compile("^@([0-9a-f-A-F]{3})\s(.*)$", re.DOTALL)
                            qqa = qqw.search(nnn[0])
                            
                            qwe = re.compile("^@(\w+)\s(.*)$", re.DOTALL)
                            qsd = qwe.search(nnn[0])
                            
                            qawe = re.compile("^#!noin\s(.*)$", re.DOTALL)
                            qasd = qawe.search(nnn[0])
                            
                            y = re.compile("^#!wiki\sstyle=&quot;((?:(?!&quot;|\n).)*)&quot;\n?\s\n(.*)$", re.DOTALL)
                            l = y.search(nnn[0])
                            
                            html = re.compile("^#!html\s(.*)$", re.DOTALL)
                            htmlr = html.search(nnn[0])
                            
                            ppp = re.compile("^#!folding\s((?:(?!\n).)*)\n?\s\n(.*)$", re.DOTALL)
                            ooo = ppp.search(nnn[0])
                            
                            if(n):
                                result = n.groups()
                                enddata = p.sub('<span class="font-size-' + result[0] + '">' + result[1] + '</span>', enddata, 1)
                            elif(a):
                                result = a.groups()
                                enddata = p.sub('<span class="font-size-small-' + result[0] + '">' + result[1] + '</span>', enddata, 1)
                            elif(b):
                                result = b.groups()
                                enddata = p.sub('<span style="color:' + result[0] + '">' + result[1] + '</span>', enddata, 1)
                            elif(c):
                                result = c.groups()
                                enddata = p.sub('<span style="color:' + result[0] + '">' + result[1] + '</span>', enddata, 1)
                            elif(d):
                                result = d.groups()
                                enddata = p.sub('<span style="color:' + result[0] + '">' + result[1] + '</span>', enddata, 1)
                            elif(qqe):
                                result = qqe.groups()
                                enddata = p.sub('<span style="background:#' + result[0] + '">' + result[1] + '</span>', enddata, 1)
                            elif(qqa):
                                result = qqa.groups()
                                enddata = p.sub('<span style="background:#' + result[0] + '">' + result[1] + '</span>', enddata, 1)
                            elif(qsd):
                                result = qsd.groups()
                                enddata = p.sub('<span style="background:' + result[0] + '">' + result[1] + '</span>', enddata, 1)
                            elif(l):
                                result = l.groups()
                                enddata = p.sub('<div style="' + result[0] + '">' + result[1] + '</div>', enddata, 1)
                            elif(htmlr):
                                result = htmlr.groups()
                                data = p.sub(result[0], enddata, 1)
                            elif(ooo):
                                result = ooo.groups()
                                enddata = p.sub("<div>" + result[0] + "<span style='float:right;'><div id='folding_" + str(jjjj + 1) + "' style='display:block;'>[<a href='javascript:void(0);' onclick='var f=document.getElementById(\"folding_" + str(jjjj) + "\");var s=f.style.display==\"block\";f.style.display=s?\"none\":\"block\";this.className=s?\"\":\"opened\";var f=document.getElementById(\"folding_" + str(jjjj + 1) + "\");var s=f.style.display==\"none\";f.style.display=s?\"block\":\"none\";var f=document.getElementById(\"folding_" + str(jjjj + 2) + "\");var s=f.style.display==\"block\";f.style.display=s?\"none\":\"block\";'>펼치기</a>]</div><div id='folding_" + str(jjjj + 2) + "' style='display:none;'>[<a href='javascript:void(0);' onclick='var f=document.getElementById(\"folding_" + str(jjjj) + "\");var s=f.style.display==\"block\";f.style.display=s?\"none\":\"block\";this.className=s?\"\":\"opened\";var f=document.getElementById(\"folding_" + str(jjjj + 1) + "\");var s=f.style.display==\"none\";f.style.display=s?\"block\":\"none\";var f=document.getElementById(\"folding_" + str(jjjj + 2) + "\");var s=f.style.display==\"block\";f.style.display=s?\"none\":\"block\";'>접기</a>]</div></a></span><div id='folding_" + str(jjjj) + "' style='display:none;'><br>" + result[1] + "</div></div>", enddata, 1)
                                jjjj = jjjj + 3
                            elif(qasd):
                                enddata = p.sub("", enddata, 1)
                            else:
                                enddata = p.sub('<code>' + nnn[0] + '</code>', enddata, 1)
                        else:
                            break
                    
                    while(True):
                        a = re.compile("<code>(((?!<\/code>).)*)<\/code>", re.DOTALL)
                        m = a.search(enddata)
                        if(m):
                            g = m.groups()
                            j = re.sub("<\/span>", "}}}", g[0])
                            j = re.sub("<\/div>", "}}}", j)
                            j = re.sub('<span class="font\-size\-(?P<in>[1-6])">', "{{{+\g<in> ", j)
                            j = re.sub('<span class="font\-size\-small\-(?P<in>[1-6])">', "{{{-\g<in> ", j)
                            j = re.sub('<span style="color:(?:#)?(?P<in>[^"]*)">', "{{{#\g<in> ", j)
                            j = re.sub('<span style="background:(?:#)?(?P<in>[^"]*)">', "{{{@\g<in> ", j)
                            j = re.sub('<div style="(?P<in>[^"]*)">', "{{{#!wiki style=&quot;\g<in>&quot;\n", j)
                            j = re.sub("(?P<in>.)", "<span>\g<in></span>", j)
                            enddata = a.sub(j, enddata, 1)
                        else:
                            break
                    
                    enddata = re.sub("<span>&</span><span>l</span><span>t</span><span>;</span>", "<span>&lt;</span>", enddata)
                    enddata = re.sub("<span>&</span><span>g</span><span>t</span><span>;</span>", "<span>&gt;</span>", enddata)
                    
                    if(results[1]):
                        a = results[1]
                        while(True):
                            g = re.search("([^= ,]*)\=([^,]*)", a)
                            if(g):
                                result = g.groups()
                                enddata = re.sub("@" + result[0] + "@", result[1], enddata)
                                a = re.sub("([^= ,]*)\=([^,]*)", "", a, 1)
                            else:
                                break                        
                    data = re.sub("\[include\(((?:(?!\)\]|,).)*)((?:,\s?(?:[^)]*))+)?\)\]", '\n<div>' + enddata + '</div>\n', data, 1)
                else:
                    curs.execute("select * from back where title = '" + pymysql.escape_string(results[0]) + "' and link = '" + pymysql.escape_string(title) + "' and type = 'include'")
                    abb = curs.fetchall()
                    if(not abb):
                        curs.execute("insert into back (title, link, type) value ('" + pymysql.escape_string(results[0]) + "', '" + pymysql.escape_string(title) + "',  'include')")
                        conn.commit()
                
                    data = re.sub("\[include\(((?:(?!\)\]|,).)*)((?:,\s?(?:[^)]*))+)?\)\]", "<a class=\"not_thing\" href=\"" + parse.quote(results[0]).replace('/','%2F') + "\">" + results[0] + "</a>", data, 1)
        else:
            break
    
    while(True):
        m = re.search('^#(?:redirect|넘겨주기)\s([^\n]*)', data)
        if(m):
            results = m.groups()
            aa = re.search("^(.*)(#(?:.*))$", results[0])
            if(aa):
                results = aa.groups()
                data = re.sub('^#(?:redirect|넘겨주기)\s([^\n]*)', '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(results[0]).replace('/','%2F') + '/from/' + parse.quote(title).replace('/','%2F') + results[1] + '" />', data, 1)
            else:
                data = re.sub('^#(?:redirect|넘겨주기)\s([^\n]*)', '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(results[0]).replace('/','%2F') + '/from/' + parse.quote(title).replace('/','%2F') + '" />', data, 1)
            curs.execute("select * from back where title = '" + pymysql.escape_string(results[0]) + "' and link = '" + pymysql.escape_string(title) + "' and type = 'redirect'")
            abb = curs.fetchall()
            if(not abb):
                curs.execute("insert into back (title, link, type) value ('" + pymysql.escape_string(results[0]) + "', '" + pymysql.escape_string(title) + "',  'redirect')")
                conn.commit()
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
            data = re.sub('(={1,6})\s?([^=]*)\s?(?:={1,6})(?:\s+)?\n', '<h' + str(wiki) + ' id="' + c + '"><a href="#toc" id="s-' + toc + '">' + toc + '.</a> ' + c + ' <span style="font-size:11px;">[<a href="/edit/' + parse.quote(title).replace('/','%2F') + '/section/' + str(i) + '">편집</a>]</span></h' + str(wiki) + '>', data, 1);
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
                curs.execute("select * from cat where title = '" + pymysql.escape_string(g[0]) + "' and cat = '" + pymysql.escape_string(title) + "'")
                abb = curs.fetchall()
                if(not abb):
                    curs.execute("insert into cat (title, cat) value ('" + pymysql.escape_string(g[0]) + "', '" + pymysql.escape_string(title) + "')")
                    conn.commit()                
                    
                if(category == ''):
                    curs.execute("select * from data where title = '" + pymysql.escape_string(g[0]) + "'")
                    exists = curs.fetchall()
                    if(exists):
                        red = ""
                    else:
                        red = 'class="not_thing"'
                        
                    category = category + '<a ' + red + ' href="/w/' + parse.quote(g[0]).replace('/','%2F') + '">' + re.sub("분류:", "", g[0]) + '</a>'
                else:
                    curs.execute("select * from data where title = '" + pymysql.escape_string(g[0]) + "'")
                    exists = curs.fetchall()
                    if(exists):
                        red = ""
                    else:
                        red = 'class="not_thing"'
                        
                    category = category + ' / ' + '<a ' + red + ' href="/w/' + parse.quote(g[0]).replace('/','%2F') + '">' + re.sub("분류:", "", g[0]) + '</a>'
            
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
                curs.execute("select * from back where title = '" + pymysql.escape_string(c[0]) + "' and link = '" + pymysql.escape_string(title) + "' and type = 'redirect'")
                abb = curs.fetchall()
                if(not abb):
                    curs.execute("insert into back (title, link, type) value ('파일:" + pymysql.escape_string(c[0]) + "', '" + pymysql.escape_string(title) + "',  'file')")
                    conn.commit()            
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
                            curs.execute("select * from data where title = '" + pymysql.escape_string(results[0]) + "'")
                            rows = curs.fetchall()
                            if(rows):
                                data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a title="' + results[0] + results[1] + '" href="/w/' + parse.quote(results[0]).replace('/','%2F') + results[1] + '">' + g + '</a>', data, 1)
                            else:
                                data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a title="' + results[0] + results[1] + '" class="not_thing" href="/w/' + parse.quote(results[0]).replace('/','%2F') + results[1] + '">' + g + '</a>', data, 1)
                            curs.execute("select * from back where title = '" + pymysql.escape_string(results[0]) + "' and link = '" + pymysql.escape_string(title) + "'")
                            rows = curs.fetchall()
                            if(not rows):
                                curs.execute("insert into back (title, link, type) value ('" + pymysql.escape_string(results[0]) + "', '" + pymysql.escape_string(title) + "', '')")
                                conn.commit()
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
                            curs.execute("select * from data where title = '" + pymysql.escape_string(results[0]) + "'")
                            rows = curs.fetchall()
                            if(rows):
                                data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a title="' + results[0] + '" href="/w/' + parse.quote(results[0]).replace('/','%2F') + '">' + results[1] + '</a>', data, 1)
                            else:
                                data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a title="' + results[0] + '" class="not_thing" href="/w/' + parse.quote(results[0]).replace('/','%2F') + '">' + results[1] + '</a>', data, 1)
                            curs.execute("select * from back where title = '" + pymysql.escape_string(results[0]) + "' and link = '" + pymysql.escape_string(title) + "'")
                            rows = curs.fetchall()
                            if(not rows):
                                curs.execute("insert into back (title, link, type) value ('" + pymysql.escape_string(results[0]) + "', '" + pymysql.escape_string(title) + "', '')")
                                conn.commit()
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
                            curs.execute("select * from data where title = '" + pymysql.escape_string(result[0]) + "'")
                            rows = curs.fetchall()
                            if(rows):
                                data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a href="/w/' + parse.quote(result[0]).replace('/','%2F') + result[1] + '">' + result[0] + result[1] + '</a>', data, 1)
                            else:
                                data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a class="not_thing" href="/w/' + parse.quote(result[0]).replace('/','%2F') + result[1] + '">' + result[0] + result[1] + '</a>', data, 1)
                            curs.execute("select * from back where title = '" + pymysql.escape_string(result[0]) + "' and link = '" + pymysql.escape_string(title) + "'")
                            rows = curs.fetchall()
                            if(not rows):
                                curs.execute("insert into back (title, link, type) value ('" + pymysql.escape_string(result[0]) + "', '" + pymysql.escape_string(title) + "', '')")
                                conn.commit()
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
                            curs.execute("select * from data where title = '" + pymysql.escape_string(result[0]) + "'")
                            rows = curs.fetchall()
                            if(rows):
                                data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a href="/w/' + parse.quote(result[0]).replace('/','%2F') + '">' + result[0] + '</a>', data, 1)
                            else:
                                data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a class="not_thing" href="/w/' + parse.quote(result[0]).replace('/','%2F') + '">' + result[0] + '</a>', data, 1)
                            curs.execute("select * from back where title = '" + pymysql.escape_string(result[0]) + "' and link = '" + pymysql.escape_string(title) + "'")
                            rows = curs.fetchall()
                            if(not rows):
                                curs.execute("insert into back (title, link, type) value ('" + pymysql.escape_string(result[0]) + "', '" + pymysql.escape_string(title) + "', '')")
                                conn.commit()
        else:
            break
            
    while(True):
        m = re.search("(http(?:s)?:\/\/(?:(?:(?:(?!\.[Jj][Pp][Gg]|\.[Pp][Nn][Gg]|\.[Gg][Ii][Ff]|\.[Jj][Pp][Ee][Gg]|#[Jj][Pp][Gg]#|#[Pp][Nn][Gg]#|#[Gg][Ii][Ff]#|#[Jj][Pp][Ee][Gg]#).)*)(?:\.[Jj][Pp][Gg]|\.[Pp][Nn][Gg]|\.[Gg][Ii][Ff]|\.[Jj][Pp][Ee][Gg])))(?:(?:(?:\?)width=((?:[0-9]*)(?:px|%)?))?(?:(?:\?|&)height=((?:[0-9]*)(?:px|%)?))?(?:(?:&)width=((?:[0-9]*)(?:px|%)?))?)?", data)
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
            data = re.sub("(http(?:s)?:\/\/(?:(?:(?:(?!\.[Jj][Pp][Gg]|\.[Pp][Nn][Gg]|\.[Gg][Ii][Ff]|\.[Jj][Pp][Ee][Gg]|#[Jj][Pp][Gg]#|#[Pp][Nn][Gg]#|#[Gg][Ii][Ff]#|#[Jj][Pp][Ee][Gg]#).)*)(?:\.[Jj][Pp][Gg]|\.[Pp][Nn][Gg]|\.[Gg][Ii][Ff]|\.[Jj][Pp][Ee][Gg])))(?:(?:(?:\?)width=((?:[0-9]*)(?:px)?))?(?:(?:\?|&)height=((?:[0-9]*)(?:px)?))?(?:(?:&)width=((?:[0-9]*)(?:px)?))?)?", "<img width='" + width + "' height='" + height + "' src='" + c + "'>", data, 1)
        else:
            break
            
    while(True):
        m = re.search("((?:(?:\s\*\s[^\n]*)\n?)+)", data)
        if(m):
            result = m.groups()
            end = str(result[0])
            end = re.sub("\s\*\s(?P<in>[^\n]*)", "<li>\g<in></li>", end)
            end = re.sub("\n", '', end)
            data = re.sub("((?:(?:\s\*\s[^\n]*)\n?)+)", '<ul id="list">' + end + '</ul>', data, 1)
        else:
            break
    
    data = re.sub('\[date\]', getnow(), data)
    
    data = re.sub("#(?P<in>[Jj][Pp][Gg]|[Pp][Nn][Gg]|[Gg][Ii][Ff]|[Jj][Pp][Ee][Gg])#", ".\g<in>", data)
    
    data = re.sub("-{4,11}", "<hr>", data)
    
    while(True):
        b = re.search("\r\n( +)", data)
        if(b):
            result = b.groups()
            tp = len(result[0])
            up = ''
            i = 0
            while(True):
                up = up + '<span id="in"></span>'
                i = i + 1
                if(i == tp):
                    break
            data = re.sub("\r\n( +)", '<br>' + up, data, 1)
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

def getip(request):
    if(session.get('Now') == True):
        ip = format(session['DREAMER'])
    else:
        if(request.headers.getlist("X-Forwarded-For")):
            ip = request.headers.getlist("X-Forwarded-For")[0]
        else:
            ip = request.remote_addr
    return ip

def getcan(ip, name):
    m = re.search("^사용자:(.*)", name)
    n = re.search("^파일:(.*)", name)
    if(m):
        g = m.groups()
        if(ip == g[0]):
            if(re.search("\.", g[0])):
                return 1
            else:
                curs.execute("select * from ban where block = '" + pymysql.escape_string(ip) + "'")
                rows = curs.fetchall()
                if(rows):
                    return 1
                else:
                    return 0
        else:
            return 1
    elif(n):
        return 1
    else:
        b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
        if(b):
            results = b.groups()
            curs.execute("select * from ban where block = '" + pymysql.escape_string(results[0]) + "' and band = 'O'")
            rowss = curs.fetchall()
            if(rowss):
                return 1
            else:
                curs.execute("select * from ban where block = '" + pymysql.escape_string(ip) + "'")
                rows = curs.fetchall()
                if(rows):
                    return 1
                else:
                    curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
                    row = curs.fetchall()
                    if(row):
                        curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
                        rows = curs.fetchall()
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
            curs.execute("select * from ban where block = '" + pymysql.escape_string(ip) + "'")
            rows = curs.fetchall()
            if(rows):
                return 1
            else:
                curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
                row = curs.fetchall()
                if(row):
                    curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
                    rows = curs.fetchall()
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

def getban(ip):
    b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
    if(b):
        results = b.groups()
        curs.execute("select * from ban where block = '" + pymysql.escape_string(results[0]) + "' and band = 'O'")
        rowss = curs.fetchall()
        if(rowss):
            return 1
        else:
            curs.execute("select * from ban where block = '" + pymysql.escape_string(ip) + "'")
            rows = curs.fetchall()
            if(rows):
                return 1
            else:
                return 0
    else:
        curs.execute("select * from ban where block = '" + pymysql.escape_string(ip) + "'")
        rows = curs.fetchall()
        if(rows):
            return 1
        else:
            return 0
        
def getdiscuss(ip, name, sub):
    b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
    if(b):
        results = b.groups()
        curs.execute("select * from ban where block = '" + pymysql.escape_string(results[0]) + "' and band = 'O'")
        rowss = curs.fetchall()
        if(rowss):
            return 1
        else:
            curs.execute("select * from ban where block = '" + pymysql.escape_string(ip) + "'")
            rows = curs.fetchall()
            if(rows):
                return 1
            else:
                curs.execute("select * from stop where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "'")
                rows = curs.fetchall()
                if(rows):
                    return 1
                else:
                    return 0
    else:
        curs.execute("select * from ban where block = '" + pymysql.escape_string(ip) + "'")
        rows = curs.fetchall()
        if(rows):
            return 1
        else:
            curs.execute("select * from stop where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "'")
            rows = curs.fetchall()
            if(rows):
                return 1
            else:
                return 0

def getnow():
    now = time.localtime()
    s = "%04d-%02d-%02d %02d:%02d:%02d" % (now.tm_year, now.tm_mon, now.tm_mday, now.tm_hour, now.tm_min, now.tm_sec)
    return s

def discuss(title, sub, date):
    curs.execute("select * from rd where title = '" + pymysql.escape_string(title) + "' and sub = '" + pymysql.escape_string(sub) + "'")
    rows = curs.fetchall()
    if(rows):
        curs.execute("update rd set date = '" + pymysql.escape_string(date) + "' where title = '" + pymysql.escape_string(title) + "' and sub = '" + pymysql.escape_string(sub) + "'")
    else:
        curs.execute("insert into rd (title, sub, date) value ('" + pymysql.escape_string(title) + "', '" + pymysql.escape_string(sub) + "', '" + pymysql.escape_string(date) + "')")
    conn.commit()
    
def block(block, end, today, blocker, why):
    curs.execute("insert into rb (block, end, today, blocker, why) value ('" + pymysql.escape_string(block) + "', '" + pymysql.escape_string(end) + "', '" + today + "', '" + pymysql.escape_string(blocker) + "', '" + pymysql.escape_string(why) + "')")
    conn.commit()

def history(title, data, date, ip, send, leng):
    curs.execute("select * from history where title = '" + pymysql.escape_string(title) + "' order by id+0 desc limit 1")
    rows = curs.fetchall()
    if(rows):
        number = int(rows[0]['id']) + 1
        curs.execute("insert into history (id, title, data, date, ip, send, leng) value ('" + str(number) + "', '" + pymysql.escape_string(title) + "', '" + pymysql.escape_string(data) + "', '" + date + "', '" + pymysql.escape_string(ip) + "', '" + pymysql.escape_string(send) + "', '" + leng + "')")
        conn.commit()
    else:
        curs.execute("insert into history (id, title, data, date, ip, send, leng) value ('1', '" + pymysql.escape_string(title) + "', '" + pymysql.escape_string(data) + "', '" + date + "', '" + pymysql.escape_string(ip) + "', '" + pymysql.escape_string(send + ' (새 문서)') + "', '" + leng + "')")
        conn.commit()

def getleng(existing, change):
    if(existing < change):
        leng = change - existing
        leng = '+' + str(leng)
    elif(change < existing):
        leng = existing - change
        leng = '-' + str(leng)
    else:
        leng = '0'
        
    return leng    
    
@app.route('/upload', methods=['GET', 'POST'])
def upload():
    app.config['MAX_CONTENT_LENGTH'] = int(data['upload']) * 1024 * 1024
    if(request.method == 'POST'):
        ip = getip(request)
        ban = getban(ip)
        
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
                        
                        curs.execute("insert into data (title, data, acl) value ('" + pymysql.escape_string('파일:' + filename) + "', '" + pymysql.escape_string('[[파일:' + filename + ']][br][br]{{{[[파일:' + filename + ']]}}}') + "', '')")
                        conn.commit()
                        
                        history('파일:' + filename, '[[파일:' + filename + ']][br][br]{{{[[파일:' + filename + ']]}}}', getnow(), ip, '파일:' + filename + ' 업로드', '0')
                        
                        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote('파일:' + filename).replace('/','%2F') + '" />'
                else:
                    return '<meta http-equiv="refresh" content="0;url=/error/15" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/14" />'
    else:
        ip = getip(request)
        ban = getban(ip)
        
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            return render_template('index.html', logo = data['name'], title = '업로드', tn = 21, number = data['upload'])
    
@app.route('/image/<path:name>')
def image(name = None):
    if(os.path.exists(os.path.join('image', name))):
        return send_file(os.path.join('image', name), mimetype='image')
    else:
        return render_template('index.html', logo = data['name'], data = '이미지 없음.', title = '이미지 보기'), 404

@app.route('/recentchanges')
def recentchanges():
    i = 0
    div = '<div>'
    
    curs.execute("select * from history order by date desc limit 50")
    rows = curs.fetchall()
    if(rows):
        admin = admincheck()
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
                curs.execute("select * from ban where block = '" + pymysql.escape_string(rows[i]['ip']) + "'")
                row = curs.fetchall()
                if(row):
                    ban = ' <a href="/ban/' + parse.quote(rows[i]['ip']).replace('/','%2F') + '">(해제)</a>'
                else:
                    ban = ' <a href="/ban/' + parse.quote(rows[i]['ip']).replace('/','%2F') + '">(차단)</a>'
            else:
                ban = ''
                
            if(re.search('\.', rows[i]['ip'])):
                ip = rows[i]['ip']
            else:
                curs.execute("select * from data where title = '사용자:" + pymysql.escape_string(rows[i]['ip']) + "'")
                row = curs.fetchall()
                if(row):
                    ip = '<a href="/w/' + parse.quote('사용자:' + rows[i]['ip']).replace('/','%2F') + '">' + rows[i]['ip'] + '</a>'
                else:
                    ip = '<a class="not_thing" href="/w/' + parse.quote('사용자:' + rows[i]['ip']).replace('/','%2F') + '">' + rows[i]['ip'] + '</a>'
                    
            if((int(rows[i]['id']) - 1) == 0):
                revert = ''
            else:
                revert = '<a href="/revert/' + parse.quote(rows[i]['title']).replace('/','%2F') + '/r/' + str(int(rows[i]['id']) - 1) + '">(되돌리기)</a>'
                
            div = div + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:33.33%;"><a href="/w/' + parse.quote(rows[i]['title']).replace('/','%2F') + '">' + title + '</a> <a href="/history/' + parse.quote(rows[i]['title']).replace('/','%2F') + '/n/1">(역사)</a> ' + revert + ' (' + leng + ')</td><td style="text-align: center;width:33.33%;">' + ip + ban + '</td><td style="text-align: center;width:33.33%;">' + rows[i]['date'] + '</td></tr><tr><td colspan="3" style="text-align: center;width:100%;">' + send + '</td></tr></tbody></table>'
            
            i = i + 1
            
        return render_template('index.html', logo = data['name'], rows = div, tn = 3, title = '최근 변경내역')
    else:
        return render_template('index.html', logo = data['name'], rows = '', tn = 3, title = '최근 변경내역')
        
@app.route('/history/<path:name>/r/<int:num>/hidden')
def hidden(name = None, num = None):
    if(ownercheck() == 1):
        curs.execute("select * from hidhi where title = '" + pymysql.escape_string(name) + "' and re = '" + pymysql.escape_string(str(num)) + "'")
        rows = curs.fetchall()
        if(rows):
            curs.execute("delete from hidhi where title = '" + pymysql.escape_string(name) + "' and re = '" + pymysql.escape_string(str(num)) + "'")
        else:
            curs.execute("insert into hidhi (title, re) value ('" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(str(num)) + "')")
        conn.commit()
        return '<meta http-equiv="refresh" content="0;url=/history/' + parse.quote(name).replace('/','%2F') + '/n/1" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/history/' + parse.quote(name).replace('/','%2F') + '/n/1" />'
        
@app.route('/record/<path:name>/n/<int:number>')
def record(name = None, number = None):
    v = number * 50
    i = v - 50
    div = '<div>'
    
    curs.execute("select * from history where ip = '" + pymysql.escape_string(name) + "' order by date desc")
    rows = curs.fetchall()
    if(rows):
        admin = admincheck()
        
        while(True):
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                if(number != 1):
                    div = div + '<br><a href="/record/' + parse.quote(name).replace('/','%2F') + '/n/' + str(number - 1) + '">(이전)'
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
                curs.execute("select * from ban where block = '" + pymysql.escape_string(rows[i]['ip']) + "'")
                row = curs.fetchall()
                if(row):
                    ban = ' <a href="/ban/' + parse.quote(rows[i]['ip']).replace('/','%2F') + '">(해제)</a>'
                else:
                    ban = ' <a href="/ban/' + parse.quote(rows[i]['ip']).replace('/','%2F') + '">(차단)</a>'
            else:
                ban = ''
                
            if(re.search('\.', rows[i]['ip'])):
                ip = rows[i]['ip']
            else:
                curs.execute("select * from data where title = '사용자:" + pymysql.escape_string(rows[i]['ip']) + "'")
                row = curs.fetchall()
                if(row):
                    ip = '<a href="/w/' + parse.quote('사용자:' + rows[i]['ip']).replace('/','%2F') + '">' + rows[i]['ip'] + '</a>'
                else:
                    ip = '<a class="not_thing" href="/w/' + parse.quote('사용자:' + rows[i]['ip']).replace('/','%2F') + '">' + rows[i]['ip'] + '</a>'
                    
            if((int(rows[i]['id']) - 1) == 0):
                revert = ''
            else:
                revert = '<a href="/revert/' + parse.quote(rows[i]['title']).replace('/','%2F') + '/r/' + str(int(rows[i]['id']) - 1) + '">(되돌리기)</a>'
                
            div = div + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:33.33%;"><a href="/w/' + parse.quote(rows[i]['title']).replace('/','%2F') + '">' + title + '</a> r' + rows[i]['id'] + ' <a href="/history/' + parse.quote(rows[i]['title']).replace('/','%2F') + '/n/1">(역사)</a> ' + revert + ' (' + leng + ')</td><td style="text-align: center;width:33.33%;">' + ip + ban +  '</td><td style="text-align: center;width:33.33%;">' + rows[i]['date'] + '</td></tr><tr><td colspan="3" style="text-align: center;width:100%;">' + send + '</td></tr></tbody></table>'
            
            if(i == v):
                div = div + '</div>'
                if(number == 1):
                    div = div + '<br><a href="/record/' + parse.quote(name).replace('/','%2F') + '/n/' + str(number + 1) + '">(다음)'
                else:
                    div = div + '<br><a href="/record/' + parse.quote(name).replace('/','%2F') + '/n/' + str(number - 1) + '">(이전) <a href="/record/' + parse.quote(name).replace('/','%2F') + '/n/' + str(number + 1) + '">(다음)'
                break
            else:
                i = i + 1
                
        return render_template('index.html', logo = data['name'], rows = div, tn = 3, title = '유저 기록')
    else:
        return render_template('index.html', logo = data['name'], rows = '', tn = 3, title = '유저 기록')
      
@app.route('/userlog/n/<int:number>')
def userlog(number = None):
    v = number * 50
    i = v - 50
    div = ''
    
    curs.execute("select * from user")
    rows = curs.fetchall()
    if(rows):
        admin = admincheck()
        
        while(True):
            try:
                a = rows[i]
            except:
                if(number != 1):
                    div = div + '<br><a href="/userlog/n/' + str(number - 1) + '">(이전)'
                break
                
            if(admin == 1):
                curs.execute("select * from ban where block = '" + pymysql.escape_string(rows[i]['id']) + "'")
                row = curs.fetchall()
                if(row):
                    ban = ' <a href="/ban/' + parse.quote(rows[i]['id']).replace('/','%2F') + '">(해제)</a>'
                else:
                    ban = ' <a href="/ban/' + parse.quote(rows[i]['id']).replace('/','%2F') + '">(차단)</a>'
            else:
                ban = ''
                
            if(re.search('\.', rows[i]['id'])):
                ip = rows[i]['id']
            else:
                curs.execute("select * from data where title = '사용자:" + pymysql.escape_string(rows[i]['id']) + "'")
                row = curs.fetchall()
                if(row):
                    ip = '<a href="/w/' + parse.quote('사용자:' + rows[i]['id']).replace('/','%2F') + '">' + rows[i]['id'] + '</a>'
                else:
                    ip = '<a class="not_thing" href="/w/' + parse.quote('사용자:' + rows[i]['id']).replace('/','%2F') + '">' + rows[i]['id'] + '</a>'
                
            div = div + '<li>' + str(i + 1) + '. ' + ip + ban + '</li>'
            
            if(i == v):
                if(number == 1):
                    div = div + '<br><a href="/userlog/n/' + str(number + 1) + '">(다음)'
                else:
                    div = div + '<br><a href="/userlog/n/' + str(number - 1) + '">(이전) <a href="/userlog/n/' + str(number + 1) + '">(다음)'
                break
            else:
                i = i + 1
                
        return render_template('index.html', logo = data['name'], data = div, title = '유저 가입 기록')
    else:
        return render_template('index.html', logo = data['name'], data = '', title = '유저 가입 기록')
        
@app.route('/backlink/<path:name>/n/<int:number>')
def backlink(name = None, number = None):
    v = number * 50
    i = v - 50
    div = ''
    restart = 0
    
    curs.execute("select * from back where title = '" + pymysql.escape_string(name) + "' order by link asc")
    rows = curs.fetchall()
    if(rows):        
        while(True):
            try:
                a = rows[i]
            except:
                if(number != 1):
                    div = div + '<br><a href="/backlink/n/' + str(number - 1) + '">(이전)'
                break
                
            if(rows[i]['type'] == 'include'):
                curs.execute("select * from back where title = '" + pymysql.escape_string(name) + "' and link = '" + pymysql.escape_string(rows[i]['link']) + "' and type = ''")
                test = curs.fetchall()
                if(test):
                    restart = 1
                    
                    curs.execute("delete from back where title = '" + pymysql.escape_string(name) + "' and link = '" + pymysql.escape_string(rows[i]['link']) + "' and type = ''")
                    conn.commit()
                
            if(not re.search('^사용자:', rows[i]['link'])):
                curs.execute("select * from data where title = '" + pymysql.escape_string(rows[i]['link']) + "'")
                row = curs.fetchall()
                if(row):
                    aa = row[0]['data']
                    aa = re.sub("(?P<in>\[include\((?P<out>(?:(?!\)\]|,).)*)((?:,\s?(?:[^)]*))+)?\)\])", "\g<in>\n\n[[\g<out>]]\n\n", aa)
                    aa = re.sub('^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)', '[[\g<in>]]', aa)
                    aa = namumark('', aa)
                    
                    if(re.search("<a(?:(?:(?!href=).)*)?href=\"\/w\/" + parse.quote(name).replace('/','%2F') + "(?:\#[^\"]*)?\">([^<]*)<\/a>", aa)):
                        div = div + '<li><a href="/w/' + parse.quote(rows[i]['link']).replace('/','%2F') + '">' + rows[i]['link'] + '</a>'
                        
                        if(rows[i]['type']):
                            div = div + ' (' + rows[i]['type'] + ')</li>'
                        else:
                            div = div + '</li>'
                            
                        if(i == v):
                            if(number == 1):
                                div = div + '<br><a href="/backlink/' + parse.quote(name).replace('/','%2F') + '/n/' + str(number + 1) + '">(다음)'
                            else:
                                div = div + '<br><a href="/backlink/' + parse.quote(name).replace('/','%2F') + '/n/' + str(number - 1) + '">(이전) <a href="/backlink/' + parse.quote(name).replace('/','%2F') + '/n/' + str(number + 1) + '">(다음)'
                                
                            break
                        else:
                            i = i + 1
                    else:
                        curs.execute("delete from back where title = '" + pymysql.escape_string(name) + "' and link = '" + pymysql.escape_string(rows[i]['link']) + "'")
                        conn.commit()
                        
                        i = i + 1
                        v = v + 1
                else:
                    curs.execute("delete from back where title = '" + pymysql.escape_string(name) + "' and link = '" + pymysql.escape_string(rows[i]['link']) + "'")
                    conn.commit()
                    
                    i = i + 1
                    v = v + 1
            else:
                curs.execute("delete from back where title = '" + pymysql.escape_string(name) + "' and link = '" + pymysql.escape_string(rows[i]['link']) + "'")
                conn.commit()
                
                i = i + 1
                v = v + 1
                
        if(restart == 1):
            return '<meta http-equiv="refresh" content="0;url=/backlink/' + parse.quote(name).replace('/','%2F') + '/n/' + str(number) + '" />'
        else:    
            return render_template('index.html', logo = data['name'], data = div, title = name, sub = '역링크')
    else:
        return render_template('index.html', logo = data['name'], data = '', title = name, sub = '역링크')

@app.route('/recentdiscuss')
def recentdiscuss():
    i = 0
    div = '<div>'
    
    curs.execute("select * from rd order by date desc limit 50")
    rows = curs.fetchall()
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
            
            div = div + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:50%;"><a href="/topic/' + parse.quote(rows[i]['title']).replace('/','%2F') + '/sub/' + parse.quote(rows[i]['sub']).replace('/','%2F') + '">' + title + '</a> (' + sub + ')</td><td style="text-align: center;width:50%;">' + rows[i]['date'] + '</td></tr></tbody></table>'
            
            i = i + 1
            
        return render_template('index.html', logo = data['name'], rows = div, tn = 12, title = '최근 토론내역')
    else:
        return render_template('index.html', logo = data['name'], rows = '', tn = 12, title = '최근 토론내역')
         
@app.route('/blocklog/n/<int:number>')
def blocklog(number = None):
    v = number * 50
    i = v - 50
    div = '<div>'
    
    curs.execute("select * from rb order by today desc")
    rows = curs.fetchall()
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
                
        return render_template('index.html', logo = data['name'], rows = div, tn = 20, title = '유저 차단 기록')
    else:
        return render_template('index.html', logo = data['name'], rows = '', tn = 20, title = '유저 차단 기록')

@app.route('/history/<path:name>/n/<int:number>', methods=['POST', 'GET'])
def gethistory(name = None, number = None):
    if(request.method == 'POST'):
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '/r/' + request.form["b"] + '/diff/' + request.form["a"] + '" />'
    else:
        select = ''
        v = number * 50
        i = v - 50
        div = '<div>'
        
        curs.execute("select * from history where title = '" + pymysql.escape_string(name) + "' order by id+0 desc")
        rows = curs.fetchall()
        if(rows):
            admin = admincheck()
            
            while(True):
                style = ''
            
                try:
                    a = rows[i]
                except:
                    div = div + '</div>'
                    
                    if(number != 1):
                        div = div + '<br><a href="/history/' + parse.quote(name).replace('/','%2F') + '/n/' + str(number - 1) + '">(이전)'
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
                    ip = rows[i]["ip"]
                else:
                    curs.execute("select * from data where title = '사용자:" + pymysql.escape_string(rows[i]['ip']) + "'")
                    row = curs.fetchall()
                    if(row):
                        ip = '<a href="/w/' + parse.quote('사용자:' + rows[i]['ip']).replace('/','%2F') + '">' + rows[i]['ip'] + '</a>'
                    else:
                        ip = '<a class="not_thing" href="/w/' + parse.quote('사용자:' + rows[i]['ip']).replace('/','%2F') + '">' + rows[i]['ip'] + '</a>'
                        
                if(admin == 1):
                    curs.execute("select * from user where id = '" + pymysql.escape_string(rows[i]['ip']) + "'")
                    row = curs.fetchall()
                    if(row):
                        if(row[0]['acl'] == 'owner' or row[0]['acl'] == 'admin'):
                            ban = ''
                        else:
                            curs.execute("select * from ban where block = '" + pymysql.escape_string(rows[i]['ip']) + "'")
                            row = curs.fetchall()
                            if(row):
                                ban = ' <a href="/ban/' + parse.quote(rows[i]['ip']).replace('/','%2F') + '">(해제)</a>'
                            else:
                                ban = ' <a href="/ban/' + parse.quote(rows[i]['ip']).replace('/','%2F') + '">(차단)</a>'
                    else:
                        curs.execute("select * from ban where block = '" + pymysql.escape_string(rows[i]['ip']) + "'")
                        row = curs.fetchall()
                        if(row):
                            ban = ' <a href="/ban/' + parse.quote(rows[i]['ip']).replace('/','%2F') + '">(해제)</a>'
                        else:
                            ban = ' <a href="/ban/' + parse.quote(rows[i]['ip']).replace('/','%2F') + '">(차단)</a>'
                    if(ownercheck() == 1):
                        curs.execute("select * from hidhi where title = '" + pymysql.escape_string(name) + "' and re = '" + pymysql.escape_string(rows[i]['id']) + "'")
                        row = curs.fetchall()
                        if(row):                            
                            ip = ip + ' (숨김)'                            
                            hidden = ' <a href="/history/' + parse.quote(name).replace('/','%2F') + '/r/' + rows[i]['id'] + '/hidden">(공개)'
                        else:
                            hidden = ' <a href="/history/' + parse.quote(name).replace('/','%2F') + '/r/' + rows[i]['id'] + '/hidden">(숨김)'
                    else:
                        curs.execute("select * from hidhi where title = '" + pymysql.escape_string(name) + "' and re = '" + pymysql.escape_string(rows[i]['id']) + "'")
                        row = curs.fetchall()
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
                    
                    curs.execute("select * from hidhi where title = '" + pymysql.escape_string(name) + "' and re = '" + pymysql.escape_string(rows[i]['id']) + "'")
                    row = curs.fetchall()
                    if(row):
                        ip = '숨김'
                        hidden = ''
                        send = '숨김'
                        ban = ''
                        style = 'display:none;'
                        v = v + 1
                    else:
                        hidden = ''                
                        
                div = div + '<table style="width: 100%;' + style + '"><tbody><tr><td style="text-align: center;width:33.33%;">r' + rows[i]['id'] + '</a> <a href="/w/' + parse.quote(rows[i]['title']).replace('/','%2F') + '/r/' + rows[i]['id'] + '">(w)</a> <a href="/w/' + parse.quote(rows[i]['title']).replace('/','%2F') + '/raw/' + rows[i]['id'] + '">(Raw)</a> <a href="/revert/' + parse.quote(rows[i]['title']).replace('/','%2F') + '/r/' + rows[i]['id'] + '">(되돌리기)</a> (' + leng + ')</td><td style="text-align: center;width:33.33%;">' + ip + ban + hidden + '</td><td style="text-align: center;width:33.33%;">' + rows[i]['date'] + '</td></tr><tr><td colspan="3" style="text-align: center;width:100%;">' + send + '</td></tr></tbody></table>'
                
                if(i == v):
                    div = div + '</div>'
                    
                    if(number == 1):
                        div = div + '<br><a href="/history/' + parse.quote(name).replace('/','%2F') + '/n/' + str(number + 1) + '">(다음)'
                    else:
                        div = div + '<br><a href="/history/' + parse.quote(name).replace('/','%2F') + '/n/' + str(number - 1) + '">(이전) <a href="/history/' + parse.quote(name).replace('/','%2F') + '/n/' + str(number + 1) + '">(다음)'
                        
                    break
                else:
                    i = i + 1
                    
            return render_template('index.html', logo = data['name'], rows = div, tn = 5, title = name, page = parse.quote(name).replace('/','%2F'), select = select, sub = '역사')
        else:
            return render_template('index.html', logo = data['name'], rows = '', tn = 5, title = name, page = parse.quote(name).replace('/','%2F'), select = select, sub = '역사')

@app.route('/search', methods=['POST'])
def search():
    curs.execute("select * from data where title = '" + pymysql.escape_string(request.form["search"]) + "'")
    rows = curs.fetchall()
    if(rows):
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(request.form["search"]).replace('/','%2F') + '" />'
    else:
        curs.execute("select * from data where title like '%" + pymysql.escape_string(request.form["search"]) + "%'")
        rows = curs.fetchall()
        if(rows):
            i = 0
            
            div = '<li>문서가 없습니다. <a href="/w/' + parse.quote(request.form["search"]).replace('/','%2F') + '">바로가기</a></li><br>'
            
            while(True):
                try:
                    div = div + '<li><a href="/w/' + parse.quote(rows[i]['title']).replace('/','%2F') + '">' + rows[i]['title'] + '</a></li>'
                except:
                    break
                    
                i = i + 1
        else:
            return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(request.form["search"]).replace('/','%2F') + '" />'
            
        return render_template('index.html', logo = data['name'], data = div, title = '검색')

@app.route('/w/<path:name>')
def w(name = None):
    i = 0
    
    curs.execute("select * from rd where title = '" + pymysql.escape_string(name) + "' order by date asc")
    rows = curs.fetchall()
    while(True):
        try:
            a = rows[i]
        except:
            topic = ""
            break
            
        curs.execute("select * from stop where title = '" + pymysql.escape_string(rows[i]['title']) + "' and sub = '" + pymysql.escape_string(rows[i]['sub']) + "' and close = 'O'")
        row = curs.fetchall()
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
        curs.execute("select * from cat where title = '" + pymysql.escape_string(name) + "' order by cat asc")
        rows = curs.fetchall()
        if(rows):
            div = ''
            i = 0
            
            while(True):
                try:
                    a = rows[i]
                except:
                    break
                    
                curs.execute("select * from data where title = '" + pymysql.escape_string(rows[i]['cat']) + "'")
                row = curs.fetchall()
                if(row):
                    aa = row[0]['data']                  
                    aa = namumark('', aa)
                    bb = re.search('<div style="width:100%;border: 1px solid #777;padding: 5px;margin-top: 1em;">분류:((?:(?!<\/div>).)*)<\/div>', aa)
                    if(bb):
                        cc = bb.groups()
                        
                        mm = re.search("^분류:(.*)", name)
                        if(mm):
                            ee = mm.groups()
                            
                            if(re.search("<a (class=\"not_thing\")? href=\"\/w\/" + parse.quote(name).replace('/','%2F') + "\">" + ee[0] + "<\/a>", cc[0])):
                                div = div + '<li><a href="/w/' + parse.quote(rows[i]['cat']).replace('/','%2F') + '">' + rows[i]['cat'] + '</a></li>'
                                
                                i = i + 1
                            else:
                                curs.execute("delete from cat where title = '" + pymysql.escape_string(name) + "' and cat = '" + pymysql.escape_string(rows[i]['cat']) + "'")
                                conn.commit()
                                
                                i = i + 1
                        else:
                            curs.execute("delete from cat where title = '" + pymysql.escape_string(name) + "' and cat = '" + pymysql.escape_string(rows[i]['cat']) + "'")
                            conn.commit()
                            
                            i = i + 1
                    else:
                        curs.execute("delete from cat where title = '" + pymysql.escape_string(name) + "' and cat = '" + pymysql.escape_string(rows[i]['cat']) + "'")
                        conn.commit()
                        
                        i = i + 1
                else:
                    curs.execute("delete from cat where title = '" + pymysql.escape_string(name) + "' and cat = '" + pymysql.escape_string(rows[i]['cat']) + "'")
                    conn.commit()
                    
                    i = i + 1
                    
            div = '<h2>분류</h2>' + div
            
            curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
            bb = curs.fetchall()
            if(bb):
                if(bb[0]['acl'] == 'admin'):
                    acl = '(관리자)'
                elif(bb[0]['acl'] == 'user'):
                    acl = '(유저)'
                else:
                    if(not acl):
                        acl = ''
                        
                enddata = namumark(name, bb[0]['data'])
                
                m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
                if(m):
                    result = m.groups()
                    left = result[0]
                else:
                    left = ''
                    
                return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = enddata + '<br>' + div, license = data['license'], tn = 1, uppage = uppage, style = style, acl = acl, topic = topic)
            else:
                return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = div, license = data['license'], tn = 1, uppage = uppage, style = style, acl = acl, topic = topic)
        else:
            return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = '분류 문서 없음', license = data['license'], tn = 1, uppage = uppage, style = style, acl = acl, topic = topic), 404
    else:
        m = re.search("^사용자:(.*)", name)
        if(m):
            g = m.groups()
            
            curs.execute("select * from user where id = '" + pymysql.escape_string(g[0]) + "'")
            rows = curs.fetchall()
            if(rows):
                if(rows[0]['acl'] == 'owner'):
                    acl = '(소유자)'
                elif(rows[0]['acl'] == 'admin'):
                    acl = '(관리자)'
                    
        curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
        rows = curs.fetchall()
        if(rows):
            if(rows[0]['acl'] == 'admin'):
                acl = '(관리자)'
            elif(rows[0]['acl'] == 'user'):
                acl = '(유저)'
            else:
                if(not acl):
                    acl = ''
                    
            enddata = namumark(name, rows[0]['data'])
            
            m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
            if(m):
                result = m.groups()
                left = result[0]
            else:
                left = ''
                
            return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = enddata, license = data['license'], tn = 1, acl = acl, left = left, uppage = uppage, style = style, topic = topic)
        else:
            return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = '문서 없음', license = data['license'], tn = 1, uppage = uppage, style = style, acl = acl, topic = topic), 404

@app.route('/w/<path:name>/from/<path:redirect>')
def redirectw(name = None, redirect = None):
    i = 0
    
    curs.execute("select * from rd where title = '" + pymysql.escape_string(name) + "' order by date asc")
    rows = curs.fetchall()
    while(True):
        try:
            a = rows[i]
        except:
            topic = ""
            
            break
            
        curs.execute("select * from stop where title = '" + pymysql.escape_string(rows[i]['title']) + "' and sub = '" + pymysql.escape_string(rows[i]['sub']) + "' and close = 'O'")
        row = curs.fetchall()
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
        curs.execute("select * from cat where title = '" + pymysql.escape_string(name) + "' order by cat asc")
        rows = curs.fetchall()
        if(rows):
            div = ''
            i = 0
            
            while(True):
                try:
                    a = rows[i]
                except:
                    break
                    
                curs.execute("select * from data where title = '" + pymysql.escape_string(rows[i]['cat']) + "'")
                row = curs.fetchall()
                if(row):
                    aa = row[0]['data']                  
                    aa = namumark('', aa)
                    
                    bb = re.search('<div style="width:100%;border: 1px solid #777;padding: 5px;margin-top: 1em;">분류:((?:(?!<\/div>).)*)<\/div>', aa)
                    if(bb):
                        cc = bb.groups()
                        
                        mm = re.search("^분류:(.*)", name)
                        
                        if(mm):
                            ee = mm.groups()
                            
                            if(re.search("<a (class=\"not_thing\")? href=\"\/w\/" + parse.quote(name).replace('/','%2F') + "\">" + ee[0] + "<\/a>", cc[0])):
                                div = div + '<li><a href="/w/' + parse.quote(rows[i]['cat']).replace('/','%2F') + '">' + rows[i]['cat'] + '</a></li>'
                                
                                i = i + 1
                            else:
                                curs.execute("delete from cat where title = '" + pymysql.escape_string(name) + "' and cat = '" + pymysql.escape_string(rows[i]['cat']) + "'")
                                conn.commit()
                                
                                i = i + 1
                        else:
                            curs.execute("delete from cat where title = '" + pymysql.escape_string(name) + "' and cat = '" + pymysql.escape_string(rows[i]['cat']) + "'")
                            conn.commit()
                            
                            i = i + 1
                    else:
                        curs.execute("delete from cat where title = '" + pymysql.escape_string(name) + "' and cat = '" + pymysql.escape_string(rows[i]['cat']) + "'")
                        conn.commit()
                        
                        i = i + 1
                else:
                    curs.execute("delete from cat where title = '" + pymysql.escape_string(name) + "' and cat = '" + pymysql.escape_string(rows[i]['cat']) + "'")
                    conn.commit()
                    
                    i = i + 1
                    
            div = '<h2>분류</h2>' + div
            
            curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
            bb = curs.fetchall()
            if(bb):
                if(bb[0]['acl'] == 'admin'):
                    acl = '(관리자)'
                elif(bb[0]['acl'] == 'user'):
                    acl = '(유저)'
                else:
                    if(not acl):
                        acl = ''
                        
                newdata = re.sub('^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)', ' * \g<in> 문서로 넘겨주기', bb[0]['data'])
                enddata = namumark(name, newdata)
                
                m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
                if(m):
                    result = m.groups()
                    left = result[0]
                else:
                    left = ''
                    
                return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = enddata + '<br>' + div, license = data['license'], tn = 1, uppage = uppage, style = style, acl = acl, topic = topic, redirect = '<a href="/w/' + parse.quote(redirect).replace('/','%2F') + '/from/' + parse.quote(name).replace('/','%2F') + '">' + redirect + '</a>에서 넘어 왔습니다.')
            else:
                return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = div, license = data['license'], tn = 1, uppage = uppage, style = style, acl = acl, topic = topic, redirect = '<a href="/w/' + parse.quote(redirect).replace('/','%2F') + '/from/' + parse.quote(name).replace('/','%2F') + '">' + redirect + '</a>에서 넘어 왔습니다.')
        else:
            return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = '분류 문서 없음', license = data['license'], tn = 1, uppage = uppage, style = style, acl = acl, topic = topic, redirect = '<a href="/w/' + parse.quote(redirect).replace('/','%2F') + '/from/' + parse.quote(name).replace('/','%2F') + '">' + redirect + '</a>에서 넘어 왔습니다.'), 404
    else:
        m = re.search("^사용자:(.*)", name)
        if(m):
            g = m.groups()
            
            curs.execute("select * from user where id = '" + pymysql.escape_string(g[0]) + "'")
            rows = curs.fetchall()
            if(rows):
                if(rows[0]['acl'] == 'owner'):
                    acl = '(소유자)'
                elif(rows[0]['acl'] == 'admin'):
                    acl = '(관리자)'
                    
        curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
        rows = curs.fetchall()
        if(rows):
            if(rows[0]['acl'] == 'admin'):
                acl = '(관리자)'
            elif(rows[0]['acl'] == 'user'):
                acl = '(유저)'
            else:
                if(not acl):
                    acl = ''
                    
            newdata = re.sub('^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)', ' * \g<in> 문서로 넘겨주기', rows[0]["data"])
            enddata = namumark(name, newdata)
            
            m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
            if(m):
                result = m.groups()
                left = result[0]
            else:
                left = ''
                
            return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = enddata, license = data['license'], tn = 1, acl = acl, left = left, uppage = uppage, style = style, topic = topic, redirect = '<a href="/w/' + parse.quote(redirect).replace('/','%2F') + '/from/' + parse.quote(name).replace('/','%2F') + '">' + redirect + '</a>에서 넘어 왔습니다.')
        else:
            return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = '문서 없음', license = data['license'], tn = 1, uppage = uppage, style = style, acl = acl, topic = topic, redirect = '<a href="/w/' + parse.quote(redirect).replace('/','%2F') + '/from/' + parse.quote(name).replace('/','%2F') + '">' + redirect + '</a>에서 넘어 왔습니다.'), 404

@app.route('/w/<path:name>/r/<int:number>')
def rew(name = None, number = None):
    curs.execute("select * from hidhi where title = '" + pymysql.escape_string(name) + "' and re = '" + pymysql.escape_string(str(number)) + "'")
    row = curs.fetchall()
    if(row):
        if(ownercheck() == 1):
            curs.execute("select * from history where title = '" + pymysql.escape_string(name) + "' and id = '" + str(number) + "'")
            rows = curs.fetchall()
            if(rows):
                enddata = namumark(name, rows[0]['data'])
                
                m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
                if(m):
                    result = m.groups()
                    left = result[0]
                else:
                    left = ''
                    
                return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = enddata, license = data['license'], tn = 6, left = left, sub = '옛 문서')
            else:
                return '<meta http-equiv="refresh" content="0;url=/history/' + parse.quote(name).replace('/','%2F') + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'
    else:
        curs.execute("select * from history where title = '" + pymysql.escape_string(name) + "' and id = '" + str(number) + "'")
        rows = curs.fetchall()
        if(rows):
            enddata = namumark(name, rows[0]['data'])
            
            m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
            if(m):
                result = m.groups()
                left = result[0]
            else:
                left = ''
                
            return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = enddata, license = data['license'], tn = 6, left = left, sub = '옛 문서')
        else:
            return '<meta http-equiv="refresh" content="0;url=/history/' + parse.quote(name).replace('/','%2F') + '" />'

@app.route('/w/<path:name>/raw/<int:number>')
def reraw(name = None, number = None):
    curs.execute("select * from hidhi where title = '" + pymysql.escape_string(name) + "' and re = '" + pymysql.escape_string(str(number)) + "'")
    row = curs.fetchall()
    if(row):
        if(ownercheck() == 1):
            curs.execute("select * from history where title = '" + pymysql.escape_string(name) + "' and id = '" + str(number) + "'")
            rows = curs.fetchall()
            if(rows):
                enddata = re.sub('<', '&lt;', rows[0]['data'])
                enddata = re.sub('>', '&gt;', enddata)
                enddata = re.sub('"', '&quot;', enddata)
                
                enddata = '<pre>' + enddata + '</pre>'
                
                return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = enddata, license = data['license'])
            else:
                return '<meta http-equiv="refresh" content="0;url=/history/' + parse.quote(name).replace('/','%2F') + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'
    else:
        curs.execute("select * from history where title = '" + pymysql.escape_string(name) + "' and id = '" + str(number) + "'")
        rows = curs.fetchall()
        if(rows):
            enddata = re.sub('<', '&lt;', rows[0]['data'])
            enddata = re.sub('>', '&gt;', enddata)
            enddata = re.sub('"', '&quot;', enddata)
            
            enddata = '<pre>' + enddata + '</pre>'
            
            return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = enddata, license = data['license'])
        else:
            return '<meta http-equiv="refresh" content="0;url=/history/' + parse.quote(name).replace('/','%2F') + '" />'

@app.route('/raw/<path:name>')
def raw(name = None):
    curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
    rows = curs.fetchall()
    if(rows):
        enddata = re.sub('<', '&lt;', rows[0]['data'])
        enddata = re.sub('>', '&gt;', enddata)
        enddata = re.sub('"', '&quot;', enddata)
        
        enddata = '<pre>' + enddata + '</pre>'
        
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = enddata, license = data['license'], tn = 7, sub = 'Raw')
    else:
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />'

@app.route('/revert/<path:name>/r/<int:number>', methods=['POST', 'GET'])
def revert(name = None, number = None):
    if(request.method == 'POST'):
        curs.execute("select * from hidhi where title = '" + pymysql.escape_string(name) + "' and re = '" + pymysql.escape_string(str(number)) + "'")
        row = curs.fetchall()
        if(row):
            if(ownercheck() == 1):        
                curs.execute("select * from history where title = '" + pymysql.escape_string(name) + "' and id = '" + str(number) + "'")
                rows = curs.fetchall()
                if(rows):
                    ip = getip(request)
                    can = getcan(ip, name)
                    
                    if(can == 1):
                        return '<meta http-equiv="refresh" content="0;url=/ban" />'
                    else:
                        today = getnow()
                        
                        curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
                        row = curs.fetchall()
                        if(row):
                            leng = getleng(len(row[0]['data']), len(rows[0]['data']))
                            
                            curs.execute("update data set data = '" + pymysql.escape_string(rows[0]['data']) + "' where title = '" + pymysql.escape_string(name) + "'")
                            conn.commit()
                        else:
                            leng = '+' + str(len(rows[0]['data']))
                            
                            curs.execute("insert into data (title, data, acl) value ('" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(rows[0]['data']) + "', '')")
                            conn.commit()
                        history(name, rows[0]['data'], today, ip, '문서를 ' + str(number) + '판으로 되돌렸습니다.', leng)
                        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />'
                else:
                    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/3" />'
        else:
            curs.execute("select * from history where title = '" + pymysql.escape_string(name) + "' and id = '" + str(number) + "'")
            rows = curs.fetchall()
            if(rows):
                ip = getip(request)
                can = getcan(ip, name)
                
                if(can == 1):
                    return '<meta http-equiv="refresh" content="0;url=/ban" />'
                else:
                    today = getnow()
                    
                    curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
                    row = curs.fetchall()
                    if(row):
                        leng = getleng(len(row[0]['data']), len(rows[0]['data']))
                        
                        curs.execute("update data set data = '" + pymysql.escape_string(rows[0]['data']) + "' where title = '" + pymysql.escape_string(name) + "'")
                        conn.commit()
                    else:
                        leng = '+' + str(len(rows[0]['data']))
                        
                        curs.execute("insert into data (title, data, acl) value ('" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(rows[0]['data']) + "', '')")
                        conn.commit()
                    history(name, rows[0]['data'], today, ip, '문서를 ' + str(number) + '판으로 되돌렸습니다.', leng)
                    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />'            
    else:
        curs.execute("select * from hidhi where title = '" + pymysql.escape_string(name) + "' and re = '" + pymysql.escape_string(str(number)) + "'")
        row = curs.fetchall()
        if(row):
            if(ownercheck() == 1):
                ip = getip(request)
                can = getcan(ip, name)
                
                if(can == 1):
                    return '<meta http-equiv="refresh" content="0;url=/ban" />'
                else:
                    curs.execute("select * from history where title = '" + pymysql.escape_string(name) + "' and id = '" + str(number) + "'")
                    rows = curs.fetchall()
                    if(rows):
                        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), r = parse.quote(str(number)).replace('/','%2F'), tn = 13, plus = '정말 되돌리시겠습니까?', sub = '되돌리기')
                    else:
                        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/3" />'
        else:            
            ip = getip(request)
            can = getcan(ip, name)
            
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                curs.execute("select * from history where title = '" + pymysql.escape_string(name) + "' and id = '" + str(number) + "'")
                rows = curs.fetchall()
                if(rows):
                    return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), r = parse.quote(str(number)).replace('/','%2F'), tn = 13, plus = '정말 되돌리시겠습니까?', sub = '되돌리기')
                else:
                    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />'
                        
@app.route('/edit/<path:name>', methods=['POST', 'GET'])
def edit(name = None):
    if(request.method == 'POST'):
        m = re.search('(?:[^A-Za-zㄱ-힣0-9 ])', request.form["send"])
        
        if(m):
            return '<meta http-equiv="refresh" content="0;url=/error/17" />'
        else:
            today = getnow()
            
            content = re.sub("\[date\(now\)\]", today, request.form["content"])
            
            curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
            rows = curs.fetchall()
            if(rows):
                if(rows[0]['data'] == content):
                    return '<meta http-equiv="refresh" content="0;url=/error/18" />'
                else:
                    ip = getip(request)
                    can = getcan(ip, name)
                    
                    if(can == 1):
                        return '<meta http-equiv="refresh" content="0;url=/ban" />'
                    else:                        
                        leng = getleng(len(rows[0]['data']), len(content))
                        history(name, content, today, ip, request.form["send"], leng)
                        
                        curs.execute("update data set data = '" + pymysql.escape_string(content) + "' where title = '" + pymysql.escape_string(name) + "'")
                        conn.commit()
            else:
                ip = getip(request)
                can = getcan(ip, name)
                
                if(can == 1):
                    return '<meta http-equiv="refresh" content="0;url=/ban" />'
                else:
                    leng = '+' + str(len(content))
                    history(name, content, today, ip, request.form["send"], leng)
                    
                    curs.execute("insert into data (title, data, acl) value ('" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(content) + "', '')")
                    conn.commit()
                    
            isin(name)
            
            return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />'
    else:
        ip = getip(request)
        can = getcan(ip, name)
        
        if(can == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            curs.execute("select * from data where title = '" + pymysql.escape_string(data["help"]) + "'")
            rows = curs.fetchall()
            if(rows):
                newdata = re.sub('^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)', ' * \g<in> 문서로 넘겨주기', rows[0]["data"])
                left = namumark(name, newdata)
            else:
                left = ''
                
            if(re.search('\.', ip)):
                notice = '비 로그인 상태 입니다. 비 로그인으로 편집시 아이피가 역사에 기록 됩니다. 편집 시 동의 함으로 간주 됩니다.'
            else:
                notice = ''
                
            curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
            rows = curs.fetchall()
            if(rows):
                return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = rows[0]['data'], tn = 2, notice = notice, left = left, sub = '편집')
            else:
                return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = '', tn = 2, notice = notice, left = left, sub = '편집')
                
@app.route('/edit/<path:name>/section/<int:number>', methods=['POST', 'GET'])
def secedit(name = None, number = None):
    if(request.method == 'POST'):
        m = re.search('(?:[^A-Za-zㄱ-힣0-9 ])', request.form["send"])
        if(m):
            return '<meta http-equiv="refresh" content="0;url=/error/17" />'
        else:
            today = getnow()
            
            content = re.sub("\[date\(now\)\]", today, request.form["content"])
            
            curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
            rows = curs.fetchall()
            if(rows):
                if(request.form["otent"] == content):
                    return '<meta http-equiv="refresh" content="0;url=/error/18" />'
                else:
                    ip = getip(request)
                    can = getcan(ip, name)
                    
                    if(can == 1):
                        return '<meta http-equiv="refresh" content="0;url=/ban" />'
                    else:                        
                        leng = getleng(len(request.form['otent']), len(content))
                        content = rows[0]['data'].replace(request.form['otent'], content)
                        history(name, content, today, ip, request.form["send"], leng)
                        
                        curs.execute("update data set data = '" + pymysql.escape_string(content) + "' where title = '" + pymysql.escape_string(name) + "'")
                        conn.commit()
                        
                    isin(name)
                    
                    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />'
    else:
        ip = getip(request)
        can = getcan(ip, name)
        
        if(can == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            curs.execute("select * from data where title = '" + pymysql.escape_string(data["help"]) + "'")
            rows = curs.fetchall()
            if(rows):
                newdata = re.sub('^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)', ' * \g<in> 문서로 넘겨주기', rows[0]["data"])
                
                left = namumark(name, newdata)
            else:
                left = ''
                
            if(re.search('\.', ip)):
                notice = '비 로그인 상태 입니다. 비 로그인으로 편집시 아이피가 역사에 기록 됩니다. 편집 시 동의 함으로 간주 됩니다.'
            else:
                notice = ''
                
            curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
            rows = curs.fetchall()
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
                    return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = gdata, tn = 2, notice = notice, left = left, section = 1, number = number, sub = '편집')
                else:
                    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />'
                
@app.route('/preview/<path:name>', methods=['POST'])
def preview(name = None):
    ip = getip(request)
    can = getcan(ip, name)
    if(can == 1):
        return '<meta http-equiv="refresh" content="0;url=/ban" />'
    else:
        if(re.search('\.', ip)):
            notice = '비 로그인 상태 입니다. 비 로그인으로 편집시 아이피가 역사에 기록 됩니다. 편집 시 동의 함으로 간주 됩니다.'
        else:
            notice = ''
            
        newdata = request.form["content"]
        newdata = re.sub('^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)', ' * \g<in> 문서로 넘겨주기', newdata)
        enddata = namumark(name, newdata)
        
        curs.execute("select * from data where title = '" + pymysql.escape_string(data["help"]) + "'")
        rows = curs.fetchall()
        if(rows):
            newdata = re.sub('^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)', ' * \g<in> 문서로 넘겨주기', rows[0]["data"])
            
            left = namumark(name, newdata)
        else:
            left = ''
            
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = request.form["content"], tn = 2, preview = 1, enddata = enddata, left = left, notice = notice, sub = '미리보기')
        
@app.route('/preview/<path:name>/section/<int:number>', methods=['POST'])
def secpreview(name = None, number = None):
    ip = getip(request)
    can = getcan(ip, name)
    if(can == 1):
        return '<meta http-equiv="refresh" content="0;url=/ban" />'
    else:
        if(re.search('\.', ip)):
            notice = '비 로그인 상태 입니다. 비 로그인으로 편집시 아이피가 역사에 기록 됩니다. 편집 시 동의 함으로 간주 됩니다.'
        else:
            notice = ''
        newdata = request.form["content"]
        newdata = re.sub('^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)', ' * \g<in> 문서로 넘겨주기', newdata)
        enddata = namumark(name, newdata)
        curs.execute("select * from data where title = '" + pymysql.escape_string(data["help"]) + "'")
        rows = curs.fetchall()
        if(rows):
            newdata = re.sub('^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)', ' * \g<in> 문서로 넘겨주기', rows[0]["data"])
            left = namumark(name, newdata)
        else:
            left = ''
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = request.form["content"], tn = 2, preview = 1, enddata = enddata, left = left, notice = notice, section = 1, number = number, odata = request.form["otent"], sub = '미리보기')

@app.route('/delete/<path:name>', methods=['POST', 'GET'])
def delete(name = None):
    if(request.method == 'POST'):
        curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
        rows = curs.fetchall()
        if(rows):
            ip = getip(request)
            can = getcan(ip, name)
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                today = getnow()
                leng = '-' + str(len(rows[0]['data']))
                history(name, '', today, ip, '문서를 삭제 했습니다.', leng)
                curs.execute("delete from data where title = '" + pymysql.escape_string(name) + "'")
                conn.commit()
                return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />'
    else:
        curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
        rows = curs.fetchall()
        if(rows):
            ip = getip(request)
            can = getcan(ip, name)
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), tn = 8, plus = '정말 삭제 하시겠습니까?', sub = '삭제')
        else:
            return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />'

@app.route('/move/<path:name>', methods=['POST', 'GET'])
def move(name = None):
    if(request.method == 'POST'):
        curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
        rows = curs.fetchall()
        if(rows):
            ip = getip(request)
            can = getcan(ip, name)
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                today = getnow()
                leng = '0'
                curs.execute("select * from history where title = '" + pymysql.escape_string(request.form["title"]) + "'")
                row = curs.fetchall()
                if(row):
                    return '<meta http-equiv="refresh" content="0;url=/error/19" />'
                else:
                    history(name, rows[0]['data'], today, ip, '<a href="/w/' + parse.quote(name).replace('/','%2F') + '">' + name + '</a> 문서를 <a href="/w/' + parse.quote(request.form["title"]).replace('/','%2F') + '">' + request.form["title"] + '</a> 문서로 이동 했습니다.', leng)
                    curs.execute("update data set title = '" + pymysql.escape_string(request.form["title"]) + "' where title = '" + pymysql.escape_string(name) + "'")
                    curs.execute("update history set title = '" + pymysql.escape_string(request.form["title"]) + "' where title = '" + pymysql.escape_string(name) + "'")
                    conn.commit()
                    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(request.form["title"]).replace('/','%2F') + '" />'
        else:
            ip = getip(request)
            can = getcan(ip, name)
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                today = getnow()
                leng = '0'
                curs.execute("select * from history where title = '" + pymysql.escape_string(request.form["title"]) + "'")
                row = curs.fetchall()
                if(row):
                     return '<meta http-equiv="refresh" content="0;url=/error/19" />'
                else:
                    history(name, '', today, ip, '<a href="/w/' + parse.quote(name).replace('/','%2F') + '">' + name + '</a> 문서를 <a href="/w/' + parse.quote(request.form["title"]).replace('/','%2F') + '">' + request.form["title"] + '</a> 문서로 이동 했습니다.', leng)
                    curs.execute("update history set title = '" + pymysql.escape_string(request.form["title"]) + "' where title = '" + pymysql.escape_string(name) + "'")
                    conn.commit()
                    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(request.form["title"]).replace('/','%2F') + '" />'
    else:
        ip = getip(request)
        can = getcan(ip, name)
        if(can == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), tn = 9, plus = '정말 이동 하시겠습니까?', sub = '이동')

@app.route('/other')
def other():
    return render_template('index.html', title = '기타 메뉴', logo = data['name'], data = '<li><a href="/titleindex">모든 문서</a></li><li><a href="/blocklog/n/1">유저 차단 기록</a></li><li><a href="/userlog/n/1">유저 가입 기록</a></li><li><a href="/upload">업로드</a></li><li><a href="/manager/1">관리자 메뉴</a></li><li><a href="/manager/6">유저 기록</a></li><br>이 오픈나무의 버전은 <a href="https://github.com/2DU/openNAMU/blob/master/version.md">1.8.2</a> 입니다.')
    
@app.route('/manager/<int:num>', methods=['POST', 'GET'])
def manager(num = None):
    if(num == 1):
        return render_template('index.html', title = '관리자 메뉴', logo = data['name'], data = '<h2 style="margin-top: 0px;">관리자 및 소유자</h2><li><a href="/manager/2">문서 ACL</a></li><li><a href="/manager/3">유저 체크</a></li><li><a href="/manager/4">유저 차단</a></li><h2>소유자</h2><li><a href="/manager/5">관리자 권한 주기</a></li>')
    elif(num == 2):
        if(request.method == 'POST'):
            return '<meta http-equiv="refresh" content="0;url=/acl/' + parse.quote(request.form["name"]).replace('/','%2F') + '" />'
        else:
            return render_template('index.html', title = 'ACL 이동', logo = data['name'], data = '<form id="usrform" method="POST" action="/manager/2"><input name="name" type="text"><br><br><button class="btn btn-primary" type="submit">이동</button></form>')
    elif(num == 3):
        if(request.method == 'POST'):
            return '<meta http-equiv="refresh" content="0;url=/check/' + parse.quote(request.form["name"]).replace('/','%2F') + '" />'
        else:
            return render_template('index.html', title = '체크 이동', logo = data['name'], data = '<form id="usrform" method="POST" action="/manager/3"><input name="name" type="text"><br><br><button class="btn btn-primary" type="submit">이동</button></form>')
    elif(num == 4):
        if(request.method == 'POST'):
            return '<meta http-equiv="refresh" content="0;url=/ban/' + parse.quote(request.form["name"]).replace('/','%2F') + '" />'
        else:
            return render_template('index.html', title = '차단 이동', logo = data['name'], data = '<form id="usrform" method="POST" action="/manager/4"><input name="name" type="text"><br><br><button class="btn btn-primary" type="submit">이동</button><br><br><span>아이피 앞 두자리 (XXX.XXX) 입력하면 대역 차단</span></form>')
    elif(num == 5):
        if(request.method == 'POST'):
            return '<meta http-equiv="refresh" content="0;url=/admin/' + parse.quote(request.form["name"]).replace('/','%2F') + '" />'
        else:
            return render_template('index.html', title = '권한 이동', logo = data['name'], data = '<form id="usrform" method="POST" action="/manager/5"><input name="name" type="text"><br><br><button class="btn btn-primary" type="submit">이동</button></form>')   
    elif(num == 6):
        if(request.method == 'POST'):
            return '<meta http-equiv="refresh" content="0;url=/record/' + parse.quote(request.form["name"]).replace('/','%2F') + '/n/1" />'
        else:
            return render_template('index.html', title = '기록 이동', logo = data['name'], data = '<form id="usrform" method="POST" action="/manager/6"><input name="name" type="text"><br><br><button class="btn btn-primary" type="submit">이동</button></form>')    
    else:
        return '<meta http-equiv="refresh" content="0;url=/" />'
        
@app.route('/titleindex')
def titleindex():
    i = 0
    div = '<div>'
    curs.execute("select * from data order by title asc")
    rows = curs.fetchall()
    if(rows):
        while(True):
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                break
                
            div = div + '<li>' + str(i + 1) + '. <a href="/w/' + parse.quote(rows[i]['title']).replace('/','%2F') + '">' + rows[i]['title'] + '</a></li>'
            
            i = i + 1
        return render_template('index.html', logo = data['name'], rows = div + '<br><span>이 위키에는 총 ' + str(i) + '개의 문서가 있습니다.</span>', tn = 4, title = '모든 문서')
    else:
        return render_template('index.html', logo = data['name'], rows = '', tn = 4, title = '모든 문서')

@app.route('/topic/<path:name>', methods=['POST', 'GET'])
def topic(name = None):
    if(request.method == 'POST'):
        return '<meta http-equiv="refresh" content="0;url=/topic/' + parse.quote(name).replace('/','%2F') + '/sub/' + parse.quote(request.form["topic"]).replace('/','%2F') + '" />'
    else:
        div = '<div>'
        i = 0
        j = 1
        curs.execute("select * from rd where title = '" + pymysql.escape_string(name) + "' order by date asc")
        rows = curs.fetchall()
        while(True):
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                break
                
            curs.execute("select * from topic where title = '" + pymysql.escape_string(rows[i]['title']) + "' and sub = '" + pymysql.escape_string(rows[i]['sub']) + "' and id = '1' order by sub asc")
            aa = curs.fetchall()
            
            indata = namumark(name, aa[0]['data'])
            
            if(aa[0]['block'] == 'O'):
                indata = '블라인드 되었습니다.'
                block = 'style="background: gainsboro;"'
            else:
                block = ''
                
            curs.execute("select * from stop where title = '" + pymysql.escape_string(rows[i]['title']) + "' and sub = '" + pymysql.escape_string(rows[i]['sub']) + "' and close = 'O'")
            row = curs.fetchall()
            if(not row):
                div = div + '<h2><a href="/topic/' + parse.quote(rows[i]['title']).replace('/','%2F') + '/sub/' + parse.quote(rows[i]['sub']).replace('/','%2F') + '">' + str(j) + '. ' + rows[i]['sub'] + '</a></h2>'
                div = div + '<table id="toron"><tbody><tr><td id="toroncolorgreen"><a href="javascript:void(0);" id="1">#1</a> ' + aa[0]['ip'] + ' <span style="float:right;">' + aa[0]['date'] + '</span></td></tr><tr><td ' + block + '>' + indata + '</td></tr></tbody></table><br>'
                j = j + 1
                
            i = i + 1
            
        return render_template('index.html', title = name, page = parse.quote(name).replace('/','%2F'), logo = data['name'], plus = div, tn = 10, list = 1, sub = '토론 목록')
        
@app.route('/topic/<path:name>/close')
def topicstoplist(name = None):
    if(request.method == 'POST'):
        return '<meta http-equiv="refresh" content="0;url=/topic/' + parse.quote(name).replace('/','%2F') + '/sub/' + parse.quote(request.form["topic"]).replace('/','%2F') + '" />'
    else:
        div = '<div>'
        i = 0
        
        curs.execute("select * from stop where title = '" + pymysql.escape_string(name) + "' and close = 'O' order by sub asc")
        rows = curs.fetchall()
        while(True):
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                break
                
            curs.execute("select * from topic where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(rows[i]['sub']) + "' and id = '1'")
            row = curs.fetchall()
            if(row):
                j = i + 1
                indata = namumark(name, row[0]['data'])
                
                if(row[0]['block'] == 'O'):
                    indata = '블라인드 되었습니다.'
                    block = 'style="background: gainsboro;"'
                else:
                    block = ''
                    
                div = div + '<h2><a href="/topic/' + parse.quote(name).replace('/','%2F') + '/sub/' + parse.quote(rows[i]['sub']).replace('/','%2F') + '">' + str((i + 1)) + '. ' + rows[i]['sub'] + '</a></h2>'
                div = div + '<table id="toron"><tbody><tr><td id="toroncolorgreen"><a href="javascript:void(0);" id="' + str(j) + '">#' + str(j) + '</a> ' + row[0]['ip'] + ' <span style="float:right;">' + row[0]['date'] + '</span></td></tr><tr><td ' + block + '>' + indata + '</td></tr></tbody></table><br>'
                
            i = i + 1
            
        return render_template('index.html', title = name, page = parse.quote(name).replace('/','%2F'), logo = data['name'], plus = div, tn = 10, sub = '토론 목록')

@app.route('/topic/<path:name>/sub/<path:sub>', methods=['POST', 'GET'])
def sub(name = None, sub = None):
    if(request.method == 'POST'):
        curs.execute("select * from topic where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' order by id+0 desc limit 1")
        rows = curs.fetchall()
        if(rows):
            number = int(rows[0]['id']) + 1
        else:
            number = 1
            
        ip = getip(request)
        ban = getdiscuss(ip, name, sub)
        admin = admincheck()
        
        if(ban == 1 and not admin == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
            rows = curs.fetchall()
            if(rows):
                if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
                    ip = ip + ' - Admin'
                    
            today = getnow()
            discuss(name, sub, today)
            
            aa = request.form["content"]
            aa = re.sub("\[\[(분류:(?:(?:(?!\]\]).)*))\]\]", "[br]", aa)
            aa = re.sub("\[date\(now\)\]", today, aa)
            
            curs.execute("insert into topic (id, title, sub, data, date, ip, block) value ('" + str(number) + "', '" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(sub) + "', '" + pymysql.escape_string(aa) + "', '" + today + "', '" + ip + "', '')")
            conn.commit()
            
            return '<meta http-equiv="refresh" content="0;url=/topic/' + parse.quote(name).replace('/','%2F') + '/sub/' + parse.quote(sub).replace('/','%2F') + '" />'
    else:
        style = ''
        
        ip = getip(request)
        ban = getdiscuss(ip, name, sub)
        admin = admincheck()
        
        if(admin == 1):
            div = '<div>'
            
            curs.execute("select * from stop where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' and close = 'O'")
            rows = curs.fetchall()
            if(rows):
                div = div + '<a href="/topic/' + parse.quote(name).replace('/','%2F') + '/sub/' + parse.quote(sub).replace('/','%2F') + '/close">(토론 열기)</a> '
            else:
                div = div + '<a href="/topic/' + parse.quote(name).replace('/','%2F') + '/sub/' + parse.quote(sub).replace('/','%2F') + '/close">(토론 닫기)</a> '
            
            curs.execute("select * from stop where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' and close = ''")
            rows = curs.fetchall()
            if(rows):
                div = div + '<a href="/topic/' + parse.quote(name).replace('/','%2F') + '/sub/' + parse.quote(sub).replace('/','%2F') + '/stop">(토론 재개)</a>'
            else:
                div = div + '<a href="/topic/' + parse.quote(name).replace('/','%2F') + '/sub/' + parse.quote(sub).replace('/','%2F') + '/stop">(토론 정지)</a>'
            
            div = div + '<br><br>'
        else:
            div = '<div>'
            
        i = 0
        
        curs.execute("select * from stop where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "'")
        rows = curs.fetchall()
        if(rows):
            if(not admin == 1):
                style = 'display:none;'            
        
        curs.execute("select * from topic where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' order by id+0 asc")
        rows = curs.fetchall()
        
        while(True):
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                break
                
            if(i == 0):
                start = rows[i]['ip']
                
            indata = namumark(name, rows[i]['data'])
            
            if(rows[i]['block'] == 'O'):
                indata = '블라인드 되었습니다.'
                block = 'style="background: gainsboro;"'
            else:
                block = ''
                
            m = re.search("^([^-]*)\s\-\s(Close|Reopen|Stop|Restart)$", rows[i]['ip'])
            if(m):
                ban = ""
            else:
                if(admin == 1):
                    curs.execute("select * from ban where block = '" + pymysql.escape_string(rows[i]['ip']) + "'")
                    row = curs.fetchall()
                    if(rows[i]['block'] == 'O'):
                        isblock = ' <a href="/topic/' + parse.quote(name).replace('/','%2F') + '/sub/' + parse.quote(sub).replace('/','%2F') + '/b/' + str(i + 1) + '">(해제)</a>'
                    else:
                        isblock = ' <a href="/topic/' + parse.quote(name).replace('/','%2F') + '/sub/' + parse.quote(sub).replace('/','%2F') + '/b/' + str(i + 1) + '">(블라인드)</a>'
                        
                    n = re.search("\- (?:Admin)$", rows[i]['ip'])
                    if(n):
                        ban = isblock
                    else:
                        if(row):
                            ban = ' <a href="/ban/' + parse.quote(rows[i]['ip']).replace('/','%2F') + '">(해제)</a>' + isblock
                        else:
                            ban = ' <a href="/ban/' + parse.quote(rows[i]['ip']).replace('/','%2F') + '">(차단)</a>' + isblock
                else:
                    ban = ""
                    
            m = re.search("([^-]*)\s\-\s(Close|Reopen|Stop|Restart|Admin)$", rows[i]['ip'])
            if(m):
                g = m.groups()
                
                curs.execute("select * from data where title = '사용자:" + pymysql.escape_string(g[0]) + "'")
                row = curs.fetchall()
                if(row):
                    ip = '<a href="/w/' + parse.quote('사용자:' + g[0]).replace('/','%2F') + '">' + g[0] + '</a> - ' + g[1]
                else:
                    ip = '<a class="not_thing" href="/w/' + parse.quote('사용자:' + g[0]).replace('/','%2F') + '">' + g[0] + '</a> - ' + g[1]
            elif(re.search("\.", rows[i]["ip"])):
                ip = rows[i]["ip"]
            else:
                curs.execute("select * from data where title = '사용자:" + pymysql.escape_string(rows[i]['ip']) + "'")
                row = curs.fetchall()
                if(row):
                    ip = '<a href="/w/' + parse.quote('사용자:' + rows[i]['ip']).replace('/','%2F') + '">' + rows[i]['ip'] + '</a>'
                else:
                    ip = '<a class="not_thing" href="/w/' + parse.quote('사용자:' + rows[i]['ip']).replace('/','%2F') + '">' + rows[i]['ip'] + '</a>'
                    
            if(rows[i]['ip'] == start):
                j = i + 1
                
                div = div + '<table id="toron"><tbody><tr><td id="toroncolorgreen"><a href="javascript:void(0);" id="' + str(j) + '">#' + str(j) + '</a> ' + ip + ban + ' <span style="float:right;">' + rows[i]['date'] + '</span></td></tr><tr><td ' + block + '>' + indata + '</td></tr></tbody></table><br>'
            else:
                j = i + 1
                
                div = div + '<table id="toron"><tbody><tr><td id="toroncolor"><a href="javascript:void(0);" id="' + str(j) + '">#' + str(j) + '</a> ' + ip + ban + ' <span style="float:right;">' + rows[i]['date'] + '</span></td></tr><tr><td ' + block + '>' + indata + '</td></tr></tbody></table><br>'
                
            i = i + 1
            
        return render_template('index.html', title = name, page = parse.quote(name).replace('/','%2F'), suburl = parse.quote(sub).replace('/','%2F'), toron = sub, logo = data['name'], rows = div, tn = 11, ban = ban, style = style, sub = '토론')

@app.route('/topic/<path:name>/sub/<path:sub>/b/<number>')
def blind(name = None, sub = None, number = None):
    if(admincheck() == 1):
        curs.execute("select * from topic where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' and id = '" + number + "'")
        row = curs.fetchall()
        if(row):
            if(row[0]['block'] == 'O'):
                curs.execute("update topic set block = '' where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' and id = '" + number + "'")
            else:
                curs.execute("update topic set block = 'O' where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' and id = '" + number + "'")
            conn.commit()
            
            return '<meta http-equiv="refresh" content="0;url=/topic/' + name + '/sub/' + sub + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/topic/' + name + '/sub/' + sub + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/error/3" />'
        
@app.route('/topic/<path:name>/sub/<path:sub>/stop')
def topicstop(name = None, sub = None):
    if(admincheck() == 1):
        ip = getip(request)
        
        curs.execute("select * from topic where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' order by id+0 desc limit 1")
        row = curs.fetchall()
        if(row):
            today = getnow()
            
            curs.execute("select * from stop where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' and close = ''")
            rows = curs.fetchall()
            if(rows):
                curs.execute("insert into topic (id, title, sub, data, date, ip, block) value ('" + pymysql.escape_string(str(int(row[0]['id']) + 1)) + "', '" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(sub) + "', 'Restart', '" + pymysql.escape_string(today) + "', '" + pymysql.escape_string(ip) + " - Restart', '')")
                curs.execute("delete from stop where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' and close = ''")
            else:
                curs.execute("insert into topic (id, title, sub, data, date, ip, block) value ('" + pymysql.escape_string(str(int(row[0]['id']) + 1)) + "', '" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(sub) + "', 'Stop', '" + pymysql.escape_string(today) + "', '" + pymysql.escape_string(ip) + " - Stop', '')")
                curs.execute("insert into stop (title, sub, close) value ('" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(sub) + "', '')")
            conn.commit()
            
            return '<meta http-equiv="refresh" content="0;url=/topic/' + name + '/sub/' + sub + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/topic/' + name + '/sub/' + sub + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/error/3" />'                
        
@app.route('/topic/<path:name>/sub/<path:sub>/close')
def topicclose(name = None, sub = None):
    if(admincheck() == 1):
        ip = getip(request)
        
        curs.execute("select * from topic where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' order by id+0 desc limit 1")
        row = curs.fetchall()
        if(row):
            today = getnow()
            
            curs.execute("select * from stop where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' and close = 'O'")
            rows = curs.fetchall()
            if(rows):
                curs.execute("insert into topic (id, title, sub, data, date, ip, block) value ('" + pymysql.escape_string(str(int(row[0]['id']) + 1)) + "', '" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(sub) + "', 'Reopen', '" + pymysql.escape_string(today) + "', '" + pymysql.escape_string(ip) + " - Reopen', '')")
                curs.execute("delete from stop where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' and close = 'O'")
            else:
                curs.execute("insert into topic (id, title, sub, data, date, ip, block) value ('" + pymysql.escape_string(str(int(row[0]['id']) + 1)) + "', '" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(sub) + "', 'Close', '" + pymysql.escape_string(today) + "', '" + pymysql.escape_string(ip) + " - Close', '')")
                curs.execute("insert into stop (title, sub, close) value ('" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(sub) + "', 'O')")
            conn.commit()
            
            return '<meta http-equiv="refresh" content="0;url=/topic/' + name + '/sub/' + sub + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/topic/' + name + '/sub/' + sub + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/error/3" />'

@app.route('/login', methods=['POST', 'GET'])
def login():
    if(request.method == 'POST'):
        ip = getip(request)
        ban = getban(ip)
        
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            curs.execute("select * from user where id = '" + pymysql.escape_string(request.form["id"]) + "'")
            rows = curs.fetchall()
            if(rows):
                if(session.get('Now') == True):
                    return '<meta http-equiv="refresh" content="0;url=/error/11" />'
                elif(bcrypt.checkpw(bytes(request.form["pw"], 'utf-8'), bytes(rows[0]['pw'], 'utf-8'))):
                    session['Now'] = True
                    session['DREAMER'] = request.form["id"]
                    
                    curs.execute("insert into login (user, ip, today) value ('" + pymysql.escape_string(request.form["id"]) + "', '" + pymysql.escape_string(ip) + "', '" + pymysql.escape_string(getnow()) + "')")
                    conn.commit()
                    
                    return '<meta http-equiv="refresh" content="0;url=/user" />'
                else:
                    return '<meta http-equiv="refresh" content="0;url=/error/13" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/12" />'
    else:
        ip = getip(request)
        ban = getban(ip)
        
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            if(session.get('Now') == True):
                return '<meta http-equiv="refresh" content="0;url=/error/11" />'
            else:
                return render_template('index.html', title = '로그인', enter = '로그인', logo = data['name'], tn = 15)
                
@app.route('/change', methods=['POST', 'GET'])
def change():
    if(request.method == 'POST'):
        ip = getip(request)
        ban = getban(ip)
        
        if(request.form["pw2"] == request.form["pw3"]):
            if(ban == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                curs.execute("select * from user where id = '" + pymysql.escape_string(request.form["id"]) + "'")
                rows = curs.fetchall()
                if(rows):
                    if(session.get('Now') == True):
                        session['Now'] = False
                        session.pop('DREAMER', None)
                        return '<meta http-equiv="refresh" content="0;url=/change" />'
                    elif(bcrypt.checkpw(bytes(request.form["pw"], 'utf-8'), bytes(rows[0]['pw'], 'utf-8'))):
                        hashed = bcrypt.hashpw(bytes(request.form["pw2"], 'utf-8'), bcrypt.gensalt())
                        
                        curs.execute("update user set pw = '" + pymysql.escape_string(hashed.decode()) + "' where id = '" + pymysql.escape_string(request.form["id"]) + "'")
                        conn.commit()
                        
                        return '<meta http-equiv="refresh" content="0;url=/login" />'
                    else:
                        return '<meta http-equiv="refresh" content="0;url=/error/10" />'
                else:
                    return '<meta http-equiv="refresh" content="0;url=/error/9" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/20" />'
    else:
        ip = getip(request)
        ban = getban(ip)
        
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            if(session.get('Now') == True):
                session['Now'] = False
                session.pop('DREAMER', None)
                return '<meta http-equiv="refresh" content="0;url=/change" />'
            else:
                return render_template('index.html', title = '비밀번호 변경', enter = '변경', logo = data['name'], tn = 15)
                
@app.route('/check/<name>')
def check(name = None, sub = None, number = None):
    curs.execute("select * from user where id = '" + pymysql.escape_string(name) + "'")
    rows = curs.fetchall()
    if(rows and rows[0]['acl'] == 'owner' or rows and rows[0]['acl'] == 'admin'):
        return '<meta http-equiv="refresh" content="0;url=/error/4" />'
    else:
        if(admincheck() == 1):
            m = re.search('(?:[0-9](?:[0-9][0-9])?\.[0-9](?:[0-9][0-9])?\.[0-9](?:[0-9][0-9])?\.[0-9](?:[0-9][0-9])?)', name)
            if(m):
                curs.execute("select * from login where ip = '" + pymysql.escape_string(name) + "' order by today desc")
                row = curs.fetchall()
                if(row):
                    i = 0
                    c = ''
                    while(True):
                        try:
                            c = c + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:33.33%;">' + row[i]['user'] + '</td><td style="text-align: center;width:33.33%;">' + row[i]['ip'] + '</td><td style="text-align: center;width:33.33%;">' + row[i]['today'] + '</td></tr></tbody></table>'
                        except:
                            break
                        i = i + 1
                    return render_template('index.html', title = '다중 검사', logo = data['name'], tn = 22, rows = c)
                else:
                    return render_template('index.html', title = '다중 검사', logo = data['name'], tn = 22, rows = '')
            else:
                curs.execute("select * from login where user = '" + pymysql.escape_string(name) + "' order by today desc")
                row = curs.fetchall()
                if(row):
                    i = 0
                    c = ''
                    while(True):
                        try:
                            c = c + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:33.33%;">' + row[i]['user'] + '</td><td style="text-align: center;width:33.33%;">' + row[i]['ip'] + '</td><td style="text-align: center;width:33.33%;">' + row[i]['today'] + '</td></tr></tbody></table>'
                        except:
                            break
                        i = i + 1
                    return render_template('index.html', title = '다중 검사', logo = data['name'], tn = 22, rows = c)
                else:
                    return render_template('index.html', title = '다중 검사', logo = data['name'], tn = 22, rows = '')
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'

@app.route('/register', methods=['POST', 'GET'])
def register():
    if(request.method == 'POST'):
        ip = getip(request)
        ban = getban(ip)
        
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
                        curs.execute("select * from user where id = '" + pymysql.escape_string(request.form["id"]) + "'")
                        rows = curs.fetchall()
                        if(rows):
                            return '<meta http-equiv="refresh" content="0;url=/error/6" />'
                        else:
                            hashed = bcrypt.hashpw(bytes(request.form["pw"], 'utf-8'), bcrypt.gensalt())
                            if(request.form["id"] == data['owner']):
                                curs.execute("insert into user (id, pw, acl) value ('" + pymysql.escape_string(request.form["id"]) + "', '" + pymysql.escape_string(hashed.decode()) + "', 'owner')")
                            else:
                                curs.execute("insert into user (id, pw, acl) value ('" + pymysql.escape_string(request.form["id"]) + "', '" + pymysql.escape_string(hashed.decode()) + "', 'user')")
                            conn.commit()
                            return '<meta http-equiv="refresh" content="0;url=/login" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/20" />'
    else:
        ip = getip(request)
        ban = getban(ip)
        
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            return render_template('index.html', title = '회원가입', enter = '회원가입', logo = data['name'], tn = 15)

@app.route('/logout')
def logout():
    session['Now'] = False
    session.pop('DREAMER', None)
    return '<meta http-equiv="refresh" content="0;url=/user" />'

@app.route('/ban/<name>', methods=['POST', 'GET'])
def ban(name = None):
    curs.execute("select * from user where id = '" + pymysql.escape_string(name) + "'")
    rows = curs.fetchall()
    if(rows and rows[0]['acl'] == 'owner' or rows and rows[0]['acl'] == 'admin'):
        return '<meta http-equiv="refresh" content="0;url=/error/4" />'
    else:
        if(request.method == 'POST'):
            if(admincheck() == 1):
                ip = getip(request)
                
                curs.execute("select * from ban where block = '" + pymysql.escape_string(name) + "'")
                row = curs.fetchall()
                if(row):
                    block(name, '해제', getnow(), ip, '')
                    
                    curs.execute("delete from ban where block = '" + pymysql.escape_string(name) + "'")
                else:
                    b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))$", name)
                    if(b):
                        block(name, request.form["end"], getnow(), ip, request.form["why"])
                        
                        curs.execute("insert into ban (block, end, why, band) value ('" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(request.form["end"]) + "', '" + pymysql.escape_string(request.form["why"]) + "', 'O')")
                    else:
                        block(name, request.form["end"], getnow(), ip, request.form["why"])
                        
                        curs.execute("insert into ban (block, end, why, band) value ('" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(request.form["end"]) + "', '" + pymysql.escape_string(request.form["why"]) + "', '')")
                conn.commit()
                
                return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']).replace('/','%2F') + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/3" />'
        else:
            if(admincheck() == 1):
                curs.execute("select * from ban where block = '" + pymysql.escape_string(name) + "'")
                row = curs.fetchall()
                if(row):
                    now = '차단 해제'
                else:
                    b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))$", name)
                    if(b):
                        now = '대역 차단'
                    else:
                        now = '차단'
                        
                return render_template('index.html', title = name, page = parse.quote(name).replace('/','%2F'), logo = data['name'], tn = 16, now = now, today = getnow(), sub = '차단')
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/3" />'

@app.route('/acl/<path:name>', methods=['POST', 'GET'])
def acl(name = None):
    if(request.method == 'POST'):
        if(admincheck() == 1):
            curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
            row = curs.fetchall()
            if(row):
                if(request.form["select"] == 'admin'):
                   curs.execute("update data set acl = 'admin' where title = '" + pymysql.escape_string(name) + "'")
                elif(request.form["select"] == 'user'):
                    curs.execute("update data set acl = 'user' where title = '" + pymysql.escape_string(name) + "'")
                else:
                    curs.execute("update data set acl = '' where title = '" + pymysql.escape_string(name) + "'")
                conn.commit()
            return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />' 
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'
    else:
        if(admincheck() == 1):
            curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
            row = curs.fetchall()
            if(row):
                if(row[0]['acl'] == 'admin'):
                    now = '관리자만'
                elif(row[0]['acl'] == 'user'):
                    now = '유저 이상'
                else:
                    now = '일반'
                return render_template('index.html', title = name, page = parse.quote(name).replace('/','%2F'), logo = data['name'], tn = 19, now = '현재 ACL 상태는 ' + now, sub = 'ACL')
            else:
                return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />' 
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'

@app.route('/admin/<name>', methods=['POST', 'GET'])
def admin(name = None):
    if(request.method == 'POST'):
        if(ownercheck() == 1):
            curs.execute("select * from user where id = '" + pymysql.escape_string(name) + "'")
            row = curs.fetchall()
            if(row):
                if(row[0]['acl'] == 'admin' or row[0]['acl'] == 'owner'):
                    curs.execute("update user set acl = 'user' where id = '" + pymysql.escape_string(name) + "'")
                else:
                    curs.execute("update user set acl = '" + pymysql.escape_string(request.form["select"]) + "' where id = '" + pymysql.escape_string(name) + "'")
                conn.commit()
                
                return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']).replace('/','%2F') + '" />'
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/5" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'
    else:
        if(ownercheck() == 1):
            curs.execute("select * from user where id = '" + pymysql.escape_string(name) + "'")
            row = curs.fetchall()
            if(row):
                if(row[0]['acl'] == 'admin' or row[0]['acl'] == 'owner'):
                    now = '권한 해제'
                else:
                    now = '권한 부여'
                return render_template('index.html', title = name, page = parse.quote(name).replace('/','%2F'), logo = data['name'], tn = 18, now = now, sub = '권한 부여')
            else:
                return '<meta http-equiv="refresh" content="0;url=/error/5" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/error/3" />'

@app.route('/ban')
def aban():
    ip = getip(request)
    
    if(getban(ip) == 1):
        curs.execute("select * from ban where block = '" + pymysql.escape_string(ip) + "'")
        rows = curs.fetchall()
        if(rows):
            if(rows[0]['end']):
                end = rows[0]['end'] + ' 까지 차단 상태 입니다. / 사유 : ' + rows[0]['why']                
                now = getnow()
                
                now = re.sub(':', '', now)
                now = re.sub('\-', '', now)
                now = re.sub(' ', '', now)
                now = int(now)
                
                day = rows[0]['end']
                day = re.sub('\-', '', day)    
                
                if(now >= int(day + '000000')):
                    curs.execute("delete from ban where block = '" + pymysql.escape_string(ip) + "'")
                    conn.commit()
                    end = '차단이 풀렸습니다. 다시 시도 해 보세요.'
            else:
                end = '영구 차단 상태 입니다. / 사유 : ' + rows[0]['why']
        else:
            b = re.search("^([0-9](?:[0-9]?[0-9]?)\.[0-9](?:[0-9]?[0-9]?))", ip)
            if(b):
                results = b.groups()
                
                curs.execute("select * from ban where block = '" + pymysql.escape_string(results[0]) + "' and band = 'O'")
                row = curs.fetchall()
                if(row):
                    if(row[0]['end']):
                        end = row[0]['end'] + ' 까지 차단 상태 입니다. / 사유 : ' + rows[0]['why']             
                        
                        now = getnow()
                        now = re.sub(':', '', now)
                        now = re.sub('\-', '', now)
                        now = re.sub(' ', '', now)
                        now = int(now)    
                        
                        day = row[0]['end']
                        day = re.sub('\-', '', day)
                        
                        if(now >= int(day + '000000')):
                            curs.execute("delete from ban where block = '" + pymysql.escape_string(results[0]) + "' and band = 'O'")
                            conn.commit()
                            end = '차단이 풀렸습니다. 다시 시도 해 보세요.'
                    else:
                        end = '영구 차단 상태 입니다. / 사유 : ' + row[0]['why']                
    else:
        end = '권한이 맞지 않는 상태 입니다.'
    return render_template('index.html', title = '권한 오류', logo = data['name'], data = end), 401
   
@app.route('/w/<path:name>/r/<int:a>/diff/<int:b>')
def diff(name = None, a = None, b = None):
    curs.execute("select * from history where id = '" + pymysql.escape_string(str(a)) + "' and title = '" + pymysql.escape_string(name) + "'")
    rows = curs.fetchall()
    if(rows):
        curs.execute("select * from history where id = '" + pymysql.escape_string(str(b)) + "' and title = '" + pymysql.escape_string(name) + "'")
        row = curs.fetchall()
        if(row):
            indata = re.sub('<', '&lt;', rows[0]['data'])
            indata = re.sub('>', '&gt;', indata)
            indata = re.sub('"', '&quot;', indata)
            
            enddata = re.sub('<', '&lt;', row[0]['data'])
            enddata = re.sub('>', '&gt;', enddata)
            enddata = re.sub('"', '&quot;', enddata)
            
            sm = difflib.SequenceMatcher(None, indata, enddata)
            c = show_diff(sm)
            
            c = '<pre>' + c + '</pre>'
            
            return render_template('index.html', title = name, logo = data['name'], data = c, sub = '비교')
        else:
            return '<meta http-equiv="refresh" content="0;url=/history/' + parse.quote(name).replace('/','%2F') + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/history/' + parse.quote(name).replace('/','%2F') + '" />'
        
@app.route('/user')
def user():
    ip = getip(request)
    
    curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
    rows = curs.fetchall()
    if(getban(ip) == 0):
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
        curs.execute("select * from data where title = '사용자:" + pymysql.escape_string(ip) + "'")
        row = curs.fetchall()
        if(row):
            ip = '<a href="/w/' + parse.quote('사용자:' + ip).replace('/','%2F') + '">' + ip + '</a>'
        else:
            ip = '<a class="not_thing" href="/w/' + parse.quote('사용자:' + ip).replace('/','%2F') + '">' + ip + '</a>'
        
    return render_template('index.html', title = '유저 메뉴', logo = data['name'], data = ip + '<br><br><span>권한 상태 : ' + acl + '<br><br><li><a href="/login">로그인</a></li><li><a href="/logout">로그아웃</a></li><li><a href="/register">회원가입</a></li><li><a href="/change">비밀번호 변경</a></li>')

@app.route('/random')
def random():
    curs.execute("select * from data order by rand() limit 1")
    rows = curs.fetchall()
    if(rows):
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(rows[0]['title']).replace('/','%2F') + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/" />'
        
@app.route('/error/<int:num>')
def error(num = None):
    if(num == 1):
        return render_template('index.html', title = '권한 오류', logo = data['name'], data = '비 로그인 상태 입니다.'), 401
    elif(num == 2):
        return render_template('index.html', title = '권한 오류', logo = data['name'], data = '이 계정이 없습니다.'), 401
    elif(num == 3):
        return render_template('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.'), 401
    elif(num == 4):
        return render_template('index.html', title = '권한 오류', logo = data['name'], data = '관리자는 차단, 검사 할 수 없습니다.'), 401
    elif(num == 5):
        return render_template('index.html', title = '유저 오류', logo = data['name'], data = '그런 계정이 없습니다.'), 401
    elif(num == 6):
        return render_template('index.html', title = '가입 오류', logo = data['name'], data = '동일한 아이디의 유저가 있습니다.'), 401
    elif(num == 7):
        return render_template('index.html', title = '가입 오류', logo = data['name'], data = '아이디는 20글자보다 짧아야 합니다.'), 401
    elif(num == 8):
        return render_template('index.html', title = '가입 오류', logo = data['name'], data = '아이디에는 한글과 알파벳과 공백만 허용 됩니다.'), 401
    elif(num == 9):
        return render_template('index.html', title = '변경 오류', logo = data['name'], data = '그런 계정이 없습니다.'), 401
    elif(num == 10):
        return render_template('index.html', title = '변경 오류', logo = data['name'], data = '비밀번호가 다릅니다.'), 401
    elif(num == 11):
        return render_template('index.html', title = '로그인 오류', logo = data['name'], data = '이미 로그인 되어 있습니다.'), 401
    elif(num == 12):
        return render_template('index.html', title = '로그인 오류', logo = data['name'], data = '그런 계정이 없습니다.'), 401
    elif(num == 13):
        return render_template('index.html', title = '로그인 오류', logo = data['name'], data = '비밀번호가 다릅니다.'), 401
    elif(num == 14):
        return render_template('index.html', title = '업로드 오류', logo = data['name'], data = 'jpg, gif, jpeg, png만 가능 합니다.'), 401
    elif(num == 15):
        return render_template('index.html', title = '업로드 오류', logo = data['name'], data = '파일 명에 . / \ * < > | : ? 가 들어 갈 수 없습니다.'), 401
    elif(num == 16):
        return render_template('index.html', title = '업로드 오류', logo = data['name'], data = '동일한 이름의 파일이 있습니다.'), 401
    elif(num == 17):
        return render_template('index.html', title = '편집 오류', logo = data['name'], data = '편집 내용 기록에는 한글과 영어와 숫자, 공백만 허용 됩니다.'), 401
    elif(num == 18):
        return render_template('index.html', title = '편집 오류', logo = data['name'], data = '내용이 원래 문서와 동일 합니다.'), 401
    elif(num == 19):
        return render_template('index.html', title = '이동 오류', logo = data['name'], data = '이동 하려는 곳에 문서가 이미 있습니다.'), 401
    elif(num == 20):
        return render_template('index.html', title = '비밀번호 오류', logo = data['name'], data = '재 확인이랑 비밀번호가 다릅니다.'), 401
    else:
        return '<meta http-equiv="refresh" content="0;url=/" />'

@app.errorhandler(404)
def uncaughtError(error):
    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']).replace('/','%2F') + '" />'

@app.errorhandler(413)
def uncaughtError(error):
    app.config['MAX_CONTENT_LENGTH'] = (1024**3)
    return error, 401

if __name__ == '__main__':
    app.run(host = '0.0.0.0', port = int(data['port']))