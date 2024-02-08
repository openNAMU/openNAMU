from .tool.func import *

def user_edit_filter(name = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check(1, None) != 1:
            return redirect('/block_log')

        if flask.request.method == 'POST':
            curs.execute(db_change('delete from user_set where name = "edit_filter" and id = ?'), [name])

            return redirect('/edit_filter/' + url_pas(name))
        else:
            curs.execute(db_change('select data from user_set where name = "edit_filter" and id = ?'), [name])
            db_data = curs.fetchall()
            p_data = db_data[0][0] if db_data else ''
            p_data = '<textarea readonly class="opennamu_textarea_500">' + html.escape(p_data) + '</textarea>'

            return easy_minify(flask.render_template(skin_check(),
                imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('edit_filter') + ')', 0])],
                data = p_data + '''
                    <hr class="main_hr">
                    <form method="post">
                        <button type="submit">''' + load_lang('delete') + '''</button>
                    </form>
                ''',
                menu = [['block_log', load_lang('return')]]
            ))