from . import tool

import datetime
import html
import re

end_data = ''
plus_data = ''
nowiki_num = ''
include_name = ''

def nowiki_js(data):
    data = data.replace('\\', '\\\\')
    data = data.replace('"', '\\"')
    data = data.replace('\r', '')
    
    data = re.sub(r'^\n', '', data)

    data = data.replace('\n', '<br>')

    return data

def link_fix(main_link, no_change = 0):
    global end_data

    main_link = main_link.replace('&#x27;', "<link_comma>")
    
    if re.search(r'^:', main_link):
        main_link = re.sub(r'^:', '', main_link)

    if no_change == 0:
        main_link = re.sub(r'^사용자:', 'user:', main_link)
        main_link = re.sub(r'^파일:', 'file:', main_link)
        main_link = re.sub(r'^분류:', 'category:', main_link)

    other_link = re.search(r'[^\\]?(#[^#]+)$', main_link)
    if other_link:
        other_link = other_link.group(1)

        main_link = re.sub(r'(#[^#]+)$', '', main_link)
    else:
        other_link = ''

    main_link = main_link.replace("<link_comma>", "&#x27;")
    main_link = main_link.replace('\\#', '%23')

    find_data = re.findall(r'<span id="((?:include_(?:[0-9]+)_)?(?:nowiki_[0-9]+))">', main_link)
    for i in find_data:
        main_link = main_link.replace('<span id="' + i + '"></span>', end_data[i])

    find_data = re.findall(r'<span id="((?:include_(?:[0-9]+)_)?(?:nowiki_[0-9]+))">', other_link)
    for i in find_data:
        other_link = other_link.replace('<span id="' + i + '"></span>', end_data[i])

    return [main_link, other_link]

def table_parser(data, cel_data, cel_num, start_data, num = 0, cel_color = {}):
    table_class = 'class="'
    div_style = 'style="'
    all_table = 'style="'
    cel_style = 'style="'
    row_style = 'style="'
    row = ''
    cel = 'colspan="' + str(round(len(start_data) / 2)) + '"'

    if not cel_num in cel_color:
        cel_color[cel_num] = ''

    cel_style += cel_color[cel_num]

    if num == 0:
        if re.search(r'^ ', cel_data) and re.search(r' $', cel_data):
            cel_style += 'text-align: center;'
        elif re.search(r'^ ', cel_data):
            cel_style += 'text-align: right;'
        elif re.search(r' $', cel_data):
            cel_style += 'text-align: left;'

    table_state = re.findall(r'&lt;((?:(?!&gt;).)+)&gt;', data)
    for in_state in table_state:
        table_state_data = re.search(r'^([^=]+)=([^=]+)$', in_state)
        if not table_state_data:
            if re.search(r'^(#(?:[0-9a-f-A-F]{3}){1,2}|\w+)$', in_state):
                table_state_data = ['bgcolor', in_state]
            else:
                table_state_data = re.search(r'([^0-9]+)([0-9]+)$', in_state)
                if table_state_data:
                    table_state_data = table_state_data.groups()
                else:
                    table_state_data = [in_state, '']
        else:
            table_state_data = table_state_data.groups()
            table_state_data = [
                table_state_data[0].replace(' ', ''),
                re.sub(r',([^,]+)$', '', table_state_data[1])
            ]

        if table_state_data[0] == 'tablewidth':
            table_data = table_state_data[1]
            div_style += 'width: ' + ((table_data + 'px') if re.search(r'^[0-9]+$', table_data) else table_data) + ';'
            all_table += 'width: 100%;'
        elif table_state_data[0] == 'tableheight':
            table_data = table_state_data[1]
            div_style += 'height: ' + ((table_data + 'px') if re.search(r'^[0-9]+$', table_data) else table_data) + ';'
        elif table_state_data[0] == 'tablealign':
            table_data = table_state_data[1]
            if table_data == 'right':
                div_style += 'float: right;'
            elif table_data == 'center':
                div_style += 'margin: auto;'
                all_table += 'margin: auto;'
        elif table_state_data[0] == 'tabletextalign':
            num = 1
            table_data = table_state_data[1]
            if table_data == 'right':
                all_table += 'text-align: right;'
            elif table_data == 'center':
                all_table += 'text-align: center;'
        elif table_state_data[0] == 'rowtextalign':
            table_data = table_state_data[1]
            if table_data == 'right':
                row_style += 'text-align: right;'
            elif table_data == 'center':
                row_style += 'text-align: center;'
            else:
                row_style += 'text-align: left;'
        elif table_state_data[0] == '-':
            cel = 'colspan="' + table_state_data[1] + '"'
        elif table_state_data[0] == '^|' or table_state_data[0] == 'v|' or table_state_data[0] == '|':
            if table_state_data[0][0] == '^':
                cel_style += 'vertical-align: top;'
            elif table_state_data[0][0] == 'v':
                cel_style += 'vertical-align: bottom;'

            row = 'rowspan="' + table_state_data[1] + '"'
        elif table_state_data[0] == 'rowbgcolor':
            table_data = table_state_data[1]
            row_style += 'background: ' + table_data + ';'
        elif table_state_data[0] == 'rowcolor':
            table_data = table_state_data[1]
            row_style += 'color: ' + table_data + ';'
        elif table_state_data[0] == 'tablebordercolor':
            table_data = table_state_data[1]
            all_table += 'border: ' + table_data + ' 2px solid;'
        elif table_state_data[0] == 'tablebgcolor':
            table_data = table_state_data[1]
            all_table += 'background: ' + table_data + ';'
        elif table_state_data[0] == 'tablecolor':
            table_data = table_state_data[1]
            all_table += 'color: ' + table_data + ';'
        elif table_state_data[0] == 'colbgcolor':
            table_data = table_state_data[1]
            table_data = table_data
            cel_color[cel_num] += 'background: ' + table_data + ';'
            cel_style += 'background: ' + table_data + ';'
        elif table_state_data[0] == 'colcolor':
            table_data = table_state_data[1]
            table_data = table_data
            cel_color[cel_num] += 'color: ' + table_data + ';'
            cel_style += 'color: ' + table_data + ';'
        elif table_state_data[0] == 'bgcolor':
            table_data = table_state_data[1]
            cel_style += 'background: ' + table_data + ';'
        elif table_state_data[0] == 'color':
            table_data = table_state_data[1]
            cel_style += 'color: ' + table_data + ';'
        elif table_state_data[0] == 'width':
            table_data = table_state_data[1]
            cel_style += 'width: ' + ((table_data + 'px') if re.search(r'^[0-9]+$', table_data) else table_data) + ';'
        elif table_state_data[0] == 'height':
            table_data = table_state_data[1]
            cel_style += 'height: ' + ((table_data + 'px') if re.search(r'^[0-9]+$', table_data) else table_data) + ';'
        elif table_state_data[0] == '(' or table_state_data[0] == ':' or table_state_data[0] == ')':
            if table_state_data[0] == '(':
                cel_style += 'text-align: right;'
            elif table_state_data[0] == ':':
                cel_style += 'text-align: center;'
            else:
                cel_style += 'text-align: left;'
        elif table_state_data[0] == 'tableclass':
            table_data = table_state_data[1]
            table_class += table_data

    div_style += '"'
    all_table += '"'
    cel_style += '"'
    row_style += '"'
    table_class += '"'

    return [all_table, row_style, cel_style, row, cel, table_class, num, div_style, cel_color]

