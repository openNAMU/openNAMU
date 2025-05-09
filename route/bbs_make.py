from .tool.func import *

async def bbs_make():   
    with get_db_connect() as conn:
        curs = conn.cursor()

        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 3)
        
        if flask.request.method == 'POST':
            curs.execute(db_change('select set_id from bbs_set where set_name = "bbs_name" order by set_id + 0 desc'))
            db_data = curs.fetchall()

            bbs_num = str(int(db_data[0][0]) + 1) if db_data else '1'
            bbs_name = flask.request.form.get('bbs_name', 'test')
            bbs_type = flask.request.form.get('bbs_type', 'comment')
            bbs_type = bbs_type if bbs_type in ['comment', 'thread'] else 'comment'

            curs.execute(db_change("insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_name', '', ?, ?)"), [bbs_num, bbs_name])
            curs.execute(db_change("insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_type', '', ?, ?)"), [bbs_num, bbs_type])

            return redirect(conn, '/bbs/main')
        else:
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'bbs_make'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        <input placeholder="''' + get_lang(conn, 'bbs_name') + '''" name="bbs_name">
                        <hr class="main_hr">
                        
                        <select name="bbs_type">
                            <option value="comment">''' + get_lang(conn, 'comment_base') + '''</option>
                            <option value="thread">''' + get_lang(conn, 'thread_base') + '''</option>
                        </select>
                        <hr class="main_hr">
                        
                        <button type="submit">''' + get_lang(conn, 'save') + '''</button>
                    </form>
                ''',
                menu = [['bbs/main', get_lang(conn, 'return')]]
            ))