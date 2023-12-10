from .tool.func import *

def recent_history_delete(name = 'Test', rev = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        num = str(rev)

        if admin_check() != 1:
            return re_error('/error/3')

        if flask.request.method == 'POST':
            admin_check(None, 'history delete ' + name + ' r' + num)

            curs.execute(db_change("delete from history where id = ? and title = ?"), [num, name])
            conn.commit()

            return redirect('/history/' + url_pas(name))
        else:
            return easy_minify(flask.render_template(skin_check(),
                imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('history_delete') + ') (r' + num + ')', 0])],
                data = '''
                    <form method="post">
                        <span>''' + load_lang('delete_warning') + '''</span>
                        <hr class="main_hr">
                        <button type="submit">''' + load_lang('delete') + '''</button>
                    </form>
                ''',
                menu = [['history/' + url_pas(name), load_lang('return')]]
            ))