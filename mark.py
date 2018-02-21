from set_mark.start import start
from set_mark.mid_pas import mid_pas
from set_mark.html_pas import html_pas
from set_mark.include_pas import include_pas
from set_mark.macro import macro
from set_mark.redirect_pas import redirect_pas
from set_mark.blockquote import blockquote
from set_mark.toc_pas import toc_pas
from set_mark.text_help import text_help
from set_mark.link import link
from set_mark.indent import indent
from set_mark.footnote import footnote
from set_mark.table import table
from set_mark.end import end

import re
import html
import sqlite3
from urllib import parse
import time
import threading

def send_p(d):
    d = html.escape(d)

    js_p = re.compile('javascript:', re.I)
    d = js_p.sub('', d)

    d = re.sub('&lt;a href="(?:[^"]*)"&gt;(?P<in>(?:(?!&lt;).)*)&lt;\/a&gt;', '<a href="' + url_pas('\g<in>') + '">\g<in></a>', d)  

    return d

def url_pas(data):
    return parse.quote(data).replace('/','%2F')
    
def plusing(conn, name, link, backtype):
    curs = conn.cursor()

    curs.execute("select title from back where title = ? and link = ? and type = ?", [link, name, backtype])
    if not curs.fetchall():
        curs.execute("insert into back (title, link, type) values (?, ?, ?)", [link, name, backtype])

def namumark(conn, title, data, num, in_c, toc_y):  
    data = start(data)
    data = html_pas(data)
    
    fol_num = 0
    a = mid_pas(data, fol_num, 0, in_c, toc_y)
    data = a[0]
    fol_num = a[1]

    a = include_pas(conn, data, title, in_c, num, toc_y, fol_num)
    data = a[0]
    category = a[1]
    fol_num = a[2]
    backlink = a[3]
    
    data = re.sub("\n##\s?([^\n]*)", "", data)    
    a = redirect_pas(data, title, backlink)
    data = a[0]
    backlink = a[1]
    
    data = blockquote(data)
    data = toc_pas(data, title, num, toc_y)
    data = text_help(data)
    data = macro(data)
    
    a = link(conn, title, data, num, category, backlink)
    data = a[0]
    category = a[1]
    backlink = a[2]
    
    data = indent(data)
    data = footnote(data, fol_num)
    data = table(data)
    data = end(data, category)

    data += '<script>function folding(num, test = 0) { var fol = document.getElementById(\'folding_\' + num); \
            if(fol.style.display == \'inline-block\' || fol.style.display == \'block\') { fol.style.display = \'none\'; } \
            else { if(num % 3 == 0 && test != 1) { fol.style.display = \'block\'; } else { fol.style.display = \'inline-block\'; } } } \
            </script>'
    
    if num == 1:        
        for d4 in backlink:
            t = threading.Thread(target = plusing, args = [conn, d4[0], d4[1], d4[2]])
            t.start()
            t.join()

        conn.commit()
        
    return data