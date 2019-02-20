from .tool.func import *

def restart_2(conn):
    curs = conn.cursor()

    if admin_check(None, 'restart') != 1:
        return re_error('/error/3')

    if flask.request.method == 'POST':
        os.execl(sys.executable, sys.executable, *sys.argv)
    else:
        print('Restart')

        return easy_minify(flask.render_template(skin_check(), 
            imp = [load_lang('wiki_restart'), wiki_set(), custom(), other2([0, 0])],
            data =  '''
                    <form method="post">
                        <button type="submit">''' + load_lang('restart') + '''</button>
                    </form>
                    ''',
            menu = [['manager', load_lang('return')]]
        ))