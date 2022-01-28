from .tool.func import *

def main_sys_shutdown():
    if admin_check() != 1:
        return re_error('/error/3')

    if flask.request.method == 'POST':
        admin_check(None, 'shutdown')

        print('----')
        print('shutdown')

        os._exit(os.EX_OK)
    else:
        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('wiki_shutdown'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = '''
                <form method="post">
                    <button type="submit">''' + load_lang('shutdown') + '''</button>
                </form>
            ''',
            menu = [['manager', load_lang('return')]]
        ))