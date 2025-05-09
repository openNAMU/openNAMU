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
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'topic_delete'), await wiki_set(), await wiki_custom(conn), wiki_css(['(#' + num + ')', 0])],
                data = '''
                    <hr class="main_hr">
                    <form method="post">
                        <button type="submit">''' + get_lang(conn, 'start') + '''</button>
                    </form>
                ''',
                menu = [['thread/' + topic_num + '/comment/' + num + '/tool', get_lang(conn, 'return')]]
            ))