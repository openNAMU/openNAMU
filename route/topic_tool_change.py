from .tool.func import *

async def topic_tool_change(topic_num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if await acl_check(tool = 'owner_auth') == 1:
            return await re_error(conn, 3)

        time = get_time()
        topic_num = str(topic_num)

        curs.execute(db_change("select title, sub from rd where code = ?"), [topic_num])
        rd_d = curs.fetchall()
        if not rd_d:
            return redirect(conn, '/')

        if flask.request.method == 'POST':
            await acl_check(tool = 'owner_auth', memo = 'move_topic (code ' + topic_num + ')')

            title_d = flask.request.form.get('title', 'test')
            sub_d = flask.request.form.get('sub', 'test')

            curs.execute(db_change("update rd set title = ?, sub = ? where code = ?"), [title_d, sub_d, topic_num])

            do_add_thread(conn, topic_num, await get_lang('topic_name_change') + ' : ' + rd_d[0][1] + ' (' + rd_d[0][0] + ') → ' + sub_d + ' (' + title_d + ')', '1')
            do_reload_recent_thread(conn, topic_num, time)

            return redirect(conn, '/thread/' + topic_num)
        else:
            return await render_template(
                await get_lang('topic_name_change'),
                '''
                    <form method="post">
                        ''' + await get_lang('document_name') + '''
                        <hr class="main_hr">
                        <input class="__ON_INPUT__" value="''' + html.escape(rd_d[0][0]) + '''" name="title" type="text">
                        <hr class="main_hr">
                        ''' + await get_lang('discussion_name') + '''
                        <hr class="main_hr">
                        <input class="__ON_INPUT__" value="''' + html.escape(rd_d[0][1]) + '''" name="sub" type="text">
                        <hr class="main_hr">
                        <button class="__ON_BUTTON__" type="submit">''' + await get_lang('save') + '''</button>
                    </form>
                ''',
                0,
                [['thread/' + topic_num + '/tool', await get_lang('return')]]
            )
