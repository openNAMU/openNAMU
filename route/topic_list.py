from .tool.func import *

async def topic_list(name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        div = ''
        tool = flask.request.args.get('tool', '')

        plus = ''
        menu = [['topic/' + url_pas(name), await get_lang('return')]]

        if tool == 'close':
            curs.execute(db_change("select code, sub from rd where title = ? and stop = 'O' order by sub asc"), [name])

            sub = await get_lang('closed_discussion')
        elif tool == 'agree':
            curs.execute(db_change("select code, sub from rd where title = ? and agree = 'O' order by sub asc"), [name])

            sub = await get_lang('agreed_discussion')
        else:
            sub = await get_lang('discussion_list')
            menu = [['w/' + url_pas(name), await get_lang('document')]]

            plus = '''
                <a href="/topic/''' + url_pas(name) + '?tool=close">(' + await get_lang('closed_discussion') + ''')</a>
                <a href="/topic/''' + url_pas(name) + '?tool=agree">(' + await get_lang('agreed_discussion') + ''')</a>
                <hr class="main_hr">
                <a href="/thread/0/''' + url_pas(name) + '''">(''' + await get_lang('make_new_topic') + ''')</a>
            '''

            curs.execute(db_change("select code, sub from rd where title = ? and stop != 'O' order by date desc"), [name])

        for data in curs.fetchall():
            div += '<h2><a href="/thread/' + data[0] + '">' + data[0] + '. ' + html.escape(data[1]) + '</a></h2>'

        if div == '':
            plus = re.sub(r'^<br>', '', plus)

        return easy_minify(flask.render_template(await skin_check(),
            imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css(['(' + sub + ')', 0])],
            data = div + plus,
            menu = menu
        ))
