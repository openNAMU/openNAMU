from . import tool

import datetime
import html
import re

def table_parser(data, cel_data, start_data, num = 0):
    table_class = 'class="'
    all_table = 'style="'
    cel_style = 'style="'
    row_style = 'style="'
    row = ''
    cel = ''

    table_width = re.search("&lt;table ?width=((?:(?!&gt;).)*)&gt;", data)
    if table_width:
        if re.search('^[0-9]+$', table_width.groups()[0]):
            all_table += 'width: ' + table_width.groups()[0] + 'px;'
        else:
            all_table += 'width: ' + table_width.groups()[0] + ';'
    
    table_height = re.search("&lt;table ?height=((?:(?!&gt;).)*)&gt;", data)
    if table_height:
        if re.search('^[0-9]+$', table_height.groups()[0]):
            all_table += 'height: ' + table_height.groups()[0] + 'px;'
        else:
            all_table += 'height: ' + table_height.groups()[0] + ';'
    
    table_align = re.search("&lt;table ?align=((?:(?!&gt;).)*)&gt;", data)
    if table_align:
        if table_align.groups()[0] == 'right':
            all_table += 'float: right;'
        elif table_align.groups()[0] == 'center':
            all_table += 'margin: auto;'
            
    table_text_align = re.search("&lt;table ?textalign=((?:(?!&gt;).)*)&gt;", data)
    if table_text_align:
        num = 1

        if table_text_align.groups()[0] == 'right':
            all_table += 'text-align: right;'
        elif table_text_align.groups()[0] == 'center':
            all_table += 'text-align: center;'

    row_table_align = re.search("&lt;row ?textalign=((?:(?!&gt;).)*)&gt;", data)
    if row_table_align:
        if row_table_align.groups()[0] == 'right':
            row_style += 'text-align: right;'
        elif row_table_align.groups()[0] == 'center':
            row_style += 'text-align: center;'
        else:
            row_style += 'text-align: left;'
    
    table_cel = re.search("&lt;-((?:(?!&gt;).)*)&gt;", data)
    if table_cel:
        cel = 'colspan="' + table_cel.groups()[0] + '"'
    else:
        cel = 'colspan="' + str(round(len(start_data) / 2)) + '"'   

    table_row = re.search("&lt;\|((?:(?!&gt;).)*)&gt;", data)
    if table_row:
        row = 'rowspan="' + table_row.groups()[0] + '"'

    row_bgcolor = re.search("&lt;rowbgcolor=(#(?:[0-9a-f-A-F]{3}){1,2}|\w+)&gt;", data)
    if row_bgcolor:
        row_style += 'background: ' + row_bgcolor.groups()[0] + ';'
        
    table_border = re.search("&lt;table ?bordercolor=(#(?:[0-9a-f-A-F]{3}){1,2}|\w+)&gt;", data)
    if table_border:
        all_table += 'border: ' + table_border.groups()[0] + ' 2px solid;'
        
    table_bgcolor = re.search("&lt;table ?bgcolor=(#(?:[0-9a-f-A-F]{3}){1,2}|\w+)&gt;", data)
    if table_bgcolor:
        all_table += 'background: ' + table_bgcolor.groups()[0] + ';'
        
    bgcolor = re.search("&lt;(?:bgcolor=)?(#(?:[0-9a-f-A-F]{3}){1,2}|\w+)&gt;", data)
    if bgcolor:
        cel_style += 'background: ' + bgcolor.groups()[0] + ';'
        
    cel_width = re.search("&lt;width=((?:(?!&gt;).)*)&gt;", data)
    if cel_width:
        cel_style += 'width: ' + cel_width.groups()[0] + 'px;'

    cel_height = re.search("&lt;height=((?:(?!&gt;).)*)&gt;", data)
    if cel_height:
        cel_style += 'height: ' + cel_height.groups()[0] + 'px;'
        
    text_right = re.search("&lt;\)&gt;", data)
    text_center = re.search("&lt;:&gt;", data)
    text_left = re.search("&lt;\(&gt;",  data)
    if text_right:
        cel_style += 'text-align: right;'
    elif text_center:
        cel_style += 'text-align: center;'
    elif text_left:
        cel_style += 'text-align: left;'
    elif num == 0:
        if re.search('^ ', cel_data) and re.search(' $', cel_data):
            cel_style += 'text-align: center;'
        elif re.search('^ ', cel_data):
            cel_style += 'text-align: right;'
        elif re.search(' $', cel_data):
            cel_style += 'text-align: left;'

    text_class = re.search("&lt;table ?class=((?:(?!&gt;).)+)&gt;", data)
    if text_class:
        table_class += text_class.groups()[0]
        
    all_table += '"'
    cel_style += '"'
    row_style += '"'
    table_class += '"'

    return [all_table, row_style, cel_style, row, cel, table_class, num]
    
