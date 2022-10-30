from .func_tool import *

class class_do_render_namumark:
    def __init__(self, curs, doc_name, doc_data, doc_set, lang_data):
        self.curs = curs
        
        self.doc_data = doc_data
        self.doc_name = doc_name
        self.doc_set = doc_set
        self.doc_include = self.doc_set['doc_include'] if 'doc_include' in self.doc_set else ''

        self.lang_data = lang_data

        self.data_temp_storage = {}
        self.data_temp_storage_count = 0

        self.data_backlink = []
        self.data_include = []

        self.data_math_count = 0
        
        self.data_toc = ''
        self.data_footnote = {}
        self.data_category = ''

        self.render_data = self.doc_data
        self.render_data = html.escape(self.render_data)
        self.render_data = '<back_br>\n' + self.render_data + '\n<front_br>'
        self.render_data_js = ''

    def get_tool_lang(self, name):
        if name in self.lang_data:
            return self.lang_data[name]
        else:
            return name + ' (RENDER LANG)'

    def get_tool_js_safe(self, data):
        data = data.replace('\n', '\\\\n')
        data = data.replace('\\', '\\\\')
        data = data.replace('"', '\\"')

        return data

    def get_tool_css_safe(self, data):
        return data.replace(';', '')

    def get_tool_data_storage(self, data_A = '', data_B = '', data_C = '', do_type = 'render'):
        self.data_temp_storage_count += 1
        if do_type == 'render':
            data_name = 'render_' + str(self.data_temp_storage_count)

            self.data_temp_storage[data_name] = data_A
            self.data_temp_storage['/' + data_name] = data_B
            self.data_temp_storage['revert_' + data_name] = data_C
        else:
            data_name = 'slash_' + str(self.data_temp_storage_count)

            self.data_temp_storage[data_name] = data_A

        return data_name

    def get_tool_data_restore(self, data, do_type = 'all'):
        storage_count = self.data_temp_storage_count * 3
        if do_type == 'all':
            storage_regex = r'<(\/?(?:render|slash)_(?:[0-9]+))>'
        elif do_type == 'render':
            storage_regex = r'<(\/?(?:render)_(?:[0-9]+))>'
        else:
            storage_regex = r'<(\/?(?:slash)_(?:[0-9]+))>'

        while 1:
            if not re.search(storage_regex, data):
                break
            if storage_count < 0:
                print('Error : render restore count overflow')

                break
            else:
                data = re.sub(storage_regex, lambda match : self.data_temp_storage[match.group(1)], data, 1)

            storage_count -= 1

        return data

    def get_tool_data_revert(self, data, do_type = 'all'):
        storage_count = self.data_temp_storage_count * 3
        if do_type == 'all':
            storage_regex = r'(?:<((slash)_(?:[0-9]+))>|<((render)_(?:[0-9]+))>(?:(?:(?!<(?:\/?render_(?:[0-9]+))>).|\n)*)<\/render_(?:[0-9]+)>)'
        elif do_type == 'render':
            storage_regex = r'<((render)_(?:[0-9]+))>(?:(?:(?!<(?:\/?render_(?:[0-9]+))>).)*)<\/render_(?:[0-9]+)>'
        else:
            storage_regex = r'<((slash)_(?:[0-9]+))>'

        while 1:
            match = re.search(storage_regex, data)
            if not match:
                break
            if storage_count < 0:
                print('Error : render restore count overflow')

                break
            else:
                match = match.groups()
                if match[1] and match[1] == 'render':
                    if ('revert_' + match[0]) in self.data_temp_storage:
                        data_revert = self.data_temp_storage['revert_' + match[0]]
                    else:
                        data_revert = ''
                else:
                    if len(match) > 3 and match[3] == 'render':
                        if ('revert_' + match[2]) in self.data_temp_storage:
                            data_revert = self.data_temp_storage['revert_' + match[2]]
                        else:
                            data_revert = ''
                    else:
                        data_revert = '\\' + self.data_temp_storage[match[0]]

                data = re.sub(storage_regex, lambda x : data_revert, data, 1)

            storage_count -= 1

        data = re.sub(r'<front_br>', '', data)
        data = re.sub(r'<back_br>', '', data)

        return data

    def get_tool_footnote_make(self):
        data = ''
        for for_a in self.data_footnote:
            if data == '':
                data += '<div class="opennamu_footnote">'
            else:
                data += '<br>'

            if len(self.data_footnote[for_a]['list']) > 1:
                data += '(' + for_a + ') '

                for for_b in self.data_footnote[for_a]['list']:
                    data += '<sup><a id="' + self.doc_include + 'fn_' + for_b + '" href="#' + self.doc_include + 'rfn_' + for_b + '">(' + for_b + ')</a></sup> '
            else:
                data += '<a id="' + self.doc_include + 'fn_' + self.data_footnote[for_a]['list'][0] + '" href="#' + self.doc_include + 'rfn_' + self.data_footnote[for_a]['list'][0] + '">(' + for_a + ') </a> '

            data += self.data_footnote[for_a]['data']

        if data != '':
            data += '</div>'

        self.data_footnote = {}

        return data

    def get_tool_px_add_check(self, data):
        if re.search(r'^[0-9]+$', data):
            return data + 'px'
        else:
            return data

    def get_tool_dark_mode_split(self, data):
        data = data.split(',')
        if len(data) == 1:
            return data[0]
        else:
            if flask.request.cookies.get('main_css_del_bold', '0') == '0':
                return data[0]
            else:
                return data[1]

    def do_render_text(self):
        # <b> function
        bold_user_set = flask.request.cookies.get('main_css_del_bold', '0')
        def do_render_text_bold(match):
            data = match.group(1)
            if bold_user_set == '0':
                data_name = self.get_tool_data_storage('<b>', '</b>', match.group(0))
            elif bold_user_set == '1':
                data_name = self.get_tool_data_storage('', '', match.group(0))
            else:
                return ''
            
            return '<' + data_name + '>' + data + '</' + data_name + '>'

        # <b>
        self.render_data = re.sub(r"&#x27;&#x27;&#x27;((?:(?!&#x27;&#x27;&#x27;).)+)&#x27;&#x27;&#x27;", do_render_text_bold, self.render_data)

        # <i> function
        def do_render_text_italic(match):
            data = match.group(1)
            data_name = self.get_tool_data_storage('<i>', '</i>', match.group(0))
            
            return '<' + data_name + '>' + data + '</' + data_name + '>'

        # <i>
        self.render_data = re.sub(r"&#x27;&#x27;((?:(?!&#x27;&#x27;).)+)&#x27;&#x27;", do_render_text_italic, self.render_data)

        # <u> function
        def do_render_text_under(match):
            data = match.group(1)
            data_name = self.get_tool_data_storage('<u>', '</u>', match.group(0))
            
            return '<' + data_name + '>' + data + '</' + data_name + '>'

        # <u>
        self.render_data = re.sub(r"__((?:(?!__).)+)__", do_render_text_under, self.render_data)
        
        # <sup> function
        def do_render_text_sup(match):
            data = match.group(1)
            data_name = self.get_tool_data_storage('<sup>', '</sup>', match.group(0))
            
            return '<' + data_name + '>' + data + '</' + data_name + '>'

        # <sup>
        self.render_data = re.sub(r"\^\^\^((?:(?!\^\^\^).)+)\^\^\^", do_render_text_sup, self.render_data)
        # <sup> 2
        self.render_data = re.sub(r"\^\^((?:(?!\^\^).)+)\^\^", do_render_text_sup, self.render_data)

        # <sub> function
        def do_render_text_sub(match):
            data = match.group(1)
            data_name = self.get_tool_data_storage('<sub>', '</sub>', match.group(0))
            
            return '<' + data_name + '>' + data + '</' + data_name + '>'
        
        # <sub>
        self.render_data = re.sub(r",,,((?:(?!,,,).)+),,,", do_render_text_sub, self.render_data)
        # <sub> 2
        self.render_data = re.sub(r",,((?:(?!,,).)+),,", do_render_text_sub, self.render_data)

        # <sub> function
        strike_user_set = flask.request.cookies.get('main_css_del_strike', '0')
        def do_render_text_strike(match):
            data = match.group(1)
            if bold_user_set == '0':
                data_name = self.get_tool_data_storage('<s>', '</s>', match.group(0))
            elif bold_user_set == '1':
                data_name = self.get_tool_data_storage('', '', match.group(0))
            else:
                return ''
            
            return '<' + data_name + '>' + data + '</' + data_name + '>'
        
        # <s>
        self.render_data = re.sub(r"--((?:(?!--).)+)--", do_render_text_strike, self.render_data)
        # <s> 2
        self.render_data = re.sub(r"~~((?:(?!~~).)+)~~", do_render_text_strike, self.render_data)
    
    def do_render_heading(self):
        toc_list = []

        # make heading base
        heading_regex = r'\n((={1,6})(#?) ?([^\n]+))\n'
        heading_count_all = len(re.findall(heading_regex, self.render_data)) * 3
        heading_stack = [0, 0, 0, 0, 0, 0]
        heading_count = 0
        while 1:
            heading_count += 1

            if not re.search(heading_regex, self.render_data):
                break
            elif heading_count_all < 0:
                print('Error : render heading count overflow')

                break
            else:
                heading_data = re.search(heading_regex, self.render_data)
                heading_data = heading_data.groups()

                heading_data_last_regex = r' ?(#?={1,6}[^=]*)$'
                heading_data_last = re.search(heading_data_last_regex, heading_data[3])
                if heading_data_last:
                    heading_data_last = heading_data_last.group(1)
                else:
                    heading_data_last = ''

                heading_data_text = re.sub(heading_data_last_regex, '', heading_data[3])

                heading_data_diff = heading_data[2] + heading_data[1]
                if heading_data_diff != heading_data_last:
                    # front != back -> restore

                    heading_data_all = heading_data[0]

                    for for_a in reversed(range(1, 7)):
                        for_a_str = str(for_a)

                        heading_restore_regex = re.compile('^={' + for_a_str + '}|={' + for_a_str + '}$')

                        heading_data_all = re.sub(heading_restore_regex, '<heading_' + for_a_str + '>', heading_data_all)

                    self.render_data = re.sub(heading_regex, '\n' + heading_data_all + '\n', self.render_data, 1)
                else:
                    heading_level = len(heading_data[1])
                    heading_level_str = str(heading_level)

                    heading_stack[heading_level - 1] += 1
                    for for_a in range(heading_level, 6):
                        heading_stack[for_a] = 0

                    heading_stack_str = '.'.join([str(for_a) for for_a in heading_stack])
                    heading_stack_str = re.sub(r'(\.0)+$', '', heading_stack_str)

                    toc_list += [['', heading_data_text]]

                    heading_data_text_fix = re.sub(r'<([^<>]*)>', '', heading_data_text)
                    
                    data_name = self.get_tool_data_storage(
                        '<h' + heading_level_str + ' id="' + heading_data_text_fix + '">', 
                        ' <sub><a href="/edit_section/' + str(heading_count) + '/' + url_pas(self.doc_name) + '">✎</a></sub></h' + heading_level_str + '>', 
                        ''
                    )

                    heading_data_complete = '' + \
                        '\n<front_br>' + \
                        '<' + data_name + '>' + \
                            '<heading_stack>' + \
                                heading_stack_str + \
                            '</heading_stack>' + \
                            ' ' + heading_data_text + \
                        '</' + data_name + '>' + \
                        '<back_br>\n' + \
                    ''

                    self.render_data = re.sub(heading_regex, heading_data_complete, self.render_data, 1)

            heading_count_all -= 1

        # heading id adjust
        heading_end_count = len(re.findall(r'<heading_stack>', self.render_data))
        for for_a in reversed(range(0, 6)):
            heading_end_stack_regex = re.compile('<heading_stack>' + ('0\\.' * for_a))

            heading_end_match_count = len(re.findall(heading_end_stack_regex, self.render_data))
            if heading_end_match_count == heading_end_count:
                self.render_data = re.sub(heading_end_stack_regex, '<heading_stack>', self.render_data)

                break

        # heading id -> inline id
        heading_id_regex = r'<heading_stack>([^<>]+)<\/heading_stack>'
        heading_id_data = re.findall(heading_id_regex, self.render_data)
        for for_a in range(len(heading_id_data)):
            self.render_data = re.sub(heading_id_regex, '<a href="#toc" id="s-' + heading_id_data[for_a] + '">' + heading_id_data[for_a] + '.</a>', self.render_data, 1)
            
            toc_list[for_a][0] = heading_id_data[for_a]

        # not heading restore
        for for_a in range(1, 7):
            for_a_str = str(for_a)

            heading_restore_regex = re.compile('<heading_' + for_a_str + '>')

            self.render_data = re.sub(heading_restore_regex, ('=' * for_a), self.render_data)

        # make toc
        if len(toc_list) == 0:
            toc_data = ''
        else:
            toc_data = '' + \
                '<div class="opennamu_TOC" id="toc">' + \
                    '<span class="opennamu_TOC_title">' + self.get_tool_lang('toc') + '</span>' + \
                    '<br>' + \
                ''

        for for_a in toc_list:
            toc_data += '' + \
                '<br>' + \
                ('<span style="margin-left: 10px;">' * for_a[0].count('.')) + \
                '<span>' + \
                    '<a href="#s-' + for_a[0] + '">' + \
                        for_a[0] + '. ' + \
                    '</a>' + \
                    for_a[1] + \
                '</span>' + \
            ''

        if toc_data != '':
            toc_data += '</div>'

        # toc replace
        self.render_data = re.sub(r'\[(목차|toc|tableofcontents)\]', lambda x : toc_data, self.render_data)

    def do_render_macro(self):
        # double macro function
        def do_render_macro_double(match):
            match_org = match
            match = match.groups()

            name_data = match[0]
            macro_split_regex = r'(?:^|,) *([^,]+)'
            macro_split_sub_regex = r'(^[^=]+) *= *([^=]+)'
            if name_data in ('youtube', 'nicovideo', 'navertv', 'kakaotv', 'vimeo'):
                data = re.findall(macro_split_regex, match[1])

                # get option
                video_code = ''
                video_start = ''
                video_end = ''
                video_width = '640px'
                video_height = '360px'
                for for_a in data:
                    data_sub = re.search(macro_split_sub_regex, for_a)
                    if data_sub:
                        data_sub = data_sub.groups()
                        if data_sub[0] == 'width':
                            video_width = self.get_tool_px_add_check(data_sub[1])
                        elif data_sub[0] == 'height':
                            video_height = self.get_tool_px_add_check(data_sub[1])
                        elif data_sub[0] == 'start':
                            video_start = data_sub[1]
                        elif data_sub[0] == 'end':
                            video_end = data_sub[1]
                        elif data_sub[0] == 'https://www.youtube.com/watch?v' and name_data == 'youtube':
                            video_code = data_sub[1]
                    else:
                        video_code = for_a

                # code to url
                if name_data == 'youtube':
                    video_code = re.sub(r'^https:\/\/youtu\.be\/', '', video_code)

                    video_code = 'https://www.youtube.com/embed/' + video_code

                    if video_start != '':
                        if video_end != '':
                            video_code += '?start=' + video_start + '&end=' + video_end
                        else:
                            video_code += '?start=' + video_start
                    else:
                        if video_end != '':
                            video_code += '?end=' + video_end
                elif name_data == 'kakaotv':
                    video_code = re.sub(r'^https:\/\/tv\.kakao\.com\/v\/', '', video_code)

                    video_code = 'https://tv.kakao.com/embed/player/cliplink/' + video_code +'?service=kakao_tv'
                elif name_data == 'navertv':
                    video_code = re.sub(r'^https:\/\/tv\.naver\.com\/v\/', '', video_code)

                    video_code = 'https://tv.naver.com/embed/' + video_code
                elif name_data == 'nicoviedo':
                    video_code = 'https://embed.nicovideo.jp/watch/' + video_code
                else:
                    video_code = 'https://player.vimeo.com/video/' + video_code

                video_width = self.get_tool_css_safe(video_width)
                video_height = self.get_tool_css_safe(video_height)

                data_name = self.get_tool_data_storage('<iframe style="width: ' + video_width + '; height: ' + video_height + ';" src="' + video_code + '" frameborder="0" allowfullscreen></iframe>', '', match_org.group(0))

                return '<' + data_name + '></' + data_name + '>'
            elif name_data == 'ruby':
                data = re.findall(macro_split_regex, match[1])

                # get option
                main_text = ''
                sub_text = ''
                color = ''
                for for_a in data:
                    data_sub = re.search(macro_split_sub_regex, for_a)
                    if data_sub:
                        data_sub = data_sub.groups()
                        if data_sub[0] == 'ruby':
                            sub_text = data_sub[1]
                        elif data_sub[0] == 'color':
                            color = data_sub[1]
                    else:
                        main_text = for_a

                main_text = self.get_tool_data_revert(main_text, do_type = 'render')
                sub_text = self.get_tool_data_revert(sub_text, do_type = 'render')

                color = self.get_tool_css_safe(color)

                # add color
                if color != '':
                    sub_text = '<span style="color:' + color + ';">' + sub_text + '</span>'

                data_name = self.get_tool_data_storage('<ruby>' + main_text + '<rp>(</rp><rt>' + sub_text + '</rt><rp>)</rp></ruby>', '', match_org.group(0))

                return '<' + data_name + '></' + data_name + '>'
            elif name_data == 'age':
                if re.search(r'^[0-9]{4}-[0-9]{2}-[0-9]{2}$', match[1]):
                    try:
                        date = datetime.datetime.strptime(match[1], '%Y-%m-%d')
                    except:
                        data_text = 'invalid date'

                    date_now = datetime.datetime.today()

                    if date > date_now:
                        data_text = 'invalid date'
                    else:
                        data_text = str((date_now - date).days // 365)
                else:
                    data_text = 'invalid date'

                data_name = self.get_tool_data_storage(data_text, '', match_org.group(0))

                return '<' + data_name + '></' + data_name + '>'
            elif name_data == 'dday':
                if re.search(r'^[0-9]{4}-[0-9]{2}-[0-9]{2}$', match[1]):
                    try:
                        date = datetime.datetime.strptime(match[1], '%Y-%m-%d')
                    except:
                        data_text = 'invalid date'

                    date_now = datetime.datetime.today()
                    
                    date_end = (date_now - date).days
                    if date_end > 0:
                        data_text = '+' + str(date_end)
                    else:
                        if date_end == 0:
                            data_text = '-' + str(date_end)
                        else:
                            data_text = str(date_end)
                else:
                    data_text = 'invalid date'

                data_name = self.get_tool_data_storage(data_text, '', match_org.group(0))

                return '<' + data_name + '></' + data_name + '>'
            elif name_data == 'pagecount':
                return '0'
            else:
                return '<macro>' + match[0] + '(' + match[1] + ')' + '</macro>'

        # double macro replace
        self.render_data = re.sub(r'\[([^[(]+)\(([^()]+)\)\]', do_render_macro_double, self.render_data)

        # single macro function
        def do_render_macro_single(match):
            match_org = match
            match = match.group(1)

            if match in ('date', 'datetime'):
                data_name = self.get_tool_data_storage(get_time(), '', match_org.group(0))

                return '<' + data_name + '></' + data_name + '>'
            elif match == 'br':
                data_name = self.get_tool_data_storage('<br>', '', match_org.group(0))

                return '<' + data_name + '></' + data_name + '>'
            elif match == 'clearfix':
                data_name = self.get_tool_data_storage('<div style="clear: both;"></div>', '', match_org.group(0))

                return '<' + data_name + '></' + data_name + '>'
            elif match == 'pagecount':
                return '0'
            else:
                return '<macro>' + match + '</macro>'

        # single macro replace
        self.render_data = re.sub(r'\[([^[\]]+)\]', do_render_macro_single, self.render_data)

        # macro safe restore
        self.render_data = re.sub(r'<macro>', '[', self.render_data)
        self.render_data = re.sub(r'<\/macro>', ']', self.render_data)

    def do_render_math(self):
        def do_render_math_sub(match):
            data = match.group(1)

            data = re.sub(r'\n', '', data)
            data = self.get_tool_data_revert(data)

            data_html = self.get_tool_js_safe(data)

            data = html.unescape(data)
            data = self.get_tool_js_safe(data)

            name_ob = 'opennamu_math_' + str(self.data_math_count)

            data_name = self.get_tool_data_storage('<span id="' + name_ob + '">', '</span>', match.group(0))

            self.render_data_js += '' + \
                'try {\n' + \
                    'katex.render("' + data + '", document.getElementById(\"' + name_ob + '\"));\n' + \
                '} catch {\n' + \
                    'if(document.getElementById(\"' + name_ob + '\")) {\n' + \
                        'document.getElementById(\"' + name_ob + '\").innerHTML = "<span style=\'color: red;\'>' + data_html + '</span>";\n' + \
                    '}\n' + \
                '}\n' + \
            ''

            self.data_math_count += 1

            return '<' + data_name + '></' + data_name + '>'

        self.render_data = re.sub(r'\[math\(((?:(?!\[math\(|\)\]).|\n)+)\)\]', do_render_math_sub, self.render_data)

    def do_render_link(self):
        link_regex = r'\[\[((?:(?!\[\[|\]\]|\||<|>).|<slash_[0-9]+>)+)(?:\|((?:(?!\[\[|\]\]|\|).)+))?\]\]'
        link_count_all = len(re.findall(link_regex, self.render_data)) * 4
        while 1:
            if not re.search(link_regex, self.render_data):
                break
            elif link_count_all < 0:
                print('Error : render heading count overflow')

                break
            else:
                # link split
                link_data = re.search(link_regex, self.render_data)
                link_data_full = link_data.group(0)
                link_data = link_data.groups()

                link_main = link_data[0]
                link_main_org = link_main

                # file link
                if re.search(r'^(파일|file|외부|out):', link_main):
                    file_width = ''
                    file_height = ''
                    file_align = ''
                    file_bgcolor = ''

                    file_split_regex = r'(?:^|&amp;) *((?:(?!&amp;).)+)'
                    file_split_sub_regex = r'(^[^=]+) *= *([^=]+)'
                    if link_data[1]:
                        data = re.findall(file_split_regex, link_data[1])
                        for for_a in data:
                            data_sub = re.search(file_split_sub_regex, for_a)
                            if data_sub:
                                data_sub = data_sub.groups()
                                if data_sub[0] == 'width':
                                    file_width = self.get_tool_px_add_check(data_sub[1])
                                elif data_sub[0] == 'height':
                                    file_height = self.get_tool_px_add_check(data_sub[1])
                                elif data_sub[0] == 'align':
                                    if data_sub[1] in ('left', 'right'):
                                        file_align = 'float:' + data_sub[1] + ';'
                                    elif data_sub[1] == 'center':
                                        file_align = 'center'
                                elif data_sub[0] == 'bgcolor':
                                    file_bgcolor = data_sub[1]

                    link_main_org = ''
                    link_sub = link_main

                    link_out_regex = r'^(외부|out):'
                    link_in_regex = r'^(파일|file):'
                    if re.search(link_out_regex, link_main):
                        link_main = re.sub(link_out_regex, '', link_main)

                        link_main = self.get_tool_data_restore(link_main, do_type = 'slash')
                        link_main = html.unescape(link_main)
                        link_main = re.sub(r'"', '&quot;', link_main)
                        
                        link_exist = ''
                    else:
                        link_main = re.sub(link_in_regex, '', link_main)

                        link_main = self.get_tool_data_restore(link_main, do_type = 'slash')
                        link_main = html.unescape(link_main)

                        self.curs.execute(db_change("select title from data where title = ?"), ['file:' + link_main])
                        db_data = self.curs.fetchall()
                        if db_data:
                            link_exist = ''
                        else:
                            link_exist = 'opennamu_not_exist_link'
                        
                        link_extension_regex = r'\.([^.]+)$'
                        link_extension = re.search(link_extension_regex, link_main)
                        if link_extension:
                            link_extension = link_extension.group(1)
                        else:
                            link_extension = 'jpg'

                        link_main = re.sub(link_extension_regex, '', link_main)
                        link_main_org = link_main

                        link_main = '/image/' + url_pas(sha224_replace(link_main)) + '.' + link_extension

                    file_width = self.get_tool_css_safe(file_width)
                    file_height = self.get_tool_css_safe(file_height)

                    file_end = '<image style="width:' + file_width + ';height:' + file_height + ';' + file_align + '" src="' + link_main + '">'
                    if file_align == 'center':
                        file_end = '<div style="text-align:center;">' + file_end + '</div>'

                    if link_exist != '':
                        data_name = self.get_tool_data_storage('<a class="' + link_exist + '" href="/upload?name=' + url_pas(link_main_org) + '">', '</a>', link_data_full)

                        self.render_data = re.sub(link_regex, lambda x : ('<' + data_name + '>' + link_sub + '</' + data_name + '>'), self.render_data, 1)
                    else:
                        data_name = self.get_tool_data_storage(file_end, '', link_data_full)

                        self.render_data = re.sub(link_regex, '<' + data_name + '></' + data_name + '>', self.render_data, 1)
                # category
                elif re.search(r'^(분류|category):', link_main):
                    link_main = re.sub(r'^(분류|category):', '', link_main)

                    if self.data_category == '':
                        self.data_category = '<div class="opennamu_category">' + self.get_tool_lang('category') + ' : '
                    else:
                        self.data_category += ' | '

                    if link_data[1]:
                        link_main += link_data[1]

                    category_blur = ''
                    if re.search(r'#blur$', link_main):
                        link_main = re.sub(r'#blur$', '', link_main)

                        category_blur = 'opennamu_category_blur'
                    
                    link_sub = link_main

                    link_main = self.get_tool_data_restore(link_main, do_type = 'slash')
                    link_main = html.unescape(link_main)

                    self.curs.execute(db_change("select title from data where title = ?"), ['category:' + link_main])
                    db_data = self.curs.fetchall()
                    if db_data:
                        link_exist = ''
                    else:
                        link_exist = 'opennamu_not_exist_link'

                    link_main = url_pas(link_main)

                    self.data_category += '<a class="' + category_blur + ' ' + link_exist + '" href="/w/category:' + link_main + '">' + link_sub + '</a>'

                    self.render_data = re.sub(link_regex, '', self.render_data, 1)
                # out link
                elif re.search(r'^(?:inter|인터):([^:]+):', link_main):
                    link_inter_name = re.search(r'^(?:inter|인터):([^:]+):', link_main)
                    link_inter_name = link_inter_name.group(1)

                    link_main = re.sub(r'^(?:inter|인터):([^:]+):', '', link_main)

                    # sharp
                    link_data_sharp_regex = r'#([^#]+)$'
                    link_data_sharp = re.search(link_data_sharp_regex, link_main)
                    if link_data_sharp:
                        link_data_sharp = link_data_sharp.group(1)
                        link_data_sharp = html.unescape(link_data_sharp)
                        link_data_sharp = '#' + url_pas(link_data_sharp)

                        link_main = re.sub(link_data_sharp_regex, '', link_main)
                    else:
                        link_data_sharp = ''

                    link_main = self.get_tool_data_restore(link_main, do_type = 'slash')
                    link_main = html.unescape(link_main)
                    link_main = url_pas(link_main)

                    self.curs.execute(db_change("select plus, plus_t from html_filter where kind = 'inter_wiki' and html = ?"), [link_inter_name])
                    db_data = self.curs.fetchall()
                    if db_data:
                        link_main = db_data[0][0] + link_main

                        # sub not exist -> sub = main
                        if link_data[1]:
                            link_sub = link_data[1]
                            link_sub_storage = ''
                        else:
                            link_sub = ''
                            link_sub_storage = link_main_org

                        if db_data[0][1] != '':
                            link_sub_storage = db_data[0][1] + link_sub_storage

                        data_name = self.get_tool_data_storage('<a class="opennamu_link_inter" href="' + link_main + link_data_sharp + '">' + link_sub_storage, '</a>', link_data_full)
                    
                        self.render_data = re.sub(link_regex, lambda x : ('<' + data_name + '>' + link_sub + '</' + data_name + '>'), self.render_data, 1)
                    else:
                        self.render_data = re.sub(link_regex, '', self.render_data, 1)
                elif re.search(r'^https?:\/\/', link_main):
                    link_main = self.get_tool_data_restore(link_main, do_type = 'slash')
                    link_main = html.unescape(link_main)
                    link_main = re.sub(r'"', '&quot;', link_main)
                    
                    # sub not exist -> sub = main
                    if link_data[1]:
                        link_sub = link_data[1]
                        link_sub_storage = ''
                    else:
                        link_sub = ''
                        link_sub_storage = link_main_org

                    data_name = self.get_tool_data_storage('<a class="opennamu_link_out" href="' + link_main + '">' + link_sub_storage, '</a>', link_data_full)

                    self.render_data = re.sub(link_regex, lambda x : ('<' + data_name + '>' + link_sub + '</' + data_name + '>'), self.render_data, 1)
                # in link
                else:
                    # sharp
                    link_data_sharp_regex = r'#([^#]+)$'
                    link_data_sharp = re.search(link_data_sharp_regex, link_main)
                    if link_data_sharp:
                        link_data_sharp = link_data_sharp.group(1)
                        link_data_sharp = html.unescape(link_data_sharp)
                        link_data_sharp = '#' + url_pas(link_data_sharp)

                        link_main = re.sub(link_data_sharp_regex, '', link_main)
                    else:
                        link_data_sharp = ''

                    # under page & fix url
                    if link_main == '../':
                        link_main = self.doc_name
                        link_main = re.sub(r'(\/[^/]+)$', '', link_main)
                    elif re.search(r'^\/', link_main):
                        link_main = re.sub(r'^\/', lambda x : (self.doc_name + '/'), link_main)
                    elif re.search(r'^:(분류|category):', link_main):
                        link_main = re.sub(r'^:(분류|category):', 'category:', link_main)
                    elif re.search(r'^:(파일|file):', link_main):
                        link_main = re.sub(r'^:(파일|file):', 'file:', link_main)
                    elif re.search(r'^사용자:', link_main):
                        link_main = re.sub(r'^사용자:', 'user:', link_main)

                    # main link fix
                    link_main = self.get_tool_data_restore(link_main, do_type = 'slash')
                    link_main = html.unescape(link_main)

                    link_exist = ''
                    if link_main != '':
                        self.curs.execute(db_change("select title from data where title = ?"), [link_main])
                        db_data = self.curs.fetchall()
                        if not db_data:
                            link_exist = 'opennamu_not_exist_link'

                    link_same = ''
                    if link_main == self.doc_name and self.doc_include == '':
                        link_same = 'opennamu_same_link'

                    link_main = url_pas(link_main)

                    if link_main != '':
                        link_main = '/w/' + link_main

                    # sub not exist -> sub = main
                    if link_data[1]:
                        link_sub = link_data[1]
                        link_sub_storage = ''
                    else:
                        link_sub = ''
                        link_sub_storage = link_main_org

                    data_name = self.get_tool_data_storage('<a class="' + link_exist + ' ' + link_same + '" href="' + link_main + link_data_sharp + '">' + link_sub_storage, '</a>', link_data_full)

                    self.render_data = re.sub(link_regex, lambda x : ('<' + data_name + '>' + link_sub + '</' + data_name + '>'), self.render_data, 1)

            link_count_all -= 1

    def do_render_slash(self):
        # slash text -> <slash_n>
        
        def do_render_slash_sub(match):
            if match.group(1) == '<':
                return '<'
            else:
                data_name = self.get_tool_data_storage(match.group(1), do_type = 'slash')

                return '<' + data_name + '>'

        self.render_data = re.sub(r'\\(&lt;|&gt;|&#x27;|&quot;|&amp;|.)', do_render_slash_sub, self.render_data)

    def do_render_include_default(self):
        def do_render_include_default_sub(match):
            match_org = match.group(0)
            match = match.groups()

            if len(match) < 3:
                match = list(match) + ['']

            if match[2] == '\\':
                return match_org
            else:
                slash_add = ''
                if match[0]:
                    if len(match[0]) % 2 == 1:
                        slash_add = '\\' * (len(match[0]) - 1)
                    else:
                        slash_add = match[0]

                return slash_add + match[2]

        self.render_data = re.sub(r'(\\+)?@([^@= ]+)=((?:\\@|[^@])+)@', do_render_include_default_sub, self.render_data)
        self.render_data = re.sub(r'(\\+)?@([^@= ]+)@', do_render_include_default_sub, self.render_data)

    def do_render_include(self):
        def do_render_include_default_sub(match):
            match_org = match.group(0)
            match = match.groups()

            if len(match) < 3:
                match = list(match) + ['']

            if match[2] == '\\':
                return match_org
            else:
                slash_add = ''
                if match[0]:
                    if len(match[0]) % 2 == 1:
                        slash_add = '\\' * (len(match[0]) - 1)
                    else:
                        slash_add = match[0]

                if match[1] in include_change_list:
                    return slash_add + include_change_list[match[1]]
                else:
                    return slash_add + match[2]

        include_num = 0
        include_regex = r'\[include\(((?:(?!\[include\(|\)\]|<\/div>).)+)\)\]'
        include_count_max = len(re.findall(include_regex, self.render_data)) * 2
        include_change_list = {}
        while 1:
            include_num += 1
            include_change_list = {}

            match = re.search(include_regex, self.render_data)
            if include_count_max < 0:
                break
            elif not match:
                break
            else:
                match_org = match.group(0)
                match = match.groups()

                macro_split_regex = r'(?:^|,) *([^,]+)'
                macro_split_sub_regex = r'(^[^=]+) *= *([^=]*)'

                include_name = ''

                data = re.findall(macro_split_regex, match[0])
                for for_a in data:
                    data_sub = re.search(macro_split_sub_regex, for_a)
                    if data_sub:
                        data_sub = data_sub.groups()
                        
                        data_sub_name = data_sub[0]
                        data_sub_data = self.get_tool_data_restore(data_sub[1], do_type = 'slash')
                        
                        data_sub_data = re.sub(r'^분류:', ':분류:', data_sub_data)
                        data_sub_data = re.sub(r'^파일:', ':파일:', data_sub_data)

                        include_change_list[data_sub_name] = data_sub_data
                    else:
                        include_name = for_a

                include_name_org = include_name
                
                include_name = self.get_tool_data_restore(include_name, do_type = 'slash')
                include_name = html.unescape(include_name)

                # load include db data
                self.curs.execute(db_change("select data from data where title = ?"), [include_name])
                db_data = self.curs.fetchall()
                if db_data:
                    # include link func
                    include_link = ''
                    if flask.request.cookies.get('main_css_include_link', '') == '1':
                        include_link = '<div><a href="/w/' + url_pas(include_name) + '">(' + include_name_org + ')</a></div>'

                    include_data = db_data[0][0]

                    # parameter replace
                    include_data = re.sub(r'(\\+)?@([^@= ]+)=((?:\\@|[^@])+)@', do_render_include_default_sub, include_data)
                    include_data = re.sub(r'(\\+)?@([^@= ]+)@', do_render_include_default_sub, include_data)

                    # remove include
                    include_data = re.sub(include_regex, '', include_data)

                    self.data_include += [[self.doc_include + 'opennamu_include_' + str(include_num), include_name, include_data]]

                    data_name = self.get_tool_data_storage('' + \
                        include_link + \
                        '<div id="' + self.doc_include + 'opennamu_include_' + str(include_num) + '"></div>' + \
                    '', '', match_org)
                else:
                    include_link = '<div><a class="opennamu_not_exist_link" href="/w/' + url_pas(include_name) + '">(' + include_name_org + ')</a></div>'

                    data_name = self.get_tool_data_storage(include_link, '', match_org)

                self.render_data = re.sub(include_regex, '<' + data_name + '></' + data_name + '>', self.render_data, 1)

            include_count_max -= 1

    def do_render_list(self):
        pass

    def do_redner_footnote(self):
        footnote_num = 0
        footnote_regex = r'(?:\[\*((?:(?!\[\*|\]| ).)+)?(?: ((?:(?!\[\*|\]).)+))?\]|\[(각주|footnote)\])'
        footnote_count_all = len(re.findall(footnote_regex, self.render_data)) * 4
        while 1:
            footnote_num += 1

            footnote_data = re.search(footnote_regex, self.render_data)
            if footnote_count_all < 0:
                break
            elif not footnote_data:
                break
            else:
                footnote_data_org = footnote_data.group(0)
                footnote_data = footnote_data.groups()
                if footnote_data[2]:
                    self.render_data = re.sub(footnote_regex, lambda x : self.get_tool_footnote_make(), self.render_data, 1)
                else:
                    if not footnote_data[0]:
                        footnote_name = str(footnote_num)
                        footnote_name_add = ''
                    else:
                        footnote_name = footnote_data[0]
                        footnote_name_add = ' (' + str(footnote_num) + ')'

                    if not footnote_data[1]:
                        footnote_text_data = ''
                    else:
                        footnote_text_data = footnote_data[1]

                    if footnote_name in self.data_footnote:
                        self.data_footnote[footnote_name]['list'] += [str(footnote_num)]
                        footnote_first = self.data_footnote[footnote_name]['list'][0]
                    
                        data_name = self.get_tool_data_storage('<sup><a id="' + self.doc_include + 'rfn_' + str(footnote_num) + '" href="#' + self.doc_include + 'fn_' + footnote_first + '">(' + footnote_name + ' (' + str(footnote_num) + ')' + ')</a></sup>', '', footnote_data_org)

                        self.render_data = re.sub(footnote_regex, '<' + data_name + '></' + data_name + '>', self.render_data, 1)
                    else:
                        self.data_footnote[footnote_name] = {}
                        self.data_footnote[footnote_name]['list'] = [str(footnote_num)]
                        self.data_footnote[footnote_name]['data'] = footnote_text_data

                        data_name = self.get_tool_data_storage('<sup><a id="' + self.doc_include + 'rfn_' + str(footnote_num) + '" href="#' + self.doc_include + 'fn_' + str(footnote_num) + '">(' + footnote_name + footnote_name_add + ')</a></sup>', '', footnote_data_org)

                        self.render_data = re.sub(footnote_regex, '<' + data_name + '></' + data_name + '>', self.render_data, 1)

            footnote_count_all -= 1

        self.render_data += '<footnote_category>'
        self.render_data += self.get_tool_footnote_make()

    def do_render_redirect(self):
        match = re.search(r'^<back_br>\n#(?:redirect|넘겨주기) ([^\n]+)', self.render_data)
        if match:
            link_data_full = match.group(0)
            link_main = match.group(1)

            # sharp
            link_data_sharp_regex = r'#([^#]+)$'
            link_data_sharp = re.search(link_data_sharp_regex, link_main)
            if link_data_sharp:
                link_data_sharp = link_data_sharp.group(1)
                link_data_sharp = html.unescape(link_data_sharp)
                link_data_sharp = '#' + url_pas(link_data_sharp)

                link_main = re.sub(link_data_sharp_regex, '', link_main)
            else:
                link_data_sharp = ''

            # under page & fix url
            if link_main == '../':
                link_main = self.doc_name
                link_main = re.sub(r'(\/[^/]+)$', '', link_main)
            elif re.search(r'^\/', link_main):
                link_main = re.sub(r'^\/', lambda x : (self.doc_name + '/'), link_main)
            elif re.search(r'^분류:', link_main):
                link_main = re.sub(r'^분류:', 'category:', link_main)
            elif re.search(r'^사용자:', link_main):
                link_main = re.sub(r'^사용자:', 'user:', link_main)

            link_main = self.get_tool_data_restore(link_main, do_type = 'slash')
            link_main = html.unescape(link_main)
            link_main = url_pas(link_main)

            if link_main != '':
                link_main = '/w_from/' + link_main

            if 'doc_from' in self.doc_set:
                data_name = self.get_tool_data_storage('<a href="' + link_main + link_data_sharp + '">(GO)</a>', '', link_data_full)
            else:
                data_name = self.get_tool_data_storage('<meta http-equiv="refresh" content="0; url=' + link_main + link_data_sharp + '">', '', link_data_full)
                
            self.render_data = '<' + data_name + '></' + data_name + '>'

    def do_render_table(self):
        self.render_data = re.sub(r'\n +\|\|', '\n||', self.render_data)

        # get_tool_dark_mode_split
        # get_tool_px_add_check
        # get_tool_css_safe
        # todo : after text render text not use to make table
        def do_render_table_parameter(cell_count, parameter, data, option = {}):
            table_parameter_all = { "div" : "", "table" : "", "tr" : "", "td" : "", "col" : "", "colspan" : "", "rowspan" : "", "data" : "" }
            
            table_align_auto = 1
            table_colspan_auto = 1

            # todo : useless parameter return
            table_parameter_regex = r'&lt;((?:(?!&lt;|&gt;).)+)&gt;'
            for table_parameter in re.findall(table_parameter_regex, parameter):
                table_parameter_split = table_parameter.split('=')
                if len(table_parameter_split) == 2:
                    table_parameter_name = table_parameter_split[0].replace(' ', '')
                    table_parameter_data = self.get_tool_css_safe(table_parameter_split[1])
                    if table_parameter_name == 'tablebgcolor':
                        table_parameter_all['table'] += 'background:' + self.get_tool_dark_mode_split(table_parameter_data) + ';'
                    elif table_parameter_name == 'tablewidth':
                        table_parameter_all['table'] += 'width:' + self.get_tool_px_add_check(table_parameter_data) + ';'
                    elif table_parameter_name == 'tableheight':
                        table_parameter_all['table'] += 'height:' + self.get_tool_px_add_check(table_parameter_data) + ';'
                    elif table_parameter_name == 'tablealign':
                        if table_parameter_data == 'right':
                            table_parameter_all['div'] += 'float:right;'
                        elif table_parameter_data == 'center':
                            table_parameter_all['div'] += 'margin:auto;'
                            table_parameter_all['table'] += 'margin:auto;'
                    elif table_parameter_name == 'tabletextalign':
                        table_parameter_all['table'] += 'text-align:' + table_parameter_data + ';'
                    elif table_parameter_name == 'tablecolor':
                        table_parameter_all['table'] += 'color:' + self.get_tool_dark_mode_split(table_parameter_data) + ';'
                    elif table_parameter_name == 'tablebordercolor':
                        table_parameter_all['table'] += 'border:2px solid ' + self.get_tool_dark_mode_split(table_parameter_data) + ';'
                    elif table_parameter_name == 'rowbgcolor':
                        table_parameter_all['tr'] += 'background:' + self.get_tool_dark_mode_split(table_parameter_data) + ';'
                    elif table_parameter_name == 'rowtextalign':
                        table_parameter_all['tr'] += 'text-align:' + table_parameter_data + ';'
                    elif table_parameter_name == 'rowcolor':
                        table_parameter_all['tr'] += 'color:' + self.get_tool_dark_mode_split(table_parameter_data) + ';'
                    elif table_parameter_name == 'colcolor':
                        table_parameter_all['col'] += 'color:' + self.get_tool_dark_mode_split(table_parameter_data) + ';'
                    elif table_parameter_name == 'colbgcolor':
                        table_parameter_all['col'] += 'background:' + self.get_tool_dark_mode_split(table_parameter_data) + ';'
                    elif table_parameter_name == 'bgcolor':
                        table_parameter_all['td'] += 'background:' + self.get_tool_dark_mode_split(table_parameter_data) + ';'
                    elif table_parameter_name == 'color':
                        table_parameter_all['td'] += 'color:' + self.get_tool_dark_mode_split(table_parameter_data) + ';'
                    elif table_parameter_name == 'width':
                        table_parameter_all['td'] += 'width:' + self.get_tool_px_add_check(table_parameter_data) + ';'
                    elif table_parameter_name == 'height':
                        table_parameter_all['td'] += 'height:' + self.get_tool_px_add_check(table_parameter_data) + ';'
                elif len(table_parameter_split) == 1:
                    if re.search(r'^-[0-9]+$', table_parameter):
                        table_colspan_auto = 0
                        table_parameter_all['colspan'] = table_parameter.replace('-', '')
                    elif re.search(r'^(\^|v)?\|[0-9]+$', table_parameter):
                        if table_parameter[0] == '^':
                            table_parameter_all['td'] += 'vertical-align: top;'
                        elif table_parameter[0] == 'v':
                            table_parameter_all['td'] += 'vertical-align: bottom;'

                        table_parameter_all['rowspan'] = re.sub(r'^|v|\|', '', table_parameter)
                    elif table_parameter in ('(', ':', ')'):
                        table_align_auto = 0
                        if table_parameter == '(':
                            table_parameter_all['td'] += 'text-align: left;'
                        elif table_parameter == ':':
                            table_parameter_all['td'] += 'text-align: center;'
                        elif table_parameter == ':':
                            table_parameter_all['td'] += 'text-align: right;'
                    else:
                        table_parameter_data = self.get_tool_css_safe(table_parameter)
                        table_parameter_all['td'] += 'background:' + self.get_tool_dark_mode_split(table_parameter_data) + ';'
            
            if table_align_auto == 1:
                if re.search(r'^ ', data):
                    data = re.sub(r'^ ', '', data)
                    if re.search(r' $', data):
                        table_parameter_all['td'] += 'text-align: center;'

                        data = re.sub(r' $', '', data)
                    else:
                        table_parameter_all['td'] += 'text-align: right;'
                else:
                    if re.search(r' $', data):
                        data = re.sub(r' $', '', data)

            if table_colspan_auto == 1:
                table_parameter_all['colspan'] = str(len(cell_count) // 2)

            table_parameter_all['data'] = data

            return table_parameter_all

        table_regex = re.compile('\n((?:(?:(?:(?:\|\|)+)|(?:\|[^|]+\|(?:\|\|)*))\n?(?:(?:(?!\|\|).)+))(?:(?:\|\||\|\|\n|(?:\|\|)+(?!\n)(?:(?:(?!\|\|).)+)\n*)*)\|\|)\n', re.DOTALL)
        table_sub_regex = r'(\n?)((?:\|\|)+)((?:&lt;(?:(?:(?!&lt;|&gt;).)+)&gt;)*)((?:\n*(?:(?:(?:(?!\|\|).)+)\n*)+)|(?:(?:(?!\|\|).)*))'
        table_caption_regex = r'^\|([^|]+)\|'
        table_count_all = len(re.findall(table_regex, self.render_data)) * 2
        while 1:
            table_data = re.search(table_regex, self.render_data)
            if table_count_all < 0:
                print('Error : render table count overflow')

                break
            elif not table_data:
                break
            else:
                table_data_org = table_data.group(0)
                table_data = table_data.group(1)
                
                table_caption = re.search(table_caption_regex, table_data)
                if table_caption:
                    table_caption = table_caption.group(1)
                    table_caption = '<caption>' + table_caption + '</caption>'

                    table_data = re.sub(table_caption_regex, '||', table_data)
                else:
                    table_caption = ''

                table_parameter = { "div" : "", "table" : "", "col" : {} }
                table_data_end = ''
                table_col_num = 0
                table_tr_change = 0
                for table_sub in re.findall(table_sub_regex, table_data):
                    table_data_in = table_sub[3]
                    table_data_in = re.sub(r'^\n+', '', table_data_in)

                    table_sub_parameter = do_render_table_parameter(table_sub[1], table_sub[2], table_data_in)

                    if table_data_end == '':
                        table_data_end += '<tr style="' + table_sub_parameter['tr'] + '">'

                    if table_sub[0] != '' and table_tr_change == 1:
                        table_col_num = 0
                        table_data_end += '</tr><tr style="' + table_sub_parameter['tr'] + '">'

                    if not table_col_num in table_parameter['col']:
                        table_parameter['col'][table_col_num] = ''

                    table_parameter['div'] += table_sub_parameter['div']
                    table_parameter['table'] += table_sub_parameter['table']
                    table_parameter['col'][table_col_num] += table_sub_parameter['col']

                    if table_sub[2] == '' and table_sub[3] == '':
                        table_tr_change = 1
                    else:
                        table_tr_change = 0
                    
                        table_data_end += '<td colspan="' + table_sub_parameter['colspan'] + '" rowspan="' + table_sub_parameter['rowspan'] + '" style="' + table_sub_parameter['td'] + table_parameter['col'][table_col_num] + '">' + table_sub_parameter['data'] + '</td>'
                    
                    table_col_num += 1

                table_data_end += '</tr>'
                table_data_end = '<table style="' + table_parameter['table'] + '">' + table_caption + table_data_end + '</table>'
                if table_parameter['div'] != '':
                    table_data_end = '<div style="' + table_parameter['div'] + '">' + table_data_end + '</div>'

                self.render_data = re.sub(table_regex, lambda x : ('\n<front_br>' + table_data_end + '<back_br>\n'), self.render_data, 1)

            table_count_all -= 1
    
    def do_render_middle(self):
        middle_regex = r'{{{([^{](?:(?!{{{|}}}).|\n)*)?(?:}|<(\/?(?:slash)_(?:[0-9]+))>)}}'
        wiki_count = 0
        syntax_count = 0
        folding_count = 0
        middle_count_all = len(re.findall(middle_regex, self.render_data)) * 10
        while 1:
            middle_data = re.search(middle_regex, self.render_data)
            if middle_count_all < 0:
                break
            elif not middle_data:
                break
            else:
                middle_data_org = middle_data.group(0)
                middle_slash = middle_data.group(2)
                if middle_slash:
                    if self.data_temp_storage[middle_slash] != '}':
                        middle_data_org = re.sub(r'<(\/?(?:slash)_(?:[0-9]+))>', '<temp_' + middle_slash + '>', middle_data_org)
                        self.render_data = re.sub(middle_regex, lambda x : middle_data_org, self.render_data, 1)
                        continue

                middle_data = middle_data.group(1)
                if not middle_data:
                    middle_data = ''

                middle_name = re.search(r'^([^ \n]+)', middle_data)
                middle_data_pass = ''
                if middle_name:
                    middle_name = middle_name.group(1)
                    if middle_name == '#!wiki':
                        if middle_slash:
                            middle_data_org = re.sub(r'<(\/?(?:slash)_(?:[0-9]+))>', '<temp_' + middle_slash + '>', middle_data_org)
                            self.render_data = re.sub(middle_regex, lambda x : middle_data_org, self.render_data, 1)
                            continue

                        wiki_regex = r'^#!wiki(?:(?: style=(&quot;(?:(?:(?!&quot;).)*)&quot;|&#x27;(?:(?:(?!&#x27;).)*)&#x27;))| [^\n]*)?\n'
                        wiki_data_style = re.search(wiki_regex, middle_data)
                        wiki_data = re.sub(wiki_regex, '', middle_data)
                        if wiki_data_style:
                            wiki_data_style = wiki_data_style.group(1)
                            if wiki_data_style:
                                wiki_data_style = wiki_data_style.replace('&#x27;', '\'')
                                wiki_data_style = wiki_data_style.replace('&quot;', '"')
                                wiki_data_style = 'style=' + wiki_data_style
                            else:
                                wiki_data_style = ''
                        else:
                            wiki_data_style = ''

                        wiki_data = self.get_tool_data_revert(wiki_data)
                        wiki_data = html.unescape(wiki_data)
                        wiki_data = re.sub('\n$', '', wiki_data)

                        self.data_include += [[self.doc_include + 'opennamu_wiki_' + str(wiki_count), self.doc_name, wiki_data, wiki_data_style]]

                        data_name = self.get_tool_data_storage('<div id="' + self.doc_include + 'opennamu_wiki_' + str(wiki_count) + '"></div>', '', middle_data_org)
                        wiki_count += 1
                    elif middle_name == '#!html':
                        if middle_slash:
                            middle_data_org = re.sub(r'<(\/?(?:slash)_(?:[0-9]+))>', '<temp_' + middle_slash + '>', middle_data_org)
                            self.render_data = re.sub(middle_regex, lambda x : middle_data_org, self.render_data, 1)
                            continue

                        data_name = self.get_tool_data_storage('', '', middle_data_org)
                    elif middle_name == '#!folding':
                        if middle_slash:
                            middle_data_org = re.sub(r'<(\/?(?:slash)_(?:[0-9]+))>', '<temp_' + middle_slash + '>', middle_data_org)
                            self.render_data = re.sub(middle_regex, lambda x : middle_data_org, self.render_data, 1)
                            continue

                        wiki_regex = r'^#!folding(?: ([^\n]*))?\n'
                        wiki_data_folding = re.search(wiki_regex, middle_data)
                        wiki_data = re.sub(wiki_regex, '', middle_data)
                        if wiki_data_folding:
                            wiki_data_folding = wiki_data_folding.group(1)
                            if not wiki_data_folding:
                                wiki_data_folding = 'test'
                        else:
                            wiki_data_folding = 'test'

                        wiki_data = self.get_tool_data_revert(wiki_data)
                        wiki_data = html.unescape(wiki_data)
                        wiki_data = re.sub('\n$', '', wiki_data)

                        self.data_include += [[self.doc_include + 'opennamu_folding_' + str(folding_count), self.doc_name, wiki_data]]

                        middle_data_pass = wiki_data_folding
                        data_name = self.get_tool_data_storage(
                            '<details><summary>',
                            '</summary><div id="' + self.doc_include + 'opennamu_folding_' + str(folding_count) + '"></div></details>', 
                            middle_data_org
                        )
                        folding_count += 1
                    elif middle_name == '#!syntax':
                        if middle_slash:
                            middle_data_org = re.sub(r'<(\/?(?:slash)_(?:[0-9]+))>', '<temp_' + middle_slash + '>', middle_data_org)
                            self.render_data = re.sub(middle_regex, lambda x : middle_data_org, self.render_data, 1)
                            continue

                        wiki_regex = r'^#!syntax(?: ([^\n]*))?\n'
                        wiki_data_syntax = re.search(wiki_regex, middle_data)
                        wiki_data = re.sub(wiki_regex, '', middle_data)
                        if wiki_data_syntax:
                            wiki_data_syntax = wiki_data_syntax.group(1)
                            if not wiki_data_syntax:
                                wiki_data_syntax = 'python'
                        else:
                            wiki_data_syntax = 'python'

                        if syntax_count == 0:
                            self.render_data_js += 'hljs.highlightAll();\n'

                        data_name = self.get_tool_data_storage('<pre id="syntax"><code class="' + wiki_data_syntax + '">' + wiki_data, '</code></pre>', middle_data_org)
                        syntax_count += 1
                    elif middle_name in ('+5', '+4', '+3', '+2', '+1'):
                        if middle_slash:
                            middle_data_org = re.sub(r'<(\/?(?:slash)_(?:[0-9]+))>', '<temp_' + middle_slash + '>', middle_data_org)
                            self.render_data = re.sub(middle_regex, lambda x : middle_data_org, self.render_data, 1)
                            continue

                        wiki_data = re.sub(r'^\+[1-5] ', '', middle_data)
                        if middle_name == '+5':
                            wiki_size = '200'
                        elif middle_name == '+4':
                            wiki_size = '180'
                        elif middle_name == '+3':
                            wiki_size = '160'
                        elif middle_name == '+2':
                            wiki_size = '140'
                        else:
                            wiki_size = '120'

                        middle_data_pass = wiki_data
                        data_name = self.get_tool_data_storage('<span style="font-size:' + wiki_size + '%">', '</span>', middle_data_org)
                    elif middle_name in ('-5', '-4', '-3', '-2', '-1'):
                        if middle_slash:
                            middle_data_org = re.sub(r'<(\/?(?:slash)_(?:[0-9]+))>', '<temp_' + middle_slash + '>', middle_data_org)
                            self.render_data = re.sub(middle_regex, lambda x : middle_data_org, self.render_data, 1)
                            continue

                        wiki_data = re.sub(r'^\-[1-5] ', '', middle_data)
                        if middle_name == '-5':
                            wiki_size = '50'
                        elif middle_name == '-4':
                            wiki_size = '60'
                        elif middle_name == '-3':
                            wiki_size = '70'
                        elif middle_name == '-2':
                            wiki_size = '80'
                        else:
                            wiki_size = '90'

                        middle_data_pass = wiki_data
                        data_name = self.get_tool_data_storage('<span style="font-size:' + wiki_size + '%">', '</span>', middle_data_org)
                    elif re.search(r'^@(?:((?:[0-9a-f-A-F]{3}){1,2})|(\w+))', middle_name):
                        if middle_slash:
                            middle_data_org = re.sub(r'<(\/?(?:slash)_(?:[0-9]+))>', '<temp_' + middle_slash + '>', middle_data_org)
                            self.render_data = re.sub(middle_regex, lambda x : middle_data_org, self.render_data, 1)
                            continue

                        wiki_color = re.search(r'^@(?:((?:[0-9a-f-A-F]{3}){1,2})|(\w+))(?:,@(?:((?:[0-9a-f-A-F]{3}){1,2})|(\w+)))?', middle_name)
                        if wiki_color:
                            wiki_color = wiki_color.groups()
                            if wiki_color[0]:
                                wiki_color = '#' + wiki_color[0]
                            else:
                                wiki_color = wiki_color[1]

                            if wiki_color[2]:
                                wiki_color += ',#' + wiki_color[0]
                            elif wiki_color[3]:
                                wiki_color = ',' + wiki_color[1]
                        else:
                            wiki_color = 'red'

                        wiki_color = self.get_tool_css_safe(wiki_color)
                        wiki_color = self.get_tool_dark_mode_split(wiki_color)

                        wiki_data = re.sub(r'^@(?:((?:[0-9a-f-A-F]{3}){1,2})|(\w+))(?:,@(?:((?:[0-9a-f-A-F]{3}){1,2})|(\w+)))? ?', '', middle_data)

                        middle_data_pass = wiki_data
                        data_name = self.get_tool_data_storage('<span style="background-color:' + wiki_color + '">', '</span>', middle_data_org)
                    elif re.search(r'^#(?:((?:[0-9a-f-A-F]{3}){1,2})|(\w+))', middle_name):
                        if middle_slash:
                            middle_data_org = re.sub(r'<(\/?(?:slash)_(?:[0-9]+))>', '<temp_' + middle_slash + '>', middle_data_org)
                            self.render_data = re.sub(middle_regex, lambda x : middle_data_org, self.render_data, 1)
                            continue

                        wiki_color = re.search(r'^#(?:((?:[0-9a-f-A-F]{3}){1,2})|(\w+))(?:,#(?:((?:[0-9a-f-A-F]{3}){1,2})|(\w+)))?', middle_name)
                        if wiki_color:
                            wiki_color = wiki_color.groups()
                            if wiki_color[0]:
                                wiki_color = '#' + wiki_color[0]
                            else:
                                wiki_color = wiki_color[1]

                            if wiki_color[2]:
                                wiki_color += ',#' + wiki_color[0]
                            elif wiki_color[3]:
                                wiki_color = ',' + wiki_color[1]
                        else:
                            wiki_color = 'red'

                        wiki_color = self.get_tool_css_safe(wiki_color)
                        wiki_color = self.get_tool_dark_mode_split(wiki_color)

                        wiki_data = re.sub(r'^#(?:((?:[0-9a-f-A-F]{3}){1,2})|(\w+))(?:,#(?:((?:[0-9a-f-A-F]{3}){1,2})|(\w+)))? ?', '', middle_data)

                        middle_data_pass = wiki_data
                        data_name = self.get_tool_data_storage('<span style="color:' + wiki_color + '">', '</span>', middle_data_org)
                    else:
                        if middle_slash:
                            middle_data += '\\'

                        data_revert = self.get_tool_data_revert(middle_data)
                        data_revert = re.sub('^\n', '', data_revert)
                        data_revert = re.sub('\n$', '', data_revert)

                        data_name = self.get_tool_data_storage(data_revert, '', middle_data_org)
                else:
                    if middle_slash:
                        middle_data += '\\'

                    data_revert = self.get_tool_data_revert(middle_data)
                    data_revert = re.sub('^\n', '', data_revert)
                    data_revert = re.sub('\n$', '', data_revert)

                    data_name = self.get_tool_data_storage(data_revert, '', middle_data_org)

                self.render_data = re.sub(middle_regex, lambda x : ('<' + data_name + '>' + middle_data_pass + '</' + data_name + '>'), self.render_data, 1)

            middle_count_all -= 1

        self.render_data = re.sub(r'<temp_(?P<in>(?:slash)_(?:[0-9]+))>', '<\g<in>>', self.render_data)

    def do_render_last(self):
        # add category
        if self.doc_include == '':
            if self.data_category != '':
                data_name = self.get_tool_data_storage(self.data_category, '</div>', '')

                if flask.request.cookies.get('main_css_category_set', '0') == '0':
                    if re.search(r'<footnote_category>', self.render_data):
                        self.render_data = re.sub(r'<footnote_category>', '<hr><' + data_name + '></' + data_name + '>', self.render_data, 1)
                    else:
                        self.render_data += '<hr><' + data_name + '></' + data_name + '>'
                else:
                    self.render_data = re.sub(r'<footnote_category>', '', self.render_data, 1)
                    self.render_data = '<' + data_name + '></' + data_name + '><hr>' + self.render_data
            else:
                self.render_data = re.sub(r'<footnote_category>', '', self.render_data, 1)
        else:
            self.render_data = re.sub(r'<footnote_category>', '', self.render_data, 1)

        # remove front_br and back_br
        self.render_data = re.sub(r'\n?<front_br>', '', self.render_data)
        self.render_data = re.sub(r'<back_br>\n?', '', self.render_data)
        
        # \n to <br>
        self.render_data = re.sub(r'\n', '<br>', self.render_data)

        # <render_n> restore
        self.render_data = self.get_tool_data_restore(self.render_data)

        self.render_data = '<div class="opennamu_render_complete">' + self.render_data + '</div>'

    def __call__(self):
        self.do_render_include_default()
        self.do_render_slash()
        self.do_render_redirect()
        self.do_render_include()
        self.do_render_middle()
        self.do_render_math()
        # self.do_render_list()
        self.do_render_table()
        self.do_render_macro()
        self.do_render_link()
        self.do_redner_footnote()
        self.do_render_text()
        self.do_render_heading()
        self.do_render_last()

        # print(self.data_temp_storage)
        # print(self.render_data)

        return [
            self.render_data, # html
            self.render_data_js, # js
            {
                'backlink' : self.data_backlink, # backlink
                'include' : self.data_include # include data
            } # other
        ]