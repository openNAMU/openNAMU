from .tool.func import *

def topic_tool_acl(topic_num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check(3) != 1:
            return re_error('/error/3')

        ip = ip_check()
        time = get_time()
        topic_num = str(topic_num)

        curs.execute(db_change("select title, sub from rd where code = ?"), [topic_num])
        rd_d = curs.fetchall()
        if not rd_d:
            return redirect('/')

        if flask.request.method == 'POST':
            admin_check(3, 'topic_acl_set (code ' + topic_num + ')')

            curs.execute(db_change("select id from topic where code = ? order by id + 0 desc limit 1"), [topic_num])
            topic_check = curs.fetchall()
            if topic_check:
                acl_data = flask.request.form.get('acl', '')
                acl_data_view = flask.request.form.get('acl_view', '')

                curs.execute(db_change("update rd set acl = ? where code = ?"), [
                    acl_data, 
                    topic_num
                ])
                
                curs.execute(db_change("select set_data from topic_set where thread_code = ? and set_name = 'thread_view_acl'"), [topic_num])
                db_data = curs.fetchall()
                if db_data:
                    curs.execute(db_change("update topic_set set set_data = ? where thread_code = ?"), [
                        acl_data_view,
                        topic_num
                    ])
                else:
                    curs.execute(db_change("insert into topic_set (thread_code, set_name, set_id, set_data) values (?, 'thread_view_acl', '1', ?)"), [
                        topic_num,
                        acl_data_view
                    ])

                do_add_thread(
                    topic_num,
                    load_lang('acl_thread_change') + ' : ' + acl_data,
                    '1'
                )
                do_reload_recent_thread(
                    topic_num, 
                    time
                )

            return redirect('/thread/' + topic_num)
        else:
            acl_list = get_acl_list()
            acl_html_list = ''
            acl_html_list_view = ''

            curs.execute(db_change("select acl from rd where code = ?"), [topic_num])
            topic_acl_get = curs.fetchall()
            for data_list in acl_list:
                if topic_acl_get and topic_acl_get[0][0] == data_list:
                    check = 'selected="selected"'
                else:
                    check = ''

                acl_html_list += '<option value="' + data_list + '" ' + check + '>' + (data_list if data_list != '' else 'normal') + '</option>'

            curs.execute(db_change("select set_data from topic_set where thread_code = ? and set_name = 'thread_view_acl'"), [topic_num])
            db_data = curs.fetchall()
            for data_list in acl_list:
                if db_data and db_data[0][0] == data_list:
                    check = 'selected="selected"'
                else:
                    check = ''

                acl_html_list_view += '<option value="' + data_list + '" ' + check + '>' + (data_list if data_list != '' else 'normal') + '</option>'

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('topic_acl_setting'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        <a href="/acl/TEST#exp">(''' + load_lang('reference') + ''')</a>
                        <h2>''' + load_lang('thread_acl') + '''</h2>
                        <select name="acl">
                            ''' + acl_html_list + '''
                        </select>
                        <h2>''' + load_lang('view_acl') + ''' (''' + load_lang('beta') + ''')</h2>
                        <select name="acl_view">
                            ''' + acl_html_list_view + '''
                        </select>
                        <hr class="main_hr">
                        <button type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['thread/' + topic_num + '/tool', load_lang('return')]]
            ))