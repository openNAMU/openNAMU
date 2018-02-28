from . import tool

import re

def table_parser(data, cel_data, start_data, num = 0):
    table_class = 'class="'
    all_table = 'style="'
    cel_style = 'style="'
    row_style = 'style="'
    row = ''
    cel = ''

    table_width = re.search("<table ?width=((?:(?!>).)*)>", data)
    if table_width:
        all_table += 'width: ' + table_width.groups()[0] + ';'
    
    table_height = re.search("<table ?height=((?:(?!>).)*)>", data)
    if table_height:
        all_table += 'height: ' + table_height.groups()[0] + ';'
    
    table_align = re.search("<table ?align=((?:(?!>).)*)>", data)
    if table_align:
        if table_align.groups()[0] == 'right':
            all_table += 'float: right;'
        elif table_align.groups()[0] == 'center':
            all_table += 'margin: auto;'
            
    table_text_align = re.search("<table ?textalign=((?:(?!>).)*)>", data)
    if table_text_align:
        num = 1
        if table_text_align.groups()[0] == 'right':
            all_table += 'text-align: right;'
        elif table_text_align.groups()[0] == 'center':
            all_table += 'text-align: center;'

    row_t_a = re.search("<row ?textalign=((?:(?!>).)*)>", data)
    if row_t_a:
        if row_t_a.groups()[0] == 'right':
            row_style += 'text-align: right;'
        elif row_t_a.groups()[0] == 'center':
            row_style += 'text-align: center;'
        else:
            row_style += 'text-align: left;'
    
    table_cel = re.search("<-((?:(?!>).)*)>", data)
    if table_cel:
        cel = 'colspan="' + table_cel.groups()[0] + '"'
    else:
        cel = 'colspan="' + str(round(len(start_data) / 2)) + '"'   

    table_row = re.search("<\|((?:(?!>).)*)>", data)
    if table_row:
        row = 'rowspan="' + table_row.groups()[0] + '"'

    row_bgcolor = re.search("<rowbgcolor=(#(?:[0-9a-f-A-F]{3}){1,2}|\w+)>", data)
    if row_bgcolor:
        row_style += 'background: ' + row_bgcolor.groups()[0] + ';'
        
    table_border = re.search("<table ?bordercolor=(#(?:[0-9a-f-A-F]{3}){1,2}|\w+)>", data)
    if table_border:
        all_table += 'border: ' + table_border.groups()[0] + ' 2px solid;'
        
    table_bgcolor = re.search("<table ?bgcolor=(#(?:[0-9a-f-A-F]{3}){1,2}|\w+)>", data)
    if table_bgcolor:
        all_table += 'background: ' + table_bgcolor.groups()[0] + ';'
        
    bgcolor = re.search("<(?:bgcolor=)?(#(?:[0-9a-f-A-F]{3}){1,2}|\w+)>", data)
    if bgcolor:
        cel_style += 'background: ' + bgcolor.groups()[0] + ';'
        
    cel_width = re.search("<width=((?:(?!>).)*)>", data)
    if cel_width:
        cel_style += 'width: ' + cel_width.groups()[0] + ';'

    cel_height = re.search("<height=((?:(?!>).)*)>", data)
    if cel_height:
        cel_style += 'height: ' + cel_height.groups()[0] + ';'
        
    text_right = re.search("<\)>", data)
    text_center = re.search("<:>", data)
    text_left = re.search("<\(>",  data)
    if text_right:
        cel_style += 'text-align: right;'
    elif text_center:
        cel_style += 'text-align: center;'
    elif text_left:
        cel_style += 'text-align: left;'
    elif num == 0:
        if re.search('^ (.*) $', cel_data):
            cel_style += 'text-align: center;'
        elif re.search('^ (.*)$', cel_data):
            cel_style += 'text-align: right;'
        elif re.search('^(.*) $', cel_data):
            cel_style += 'text-align: left;'

    text_class = re.search("<table ?class=((?:(?!>).)+)>", data)
    if text_class:
        table_class += text_class.groups()[0]
        
    all_table += '"'
    cel_style += '"'
    row_style += '"'
    table_class += '"'

    return [all_table, row_style, cel_style, row, cel, table_class, num]

