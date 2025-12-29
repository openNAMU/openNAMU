from .tool.func import *

async def topic_comment_delete(topic_num = 1, num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if await acl_check(tool = 'owner_auth') == 1:
            return await re_error(conn, 3)

        topic_num = str(topic_num)
        num = str(num)

        if flask.request.method == 'POST':
            curs.execute(db_change("delete from topic where code = ? and id = ?"), [topic_num, num])

            return redirect(conn, '/thread/' + topic_num)
        else:
            return await render_template(
                await get_lang('topic_delete'),
                '''
                    <hr class="main_hr">
                    <form method="post">
                        <button class="__ON_BUTTON__" type="submit">''' + await get_lang('start') + '''</button>
                    </form>
                ''',
                '(#' + num + ')',
                [['thread/' + topic_num + '/comment/' + num + '/tool', await get_lang('return')]]
            )
