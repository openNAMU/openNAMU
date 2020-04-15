from . import tool

import datetime
import html
import re

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

    class head_render:
        def __init__(self):
            pass

        def __call__(self, match):
            head_len = str(len(match[1]))
            head_data = match[2]

            return '<h' + head_len + '>' + head_data + '</h' + head_len + '>'

    head_r = r'\n(#{1,6}) ?([^\n]+)'
    head_do = head_render()
    data = re.sub(head_r, head_do, data)

    class link_render:
        def __init__(self):
            pass

        def __call__(self, match):
            if re.search(r'^http(s)?:\/\/', match[2], flags = re.I):
                return '<a id="out_link" href="' + match[2] + '">' + match[1] + '</a>'
            else:
                return '<a class="' + include_num + 'link_finder" href="/w/' + match[2] + '">' + match[1] + '</a>'

    link_r = r'\[((?:(?!\]\().)+)\]\(([^\]]+)\)'
    link_do = link_render()
    data = re.sub(link_r, link_do, data)
    
    return [data, plus_data, backlink]