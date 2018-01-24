import re
from urllib import parse

def url_pas(data):
    return(parse.quote(data).replace('/','%2F'))

def toc_pas(data, title, num, toc_y):
    if(not re.search('\[목차\]', data)):
        if(not re.search('\[목차\(없음\)\]', data)):
            data = re.sub("(?P<in>(={1,6})\s?([^=]*)\s?(?:={1,6})(?:\s+)?\n)", "[목차]\n\g<in>", data, 1)
        else:
            data = re.sub("\[목차\(없음\)\]", "", data)
        
    data = re.sub("(\n)(?P<in>\r\n(={1,6})\s?([^=]*)\s?(?:={1,6})(?:\s+)?\n)", "\g<in>", data)

    i = [0, 0, 0, 0, 0, 0, 0]
    last = 0
    toc_c = -1
    toc_d = -1
    span = ''
    rtoc = '<div id="toc"><span id="toc-name">목차</span><br><br>'
    while(1):
        i[0] += 1
        m = re.search('(={1,6})\s?([^=]*)\s?(?:={1,6})(?:\s+)?\r\n', data)
        if(m):
            result = m.groups()
            wiki = len(result[0])
            if(last < wiki):
                last = wiki
            else:
                last = wiki
                for a in range(wiki + 1, 7):
                    i[a] = 0
            
            i[wiki] += 1

            toc = str(i[1]) + '.' + str(i[2]) + '.' + str(i[3]) + '.' + str(i[4]) + '.' + str(i[5]) + '.' + str(i[6]) + '.'

            toc = re.sub("(?P<in>[0-9]0(?:[0]*)?)\.", '\g<in>#.', toc)

            toc = re.sub("0\.", '', toc)
            toc = re.sub("#\.", '.', toc)
            toc = re.sub("\.$", '', toc)

            if(toc_c == -1):
                margin = ''
                toc_c = toc.count('.')
            else:
                toc_d = toc.count('.')
                if(toc_c == toc_d):
                    margin = 'style="margin-top: 30px;"'
                else:
                    if(toc_d < toc_c):
                        margin = 'style="margin-top: 30px;"'
                    else:
                        margin = 'style="margin-top: 15px;"'
                    
                    toc_c = toc_d

            t = toc.count('.')
            span = '<span style="margin-left: 5px;"></span>' * t

            rtoc += span + '<a href="#s-' + toc + '">' + toc + '</a>. ' + result[1] + '<br>'

            c = re.sub(" $", "", result[1])
            d = c
            c = re.sub("\[\[(([^|]*)\|)?(?P<in>[^\]]*)\]\]", "\g<in>", c)

            edit_d = ''
            if(toc_y == 1):
                edit_d = ' <span style="font-size: 12px;"><a href="/edit/' + url_pas(title) + '?section=' + str(i[0]) + '">(편집)</a></span>'

            data = re.sub('(={1,6})\s?([^=]*)\s?(?:={1,6})(?:\s+)?\n', '<tablenobr><h' + str(wiki) + ' id="' + c + '" ' + margin + '><a href="#toc" id="s-' + toc + '">' \
                                                                         + toc + '.<span style="margin-left: 5px;"></span></a> ' \
                                                                         + d + edit_d + '</h' + str(wiki) + '><hr id="under_bar" style="margin-top: -5px;">\n', data, 1)
        else:
            rtoc += '</div>'
            
            break
    
    data = re.sub("\[목차\]", rtoc, data)

    return(data)