def start(conn, data, title):
    # DB 지정
    curs = conn.cursor()

    # 맨 앞과 끝에 개행 문자 추가
    data = '\r\n' + data + '\r\n'

    while 1:
        include = re.search('\[include\(((?:(?!\)\]).)+)\)\]', data)
        if include:
            include = include.groups()[0]

            include_data = re.search('^((?:(?!,).)+)', include)
            if include_data:
                include_data = include_data.groups()[0]
            else:
                include_data = 'Test'

            include = re.sub('^((?:(?!,).)+)', '', include)

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

                data = re.sub('\[include\(((?:(?!\)\]).)+)\)\]', '\r\n<span id="include"></span>' + include_parser + '<span id="include"></span>\r\n', data, 1)
            else:
                data = re.sub('\[include\(((?:(?!\)\]).)+)\)\]', '[[' + include + ']]', data, 1)

        else:
            break

    # 텍스트 꾸미기 문법
    data = re.sub('\'\'\'(?P<in>(?:(?!\'\'\').)+)\'\'\'', '<b>\g<in></b>', data)
    data = re.sub('\'\'(?P<in>(?:(?!\'\').)+)\'\'', '<i>\g<in></i>', data)

    data = re.sub('~~(?P<in>(?:(?!~~).)+)~~', '<s>\g<in></s>', data)
    data = re.sub('--(?P<in>(?:(?!~~).)+)--', '<s>\g<in></s>', data)

    data = re.sub('__(?P<in>(?:(?!__).)+)__', '<u>\g<in></u>', data)
    
    data = re.sub('^^(?P<in>(?:(?!^^).)+)^^', '<sup>\g<in></sup>', data)
    data = re.sub(',,(?P<in>(?:(?!,,).)+),,', '<sub>\g<in></sub>', data)

    # 넘겨주기 변환
    data = re.sub('\r\n#(?:redirect|넘겨주기) (?P<in>(?:(?!\r\n).)+)\r\n', '<meta http-equiv="refresh" content="0; url=/w/\g<in>?froms=' + tool.url_pas(title) + '">', data)

    # 각주 처리
    footnote_number = 0
    footnote_all = '\r\n<hr><ul id="footnote_data">'
    while 1:
        footnote = re.search('(?:\[\*((?:(?! ).)*) ((?:(?!\]).)+)\]|(\[각주\]))', data)
        if footnote:
            footnote_data = footnote.groups()

            if footnote_data[2]:
                footnote_all += '</ul>'

                data = re.sub('(?:\[\*((?:(?! ).)*) ((?:(?!\]).)+)\]|(\[각주\]))', footnote_all, data, 1)

                footnote_all = '\r\n<hr><ul id="footnote_data">'
            else:
                footnote = footnote_data[1]
                footnote_name = footnote_data[0]

                footnote_number += 1

                if not footnote_name:
                    footnote_name = str(footnote_number)

                footnote_all += '<li><a href="#rfn-' + str(footnote_number) + '" id="fn-' + str(footnote_number) + '">(' + footnote_name + ')</a> ' + footnote + '</li>'

                data = re.sub('(?:\[\*((?:(?! ).)*) ((?:(?!\]).)+)\]|(\[각주\]))', '<sup><a href="#fn-' + str(footnote_number) + '" id="rfn-' + str(footnote_number) + '">(' + footnote_name + ')</a></sup>', data, 1)
        else:
            break

    footnote_all += '</ul>'

    if footnote_all == '\r\n<hr><ul id="footnote_data"></ul>':
        footnote_all = ''

    data = re.sub('\r\n$', footnote_all, data)

    # [목차(없음)] 처리
    if not re.search('\[목차\(없음\)\]\r\n', data):
        if not re.search('\[목차\]', data):
            data = re.sub('\r\n(?P<in>={1,6}) ?(?P<out>(?:(?!=).)+) ?={1,6}\r\n', '\r\n[목차]\r\n\g<in> \g<out> \g<in>\r\n', data, 1)
    else:
        data = re.sub('\[목차\(없음\)\]\r\n', '', data)

    # 문단 문법
    toc_full = 0
    toc_top_stack = 6
    toc_stack = [0, 0, 0, 0, 0, 0]
    edit_number = 0
    toc_data = '<div id="toc"><span style="font-size: 18px;">목차</span>\r\n\r\n'
    while 1:
        toc = re.search('\r\n(={1,6}) ?((?:(?!=).)+) ?={1,6}\r\n', data)
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
            
            data = re.sub('\r\n(={1,6}) ?((?:(?!=).)+) ?={1,6}\r\n', '\r\n<h' + toc_number + '><a href="">' + all_stack + '</a> ' + toc[1] + ' <span style="font-size: 12px"><a href="/edit/' + tool.url_pas(title) + '?section=' + str(edit_number) + '">(편집)</a></span></h' + toc_number + '><hr id="under_bar" style="margin-top: -5px; margin-bottom: 10px;">\r\n', data, 1)
            toc_data += '<span style="margin-left: ' + str((toc_full - toc_top_stack) * 10) + 'px"><a href="">' + all_stack + '</a> ' + toc[1] + '</span>\r\n'
        else:
            break

    toc_data += '</div>'
    
    data = re.sub('\[목차\]', toc_data, data)

    while 1:
        hr = re.search('\r\n-{4,9}\r\n', data)
        if hr:
            data = re.sub('\r\n-{4,9}\r\n', '<hr>', data, 1)
        else:
            break

    data += '\r\n'

    # 일부 매크로 처리
    data = tool.savemark(data)

    data = re.sub("\[br\]", '\r\n', data)
    data = re.sub("\[anchor\((?P<in>(?:(?!\)\]).)+)\)\]", '<span id="\g<in>"></span>', data)          
    data = re.sub("\[nicovideo\((?P<in>(?:(?!,|\)\]).)+)(?:(?:(?!\)\]).)*)\)\]", "[[http://embed.nicovideo.jp/watch/\g<in>|\g<in>]]", data)
    data = re.sub('\[ruby\((?P<in>(?:(?!,).)+)\, ?(?P<out>(?:(?!\)\]).)+)\)\]', '<ruby>\g<in><rp>(</rp><rt>\g<out></rt><rp>)</rp></ruby>', data)

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
                
            data = re.sub('\[(youtube|kakaotv)\(((?:(?!\)\]).)+)\)\]', '<iframe width="' + video_width + '" height="' + video_height + '" src="' + video_src + '" allowfullscreen></iframe>', data, 1)
        else:
            break

    # 인용문 구현
    while 1:
        block = re.search('(\r\n(?:> ?(?:(?:(?!\r\n).)+)\r\n)+)', data)
        if block:
            block = block.groups()[0]

            block = re.sub('^\r\n> ?', '', block)
            block = re.sub('\r\n> ?', '\r\n', block)
            block = re.sub('\r\n$', '', block)

            data = re.sub('(\r\n(?:> ?(?:(?:(?!\r\n).)+)\r\n)+)', '<blockquote>' + block + '</blockquote>\r\n', data, 1)
        else:
            break

    # 리스트 구현
    while 1:
        li = re.search('(\r\n(?:(?: *)\* ?(?:(?:(?!\r\n).)+)\r\n)+)', data)
        if li:
            li = li.groups()[0]

            while 1:
                sub_li = re.search('\r\n(?:( *)\* ?((?:(?!\r\n).)+))', li)
                if sub_li:
                    sub_li = sub_li.groups()

                    # 앞의 공백 만큼 margin 먹임
                    if len(sub_li[0]) == 0:
                        margin = 20
                    else:
                        margin = len(sub_li[0]) * 20

                    li = re.sub('\r\n(?:( *)\* ?((?:(?!\r\n).)+))', '<li style="margin-left: ' + str(margin) + 'px">' + sub_li[1] + '</li>', li, 1)
                else:
                    break

            data = re.sub('(\r\n(?:(?: *)\* ?(?:(?:(?!\r\n).)+)\r\n)+)', '\r\n\r\n<ul>' + li + '</ul>\r\n', data, 1)
        else:
            break

    # 들여쓰기 구현
    while 1:
        indent = re.search('\r\n( +)', data)
        if indent:
            indent = len(indent.groups()[0])
            
            # 앞에 공백 만큼 margin 먹임
            margin = '<span style="margin-left: 20px;"></span>' * indent

            data = re.sub('\r\n( +)', '\r\n' + margin, data, 1)
        else:
            break

    # 표 처리
    while 1:
        table = re.search('((?:(?:(?:(?:\|\|)+(?:(?:(?!\|\|).(?:\r\n)*)+))+)\|\|(?:\r\n)?)+)', data)
        if table:
            table = table.groups()[0]
            
            # return [all_table, row_style, cel_style, row, cel, table_class, num]
            while 1:
                all_table = re.search('^((?:\|\|)+)((?:<(?:(?:(?!>).)+)>)*)((?:(?!\|\||<\/td>).)+)', table)
                if all_table:
                    all_table = all_table.groups()

                    return_table = table_parser(all_table[1], all_table[2], all_table[0])
                    number = return_table[6]

                    table = re.sub('^\|\|((?:<(?:(?:(?!>).)+)>)*)', '<table ' + return_table[5] + ' ' + return_table[0] + '><tbody><tr ' + return_table[1] + '><td ' + return_table[2] + ' ' + return_table[3] + ' ' + return_table[4] + '>', table, 1)
                else:
                    break

            table = re.sub('\|\|\r\n$', '</td></tr></tbody></table>', table)

            while 1:
                row_table = re.search('\|\|\r\n((?:\|\|)+)((?:<(?:(?:(?!>).)+)>)*)((?:(?!\|\||<\/td>).)+)', table)
                if row_table:
                    row_table = row_table.groups()

                    return_table = table_parser(row_table[1], row_table[2], row_table[0], number)

                    table = re.sub('\|\|\r\n((?:\|\|)+)((?:<(?:(?:(?!>).)+)>)*)', '</td></tr><tr ' + return_table[1] + '><td ' + return_table[2] + ' ' + return_table[3] + ' ' + return_table[4] + '>', table, 1)
                else:
                    break

            while 1:
                cel_table = re.search('((?:\|\|)+)((?:<(?:(?:(?!>).)+)>)*)((?:(?!\|\||<\/td>).)+)', table)
                if cel_table:
                    cel_table = cel_table.groups()

                    return_table = table_parser(cel_table[1], cel_table[2], cel_table[0], number)

                    table = re.sub('((?:\|\|)+)((?:<(?:(?:(?!>).)+)>)*)', '</td><td ' + return_table[2] + ' ' + return_table[3] + ' ' + return_table[4] + '>', table, 1)
                else:
                    break

            data = re.sub('((?:(?:(?:(?:\|\|)+(?:(?:(?!\|\|).(?:\r\n)*)+))+)\|\|(?:\r\n)?)+)', table, data, 1)
        else:
            break

    # 링크 관련 문법 구현
    category = '\r\n<div id="cate">분류: '
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
                width = re.search('width=((?:(?!,).)+)', see_link)
                if width:
                    file_width = width.groups()[0]
                else:
                    file_width = 'auto'
                
                height = re.search('height=((?:(?!,).)+)', see_link)
                if height:
                    file_height = height.groups()[0]
                else:
                    file_height = 'auto'

                align = re.search('align=((?:(?!,).)+)', see_link)
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
                else:
                    file_data = re.search('^파일:((?:(?!\.).)+)\.(.+)$', main_link)
                    if file_data:
                        file_data = file_data.groups()

                        file_name = file_data[0]
                        file_end = file_data[1]
                    else:
                        file_name = 'TEST'
                        file_end = 'jpg'

                    file_src = '/image/' + tool.sha224(file_name) + '.' + file_end
                    file_alt = '파일:' + file_name + '.' + file_end
                
                data = re.sub('\[\[((?:(?!\[\[|\]\]).)+)\]\]', '<span style="' + file_align + '"><img width="' + file_width + '" height="' + file_height + '" alt="' + file_alt + '" src="' + file_src + '"></span>', data, 1)
            elif re.search('^분류:', main_link):
                see_link = re.sub('#include', '', see_link)
                main_link = re.sub('#include', '', main_link)
                
                if re.search('#blur', main_link):
                    see_link = '스포일러'
                    link_id = 'id="inside"'

                    main_link = re.sub('#blur', '', main_link)
                else:
                    link_id = ''

                category += '<a ' + link_id + ' href="' + tool.url_pas(main_link) + '">' + re.sub('^분류:', '', see_link) + '</a> / '

                data = re.sub('\[\[((?:(?!\[\[|\]\]).)+)\]\]', '', data, 1)
            elif re.search('^wiki:', main_link):
                data = re.sub('\[\[((?:(?!\[\[|\]\]).)+)\]\]', '<a id="inside" href="/' + tool.url_pas(re.sub('^wiki:', '', main_link)) + '">' + see_link + '</a>', data, 1)
            elif re.search('^http(s)?:\/\/', main_link):
                data = re.sub('\[\[((?:(?!\[\[|\]\]).)+)\]\]', '<a class="out_link" rel="nofollow" href="' + main_link + '">' + see_link + '</a>', data, 1)
            else:
                if re.search('^:', main_link):
                    main_link = re.sub('^:', '', main_link)

                curs.execute("select title from data where title = ?", [main_link])
                if not curs.fetchall():
                    link_class = 'class="not_thing"'
                else:
                    link_class = ''

                data = re.sub('\[\[((?:(?!\[\[|\]\]).)+)\]\]', '<a ' + link_class + ' href="/w/' + tool.url_pas(main_link) + '">' + see_link + '</a>', data, 1)
        else:
            break

    category += '</div>'
    category = re.sub(' / <\/div>$', '</div>', category)

    if category == '\r\n<div id="cate">분류: </div>':
        category = ''

    data += category
    
    # 마지막 처리
    data = re.sub('(?P<in><hr id="under_bar" style="margin-top: -5px; margin-bottom: 10px;">)(\r\n)+', '\g<in>', data)
    data = re.sub('<\/ul>\r\n\r\n', '</ul>\r\n', data)
    data = re.sub('^(\r\n)+', '', data)
    data = re.sub('(\r\n)+$', '', data)
    data = re.sub('(\r\n)?<span id="include"><\/span>(\r\n)?(<span style="margin-left: 20px;"><\/span>)?', '', data)
    data = re.sub('\r\n', '<br>', data)

    return data