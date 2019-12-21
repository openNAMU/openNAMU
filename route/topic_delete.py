from .tool.func import *

def topic_delete_2(conn, topic_num):
    curs = conn.cursor()

    if admin_check(None) != 1:
        return re_error('/error/3')

    topic_change_data = topic_change(topic_num)
    name = topic_change_data[0]
    sub = topic_change_data[1]

    if flask.request.method == 'POST':
        curs.execute(db_change("delete from topic where title = ? and sub = ?"), [name, sub])
        curs.execute(db_change("delete from rd where title = ? and sub = ?"), [name, sub])
        conn.commit()

        return redirect('/topic/' + url_pas(name))
    else:
        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('topic_delete'), wiki_set(), custom(), other2([0, 0])],
            data = '''
                <hr class=\"main_hr\">
                <form method="post">
                    <button type="submit">''' + load_lang('start') + '''</button>
                </form>
            ''',
            menu = [['thread/' + str(topic_num) + '/tool', load_lang('return')]]
        ))
