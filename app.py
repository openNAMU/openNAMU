from flask import Flask, request, session, render_template, send_file
from werkzeug import secure_filename
app = Flask(__name__)

from urllib import parse
import json
import pymysql
import time
import re
import bcrypt
import os
import difflib

json_data = open('set.json').read()
data = json.loads(json_data)

print('오픈나무 시작 포트 : ' + data['port'])

import logging
log = logging.getLogger('werkzeug')
log.setLevel(logging.ERROR)

app.config['MAX_CONTENT_LENGTH'] = int(data['upload']) * 1024 * 1024

conn = pymysql.connect(host = data['host'], user = data['user'], password = data['pw'], db = data['db'], charset = 'utf8mb4')
curs = conn.cursor(pymysql.cursors.DictCursor)

app.secret_key = data['key']

ALLOWED_EXTENSIONS = set(['png', 'jpg', 'jpeg', 'gif'])

def show_diff(seqm):
    output= []
    for opcode, a0, a1, b0, b1 in seqm.get_opcodes():
        if opcode == 'equal':
            output.append(seqm.a[a0:a1])
        elif opcode == 'insert':
            output.append("<span style='background:#CFC;'>" + seqm.b[b0:b1] + "</span>")
        elif opcode == 'delete':
            output.append("<span style='background:#FDD;'>" + seqm.a[a0:a1] + "</span>")
    return ''.join(output)

def allowed_file(filename):
    return '.' in filename and \
           filename.rsplit('.', 1)[1] in ALLOWED_EXTENSIONS
           
def admincheck():
    if(session.get('Now') == True):
        ip = getip(request)
        curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
        rows = curs.fetchall()
        if(rows):
            if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
                return 1

