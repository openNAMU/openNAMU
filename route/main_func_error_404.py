from .tool.func import *

def main_func_error_404(e = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if flask.request.path == '/':
            curs.execute(db_change('select data from other where name = "frontpage"'))
            db_data = curs.fetchall()
            db_data = db_data[0][0] if db_data and db_data[0][0] != '' else 'FrontPage'
            
            return redirect(conn, '/w/' + url_pas(db_data))
        else:
            curs.execute(db_change('select data from other where name = "manage_404_page"'))
            db_data = curs.fetchall()
            db_data = db_data[0][0] if db_data else ''

            if os.path.exists('404.html') and db_data == '404_file':
                return open('404.html', encoding = 'utf8').read(), 404
            else:
                curs.execute(db_change('select data from other where name = "manage_404_page_content"'))
                db_data = curs.fetchall()
                db_data = db_data[0][0] if db_data and db_data[0][0] != '' else ''

                if db_data != '':
                    return easy_minify(conn, flask.render_template(skin_check(conn),
                        imp = ['404', wiki_set(conn), wiki_custom(conn), wiki_css([0, 0])],
                        data = db_data,
                        menu = 0
                    )), 404
                else:
                    return re_error(conn, 46)