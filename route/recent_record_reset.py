from .tool.func import *

def recent_record_reset(name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check() != 1:
            return re_error('/error/3')

        if flask.request.method == 'POST':
            admin_check(None, 'record reset ' + name)

            curs.execute(db_change("delete from history where ip = ?"), [name])
            conn.commit()

            return redirect('/record/' + url_pas(name))
        else:
            return easy_minify(flask.render_template(skin_check(),
                imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('record_reset') + ')', 0])],
                data = '''
                    <form method="post">
                        <span>''' + load_lang('delete_warning') + '''</span>
                        <hr class="main_hr">
                        <button type="submit">''' + load_lang('reset') + '''</button>
                    </form>
                ''',
                menu = [['record/' + url_pas(name), load_lang('return')]]
            ))