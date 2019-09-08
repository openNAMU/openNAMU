from .tool.func import *

def recent_history_delete_2(conn, name):
    curs = conn.cursor()

    num = str(int(number_check(flask.request.args.get('num', '1'))))

    if admin_check() != 1:
        return re_error('/error/3')

    if flask.request.method == 'POST':
        admin_check(None, 'history delete r' + num)

        curs.execute("delete from history where id = ? and title = ?", [num, name])
        conn.commit()

        return redirect('/history/' + url_pas(name))
    else:
        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2(['(r' + num + ')', 0])],
            data = '''
                <form method="post">
                    <button type="submit">''' + load_lang('history_delete') + '''</button>
                </form>
            ''',
            menu = [['history/' + url_pas(name), load_lang('return')]]
        ))