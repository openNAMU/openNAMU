import re
from urllib import parse

def end(data, category):    
    if category:
        data += '<div style="margin-top: 30px;" id="cate">분류: ' + category + '</div>'
            
    data = re.sub("\r\n(?P<in><h[0-6])", "\g<in>", data)
    data = re.sub("(\n#no-br#|#no-br#\n|#no-br#)", "", data)
    data = re.sub("&lt;space&gt;", " ", data)

    com = re.compile('#base64#((?:(?!#\/base64#).)+)#\/base64#', re.DOTALL)
    while 1:
        m = com.search(data)
        if m:
            data = com.sub(parse.unquote(m.groups()[0]).replace('&#95;', '_').replace('\r\n', '<br>'), data, 1)
        else:
            break

    com3 = re.compile('(?:#mid#|#\/mid2#)((?:(?!#\/mid#|#\/mid2#).)+)(?:#\/mid#|#\/mid2#)', re.DOTALL)
    m = com3.search(data)
    while 1:
        m = com3.search(data)
        if m:
            data = com3.sub('{{{' + m.groups()[0] + '}}}', data, 1)
        else:
            break

    data = re.sub('<\/blockquote>(?:(?:\r)?\n){2}<blockquote>', '</blockquote><blockquote>', data)
    data = re.sub('<\/blockquote>(?:(?:\r)?\n)<br><blockquote>', '</blockquote><blockquote>', data)

    data = re.sub('\n', '<br>', data)

    data = re.sub('<br><ul id="list">', '<ul id="list">', data)
    data = re.sub('<\/ul>\r<br>', '</ul>', data)
    data = re.sub('<\/table>\r<br><ul ', '</table><ul ', data)
    data = re.sub('<hr id="under_bar"([^>]*)>(\r)?<br>', '<hr id="under_bar" style="margin-top: -5px;">', data)
    data = re.sub('&lt;isbr&gt;', '\r\n', data)
    data = re.sub('^(?:<br>|\r|\n| )+', '', data)
    data = re.sub('^<div style="margin-top: 30px;" id="cate">', '<div id="cate">', data)        
    data = re.sub('&amp;#92;', '&#92;', data)

    if re.search('\[math\]((?:(?!\[\/math\]).)*)\[\/math\]', data):
        data += '<script type="text/x-mathjax-config">MathJax.Hub.Config({ extensions: ["tex2jax.js", "AMSmath.js"], jax: ["input/TeX", "output/HTML-CSS"], tex2jax: { inlineMath: [ [\'[math]\',\'[/math]\'] ],processEscapes: true }, }); </script>'
        data += '<script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.3/MathJax.js?config=TeX-AMS-MML_HTMLorMML"></script>'

    return data