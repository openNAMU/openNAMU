import re

def table_p(d, d2, d3, num = 0):
    table_class = 'class="'
    alltable = 'style="'
    celstyle = 'style="'
    rowstyle = 'style="'
    row = ''
    cel = ''

    table_w = re.search("&lt;table\s?width=((?:(?!&gt;).)*)&gt;", d)
    table_h = re.search("&lt;table\s?height=((?:(?!&gt;).)*)&gt;", d)
    table_a = re.search("&lt;table\s?align=((?:(?!&gt;).)*)&gt;", d)
    if table_w:
        alltable += 'width: ' + table_w.groups()[0] + ';'
    if table_h:
        alltable += 'height: ' + table_h.groups()[0] + ';'
    if table_a:
        if table_a.groups()[0] == 'right':
            alltable += 'float: right;'
        elif table_a.groups()[0] == 'center':
            alltable += 'margin: auto;'
            
    table_t_a = re.search("&lt;table\s?textalign=((?:(?!&gt;).)*)&gt;", d)
    if table_t_a:
        num = 1
        if table_t_a.groups()[0] == 'right':
            alltable += 'text-align: right;'
        elif table_t_a.groups()[0] == 'center':
            alltable += 'text-align: center;'

    row_t_a = re.search("&lt;row\s?textalign=((?:(?!&gt;).)*)&gt;", d)
    if row_t_a:
        if row_t_a.groups()[0] == 'right':
            rowstyle += 'text-align: right;'
        elif row_t_a.groups()[0] == 'center':
            rowstyle += 'text-align: center;'
        else:
            rowstyle += 'text-align: left;'
    
    table_cel = re.search("&lt;-((?:(?!&gt;).)*)&gt;", d)
    if table_cel:
        cel = 'colspan="' + table_cel.groups()[0] + '"'
    else:
        cel = 'colspan="' + str(round(len(d2) / 2)) + '"'   

    table_row = re.search("&lt;\|((?:(?!&gt;).)*)&gt;", d)
    if table_row:
        row = 'rowspan="' + table_row.groups()[0] + '"'

    row_bgcolor_2 = re.search("&lt;rowbgcolor=(#(?:[0-9a-f-A-F]{3}){1,2})&gt;", d)
    row_bgcolor_3 = re.search("&lt;rowbgcolor=(\w+)&gt;", d)
    if row_bgcolor_2:
        rowstyle += 'background: ' + row_bgcolor_2.groups()[0] + ';'
    elif row_bgcolor_3:
        rowstyle += 'background: ' + row_bgcolor_3.groups()[0] + ';'
        
    table_border_2 = re.search("&lt;table\s?bordercolor=(#(?:[0-9a-f-A-F]{3}){1,2})&gt;", d)
    table_border_3 = re.search("&lt;table\s?bordercolor=(\w+)&gt;", d)
    if table_border_2:
        alltable += 'border: ' + table_border_2.groups()[0] + ' 2px solid;'
    elif table_border_3:
        alltable += 'border: ' + table_border_3.groups()[0] + ' 2px solid;'
        
    table_bgcolor_2 = re.search("&lt;table\s?bgcolor=(#(?:[0-9a-f-A-F]{3}){1,2})&gt;", d)
    table_bgcolor_3 = re.search("&lt;table\s?bgcolor=(\w+)&gt;", d)
    if table_bgcolor_2:
        alltable += 'background: ' + table_bgcolor_2.groups()[0] + ';'
    elif table_bgcolor_3:
        alltable += 'background: ' + table_bgcolor_3.groups()[0] + ';'
        
    bgcolor_2 = re.search("&lt;(?:bgcolor=)?(#(?:[0-9a-f-A-F]{3}){1,2})&gt;", d)
    bgcolor_3 = re.search("&lt;(?:bgcolor=)?(\w+)&gt;", d)
    if bgcolor_2:
        celstyle += 'background: ' + bgcolor_2.groups()[0] + ';'
    elif bgcolor_3:
        celstyle += 'background: ' + bgcolor_3.groups()[0] + ';'
        
    n_width = re.search("&lt;width=((?:(?!&gt;).)*)&gt;", d)
    n_height = re.search("&lt;height=((?:(?!&gt;).)*)&gt;", d)
    if n_width:
        celstyle += 'width: ' + n_width.groups()[0] + ';'
    if n_height:
        celstyle += 'height: ' + n_height.groups()[0] + ';'
        
    text_right = re.search("&lt;\)&gt;", d)
    text_center = re.search("&lt;:&gt;", d)
    text_left = re.search("&lt;\(&gt;",  d)
    if text_right:
        celstyle += 'text-align: right;'
    elif text_center:
        celstyle += 'text-align: center;'
    elif text_left:
        celstyle += 'text-align: left;'
    elif num == 0:
        if re.search('^ (.*) $', d3):
            celstyle += 'text-align: center;'
        elif re.search('^ (.*)$', d3):
            celstyle += 'text-align: right;'
        elif re.search('^(.*) $', d3):
            celstyle += 'text-align: left;'

    text_class = re.search("&lt;table\s?class=((?:(?!&gt;).)+)&gt;", d)
    if text_class:
        d = text_class.groups()
        table_class += d[0]
        
    alltable += 'margin-top: 10px;"'
    celstyle += '"'
    rowstyle += '"'
    table_class += '"'

    return [alltable, rowstyle, celstyle, row, cel, table_class, num]

