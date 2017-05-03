from code import *

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

def 세이브마크(데이터):
    데이터 = re.sub("\[date\(now\)\]", 시간(), 데이터)
    if(not re.search("\.", 아이피_확인())):
        이름 = '[[사용자:' + 아이피_확인() + '|' + 아이피_확인() + ']]'
    else:
        이름 = 아이피_확인()
    데이터 = re.sub("\[name\]", 이름, 데이터)

    return 데이터

def HTML_파싱(데이터):
    while(True):
        있나 = re.search("<((div|span|embed|iframe)(?:[^>]*))>", 데이터)
        
        if(있나):
            분리 = 있나.groups()

            if(re.search("<(\/" + 분리[1] + ")>", 데이터)):
                XSS = re.search('src=(?:"|\')http(?:s)?:\/\/([^\/]*)\/(?:[^"\']*)(?:"|\')', 분리[0])
                
                if(XSS):
                    확인 = XSS.groups()
                    
                    if(확인[0] == "www.youtube.com" or 확인[0] == "serviceapi.nmv.naver.com" or 확인[0] == "tv.kakao.com" or 확인[0] == "tvple.com"):
                        임시_저장 = 분리[0]
                    else:
                        임시_저장 = re.sub('src=(?:"|\')([^"\']*)(?:"|\')', '', 분리[0])
                else:
                    임시_저장 = 분리[0]
                
                임시_저장 = re.sub('(?:"|\')', '#.#', 임시_저장)
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

def 역링크_추가(이름, 링크, 값):
    DB_실행("select title from back where title = '" + DB_인코딩(링크) + "' and link = '" + DB_인코딩(이름) + "' and type = '" + 값 + "'")
    있나 = DB_가져오기()
    if(not 있나):
        DB_실행("insert into back (title, link, type) value ('" + DB_인코딩(링크) + "', '" + DB_인코딩(이름) + "',  'include')")
        DB_갱신()

def 분류_추가(이름, 링크):
    DB_실행("select title from cat where title = '" + DB_인코딩(링크) + "' and cat = '" + DB_인코딩(이름) + "'")
    있나 = DB_가져오기()
    if(not 있나):
        DB_실행("insert into cat (title, cat) value ('" + DB_인코딩(링크) + "', '" + DB_인코딩(이름) + "')")
        DB_갱신()

