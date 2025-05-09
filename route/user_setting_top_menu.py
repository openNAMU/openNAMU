from .tool.func import *

async def user_setting_top_menu():
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if (await ban_check(ip))[0] == 1:
            return await re_error(conn, 0)

        if ip_or_user(ip) == 1:
            return redirect(conn, '/login')
        
        if flask.request.method == 'POST':
            curs.execute(db_change("select data from user_set where name = 'top_menu' and id = ?"), [ip])
            if curs.fetchall():
                curs.execute(db_change("update user_set set data = ? where name = 'top_menu' and id = ?"), [flask.request.form.get('content', ''), ip])
            else:
                curs.execute(db_change("insert into user_set (name, data, id) values ('top_menu', ?, ?)"), [flask.request.form.get('content', ''), ip])

            return redirect(conn, '/change/top_menu')
        else:
            curs.execute(db_change("select data from user_set where name = 'top_menu' and id = ?"), [ip])
            db_data = curs.fetchall()
            db_data = db_data[0][0] if db_data else ''
            
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'user_added_menu'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
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