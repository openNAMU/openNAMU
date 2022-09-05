from .func_tool import *

class class_do_render_namumark:
    def __init__(self, curs, doc_name, doc_data, doc_include, lang_data):
        self.curs = curs
        
        self.doc_data = doc_data
        self.doc_name = doc_name
        self.doc_include = doc_include
        
        self.lang_data = lang_data

        self.data_temp_storage = {}
        self.data_temp_storage_count = 0
        self.data_backlink = []
        
        self.data_toc = ''
        self.data_footnote = ''
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

    def get_tool_temp_storage(self, data_A = '', data_B = ''):
        self.data_temp_storage_count += 1
        data_name = 'render_' + str(self.data_temp_storage_count)

        self.data_temp_storage[data_name] = data_A
        self.data_temp_storage['/' + data_name] = data_B

        return data_name

    def get_tool_data_restore(self, data):
        storage_count = self.data_temp_storage_count * 3
        storage_regex = r'<(\/?render_(?:[0-9]+))>'

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

    def do_render_text(self):
        # <b> function
        def do_render_text_bold(match):
            data = match.group(1)
            data_name = self.get_tool_temp_storage('<b>', '</b>')
            
            return '<' + data_name + '>' + data + '</' + data_name + '>'

        # <b>
        self.render_data = re.sub(r"&#x27;&#x27;&#x27;((?:(?!&#x27;&#x27;&#x27;).)+)&#x27;&#x27;&#x27;", do_render_text_bold, self.render_data)

        # <i> function
        def do_render_text_italic(match):
            data = match.group(1)
            data_name = self.get_tool_temp_storage('<i>', '</i>')
            
            return '<' + data_name + '>' + data + '</' + data_name + '>'

        # <i>
        self.render_data = re.sub(r"&#x27;&#x27;((?:(?!&#x27;&#x27;).)+)&#x27;&#x27;", do_render_text_italic, self.render_data)

        # <u> function
        def do_render_text_under(match):
            data = match.group(1)
            data_name = self.get_tool_temp_storage('<u>', '</u>')
            
            return '<' + data_name + '>' + data + '</' + data_name + '>'

        # <u>
        self.render_data = re.sub(r"__((?:(?!__).)+)__", do_render_text_under, self.render_data)
        
        # <sup> function
        def do_render_text_sup(match):
            data = match.group(1)
            data_name = self.get_tool_temp_storage('<sup>', '</sup>')
            
            return '<' + data_name + '>' + data + '</' + data_name + '>'

        # <sup>
        self.render_data = re.sub(r"\^\^\^((?:(?!\^\^\^).)+)\^\^\^", do_render_text_sup, self.render_data)
        # <sup> 2
        self.render_data = re.sub(r"\^\^((?:(?!\^\^).)+)\^\^", do_render_text_sup, self.render_data)

        # <sub> function
        def do_render_text_sub(match):
            data = match.group(1)
            data_name = self.get_tool_temp_storage('<sub>', '</sub>')
            
            return '<' + data_name + '>' + data + '</' + data_name + '>'
        
        # <sub>
        self.render_data = re.sub(r",,,((?:(?!,,,).)+),,,", do_render_text_sub, self.render_data)
        # <sub> 2
        self.render_data = re.sub(r",,((?:(?!,,).)+),,", do_render_text_sub, self.render_data)

        # <sub> function
        def do_render_text_strike(match):
            data = match.group(1)
            data_name = self.get_tool_temp_storage('<s>', '</s>')
            
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
        while 1:
            if not re.search(heading_regex, self.render_data):
                break
            elif heading_count_all < 0:
                print('Error : render heading count overflow')

                break
            else:
                heading_data = re.search(heading_regex, self.render_data)
                heading_data = heading_data.groups()

                heading_data_last_regex = r' ?(#?={1,6})$'
                heading_data_last = re.search(heading_data_last_regex, heading_data[3])
                heading_data_last = heading_data_last.group(1)
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
                    
                    heading_html_name = self.get_tool_temp_storage(
                        '<h' + heading_level_str + '>',
                        '</h' + heading_level_str + '>'
                    )
                    heading_data_complete = '' + \
                        '\n<front_br>' + \
                        '<' + heading_html_name + '>' + \
                            '<heading_stack>' + \
                                heading_stack_str + \
                            '</heading_stack>' + \
                            ' ' + heading_data_text + \
                        '</' + heading_html_name + '>' + \
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
        self.render_data = re.sub(r'\[(목차|toc|tableofcontents)\]', toc_data, self.render_data)

    def do_render_macro(self):
        # double macro function
        def do_render_macro_double(match):
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
                            if re.search(r'^[0-9]+$', data_sub[1]):
                                video_width = data_sub[1] + 'px'
                            else:
                                video_width = data_sub[1]
                        elif data_sub[0] == 'height':
                            if re.search(r'^[0-9]+$', data_sub[1]):
                                video_height = data_sub[1] + 'px'
                            else:
                                video_height = data_sub[1]
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

                return '<iframe style="width: ' + video_width + '; height: ' + video_height + ';" src="' + video_code + '" frameborder="0" allowfullscreen></iframe>'
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

                # add color
                if color != '':
                    sub_text = '<span style="color:' + color + ';">' + sub_text + '</span>'

                return '' + \
                    '<ruby>' + \
                        main_text + \
                        '<rp>(</rp>' + \
                        '<rt>' + \
                            sub_text + \
                        '</rt>' + \
                        '<rp>)</rp>' + \
                    '</ruby>' + \
                ''
            elif name_data == 'age':
                if re.search(r'^[0-9]{4}-[0-9]{2}-[0-9]{2}$', match[1]):
                    try:
                        date = datetime.datetime.strptime(match[1], '%Y-%m-%d')
                    except:
                        return 'invalid date'

                    date_now = datetime.datetime.today()

                    if date > date_now:
                        return 'invalid date'
                    else:
                        return str((date_now - date).days // 365)
                else:
                    return 'invalid date'
            elif name_data == 'dday':
                if re.search(r'^[0-9]{4}-[0-9]{2}-[0-9]{2}$', match[1]):
                    try:
                        date = datetime.datetime.strptime(match[1], '%Y-%m-%d')
                    except:
                        return 'invalid date'

                    date_now = datetime.datetime.today()
                    
                    date_end = (date_now - date).days
                    if date_end > 0:
                        return '+' + str(date_end)
                    else:
                        if date_end == 0:
                            return '-' + str(date_end)
                        else:
                            return str(date_end)
                else:
                    return 'invalid date'
            else:
                return '<macro>' + match[0] + '(' + match[1] + ')' + '</macro>'

        # double macro replace
        self.render_data = re.sub(r'\[([^[(]+)\(([^()]+)\)\]', do_render_macro_double, self.render_data)

        # single macro function
        def do_render_macro_single(match):
            match = match.group(1)

            if match in ('date', 'datetime'):
                return get_time()
            elif match == 'br':
                return '<br>'
            elif match == 'clearfix':
                return '<div style="clear: both;"></div>'
            else:
                return '<macro>' + match + '</macro>'

        # single macro replace
        self.render_data = re.sub(r'\[([^[\]]+)\]', do_render_macro_single, self.render_data)

        # macro safe restore
        self.render_data = re.sub(r'<macro>', '[', self.render_data)
        self.render_data = re.sub(r'<\/macro>', ']', self.render_data)

    def do_render_last(self):
        # remove front_br and back_br
        self.render_data = re.sub(r'\n?<front_br>', '', self.render_data)
        self.render_data = re.sub(r'<back_br>\n?', '', self.render_data)
        
        # \n to <br>
        self.render_data = re.sub(r'\n', '<br>', self.render_data)

        # <render_n> restore
        self.render_data = self.get_tool_data_restore(self.render_data)

    def __call__(self):
        self.do_render_text()
        self.do_render_macro()
        self.do_render_heading()
        self.do_render_last()

        # print(self.data_temp_storage)
        # print(self.render_data)

        return [
            self.render_data, # HTML
            self.render_data_js, # JS
            [] # Other
        ]