def 나무마크(title, data):
    data = HTML_파싱(data)

    접기_숫자 = 0
    임시_저장 = 중괄호_문법(data, 접기_숫자, False)
    
    data = 임시_저장[0]
    접기_숫자 = 임시_저장[1]
    
    data = re.sub("\[anchor\((?P<in>[^\[\]]*)\)\]", '<span id="\g<in>"></span>', data)
    data = 세이브마크(data)
    
    while(True):
        m = re.search("\[include\(((?:(?!\)\]|,).)*)((?:,\s?(?:[^)]*))+)?\)\]", data)
        if(m):
            results = m.groups()
            if(results[0] == title):
                data = re.sub("\[include\(((?:(?!\)\]|,).)*)((?:,\s?(?:[^)]*))+)?\)\]", "<b>" + results[0] + "</b>", data, 1)
            else:
                DB_실행("select * from data where title = '" + DB_인코딩(results[0]) + "'")
                틀_내용 = DB_가져오기()
                
                역링크_추가(title, results[0], 'include')
                if(틀_내용):                        
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

                    data = re.sub("\[include\(((?:(?!\)\]|,).)*)((?:,\s?(?:[^)]*))+)?\)\]", '\n#nobr#<div>' + 틀_데이터 + '</div>\n#nobr#', data, 1)
                else:
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
            
            역링크_추가(title, results[0], 'redirect')
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
        
    data = re.sub("(\n)(?P<in>\r\n(={1,6})\s?([^=]*)\s?(?:={1,6})(?:\s+)?\n)", "\g<in>", data)
    
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
                분류_추가(title, g[0])
                    
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

            if(c):
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
            else:
                img = re.sub("\.(?P<in>[Jj][Pp][Gg]|[Pp][Nn][Gg]|[Gg][Ii][Ff]|[Jj][Pp][Ee][Gg])", "#\g<in>#", c[0])
                data = re.sub("\[\[파일:((?:(?!\]\]|\?).)*)(?:\?((?:(?!\]\]).)*))?\]\]", "<a href='/w/파일:" + img + "'><img src='/image/" + img + "'></a>", data, 1)

            if(not re.search("^파일:([^\n]*)", title)):
                역링크_추가(title, '파일' + c[0], 'file')            
        else:
            break
    
    data = re.sub("\[br\]",'<br>', data)
    
    while(True):
        문법_컴파일 = re.compile("\[youtube\(((?:(?!,|\)\]).)*)(?:,(?:\s)?)?(?:width=((?:(?!,|\)\]).)*))?(?:,(?:\s)?)?(?:height=((?:(?!,|\)\]).)*))?\)\]")
        m = 문법_컴파일.search(data)
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
                height = result[2]
                width = '560'
            else:
                width = '560'
                height = '315'

            data = 문법_컴파일.sub('<iframe width="' + width + '" height="' + height + '" src="https://www.youtube.com/embed/' + result[0] + '" frameborder="0" allowfullscreen></iframe>', data, 1)
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
                            DB_실행("select title from data where title = '" + DB_인코딩(results[0]) + "'")
                            있나 = DB_가져오기()
                            if(있나):
                                클래스 = ''
                            else:
                                클래스 = 'not_thing'
                                
                            data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a title="' + results[0] + results[1] + '" class="' + 클래스 + '" href="/w/' + URL_인코딩(results[0]) + results[1] + '">' + g + '</a>', data, 1)
                            
                            역링크_추가(title, results[0], '')
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
                            DB_실행("select title from data where title = '" + DB_인코딩(results[0]) + "'")
                            있나 = DB_가져오기()
                            if(있나):
                                클래스 = ''
                            else:
                                클래스 = 'not_thing'

                            data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a title="' + results[0] + '" class="' + 클래스 + '" href="/w/' + URL_인코딩(results[0]) + '">' + results[1] + '</a>', data, 1)

                            역링크_추가(title, results[0], '')
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
                            DB_실행("select title from data where title = '" + DB_인코딩(result[0]) + "'")
                            있나 = DB_가져오기()
                            if(있나):
                                클래스 = ''
                            else:
                                클래스 = 'not_thing'

                            data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a href="/w/' + URL_인코딩(result[0]) + result[1] + '" class="' + 클래스 + '">' + result[0] + result[1] + '</a>', data, 1)
                            
                            역링크_추가(title, result[0], '')
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
                            DB_실행("select title from data where title = '" + DB_인코딩(result[0]) + "'")
                            있나 = DB_가져오기()
                            if(있나):
                                클래스 = ''
                            else:
                                클래스 = 'not_thing'

                            data = re.sub('\[\[(((?!\]\]).)*)\]\]', '<a href="/w/' + URL_인코딩(result[0]) + '" class="' + 클래스 + '">' + result[0] + '</a>', data, 1)
                            
                            역링크_추가(title, result[0], '')
        else:
            break
            
    while(True):
        문법_컴파일 = re.compile("(http(?:s)?:\/\/(?:(?:(?:(?!\.[Jj][Pp][Gg]|\.[Pp][Nn][Gg]|\.[Gg][Ii][Ff]|\.[Jj][Pp][Ee][Gg]|#[Jj][Pp][Gg]#|#[Pp][Nn][Gg]#|#[Gg][Ii][Ff]#|#[Jj][Pp][Ee][Gg]#|<\/(?:[^>]*)>).)*)(?:\.[Jj][Pp][Gg]|\.[Pp][Nn][Gg]|\.[Gg][Ii][Ff]|\.[Jj][Pp][Ee][Gg])))(?:(?:(?:\?)width=((?:[0-9]*)(?:px|%)?))?(?:(?:\?|&)height=((?:[0-9]*)(?:px|%)?))?)")
        m = 문법_컴파일.search(data)
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
                height = result[2]
                width = ''
            else:
                width = ''
                height = ''

            c = result[0]
            c = re.sub("\.(?P<in>[Jj][Pp][Gg]|[Pp][Nn][Gg]|[Gg][Ii][Ff]|[Jj][Pp][Ee][Gg])", "#\g<in>#", c)

            data = 문법_컴파일.sub("<img width='" + width + "' height='" + height + "' src='" + c + "'>", data, 1)
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
                        alltable = 'style="'
                        celstyle = 'style="'
                        rowstyle = 'style="'

                        q = re.search("&lt;table\s?width=((?:(?!&gt;).)*)&gt;", result[1])
                        w = re.search("&lt;table\s?height=((?:(?!&gt;).)*)&gt;", result[1])
                        e = re.search("&lt;table\s?align=((?:(?!&gt;).)*)&gt;", result[1])
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
    
    data = re.sub("(\n#nobr#|#nobr#\n|#nobr#)", "", data)
    data = re.sub('<\/blockquote>((\r)?\n){2}<blockquote>', '</blockquote><br><blockquote>', data)
    data = re.sub('\n', '<br>', data)
    data = re.sub('^<br>', '', data)
    
    return data