def namumark(title, data):
    data = re.sub('<', '&lt;', data)
    data = re.sub('>', '&gt;', data)
    data = re.sub('"', '&quot;', data)
    
    while True:
        m = re.search("\[include\((((?!\)\]).)*)\)\]", data)
        if(m):
            results = m.groups()
            if(results[0] == title):
                data = re.sub("\[include\((((?!\)\]).)*)\)\]", "<b>" + results[0] + "</b>", data, 1)
            else:
                curs.execute("select * from data where title = '" + pymysql.escape_string(results[0]) + "'")
                rows = curs.fetchall()
                if(rows):
                    enddata = rows[0]['data']
                    enddata = re.sub("\[include\((((?!\)\]).)*)\)\]", "", enddata)
                    enddata = re.sub('<', '&lt;', enddata)
                    enddata = re.sub('>', '&gt;', enddata)
                    enddata = re.sub('"', '&quot;', enddata)
                    data = re.sub("\[include\((((?!\)\]).)*)\)\]", enddata, data, 1)
                else:
                    data = re.sub("\[include\((((?!\)\]).)*)\)\]", "[[" + results[0] + "]]", data, 1)
        else:
            break
    
    while True:
        m = re.search('^#(?:redirect|넘겨주기)\s([^\n]*)', data)
        if(m):
            results = m.groups()
            data = re.sub('^#(?:redirect|넘겨주기)\s([^\n]*)', '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(results[0]).replace('/','%2F') + '/redirect/' + parse.quote(title).replace('/','%2F') + '" />', data, 1)
        else:
            break
    
    
    
    data = '\n' + data + '\n'
    
    while True:
        m = re.search("\n&gt;\s?((?:[^\n]*)(?:(?:(?:(?:\n&gt;\s?)(?:[^\n]*))+)?))", data)
        if(m):
            result = m.groups()
            blockquote = result[0]
            blockquote = re.sub("\n&gt;\s?", "\n", blockquote)
            data = re.sub("\n&gt;\s?((?:[^\n]*)(?:(?:(?:(?:\n&gt;\s?)(?:[^\n]*))+)?))", "\n<blockquote>" + blockquote + "</blockquote>", data, 1)
        else:
            break

    h0c = 0;
    h1c = 0;
    h2c = 0;
    h3c = 0;
    h4c = 0;
    h5c = 0;
    last = 0;
    rtoc = '<div id="toc"><span id="toc-name">목차</span><br><br>'
    while True:
        m = re.search('(={1,6})\s?([^=]*)\s?(?:={1,6})(?:\s+)?\n', data)
        if(m):
            result = m.groups()
            wiki = len(result[0])
            if(last < wiki):
                last = wiki
            else:
                last = wiki;
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
            data = re.sub('(={1,6})\s?([^=]*)\s?(?:={1,6})(?:\s+)?\n', '<h' + str(wiki) + '><a href="#toc" id="s-' + toc + '">' + toc + '.</a> ' + result[1] + '</h' + str(wiki) + '>', data, 1);
        else:
            rtoc = rtoc + '</div>'
            break
    
    data = re.sub("\[목차\]", rtoc, data)
    
    data = re.sub("\[\[분류:(((?!\]\]).)*)\]\]", '', data)
    
    while True:
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
            
            y = re.compile("^#!wiki\sstyle=&quot;((?:(?!&quot;|\n).)*)&quot;\n?\s\n(.*)$", re.DOTALL)
            l = y.search(results[0])
            
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
            elif(l):
                result = l.groups()
                data = p.sub('<div style="' + result[0] + '">' + result[1] + '</div>', data, 1)
            else:
                data = p.sub(results[0], data, 1)
        else:
            break

    data = re.sub("'''(?P<in>.+?)'''(?!')", '<b>\g<in></b>', data)
    data = re.sub("''(?P<in>.+?)''(?!')", '<i>\g<in></i>', data)
    data = re.sub('~~(?P<in>.+?)~~(?!~)', '<s>\g<in></s>', data)
    data = re.sub('--(?P<in>.+?)--(?!-)', '<s>\g<in></s>', data)
    data = re.sub('__(?P<in>.+?)__(?!_)', '<u>\g<in></u>', data)
    data = re.sub('\^\^(?P<in>.+?)\^\^(?!\^)', '<sup>\g<in></sup>', data)
    data = re.sub(',,(?P<in>.+?),,(?!,)', '<path:sub>\g<in></sub>', data)
    
    data = re.sub('{{\|(?P<in>(?:(?:(?:(?!\|}}).)*)(?:\n?))+)\|}}', '<table><tbody><tr><td>\g<in></td></tr></tbody></table>', data)
    
    data = re.sub("##\s?(?P<in>[^\n]*)\n", "<div style='display:none;'>\g<in></div>", data);
    
    while True:
        m = re.search("\[\[파일:((?:(?!\]\]|\?).)*)(?:\?((?:(?!\]\]).)*))?\]\]", data)
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
                img = re.sub("\.(?P<in>jpg|png|gif|jpeg)", "#\g<in>#", c[0])
                data = re.sub("\[\[파일:((?:(?!\]\]|\?).)*)(?:\?((?:(?!\]\]).)*))?\]\]", '<img src="/image/' + img + '" width="' + width + '" height="' + height + '">', data, 1)
            else:
                img = re.sub("\.(?P<in>jpg|png|gif|jpeg)", "#\g<in>#", c[0])
                data = re.sub("\[\[파일:((?:(?!\]\]|\?).)*)(?:\?((?:(?!\]\]).)*))?\]\]", "<img src='/image/" + img + "'>", data, 1)
        else:
            break
    
    data = re.sub("\[br\]",'<br>', data);
    
    while True:
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
            
    while True:
        m = re.search("(http(?:s)?:\/\/(?:(?:(?:(?!\.jpg|\.png|\.gif|\.jpeg|#jpg#|#png#|#gif#|#jpeg#).)*)(?:\.jpg|\.png|\.gif|\.jpeg)))(?:(?:(?:\?)width=((?:[0-9]*)(?:px)?))?(?:(?:\?|&)height=((?:[0-9]*)(?:px)?))?(?:(?:&)width=((?:[0-9]*)(?:px)?))?)?", data)
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
            c = re.sub("\.(?P<in>jpg|png|gif|jpeg)", "#\g<in>#", c)
            data = re.sub("(http(?:s)?:\/\/(?:(?:(?:(?!\.jpg|\.png|\.gif|\.jpeg|#jpg#|#png#|#gif#|#jpeg#).)*)(?:\.jpg|\.png|\.gif|\.jpeg)))(?:(?:(?:\?)width=((?:[0-9]*)(?:px)?))?(?:(?:\?|&)height=((?:[0-9]*)(?:px)?))?(?:(?:&)width=((?:[0-9]*)(?:px)?))?)?", "<img width='" + width + "' height='" + height + "' src='" + c + "'>", data, 1)
        else:
            break
            
    while True:
        m = re.search("\[\[(((?!\]\]).)*)\]\]", data)
        if(m):
            result = m.groups()
            a = re.search("(((?!\|).)*)\|(.*)", result[0])
            if(a):
                results = a.groups()
                b = re.search("^http(?:s)?:\/\/", results[0])
                if(b):
                    c = re.search("(?:\.jpg|\.png|\.gif|\.jpeg)", results[0])
                    if(c):
                        img = results[0]
                        img = re.sub("\.(?P<in>jpg|png|gif|jpeg)", "#\g<in>#", img)
                        data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a class="out_link" href="' + img + '">' + results[2] + '</a>', data, 1)
                    else:
                        data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a class="out_link" href="' + results[0] + '">' + results[2] + '</a>', data, 1)
                else:
                    if(results[0] == title):
                        data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<b>' + results[2] + '</b>', data, 1)
                    else:
                        curs.execute("select * from data where title = '" + pymysql.escape_string(results[0]) + "'")
                        rows = curs.fetchall()
                        if(rows):
                            data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a title="' + results[0] + '" href="/w/' + parse.quote(results[0]).replace('/','%2F') + '">' + results[2] + '</a>', data, 1)
                        else:
                            data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a title="' + results[0] + '" class="not_thing" href="/w/' + parse.quote(results[0]).replace('/','%2F') + '">' + results[2] + '</a>', data, 1)
            else:
                b = re.search("^http(?:s)?:\/\/", result[0])
                if(b):
                    c = re.search("(?:\.jpg|\.png|\.gif|\.jpeg)", result[0])
                    if(c):
                        img = result[0]
                        img = re.sub("\.(?P<in>jpg|png|gif|jpeg)", "#\g<in>#", img)
                        data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a class="out_link" href="' + img + '">' + result[0] + '</a>', data, 1)
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
        else:
            break
            
    while True:
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
    data = re.sub("\[anchor\((?P<in>[^\[\]]*)\)\]", '<span id="\g<in>"></span>', data)
    
    data = re.sub("#(?P<in>jpg|png|gif|jpeg)#", ".\g<in>", data)
    
    data = re.sub("-{4,11}", "<hr>", data)
    
    while True:
        b = re.search("\r\n( +)", data)
        if(b):
            result = b.groups()
            tp = len(result[0])
            up = ''
            i = 0
            while True:
                up = up + '<span id="in"></span>'
                i = i + 1
                if(i == tp):
                    break
            data = re.sub("\r\n( +)", '<br>' + up, data, 1)
        else:
            break
    
    a = 1
    tou = "<hr id='footnote'><div class='wiki-macro-footnote'><br>"
    while True:
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
    
    data = re.sub("\[각주\](?:(?:(?:<br>+)*(?:\s+)*(?:\r+)*(?:\n+))+)?$", "", data)
    data = re.sub("\[각주\]", "<br>" + tou, data)
    data = data + tou;
    
    while True:
        m = re.search("(\|\|(?:(?:(?:.*)\n?)\|\|)+)", data)
        if(m):
            results = m.groups()
            table = results[0]
            while True:
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
                        else:
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
            
            while True:
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
                        else:
                            celstyle = celstyle + 'text-align:left;'

                        celstyle = celstyle + '"'
                        rowstyle = rowstyle + '"'
                        
                        table = re.sub("\|\|\r\n(\|\|(?:(?:\|\|)+)?)((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)?", "</td></tr><tr " + rowstyle + "><td " + cel + " " + row + " " + celstyle + ">", table, 1)
                    else:
                        cel = 'colspan="' + str(round(len(result[0]) / 2)) + '"'
                        table = re.sub("\|\|\r\n(\|\|(?:(?:\|\|)+)?)((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)?", "</td></tr><tr><td " + cel + ">", table, 1)
                else:
                    break

            while True:
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
                        else:
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
    else:
        b = re.search("^([0-9](?:[0-9][0-9])?\.[0-9](?:[0-9][0-9])?)", ip)
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
    b = re.search("^([0-9](?:[0-9][0-9])?\.[0-9](?:[0-9][0-9])?)", ip)
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
    b = re.search("^([0-9](?:[0-9][0-9])?\.[0-9](?:[0-9][0-9])?)", ip)
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

def recent(title, ip, today, send, leng):
    if(not re.search("^사용자:(.*)", title)):
        curs.execute("insert into rc (title, date, ip, send, leng, back) value ('" + pymysql.escape_string(title) + "', '" + today + "', '" + ip + "', '" + pymysql.escape_string(send) + "', '" + leng + "', '')")
        conn.commit()

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
        curs.execute("insert into history (id, title, data, date, ip, send, leng) value ('1', '" + pymysql.escape_string(title) + "', '" + pymysql.escape_string(data) + "', '" + date + "', '" + pymysql.escape_string(ip) + "', '" + pymysql.escape_string(send) + "', '" + leng + "')")
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
    return leng;
    
    
@app.route('/upload', methods=['GET', 'POST'])
def upload():
    if(request.method == 'POST'):
        ip = getip(request)
        ban = getban(ip)
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            file = request.files['file']
            if(file and allowed_file(file.filename)):
                filename = secure_filename(file.filename)
                if(os.path.exists(os.path.join('image', filename))):
                    return render_template('index.html', logo = data['name'], title = '업로드', data = '동일한 이름의 파일이 있습니다.')
                else:
                    file.save(os.path.join('image', filename))
                    return render_template('index.html', logo = data['name'], title = '업로드', data = '완료 되었습니다.')
            else:
                return render_template('index.html', logo = data['name'], title = '업로드', data = 'jpg gif jpeg png만 가능 합니다.')
    else:
        ip = getip(request)
        ban = getban(ip)
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            return render_template('index.html', logo = data['name'], title = '업로드', tn = 21, number = data['upload'])
    
@app.route('/image/<path:name>')
def image(name = None):
    filename = secure_filename(name)
    return send_file(os.path.join('image', filename), mimetype='image')
    
@app.route('/')
@app.route('/w/')
def redirect():
    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']).replace('/','%2F') + '" />'

@app.route('/recentchanges')
def recentchanges():
    i = 0
    div = '<div>'
    curs.execute("select * from rc order by date desc limit 50")
    rows = curs.fetchall()
    if(rows):
        admin = admincheck()
        while True:
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                break
            if(rows[i]['send']):
                send = rows[i]['send']
                send = re.sub('<', '&lt;', send)
                send = re.sub('>', '&gt;', send)
                send = re.sub('&lt;a href="\/w\/(?P<in>[^"]*)"&gt;(?P<out>[^&]*)&lt;\/a&gt;', '<a href="/w/\g<in>">\g<out></a>', send)
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
                curs.execute("select * from user where id = '" + pymysql.escape_string(rows[i]['ip']) + "'")
                row = curs.fetchall()
                if(row):
                    if(row[0]['acl'] == 'owner' or row[0]['acl'] == 'admin'):
                        ip = rows[i]['ip']
                    else:
                        curs.execute("select * from ban where block = '" + pymysql.escape_string(rows[i]['ip']) + "'")
                        row = curs.fetchall()
                        if(row):
                            ip = rows[i]['ip'] + ' <a href="/ban/' + parse.quote(rows[i]['ip']).replace('/','%2F') + '">(해제)</a>'
                        else:
                            ip = rows[i]['ip'] + ' <a href="/ban/' + parse.quote(rows[i]['ip']).replace('/','%2F') + '">(차단)</a>'
                else:
                    curs.execute("select * from ban where block = '" + pymysql.escape_string(rows[i]['ip']) + "'")
                    row = curs.fetchall()
                    if(row):
                        ip = rows[i]['ip'] + ' <a href="/ban/' + parse.quote(rows[i]['ip']).replace('/','%2F') + '">(해제)</a>'
                    else:
                        ip = rows[i]['ip'] + ' <a href="/ban/' + parse.quote(rows[i]['ip']).replace('/','%2F') + '">(차단)</a>'
            else:
                ip = rows[i]['ip']
            div = div + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:33.33%;"><a href="/w/' + parse.quote(rows[i]['title']).replace('/','%2F') + '">' + title + '</a> <a href="/history/' + parse.quote(rows[i]['title']).replace('/','%2F') + '/n/1">(역사)</a> (' + leng + ')</td><td style="text-align: center;width:33.33%;">' + ip + '</td><td style="text-align: center;width:33.33%;">' + rows[i]['date'] + '</td></tr><tr><td colspan="3" style="text-align: center;width:100%;">' + send + '</td></tr></tbody></table>'
            i = i + 1
        return render_template('index.html', logo = data['name'], rows = div, tn = 3, title = '최근 변경내역')
    else:
        return render_template('index.html', logo = data['name'], rows = '', tn = 3, title = '최근 변경내역')

@app.route('/recentdiscuss')
def recentdiscuss():
    i = 0
    div = '<div>'
    curs.execute("select * from rd order by date desc limit 50")
    rows = curs.fetchall()
    if(rows):
        while True:
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
         
@app.route('/recentblock')
def recentblock():
    i = 0
    div = '<div>'
    curs.execute("select * from rb order by today desc limit 50")
    rows = curs.fetchall()
    if(rows):
        while True:
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                break
            why = rows[i]['why']
            why = re.sub('<', '&lt;', why)
            why = re.sub('>', '&gt;', why)
            div = div + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:20%;">' + rows[i]['block'] + '</a></td><td style="text-align: center;width:20%;">' + rows[i]['blocker'] + '</td><td style="text-align: center;width:20%;">' + rows[i]['end'] + '</td><td style="text-align: center;width:20%;">' + rows[i]['why'] + '</td><td style="text-align: center;width:20%;">' + rows[i]['today'] + '</td></tr></tbody></table>'
            i = i + 1
        return render_template('index.html', logo = data['name'], rows = div, tn = 20, title = '최근 차단내역')
    else:
        return render_template('index.html', logo = data['name'], rows = '', tn = 20, title = '최근 차단내역')

@app.route('/history/<path:name>/n/<int:number>', methods=['POST', 'GET'])
def gethistory(name = None, number = None):
    if(request.method == 'POST'):
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '/r/' + request.form["a"] + '/diff/' + request.form["b"] + '" />'
    else:
        v = number * 50
        i = v - 50
        div = '<div>'
        curs.execute("select * from history where title = '" + pymysql.escape_string(name) + "' order by id+0 desc")
        rows = curs.fetchall()
        if(rows):
            admin = admincheck()
            while True:
                try:
                    a = rows[i]
                except:
                    div = div + '</div>'
                    if(number != 1):
                        div = div + '<br><a href="/history/' + parse.quote(name).replace('/','%2F') + '/n/' + str(number - 1) + '">(이전)'
                    break
                if(rows[i]['send']):
                    send = rows[i]['send']
                    send = re.sub('<', '&lt;', send)
                    send = re.sub('>', '&gt;', send)
                    send = re.sub('&lt;a href="\/w\/(?P<in>[^"]*)"&gt;(?P<out>[^&]*)&lt;\/a&gt;', '<a href="/w/\g<in>">\g<out></a>', send)
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
                if(admin == 1):
                    curs.execute("select * from user where id = '" + pymysql.escape_string(rows[i]['ip']) + "'")
                    row = curs.fetchall()
                    if(row):
                        if(row[0]['acl'] == 'owner' or row[0]['acl'] == 'admin'):
                            ip = rows[i]['ip']
                        else:
                            curs.execute("select * from ban where block = '" + pymysql.escape_string(rows[i]['ip']) + "'")
                            row = curs.fetchall()
                            if(row):
                                ip = rows[i]['ip'] + ' <a href="/ban/' + parse.quote(rows[i]['ip']).replace('/','%2F') + '">(해제)</a>'
                            else:
                                ip = rows[i]['ip'] + ' <a href="/ban/' + parse.quote(rows[i]['ip']).replace('/','%2F') + '">(차단)</a>'
                    else:
                        curs.execute("select * from ban where block = '" + pymysql.escape_string(rows[i]['ip']) + "'")
                        row = curs.fetchall()
                        if(row):
                            ip = rows[i]['ip'] + ' <a href="/ban/' + parse.quote(rows[i]['ip']).replace('/','%2F') + '">(해제)</a>'
                        else:
                            ip = rows[i]['ip'] + ' <a href="/ban/' + parse.quote(rows[i]['ip']).replace('/','%2F') + '">(차단)</a>'
                else:
                    ip = rows[i]['ip']
                div = div + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:33.33%;">r' + rows[i]['id'] + '</a> <a href="/w/' + parse.quote(rows[i]['title']).replace('/','%2F') + '/r/' + rows[i]['id'] + '">(w)</a> <a href="/w/' + parse.quote(rows[i]['title']).replace('/','%2F') + '/raw/' + rows[i]['id'] + '">(Raw)</a> <a href="/revert/' + parse.quote(rows[i]['title']).replace('/','%2F') + '/r/' + rows[i]['id'] + '">(되돌리기)</a> (' + leng + ')</td><td style="text-align: center;width:33.33%;">' + ip + '</td><td style="text-align: center;width:33.33%;">' + rows[i]['date'] + '</td></tr><tr><td colspan="3" style="text-align: center;width:100%;">' + send + '</td></tr></tbody></table>'
                if(i == v):
                    div = div + '</div>'
                    if(number == 1):
                        div = div + '<br><a href="/history/' + parse.quote(name).replace('/','%2F') + '/n/' + str(number + 1) + '">(다음)'
                    else:
                        div = div + '<br><a href="/history/' + parse.quote(name).replace('/','%2F') + '/n/' + str(number - 1) + '">(이전) <a href="/history/' + parse.quote(name).replace('/','%2F') + '/n/' + str(number + 1) + '">(다음)'
                    break
                else:
                    i = i + 1
            return render_template('index.html', logo = data['name'], rows = div, tn = 5, title = name, page = parse.quote(name).replace('/','%2F'))
        else:
            return render_template('index.html', logo = data['name'], rows = '', tn = 5, title = name, page = parse.quote(name).replace('/','%2F'))

@app.route('/search', methods=['POST', 'GET'])
def search():
    if(request.method == 'POST'):
        curs.execute("select * from data where title = '" + pymysql.escape_string(request.form["search"]) + "'")
        rows = curs.fetchall()
        if(rows):
            return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(request.form["search"]).replace('/','%2F') + '" />'
        else:
            curs.execute("select * from data where title like '%" + pymysql.escape_string(request.form["search"]) + "%'")
            rows = curs.fetchall()
            div = ''
            if(rows):
                i = 0
                div = div + '<li>문서가 없습니다. <a href="/w/' + parse.quote(request.form["search"]).replace('/','%2F') + '">바로가기</a></li><br>'
                while True:
                    try:
                        div = div + '<li><a href="/w/' + parse.quote(rows[i]['title']).replace('/','%2F') + '">' + rows[i]['title'] + '</a></li>'
                    except:
                        break
                    i = i + 1
            else:
                div = div + '<li>문서가 없습니다. <a href="/w/' + parse.quote(request.form["search"]).replace('/','%2F') + '">바로가기</a></li>'
            return render_template('index.html', logo = data['name'], data = div, title = '검색')
    else:
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']).replace('/','%2F') + '" />'

@app.route('/w/<path:name>')
def w(name = None):
    curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
    rows = curs.fetchall()
    if(rows):
        if(rows[0]['acl'] == 'admin'):
            acl = '(관리자)'
        elif(rows[0]['acl'] == 'user'):
            acl = '(유저)'
        else:
            acl = ''
        enddata = namumark(name, rows[0]['data'])
        m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
        if(m):
            result = m.groups()
            left = result[0]
        else:
            left = ''
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = enddata, license = data['license'], tn = 1, acl = acl, left = left)
    else:
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = '문서 없음', license = data['license'], tn = 1)

