from .tool.func import *

async def main_setting_top_menu():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 0)
        
        if flask.request.method == 'POST':
            curs.execute(db_change("select name from other where name = 'top_menu'"))
            if curs.fetchall():
                curs.execute(db_change("update other set data = ? where name = 'top_menu'"), [flask.request.form.get('content', '')])
            else:
                curs.execute(db_change("insert into other (name, data, coverage) values ('top_menu', ?, '')"), [flask.request.form.get('content', '')])

            await acl_check(tool = 'owner_auth', memo = 'edit_set (top_menu)')

            return redirect(conn, '/setting/top_menu')
        else:
            curs.execute(db_change("select data from other where name = 'top_menu'"))
            db_data = curs.fetchall()
            db_data = db_data[0][0] if db_data else ''
            
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'top_menu_setting'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = '''
                    <span>
                        EX)
                        <br>
                        ONTS
                        <br>
                        https://2du.pythonanywhere.com/
                        <br>
                        FrontPage
                        <br>
                        /w/FrontPage
                    </span>
                    <hr class="main_hr">
                    ''' + get_lang(conn, 'not_support_skin_warning') + '''
                    <hr class="main_hr">
                    <form method="post">
                        <textarea class="opennamu_textarea_500" placeholder="''' + get_lang(conn, 'enter_top_menu_setting') + '''" name="content" id="content">''' + html.escape(db_data) + '''</textarea>
                        <hr class="main_hr">
                        <button id="opennamu_save_button" type="submit">''' + get_lang(conn, 'save') + '''</button>
                    </form>
                ''',
                menu = [['setting', get_lang(conn, 'return')]]
            ))