def table_start(data):
    data = re.sub(r'\n( +)\|\|', '\n||', data)
    data = re.sub(r'\|\|( +)\n', '||\n', data)
    data = re.sub(r'(\|\|)+\n', '||\n', data)

    while 1:
        cel_num = 0
        table_num = 0
        table_end = ''
        cel_color = {}

        table = re.search(r'\n((?:(?:(?:(?:\|\||\|[^|]+\|)+(?:(?:(?!\|\|).\n*)*))+)\|\|(?:\n)?)+)', data)
        if table:
            table = re.sub(r'(\|\|)+\n', '||\n', table.group(1))
            
            table_caption = re.search(r'^\|([^|]+)\|', table)
            if table_caption:
                table_caption = '<caption>' + table_caption.group(1) + '</caption>'
                
                table = re.sub(r'^\|([^|]+)\|', '||', table)
            else:
                table_caption = ''
            
            table = '\n' + table
            
            table_cel = re.findall(r'(\n(?:(?:\|\|)+)|\|\|\n(?:(?:\|\|)+)|(?:(?:\|\|)+))((?:(?:(?!\n|\|\|).)+\n*)+)', table)
            for i in table_cel:
                cel_plus = re.search(r'^((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)', i[1])
                cel_plus = cel_plus.group(1) if cel_plus else ''
                cel_data = re.sub(r'^((?:&lt;(?:(?:(?!&gt;).)*)&gt;)+)', '', i[1])

                if re.search(r'^\n', i[0]):
                    cel_num = 1

                    cel_plus = table_parser(
                        cel_plus, 
                        cel_data,
                        cel_num,
                        re.sub(r'^\n', '', i[0]),
                        table_num,
                        cel_color
                    )
                    cel_color = cel_plus[8]
                    table_num = cel_plus[6]

                    table_end += '' + \
                        '<div class="table_safe" ' + cel_plus[7] + '>' + \
                            '<table ' + cel_plus[5] + ' ' + cel_plus[0] + '>' + \
                                table_caption + \
                                '<tr ' + cel_plus[1] + '>' + \
                                    '<td ' + cel_plus[2] + ' ' + cel_plus[3] + ' ' + cel_plus[4] + '>' + \
                                        cel_data
                                        
                elif re.search(r'\n', i[0]):
                    cel_num = 1

                    cel_plus = table_parser(
                        cel_plus, 
                        cel_data,
                        cel_num,
                        re.sub(r'^\|\|\n', '', i[0]),
                        table_num,
                        cel_color
                    )
                    cel_color = cel_plus[8]

                    table_end += '' + \
                            '</td>' + \
                        '</tr>' + \
                        '<tr ' + cel_plus[1] + '>' + \
                            '<td ' + cel_plus[2] + ' ' + cel_plus[3] + ' ' + cel_plus[4] + '>' + \
                                cel_data
                else:
                    cel_num += 1

                    cel_plus = table_parser(
                        cel_plus, 
                        cel_data,
                        cel_num,
                        re.sub(r'^\|\|\n', '', i[0]),
                        table_num,
                        cel_color
                    )
                    cel_color = cel_plus[8]

                    table_end += '' + \
                        '</td>' + \
                        '<td ' + cel_plus[2] + ' ' + cel_plus[3] + ' ' + cel_plus[4] + '>' + \
                            cel_data

            table_end += '' + \
                            '</td>' + \
                        '</tr>' + \
                    '</table>' + \
                '</div>' + \
            ''

            data = re.sub(r'\n((?:(?:(?:(?:\|\||\|[^|]+\|)+(?:(?:(?!\|\|).\n*)*))+)\|\|(?:\n)?)+)', '\n' + table_end + '\n', data, 1)
        else:
            break

    return data.replace('||', '<no_table>')