@app.route('/w/<path:name>/redirect/<redirect>')
def redirectw(name = None, redirect = None):
    curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
    rows = curs.fetchall()
    if(rows):
        if(rows[0]['acl'] == 'admin'):
            acl = '(관리자)'
        elif(rows[0]['acl'] == 'user'):
            acl = '(유저)'
        else:
            acl = ''
        newdata = rows[0]['data']
        newdata = re.sub('^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)', ' * \g<in> 문서로 넘겨주기', newdata)
        enddata = namumark(name, newdata)
        m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
        if(m):
            result = m.groups()
            left = result[0]
        else:
            left = ''
        test = redirect
        redirect = re.sub('<', '&lt;', redirect)
        redirect = re.sub('>', '&gt;', redirect)
        redirect = re.sub('"', '&quot;', redirect)
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = enddata, license = data['license'], tn = 1, redirect = '<a href="/edit/' + parse.quote(test).replace('/','%2F') + '">' + redirect + '</a>에서 넘어 왔습니다.', left = left, acl = acl)
    else:
        test = redirect
        redirect = re.sub('<', '&lt;', redirect)
        redirect = re.sub('>', '&gt;', redirect)
        redirect = re.sub('"', '&quot;', redirect)
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = '문서 없음', license = data['license'], tn = 1, redirect = '<a href="/edit/' + parse.quote(test).replace('/','%2F') + '">' + redirect + '</a>에서 넘어 왔습니다.')

