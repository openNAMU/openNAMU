from .tool.func import *

def recent_discuss(tool):
    with get_db_connect() as conn:
        curs = conn.cursor()

        div = ''
        admin_auth = admin_check(3)

        if tool == 'normal':
            div += '<a href="/recent_discuss/close">(' + load_lang('close_discussion') + ')</a> '
            div += '<a href="/recent_discuss/open">(' + load_lang('open_discussion_list') + ')</a>'

            m_sub = 0
        elif tool == 'close':
            div += '<a href="/recent_discuss">(' + load_lang('normal') + ')</a>'

            m_sub = ' (' + load_lang('closed') + ')'
        else:
            div += '<a href="/recent_discuss">(' + load_lang('normal') + ')</a>'

            m_sub = ' (' + load_lang('open_discussion_list') + ')'

        div +=  '''
                <hr class="main_hr">
                <table id="main_table_set">
                    <tbody>
                        <tr id="main_table_top_tr">
                            <td id="main_table_width">''' + load_lang('discussion_name') + '''</td>
                            <td id="main_table_width">''' + load_lang('editor') + '''</td>
                            <td id="main_table_width">''' + load_lang('time') + '''</td>
                        </tr>
                '''

        if tool == 'normal':
            curs.execute(db_change("select title, sub, date, code from rd where not stop = 'O' order by date desc limit 50"))
        elif tool == 'close':
            curs.execute(db_change("select title, sub, date, code from rd where stop = 'O' order by date desc limit 50"))
        else:
            curs.execute(db_change('select title, sub, date, code from rd where stop != "O" order by date asc limit 50'))

        db_data = curs.fetchall()

        last_editor = []
        for for_a in db_data:
            curs.execute(db_change("select ip from topic where code = ? order by id + 0 desc limit 1"), [for_a[3]])
            db_data_2 = curs.fetchall()
            if db_data_2:
                last_editor += [db_data_2[0][0]]
            else:
                last_editor += ['']

        last_editor_ip_pas = ip_pas(last_editor)

        count = 0
        for data in db_data:
            div += '' + \
                '<tr>' + \
                    '<td>' + \
                        '<a href="/thread/' + data[3] + '">' + html.escape(data[1]) + '</a> ' + \
                        '<a href="/topic/' + url_pas(data[0]) + '">(' + html.escape(data[0]) + ')</a>' + \
                        (' <a href="/thread/' + data[3] + '/tool">(' + load_lang('tool') + ')</a>' if admin_auth == 1 else '') + \
                    '</td>' + \
                    '<td>' + last_editor_ip_pas[last_editor[count]] + '</td>' + \
                    '<td>' + data[2] + '</td>' + \
                '</tr>' + \
            ''

            count += 1

        div += '' + \
                '</tbody>' + \
            '</table>' + \
        ''

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('recent_discussion'), wiki_set(), wiki_custom(), wiki_css([m_sub, 0])],
            data = div,
            menu = [['other', load_lang('return')]]
        ))