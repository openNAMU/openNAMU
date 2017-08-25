from bottle import request, app
from bottle.ext import beaker
from urllib import parse
import json
import sqlite3
import time
import re
import hashlib
import html

json_data = open('set.json').read()
set_data = json.loads(json_data)

conn = sqlite3.connect(set_data['db'] + '.db')
curs = conn.cursor()

session_opts = {
    'session.type': 'file',
    'session.data_dir': './app_session/',
    'session.auto': True
}

app = beaker.middleware.SessionMiddleware(app(), session_opts)

def get_time():
    now = time.localtime()
    date = "%04d-%02d-%02d %02d:%02d:%02d" % (now.tm_year, now.tm_mon, now.tm_mday, now.tm_hour, now.tm_min, now.tm_sec)

    return(date)
    
def ip_check():
    session = request.environ.get('beaker.session')
    if(session.get('Now') == 1):
        ip = format(session['DREAMER'])
    else:
        if(request.environ.get('HTTP_X_FORWARDED_FOR')):
            ip = request.environ.get('HTTP_X_FORWARDED_FOR')
        else:
            ip = request.environ.get('REMOTE_ADDR')

    return(ip)

def url_pas(data):
    return(parse.quote(data).replace('/','%2F'))

def sha224(data):
    return(hashlib.sha224(bytes(data, 'utf-8')).hexdigest())

def savemark(data):
    data = re.sub("\[date\(now\)\]", get_time(), data)
    
    if(not re.search("\.", ip_check())):
        name = '[[사용자:' + ip_check() + '|' + ip_check() + ']]'
    else:
        name = ip_check()
        
    data = re.sub("\[name\]", name, data)

    return(data)

def send_p(d):
    d = html.escape(d)

    js_p = re.compile('javascript:', re.I)
    d = js_p.sub('', d)

    d = re.sub('&lt;a href="(?:[^"]*)"&gt;(?P<in>(?:(?!&lt;).)*)&lt;\/a&gt;', '<a href="' + url_pas('\g<in>') + '">\g<in></a>', d)

    return(d)

def table_p(d, d2):
    alltable = 'style="'
    celstyle = 'style="'
    rowstyle = 'style="'
    row = ''
    cel = ''

    q = re.search("&lt;table\s?width=((?:(?!&gt;).)*)&gt;", d)
    w = re.search("&lt;table\s?height=((?:(?!&gt;).)*)&gt;", d)
    e = re.search("&lt;table\s?align=((?:(?!&gt;).)*)&gt;", d)
    if(q):
        resultss = q.groups()
        alltable += 'width:' + resultss[0] + ';'
    if(w):
        resultss = w.groups()
        alltable += 'height:' + resultss[0] + ';'
    if(e):
        resultss = e.groups()
        if(resultss[0] == 'right'):
            alltable += 'float:right;'
        elif(resultss[0] == 'center'):
            alltable += 'margin:auto;'
            
    ee = re.search("&lt;table\s?textalign=((?:(?!&gt;).)*)&gt;", d)
    if(ee):
        resultss = ee.groups()
        if(resultss[0] == 'right'):
            alltable += 'text-align:right;'
        elif(resultss[0] == 'center'):
            alltable += 'text-align:center;'
    
    r = re.search("&lt;-((?:(?!&gt;).)*)&gt;", d)
    if(r):
        resultss = r.groups()
        cel = 'colspan="' + resultss[0] + '"'
    else:
        cel = 'colspan="' + str(round(len(d2) / 2)) + '"'   

    t = re.search("&lt;\|((?:(?!&gt;).)*)&gt;", d)
    if(t):
        resultss = t.groups()
        row = 'rowspan="' + resultss[0] + '"'

    ba = re.search("&lt;rowbgcolor=(#[0-9a-f-A-F]{6})&gt;", d)
    bb = re.search("&lt;rowbgcolor=(#[0-9a-f-A-F]{3})&gt;", d)
    bc = re.search("&lt;rowbgcolor=(\w+)&gt;", d)
    if(ba):
        resultss = ba.groups()
        rowstyle += 'background:' + resultss[0] + ';'
    elif(bb):
        resultss = bb.groups()
        rowstyle += 'background:' + resultss[0] + ';'
    elif(bc):
        resultss = bc.groups()
        rowstyle += 'background:' + resultss[0] + ';'
        
    z = re.search("&lt;table\s?bordercolor=(#[0-9a-f-A-F]{6})&gt;", d)
    x = re.search("&lt;table\s?bordercolor=(#[0-9a-f-A-F]{3})&gt;", d)
    c = re.search("&lt;table\s?bordercolor=(\w+)&gt;", d)
    if(z):
        resultss = z.groups()
        alltable += 'border:' + resultss[0] + ' 2px solid;'
    elif(x):
        resultss = x.groups()
        alltable += 'border:' + resultss[0] + ' 2px solid;'
    elif(c):
        resultss = c.groups()
        alltable += 'border:' + resultss[0] + ' 2px solid;'
        
    aq = re.search("&lt;table\s?bgcolor=(#[0-9a-f-A-F]{6})&gt;", d)
    aw = re.search("&lt;table\s?bgcolor=(#[0-9a-f-A-F]{3})&gt;", d)
    ae = re.search("&lt;table\s?bgcolor=(\w+)&gt;", d)
    if(aq):
        resultss = aq.groups()
        alltable += 'background:' + resultss[0] + ';'
    elif(aw):
        resultss = aw.groups()
        alltable += 'background:' + resultss[0] + ';'
    elif(ae):
        resultss = ae.groups()
        alltable += 'background:' + resultss[0] + ';'
        
    j = re.search("&lt;bgcolor=(#[0-9a-f-A-F]{6})&gt;", d)
    k = re.search("&lt;bgcolor=(#[0-9a-f-A-F]{3})&gt;", d)
    l = re.search("&lt;bgcolor=(\w+)&gt;", d)
    if(j):
        resultss = j.groups()
        celstyle += 'background:' + resultss[0] + ';'
    elif(k):
        resultss = k.groups()
        celstyle += 'background:' + resultss[0] + ';'
    elif(l):
        resultss = l.groups()
        celstyle += 'background:' + resultss[0] + ';'
        
    aa = re.search("&lt;(#[0-9a-f-A-F]{6})&gt;", d)
    ab = re.search("&lt;(#[0-9a-f-A-F]{3})&gt;", d)
    ac = re.search("&lt;(\w+)&gt;", d)
    if(aa):
        resultss = aa.groups()
        celstyle += 'background:' + resultss[0] + ';'
    elif(ab):
        resultss = ab.groups()
        celstyle += 'background:' + resultss[0] + ';'
    elif(ac):
        resultss = ac.groups()
        celstyle += 'background:' + resultss[0] + ';'
        
    qa = re.search("&lt;width=((?:(?!&gt;).)*)&gt;", d)
    qb = re.search("&lt;height=((?:(?!&gt;).)*)&gt;", d)
    if(qa):
        resultss = qa.groups()
        celstyle += 'width:' + resultss[0] + ';'
    if(qb):
        resultss = qb.groups()
        celstyle += 'height:' + resultss[0] + ';'
        
    i = re.search("&lt;\)&gt;", d)
    o = re.search("&lt;:&gt;", d)
    p = re.search("&lt;\(&gt;",  d)
    if(i):
        celstyle += 'text-align:right;'
    elif(o):
        celstyle += 'text-align:center;'
    elif(p):
        celstyle += 'text-align:left;'
        
    alltable += '"'
    celstyle += '"'
    rowstyle += '"'

    return([alltable, rowstyle, celstyle, row, cel])