@app.route('/w/<path:name>/r/<number>')
def rew(name = None, number = None):
    curs.execute("select * from history where title = '" + pymysql.escape_string(name) + "' and id = '" + number + "'")
    rows = curs.fetchall()
    if(rows):
        enddata = namumark(name, rows[0]['data'])
        m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
        if(m):
            result = m.groups()
            left = result[0]
        else:
            left = ''
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = enddata, license = data['license'], tn = 6, left = left)
    else:
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = '문서 없음', license = data['license'], tn = 6)

@app.route('/w/<path:name>/raw/<number>')
def reraw(name = None, number = None):
    curs.execute("select * from history where title = '" + pymysql.escape_string(name) + "' and id = '" + number + "'")
    rows = curs.fetchall()
    if(rows):
        enddata = re.sub('<', '&lt;', rows[0]['data'])
        enddata = re.sub('>', '&gt;', enddata)
        enddata = re.sub('"', '&quot;', enddata)
        enddata = re.sub("\n", '<br>', enddata)
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = enddata, license = data['license'])
    else:
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = '문서 없음', license = data['license'])

@app.route('/raw/<path:name>')
def raw(name = None):
    curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
    rows = curs.fetchall()
    if(rows):
        enddata = re.sub('<', '&lt;', rows[0]['data'])
        enddata = re.sub('>', '&gt;', enddata)
        enddata = re.sub('"', '&quot;', enddata)
        enddata = re.sub("\n", '<br>', enddata)
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = enddata, license = data['license'], tn = 7)
    else:
        return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = '문서 없음', license = data['license'], tn = 7)

