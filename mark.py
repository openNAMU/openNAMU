from set_mark.start import *
from set_mark.html_pas import *
from set_mark.mid_pas import *
from set_mark.include import *
from set_mark.macro import *
from set_mark.redirect_pas import *
from set_mark.blockquote import *
from set_mark.toc import *
from set_mark.text_help import *
from set_mark.link import *
from set_mark.indent import *
from set_mark.footnote import *
from set_mark.table import *
from set_mark.end import *
import re
import asyncio

def send_p(d):
    d = html.escape(d)

    js_p = re.compile('javascript:', re.I)
    d = js_p.sub('', d)

    d = re.sub('&lt;a href="(?:[^"]*)"&gt;(?P<in>(?:(?!&lt;).)*)&lt;\/a&gt;', '<a href="' + url_pas('\g<in>') + '">\g<in></a>', d)  

    return(d)
    
async def plusing(name, link, backtype):
    curs.execute("select title from back where title = ? and link = ? and type = ?", [link, name, backtype])
    if(not curs.fetchall()):
        curs.execute("insert into back (title, link, type) values (?, ?,  ?)", [link, name, backtype])
                        
    return('')

def namumark(title, data, num, in_c, toc_y):    
    data = start(data)
    data = html_pas(data)
    
    fol_num = 0
    data = mid_pas(data, fol_num, 0, in_c)

    a = include(data, title, in_c, num, toc_y, fol_num)
    data = a[0]
    category = a[1]
    fol_num = a[2]
    backlink = a[3]
    
    data = re.sub("\r\n##\s?([^\n]*)\r\n", "\r\n", data)
    data = savemark(data)
    
    a = redirect_pas(data, backlink)
    data = a[0]
    backlink = a[1]
    
    data = blockquote(data)
    data = toc_pas(data, title, num, toc_y)
    data = text_help(data)
    data = macro(data)
    
    a = link(title, data, num, category, backlink)
    data = a[0]
    category = a[1]
    backlink = a[2]

    
    data = indent(data)
    data = footnote(data, fol_num)
    data = table(data)
    data = end(data, category)
    
    if(num == 1):        
        asyncio.set_event_loop(asyncio.new_event_loop())
        loop = asyncio.get_event_loop()
        for d4 in back_list:
            loop.run_until_complete(plusing(d4[0], d4[1], d4[2]))
        loop.close()
        conn.commit()
        
    return(data)