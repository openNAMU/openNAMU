from . import start
from . import mid_pas
from . import html_pas
from . import include_pas
from . import macro
from . import redirect_pas
from . import blockquote
from . import toc_pas
from . import text_help
from . import link
from . import indent
from . import footnote
from . import table
from . import end
import re
import html
import sqlite3
from urllib import parse
import time
import asyncio

def send_p(d):
    d = html.escape(d)

    js_p = re.compile('javascript:', re.I)
    d = js_p.sub('', d)

    d = re.sub('&lt;a href="(?:[^"]*)"&gt;(?P<in>(?:(?!&lt;).)*)&lt;\/a&gt;', '<a href="' + url_pas('\g<in>') + '">\g<in></a>', d)  

    return(d)

def url_pas(data):
    return(parse.quote(data).replace('/','%2F'))
    
async def plusing(conn, name, link, backtype):
    curs = conn.cursor()
    curs.execute("select title from back where title = ? and link = ? and type = ?", [link, name, backtype])
    if(not curs.fetchall()):
        curs.execute("insert into back (title, link, type) values (?, ?,  ?)", [link, name, backtype])

def namumark(conn, title, data, num, in_c, toc_y):  
    data = start.start(data)
    data = html_pas.html_pas(data)
    
    fol_num = 0
    a = mid_pas.mid_pas(data, fol_num, 0, in_c)
    data = a[0]
    fol_num = a[1]

    a = include_pas.include_pas(conn, data, title, in_c, num, toc_y, fol_num)
    data = a[0]
    category = a[1]
    fol_num = a[2]
    backlink = a[3]
    
    data = re.sub("\r\n##\s?([^\n]*)\r\n", "\r\n", data)    
    a = redirect_pas.redirect_pas(data, title, backlink)
    data = a[0]
    backlink = a[1]
    
    data = blockquote.blockquote(data)
    data = toc_pas.toc_pas(data, title, num, toc_y)
    data = text_help.text_help(data)
    data = macro.macro(data)
    
    a = link.link(conn, title, data, num, category, backlink)
    data = a[0]
    category = a[1]
    backlink = a[2]
    
    data = indent.indent(data)
    data = footnote.footnote(data, fol_num)
    data = table.table(data)
    data = end.end(data, category)
    
    if(num == 1):        
        asyncio.set_event_loop(asyncio.new_event_loop())
        loop = asyncio.get_event_loop()
        for d4 in backlink:
            loop.run_until_complete(plusing(conn, d4[0], d4[1], d4[2]))
        loop.close()
        conn.commit()
        
    return(data)