def table_start(data):
    while 1:
        table = re.search('\n((?:(?:(?:(?:\|\|)+(?:(?:(?!\|\|).(?:\n)*)*))+)\|\|(?:\n)?)+)', data)
        if table:
            table = table.groups()[0]
            
            # return [all_table, row_style, cel_style, row, cel, table_class, num]
            while 1:
                all_table = re.search('^((?:\|\|)+)((?:&lt;(?:(?:(?!&gt;).)+)&gt;)*)\n*((?:(?!\|\|).\n*)*)', table)
                if all_table:
                    all_table = all_table.groups()
                    
                    return_table = table_parser(all_table[1], all_table[2], all_table[0])
                    
                    number = return_table[6]
                    
                    table = re.sub('^((?:\|\|)+)((?:&lt;(?:(?:(?!&gt;).)+)&gt;)*)\n*', '\n<table ' + return_table[5] + ' ' + return_table[0] + '><tbody><tr ' + return_table[1] + '><td ' + return_table[2] + ' ' + return_table[3] + ' ' + return_table[4] + '>', table, 1)
                else:
                    break
                    
            table = re.sub('\|\|\n?$', '</td></tr></tbody></table>', table)

            while 1:
                row_table = re.search('\|\|\n((?:\|\|)+)((?:&lt;(?:(?:(?!&gt;).)+)&gt;)*)\n*((?:(?!\|\||<\/td>).\n*)*)', table)
                if row_table:
                    row_table = row_table.groups()
                    
                    return_table = table_parser(row_table[1], row_table[2], row_table[0], number)
                    
                    table = re.sub('\|\|\n((?:\|\|)+)((?:&lt;(?:(?:(?!&gt;).)+)&gt;)*)\n*', '</td></tr><tr ' + return_table[1] + '><td ' + return_table[2] + ' ' + return_table[3] + ' ' + return_table[4] + '>', table, 1)
                else:
                    break

            while 1:
                cel_table = re.search('((?:\|\|)+)((?:&lt;(?:(?:(?!&gt;).)+)&gt;)*)\n*((?:(?:(?!\|\||<\/td>).)|\n)*\n*)', table)
                if cel_table:
                    cel_table = cel_table.groups()
                    
                    return_table = table_parser(cel_table[1], re.sub('\n', ' ', cel_table[2]), cel_table[0], number)
                    
                    table = re.sub('((?:\|\|)+)((?:&lt;(?:(?:(?!&gt;).)+)&gt;)*)\n*', '</td><td ' + return_table[2] + ' ' + return_table[3] + ' ' + return_table[4] + '>', table, 1)
                else:
                    break

            data = re.sub('\n((?:(?:(?:(?:\|\|)+(?:(?:(?!\|\|).(?:\n)*)*))+)\|\|(?:\n)?)+)', table, data, 1)
        else:
            break
            
    return data
    
def link_fix(main_link):
    if re.search('^:', main_link):
        main_link = re.sub('^:', '', main_link)

    other_link = re.search('(#.+)$', main_link)
    if other_link:
        other_link = other_link.groups()[0]

        main_link = re.sub('(#.+)$', '', main_link)
    else:
        other_link = ''
        
    return [main_link, other_link]

