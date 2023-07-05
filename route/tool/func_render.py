from .func_tool import *

from .func_render_namumark import class_do_render_namumark

# 커스텀 마크 언젠간 다시 추가 예정

class class_do_render:
    def __init__(self, conn, lang_data = {}):
        self.conn = conn

        if lang_data == '{}':
            lang_data = {
                'toc' : 'toc',
                'category' : 'category'
            }

        self.lang_data = lang_data

    def do_render(self, doc_name, doc_data, data_type, data_in):
        curs = self.conn.cursor()

        doc_set = {}
        if data_type == 'from':
            doc_set['doc_from'] = 'O'
            data_type = 'view'
        
        data_in = (data_in + '_') if data_in != '' else ''
        doc_set['doc_include'] = data_in
        rep_data = ''

        if rep_data == '' and doc_name != '':
            curs.execute(db_change("select set_data from data_set where doc_name = ? and set_name = 'document_markup'"), [doc_name])
            db_data = curs.fetchall()
            if db_data and db_data[0][0] != '' and db_data[0][0] != 'normal':
                rep_data = db_data[0][0]

        if rep_data == '':
            curs.execute(db_change('select data from other where name = "markup"'))
            db_data = curs.fetchall()
            rep_data = db_data[0][0] if db_data else 'namumark'

        if rep_data == 'namumark' or rep_data == 'namumark_beta':
            data_end = class_do_render_namumark(
                curs,
                doc_name,
                doc_data,
                doc_set,
                self.lang_data
            )()
        elif rep_data == 'raw':
            data_end = [
                html.escape(doc_data).replace('\n', '<br>'), 
                '', 
                {}
            ]
        else:
            data_end = [
                doc_data, 
                '', 
                {}
            ]

        if data_type == 'thread' or data_type == 'api_thread':
            data_end[0] = re.sub(
                r'&lt;topic_a&gt;(?P<in>(?:(?!&lt;\/topic_a&gt;).)+)&lt;\/topic_a&gt;',
                '<a href="\g<in>">\g<in></a>',
                data_end[0]
            )
            data_end[0] = re.sub(
                r'&lt;topic_call&gt;@(?P<in>(?:(?!&lt;\/topic_call&gt;).)+)&lt;\/topic_call&gt;',
                '<a href="/w/user:\g<in>">@\g<in></a>',
                data_end[0]
            )

        if data_type == 'backlink':
            if 'backlink' in data_end[2]:
                backlink = data_end[2]['backlink']
            else:
                backlink = []

            if backlink != []:
                curs.executemany(db_change("insert into back (link, title, type) values (?, ?, ?)"), data_end[2]['backlink'])
                curs.execute(db_change("delete from back where title = ? and type = 'no'"), [doc_name])

            self.conn.commit()

        return [
            data_end[0], 
            data_end[1],
            data_end[2]
        ]