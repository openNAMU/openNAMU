from .tool.func import *

async def topic_list(name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        div = ''
        tool = flask.request.args.get('tool', '')

        plus = ''
        menu = [['topic/' + url_pas(name), get_lang(conn, 'return')]]

        if tool == 'close':
            curs.execute(db_change("select code, sub from rd where title = ? and stop = 'O' order by sub asc"), [name])

            sub = get_lang(conn, 'closed_discussion')
        elif tool == 'agree':
            curs.execute(db_change("select code, sub from rd where title = ? and agree = 'O' order by sub asc"), [name])

            sub = get_lang(conn, 'agreed_discussion')
        else:
            sub = get_lang(conn, 'discussion_list')
            menu = [['w/' + url_pas(name), get_lang(conn, 'document')]]

            plus = '''
                <a href="/topic/''' + url_pas(name) + '?tool=close">(' + get_lang(conn, 'closed_discussion') + ''')</a>
                <a href="/topic/''' + url_pas(name) + '?tool=agree">(' + get_lang(conn, 'agreed_discussion') + ''')</a>
                <hr class="main_hr">
                <a href="/thread/0/''' + url_pas(name) + '''">(''' + get_lang(conn, 'make_new_topic') + ''')</a>
            '''

            curs.execute(db_change("select code, sub from rd where title = ? and stop != 'O' order by date desc"), [name])

        for data in curs.fetchall():
            div += '<h2><a href="/thread/' + data[0] + '">' + data[0] + '. ' + html.escape(data[1]) + '</a></h2>'

        if div == '':
            plus = re.sub(r'^<br>', '', plus)

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css(['(' + sub + ')', 0])],
            data = div + plus,
            menu = menu
        ))
