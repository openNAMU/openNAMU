from .tool.func import *

def user_edit_filter(name = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check(conn, 1, None) != 1:
            return redirect(conn, '/block_log')

        if flask.request.method == 'POST':
            curs.execute(db_change('delete from user_set where name = "edit_filter" and id = ?'), [name])

            return redirect(conn, '/edit_filter/' + url_pas(name))
        else:
            curs.execute(db_change('select data from user_set where name = "edit_filter" and id = ?'), [name])
            db_data = curs.fetchall()
            p_data = db_data[0][0] if db_data else ''
            p_data = '<textarea readonly class="opennamu_textarea_500">' + html.escape(p_data) + '</textarea>'

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [name, wiki_set(conn), wiki_custom(conn), wiki_css(['(' + get_lang(conn, 'edit_filter') + ')', 0])],
                data = p_data + '''
                    <hr class="main_hr">
                    <form method="post">
                        <button type="submit">''' + get_lang(conn, 'delete') + '''</button>
                    </form>
                ''',
                menu = [['block_log', get_lang(conn, 'return')]]
            ))