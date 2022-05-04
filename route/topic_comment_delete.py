from .tool.func import *

def topic_comment_delete(topic_num = 1, num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check(None) != 1:
            return re_error('/error/3')

        topic_num = str(topic_num)
        num = str(num)

        if flask.request.method == 'POST':
            curs.execute(db_change("delete from topic where code = ? and id = ?"), [topic_num, num])
            conn.commit()

            return redirect('/thread/' + topic_num)
        else:
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('topic_delete'), wiki_set(), wiki_custom(), wiki_css(['(#' + num + ')', 0])],
                data = '''
                    <hr class="main_hr">
                    <form method="post">
                        <button type="submit">''' + load_lang('start') + '''</button>
                    </form>
                ''',
                menu = [['thread/' + topic_num + '/comment/' + num + '/tool', load_lang('return')]]
            ))