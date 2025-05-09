from .func_tool import *

from typing import Any

class class_do_render_namumark:
    def __init__(
        self,
        conn,
        doc_name,
        doc_data,
        doc_set,
        lang_data,
        do_type = 'exter',
        parameter = {},
        parent = None
    ):
        self.conn = conn
        self.curs = self.conn.cursor()

        self.doc_data = doc_data.replace('\r', '')
        self.doc_name = doc_name
        self.doc_set = doc_set

        self.do_type = do_type
        self.parameter = parameter
        self.parent = parent

        self.lang_data = lang_data
        try:
            self.ip = ip_check()
        except:
            self.ip = '0.0.0.0'

        try:
            if 'main_css_bold' in flask.session:
                pass

            self.flask_session = flask.session
        except:
            self.flask_session = ''

        try:
            self.darkmode = flask.request.cookies.get('main_css_darkmode', '0')
            if self.darkmode == 'default':
                self.darkmode = '0'
        except:
            self.darkmode = '0'

        self.data_temp_storage = {}
        self.data_temp_storage_count = 0

        self.data_backlink : dict[str, Any] = {}

        self.data_math_count = 0
        self.data_redirect = 0
        self.link_count = 0

        self.data_toc = ''
        self.data_footnote = {}
        self.data_footnote_all = {}
        self.data_category = ''
        self.data_category_list = []

        self.render_data = self.doc_data
        if self.do_type == 'exter':
            self.render_data = html.escape(self.render_data)

        self.render_data = '<back_br>\n' + self.render_data + '\n<front_br>'
        self.render_data_js = ''

        self.curs.execute(db_change('select data from other where name = "link_case_insensitive"'))
        db_data = self.curs.fetchall()
        self.link_case_insensitive = ' collate nocase' if db_data and db_data[0][0] != '' else ''

    def get_tool_lang(self, name):
        if name in self.lang_data:
            return self.lang_data[name]
        else:
            return name + ' (RENDER LANG)'

    def get_tool_js_safe(self, data):
        data = data.replace('\n', '\\\\n')
        data = data.replace('\\', '\\\\')
        data = data.replace("'", "\\'")
        data = data.replace('"', '\\"')

        return data

    def get_tool_css_safe(self, data):
        return data.replace(';', '')

    def get_tool_data_storage(self, data_A = '', data_B = '', data_C = '', do_type = 'render'):        
        self.data_temp_storage_count += 1
        if do_type == 'render':
            data_name = 'render_' + str(self.data_temp_storage_count) + '_' + self.doc_set['doc_include']

            self.data_temp_storage[data_name] = data_A
            self.data_temp_storage['/' + data_name] = data_B
            self.data_temp_storage['revert_' + data_name] = data_C
        else:
            data_name = 'slash_' + str(self.data_temp_storage_count) + '_' + self.doc_set['doc_include']

            self.data_temp_storage[data_name] = data_A

        return data_name

    def get_tool_data_restore(self, data, do_type = 'all'):
        storage_count = self.data_temp_storage_count * 3
        if do_type == 'all':
            storage_regex = r'<(\/?(?:render|slash)_(?:[0-9]+)(?:[^<>]+))>'
        elif do_type == 'render':
            storage_regex = r'<(\/?(?:render)_(?:[0-9]+)(?:[^<>]+))>'
        else:
            storage_regex = r'<(\/?(?:slash)_(?:[0-9]+)(?:[^<>]+))>'

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
            storage_regex = r'(?:<((slash)_(?:[0-9]+)(?:[^<>]+))>|<((render)_(?:[0-9]+)(?:[^<>]+))>(?:(?:(?!<(?:\/?render_(?:[0-9]+)(?:[^<>]+))>).|\n)*)<\/render_(?:[0-9]+)(?:[^<>]+)>)'
        elif do_type == 'render':
            storage_regex = r'<((render)_(?:[0-9]+)(?:[^<>]+))>(?:(?:(?!<(?:\/?render_(?:[0-9]+)(?:[^<>]+))>).)*)<\/render_(?:[0-9]+)(?:[^<>]+)>'
        else:
            storage_regex = r'<((slash)_(?:[0-9]+)(?:[^<>]+))>'

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

        data = data.replace('<front_br>', '')
        data = data.replace('<back_br>', '')

        return data

    def get_tool_footnote_make(self):    
        data = ''
        for for_a in self.data_footnote:
            if data == '':
                data += '<div class="opennamu_footnote">'
            else:
                data += '<hr class="main_hr">'

            if len(self.data_footnote[for_a]['list']) > 1:
                data += '(' + for_a + ') '

                for for_b in self.data_footnote[for_a]['list']:
                    data += '<sup><a id="' + self.doc_set['doc_include'] + 'fn_' + for_b + '" href="#' + self.doc_set['doc_include'] + 'rfn_' + for_b + '">(' + for_b + ')</a></sup> '
            else:
                data += '<a id="' + self.doc_set['doc_include'] + 'fn_' + self.data_footnote[for_a]['list'][0] + '" href="#' + self.doc_set['doc_include'] + 'rfn_' + self.data_footnote[for_a]['list'][0] + '">(' + for_a + ') </a> '

            data += '<footnote_title id="' + self.doc_set['doc_include'] + 'fn_' + self.data_footnote[for_a]['list'][0] + '_title">' + self.data_footnote[for_a]['data'] + '</footnote_title>'

        if data != '':
            data += '</div>'

        self.data_footnote_all.update(self.data_footnote)
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
            if self.darkmode == '0':
                return data[0]
            else:
                return data[1]

    def get_tool_link_fix_sub(self, link_sub):
        if re.search(r'^:(분류|category):', link_sub, flags = re.I):
            link_sub = re.sub(r'^:(?P<in>분류|category):', '\\g<in>:', link_sub, flags = re.I)
        elif re.search(r'^:(파일|file):', link_sub, flags = re.I):
            link_sub = re.sub(r'^:(?P<in>파일|file):', '\\g<in>:', link_sub, flags = re.I)

        return link_sub

    def get_tool_link_fix(self, link_main, do_type = 'link'):
        if do_type == 'link':
            if link_main.startswith('../'):
                # ../하위문서 처리
                levels_up = link_main.count('../')  # ../ 개수 계산
                remaining_path = link_main.lstrip('../')  # ../ 제거 후 나머지 경로 추출

                # 현재 문서 이름에서 상위 경로로 이동
                link_main = self.doc_name
                for _ in range(levels_up):
                    link_main = re.sub(r'(\/[^/]+)$', '', link_main)

                # 상위 경로에 나머지 하위 문서 경로 추가
                if remaining_path:
                    link_main = f"{link_main}/{remaining_path}"
            elif re.search(r'^\/', link_main):
                # 절대 하위 문서 이동 처리
                link_main = re.sub(r'^\/', lambda x: (self.doc_name + '/'), link_main)
            elif re.search(r'^:(분류|category):', link_main, flags=re.I):
                # 분류 문서 처리
                link_main = re.sub(r'^:(분류|category):', 'category:', link_main, flags=re.I)
            elif re.search(r'^:(파일|file):', link_main, flags=re.I):
                # 파일 문서 처리
                link_main = re.sub(r'^:(파일|file):', 'file:', link_main, flags=re.I)
            elif re.search(r'^사용자:', link_main, flags=re.I):
                # 사용자 문서 처리
                link_main = re.sub(r'^사용자:', 'user:', link_main, flags=re.I)
        else:
            # 리다이렉트 처리 (do_type == 'redirect')
            if link_main.startswith('../'):
                # ../하위문서 처리 (리다이렉트 시)
                levels_up = link_main.count('../')
                remaining_path = link_main.lstrip('../')

                link_main = self.doc_name
                for _ in range(levels_up):
                    link_main = re.sub(r'(\/[^/]+)$', '', link_main)

                if remaining_path:
                    link_main = f"{link_main}/{remaining_path}"
            elif re.search(r'^\/', link_main):
                link_main = re.sub(r'^\/', lambda x: (self.doc_name + '/'), link_main)
            elif re.search(r'^분류:', link_main):
                link_main = re.sub(r'^분류:', 'category:', link_main)
            elif re.search(r'^사용자:', link_main):
                link_main = re.sub(r'^사용자:', 'user:', link_main)

        return link_main
    
    def do_inter_render(self, data, doc_include):
        doc_set = dict(self.doc_set)
        doc_set['doc_include'] = doc_include

        data_end = class_do_render_namumark(self.conn, self.doc_name, data, doc_set, self.lang_data, do_type = 'inter')()

        self.render_data_js += data_end[1]
        self.data_category_list += data_end[2]['category']
        self.data_backlink = dict(self.data_backlink, **data_end[2]['backlink_dict'])
        self.data_temp_storage = dict(self.data_temp_storage, **data_end[2]['temp_storage'][0])
        self.data_temp_storage_count += data_end[2]['temp_storage'][1]
        self.link_count += data_end[2]['link_count']

        return data_end[0]

    # Render
    def do_render_text(self):
        # <b> function
        bold_user_set = get_main_skin_set(self.conn, self.flask_session, 'main_css_bold', self.ip)

        def do_render_text_bold(match):
            data = match.group(1)
            if bold_user_set == 'delete':
                return ''
            elif bold_user_set == 'change':
                data_name = self.get_tool_data_storage('', '', match.group(0))
            else:
                data_name = self.get_tool_data_storage('<b>', '</b>', match.group(0))

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

        # <s> function
        strike_user_set = get_main_skin_set(self.conn, self.flask_session, 'main_css_strike', self.ip)

        def do_render_text_strike(match):
            data = match.group(1)
            if strike_user_set == 'delete':
                return ''
            elif strike_user_set == 'change':
                data_name = self.get_tool_data_storage('', '', match.group(0))
            else:
                data_name = self.get_tool_data_storage('<s>', '</s>', match.group(0))

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

            heading_data = re.search(heading_regex, self.render_data)
            if not heading_data:
                break
            elif heading_count_all < 0:
                print('Error : render heading count overflow')

                break
            else:
                heading_data_org = heading_data.group(0)
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

                    heading_folding = ['⊖', 'block', '1']
                    if heading_data[2]:
                        heading_folding = ['⊕', 'none', '0.5']

                    heading_id_name = 'edit_load_' + str(heading_count)
                    if self.doc_set["doc_type"] != "view":
                        heading_id_name = self.doc_set['doc_include'] + 'edit_load_' + str(heading_count)

                    data_name = self.get_tool_data_storage(
                        '<h' + heading_level_str + '><span id="' + self.doc_set['doc_include'] + 'opennamu_heading_' + str(heading_count) + '_sub" style="opacity: ' + heading_folding[2] + '">', 
                            ' <sub>' + \
                                '<a id="' + heading_id_name + '" href="/edit_section/' + str(heading_count) + '/' + url_pas(self.doc_name) + '">✎</a> ' + \
                                '<a href="javascript:void(0);" onclick="javascript:opennamu_heading_folding(\'' + self.doc_set['doc_include'] + 'opennamu_heading_' + str(heading_count) + '\', this);">' + \
                                    heading_folding[0] + \
                                '</a>'
                            '</sub>' + \
                        '</span></h' + heading_level_str + '>', 
                        heading_data_org
                    )

                    heading_data_complete = '' + \
                        '\n<front_br>' + \
                        ('</div>' if heading_count != 1 else '') + \
                        '<' + data_name + '>' + \
                            '<heading_stack>' + \
                                heading_stack_str + \
                            '</heading_stack>' + \
                            ' ' + heading_data_text + \
                        '</' + data_name + '>' + \
                        '<div id="' + self.doc_set['doc_include'] + 'opennamu_heading_' + str(heading_count) + '" style="display: ' + heading_folding[1] + ';">' + \
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
                ('<span style="margin-left: 10px;"></span>' * for_a[0].count('.')) + \
                '<span class="opennamu_TOC_list">' + \
                    '<a href="#s-' + for_a[0] + '">' + \
                        for_a[0] + '. ' + \
                    '</a>' + \
                    '<toc_inside>' + for_a[1] + '</toc_inside>' + \
                '</span>' + \
            ''

        if toc_data != '':
            toc_data += '</div>'

            self.data_toc = toc_data
            self.render_data += '<toc_data>' + toc_data + '</toc_data>'
        else:
            self.data_toc = ''

    def do_render_macro(self):
        # double macro function
        def do_render_macro_double(match):
            match_org = match
            match = match.groups()

            name_data = match[0]
            name_data = name_data.lower()

            macro_split_regex = r'(?:^|,) *([^,]+)'
            macro_split_sub_regex = r'(^[^=]+) *= *([^=]+)'
            if name_data in ('youtube', 'nicovideo', 'navertv', 'kakaotv', 'vimeo', 'instagram', 'twitter', 'tiktok', 'facebook'):
                data = re.findall(macro_split_regex, match[1])

                # get option
                video_code = ''
                video_start = ''
                video_end = ''
                
                video_width = '640px'
                video_height = '360px'
                if name_data == 'instagram'or name_data == 'tiktok':
                    video_width = '360px'
                    video_height = '480px'
                elif name_data == 'facebook':
                    video_width = '500px'
                    video_height = '616px'
                elif name_data == 'twitter':
                    video_width = '480px'
                    video_height = '480px'

                for for_a in data:
                    data_sub = re.search(macro_split_sub_regex, for_a)
                    if data_sub:
                        data_sub = data_sub.groups()
                        data_sub = [data_sub[0].lower(), data_sub[1]]

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
                elif name_data == 'instagram':
                    video_code = re.sub(r'^https:\/\/www\.instagram\.com\/p\/', '', video_code)

                    video_code = 'https://www.instagram.com/p/' + video_code +'/embed/'
                elif name_data == 'facebook':
                    video_code = 'https://www.facebook.com/plugins/post.php?href=' + video_code + '&width=' + video_width + '&height=' + video_height
                elif name_data == 'tiktok':
                    video_code = 'https://www.tiktok.com/embed/v2/' + video_code
                elif name_data == 'twitter':
                    video_code = 'https://twitframe.com/show?url=' + video_code

                    if self.darkmode == '1':
                        video_code += '&theme=dark'
                elif name_data == 'kakaotv':
                    video_code = re.sub(r'^https:\/\/tv\.kakao\.com\/v\/', '', video_code)

                    video_code = 'https://tv.kakao.com/embed/player/cliplink/' + video_code +'?service=kakao_tv'
                elif name_data == 'navertv':
                    video_code = re.sub(r'^https:\/\/tv\.naver\.com\/v\/', '', video_code)

                    video_code = 'https://tv.naver.com/embed/' + video_code
                elif name_data == 'nicovideo':
                    video_code = 'https://embed.nicovideo.jp/watch/' + video_code
                else:
                    video_code = 'https://player.vimeo.com/video/' + video_code

                video_width = self.get_tool_css_safe(video_width)
                video_height = self.get_tool_css_safe(video_height)

                data_name = self.get_tool_data_storage(
                    '<iframe style="width: ' + video_width + '; height: ' + video_height + ';" src="' + video_code + '" frameborder="0" allowfullscreen>',
                    '</iframe>',
                    match_org.group(0)
                )

                return '<' + data_name + '></' + data_name + '>'
            elif name_data == 'toc':
                return '<toc_no_auto>'
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
                        data_sub = [data_sub[0].lower(), data_sub[1]]

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
            elif name_data == 'anchor':
                main_text = self.get_tool_data_revert(match[1], do_type = 'render')

                data_name = self.get_tool_data_storage('<span id="' + main_text + '">', '</span>', match_org.group(0))

                return '<' + data_name + '></' + data_name + '>'
            elif name_data == 'username':
                data = re.findall(macro_split_regex, match[1])

                # get option
                user_name = ''
                render = 1
                for for_a in data:
                    data_sub = re.search(macro_split_sub_regex, for_a)
                    if data_sub:
                        data_sub = data_sub.groups()
                        data_sub = [data_sub[0].lower(), data_sub[1]]

                        if data_sub[0] == 'load_name':
                            if data_sub[1] == '1':
                                user_name = self.ip
                        elif data_sub[0] == 'render':
                            if data_sub[1] == '0':
                                render = 0
                    else:
                        user_name = for_a

                data_name = self.get_tool_data_storage('<span class="' + ('opennamu_render_ip' if render == 1 else '') + '">' + user_name + '</span>', '', match_org.group(0))

                return '<' + data_name + '></' + data_name + '>'
            elif name_data == 'timeif':
                data = re.findall(macro_split_regex, match[1])

                main_text = ''
                before_text = ''
                after_text = ''
                for for_a in data:
                    data_sub = re.search(macro_split_sub_regex, for_a)
                    if data_sub:
                        data_sub = data_sub.groups()
                        data_sub = [data_sub[0].lower(), data_sub[1]]

                        if data_sub[0] == 'before':
                            before_text = data_sub[1]
                        elif data_sub[0] == 'after':
                            after_text = data_sub[1]
                    else:
                        main_text = for_a

                if re.search(r'^[0-9]{4}-[0-9]{2}-[0-9]{2}$', main_text):
                    try:
                        date = datetime.datetime.strptime(main_text, '%Y-%m-%d')
                        data_text = ''
                    except:
                        data_text = 'invalid date'

                    date_now = datetime.datetime.today()

                    if data_text == '':
                        if date > date_now:
                            data_text = before_text
                        else:
                            data_text = after_text
                else:
                    data_text = 'invalid date'

                return data_text
            elif name_data == 'age':
                if re.search(r'^[0-9]{4}-[0-9]{2}-[0-9]{2}$', match[1]):
                    try:
                        date = datetime.datetime.strptime(match[1], '%Y-%m-%d')
                        data_text = ''
                    except:
                        data_text = 'invalid date'

                    date_now = datetime.datetime.today()

                    if data_text == '':
                        if date > date_now:
                            data_text = 'invalid date'
                        else:
                            data_text = str((date_now - date).days // 365)
                else:
                    data_text = 'invalid date'

                data_name = self.get_tool_data_storage(data_text, '', match_org.group(0))

                return '<' + data_name + '></' + data_name + '>'
            elif name_data in ('dday', 'dmonth', 'dyear'):
                if re.search(r'^[0-9]{4}-[0-9]{2}-[0-9]{2}$', match[1]):
                    try:
                        date = datetime.datetime.strptime(match[1], '%Y-%m-%d')
                        data_text = ''
                    except:
                        data_text = 'invalid date'

                    date_now = datetime.datetime.today()
                    if data_text == '':
                        if name_data == 'dday':
                            date_end = (date_now - date).days
                        elif name_data == 'dyear':
                            date_end = (date_now - date)
                            if date_end.days > 0:
                                date_end = datetime.date.min + date_end
                                date_end = date_end.year - 1
                            else:
                                date_end = datetime.date.min - date_end
                                date_end = 0 - (date_end.year - 1)
                        else:
                            date_end = (date_now - date)
                            if date_end.days > 0:
                                date_end = datetime.date.min + date_end
                                date_end = (date_end.year - 1) * 12 + (date_end.month - 1)
                            else:
                                date_end = datetime.date.min - date_end
                                date_end = 0 - ((date_end.year - 1) * 12 + (date_end.month - 1))

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
            elif name_data == 'joke':
                data_name = self.get_tool_data_storage('<span class="opennamu_joke">', '</span>', match_org.group(0))

                return '<' + data_name + '>' + match[1] + '</' + data_name + '>'
            elif name_data == 'pagecount':
                return '0'
            elif name_data == 'lastedit':
                link_main = match[1]
                data_view = ''

                data = re.findall(macro_split_regex, match[1])
                for for_a in data:
                    data_sub = re.search(macro_split_sub_regex, for_a)
                    if data_sub:
                        data_sub = data_sub.groups()
                        data_sub = [data_sub[0].lower(), data_sub[1]]

                        if data_sub[0] == 'view':
                            if data_sub[1] == 'full':
                                data_view = '1'
                    else:
                        link_main = for_a
                        
                link_main = self.get_tool_link_fix(link_main, 'redirect')

                link_main = self.get_tool_data_restore(link_main, do_type = 'slash')
                link_main = html.unescape(link_main)

                self.curs.execute(db_change("select set_data from data_set where doc_name = ? and set_name = 'last_edit'"), [link_main])
                db_data = self.curs.fetchall()
                if db_data:
                    date_data = db_data[0][0]
                    if data_view != '1':
                        date_data = date_data.split()[0]

                    return date_data
                else:
                    return '0'
            else:
                return '<macro>' + match[0] + '(' + match[1] + ')' + '</macro>'

        # double macro replace
        self.render_data = re.sub(r'\[([^[(\]]+)\(((?:(?!\)\]).)+)\)\]', do_render_macro_double, self.render_data)

        # single macro function
        def do_render_macro_single(match):
            match_org = match

            match = match.group(1)
            match = match.lower()

            if match in ('date', 'datetime'):
                data_name = self.get_tool_data_storage(get_time(), '', match_org.group(0))

                return '<' + data_name + '></' + data_name + '>'
            elif match == 'br':
                data_name = self.get_tool_data_storage('<br>', '', match_org.group(0))

                return '<' + data_name + '></' + data_name + '>'
            elif match == 'clearfix':
                data_name = self.get_tool_data_storage('<div style="clear: both;"></div>', '', match_org.group(0))

                return '<' + data_name + '></' + data_name + '>'
            elif match in ('목차', 'toc', 'tableofcontents'):
                return '<toc_need_part>'
            elif match == 'pagecount':
                self.curs.execute(db_change('select data from other where name = "count_all_title"'))
                db_data = self.curs.fetchall()
                if db_data:
                    return db_data[0][0]
                else:
                    return '0'
            else:
                return '<macro>' + match_org.group(1) + '</macro>'

        # single macro replace
        self.render_data = re.sub(r'\[([^[\]]+)\]', do_render_macro_single, self.render_data)

        # macro safe restore
        self.render_data = self.render_data.replace('<macro>', '[')
        self.render_data = self.render_data.replace('</macro>', ']')

    def do_render_math(self):
        def do_render_math_sub(match):
            data = match.group(1)

            data = data.replace('\n', '')
            data = self.get_tool_data_revert(data)

            data_html = self.get_tool_js_safe(data)

            data = html.unescape(data)
            data = self.get_tool_js_safe(data)

            name_ob = self.doc_set['doc_include'] + 'opennamu_math_' + str(self.data_math_count)

            data_name = self.get_tool_data_storage('<span id="' + name_ob + '">' + data_html, '</span>', match.group(0))

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

        math_regex = re.compile(r'\[math\(((?:(?!\[math\(|\)\]).|\n)+)\)\]', re.I)
        math_regex_2 = re.compile(r'&lt;math&gt;((?:(?!&lt;math&gt;|&lt;\/math&gt;).)+)&lt;\/math&gt;', re.I)

        self.render_data = re.sub(math_regex_2, do_render_math_sub, self.render_data)
        self.render_data = re.sub(math_regex, do_render_math_sub, self.render_data)

    def do_render_link(self):
        self.render_data = self.render_data.replace('[[]]', '')

        link_regex = r'\[\[((?:(?!\[\[|\]\]|\||<|>).|<(?:\/?(?:slash)_(?:[0-9]+)(?:[^<>]+))>)+)(?:\|((?:(?!\[\[|\]\]|\|).)+))?\]\](\n?)'
        image_count = 0
        link_count_all = len(re.findall(link_regex, self.render_data)) * 4
        while 1:
            link_data = re.search(link_regex, self.render_data)
            if not link_data:
                break
            elif link_count_all < 0:
                print('Error : render link count overflow')
                break
            else:
                # link split
                link_data_full = link_data.group(0)
                link_data = link_data.groups()

                link_main = link_data[0]
                link_main_org = link_main

                # file link
                if re.search(r'^(파일|file|외부|out):', link_main, flags = re.I):
                    file_width = ''
                    file_height = ''
                    file_align = ''
                    file_bgcolor = ''
                    file_turn = ''
                    file_radius = ''
                    file_rendering = ''

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
                                    if data_sub[1] in ('center', 'left', 'right'):
                                        file_align = data_sub[1]
                                elif data_sub[0] == 'bgcolor':
                                    file_bgcolor = data_sub[1]
                                elif data_sub[0] == 'theme':
                                    if data_sub[1] == 'dark':
                                        file_turn = 'dark'
                                    elif data_sub[1] == 'light':
                                        file_turn = 'light'
                                elif data_sub[0] == 'border-radius':
                                    file_radius = self.get_tool_px_add_check(data_sub[1])
                                elif data_sub[0] == 'rendering':
                                    if data_sub[1] == 'pixelated':
                                        file_rendering = 'pixelated'

                    link_main_org = ''
                    link_sub = link_main
                    file_out = 0

                    link_out_regex = re.compile('^(외부|out):', re.I)
                    link_in_regex = re.compile('^(파일|file):', re.I)
                    if re.search(link_out_regex, link_main):
                        link_main = re.sub(link_out_regex, '', link_main)

                        link_main = self.get_tool_data_restore(link_main, do_type = 'slash')
                        link_main = html.unescape(link_main)
                        link_main = link_main.replace('"', '&quot;')
                        
                        link_exist = ''
                        file_out = 1
                    else:
                        link_main = re.sub(link_in_regex, '', link_main)

                        link_main = self.get_tool_data_restore(link_main, do_type = 'slash')
                        link_main = html.unescape(link_main)

                        if not ('file:' + link_main) in self.data_backlink:
                            self.data_backlink['file:' + link_main] = {}

                        self.curs.execute(db_change("select title from data where title = ?"), ['file:' + link_main])
                        db_data = self.curs.fetchall()
                        if db_data:
                            link_exist = ''
                        else:
                            link_exist = 'opennamu_not_exist_link'
                            self.data_backlink['file:' + link_main]['no'] = ''

                        self.curs.execute(db_change('select id from history where title = ? order by date desc limit 1'), ['file:' + link_main])
                        db_data = self.curs.fetchall()
                        rev = db_data[0][0] if db_data else '1' 

                        self.data_backlink['file:' + link_main]['file'] = ''

                        link_extension_regex = r'\.([^.]+)$'
                        link_extension = re.search(link_extension_regex, link_main)
                        if link_extension:
                            link_extension = link_extension.group(1)
                        else:
                            link_extension = 'jpg'

                        link_main = re.sub(link_extension_regex, '', link_main)
                        link_main_org = link_main

                        link_main = '/image/' + url_pas(sha224_replace(link_main)) + '.' + link_extension + '.cache_v' + rev

                    if file_width != '':
                        file_width = 'width:' + self.get_tool_css_safe(file_width) + ';'
                    
                    if file_height != '':
                        file_height = 'height:' + self.get_tool_css_safe(file_height) + ';'

                    file_align_style = ''
                    if file_align in ('left', 'right'):
                        file_align_style = 'float:' + file_align + ';'

                    if file_bgcolor != '':
                        file_bgcolor = 'background:' + self.get_tool_css_safe(file_bgcolor) + ';'

                    if file_radius != '':
                        file_radius = 'border-radius:' + self.get_tool_css_safe(file_radius) + ';'

                    if file_rendering != '':
                        file_rendering = 'image-rendering:' + self.get_tool_css_safe(file_rendering) + ';'

                    file_style = file_width + file_height + file_align_style + file_bgcolor + file_radius + file_rendering

                    image_set = get_main_skin_set(self.conn, self.flask_session, 'main_css_image_set', self.ip)
                    if image_set == 'new_click' or image_set == 'click':
                        file_end = '<img style="' + file_style + '" id="opennamu_image_' + str(image_count) + '" alt="' + link_sub + '" src="">'
                    else:
                        file_end = '<img style="' + file_style + '" alt="' + link_sub + '" src="' + link_main + '">'

                    if file_align == 'center':
                        file_end = '<div style="text-align:center;">' + file_end + '</div>'

                    if link_exist != '':
                        data_name = self.get_tool_data_storage('<a class="' + link_exist + '" title="' + link_sub + '" href="/upload?name=' + url_pas(link_main_org) + '">(' + link_sub, ')</a>', link_data_full)
                        self.render_data = re.sub(link_regex, '<' + data_name + '></' + data_name + '>', self.render_data, 1)
                    else:
                        file_pass = 0
                        if file_turn != '':
                            if file_turn == 'dark' and self.darkmode == '1':
                                file_pass = 1
                            elif file_turn == 'light' and self.darkmode == '0':
                                file_pass = 1
                        else:
                            file_pass = 1

                        if file_pass == 1:
                            if file_out == 0:
                                file_link = '/w/file:' + url_pas(link_main_org) + '.' + url_pas(link_extension)
                            else:
                                file_link = link_main

                            if image_set == 'new_click':
                                data_name = self.get_tool_data_storage('<a title="' + link_sub + '" id="opennamu_image_' + str(image_count) + '_link" href="javascript:void(0);">' + file_end, '</a>', link_data_full)
                                self.render_data_js += '''
                                    document.getElementById("opennamu_image_''' + str(image_count) + '''_link").addEventListener("click", function(e) {
                                        document.getElementById("opennamu_image_''' + str(image_count) + '''").src = "''' + self.get_tool_js_safe(link_main) + '''";
                                        setTimeout(function() {
                                            document.getElementById("opennamu_image_''' + str(image_count) + '''_link").href = "''' + self.get_tool_js_safe(file_link) + '''";
                                        }, 100);
                                    });\n
                                '''
                                image_count += 1
                            else:
                                data_name = self.get_tool_data_storage('<a title="' + link_sub + '" href="' + file_link + '">' + file_end, '</a>', link_data_full)
                        else:
                            data_name = self.get_tool_data_storage('', '', link_data_full)
                        
                        self.render_data = re.sub(link_regex, '<' + data_name + '></' + data_name + '>' + link_data[2], self.render_data, 1)
                # category
                elif re.search(r'^(분류|category):', link_main, flags = re.I):
                    link_main = re.sub(r'^(분류|category):', '', link_main, flags = re.I)

                    category_blur = ''
                    if re.search(r'#blur$', link_main, flags = re.I):
                        link_main = re.sub(r'#blur$', '', link_main, flags = re.I)
                        category_blur = 'opennamu_category_blur'
                    
                    link_sub = link_main
                    link_view = ''
                    if len(link_data) > 1 and link_data[1]:
                        link_view = link_data[1]

                    link_main = self.get_tool_data_restore(link_main, do_type = 'slash')
                    link_main = html.unescape(link_main)

                    if not link_main in self.data_category_list:
                        self.data_category_list += [link_main]
                        
                        if not ('category:' + link_main) in self.data_backlink:
                            self.data_backlink['category:' + link_main] = {}

                        self.curs.execute(db_change("select title from data where title = ?"), ['category:' + link_main])
                        db_data = self.curs.fetchall()
                        if db_data:
                            link_exist = ''
                        else:
                            link_exist = 'opennamu_not_exist_link'
                            self.data_backlink['category:' + link_main]['no'] = ''

                        self.data_backlink['category:' + link_main]['cat'] = ''
                        
                        if link_view != '':
                            self.data_backlink['category:' + link_main]['cat_view'] = link_view
                        
                        if category_blur != '':
                            self.data_backlink['category:' + link_main]['cat_blur'] = ''

                        link_main = url_pas(link_main)

                        if self.data_category == '':
                            self.data_category = '' + \
                                '<div class="opennamu_category" id="cate">' + \
                                    '<a class="opennamu_category_button" href="javascript:opennamu_do_category_spread();"> (+)</a>' + \
                                    self.get_tool_lang('category') + ' : ' + \
                                '' + \
                            ''
                        else:
                            self.data_category += ' | '

                        self.data_category += '<a class="' + category_blur + ' ' + link_exist + '" title="' + link_sub + '" href="/w/category:' + link_main + '">' + link_sub + '</a>'

                    self.render_data = re.sub(link_regex, '', self.render_data, 1)
                # inter link
                elif re.search(r'^(?:inter|인터):([^:]+):', link_main, flags = re.I):
                    link_inter_regex = re.compile('^(?:inter|인터):([^:]+):', flags = re.I)

                    link_inter_name = re.search(link_inter_regex, link_main)
                    if link_inter_name:
                        link_inter_name = link_inter_name.group(1)
                    else:
                        link_inter_name = ''

                    link_main = re.sub(link_inter_regex, '', link_main)
                    link_title = link_inter_name + ':' + link_main

                    # sharp
                    link_main = link_main.replace('&#x27;', '<link_single>')
                    link_data_sharp_regex = r'#([^#]+)$'
                    link_data_sharp = re.search(link_data_sharp_regex, link_main)
                    if link_data_sharp:
                        link_data_sharp = link_data_sharp.group(1)
                        link_data_sharp = html.unescape(link_data_sharp)
                        link_data_sharp = '#' + url_pas(link_data_sharp)

                        link_main = re.sub(link_data_sharp_regex, '', link_main)
                    else:
                        link_data_sharp = ''
                    
                    link_main = link_main.replace('<link_single>', '&#x27;')

                    # main link fix
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
                            link_sub_storage = re.sub(link_inter_regex, '', link_sub_storage)

                        link_inter_icon = link_inter_name + ':'
                        if db_data[0][1] != '':
                            link_inter_icon = db_data[0][1]

                        link_sub_storage = link_inter_icon + link_sub_storage

                        self.curs.execute(db_change("select plus_t from html_filter where kind = 'inter_wiki_sub' and html = ?"), [link_inter_name])
                        db_data = self.curs.fetchall()
                        if db_data and db_data[0][0] == 'under_bar':
                            link_main = link_main.replace('%20', '_')

                        data_name = self.get_tool_data_storage('<a class="opennamu_link_inter" title="' + link_title + '" href="' + link_main + link_data_sharp + '">' + link_sub_storage, '</a>', link_data_full)
                    
                        add_str = ''
                        if link_data[2]:
                            add_str = link_data[2]

                        self.render_data = re.sub(link_regex, lambda x : ('<' + data_name + '>' + link_sub + '</' + data_name + '>' + add_str), self.render_data, 1)
                    else:
                        self.render_data = re.sub(link_regex, link_data[2], self.render_data, 1)
                # out link
                elif re.search(r'^https?:\/\/', link_main, flags = re.I):
                    link_main = self.get_tool_data_restore(link_main, do_type = 'slash')
                    link_title = link_main
                    link_main = html.unescape(link_main)

                    link_main = link_main.replace('"', '&quot;')
                    link_main = link_main.replace('<', '&lt;')
                    link_main = link_main.replace('>', '&gt;')

                    # sub not exist -> sub = main
                    if link_data[1]:
                        link_sub = link_data[1]
                        link_sub_storage = ''
                    else:
                        link_sub = ''
                        link_sub_storage = link_main_org

                    domain = ''
                    try:
                        domain = urllib.parse.urlparse(link_main).netloc
                    except:
                        pass

                    link_inter_icon = ''
                    link_class = 'opennamu_link_out'

                    self.curs.execute(db_change("select html, plus_t from html_filter where kind = 'outer_link' and plus = ?"), [domain])
                    db_data = self.curs.fetchall()
                    if db_data:
                        if db_data[0][1] != '':
                            if re.search(r'<|>', db_data[0][1]):
                                link_inter_icon = db_data[0][1]
                                link_class = 'opennamu_link_inter'
                            else:
                                if self.get_tool_data_restore(link_sub).find('"' + db_data[0][1] + '"') != -1:
                                    link_inter_icon = ''
                                    link_class = 'opennamu_link_inter'
                                else:
                                    link_inter_icon = '<img src="' + db_data[0][1] + '">'
                                    link_class = 'opennamu_link_inter'
                        else:
                            link_inter_icon = db_data[0][0] + ':'
                            link_class = 'opennamu_link_inter'

                    add_str = ''
                    if link_data[2]:
                        add_str = link_data[2]

                    data_name = self.get_tool_data_storage('<a class="' + link_class + '" target="_blank" title="' + link_title + '" href="' + link_main + '">' + link_inter_icon + link_sub_storage, '</a>', link_data_full)
                    self.render_data = re.sub(link_regex, lambda x : ('<' + data_name + '>' + link_sub + '</' + data_name + '>' + add_str), self.render_data, 1)
                # in link
                else:
                    # under page & fix url
                    link_main = self.get_tool_link_fix(link_main)

                    # sharp
                    link_main = link_main.replace('&#x27;', '<link_single>')
                    link_data_sharp_regex = r'#([^#]+)$'
                    link_data_sharp = re.search(link_data_sharp_regex, link_main)
                    if link_data_sharp:
                        link_data_sharp = link_data_sharp.group(1)
                        link_data_sharp = html.unescape(link_data_sharp)
                        link_data_sharp = '#' + url_pas(link_data_sharp)

                        link_main = re.sub(link_data_sharp_regex, '', link_main)
                    else:
                        link_data_sharp = ''
                    
                    link_main = link_main.replace('<link_single>', '&#x27;')

                    # main link fix
                    link_main = self.get_tool_data_restore(link_main, do_type = 'slash')
                    link_main = html.unescape(link_main)

                    # link_title
                    link_title = html.escape(link_main + link_data_sharp)

                    link_exist = ''
                    if link_main != '':
                        self.curs.execute(db_change("select title from data where title = ?" + self.link_case_insensitive), [link_main])
                        db_data = self.curs.fetchall()
                        if not db_data:
                            if not link_main in self.data_backlink:
                                self.data_backlink[link_main] = {}

                            self.data_backlink[link_main]['no'] = ''
                            link_exist = 'opennamu_not_exist_link'
                        else:
                            link_main = db_data[0][0]
                            if not link_main in self.data_backlink:
                                self.data_backlink[link_main] = {}
                        
                        self.data_backlink[link_main][''] = ''

                    link_same = ''
                    if link_main == self.doc_name:
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
                        link_sub_storage = self.get_tool_link_fix_sub(link_main_org)

                    self.link_count += 1

                    add_str = ''
                    if link_data[2]:
                        add_str = link_data[2]

                    data_name = self.get_tool_data_storage('<a class="' + link_exist + ' ' + link_same + '" title="' + link_title + '" href="' + link_main + link_data_sharp + '">' + link_sub_storage, '</a>', link_data_full)
                    self.render_data = re.sub(link_regex, lambda x : ('<' + data_name + '>' + link_sub + '</' + data_name + '>' + add_str), self.render_data, 1)

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

                if match[1] in self.parameter:
                    return slash_add + self.parameter[match[1]]
                else:
                    return slash_add + match[2]

        self.render_data = re.sub(r'(\\+)?@([ㄱ-힣a-zA-Z0-9_]+)=((?:\\@|[^@\n])+)@', do_render_include_default_sub, self.render_data)
        self.render_data = re.sub(r'(\\+)?@([ㄱ-힣a-zA-Z0-9_]+)@', do_render_include_default_sub, self.render_data)

    def do_render_include(self):
        include_num = 0
        include_set_data = get_main_skin_set(self.conn, self.flask_session, 'main_css_include_link', self.ip)
        include_regex = re.compile(r'\[include\(((?:(?!\[include\(|\)\]|<\/div>).)+)\)\](\n?)', re.I)
        include_count_max = len(re.findall(include_regex, self.render_data)) * 2
        while 1:
            include_num += 1
            include_change_list = {}

            match = re.search(include_regex, self.render_data)
            if include_count_max < 0:
                break
            elif not match:
                break
            else:
                if self.doc_set['doc_type'] == "include":
                    self.render_data = re.sub(include_regex, '', self.render_data, 1)
                else:
                    match_org = match.group(0)
                    match = match.groups()

                    macro_split_regex = r'(?:^|,) *([^,]+)'
                    macro_split_sub_regex = r'^([^=]+) *= *(.*)$'

                    include_name = ''

                    data = re.findall(macro_split_regex, match[0])
                    for for_a in data:
                        data_sub = re.search(macro_split_sub_regex, for_a)
                        if data_sub:
                            data_sub = data_sub.groups()
                            
                            data_sub_name = data_sub[0]
                            data_sub_data = self.get_tool_data_restore(data_sub[1], do_type = 'slash')
                            data_sub_data = html.unescape(data_sub_data)
                            
                            data_sub_data = re.sub(r'^(?P<in>분류|category):', ':\\g<in>:', data_sub_data)
                            data_sub_data = re.sub(r'^(?P<in>파일|file):', ':\\g<in>:', data_sub_data)

                            include_change_list[data_sub_name] = data_sub_data
                        else:
                            include_name = for_a

                    include_name_org = include_name
                    
                    include_name = self.get_tool_data_restore(include_name, do_type = 'slash')
                    include_name = html.unescape(include_name)

                    if not include_name in self.data_backlink:
                        self.data_backlink[include_name] = {}

                    self.data_backlink[include_name]['include'] = ''

                    # load include db data
                    self.curs.execute(db_change("select data from data where title = ?"), [include_name])
                    db_data = self.curs.fetchall()
                    if db_data:
                        # include link func
                        include_link = ''
                        if include_set_data == 'use':
                            include_link = '<a href="/w/' + url_pas(include_name) + '">(' + include_name_org + ')</a><br>'

                        include_data = ''
                        if self.parent:
                            include_data_tmp = self.parent(
                                self.conn,
                                doc_name = self.doc_name,
                                doc_data = db_data[0][0], 
                                data_type = 'api_include',
                                parameter = include_change_list
                            )

                            include_data = include_data_tmp[0] + '<script>window.addEventListener("DOMContentLoaded", function() {' + include_data_tmp[1] + '});</script>'

                        data_name = self.get_tool_data_storage(include_link + include_data, '', match_org)
                    else:
                        self.data_backlink[include_name]['no'] = ''

                        include_link = '<a class="opennamu_not_exist_link" href="/w/' + url_pas(include_name) + '">(' + include_name_org + ')</a>'
                        data_name = self.get_tool_data_storage(include_link, '', match_org)

                    self.render_data = re.sub(include_regex, '<' + data_name + '></' + data_name + '>' + match[1], self.render_data, 1)

            include_count_max -= 1

    def do_redner_footnote(self):
        footnote_num = 0

        footnote_set = get_main_skin_set(self.conn, self.flask_session, 'main_css_footnote_set', self.ip)
        footnote_number_set = get_main_skin_set(self.conn, self.flask_session, 'main_css_footnote_number', self.ip)
        footnote_number_view_set = get_main_skin_set(self.conn, self.flask_session, 'main_css_view_real_footnote_num', self.ip)

        footnote_regex = re.compile(r'(?:\[\*((?:(?!\[\*|\]| ).)+)?(?: ((?:(?!\[\*|\]).)+))?\]|\[(각주|footnote)\])', re.I)
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
                    footnote_num_str = str(footnote_num)

                    if not footnote_data[0]:
                        footnote_name = footnote_num_str
                        footnote_name_add = ''
                    else:
                        footnote_name = footnote_data[0]
                        footnote_name_add = ' (' + footnote_num_str + ')'

                    if not footnote_data[1]:
                        footnote_text_data = ''
                    else:
                        footnote_text_data = footnote_data[1]

                    fn = ''
                    rfn = ''
                    foot_v_name = ''

                    if footnote_name in self.data_footnote_all or footnote_name in self.data_footnote:
                        if footnote_name in self.data_footnote:
                            self.data_footnote[footnote_name]['list'] += [footnote_num_str]
                            footnote_first = self.data_footnote[footnote_name]['list'][0]
                        else:
                            self.data_footnote[footnote_name] = {}
                            self.data_footnote[footnote_name]['list'] = [footnote_num_str]
                            self.data_footnote[footnote_name]['data'] = footnote_text_data
                            footnote_first = self.data_footnote_all[footnote_name]['list'][0]

                        fn = self.doc_set['doc_include'] + 'fn_' + footnote_first
                        rfn = self.doc_set['doc_include'] + 'rfn_' + footnote_num_str

                        if footnote_number_set == 'only_number':
                            foot_v_name += footnote_first
                        else:
                            foot_v_name += footnote_name
                            
                        if footnote_number_view_set == 'on':
                            foot_v_name += ' (' + footnote_num_str + ')'
                    else:
                        self.data_footnote[footnote_name] = {}
                        self.data_footnote[footnote_name]['list'] = [footnote_num_str]
                        self.data_footnote[footnote_name]['data'] = footnote_text_data

                        fn = self.doc_set['doc_include'] + 'fn_' + footnote_num_str
                        rfn = self.doc_set['doc_include'] + 'rfn_' + footnote_num_str

                        if footnote_number_set == 'only_number':
                            foot_v_name += footnote_num_str
                        else:
                            foot_v_name += footnote_name
                            
                        if footnote_number_view_set == 'on':
                            foot_v_name += footnote_name_add

                    if footnote_set == 'spread':
                        data_name = self.get_tool_data_storage(
                            '<sup>' + \
                                '<a fn_target="' + fn + '" id="' + rfn + '" href="javascript:void(0);">(' + foot_v_name + ')</a>' + \
                            '</sup>' + \
                            '<span class="opennamu_spead_footnote" id="' + rfn + '_load" style="display: none;"></span>',
                            '',
                            footnote_data_org
                        )
                        self.render_data_js += 'document.getElementById("' + rfn + '").addEventListener("click", function() { opennamu_do_footnote_spread("' + rfn + '", "' + fn + '"); });\n'
                    elif footnote_set == 'popup':
                        data_name = self.get_tool_data_storage(
                            '<sup>' + \
                                '<a fn_target="' + fn + '" id="' + rfn + '" href="javascript:void(0);">(' + foot_v_name + ')</a>' + \
                            '</sup>' + \
                            '<span class="opennamu_spead_footnote" id="' + rfn + '_load" style="display: none;"></span>',
                            '',
                            footnote_data_org
                        )
                        self.render_data_js += 'document.getElementById("' + rfn + '").addEventListener("click", function() { opennamu_do_footnote_spread("' + rfn + '", "' + fn + '"); });\n'
                    elif footnote_set == 'popover':
                        data_name = self.get_tool_data_storage(
                            '<span id="' + rfn + '_over">' + \
                                '<sup>' + \
                                    '<a fn_target="' + fn + '" id="' + rfn + '" href="javascript:void(0);">(' + foot_v_name + ')</a>' + \
                                '</sup>' + \
                                '<span class="opennamu_popup_footnote" id="' + rfn + '_load" style="display: none;"></span>' + \
                            '</span>',
                            '',
                            footnote_data_org
                        )
                        self.render_data_js += 'document.getElementById("' + rfn + '_over").addEventListener("click", function() { opennamu_do_footnote_popover("' + rfn + '", "' + fn + '", undefined, "open"); });\n'
                        self.render_data_js += 'document.addEventListener("click", function() { opennamu_do_footnote_popover("' + rfn + '", "' + fn + '", undefined, "close"); });\n'
                    else:
                        data_name = self.get_tool_data_storage('<sup><a fn_target="' + fn + '" id="' + rfn + '" href="#' + fn + '">(' + foot_v_name + ')</a></sup>', '', footnote_data_org)

                    self.render_data = re.sub(footnote_regex, '<' + data_name + '></' + data_name + '>', self.render_data, 1)

            footnote_count_all -= 1

        self.render_data += '<footnote_category>'
        self.render_data += self.get_tool_footnote_make()

    def do_render_redirect(self):
        match = re.search(r'^<back_br>\n#(?:redirect|넘겨주기) ([^\n]+)', self.render_data, flags = re.I)
        if match:
            if self.doc_set['doc_type'] == 'view':
                link_data_full = match.group(0)
                link_main = match.group(1)

                link_inter_name = ''

                link_inter_regex = re.compile('^(?:inter|인터):([^:]+):', flags = re.I)
                inter_check = re.search(link_inter_regex, link_main)
                if not inter_check:
                    # under page & fix url
                    link_main = self.get_tool_link_fix(link_main, 'redirect')
                else:
                    link_inter_name = inter_check.group(1)
                    link_main = re.sub(link_inter_regex, '', link_main)

                # sharp
                link_main = link_main.replace('&#x27;', '<link_single>')
                link_data_sharp_regex = r'#([^#]+)$'
                link_data_sharp = re.search(link_data_sharp_regex, link_main)
                if link_data_sharp:
                    link_data_sharp = link_data_sharp.group(1)
                    link_data_sharp = html.unescape(link_data_sharp)
                    link_data_sharp = '#' + url_pas(link_data_sharp)

                    link_main = re.sub(link_data_sharp_regex, '', link_main)
                else:
                    link_data_sharp = ''
                
                link_main = link_main.replace('<link_single>', '&#x27;')

                if not inter_check:
                    # main link fix
                    link_main = self.get_tool_data_restore(link_main, do_type = 'slash')
                    link_main = html.unescape(link_main)

                    self.curs.execute(db_change("select title from data where title = ?" + self.link_case_insensitive), [link_main])
                    db_data = self.curs.fetchall()
                    if not db_data:
                        if not link_main in self.data_backlink:
                            self.data_backlink[link_main] = {}

                        self.data_backlink[link_main]['no'] = ''
                    else:
                        link_main = db_data[0][0]
                        if not link_main in self.data_backlink:
                            self.data_backlink[link_main] = {}

                    self.data_backlink[link_main]['redirect'] = link_data_sharp

                    link_main = url_pas(link_main)
                    if link_main != '':
                        link_main = '/w_from/' + link_main

                    self.data_redirect = 1
                    data_name = self.get_tool_data_storage('<a href="' + link_main + link_data_sharp + '">(GO)</a>', '', link_data_full)

                    self.render_data = '<' + data_name + '></' + data_name + '>'
                else:
                    self.curs.execute(db_change("select plus, plus_t from html_filter where kind = 'inter_wiki' and html = ?"), [link_inter_name])
                    db_data = self.curs.fetchall()
                    if db_data:
                        link_main = url_pas(link_main)
                        link_main = db_data[0][0] + link_main

                        link_sub_storage = match.group(1)
                        link_sub_storage = re.sub(link_inter_regex, '', link_sub_storage)

                        link_inter_icon = link_inter_name + ':'
                        if db_data[0][1] != '':
                            link_inter_icon = db_data[0][1]

                        link_sub_storage = link_inter_icon + link_sub_storage

                        self.curs.execute(db_change("select plus_t from html_filter where kind = 'inter_wiki_sub' and html = ?"), [link_inter_name])
                        db_data = self.curs.fetchall()
                        if db_data and db_data[0][0] == 'under_bar':
                            link_main = link_main.replace('%20', '_')

                        self.data_redirect = 1
                        if 'doc_from' in self.doc_set:
                            data_name = self.get_tool_data_storage('<a href="' + link_main + link_data_sharp + '">(GO)</a>', '', link_data_full)
                        else:
                            data_name = self.get_tool_data_storage('<meta http-equiv="refresh" content="5; url=' + link_main + link_data_sharp + '">', link_sub_storage + ' - After 5s', link_data_full)
                    
                        self.render_data = '<' + data_name + '></' + data_name + '>'
                    else:
                        self.data_redirect = 1
                        self.render_data = ''
            else:
                self.data_redirect = 1
                self.render_data = ''

    def do_render_table(self):
        self.render_data = re.sub(r'\n +\|\|', '\n||', self.render_data)

        # get_tool_dark_mode_split
        # get_tool_px_add_check
        # get_tool_css_safe
        # todo : after text render text not use to make table
        def do_render_table_parameter(cell_count, parameter, data, option = {}):
            table_parameter_all = { "div" : "", "class" : "", "table" : "", "tr" : "", "td" : "", "col" : "", "colspan" : "", "rowspan" : "", "data" : "" }
            
            table_align_auto = 1
            table_colspan_auto = 1

            do_any_thing = ''

            table_parameter_regex = r'&lt;((?:(?!&lt;|&gt;).)+)&gt;'
            for table_parameter in re.findall(table_parameter_regex, parameter):
                table_parameter_split = table_parameter.split('=')
                if len(table_parameter_split) == 2:
                    table_parameter_name = table_parameter_split[0].replace(' ', '')
                    table_parameter_name = table_parameter_name.lower()
                    
                    table_parameter_data = self.get_tool_css_safe(table_parameter_split[1])

                    if table_parameter_name == 'tablebgcolor':
                        table_parameter_all['table'] += 'background:' + self.get_tool_dark_mode_split(table_parameter_data) + ';'
                    elif table_parameter_name == 'tablewidth':
                        table_parameter_all['div'] += 'width:' + self.get_tool_px_add_check(table_parameter_data) + ';'
                        table_parameter_all['table'] += 'width:100%;'
                    elif table_parameter_name == 'tableheight':
                        table_parameter_all['table'] += 'height:' + self.get_tool_px_add_check(table_parameter_data) + ';'
                    elif table_parameter_name == 'tablealign':
                        if table_parameter_data == 'right':
                            table_parameter_all['div'] += 'float:right;'
                        elif table_parameter_data == 'center':
                            table_parameter_all['div'] += 'margin:auto;'
                            table_parameter_all['table'] += 'margin:auto;'
                    elif table_parameter_name == 'tableclass':
                        table_parameter_all['class'] = table_parameter_split[1]
                    elif table_parameter_name == 'tabletextalign':
                        table_parameter_all['table'] += 'text-align:' + table_parameter_data + ' !important;'
                    elif table_parameter_name == 'tablecolor':
                        table_parameter_all['table'] += 'color:' + self.get_tool_dark_mode_split(table_parameter_data) + ';'
                    elif table_parameter_name == 'tablebordercolor':
                        table_parameter_all['table'] += 'border:2px solid ' + self.get_tool_dark_mode_split(table_parameter_data) + ';'
                    elif table_parameter_name == 'rowbgcolor':
                        table_parameter_all['tr'] += 'background:' + self.get_tool_dark_mode_split(table_parameter_data) + ';'
                    elif table_parameter_name == 'rowtextalign':
                        table_parameter_all['tr'] += 'text-align:' + table_parameter_data + ' !important;'
                    elif table_parameter_name == 'rowcolor':
                        table_parameter_all['tr'] += 'color:' + self.get_tool_dark_mode_split(table_parameter_data) + ';'
                    elif table_parameter_name == 'colcolor':
                        table_parameter_all['col'] += 'color:' + self.get_tool_dark_mode_split(table_parameter_data) + ';'
                    elif table_parameter_name == 'colbgcolor':
                        table_parameter_all['col'] += 'background:' + self.get_tool_dark_mode_split(table_parameter_data) + ';'
                    elif table_parameter_name == 'coltextalign':
                        table_parameter_all['col'] += 'text-align:' + table_parameter_data + ' !important;'
                    elif table_parameter_name == 'bgcolor':
                        table_parameter_all['td'] += 'background:' + self.get_tool_dark_mode_split(table_parameter_data) + ';'
                    elif table_parameter_name == 'color':
                        table_parameter_all['td'] += 'color:' + self.get_tool_dark_mode_split(table_parameter_data) + ';'
                    elif table_parameter_name == 'width':
                        table_parameter_all['td'] += 'width:' + self.get_tool_px_add_check(table_parameter_data) + ';'
                    elif table_parameter_name == 'height':
                        table_parameter_all['td'] += 'height:' + self.get_tool_px_add_check(table_parameter_data) + ';'
                    else:
                        do_any_thing += '&lt;' + table_parameter + '&gt;'
                elif len(table_parameter_split) == 1:
                    if table_parameter == 'keepall':
                        table_parameter_all['td'] += 'word-break: keep-all !important;'
                    elif table_parameter == 'rowkeepall':
                        table_parameter_all['tr'] += 'word-break: keep-all !important;'
                    elif table_parameter == 'colkeepall':
                        table_parameter_all['col'] += 'word-break: keep-all !important;'
                    elif table_parameter == 'nopad':
                        table_parameter_all['td'] += 'padding: 0 !important;'
                    elif re.search(r'^-[0-9]+$', table_parameter):
                        table_colspan_auto = 0
                        table_parameter_all['colspan'] = re.sub(r'[^0-9]+', '', table_parameter)
                    elif re.search(r'^(\^|v)?\|[0-9]+$', table_parameter):
                        if table_parameter[0] == '^':
                            table_parameter_all['td'] += 'vertical-align: top;'
                        elif table_parameter[0] == 'v':
                            table_parameter_all['td'] += 'vertical-align: bottom;'

                        table_parameter_all['rowspan'] = re.sub(r'[^0-9]+', '', table_parameter)
                    elif table_parameter in ('(', ':', ')'):
                        table_align_auto = 0
                        if table_parameter == '(':
                            table_parameter_all['td'] += 'text-align: left !important;'
                        elif table_parameter == ':':
                            table_parameter_all['td'] += 'text-align: center !important;'
                        elif table_parameter == ')':
                            table_parameter_all['td'] += 'text-align: right !important;'
                    elif re.search(r'^(?:(?:#((?:[0-9a-f-A-F]{3}){1,2}))|(\w+))$', table_parameter):
                        table_parameter_data = self.get_tool_css_safe(table_parameter)
                        table_parameter_all['td'] += 'background:' + self.get_tool_dark_mode_split(table_parameter_data) + ';'
                    else:
                        do_any_thing += '&lt;' + table_parameter + '&gt;'
                else:
                    do_any_thing += '&lt;' + table_parameter + '&gt;'
            
            if table_align_auto == 1:
                if re.search(r'^ ', data):
                    data = re.sub(r'^ ', '', data)
                    if re.search(r' $', data):
                        table_parameter_all['td'] += 'text-align: center;'

                        data = re.sub(r' $', '', data)
                    else:
                        table_parameter_all['td'] += 'text-align: right;'
                else:
                    table_parameter_all['td'] += 'text-align: left;'
                    if re.search(r' $', data):
                        data = re.sub(r' $', '', data)
            else:
                data = re.sub(r'^ ', '', data)
                data = re.sub(r' $', '', data)

            if table_colspan_auto == 1:
                table_parameter_all['colspan'] = str(len(cell_count) // 2)

            table_parameter_all['data'] = do_any_thing + data

            return table_parameter_all

        table_regex = re.compile(r'\n((?:(?:(?:(?:\|\|)+)|(?:\|[^|]+\|(?:\|\|)*))\n?(?:(?:(?!\|\|).)+))(?:(?:\|\||\|\|\n|(?:\|\|)+(?!\n)(?:(?:(?!\|\|).)+)\n*)*)\|\|)\n', re.DOTALL)
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

                table_parameter = { "div" : "", "class" : "", "table" : "", "tr" : "", "td" : "", "col" : {}, "rowspan" : {} }
                table_data_end = ''
                table_col_num = 0
                table_tr_change = 0
                for table_sub in re.findall(table_sub_regex, table_data):
                    table_data_in = table_sub[3]
                    table_data_in = re.sub(r'^\n+', '', table_data_in)

                    if table_sub[0] != '' and table_tr_change == 1:
                        table_col_num = 0
                        table_data_end += '<tr style="' + table_parameter["tr"] + '">' + table_parameter["td"] + '</tr>'
                        
                        table_parameter["tr"] = ""
                        table_parameter["td"] = ""

                    table_sub_parameter = do_render_table_parameter(table_sub[1], table_sub[2], table_data_in)
                    table_parameter["tr"] += table_sub_parameter['tr']

                    if not table_col_num in table_parameter['rowspan']:
                        table_parameter['rowspan'][table_col_num] = 0
                    else:
                        if table_parameter['rowspan'][table_col_num] != 0:
                            table_parameter['rowspan'][table_col_num] -= 1
                            table_col_num += int(table_sub_parameter['colspan'])

                    if table_sub_parameter['rowspan'] != '':
                        rowspan_int = int(table_sub_parameter['rowspan'])
                        if rowspan_int > 1:
                            table_parameter['rowspan'][table_col_num] = rowspan_int - 1

                    if not table_col_num in table_parameter['col']:
                        table_parameter['col'][table_col_num] = ''

                    table_parameter['div'] += table_sub_parameter['div']
                    table_parameter['class'] = table_sub_parameter['class'] if table_sub_parameter['class'] != '' else table_parameter['class']
                    table_parameter['table'] += table_sub_parameter['table']
                    table_parameter['col'][table_col_num] += table_sub_parameter['col']

                    if table_sub[2] == '' and table_sub[3] == '':
                        table_tr_change = 1
                    else:
                        table_tr_change = 0
                    
                        table_parameter["td"] += '<td colspan="' + table_sub_parameter['colspan'] + '" rowspan="' + table_sub_parameter['rowspan'] + '" style="' + table_parameter['col'][table_col_num] + table_sub_parameter['td'] + '"><back_br>\n' + table_sub_parameter['data'] + '\n<front_br></td>'
                    
                    table_col_num += int(table_sub_parameter['colspan'])
                else:
                    table_data_end += '<tr style="' + table_parameter["tr"] + '">' + table_parameter["td"] + '</tr>'

                table_data_end = '<table class="' + table_parameter['class'] + '" style="' + table_parameter['table'] + '">' + table_caption + table_data_end + '</table>'
                table_data_end = '<div class="table_safe" style="' + table_parameter['div'] + '">' + table_data_end + '</div>'

                self.render_data = re.sub(table_regex, lambda x : ('\n' + table_data_end + '\n'), self.render_data, 1)

            table_count_all -= 1
    
    def do_render_middle(self):
        wiki_count = 0
        html_count = 0
        syntax_count = 0
        folding_count = 0

        inter_count = 0
        inter_data = {}

        middle_regex = r'{{{([^{](?:(?!{{{|}}}).|\n)*)?(?:}|<(\/?(?:slash)_(?:[0-9]+)(?:[^<>]+))>)}}'
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
                        middle_data_org = re.sub(r'<(\/?(?:slash)_(?:[0-9]+)(?:[^<>]+))>', '<temp_' + middle_slash + '>', middle_data_org)
                        self.render_data = re.sub(middle_regex, lambda x : middle_data_org, self.render_data, 1)
                        continue

                middle_data = middle_data.group(1)
                if not middle_data:
                    middle_data = ''

                middle_name = re.search(r'^([^ \n]+)', middle_data)
                middle_data_pass = ''
                middle_data_add = ''
                if middle_name:
                    middle_name = middle_name.group(1)
                    middle_name = middle_name.lower()
                    
                    if middle_name == "#!if":
                        if middle_slash:
                            middle_data_org = re.sub(r'<(\/?(?:slash)_(?:[0-9]+)(?:[^<>]+))>', '<temp_' + middle_slash + '>', middle_data_org)
                            self.render_data = re.sub(middle_regex, lambda x : middle_data_org, self.render_data, 1)
                            continue

                        if_data = re.sub(r'^#!if *', '', middle_data)

                        var_name = ""
                        if_state = ""
                        parameter = ""
                        
                        if_regex = r"([^ ]+)(?: *)(!=|==)(?: *)([^\n]+)"
                        if_data_get = re.search(if_regex, if_data)
                        if if_data_get:
                            if_data_get_data = if_data_get.groups()

                            var_name = if_data_get_data[0]
                            if_state = if_data_get_data[1]
                            parameter = if_data_get_data[2]

                            if_data = re.sub(if_regex, '', if_data)

                        var_data = None
                        if var_name in self.parameter:
                            var_data = self.parameter[var_name]

                        result_data = False

                        if if_state == "==":
                            if parameter == "null" and var_data == None:
                                result_data = True
                            else:
                                parameter = re.sub(r'(^&quot;|&quot;$)', '', parameter)
                                if parameter == var_data:
                                    result_data = True
                        else:
                            if parameter == "null" and var_data != None:
                                result_data = True
                            else:
                                parameter = re.sub(r'(^&quot;|&quot;$)', '', parameter)
                                if parameter != var_data:
                                    result_data = True

                        if result_data:
                            middle_data_pass = if_data
                        else:
                            middle_data_pass = ""
                        
                        data_name = self.get_tool_data_storage('', '', middle_data_org)
                    elif middle_name == '#!wiki':
                        if middle_slash:
                            middle_data_org = re.sub(r'<(\/?(?:slash)_(?:[0-9]+)(?:[^<>]+))>', '<temp_' + middle_slash + '>', middle_data_org)
                            self.render_data = re.sub(middle_regex, lambda x : middle_data_org, self.render_data, 1)
                            continue

                        wiki_data = re.sub(r'^#!wiki *', '', middle_data)

                        wiki_regex = re.compile('^(?:(?:style=(&quot;(?:(?:(?!&quot;).)*)&quot;|&#x27;(?:(?:(?!&#x27;).)*)&#x27;)))(?:\n| +)', re.I)
                        wiki_dark_regex = re.compile('^(?:(?:dark-style=(&quot;(?:(?:(?!&quot;).)*)&quot;|&#x27;(?:(?:(?!&#x27;).)*)&#x27;)))(?:\n| +)', re.I)

                        wiki_data_style_data = ''
                        while 1:
                            dark_only = 0
                            wiki_data_style = re.search(wiki_regex, wiki_data)
                            if wiki_data_style:
                                wiki_data = re.sub(wiki_regex, '', wiki_data)
                            else:
                                wiki_data_style = re.search(wiki_dark_regex, wiki_data)
                                if wiki_data_style:
                                    dark_only = 1
                                    wiki_data = re.sub(wiki_dark_regex, '', wiki_data)
                                else:
                                    break
                        
                            if wiki_data_style:
                                wiki_data_style = wiki_data_style.group(1)
                                if wiki_data_style:
                                    wiki_data_style = wiki_data_style.replace('&#x27;', '')
                                    wiki_data_style = wiki_data_style.replace('&quot;', '')

                                    if dark_only == 1 and self.darkmode == '1':
                                        wiki_data_style_data += wiki_data_style
                                    elif dark_only == 0:
                                        wiki_data_style_data += wiki_data_style
                                else:
                                    wiki_data_style = ''
                            else:
                                wiki_data_style = ''

                        find_regex = [
                            r" *box-shadow *: *(([^,;]*)(,|;|$)){10,}",
                            r" *url\([^()]*\)",
                            r" *linear-gradient\((([^(),]+)(,|\))){10,}",
                            r" *position *: *"
                        ]
                        for for_a in find_regex:
                            if re.search(for_a, wiki_data_style_data, flags = re.I):
                                wiki_data_style_data = ""
                                break

                        wiki_data = self.get_tool_data_revert(wiki_data)
                        wiki_data = re.sub('(^\n|\n$)', '', wiki_data)

                        inter_data["inter_data_" + str(inter_count)] = wiki_data
                        middle_data_pass = '<inter_data_' + str(inter_count) + '>'
                        inter_count += 1

                        data_name = self.get_tool_data_storage('<div style="' + wiki_data_style_data + '">', '</div>', middle_data_org)
                        wiki_count += 1
                    elif middle_name == '#!html':
                        if middle_slash:
                            middle_data_org = re.sub(r'<(\/?(?:slash)_(?:[0-9]+)(?:[^<>]+))>', '<temp_' + middle_slash + '>', middle_data_org)
                            self.render_data = re.sub(middle_regex, lambda x : middle_data_org, self.render_data, 1)
                            continue

                        html_data = re.sub(r'^#!html( |\n)', '', middle_data)
                        if middle_slash:
                            html_data += '\\'

                        data_revert = self.get_tool_data_revert(html_data)
                        data_revert = re.sub(r'^\n', '', data_revert)
                        data_revert = re.sub(r'\n$', '', data_revert)
                        data_revert = data_revert.replace('&amp;nbsp;', '&nbsp;')

                        self.render_data_js += 'opennamu_do_render_html("' + self.doc_set['doc_include'] + 'opennamu_wiki_' + str(html_count) + '");\n'

                        data_name = self.get_tool_data_storage('<span id="' + self.doc_set['doc_include'] + 'opennamu_wiki_' + str(html_count) + '">' + data_revert, '</span>', middle_data_org)
                        html_count += 1
                    elif middle_name == '#!folding':
                        if middle_slash:
                            middle_data_org = re.sub(r'<(\/?(?:slash)_(?:[0-9]+)(?:[^<>]+))>', '<temp_' + middle_slash + '>', middle_data_org)
                            self.render_data = re.sub(middle_regex, lambda x : middle_data_org, self.render_data, 1)
                            continue

                        wiki_regex = re.compile('^#!folding(?: ([^\n]*))?\n', re.I)
                        wiki_data_folding = re.search(wiki_regex, middle_data)
                        wiki_data = re.sub(wiki_regex, '', middle_data)
                        if wiki_data_folding:
                            wiki_data_folding = wiki_data_folding.group(1)
                            if not wiki_data_folding:
                                wiki_data_folding = 'test'
                        else:
                            wiki_data_folding = 'test'

                        wiki_data = self.get_tool_data_revert(wiki_data)
                        wiki_data = re.sub('(^\n|\n$)', '', wiki_data)

                        inter_data["inter_data_" + str(inter_count)] = wiki_data
                        wiki_data_end = '<inter_data_' + str(inter_count) + '>'
                        inter_count += 1

                        middle_data_pass = wiki_data_folding
                        data_name = self.get_tool_data_storage('<details><summary>', '</summary><div class="opennamu_folding">', middle_data_org)

                        data_name_2 = self.get_tool_data_storage('', '</div></details>', '')
                        middle_data_add = '<' + data_name_2 + '>' + wiki_data_end + '</' + data_name_2 + '>'

                        folding_count += 1
                    elif middle_name == '#!syntax':
                        if middle_slash:
                            middle_data_org = re.sub(r'<(\/?(?:slash)_(?:[0-9]+)(?:[^<>]+))>', '<temp_' + middle_slash + '>', middle_data_org)
                            self.render_data = re.sub(middle_regex, lambda x : middle_data_org, self.render_data, 1)
                            continue

                        wiki_regex = re.compile('^#!syntax(?: ([^\n]*))?\n', re.I)
                        wiki_data_syntax = re.search(wiki_regex, middle_data)
                        wiki_data = re.sub(wiki_regex, '', middle_data)
                        if wiki_data_syntax:
                            wiki_data_syntax = wiki_data_syntax.group(1)
                            if not wiki_data_syntax:
                                wiki_data_syntax = 'python'
                            else:
                                if wiki_data_syntax == 'asm' or wiki_data_syntax == 'assembly':
                                    wiki_data_syntax = 'x86arm'
                        else:
                            wiki_data_syntax = 'python'

                        if syntax_count == 0:
                            self.render_data_js += 'hljs.highlightAll();\n'
                            
                        if wiki_data.count('\n') > 10:
                            self.render_data_js += 'hljs.lineNumbersBlock(document.getElementById("opennamu_syntax_' + str(syntax_count) + '"));\n'

                        data_name = self.get_tool_data_storage('<pre id="syntax"><code class="' + wiki_data_syntax + '" id="opennamu_syntax_' + str(syntax_count) + '">' + wiki_data, '</code></pre>', middle_data_org)
                        
                        syntax_count += 1
                    elif middle_name in ('#!dark', '#!white'):
                        if middle_slash:
                            middle_data_org = re.sub(r'<(\/?(?:slash)_(?:[0-9]+)(?:[^<>]+))>', '<temp_' + middle_slash + '>', middle_data_org)
                            self.render_data = re.sub(middle_regex, lambda x : middle_data_org, self.render_data, 1)
                            continue

                        wiki_data = re.sub(r'^#!(dark|white)( |\n)', '', middle_data)
                        if middle_name == '#!dark' and self.darkmode == '1':
                            middle_data_pass = wiki_data
                        elif middle_name == '#!white' and self.darkmode == '0':
                            middle_data_pass = wiki_data
                        else:
                            middle_data_pass = ''
                        
                        data_name = self.get_tool_data_storage('', '', middle_data_org)
                    elif middle_name in ('+5', '+4', '+3', '+2', '+1', '-5', '-4', '-3', '-2', '-1'):
                        if middle_slash:
                            middle_data_org = re.sub(r'<(\/?(?:slash)_(?:[0-9]+)(?:[^<>]+))>', '<temp_' + middle_slash + '>', middle_data_org)
                            self.render_data = re.sub(middle_regex, lambda x : middle_data_org, self.render_data, 1)
                            continue

                        wiki_data = re.sub(r'^(\+|\-)[1-5]( |\n)', '', middle_data)
                        if middle_name == '+5':
                            wiki_size = '200'
                        elif middle_name == '+4':
                            wiki_size = '180'
                        elif middle_name == '+3':
                            wiki_size = '160'
                        elif middle_name == '+2':
                            wiki_size = '140'
                        elif middle_name == '+1':
                            wiki_size = '120'
                        elif middle_name == '-5':
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
                            middle_data_org = re.sub(r'<(\/?(?:slash)_(?:[0-9]+)(?:[^<>]+))>', '<temp_' + middle_slash + '>', middle_data_org)
                            self.render_data = re.sub(middle_regex, lambda x : middle_data_org, self.render_data, 1)
                            continue

                        wiki_color = re.search(r'^@(?:((?:[0-9a-f-A-F]{3}){1,2})|(\w+))(,@(?:((?:[0-9a-f-A-F]{3}){1,2})|(\w+)))?', middle_name)
                        wiki_color_data = ''
                        if wiki_color:
                            wiki_color = wiki_color.groups()
                            if wiki_color[0]:
                                wiki_color_data += '#' + wiki_color[0]
                            else:
                                wiki_color_data += wiki_color[1]

                            if wiki_color[2]:
                                if wiki_color[3]:
                                    wiki_color_data += ',#' + wiki_color[3]
                                elif wiki_color[4]:
                                    wiki_color_data += ',' + wiki_color[4]
                        else:
                            wiki_color_data += 'red'

                        wiki_color = self.get_tool_css_safe(wiki_color_data)
                        wiki_color = self.get_tool_dark_mode_split(wiki_color)

                        wiki_data = re.sub(r'^@(?:((?:[0-9a-f-A-F]{3}){1,2})|(\w+))(?:,@(?:((?:[0-9a-f-A-F]{3}){1,2})|(\w+)))?( |\n)', '', middle_data)

                        middle_data_pass = wiki_data
                        data_name = self.get_tool_data_storage('<span style="background-color:' + wiki_color + '">', '</span>', middle_data_org)
                    elif re.search(r'^#(?:((?:[0-9a-f-A-F]{3}){1,2})|(\w+))', middle_name):
                        if middle_slash:
                            middle_data_org = re.sub(r'<(\/?(?:slash)_(?:[0-9]+)(?:[^<>]+))>', '<temp_' + middle_slash + '>', middle_data_org)
                            self.render_data = re.sub(middle_regex, lambda x : middle_data_org, self.render_data, 1)
                            continue

                        wiki_color = re.search(r'^#(?:((?:[0-9a-f-A-F]{3}){1,2})|(\w+))(,#(?:((?:[0-9a-f-A-F]{3}){1,2})|(\w+)))?', middle_name)
                        wiki_color_data = ''
                        if wiki_color:
                            wiki_color = wiki_color.groups()
                            if wiki_color[0]:
                                wiki_color_data += '#' + wiki_color[0]
                            else:
                                wiki_color_data += wiki_color[1]

                            if wiki_color[2]:
                                if wiki_color[3]:
                                    wiki_color_data += ',#' + wiki_color[3]
                                elif wiki_color[4]:
                                    wiki_color_data += ',' + wiki_color[4]
                        else:
                            wiki_color_data += 'red'

                        wiki_color = self.get_tool_css_safe(wiki_color_data)
                        wiki_color = self.get_tool_dark_mode_split(wiki_color)

                        wiki_data = re.sub(r'^#(?:((?:[0-9a-f-A-F]{3}){1,2})|(\w+))(?:,#(?:((?:[0-9a-f-A-F]{3}){1,2})|(\w+)))?( |\n)', '', middle_data)

                        middle_data_pass = wiki_data
                        data_name = self.get_tool_data_storage('<span style="color:' + wiki_color + '">', '</span>', middle_data_org)
                    else:
                        if middle_slash:
                            middle_data += '\\'

                        data_revert = self.get_tool_data_revert(middle_data)
                        data_revert = re.sub(r'^\n', '', data_revert)
                        data_revert = re.sub(r'\n$', '', data_revert)

                        data_name = self.get_tool_data_storage(data_revert, '', middle_data_org)
                else:
                    if middle_slash:
                        middle_data += '\\'

                    data_revert = self.get_tool_data_revert(middle_data)
                    data_revert = re.sub(r'^\n', '', data_revert)
                    data_revert = re.sub(r'\n$', '', data_revert)

                    data_name = self.get_tool_data_storage(data_revert, '', middle_data_org)

                self.render_data = re.sub(middle_regex, lambda x : ('<' + data_name + '>' + middle_data_pass + '</' + data_name + '>' + middle_data_add), self.render_data, 1)

            middle_count_all -= 1

        inter_data_regex = r'<(inter_data_[0-9]+)>'
        class do_render_middle_replace_inter_class:
            def __init__(self, func, doc_set):
                self.inter_count = 0
                self.do_inter_render = func
                self.doc_set = doc_set

            def replace_sub(self, match):
                data = inter_data[match[1]]
                data = re.sub(inter_data_regex, self.replace_sub, data)

                return data

            def replace(self, match):
                data = inter_data[match[1]]
                data = re.sub(inter_data_regex, self.replace_sub, data)
                
                data = self.do_inter_render(data, self.doc_set['doc_include'] + 'opennamu_inter_render_' + str(self.inter_count))
                data = re.sub(r'\|\|', '<no_td>', data)
                
                self.inter_count += 1

                return data

        replace_inter_data_top = do_render_middle_replace_inter_class(self.do_inter_render, self.doc_set)
        self.render_data = re.sub(inter_data_regex, replace_inter_data_top.replace, self.render_data)
        self.render_data = re.sub(r'<temp_(?P<in>(?:slash)_(?:[0-9]+)(?:[^<>]+))>', '<\\g<in>>', self.render_data)

    def do_render_hr(self):
        hr_regex = r'\n-{4,9}\n'
        hr_count_max = len(re.findall(hr_regex, self.render_data)) * 3
        while 1:
            hr_data = re.search(hr_regex, self.render_data)
            if hr_count_max < 0:
                break
            elif not hr_data:
                break
            else:
                self.render_data = re.sub(hr_regex, '\n<front_br><hr><back_br>\n', self.render_data, 1)

            hr_count_max -= 1

    def do_render_list(self):
        # 인용문
        quote_regex = r'((?:\n&gt; *[^\n]*)+)\n'
        quote_count = 0
        quote_count_max = len(re.findall(quote_regex, self.render_data)) * 10
        while 1:
            quote_data = re.search(quote_regex, self.render_data)
            if quote_count_max < 0:
                break
            elif not quote_data:
                break
            else:
                quote_data_org = quote_data.group(0)
                
                quote_data = quote_data.group(1)
                quote_data = re.sub(r'\n&gt; *(?P<in>[^\n]*)', '\\g<in>\n', quote_data)
                quote_data = re.sub(r'\n$', '', quote_data)
                quote_data = self.get_tool_data_revert(quote_data)

                quote_data_end = self.do_inter_render(quote_data, self.doc_set['doc_include'] + 'opennamu_quote_' + str(quote_count))
                data_name = self.get_tool_data_storage('<div>', '</div>', quote_data_org)

                self.render_data = re.sub(quote_regex, lambda x : ('\n<front_br><hr class="mini_hr"><blockquote><back_br>\n<' + data_name + '>' + quote_data_end + '</' + data_name + '><front_br></blockquote><hr class="mini_hr"><back_br>\n'), self.render_data, 1)

            quote_count_max -= 1
            quote_count += 1

        # 리스트 공통 파트
        def int_to_alpha(num):
            alpha_list = string.ascii_lowercase
            alpha_len = len(alpha_list)
            end_text = ''

            while num:
                end_text = alpha_list[num % alpha_len - 1] + end_text
                num = num // alpha_len

            return end_text
 
        # https://www.geeksforgeeks.org/python-program-to-convert-integer-to-roman/
        def int_to_roman(number):
            num = [1, 4, 5, 9, 10, 40, 50, 90, 100, 400, 500, 900, 1000]
            sym = ["I", "IV", "V", "IX", "X", "XL", "L", "XC", "C", "CD", "D", "CM", "M"]
            i = 12
            end_text = ''

            while number:
                div = number // num[i]
                number %= num[i]
        
                while div:
                    end_text += sym[i]
                    div -= 1

                i -= 1

            return end_text
            
        list_view_set = get_main_skin_set(self.conn, self.flask_session, 'main_css_list_view_change', self.ip)
        
        list_style = {
            1 : 'opennamu_list_1',
            2 : 'opennamu_list_2',
            3 : 'opennamu_list_3',
            4 : 'opennamu_list_4'
        }
        class do_render_list_int_to:
            def __init__(self, list_view_set = ''):
                self.list_num = {}
                self.list_view_set = list_view_set

            def __call__(self, match):
                if match.group(4):
                    list_data = match.group(5)
                    list_len = len(match.group(1))
                    if list_len == 0:
                        list_len = 1

                    list_style_data = 'opennamu_list_5'
                    if list_len in list_style:
                        list_style_data = list_style[list_len]

                    return '<li style="margin-left: ' + str((list_len - 1) * 20) + 'px;" class="' + list_style_data + '">' + list_data + '</li>'
                else:
                    list_type = match.group(2)

                    do_type = 'int'
                    if list_type == 'a':
                        do_type = 'alpha_small'
                    elif list_type == 'A':
                        do_type = 'alpha_big'
                    elif list_type == 'i':
                        do_type = 'roman_small'
                    elif list_type == 'I':
                        do_type = 'roman_big'

                    if not do_type in self.list_num:
                        self.list_num[do_type] = []
                    
                    for for_a in self.list_num:
                        if for_a != do_type:
                            self.list_num[for_a] = []

                    list_data = match.group(5)
                    list_start = match.group(3)
                    list_len = len(match.group(1))
                    if list_len == 0:
                        list_len = 1

                    if len(self.list_num[do_type]) >= list_len:
                        self.list_num[do_type][list_len - 1] += 1

                        for for_a in range(list_len, len(self.list_num[do_type])):
                            self.list_num[do_type][for_a] = 0
                    else:
                        self.list_num[do_type] += [1] * (list_len - len(self.list_num[do_type]))

                    if list_start:
                        self.list_num[do_type][list_len - 1] = int(list_start)

                    if do_type == 'int':
                        if self.list_view_set == 'on':
                            change_text = str('-'.join([str(for_a) for for_a in self.list_num[do_type] if for_a != 0]))
                        else:
                            change_text = str(self.list_num[do_type][list_len - 1])
                    elif do_type == 'roman_big':
                        change_text = int_to_roman(self.list_num[do_type][list_len - 1]).upper()
                    elif do_type == 'roman_small':
                        change_text = int_to_roman(self.list_num[do_type][list_len - 1]).lower()
                    elif do_type == 'alpha_big':
                        change_text = int_to_alpha(self.list_num[do_type][list_len - 1]).upper()
                    else:
                        change_text = int_to_alpha(self.list_num[do_type][list_len - 1]).lower()

                    return '<li style="margin-left: ' + str((list_len - 1) * 20) + 'px;" class="opennamu_list_none">' + change_text + '. ' + list_data + '</li>'

        list_regex = r'((?:\n( *)(?:(\*)) ?([^\n]*))+|(?:\n( *)(?:(1|a|A|I|i)\.(?:#([0-9]*))?) ?([^\n]*)){2,})\n'
        list_count_max = len(re.findall(list_regex, self.render_data)) * 3
        while 1:
            list_data = re.search(list_regex, self.render_data)
            if list_count_max < 0:
                break
            elif not list_data:
                break
            else:
                list_data = list_data.group(1)
                list_sub_regex = r'\n( *)(?:(1|a|A|I|i)\.(?:#([0-9]*))?|(\*)) ?([^\n]*)'

                list_class = do_render_list_int_to(list_view_set)
                list_data_str = re.sub(list_sub_regex, list_class, list_data)

                self.render_data = re.sub(list_regex, lambda x : ('\n<front_br><ul>' + list_data_str + '</ul><back_br>\n'), self.render_data, 1)

            list_count_max -= 1

    def do_render_remark(self):
        self.render_data = re.sub(r'\n##[^\n]+', '\n<front_br>', self.render_data)

    def do_render_last(self):
        self.render_data = self.render_data.replace('<no_br>', '\n')
        self.render_data = self.render_data.replace('<no_td>', '||')

        # add category
        if self.doc_set["doc_type"] == 'view':
            if self.data_category != '':
                data_name = self.get_tool_data_storage(self.data_category, '</div>', '')

                category_set_data = get_main_skin_set(self.conn, self.flask_session, 'main_css_category_set', self.ip)
                if category_set_data == 'bottom':
                    if re.search(r'<footnote_category>', self.render_data):
                        self.render_data = self.render_data.replace('<footnote_category>', '<hr><' + data_name + '></' + data_name + '>', 1)
                    else:
                        self.render_data += '<hr><' + data_name + '></' + data_name + '>'
                else:
                    self.render_data = self.render_data.replace('<footnote_category>', '', 1)
                    self.render_data = '<' + data_name + '></' + data_name + '><hr class="main_hr">' + self.render_data
            else:
                self.render_data = self.render_data.replace(r'<footnote_category>', '', 1)
        else:
            self.render_data = self.render_data.replace(r'<footnote_category>', '', 1)

        # remove front_br and back_br
        self.render_data = re.sub(r'\n?<front_br>', '', self.render_data)
        self.render_data = re.sub(r'<back_br>\n?', '', self.render_data)
        
        # \n to <br>
        self.render_data = self.render_data.replace('\n', '<br>')

        # <render_n> restore
        self.render_data = self.get_tool_data_restore(self.render_data)

        # a fix
        self.temp_a_link_count = 0
        def do_render_last_a_link(match):
            data = match.group(1)
            if data == '</a>':
                if self.temp_a_link_count == 0:
                    return ''
                elif self.temp_a_link_count > 1:
                    self.temp_a_link_count -= 1
                    
                    return ''
                else:
                    self.temp_a_link_count -= 1
                    
                    return match.group(0)
            else:
                if self.temp_a_link_count > 0:
                    self.temp_a_link_count += 1
                    
                    return ''
                else:
                    self.temp_a_link_count += 1
                    
                    return match.group(0)
            
        self.render_data = re.sub(r'(<a(?: [^<>]*)?>|<\/a>)', do_render_last_a_link, self.render_data)
        
        def do_render_last_toc_filter(match):
            data = match.group(1).split(' ')
            if data[0] == 'a' or data[0] == '/a':
                if len(data) > 3 and data[3] != 'href="javascript:void(0);"':
                    return '<' + match[1] + '>'
                else:
                    return ''
            else:
                return ''

        # add toc
        def do_render_last_toc(match):
            data = match.group(1)
            
            data_sub = re.sub(r'<([^<>]*)>', '', data)
            data = re.sub(r'<([^<>]*)>', do_render_last_toc_filter, data)

            heading_regex = r'<h([1-6])>'
            heading_data = re.search(heading_regex, self.render_data)
            if heading_data:
                heading_data = heading_data.group(1)
                self.render_data = re.sub(heading_regex, lambda x : ('<h' + heading_data + ' id="' + data_sub + '">'), self.render_data, 1)
            
            return data

        if self.data_toc != '':
            self.render_data += '</div>'
            toc_search_regex = r'<toc_data>((?:(?!<toc_data>|<\/toc_data>).)*)<\/toc_data>'

            toc_data_on = 0

            toc_data = re.search(toc_search_regex, self.render_data)
            toc_data = toc_data.group(1) if toc_data else ''

            self.data_toc = toc_data
            self.data_toc = re.sub(r'<toc_inside>((?:(?!<toc_inside>|<\/toc_inside>).)*)<\/toc_inside>', do_render_last_toc, self.data_toc)

            toc_set_data = get_main_skin_set(self.conn, self.flask_session, 'main_css_toc_set', self.ip)

            self.render_data = re.sub(toc_search_regex, '', self.render_data)
            if toc_set_data == 'off':
                self.render_data = self.render_data.replace('<toc_need_part>', '')
            else:
                if re.search(r'<toc_need_part>', self.render_data):
                    toc_data_on = 1

                self.render_data = self.render_data.replace('<toc_need_part>', self.data_toc, 20)
                self.render_data = self.render_data.replace('<toc_need_part>', '')

            if self.doc_set["doc_type"] != 'view' or re.search(r'<toc_no_auto>', self.render_data) or toc_set_data == 'half_off' or toc_set_data == 'off' or toc_data_on == 1:
                self.render_data = self.render_data.replace('<toc_no_auto>', '')
            else:
                self.render_data = re.sub(r'(?P<in><h[1-6] id="[^"]*">)', '<br>' + self.data_toc + '\\g<in>', self.render_data, 1)
        else:
            self.render_data = self.render_data.replace('<toc_need_part>', '')
            self.render_data = self.render_data.replace('<toc_no_auto>', '')

        def do_render_last_footnote(match):
            match = match.group(1)

            find_regex = re.compile(r'<footnote_title id="' + match + r'_title">((?:(?!<footnote_title|<\/footnote_title>).)*)<\/footnote_title>')
            find_data = re.search(find_regex, self.render_data)
            if find_data:
                find_data = find_data.group(1)
                find_data = re.sub(r'<[^<>]*>', '', find_data)
            else:
                find_data = ''

            return '<a title="' + find_data + '"'

        self.render_data = re.sub(r'<a fn_target="([^"]+)"', do_render_last_footnote, self.render_data)

        self.render_data_js += '''
            document.querySelectorAll('details').forEach((el) => {
                new Accordion(el);
            });
            if(window.location.hash !== '' && document.getElementById(window.location.hash.replace(/^#/, ''))) {
                document.getElementById(window.location.hash.replace(/^#/, '')).focus();
            }\n
            opennamu_do_ip_render();\n
        '''

    def __call__(self):
        self.do_render_remark()
        self.do_render_include_default()
        self.do_render_slash()

        if self.do_type == 'exter':
            self.do_render_redirect()
        
        if self.data_redirect == 0:
            self.do_render_middle()
            self.do_render_include()
            self.do_render_math()
            self.do_render_table()
            self.do_render_list()
            self.do_render_macro()
            self.do_render_link()
            self.do_render_text()
            self.do_render_hr()

            if self.do_type == 'exter':
                self.do_redner_footnote()
                self.do_render_heading()
            
        if self.do_type == 'exter':
            self.do_render_last()
        else:
            self.render_data = self.render_data.replace(r'\|\|', '<no_td>')
            self.render_data = self.render_data.replace('\n', '<no_br>')

        data_backlink_dict = self.data_backlink
        data_backlink_list = [[self.doc_name, for_a, for_b, self.data_backlink[for_a][for_b]] for for_a in self.data_backlink for for_b in self.data_backlink[for_a]]

        # 여기 수정시 do_inter_render도 수정 필요
        return [
            self.render_data, # html
            self.render_data_js, # js
            {
                'backlink' : data_backlink_list, # backlink
                'backlink_dict' : data_backlink_dict,
                'footnote' : self.data_footnote_all, # footnote
                'category' : self.data_category_list,
                'temp_storage' : [self.data_temp_storage, self.data_temp_storage_count],
                'link_count' : self.link_count,
                'redirect' : self.data_redirect
            } # other
        ]