@app.route('/revert/<path:name>/r/<number>', methods=['POST', 'GET'])
def revert(name = None, number = None):
    if(request.method == 'POST'):
        curs.execute("select * from history where title = '" + pymysql.escape_string(name) + "' and id = '" + number + "'")
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
                recent(name, ip, today, '문서를 ' + number + '판으로 되돌렸습니다.', leng)
                history(name, rows[0]['data'], today, ip, '문서를 ' + number + '판으로 되돌렸습니다.', leng)
                return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />'
        else:
            return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />'
    else:
        ip = getip(request)
        can = getcan(ip, name)
        if(can == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            curs.execute("select * from history where title = '" + pymysql.escape_string(name) + "' and id = '" + number + "'")
            rows = curs.fetchall()
            if(rows):
                return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), r = parse.quote(number).replace('/','%2F'), tn = 13, plus = '정말 되돌리시겠습니까?')
            else:
                return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />'

@app.route('/edit/<path:name>', methods=['POST', 'GET'])
def edit(name = None):
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
                leng = getleng(len(rows[0]['data']), len(request.form["content"]))
                recent(name, ip, today, request.form["send"], leng)
                history(name, request.form["content"], today, ip, request.form["send"], leng)
                curs.execute("update data set data = '" + pymysql.escape_string(request.form["content"]) + "' where title = '" + pymysql.escape_string(name) + "'")
                conn.commit()
        else:
            ip = getip(request)
            can = getcan(ip, name)
            if(can == 1):
                return '<meta http-equiv="refresh" content="0;url=/ban" />'
            else:
                today = getnow()
                leng = '+' + str(len(request.form["content"]))
                recent(name, ip, today, request.form["send"], leng)
                history(name, request.form["content"], today, ip, request.form["send"], leng)
                curs.execute("insert into data (title, data, acl) value ('" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(request.form["content"]) + "', '')")
                conn.commit()
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />'
    else:
        ip = getip(request)
        can = getcan(ip, name)
        if(can == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            curs.execute("select * from data where title = '" + pymysql.escape_string(name) + "'")
            rows = curs.fetchall()
            if(rows):
                return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = rows[0]['data'], tn = 2)
            else:
                return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = '', tn = 2)
                
