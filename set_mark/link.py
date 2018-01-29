import sqlite3
import re
from urllib import parse
import hashlib

def url_pas(data):
    return parse.quote(data).replace('/','%2F')

def sha224(data):
    return hashlib.sha224(bytes(data, 'utf-8')).hexdigest()

def link(conn, title, data, num, category, backlink):
    curs = conn.cursor()
    data = data.replace('&#92;', '\\')
    
    m = re.findall("\[\[(분류:(?:(?:(?!\]\]|#).)+))((?:#(?:(?:(?!#|\]\]).)+))+)?\]\]", data)
    for g in m:
        if title != g[0]:
            if num == 1:
                backlink += [[title, g[0], 'cat']]
                
            curs.execute("select title from data where title = ?", [g[0]])
            if curs.fetchall():
                red = ""
            else:
                red = 'class="not_thing"'

            if(category != ''):    
                category += ' / '                

            style = ''
            if g[1]:
                if re.search('#blur', g[1]):
                    style = ' style="filter: blur(3px);" onmouseover="this.style.filter=\'none\';" onmouseout="this.style.filter=\'blur(3px)\';"'
                
            category += '<a ' + red + ' ' + style + '" href="/w/' + url_pas(g[0]) + '">' + re.sub("분류:", "", g[0]) + '</a>'
        
        data = re.sub("\[\[(분류:(?:(?:(?!\]\]|#).)+))((?:#(?:(?:(?!#|\]\]).)+))+)?\]\]", '', data, 1)
    
    test = re.findall('\[\[wiki:([^|\]]+)(?:\|([^\]]+))?\]\]', data)
    for wiki in test:
        if wiki[1]:
            out = wiki[1]
        else:
            out = wiki[0]

        data = re.sub('\[\[wiki:([^|\]]+)(?:\|([^\]]+))?\]\]', '<a id="inside" href="/' + wiki[0] + '">' + out + '</a>', data, 1)
     
    data = re.sub("\[\[(?::(?P<in>(?:분류|파일):(?:(?:(?!\]\]).)*)))\]\]", "[[\g<in>]]", data)

    a = re.findall('\[\[\.\.\/(\|(?:(?!]]).)+)?]]', data)
    for i in a:
        b = re.search('(.*)\/', title)
        if b:
            m = b.groups()
            if i:
                data = re.sub('\[\[\.\.\/(\|((?!]]).)+)?]]', '[[' + m[0] + i + ']]', data, 1)
            else:
                data = re.sub('\[\[\.\.\/(\|((?!]]).)+)?]]', '[[' + m[0] + ']]', data, 1)
        else:
            if i:
                data = re.sub('\[\[\.\.\/(\|((?!]]).)+)?]]', '[[' + title + i + ']]', data, 1)
            else:
                data = re.sub('\[\[\.\.\/(\|((?!]]).)+)?]]', '[[' + title + ']]', data, 1)

    data = re.sub('\[\[(?P<in>\/(?:(?!]]|\|).)+)(?P<out>\|(?:(?:(?!]]).)+))?]]', '[[' + title + '\g<in>\g<out>]]', data)
                
    link = re.compile('\[\[((?:(?!\[\[|\]\]|\|).)*)(?:\|((?:(?!\[\[|\]\]).)*))?\]\]')
    while 1:
        l_d = link.search(data)
        if l_d:
            d = l_d.groups()
            if re.search('^(?:파일|외부):', d[0]):
                width = ''
                height = ''
                align = ''
                span = ['', '']
                
                try:        
                    w_d = re.search('width=([0-9]+(?:[a-z%]+)?)', d[1])
                    if w_d:
                        width = 'width="' + w_d.groups()[0] + '" '
                    
                    h_d = re.search('height=([0-9]+(?:[a-z%]+)?)', d[1])
                    if h_d:
                        height = 'height="' + h_d.groups()[0] + '" '
                        
                    a_d = re.search('align=(center|right)', d[1])
                    if a_d:
                        span[0] = '<span style="display: block; text-align: ' + a_d.groups()[0] + ';">'
                        span[1] = '</span>'
                except:
                    pass
                    
                f_d = re.search('^파일:([^.]+)\.(.+)$', d[0])
                if f_d:
                    if not re.search("^파일:([^\n]*)", title):
                        if num == 1:
                            backlink += [[title, d[0], 'file']]

                    file_name = f_d.groups()

                    curs.execute("select title from data where title = ?", ['파일:' + file_name[0] + '.' + file_name[1]])
                    if not curs.fetchall():
                        img = '<a class="not_thing" href="/w/' + url_pas('파일:' + file_name[0] + '.' + file_name[1]) + '">파일:' + file_name[0] + '.' + file_name[1] + '</a>'
                    else:
                        img = span[0] + '<img src="/image/' + sha224(file_name[0]) + '.' + file_name[1] + '" ' + width + height + '>' + span[1]
                    
                    data = link.sub(img, data, 1)
                else:
                    img = span[0] + '<img src="' + re.sub('^외부:', '', d[0]) + '" ' + width + height + '>' + span[1]
                    data = link.sub(img, data, 1)
                                    
            elif re.search('^https?:\/\/', re.sub('<([^>]*)>', '', d[0])):
                view = d[0]
                try:
                    if re.search('(.+)', d[1]):
                        view = d[1]
                except:
                    pass
                
                data = link.sub('<a class="out_link" rel="nofollow" href="' + re.sub('<([^>]*)>', '', d[0]) + '">' + view + '</a>', data, 1)
            else:
                view = d[0].replace('\\\\', '<slash>').replace('\\', '').replace('<slash>', '\\')
                try:
                    if re.search('(.+)', d[1]):
                        view = d[1].replace('\\\\', '<slash>').replace('\\', '').replace('<slash>', '\\')
                except:
                    pass        
                    
                sh = ''
                s_d = re.search('#((?:(?!x27;|#).)+)$', d[0])
                if s_d:
                    href = re.sub('#((?:(?!x27;|#).)+)$', '', d[0])
                    sh = '#' + s_d.groups()[0]
                else:
                    href = d[0]
                    
                if d[0] == title:
                    data = link.sub('<b>' + view + '</b>', data, 1)
                elif re.search('^#', d[0]):
                    data = link.sub('<a title="' + sh + '" href="' + sh + '">' + view + '</a>', data, 1)
                else:                    
                    a = re.sub('<([^>]*)>', '', href.replace('&#x27;', "'").replace('&quot;', '"').replace('\\\\', '<slash>').replace('\\', '').replace('<slash>', '\\'))
                    
                    if num == 1:
                        backlink += [[title, a, '']]
                    
                    curs.execute("select title from data where title = ?", [a])
                    if not curs.fetchall():
                        no = 'class="not_thing"'
                        
                        if num == 1:
                            backlink += [[title, a, 'no']]
                    else:
                        no = ''
                    
                    data = link.sub('<a ' + no + ' title="' + re.sub('<([^>]*)>', '', href) + sh + '" href="/w/' + url_pas(a) + sh + '">' + view.replace('\\', '\\\\') + '</a>', data, 1)
        else:
            break
            
    data = data.replace('\\', '&#92;')

    return [data, category, backlink]