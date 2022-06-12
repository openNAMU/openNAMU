from .tool.func import *

def topic_list(name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        div = ''
        tool = flask.request.args.get('tool', '')

        plus = ''
        menu = [['topic/' + url_pas(name), load_lang('return')]]

        if tool == 'close':
            curs.execute(db_change("select code, sub from rd where title = ? and stop = 'O' order by sub asc"), [name])

            sub = load_lang('closed_discussion')
        elif tool == 'agree':
            curs.execute(db_change("select code, sub from rd where title = ? and agree = 'O' order by sub asc"), [name])

            sub = load_lang('agreed_discussion')
        else:
            sub = load_lang('discussion_list')
            menu = [['w/' + url_pas(name), load_lang('document')]]

            plus = '''
                <a href="/topic/''' + url_pas(name) + '?tool=close">(' + load_lang('closed_discussion') + ''')</a>
                <a href="/topic/''' + url_pas(name) + '?tool=agree">(' + load_lang('agreed_discussion') + ''')</a>
                <hr class="main_hr">
                <a href="/thread/0">(''' + load_lang('make_new_topic') + ''')</a>
            '''

            curs.execute(db_change("select code, sub from rd where title = ? and stop != 'O' order by date desc"), [name])

        for data in curs.fetchall():
            curs.execute(db_change("select id from topic where code = ? order by id + 0 desc limit 1"), [data[0]])
            t_data = curs.fetchall()

            div += '''
                <h2><a href="/thread/''' + data[0] + '">' + html.escape(data[0] + '. ' + data[1]) + '''</a></h2>
                <div id="topic_pre_''' + data[0] + '''"></div>
                <div id="topic_back_pre_''' + data[0] + '''"></div>
                <script>
                    opennamu_do_thread_make(''' + data[0] + ', "list", "/normal/1", "topic_pre_' + data[0] + '''");
                    if(''' + t_data[0][0] + ''' !== 1) {
                        opennamu_do_thread_make(''' + data[0] + ', "list", "/normal/' + t_data[0][0] + '", "topic_back_pre_' + data[0] + '''");
                    }
                </script>
            '''

        if div == '':
            plus = re.sub(r'^<br>', '', plus)

        return easy_minify(flask.render_template(skin_check(),
            imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + sub + ')', 0])],
            data = div + plus,
            menu = menu
        ))
