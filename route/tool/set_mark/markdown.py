from . import tool

import datetime
import html
import re

class head_render:
    def __init__(self):
        pass

    def __call__(self, match):
        head_len = str(len(match[1]))
        head_data = match[2]

        return '<h' + head_len + '>' + head_data + '</h' + head_len + '>'

class link_render:
    def __init__(self):
        pass

    def __call__(self, match):
        if match[1] == '!':
            if re.search(r'^http(s)?:\/\/', match[3], flags = re.I):
                return '<img alt="' + match[2] + '" src="' + match[3] + '">'
            else:
                file_name = re.search(r'^([^.]+)\.([^.]+)$', match[3])
                if file_name:
                    file_end = file_name.group(2)
                    file_name = file_name.group(1)
                else:
                    file_name = 'Test'
                    file_end = 'jpg'

                file_src = '/image/' + tool.sha224_replace(file_name) + '.' + file_end
                file_alt = 'file:' + file_name + '.' + file_end

                return '' + \
                    '<img class="' + include_num + 'file_finder_1" alt="' + match[2] + '" src="' + file_src + '">' + \
                    '<a class="' + include_num + 'file_finder_2" id="not_thing" href="/upload?name=' + tool.url_pas(file_name) + '">' + file_alt + '</a>' + \
                ''
        else:
            if re.search(r'^http(s)?:\/\/', match[3], flags = re.I):
                return '<a id="out_link" href="' + match[3] + '">' + match[2] + '</a>'
            else:
                return '<a class="' + include_num + 'link_finder" href="/w/' + match[3] + '">' + match[2] + '</a>'

def markdown(conn, data, title, include_num):
    backlink = []
    include_num = include_num + '_' if include_num else ''
    plus_data = '' + \
        'get_link_state("' + include_num + '");\n' + \
        'get_file_state("' + include_num + '");\n' + \
    ''

    data = html.escape(data)
    data = data.replace('\r\n', '\n')
    data = '\n' + data

    head_r = r'\n(#{1,6}) ?([^\n]+)'
    head_do = head_render()
    data = re.sub(head_r, head_do, data)

    link_r = r'(!)?\[((?:(?!\]\().)+)\]\(([^\]]+)\)'
    link_do = link_render()
    data = re.sub(link_r, link_do, data)

    data = re.sub(r'\*\*((?:(?!\*\*).)+)\*\*', '<b>\1</b>', data)
    data = re.sub(r'__((?:(?!__).)+)__', '<i>\1</i>', data)

    
    return [data, plus_data, backlink]