def html_pas(data):
    data = re.sub("%phtml%(?P<in>(?:\/)?(?:a|div|span|embed|iframe)(?:\s[^%]*)?)%phtml%", "<\g<in>>", data)   
    while(1):
        y = re.search("<((div|span|embed|iframe)(?:\s[^>]*))>", data)
        if(y):
            b = y.groups()

            if(re.search("<(\/" + b[1] + ")>", data)):
                xss_test = re.search('src=(?:"|\')?(http(s)?:\/\/([^\/]*)\/(?:[^"\' ]*))(?:"|\')?', b[0])
                
                if(xss_test):
                    check = xss_test.groups()
                    
                    if(check[2] == "www.youtube.com" or check[2] == "serviceapi.nmv.naver.com" or check[2] == "tv.kakao.com" or check[2] == "tvple.com"):
                        a = b[0]
                    else:
                        a = re.sub('src=(?:"|\')([^"\']*)(?:"|\')', '', b[0])
                else:
                    a = b[0]

                try:
                    if(check[1] != None):
                        data = re.sub("<((?:\/)?" + b[1] + "(?:\s[^>]*))>", "%phtml%" + a + "%phtml%", data, 1)
                        data = re.sub("<\/" + b[1] + ">", "%phtml%/" + b[1] + "%phtml%", data, 1)
                    else:
                        data = re.sub("<((?:\/)?" + b[1] + "(?:\s[^>]*))>", "[[" + check[0] + "]]", data, 1)
                        data = re.sub("<\/" + b[1] + ">", "", data, 1)
                except:
                    data = re.sub("<((?:\/)?" + b[1] + "(?:\s[^>]*))>", "%phtml%" + a + "%phtml%", data, 1)
                    data = re.sub("<\/" + b[1] + ">", "%phtml%/" + b[1] + "%phtml%", data, 1)
            else:
                data = re.sub("<((?:\/)?" + b[1] + "(?:\s[^>]*))>", '&lt;' + b[0] + '&gt;', data, 1)
                
                break
        else:
            break

    data = html.escape(data)

    js_p = re.compile('javascript:', re.I)
    data = js_p.sub('', data)
    
    data = re.sub("%phtml%(?P<in>(?:\/)?(?:a|div|span|embed|iframe)(?:\s[^%]*)?)%phtml%", "<\g<in>>", data)
    return(data)
    
