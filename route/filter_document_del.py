from .tool.func import *

def filter_document_del(name = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check() != 1:
            return re_error('/error/3')

        curs.execute(db_change("delete from html_filter where html = ? and kind = 'document'"), [name])
        conn.commit()

        return redirect('/filter/document/list')