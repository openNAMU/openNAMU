from .tool.func import *

def delete_admin_group_2(conn, name):
    curs = conn.cursor()

    if admin_check() != 1:
        return re_error('/error/3')
    
    if flask.request.method == 'POST':
        admin_check(None, 'alist del ' + name)
        curs.execute(db_change("delete from alist where name = ?"), [name])

        return redirect('/give_log')
    
    return easy_minify(flask.render_template(skin_check(),
        imp = [load_lang("delete_admin_group"), wiki_set(), custom(), other2(['(' + name + ')', 0])],
        data = '''
            <form method=post>
                <button type=submit>''' + load_lang('start') + '''</button>
            </form>
        ''',
        menu = [['give_log', load_lang('return')]]
    ))  