def middle_parser(data):
    global end_data
    global plus_data
    global nowiki_num
    global include_name

    middle_stack = 0
    middle_list = []
    middle_num = 0

    html_num = 0
    syntax_num = 0
    folding_num = 0

    middle_re = re.compile(r'(?:{{{((?:(?:(?! |{{{|}}}|&lt;).)*) ?)|(}}}))')
    middle_all_data = middle_re.findall(data)
    for middle_data in middle_all_data:
        if not middle_data[1]:
            if middle_stack > 0:
                middle_stack += 1

                data = re.sub(r'(?:{{{((?:(?! |{{{|}}}|&lt;).)*)(?P<in> ?)|(}}}))', '<middle_start>' + middle_data[0] + '\g<in>', data, 1)
            else:
                if re.search(r'^(#|@|\+|\-)', middle_data[0]) and not re.search(r'^(#|@|\+|\-){2}|(#|@|\+|\-)\\', middle_data[0]):
                    if re.search(r'^(#(?:[0-9a-f-A-F]{3}){1,2})', middle_data[0]):
                        middle_search = re.search(r'^(#(?:[0-9a-f-A-F]{3}){1,2})', middle_data[0])
                        middle_list += ['span']

                        data = middle_re.sub('<span style="color: ' + middle_search.group(1) + ';">', data, 1)
                    elif re.search(r'^(?:#(\w+))', middle_data[0]):
                        middle_search = re.search(r'^(?:#(\w+))', middle_data[0])
                        middle_list += ['span']

                        data = middle_re.sub('<span style="color: ' + middle_search.group(1) + ';">', data, 1)
                    elif re.search(r'^(?:@((?:[0-9a-f-A-F]{3}){1,2}))', middle_data[0]):
                        middle_search = re.search(r'^(?:@((?:[0-9a-f-A-F]{3}){1,2}))', middle_data[0])
                        middle_list += ['span']

                        data = middle_re.sub('<span style="background: #' + middle_search.group(1) + ';">', data, 1)
                    elif re.search(r'^(?:@(\w+))', middle_data[0]):
                        middle_search = re.search(r'^(?:@(\w+))', middle_data[0])
                        middle_list += ['span']

                        data = middle_re.sub('<span style="background: ' + middle_search.group(1) + ';">', data, 1)
                    elif re.search(r'^(\+|-)([1-5])', middle_data[0]):
                        middle_search = re.search(r'^(\+|-)([1-5])', middle_data[0])
                        middle_search = middle_search.groups()
                        if middle_search[0] == '+':
                            font_size = str(int(middle_search[1]) * 20 + 100)
                        else:
                            font_size = str(100 - int(middle_search[1]) * 10)

                        middle_list += ['span']

                        data = middle_re.sub('<span style="font-size: ' + font_size + '%;">', data, 1)
                    elif re.search(r'^#!wiki', middle_data[0]):
                        middle_data_2 = re.search(r'{{{#!wiki(?: style=(?:&quot;|&#x27;)((?:(?!&quot;|&#x27;).)*)(?:&quot;|&#x27;))?(?: *)\n?', data)
                        if middle_data_2:
                            middle_data_2 = middle_data_2.groups()
                        else:
                            middle_data_2 = ['']

                        middle_list += ['div_1']

                        data = re.sub(
                            r'{{{#!wiki(?: style=(?:&quot;|&#x27;)((?:(?!&quot;|&#x27;).)*)(?:&quot;|&#x27;))?(?: *)\n?',
                            '<div_1 style="' + str(middle_data_2[0] if middle_data_2[0] else '') + '">',
                            data,
                            1
                        )
                    elif re.search(r'^#!syntax', middle_data[0]):
                        syntax_re = re.compile(r'{{{#!syntax ?((?:(?!\n|{{{|}}}).)*)\n?')
                        middle_data_2 = syntax_re.search(data)
                        if middle_data_2:
                            middle_data_2 = middle_data_2.groups()
                        else:
                            middle_data_2 = ['python']

                        if syntax_num == 0:
                            plus_data += 'hljs.initHighlightingOnLoad();\n'

                            syntax_num = 1

                        middle_list += ['pre']

                        data = syntax_re.sub(
                            '<pre id="syntax"><code class="' + middle_data_2[0] + '">',
                            data,
                            1
                        )
                    elif re.search(r'^#!folding', middle_data[0]):
                        middle_list += ['div_dd']

                        folding_data = re.search(r'{{{#!folding ?((?:(?!\n).)*)\n?', data)
                        if folding_data:
                            folding_data = folding_data.groups()
                        else:
                            folding_data = ['Test']

                        plus_data += '' + \
                            'if(document.getElementById("get_' + include_name + 'folding_' + str(folding_num) + '")) { ' + \
                                'document.getElementById("get_' + include_name + 'folding_' + str(folding_num) + '").innerHTML = ' + \
                                    '"' + nowiki_js(folding_data[0]) + '"; ' + \
                                '' + \
                            '}' + \
                            '\n' + \
                        ''
                        data = re.sub(
                            r'{{{#!folding ?((?:(?!\n).)*)\n?', '' + \
                            '<div>' + \
                                '<div style="display: inline-block;">' + \
                                    '<b>' + \
                                        '<a href="javascript:void(0);" ' + \
                                            'onclick="do_open_folding(\'' + include_name + 'folding_' + str(folding_num) + '\');" ' + \
                                            'id="get_' + include_name + 'folding_' + str(folding_num) + '">' + \
                                        '</a>' + \
                                    '</b>' + \
                                '</div_2>' + \
                                '<div id="' + include_name + 'folding_' + str(folding_num) + '" style="display: none;">' + \
                                    '<div_1 style="">\n',
                            data,
                            1
                        )

                        folding_num += 1
                    elif re.search(r'^#!html', middle_data[0]):
                        middle_list += ['span']

                        html_num += 1

                        data = middle_re.sub('<span id="' + include_name + 'render_contect_' + str(html_num) + '">', data, 1)
                    else:
                        middle_list += ['span']

                        data = middle_re.sub('<span>', data, 1)
                else:
                    middle_list += ['code']

                    middle_stack += 1

                    data = middle_re.sub('<code>' + middle_data[0].replace('\\', '\\\\'), data, 1)

                middle_num += 1
        else:
            if middle_list == []:
                data = middle_re.sub('<middle_end>', data, 1)
            else:
                if middle_stack > 0:
                    middle_stack -= 1

                if middle_stack > 0:
                    data = middle_re.sub('<middle_end>', data, 1)
                else:
                    if middle_num > 0:
                        middle_num -= 1

                    if middle_list[middle_num] == 'div_dd':
                        data = middle_re.sub('</div_1></div_2></div_2>', data, 1)
                    elif middle_list[middle_num] == 'pre':
                        data = middle_re.sub('</code></pre>', data, 1)
                    else:
                        data = middle_re.sub('</' + middle_list[middle_num] + '>', data, 1)

                    del middle_list[middle_num]

    while 1:
        if middle_list == []:
            break
        else:
            if middle_stack > 0:
                middle_stack -= 1

            if middle_stack > 0:
                data += '<middle_end>'
            else:
                if middle_num > 0:
                    middle_num -= 1

                if middle_list[middle_num] == '2div':
                    data += '</div_1></div_2></div_2>'
                elif middle_list[middle_num] == 'pre':
                    data += '</code></pre>'
                else:
                    data += '</' + middle_list[middle_num] + '>'

                del middle_list[middle_num]

    data = data.replace('<middle_start>', '{{{')
    data = data.replace('<middle_end>', '}}}')

    while 1:
        nowiki_data = re.search(r'<code>((?:(?:(?!<\/code>).)*\n*)*)<\/code>', data)
        if nowiki_data:
            nowiki_data = nowiki_data.groups()

            nowiki_num += 1
            end_data[include_name + 'nowiki_' + str(nowiki_num)] = nowiki_data[0]
            plus_data += '' + \
                'if(document.getElementById("' + include_name + 'nowiki_' + str(nowiki_num) + '")) {\n' + \
                    'document.getElementById("' + include_name + 'nowiki_' + str(nowiki_num) + '").innerHTML = "' + nowiki_js(nowiki_data[0]) + '";\n' + \
                '}\n' + \
            ''

            data = re.sub(
                r'<code>((?:(?:(?!<\/code>).)*\n*)*)<\/code>',
                '<span id="' + include_name + 'nowiki_' + str(nowiki_num) + '"></span>',
                data,
                1
            )
        else:
            break

    while 1:
        syntax_data = re.search(
            r'<code class="([^"\'>]+)">((?:(?:(?:(?!<\/code>|<span ).)*)\n*)+)<\/code>', 
            data
        )
        if syntax_data:
            syntax_data = syntax_data.groups()

            nowiki_num += 1
            end_data[include_name + 'nowiki_' + str(nowiki_num)] = syntax_data[1]
            plus_data += '' + \
                'if(document.getElementById("' + include_name + 'nowiki_' + str(nowiki_num) + '")) {\n' + \
                    'document.getElementById("' + include_name + 'nowiki_' + str(nowiki_num) + '").innerHTML = "' + nowiki_js(syntax_data[1]) + '";\n' + \
                '}\n' + \
            ''

            data = re.sub(
                r'<code class="([^"\'>]+)">((?:(?:(?:(?!<\/code>|<span ).)*)\n*)+)<\/code>', 
                '<code class="' + syntax_data[0] + '"><span id="' + include_name + 'nowiki_' + str(nowiki_num) + '"></span></code>',
                data,
                1
            )
        else:
            break

    return data

