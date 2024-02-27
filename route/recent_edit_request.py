from .tool.func import *

from .recent_change import recent_change_send_render

from .go_api_recent_edit_request import api_recent_edit_request

def recent_edit_request(db_set):
    with get_db_connect() as conn:
        div = ''
        div += '''
            <table id="main_table_set">
                <tbody>
                    <tr id="main_table_top_tr">
                        <td id="main_table_width">''' + load_lang('discussion_name') + '''</td>
                        <td id="main_table_width">''' + load_lang('editor') + '''</td>
                        <td id="main_table_width">''' + load_lang('time') + '''</td>
                    </tr>
        '''

        all_list = json.loads(api_recent_edit_request(db_set).data)
        for data in all_list:
            if re.search(r"\+", data[5]):
                leng = '<span style="color:green;">(' + data[5] + ')</span>'
            elif re.search(r"\-", data[5]):
                leng = '<span style="color:red;">(' + data[5] + ')</span>'
            else:
                leng = '<span style="color:gray;">(' + data[5] + ')</span>'

            send = data[4]
            ip = data[6]
            date = data[2]

            title = '<a href="/edit_request/' + url_pas(data[0]) + '">' + html.escape(data[0]) + '</a> '
            title += '<a href="/history/' + url_pas(data[0]) + '">(r' + data[1] + ')</a> '

            div += '''
                <tr>
                    <td>''' + title + ' ' + leng + '''</td>
                    <td>''' + ip + '''</td>
                    <td>''' + date + '''</td>
                </tr>
                <tr>
                    <td colspan="3">''' + recent_change_send_render(html.escape(send)) + '''</td>
                </tr>
            '''

        div += '' + \
                '</tbody>' + \
            '</table>' + \
        ''

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('recent_edit_request'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = div,
            menu = [['recent_change', load_lang('return')]]
        ))