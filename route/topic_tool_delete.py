from .tool.func import *

async def topic_tool_delete(topic_num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if await acl_check(tool = 'owner_auth') == 1:
            return await re_error(conn, 3)

        topic_num = str(topic_num)

        if flask.request.method == 'POST':
            await acl_check(tool = 'owner_auth', memo = 'delete topic (' + topic_num + ')')

            curs.execute(db_change("delete from topic where code = ?"), [topic_num])
            curs.execute(db_change("delete from rd where code = ?"), [topic_num])

            return redirect(conn, '/')
        else:
            return easy_minify(flask.render_template(await skin_check(),
                imp = [await get_lang('topic_delete'), await wiki_set(), await wiki_custom(), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        <span>''' + await get_lang('delete_warning') + '''</span>
                        <hr class="main_hr">
                        <button type="submit">''' + await get_lang('delete') + '''</button>
                    </form>
                ''',
                menu = [['thread/' + topic_num + '/tool', await get_lang('return')]]
            ))