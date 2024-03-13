from .tool.func import *

def topic_tool_change(topic_num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check(conn, None) != 1:
            return re_error(conn, '/error/3')

        ip = ip_check()
        time = get_time()
        topic_num = str(topic_num)

        curs.execute(db_change("select title, sub from rd where code = ?"), [topic_num])
        rd_d = curs.fetchall()
        if not rd_d:
            return redirect(conn, '/')

        if flask.request.method == 'POST':
            admin_check(conn, None, 'move_topic (code ' + topic_num + ')')

            curs.execute(db_change("select id from topic where code = ? order by id + 0 desc limit 1"), [topic_num])
            topic_check = curs.fetchall()
            if topic_check:
                title_d = flask.request.form.get('title', 'test')
                sub_d = flask.request.form.get('sub', 'test')

                curs.execute(db_change("update rd set title = ?, sub = ? where code = ?"), [
                    title_d,
                    sub_d,
                    topic_num
                ])

                do_add_thread(conn, 
                    topic_num,
                    get_lang(conn, 'topic_name_change') + ' : ' + sub_d + ' (' + title_d + ')',
                    '1'
                )
                do_reload_recent_thread(conn, 
                    topic_num, 
                    time
                )

            return redirect(conn, '/thread/' + topic_num)
        else:
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'topic_name_change'), wiki_set(conn), wiki_custom(conn), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        ''' + get_lang(conn, 'document_name') + '''
                        <hr class="main_hr">
                        <input value="''' + rd_d[0][0] + '''" name="title" type="text">
                        <hr class="main_hr">
                        ''' + get_lang(conn, 'discussion_name') + '''
                        <hr class="main_hr">
                        <input value="''' + rd_d[0][1] + '''" name="sub" type="text">
                        <hr class="main_hr">
                        <button type="submit">''' + get_lang(conn, 'save') + '''</button>
                    </form>
                ''',
                menu = [['thread/' + topic_num + '/tool', get_lang(conn, 'return')]]
            ))