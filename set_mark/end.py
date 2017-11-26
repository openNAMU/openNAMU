import re

def end(data, category):    
    if(category):
        data += '<div style="margin-top: 30px;" id="cate">분류: ' + category + '</div>'
            
    data = re.sub("\r\n(?P<in><h[0-6])", "\g<in>", data)
    data = re.sub("(\n<nobr>|<nobr>\n|<nobr>)", "", data)
    data = re.sub("#no#(?P<in>.)#\/no#", "\g<in>", data)
    data = re.sub("&lt;space&gt;", " ", data)

    data = re.sub('<\/blockquote>(?:(?:\r)?\n){2}<blockquote>', '</blockquote><blockquote>', data)
    data = re.sub('<\/blockquote>(?:(?:\r)?\n)<br><blockquote>', '</blockquote><blockquote>', data)

    data = re.sub('\n', '<br>', data)
    data = re.sub('<hr style="margin-top: -5px;"><br>', '<hr style="margin-top: -5px;">', data)
    data = re.sub('&lt;isbr&gt;', '\r\n', data)
    data = re.sub('^(?:<br>|\r|\n| )+', '', data)
    data = re.sub('^<div style="margin-top: 30px;" id="cate">', '<div id="cate">', data)        
    data = re.sub('&amp;#92;', '&#92;', data)
        
    return(data)