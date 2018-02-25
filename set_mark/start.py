from . import tool

import re

def start(data, title):
    # 맨 앞과 끝에 개행 문자 추가
    data = '\r\n' + data + '\r\n'

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

    # 문단 문법
    toc_full = 0
    toc_top_stack = 6
    toc_stack = [0, 0, 0, 0, 0, 0]
    toc_data = '<div id="toc"><span style="font-size: 18px;">목차</span>\r\n\r\n'
    while 1:
        toc = re.search('\r\n(={1,6}) ?((?:(?!=).)+) ?={1,6}\r\n', data)
        if toc:
            toc = toc.groups()
            toc_number = len(toc[0])

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
            
            data = re.sub('\r\n(={1,6}) ?((?:(?!=).)+) ?={1,6}\r\n', '\r\n<h' + toc_number + '><a href="">' + all_stack + '</a> ' + toc[1] + '</h' + toc_number + '><hr id="under_bar" style="margin-top: -5px; margin-bottom: 10px;">\r\n', data, 1)
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

            data = re.sub('(\r\n(?:(?: *)\* ?(?:(?:(?!\r\n).)+)\r\n)+)', '<ul>' + li + '</ul>', data, 1)
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
    
    # 마지막 처리
    data = re.sub('(?P<in><hr id="under_bar" style="margin-top: -5px; margin-bottom: 10px;">)\r\n', '\g<in>', data)
    data = re.sub('\r\n', '<br>', data)

    return data