def table(data):
    data += '\r\n'
    data = re.sub('(\r+)', '\r', data)
    data = re.sub("(?:\|\|\r\n)", "#table-start##table-no-br#", data)
        
    while 1:
        y = re.search("(\|\|(?:(?:(?!\|\|).)+(?:\n*))+)", data)
        if y:
            a = y.groups()
            
            mid_data = re.sub("\|\|", "#table-start#", a[0])
            mid_data = re.sub("\r\n", "<br>", mid_data)
            
            data = re.sub("(\|\|(?:(?:(?!\|\|).)+(?:\n*))+)", mid_data, data, 1)
        else:
            break
            
    data = re.sub("#table-start#", "||", data)
    data = re.sub("#table-no-br#", "\r\n", data)
    
    while 1:
        m = re.search("(?:\n|<br>)(\|\|(?:(?:(?:.*)\n?)\|\|)+)", data)
        if m:
            results = m.groups()
            table = results[0]

            while 1:
                a = re.search("^(\|\|(?:(?:\|\|)*))((?:&lt;(?:(?:(?!&gt;).)*)&gt;)*)((?:(?!\|\||<\/td>).)*)", table)
                if a:
                    row = ''
                    cel = ''
                    celstyle = ''
                    rowstyle = ''
                    alltable = ''
                    table_d = ''
                    num = 0

                    result = a.groups()
                    if result[1]:
                        table_d = table_p(result[1], result[0], result[2])
                        alltable = table_d[0]
                        rowstyle = table_d[1]
                        celstyle = table_d[2]
                        row = table_d[3]
                        cel = table_d[4]
                        table_class = table_d[5]
                        num = table_d[6]
                            
                        table = re.sub("^(\|\|(?:(?:\|\|)*))((?:&lt;(?:(?:(?!&gt;).)*)&gt;)*)", "<table " + table_class + " " + alltable + "><tbody><tr " + rowstyle + "><td " + cel + " " + row + " " + celstyle + ">", table, 1)
                    else:
                        cel = 'colspan="' + str(round(len(result[0]) / 2)) + '"'

                        if re.search('^ (.*) $', result[2]):
                            celstyle += 'text-align: center;'
                        elif re.search('^ (.*)$', result[2]):
                            celstyle += 'text-align: right;'
                        elif re.search('^(.*) $', result[2]):
                            celstyle += 'text-align: left;'

                        table = re.sub("^(\|\|(?:(?:\|\|)*))((?:&lt;(?:(?:(?!&gt;).)*)&gt;)*)", "<table style='margin-top: 10px;'><tbody><tr><td " + cel + " style='" + celstyle + "'>", table, 1)
                else:
                    break
                    
            table = re.sub("\|\|$", "</td></tr></tbody></table>", table)
            
            while 1:
                b = re.search("\|\|\r\n(\|\|(?:(?:\|\|)*))((?:&lt;(?:(?:(?!&gt;).)*)&gt;)*)((?:(?!\|\||<\/td>).)*)", table)
                if b:
                    row = ''
                    cel = ''
                    celstyle = ''
                    rowstyle = ''
                    table_d = ''

                    result = b.groups()
                    if result[1]:
                        table_d = table_p(result[1], result[0], result[2], num)
                        rowstyle = table_d[1]
                        celstyle = table_d[2]
                        row = table_d[3]
                        cel = table_d[4]
                        
                        table = re.sub("\|\|\r\n(\|\|(?:(?:\|\|)*))((?:&lt;(?:(?:(?!&gt;).)*)&gt;)*)", "</td></tr><tr " + rowstyle + "><td " + cel + " " + row + " " + celstyle + ">", table, 1)
                    else:
                        cel = 'colspan="' + str(round(len(result[0]) / 2)) + '"'

                        if re.search('^ (.*) $', result[2]):
                            celstyle += 'text-align: center;'
                        elif re.search('^ (.*)$', result[2]):
                            celstyle += 'text-align: right;'
                        elif re.search('^(.*) $', result[2]):
                            celstyle += 'text-align: left;'

                        table = re.sub("\|\|\r\n(\|\|(?:(?:\|\|)*))((?:&lt;(?:(?:(?!&gt;).)*)&gt;)*)", "</td></tr><tr><td " + cel + " style='" + celstyle + "'>", table, 1)
                else:
                    break

            while 1:
                c = re.search("(\|\|(?:(?:\|\|)*))((?:&lt;(?:(?:(?!&gt;).)*)&gt;)*)((?:(?!\|\||<\/td>).)*)", table)
                if c:
                    row = ''
                    cel = ''
                    celstyle = ''
                    table_d = ''

                    result = c.groups()
                    if result[1]:
                        table_d = table_p(result[1], result[0], result[2], num)
                        celstyle = table_d[2]
                        row = table_d[3]
                        cel = table_d[4]

                        table = re.sub("(\|\|(?:(?:\|\|)*))((?:&lt;(?:(?:(?!&gt;).)*)&gt;)*)", "</td><td " + cel + " " + row + " " + celstyle + ">", table, 1)
                    else:
                        cel = 'colspan="' + str(round(len(result[0]) / 2)) + '"'

                        if re.search('^ (.*) $', result[2]):
                            celstyle += 'text-align: center;'
                        elif re.search('^ (.*)$', result[2]):
                            celstyle += 'text-align: right;'
                        elif re.search('^(.*) $', result[2]):
                            celstyle += 'text-align: left;'

                        table = re.sub("(\|\|(?:(?:\|\|)*))((?:&lt;(?:(?:(?!&gt;).)*)&gt;)*)", "</td><td " + cel + " style='" + celstyle + "'>", table, 1)
                else:
                    break
            
            data = re.sub("(?:\n|<br>)(\|\|(?:(?:(?:.*)\n?)\|\|)+)", table, data, 1)
        else:
            break
        
    return data