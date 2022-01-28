from .tool.func import *

def main_sys_restart():
    if admin_check() != 1:
        return re_error('/error/3')

    if flask.request.method == 'POST':
        admin_check(None, 'restart')

        print('----')
        print('Restart')

        try:
            os.execl(sys.executable, sys.executable, *sys.argv)
        except:
            pass
        
        try:
            os.execl(sys.executable, '"' + sys.executable + '"', *sys.argv)
        except:
            return re_error('/error/33')
    else:
        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('wiki_restart'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = '''
                <form method="post">
                    <button type="submit">''' + load_lang('restart') + '''</button>
                </form>
            ''',
            menu = [['manager', load_lang('return')]]
        ))