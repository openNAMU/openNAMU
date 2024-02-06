from .tool.func import *

def main_func_error_404(e = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()
        
        if os.path.exists('404.html') and flask.request.path != '/':
            return open('404.html', encoding = 'utf8').read(), 404
        else:
            curs.execute(db_change('select data from other where name = "frontpage"'))
            db_data = curs.fetchall()
            db_data = db_data[0][0] if db_data and db_data[0][0] != '' else 'FrontPage'
            
            return redirect('/w/' + url_pas(db_data))