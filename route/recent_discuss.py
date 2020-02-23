from .tool.func import *

def recent_discuss_2(conn):
    curs = conn.cursor()

    div = ''

    if flask.request.args.get('what', 'normal') == 'normal':
        div += '<a href="/recent_discuss?what=close">(' + load_lang('close_discussion') + ')</a>'

        m_sub = 0
    else:
        div += '<a href="/recent_discuss">(' + load_lang('open_discussion') + ')</a>'

        m_sub = ' (' + load_lang('closed') + ')'

    div +=  '''
            <hr class=\"main_hr\">
            <table id="main_table_set">
                <tbody>
                    <tr>
                        <td id="main_table_width_half">''' + load_lang('discussion_name') + '''</td>
                        <td id="main_table_width_half">''' + load_lang('time') + '''</td>
                    </tr>
            '''

    if m_sub == 0:
        curs.execute(db_change("select title, sub, date from rd where not stop = 'O' order by date desc limit 50"))
    else:
        curs.execute(db_change("select title, sub, date from rd where stop = 'O' order by date desc limit 50"))

    for data in curs.fetchall():
        curs.execute(db_change("select code from topic where id = '1' and title = ? and sub = ?"), [data[0], data[1]])
        get_code = curs.fetchall()
        if get_code and get_code[0][0] != '':
            get_code = get_code[0][0]
        else:
            get_code = '1'

        title = html.escape(data[0])
        sub = html.escape(data[1])

        div += '<tr><td><a href="/thread/' + get_code + '">' + sub + '</a> (' + title + ')</td><td>' + data[2] + '</td></tr>'

    div += '</tbody></table>'

    return easy_minify(flask.render_template(skin_check(),
        imp = [load_lang('recent_discussion'), wiki_set(), custom(), other2([m_sub, 0])],
        data = div,
        menu = 0
    ))