def mid_pas(data, fol_num, include):
    while(1):
        com = re.compile("{{{((?:(?!{{{)(?!}}}).)*)}}}", re.DOTALL)
        y = com.search(data)
        
        if(y):
            a = y.groups()
            
            big_a = re.compile("^\+([1-5])\s(.*)$", re.DOTALL)
            big = big_a.search(a[0])
            
            small_a = re.compile("^\-([1-5])\s(.*)$", re.DOTALL)
            small = small_a.search(a[0])
            
            color_a = re.compile("^(#[0-9a-f-A-F]{6})\s(.*)$", re.DOTALL)
            color = color_a.search(a[0])
            
            color_b = re.compile("^(#[0-9a-f-A-F]{3})\s(.*)$", re.DOTALL)
            color_2 = color_b.search(a[0])
            
            color_c = re.compile("^#(\w+)\s(.*)$", re.DOTALL)
            color_3 = color_c.search(a[0])
            
            back_a = re.compile("^@([0-9a-f-A-F]{6})\s(.*)$", re.DOTALL)
            back = back_a.search(a[0])
            
            back_b = re.compile("^@([0-9a-f-A-F]{3})\s(.*)$", re.DOTALL)
            back_2 = back_b.search(a[0])
            
            back_c = re.compile("^@(\w+)\s(.*)$", re.DOTALL)
            back_3 = back_c.search(a[0])
            
            include_out_a = re.compile("^#!noin\s(.*)$", re.DOTALL)
            include_out = include_out_a.search(a[0])
            
            div_a = re.compile("^#!wiki\sstyle=(?:&quot;|&apos;)((?:(?!&quot;|&apos;).)*)(?:&quot;|&apos;)\r\n(.*)$", re.DOTALL)
            div = div_a.search(a[0])
            
            html_a = re.compile("^#!html\s(.*)$", re.DOTALL)
            html = html_a.search(a[0])
            
            fol_a = re.compile("^#!folding\s((?:(?!\n).)*)\n?\s\n(.*)$", re.DOTALL)
            fol = fol_a.search(a[0])
            
            syn_a = re.compile("^#!syntax\s([^\n]*)\r\n(.*)$", re.DOTALL)
            syn = syn_a.search(a[0])
            
            if(big):
                result = big.groups()
                data = com.sub('<span class="font-size-' + result[0] + '">' + result[1] + '</span>', data, 1)
            elif(small):
                result = small.groups()
                data = com.sub('<span class="font-size-small-' + result[0] + '">' + result[1] + '</span>', data, 1)
            elif(color):
                result = color.groups()
                data = com.sub('<span style="color:' + result[0] + '">' + result[1] + '</span>', data, 1)
            elif(color_2):
                result = color_2.groups()
                data = com.sub('<span style="color:' + result[0] + '">' + result[1] + '</span>', data, 1)
            elif(color_3):
                result = color_3.groups()
                data = com.sub('<span style="color:' + result[0] + '">' + result[1] + '</span>', data, 1)
            elif(back):
                result = back.groups()
                data = com.sub('<span style="background:#' + result[0] + '">' + result[1] + '</span>', data, 1)
            elif(back_2):
                result = back_2.groups()
                data = com.sub('<span style="background:#' + result[0] + '">' + result[1] + '</span>', data, 1)
            elif(back_3):
                result = back_3.groups()
                data = com.sub('<span style="background:' + result[0] + '">' + result[1] + '</span>', data, 1)
            elif(div):
                result = div.groups()
                data = com.sub('<div style="' + result[0] + '">' + result[1] + '</div>', data, 1)
            elif(html):
                result = html.groups()
                data = com.sub(result[0], data, 1)
            elif(fol):
                result = fol.groups()
                data = com.sub( "<div> \
                                    " + result[0] + " \
                                    <div id='folding_" + str(fol_num + 1) + "' style='display: inline-block;'> \
                                        [<a href='javascript:void(0);' onclick='folding(" + str(fol_num + 1) + "); folding(" + str(fol_num + 2) + "); folding(" + str(fol_num) + ");'>펼치기</a>] \
                                    </div> \
                                    <div id='folding_" + str(fol_num + 2) + "' style='display: none;'> \
                                        [<a href='javascript:void(0);' onclick='folding(" + str(fol_num + 1) + "); folding(" + str(fol_num + 2) + "); folding(" + str(fol_num) + ");'>접기</a>] \
                                    </div> \
                                    <div id='folding_" + str(fol_num) + "' style='display: none;'> \
                                        <br> \
                                        " + result[1] + " \
                                    </div> \
                                </div>", data, 1)
                
                fol_num += 3
            elif(syn):
                result = syn.groups()
                data = com.sub('<pre id="syntax"><code class="' + result[0] + '">' + re.sub('\r\n', '<isbr>', re.sub(' ', '<space>', result[1])) + '</code></pre>', data, 1)
            elif(html):
                result = html.groups()
                data = com.sub(result[0], data, 1)
            elif(include_out):
                if(include == 1):
                    data = com.sub("", data, 1)
                else:
                    result = include_out.groups()
                    data = com.sub(result[0], data, 1)
            else:
                data = com.sub('<code>' + a[0] + '</code>', data, 1)
        else:
            break
            
    while(1):
        com = re.compile("<code>(((?!<\/code>).)*)<\/code>", re.DOTALL)
        y = com.search(data)
        if(y):
            a = y.groups()
            
            mid_data = re.sub("<\/span>", "}}}", a[0])
            mid_data = re.sub("<\/div>", "}}}", mid_data)
            mid_data = re.sub('<span class="font\-size\-(?P<in>[1-6])">', "{{{+\g<in> ", mid_data)
            mid_data = re.sub('<span class="font\-size\-small\-(?P<in>[1-6])">', "{{{-\g<in> ", mid_data)
            mid_data = re.sub('<span style="color:(?:#)?(?P<in>[^"]*)">', "{{{#\g<in> ", mid_data)
            mid_data = re.sub('<span style="background:(?:#)?(?P<in>[^"]*)">', "{{{@\g<in> ", mid_data)
            mid_data = re.sub('<div style="(?P<in>[^"]*)">', "{{{#!wiki style=\"\g<in>\"\n", mid_data)
            mid_data = re.sub("(?P<in>.)", "<nowiki>\g<in></nowiki>", mid_data)
            
            data = com.sub(mid_data, data, 1)
        else:
            break
            
    return([data, fol_num])

