import re

def indent(data):
    while 1:
        m = re.search("(\n(?:(?:( *)\* ?(?:[^\n]*))\n?)+)", data)
        if m:
            end = m.groups()[0]

            while 1:
                isspace = re.search("\n( *)\* ?([^\n]*)", end)
                if isspace:
                    spacebar = isspace.groups()
                    if len(spacebar[0]) == 0:
                        up = 20
                    else:
                        up = len(spacebar[0]) * 20

                    end = re.sub("\n( *)\* ?([^\n]*)", '<li style="margin-left: ' + str(up) + 'px">' + spacebar[1] + '</li>', end, 1)
                else:
                    break

            end = re.sub("\n", '', end)
            data = re.sub("(\n(?:(?:( *)\* ?(?:[^\n]*))\n?)+)", '\r\n<ul style="margin-top: 10px; margin-bottom: 10px;" id="list">' + end + '</ul>\r\n', data, 1)
        else:
            break

    while 1:
        b = re.search("(<\/h[0-9]>|\n)( +)", data)
        if b:
            result = b.groups()
            up = re.sub(' ', '<span id="in"></span>', result[1])

            if re.search('<\/h[0-9]>', result[0]):
                data = re.sub("(?P<in>\/h[0-9]>)( +)", '\g<in>' + up, data, 1)
            else:
                data = re.sub("(?:\n)( +)", '<br>' + up, data, 1)
        else:
            break
            
    return data