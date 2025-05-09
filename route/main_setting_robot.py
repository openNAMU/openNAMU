from .tool.func import *

async def main_setting_robot():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 0)

        curs.execute(db_change("select data from other where name = 'robot'"))
        db_data = curs.fetchall()
        if db_data:
            data = db_data[0][0]
        else:
            data = ''

        curs.execute(db_change("select data from other where name = 'robot_default'"))
        db_data_2 = curs.fetchall()
        if db_data_2 and db_data_2[0][0] != '':
            default_data = 'checked'
        else:
            default_data = ''
        
        if flask.request.method == 'POST':
            if db_data:
                curs.execute(db_change("update other set data = ? where name = 'robot'"), [flask.request.form.get('content', '')])
            else:
                curs.execute(db_change("insert into other (name, data, coverage) values ('robot', ?, '')"), [flask.request.form.get('content', '')])

            if db_data_2:
                curs.execute(db_change("update other set data = ? where name = 'robot_default'"), [flask.request.form.get('default', '')])
            else:
                curs.execute(db_change("insert into other (name, data, coverage) values ('robot_default', ?, '')"), [flask.request.form.get('default', '')])

            await acl_check(tool = 'owner_auth', memo = 'edit_set (robot)')

            return redirect(conn, '/setting/robot')
        else:
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = ['robots.txt', await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = '''
                    <a href="/robots.txt">(''' + get_lang(conn, 'view') + ''')</a>
                    <hr class="main_hr">
                    <form method="post">
                        <textarea class="opennamu_textarea_500" name="content">''' + html.escape(data) + '''</textarea>
                        <hr class="main_hr">
                        <label><input type="checkbox" name="default" ''' + default_data + '''> ''' + get_lang(conn, 'default') + '''</label>
                        <hr class="main_hr">
                        <button id="opennamu_save_button" type="submit">''' + get_lang(conn, 'save') + '''</button>
                    </form>
                ''',
                menu = [['setting', get_lang(conn, 'return')]]
            ))