def toc_pas(data, title):
    i = [0, 0, 0, 0, 0, 0, 0]
    last = 0
    span = ''
    rtoc = '<div id="toc"><span id="toc-name">목차</span><br><br>'
    while(1):
        i[0] += 1
        m = re.search('(={1,6})\s?([^=]*)\s?(?:={1,6})(?:\s+)?\n', data)
        if(m):
            result = m.groups()
            wiki = len(result[0])
            if(last < wiki):
                last = wiki
            else:
                last = wiki
                if(wiki == 1):
                    i[2] = 0
                    i[3] = 0
                    i[4] = 0
                    i[5] = 0
                    i[6] = 0
                elif(wiki == 2):
                    i[3] = 0
                    i[4] = 0
                    i[5] = 0
                    i[6] = 0
                elif(wiki == 3):
                    i[4] = 0
                    i[5] = 0
                    i[6] = 0
                elif(wiki == 4):
                    i[5] = 0
                    i[6] = 0
                elif(wiki == 5):
                    i[6] = 0
            if(wiki == 1):
                i[1] += 1
            elif(wiki == 2):
                i[2] += 1
            elif(wiki == 3):
                i[3] += 1
            elif(wiki == 4):
                i[4] += 1
            elif(wiki == 5):
                i[5] += 1
            else:
                i[6] += 1

            toc = str(i[1]) + '.' + str(i[2]) + '.' + str(i[3]) + '.' + str(i[4]) + '.' + str(i[5]) + '.' + str(i[6]) + '.'

            toc = re.sub("(?P<in>[0-9]0(?:[0]*)?)\.", '\g<in>#.', toc)

            toc = re.sub("0\.", '', toc)
            toc = re.sub("#\.", '.', toc)
            toc = re.sub("\.$", '', toc)

            test = re.search('([0-9]*)(\.([0-9]*))?(\.([0-9]*))?(\.([0-9]*))?(\.([0-9]*))?', toc)
            if(test):
                g = test.groups()
                if(g[4]):
                    span = '<span id="out"></span>' * 4
                elif(g[3]):
                    span = '<span id="out"></span>' * 3
                elif(g[2]):
                    span = '<span id="out"></span>' * 2
                elif(g[1]):
                    span = '<span id="out"></span>'
                else:
                    span = ''

            rtoc += span + '<a href="#s-' + toc + '">' + toc + '</a>. ' + result[1] + '<br>'

            c = re.sub(" $", "", result[1])
            d = c
            c = re.sub("\[\[(([^|]*)\|)?(?P<in>[^\]]*)\]\]", "\g<in>", c)

            data = re.sub('(={1,6})\s?([^=]*)\s?(?:={1,6})(?:\s+)?\n', '<tablenobr><h' + str(wiki) + ' id="' + c + '"><a href="#toc" id="s-' + toc + '">' + toc + '.</a> ' + d + ' <span style="font-size:11px;">[<a href="/edit/' + url_pas(title) + '/section/' + str(i[0]) + '">편집</a>]</span></h' + str(wiki) + '>', data, 1)
        else:
            rtoc += '</div>'
            
            break
    
    data = re.sub("\[목차\]", rtoc, data)

    return(data)

def backlink_plus(name, link, backtype, num):
    if(num == 1):       
        curs.execute("select title from back where title = ? and link = ? and type = ?", [link, name, backtype])
        y = curs.fetchall()
        if(not y):
            curs.execute("insert into back (title, link, type) values (?, ?,  ?)", [link, name, backtype])
            conn.commit()

def cat_plus(name, link, num):
    if(num == 1):        
        curs.execute("select title from cat where title = ? and cat = ?", [link, name])
        y = curs.fetchall()
        if(not y):
            curs.execute("insert into cat (title, cat) values (?, ?)", [link, name])
            conn.commit()

