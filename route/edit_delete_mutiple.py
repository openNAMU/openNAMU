from .tool.func import *
from . import edit_delete

def edit_delete_mutiple_2(conn):
    curs = conn.cursor()

    ip = ip_check()
    if admin_check() != 1:
        return re_error('/ban')

    if flask.request.method == 'POST':
        all_title = re.findall(r'([^\n]+)\n', flask.request.form.get('content', '').replace('\r\n', '\n') + '\n')
        for name in all_title:
            edit_delete.edit_delete_2(conn, name)

        return redirect('/recent_changes')
    else:
        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('many_delete'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = '''
                <form method="post">
                    <textarea rows="25" placeholder="''' + load_lang('many_delete_help') + '''" name="content"></textarea>
                    <hr class=\"main_hr\">
                    <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                    <hr class=\"main_hr\">
                    <button type="submit">''' + load_lang('delete') + '''</button>
                </form>
            ''',
            menu = [['manager/1', load_lang('return')]]
        ))