def namumark(conn, data, title, include_num):
    curs = conn.cursor()

    global plus_data
    global end_data
    global nowiki_num
    global include_name

    nowiki_num = 0
    data = '\n' + data + '\n'
    include_name = (include_num + '_') if include_num else ''
    now_time = tool.get_time().split()[0]
    plus_data = ''

    backlink = []
    end_data = {}

    data = re.sub(r'@((?:(?!(?:=|{{{|}}}|\[\[|\]\]|@)).)+)=(?P<in>(?:(?!(?:=|{{{|}}}|\[\[|\]\]|@)).)+)@', '\g<in>', data)
    data = re.sub(r'<math>(?P<in>(?:(?!<\/math>).)+)<\/math>', '[math(\g<in>)]', data)

    data = html.escape(data)
    data = data.replace('\r\n', '\n')

    # 테이블 앞에 공백 있는 경우 처리 필요
    # 테이블 col 옵션 적용시 셀 병합 고려 필요
    data = re.sub(r'\n +\|\|', '\n||', data)

    math_re = re.compile(r'\[math\(((?:(?!\)\]).)+)\)\]', re.I)
    while 1:
        math = math_re.search(data)
        if math:
            math = math.group(1)
            math = math.replace('{', '<math_mid_1>')
            math = math.replace('}', '<math_mid_2>')
            math = math.replace('\\', '<math_slash>')

            data = math_re.sub('<math>' + math + '</math>', data, 1)
        else:
            break

    data = data.replace('\\{', '<break_middle>')
    data = middle_parser(data)
    data = data.replace('<break_middle>', '\\{')

    first = 0
    math_re = re.compile(r'<math>((?:(?!<\/math>).)+)<\/math>', re.I)
    while 1:
        math = math_re.search(data)
        if math:
            math = math.group(1)
            math = math.replace('<math_mid_1>', '{')
            math = math.replace('<math_mid_2>', '}')
            math = math.replace('<math_slash>', '\\')

            first += 1
            data = math_re.sub('<span id="math_' + str(first) + '"></span>', data, 1)

            plus_data += '' + \
                'try {\n' + \
                    'katex.render("' + nowiki_js(html.unescape(math)) + '", document.getElementById(\"' + include_name + 'math_' + str(first) + '\"));\n' + \
                '} catch {\n' + \
                    'document.getElementById(\"' + include_name + 'math_' + str(first) + '\").innerHTML = "<span style=\'color: red;\'>' + nowiki_js(math) + '</span>";\n' + \
                '}\n' + \
            ''
        else:
            break

    num = 0
    while 1:
        one_nowiki = re.search(r'(?:\\)(.)', data)
        if one_nowiki:
            one_nowiki = one_nowiki.groups()

            nowiki_num += 1
            end_data[include_name + 'nowiki_' + str(nowiki_num)] = one_nowiki[0]
            plus_data += '' + \
                'if(document.getElementById("' + include_name + 'nowiki_' + str(nowiki_num) + '")) {\n' + \
                    'document.getElementById("' + include_name + 'nowiki_' + str(nowiki_num) + '").innerHTML = "' + nowiki_js(one_nowiki[0]) + '";\n' + \
                '}\n' + \
            ''

            data = re.sub(r'(?:\\)(.)', '<span id="' + include_name + 'nowiki_' + str(nowiki_num) + '"></span>', data, 1)
        else:
            break

    include_re = re.compile(r'\[include\(((?:(?!\)\]).)+)\)\]', re.I)
    if include_name == '':
        i = 0
        while 1:
            i += 1

            include = include_re.search(data)
            if include:
                include = include.group(1)

                include_data = re.search(r'^((?:(?!,).)+)', include)
                if include_data:
                    include_data = include_data.group(1)
                else:
                    include_data = 'Test'

                include_link = include_data
                backlink += [[title, include_link, 'include']]

                data = include_re.sub('' + \
                    '<a id="' + include_name + 'include_link" class="include_' + str(i) + '" href="/w/' + tool.url_pas(include_link) + '">(' + include_link + ')</a>' + \
                    '<div id="' + include_name + 'include_' + str(i) + '"></div>' + \
                '', data, 1)

                include_plus_data = []
                while 1:
                    include_plus = re.search(r', ?((?:(?!=).)+)=((?:(?!,).)+)', include)
                    if include_plus:
                        include_plus = include_plus.groups()

                        include_data_set = include_plus[1]
                        find_data = re.findall(r'<span id="((?:include_(?:[0-9]+)_)?(?:nowiki_[0-9]+))">', include_data_set)
                        for j in find_data:
                            include_data_set = include_data_set.replace('<span id="' + j + '"></span>', end_data[j])

                        include_plus_data += [[include_plus[0], include_data_set]]

                        include = re.sub(r', ?((?:(?!=).)+)=((?:(?!,).)+)', '', include, 1)
                    else:
                        break

                plus_data += 'load_include("' + include_link + '", "' + include_name + 'include_' + str(i) + '", ' + str(include_plus_data) + ');\n'
            else:
                break
    else:
        data = include_re.sub('', data)

    data = data.replace('&amp;', '&')
    data = re.sub(r'\n##[^\n]+', '', data)

    div_1_re = re.compile(r'<div_1 ([^>]+)>((?:(?!<div_1|<\/div_1>).|\n*)+)<\/div_1>')
    while 1:
        wiki_table_data = div_1_re.search(data)
        if wiki_table_data:
            wiki_table_data = wiki_table_data.groups()
            if re.search(r'\|\|', wiki_table_data[1]):
                end_parser = table_start('\n' + wiki_table_data[1] + '\n')
                end_parser = re.sub(r'^\n', '', end_parser)
                end_parser = re.sub(r'\n$', '', end_parser)
            else:
                end_parser = wiki_table_data[1]

            data = div_1_re.sub('<div id="wiki_div" ' + wiki_table_data[0] + '>' + end_parser + '</div_2>', data, 1)
        else:
            break

    data = data.replace('</div_2>', '</div>')
    data = data.replace('</td>', '</td_1>')

    data += '\n'
    data = data.replace('\\', '&#92;')

    redirect_re = re.compile(r'\n#(?:redirect|넘겨주기) ([^\n]+)', re.I)
    redirect = redirect_re.search(data)
    if redirect:
        redirect = redirect.group(1)

        return_link = link_fix(redirect)
        main_link = html.unescape(return_link[0])
        other_link = return_link[1]

        backlink += [[title, main_link, 'redirect']]
        
        plus_data += '' + \
            'var get_link = window.location.search.match(/(?:\?|&)from=([^&]+)/);\n' + \
            'var get_link_2 = window.location.pathname.match(/^\/w\//);\n' + \
            'if(!get_link && get_link_2) {\n' + \
                'window.location.href = "/w/' + tool.url_pas(main_link) + '?from=' + tool.url_pas(title) + other_link + '";\n' + \
            '}\n' + \
        ''
        data = redirect_re.sub('\nredirect to ' + html.escape(main_link) + other_link, data, 1)

    no_toc_re = re.compile(r'\[(?:목차|toc)\(no\)\]', re.I)
    toc_re = re.compile(r'\[(?:목차|toc)\]', re.I)
    if not no_toc_re.search(data) and not toc_re.search(data):
        data = re.sub(r'(?P<in>\n(={1,6})(#)? ?((?:(?!(?: #=| =)).)+) ?#?(?:=+)\n)', '\n[toc]\g<in>', data, 1)
        auto_toc = 1
    else:
        data = no_toc_re.sub('', data)
        auto_toc = 0

    data = '<div class="all_in_data" id="in_data_0">' + data

    toc_stack = [0, 0, 0, 0, 0, 0]
    toc_num = 0
    toc_data = '' + \
        '<div id="toc">' + \
            '<span id="toc_title">TOC</span>' + \
            '<br>' + \
            '<br>' + \
    ''
    edit_num = 0
    toc_head_re = re.compile(r'\n(={1,6})(#)? ?((?:(?!(?: #=| =)).)+) ?#?(?:=+)\n')
    while 1:
        toc = toc_head_re.search(data)
        if toc:
            toc = toc.groups()
            edit_num += 1

            toc_len_num = len(toc[0])
            toc_len_str = str(toc_len_num)
            toc_len_num -= 1
            toc_stack[toc_len_num] += 1
            for i in range(toc_len_num + 1, 6):
                toc_stack[i] = 0

            edit_num_str = str(edit_num)
            toc_level_str = '.'.join([str(i) for i in toc_stack if i != 0])
            toc_fol = '+' if toc[1] else '-'

            data = toc_head_re.sub(
                '\n' + \
                '</div>'
                '<h' + toc_len_str + ' id="s-' + toc_level_str + '">' + \
                    '<a href="#toc">' + toc_level_str + '.</a> ' + toc[2] + ' ' + \
                    '<span style="font-size: 12px">' + \
                        '<a id="edit_load_' + edit_num_str + '" href="/edit/' + tool.url_pas(title) + '?section=' + edit_num_str + '">(Edit)</a>' + \
                        ' ' + \
                        '<a href="javascript:void(0);" onclick="do_open_folding(\'in_data_' + edit_num_str + '\', this);">' + \
                            '(' + toc_fol + ')' + \
                        '</a>' + \
                    '</span>' + \
                '</h' + toc_len_str + '>' + \
                '<div class="all_in_data"' + (' style="display: none;"' if toc_fol == '+' else '') + ' id="in_data_' + edit_num_str + '">' + \
                    '\n',
                data,
                1
            )

            toc_main_data = toc[2]
            toc_main_data = re.sub(r'\[\*((?:(?! |\]).)*)(?: ((?:(?!(\[\*(?:(?:(?!\]).)+)\]|\])).)+))?\]', '', toc_main_data)
            toc_main_data = re.sub(r'<span id="math_[0-9]"><\/span>', '(Math)', toc_main_data)

            toc_data += '' + \
                '<span style="margin-left:' + str(len(re.findall(r'\.', toc_level_str)) * 10) + 'px;">' + \
                    '<a href="#s-' + toc_level_str + '">' + toc_level_str + '.</a> ' + toc_main_data + \
                '</span>' + \
                '<br>' + \
            ''
        else:
            break

    toc_data += '</div>'
    if auto_toc == 1:
        data = toc_re.sub('<div id="auto_toc">' + toc_data + '</div>', data, 1)
    else:
        data = toc_re.sub(toc_data, data)
    
    macro_re = re.compile(r'\[([^[(]+)\(((?:(?!\[|\)]).)+)\)\]')
    macro_data = macro_re.findall(data)
    for i in macro_data:
        macro_name = i[0].lower()
        if macro_name == 'youtube' or macro_name == 'kakaotv' or macro_name == 'nicovideo':
            width = re.search(r', ?width=((?:(?!,).)+)', i[1])
            if width:
                video_width = width.group(1)
                if re.search(r'^[0-9]+$', video_width):
                    video_width += 'px'
            else:
                video_width = '560px'

            height = re.search(r', ?height=((?:(?!,).)+)', i[1])
            if height:
                video_height = height.group(1)
                if re.search(r'^[0-9]+$', video_height):
                    video_height += 'px'
            else:
                video_height = '315px'

            code = re.search(r'^((?:(?!,).)+)', i[1])
            if code:
                video_code = code.group(1)
            else:
                video_code = ''

            video_start = ''

            if macro_name == 'youtube':
                start = re.search(r', ?(start=(?:(?!,).)+)', i[1])
                if start:
                    video_start = '?' + start.group(1)

                video_code = re.sub(r'^https:\/\/www\.youtube\.com\/watch\?v=', '', video_code)
                video_code = re.sub(r'^https:\/\/youtu\.be\/', '', video_code)

                video_src = 'https://www.youtube.com/embed/' + video_code
            elif macro_name == 'kakaotv':
                video_code = re.sub(r'^https:\/\/tv\.kakao\.com\/channel\/9262\/cliplink\/', '', video_code)
                video_code = re.sub(r'^http:\/\/tv\.kakao\.com\/v\/', '', video_code)

                video_src = 'https://tv.kakao.com/embed/player/cliplink/' + video_code +'?service=kakao_tv'
            else:
                video_src = 'https://embed.nicovideo.jp/watch/' + video_code

            data = macro_re.sub(
                '<iframe style="width: ' + video_width + '; height: ' + video_height + ';" src="' + video_src + video_start + '" frameborder="0" allowfullscreen></iframe>', 
                data, 
                1
            )
        elif macro_name == 'anchor':
            data = macro_re.sub('<span id="' + i[1] + '"></span>', data, 1)
        elif macro_name == 'ruby':
            ruby_code = re.search(r'^([^,]+)', i[1])
            if ruby_code:
                ruby_code = ruby_code.group(1)
            else:
                ruby_code = 'Test'

            ruby_top = re.search(r'ruby=([^,]+)', i[1], flags = re.I)
            if ruby_top:
                ruby_top = ruby_top.group(1)
            else:
                ruby_top = 'Test'

            ruby_color = re.search(r'color=([^,]+)', i[1], flags = re.I)
            if ruby_color:
                ruby_color = 'color: ' + ruby_color.group(1) + ';'
            else:
                ruby_color = ''

            ruby_data = '' + \
                '<ruby>' + \
                    ruby_code \
                    + '<rp>(</rp>' + \
                    '<rt style="' + ruby_color + '">' + ruby_top + '</rt>' + \
                    '<rp>)</rp>' + \
                '</ruby>' + \
            ''

            data = macro_re.sub(ruby_data, data, 1)
        elif macro_name == 'age' or macro_name == 'dday':
            try:
                old = datetime.datetime.strptime(now_time, '%Y-%m-%d')
                will = datetime.datetime.strptime(i[1], '%Y-%m-%d')

                e_data = old - will

                if macro_name == 'age':
                    data = macro_re.sub(str(int(e_data.days / 365)), data, 1)
                else:
                    data = macro_re.sub((str(e_data.days) if re.search(r'^-', str(e_data.days)) else ('+' + str(e_data.days))), data, 1)
            except:
                data = macro_re.sub('age-dday-error', data, 1)
        else:
            data = macro_re.sub('<macro_start>' + i[0] + '<macro_middle>' + i[1] + '<macro_end>', data, 1)
            
    data = data.replace('<macro_start>', '[')
    data = data.replace('<macro_middle>', '(')
    data = data.replace('<macro_end>', ')]')

    blockquote_re = re.compile(r'(\n(?:&gt; ?(?:[^\n]*)\n)+)')
    while 1:
        blockquote = blockquote_re.search(data)
        if blockquote:
            blockquote = re.sub(r'^\n&gt; ?', '', blockquote.group(1))
            blockquote = re.sub(r'\n&gt; ?', '\n', blockquote)
            blockquote = re.sub(r'\n$', '', blockquote)

            data = blockquote_re.sub('\n<blockquote>' + blockquote + '</blockquote>\n', data, 1)
        else:
            break

    hr_re = re.compile(r'\n-{4,9}\n')
    while 1:
        if hr_re.search(data):
            data = hr_re.sub('\n<hr>\n', data, 1)
        else:
            break

    data = re.sub(r'(?P<in>\n +\* ?(?:(?:(?!\|\|).)+))\|\|', '\g<in>\n ||', data)
    data = re.sub(r'(?P<in><div id="folding_(?:[0-9]+)" style="display: none;"><div style="">|<blockquote>)(?P<out> )?\* ', '\g<in>\n\g<out>* ', data)

    li_re = re.compile(r'(\n(?:(?: *)\* (?:[^\n]+)\n)+)')
    while 1:
        li_data = li_re.search(data)
        if li_data:
            li_data = li_data.group(1)
            li_end_data = ''

            sub_li = re.findall(r'( *)\* ([^\n]+)\n', li_data)
            for i in sub_li:
                li_end_data += '<li style="margin-left: ' + str(20 if len(i[0]) == 0 else (len(i[0]) * 20)) + 'px;">' + i[1] + '</li>'

            data = li_re.sub('\n\n<ul>' + li_end_data + '</ul>\n', data, 1)
        else:
            break

    data = re.sub(r'<\/ul>\n \|\|', '</ul>||', data)
    data = re.sub(r'\|\|</blockquote>', '</blockquote>||', data)

    while 1:
        indent = re.search(r'\n( +)', data)
        if indent:
            data = re.sub(r'\n( +)', '\n' + ('<span style="margin-left: 20px;"></span>' * len(indent.group(1))), data, 1)
        else:
            break

    data = table_start(data)

    category = ''
    link_re = re.compile('\[\[((?:(?!\[\[|\]\]|<\/td>).)+)\]\]')
    category_re = re.compile(r'^(?:category|분류):', re.I)
    e_link_id = 0
    while 1:
        link = link_re.search(data)
        if link:
            link = link.group(1)
            str_e_link_id = str(e_link_id)
            e_link_id += 1

            link_split = re.search(r'((?:(?!\|).)+)(?:\|((?:(?!\|).)+))', link)
            if link_split:
                link_split = link_split.groups()

                main_link = link_split[0]
                see_link = link_split[1]
                inter_same = 0
            else:
                main_link = link
                see_link = link
                inter_same = 1

            if re.search(r'^((?:file|파일)|(?:out|외부)):', main_link):
                file_style = ''

                file_width = re.search(r'width=((?:(?!&).)+)', see_link)
                if file_width:
                    file_width = file_width.group(1)
                    if re.search(r'px$', file_width):
                        file_style += 'width: ' + file_width + ';'
                    else:
                        file_style += 'width: ' + file_width + 'px;'

                file_height = re.search(r'height=((?:(?!&).)+)', see_link)
                if file_height:
                    file_height = file_height.group(1)
                    if re.search(r'px$', file_height):
                        file_style += 'height: ' + file_height + ';'
                    else:
                        file_style += 'height: ' + file_height + 'px;'

                file_align = re.search(r'align=((?:(?!&).)+)', see_link)
                if file_align:
                    file_align = file_align.group(1)
                    if file_align == 'center':
                        file_align = 'display: block; text-align: center;'
                    else:
                        file_align = 'float: ' + file_align + ';'
                else:
                    file_align = ''

                file_color = re.search(r'bgcolor=((?:(?!&).)+)', see_link)
                if file_color:
                    file_color = 'background: ' + file_color.group(1) + '; display: inline-block;'
                else:
                    file_color = ''

                file_alt = re.search(r'alt=((?:(?!&).)+)', see_link)
                if file_alt:
                    file_alt = file_alt.group(1)
                else:
                    file_alt = ''

                if re.search(r'^(?:out|외부):', main_link):
                    file_src = re.sub(r'^(?:out|외부):', '', main_link)

                    file_alt = main_link if file_alt == '' else file_alt
                    exist = 'Yes'
                else:
                    file_data = re.search(r'^(?:file|파일):((?:(?!\.).)+)\.(.+)$', main_link)
                    if file_data:
                        file_data = file_data.groups()
                        file_name = file_data[0]
                        file_end = file_data[1]

                        if main_link != title:
                            backlink += [[title, main_link, 'file']]
                    else:
                        file_name = 'TEST'
                        file_end = 'jpg'

                    file_src = '/image/' + tool.sha224_replace(file_name) + '.' + file_end
                    file_alt = ('file:' + file_name + '.' + file_end) if file_alt == '' else file_alt
                    exist = None

                data = link_re.sub(
                    '<span style="' + file_align + '">' + \
                        '<span  style="' + file_color + '" ' + \
                                'class="' + include_name + 'file_finder" ' + \
                                'under_style="' + file_style + '" ' + \
                                'under_alt="' + file_alt + '" ' + \
                                'under_src="' + file_src + '" ' + \
                                'under_href="' + ("out_link" if exist else '/upload?name=' + tool.url_pas(file_name)) + '">' + \
                        '</span>' + \
                    '</span>',
                    data,
                    1
                )
            elif category_re.search(main_link):
                if include_name == '':
                    if category == '':
                        category += '<div id="cate_all"><div id="cate">Category : '

                    main_link = category_re.sub('category:', main_link)
                    link_id = ''

                    curs.execute(tool.db_change("select title from data where title = ?"), [main_link])
                    if re.search(r'#blur', main_link):
                        link_id = ' hidden_link'
                        main_link = main_link.replace('#blur', '')
                        see_link = see_link.replace('#blur', '')

                    backlink += [[title, main_link, 'cat']]
                    category += '' + \
                        '<a class="' + include_name + 'link_finder' + link_id + '" ' + \
                            'href="/w/' + tool.url_pas(main_link) + '">' + \
                            category_re.sub('', see_link) + \
                        '</a> | ' + \
                    ''

                data = link_re.sub('', data, 1)
            elif re.search(r'^inter:((?:(?!:).)+):', main_link):
                inter_data = re.search(r'^inter:((?:(?!:).)+):((?:(?!\]\]).)+)', main_link)
                inter_data = inter_data.groups()

                curs.execute(tool.db_change('select plus, plus_t from html_filter where html = ? and kind = "inter_wiki"'), [inter_data[0]])
                inter = curs.fetchall()
                if inter:
                    return_link = link_fix(inter_data[1], 1)
                    main_link = html.unescape(return_link[0])
                    other_link = return_link[1]

                    inter_view = inter[0][1] if inter[0][1] != '' else (inter_data[0] + ':')

                    data = link_re.sub('' + \
                        '<a id="inside" ' + \
                            'name="' + include_name + 'set_link_' + str_e_link_id + '" ' + \
                            'href="">' + inter_view + see_link + '</a>' + \
                        '', 
                        data, 
                        1
                    )
                    plus_data += "" + \
                        "document.getElementsByName('" + include_name + "set_link_" + str_e_link_id + "')[0].href = '" + \
                            (inter[0][0] + tool.url_pas(main_link) + other_link).replace('\'', '\\\'') + "';" + \
                        "\n" + \
                    ""
                    if inter_same == 1:
                        plus_data += "" + \
                            "document.getElementsByName('" + include_name + "set_link_" + str_e_link_id + "')[0].innerHTML = '" + \
                                (inter_view + main_link + other_link).replace('\'', '\\\'') + "';" + \
                            "\n" + \
                        ""
                else:
                    data = link_re.sub('', data, 1)
            elif re.search(r'^(\/(?:.+))$', main_link):
                under_title = re.search(r'^(\/(?:.+))$', main_link)
                under_title = under_title.group(1)

                if see_link != main_link:
                    data = link_re.sub('[[' + title + under_title + '|' + see_link + ']]', data, 1)
                else:
                    data = link_re.sub('[[' + title + under_title + ']]', data, 1)
            elif re.search(r'^http(s)?:\/\/', main_link):
                data = link_re.sub('' + \
                    '<a id="out_link" ' + \
                        'name="' + include_name + 'set_link_' + str_e_link_id + '" ' + \
                        'rel="nofollow" ' + \
                        'href="">' + see_link + '</a>' + \
                    '', 
                    data, 
                    1
                )

                plus_data += "" + \
                    "document.getElementsByName('" + include_name + "set_link_" + str_e_link_id + "')[0].href = '" + \
                        main_link.replace('\'', '\\\'') + "';" + \
                    "\n" + \
                ""
                if inter_same == 1:
                    plus_data += "" + \
                        "document.getElementsByName('" + include_name + "set_link_" + str_e_link_id + "')[0].innerHTML = '" + \
                            main_link.replace('\'', '\\\'') + "';" + \
                        "\n" + \
                    ""
            else:
                return_link = link_fix(main_link)
                main_link = html.unescape(return_link[0])
                other_link = return_link[1]

                if re.search(r'^\/', main_link):
                    main_link = re.sub(r'^\/', title + '/', main_link)
                elif re.search(r'\.\.\/\/', main_link):
                    main_link = re.sub(r'\.\.\/\/', '/', main_link)
                elif re.search(r'^\.\.\/', main_link):
                    main_link = re.sub(r'^\.\.\/', re.sub(r'(?P<in>.+)\/.*$', '\g<in>', title), main_link)

                if main_link != title and main_link != '':
                    backlink += [[title, main_link, '']]

                    curs.execute(tool.db_change("select title from data where title = ?"), [main_link])
                    if not curs.fetchall():
                        backlink += [[title, main_link, 'no']]

                    data = link_re.sub('' + \
                        '<a class="' + include_name + 'link_finder" ' + \
                            'name="' + include_name + 'set_link_' + str_e_link_id + '" ' + \
                            'title="" ' + \
                            'href="">' + see_link + '</a>' + \
                        '', 
                        data,
                        1
                    )

                    plus_data += "" + \
                        "document.getElementsByName('" + include_name + "set_link_" + str_e_link_id + "')[0].href = '" + \
                            ('/w/' + tool.url_pas(main_link) + other_link).replace('\'', '\\\'') + "';" + \
                        "\n" + \
                    ""
                    plus_data += "" + \
                        "document.getElementsByName('" + include_name + "set_link_" + str_e_link_id + "')[0].title = '" + \
                            (html.escape(main_link) + other_link).replace('\'', '\\\'') + "';" + \
                        "\n" + \
                    ""
                    if inter_same == 1:
                        plus_data += "" + \
                            "document.getElementsByName('" + include_name + "set_link_" + str_e_link_id + "')[0].innerHTML = '" + \
                                (html.escape(main_link) + other_link).replace('\'', '\\\'') + "';" + \
                            "\n" + \
                        ""
                else:
                    if re.search(r'^#', other_link):
                        data = link_re.sub(
                            '<a title="' + other_link + '" href="' + other_link + '">' + (other_link if see_link == other_link else see_link) + '</a>',
                            data,
                            1
                        )
                    else:
                        data = link_re.sub('<b>' + see_link + '</b>', data, 1)
        else:
            break

    if re.search(r'\[pagecount(?:\([^()]+\))?\]', data, flags = re.I):
        plus_data += 'page_count();\n'
        data = re.sub(r'\[pagecount(?:\([^()]+\))?\]', '<span class="all_page_count"></span>', data, flags = re.I)

    data = re.sub(r'\[date\]', now_time, data, flags = re.I)
    data = re.sub(r'\[clearfix\]', '<div style="clear:both"></div>', data, flags = re.I)
    data = re.sub(r'\[br\]', '<br>', data, flags = re.I)

    # 각주 기능
    footnote_number = 0
    footnote_all = []
    footnote_dict = {}
    footnote_re = {}
    footdata_all = '<ul id="footnote_data">'
    re_footnote = re.compile(r'(?:\[\*((?:(?! |\]).)*)(?: ((?:(?!(?:\[\*|\])).)+))?\]|(\[(?:각주|footnote)\]))')
    while 1:
        footnote = re_footnote.search(data)
        if footnote:
            footnote_data = footnote.groups()
            if footnote_data[2]:
                footnote_all.sort()

                for footdata in footnote_all:
                    if footdata[2] == 0:
                        footdata_in = ''
                    else:
                        footdata_in = footdata[2]

                    footdata_all += '' + \
                        '<li>' + \
                            '<a href="javascript:do_open_foot(\'' + include_name + 'fn-' + str(footdata[0]) + '\', 1);" ' + \
                                'id="' + include_name + 'cfn-' + str(footdata[0]) + '">' + \
                                '(' + footdata[1] + ')' + \
                            '</a> <span id="' + include_name + 'fn-' + str(footdata[0]) + '">' + footdata_in + '</span>' + \
                        '</li>' + \
                    ''

                data = re_footnote.sub(footdata_all + '</ul>', data, 1)

                footnote_all = []
                footdata_all = '<ul id="footnote_data">'
            else:
                footnote = footnote_data[1]
                footnote_name = footnote_data[0]
                if footnote_name and not footnote:
                    if footnote_name in footnote_dict:
                        footnote_re[footnote_name] += 1

                        foot_plus_num = str(footnote_re[footnote_name])
                        footshort = footnote_dict[footnote_name] + '.' + foot_plus_num

                        footnote_all += [[float(footshort), footshort, 0]]

                        data = re_footnote.sub('' + \
                            '<sup>' + \
                                '<a href="javascript:do_open_foot(\'' + include_name + 'fn-' + footshort + '\', 0);" ' + \
                                    'id="' + include_name + 'rfn-' + footshort + '">' + \
                                    '(' + footnote_name + ')' + \
                                '</a>' + \
                            '</sup><span id="' + include_name + 'dfn-' + footshort + '"></span>' + \
                        '', data, 1)
                    else:
                        data = re_footnote.sub('<sup><a href="javascript:void(0);">(' + footnote_name + ')</a></sup>', data, 1)
                else:
                    footnote_number += 1

                    if not footnote_name:
                        footnote_name = str(footnote_number)

                    footnote_dict.update({ footnote_name : str(footnote_number) })

                    if not footnote_name in footnote_re:
                        footnote_re.update({ footnote_name : 0 })
                    else:
                        footnote_re[footnote_name] += 1

                    footnote_all += [[footnote_number, footnote_name, footnote]]

                    data = re_footnote.sub('' + \
                        '<sup>' + \
                            '<a href="javascript:do_open_foot(\'' + include_name + 'fn-' + str(footnote_number) + '\', 0);" ' + \
                                'id="' + include_name + 'rfn-' + str(footnote_number) + '">' + \
                                '(' + footnote_name + ')' + \
                            '</a>' + \
                        '</sup><span id="' + include_name + 'dfn-' + str(footnote_number) + '"></span>' + \
                    '', data, 1)
        else:
            break

    data = re.sub(r'\n+$', '', data)

    footnote_all.sort()

    for footdata in footnote_all:
        if footdata[2] == 0:
            footdata_in = ''
        else:
            footdata_in = str(footdata[2])

        footdata_all += '' + \
            '<li>' + \
                '<a href="javascript:do_open_foot(\'' + include_name + 'fn-' + str(footdata[0]) + '\', 1);" ' + \
                    'id="' + include_name + 'cfn-' + str(footdata[0]) + '">' + \
                    '(' + str(footdata[1]) + ')' + \
                '</a> <span id="' + include_name + 'fn-' + str(footdata[0]) + '">' + footdata_in + '</span>' + \
            '</li>' + \
        ''

    footdata_all += '</ul>'
    footdata_all = '' if footdata_all == '<ul id="footnote_data"></ul>' else footdata_all
    footdata_all = '</div>' + footdata_all    

    data += footdata_all

    # 기본 꾸미기 문법
    data = re.sub(r'&#x27;&#x27;&#x27;(?P<in>((?!&#x27;&#x27;&#x27;).)+)&#x27;&#x27;&#x27;', '<b>\g<in></b>', data)
    data = re.sub(r'&#x27;&#x27;(?P<in>((?!&#x27;&#x27;).)+)&#x27;&#x27;', '<i>\g<in></i>', data)
    data = re.sub(r'~~(?P<in>(?:(?!~~).)+)~~', '<s>\g<in></s>', data)
    data = re.sub(r'--(?P<in>(?:(?!--).)+)--', '<s>\g<in></s>', data)
    data = re.sub(r'__(?P<in>(?:(?!__).)+)__', '<u>\g<in></u>', data)
    data = re.sub(r'\^\^(?P<in>(?:(?!\^\^).)+)\^\^', '<sup>\g<in></sup>', data)
    data = re.sub(r',,(?P<in>(?:(?!,,).)+),,', '<sub>\g<in></sub>', data)

    # 최종 처리
    data += (re.sub(r' \| $', '', category) + '</div></div>') if category != '' else ''
    data = data.replace('<no_table>', '||')
    data = data.replace('</td_1>', '</td>')
    data = re.sub(r'<\/ul>\n?', '</ul>', data)
    data = re.sub(r'<\/pre>\n?', '</pre>', data)
    data = re.sub(r'(?P<in><div class="all_in_data"(?:(?:(?!id=).)+)? id="in_data_([^"]+)">)(\n)+', '\g<in>', data)
    data = data.replace('\n\n<ul>', '\n<ul>')
    data = data.replace('</ul>\n\n', '</ul>')
    data = re.sub(r'^(\n)+', '', data)
    data = re.sub(r'(\n)+<ul id="footnote_data">', '<ul id="footnote_data">', data)
    data = re.sub(r'(?P<in><td(((?!>).)*)>)\n', '\g<in>', data)
    data = re.sub(r'(\n)?<hr>(\n)?', '<hr>', data)
    data = data.replace('</ul>\n\n<ul>', '</ul>\n<ul>')
    data = data.replace('</ul>\n<ul>', '</ul><ul>')
    data = data.replace('\n</ul>', '</ul>')
    data = data.replace('\n', '<br>')

    # 파일, 링크, HTML 작동 JS
    plus_data += '' + \
        'get_link_state("' + include_name + '");\n' + \
        'get_file_state("' + include_name + '");\n' + \
    ''
    plus_data = 'render_html("' + include_name + 'render_contect");\n' + plus_data

    return [data, plus_data, backlink]