def namumark(title, data, num):    
    data = html_pas(data)

    b = 0
    a = mid_pas(data, b, 0)
    
    data = a[0]
    b = a[1]
    
    data = re.sub("\[anchor\((?P<in>[^\[\]]*)\)\]", '<span id="\g<in>"></span>', data)
    data = savemark(data)
    
    include = re.compile("\[include\(((?:(?!\)\]|,).)*)((?:,\s?(?:[^)]*))+)?\)\]")
    while(1):
        m = include.search(data)
        if(m):
            results = m.groups()
            if(results[0] == title):
                data = include.sub("<b>" + results[0] + "</b>", data, 1)
            else:
                curs.execute("select data from data where title = ?", [results[0]])
                in_con = curs.fetchall()
                
                backlink_plus(title, results[0], 'include', num)
                if(in_con):                        
                    in_data = in_con[0][0]
                    in_data = include.sub("", in_data)
                    
                    in_data = html_pas(in_data)
                    in_data = mid_pas(in_data, b, 1)[0]
                    
                    if(results[1]):
                        a = results[1]
                        while(1):
                            g = re.search("([^= ,]*)\=([^,]*)", a)
                            if(g):
                                result = g.groups()
                                in_data = re.sub("@" + result[0] + "@", result[1], in_data)
                                a = re.sub("([^= ,]*)\=([^,]*)", "", a, 1)
                            else:
                                break       

                    in_data = toc_pas(in_data, results[0])
                                
                    data = include.sub('\n<nobr><a href="/w/' + url_pas(results[0]) + '">[' + results[0] + ' 이동]</a><div>' + in_data + '</div><nobr>\n', data, 1)
                else:
                    data = include.sub("<a class=\"not_thing\" href=\"" + url_pas(results[0]) + "\">" + results[0] + "</a>", data, 1)
        else:
            break
    
    while(1):
        m = re.search('^#(?:redirect|넘겨주기) ([^\n]*)$', data)
        if(m):
            results = m.groups()
            g = re.sub("\\\#", "<sharp>", results[0])
            aa = re.search("^([^\n]*)(#(?:[^\n]*))$", g)
            if(aa):
                results = aa.groups()
                nosharp = re.sub("<sharp>", "#", results[0])
                data = re.sub('^#(?:redirect|넘겨주기) ([^\n]*)$', '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(nosharp) + '/from/' + url_pas(title) + results[1] + '" />', data, 1)
            else:
                nosharp = re.sub("<sharp>", "#", g)
                data = re.sub('^#(?:redirect|넘겨주기) ([^\n]*)$', '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(nosharp) + '/from/' + url_pas(title) + '" />', data, 1)
            
            backlink_plus(title, results[0], 'redirect', num)
        else:
            break
            
    data = '\n' + data + '\n'
    data = re.sub("\[nicovideo\((?P<in>[^,)]*)(?:(?:,(?:[^,)]*))+)?\)\]", "[[http://embed.nicovideo.jp/watch/\g<in>]]", data)
    
    while(1):
        m = re.search("\n&gt;\s?((?:[^\n]*)(?:(?:(?:(?:\n&gt;\s?)(?:[^\n]*))+)?))", data)
        if(m):
            result = m.groups()
            blockquote = result[0]
            blockquote = re.sub("\n&gt;\s?", "\n", blockquote)
            data = re.sub("\n&gt;\s?((?:[^\n]*)(?:(?:(?:(?:\n&gt;\s?)(?:[^\n]*))+)?))", "\n<blockquote>" + blockquote + "</blockquote>", data, 1)
        else:
            break
    
    if(not re.search('\[목차\]', data)):
        if(not re.search('\[목차\(없음\)\]', data)):
            data = re.sub("(?P<in>(={1,6})\s?([^=]*)\s?(?:={1,6})(?:\s+)?\n)", "[목차]\n\g<in>", data, 1)
        else:
            data = re.sub("\[목차\(없음\)\]", "", data)
        
    data = re.sub("(\n)(?P<in>\r\n(={1,6})\s?([^=]*)\s?(?:={1,6})(?:\s+)?\n)", "\g<in>", data)
    
    data = toc_pas(data, title)
    
    category = ''
    while(1):
        m = re.search("\[\[(분류:(?:(?:(?!\]\]).)*))\]\]", data)
        if(m):
            g = m.groups()
            
            if(title != g[0]):
                cat_plus(title, g[0], num)
                    
                if(category == ''):
                    curs.execute("select title from data where title = ?", [g[0]])
                    exists = curs.fetchall()
                    if(exists):
                        red = ""
                    else:
                        red = 'class="not_thing"'
                        
                    category += '<a ' + red + ' href="/w/' + url_pas(g[0]) + '">' + re.sub("분류:", "", g[0]) + '</a>'
                else:
                    curs.execute("select title from data where title = ?", [g[0]])
                    exists = curs.fetchall()
                    if(exists):
                        red = ""
                    else:
                        red = 'class="not_thing"'
                        
                    category += ' / ' + '<a ' + red + ' href="/w/' + url_pas(g[0]) + '">' + re.sub("분류:", "", g[0]) + '</a>'
            
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
    
    data = re.sub('{{\|(?P<in>(?:(?:(?:(?!\|}}).)*)(?:\n?))+)\|}}', '<table> \
                                                                        <tbody> \
                                                                            <tr> \
                                                                                <td> \
                                                                                    \g<in> \
                                                                                </td> \
                                                                            </tr> \
                                                                        </tbody> \
                                                                    </table>', data)
    
    data = re.sub('\[ruby\((?P<in>[^\,]*)\,\s?(?P<out>[^\)]*)\)\]', '<ruby> \
                                                                        \g<in> \
                                                                        <rp>(</rp> \
                                                                        <rt>\g<out></rt> \
                                                                        <rp>)</rp> \
                                                                    </ruby>', data)
    
    data = re.sub("##\s?(?P<in>[^\n]*)\n", "<div style='display:none;'>\g<in></div>", data)
    
    test = re.findall('\[\[wiki:([^|\]]+)(?:\|([^\]]+))?\]\]', data)
    if(test):
        for wiki in test:
            if(wiki[1]):
                data = re.sub('\[\[wiki:([^|\]]+)(?:\|([^\]]+))?\]\]', '<a id="inside" href="/' + wiki[0] + '">' + wiki[1] + '</a>', data, 1)
            else:
                data = re.sub('\[\[wiki:([^|\]]+)(?:\|([^\]]+))?\]\]', '<a id="inside" href="/' + wiki[0] + '">' + wiki[0] + '</a>', data, 1)
    
    while(1):
        m = re.search("\[\[(파일|외부):((?:(?!\]\]|\|).)*)(?:\|((?:(?!\]\]).)*))?\]\]", data)
        if(m):
            img_d = m.groups()

            if(img_d):
                if(img_d[0] == '파일' and not re.search("^파일:([^\n]*)", title)):
                    backlink_plus(title, '파일:' + img_d[1], 'file', num)

                width = ''
                height = ''
                if(img_d[2]):
                    width_r = re.search("width=([^ \n&]*)", img_d[2])
                    height_r = re.search("height=([^ \n&]*)", img_d[2])

                    if(width_r):
                        width_d = width_r.groups()
                        width = width_d[0]

                    if(height_r):
                        height_d = height_r.groups()
                        height = height_d[0]

                if(img_d[0] == '파일'):
                    data = re.sub("\[\[(파일|외부):((?:(?!\]\]|\|).)*)(?:\|((?:(?!\]\]).)*))?\]\]", '<a href="/w/파일:' + img_d[1] + '"><img src="/image/' + sha224(img_d[1]) + '" width="' + width + '" height="' + height + '"></a>', data, 1)
                else:
                    data = re.sub("\[\[(파일|외부):((?:(?!\]\]|\|).)*)(?:\|((?:(?!\]\]).)*))?\]\]", '<img src="' + img_d[1] + '" width="' + width + '" height="' + height + '">', data, 1)
            else:
                break            
        else:
            break
    
    data = re.sub("\[br\]",'<br>', data)
    
    while(1):
        com = re.compile("\[youtube\(([^, )]*)(,[^)]*)?\)\]")
        m = com.search(data)
        if(m):
            src = ''
            width = '560'
            height = '315'
            time = '0'
            
            result = m.groups()
            if(result[0]):
                yudt = re.search('(?:\?v=(.*)|\/([^/?]*)|^([a-zA-Z0-9\-]*))$', result[0])
                if(yudt):
                    if(yudt.groups()[0]):
                        src = yudt.groups()[0]
                    elif(yudt.groups()[1]):
                        src = yudt.groups()[1]
                    elif(yudt.groups()[2]):
                        src = yudt.groups()[2]
                else:
                    src = ''
                    
            if(result[1]):
                mdata = re.search('width=([0-9]*)', result[1])
                if(mdata):
                    width = mdata.groups()[0]
                
                mdata = re.search('height=([0-9]*)', result[1])
                if(mdata):
                    height = mdata.groups()[0]
                    
                mdata = re.search('time=([0-9]*)', result[1])
                if(mdata):
                    time = mdata.groups()[0]

            data = com.sub('<iframe width="' + width + '" height="' + height + '" src="https://www.youtube.com/embed/' + src + '?start=' + time + '" frameborder="0" allowfullscreen></iframe><br>', data, 1)
        else:
            break
     
    data = re.sub("\[\[(?::(?P<in>(?:분류|파일):(?:(?:(?!\]\]).)*)))\]\]", "[[\g<in>]]", data)

    a = re.findall('\[\[\.\.\/(\|(?:[^\]]*))?\]\]', data)
    for i in a:
        b = re.search('(.*)\/', title)
        if(b):
            m = b.groups()
            if(i):
                data = re.sub('\[\[\.\.\/(\|(?:[^\]]*))?\]\]', '[[' + m[0] + i + ']]', data, 1)
            else:
                data = re.sub('\[\[\.\.\/(\|(?:[^\]]*))?\]\]', '[[' + m[0] + ']]', data, 1)
        else:
            if(i):
                data = re.sub('\[\[\.\.\/(\|(?:[^\]]*))?\]\]', '[[' + title + i + ']]', data, 1)
            else:
                data = re.sub('\[\[\.\.\/(\|(?:[^\]]*))?\]\]', '[[' + title + ']]', data, 1)

    data = re.sub('\[\[(?P<in>\/[^\]|]*)(?P<out>\|(?:[^\]]*))?\]\]', '[[' + title + '\g<in>\g<out>]]', data)
                
    while(1):
        m = re.search("\[\[(((?!\]\]).)*)\]\]", data)
        if(m):
            result = m.groups()
            rep = result[0]
            rep = re.sub("\\\#", "<sharp>", rep)
            a = re.search("^((?:(?!\|).)*)(?:\|(.*))?$", rep)
            
            results = a.groups()
            
            aa = re.search("^([^#]*)(#(?:.*))?$", results[0])
            if(results[1]):
                g = re.sub("<sharp>", "#", results[1])
            else:
                g = re.sub("<sharp>", "#", results[0])
            
            results = aa.groups()
            
            if(not results[1]):
                sharp = ''
            else:
                sharp = results[1]
                
            b = re.search("^http(?:s)?:\/\/", results[0])
            if(b):
                comp = re.compile("(?:\.(?:jpg|png|gif|jpeg))", re.I)
                c = comp.search(results[0])
                if(c):                    
                    data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a class="out_link" rel="nofollow" target="_blank" href="' + results[0] + sharp + '">' + g + '</a>', data, 1)
                else:
                    data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a class="out_link" rel="nofollow" target="_blank" href="' + results[0] + sharp + '">' + g + '</a>', data, 1)
            else:
                if(results[0] == title):
                    data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<b>' + g + '</b>', data, 1)
                else:
                    nosharp = re.sub("<sharp>", "#", results[0])
                    
                    curs.execute("select title from data where title = ?", [nosharp])
                    y = curs.fetchall()
                    if(y):
                        clas = ''
                    else:
                        clas = 'not_thing'
                    
                    
                    data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a title="' + re.sub('#', '\#', nosharp) + sharp + '" class="' + clas + '" href="/w/' + url_pas(nosharp) + sharp + '">' + g + '</a>', data, 1)

                    backlink_plus(title, results[0], '', num)
        else:
            break
            
    while(1):
        m = re.search("((?:(?:( +)\*\s(?:[^\n]*))\n?)+)", data)
        if(m):
            result = m.groups()
            end = str(result[0])

            while(1):
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
    
    now_time = get_time()
    data = re.sub('\[date\]', now_time, data)
    
    time_data = re.search('^([0-9]{4})-([0-9]{2})-([0-9]{2})', now_time)
    time = time_data.groups()
    
    age_data = re.findall('\[age\(([0-9]{4})-([0-9]{2})-([0-9]{2})\)\]', data)
    for age in age_data:
        year = int(time[0]) - int(age[0])
        if(age[1] > time[1]):
            year -= 1
        elif(age[1] == time[1]):
            if(age[2] > time[2]):
                year -= 1
                
        data = re.sub('\[age\(([0-9]{4})-([0-9]{2})-([0-9]{2})\)\]', str(year), data, 1)
    
    data = re.sub("-{4,11}", "<hr>", data)
    
    while(1):
        b = re.search("(<\/h[0-9]>|\n)( +)", data)
        if(b):
            result = b.groups()
            up = re.sub(' ', '<span id="in"></span>', result[1])

            if(re.search('<\/h[0-9]>', result[0])):
                data = re.sub("(?P<in>\/h[0-9]>)( +)", '\g<in>' + up, data, 1)
            else:
                data = re.sub("(?:\n)( +)", '<br>' + up, data, 1)
        else:
            break
    
    a = 1
    tou = "<hr id='footnote'><div><br>"
    namu = []
    while(1):
        b = re.search("\[\*([^\s]*)(?:\s(((?!\[|\]).)*))?\]", data)
        if(b):
            results = b.groups()
            if(not results[1] and results[0]):
                i = 0
                
                while(1):
                    try:
                        if(namu[i] == results[0]):
                            none_this = 0

                            break
                        else:
                            i += 2
                    except:
                        none_this = 1

                        break
                        
                if(none_this == 0):
                    data = re.sub("\[\*([^\s]*)(?:\s(((?!\[|\]).)*))?\]", "<sup><a title=\"" + namu[i + 1] + "\" id=\"rfn-" + str(a) + "\" href=\"#fn-" + str(a) + "\">[" + results[0] + "]</a></sup>", data, 1)
                else:
                    data = re.sub("\[\*([^\s]*)(?:\s(((?!\[|\]).)*))?\]", "<sup><a id=\"rfn-" + str(a) + "\" href=\"#fn-" + str(a) + "\">[" + results[0] + "]</a></sup>", data, 1)
            elif(results[0]):                
                namu += [results[0]]
                namu += [results[1]]
                
                c = results[1]
                c = re.sub("<(?:[^>]*)>", '', c)

                tou += "<span id='footnote-list'><a href=\"#rfn-" + str(a) + "\" id=\"fn-" + str(a) + "\">[" + results[0] + "]</a> " + results[1] + "</span><br>"
                data = re.sub("\[\*([^\s]*)(?:\s(((?!\[|\]).)*))?\]", "<sup><a title=\"" + c + "\" id=\"rfn-" + str(a) + "\" href=\"#fn-" + str(a) + "\">[" + results[0] + "]</a></sup>", data, 1)     

                a += 1
            else:
                c = results[1]
                c = re.sub("<(?:[^>]*)>", '', c)
                
                tou += "<span id='footnote-list'><a href=\"#rfn-" + str(a) + "\" id=\"fn-" + str(a) + "\">[" + str(a) + "]</a> " + results[1] + "</span><br>"
                data = re.sub("\[\*([^\s]*)(?:\s(((?!\[|\]).)*))?\]", '<sup><a title="' + c + '" id="rfn-' + str(a) + '" href="#fn-' + str(a) + '">[' + str(a) + ']</a></sup>', data, 1)

                a += 1
        else:
            tou += '</div>'

            if(tou == "<hr id='footnote'><div><br></div>"):
                tou = ""

            break
    
    data = re.sub("\[각주\](?:(?:<br>| |\r|\n)+)?$", "", data)
    data = re.sub("(?:(?:<br>| |\r|\n)+)$", "", data)
    data = re.sub("\[각주\]", "<br>" + tou, data)
    data += tou
    
    if(category):
        data += '<div id="cate">분류: ' + category + '</div>'
    
    data = re.sub("(?:\|\|\r\n)", "#table#<tablenobr>", data)
        
    while(1):
        y = re.search("(\|\|(?:(?:(?:(?:(?!\|\|).)*)(?:\n?))+))", data)
        if(y):
            a = y.groups()
            
            mid_data = re.sub("\|\|", "#table#", a[0])
            mid_data = re.sub("\r\n", "<br>", mid_data)
            
            data = re.sub("(\|\|((?:(?:(?:(?!\|\|).)*)(?:\n?))+))", mid_data, data, 1)
        else:
            break
            
    data = re.sub("#table#", "||", data)
    data = re.sub("<tablenobr>", "\r\n", data)
    
    while(1):
        m = re.search("(\|\|(?:(?:(?:.*)\n?)\|\|)+)", data)
        if(m):
            results = m.groups()
            table = results[0]
            while(1):
                a = re.search("^(\|\|(?:(?:\|\|)+)?)((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)?", table)
                if(a):
                    row = ''
                    cel = ''
                    celstyle = ''
                    rowstyle = ''
                    alltable = ''
                    table_d = ''

                    result = a.groups()
                    if(result[1]):
                        table_d = table_p(result[1], result[0])
                        alltable = table_d[0]
                        rowstyle = table_d[1]
                        celstyle = table_d[2]
                        row = table_d[3]
                        cel = table_d[4]
                            
                        table = re.sub("^(\|\|(?:(?:\|\|)+)?)((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)?",   "<table " + alltable + "> \
                                                                                                        <tbody> \
                                                                                                            <tr " + rowstyle + "> \
                                                                                                                <td " + cel + " " + row + " " + celstyle + ">", table, 1)
                    else:
                        cel = 'colspan="' + str(round(len(result[0]) / 2)) + '"'
                        table = re.sub("^(\|\|(?:(?:\|\|)+)?)((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)?",   "<table> \
                                                                                                        <tbody> \
                                                                                                            <tr> \
                                                                                                                <td " + cel + ">", table, 1)
                else:
                    break
                    
            table = re.sub("\|\|$",                 "</td> \
                                                </tr> \
                                            </tbody> \
                                        </table>", table)
            
            while(1):
                b = re.search("\|\|\r\n(\|\|(?:(?:\|\|)+)?)((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)?", table)
                if(b):
                    row = ''
                    cel = ''
                    celstyle = ''
                    rowstyle = ''
                    table_d = ''

                    result = b.groups()
                    if(result[1]):
                        table_d = table_p(result[1], result[0])
                        rowstyle = table_d[1]
                        celstyle = table_d[2]
                        row = table_d[3]
                        cel = table_d[4]
                        
                        table = re.sub("\|\|\r\n(\|\|(?:(?:\|\|)+)?)((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)?",        "</td> \
                                                                                                            </tr> \
                                                                                                            <tr " + rowstyle + "> \
                                                                                                                <td " + cel + " " + row + " " + celstyle + ">", table, 1)
                    else:
                        cel = 'colspan="' + str(round(len(result[0]) / 2)) + '"'
                        table = re.sub("\|\|\r\n(\|\|(?:(?:\|\|)+)?)((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)?",        "</td> \
                                                                                                            </tr> \
                                                                                                            <tr> \
                                                                                                                <td " + cel + ">", table, 1)
                else:
                    break

            while(1):
                c = re.search("(\|\|(?:(?:\|\|)+)?)((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)?", table)
                if(c):
                    row = ''
                    cel = ''
                    celstyle = ''
                    table_d = ''

                    result = c.groups()
                    if(result[1]):
                        table_d = table_p(result[1], result[0])
                        celstyle = table_d[2]
                        row = table_d[3]
                        cel = table_d[4]

                        table = re.sub("(\|\|(?:(?:\|\|)+)?)((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)?",    "</td> \
                                                                                                    <td " + cel + " " + row + " " + celstyle + ">", table, 1)
                    else:
                        cel = 'colspan="' + str(round(len(result[0]) / 2)) + '"'
                        table = re.sub("(\|\|(?:(?:\|\|)+)?)((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)?",    "</td> \
                                                                                                    <td " + cel + ">", table, 1)
                else:
                    break
            
            data = re.sub("(\|\|(?:(?:(?:.*)\n?)\|\|)+)", table, data, 1)
        else:
            break
            
    data = re.sub("\r\n(?P<in><h[0-6])", "\g<in>", data)
    data = re.sub("(\n<nobr>|<nobr>\n|<nobr>)", "", data)
    data = re.sub("<nowiki>(?P<in>.)<\/nowiki>", "\g<in>", data)
    data = re.sub("<space>", " ", data)

    data = re.sub('<\/blockquote>(?:(?:\r)?\n){2}<blockquote>', '</blockquote><blockquote>', data)
    data = re.sub('<\/blockquote>(?:(?:\r)?\n)<br><blockquote>', '</blockquote><blockquote>', data)
    data = re.sub('\n', '<br>', data)
    data = re.sub('<isbr>', '\r\n', data)
    data = re.sub('^<br>', '', data)

    data += "<script> \
                function folding(num) { \
                    var fol = document.getElementById('folding_' + num); \
                    if(fol.style.display == 'inline-block' || fol.style.display == 'block') { \
                        fol.style.display = 'none'; \
                    } else { \
                        if(num % 3 == 0) { \
                            fol.style.display = 'block'; \
                        } else { \
                            fol.style.display = 'inline-block'; \
                        } \
                    } \
                } \
            </script>"
    
    
    return(data)
