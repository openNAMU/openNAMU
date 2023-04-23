from .tool.func import *

def main_setting_top_menu():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check() != 1:
            return re_error('/ban')
        
        if flask.request.method == 'POST':
            curs.execute(db_change("select name from other where name = 'top_menu'"))
            if curs.fetchall():
                curs.execute(db_change("update other set data = ? where name = 'top_menu'"), [flask.request.form.get('content', '')])
            else:
                curs.execute(db_change("insert into other (name, data, coverage) values ('top_menu', ?, '')"), [flask.request.form.get('content', '')])

            conn.commit()

            admin_check(None, 'edit_set (top_menu)')

            return redirect('/setting/top_menu')
        else:
            curs.execute(db_change("select data from other where name = 'top_menu'"))
            db_data = curs.fetchall()
            db_data = db_data[0][0] if db_data else ''
            
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('top_menu_setting'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
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
                    ''' + load_lang('not_support_skin_warning') + '''
                    <hr class="main_hr">
                    <form method="post">
                        <textarea class="opennamu_textarea_500" placeholder="''' + load_lang('enter_top_menu_setting') + '''" name="content" id="content">''' + html.escape(db_data) + '''</textarea>
                        <hr class="main_hr">
                        <button id="opennamu_save_button" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['setting', load_lang('return')]]
            ))