def start(conn, data, title):
    # DB 지정
    curs = conn.cursor()

    # 초기 설정
    data = '\n' + data + '\n'
    backlink = []
    plus_data = '<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.12.0/styles/default.min.css"><script src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/9.12.0/highlight.min.js"></script><script>hljs.initHighlightingOnLoad(); function folding(num, test = 0) { var fol = document.getElementById(\'folding_\' + num); if(fol.style.display == \'inline-block\' || fol.style.display == \'block\') { fol.style.display = \'none\'; } else { if(num % 3 == 0 && test != 1) { fol.style.display = \'block\'; } else { fol.style.display = \'inline-block\'; } } }</script>'
    end_data= []
    
    # XSS 이스케이프
    data = html.escape(data)

    # 포함 문법 처리
    while 1:
        include = re.search('\[include\(((?:(?!\)\]).)+)\)\]', data)
        if include:
            include = include.groups()[0]
    
            include_data = re.search('^((?:(?!,).)+)', include)
            if include_data:
                include_data = include_data.groups()[0]
            else:
                include_data = 'Test'

            include_link = include_data

            backlink += [[title, include_link, 'include']]

            include = re.sub('^((?:(?!,).)+)', '', include)
            
            # 틀 NoWiki
            num = 0
            while 1:
                include_one_nowiki = re.search('(?:\\\\){2}(.)', include)
                if include_one_nowiki:
                    include_one_nowiki = include_one_nowiki.groups()

                    num += 1

                    end_data += [['include_one_nowiki_' + str(num), include_one_nowiki[0]]]

                    include = re.sub('(?:\\\\){2}(.)', '<span id="include_one_nowiki_' + str(num) + '"></span>', include, 1)
                else:
                    break

            curs.execute("select data from data where title = ?", [include_data])
            include_data = curs.fetchall()
            if include_data:
                include_parser = include_data[0][0]

                while 1:
                    include_plus = re.search(', ?((?:(?!=).)+)=((?:(?!,).)+)', include)
                    if include_plus:
                        include_plus = include_plus.groups()
                        include_parser = re.sub('@' + include_plus[0] + '@', include_plus[1], include_parser)

                        include = re.sub(', ?((?:(?!=).)+)=((?:(?!,).)+)', '', include, 1)
                    else:
                        break

                include_parser = re.sub('\[\[분류:(((?!\]\]|#include).)+)\]\]', '', include_parser)
                include_parser = html.escape(include_parser)

                data = re.sub('\[include\(((?:(?!\)\]).)+)\)\]', '<include>\n<a id="include_link" href="/w/' + tool.url_pas(include_link) + '">[' + include_link + ']</a>\n' + include_parser + '\n</include>', data, 1)
            else:
                data = re.sub('\[include\(((?:(?!\)\]).)+)\)\]', '<a id="not_thing" href="/w/' + tool.url_pas(include_link) + '">' + include_link + '</a>', data, 1)
        else:
            break

    # 개행 정리
    data = re.sub('\r\n', '\n', data)

    # 기타 처리
    data = re.sub('&amp;', '&', data)

    # HTML 허용
    curs.execute('select html from html_filter')
    html_db = curs.fetchall()

    src_list = ["www.youtube.com", "serviceapi.nmv.naver.com", "tv.kakao.com", "www.google.com", "serviceapi.rmcnmv.naver.com"]
    html_list = ['div', 'span', 'embed', 'iframe', 'ruby', 'rp', 'rt']
    
    html_data = re.findall('&lt;(\/)?((?:(?!&gt;| ).)+)( (?:(?:(?!&gt;).)+)?)?&gt;', data)
    for in_data in html_data:
        if in_data[0] == '':
            if in_data[1] in html_list or in_data[1] in html_db[0]:
                if re.search('&lt;\/' + in_data[1] + '&gt;', data):
                    src = re.search('src=([^ ]*)', in_data[2])
                    if src:
                        v_src = re.search('http(?:s)?:\/\/([^/\'" ]*)', src.groups()[0])
                        if v_src:
                            if not v_src.groups()[0] in src_list:
                                and_data = re.sub('&#x27;', '\'', re.sub('&quot;', '"', re.sub('src=([^ ]*)', '', in_data[2])))
                            else:
                                and_data = re.sub('&#x27;', '\'', re.sub('&quot;', '"', in_data[2]))
                        else:
                            and_data = re.sub('&#x27;', '\'', re.sub('&quot;', '"', re.sub('src=([^ ]*)', '', in_data[2])))
                    else:
                        and_data = re.sub('&#x27;', '\'', re.sub('&quot;', '"', in_data[2]))
                        

                    data = data.replace('&lt;' + in_data[1] + in_data[2] + '&gt;', '<' + in_data[1] + and_data + '>', 1)
                    data = re.sub('&lt;\/' + in_data[1] + '&gt;', '</' + in_data[1] + '>', data, 1)

    position = re.compile('position', re.I)
    data = position.sub('', data)

    # 표 정리
    data = re.sub('\n( +)\|\|', '\n||', data)
    data = re.sub('\|\|( +)\n', '||\n', data)

    # 주석 처리
    data = re.sub('\n##(((?!\n).)+)', '', data)

    # 중괄호 문법 처리
    middle_stack = 0
    middle_list = []
    middle_number = 0
    fol_num = 0
    while 1:
        middle_data = re.search('(?:{{{((?:(?! |{{{|}}}).)*) ?|(}}}))', data)
        if middle_data:
            middle_data = middle_data.groups()
            if not middle_data[1]:
                if middle_stack > 0:
                    middle_stack += 1
                    
                    data = re.sub('(?:{{{((?:(?! |{{{|}}}).)*)(?P<in> ?)|(}}}))', '&#123;&#123;&#123;' + middle_data[0] + '\g<in>', data, 1)
                else:
                    if re.search('^(#|@|\+|\-)', middle_data[0]):
                        middle_search = re.search('^(#(?:[0-9a-f-A-F]{3}){1,2})', middle_data[0])
                        if middle_search:                            
                            middle_list += ['span']
                            
                            data = re.sub('(?:{{{((?:(?! |{{{|}}}).)*) ?|(}}}))', '<span style="color: ' + middle_search.groups()[0] + ';">', data, 1)
                        else:
                            middle_search = re.search('^(?:#(\w+))', middle_data[0])
                            if middle_search:
                                middle_list += ['span']
                                
                                data = re.sub('(?:{{{((?:(?! |{{{|}}}).)*) ?|(}}}))', '<span style="color: ' + middle_search.groups()[0] + ';">', data, 1)
                            else:
                                middle_search = re.search('^(?:@((?:[0-9a-f-A-F]{3}){1,2}))', middle_data[0])
                                if middle_search:
                                    middle_list += ['span']
                                    
                                    data = re.sub('(?:{{{((?:(?! |{{{|}}}).)*) ?|(}}}))', '<span style="background: #' + middle_search.groups()[0] + ';">', data, 1)
                                else:
                                    middle_search = re.search('^(?:@(\w+))', middle_data[0])
                                    if middle_search:
                                        middle_list += ['span']
                                        
                                        data = re.sub('(?:{{{((?:(?! |{{{|}}}).)*) ?|(}}}))', '<span style="background: ' + middle_search.groups()[0] + ';">', data, 1)
                                    else:
                                        middle_search = re.search('^(\+|-)([1-5])', middle_data[0])
                                        if middle_search:
                                            middle_search = middle_search.groups()
                                            if middle_search[0] == '+':
                                                font_size = str(int(middle_search[1]) * 20 + 100)
                                            else:
                                                font_size = str(100 - int(middle_search[1]) * 10)

                                            middle_list += ['span']
                                            
                                            data = re.sub('(?:{{{((?:(?! |{{{|}}}).)*) ?|(}}}))', '<span style="font-size: ' + font_size + '%;">', data, 1)
                                        else:
                                            middle_search = re.search('^#!wiki', middle_data[0])
                                            if middle_search:
                                                middle_data_2 = re.search('{{{#!wiki(?: style=(?:&quot;|&#x27;)((?:(?!&quot;|&#x27;).)*)(?:&quot;|&#x27;))?\n?', data)
                                                if middle_data_2:
                                                    middle_data_2 = middle_data_2.groups()
                                                else:
                                                    middle_data_2 = ['']

                                                middle_list += ['div']
                                                
                                                data = re.sub('{{{#!wiki(?: style=(?:&quot;|&#x27;)((?:(?!&quot;|&#x27;).)*)(?:&quot;|&#x27;))?\n?', '<div id="wiki_div" style="' + str(middle_data_2[0]) + '">', data, 1)
                                            else:
                                                middle_search = re.search('^#!syntax', middle_data[0])
                                                if middle_search:                                                
                                                    middle_data_2 = re.search('{{{#!syntax ((?:(?!\n).)+)\n?', data)
                                                    if middle_data_2:
                                                        middle_data_2 = middle_data_2.groups()
                                                    else:
                                                        middle_data_2 = ['python']

                                                    middle_list += ['pre']
                                                    
                                                    data = re.sub('{{{#!syntax ((?:(?!\n).)+)\n?', '<pre id="syntax"><code class="' + middle_data_2[0] + '">', data, 1)
                                                else:
                                                    middle_search = re.search('^#!html', middle_data[0])
                                                    if middle_search:
                                                        middle_list += ['span']
                                                        
                                                        data = re.sub('(?:{{{((?:(?! |{{{|}}}).)*) ?|(}}}))', '<span>', data, 1)
                                                    else:
                                                        middle_search = re.search('^#!folding', middle_data[0])
                                                        if middle_search:
                                                            middle_list += ['2div']
                                                            
                                                            folding_data = re.search('{{{#!folding ?((?:(?!\n).)*)\n?', data)
                                                            if folding_data:
                                                                folding_data = folding_data.groups()
                                                            else:
                                                                folding_data = ['Test']
                                                            
                                                            data = re.sub('{{{#!folding ?((?:(?!\n).)*)\n?', "<div>" + str(folding_data[0]) + " <div id='folding_" + str(fol_num + 1) + "' style='display: inline-block;'><a href='javascript:void(0);' onclick='folding(" + str(fol_num + 1) + "); folding(" + str(fol_num + 2) + "); folding(" + str(fol_num) + ");'>[펼치기]</a></div_end><div id='folding_" + str(fol_num + 2) + "' style='display: none; '><a href='javascript:void(0);' onclick='folding(" + str(fol_num + 1) + "); folding(" + str(fol_num + 2) + "); folding(" + str(fol_num) + ");'>[접기]</a></div_end><div id='folding_" + str(fol_num) + "' style='display: none;'>\n", data, 1)
                                                            
                                                            fol_num += 3
                                                        else:
                                                            middle_list += ['span']

                                                            data = re.sub('(?:{{{((?:(?! |{{{|}}}).)*) ?|(}}}))', '<span>', data, 1)
                    else:
                        middle_list += ['code']
                        
                        middle_stack += 1
                        
                        data = re.sub('(?:{{{((?:(?! |{{{|}}}).)*)|(}}}))', '<code>' + middle_data[0].replace('\\', '\\\\'), data, 1)
                
                    middle_number += 1
            else:
                if middle_list == []:
                    data = re.sub('(?:{{{((?:(?! |{{{|}}}).)*) ?|(}}}))', '&#125;&#125;&#125;', data, 1)
                else:
                    if middle_stack > 0:
                        middle_stack -= 1

                    if middle_stack > 0:
                        data = re.sub('(?:{{{((?:(?! |{{{|}}}).)*) ?|(}}}))', '&#125;&#125;&#125;', data, 1)
                    else:                    
                        if middle_number > 0:
                            middle_number -= 1
                            
                        if middle_list[middle_number] == '2div':
                            data = re.sub('(?:{{{((?:(?! |{{{|}}}).)*) ?|(}}}))', '</div_end></div_end>', data, 1)
                        elif middle_list[middle_number] == 'pre':
                            data = re.sub('(?:{{{((?:(?! |{{{|}}}).)*) ?|(}}}))', '</code></pre>', data, 1)
                        else:
                            data = re.sub('(?:{{{((?:(?! |{{{|}}}).)*) ?|(}}}))', '</' + middle_list[middle_number] + '>', data, 1)
                        
                        del(middle_list[middle_number])
        else:
            break

    # NoWiki 처리
    num = 0
    while 1:
        nowiki_data = re.search('<code>((?:(?:(?!<\/code>).)*\n*)*)<\/code>', data)
        if nowiki_data:
            nowiki_data = nowiki_data.groups()

            num += 1

            end_data += [['nowiki_' + str(num), nowiki_data[0]]]

            data = re.sub('<code>((?:(?:(?!<\/code>).)*\n*)*)<\/code>', '<span id="nowiki_' + str(num) + '"></span>', data, 1)
        else:
            break

    # Syntax 처리
    num = 0
    while 1:
        syntax_data = re.search('<code class="((?:(?!").)+)">((?:(?:(?:(?!<\/code>|<span id="syntax_)).)+\n*)+)<\/code>', data)
        if syntax_data:
            syntax_data = syntax_data.groups()

            num += 1

            end_data += [['syntax_' + str(num), syntax_data[1]]]

            data = re.sub('<code class="((?:(?!").)+)">((?:(?:(?:(?!<\/code>|<span id="syntax_)).)+\n*)+)<\/code>', '<code class="' + syntax_data[0] + '"><span id="syntax_' + str(num) + '"></span></code>', data, 1)
        else:
            break
           
    # 이중 표 처리
    while 1:
        wiki_table_data = re.search('<div id="wiki_div" ((?:(?!>).)+)>((?:(?!<div id="wiki_div"|<\/div>).\n*)+)<\/div>', data)
        if wiki_table_data:
            wiki_table_data = wiki_table_data.groups()
            if re.search('\|\|', wiki_table_data[1]):
                end_parser = re.sub('\n$', '', re.sub('^\n', '', table_start('\n' + wiki_table_data[1] + '\n')))
            else:
                end_parser = wiki_table_data[1]

            data = re.sub('<div id="wiki_div" ((?:(?!>).)+)>((?:(?!<div id="wiki_div"|<\/div>).\n*)+)<\/div>', '<div ' + wiki_table_data[0] + '>' + end_parser + '</div_end>', data, 1)
        else:
            break
            
    data = re.sub('<\/div_end>', '</div>', data)
    data = re.sub('<\/td>', '</td_end>', data)
    
    # 수식 처리
    first = 0
    while 1:
        math = re.search('&lt;math&gt;((?:(?!&lt;\/math&gt;).)+)&lt;\/math&gt;', data)
        if math:
            if first == 0:
                plus_data += '<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/KaTeX/0.9.0/katex.min.css" integrity="sha384-TEMocfGvRuD1rIAacqrknm5BQZ7W7uWitoih+jMNFXQIbNl16bO8OZmylH/Vi/Ei" crossorigin="anonymous"><script src="https://cdnjs.cloudflare.com/ajax/libs/KaTeX/0.9.0/katex.min.js" integrity="sha384-jmxIlussZWB7qCuB+PgKG1uLjjxbVVIayPJwi6cG6Zb4YKq0JIw+OMnkkEC7kYCq" crossorigin="anonymous"></script>'

            math = math.groups()[0]
            
            first += 1
            
            data = re.sub('&lt;math&gt;((?:(?!&lt;\/math&gt;).)+)&lt;\/math&gt;', '<span id="math_' + str(first) + '"></span>', data, 1)

            plus_data += '<script>katex.render("' + math.replace('\\', '\\\\').replace('&lt;', '<').replace('&gt;', '>') +'", document.getElementById("math_' + str(first) + '"));</script>'
        else:
            break
            
    # 한 글자 NoWiki
    num = 0
    while 1:
        one_nowiki = re.search('(?:\\\\)(.)', data)
        if one_nowiki:
            one_nowiki = one_nowiki.groups()

            num += 1

            end_data += [['one_nowiki_' + str(num), one_nowiki[0]]]

            data = re.sub('(?:\\\\)(.)', '<span id="one_nowiki_' + str(num) + '"></span>', data, 1)
        else:
            break

    # 추가 이스케이프
    data = data.replace('\\', '&#92;')

    # 텍스트 꾸미기 문법
    data = re.sub('&#x27;&#x27;&#x27;(?P<in>((?!&#x27;&#x27;&#x27;).)+)&#x27;&#x27;&#x27;', '<b>\g<in></b>', data)
    data = re.sub('&#x27;&#x27;(?P<in>((?!&#x27;&#x27;).)+)&#x27;&#x27;', '<i>\g<in></i>', data)
    data = re.sub('~~(?P<in>(?:(?!~~).)+)~~', '<s>\g<in></s>', data)
    data = re.sub('--(?P<in>(?:(?!~~).)+)--', '<s>\g<in></s>', data)
    data = re.sub('__(?P<in>(?:(?!__).)+)__', '<u>\g<in></u>', data)
    data = re.sub('\^\^(?P<in>(?:(?!\^\^).)+)\^\^', '<sup>\g<in></sup>', data)
    data = re.sub(',,(?P<in>(?:(?!,,).)+),,', '<sub>\g<in></sub>', data)

    # 넘겨주기 변환
    redirect = re.search('\n#(?:redirect|넘겨주기) ((?:(?!\n).)+)\n', data)
    if redirect:
        redirect = redirect.groups()[0]
        
        return_link = link_fix(redirect)
        main_link = return_link[0]
        other_link = return_link[1]
        
        backlink += [[title, main_link, 'redirect']]
        
        data = re.sub('\n#(?:redirect|넘겨주기) (?P<in>(?:(?!\n).)+)\n', '<meta http-equiv="refresh" content="0; url=/w/' + tool.url_pas(main_link) + '?froms=' + tool.url_pas(title) + other_link + '">', data, 1)

    # [목차(없음)] 처리
    if not re.search('\[목차\(없음\)\]\n', data):
        if not re.search('\[목차\]', data):
            data = re.sub('\n(?P<in>={1,6}) ?(?P<out>(?:(?!=).)+) ?={1,6}\n', '\n[목차]\n\g<in> \g<out> \g<in>\n', data, 1)
    else:
        data = re.sub('\[목차\(없음\)\]\n', '', data)
        
    # 문단 문법
    toc_full = 0
    toc_top_stack = 6
    toc_stack = [0, 0, 0, 0, 0, 0]
    edit_number = 0
    toc_data = '<div id="toc"><span style="font-size: 18px;">목차</span>\n\n'
    while 1:
        toc = re.search('\n(={1,6}) ?((?:(?!\n).)+) ?\n', data)
        if toc:
            toc = toc.groups()
            
            toc_number = len(toc[0])
            edit_number += 1

            # 더 크면 그 전 스택은 초기화
            if toc_full > toc_number:
                for i in range(toc_number, 6):
                    toc_stack[i] = 0

            if toc_top_stack > toc_number:
                toc_top_stack = toc_number
                    
            toc_full = toc_number        
            toc_stack[toc_number - 1] += 1
            toc_number = str(toc_number)
            all_stack = ''

            # 스택 합치기
            for i in range(0, 6):
                all_stack += str(toc_stack[i]) + '.'

            all_stack = re.sub('0.', '', all_stack)
            
            data = re.sub('\n(={1,6}) ?((?:(?!\n).)+) ?\n', '\n<h' + toc_number + ' id="s-' + re.sub('\.$', '', all_stack) + '"><a href="#toc">' + all_stack + '</a> ' + re.sub('=*$', '', toc[1]) + ' <span style="font-size: 12px"><a href="/edit/' + tool.url_pas(title) + '?section=' + str(edit_number) + '">(편집)</a></span></h' + toc_number + '>\n', data, 1)
            
            toc_data += '<span style="margin-left: ' + str((toc_full - toc_top_stack) * 10) + 'px;"><a href="#s-' + re.sub('\.$', '', all_stack) + '">' + all_stack + '</a> ' + re.sub('\[\*((?:(?! |\]).)*)(?: ((?:(?!\]).)+))?\]', '', re.sub('=*$', '', toc[1])) + '</span>\n'
        else:
            break

    toc_data += '</div>'
    
    data = re.sub('\[목차\]', toc_data, data)
    
    while 1:
        hr = re.search('\n-{4,9}\n', data)
        if hr:
            data = re.sub('\n-{4,9}\n', '\n<hr>\n', data, 1)
        else:
            break

    data += '\n'

    # 일부 매크로 처리
    data = tool.savemark(data)
    
    data = re.sub("\[anchor\((?P<in>(?:(?!\)\]).)+)\)\]", '<span id="\g<in>"></span>', data)          
    data = re.sub("\[nicovideo\((?P<in>(?:(?!,|\)\]).)+)(?:(?:(?!\)\]).)*)\)\]", "[[http://embed.nicovideo.jp/watch/\g<in>|\g<in>]]", data)
    data = re.sub('\[ruby\((?P<in>(?:(?!,).)+)\, ?(?P<out>(?:(?!\)\]).)+)\)\]', '<ruby>\g<in><rp>(</rp><rt>\g<out></rt><rp>)</rp></ruby>', data)

    # 원래 코드 재탕
    now_time = tool.get_time()

    data = re.sub('\[date\]', now_time, data)
    
    time_data = re.search('^([0-9]{4}-[0-9]{2}-[0-9]{2})', now_time)
    time = time_data.groups()
    
    age_data = re.findall('\[age\(([0-9]{4}-[0-9]{2}-[0-9]{2})\)\]', data)
    for age in age_data:
        old = datetime.datetime.strptime(time[0], '%Y-%m-%d')
        will = datetime.datetime.strptime(age, '%Y-%m-%d')
        
        e_data = old - will
        
        data = re.sub('\[age\(([0-9]{4})-([0-9]{2})-([0-9]{2})\)\]', str(int(int(e_data.days) / 365)), data, 1)

    dday_data = re.findall('\[dday\(([0-9]{4}-[0-9]{2}-[0-9]{2})\)\]', data)
    for dday in dday_data:
        old = datetime.datetime.strptime(time[0], '%Y-%m-%d')
        will = datetime.datetime.strptime(dday, '%Y-%m-%d')
        
        e_data = old - will
        
        if re.search('^-', str(e_data.days)):
            e_day = str(e_data.days)
        else:
            e_day = '+' + str(e_data.days)

        data = re.sub('\[dday\(([0-9]{4}-[0-9]{2}-[0-9]{2})\)\]', e_day, data, 1)

    # 유튜브, 카카오 티비 처리
    while 1:
        video = re.search('\[(youtube|kakaotv)\(((?:(?!\)\]).)+)\)\]', data)
        if video:
            video = video.groups()
            
            width = re.search(', ?width=((?:(?!,).)+)', video[1])
            if width:
                video_width = width.groups()[0]
            else:
                video_width = '560'
            
            height = re.search(', ?height=((?:(?!,).)+)', video[1])
            if height:
                video_height = height.groups()[0]
            else:
                video_height = '315'

            code = re.search('^(((?!,).)+)', video[1])
            if code:
                video_code = code.groups()[0]
            else:
                if video[0] == 'youtube':
                    video_code = 'BQ5PcIUcdUE'
                else:
                    video_code = '66861302'

            if video[0] == 'youtube':
                video_code = re.sub('^https:\/\/www\.youtube\.com\/watch\?v=', '', video_code)
                video_code = re.sub('^https:\/\/youtu\.be\/', '', video_code)
                
                video_src = 'https://www.youtube.com/embed/' + video_code
            else:
                video_code = re.sub('^https:\/\/tv\.kakao\.com\/channel\/9262\/cliplink\/', '', video_code)
                video_code = re.sub('^http:\/\/tv\.kakao\.com\/v\/', '', video_code)
                
                video_src = 'https://tv.kakao.com/embed/player/cliplink/' + video_code +'?service=kakao_tv'
                
            data = re.sub('\[(youtube|kakaotv)\(((?:(?!\)\]).)+)\)\]', '<iframe width="' + video_width + '" height="' + video_height + '" src="' + video_src + '" allowfullscreen frameborder="0"></iframe>', data, 1)
        else:
            break

    # 인용문 구현
    while 1:
        block = re.search('(\n(?:&gt; ?(?:(?:(?!\n).)+)?\n)+)', data)
        if block:
            block = block.groups()[0]
            
            block = re.sub('^\n&gt; ?', '', block)
            block = re.sub('\n&gt; ?', '\n', block)
            block = re.sub('\n$', '', block)
            
            data = re.sub('(\n(?:&gt; ?(?:(?:(?!\n).)+)?\n)+)', '\n<blockquote>' + block + '</blockquote>\n', data, 1)
        else:
            break

    # 리스트 구현
    while 1:
        li = re.search('(\n(?:(?: *)\* ?(?:(?:(?!\n).)+)\n)+)', data)
        if li:
            li = li.groups()[0]
            while 1:
                sub_li = re.search('\n(?:( *)\* ?((?:(?!\n).)+))', li)
                if sub_li:
                    sub_li = sub_li.groups()

                    # 앞의 공백 만큼 margin 먹임
                    if len(sub_li[0]) == 0:
                        margin = 20
                    else:
                        margin = len(sub_li[0]) * 20

                    li = re.sub('\n(?:( *)\* ?((?:(?!\n).)+))', '<li style="margin-left: ' + str(margin) + 'px;">' + sub_li[1] + '</li>', li, 1)
                else:
                    break

            data = re.sub('(\n(?:(?: *)\* ?(?:(?:(?!\n).)+)\n)+)', '\n\n<ul>' + li + '</ul>\n', data, 1)
        else:
            break

    # 들여쓰기 구현
    while 1:
        indent = re.search('\n( +)', data)
        if indent:
            indent = len(indent.groups()[0])
            
            # 앞에 공백 만큼 margin 먹임
            margin = '<span style="margin-left: 20px;"></span>' * indent
            
            data = re.sub('\n( +)', '\n' + margin, data, 1)
        else:
            break

    # 표 처리
    data = table_start(data)

    # 하위 문서
    while 1:
        under_link = re.search('\[\[\.\.\/(?:\|((?:(?!]]).)+))?]]', data)
        if under_link:
            under_link = under_link.groups()

            title_data = re.search('^(.+)\/(?:(?:(?!\/).)+)$', title)
            if title_data:
                title_data = title_data.groups()[0]

                if under_link[0]:
                    data = re.sub('\[\[\.\.\/(?:\|((?:(?!]]).)+))?]]', '[[' + title_data + '|' + under_link[0] + ']]', data, 1)
                else:
                    data = re.sub('\[\[\.\.\/(?:\|((?:(?!]]).)+))?]]', '[[' + title_data + ']]', data, 1)
            else:
                if under_link[0]:
                    data = re.sub('\[\[\.\.\/(?:\|((?:(?!]]).)+))?]]', '[[' + title + '|' + under_link[0] + ']]', data, 1)
                else:
                    data = re.sub('\[\[\.\.\/(?:\|((?:(?!]]).)+))?]]', '[[' + title + ']]', data, 1)
        else:
            break

    # 링크 관련 문법 구현
    category = '\n<hr><div id="cate">분류: '
    while 1:
        link = re.search('\[\[((?:(?!\[\[|\]\]).)+)\]\]', data)
        if link:
            link = link.groups()[0]
            
            link_split = re.search('((?:(?!\|).)+)(?:\|((?:(?!\|).)+))', link)
            if link_split:
                link_split = link_split.groups()
                
                main_link = link_split[0]
                see_link = link_split[1]
            else:
                main_link = link
                see_link = link

            if re.search('^(파일|외부):', main_link):
                file_style = ''

                width = re.search('width=((?:(?!&).)+)', see_link)
                if width:
                    file_width = width.groups()[0]
                    if re.search('px$', file_width):
                        file_style += 'width: ' + file_width + ';'
                    else:
                        file_style += 'width: ' + file_width + 'px;'
                
                height = re.search('height=((?:(?!&).)+)', see_link)
                if height:
                    file_height = height.groups()[0]
                    if re.search('px$', file_height):
                        file_style += 'height: ' + file_height + ';'
                    else:
                        file_style += 'height: ' + file_height + 'px;'

                align = re.search('align=((?:(?!&).)+)', see_link)
                if align:
                    file_align = align.groups()[0]
                    if file_align == 'center':
                        file_align = 'display: block; text-align: center;'
                    else:
                        file_align = 'float: ' + file_align + ';'
                else:
                    file_align = ''

                if re.search('^외부:', main_link):
                    file_src = re.sub('^외부:', '', main_link)
            
                    file_alt = main_link
                    exist = 'Yes'
                else:
                    file_data = re.search('^파일:((?:(?!\.).)+)\.(.+)$', main_link)
                    if file_data:
                        file_data = file_data.groups()
                        file_name = file_data[0]
                        file_end = file_data[1]

                        backlink += [[title, main_link, 'file']]
                    else:
                        file_name = 'TEST'
                        file_end = 'jpg'

                    file_src = '/image/' + tool.sha224(file_name) + '.' + file_end
                    file_alt = '파일:' + file_name + '.' + file_end

                    curs.execute("select title from data where title = ?", [file_alt])
                    exist = curs.fetchall()
                
                if exist:
                    data = re.sub('\[\[((?:(?!\[\[|\]\]).)+)\]\]', '<span style="' + file_align + '"><img style="' + file_style + '" alt="' + file_alt + '" src="' + file_src + '"></span>', data, 1)
                else:
                    data = re.sub('\[\[((?:(?!\[\[|\]\]).)+)\]\]', '<a id="not_thing" href="/w/' + tool.url_pas(file_alt) + '">' + file_alt + '</a>', data, 1)
            elif re.search('^분류:', main_link):
                see_link = re.sub('#include', '', see_link)
                main_link = re.sub('#include', '', main_link)

                if re.search('#blur', main_link):
                    see_link = '스포일러'
                    link_id = 'id="inside"'
                    
                    main_link = re.sub('#blur', '', main_link)
                else:
                    link_id = ''

                backlink += [[title, main_link, 'cat']]

                category += '<a ' + link_id + ' href="' + tool.url_pas(main_link) + '">' + re.sub('^분류:', '', see_link) + '</a> / '
                data = re.sub('\[\[((?:(?!\[\[|\]\]).)+)\]\]', '', data, 1)
            elif re.search('^wiki:', main_link):
                data = re.sub('\[\[((?:(?!\[\[|\]\]).)+)\]\]', '<a id="inside" href="/' + tool.url_pas(re.sub('^wiki:', '', main_link)) + '">' + see_link + '</a>', data, 1)
            elif re.search('^inter:((?:(?!:).)+):', main_link):
                inter_data = re.search('^inter:((?:(?!:).)+):((?:(?!\]\]|\|).)+)', main_link)
                inter_data = inter_data.groups()

                curs.execute('select link from inter where title = ?', [inter_data[0]])
                inter = curs.fetchall()
                if inter:
                    if see_link != main_link:
                        data = re.sub('\[\[((?:(?!\[\[|\]\]).)+)\]\]', '<a id="inside" href="' + inter[0][0] + inter_data[1] + '">' + inter_data[0] + ':' + see_link + '</a>', data, 1)
                    else:
                        data = re.sub('\[\[((?:(?!\[\[|\]\]).)+)\]\]', '<a id="inside" href="' + inter[0][0] + inter_data[1] + '">' + inter_data[0] + ':' + inter_data[1] + '</a>', data, 1)
                else:
                    data = re.sub('\[\[((?:(?!\[\[|\]\]).)+)\]\]', '인터위키 정보 없음', data, 1)
            elif re.search('^\/', main_link):
                under_title = re.search('^(\/(?:.+))$', main_link)
                under_title = under_title.groups()[0]

                if see_link != main_link:
                    data = re.sub('\[\[((?:(?!\[\[|\]\]).)+)\]\]', '[[' + title + under_title + '|' + see_link + ']]', data, 1)
                else:
                    data = re.sub('\[\[((?:(?!\[\[|\]\]).)+)\]\]', '[[' + title + under_title + ']]', data, 1)
            elif re.search('^http(s)?:\/\/', main_link):
                data = re.sub('\[\[((?:(?!\[\[|\]\]).)+)\]\]', '<a id="out_link" rel="nofollow" href="' + main_link + '">' + see_link + '</a>', data, 1)
            else:
                return_link = link_fix(main_link)
                main_link = return_link[0]
                other_link = return_link[1]
                
                if not re.search('^\|', main_link):
                    if main_link != title:
                        if main_link != '':
                            curs.execute("select title from data where title = ?", [main_link])
                            if not curs.fetchall():
                                link_id = 'id="not_thing"'

                                backlink += [[title, main_link, 'no']]
                            else:
                                link_id = ''
                        
                            backlink += [[title, main_link, '']]
                            
                            data = re.sub('\[\[((?:(?!\[\[|\]\]).)+)\]\]', '<a ' + link_id + ' href="/w/' + tool.url_pas(main_link) + other_link + '">' + see_link + '</a>', data, 1)
                        else:
                            data = re.sub('\[\[((?:(?!\[\[|\]\]).)+)\]\]', '<a href="' + other_link + '">' + see_link + '</a>', data, 1)
                    else:
                        data = re.sub('\[\[((?:(?!\[\[|\]\]).)+)\]\]', '<b>' + see_link + '</b>', data, 1)
                else:
                    data = re.sub('\[\[((?:(?!\[\[|\]\]).)+)\]\]', '&#91;&#91;' + link + '&#93;&#93;', data, 1)
        else:
            break

    # br 처리
    data = re.sub("\[br\]", '<br>', data)
            
    # 각주 처리
    footnote_number = 0
    footnote_all = '\n<hr><ul id="footnote_data">'
    footnote_dict = {}
    while 1:
        footnote = re.search('(?:\[\*((?:(?! |\]).)*)(?: ((?:(?!\]).)+))?\]|(\[각주\]))', data)
        if footnote:
            footnote_data = footnote.groups()
            if footnote_data[2]:
                footnote_all += '</ul>'
                
                data = re.sub('(?:\[\*((?:(?! ).)*) ((?:(?!\]).)+)\]|(\[각주\]))', footnote_all, data, 1)
                
                footnote_all = '\n<hr><ul id="footnote_data">'
            else:
                footnote = footnote_data[1]
                footnote_name = footnote_data[0]
                if footnote_name and not footnote:
                    data = re.sub('(?:\[\*((?:(?! |\]).)*)(?: ((?:(?!\]).)+))?\]|(\[각주\]))', '<sup><a href="#fn-' + footnote_dict[footnote_name] + '" id="rfn-' + footnote_dict[footnote_name] + '">(' + footnote_name + ')</a></sup>', data, 1)
                else:
                    footnote_number += 1

                    if not footnote_name:
                        footnote_name = str(footnote_number)
                    else:
                        footnote_dict.update({ footnote_name : str(footnote_number) })

                    footnote_all += '<li><a href="#rfn-' + str(footnote_number) + '" id="fn-' + str(footnote_number) + '">(' + footnote_name + ')</a> ' + footnote + '</li>'
                    
                    data = re.sub('(?:\[\*((?:(?! |\]).)*)(?: ((?:(?!\]).)+))?\]|(\[각주\]))', '<sup><a href="#fn-' + str(footnote_number) + '" id="rfn-' + str(footnote_number) + '">(' + footnote_name + ')</a></sup>', data, 1)
        else:
            break

    data = re.sub('\n+$', '', data)

    footnote_all += '</ul>'
    if footnote_all == '\n<hr><ul id="footnote_data"></ul>':
        footnote_all = ''

    
    data = re.sub('\n$', footnote_all, data + '\n', 1)

    # 분류 마지막 처리
    category += '</div>'
    category = re.sub(' / <\/div>$', '</div>', category)

    if category == '\n<hr><div id="cate">분류: </div>':
        category = ''

    data += category
    
    # NoWiki 마지막 처리
    for re_data in end_data:
        print(end_data)
        if re.search('\n', re_data[1]):
            data = data.replace('<span id="' + re_data[0] + '"></span>', '<pre>' + re.sub('^\n', '', re_data[1]) + '</pre>')
        else:
            data = data.replace('<span id="' + re_data[0] + '"></span>', '<code>' + re_data[1] + '</code>')

        data = data.replace(tool.url_pas('<span id="' + re_data[0] + '"></span>'), tool.url_pas(re_data[1]))

    while 1:
        if re.search('href="((?:(?!").)+)">((?:(?!>).)*)<code>((?:(?!<).)+)</code>', data):
            data = re.sub('href="(?P<one>(?:(?!").)+)">(?P<two>(?:(?!>).)*)<code>(?P<three>(?:(?!<).)+)</code>', 'href="\g<one>">\g<two>\g<three>', data, 1)
        else:
            break
    
    # 마지막 처리
    data = re.sub('<\/td_end>', '</td>', data)
    data = re.sub('<include>\n', '', data)
    data = re.sub('\n<\/include>', '', data)

    data = re.sub('(?P<in><\/h[0-9]>)(\n)+', '\g<in>', data)
    data = re.sub('\n\n<ul>', '\n<ul>', data)
    data = re.sub('<\/ul>\n\n', '</ul>\n', data)
    data = re.sub('^(\n)+', '', data)
    data = re.sub('(\n)+<hr><ul id="footnote_data">', '<hr><ul id="footnote_data">', data)
    data = re.sub('(?P<in><td(((?!>).)*)>)\n', '\g<in>', data)
    data = re.sub('(\n)?<hr>(\n)?', '<hr>', data)
    data = re.sub('<\/ul>\n\n<ul>', '</ul>\n<ul>', data)
    data = re.sub('\n', '<br>', data)

    return [data, plus_data, backlink]
