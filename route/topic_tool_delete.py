from .tool.func import *

def topic_tool_delete(topic_num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check(None) != 1:
            return re_error('/error/3')

        topic_num = str(topic_num)

        if flask.request.method == 'POST':
            curs.execute(db_change("delete from topic where code = ?"), [topic_num])
            curs.execute(db_change("delete from rd where code = ?"), [topic_num])
            conn.commit()

            return redirect('/')
        else:
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('topic_delete'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        <span>''' + load_lang('delete_warning') + '''</span>
                        <hr class="main_hr">
                        <button type="submit">''' + load_lang('delete') + '''</button>
                    </form>
                ''',
                menu = [['thread/' + topic_num + '/tool', load_lang('return')]]
            ))