@app.route('/preview/<path:name>', methods=['POST'])
def preview(name = None):
    ip = getip(request)
    can = getcan(ip, name)
    if(can == 1):
        return '<meta http-equiv="refresh" content="0;url=/ban" />'
    else:
        newdata = request.form["content"]
        newdata = re.sub('^#(?:redirect|넘겨주기)\s(?P<in>[^\n]*)', ' * \g<in> 문서로 넘겨주기', newdata)
        enddata = namumark(name, newdata)
        m = re.search('<div id="toc">((?:(?!\/div>).)*)<\/div>', enddata)
        if(m):
            result = m.groups()
            left = result[0]
        else:
            left = ''
        m = re.search("^사용자:(.*)", name)
        if(m):
            g = m.groups()
            if(ip == g[0]):
                if(re.search("\.", g[0])):
                    return render_template('index.html', title = '사문 오류', logo = data['name'], data = '사문을 사용하려면 로그인 해야 합니다.')
                else:
                    return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = request.form["content"], tn = 2, preview = 1, enddata = enddata, left = left)
            else:
                return render_template('index.html', title = '사문 오류', logo = data['name'], data = '본인 사문이 아닙니다.')
        else:
            return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), data = request.form["content"], tn = 2, preview = 1, enddata = enddata, left = left)

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
                recent(name, ip, today, '문서를 삭제 했습니다.', leng)
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
                return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), tn = 8, plus = '정말 삭제 하시겠습니까?')
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
                    return render_template('index.html', title = '이동 오류', logo = data['name'], data = '이동 하려는 곳에 문서가 이미 있습니다.')
                else:
                    recent(name, ip, today, '문서를 <a href="/w/' + pymysql.escape_string(parse.quote(request.form["title"]).replace('/','%2F')) + '">' + pymysql.escape_string(request.form["title"]) + '</a> 문서로 이동 했습니다.', leng)
                    history(name, rows[0]['data'], today, ip, '<a href="/w/' + pymysql.escape_string(parse.quote(name).replace('/','%2F')) + '">' + pymysql.escape_string(name) + '</a> 문서를 <a href="/w/' + pymysql.escape_string(parse.quote(request.form["title"]).replace('/','%2F')) + '">' + pymysql.escape_string(request.form["title"]) + '</a> 문서로 이동 했습니다.', leng)
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
                     return render_template('index.html', title = '이동 오류', logo = data['name'], data = '이동 하려는 곳에 문서가 이미 있습니다.')
                else:
                    recent(name, ip, today, '문서를 <a href="/w/' + pymysql.escape_string(parse.quote(request.form["title"]).replace('/','%2F')) + '">' + pymysql.escape_string(request.form["title"]) + '</a> 문서로 이동 했습니다.', leng)
                    history(name, rows[0]['data'], today, ip, '<a href="/w/' + pymysql.escape_string(parse.quote(name).replace('/','%2F')) + '">' + pymysql.escape_string(name) + '</a> 문서를 <a href="/w/' + pymysql.escape_string(parse.quote(request.form["title"]).replace('/','%2F')) + '">' + pymysql.escape_string(request.form["title"]) + '</a> 문서로 이동 했습니다.', leng)
                    curs.execute("update history set title = '" + pymysql.escape_string(request.form["title"]) + "' where title = '" + pymysql.escape_string(name) + "'")
                    conn.commit()
                    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(request.form["title"]).replace('/','%2F') + '" />'
    else:
        ip = getip(request)
        can = getcan(ip, name)
        if(can == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            return render_template('index.html', title = name, logo = data['name'], page = parse.quote(name).replace('/','%2F'), tn = 9, plus = '정말 이동 하시겠습니까?')

@app.route('/setup')
def setup():
    curs.execute("create table if not exists data(title text not null, data longtext not null, acl text not null)")
    curs.execute("create table if not exists history(id text not null, title text not null, data longtext not null, date text not null, ip text not null, send text not null, leng text not null)")
    curs.execute("create table if not exists rc(title text not null, date text not null, ip text not null, send text not null, leng text not null, back text not null)")
    curs.execute("create table if not exists rd(title text not null, sub text not null, date text not null)")
    curs.execute("create table if not exists user(id text not null, pw text not null, acl text not null)")
    curs.execute("create table if not exists ban(block text not null, end text not null, why text not null, band text not null)")
    curs.execute("create table if not exists topic(id text not null, title text not null, sub text not null, data longtext not null, date text not null, ip text not null, block text not null)")
    curs.execute("create table if not exists stop(title text not null, sub text not null, close text not null)")
    curs.execute("create table if not exists rb(block text not null, end text not null, today text not null, blocker text not null, why text not null)")
    curs.execute("create table if not exists login(user text not null, ip text not null, today text not null)")
    return render_template('index.html', title = '설치 완료', logo = data['name'], data = '문제 없었음')

@app.route('/other')
def other():
    return render_template('index.html', title = '기타 메뉴', logo = data['name'], data = '<li><a href="/titleindex">모든 문서</a><li><a href="/grammar">문법 설명</a></li><li><a href="/version">버전</a></li><li><a href="/recentblock">최근 차단내역</a></li><li><a href="/upload">업로드</a></li>')

@app.route('/titleindex')
def titleindex():
    i = 0
    div = '<div>'
    curs.execute("select * from data")
    rows = curs.fetchall()
    if(rows):
        while True:
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                break
            div = div + '<li><a href="/w/' + parse.quote(rows[i]['title']).replace('/','%2F') + '">' + rows[i]['title'] + '</a></li>'
            i = i + 1
        curs.execute("select TABLE_ROWS from information_schema.tables where table_name = 'data';")
        row = curs.fetchall()
        return render_template('index.html', logo = data['name'], rows = div + '<br><span>이 위키에는 총 ' + str(row[0]['TABLE_ROWS']) + '개의 문서가 있습니다.</span>', tn = 4, title = '모든 문서')
    else:
        return render_template('index.html', logo = data['name'], rows = '', tn = 4, title = '모든 문서')

@app.route('/topic/<path:name>', methods=['POST', 'GET'])
def topic(name = None):
    if(request.method == 'POST'):
        return '<meta http-equiv="refresh" content="0;url=/topic/' + parse.quote(name).replace('/','%2F') + '/sub/' + parse.quote(request.form["topic"]).replace('/','%2F') + '" />'
    else:
        div = '<div>'
        i = 0
        curs.execute("select * from topic where title = '" + pymysql.escape_string(name) + "' order by sub asc")
        rows = curs.fetchall()
        while True:
            try:
                a = rows[i]
            except:
                div = div + '</div>'
                break
            j = i + 1
            indata = namumark(name, rows[i]['data'])
            if(rows[i]['block'] == 'O'):
                indata = '블라인드 되었습니다.'
                block = 'style="background: gainsboro;"'
            else:
                block = ''
            if(i == 0):
                sub = rows[i]['sub']
                curs.execute("select * from stop where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' and close = 'O'")
                row = curs.fetchall()
                if(not row):
                    div = div + '<h2><a href="/topic/' + parse.quote(name).replace('/','%2F') + '/sub/' + parse.quote(rows[i]['sub']).replace('/','%2F') + '">' + str((i + 1)) + '. ' + rows[i]['sub'] + '</a></h2>'
                    div = div + '<table id="toron"><tbody><tr><td id="toroncolorgreen"><a href="javascript:void(0);" id="' + str(j) + '">#' + str(j) + '</a> ' + rows[i]['ip'] + ' <span style="float:right;">' + rows[i]['date'] + '</span></td></tr><tr><td ' + block + '>' + indata + '</td></tr></tbody></table><br>'
            else:
                if(not sub == rows[i]['sub']):
                    sub = rows[i]['sub']
                    curs.execute("select * from stop where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' and close = 'O'")
                    row = curs.fetchall()
                    if(not row):
                        div = div + '<h2><a href="/topic/' + parse.quote(name).replace('/','%2F') + '/sub/' + parse.quote(rows[i]['sub']).replace('/','%2F') + '">' + str((i + 1)) + '. ' + rows[i]['sub'] + '</a></h2>'
                        div = div + '<table id="toron"><tbody><tr><td id="toroncolorgreen"><a href="javascript:void(0);" id="' + str(j) + '">#' + str(j) + '</a> ' + rows[i]['ip'] + ' <span style="float:right;">' + rows[i]['date'] + '</span></td></tr><tr><td ' + block + '>' + indata + '</td></tr></tbody></table><br>'
            i = i + 1
        return render_template('index.html', title = name, page = parse.quote(name).replace('/','%2F'), logo = data['name'], plus = div, tn = 10, list = 1)
        
@app.route('/topic/<path:name>/close')
def topicstoplist(name = None):
    if(request.method == 'POST'):
        return '<meta http-equiv="refresh" content="0;url=/topic/' + parse.quote(name).replace('/','%2F') + '/sub/' + parse.quote(request.form["topic"]).replace('/','%2F') + '" />'
    else:
        div = '<div>'
        i = 0
        curs.execute("select * from stop where title = '" + pymysql.escape_string(name) + "' and close = 'O' order by sub asc")
        rows = curs.fetchall()
        while True:
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
        return render_template('index.html', title = name, page = parse.quote(name).replace('/','%2F'), logo = data['name'], plus = div, tn = 10)

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
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
            rows = curs.fetchall()
            if(rows):
                if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
                    ip = ip + ' - Admin'
            today = getnow()
            discuss(name, sub, today)
            curs.execute("insert into topic (id, title, sub, data, date, ip, block) value ('" + str(number) + "', '" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(sub) + "', '" + pymysql.escape_string(request.form["content"]) + "', '" + today + "', '" + ip + "', '')")
            conn.commit()
            return '<meta http-equiv="refresh" content="0;url=/topic/' + parse.quote(name).replace('/','%2F') + '/sub/' + parse.quote(sub).replace('/','%2F') + '" />'
    else:
        ip = getip(request)
        ban = getdiscuss(ip, name, sub)
        admin = admincheck()
        if(admin == 1):
            div = '<div>' + '<a href="/topic/' + parse.quote(name).replace('/','%2F') + '/sub/' + parse.quote(sub).replace('/','%2F') + '/close">(토론 닫기 및 열기)</a>' + ' <a href="/topic/' + parse.quote(name).replace('/','%2F') + '/sub/' + parse.quote(sub).replace('/','%2F') + '/stop">(토론 정지 및 재개)</a><br><br>'
        else:
            div = '<div>'
        i = 0
        curs.execute("select * from topic where title = '" + pymysql.escape_string(name) + "' and sub = '" + pymysql.escape_string(sub) + "' order by id+0 asc")
        rows = curs.fetchall()
        while True:
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
            m = re.search("\- (?:Close|Reopen|Stop|Restart)$", rows[i]['ip'])
            if(m):
                ip = rows[i]['ip']
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
                        ip = rows[i]['ip'] + isblock
                    else:
                        if(row):
                            ip = rows[i]['ip'] + ' <a href="/ban/' + parse.quote(rows[i]['ip']).replace('/','%2F') + '">(해제)</a>' + isblock
                        else:
                            ip = rows[i]['ip'] + ' <a href="/ban/' + parse.quote(rows[i]['ip']).replace('/','%2F') + '">(차단)</a>' + isblock
                else:
                    ip = rows[i]['ip']
            if(rows[i]['ip'] == start):
                j = i + 1
                div = div + '<table id="toron"><tbody><tr><td id="toroncolorgreen"><a href="javascript:void(0);" id="' + str(j) + '">#' + str(j) + '</a> ' + ip + ' <span style="float:right;">' + rows[i]['date'] + '</span></td></tr><tr><td ' + block + '>' + indata + '</td></tr></tbody></table><br>'
            else:
                j = i + 1
                div = div + '<table id="toron"><tbody><tr><td id="toroncolor"><a href="javascript:void(0);" id="' + str(j) + '">#' + str(j) + '</a> ' + ip + ' <span style="float:right;">' + rows[i]['date'] + '</span></td></tr><tr><td ' + block + '>' + indata + '</td></tr></tbody></table><br>'
            i = i + 1
        return render_template('index.html', title = name, page = parse.quote(name).replace('/','%2F'), suburl = parse.quote(sub).replace('/','%2F'), sub = sub, logo = data['name'], rows = div, tn = 11, ban = ban)

@app.route('/topic/<path:name>/sub/<path:sub>/b/<number>')
def blind(name = None, sub = None, number = None):
    if(session.get('Now') == True):
        ip = getip(request)
        curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
        rows = curs.fetchall()
        if(rows):
            if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
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
                return render_template('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.')
        else:
            return render_template('index.html', title = '권한 오류', logo = data['name'], data = '계정이 없습니다.')
    else:
        return render_template('index.html', title = '권한 오류', logo = data['name'], data = '비 로그인 상태 입니다.')
        
@app.route('/topic/<path:name>/sub/<path:sub>/stop')
def topicstop(name = None, sub = None):
    if(session.get('Now') == True):
        ip = getip(request)
        curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
        rows = curs.fetchall()
        if(rows):
            if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
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
                return render_template('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.')
        else:
            return render_template('index.html', title = '권한 오류', logo = data['name'], data = '계정이 없습니다.')
    else:
        return render_template('index.html', title = '권한 오류', logo = data['name'], data = '비 로그인 상태 입니다.')
        
@app.route('/topic/<path:name>/sub/<path:sub>/close')
def topicclose(name = None, sub = None):
    if(session.get('Now') == True):
        ip = getip(request)
        curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
        rows = curs.fetchall()
        if(rows):
            if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
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
                return render_template('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.')
        else:
            return render_template('index.html', title = '권한 오류', logo = data['name'], data = '계정이 없습니다.')
    else:
        return render_template('index.html', title = '권한 오류', logo = data['name'], data = '비 로그인 상태 입니다.')

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
                    return render_template('index.html', title = '로그인 오류', logo = data['name'], data = '이미 로그인 되어 있습니다.')
                elif(bcrypt.checkpw(bytes(request.form["pw"], 'utf-8'), bytes(rows[0]['pw'], 'utf-8'))):
                    session['Now'] = True
                    session['DREAMER'] = request.form["id"]
                    curs.execute("insert into login (user, ip, today) value ('" + pymysql.escape_string(request.form["id"]) + "', '" + pymysql.escape_string(ip) + "', '" + pymysql.escape_string(getnow()) + "')")
                    conn.commit()
                    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']).replace('/','%2F') + '" />'
                else:
                    return render_template('index.html', title = '로그인 오류', logo = data['name'], data = '비밀번호가 다릅니다.')
            else:
                return render_template('index.html', title = '로그인 오류', logo = data['name'], data = '없는 계정 입니다.')
    else:
        ip = getip(request)
        ban = getban(ip)
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            if(session.get('Now') == True):
                return render_template('index.html', title = '로그인 오류', logo = data['name'], data = '이미 로그인 되어 있습니다.')
            else:
                return render_template('index.html', title = '로그인', enter = '로그인', logo = data['name'], tn = 15)
                
@app.route('/check/<path:name>')
def check(name = None, sub = None, number = None):
    curs.execute("select * from user where id = '" + pymysql.escape_string(name) + "'")
    rows = curs.fetchall()
    if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
        return render_template('index.html', title = '차단 오류', logo = data['name'], data = '관리자는 검사 할 수 없습니다.')
    else:
        if(admincheck() == 1):
            m = re.search('(?:[0-9](?:[0-9][0-9])?\.[0-9](?:[0-9][0-9])?\.[0-9](?:[0-9][0-9])?\.[0-9](?:[0-9][0-9])?)', name)
            if(m):
                curs.execute("select * from login where ip = '" + pymysql.escape_string(name) + "' order by today desc")
                row = curs.fetchall()
                if(row):
                    i = 0
                    c = ''
                    while True:
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
                    while True:
                        try:
                            c = c + '<table style="width: 100%;"><tbody><tr><td style="text-align: center;width:33.33%;">' + row[i]['user'] + '</td><td style="text-align: center;width:33.33%;">' + row[i]['ip'] + '</td><td style="text-align: center;width:33.33%;">' + row[i]['today'] + '</td></tr></tbody></table>'
                        except:
                            break
                        i = i + 1
                    return render_template('index.html', title = '다중 검사', logo = data['name'], tn = 22, rows = c)
                else:
                    return render_template('index.html', title = '다중 검사', logo = data['name'], tn = 22, rows = '')
        else:
            return render_template('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.')

@app.route('/register', methods=['POST', 'GET'])
def register():
    if(request.method == 'POST'):
        ip = getip(request)
        ban = getban(ip)
        if(ban == 1):
            return '<meta http-equiv="refresh" content="0;url=/ban" />'
        else:
            m = re.search('(?:[^A-Za-zㄱ-힣0-9 ])', request.form["id"])
            if(m):
                return render_template('index.html', title = '회원가입 오류', logo = data['name'], data = '아이디에는 한글과 알파벳 공백만 허용 됩니다.')
            else:
                if(len(request.form["id"]) > 20):
                    return render_template('index.html', title = '회원가입 오류', logo = data['name'], data = '아이디는 20글자보다 짧아야 합니다.')
                else:
                    curs.execute("select * from user where id = '" + pymysql.escape_string(request.form["id"]) + "'")
                    rows = curs.fetchall()
                    if(rows):
                        return render_template('index.html', title = '회원가입 오류', logo = data['name'], data = '동일한 아이디의 유저가 있습니다.')
                    else:
                        hashed = bcrypt.hashpw(bytes(request.form["pw"], 'utf-8'), bcrypt.gensalt())
                        if(request.form["id"] == data['owner']):
                            curs.execute("insert into user (id, pw, acl) value ('" + pymysql.escape_string(request.form["id"]) + "', '" + pymysql.escape_string(hashed.decode()) + "', 'owner')")
                        else:
                            curs.execute("insert into user (id, pw, acl) value ('" + pymysql.escape_string(request.form["id"]) + "', '" + pymysql.escape_string(hashed.decode()) + "', 'user')")
                        conn.commit()
                        return '<meta http-equiv="refresh" content="0;url=/login" />'
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
    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']).replace('/','%2F') + '" />'

@app.route('/ban/<path:name>', methods=['POST', 'GET'])
def ban(name = None):
    curs.execute("select * from user where id = '" + pymysql.escape_string(name) + "'")
    rows = curs.fetchall()
    if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
        return render_template('index.html', title = '차단 오류', logo = data['name'], data = '관리자는 차단 할 수 없습니다.')
    else:
        if(request.method == 'POST'):
            if(admincheck() == 1):
                curs.execute("select * from ban where block = '" + pymysql.escape_string(name) + "'")
                row = curs.fetchall()
                if(row):
                    block(name, '해제', getnow(), ip, '')
                    curs.execute("delete from ban where block = '" + pymysql.escape_string(name) + "'")
                else:
                    b = re.search("^([0-9](?:[0-9][0-9])?\.[0-9](?:[0-9][0-9])?)$", name)
                    if(b):
                        block(name, request.form["end"], getnow(), ip, request.form["why"])
                        curs.execute("insert into ban (block, end, why, band) value ('" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(request.form["end"]) + "', '" + pymysql.escape_string(request.form["why"]) + "', 'O')")
                    else:
                        block(name, request.form["end"], getnow(), ip, request.form["why"])
                        curs.execute("insert into ban (block, end, why, band) value ('" + pymysql.escape_string(name) + "', '" + pymysql.escape_string(request.form["end"]) + "', '" + pymysql.escape_string(request.form["why"]) + "', '')")
                conn.commit()
                return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']).replace('/','%2F') + '" />'
            else:
                return render_template('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.')
        else:
            if(admincheck() == 1):
                curs.execute("select * from ban where block = '" + pymysql.escape_string(name) + "'")
                row = curs.fetchall()
                if(row):
                    now = '차단 해제'
                else:
                    b = re.search("^([0-9](?:[0-9][0-9])?\.[0-9](?:[0-9][0-9])?)$", name)
                    if(b):
                        now = '대역 차단'
                    else:
                        now = '차단'
                return render_template('index.html', title = name, page = parse.quote(name).replace('/','%2F'), logo = data['name'], tn = 16, now = now, today = getnow())
            else:
                return render_template('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.')

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
            return render_template('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.')
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
                return render_template('index.html', title = name, page = parse.quote(name).replace('/','%2F'), logo = data['name'], tn = 19, now = '현재 ACL 상태는 ' + now)
            else:
                return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(name).replace('/','%2F') + '" />' 
        else:
            return render_template('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.')

@app.route('/admin/<path:name>', methods=['POST', 'GET'])
def admin(name = None):
    if(request.method == 'POST'):
        if(session.get('Now') == True):
            ip = getip(request)
            curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
            rows = curs.fetchall()
            if(rows):
                if(rows[0]['acl'] == 'owner' or rows[0]['acl'] == 'admin'):
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
                        return render_template('index.html', title = '사용자 오류', logo = data['name'], data = '계정이 없습니다.')
                else:
                    return render_template('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.')
            else:
                return render_template('index.html', title = '권한 오류', logo = data['name'], data = '계정이 없습니다.')
        else:
            return render_template('index.html', title = '권한 오류', logo = data['name'], data = '비 로그인 상태 입니다.')
    else:
        if(session.get('Now') == True):
            ip = getip(request)
            curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
            rows = curs.fetchall()
            if(rows):
                if(rows[0]['acl'] == 'owner'):
                    curs.execute("select * from user where id = '" + pymysql.escape_string(name) + "'")
                    row = curs.fetchall()
                    if(row):
                        if(row[0]['acl'] == 'admin' or row[0]['acl'] == 'owner'):
                            now = '권한 해제'
                        else:
                            now = '권한 부여'
                        return render_template('index.html', title = name, page = parse.quote(name).replace('/','%2F'), logo = data['name'], tn = 18, now = now)
                    else:
                        return render_template('index.html', title = '사용자 오류', logo = data['name'], data = '계정이 없습니다.')
                else:
                    return render_template('index.html', title = '권한 오류', logo = data['name'], data = '권한이 모자랍니다.')
            else:
                return render_template('index.html', title = '권한 오류', logo = data['name'], data = '계정이 없습니다.')
        else:
            return render_template('index.html', title = '권한 오류', logo = data['name'], data = '비 로그인 상태 입니다.')

@app.route('/grammar')
def grammar():
    return render_template('index.html', title = '문법 설명', logo = data['name'], tn = 17)

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
            b = re.search("^([0-9](?:[0-9][0-9])?\.[0-9](?:[0-9][0-9])?)", ip)
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
    
    return render_template('index.html', title = '권한 오류', logo = data['name'], data = end)
   
@app.route('/w/<path:name>/r/<a>/diff/<b>')
def diff(name = None, a = None, b = None):
    curs.execute("select * from history where id = '" + pymysql.escape_string(a) + "' and title = '" + pymysql.escape_string(name) + "'")
    rows = curs.fetchall()
    if(rows):
        curs.execute("select * from history where id = '" + pymysql.escape_string(b) + "' and title = '" + pymysql.escape_string(name) + "'")
        row = curs.fetchall()
        if(row):
            indata = re.sub('<', '&lt;', rows[0]['data'])
            indata = re.sub('>', '&gt;', indata)
            indata = re.sub('"', '&quot;', indata)
            indata = re.sub('\n', '<br>', indata)
            enddata = re.sub('<', '&lt;', row[0]['data'])
            enddata = re.sub('>', '&gt;', enddata)
            enddata = re.sub('"', '&quot;', enddata)
            enddata = re.sub('\n', '<br>', enddata)
            sm = difflib.SequenceMatcher(None, indata, enddata)
            c = show_diff(sm)
            return render_template('index.html', title = 'Diff', logo = data['name'], data = c)
        else:
            return render_template('index.html', title = 'Diff 오류', logo = data['name'], data = '<a href="/w/' + name + '">이 리비전이나 문서가 없습니다.</a>')
    else:
        return render_template('index.html', title = 'Diff 오류', logo = data['name'], data = '<a href="/w/' + name + '">이 리비전이나 문서가 없습니다.</a>')

@app.route('/version')
def version():
    return render_template('index.html', title = '버전', logo = data['name'], tn = 14)

@app.route('/user')
def user():
    ip = getip(request)
    curs.execute("select * from user where id = '" + pymysql.escape_string(ip) + "'")
    rows = curs.fetchall()
    if(getban(ip) == 0):
        if(rows):
            if(rows[0]['acl'] == 'admin' or rows[0]['acl'] == 'owner'):
                acl = '관리자'
            else:
                acl = '유저'
        else:
            acl = '일반'
    else:
        acl = '차단'
    return render_template('index.html', title = '유저 메뉴', logo = data['name'], data = ip + '<br><br><span>권한 상태 : ' + acl + '<br><br><li><a href="/login">로그인</a></li><li><a href="/logout">로그아웃</a></li><li><a href="/register">회원가입</a></li>')

@app.route('/random')
def random():
    curs.execute("select * from data order by rand() limit 1")
    rows = curs.fetchall()
    if(rows):
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(rows[0]['title']).replace('/','%2F') + '" />'
    else:
        return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']).replace('/','%2F') + '" />'

@app.errorhandler(404)
def uncaughtError(error):
    return '<meta http-equiv="refresh" content="0;url=/w/' + parse.quote(data['frontpage']).replace('/','%2F') + '" />'

if __name__ == '__main__':
    app.run(host = '0.0.0.0', port = int(data['port']))
